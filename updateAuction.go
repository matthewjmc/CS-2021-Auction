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

///// Struct Declaration
//	test := User{
//		firstName: "Tagun",
//		lastName:  "Jivasitthikul",
//		accountID: 15235,
//		bidTime:   func() string { return time.Now().Format("2006-01-02 15:04:05.00000000000") }
//	}

///// Struct is a creation of User Input Data Type.

func main() {
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

func createUser(username User, first string, last string, id uint32) User {

	/*username = User{
		firstName: first,
		lastName:  last,
		accountID: 10,
	}*/
	/*fmt.Println("This is the account id", username.accountID)
	fmt.Println("This is the first name", username.firstName)
	fmt.Println("this is the last name : ", username.lastName)*/

	username.firstName = first
	username.lastName = last
	username.accountID = id // need some algorithm to uniquely randomize username id.

	/*fmt.Println("This is the account id", username.accountID)
	fmt.Println("This is the first name", username.firstName)
	fmt.Println("this is the last name", username.lastName)*/

	return username
}

// func createUser( username string ,

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

func createAuction(auction Auction, auctioneer User, id uint32) Auction {

	var itemName string
	fmt.Println("Enter your item name for the auction :")
	fmt.Scanln(&itemName)

	// Consideration of itemName defection
	//if reflect.TypeOf(itemName) != reflect.TypeOf("") {
	//	fmt.Println("The item name must be a string")
	//	fmt.Scanf("Please re-enter your bidding item : %s", itemNam)
	//}

	var bidStep uint32
	fmt.Println("Enter your bidding step :")
	fmt.Scanln(&bidStep)

	// Consideration of bidding step defection
	//var integer32 uint32 = 4
	//if reflect.TypeOf(itemName) != reflect.TypeOf(integer32) {
	//	fmt.Println("The bidding step must be a positive integer")
	//	fmt.Scanf("Please re-enter your bidding step : %s", bidStep)
	//}

	var initBid uint32
	fmt.Println("Enter your initial starting price :")
	fmt.Scanln(&initBid)

	var duration time.Duration
	fmt.Println("Enter your auction duration in hours :")
	fmt.Scanln(&duration)

	auction = Auction{
		auctionID:     id,
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

func createBid(bidAction Bid, user User, price uint32, id uint32) Bid {
	bidAction = Bid{
		biddingID: id,
		bidderID:  user.accountID,
		bidPrice:  price,
		bidTime:   time.Now(),
	}
	return bidAction
}

/*func updateAuction(b Bid, a Auction) Auction {
	var difference = b.bidPrice - a.currMaxBid
	if b.bidPrice > a.currMaxBid && difference >= a.bidStep && b.bidTime.After(a.latestBidTime) {

	}

}*/

// hihihihi

func updateDB(x interface{}) string {
	// get the input items to be transferred to the database
	return reflect.TypeOf(x).String()
}

func displayAction(x interface{}) string {
	// get the input items to be transferred to the database
	return reflect.TypeOf(x).String()
}
