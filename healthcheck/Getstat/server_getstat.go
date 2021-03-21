package main

import (
	"bufio"
	"fmt"
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
