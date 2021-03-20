package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
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

func main() {
	serverInit()
}

func serverInit() {
	var wg sync.WaitGroup                      //Ensure Data Integrity
	stream, err := net.Listen("tcp", ":19530") //Listen at port 19530
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stream.Close()
	for {
		con, err := stream.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		wg.Add(1)
		go requestHandle(con, &wg)
		wg.Wait()
	}
}

func requestHandle(con net.Conn, wg *sync.WaitGroup) { //Check make Sure other thread does not RW Same Data
	var state bool = false //Check if User has been registered
	var received Package   //Data Received From User to be decoded to Struct
	defer con.Close()
	for {
		rawdata, err := bufio.NewReader(con).ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println(err)
			break
		}
		if !state {
			json.Unmarshal([]byte(rawdata), &received)
			addUsr(con, received.AuctionID, received.UserID)
			wg.Done()
			state = true
		}
		// else if received.Command == "Update" {
		// 	_updateClient(received.AuctionID,)
		// }

	}
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
