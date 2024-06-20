package wrapper

import "io"

type RequestType string

const (
	info RequestType = "info"
	data RequestType = "data"
)

type Request struct {
	Type      RequestType       `json:"type"`
	Origin    string            `json:"origin"`
	Host      string            `json:"host"`
	Resource  string            `json:"resource"`
	Additions map[string]string `json:"additions"`
	Data      *io.Reader        `json:"data"`
}
