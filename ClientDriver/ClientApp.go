package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	"encoding/json"
)

type Package struct {
	AuctionID uint64
	UserID    uint64
	Command   string
	Data      struct {
		Item  string
		Value uint64
	}
}

func main() {
	data := Package{}
	var userIn uint64
	var aID uint64
	var uID uint64
	fmt.Println("Please Enter User ID:")
	fmt.Scanf("%d", &uID)
	fmt.Println("\nWhat Would you like to do?\n\t1--> Create an Auction\n\t2--> Join an Auction")
	fmt.Scanf("%d", &userIn)
	switch userIn {
	case 1:
		data.Command = "create"
		fmt.Println("Creating Auction!!!")
	case 2:
		data.Command = "join"
		fmt.Println("Please Enter Auction To Join:")
		fmt.Scanf("%d", &aID)

	}
	returnVal := handleCon(data)
	if returnVal.Command == "AucCreated" {
		data = Package{}
		data.AuctionID = returnVal.Data.Value
		data.UserID = uID
		data.Command = "join"
		//fmt.Println(data)
		openCon(data)

	} else if data.Command == "join" {
		data = Package{}
		data.AuctionID = aID
		data.UserID = uID
		data.Command = "join"
		//fmt.Println(data)
		openCon(data)
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
	jsonData := jsonify(data)
	fmt.Fprintf(connection, string(jsonData)+"\n")
	rawdata, err := bufio.NewReader(connection).ReadString('\n')
	json.Unmarshal([]byte(rawdata), &received)

	return received
}

func openCon(data Package) {
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
	fmt.Fprintf(connection, string(jsonify(data))+"\n")
	go readInput(connection, data.AuctionID, data.UserID)
	for {
		received := Package{}
		//fmt.Println("Waiting For Data")
		rawdata, err := bufio.NewReader(connection).ReadString('\n')
		//fmt.Println(rawdata)
		if err != nil {
			fmt.Println(err)
			return
		}
		json.Unmarshal([]byte(rawdata), &received)
		if received.Command == "usrjoin" {
			clearScreen()
			fmt.Printf("User %d has Joined Room No.%d\n", received.UserID, received.AuctionID)
		} else if received.Command == "curPrice" {
			clearScreen()
			fmt.Printf("Room No.%d\n", received.AuctionID)
			fmt.Printf("Current Price -----> %d\n", received.Data.Value)
		}
	}

}

func jsonify(data Package) []byte {
	var jsonData []byte
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	return jsonData
}

func readInput(con net.Conn, aID uint64, uID uint64) {
	temp := Package{}
	temp.UserID = uID
	temp.AuctionID = aID
	temp.Command = "bid"
	temp.Data.Item = "price"
	for {
		var price uint64
		fmt.Scanf("%d", &price)
		temp.Data.Value = price
		jsondata := jsonify(temp)
		clearScreen()
		fmt.Fprintf(con, string(jsondata)+"\n")
		//fmt.Println("Data Sent")
	}
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
