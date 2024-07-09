package protocol

import (
	"fmt"
	"net"
	"sync"
)

type FerxesServer struct{}

func (fs *FerxesServer) NewServer(port int, handler func(request Request, response Response)) error {
	var wg sync.WaitGroup
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}
	wg.Add(1)

	go func() {
		defer wg.Done()
		fs.acceptConnections(&listener, &handler)
	}()

	return nil
}

func (fs *FerxesServer) acceptConnections(listener *net.Listener, handler *func(request Request, response Response)) {
	listen := *listener
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go fs.handleConnection(&conn, handler)
	}
}

func (fs *FerxesServer) handleConnection(conn *net.Conn, handler *func(request Request, response Response)) {
	connection := *conn
	handle := *handler
	defer connection.Close()
	request := Request{}
	request.ConvertFrom(conn)

	response := Response{
		conn:      conn,
		Additions: make(map[string]string),
	}
	handle(request, response)
}

func (rs *Response) Send() error {
	err := rs.ConvertTo(1, rs.conn)
    if err != nil {
        return err
    }
    return nil
}
