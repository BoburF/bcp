package protocol

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
)

type Request struct {
	Resource  string
	Additions map[string]string
	Data      io.Reader
}

func (rq *Request) ConvertTo(version int, conn *net.Conn) error {
	switch version {
	case 1:
		return rq.v1ConvertTo(conn)
	default:
		return errors.New("unsupported version")
	}
}

func (rq *Request) v1ConvertTo(conn *net.Conn) error {
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

func (rq *Request) ConvertFrom(version int, conn *net.Conn) error {
	reader := bufio.NewReader(*conn)

	firstByte, err := reader.ReadByte()
	if err != nil {
		return err
	}
	if firstByte != 'v' {
		return err
	}

	var versionBytes []byte
	for i := 0; i < 3; i++ {
		b, err := reader.ReadByte()
		if err != nil {
			return err
		}
		if b == '\x00' {
			break
		}
		if b < '0' || b > '9' {
			return errors.New("invalid character in version number")
		}
		versionBytes = append(versionBytes, b)
        if i == 3 && b != '\x00' {
            return errors.New("protocol version problem")
        }
	}
	if len(versionBytes) == 0 {
		return errors.New("version number not found")
	}

	versionNumber, err := strconv.Atoi(string(versionBytes))
	if err != nil {
		return err
	}

	switch versionNumber {
	case 1:
		return rq.v1ConvertFrom(reader)
	default:
		return errors.New("unsupported version")
	}
}

func (rq *Request) v1ConvertFrom(conn *bufio.Reader) error {
	return nil
}
