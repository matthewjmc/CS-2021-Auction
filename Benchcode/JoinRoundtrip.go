package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
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
	Time []time.Time
}

var serverIP = "com1.mcmullin.org:19530"

func main() {
	arguments := os.Args
	args, _ := strconv.Atoi(arguments[1])
	n := 1
	for n < args+1 {
		go Bench()
		n++
	}
	for {

	}
}

func Bench() {
	data := Package{}
	// aID := rand.Uint64()
	uID := rand.Uint64()
	data.Command = "create"
	data.UserID = uID
	returnVal := handleCon(data)
	//fmt.Println(returnVal.Data.Value)
	if returnVal.Command == "AucCreated" {
		data = Package{}
		data.AuctionID = returnVal.Data.Value
		data.UserID = uID
		data.Command = "join"
		//openCon(data)
		//fmt.Println("Connection Established")

	}
	// else if data.Command == "join" {
	// 	data = Package{}
	// 	data.AuctionID = aID
	// 	data.UserID = uID
	// 	data.Command = "join"
	// 	openCon(data)
	// }

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
	//start := []time.Time{time.Now()} //
	//data.Time = start
	jsonData := jsonify(data)
	fmt.Fprintf(connection, string(jsonData)+"\n")

	
	rawdata, err := bufio.NewReader(connection).ReadString('\n')
	json.Unmarshal([]byte(rawdata), &received)
	fmt.Println(rawdata)
	//duration := time.Since(start) //
	//fmt.Println("Round Joined Time = ,", duration) //
	return received
}

func openCon(data Package) {

	connection, _ := net.Dial("tcp4", serverIP)
	timeJoin := time.Now()            //JoinTest
	data.Time = []time.Time{timeJoin} //JoinTest
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

	for {
		received := Package{}
		//price := rand.Uint64()
		rawdata, err := bufio.NewReader(connection).ReadString('\n')
		//fmt.Println(rawdata)

		if err != nil {
			fmt.Println(err)
			return
		}
		json.Unmarshal([]byte(rawdata), &received) //JoinTest
		// fmt.Println(received.Time[1])
		//fmt.Println(received.Time)
		fmt.Println("JoinRoundTime", time.Now().Sub(received.Time[0])) //JoinRoundTest

		// if len(received.Time) > 0 {
		// 	fmt.Println("Response Time:,", float64(time.Since(received.Time[0])/time.Millisecond))
		// } //Round Trip
		//time.Sleep(1 * time.Second)
		// temp := Package{}
		// temp.UserID = data.UserID
		// temp.AuctionID = data.AuctionID
		// temp.Command = "bid"
		// temp.Data.Item = "price"
		// temp.Time = []time.Time{time.Now()}
		// temp.Data.Value = price
		// jsondata := jsonify(temp)
		// fmt.Fprintf(connection, string(jsondata)+"\n")
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
