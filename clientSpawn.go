package main

import (
	"net"
)

func main() {
	for i := 1; i < 53000; i++ {
		go spawnConn()
	}
}

func spawnConn() {
	conn, err := net.Dial("tcp", "10.0.59.136:19530")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}
