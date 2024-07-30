package protocol

import (
	"fmt"
	"net"
)

func NewServer(host string, port int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
}
