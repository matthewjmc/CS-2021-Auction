//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Spawn mock up user accounts with some basic information.
func mockUserCreate() []User {

	var johnd = createUser("JohnD", "John Doe", 9921)
	var alanr = createUser("AlanR", "Alan Rogers", 1547)
	var stepb = createUser("StephenB", "Stephen Browns", 7812)
	var geors = createUser("GeorgeS", "George Samuels", 3443)
	mockUserArray := []User{johnd, alanr, stepb, geors}
	return mockUserArray
}

// A function used to start a simple mockup auctioning test case for the actual logic.
func simpleMockTest() {

	var userArray = mockUserCreate()

	var auctionReport = createAuction(userArray[0], randomize(100, 500), randomize(10, 100), 1)
	testAuction := auctionReport.createdAuction
	// tagun9921 creates an auction.
	// The testing auction has the initial bid of 200, bid steps of 50 and duration of 1 hour.
	// The possible first bid suppose to have at least 260
	fmt.Println("\nThe initial bidding price is", testAuction.currMaxBid, "with a bidding step of", testAuction.bidStep)
	fmt.Println("This testAuction is being hosted by", testAuction.currWinnerName, "with the Auction ID of", testAuction.auctionID)

	var bid1 = createBid(userArray[1], 600)
	testAuction.updateAuctionWinner(bid1)
	fmt.Println("As", bid1.bidderUsername, "bids with", bid1.bidPrice, ", now the current winner is", testAuction.currMaxBid, "with", testAuction.currWinnerName)

	var bid2 = createBid(userArray[2], 5000) /*randomize(100, 2000)*/
	testAuction.updateAuctionWinner(bid2n)
	fmt.Println("As", bid2.bidderUsername, "bids with", bid2.bidPrice, ", now the current winner is", testAuction.currMaxBid, "with", testAuction.currWinnerName)

	var bid3 = createBid(userArray[3], 1500)
	testAuction = updateAuctionWinner(bid3, testAuction)
	//fmt.Println("As", bid3.bidderUsername, "bids with", bid3.bidPrice, ", now the current winner is", testAuction.currMaxBid, "with", testAuction.currWinnerName)
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
	auctionArray[auctionNumberRandom] = updateAuctionWinner(createBid(userArray[bidUserRandom], randomize(1, 10000)), auctionArray[auctionNumberRandom])
}

func mockMultiBidding(userArray []User, auction Auction) Auction {

	// updateAuctionWinner() for 1 updates : auction = updateAuctionWinner(createBid(userArray[bidUserRandom], randomize(1, 10000)), auction)

	// datetime format for bidTime parameter is --- 2021-03-08 00:17:12.0959143 +0700 +07 m=+0.002279301 ---

	updatedAuction1 := make(chan Auction)
	updatedAuction2 := make(chan Auction)
	var resultAuction Auction

	go func() {
		fmt.Println("First bidding is being made")
		newBidder1 := userArray[randomize(0, len(userArray)-1)]
		updatedAuction2 <- updateAuctionWinner(createBid(newBidder1, randomize(0, 100000)), auction)
	}()
	go func() {
		fmt.Println("Second bidding is being made")
		newBidder2 := userArray[randomize(0, len(userArray)-1)]
		updatedAuction1 <- updateAuctionWinner(createBid(newBidder2, randomize(0, 100000)), auction)
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

	fmt.Println("This marks the end of mock data setup")
	mockUpEnding := time.Now()
	fmt.Println("The initial time captured before the spawning is at", mockUpStart)
	fmt.Println("The time captured after completing the spawning is at", mockUpEnding)
	result := mockMultiBidding(userArray, auctionArray[randomize(0, len(auctionArray)-1)])
	fmt.Println(result.currMaxBid)
}