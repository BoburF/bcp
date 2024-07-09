package main

import (
	"bufio"
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

	resource := "/"
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

	response := protocol.Response{}
	err = response.ConvertFrom(&conn)
	if err != nil {
		log.Println("Error from connection:", err)
		return
	}
	log.Println("response is from dialing:", response)
	reader := bufio.NewReader(response.Data)
	line, _, err := reader.ReadLine()
	if err != nil {
		log.Println("Error from connection:", err)
		return
	}
    log.Println(string(line))
}

func server() {
	server := protocol.FerxesServer{}
	err := server.NewServer(2323, handler)
	if err != nil {
		log.Println("creating server err:", err)
	}
	wg.Done()
}

func handler(req protocol.Request, res protocol.Response) {
	if req.Resource == "/" {
		res.Additions["bobur"] = "abdullayev"
		res.Data = strings.NewReader("Bobur zo'r bolasanda")
		err := res.Send()
		if err != nil {
			panic("Server can't write response")
		}
	}
}
