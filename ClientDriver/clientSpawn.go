package main

import (
	"fmt"
	"net"
	"time"
	"bufio"
	"runtime"
	"os"
	"strconv"
	"encoding/json"
	"math/rand"
)
type Package struct{
	AuctionID int
	UserID int
	Command  string
	Data struct {
		Item string
		Value int
	}
}

func main() {
	arguments := os.Args
	args,_ := strconv.Atoi(arguments[1])
	runtime.GOMAXPROCS(4) //Use 4 Cores
	n:=1
	for n<args{
		user := Package{}
		user.AuctionID = rand.Intn(1500 - 1) + 1
		user.UserID = rand.Intn(2000 - 1) + 1
		user.Data.Item = "Price"
		user.Data.Value = rand.Intn(100000 - 10) + 10
		go handleCon(user)
		n++
	}
	for {

	}
}

func handleCon(data Package) {
	connection, _ := net.Dial("tcp", "10.0.59.139:19530")
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
	if err != nil {
    	fmt.Println(err)
	}
	//Transmit and Receive
	for {
		fmt.Fprintf(connection,string(jsonData)+"\n")
		data, _ := bufio.NewReader(connection).ReadString('\n')
		fmt.Println("From -->", data)
		time.Sleep(1*time.Second)
	}
	connection.Close()
}
