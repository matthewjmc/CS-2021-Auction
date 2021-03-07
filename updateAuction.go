package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

///// Channel Creation
//	c := make(chan int) // value of c is a point which the channel is located.
//	fmt.Printf("type of c is %T\n", c) // %T is to provide the type

///// Datetime Calling
//	currentTime := time.Now()
//	fmt.Printf("The datetime data type is %T\n", currentTime.Format("2006-01-02 15:04:05.000000000"))

///// Struct is a creation of User Input Data Type.

func main() {

	simpleMockTest()

}

// User contains a user's information for every other implementation.
type User struct {
	firstName string
	lastName  string
	accountID uint64
}

func createUser(first string, last string, id uint64) User {

	account := User{firstName: first}
	account.lastName = last
	account.accountID = id // need some algorithm to uniquely randomize username id.

	return account
}

// Auction stores all information used to declare the auction current status.
type Auction struct {
	auctionID      uint64
	auctioneerID   uint64
	itemName       string
	currWinnerID   uint64
	currWinnerName string
	currMaxBid     uint64
	bidStep        uint64
	latestBidTime  time.Time
	startTime      time.Time
	endTime        time.Time
	actionCount    uint64
}

func createAuction(auctioneer User) Auction {

	//var id uint64 = 1111111111
	var id = randomize(0, 999999)
	var itemName string = "test"
	var initBid uint64 = 200
	var bidStep uint64 = 50
	var duration time.Duration = 1

	/*
		var itemName string
		fmt.Println("Enter your item name for the auction :")
		fmt.Scanln(&itemName)

		var bidStep uint64
		fmt.Println("Enter your bidding step :")
		fmt.Scanln(&bidStep)

		var initBid uint64
		fmt.Println("Enter your initial starting price :")
		fmt.Scanln(&initBid)

		var duration time.Duration
		fmt.Println("Enter your auction duration in hours :")
		fmt.Scanln(&duration)
	*/

	auction := Auction{}
	auction = Auction{
		auctionID:      id,
		auctioneerID:   auctioneer.accountID,
		itemName:       itemName,
		currWinnerID:   auctioneer.accountID,
		currWinnerName: auctioneer.firstName,
		currMaxBid:     initBid,
		bidStep:        bidStep,
		latestBidTime:  time.Now(),
		startTime:      time.Now(),
		endTime:        time.Now().Add(duration * time.Hour),
		actionCount:    0,
	}

	return auction
}

// Bid is a datatype used to store bid interactions containing the bidding information.
type Bid struct {
	biddingID  uint64
	bidderID   uint64
	bidderName string
	bidPrice   uint64
	bidTime    time.Time
}

func createBid(user User, price uint64) Bid {
	id := rand.Uint64()
	bid := Bid{}
	bid = Bid{
		biddingID:  id,
		bidderID:   user.accountID,
		bidderName: user.firstName,
		bidPrice:   price,
		bidTime:    time.Now(),
	}
	return bid
}

func updateAuction(b Bid, a Auction) Auction {

	//var start time.Time = time.Now()
	//fmt.Println("The initial time that the auction update process caused by", b.biddingID, "is being called is", time.Now())

	if b.bidTime.After(a.endTime) {
		fmt.Printf("Auction %d has already end. Bid placement is canceled\n", a.auctionID)
		return a
	}

	if (b.bidPrice > a.currMaxBid) && (b.bidPrice-a.currMaxBid) >= a.bidStep {
		a.currMaxBid = b.bidPrice
		a.currWinnerID = b.bidderID
		a.latestBidTime = b.bidTime
		a.currWinnerName = b.bidderName
		a.actionCount++
	}

	//var end time.Time = time.Now()
	//fmt.Println("The final time that the auction update process caused by", b.biddingID, "is ended is", time.Now())

	//fmt.Println(end.Sub(start))
	time.Sleep(1 * time.Millisecond)
	return a
}

func updateDB(x interface{}) string {
	// get the input items to be transferred to the database
	return reflect.TypeOf(x).String()
}

func displayAction(x interface{}) string {
	// get the input items to be transferred through TCP sockets
	return reflect.TypeOf(x).String()
}

func mockUserCreate() []User {

	var tagun9921 = createUser("Tagun", "Jivasitthikul", 9921)
	var lengeiei = createUser("Katisak", "Jiangjaturapat", 1547)
	var mattfatt = createUser("Matthew", "McMullin", 7812)
	var luckyS = createUser("Vorachat", "Somsuay", 3443)
	mockUserArray := []User{tagun9921, lengeiei, mattfatt, luckyS}
	return mockUserArray
}

func simpleMockTest() {

	var userArray = mockUserCreate()
	//fmt.Println("\n", userArray)
	//fmt.Println(userArray[0])

	var testAuction = createAuction(userArray[0]) // tagun9921 creates an auction.
	// The testing auction has the initial bid of 200, bid steps of 50 and duration of 1 hour.
	// The possible first bid suppose to have at least 260
	fmt.Println("\nThe initial bidding price is", testAuction.currMaxBid, "with a bidding step of", testAuction.bidStep)
	fmt.Println("This testAuction is being hosted by", testAuction.currWinnerName, "with the Auction ID of", testAuction.auctionID)

	var bid1 = createBid(userArray[1], randomize(100, 2000)) // lengeiei , with the condition of winning the auction
	testAuction = updateAuction(bid1, testAuction)
	fmt.Println("As", bid1.bidderName, "bids with", bid1.bidPrice, ", now the current winner is", testAuction.currMaxBid, "with", testAuction.currWinnerName)

	var bid2 = createBid(userArray[2], randomize(100, 2000)) // matt , with the condition of bidding lesser than the step.
	testAuction = updateAuction(bid2, testAuction)

	fmt.Println("As", bid2.bidderName, "bids with", bid2.bidPrice, ", now the current winner is", testAuction.currMaxBid, "with", testAuction.currWinnerName)

	var bid3 = createBid(userArray[3], randomize(100, 2000)) // lengeiei , with the condition of winning the auction
	testAuction = updateAuction(bid3, testAuction)

	fmt.Println("As", bid3.bidderName, "bids with", bid3.bidPrice, ", now the current winner is", testAuction.currMaxBid, "with", testAuction.currWinnerName)

	//winnerResult <- testAuction.currWinnerName

}

func randomize(min int, max int) uint64 {
	rand.Seed(time.Now().UnixNano())
	var check int = rand.Intn(max-min+1) + min
	//fmt.Println(check)
	random := uint64(check)

	return random
}
