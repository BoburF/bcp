package protocol

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"

	"ferxes.uz/bcp/utils"
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

	buffer.WriteString(fmt.Sprintf("%d\x00", len(rq.Additions)))
	for key, value := range rq.Additions {
		buffer.WriteString(fmt.Sprintf("%d\x00%s\x00%d\x00%s\x00", len(key), key, len(value), value))
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

func (rq *Request) ConvertFrom(conn *net.Conn) error {
	firstByte := make([]byte, 1)
	readBytes := *conn
	_, err := readBytes.Read(firstByte)
	if err != nil {
		return err
	}
	if firstByte[0] != 'v' {
		return err
	}

	var reader io.Reader = *conn
	versionBytes, err := utils.ReadInteger(&reader)

	switch versionBytes {
	case 1:
		return rq.v1ConvertFrom(conn)
	default:
		return errors.New("unsupported version")
	}
}

func (rq *Request) v1ConvertFrom(conn *net.Conn) error {
	var reader io.Reader = *conn

	resourceLength, err := utils.ReadInteger(&reader)
	if err != nil {
		return err
	}

	resourceBuff := make([]byte, resourceLength)
	_, err = reader.Read(resourceBuff)
	if err != nil {
		return err
	}
	rq.Resource = string(resourceBuff)

	err = utils.CheckForNullByte(conn)
	if err != nil {
		return err
	}

	additionsLength, err := utils.ReadInteger(&reader)
	if err != nil {
		return err
	}
	
    rq.Additions = make(map[string]string)
	for i := 0; i < additionsLength; i++ {
		keyLength, err := utils.ReadInteger(&reader)
		if err != nil {
			return err
		}

		keyBuff := make([]byte, keyLength)
		_, err = reader.Read(keyBuff)
		if err != nil {
			return err
		}
		err = utils.CheckForNullByte(conn)
		if err != nil {
			return err
		}

		valueLength, err := utils.ReadInteger(&reader)
		if err != nil {
			return err
		}

		valueBuff := make([]byte, valueLength)
		_, err = reader.Read(valueBuff )
		if err != nil {
			return err
		}
		err = utils.CheckForNullByte(conn)
		if err != nil {
			return err
		}

        rq.Additions[string(keyBuff)] = string(valueBuff)
	}

    rq.Data = reader

	return nil
}
