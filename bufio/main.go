package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	var inputString *bufio.Reader
	inputString = bufio.NewReader(os.Stdin)
	for {
		msg, _ := inputString.ReadString('\n')
		result := (msg == "quit\n")
		fmt.Println(result)

		if result {
			os.Exit(0)
		}
	}
}
