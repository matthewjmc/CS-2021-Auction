package main

import (
	. "CS-2021-Auction/AuctionSystem"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

/*
func main() {
	DatabaseInit()
	u := UserAllocate()
	CreateUserMain(u, 8888, "nonthicha")
}
*/

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 3 Main Functions for business logic usage.
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var wg sync.WaitGroup

func MakeBidMain(u *UserHashTable, h *AuctionHashTable, uid uint64, targetid uint64, placeVal uint64, bidId uint64) (bool, uint64) {
	if !u.SearchUserIDHashTable(uid) {
		return false, 1 // code 1 , the user has not been found within the system.
	} else if !h.SearchAuctIDHashTable(targetid) {
		return false, 2 // code 2 , the auction has not been found within the system.
	} else {
		currUser := *u.AccessUserHash(uid)
		newBid := CreateBid(currUser, placeVal, bidId, time.Now().Format(time.RFC3339Nano))
		target := h.AccessHashAuction(targetid)
		InsertBidToDB(newBid, *target)
		if !target.UpdateAuctionWinner(newBid) {
			return false, 3 // code 3 , the auction has not been updated due to the losing auctionconditions.
		} else {
			h.AuctionHashAccessUpdate(*target)
			UpdateAuctionInDB(*target)
			return true, 0
		}
	}
}

func CreateUserMain(h *UserHashTable, uid uint64, name string) (bool, uint64) {
	if !h.SearchUserIDHashTable(uid) {
		newUser := CreateUser("username"+fmt.Sprint(uid), name, uid)
		h.InsertUserToHash(newUser)
		InsertUserToDB(newUser)
		return true, 0 // code 0 , the user has not been found within the system, creating new user object.
	} else {
		return false, 1 // code 1 , the user has been found in the system.
	}
}

func CreateAuctionMain(U *UserHashTable, A *AuctionHashTable, uid uint64, aid uint64, initial uint64, step uint64, duration time.Duration, itemName string) (bool, uint64) {
	user := U.AccessUserHash(uid)
	if !A.SearchAuctIDHashTable(aid) {
		newAuction := CreateAuction(*user, initial, step, aid, duration, itemName)
		A.InsertAuctToHash(&newAuction)
		InsertAuctionToDB(newAuction)
		return true, 0 // code 0, auction has not been found within the system, creating new auction object.
	}
	return false, 1 // error code 1 , auction has been found in the system.
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Testing functions.
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func test_database_conn() {
	user := User{
		AccountID: 1,
		Username:  "testUsername",
		Fullname:  "testFullname",
	}
	auction := Auction{
		AuctionID:      111111,
		AuctioneerID:   1,
		ItemName:       "test",
		CurrWinnerID:   1,
		CurrWinnerName: "test2",
		CurrMaxBid:     100,
		BidStep:        200,
		LatestBidTime:  time.Now().Format(time.RFC3339Nano),
		StartTime:      time.Now().Format(time.RFC3339Nano),
		EndTime:        time.Now().Add(1 * time.Hour).Format(time.RFC3339Nano),
	}
	bid := Bid{
		BiddingID:      111,
		BidderID:       user.AccountID,
		BidderUsername: user.Username,
		BidPrice:       10000,
		BidTime:        time.Now().Format(time.RFC3339Nano),
	}
	InsertUserToDB(user)
	InsertAuctionToDB(auction)
	InsertBidToDB(bid, auction)
}

func test_database_transaction() {
	DatabaseInit()
	a := AuctionAllocate()
	u := UserAllocate()
	CreateUserMain(u, 9921, "tagun")   // db checked!
	CreateUserMain(u, 1338, "matthew") // db checked!
	CreateUserMain(u, 7777, "katisak") // db checked!
	CreateAuctionMain(u, a, 9921, 111111, 100, 100, 1, "verygooditem1")
	MakeBidMain(u, a, 7777, 111111, 120, 1)
	MakeBidMain(u, a, 1338, 111111, 400, 2)
	MakeBidMain(u, a, 7777, 111111, 1000, 3)
	MakeBidMain(u, a, 9921, 222222, 1000, 4)
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
