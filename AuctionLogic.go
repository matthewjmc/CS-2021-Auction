package main

import (
	. "CS-2021-Auction/AuctionSystem"
	"fmt"
	"sync"
	"time"
)

// AuctionLogic is a file composed of Tagun's and Katisak's code altogether.
// The file is also being separated into 3 files which are auction_method, auction_timeline, data_structure.

/*
func main() {
	//multipleUserTest()
	//testTimeFormat()
	A := AuctionAllocate()
	U := UserAllocate()
	// modification of memory allocation to be dynamically allocating.

	mainTimeline(A, U)

}*/
type Data struct {
	command      string
	uid          uint64
	fullname     string
	aid          uint64
	itemname     string
	biddingValue uint64
	biddingStep  uint64
}

func mainTimeline(A *AuctionHashTable, U *UserHashTable, instructions Data) {

	command := instructions.command

	if command == "User" || command == "user" {
		user_report := make(chan User)
		report_log := make(chan string)
		go createUserMain(U, user_report, report_log, instructions.uid, instructions.fullname)
		newUser := <-user_report
		log := <-report_log
		fmt.Println(log, newUser)

	} else if command == "Auction" || command == "auction" {
		report_auction := make(chan Auction)
		report_log := make(chan string)
		go createAuctionMain(U, A, report_auction, report_log, instructions.uid, instructions.aid, instructions.biddingValue, instructions.biddingStep)
		newAuction := <-report_auction
		log := <-report_log
		fmt.Println(newAuction, log)
	} else if command == "bid" {

		if !A.SearchAuctIDHashTable(instructions.aid) {
			fmt.Println("The auction has not been found within the memory")
		} else {
			report_price := make(chan uint64)
			report_log := make(chan string)
			go makeBidMain(U, A, report_price, report_log, instructions.uid, instructions.aid, instructions.biddingValue)
			finalAuction := <-report_price
			log := <-report_log
			fmt.Println(finalAuction, log)
		}
	} else if command == "searchAuction" {
		if A.SearchAuctIDHashTable(instructions.aid) {
			fmt.Println("Auction", instructions.aid, " is found within the system")
		} else {
			fmt.Println("Auction", instructions.aid, " is not found within the system")
		}
	} else if command == "deleteAuction" {
		if A.AuctionHashAccessDelete(instructions.aid) {
			fmt.Println("Auction", instructions.aid, " has been deleted for the system")
		} else {
			fmt.Println("Auction", instructions.aid, " is not found within the system")
		}
	} else if command == "searchUser" {
		if U.SearchUserIDHashTable(instructions.uid) {
			fmt.Println("Auction", instructions.aid, " is found within the system")
		} else {
			fmt.Println("Auction", instructions.aid, " is not found within the system")
		}
	} else if command == "deleteUser" {
		if U.UserHashAccessDelete(instructions.uid) {
			fmt.Println("Auction", instructions.aid, " has been deleted for the system")
		} else {
			fmt.Println("Auction", instructions.aid, " is not found within the system")
		}
	}

}

var wg sync.WaitGroup

func makeBidMain(u *UserHashTable, h *AuctionHashTable, report_price chan uint64, report_log chan string, uid uint64, targetid uint64, placeVal uint64) {

	bidTime := time.Now().Format(time.RFC3339Nano)
	if u.SearchUserIDHashTable(uid) == false {
		fmt.Println("The user could not be located within the system.")
	} else { // for testing
		currUser := *u.AccessUserHash(uid)

		newBid := CreateBid(currUser, placeVal, bidTime)

		// access for auction object to be updated at the target variable.
		target := h.AccessHashAuction(targetid)
		//fmt.Println("Previous Winner:", target.currWinnerName)
		target.UpdateAuctionWinner(newBid)
		h.AuctionHashAccessUpdate(*target)
		//fmt.Println("Current Winner:", target.currWinnerName)
		report_price <- target.CurrMaxBid // This line is used to notate new user created.
		report_log <- "auction has been updated completely"
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func createUserMain(h *UserHashTable, report chan User, report_log chan string, uid uint64, name string) {
	if !h.SearchUserIDHashTable(uid) {
		newUser := CreateUser("username"+fmt.Sprint(uid), name, uid)
		report <- newUser
		h.InsertUserToHash(newUser)
		report_log <- "account has been created completely"
	} else {
		fmt.Println("The user has already registered into the system")
	}
}

func createAuctionMain(U *UserHashTable, A *AuctionHashTable, auction chan Auction, report_log chan string, uid uint64, aid uint64, initial uint64, step uint64) {
	user := U.AccessUserHash(uid)
	if !A.SearchAuctIDHashTable(aid) {
		newAuction := CreateAuction(*user, initial, step, aid)
		auction <- *newAuction.CreatedAuction // This line is used to notate new user created.
		report_log <- "auction has been created completely"
	} else {
		fmt.Println("The Auction has already been created.")
	}
}
