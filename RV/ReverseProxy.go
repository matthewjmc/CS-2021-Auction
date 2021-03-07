package main

import (
	"io"
	"log"
	"net"
)

func main() {
	front := "localhost:7000"
	back := "localhost:9000"

	In, err := net.Listen("tcp", front)
	if err != nil {
		log.Fatalf("failed to setup listener: %v", err)
	}

	log.Printf("listening on %s, balancing %s", front, back)
	log.Println("//////////////////LISTENING//////////////////")

	for {
		Inconn, err := In.Accept()
		if err != nil {
			log.Printf("failed to accept: %s", err)
			continue
		}
		go ReProx(Inconn, back)
	}
}

func ReProx(src net.Conn, server string) {
	dst, err := net.Dial("tcp", server)
	if err != nil {
		src.Close()
		log.Printf("failed to dial %s: %s", server, err)
		return
	}

	io.Copy(dst, src)
	io.Copy(src, dst)

}
