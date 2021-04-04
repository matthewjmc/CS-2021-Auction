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

	db, err := sql.Open("mysql", Server_conn)
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
	fmt.Println(u.AccessUserHash(9921))
	time.Sleep(1 * time.Second)
	MakeBidMain(u, a, 3257, 555555, 550, 3, db)
	time.Sleep(1 * time.Second)
	MakeBidMain(u, a, 1777, 555555, 350, 2, db)
	time.Sleep(1 * time.Second)
	fmt.Println(u.AccessUserHash(9921))
	MakeBidMain(u, a, 3257, 222222, 300, 4, db)
	time.Sleep(1 * time.Second)
	u.SearchUserIDHashTable(9921)
	a.SearchAuctIDHashTable(555555)
	fmt.Println(u.AccessUserHash(9921))
	fmt.Println(a.AccessHashAuction(555555))
	time.Sleep(1 * time.Second)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 3 Main Functions for business logic usage.
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func MakeBidMain(u *UserHashTable, h *AuctionHashTable, uid uint64, targetid uint64, placeVal uint64, bidId uint64, db *sql.DB) (bool, Auction, Bid) {
	currUser := *u.AccessUserHash(uid)
	newBid := CreateBid(currUser, placeVal, time.Now().Format(time.RFC3339Nano))
	target := h.AccessHashAuction(targetid)
	go InsertBidToDB(newBid, targetid, db)
	if !target.UpdateAuctionWinner(newBid) {
		return false, *target, newBid
	} else {
		go UpdateAuctionInDB(*target, db)
		go h.AuctionHashAccessUpdate(*target)
		return true, *target, newBid
	}
}

func CreateUserMain(h *UserHashTable, uid uint64, name string, db *sql.DB) (bool, User) {
	if !h.SearchUserIDHashTable(uid) {
		newUser := CreateUser("username"+fmt.Sprint(uid), name, uid)
		go InsertUserToDB(newUser, db)
		go h.InsertUserToHash(newUser)
		return true, newUser
	}
	return false, User{}
}

func CreateAuctionMain(U *UserHashTable, A *AuctionHashTable, uid uint64, aid uint64, initial uint64, step uint64, duration time.Duration, itemName string, db *sql.DB) (bool, Auction) {
	user := U.AccessUserHash(uid)
	if !A.SearchAuctIDHashTable(aid) {
		newAuction := CreateAuction(*user, initial, step, aid, duration, itemName)
		go InsertAuctionToDB(newAuction, db)
		go A.InsertAuctToHash(&newAuction)
		return true, newAuction
	}
	return false, Auction{}
}
