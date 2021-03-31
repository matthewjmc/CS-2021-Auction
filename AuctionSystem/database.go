package AuctionSystem

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func DatabaseInit() {

	db, err := sql.Open("mysql", "auction:Helloworld1@tcp(db.mcmullin.org)/")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = db.Exec("CREATE DATABASE auction_system")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = db.Exec("USE auction_system")
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt, err := db.Prepare("CREATE Table user_table( AccountID int UNSIGNED NOT NULL UNIQUE PRIMARY KEY, first_name varchar(20) NOT NULL, last_name varchar(20) NOT NULL )")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt2, err := db.Prepare("CREATE Table bid_table( biddingID int UNSIGNED NOT NULL UNIQUE PRIMARY KEY, bidderID int UNSIGNED NOT NULL UNIQUE, bidderUsername varchar(30) NOT NULL, bidPrice int UNSIGNED NOT NULL, bidTime varchar(50) NOT NULL,FOREIGN KEY (bidderID) REFERENCES user_table(AccountID) );")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt2.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt3, err := db.Prepare("CREATE Table auction_table( AuctionID int UNSIGNED NOT NULL UNIQUE PRIMARY KEY,AuctioneerID int UNSIGNED NOT NULL,ItemName varchar(30) NOT NULL, CurrWinnerID int UNSIGNED NOT NULL, CurrWinnerName varchar(30), CurrMaxBid int UNSIGNED NOT NULL, BidStep int UNSIGNED NOT NULL, LatestBidTime varchar(50) NOT NULL, StartTime varchar(50) NOT NULL, EndTime varchar(50) NOT NULL, FOREIGN KEY (AuctioneerID) references user_table(AccountID), FOREIGN KEY (CurrWinnerID) references user_table(AccountID))")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt3.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
}

// Server : db.mcmullin.org:3306
// username : auction
// password : " first result usually used in programming world as an intro to everylanguage without spacing ,1"

func InsertAuctionToDB(auction Auction) bool {

	db, err := sql.Open("mysql", "auction:Helloworld1@tcp(db.mcmullin.org)/auction_system")
	if err != nil {
		panic(err.Error())
		return false
	}
	defer db.Close()
	query := "INSERT INTO auction_table(AuctionID,AuctioneerID,ItemName,CurrWinnerID,CurrWinnerName,CurrMaxBid,BidStep,LatestBidTime,StartTime,EndTime) VALUES ("
	auctionId := strconv.FormatUint(auction.AuctionID, 10)
	auctioneerId := "," + strconv.FormatUint(auction.AuctioneerID, 10)
	itemName := "," + auction.ItemName
	currWinnerID := "," + strconv.FormatUint(auction.CurrWinnerID, 10)
	currWinnerName := "," + auction.CurrWinnerName
	currMaxBid := "," + strconv.FormatUint(auction.CurrMaxBid, 10)
	bidStep := "," + strconv.FormatUint(auction.BidStep, 10)
	latestBidTime := "," + fmt.Sprint(auction.LatestBidTime)
	startTime := "," + auction.StartTime
	EndTime := "," + auction.EndTime

	fmt.Println(latestBidTime)

	exec := query + auctionId + auctioneerId + itemName + currWinnerID + currWinnerName + currMaxBid + bidStep + latestBidTime + startTime + EndTime + ")"
	fmt.Println(exec)
	insert, err := db.Query(exec)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	return true
}

func InsertUserToDB(user User) bool {

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

func InsertBidToDB(bid Bid) bool {

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
