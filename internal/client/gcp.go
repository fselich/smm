package client

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
)

type Gcp struct {
	projectID   string
	ctx         context.Context
	client      *secretmanager.Client
	secretInfos []SecretInfo
	cancel      context.CancelFunc
}

func NewGcp(projectID string) (*Gcp, error) {
	ctx, cancel := context.WithCancel(context.Background())
	gcp := &Gcp{projectID: projectID, ctx: ctx, cancel: cancel}
	err := gcp.gcpConnect()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to GCP Secret Manager: %w", err)
	}

	gcp.secretInfos, err = gcp.fetchSecretInfos()
	return gcp, nil
}

func (g *Gcp) gcpConnect() error {

	var err error
	g.client, err = secretmanager.NewClient(g.ctx)

	if err != nil {
		return err
	}
	return nil
}

func (g *Gcp) Secrets() ([]SecretInfo, error) {
	if g.secretInfos == nil {
		secretInfos, err := g.fetchSecretInfos()
		if err != nil {
			return nil, err
		}
		g.secretInfos = secretInfos
	}
	return g.secretInfos, nil
}

func (g *Gcp) fetchSecretInfos() ([]SecretInfo, error) {
	listSecretsReq := &secretmanagerpb.ListSecretsRequest{
		Parent: fmt.Sprintf("projects/%s", g.projectID),
	}

	listSecrets := g.client.ListSecrets(g.ctx, listSecretsReq)

	var secretInfos []SecretInfo

	for {
		secretData, err := listSecrets.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		name := filepath.Base(secretData.Name)

		secretInfo := SecretInfo{
			Name:        name,
			FullPath:    secretData.Name,
			CreateTime:  secretData.CreateTime.AsTime(),
			Labels:      secretData.Labels,
			Annotations: secretData.Annotations,
		}

		secretInfos = append(secretInfos, secretInfo)
	}

	return secretInfos, nil
}

func (g *Gcp) GetSecretVersions(secretName string) ([]Version, error) {
	req := &secretmanagerpb.ListSecretVersionsRequest{
		Parent: fmt.Sprintf("%s", secretName),
	}

	var versions []Version
	it := g.client.ListSecretVersions(g.ctx, req)
	for {
		resp, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to list secret versions: %w", err)
		}

		versionParts := strings.Split(resp.Name, "/")
		versionNumber, _ := strconv.Atoi(versionParts[len(versionParts)-1])

		version := Version{
			Name:      filepath.Base(secretName),
			FullPath:  secretName,
			State:     resp.State.String(),
			Version:   versionNumber,
			CreatedAt: resp.CreateTime.AsTime(),
		}
		versions = append(versions, version)
	}

	return versions, nil
}

func (g *Gcp) GetSecret(secretName string) ([]byte, error) {
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("%s/versions/latest", secretName),
	}

	result, err := g.client.AccessSecretVersion(g.ctx, accessRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version %q: %w", secretName, err)
	}

	return result.Payload.Data, nil
}

func (g *Gcp) GetSecretVersion(secretName, version string) ([]byte, error) {
	name := fmt.Sprintf("%s/versions/%s", secretName, version)
	log.Info().Msgf("Fetching secret version: %s", name)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := g.client.AccessSecretVersion(g.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %w", err)
	}

	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		return nil, fmt.Errorf("data corruption detected: expected crc32c of %d but got %d", checksum, *result.Payload.DataCrc32C)
	}

	return result.Payload.Data, nil
}

func (g *Gcp) AddSecretVersion(secretName string, payload []byte) error {
	parent := fmt.Sprintf("projects/%s/secrets/%s", g.projectID, secretName)

	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(payload, crc32c))

	req := &secretmanagerpb.AddSecretVersionRequest{
		Parent: parent,
		Payload: &secretmanagerpb.SecretPayload{
			Data:       payload,
			DataCrc32C: &checksum,
		},
	}

	result, err := g.client.AddSecretVersion(g.ctx, req)
	if err != nil {
		return fmt.Errorf("failed to add secret version: %w", err)
	}
	log.Info().Msgf("Added secret version: %s\n", result.Name)
	return nil
}

func (g *Gcp) SearchInSecrets(query string) ([]SecretInfo, error) {
	secretInfos, err := g.fetchSecretInfos()
	if err != nil {
		return nil, err
	}

	var foundSecrets []SecretInfo
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, secretInfo := range secretInfos {
		wg.Add(1)
		go func(secretInfo SecretInfo) {
			defer wg.Done()
			secretData, err := g.GetSecret(secretInfo.FullPath)
			if err != nil {
				log.Error().Err(err).Str("secret", secretInfo.FullPath).Msg("failed to get secret during search")
				return
			}

			if strings.Contains(string(secretData), query) {
				mu.Lock()
				foundSecrets = append(foundSecrets, secretInfo)
				mu.Unlock()
			}
		}(secretInfo)
	}

	wg.Wait()
	return foundSecrets, nil
}

func (g *Gcp) GetSecretInfo(fullPath string) (SecretInfo, error) {
	for _, secretInfo := range g.secretInfos {
		if secretInfo.FullPath == fullPath {
			return secretInfo, nil
		}
	}

	req := &secretmanagerpb.GetSecretRequest{
		Name: fullPath,
	}

	secret, err := g.client.GetSecret(g.ctx, req)
	if err != nil {
		return SecretInfo{}, fmt.Errorf("failed to get secret info: %w", err)
	}

	name := filepath.Base(secret.Name)

	secretInfo := SecretInfo{
		Name:        name,
		FullPath:    secret.Name,
		CreateTime:  secret.CreateTime.AsTime(),
		Labels:      secret.Labels,
		Annotations: secret.Annotations,
	}

	return secretInfo, nil
}
