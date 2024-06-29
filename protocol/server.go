package protocol

import (
	"fmt"
	"net"
)

type FerxesServer struct{}

func (fs *FerxesServer) newServer(port int) (net.Listener, error){
  return net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
}
