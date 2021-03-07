package main

import (
	"fmt"
	"net"
	"time"
	"bufio"
	"runtime"
	"os"
	"strconv"
)

func main() {
	arguments := os.Args
	args,_ := strconv.Atoi(arguments[1])
	runtime.GOMAXPROCS(4)

	n:=1
	for n<args{
		go handleCon()
		n++
	}
	for {

	}
}

func handleCon() {
	c, _ := net.Dial("tcp", "10.0.59.139:19530")
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
	c.Close()
}
