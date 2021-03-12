package main

import (
	"bufio"
	"fmt"
	"net"
	"encoding/json"
)
type Package struct{  // Data Sent and Received From user
	AuctionID int
	UserID int
	Command  string
	Data struct {
		Item string
		Value int
	}
}

type Auction struct{ //Auctions Running at one time
	AuctionID int
	ConnectedClients []net.Conn
}

var aucSessions = []Auction{} //All Connected Auction

func main() {
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
		go requestHandle(con)
	}
}

func requestHandle(con net.Conn) {
	var state bool = false
	var received Package //Data Received From User to be decoded to Struct

	for {
		rawdata,err := bufio.NewReader(con).ReadString('\n')
		if !state{
			json.Unmarshal([]byte(rawdata), &received)
			addUsr(con,received.AuctionID,received.UserID)
			state = true
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		go returnData(con)
	}
	con.Close()
}

func returnData(con net.Conn){
	fmt.Fprintf(con, con.RemoteAddr().String()+"\n")
}

func addUsr(con net.Conn, aID int, uID int){
	exists,index := _aucExists(aID)
	if len(aucSessions) == 0 && !exists{
		fmt.Println(aID)
		aucSessions = append(aucSessions,
			Auction{
				AuctionID:aID,
				ConnectedClients:[]net.Conn{con}})
	} else{
		if exists{
			aucSessions[index].ConnectedClients = append(aucSessions[index].ConnectedClients,con)
		}else{
			aucSessions = append(aucSessions,
				Auction{
					AuctionID:aID,
					ConnectedClients:[]net.Conn{con}})
		}
		//fmt.Println("Line85:",aID)
	}
	
	//fmt.Println(aucSessions)
	
}

func _aucExists(aID int) (bool,int){
	for i :=0 ; i < len(aucSessions); i++{
		if aucSessions[i].AuctionID == aID{
			return true,i
		}
	}
	return false,0
}
