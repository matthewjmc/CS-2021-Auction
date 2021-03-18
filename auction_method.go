package main

import (
	"fmt"
	"math/rand"
	"time"
)

// This file composed of the methods used within the main timeline function to do different tasks.

func createAuction(auctioneer User, initBid uint64, bidStep uint64, id uint64) *auctionReport {

	var itemName string = "testItem"
	var duration time.Duration = 1

	auction := Auction{}
	auction = Auction{
		auctionID:      id,
		auctioneerID:   auctioneer.accountID,
		itemName:       itemName,
		currWinnerID:   auctioneer.accountID,
		currWinnerName: auctioneer.fullname,
		currMaxBid:     initBid,
		bidStep:        bidStep,
		latestBidTime:  time.Now().Format(time.RFC3339Nano),
		startTime:      time.Now().Format(time.RFC3339Nano),
		endTime:        time.Now().Add(duration * time.Hour).Format(time.RFC3339Nano),
		actionCount:    0,
	}
	result := auctionReport{
		createdAuction:     &auction,
		created_auction_id: id,
	}
	return &result
}

// Create bidding to be used to update the auction.
func createBid(user User, price uint64) Bid {
	id := rand.Uint64()
	bid := Bid{}
	bid = Bid{
		biddingID:      id,
		bidderID:       user.accountID,
		bidderUsername: user.username,
		bidPrice:       price,
		bidTime:        time.Now().Format(time.RFC3339Nano),
	}
	return bid
}

// Get the bidding processes created and compare it with the current auction.
func (a *Auction) updateAuctionWinner(b Bid) string {

	//fmt.Println("bid time ", b.bidTime)
	//fmt.Println("end time", a.endTime)

	if b.bidTime > a.endTime {
		return "The auction has already ended"
	}

	if (b.bidPrice > a.currMaxBid) && (b.bidPrice-a.currMaxBid) >= a.bidStep {
		a.currMaxBid = b.bidPrice
		a.currWinnerID = b.bidderID
		a.latestBidTime = b.bidTime
		a.currWinnerName = b.bidderUsername
	}

	time.Sleep(1 * time.Millisecond)
	report := fmt.Sprint(a.currWinnerID) + "is now the winner of auction" + fmt.Sprint(a.auctionID)

	return report

	// where a is the updated auction.
}

// used to randomize integers for different test cases.
func randomize(min int, max int) uint64 {
	rand.Seed(time.Now().UnixNano())
	var check int = rand.Intn(max-min+1) + min
	//fmt.Println(check)
	random := uint64(check)

	return random
}

//	c := make(chan int) // value of c is a point which the channel is located.
//	fmt.Printf("type of c is %T\n", c) // %T is to provide the type

func createUserMain(h *userHashTable, report chan User, report_log chan string) {

	count := randomize(1, 1000000)
	newUser := createUser("testUsername"+fmt.Sprint(count), "test"+fmt.Sprint(count), randomize(100000, 999999))

	h.insertUserToHash(newUser)
	report <- newUser // This line is used to notate new user created.
	report_log <- "account has been created completely"

}

func createAuctionMain(A *auctionHashTable, report chan Auction, report_log chan string) {

	count := randomize(1, 1000000)
	newUser := createUser("testUsername"+fmt.Sprint(count), "test"+fmt.Sprint(count), randomize(100000, 999999))
	newAuction := createAuction(newUser, randomize(100, 10000), randomize(100, 1000), 992129)

	A.insertAuctToHash(newAuction.createdAuction)

	report <- *newAuction.createdAuction // This line is used to notate new user created.
	report_log <- "auction has been created completely"

}

func makeBidMain(h *auctionHashTable, report chan Auction, report_log chan string, targetid uint64) {

	count := randomize(1, 1000000)                                                                               // for testing
	newUser := createUser("testUsername"+fmt.Sprint(count), "test"+fmt.Sprint(count), randomize(100000, 999999)) // for testing

	newBid := createBid(newUser, randomize(100, 9999))

	// access for auction object to be updated at the target variable.
	target := h.accessHashAuction(targetid)
	fmt.Println("Previous Winner:", target.currWinnerName)
	target.updateAuctionWinner(newBid)
	h.auctionHashAccessUpdate(*target)
	fmt.Println("Current Winner:", target.currWinnerName)
	report <- *target // This line is used to notate new user created.
	report_log <- "auction has been updated completely"

}
