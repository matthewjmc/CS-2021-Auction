package main

import (
	. "CS-2021-Auction/AuctionSystem"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	testingFinal()
}

func testingFinal() {
	ServerDatabaseInit()
	a := AuctionAllocate()
	u := UserAllocate()

	db, err := sql.Open("mysql", "server_username:server_password@tcp(server_addr)/auction_system")
	if err != nil {
		panic(err.Error())
	}
	CreateUserMain(u, 1777, "matthew", db)
	CreateUserMain(u, 3257, "maeluenie", db)
	CreateUserMain(u, 9921, "tagun", db)
	time.Sleep(1 * time.Second)
	CreateAuctionMain(u, a, 9921, 555555, 100, 100, 1, "GOGOGOPHER", db)
	time.Sleep(1 * time.Second)
	MakeBidMain(u, a, 9921, 555555, 300, 1, db)
	time.Sleep(1 * time.Second)
	MakeBidMain(u, a, 3257, 555555, 550, 3, db)
	time.Sleep(1 * time.Second)
	MakeBidMain(u, a, 1777, 555555, 350, 2, db)
	time.Sleep(1 * time.Second)
	MakeBidMain(u, a, 3257, 222222, 300, 4, db)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 3 Main Functions for business logic usage.
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func MakeBidMain(u *UserHashTable, h *AuctionHashTable, uid uint64, targetid uint64, placeVal uint64, bidId uint64, db *sql.DB) {
	currUser := *u.AccessUserHash(uid)
	newBid := CreateBid(currUser, placeVal, bidId, time.Now().Format(time.RFC3339Nano))
	target := h.AccessHashAuction(targetid)
	go InsertBidToDB(newBid, targetid, db)
	if !target.UpdateAuctionWinner(newBid) {
	} else {
		go UpdateAuctionInDB(*target, db)
		h.AuctionHashAccessUpdate(*target)
	}
}

func CreateUserMain(h *UserHashTable, uid uint64, name string, db *sql.DB) {
	if !h.SearchUserIDHashTable(uid) {
		newUser := CreateUser("username"+fmt.Sprint(uid), name, uid)
		go InsertUserToDB(newUser, db)
		h.InsertUserToHash(newUser)
	}
}

// Average time to create and store new user into both server-side caching and database is approximately 7.5 ms. ( 0.0075 second )
// Considering actual networking latencies, inserting the information onto the actual database is approximately 155 ms. ( 0.155 second )

func CreateAuctionMain(U *UserHashTable, A *AuctionHashTable, uid uint64, aid uint64, initial uint64, step uint64, duration time.Duration, itemName string, db *sql.DB) {
	user := U.AccessUserHash(uid)
	if !A.SearchAuctIDHashTable(aid) {
		newAuction := CreateAuction(*user, initial, step, aid, duration, itemName)
		go InsertAuctionToDB(newAuction, db)
		A.InsertAuctToHash(&newAuction)
	}
}
