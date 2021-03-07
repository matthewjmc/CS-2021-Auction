package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

var count int = 0

func main() {
	checkEvery(15*time.Millisecond, checkConnection1)
	checkEvery(20*time.Millisecond, checkConnection2)
}

func checkEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		checkConnection1(x)
		checkConnection2(x)
	}
}

func checkConnection1(t time.Time) {
	conn, _ := net.Dial("tcp", ":1123")
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
				fmt.Println("unexpected data:", buf[:n])
			}
		}
	}()
	select {
	case err := <-notify:
		fmt.Println(time.Now(), "connection1 dropped:", err)
		return
	case <-time.After(time.Second * 1):
		fmt.Println(time.Now(), "timeout1, still alive")
	}
	defer conn.Close()
}

func checkConnection2(t time.Time) {
	conn, _ := net.Dial("tcp", ":6250")
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
				fmt.Println("unexpected data:", buf[:n])
			}
		}
	}()
	select {
	case err := <-notify:
		fmt.Println(time.Now(), "connection2 dropped:", err)
		return
	case <-time.After(time.Second * 1):
		fmt.Println(time.Now(), "timeout2, still alive")
	}
	defer conn.Close()
}
