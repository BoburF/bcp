package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	bcp_client "github.com/BoburF/bcp/bcp/client"
	bcp "github.com/BoburF/bcp/bcp/server"
)

func main() {
	go server()
	time.Sleep(2 * time.Second)

	request := bcp_client.Request{
		Resource:  "/",
		Additions: make(map[string]string),
		Data:      nil,
	}
	response, err := bcp_client.Dial("localhost", 2323, request)
	if err != nil {
		log.Println("Dial:", err)
		return
	}
	defer response.Connection.Close()

	file, err := os.Create(fmt.Sprintf("sample_received.%s", response.Additions["format"]))
	if err != nil {
		log.Println("Response:", err)
		return
	}

	n, err := io.Copy(file, response.Data)
	if err != nil {
		log.Println("File:", err)
		return
	}

	log.Println("this many bytes:", n)
}

func server() {
	err := bcp.NewServer("0.0.0.0", 2323, handler)
	if err != nil {
		log.Println("creating server err:", err)
	}
}

func handler(req bcp.Request, res bcp.Response) {
	if req.Resource == "/" {
		res.Additions["bobur"] = "abdullayev"
		res.Additions["format"] = "txt"

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
