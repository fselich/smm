package client

type Client interface {
	GetSecretVersions(secretName string) []version
	GetSecret(secretName string) []byte
	GetSecretVersion(secretName, version string) []byte
	AddSecretVersion(secretName string, payload []byte) error
	SearchInSecrets(query string) []string
	Secrets() []string
}
