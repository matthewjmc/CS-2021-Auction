package AuctionSystem

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func DatabaseInit() {
	db, debug := sql.Open("mysql", "username:password@tcp(db_server)/") // for initialization
	if debug != nil {
		fmt.Println(debug.Error())
	}
	_, debug = db.Exec("CREATE DATABASE auction_system")
	if debug != nil {
		fmt.Println(debug.Error())
	}
	_, debug = db.Exec("USE auction_system")
	if debug != nil {
		fmt.Println(debug.Error())
	}
	statement, debug := db.Prepare("CREATE Table user_table( AccountID int UNSIGNED NOT NULL UNIQUE PRIMARY KEY, Username varchar(20) NOT NULL, Fullname varchar(20) NOT NULL )")
	if debug != nil {
		fmt.Println(debug.Error())
	}
	_, debug = statement.Exec()
	if debug != nil {
		fmt.Println(debug.Error())
	}
	statement2, debug := db.Prepare("CREATE Table auction_table( AuctionID int UNSIGNED NOT NULL UNIQUE PRIMARY KEY,AuctioneerID int UNSIGNED NOT NULL,ItemName varchar(30) NOT NULL, CurrWinnerID int UNSIGNED NOT NULL, CurrWinnerName varchar(30), CurrMaxBid int UNSIGNED NOT NULL, BidStep int UNSIGNED NOT NULL, LatestBidTime varchar(50) NOT NULL, StartTime varchar(50) NOT NULL, EndTime varchar(50) NOT NULL, FOREIGN KEY (AuctioneerID) references user_table(AccountID), FOREIGN KEY (CurrWinnerID) references user_table(AccountID))")
	if debug != nil {
		fmt.Println(debug.Error())
	}
	_, debug = statement2.Exec()
	if debug != nil {
		fmt.Println(debug.Error())
	}
	statement3, debug := db.Prepare("CREATE Table bid_table( BiddingID int UNSIGNED NOT NULL UNIQUE PRIMARY KEY, BidderID int UNSIGNED NOT NULL, BidderUsername varchar(30) NOT NULL, BidPrice int UNSIGNED NOT NULL, BidTime varchar(50) NOT NULL , AuctionID int UNSIGNED NOT NULL, FOREIGN KEY (BidderID) REFERENCES user_table(AccountID), FOREIGN KEY (AuctionID) references auction_table(AuctionID) );")
	if debug != nil {
		fmt.Println(debug.Error())
	}
	_, debug = statement3.Exec()
	if debug != nil {
		fmt.Println(debug.Error())
	}

	defer db.Close()
}

func InsertAuctionToDB(auction Auction) bool {
	db, err := sql.Open("mysql", "username:password@tcp(db_server)/db_name")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	query := "INSERT INTO auction_table(AuctionID,AuctioneerID,ItemName,CurrWinnerID,CurrWinnerName,CurrMaxBid,BidStep,LatestBidTime,StartTime,EndTime) VALUES ("
	auctionId := strconv.FormatUint(auction.AuctionID, 10)
	auctioneerId := "," + strconv.FormatUint(auction.AuctioneerID, 10)
	itemName := "," + "\"" + auction.ItemName + "\""
	currWinnerID := "," + strconv.FormatUint(auction.CurrWinnerID, 10)
	currWinnerName := "," + "\"" + auction.CurrWinnerName + "\""
	currMaxBid := "," + strconv.FormatUint(auction.CurrMaxBid, 10)
	bidStep := "," + strconv.FormatUint(auction.BidStep, 10)
	latestBidTime := "," + "\"" + auction.LatestBidTime + "\""
	startTime := "," + "\"" + auction.StartTime + "\""
	EndTime := "," + "\"" + auction.EndTime + "\""
	exec := query + auctionId + auctioneerId + itemName + currWinnerID + currWinnerName + currMaxBid + bidStep + latestBidTime + startTime + EndTime + ")"
	insert, err := db.Query(exec)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	return true
}

func UpdateAuctionInDB(auction Auction) bool {

	db, err := sql.Open("mysql", "username:password@tcp(db_server)/db_name")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	update, err := db.Prepare("UPDATE auction_table SET CurrMaxBid = ? , CurrWinnerID = ? , CurrWinnerName = ? , LatestBidTime = ? WHERE AuctionID = ?")
	if err != nil {
		panic(err.Error())
	}

	update.Exec(auction.CurrMaxBid, auction.CurrWinnerID, auction.CurrWinnerName, auction.LatestBidTime, auction.AuctionID)
	defer update.Close()
	return true
}

func InsertUserToDB(user User) bool {

	db, err := sql.Open("mysql", "username:password@tcp(db_server)/db_name")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	query := "INSERT INTO user_table(AccountID,Username,Fullname) VALUES ("
	accountId := strconv.FormatUint(user.AccountID, 10)
	username := "," + "\"" + user.Username + "\""
	fullname := "," + "\"" + user.Fullname + "\""

	exec := query + accountId + username + fullname + ")"
	insert, err := db.Query(exec)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	return true
}

func InsertBidToDB(bid Bid, target Auction) bool {

	db, err := sql.Open("mysql", "username:password@tcp(db_server)/db_name")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	query := "INSERT INTO bid_table(BiddingID,BidderID,BidderUsername,BidPrice,BidTime,AuctionID) VALUES ("

	bidId := strconv.FormatUint(bid.BiddingID, 10)
	bidderId := "," + strconv.FormatUint(bid.BidderID, 10)
	bidderName := "," + "\"" + bid.BidderUsername + "\""
	bidPrice := "," + strconv.FormatUint(bid.BidPrice, 10)
	bidTime := "," + "\"" + bid.BidTime + "\""
	AuctionId := "," + strconv.FormatUint(target.AuctionID, 10)
	exec := query + bidId + bidderId + bidderName + bidPrice + bidTime + AuctionId + ")"
	insert, err := db.Query(exec)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	return true
}