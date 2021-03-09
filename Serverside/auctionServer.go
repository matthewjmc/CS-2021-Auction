package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	stream, err := net.Listen("tcp", ":19530") //Listen at port 19530
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stream.Close()
	var count int64 = 1
	for {
		con, err := stream.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println("Conn Number:", count,"-->",con.RemoteAddr().String())
		go requestHandle(con, count)
		count++
		fmt.Println(count)
	}
}

func requestHandle(con net.Conn, cnt int64) {
	for {
		_, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		go returnData(con)
	}
	con.Close()
}

func returnData(con net.Conn){
	fmt.Fprintf(con, con.RemoteAddr().String()+"\n")
}