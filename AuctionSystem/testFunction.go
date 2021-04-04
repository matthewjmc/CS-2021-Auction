package AuctionSystem

import (
	"database/sql"
	"fmt"
	"time"
)

func TestLocalHost() {
	LocalHostDatabaseInit()
	a := AuctionAllocate()
	u := UserAllocate()

	db, err := sql.Open("mysql", Local_conn)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("=========================================================================================================================")
	fmt.Println("Time used to create users  =================================================================")
	createUserMain_test(u, 1777, "matthew", db)
	createUserMain_test(u, 3257, "maeluenie", db)
	createUserMain_test(u, 9921, "tagun", db)
	fmt.Println("Time used to create auctions  ==============================================================")
	createAuctionMain_test(u, a, 9921, 555555, 100, 100, 1, "GOGOGOPHER", db)
	fmt.Println("Time used to create bidding and update the auction  ========================================")
	makeBidMain_test(u, a, 9921, 555555, 300, 1, db)
	makeBidMain_test(u, a, 3257, 555555, 550, 3, db)
	fmt.Println("Time used to complete a bidding without WINNING the auction  ===============================")
	makeBidMain_test(u, a, 3257, 555555, 350, 2, db)
	fmt.Println("Time used to create a bidding that could not update any auction  ===========================")
	makeBidMain_test(u, a, 3257, 222222, 300, 4, db)
	fmt.Println("=========================================================================================================================")
}

func TestServer() {
	ServerDatabaseInit()
	a := AuctionAllocate()
	u := UserAllocate()

	db, err := sql.Open("mysql", Server_init)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("=========================================================================================================================")
	fmt.Println("Time used to create users  =================================================================")
	createUserMain_test(u, 1777, "matthew", db)
	createUserMain_test(u, 3257, "maeluenie", db)
	createUserMain_test(u, 9921, "tagun", db)
	fmt.Println("Time used to create auctions  ==============================================================")
	createAuctionMain_test(u, a, 9921, 555555, 100, 100, 1, "GOGOGOPHER", db)
	fmt.Println("Time used to create bidding and update the auction  ========================================")
	makeBidMain_test(u, a, 9921, 555555, 300, 1, db)
	makeBidMain_test(u, a, 3257, 555555, 550, 3, db)
	fmt.Println("Time used to complete a bidding without WINNING the auction  ===============================")
	makeBidMain_test(u, a, 3257, 555555, 350, 2, db)
	fmt.Println("Time used to create a bidding that could not update any auction  ===========================")
	makeBidMain_test(u, a, 3257, 222222, 300, 4, db)
	fmt.Println("=========================================================================================================================")
}

func makeBidMain_test(u *UserHashTable, h *AuctionHashTable, uid uint64, targetid uint64, placeVal uint64, bidId uint64, db *sql.DB) (bool, uint64, bool) {

	init := time.Now()

	currUser := *u.AccessUserHash(uid)
	newBid := CreateBid(currUser, placeVal, time.Now().Format(time.RFC3339Nano))
	target := h.AccessHashAuction(targetid)

	bidResult := make(chan bool)
	bidDur := make(chan time.Duration)
	go insertBid_test(newBid, targetid, db, bidResult, bidDur)
	fmt.Println("Insert Bid to DB:", <-bidDur)

	if !target.UpdateAuctionWinner(newBid) {
		final := time.Now()
		fmt.Println("makeBidMain_test() , losing bid :", final.Sub(init))

		return false, 3, <-bidResult // code 3 , the auction has not been updated due to the losing auction conditions.

	} else {
		auctDur := make(chan time.Duration)
		auctResult := make(chan bool)
		go updateAuction_test(*target, db, auctResult, auctDur)

		h.AuctionHashAccessUpdate(*target)
		final := time.Now()
		fmt.Println("Update Auction to DB:", <-auctDur)
		fmt.Println("makeBidMain_test() , winning bid :", final.Sub(init))

		return <-auctResult, 0, <-bidResult
	}

}

