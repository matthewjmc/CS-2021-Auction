package main

import (
	"fmt"
	"net"
	"sync"
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

var response []time.Duration
var serverIP = "143.198.212.115:8888"

func main() {
	arguments := os.Args
	args, _ := strconv.Atoi(arguments[1])
	var wg sync.WaitGroup
	n := 0
	for n < args {
		user := Package{}
		user.AuctionID = rand.Intn(10-1) + 1
		user.UserID = rand.Intn(2000-1) + 1
		user.Data.Item = "bid"
		user.Data.Value = rand.Intn(100000-10) + 10
		wg.Add(1)
		go handleCon(user, n, &wg)
		n++
	}
	wg.Wait()
	fmt.Println("Main: Completed")

}

func handleCon(data Package, n int, wg *sync.WaitGroup) {
	connection, _ := net.Dial("tcp", serverIP)
	defer connection.Close()
	defer wg.Done()
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
	start := 1
	for start < n {
		start := time.Now()
		fmt.Fprintf(connection, string(jsonData)+"\n")
		//data, _ := bufio.NewReader(connection).ReadString('\n')
		//fmt.Println("From -->", data)
		elapsed := time.Since(start)
		fmt.Println(elapsed)
		//response = append(response, elapsed)
		time.Sleep(1 * time.Second)
		n++
	}

}

func sum(array []float64) float64 {
	var result float64 = 0
	for _, v := range array {
		result += v
	}
	return result
}
