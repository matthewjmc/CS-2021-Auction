package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// this will locate on the loadbalance server to listen for cpu usage from S1/S2
// every period of time set
func main() {
	for {
		conn, err := net.Listen("tcp4", ":19530")
		if err != nil {
			fmt.Println(err)
		}
		c, err := conn.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			netData, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			if strings.TrimSpace(string(netData)) == "STOP" {
				fmt.Println("Exiting TCP server!")
				return
			}
			fmt.Print(string(netData))
			t := time.Now()
			myTime := t.Format(time.RFC3339) + "\n"
			c.Write([]byte(myTime))
		}
	}
}

func CheckEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		CheckAlive(x)
	}
}

func CheckAlive(t time.Time) {
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
