package bcp_client

import (
	"fmt"
	"net"

	"github.com/BoburF/bcp/protocol"
)

func Dial(host string, port int, request Request) (Response, error) {
	response := Response{}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return response, err
	}

	requestProtocol := protocol.Request{
		Resource:  request.Resource,
		Additions: request.Additions,
		Data:      request.Data,
	}
	err = requestProtocol.ConvertTo(1, &conn)
	if err != nil {
		return response, err
	}

	responseProtocol := protocol.Response{}
	err = responseProtocol.ConvertFrom(&conn)
	if err != nil {
		return response, err
	}

  response.Additions = responseProtocol.Additions
  response.Data = responseProtocol.Data
  response.Connection = conn
	return response, nil
}
