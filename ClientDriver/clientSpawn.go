package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	for {
		c, err := net.Dial("tcp", "10.0.59.136:19530")
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleCon(c)
	}
}

func handleCon(c net.Conn) {
	for {
		fmt.Fprintf(c, time.Nanosecond.String()+"\n") //Print to server
		time.Sleep(1 * time.Second)
	}
}
