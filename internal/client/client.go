package client

type Client interface {
	GetSecretVersions(secretName string) ([]Version, error)
	GetSecret(secretName string) ([]byte, error)
	GetSecretVersion(secretName, version string) ([]byte, error)
	AddSecretVersion(secretName string, payload []byte) error
	SearchInSecrets(query string) ([]SecretInfo, error)
	Secrets() ([]SecretInfo, error)
	GetSecretInfo(fullPath string) (SecretInfo, error)
}
