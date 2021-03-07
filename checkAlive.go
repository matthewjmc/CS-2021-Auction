package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func main() {

	servers := []struct {
		protocol string
		addr     string
	}{
		{"tcp", ":1123"},
		{"tcp", ":6250"},
	}
	for _, serv := range servers {
		conn, _ := net.Dial(serv.protocol, serv.addr)
		err := conn.(*net.TCPConn).SetKeepAlive(true)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = conn.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)
		if err != nil {
			fmt.Println(err)
			return
		}
		notify := make(chan error)

		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := conn.Read(buf)
				if err != nil {
					notify <- err
					if io.EOF == err {
						close(notify)
						return
					}
				}
				if n > 0 {
					fmt.Println("unexpected data: %s", buf[:n])
				}
			}
		}()

		for {
			select {
			case err := <-notify:
				fmt.Println("connection dropped message", err)
				return
			case <-time.After(time.Second * 1):
				fmt.Println("timeout 1, still alive")
			}
		}
	}

}
