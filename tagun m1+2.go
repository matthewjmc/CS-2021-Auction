type User struct {
	username  string
	fullname  string
	accountID uint64
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
	latestBidTime  string
	startTime      string
	endTime        string
	actionCount    uint64
}

// Bid is a datatype used to store bid interactions containing the bidding information.
type Bid struct {
	biddingID      uint64
	bidderID       uint64
	bidderUsername string
	bidPrice       uint64
	bidTime        string
}

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

// Create new users into the system
func createUser(username string, fullname string, id uint64) User {

	account := User{username: username}
	account.fullname = fullname
	account.accountID = id // need some algorithm to uniquely randomize username id
	// For first milestone, a counter is used to notate the number of users.
	return account
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

///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////Milestone 2/////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

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

