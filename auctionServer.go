package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/matthewjmc/CS-2021-Auction/AuctionSystem"
)

type Package struct { // Data Sent and Received From user
	AuctionID uint64
	UserID    uint64
	Command   string
	Data      struct {
		Item  string
		Value uint64
	}
	Time []time.Time
}

type Temp struct {
	AuctionID uint64
	UserID    uint64
}

type Auction struct { //Auctions Running at one time
	AuctionID        uint64
	ConnectedClients []net.Conn
}

var aucSessions = []Auction{} //All Connected Auction
var A = AuctionSystem.AuctionAllocate()
var U = AuctionSystem.UserAllocate()

func main() {
	serverInit()
}

func serverInit() {

	var wg sync.WaitGroup                       //Ensure Data Integrity
	stream, err := net.Listen("tcp4", ":19530") //Listen at port 19530
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer stream.Close()
	n := 1
	for {
		con, err := stream.Accept()
		if err != nil {
			//fmt.Println(err)
			return
		}
		wg.Add(1)
		go requestHandle(con, &wg)
		n++
		wg.Wait()
	}
}

func requestHandle(con net.Conn, wg *sync.WaitGroup) { //Check make Sure other thread does not RW Same Data
	var loggedIn bool = false //Check if User has been registered
	var received Package      //Data Received From User to be decoded to Struct
	defer con.Close()
	for {
		rawdata, err := bufio.NewReader(con).ReadString('\n')
		json.Unmarshal([]byte(rawdata), &received)
		if err != nil {
			////fmt.Println(err)
			return
		}
		if received.Command == "create" {
			tmp := Package{}
			tmp.Data.Item = "AuctionID"
			tmp.Data.Value = _generateAucID()
			tmp.Command = "AucCreated"
			var jsonData []byte
			jsonData, err = json.Marshal(tmp)
			returnData(con, string(jsonData))
			//fmt.Printf("User %d Has Created Auction %d\n", received.UserID, tmp.Data.Value)
			wg.Done()

			//Update Data in Cache
			aChan := make(chan AuctionSystem.Auction)
			sChan := make(chan string)
			AuctionSystem.CreateAuctionMain(U, A, aChan, sChan, received.UserID, tmp.Data.Value, 100, 25) //Last 2 params  InitBid, StepSize

		} else if !loggedIn && received.Command == "join" {
			addUsr(con, received.AuctionID, received.UserID)
			wg.Done()
			fmt.Printf("User %d has Joined Auction %d\n", received.UserID, received.AuctionID)
			loggedIn = true

			tmp := Package{}
			tmp.Command = "Success"
			var jsonData []byte
			jsonData, err = json.Marshal(tmp)
			returnData(con, string(jsonData))
			go _updateUsers(received.AuctionID, received.UserID)

			// uChan := make(chan uint64)
			// sChan := make(chan string)
			AuctionSystem.CreateUserMain(U, received.UserID, "Demo")
			//AuctionSystem.CreateUserMain(U, uChan, sChan, received.UserID, "Demo")

		} else if loggedIn {
			switch received.Command {
			case "bid":
				//fmt.Println("User Requesting to Bid")
				go _updateClient(received.AuctionID, received.UserID, received.Data.Value, received.Time)
				// pChan := make(chan uint64)
				// sChan := make(chan string)
				// AuctionSystem.MakeBidMain(U, A, pChan, sChan, received.UserID, received.AuctionID, received.Data.Value)
				AuctionSystem.MakeBidMain(U, A, received.UserID, received.AuctionID, received.Data.Value)
			}
		}
	}
}

func returnData(con net.Conn, data string) {
	fmt.Fprintf(con, string(data)+"\n") //Fix this
}

func addUsr(con net.Conn, aID uint64, uID uint64) {
	exists, index := _aucExists(aID)

	if len(aucSessions) == 0 && !exists {
		aucSessions = append(aucSessions,
			Auction{
				AuctionID:        aID,
				ConnectedClients: []net.Conn{con}})
	} else {
		if exists {
			aucSessions[index].ConnectedClients = append(aucSessions[index].ConnectedClients, con)
		} else {
			aucSessions = append(aucSessions,
				Auction{
					AuctionID:        aID,
					ConnectedClients: []net.Conn{con}})
		}
	}
}

func _aucExists(aID uint64) (bool, int) {
	for i := 0; i < len(aucSessions); i++ {
		if aucSessions[i].AuctionID == aID {
			return true, i
		}
	}
	return false, 0
}

func _updateClient(aID uint64, uID uint64, price uint64, sTime []time.Time) {
	var temp Package
	found, index := _aucExists(aID)
	if found {
		temp.UserID = uID
		temp.Command = "curPrice"
		temp.AuctionID = aID
		temp.Data.Value = price
		temp.Time = append(sTime, time.Now())
		jsonData, err := json.Marshal(temp)
		if err != nil {
			//fmt.Println(err)
		}
		auc := aucSessions[index]
		for i := 0; i < len(auc.ConnectedClients); i++ {
			fmt.Fprintf(auc.ConnectedClients[i], string(jsonData)+"\n")
		}
	}
}

func _updateUsers(aID uint64, uID uint64) {
	time.Sleep(1 * time.Second)
	var temp Package
	found, index := _aucExists(aID)
	if found {
		temp.Command = "usrjoin"
		temp.UserID = uID
		temp.AuctionID = aID
		jsonData, err := json.Marshal(temp)
		if err != nil {
			//fmt.Println(err)
		}
		auc := aucSessions[index]
		for i := 0; i < len(auc.ConnectedClients); i++ {
			fmt.Fprintf(auc.ConnectedClients[i], string(jsonData)+"\n")
		}
	}
}

func _generateAucID() uint64 {
	aucID := rand.Uint64()
	exist, _ := _aucExists(aucID)
	for exist {
		aucID := rand.Uint64()
		exist, _ = _aucExists(aucID)
	}
	return aucID
}
