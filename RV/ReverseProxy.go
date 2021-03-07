package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"
)

func main() {
	for {
		In, err := net.Listen("tcp", "localhost:7777") //prox

		if err != nil {
			log.Fatalf("failed to setup listener %v", err)
		}

		log.Println("Accepted from," + "localhost:7777")
		log.Println("//////////////WATING FOR CONNECTION//////////////")

		for {

			Inconn, err := In.Accept()
			if err != nil {
				log.Fatalf("failed to accept listener %v", err)
			}
			log.Print("Connection Accepted")

			go reverse_proxy(Inconn)
		}
	}
}

func reverse_proxy(src net.Conn) {
	//setup reader
	defer src.Close()
	err := src.SetDeadline(time.Now().Add(5 * time.Second))

	if err != nil {
		log.Fatalf("No Reverse Proxy deadline %v", err)
	}

	//Buffer
	requestBuf := make([]byte, 512)
	reqLen, _ := src.Read(requestBuf)

	_ = reqLen

	log.Println("We Sent Request:", string(requestBuf))

	//read the request

	dst, _ := net.Dial("tcp", "localhost:9000")

	defer dst.Close()
	err = dst.SetDeadline(time.Now().Add(5 * time.Second))

	if err != nil {
		log.Fatalf("No Reverse Proxy deadline %v", err)
	}

	log.Print("ReverseProxy Connected")

	var backendBuf bytes.Buffer
	ndst := io.TeeReader(dst, &backendBuf)
	io.Copy(dst, bytes.NewReader(requestBuf)) //server, prox
	io.Copy(src, ndst)

	//read response
	backendBytes, _ := ioutil.ReadAll(&backendBuf)
	log.Println("We got response:", string(backendBytes))
}
