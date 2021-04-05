package load_balance

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"time"

	"fmt"
)

// Data Sent and Received From user (Matthew's)
type Package struct {
	AuctionID uint64
	UserID    uint64
	Command   string
	Data      struct {
		Item  string
		Value uint64
	}
	Time []time.Time
}

func copy(src net.Conn, dst net.Conn, stop chan bool) {
	io.Copy(dst, src)
	dst.Close()
	src.Close()
	// fmt.Println("Closed")
	stop <- true
	return
}

func ReProx(src net.Conn, server string, init Package) {

	dst, err := net.Dial("tcp4", server)
	if err != nil {
		src.Close()
		log.Printf("failed to dial %s: %s", server, err)
		return
	}
	stop := make(chan bool)
	var jsonData []byte
	jsonData, err = json.Marshal(init)
	fmt.Fprintf(dst, string(jsonData)+"\n")

	go copy(dst, src, stop)
	go copy(src, dst, stop)
	select {
	case <-stop:
		//fmt.Println("Connections Closed")
		return
	}
	//fmt.Println(stop)
}
