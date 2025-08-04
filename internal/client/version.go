package client

import "time"

type version struct {
	Name      string
	State     string
	Version   int
	FullPath  string
	CreatedAt time.Time
}
