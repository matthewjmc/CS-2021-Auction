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
				fmt.Println("unexpected data:", buf[:n])
			}
		}
	}()

	for {
		select {
		case err := <-notify:
			if io.EOF == err {
				fmt.Println("connection dropped\n error:", err)
				break
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
	}
	for _, serv := range servers {
		ln, _ := net.Listen(serv.protocol, serv.addr)
		fmt.Println(serv)
		for {
			conn, _ := ln.Accept()
			go handleConnection(conn)
		}
	}
}
