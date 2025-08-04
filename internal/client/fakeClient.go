package client

import "time"

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

func (f FakeClient) GetSecret(secretName string) []byte {
	secret := []byte("Im a fake secret of " + secretName)
	return secret
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
	secret1 := "secret-1"
	secret2 := "secret-2"
	return []string{secret1, secret2}
}
