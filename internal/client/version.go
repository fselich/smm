package client

import "time"

type Version struct {
	Name      string
	State     string
	Version   int
	FullPath  string
	CreatedAt time.Time
}
