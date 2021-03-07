package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	notify := make(chan error)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				notify <- err
				return
			}
			if n > 0 {
				fmt.Println("unexpected data: %s", buf[:n])
			}
		}
	}()

	for {
		select {
		case err := <-notify:
			if io.EOF == err {
				fmt.Println("connection dropped\n error:", err)
				return
			}

		case <-time.After(time.Second * 1):
			fmt.Println("checked, still alive")
		}
	}
}

func main() {

	servers := []struct {
		protocol string
		addr     string
	}{
		{"tcp", ":1123"},
		{"tcp", ":6250"},
	}
	fmt.Println("Launching server...")

	for _, serv := range servers {
		ln1, _ := net.Listen(serv.protocol, serv.addr)
		fmt.Println(serv.addr)
		for {
			conn, _ := ln1.Accept()
			go handleConnection(conn)
		}
	}

}
