package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type RequestType string

const (
	info RequestType = "info"
	data RequestType = "data"
)

type Request struct {
	Type      RequestType       `json:"type"`
	Origin    string            `json:"origin"`
	Host      string            `json:"host"`
	Resource  string            `json:"resource"`
	Additions map[string]string `json:"additions"`
	Data      *io.Reader        `json:"data"`
}

func (request *Request) ProtocolFormatWrite(conn *net.Conn) error {
	reqEssential := fmt.Sprintf("%s\x00%d\x00%s\x00%d\x00%s\x00%d\x00%s\x00%d",
		request.Type,
		len(request.Origin), request.Origin,
		len(request.Host), request.Host,
		len(request.Resource), request.Resource,
		len(request.Additions),
	)

	var additions strings.Builder

	for key, value := range request.Additions {
		additions.WriteString("\x00")
		additions.WriteString(fmt.Sprint(len(key)))
		additions.WriteString(":")
		additions.WriteString(strings.ToLower(key))
		additions.WriteString(":")
		additions.WriteString(fmt.Sprint(len(value)))
		additions.WriteString(":")
		additions.WriteString(strings.ToLower(value))
	}
	reqEssential += additions.String()

	connection := *conn
	_, err := connection.Write([]byte(reqEssential))
	if err != nil {
		return err
	}

	if request.Type == "data" {
		var requestData strings.Builder
		requestData.WriteString("\x00")
		_, err := io.Copy(connection, *request.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (request *Request) ProtocolFormatRead(conn *net.Conn) error {
	reader := bufio.NewReader(*conn)

	requestTypeBytes := make([]byte, 4)
	_, err := io.ReadFull(reader, requestTypeBytes)
	if err != nil {
		return err
	}
	request.Type = RequestType(requestTypeBytes)

	if nullByte, _ := reader.ReadByte(); nullByte != byte('\x00') {
		return errors.New("Format of request is wrong. Request type should be 4 byte length." + string(nullByte))
	}

	for i := 0; i < 3; i++ {
		length, err := reader.ReadString('\x00')
		if err != nil {
			return err
		}
		reader.UnreadByte()
		length = length[:len(length)-1]

		lengthInt, err := strconv.ParseInt(length, 10, 64) // base 10, 64-bit integer
		if err != nil {
			return err
		}
		if nullByte, _ := reader.ReadByte(); nullByte != byte('\x00') {
			return errors.New("Format of request is wrong. First should be a 4 byte length of request type.")
		}

		requestEssential := make([]byte, lengthInt)
		_, err = io.ReadFull(reader, requestEssential)
		if err != nil {
			return err
		}
		if nullByte, _ := reader.ReadByte(); nullByte != byte('\x00') {
			return errors.New("Format of request is wrong. First should be a 4 byte length of request type.")
		}

		switch i {
		case 0:
			request.Origin = string(requestEssential)
		case 1:
			request.Host = string(requestEssential)
		case 2:
			request.Resource = string(requestEssential)
		}
	}

	additionalsLength, err := reader.ReadString('\x00')
	if err != nil {
		return err
	}
	reader.UnreadByte()
	additionalsLength = additionalsLength[:len(additionalsLength)-1]

	additionsLengthInt, err := strconv.ParseInt(additionalsLength, 10, 64) // base 10, 64-bit integer
	if err != nil {
		return err
	}

	if additionsLengthInt > 0 {
		request.Additions = make(map[string]string)
	}
	if nullByte, _ := reader.ReadByte(); nullByte != byte('\x00') {
		return errors.New("Format of request is wrong. Additions header.")
	}

	for i := 0; i < int(additionsLengthInt); i++ {
		keyLength, err := reader.ReadString(':')
		if err != nil {
			return err
		}
		reader.UnreadByte()
		keyLength = keyLength[:len(keyLength)-1]

		if nullByte, _ := reader.ReadByte(); nullByte != byte(':') {
			return errors.New("Format of request is wrong. Check the Additions")
		}

		keyLengthInt, err := strconv.ParseInt(keyLength, 10, 64) // base 10, 64-bit integer
		if err != nil {
			return err
		}
		key := make([]byte, keyLengthInt)
		_, err = io.ReadFull(reader, key)
		if err != nil {
			return err
		}

		if nullByte, _ := reader.ReadByte(); nullByte != byte(':') {
			return errors.New("Format of request is wrong. Check the Additions")
		}

		valueLength, err := reader.ReadString(':')
		if err != nil {
			return err
		}
		reader.UnreadByte()
		valueLength = valueLength[:len(valueLength)-1]

		valueLengthInt, err := strconv.ParseInt(valueLength, 10, 64) // base 10, 64-bit integer
		if err != nil {
			return err
		}

		if nullByte, _ := reader.ReadByte(); nullByte != byte(':') {
			return errors.New("Format of request is wrong. Check the Additions")
		}

		value := make([]byte, valueLengthInt)
		_, err = io.ReadFull(reader, value)
		if err != nil {
			return err
		}

		request.Additions[string(key)] = string(value)
	}

	if request.Type == RequestType("data") {
		if nullByte, _ := reader.ReadByte(); nullByte != byte('\x00') {
			return errors.New("Format of request is wrong. First should be a 4 byte length of request type.")
		}
		var dataReader io.Reader = reader
		request.Data = &dataReader
	}

	return nil
}
