package main

import (
	"fmt"
)

// auction_timeline is a file related to the actual timeline for main function calling.

func main() {
	//multipleUserTest()
	//testTimeFormat()
	A := auctionAllocate()
	U := userAllocate()
	// modification of memory allocation to be dynamically allocating.

	for {
		mainTimeline(A, U)
	}
}
func mainTimeline(A *auctionHashTable, U *userHashTable) {
	var command string

	fmt.Println("Please state your action.")
	fmt.Scanln(&command)

	if command == "Create" || command == "create" {

		var createcommand string
		fmt.Println("What would you like to create?")
		fmt.Scanln(&createcommand)

		if createcommand == "User" || createcommand == "user" {
			report := make(chan User)
			report_log := make(chan string)
			go createUserMain(U, report, report_log) // possible user spawning algorithm could be used to pass the users into the function for an easier goroutine.
			// newUser := <-report
			log := <-report_log
			fmt.Println(log)

		} else if createcommand == "Auction" || createcommand == "auction" {
			report := make(chan Auction)
			report_log := make(chan string)
			go createAuctionMain(A, report, report_log) // possible user spawning algorithm could be used to pass the users into the function for an easier goroutine.
			newAuction := <-report
			log := <-report_log
			fmt.Println(newAuction, log)
			//A.searchAuctIDHashTable(newAuction.auctionID)
		}

	} else if command == "bid" {

		//newUser := createUser("tagun9921", "tagun", 9921) // for actual mock-up user, a selection for each timeline iteration must be done.

		var targetedAuctionID uint64
		fmt.Println("What is your target auction ID in the system?")
		fmt.Scanln(&targetedAuctionID)

		if !A.searchAuctIDHashTable(targetedAuctionID) {
			fmt.Println("The auction has not been found within the memory")
		} else {
			// targetAuction := createAuction(newUser, randomize(100, 10000), randomize(100, 1000), 992129) initially, used to
			report := make(chan Auction)
			report_log := make(chan string)
			go makeBidMain(A, report, report_log, 992129) // possible user spawning algorithm could be used to pass the users into the function for an easier goroutine.
			finalAuction := <-report
			log := <-report_log
			fmt.Println(finalAuction, log)
		}
	} else if command == "search" {
		A.searchAuctIDHashTable(992129)
	}
}
