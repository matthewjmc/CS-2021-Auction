package main

import (
	"fmt"
	"math/rand"
	"time"
)

///// Channel Creation
//	c := make(chan int) // value of c is a point which the channel is located.
//	fmt.Printf("type of c is %T\n", c) // %T is to provide the type

//	fmt.Printf("The datetime data type is %T\n", currentTime.Format("2006-01-02 15:04:05.000000000"))

func main() {
	//multipleUserTest()
	simpleMockTest()
}

// User contains a user's information for every other implementation.
type User struct {
	firstName string
	lastName  string
	accountID uint64
}

// Create new users into the system
func createUser(first string, last string, id uint64) User {

	account := User{firstName: first}
	account.lastName = last
	account.accountID = id // need some algorithm to uniquely randomize username id
	// For first milestone, a counter is used to notate the number of users.
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

// Create new auction struct into the system
func createAuction(auctioneer User, initBid uint64, bidStep uint64, id uint64) Auction {

	var itemName string = "testItem"
	var duration time.Duration = 1

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

// Create bidding to be used to update the auction.
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

// Get the bidding processes created and compare it with the current auction.
func updateAuction(b Bid, a Auction) Auction {

	if b.bidTime.After(a.endTime) {
		fmt.Printf("Auction %d has already end. Bid placement is canceled\n", a.auctionID)
		return a
	}

	// The printing results are used to debug different possibilities.
	fmt.Println("The initial winning bid is bidded by user with the name of", a.currWinnerName, "with the price of", a.currMaxBid)
	fmt.Println("The new incoming bid is bidded by user with the name of", b.bidderName, "with the price of", b.bidPrice)
	fmt.Println("The previous bid was made at", a.latestBidTime, "while the new bid is bidded at time", b.bidTime)
	fmt.Println("New bid is bidded after the previous bid:", b.bidTime.After(a.latestBidTime))

	if (b.bidPrice > a.currMaxBid) && (b.bidPrice-a.currMaxBid) >= a.bidStep {
		a.currMaxBid = b.bidPrice
		a.currWinnerID = b.bidderID
		a.latestBidTime = b.bidTime
		a.currWinnerName = b.bidderName
		a.actionCount++
		//fmt.Println("Now,", a.currWinnerName, "is winning auction with the ID", a.auctionID)
		//fmt.Println("Winning bid price:", a.currMaxBid)
		//fmt.Println("Winning bidder :", a.currWinnerName)
	}

	time.Sleep(1 * time.Millisecond)

	return a // where a is the updated auction.
}

/*
func updateDB(x interface{}) string {
	// get the input items to be transferred to the database
	return reflect.TypeOf(x).String()
}

func displayAction(x interface{}) string {
	// get the input items to be transferred through TCP sockets
	return reflect.TypeOf(x).String()
}
*/

// used to randomize integers for different test cases.
func randomize(min int, max int) uint64 {
	rand.Seed(time.Now().UnixNano())
	var check int = rand.Intn(max-min+1) + min
	//fmt.Println(check)
	random := uint64(check)

	return random
}

// mockUserCreate() and simpleMockTest() are used to test a simple test case for the auction system.

// Spawn mock up user accounts with some basic information.
func mockUserCreate() []User {

	var johnd = createUser("John", "Doe", 9921)
	var alanr = createUser("Alan", "Rogers", 1547)
	var stepb = createUser("Stephen", "Browns", 7812)
	var geors = createUser("George", "Samuels", 3443)
	mockUserArray := []User{johnd, alanr, stepb, geors}
	return mockUserArray
}

// A function used to start a simple mockup auctioning test case for the actual logic.
func simpleMockTest() {

	var userArray = mockUserCreate()

	var testAuction = createAuction(userArray[0], randomize(100, 500), randomize(10, 100), 1)
	// tagun9921 creates an auction.
	// The testing auction has the initial bid of 200, bid steps of 50 and duration of 1 hour.
	// The possible first bid suppose to have at least 260
	fmt.Println("\nThe initial bidding price is", testAuction.currMaxBid, "with a bidding step of", testAuction.bidStep)
	fmt.Println("This testAuction is being hosted by", testAuction.currWinnerName, "with the Auction ID of", testAuction.auctionID)

	var bid1 = createBid(userArray[1], 600)
	testAuction = updateAuction(bid1, testAuction)
	fmt.Println("As", bid1.bidderName, "bids with", bid1.bidPrice, ", now the current winner is", testAuction.currMaxBid, "with", testAuction.currWinnerName)

	var bid2 = createBid(userArray[2], 5000) /*randomize(100, 2000)*/
	testAuction = updateAuction(bid2, testAuction)
	fmt.Println("As", bid2.bidderName, "bids with", bid2.bidPrice, ", now the current winner is", testAuction.currMaxBid, "with", testAuction.currWinnerName)

	var bid3 = createBid(userArray[3], 1500)
	testAuction = updateAuction(bid3, testAuction)
	//fmt.Println("As", bid3.bidderName, "bids with", bid3.bidPrice, ", now the current winner is", testAuction.currMaxBid, "with", testAuction.currWinnerName)
	//fmt.Println(testAuction.startTime)

	fmt.Println(testAuction.currWinnerName)
}

// The functions below are used for testing for multiple users handling.
// creation of multiple users, with those incremented counts being used to provide a unique identification.
func multiUserCreate() []User {

	var count uint64 = 0
	mockUserArray := []User{}

	for count = 0; count <= 100000; count++ {
		mockUserArray = append(mockUserArray, createUser("username"+fmt.Sprint(count), "b", count))
	}
	return mockUserArray
}

func multiAuctionCreate(userArray []User) []Auction {

	var count uint64 = 0
	mockAuctArray := []Auction{}

	for count = 0; count <= 100; count++ {
		time.Sleep(1 * time.Millisecond)
		mockAuctArray = append(mockAuctArray, createAuction(userArray[randomize(0, len(userArray))], randomize(100, 200), randomize(1, 100), count))
	}
	return mockAuctArray
}

func mockBidding(userArray []User, auctionArray []Auction) {
	numberUser := len(userArray) - 1
	numberAuction := len(auctionArray) - 1
	bidUserRandom := randomize(0, numberUser)
	auctionNumberRandom := randomize(0, numberAuction)
	auctionArray[auctionNumberRandom] = updateAuction(createBid(userArray[bidUserRandom], randomize(1, 10000)), auctionArray[auctionNumberRandom])
}

func mockMultiBidding(userArray []User, auction Auction) Auction {

	// updateAuction() for 1 updates : auction = updateAuction(createBid(userArray[bidUserRandom], randomize(1, 10000)), auction)

	// datetime format for bidTime parameter is --- 2021-03-08 00:17:12.0959143 +0700 +07 m=+0.002279301 ---

	updatedAuction1 := make(chan Auction)
	updatedAuction2 := make(chan Auction)
	var resultAuction Auction

	go func() {
		fmt.Println("First bidding is being made")
		newBidder1 := userArray[randomize(0, len(userArray)-1)]
		updatedAuction2 <- updateAuction(createBid(newBidder1, randomize(0, 100000)), auction)
	}()
	go func() {
		fmt.Println("Second bidding is being made")
		newBidder2 := userArray[randomize(0, len(userArray)-1)]
		updatedAuction1 <- updateAuction(createBid(newBidder2, randomize(0, 100000)), auction)
	}()

	updatedResult1, updatedResult2 := <-updatedAuction1, <-updatedAuction2
	time.Sleep(1 * time.Millisecond)

	/*fmt.Println(updatedResult1)
	fmt.Println(updatedResult2)*/

	if updatedResult1.currMaxBid > updatedResult2.currMaxBid {
		resultAuction = updatedResult1
	} else {
		resultAuction = updatedResult2
	}

	fmt.Println(resultAuction)
	return resultAuction
}

func multipleUserTest() {
	fmt.Println("\nThis line marks the creation of user mock data creation.")

	mockUpStart := time.Now()
	userArray := multiUserCreate()
	auctionArray := multiAuctionCreate(userArray)

	/*
		fmt.Println("\nThe user array with 10000 users is listed below.")
		fmt.Println(userArray)
		fmt.Println("\nThe auction array with 100 auctions is listed below.")
		for i := 1; i < len(auctionArray); i++ {
			fmt.Println(auctionArray[i])
		}
	*/

	fmt.Println("This marks the end of mock data setup\n")
	mockUpEnding := time.Now()
	fmt.Println("The initial time captured before the spawning is at", mockUpStart)
	fmt.Println("The time captured after completing the spawning is at", mockUpEnding)
	result := mockMultiBidding(userArray, auctionArray[randomize(0, len(auctionArray)-1)])
	fmt.Println(result.currMaxBid)
}

/*
func main() {

    messages := make(chan string)

    go func() { messages <- "ping" }()

    msg := <-messages
    fmt.Println(msg)
}
*/

// youtube video link for milestone 2 preparation...
// Goroutine Synchronization : Golang Practical Programming.
