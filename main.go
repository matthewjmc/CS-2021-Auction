package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	//"time"
	rd "load_balance/redis"
	rv "load_balance/reverseproxy"
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
		//fmt.Println(n)
		wg.Add(1)
		go handleConnection(con, &wg)
		n++
		wg.Wait()
	}

}

func handleConnection(con net.Conn, wg *sync.WaitGroup) {
	var received rv.Package
	//defer con.Close()
	buffer := make([]byte, 1024)
	//fmt.Println("start")
	n, err := con.Read(buffer)
	rawdata := string(buffer[:n])
	//fmt.Println(rawdata)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal([]byte(rawdata), &received)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(received)
	// fmt.Println(reflect.TypeOf(received))
	IP, Init := rd.CommandFunction(received)

	wg.Done()
	//Init.Time = append([]time.Time{time.Now()})
	//fmt.Println(Init)
	rv.ReProx(con, IP, Init)
	//fmt.Println("end")
}
