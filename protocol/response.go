package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	CacheTime time.Duration     `json:"cacheTime"`
	ServedBy  string            `json:"servedBy"`
	Additions map[string]string `json:"additions"`
	Data      *io.Reader        `json:"data"`
}

func (response *Response) ProtocolFormatWrite(conn *net.Conn) error {
	resEssential := fmt.Sprintf("%d\x00%s\x00%d\x00%s\x00%d",
		len(response.CacheTime.String()), response.CacheTime,
		len(response.ServedBy), response.ServedBy,
		len(response.Additions),
	)
	var additions strings.Builder

	for key, value := range response.Additions {
		additions.WriteString("\x00")
		additions.WriteString(fmt.Sprint(len(key)))
		additions.WriteString(":")
		additions.WriteString(strings.ToLower(key))
		additions.WriteString(":")
		additions.WriteString(fmt.Sprint(len(value)))
		additions.WriteString(":")
		additions.WriteString(strings.ToLower(value))
	}
	resEssential += additions.String()

	connection := *conn
	_, err := connection.Write([]byte(resEssential))
	if err != nil {
		return err
	}

	var responseData strings.Builder
	responseData.WriteString("\x00")
	_, err = io.Copy(connection, *response.Data)
	if err != nil {
		return err
	}

	return nil
}

func (response *Response) ProtocolFormatReader(conn *net.Conn) error {
	reader := bufio.NewReader(*conn)
	for i := 0; i < 2; i++ {
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

		responseEssential := make([]byte, lengthInt)
		_, err = io.ReadFull(reader, responseEssential)
		if err != nil {
			return err
		}
		if nullByte, _ := reader.ReadByte(); nullByte != byte('\x00') {
			return errors.New("Format of request is wrong. First should be a 4 byte length of request type.")
		}

		switch i {
		case 0:
			durationStr := string(responseEssential)
			response.CacheTime, err = time.ParseDuration(durationStr)
			if err != nil {
				return err
			}
		case 1:
			response.ServedBy = string(responseEssential)
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

	log.Println("Length of additionals in response:", additionalsLength)

	if additionsLengthInt > 0 {
		response.Additions = make(map[string]string)
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

		response.Additions[string(key)] = string(value)
	}

	if nullByte, _ := reader.ReadByte(); nullByte != byte('\x00') {
		return errors.New("Format of request is wrong. First should be a 4 byte length of request type.")
	}
	var dataReader io.Reader = reader
	response.Data = &dataReader

	return nil
}
