package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
	"time"
	//"sync"
	//"reflect"
	//rv "load_balance/reverseproxy"
)

var S1_Usage, S2_Usage string

// this will locate on the loadbalance server to listen for cpu usage from S1/S2
// every period of time set
func main() {
	//var wg sync.WaitGroup
	conn, err := net.Listen("tcp4", ":19530")
	if err != nil {
		fmt.Println(err)
	}
	for {
		c, err := conn.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		//wg.Add(1)
		go getStat(c)
		//wg.Wait()
	}

}

func getStat(conn net.Conn) {
	defer conn.Close()
	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Exiting TCP server!")
		}

		usage := string(netData)
		fmt.Println(usage)
		if usage[12] == 116 {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(usage, -1)
			for _, element := range submatchall {
				S2_Usage = element
				fmt.Println(S2_Usage)
			}
		}
		if usage[12] == 111 {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(usage, -1)
			for _, element := range submatchall {
				S1_Usage = element
				fmt.Println(S1_Usage)
			}
		}

		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		conn.Write([]byte(myTime))
		continue
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
