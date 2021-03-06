package main

import (
	"fmt"
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
	var leng, tagun9921, matty User
	var auction Auction
	var bid1, bid2, bid3 Bid

	leng = createUser(leng, "Katisak", "Jiangjaturapat", 111)
	tagun9921 = createUser(tagun9921, "Tagun", "Jivasitthikul", 222)
	matty = createUser(matty, "Matthew", "McMullin", 333)

	auction = createAuction(auction, leng, 1)

	bid1 = createBid(bid1, tagun9921, 500, 11)
	fmt.Println("wait2")
	auction = updateAuction(bid1, auction)
	fmt.Println(auction.currMaxBid, auction.currWinnerID)

	bid2 = createBid(bid2, matty, 1000, 22)
	fmt.Println("wait3")
	auction = updateAuction(bid2, auction)
	fmt.Println(auction.currMaxBid, auction.currWinnerID)

	bid3 = createBid(bid3, tagun9921, 2000, 33)
	fmt.Println("wait4")
	auction = updateAuction(bid3, auction)
	fmt.Println(auction.currMaxBid, auction.currWinnerID)

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

func updateAuction(b Bid, a Auction) Auction {

	if b.bidPrice-a.currMaxBid >= a.bidStep {

		if b.bidPrice > a.currMaxBid {

			if b.bidTime.After(a.latestBidTime) {

				a.currMaxBid = b.bidPrice
				a.currWinnerID = b.bidderID
				a.latestBidTime = b.bidTime
				a.actionCount++

			}

		} else {
			fmt.Println("The bidded price is lesser than the current bid price")
		}

	} else {
		fmt.Println("The bidded price is lesser than the bid step")
	}

	return a
}

// hihihihi
