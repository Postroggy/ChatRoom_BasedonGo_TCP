package main

import "fmt"

func main() {

	var mp12 map[string]int
	mp12["A"] = 1
	mp12["B"] = 2
	mp12["C"] = 3

	for name,_ := range mp12 {
		fmt.Println(name)
	}
}
