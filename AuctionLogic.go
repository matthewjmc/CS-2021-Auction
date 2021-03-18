package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// AuctionLogic is a file composed of Tagun's and Katisak's code altogether.
// The file is also being separated into 3 files which are auction_method, auction_timeline, data_structure.

func main() {
	//multipleUserTest()
	//testTimeFormat()
	A := auctionAllocate()
	U := userAllocate()
	// modification of memory allocation to be dynamically allocating.

	for {
		mainTimeline(A, U)
	}
}

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

var wg sync.WaitGroup

const ArraySize = 1000

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

type auctionReport struct {
	createdAuction     *Auction
	created_auction_id uint64
}

func hashAuction(targetID uint64) uint64 {
	return targetID % ArraySize
}
func hashUser(key User) uint64 {
	return key.accountID % ArraySize
}

// Functions to insert the data into a hash table and into the linked list nodes of their corresponding datatypes.

// A behavior of a hash table object used to insert an auction into a hash function to properly placed it at the correct index.
func (h *auctionHashTable) insertAuctToHash(auction *Auction) {
	index := hashAuction(auction.auctionID)
	h.array[index].insertAuctToLinkedList(*auction)
}

// Continuation of hash function insertion to place it within a linked list as a node.
func (b *auctionLinkedList) insertAuctToLinkedList(auction Auction) {
	if !b.searchAuctIDLinkedList(auction.auctionID) {
		newNode := &auctionNode{key: auction}
		newNode.next = b.head
		b.head = newNode
		fmt.Println("The auction has been inserted properly.")
	} else {
		fmt.Println("The created auction already exists")
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
func (h *auctionHashTable) searchAuctIDHashTable(auctionid uint64) bool {
	index := hashAuction(auctionid)
	return h.array[index].searchAuctIDLinkedList(auctionid)
}

// Continuation of seachAuctIDHashTable() function to continue the search within the linked list at the hash index location.
func (b *auctionLinkedList) searchAuctIDLinkedList(auctionid uint64) bool { //For search the auction by using auctionID
	currentNode := b.head
	for currentNode != nil {
		if currentNode.key.auctionID == auctionid {
			return true
		}
		currentNode = currentNode.next
	}
	fmt.Println("There is no function with that ID in the memory.")
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
	fmt.Println("There is no auction with that ID in the memory.")
	return false
}

// A behavior of a hash table object used to search of an auction object within the hash table using auction name of each auction.
func (h *auctionHashTable) searchAuctNameInHash(key Auction) bool {
	index := hashAuction(key.auctionID)
	return h.array[index].searchAuctNameInLinkedList(key)
}

// Continuation of seachAuctNameHashTable() function to continue the search within the linked list at the hash index location.
func (b *auctionLinkedList) searchAuctNameInLinkedList(k Auction) bool { //For checking when updated
	currentNode := b.head
	for currentNode != nil {
		if currentNode.key == k {
			return true
		}
		currentNode = currentNode.next
	}
	fmt.Println("There is no auction with that name in the memory.")
	return false
}

func (h *auctionHashTable) accessHashAuction(auctionID uint64) *Auction {

	index := hashAuction(auctionID)
	return h.array[index].accessLinkedListAuction(auctionID)
}

func (b *auctionLinkedList) accessLinkedListAuction(auctionID uint64) *Auction { //For checking when updated
	currentNode := b.head
	for currentNode != nil {
		if currentNode.key.auctionID == auctionID {
			fmt.Println("The auction is being accessed")
			return &currentNode.key
		}
		currentNode = currentNode.next
	}
	return &Auction{}
}

// A behavior of a hash table object used to delete a user within the table.
func (h *auctionHashTable) auctionHashAccessDelete(key Auction) {
	index := hashAuction(key.auctionID)
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
	index := hashAuction(key.auctionID)
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

// used to randomize integers for different test cases.
func randomize(min int, max int) uint64 {
	rand.Seed(time.Now().UnixNano())
	var check int = rand.Intn(max-min+1) + min
	//fmt.Println(check)
	random := uint64(check)

	return random
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
