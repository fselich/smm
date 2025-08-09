package client

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/jaswdr/faker/v2"
)

type FakeClient struct {
}

func NewFakeClient(projectId string) (FakeClient, error) {
	return FakeClient{}, nil
}

func seedFromSecretName(secretName string) int64 {
	hasher := sha256.New()
	hasher.Write([]byte(secretName))
	hash := hasher.Sum(nil)

	seed := int64(0)
	for i := 0; i < 8 && i < len(hash); i++ {
		seed = seed*256 + int64(hash[i])
	}
	return seed
}

func isSecretEnvType(secretName string) bool {
	seed := seedFromSecretName(secretName)
	source := rand.NewPCG(uint64(seed), uint64(seed>>32))
	rng := rand.New(source)
	return rng.IntN(2) == 0
}

func (f FakeClient) GetSecretVersions(secretName string) ([]Version, error) {
	seed := seedFromSecretName(secretName + "_versions")
	source := rand.NewPCG(uint64(seed), uint64(seed>>32))
	rng := rand.New(source)

	numVersions := 1 + rng.IntN(5)
	versions := make([]Version, numVersions)

	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < numVersions; i++ {
		versionNum := i + 1
		timeOffset := time.Duration(rng.Int64N(int64(time.Hour * 24 * 30)))

		versions[i] = Version{
			Name:      fmt.Sprintf("%s-version-%d", secretName, versionNum),
			State:     "enabled",
			Version:   versionNum,
			FullPath:  fmt.Sprintf("projects/test-project/secrets/%s/versions/%d", secretName, versionNum),
			CreatedAt: baseTime.Add(timeOffset),
		}
	}

	return versions, nil
}

func (f FakeClient) createEnvSecret(secretName string) []byte {
	seed := seedFromSecretName(secretName)
	source := rand.NewPCG(uint64(seed), uint64(seed>>32))
	rng := rand.New(source)
	fk := faker.NewWithSeed(source)

	var envLines []string

	numVars := 3 + rng.IntN(10)
	for i := 0; i < numVars; i++ {
		key := strings.ToUpper(fk.Lorem().Word()) + fmt.Sprintf("_%d", rng.IntN(1000))
		value := "\"" + fk.Lorem().Word() + "\""
		envLines = append(envLines, fmt.Sprintf("%s=%s", key, value))
	}

	return []byte(strings.Join(envLines, "\n"))
}

func (f FakeClient) createJsonSecret(secretName string) []byte {
	seed := seedFromSecretName(secretName)
	source := rand.NewPCG(uint64(seed), uint64(seed>>32))
	fk := faker.NewWithSeed(source)
	js := fk.Json()

	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, []byte(js.String()), "", "  ")
	return prettyJSON.Bytes()
}

func (f FakeClient) GetSecret(secretName string) ([]byte, error) {
	if isSecretEnvType(secretName) {
		return f.createEnvSecret(secretName), nil
	}
	return f.createJsonSecret(secretName), nil
}

func (f FakeClient) createEnvSecretVersion(secretName, version string) []byte {
	versionSeed := seedFromSecretName(secretName + "_v" + version)
	versionSource := rand.NewPCG(uint64(versionSeed), uint64(versionSeed>>32))
	versionRng := rand.New(versionSource)
	fk := faker.NewWithSeed(versionSource)

	var envLines []string
	numVars := 3 + versionRng.IntN(10)
	for i := 0; i < numVars; i++ {
		key := strings.ToUpper(fk.Lorem().Word()) + fmt.Sprintf("_V%s_%d", version, versionRng.IntN(1000))
		value := "\"" + fk.Lorem().Word() + "\""
		envLines = append(envLines, fmt.Sprintf("%s=%s", key, value))
	}
	return []byte(strings.Join(envLines, "\n"))
}

func (f FakeClient) createJsonSecretVersion(secretName, version string) []byte {
	versionSeed := seedFromSecretName(secretName + "_v" + version)
	versionSource := rand.NewPCG(uint64(versionSeed), uint64(versionSeed>>32))
	fk := faker.NewWithSeed(versionSource)

	js := fk.Json()
	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, []byte(js.String()), "", "  ")
	return prettyJSON.Bytes()
}

