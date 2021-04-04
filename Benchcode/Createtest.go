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

var serverIP = "load.mcmullin.org:19530"

func main() {
	arguments := os.Args
	args, _ := strconv.Atoi(arguments[1])
	n := 1
	for n < args+1 {
		go Bench(n)
		n++
		
	}
	for {

	}
}

func Bench(n int) {
	rd := rand.Intn(5)
	time.Sleep(time.Duration(rd)*time.Millisecond)

	data := Package{}
	uID := rand.Uint64()
	data.Command = "create"
	data.UserID = uID
	returnVal := handleCon(data,n)

	if returnVal.Command == "AucCreated" {
		data = Package{}
		data.AuctionID = returnVal.Data.Value
		data.UserID = uID
		data.Command = "join"
		//openCon(data)
		//fmt.Println("Connection Established")

	}
}
	
func handleCon(data Package,n int) Package {
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
	data.Time = append([]time.Time{time.Now()}) //
	jsonData := jsonify(data)
	
	fmt.Fprintf(connection, string(jsonData)+"\n")

	
	rawdata, err := bufio.NewReader(connection).ReadString('\n')
	json.Unmarshal([]byte(rawdata), &received)
	//fmt.Println(rawdata)
	duration := time.Since(received.Time[0]) //
	//fmt.Println("Round Joined Time = ,", duration) //
	fmt.Println(n,",", float64(duration/time.Millisecond))
	return received
}

func jsonify(data Package) []byte {
	var jsonData []byte
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	return jsonData
}
