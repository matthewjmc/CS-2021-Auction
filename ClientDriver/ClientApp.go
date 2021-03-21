package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"time"

	"encoding/json"
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

func main() {
	data := Package{}
	var userIn int
	uID := rand.Intn(100000-10) + 10
	fmt.Println("\nWhat Would you like to do?\n\t1--> Create an Auction\n\t2--> Join an Auction")
	fmt.Scanf("%d", &userIn)
	switch userIn {
	case 1:
		data.Command = "create"
		fmt.Println("Creating Auction")
	case 2:
		data.Command = "join"
		fmt.Println("Joining Auction")
	}
	data.UserID = uID
	returnVal := handleCon(data)
	if returnVal.Command == "AucCreated" {
		data.AuctionID = returnVal.Data.Value
		data.UserID = returnVal.UserID
		data.Command = "join"
		returnVal = handleCon(data)
	}
	for {

	}
}

func handleCon(data Package) Package {
	received := Package{}
	connection, _ := net.Dial("tcp4", "167.99.67.7:19530")
	defer connection.Close()
	err := connection.(*net.TCPConn).SetKeepAlive(true)
	if err != nil {
		fmt.Println(err)
	}
	err = connection.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		fmt.Println(err)
	}
	//Convert Struct to JSON Document
	var jsonData []byte
	jsonData, err = json.Marshal(data)
	//fmt.Println(jsonData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(connection, string(jsonData)+"\n")
	rawdata, err := bufio.NewReader(connection).ReadString('\n')
	json.Unmarshal([]byte(rawdata), &received)

	return received
}
