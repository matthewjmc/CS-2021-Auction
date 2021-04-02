package AuctionSystem

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func DatabaseInit() {

	db, debug := sql.Open("mysql", "auction:Helloworld1@tcp(db.mcmullin.org)/")
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
	db, err := sql.Open("mysql", "auction:Helloworld1@tcp(db.mcmullin.org)/auction_system")
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

	db, err := sql.Open("mysql", "auction:Helloworld1@tcp(db.mcmullin.org)/auction_system")
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

	db, err := sql.Open("mysql", "auction:Helloworld1@tcp(db.mcmullin.org)/auction_system")
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

	db, err := sql.Open("mysql", "auction:Helloworld1@tcp(db.mcmullin.org)/auction_system")
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func UserFromDBtoHash(u *UserHashTable) (bool, uint64) {

	db, err := sql.Open("mysql", "auction:Helloworld1@tcp(db.mcmullin.org)/auction_system")
	if err != nil {
		return false, 1 // code 1 : cannot connect to the database.
		panic(err.Error())
	}
	defer db.Close()
	user_result, _ := db.Query("SELECT * FROM user_table")

	for user_result.Next() {
		var id uint64
		var username string
		var fullname string
		err = user_result.Scan(&id, &username, &fullname)
		if err != nil {
			return false, 3 // code 3 : failed to read information from the query.
			panic(err)
		}
		user := CreateUser(username, fullname, id)
		u.InsertUserToHash(user)

		// fmt.Println(user)
		// fmt.Println(id, username, fullname)  Can be used to debug whether the results are taken out properly or not.
	}
	err = user_result.Err()
	if err != nil {
		return false, 2 // code 2 : error from database transaction.
		panic(err)
	}
	defer user_result.Close()
	return true, 0 // code 0 : the program has retrieved the data properly.
}

func AuctionFromDBtoHash(a *AuctionHashTable) (bool, uint64) {

	db, err := sql.Open("mysql", "auction:Helloworld1@tcp(db.mcmullin.org)/auction_system")
	if err != nil {
		return false, 1 // code 1 : cannot connect to the database.
		panic(err.Error())
	}
	defer db.Close()

	auct_result, _ := db.Query("SELECT * FROM auction_table")

	for auct_result.Next() {
		var auction_id, auctioneer_id, currwin_id, currmax_bid, step uint64
		var itemname, currwinname, lastbidtime, starttime, endtime string
		err = auct_result.Scan(&auction_id, &auctioneer_id, &itemname, &currwin_id, &currwinname, &currmax_bid, &step, &lastbidtime, &starttime, &endtime)
		if err != nil {
			return false, 3 // code 3 : failed to read information from the query.
			panic(err)
		}
		auction := Auction{
			AuctionID:      auction_id,
			AuctioneerID:   auctioneer_id,
			ItemName:       itemname,
			CurrWinnerID:   currwin_id,
			CurrWinnerName: currwinname,
			CurrMaxBid:     currmax_bid,
			BidStep:        step,
			LatestBidTime:  lastbidtime,
			StartTime:      starttime,
			EndTime:        endtime,
		}
		a.InsertAuctToHash(&auction)

		//fmt.Println(auction)
		// fmt.Println(id, username, fullname)  Can be used to debug whether the results are taken out properly or not.
	}
	err = auct_result.Err()
	if err != nil {
		return false, 2 // code 2 : error from database transaction.
		panic(err)
	}
	defer auct_result.Close()

	return true, 0 // code 0 : the program has retrieved the data properly.
}
