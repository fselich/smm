package gcp

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
	"hash/crc32"
	"path/filepath"
	"strconv"
	"strings"
)

type Gcp struct {
	projectID string
	ctx       context.Context
	client    *secretmanager.Client
	secrets   []string
}

type version struct {
	Name     string
	State    string
	Version  int
	FullPath string
}

func NewGcp(projectID string) (*Gcp, error) {
	ctx := context.Background()
	log.Info().Msg("Connecting to GCP")
	log.Info().Msgf("Project ID: %s", projectID)
	gcp := &Gcp{projectID: projectID, ctx: ctx}
	err := gcp.gcpConnect()
	gcp.secrets, err = gcp.fetchSecrets()
	return gcp, err
}

func (g *Gcp) gcpConnect() error {

	var err error
	g.client, err = secretmanager.NewClient(g.ctx)

	if err != nil {
		return err
	}
	return nil
}

func (g *Gcp) Secrets() []string {
	return g.secrets
}

func (g *Gcp) fetchSecrets() ([]string, error) {
	listSecretsReq := &secretmanagerpb.ListSecretsRequest{
		Parent: fmt.Sprintf("projects/%s", g.projectID),
	}

	listSecrets := g.client.ListSecrets(g.ctx, listSecretsReq)

	var secrets []string
	for {
		secret, err := listSecrets.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}

		secrets = append(secrets, secret.Name)
	}

	return secrets, nil
}

func (g *Gcp) GetSecretVersions(secretName string) []version {
	req := &secretmanagerpb.ListSecretVersionsRequest{
		Parent: fmt.Sprintf("%s", secretName),
	}

	var versions []version
	it := g.client.ListSecretVersions(g.ctx, req)
	for {
		resp, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}

		if err != nil {
			log.Error().Msgf("failed to list secret versions: %v", err)
		}

		versionParts := strings.Split(resp.Name, "/")
		versionNumber, _ := strconv.Atoi(versionParts[len(versionParts)-1])

		version := version{
			Name:     filepath.Base(secretName),
			FullPath: secretName,
			State:    resp.State.String(),
			Version:  versionNumber,
		}
		versions = append(versions, version)
	}

	return versions
}

func (g *Gcp) GetSecret(secretName string) []byte {
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("%s/versions/latest", secretName),
	}

	result, err := g.client.AccessSecretVersion(g.ctx, accessRequest)
	if err != nil {
		log.Fatal().Msgf("failed to access secret version: %v", err)
	}

	return result.Payload.Data
}

func (g *Gcp) GetSecretVersion(secretName, version string) []byte {
	name := fmt.Sprintf("%s/versions/%s", secretName, version)
	log.Info().Msgf("Fetching secret version: %s", name)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := g.client.AccessSecretVersion(g.ctx, req)
	if err != nil {
		log.Fatal().Msgf("failed to access secret version: %v", err)
	}

	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		log.Fatal().Msgf("data corruption detected: expected crc32c of %d but got %d", checksum, *result.Payload.DataCrc32C)
	}

	return result.Payload.Data
}

func (g *Gcp) AddSecretVersion(secretName string, payload []byte) error {
	log.Info().Msgf("[FAKE] Adding secret version: %s", secretName)
	return nil
	//parent := fmt.Sprintf("projects/%s/secrets/%s", g.projectID, secretName)
	//
	//crc32c := crc32.MakeTable(crc32.Castagnoli)
	//checksum := int64(crc32.Checksum(payload, crc32c))
	//
	//req := &secretmanagerpb.AddSecretVersionRequest{
	//	Parent: parent,
	//	Payload: &secretmanagerpb.SecretPayload{
	//		Data:       payload,
	//		DataCrc32C: &checksum,
	//	},
	//}
	//
	//result, err := g.client.AddSecretVersion(g.ctx, req)
	//if err != nil {
	//	return fmt.Errorf("failed to add secret version: %w", err)
	//}
	//log.Info().Msgf("Added secret version: %s\n", result.Name)
	//return nil
}
