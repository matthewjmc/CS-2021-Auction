package main

import (
	//"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
//Ticker for AutoUpdate Database
var ticker = time.NewTicker(1 * time.Minute)

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
	db, err := sql.Open("mysql", "auction:Helloworld1@tcp(db.mcmullin.org:3306)/auction_system")
	if err != nil {
		fmt.Println("DB Error:", err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(200000)
	db.SetMaxIdleConns(200000)
	//AutoUpdate
	done := make(chan bool)
	go dbUpdate(db, done)
	defer ticker.Stop()
	for {
		con, err := stream.Accept()
		//fmt.Println(con.RemoteAddr())
		if err != nil {
			log.Println(err)
			return
		}
		wg.Add(1)
		go requestHandle(con, &wg)
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
		if err != nil {
			fmt.Println(err)
		}

		json.Unmarshal([]byte(rawdata), &received)

		if err != nil {
			fmt.Println(err)
			return
		}

		if received.Command == "create" {
			aucID := received.Data.Value
			// aucID := _generateAucID()
			state, _ := AuctionSystem.CreateAuctionMain(U, A, received.UserID, aucID, 100, 25, 1*time.Hour, "Demo")
			// fmt.Println(received.Time)
			// fmt.Println(received)
			if state {
				tmp := Package{}
				tmp.Data.Item = "AuctionID"
				tmp.Data.Value = aucID
				tmp.Command = "AucCreated"

				tmp.Time = received.Time

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
				//tmp.Time = append(received.Time, time.Now()) //test join time
				//tmp.Time = append(received.Time) //test join round time

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
	} else {
		temp := Auction{
			AuctionID:        aID,
			ConnectedClients: []net.Conn{con},
		}
		hashTable[aID] = temp
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

func dbUpdate(db *sql.DB, done chan bool) {
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			fmt.Println("Pushing to DB")
			for aucID, _ := range hashTable {
				data := AuctionSystem.AccessHashAuctionCalling(A, aucID)
				if data != nil {
					query := `INSERT INTO auction(aucID,userID,itemName,currWinnerID,currWinnerName,currMaxBid,bidStep,latestBidTime,startTime,endTime) 
							values(?,?,?,?,?,?,?,?,?,?)
							ON DUPLICATE KEY UPDATE 
							aucID=?,
							userID=?,
							itemName=?,
							currWinnerID=?,
							currWinnerName=?,
							currMaxBid=?,
							bidStep=?,
							latestBidTime=?,
							startTime=?,
							endTime=?`
					stmt, err := db.Prepare(query)
					if err != nil {
						fmt.Println("Prepare Error:", err)
					}
					_, err = stmt.Exec(
						data.AuctionID,
						data.AuctioneerID,
						data.ItemName,
						data.CurrWinnerID,
						data.CurrWinnerName,
						data.CurrMaxBid,
						data.BidStep,
						data.LatestBidTime,
						data.StartTime,
						data.EndTime,
						data.AuctionID,
						data.AuctioneerID,
						data.ItemName,
						data.CurrWinnerID,
						data.CurrWinnerName,
						data.CurrMaxBid,
						data.BidStep,
						data.LatestBidTime,
						data.StartTime,
						data.EndTime)
					if err != nil {
						fmt.Println("Query Error:", err)
					}
				}
			}
			fmt.Println("Done pushing DB")
		}
	}
}
