package main

import (
	"net"
)

func main() {
	i := 1
	for i < 52000 {
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
