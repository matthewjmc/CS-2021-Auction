package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	// "github.com/matthewjmc/CS-2021-Auction/AuctionSystem"
)

type Package struct { // Data Sent and Received From user
	AuctionID int
	UserID    int
	Command   string
	Data      struct {
		Item  string
		Value int
	}
}

type Auction struct { //Auctions Running at one time
	AuctionID        int
	ConnectedClients []net.Conn
}

var aucSessions = []Auction{} //All Connected Auction

// mainTimeline(A *AuctionHashTable, U *UserHashTable, command string)

func main() {
	// A := AuctionSystem.AuctionAllocate()
	// U := AuctionSystem.UserAllocate()
	serverInit()
	// temp := AuctionSystem.MainTimeline()

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
		if err != nil {
			fmt.Println(err)
			return
		}
		wg.Add(1)
		go requestHandle(con, &wg)
		fmt.Println(n)
		n++
		wg.Wait()
	}
}

func requestHandle(con net.Conn, wg *sync.WaitGroup) { //Check make Sure other thread does not RW Same Data
	var state bool = false //Check if User has been registered
	var received Package   //Data Received From User to be decoded to Struct
	defer con.Close()
	for {
		rawdata, err := bufio.NewReader(con).ReadString('\n')
		//fmt.Println(rawdata)
		if err != nil {
			fmt.Println(err)
			return
		}
		if !state {
			json.Unmarshal([]byte(rawdata), &received)
			//fmt.Println(received)
			addUsr(con, received.AuctionID, received.UserID)
			wg.Done()
			state = true
		}
		// else if received.Command == "Update" {
		// 	_updateClient(received.AuctionID,)
		// }
	}
	//con.Close()
}

func addUsr(con net.Conn, aID int, uID int) {
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

func _aucExists(aID int) (bool, int) {
	for i := 0; i < len(aucSessions); i++ {
		if aucSessions[i].AuctionID == aID {
			return true, i
		}
	}
	return false, 0
}

func _updateClient(aID int, uID int, price int) {
	var temp Package
	found, index := _aucExists(aID)
	if found {
		temp.UserID = uID
		temp.Data.Item = "Price"
		temp.Data.Value = price
		jsonData, err := json.Marshal(temp)
		if err != nil {
			fmt.Println(err)
		}
		auc := aucSessions[index]
		for i := 0; i < len(auc.ConnectedClients); i++ {
			fmt.Fprintf(auc.ConnectedClients[i], string(jsonData))
		}

	}

}

func _updateServerInfo() {

}