func (f FakeClient) GetSecretVersion(secretName, version string) ([]byte, error) {
	if isSecretEnvType(secretName) {
		return f.createEnvSecretVersion(secretName, version), nil
	}
	return f.createJsonSecretVersion(secretName, version), nil
}

func (f FakeClient) AddSecretVersion(secretName string, payload []byte) error {
	return nil
}

func (f FakeClient) SearchInSecrets(query string) ([]SecretInfo, error) {
	seed := seedFromSecretName("search_" + query)
	source := rand.NewPCG(uint64(seed), uint64(seed>>32))
	rng := rand.New(source)
	fk := faker.NewWithSeed(source)

	numResults := 2 + rng.IntN(5)
	results := make([]SecretInfo, numResults)

	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < numResults; i++ {
		secretName := fmt.Sprintf("%s-%s-secret", query, fk.Lorem().Word())
		timeOffset := time.Duration(rng.Int64N(int64(time.Hour * 24 * 365)))
		
		results[i] = SecretInfo{
			Name:        secretName,
			FullPath:    fmt.Sprintf("projects/test-project/secrets/%s", secretName),
			CreateTime:  baseTime.Add(timeOffset),
			Labels:      map[string]string{"environment": "test", "type": "search-result"},
			Annotations: map[string]string{"description": fmt.Sprintf("Search result for: %s", query)},
		}
	}

	return results, nil
}

func (f FakeClient) Secrets() ([]SecretInfo, error) {
	seed := int64(12345)
	source := rand.NewPCG(uint64(seed), uint64(seed>>32))
	rng := rand.New(source)
	fk := faker.NewWithSeed(source)

	secrets := make([]SecretInfo, 30)
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i <= 29; i++ {
		secretName := fmt.Sprintf("%s-secret", fk.Lorem().Word())
		timeOffset := time.Duration(rng.Int64N(int64(time.Hour * 24 * 365)))
		
		secrets[i] = SecretInfo{
			Name:        secretName,
			FullPath:    fmt.Sprintf("projects/test-project/secrets/%s", secretName),
			CreateTime:  baseTime.Add(timeOffset),
			Labels:      map[string]string{"environment": "test", "team": fk.Company().Name()},
			Annotations: map[string]string{"description": fk.Lorem().Sentence(5)},
		}
	}
	return secrets, nil
}

func (f FakeClient) GetSecretInfo(fullPath string) (SecretInfo, error) {
	// Extract secret name from full path
	parts := strings.Split(fullPath, "/")
	if len(parts) < 4 {
		return SecretInfo{}, fmt.Errorf("invalid secret path: %s", fullPath)
	}
	secretName := parts[len(parts)-1]
	
	seed := seedFromSecretName(fullPath)
	source := rand.NewPCG(uint64(seed), uint64(seed>>32))
	rng := rand.New(source)
	fk := faker.NewWithSeed(source)
	
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	timeOffset := time.Duration(rng.Int64N(int64(time.Hour * 24 * 365)))
	
	labels := map[string]string{
		"environment": "test",
		"team":        fk.Company().Name(),
	}
	
	// Add some random labels
	if rng.IntN(2) == 0 {
		labels["type"] = "api-key"
	}
	if rng.IntN(2) == 0 {
		labels["region"] = "us-central1"
	}
	
	annotations := map[string]string{
		"description": fk.Lorem().Sentence(5),
	}
	
	// Add some random annotations
	if rng.IntN(2) == 0 {
		annotations["owner"] = fk.Person().Name()
	}
	if rng.IntN(2) == 0 {
		annotations["last-rotated"] = time.Now().AddDate(0, -rng.IntN(12), -rng.IntN(30)).Format("2006-01-02")
	}
	
	return SecretInfo{
		Name:        secretName,
		FullPath:    fullPath,
		CreateTime:  baseTime.Add(timeOffset),
		Labels:      labels,
		Annotations: annotations,
	}, nil
}
