package wrapper_server

import (
	"fmt"
	"net"
	"sync"

	"ferxes.uz/bcp/protocol"
)

func NewServer(port int, handler func(request Request, response Response)) error {
	var wg sync.WaitGroup
	listener, err := protocol.NewServer(port)
	if err != nil {
		return err
	}
	wg.Add(1)

	go func() {
		defer wg.Done()
		acceptConnections(&listener, handler)
	}()

	return nil
}

func acceptConnections(listener *net.Listener, handler func(request Request, response Response)) {
	listen := *listener
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(&conn, handler)
	}
}

func handleConnection(conn *net.Conn, handler func(request Request, response Response)) {
	connection := *conn
	defer connection.Close()
	request := protocol.Request{}
	request.ConvertFrom(conn)

	response := protocol.Response{
        Conn: conn,
		Additions: make(map[string]string),
	}

    req, res := handleTheHandler(request, response)

    handler(req, res)
}

func handleTheHandler(request protocol.Request, response protocol.Response) (Request, Response) {
	req := Request{
		Resource:  request.Resource,
		Additions: request.Additions,
		Data:      request.Data,
	}

	res := Response{
		res:       response,
		Additions: response.Additions,
		Data:      response.Data,
	}

	return req, res
}
