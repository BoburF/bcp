package wrapper

import (
	"fmt"
	"log"
	"net"

	"ferxes.uz/bcp/protocol"
)

type Server struct {
	server net.Listener
}

func (s *Server) newServer(handler func(req Request, response Response)) {
	handleConnection := func(conn *net.Conn) {
		connection := *conn
		defer connection.Close()

		for {
			requestProtocol := protocol.Request{}
			err := requestProtocol.ProtocolFormatRead(conn)
			if err != nil {
				log.Println("Error occured when processing request due to ", err)
				return
			}

            request := Request{
                Type: RequestType(requestProtocol.Type),
                Origin: requestProtocol.Origin,
                Host: requestProtocol.Host,
                Resource: requestProtocol.Resource,
                Data: requestProtocol.Data,
            }
            response := Response{conn: conn}

            handler(request, response)

			connection.Close()
			break
		}
	}

	for {
		conn, err := s.server.Accept()
		if err != nil {
			log.Println("Error accepting connections due to ", err)
			return
		}
		go handleConnection(&conn)
	}
}

func (s *Server) Listen(port int) {
	server, err := net.Listen("tcp", "localhost:"+fmt.Sprint(port))
	if err != nil {
		log.Println("Can't create server due to ", err)
		return
	}
	s.server = server
	log.Println("Server listening on port:", port)
}
