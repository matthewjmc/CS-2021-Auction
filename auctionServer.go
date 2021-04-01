package main

import (
	//"bufio"
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
var A = AuctionSystem.AuctionAllocate()
var U = AuctionSystem.UserAllocate()

var hashTable = make(map[uint64]Auction) //Hash Table to Storing Current Auction Data

func main() {
	serverInit()
}

func serverInit() {

	var wg sync.WaitGroup                       //Ensure Data Integrity
	stream, err := net.Listen("tcp4", ":19530") //Listen at port 19530
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stream.Close()
	n := 1
	for {
		con, err := stream.Accept()
		//fmt.Println(con.RemoteAddr())
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
	buffer := make([]byte, 1024)
	for {
		n, err := con.Read(buffer)
		rawdata := string(buffer[:n])
		//rawdata, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal([]byte(rawdata), &received)
		if err != nil {
			fmt.Println(err)
			return
		}
		if received.Command == "create" {
			aucID := _generateAucID()
			state, _ := AuctionSystem.CreateAuctionMain(U, A, received.UserID, aucID, 100, 25, 1*time.Hour, "Demo")
			if state {
				tmp := Package{}
				tmp.Data.Item = "AuctionID"
				tmp.Data.Value = aucID
				tmp.Command = "AucCreated"
				var jsonData []byte
				jsonData, err = json.Marshal(tmp)
				returnData(con, string(jsonData))
				wg.Done()
			} else {
				fmt.Println("Unable to Add Auction")
			}

		} else if !loggedIn && received.Command == "join" {
			state, _ := AuctionSystem.CreateUserMain(U, received.UserID, "Demo")
			defer removeUser(con, received.AuctionID)
			if state {
				addUsr(con, received.AuctionID, received.UserID)
				wg.Done()
				loggedIn = true
				tmp := Package{}
				tmp.Command = "Success"
				var jsonData []byte
				jsonData, err = json.Marshal(tmp)
				returnData(con, string(jsonData))
				go _updateUsers(received.AuctionID, received.UserID)
			} else {
				fmt.Println("Unable to Create User")
			}

		} else if loggedIn {
			switch received.Command {
			case "bid":
				state, err := AuctionSystem.MakeBidMain(U, A, received.UserID, received.AuctionID, received.Data.Value)
				if err != 0 {
					switch err {
					case 1:
						fmt.Println("User not found in System")
					case 2:
						fmt.Println("Auction Not Found in System")
					}
				} else if state {
					go _updateClient(received.AuctionID, received.UserID, received.Data.Value, received.Time)
				}
			}
		}
	}
}

func returnData(con net.Conn, data string) {
	fmt.Fprintf(con, string(data)+"\n")
}

func addUsr(con net.Conn, aID uint64, uID uint64) {
	if _aucExists(aID) {
		temp := hashTable[aID]
		temp.ConnectedClients = append(temp.ConnectedClients, con)
		hashTable[aID] = temp
		//fmt.Println(hashTable[aID])
	} else {
		temp := Auction{
			AuctionID:        aID,
			ConnectedClients: []net.Conn{con},
		}
		hashTable[aID] = temp
		//fmt.Println(hashTable[aID])
	}
	fmt.Println("Number of Rooms Currently", len(hashTable))
}

func _aucExists(aID uint64) bool {
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

func removeUser(con net.Conn, aID uint64) {
	hash := hashTable[aID]
	for i, val := range hash.ConnectedClients {
		if val == con {
			hash.ConnectedClients[i] = hash.ConnectedClients[len(hash.ConnectedClients)-1]
			hash.ConnectedClients[len(hash.ConnectedClients)-1] = nil
			hash.ConnectedClients = hash.ConnectedClients[:len(hash.ConnectedClients)-1]
		}
	}
}
