package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	rv "RV/ReverseProxy"
	rd "redis-go/redis"
)

func main() {
	var wg sync.WaitGroup

	stream, err := net.Listen("tcp4", ":19530") //Listen at port 19530
	if err != nil {
		fmt.Println(err)
		return
	}

	defer stream.Close()
	n := 1
	for {
		con, err := stream.Accept()
		//Bench test load to sever
		if err != nil {
			log.Println(err)
			return
		}

		wg.Add(1)
		go handleConnection(con, &wg)
		n++
		wg.Wait()
	}

}

func handleConnection(con net.Conn, wg *sync.WaitGroup) {
	var received rv.Package
	buffer := make([]byte, 1024)
	n, err := con.Read(buffer)
	rawdata := string(buffer[:n])
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal([]byte(rawdata), &received)
	if err != nil {
		fmt.Println(err)
		return
	}

	IP, Init := rd.CommandFunction(received)
	wg.Done()
	rv.ReProx(con, IP, Init)

}
