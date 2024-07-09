package utils

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ReadInteger(reader *io.Reader) (int, error) {
	var strBuilder strings.Builder
	buf := make([]byte, 1)

    read := *reader
	for {
		_, err := read.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}
		
		if buf[0] < '0' || buf[0] > '9' {
			break
		}
		
		strBuilder.WriteByte(buf[0])
	}

	if strBuilder.Len() == 0 {
		return 0, fmt.Errorf("no digits were found")
	}

	return strconv.Atoi(strBuilder.String())
}
