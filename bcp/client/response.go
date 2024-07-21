package bcp_client

import (
	"io"
	"net"
)

type Response struct {
	Additions map[string]string
	Data      io.Reader
  Connection net.Conn
}
