package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"ferxes.uz/bcp/protocol"
)

const (
	port = 4343
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)
	go serverHandle()

	additions := make(map[string]string)
	additions["key"] = "value"
	request := protocol.Request{
		Type:      "info",
		Origin:    "bobur.uz",
		Host:      "server.host",
		Resource:  "/",
		Additions: additions,
	}

	response, err := protocol.Dial(request, port)
	if err != nil {
		log.Println("Error during dialing due to ", err)
		return
	}

	rsp, err := json.Marshal(response)
	if err != nil {
		log.Println("json error ", err)
	}
	log.Println("Response is this ", string(rsp))

	wg.Wait()
}

func serverHandle() {
	server, err := net.Listen("tcp", "localhost:"+fmt.Sprint(port))
	if err != nil {
		log.Println("Can't create server due to ", err)
		return
	}
	log.Println("Server listening on port:", port)
	wg.Done()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("Error accepting connections due to ", err)
			return
		}

		wg.Add(1)
		handleConnection(&conn)
	}
}

func handleConnection(conn *net.Conn) {
	connection := *conn
	log.Println("Connection accepted")
	defer connection.Close()
	defer wg.Done()

	for {
		requestProtocol := protocol.Request{}
		err := requestProtocol.ProtocolFormatRead(conn)
		if err != nil {
			log.Println("Error occured when processing request due to ", err)
			return
		}

		rp, err := json.Marshal(requestProtocol)
		if err != nil {
			log.Println("json error ", err)
		}

		log.Println("Request is this ", string(rp))

		additions := make(map[string]string)
		additions["key"] = "value"
		var stringReader io.Reader = strings.NewReader("Boburbyte")
		response := protocol.Response{
			CacheTime: 1 * time.Millisecond,
			ServedBy:  "boburs first server",
			Additions: additions,
			Data:      &stringReader,
		}

		response.ProtocolFormatWrite(conn)
		connection.Close()
		break
	}
}
