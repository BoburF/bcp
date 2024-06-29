package protocol

import (
	"errors"
	"io"
)

type Response struct {
	Additions map[string]string
	Data      io.Reader
}

func (rs *Response) ConverTo(version int) (io.Writer, error) {
	switch version {
	case 1:
		return rs.v1ConverTo()
	default:
		return nil, errors.New("unsupported version")
	}
}

func (rs *Response) v1ConverTo() (io.Writer, error) {
	return nil, nil
}
