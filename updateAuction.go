package main

import (
	"fmt"
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

/*

Tagun := User{
		firstName: "Tagun",
		lastName:  "Jivasitthikul",
		accountID: 1,
	}
	var bid1 Bid
	createBid(bid1, Tagun, 3000, 1)
	go fmt.Println(updateDB(Tagun))
	go fmt.Println(updateDB(bid1))

	fmt.Scanln()

*/

func main() {
	simpleMockTest()
}

type User struct {
	firstName string
	lastName  string
	accountID uint32
}

/*test := User{
	firstName: "Tagun",
	lastName:  "Jivasitthikul",
	accountID: 15235,
}
createAuction(test)*/

func createUser(first string, last string, id uint32) User {

	/*username = User{
		firstName: first,
		lastName:  last,
		accountID: 10,
	}*/
	/*fmt.Println("This is the account id", username.accountID)
	fmt.Println("This is the first name", username.firstName)
	fmt.Println("this is the last name : ", username.lastName)*/

	account := User{firstName: first}
	account.lastName = last
	account.accountID = id // need some algorithm to uniquely randomize username id.

	/*fmt.Println("This is the account id", username.accountID)
	fmt.Println("This is the first name", username.firstName)
	fmt.Println("this is the last name", username.lastName)*/

	return account
}

type Auction struct {
	auctionID     uint32
	auctioneerID  uint32
	itemName      string
	currWinnerID  uint32
	currMaxBid    uint32
	bidStep       uint32
	latestBidTime time.Time
	startTime     time.Time
	endTime       time.Time
	actionCount   uint32
}

func createAuction(auctioneer User) Auction {

	var id uint32 = 1111111111
	var itemName string = "test"
	var initBid uint32 = 200
	var bidStep uint32 = 50
	var duration time.Duration = 1
	/*
		var itemName string
		fmt.Println("Enter your item name for the auction :")
		fmt.Scanln(&itemName)

		var bidStep uint32
		fmt.Println("Enter your bidding step :")
		fmt.Scanln(&bidStep)

		var initBid uint32
		fmt.Println("Enter your initial starting price :")
		fmt.Scanln(&initBid)

		var duration time.Duration
		fmt.Println("Enter your auction duration in hours :")
		fmt.Scanln(&duration)
	*/

	auction := Auction{auctionID: id}
	auction = Auction{
		auctioneerID:  auctioneer.accountID,
		itemName:      itemName,
		currWinnerID:  auctioneer.accountID,
		currMaxBid:    initBid,
		bidStep:       bidStep,
		latestBidTime: time.Now(),
		startTime:     time.Now(),
		endTime:       time.Now().Add(duration * time.Hour),
		actionCount:   0,
	}
	return auction
}

type Bid struct {
	biddingID uint32
	bidderID  uint32
	bidPrice  uint32
	bidTime   time.Time
}

func createBid(user User, price uint32) Bid {
	var id uint32 = 2344
	bid := Bid{}
	bid = Bid{
		biddingID: id,
		bidderID:  user.accountID,
		bidPrice:  price,
		bidTime:   time.Now(),
	}
	return bid
}

func updateAuction(b Bid, a Auction) Auction {

	var start time.Time = time.Now()
	//fmt.Println("The initial time that the auction update process caused by", b.biddingID, "is being called is", time.Now())

	if b.bidTime.After(a.endTime) {
		fmt.Printf("Auction %d has already end. Bid placement is canceled\n", a.auctionID)
		return a
	}

	var difference = b.bidPrice - a.currMaxBid
	if b.bidPrice > a.currMaxBid && difference >= a.bidStep && b.bidTime.After(a.latestBidTime) {
		a.currMaxBid = b.bidPrice
		a.currWinnerID = b.bidderID
		a.latestBidTime = b.bidTime
	}

	var end time.Time = time.Now()
	//fmt.Println("The final time that the auction update process caused by", b.biddingID, "is ended is", time.Now())

	fmt.Println(end.Sub(start))

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

	fmt.Println(tagun9921)
	fmt.Println(lengeiei)
	fmt.Println(mattfatt)
	fmt.Println(luckyS)

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

	var bid1 = createBid(userArray[1], 300)
	//fmt.Println(bid1.biddingID, bid1.bidPrice, bid1.bidTime, bid1.bidderID)

	updateAuction(bid1, testAuction)
}
