package main

import (
	"log"
	"strings"

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
		res.Data = strings.NewReader("Bobur zo'r bolasanda")
		err := res.Send()
		if err != nil {
			panic("Server can't write response")
		}
	}
}
