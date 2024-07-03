package main

import (
	"log"
	"net"
	"strings"
	"sync"

	"ferxes.uz/bcp/protocol"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)
	go server()
	wg.Wait()

	conn, err := net.Dial("tcp", "localhost:2323")
	if err != nil {
		log.Println("Error dialing server:", err)
		return
	}

	resource := "example"
	additions := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	data := strings.NewReader("sample data")
	request := &protocol.Request{
		Resource:  resource,
		Additions: additions,
		Data:      data,
	}

	err = request.ConvertTo(1, &conn)
	if err != nil {
		log.Println("Error writing to connection:", err)
		return
	}

	buff := make([]byte, 1024)
	_, err = conn.Read(buff)
	if err != nil {
		log.Println("Error reading from a connection:", err)
		return
	}
	log.Println(string(buff))
}

func server() {
	server := protocol.FerxesServer{}
	listener, err := server.NewServer(2323)
	if err != nil {
		log.Println("Error creating server:", err)
		return
	}
	defer listener.Close()
	wg.Done()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connections:", err)
			return
		}
		defer conn.Close()

		buff := make([]byte, 1024)
		_, err = conn.Read(buff)
		if err != nil {
			log.Println("Error reading from a connection:", err)
			return
		}
		log.Println(string(buff))

		log.Println("Readed str:", string(buff))
		_, err = conn.Write(buff)
		if err != nil {
			log.Println("Error writing to connection:", err)
			return
		}
	}
}
