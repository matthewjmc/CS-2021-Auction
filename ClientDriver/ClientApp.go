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
		Value []uint64
	}
}

var serverIP = "load.mcmullin.org:19530"

func main() {
	data := Package{}
	var userIn uint64
	var aID uint64
	var uID uint64
	var initPrice uint64
	var stepSize uint64
	var returnVal Package
	fmt.Println("Please Enter User ID:")
	fmt.Scanf("%d", &uID)
	fmt.Println("\nWhat Would you like to do?\n\t1--> Create an Auction\n\t2--> Join an Auction")
	fmt.Scanf("%d", &userIn)
	switch userIn {
	case 1:
		data.Command = "create"
		data.UserID = uID
		fmt.Println("What would you like starting price to be?")
		fmt.Scanf("%d", &initPrice)
		fmt.Println("Step size for the bidding to be?")
		fmt.Scanf("%d", &stepSize)
		data.Data.Value = []uint64{initPrice, stepSize}
		fmt.Println("Creating Auction!!!")
		returnVal = handleCon(data)
	case 2:
		data.Command = "join"
		data.UserID = uID
		fmt.Println("Please Enter Auction To Join:")
		fmt.Scanf("%d", &aID)
		data.AuctionID = aID

	}
	//fmt.Println(data)
	if returnVal.Command == "AucCreated" {
		data = Package{}
		data.AuctionID = returnVal.Data.Value[0]
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
	connection, _ := net.Dial("tcp4", serverIP)
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
	var stepSize uint64
	connection, _ := net.Dial("tcp4", serverIP)
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

	temp := recData(connection)
	if temp.Command == "Success" {
		stepSize = temp.Data.Value[0]
	}

	go readInput(connection, data.AuctionID, data.UserID)
	for {
		received := Package{}
		rawdata, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		json.Unmarshal([]byte(rawdata), &received)
		if received.Command == "usrjoin" {
			clearScreen()
			fmt.Printf("User %d has Joined Room No.%d\nBid Step Size -->\t%d\nCurrent Price ----->\t%d\n", received.UserID, received.AuctionID, stepSize, received.Data.Value[0])
		} else if received.Command == "curPrice" {
			clearScreen()
			fmt.Printf("Room No.%d\n", received.AuctionID)
			fmt.Printf("Bid Step Size -->\t%d\nCurrent Price ----->\t%d\n", stepSize, received.Data.Value[0])
		} else if received.Command == "invalidBid" {
			clearScreen()
			fmt.Printf("Room No.%d\nBid Step Size -->%d\n", received.AuctionID, stepSize)
			fmt.Println("Please Place a Valid Bid")
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

func recData(connection net.Conn) Package {
	received := Package{}
	rawdata, err := bufio.NewReader(connection).ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal([]byte(rawdata), &received)
	return received
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
		temp.Data.Value = []uint64{price}
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
