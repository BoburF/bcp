package wrapper_server

import "io"

type Request struct {
	Resource  string
	Additions map[string]string
	Data      io.Reader
}
