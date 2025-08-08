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

func (f FakeClient) GetSecretVersions(secretName string) []version {
	seed := seedFromSecretName(secretName + "_versions")
	source := rand.NewPCG(uint64(seed), uint64(seed>>32))
	rng := rand.New(source)

	numVersions := 1 + rng.IntN(5)
	versions := make([]version, numVersions)

	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < numVersions; i++ {
		versionNum := i + 1
		timeOffset := time.Duration(rng.Int64N(int64(time.Hour * 24 * 30)))

		versions[i] = version{
			Name:      fmt.Sprintf("%s-version-%d", secretName, versionNum),
			State:     "enabled",
			Version:   versionNum,
			FullPath:  fmt.Sprintf("projects/test-project/secrets/%s/versions/%d", secretName, versionNum),
			CreatedAt: baseTime.Add(timeOffset),
		}
	}

	return versions
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

func (f FakeClient) GetSecret(secretName string) []byte {
	if isSecretEnvType(secretName) {
		return f.createEnvSecret(secretName)
	}
	return f.createJsonSecret(secretName)
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

func (f FakeClient) GetSecretVersion(secretName, version string) []byte {
	if isSecretEnvType(secretName) {
		return f.createEnvSecretVersion(secretName, version)
	}
	return f.createJsonSecretVersion(secretName, version)
}

func (f FakeClient) AddSecretVersion(secretName string, payload []byte) error {
	return nil
}

func (f FakeClient) SearchInSecrets(query string) []string {
	seed := seedFromSecretName("search_" + query)
	source := rand.NewPCG(uint64(seed), uint64(seed>>32))
	rng := rand.New(source)
	fk := faker.NewWithSeed(source)

	numResults := 2 + rng.IntN(5)
	results := make([]string, numResults)

	for i := 0; i < numResults; i++ {
		secretName := fmt.Sprintf("%s-%s-secret", query, fk.Lorem().Word())
		results[i] = secretName
	}

	return results
}

func (f FakeClient) Secrets() []string {
	seed := int64(12345)
	source := rand.NewPCG(uint64(seed), uint64(seed>>32))
	fk := faker.NewWithSeed(source)

	secrets := make([]string, 30)
	for i := 0; i <= 29; i++ {
		secret := fmt.Sprintf("%s-secret", fk.Lorem().Word())
		secrets[i] = secret
	}
	return secrets
}
