package main

import (
	"fmt"
	"log"
	"net"
	"zhenghe/server"
)

func main() {
	s := server.NewServer()
	go s.Run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("server started on :8888")

	for {
		conn, err := listener.Accept()
		fmt.Println("ok,成功连接！")
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}

		go s.NewClient(conn)
	}
}