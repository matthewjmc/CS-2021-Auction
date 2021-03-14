package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

///// Channel Creation
//	c := make(chan int) // value of c is a point which the channel is located.
//	fmt.Printf("type of c is %T\n", c) // %T is to provide the type

//	fmt.Printf("The datetime data type is %T\n", currentTime.Format("2006-01-02 15:04:05.000000000"))

func main() {
	//multipleUserTest()
	//testTimeFormat()
	A := auctionAllocate()
	U := userAllocate()

	for {
		mainTimeline(A, U)
	}
}

var wg sync.WaitGroup

const ArraySize = 5

// User contains a user's information for every other implementation.
type User struct {
	username  string
	fullname  string
	accountID uint64
}

// Create new users into the system
func createUser(username string, fullname string, id uint64) User {

	account := User{username: username}
	account.fullname = fullname
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

type auctionHashTable struct {
	array [ArraySize]*auctionLinkedList
}

type auctionLinkedList struct {
	head *auctionNode
}

type auctionNode struct {
	key  Auction
	next *auctionNode
}

type userHashTable struct {
	array [ArraySize]*userLinkedList
}

type userLinkedList struct {
	head *userNode
}

type userNode struct {
	key  User
	next *userNode
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
		currWinnerName: auctioneer.fullname,
		currMaxBid:     initBid,
		bidStep:        bidStep,
		latestBidTime:  time.Now().Format(time.RFC3339Nano),
		startTime:      time.Now().Format(time.RFC3339Nano),
		endTime:        time.Now().Add(duration * time.Hour).Format(time.RFC3339Nano),
		actionCount:    0,
	}
	return auction
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
func updateAuctionWinner(b Bid, a Auction) Auction {
	fmt.Println("bid time ", b.bidTime)
	fmt.Println("end time", a.endTime)
	if b.bidTime > a.endTime {
		fmt.Printf("Auction %d has already end. Bid placement is canceled\n", a.auctionID)
		return a
	}

	// The printing results are used to debug different possibilities.
	//fmt.Println("The initial winning bid is bidded by user with the name of", a.currWinnerName, "with the price of", a.currMaxBid)
	//fmt.Println("The new incoming bid is bidded by user with the name of", b.bidderName, "with the price of", b.bidPrice)
	//fmt.Println("The previous bid was made at", a.latestBidTime, "while the new bid is bidded at time", b.bidTime)
	//fmt.Println("New bid is bidded after the previous bid:", time_test)

	if (b.bidPrice > a.currMaxBid) && (b.bidPrice-a.currMaxBid) >= a.bidStep {
		a.currMaxBid = b.bidPrice
		a.currWinnerID = b.bidderID
		a.latestBidTime = b.bidTime
		a.currWinnerName = b.bidderUsername
		a.actionCount++
	}

	time.Sleep(1 * time.Millisecond)

	return a // where a is the updated auction.
}

// used to randomize integers for different test cases.
func randomize(min int, max int) uint64 {
	rand.Seed(time.Now().UnixNano())
	var check int = rand.Intn(max-min+1) + min
	//fmt.Println(check)
	random := uint64(check)

	return random
}

// Data Structure and Storage Function Declaration by Katisak in Milestone 1.
// Data Structure and Storage Function Declaration by Katisak in Milestone 1.
// Data Structure and Storage Function Declaration by Katisak in Milestone 1.
// Data Structure and Storage Function Declaration by Katisak in Milestone 1.
// Data Structure and Storage Function Declaration by Katisak in Milestone 1.
// Data Structure and Storage Function Declaration by Katisak in Milestone 1.
// Data Structure and Storage Function Declaration by Katisak in Milestone 1.

// Functions used to find the hash index of an object
func hashAuction(key Auction) uint64 {
	return key.auctionID % ArraySize
}
func hashUser(key User) uint64 {
	return key.accountID % ArraySize
}

// Functions to insert the data into a hash table and into the linked list nodes of their corresponding datatypes.

// A behavior of a hash table object used to insert an auction into a hash function to properly placed it at the correct index.
func (h *auctionHashTable) insertAuctToHash(auction Auction) {
	index := hashAuction(auction)

	h.array[index].insertAuctToLinkedList(auction)
}

// Continuation of hash function insertion to place it within a linked list as a node.
func (b *auctionLinkedList) insertAuctToLinkedList(auction Auction) {
	if !b.searchAuctIDLinkedList(auction) {
		newNode := &auctionNode{key: auction}
		newNode.next = b.head
		b.head = newNode
		fmt.Println("The auction has been inserted properly.")
		//fmt.Println(k)
	} else {
		//fmt.Println(k, "already exists")
	}
}

// A behavior of a hash table object used to insert a user into a hash function to properly placed it at the correct index.
func (h *userHashTable) insertUserToHash(user User) {
	index := hashUser(user)
	h.array[index].insertUserToLinkedList(user)
}

// Continuation of hash function insertion to place it within a linked list as a node.
func (b *userLinkedList) insertUserToLinkedList(user User) {
	if !b.searchUserIDLinkedList(user) {
		newNode := &userNode{key: user}
		newNode.next = b.head
		b.head = newNode
		//fmt.Println(k)
	} else {
		//fmt.Println(k, "already exists")
	}
}

// A behavior of a hash table object used to search of an auction object within the hash table using auction ID of each auction.
func (h *auctionHashTable) searchAuctIDHashTable(auction Auction) bool {
	index := hashAuction(auction)
	return h.array[index].searchAuctIDLinkedList(auction)
}

// Continuation of seachAuctIDHashTable() function to continue the search within the linked list at the hash index location.
func (b *auctionLinkedList) searchAuctIDLinkedList(auction Auction) bool { //For search the auction by using auctionID
	currentNode := b.head
	temp := auction.auctionID
	for currentNode != nil {
		if currentNode.key.auctionID == temp {
			fmt.Println("The auction is found in the repository., searched by ID")
			return true
		}
		currentNode = currentNode.next
	}
	return false
}

// A behavior of a hash table object used to search of an user object within the hash table using account ID of each auction.
func (h *userHashTable) searchUserIDHashTable(user User) bool {
	index := hashUser(user)
	return h.array[index].searchUserIDLinkedList(user)
}

// Continuation of seachUserIDHashTable() function to continue the search within the linked list at the hash index location.
func (b *userLinkedList) searchUserIDLinkedList(user User) bool { //For search the user by using accouintID
	currentNode := b.head
	temp := user.accountID
	for currentNode != nil {
		if currentNode.key.accountID == temp {
			return true
		}
		currentNode = currentNode.next
	}
	return false
}

// A behavior of a hash table object used to search of an auction object within the hash table using auction name of each auction.
func (h *auctionHashTable) searchAuctNameInHash(key Auction) bool {
	index := hashAuction(key)
	return h.array[index].searchAuctNameInLinkedList(key)
}

// Continuation of seachAuctNameHashTable() function to continue the search within the linked list at the hash index location.
func (b *auctionLinkedList) searchAuctNameInLinkedList(k Auction) bool { //For checking when updated
	currentNode := b.head
	for currentNode != nil {
		if currentNode.key == k {
			fmt.Println("The auction is found in the repository. searched by name")
			return true
		}
		currentNode = currentNode.next
	}
	return false
}

// A behavior of a hash table object used to delete a user within the table.
func (h *auctionHashTable) auctionHashAccessDelete(key Auction) {
	index := hashAuction(key)
	h.array[index].deleteAuctionInLinkedList(key)
}
func (b *auctionLinkedList) deleteAuctionInLinkedList(k Auction) {

	if b.head.key.auctionID == k.auctionID {
		b.head = b.head.next
		return
	}
	previousNode := b.head
	for previousNode.next != nil {
		if previousNode.next.key.auctionID == k.auctionID {
			previousNode.next = previousNode.next.next
			return
		}
		previousNode = previousNode.next
	}
}

func (h *userHashTable) userHashAccessDelete(key User) {
	index := hashUser(key)
	h.array[index].deleteUserInLinkedList(key)
}

func (b *userLinkedList) deleteUserInLinkedList(k User) {

	if b.head.key.accountID == k.accountID {
		b.head = b.head.next
		return
	}
	previousNode := b.head
	for previousNode.next != nil {
		if previousNode.next.key.accountID == k.accountID {
			previousNode.next = previousNode.next.next
			return
		}
		previousNode = previousNode.next
	}
}

func (h *auctionHashTable) auctionHashAccessUpdate(key Auction) {
	index := hashAuction(key)
	h.array[index].updateAuctionInLinkedList(key)
}

func (b *auctionLinkedList) updateAuctionInLinkedList(k Auction) { //update auction
	currentNode := b.head
	temp := k.auctionID
	for currentNode != nil {
		if currentNode.key.auctionID == temp {
			currentNode.key = k
			fmt.Println(currentNode.key)
			fmt.Println("updateAuction completed")
			return
		}
		currentNode = currentNode.next
	}
}

func (h *userHashTable) userHashAccessUpdate(key User) {
	index := hashUser(key)
	h.array[index].updateUserInLinkedList(key)
}

func (b *userLinkedList) updateUserInLinkedList(k User) { //update user
	currentNode := b.head
	temp := k.accountID
	for currentNode != nil {
		if currentNode.key.accountID == temp {
			currentNode.key = k
			fmt.Println(currentNode.key)
			return
		}
		currentNode = currentNode.next
	}
}

// Hash block allocation.
func auctionAllocate() *auctionHashTable {
	result := &auctionHashTable{}
	for i := range result.array {
		result.array[i] = &auctionLinkedList{}
	}
	return result
}

// Hash block allocation.
func userAllocate() *userHashTable {
	result := &userHashTable{}
	for i := range result.array {
		result.array[i] = &userLinkedList{}
	}
	return result
}

// MOCK UP DATA CREATIONS ARE CREATED BELOW.
// MOCK UP DATA CREATIONS ARE CREATED BELOW.
// MOCK UP DATA CREATIONS ARE CREATED BELOW.
// MOCK UP DATA CREATIONS ARE CREATED BELOW.
// MOCK UP DATA CREATIONS ARE CREATED BELOW.
// MOCK UP DATA CREATIONS ARE CREATED BELOW.
// MOCK UP DATA CREATIONS ARE CREATED BELOW.

// mockUserCreate() and simpleMockTest() are used to test a simple test case for the auction system.

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

	var testAuction = createAuction(userArray[0], randomize(100, 500), randomize(10, 100), 1)
	// tagun9921 creates an auction.
	// The testing auction has the initial bid of 200, bid steps of 50 and duration of 1 hour.
	// The possible first bid suppose to have at least 260
	fmt.Println("\nThe initial bidding price is", testAuction.currMaxBid, "with a bidding step of", testAuction.bidStep)
	fmt.Println("This testAuction is being hosted by", testAuction.currWinnerName, "with the Auction ID of", testAuction.auctionID)

	var bid1 = createBid(userArray[1], 600)
	testAuction = updateAuctionWinner(bid1, testAuction)
	fmt.Println("As", bid1.bidderUsername, "bids with", bid1.bidPrice, ", now the current winner is", testAuction.currMaxBid, "with", testAuction.currWinnerName)

	var bid2 = createBid(userArray[2], 5000) /*randomize(100, 2000)*/
	testAuction = updateAuctionWinner(bid2, testAuction)
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

//// Milestone 2
//// Milestone 2
//// Milestone 2
//// Milestone 2
//// Milestone 2
//// Milestone 2
//// Milestone 2
//// Milestone 2
//// Milestone 2
//// Milestone 2
//// Milestone 2
//// Milestone 2

func mainTimeline(A *auctionHashTable, U *userHashTable) {
	var command string

	fmt.Println("Please state your action.")
	fmt.Scanln(&command)
	testDummy := Auction{auctionID: 992129} // using for test purposes

	if command == "User" {
		report := make(chan User)
		report_log := make(chan string)
		go createUserMain(U, report, report_log) // possible user spawning algorithm could be used to pass the users into the function for an easier goroutine.
		newUser := <-report
		log := <-report_log
		fmt.Println(newUser, log)
	} else if command == "Auction" {
		report := make(chan Auction)
		report_log := make(chan string)
		go createAuctionMain(A, report, report_log) // possible user spawning algorithm could be used to pass the users into the function for an easier goroutine.
		newAuction := <-report
		log := <-report_log
		fmt.Println(newAuction, log)
	} else if command == "bid" {
		newUser := createUser("tagun9921", "tagun", 9921)
		newAuction := createAuction(newUser, randomize(100, 10000), randomize(100, 1000), 992129)
		report := make(chan Auction)
		report_log := make(chan string)
		go makeBidMain(A, report, report_log, newAuction) // possible user spawning algorithm could be used to pass the users into the function for an easier goroutine.
		finalAuction := <-report
		log := <-report_log
		fmt.Println(finalAuction, log)
	} else if command == "insert" {
		A.insertAuctToHash(testDummy)
	} else if command == "seek" {
		A.searchAuctNameInHash(testDummy)
	}
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

func createAuctionMain(h *auctionHashTable, report chan Auction, report_log chan string) {

	count := randomize(1, 1000000)
	newUser := createUser("testUsername"+fmt.Sprint(count), "test"+fmt.Sprint(count), randomize(100000, 999999))
	newAuction := createAuction(newUser, randomize(100, 10000), randomize(100, 1000), randomize(100000, 999999))

	h.insertAuctToHash(newAuction)

	report <- newAuction // This line is used to notate new user created.
	report_log <- "auction has been created completely"

}

func makeBidMain(h *auctionHashTable, report chan Auction, report_log chan string, target Auction) {

	count := randomize(1, 1000000)
	newUser := createUser("testUsername"+fmt.Sprint(count), "test"+fmt.Sprint(count), randomize(100000, 999999))
	newBid := createBid(newUser, randomize(100, 9999))

	// access for auction object to be updated at the target variable.

	newAuction := updateAuctionWinner(newBid, target)
	h.auctionHashAccessUpdate(newAuction)
	h.searchAuctIDHashTable(newAuction)

	report <- newAuction // This line is used to notate new user created.
	report_log <- "auction has been updated completely"

}
