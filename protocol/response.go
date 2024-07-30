package protocol

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/BoburF/bcp/utils"
)

type Response struct {
	Conn      *net.Conn
	Additions map[string]string
	Data      io.Reader
}

func (rs *Response) ConvertTo(version int) error {
	switch version {
	case 1:
		return rs.v1ConvertTo(rs.Conn)
	default:
		return errors.New("unsupported version")
	}
}

func (rs *Response) v1ConvertTo(conn *net.Conn) error {
	var buffer bytes.Buffer

	versionInfo := "v1\x00"
	buffer.WriteString(versionInfo)

	buffer.WriteString(fmt.Sprintf("%d\x00", len(rs.Additions)))
	for key, value := range rs.Additions {
		buffer.WriteString(fmt.Sprintf("%d\x00%s\x00%d\x00%s\x00", len(key), key, len(value), value))
	}

	_, err := io.Copy(*conn, &buffer)
	if err != nil {
		return err
	}

	if rs.Data != nil {
		_, err := io.Copy(*conn, rs.Data)
		if err != nil {
			return err
		}
	}
    log.Println("Reqponse is written")

	return nil
}

func (rs *Response) ConvertFrom(conn *net.Conn) error {
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
		return rs.v1ConvertFrom(conn)
	default:
		return errors.New("unsupported version")
	}
}

func (rs *Response) v1ConvertFrom(conn *net.Conn) error {
	var reader io.Reader = *conn

	additionsLength, err := utils.ReadInteger(&reader)
	if err != nil {
		return err
	}

	rs.Additions = make(map[string]string)
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
		_, err = reader.Read(valueBuff)
		if err != nil {
			return err
		}
		err = utils.CheckForNullByte(conn)
		if err != nil {
			return err
		}

		rs.Additions[string(keyBuff)] = string(valueBuff)
	}

	rs.Data = reader

	return nil
}
