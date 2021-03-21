package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// use notify to return error
	notify := make(chan error)
	go func() {
		// create buf with all 0 values
		buf := make([]byte, 1024)
		for {
			// read through buf
			n, err := conn.Read(buf)
			// if can't read then print error
			if err != nil {
				notify <- err
				return
			}
			// if find 1 in buf
			if n > 0 {
				fmt.Println("unexpected data:", buf[:n])
			}
		}
	}()

	for {
		select {
		// define case of error
		case err := <-notify:
			// if end-of-file then break
			if io.EOF == err {
				fmt.Println(time.Now(), "connection dropped\n error:", err)
				break
			}
		// define successful case when check every seconds
		case <-time.After(time.Second * 1):
			fmt.Println(time.Now(), "checked, alive")
		}
	}
}

func main() {
	// create a struct to store server info
	servers := []struct {
		protocol string
		addr     string
	}{
		{"tcp", ":1123"},
		{"tcp", ":6250"},
	}
	for {
		// set listen on port
		ln, _ := net.Listen(servers[0].protocol, servers[0].addr)
		ln2, _ := net.Listen(servers[1].protocol, servers[1].addr)
		fmt.Println(ln, ln2)
		for {
			conn1, _ := ln.Accept()
			handleConnection(conn1)
			conn2, _ := ln2.Accept()
			go handleConnection(conn2)
		}
	}
}