func createUserMain_test(h *UserHashTable, uid uint64, name string, db *sql.DB) (bool, uint64) {
	init := time.Now()
	if !h.SearchUserIDHashTable(uid) {
		newUser := CreateUser("username"+fmt.Sprint(uid), name, uid)
		result := make(chan bool)
		dur := make(chan time.Duration)
		go insertUser_test(newUser, db, result, dur)
		taken := <-dur
		h.InsertUserToHash(newUser)
		final := time.Now()
		fmt.Println("Insert User to DB :", taken)
		fmt.Println("createUserMain_test() :", final.Sub(init))
		return <-result, 0 // code 0 , the user has not been found within the system, creating new user object.
	} else {
		return false, 1 // code 1 , the user has been found in the system.
	}
}

// Average time to create and store new user into both server-side caching and database is approximately 7.5 ms. ( 0.0075 second )
// Considering actual networking latencies, inserting the information onto the actual database is approximately 155 ms. ( 0.155 second )

func createAuctionMain_test(U *UserHashTable, A *AuctionHashTable, uid uint64, aid uint64, initial uint64, step uint64, duration time.Duration, itemName string, db *sql.DB) (bool, uint64) {
	init := time.Now()
	user := U.AccessUserHash(uid)
	if !A.SearchAuctIDHashTable(aid) {
		newAuction := CreateAuction(*user, initial, step, aid, duration, itemName)
		result := make(chan bool)
		taken := make(chan time.Duration)
		go insertAuction_test(newAuction, db, result, taken)
		dur := <-taken
		A.InsertAuctToHash(&newAuction)
		final := time.Now()
		fmt.Println("Insert Auction to DB:", dur)
		fmt.Println("createAuctionMain_test() :", final.Sub(init))
		return <-result, 0 // code 0, auction has not been found within the system, creating new auction object.
	}
	return false, 1 // error code 1 , auction has been found in the system.
}

func insertAuction_test(auction Auction, db *sql.DB, result chan bool, taken chan time.Duration) {
	init := time.Now()
	query, err := db.Prepare("INSERT INTO auction_table VALUES (?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	query.Exec(auction.AuctionID, auction.AuctioneerID, auction.ItemName, auction.CurrWinnerID, auction.CurrWinnerName, auction.CurrMaxBid, auction.BidStep, auction.LatestBidTime, auction.StartTime, auction.EndTime)
	defer query.Close()

	taken <- time.Since(init)
	result <- true
}

func updateAuction_test(auction Auction, db *sql.DB, result chan bool, taken chan time.Duration) {

	init := time.Now()

	update, err := db.Prepare("UPDATE auction_table SET CurrMaxBid = ? , CurrWinnerID = ? , CurrWinnerName = ? , LatestBidTime = ? WHERE AuctionID = ?")
	if err != nil {
		panic(err.Error())
	}

	update.Exec(auction.CurrMaxBid, auction.CurrWinnerID, auction.CurrWinnerName, auction.LatestBidTime, auction.AuctionID)
	defer update.Close()

	taken <- time.Since(init)
	result <- true
}

func insertUser_test(user User, db *sql.DB, result chan bool, taken chan time.Duration) {

	init := time.Now()

	query, err := db.Prepare("INSERT INTO user_table VALUES (?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	query.Exec(user.AccountID, user.Username, user.Fullname)
	defer query.Close()

	taken <- time.Since(init)
	result <- true
}

func insertBid_test(bid Bid, target uint64, db *sql.DB, result chan bool, taken chan time.Duration) {

	init := time.Now()

	query, err := db.Prepare("INSERT INTO bid_table VALUES ( ? , ? , ? , ? , ? )")
	if err != nil {
		panic(err)
	}

	query.Exec(bid.BidderID, bid.BidderUsername, bid.BidPrice, bid.BidTime, target)
	defer query.Close()

	taken <- time.Since(init)
	result <- true
}
