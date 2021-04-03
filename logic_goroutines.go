package main

import (
	. "CS-2021-Auction/AuctionSystem"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	DatabaseInit()
	a := AuctionAllocate()
	u := UserAllocate()

	db, err := sql.Open("mysql", "auction:Helloworld1@tcp(db.mcmullin.org)/auction_system")
	if err != nil {
		panic(err.Error())
	}
	test_database_transaction1(u, a, db)
	//test_database_transaction2(u, a, db)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 3 Main Functions for business logic usage.
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func MakeBidMain(u *UserHashTable, h *AuctionHashTable, uid uint64, targetid uint64, placeVal uint64, bidId uint64, db *sql.DB) (bool, uint64, bool) {
	if !u.SearchUserIDHashTable(uid) {
		return false, 1, false // code 1 , the user has not been found within the system.
	} else if !h.SearchAuctIDHashTable(targetid) {
		return false, 2, false // code 2 , the auction has not been found within the system.
	} else {
		currUser := *u.AccessUserHash(uid)
		newBid := CreateBid(currUser, placeVal, bidId, time.Now().Format(time.RFC3339Nano))
		target := h.AccessHashAuction(targetid)
		go InsertBidToDB(newBid, targetid, db)
		if !target.UpdateAuctionWinner(newBid) {
			return false, 3, true // code 3 , the auction has not been updated due to the losing auctionconditions.
		} else {
			go UpdateAuctionInDB(*target, db)
			h.AuctionHashAccessUpdate(*target)
			return true, 0, true
		}
	}
}

func CreateUserMain(h *UserHashTable, uid uint64, name string, db *sql.DB) (bool, uint64) {
	if !h.SearchUserIDHashTable(uid) {
		newUser := CreateUser("username"+fmt.Sprint(uid), name, uid)
		h.InsertUserToHash(newUser)
		go InsertUserToDB(newUser, db)
		return true, 0 // code 0 , the user has not been found within the system, creating new user object.
	} else {
		return false, 1 // code 1 , the user has been found in the system.
	}
}

func CreateAuctionMain(U *UserHashTable, A *AuctionHashTable, uid uint64, aid uint64, initial uint64, step uint64, duration time.Duration, itemName string, db *sql.DB) (bool, uint64) {
	user := U.AccessUserHash(uid)
	if !A.SearchAuctIDHashTable(aid) {
		newAuction := CreateAuction(*user, initial, step, aid, duration, itemName)
		go InsertAuctionToDB(newAuction, db)
		A.InsertAuctToHash(&newAuction)
		return true, 0 // code 0, auction has not been found within the system, creating new auction object.
	}
	return false, 1 // error code 1 , auction has been found in the system.
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Testing functions.
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func test_database_transaction1(u *UserHashTable, a *AuctionHashTable, db *sql.DB) {
	init := time.Now().Format(time.RFC3339Nano)
	fmt.Println("Time completed with goroutine logic:", init)
	CreateUserMain(u, 9921, "tagun", db) // db checked!
	//CreateUserMain(u, 1338, "matthew", db) // db checked!
	//CreateUserMain(u, 7777, "katisak", db) // db checked!
	//CreateAuctionMain(u, a, 9921, 111111, 100, 100, 1, "verygooditem1", db)
	//MakeBidMain(u, a, 7777, 111111, 120, 1, db)
	//MakeBidMain(u, a, 1338, 111111, 400, 2, db)
	//MakeBidMain(u, a, 7777, 111111, 6746, 3, db)
	//MakeBidMain(u, a, 9921, 111111, 500, 4, db)
	end := time.Now().Format(time.RFC3339Nano)
	fmt.Println("Time completed with goroutine logic:", end)
}

func test_retrieve_fromDB() {
	u := UserAllocate()
	a := AuctionAllocate()
	UserFromDBtoHash(u)
	u.SearchUserIDHashTable(9921)
	u.SearchUserIDHashTable(7777)
	u.SearchUserIDHashTable(1338)
	fmt.Println(u.SearchUserIDHashTable(3333))
	fmt.Println(*u.AccessUserHash(9921))
	fmt.Println(*u.AccessUserHash(7777))
	fmt.Println(*u.AccessUserHash(1338))
	AuctionFromDBtoHash(a)
	fmt.Println(a.SearchAuctIDHashTable(111111))
	fmt.Println(a.AccessHashAuction(111111))
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 3 Main Functions for business logic usage without concurrencies.
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func MakeBidMain_Original(u *UserHashTable, h *AuctionHashTable, uid uint64, targetid uint64, placeVal uint64, bidId uint64, db *sql.DB) (bool, uint64, bool) {
	if !u.SearchUserIDHashTable(uid) {
		return false, 1, false // code 1 , the user has not been found within the system.
	} else if !h.SearchAuctIDHashTable(targetid) {
		return false, 2, false // code 2 , the auction has not been found within the system.
	} else {
		currUser := *u.AccessUserHash(uid)
		newBid := CreateBid(currUser, placeVal, bidId, time.Now().Format(time.RFC3339Nano))
		target := h.AccessHashAuction(targetid)
		InsertBidToDB(newBid, target.AuctionID, db)
		if !target.UpdateAuctionWinner(newBid) {
			return false, 3, true // code 3 , the auction has not been updated due to the losing auctionconditions.
		} else {
			UpdateAuctionInDB(*target, db)
			h.AuctionHashAccessUpdate(*target)
			return true, 0, true
		}
	}
}

func CreateUserMain_Original(h *UserHashTable, uid uint64, name string, db *sql.DB) (bool, uint64) {
	if !h.SearchUserIDHashTable(uid) {
		newUser := CreateUser("username"+fmt.Sprint(uid), name, uid)
		h.InsertUserToHash(newUser)
		InsertUserToDB(newUser, db)
		return true, 0 // code 0 , the user has not been found within the system, creating new user object.
	} else {
		return false, 1 // code 1 , the user has been found in the system.
	}
}

func CreateAuctionMain_Original(U *UserHashTable, A *AuctionHashTable, uid uint64, aid uint64, initial uint64, step uint64, duration time.Duration, itemName string, db *sql.DB) (bool, uint64) {
	user := U.AccessUserHash(uid)
	if !A.SearchAuctIDHashTable(aid) {
		newAuction := CreateAuction(*user, initial, step, aid, duration, itemName)
		A.InsertAuctToHash(&newAuction)
		InsertAuctionToDB(newAuction, db)
		return true, 0 // code 0, auction has not been found within the system, creating new auction object.
	}
	return false, 1 // error code 1 , auction has been found in the system.
}

func test_database_transaction2(u *UserHashTable, a *AuctionHashTable, db *sql.DB) {
	init := time.Now().Format(time.RFC3339Nano)
	fmt.Println("Time completed with none goroutine logic:", init)
	CreateUserMain_Original(u, 2222, "nonthicha", db) // db checked!
	CreateUserMain_Original(u, 4444, "maeluenie", db) // db checked!
	CreateUserMain_Original(u, 8888, "vorachat", db)  // db checked!
	CreateAuctionMain_Original(u, a, 8888, 999999, 100, 100, 9, "badbadGopher", db)
	MakeBidMain_Original(u, a, 7777, 999999, 120, 5, db)
	MakeBidMain_Original(u, a, 1338, 999999, 400, 6, db)
	MakeBidMain_Original(u, a, 7777, 999999, 9483, 7, db)
	MakeBidMain_Original(u, a, 9921, 999999, 1000, 8, db)
	end := time.Now().Format(time.RFC3339Nano)
	fmt.Println("Time completed with none goroutine logic:", end)
}
