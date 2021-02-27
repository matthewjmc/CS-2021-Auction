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
//		first_name: "Tagun",
//		last_name:  "Jivasitthikul",
//		account_id: 15235,
//		bid_time:   func() string { return time.Now().Format("2006-01-02 15:04:05.00000000000") }
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
	fmt.Println(auction.curr_max_bid, auction.curr_winner_id)

	bid2 = createBid(bid2, matty, 1000, 22)
	fmt.Println("wait3")
	auction = updateAuction(bid2, auction)
	fmt.Println(auction.curr_max_bid, auction.curr_winner_id)

	bid3 = createBid(bid3, tagun9921, 2000, 33)
	fmt.Println("wait4")
	auction = updateAuction(bid3, auction)
	fmt.Println(auction.curr_max_bid, auction.curr_winner_id)

}

type User struct {
	first_name string
	last_name  string
	account_id uint32

	/*test := User{
		first_name: "Tagun",
		last_name:  "Jivasitthikul",
		account_id: 15235,
	}
	createAuction(test)*/

}

func createUser(username User, first string, last string, id uint32) User {

	/*username = User{
		first_name: first,
		last_name:  last,
		account_id: 10,
	}*/
	/*fmt.Println("This is the account id", username.account_id)
	fmt.Println("This is the first name", username.first_name)
	fmt.Println("this is the last name : ", username.last_name)*/

	username.first_name = first
	username.last_name = last
	username.account_id = id // need some algorithm to uniquely randomize username id.

	/*fmt.Println("This is the account id", username.account_id)
	fmt.Println("This is the first name", username.first_name)
	fmt.Println("this is the last name", username.last_name)*/

	return username
}

// func createUser( username string ,

type Auction struct {
	auction_id      uint32
	auctioneer_id   uint32
	item_name       string
	curr_winner_id  uint32
	curr_max_bid    uint32
	bid_step        uint32
	latest_bid_time time.Time
	start_time      time.Time
	end_time        time.Time
	action_count    uint32
}

func createAuction(auction Auction, auctioneer User, id uint32) Auction {

	var item_name string
	fmt.Println("Enter your item name for the auction :")
	fmt.Scanln(&item_name)

	// Consideration of item_name defection
	//if reflect.TypeOf(item_name) != reflect.TypeOf("") {
	//	fmt.Println("The item name must be a string")
	//	fmt.Scanf("Please re-enter your bidding item : %s", item_name)
	//}

	var bid_step uint32
	fmt.Println("Enter your bidding step :")
	fmt.Scanln(&bid_step)

	// Consideration of bidding step defection
	//var integer32 uint32 = 4
	//if reflect.TypeOf(item_name) != reflect.TypeOf(integer32) {
	//	fmt.Println("The bidding step must be a positive integer")
	//	fmt.Scanf("Please re-enter your bidding step : %s", bid_step)
	//}

	var init_bid uint32
	fmt.Println("Enter your initial starting price :")
	fmt.Scanln(&init_bid)

	var duration time.Duration
	fmt.Println("Enter your auction duration in hours :")
	fmt.Scanln(&duration)

	auction = Auction{
		auction_id:      id,
		auctioneer_id:   auctioneer.account_id,
		item_name:       item_name,
		curr_winner_id:  auctioneer.account_id, //
		curr_max_bid:    init_bid,              //
		bid_step:        bid_step,
		latest_bid_time: time.Now(), //
		start_time:      time.Now(),
		end_time:        time.Now().Add(duration * time.Hour),
		action_count:    0,
	}
	return auction
}

type Bid struct {
	bidding_id uint32
	bidder_id  uint32
	bid_price  uint32
	bid_time   time.Time
}

func createBid(bid_action Bid, user User, price uint32, id uint32) Bid {
	bid_action = Bid{
		bidding_id: id,
		bidder_id:  user.account_id,
		bid_price:  price,
		bid_time:   time.Now(),
	}
	return bid_action
}

func updateAuction(b Bid, a Auction) Auction {

	if b.bid_price-a.curr_max_bid >= a.bid_step {

		if b.bid_price > a.curr_max_bid {

			if b.bid_time.After(a.latest_bid_time) {

				a.curr_max_bid = b.bid_price
				a.curr_winner_id = b.bidder_id
				a.latest_bid_time = b.bid_time
				a.action_count += 1

			}

		} else {
			fmt.Println("The bidded price is lesser than the current bid price")
		}

	} else {
		fmt.Println("The bidded price is lesser than the bid step")
	}

	return a
}
