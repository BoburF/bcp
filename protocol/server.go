package protocol

import (
	"fmt"
	"net"
)

func NewServer(port int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
}
