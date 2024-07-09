package wrapper_server

import (
	"io"

	"ferxes.uz/bcp/protocol"
)

type Response struct {
	res       protocol.Response
	Additions map[string]string
	Data      io.Reader
}

func (rs *Response) Send() error {
	return rs.res.ConvertTo(1)
}
