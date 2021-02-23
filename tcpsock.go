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

	for {
		con, err := stream.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go requestHandle(con)
	}

}

func requestHandle(con net.Conn) {
	for {
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(data)
	}
	con.Close()
}
