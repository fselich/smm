package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jaswdr/faker/v2"
	"math/rand"
	"strings"
	"time"
)

type FakeClient struct {
}

func NewFakeClient(projectId string) (FakeClient, error) {
	return FakeClient{}, nil
}

func (f FakeClient) GetSecretVersions(secretName string) []version {
	v := version{
		Name:      "test-version",
		State:     "enabled",
		Version:   1,
		FullPath:  "projects/test-project/secrets/test-secret/versions/1",
		CreatedAt: time.Now(),
	}

	return []version{v}
}

func (f FakeClient) createEnvSecret() []byte {
	fk := faker.New()
	var envLines []string

	numVars := 3 + rand.Intn(10)
	for i := 0; i < numVars; i++ {
		key := strings.ToUpper(fk.Lorem().Word()) + fmt.Sprintf("_%d", rand.Intn(1000))
		value := "\"" + fk.Lorem().Word() + "\""
		envLines = append(envLines, fmt.Sprintf("%s=%s", key, value))
	}

	return []byte(strings.Join(envLines, "\n"))
}

func (f FakeClient) createJsonSecret() []byte {
	fk := faker.New()
	js := fk.Json()

	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, []byte(js.String()), "", "  ")
	return prettyJSON.Bytes()
}

func (f FakeClient) GetSecret(secretName string) []byte {
	if rand.Intn(2) == 0 {
		return f.createEnvSecret()
	}
	return f.createJsonSecret()
}

func (f FakeClient) GetSecretVersion(secretName, version string) []byte {
	secret := []byte("Im a fake secret version")
	return secret

}

func (f FakeClient) AddSecretVersion(secretName string, payload []byte) error {
	return nil
}

func (f FakeClient) SearchInSecrets(query string) []string {
	secret1 := "Im a fake secret with query " + query
	secret2 := "Im another fake secret with query " + query
	return []string{secret1, secret2}
}

func (f FakeClient) Secrets() []string {
	fk := faker.New()
	secrets := make([]string, 10)
	for i := 0; i <= 9; i++ {
		secret := fmt.Sprintf("%s-secret", fk.Lorem().Word())
		secrets[i] = secret
	}
	return secrets
}
