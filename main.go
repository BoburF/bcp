package main

import (
	"log"
	"os"

	wrapper_server "ferxes.uz/bcp/wrapper/server"
)

func main() {
	server()
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
