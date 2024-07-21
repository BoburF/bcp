package wrapper_dial

import (
	"io"
	"net"
)

type Response struct {
	Additions map[string]string
	Data      io.Reader
  Connection net.Conn
}
