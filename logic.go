package main

import (
	. "CS-2021-Auction/AuctionSystem"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// AuctionLogic is a file composed of Tagun's and Katisak's code altogether.
// The file is also being separated into 3 files which are auction_method, auction_timeline, data_structure.
/*
func main() {
	//multipleUserTest()
	//testTimeFormat()
	A := AuctionAllocate()
	U := UserAllocate()
	// modification of memory allocation to be dynamically allocating.
	// mainTimeline(A, U, )
}
*/

type Data struct {
	command      string
	uid          uint64
	fullname     string
	aid          uint64
	itemname     string
	biddingValue uint64
	biddingStep  uint64
	duration     time.Duration
}

func mainTimeline(A *AuctionHashTable, U *UserHashTable, instructions Data) {

	command := instructions.command

	if command == "User" || command == "user" {
		user_report := make(chan User)
		report_log := make(chan string)
		go CreateUserMain(U, instructions.uid, instructions.fullname)
		newUser := <-user_report
		log := <-report_log
		fmt.Println(log, newUser)

	} else if command == "Auction" || command == "auction" {
		report_auction := make(chan Auction)
		report_log := make(chan string)
		go CreateAuctionMain(U, A, instructions.uid, instructions.aid, instructions.biddingValue, instructions.biddingStep, instructions.duration, instructions.itemname)
		newAuction := <-report_auction
		log := <-report_log
		fmt.Println(newAuction, log)
	} else if command == "bid" {

		if !A.SearchAuctIDHashTable(instructions.aid) {
			fmt.Println("The auction has not been found within the memory")
		} else {
			report_price := make(chan uint64)
			report_log := make(chan string)
			go MakeBidMain(U, A, instructions.uid, instructions.aid, instructions.biddingValue)
			finalAuction := <-report_price
			log := <-report_log
			fmt.Println(finalAuction, log)
		}
	} else if command == "searchAuction" {
		if A.SearchAuctIDHashTable(instructions.aid) {
			fmt.Println("Auction", instructions.aid, " is found within the system")
		} else {
			fmt.Println("Auction", instructions.aid, " is not found within the system")
		}
	} else if command == "deleteAuction" {
		if A.AuctionHashAccessDelete(instructions.aid) {
			fmt.Println("Auction", instructions.aid, " has been deleted for the system")
		} else {
			fmt.Println("Auction", instructions.aid, " is not found within the system")
		}
	} else if command == "searchUser" {
		if U.SearchUserIDHashTable(instructions.uid) {
			fmt.Println("Auction", instructions.aid, " is found within the system")
		} else {
			fmt.Println("Auction", instructions.aid, " is not found within the system")
		}
	} else if command == "deleteUser" {
		if U.UserHashAccessDelete(instructions.uid) {
			fmt.Println("Auction", instructions.aid, " has been deleted for the system")
		} else {
			fmt.Println("Auction", instructions.aid, " is not found within the system")
		}
	}

}

var wg sync.WaitGroup

func MakeBidMain(u *UserHashTable, h *AuctionHashTable, uid uint64, targetid uint64, placeVal uint64) (bool, uint64) {
	if !u.SearchUserIDHashTable(uid) {
		return false, 1 // code 1 , the user has not been found within the system.
	} else if !h.SearchAuctIDHashTable(targetid) {
		return false, 2 // code 2 , the auction has not been found within the system.
	} else {
		bidTime := time.Now().Format(time.RFC3339Nano)
		currUser := *u.AccessUserHash(uid)
		newBid := CreateBid(currUser, placeVal, bidTime)
		target := h.AccessHashAuction(targetid)
		target.UpdateAuctionWinner(newBid)
		h.AuctionHashAccessUpdate(*target)
		return true, 0 // code 0 . the bid has been made and updated properly.
	}
}

func CreateUserMain(h *UserHashTable, uid uint64, name string) (bool, uint64) {
	if !h.SearchUserIDHashTable(uid) {
		newUser := CreateUser("username"+fmt.Sprint(uid), name, uid)
		h.InsertUserToHash(newUser)
		return true, 0 // code 0 , the user has not been found within the system, creating new user object.
	} else {
		return false, 1 // code 1 , the user has been found in the system.
	}
}

func CreateAuctionMain(U *UserHashTable, A *AuctionHashTable, uid uint64, aid uint64, initial uint64, step uint64, duration time.Duration, itemName string) (bool, uint64) {
	user := U.AccessUserHash(uid)
	if !A.SearchAuctIDHashTable(aid) {
		newAuction := CreateAuction(*user, initial, step, aid, duration, itemName)
		A.InsertAuctToHash(newAuction.CreatedAuction)
		return true, 0 // code 0, auction has not been found within the system, creating new auction object.
	} else {
		return false, 1 // error code 1 , auction has been found in the system.
	}
}

func _insertAuctionToDB(auction Auction) bool {

	db, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/database_name")

	if err != nil {
		panic(err.Error())
	}
	// error handler whether what causes the error regarding the connection to the database.
	defer db.Close()
	// perform a db.Query CRUD commands inputted.
	insert, err := db.Query("INSERT INTO auction_table VALUES ( 2, 'T' )")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	// error handler whether what causes the error regarding the connection to the database.
	return true
}

func _insertUserToDB(user User) bool {

	db, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/database_name")

	if err != nil {
		panic(err.Error())
	}
	// error handler whether what causes the error regarding the connection to the database.
	defer db.Close()
	// perform a db.Query CRUD commands inputted.
	insert, err := db.Query("INSERT INTO user_table VALUES ( 2, 'T' )")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	// error handler whether what causes the error regarding the connection to the database.
	return true
}

func _insertBidToDB(bid Bid) bool {

	db, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/database_name")

	if err != nil {
		panic(err.Error())
	}
	// error handler whether what causes the error regarding the connection to the database.
	defer db.Close()
	// perform a db.Query CRUD commands inputted.
	insert, err := db.Query("INSERT INTO bidding_table VALUES ( 2, 'T' )")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	// error handler whether what causes the error regarding the connection to the database.
	return true
}

func main() {

	db, err := sql.Open("mysql", "auction:Helloworld1@tcp(db.mcmullin.org)/")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Database created successfully")
	}
	_, err = db.Exec("CREATE DATABASE auction_system")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully created database..")
	}
	_, err = db.Exec("USE auction_system")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("DB selected successfully..")
	}

	stmt, err := db.Prepare("CREATE Table user_table(AccountID int UNSIGNED NOT NULL UNIQUE PRIMARY KEY, first_name varchar(20) NOT NULL, last_name varchar(20) NOT NULL) );")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("User Table created successfully..")
	}

	stmt2, err := db.Prepare("CREATE Table bid_table(biddingID int UNSIGNED NOT NULL UNIQUE PRIMARY KEY, bidderID int UNSIGNED NOT NULL UNIQUE, bidderUsername varchar(30) NOT NULL, bidPrice int UNSIGNED NOT NULL, bidTime varchar(30) NOT NULL) FOREIGN KEY (bidderID) reference user_table(AccountID) );")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt2.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Bid Table created successfully..")
	}

	stmt3, err := db.Prepare("CREATE Table auction_table(AuctionID int UNSIGNED NOT NULL UNIQUE PRIMARY KEY, AuctioneerID int UNSIGNED NOT NULL, ItemName varchar(30) NOT NULL, CurrWinnerID int UNSIGNED NOT NULL,CurrWinnerName varchar(30),CurrMaxBid int UNSIGNED NOT NULL,BidStep int UNSIGNED NOT NULL,LatestBidTime varchar(30) NOT NULL, StartTime varchar(30) NOT NULL, EndTime varchar(30) NOT NULL, FOREIGN KEY (AuctioneerID) reference user_table(AccountID) FOREIGN KEY (CurrWinnerID) reference user_table(AccountID) );")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt3.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Table created successfully..")
	}

	defer db.Close()
}

// Server : db.mcmullin.org:3306
// username : auction
// password : Helloworld1
