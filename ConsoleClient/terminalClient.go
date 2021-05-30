package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	var inputString *bufio.Reader
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
	}
	service := os.Args[1]
	// 绑定
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)
	// 连接
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	rAddr := conn.RemoteAddr()
	fmt.Println("connected!")
	for {
		//接收
		go ReplyFromServer(rAddr, conn)

		// 发送
		inputString = bufio.NewReader(os.Stdin)
		msg, err := inputString.ReadString('\n')
		checkError(err)
		_, err = conn.Write([]byte(msg))
		checkError(err)
		if msg == "/quit\n" {
			os.Exit(0)
		}
	}
	conn.Close()
	os.Exit(0)
}
func ReplyFromServer(rAddr net.Addr, conn *net.TCPConn) {
	var buf [512]byte
	n, err := conn.Read(buf[0:])
	checkError(err)
	fmt.Println("Reply:", string(buf[0:n]))
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
