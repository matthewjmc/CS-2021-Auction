package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
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

/* Global Variables */
var aucSessions = []Auction{} //All Connected Auction
var A = AuctionSystem.AuctionAllocate()
var U = AuctionSystem.UserAllocate()

//var hashTable =  map[uint64]Auction
var hashTable = make(map[uint64]Auction) //Hash Table to Storing Current Auction Data

func main() {
	serverInit()
}

func serverInit() {

	var wg sync.WaitGroup                       //Ensure Data Integrity
	stream, err := net.Listen("tcp4", ":19530") //Listen at port 19530
	if err != nil {
		log.Println(err)
		return
	}
	defer stream.Close()
	n := 1
	for {
		con, err := stream.Accept()
		if err != nil {
			log.Println(err)
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
			wg.Done()

			//Update Data in Cache
			aChan := make(chan AuctionSystem.Auction)
			sChan := make(chan string)
			AuctionSystem.CreateAuctionMain(U, A, aChan, sChan, received.UserID, tmp.Data.Value, 100, 25)

		} else if !loggedIn && received.Command == "join" {
			addUsr(con, received.AuctionID, received.UserID)
			wg.Done()
			// fmt.Printf(" %d has Joined Auction %d\n", received.UserID, received.AuctionID)
			loggedIn = true
			tmp := Package{}
			tmp.Command = "Success"
			var jsonData []byte
			jsonData, err = json.Marshal(tmp)
			returnData(con, string(jsonData))
			go _updateUsers(received.AuctionID, received.UserID)

			AuctionSystem.CreateUserMain(U, received.UserID, "Demo")

		} else if loggedIn {
			switch received.Command {
			case "bid":
				go _updateClient(received.AuctionID, received.UserID, received.Data.Value, received.Time)
				AuctionSystem.MakeBidMain(U, A, received.UserID, received.AuctionID, received.Data.Value)
			}
		}
	}
}

func returnData(con net.Conn, data string) {
	fmt.Fprintf(con, string(data)+"\n")
}

func addUsr(con net.Conn, aID uint64, uID uint64) {
	// exists, index := _aucExists(aID)

	// if len(aucSessions) == 0 && !exists {
	// 	aucSessions = append(aucSessions,
	// 		Auction{
	// 			AuctionID:        aID,
	// 			ConnectedClients: []net.Conn{con}})
	// } else {
	// 	if exists {
	// 		aucSessions[index].ConnectedClients = append(aucSessions[index].ConnectedClients, con)
	// 	} else {
	// 		aucSessions = append(aucSessions,
	// 			Auction{
	// 				AuctionID:        aID,
	// 				ConnectedClients: []net.Conn{con}})
	// 	}
	// }
	if _aucExists(aID) {
		temp := hashTable[aID]
		temp.ConnectedClients = append(temp.ConnectedClients, con)
		hashTable[aID] = temp
	} else {
		temp := Auction{
			AuctionID:        aID,
			ConnectedClients: []net.Conn{con},
		}
		hashTable[aID] = temp
	}
}

func _aucExists(aID uint64) bool {
	// for i := 0; i < len(aucSessions); i++ {
	// 	if aucSessions[i].AuctionID == aID {
	// 		return true, i
	// 	}
	// }
	// return false, 0
	if _, ok := hashTable[aID]; ok {
		return true
	}
	return false
}

func _updateClient(aID uint64, uID uint64, price uint64, sTime []time.Time) {
	var temp Package
	found := _aucExists(aID)
	if found {
		temp.UserID = uID
		temp.Command = "curPrice"
		temp.AuctionID = aID
		temp.Data.Value = price
		temp.Time = append(sTime, time.Now())
		jsonData, err := json.Marshal(temp)
		if err != nil {
			log.Println(err)
		}
		auc := hashTable[aID]
		for i := 0; i < len(auc.ConnectedClients); i++ {
			fmt.Fprintf(auc.ConnectedClients[i], string(jsonData)+"\n")
		}
		// auc := aucSessions[index]
		// for i := 0; i < len(auc.ConnectedClients); i++ {
		// 	fmt.Fprintf(auc.ConnectedClients[i], string(jsonData)+"\n")
		// }
	}
}

func _updateUsers(aID uint64, uID uint64) {
	time.Sleep(1 * time.Second)
	var temp Package
	found := _aucExists(aID)
	if found {
		temp.Command = "usrjoin"
		temp.UserID = uID
		temp.AuctionID = aID
		jsonData, err := json.Marshal(temp)
		if err != nil {
			log.Println(err)
		}
		//auc := aucSessions[index]
		auc := hashTable[aID]
		for i := 0; i < len(auc.ConnectedClients); i++ {
			fmt.Fprintf(auc.ConnectedClients[i], string(jsonData)+"\n")
		}
	}
}

func _generateAucID() uint64 {
	aucID := rand.Uint64()
	exist := _aucExists(aucID)
	for exist {
		aucID := rand.Uint64()
		exist = _aucExists(aucID)
	}
	return aucID
}
