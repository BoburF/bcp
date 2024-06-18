package protocol

import (
	"fmt"
	"log"
	"net"
)

func Dial(request Request, port int) (Response, error) {
	conn, err := net.Dial("tcp", "localhost:"+fmt.Sprint(port))
	if err != nil {
		return Response{}, err
	}

	err = request.ProtocolFormatWrite(&conn)
	if err != nil {
		return Response{}, err
	}
	log.Println("Request written")

	response := Response{}
	response.ProtocolFormatReader(&conn)

	return response, nil
}
