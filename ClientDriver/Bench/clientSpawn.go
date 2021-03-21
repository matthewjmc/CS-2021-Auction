package main

import (
	"fmt"
	"net"
	"time"

	//"bufio"
	"encoding/json"
	"math/rand"
	"os"
	"runtime"
	"strconv"
)

type Package struct {
	AuctionID int
	UserID    int
	Command   string
	Data      struct {
		Item  string
		Value int
	}
}

var serverIP = "10.104.0.9:19530"

func main() {
	arguments := os.Args
	args, _ := strconv.Atoi(arguments[1])
	runtime.GOMAXPROCS(4) //Use 4 Cores
	n := 0
	for n < args {
		user := Package{}
		user.AuctionID = rand.Intn(10-1) + 1
		user.UserID = rand.Intn(2000-1) + 1
		user.Data.Item = "Price"
		user.Data.Value = rand.Intn(100000-10) + 10

		go handleCon(user)
		n++
		fmt.Println(n)
	}
	for {

	}
}

func handleCon(data Package) {
	connection, _ := net.Dial("tcp", "10.0.59.139:19530")
	defer connection.Close()
	err := connection.(*net.TCPConn).SetKeepAlive(true)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = connection.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	//Convert Struct to JSON Document
	var jsonData []byte
	jsonData, err = json.Marshal(data)
	fmt.Println("AuctionID:", data.AuctionID)
	if err != nil {
		fmt.Println(err)
	}
	//Transmit and Receive
	for {
		fmt.Fprintf(connection, string(jsonData)+"\n")
		//data, _ := bufio.NewReader(connection).ReadString('\n')
		//fmt.Println("From -->", data)
		time.Sleep(1 * time.Second)
	}

}
