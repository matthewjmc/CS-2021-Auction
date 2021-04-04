package AuctionSystem

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const Local_init string = "local_username:local_password@tcp(localhost:3306)/"
const Local_conn string = "local_username:local_password@tcp(localhost:3306)/auction_system"
const Server_init string = "auction:Helloworld1@tcp(db.mcmullin.org)/"
const Server_conn string = "auction:Helloworld1@tcp(db.mcmullin.org)/auction_system"

// _____DatabaseInit() functions are used to initialize the database along with its relationships upon each tables.
func LocalHostDatabaseInit() {

	db, debug := sql.Open("mysql", Local_init)
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

// _____DatabaseInit() functions are used to initialize the database along with its relationships upon each tables.
func ServerDatabaseInit() {

	db, debug := sql.Open("mysql", Server_init)
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
	statement, debug := db.Prepare("CREATE Table user_table( AccountID BIGINT UNSIGNED NOT NULL UNIQUE PRIMARY KEY, Username varchar(40) NOT NULL, Fullname varchar(30) NOT NULL )")
	if debug != nil {
		fmt.Println(debug.Error())
	}
	_, debug = statement.Exec()
	if debug != nil {
		fmt.Println(debug.Error())
	}
	statement2, debug := db.Prepare("CREATE Table auction_table( AuctionID BIGINT UNSIGNED NOT NULL UNIQUE PRIMARY KEY,AuctioneerID BIGINT UNSIGNED NOT NULL,ItemName varchar(40) NOT NULL, CurrWinnerID BIGINT UNSIGNED NOT NULL, CurrWinnerName varchar(40), CurrMaxBid BIGINT UNSIGNED NOT NULL, BidStep BIGINT UNSIGNED NOT NULL, LatestBidTime varchar(50) NOT NULL, StartTime varchar(50) NOT NULL, EndTime varchar(50) NOT NULL)")
	if debug != nil {
		fmt.Println(debug.Error())
	}
	_, debug = statement2.Exec()
	if debug != nil {
		fmt.Println(debug.Error())
	}
	statement3, debug := db.Prepare("CREATE Table bid_table( BiddingID BIGINT UNSIGNED AUTO_INCREMENT UNIQUE PRIMARY KEY, BidderID BIGINT UNSIGNED NOT NULL, BidderUsername varchar(30) NOT NULL, BidPrice BIGINT UNSIGNED NOT NULL, BidTime varchar(50) NOT NULL , AuctionID BIGINT UNSIGNED NOT NULL)")
	if debug != nil {
		fmt.Println(debug.Error())
	}
	_, debug = statement3.Exec()
	if debug != nil {
		fmt.Println(debug.Error())
	}

	defer db.Close()
}

func InsertAuctionToDB(auction Auction, db *sql.DB) bool {
	query, err := db.Prepare("INSERT INTO auction_table VALUES (?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	_, err = query.Exec(auction.AuctionID, auction.AuctioneerID, auction.ItemName, auction.CurrWinnerID, auction.CurrWinnerName, auction.CurrMaxBid, auction.BidStep, auction.LatestBidTime, auction.StartTime, auction.EndTime)
	if err != nil {
		fmt.Println("Insert Auction:", err)
	}
	//fmt.Println("Insert Auction:", result)
	defer query.Close()

	return true
}

func UpdateAuctionInDB(auction Auction, db *sql.DB) bool {

	update, err := db.Prepare("UPDATE auction_table SET CurrMaxBid = ? , CurrWinnerID = ? , CurrWinnerName = ? , LatestBidTime = ? WHERE AuctionID = ?")
	if err != nil {
		panic(err.Error())
	}

	update.Exec(auction.CurrMaxBid, auction.CurrWinnerID, auction.CurrWinnerName, auction.LatestBidTime, auction.AuctionID)
	defer update.Close()

	return true
}

func InsertUserToDB(user User, db *sql.DB) bool {

	query, err := db.Prepare("INSERT INTO user_table VALUES (?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	_, err = query.Exec(user.AccountID, user.Username, user.Fullname)
	if err != nil {
		fmt.Println("InsertUsers:", err)
	}
	//fmt.Println("InserUser:", result)
	defer query.Close()

	return true
}

func InsertBidToDB(bid Bid, target uint64, db *sql.DB) bool {

	query, err := db.Prepare("INSERT INTO bid_table(BidderID,BidderUsername,BidPrice,BidTime,AuctionID) VALUES ( ? , ? , ? , ? , ? )")
	if err != nil {
		panic(err)
	}

	_, err = query.Exec(bid.BidderID, bid.BidderUsername, bid.BidPrice, bid.BidTime, target)
	if err != nil {
		fmt.Println("InsertBid:", err)
	}
	//fmt.Println("InsertBid:", result)
	defer query.Close()

	return true
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func UserFromDBtoHash(u *UserHashTable) (bool, uint64) {
	db, err := sql.Open("mysql", "server_username:server_password@tcp(server_addr)/auction_system")
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

	db, err := sql.Open("mysql", "server_username:server_password@tcp(server_addr)/auction_system")
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