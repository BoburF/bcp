package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"

	"ferxes.uz/bcp/protocol"
	wrapper_server "ferxes.uz/bcp/wrapper/server"
)

func main() {
	go server()
	time.Sleep(2 * time.Second)

	conn, err := net.Dial("tcp", "localhost:2323")
	if err != nil {
		log.Println(err)
		panic("Creating dial")
	}

	request := protocol.Request{
		Resource:  "/",
		Additions: make(map[string]string),
		Data:      nil,
	}

	err = request.ConvertTo(1, &conn)
	if err != nil {
		log.Println(err)
		panic("Error request convertTo")
	}

	response := protocol.Response{}

	err = response.ConvertFrom(&conn)
	if err != nil {
		log.Println(err)
		panic("Error response convertFrom")
	}

	log.Println(response)

	file, err := os.Create("sample_transfered.txt")
	if err != nil {
		log.Println(err)
		panic("Error file opening")
	}

	_, err = io.Copy(file, response.Data)
	if err != nil {
		log.Println(err)
		panic("Error file copying")
	}
}

func server() {
	err := wrapper_server.NewServer(2323, handler)
	if err != nil {
		log.Println("creating server err:", err)
	}
}

func handler(req wrapper_server.Request, res wrapper_server.Response) {
	if req.Resource == "/" {
		res.Additions["bobur"] = "abdullayev"

		file, err := os.Open("sample.txt")
		if err != nil {
			panic("Server can't open file")
		}

		res.Data = file
		err = res.Send()
		if err != nil {
			panic("Server can't write response")
		}
	}
}
