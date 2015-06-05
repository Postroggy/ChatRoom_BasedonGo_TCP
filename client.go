package main

import (
	"./codec"
	"bufio"
	"fmt"
	"net"
	"time"
)

var quitSemaphore chan bool

func main() {
	for i := 0; i < 5000; i++ {
		go openConn()
	}
	// b := []byte("time\n")
	// conn.Write(b)
	var msg string
	fmt.Scanln(&msg)
	<-quitSemaphore
}

func openConn() {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")

	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	fmt.Println("connected!")

	go onMessageRecived(conn)
	go sendMessage(conn)
	<-quitSemaphore
}

func sendMessage(conn *net.TCPConn) {
	for {
		time.Sleep(1 * time.Second)
		// var msg string
		// fmt.Scanln(&msg)
		b, _ := codec.Encode(string(time.Now().Year()))
		conn.Write(b)
	}
}

func onMessageRecived(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := codec.Decode(reader)
		fmt.Println(msg)
		if err != nil {
			quitSemaphore <- true
			break
		}
	}
}
