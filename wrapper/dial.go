package wrapper

import (
	"fmt"
	"io"
	"net"

	"ferxes.uz/bcp/protocol"
)

func Dial(request Request, port int) (io.Reader, error) {
	conn, err := net.Dial("tcp", request.Host+fmt.Sprint(port))
	if err != nil {
		return nil, err
	}

    requestProtocol := protocol.Request{
        Type: protocol.RequestType(request.Type),
        Origin: request.Origin,
        Host: request.Host,
        Resource: request.Resource,
        Additions: request.Additions,
        Data: request.Data,
    }

	err = requestProtocol.ProtocolFormatWrite(&conn)
	if err != nil {
		return nil, err
	}

	responseProtocol := protocol.Response{}
	responseProtocol.ProtocolFormatReader(&conn)

	return *responseProtocol.Data, nil
}
