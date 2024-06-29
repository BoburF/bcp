package protocol

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
)

type Request struct {
	Resource  string
	Additions map[string]string
	Data      io.Reader
}

func (rq *Request) ConverTo(version int, conn *net.Conn) error {
	switch version {
	case 1:
		return rq.v1ConverTo(conn)
	default:
		return errors.New("unsupported version")
	}
}

func (rq *Request) v1ConverTo(conn *net.Conn) error {
	var buffer bytes.Buffer

	versionInfo := "v1\x00"
	buffer.WriteString(versionInfo)

	buffer.WriteString(fmt.Sprintf("%d\x00%s\x00", len(rq.Resource), rq.Resource))

	for key, value := range rq.Additions {
		buffer.WriteString(fmt.Sprintf("%d\x00%s\x00%d\x00%s", len(key), key, len(value), value))
	}

	_, err := io.Copy(*conn, &buffer)
	if err != nil {
		return err
	}

	if rq.Data != nil {
		_, err := io.Copy(*conn, rq.Data)
		if err != nil {
			return err
		}
	}

	return nil
}
