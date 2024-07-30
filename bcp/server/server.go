package bcp

import (
	"fmt"
	"log"
	"net"

	"github.com/BoburF/bcp/protocol"
)

func NewServer(port int, handler func(request Request, response Response)) error {
	listener, err := protocol.NewServer(port)
	defer listener.Close()
	if err != nil {
		return err
	}
	log.Println("Server started working on port:", port)

	acceptConnections(listener, handler)

	return nil
}

func acceptConnections(listener net.Listener, handler func(request Request, response Response)) {
	for {
		conn, err := listener.Accept()
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
		Conn:      conn,
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
