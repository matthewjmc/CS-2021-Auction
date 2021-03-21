package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"encoding/json"
	"math/rand"
	"os"
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

var timeouts = 0
var serverIP = "10.104.0.9:19530"

func main() {
	LOG_FILE := "logs/bench.txt"
	// open log file
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	arguments := os.Args
	args, _ := strconv.Atoi(arguments[1])
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
	connection, _ := net.Dial("tcp", serverIP)
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
	if err != nil {
		fmt.Println(err)
	}
	//Transmit and Receive
	for {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			timeouts++
		} else if err != nil {
			log.Println(err)
		}

		fmt.Fprintf(connection, string(jsonData)+"\n")
		//data, _ := bufio.NewReader(connection).ReadString('\n')
		//fmt.Println("From -->", data)
		time.Sleep(1 * time.Second)
	}

}
