package utils

import (
	"errors"
	"io"
	"net"
)

func CheckForNullByte(conn *net.Conn) error {
	var reader io.Reader = *conn
	nullBuff := make([]byte, 1)

    _, err := reader.Read(nullBuff)
	if err != nil {
		return err
	}
	if nullBuff[0] != '\x00' {
		return errors.New("Format of request is wrong")
	}

    return nil
}
