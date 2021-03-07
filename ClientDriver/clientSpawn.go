package main

import (
	"fmt"
	"net"
	"time"
	"bufio"
)

func main() {
	handleCon()

}

func handleCon() {
	c, _ := net.Dial("tcp", "10.0.59.140:7777")
	err := c.(*net.TCPConn).SetKeepAlive(true)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		fmt.Fprintf(c,"TestTestTest\n")
		data, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Println("From -->", data)
		time.Sleep(1*time.Second)
		
		
	}
}
