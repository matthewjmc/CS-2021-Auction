package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	c, err := net.Dial("tcp", "10.0.59.136:19530")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		//reader := bufio.NewReader(os.Stdin)
		//text, _ := reader.ReadString('\n')

		fmt.Fprintf(c, time.Nanosecond.String()+"\n") //Print to server

	}
}
