package main

import (
	. "CS-2021-Auction/AuctionSystem"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	DatabaseInit()
	a := AuctionAllocate()
	u := UserAllocate()
	test_database_transaction1(u, a)
	test_database_transaction2(u, a)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 3 Main Functions for business logic usage.
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var wg sync.WaitGroup

func MakeBidMain(u *UserHashTable, h *AuctionHashTable, uid uint64, targetid uint64, placeVal uint64, bidId uint64) (bool, uint64, bool) {
	if !u.SearchUserIDHashTable(uid) {
		return false, 1, false // code 1 , the user has not been found within the system.
	} else if !h.SearchAuctIDHashTable(targetid) {
		return false, 2, false // code 2 , the auction has not been found within the system.
	} else {
		currUser := *u.AccessUserHash(uid)
		newBid := CreateBid(currUser, placeVal, bidId, time.Now().Format(time.RFC3339Nano))
		target := h.AccessHashAuction(targetid)
		bid_report := make(chan bool)
		go InsertBidToDB(newBid, *target, bid_report)
		bidstmt_result := <-bid_report
		if !target.UpdateAuctionWinner(newBid) {
			return false, 3, bidstmt_result // code 3 , the auction has not been updated due to the losing auctionconditions.
		} else {
			update_report := make(chan bool)
			go UpdateAuctionInDB(*target, update_report)
			update_result := <-update_report
			h.AuctionHashAccessUpdate(*target)
			return update_result, 0, bidstmt_result
		}
	}
}

func CreateUserMain(h *UserHashTable, uid uint64, name string) (bool, uint64) {
	if !h.SearchUserIDHashTable(uid) {
		newUser := CreateUser("username"+fmt.Sprint(uid), name, uid)
		h.InsertUserToHash(newUser)
		report := make(chan bool)
		go InsertUserToDB(newUser, report)
		stmt_result := <-report
		return stmt_result, 0 // code 0 , the user has not been found within the system, creating new user object.
	} else {
		return false, 1 // code 1 , the user has been found in the system.
	}
}

func CreateAuctionMain(U *UserHashTable, A *AuctionHashTable, uid uint64, aid uint64, initial uint64, step uint64, duration time.Duration, itemName string) (bool, uint64) {
	user := U.AccessUserHash(uid)
	if !A.SearchAuctIDHashTable(aid) {
		newAuction := CreateAuction(*user, initial, step, aid, duration, itemName)
		report := make(chan bool)
		go InsertAuctionToDB(newAuction, report)
		stmt_result := <-report
		A.InsertAuctToHash(&newAuction)
		return stmt_result, 0 // code 0, auction has not been found within the system, creating new auction object.
	}
	return false, 1 // error code 1 , auction has been found in the system.
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Testing functions.
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func test_database_transaction1(u *UserHashTable, a *AuctionHashTable) {
	init := time.Now()
	CreateUserMain(u, 9921, "tagun")   // db checked!
	CreateUserMain(u, 1338, "matthew") // db checked!
	CreateUserMain(u, 7777, "katisak") // db checked!
	CreateAuctionMain(u, a, 9921, 111111, 100, 100, 1, "verygooditem1")
	/* MakeBidMain(u, a, 7777, 111111, 120, 1)
	MakeBidMain(u, a, 1338, 111111, 400, 2)
	MakeBidMain(u, a, 7777, 111111, 1000, 3)
	MakeBidMain(u, a, 9921, 222222, 1000, 4) */
	fmt.Println("Time taken with goroutine logic:", time.Since(init))
}

func test_database_transaction2(u *UserHashTable, a *AuctionHashTable) {
	init := time.Now()
	CreateUserMain_Original(u, 2222, "nonthicha") // db checked!
	CreateUserMain_Original(u, 4444, "maeluenie") // db checked!
	CreateUserMain_Original(u, 8888, "vorachat")  // db checked!
	CreateAuctionMain_Original(u, a, 8888, 999999, 100, 100, 9, "badbadGopher")
	/* MakeBidMain_Original(u, a, 7777, 111111, 120, 1)
	MakeBidMain_Original(u, a, 1338, 111111, 400, 2)
	MakeBidMain_Original(u, a, 7777, 111111, 1000, 3)
	MakeBidMain_Original(u, a, 9921, 222222, 1000, 4) */
	fmt.Println("Time taken for the none goroutine logic:", time.Since(init))
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

func MakeBidMain_Original(u *UserHashTable, h *AuctionHashTable, uid uint64, targetid uint64, placeVal uint64, bidId uint64) (bool, uint64) {
	if !u.SearchUserIDHashTable(uid) {
		return false, 1 // code 1 , the user has not been found within the system.
	} else if !h.SearchAuctIDHashTable(targetid) {
		return false, 2 // code 2 , the auction has not been found within the system.
	} else {
		currUser := *u.AccessUserHash(uid)
		newBid := CreateBid(currUser, placeVal, bidId, time.Now().Format(time.RFC3339Nano))
		target := h.AccessHashAuction(targetid)
		InsertBidToDB_Original(newBid, *target)
		if !target.UpdateAuctionWinner(newBid) {
			return false, 3 // code 3 , the auction has not been updated due to the losing auctionconditions.
		} else {
			h.AuctionHashAccessUpdate(*target)
			UpdateAuctionInDB_Original(*target)
			return true, 0
		}
	}
}

func CreateUserMain_Original(h *UserHashTable, uid uint64, name string) (bool, uint64) {
	if !h.SearchUserIDHashTable(uid) {
		newUser := CreateUser("username"+fmt.Sprint(uid), name, uid)
		h.InsertUserToHash(newUser)
		InsertUserToDB_Original(newUser)
		return true, 0 // code 0 , the user has not been found within the system, creating new user object.
	} else {
		return false, 1 // code 1 , the user has been found in the system.
	}
}

func CreateAuctionMain_Original(U *UserHashTable, A *AuctionHashTable, uid uint64, aid uint64, initial uint64, step uint64, duration time.Duration, itemName string) (bool, uint64) {
	user := U.AccessUserHash(uid)
	if !A.SearchAuctIDHashTable(aid) {
		newAuction := CreateAuction(*user, initial, step, aid, duration, itemName)
		A.InsertAuctToHash(&newAuction)
		InsertAuctionToDB_Original(newAuction)
		return true, 0 // code 0, auction has not been found within the system, creating new auction object.
	}
	return false, 1 // error code 1 , auction has been found in the system.
}
