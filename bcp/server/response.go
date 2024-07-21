package bcp

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
    rs.res.Additions = rs.Additions
    rs.res.Data = rs.Data
	return rs.res.ConvertTo(1)
}
