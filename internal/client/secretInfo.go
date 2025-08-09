package client

import "time"

type SecretInfo struct {
	Name        string
	FullPath    string
	CreateTime  time.Time
	Labels      map[string]string
	Annotations map[string]string
}