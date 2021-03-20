package AuctionSystem

import (
	. "CS-2021-Auction/AuctionSystem"
	"fmt"
	"reflect"
	"sync"
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

//
func mainTimeline(A *AuctionHashTable, U *UserHashTable, command string) {

	var searchID uint64

	if command == "User" || command == "user" {
		report := make(chan User)
		report_log := make(chan string)
		go createUserMain(U, report, report_log) // possible user spawning algorithm could be used to pass the users into the function for an easier goroutine.
		newUser := <-report
		log := <-report_log
		fmt.Println(log, newUser)

	} else if command == "Auction" || command == "auction" {
		report_id := make(chan uint64)
		report_log := make(chan string)
		go createAuctionMain(A, report_id, report_log) // possible user spawning algorithm could be used to pass the users into the function for an easier goroutine.
		newAuction := <-report_id
		log := <-report_log
		fmt.Println(newAuction, log)
		//A.searchAuctIDHashTable(newAuction.auctionID)

	} else if command == "bid" {

		var targetedAuctionID uint64
		fmt.Println("What is your target auction ID in the system?")
		fmt.Scanln(&targetedAuctionID)

		if !A.SearchAuctIDHashTable(targetedAuctionID) {
			fmt.Println("The auction has not been found within the memory")
		} else {
			// targetAuction := createAuction(newUser, randomize(100, 10000), randomize(100, 1000), 992129) initially, used to
			report_price := make(chan uint64)
			report_log := make(chan string)
			go makeBidMain(A, report_price, report_log, searchID) // possible user spawning algorithm could be used to pass the users into the function for an easier goroutine.
			finalAuction := <-report_price
			log := <-report_log
			fmt.Println(finalAuction, log)
		}

	} else if command == "search" {
		fmt.Println(A.SearchAuctIDHashTable(searchID))
	}

}

var wg sync.WaitGroup

func makeBidMain(h *AuctionHashTable, report_price chan uint64, report_log chan string, targetid uint64) {

	count := Randomize(1, 1000000)                                                                               // for testing
	newUser := CreateUser("testUsername"+fmt.Sprint(count), "test"+fmt.Sprint(count), Randomize(100000, 999999)) // for testing

	newBid := CreateBid(newUser, Randomize(100, 9999))

	// access for auction object to be updated at the target variable.
	target := h.AccessHashAuction(targetid)
	//fmt.Println("Previous Winner:", target.currWinnerName)
	fmt.Println("this is type of target:", reflect.TypeOf(target))
	fmt.Println("this is type of *target:", reflect.TypeOf(*target))

	target.UpdateAuctionWinner(newBid)
	h.AuctionHashAccessUpdate(*target)
	//fmt.Println("Current Winner:", target.currWinnerName)
	report_price <- target.CurrMaxBid // This line is used to notate new user created.
	report_log <- "auction has been updated completely"

}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func createUserMain(h *UserHashTable, report chan User, report_log chan string) {

	count := Randomize(1, 1000000)
	newUser := CreateUser("testUsername"+fmt.Sprint(count), "test"+fmt.Sprint(count), Randomize(100000, 999999))

	h.InsertUserToHash(newUser)

	report <- newUser // This line is used to notate new user created.
	report_log <- "account has been created completely"

}

func createAuctionMain(A *AuctionHashTable, report chan uint64, report_log chan string) {

	count := Randomize(1, 1000000)
	newUser := CreateUser("testUsername"+fmt.Sprint(count), "test"+fmt.Sprint(count), Randomize(100000, 999999))
	newAuction := CreateAuction(newUser, Randomize(100, 10000), Randomize(100, 1000), 992129)

	A.InsertAuctToHash(newAuction.CreatedAuction)

	report <- newAuction.CreatedID // This line is used to notate new user created.
	report_log <- "auction has been created completely"

}
