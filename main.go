package AuctionSystem

import (
	. "CS-2021-Auction/AuctionSystem"
	"fmt"
	"sync"
	"time"
)

// AuctionLogic is a file composed of Tagun's and Katisak's code altogether.
// The file is also being separated into 3 files which are auction_method, auction_timeline, data_structure.

func main() {
	//multipleUserTest()
	//testTimeFormat()
	A := AuctionAllocate()
	U := UserAllocate()
	// modification of memory allocation to be dynamically allocating.
	var uid uint64
	fmt.Println("What is your user ID")
	fmt.Scanln(&uid)
	for {
		mainTimeline(A, U, uid)
	}
}

func mainTimeline(A *AuctionHashTable, U *UserHashTable, uid uint64) {
	var command string

	fmt.Println("Command?")
	fmt.Scanln(&command)

	if command == "User" || command == "user" {
		if U.SearchUserIDHashTable(uid) {
			fmt.Println("The account has already been occupied")
		} else {
			userid_report := make(chan uint64)
			report_log := make(chan string)
			var fullname string
			fmt.Println("fullname")
			fmt.Scanln(&fullname)
			go createUserMain(U, userid_report, report_log, uid, fullname)
			newUser := <-userid_report
			log := <-report_log
			fmt.Println(log, newUser)
		}
	} else if command == "Auction" || command == "auction" {
		report_auction := make(chan Auction)
		report_log := make(chan string)
		var initbid, step uint64
		fmt.Println("enter initbid")
		fmt.Scanln(&initbid)
		fmt.Println("enter step")
		fmt.Scanln(&step)
		go createAuctionMain(U, A, report_auction, report_log, uid, Randomize(100, 999), initbid, step)
		newAuction := <-report_auction
		log := <-report_log
		fmt.Println(log)
		fmt.Println("New created auction is having the id of", newAuction.AuctionID)
		fmt.Println(newAuction)
	} else if command == "bid" {
		var tid, price uint64
		fmt.Println("What is your target auction id")
		fmt.Scanln(&tid)
		if !A.SearchAuctIDHashTable(tid) {
			fmt.Println("The auction has not been found within the memory")
		} else {
			fmt.Println("price?")
			fmt.Scanln(&price)
			report_price := make(chan uint64)
			report_log := make(chan string)
			go makeBidMain(U, A, report_price, report_log, uid, tid, price)
			finalAuction := <-report_price
			log := <-report_log
			fmt.Println(finalAuction, log)
		}
		time.Sleep(1 * time.Millisecond)
	}

}

var wg sync.WaitGroup

func makeBidMain(u *UserHashTable, h *AuctionHashTable, report_price chan uint64, report_log chan string, uid uint64, targetid uint64, placeVal uint64) {

	bidTime := time.Now().Format(time.RFC3339Nano)
	if !u.SearchUserIDHashTable(uid) {
		report_price <- 0
		report_log <- "The user could not be located within the system."
	}
	if !h.SearchAuctIDHashTable(targetid) {
		report_price <- 0
		report_log <- "The auction could not be located within the system."
	}
	currUser := *u.AccessUserHash(uid)
	newBid := CreateBid(currUser, placeVal, bidTime)
	target := h.AccessHashAuction(targetid)
	target.UpdateAuctionWinner(newBid)
	h.AuctionHashAccessUpdate(*target)
	fmt.Println("printed from the function, The auction has been updated successfully")
	report_price <- target.CurrMaxBid
	report_log <- "The auction has been updated completely"
}

func createUserMain(h *UserHashTable, report chan uint64, report_log chan string, uid uint64, name string) {
	if !h.SearchUserIDHashTable(uid) {
		newUser := CreateUser("username"+fmt.Sprint(uid), name, uid)
		h.InsertUserToHash(newUser)
		fmt.Println("New user has been inserted properly")
		report <- newUser.AccountID
		report_log <- "account has been created completely"
	} else {
		fmt.Println("The user has already registered into the system")
		report <- uid
		report_log <- "account has already been created"
	}
}

func createAuctionMain(U *UserHashTable, A *AuctionHashTable, auction chan Auction, report_log chan string, uid uint64, aid uint64, initial uint64, step uint64) {
	user := U.AccessUserHash(uid)
	if !A.SearchAuctIDHashTable(aid) {
		newAuction := CreateAuction(*user, initial, step, aid)
		A.InsertAuctToHash(newAuction.CreatedAuction)
		auction <- *newAuction.CreatedAuction
		report_log <- "auction has been created completely"
	} else {
		auction <- *A.AccessHashAuction(aid)
		report_log <- "auction has already been declared"
	}
}
