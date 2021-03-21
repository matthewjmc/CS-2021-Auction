package AuctionSystem

import (
	"fmt"
	"math/rand"
	"time"
)

// Auction stores all information used to declare the auction current status.
type Auction struct {
	AuctionID      uint64
	AuctioneerID   uint64
	ItemName       string
	CurrWinnerID   uint64
	CurrWinnerName string
	CurrMaxBid     uint64
	BidStep        uint64
	LatestBidTime  string
	StartTime      string
	EndTime        string
	ActionCount    uint64
}

type AuctionHashTable struct {
	array [ArraySize]*AuctionLinkedList
}

type AuctionLinkedList struct {
	head *AuctionNode
}

type AuctionNode struct {
	key  Auction
	next *AuctionNode
}

type AuctionReport struct {
	CreatedAuction *Auction
	CreatedID      uint64
}

func HashAuction(targetID uint64) uint64 {
	return targetID % ArraySize
}

// A behavior of a hash table object used to insert an auction into a hash function to properly placed it at the correct index.
func (h *AuctionHashTable) InsertAuctToHash(auction *Auction) {
	index := HashAuction(auction.AuctionID)
	h.array[index].insertAuctToLinkedList(*auction)
}

// Continuation of hash function insertion to place it within a linked list as a node.
func (b *AuctionLinkedList) insertAuctToLinkedList(auction Auction) {
	if !b.searchAuctIDLinkedList(auction.AuctionID) {
		newNode := &AuctionNode{key: auction}
		newNode.next = b.head
		b.head = newNode
		fmt.Println("The auction has been inserted properly.")
	} else {
		fmt.Println("The created auction already exists")
	}
}

// A behavior of a hash table object used to search of an auction object within the hash table using auction ID of each auction.
func (h *AuctionHashTable) SearchAuctIDHashTable(auctionid uint64) bool {
	index := HashAuction(auctionid)
	return h.array[index].searchAuctIDLinkedList(auctionid)
}

// Continuation of seachAuctIDHashTable() function to continue the search within the linked list at the hash index location.
func (b *AuctionLinkedList) searchAuctIDLinkedList(auctionid uint64) bool { //For search the auction by using auctionID
	currentNode := b.head
	for currentNode != nil {
		if currentNode.key.AuctionID == auctionid {
			return true
		}
		currentNode = currentNode.next
	}
	return false
}

// Hash block allocation.
func AuctionAllocate() *AuctionHashTable {
	result := &AuctionHashTable{}
	for i := range result.array {
		result.array[i] = &AuctionLinkedList{}
	}
	return result
}

func (h *AuctionHashTable) AuctionHashAccessUpdate(key Auction) {
	index := HashAuction(key.AuctionID)
	h.array[index].updateAuctionInLinkedList(key)
}

func (b *AuctionLinkedList) updateAuctionInLinkedList(k Auction) { //update auction
	currentNode := b.head
	temp := k.AuctionID
	for currentNode != nil {
		if currentNode.key.AuctionID == temp {
			currentNode.key = k
			fmt.Println(currentNode.key)
			fmt.Println("updateAuction completed")
			return
		}
		currentNode = currentNode.next
	}
}
func (h *AuctionHashTable) AccessHashAuction(auctionID uint64) *Auction {

	index := HashAuction(auctionID)
	return h.array[index].accessLinkedListAuction(auctionID)
}

func (b *AuctionLinkedList) accessLinkedListAuction(auctionID uint64) *Auction { //For checking when updated
	currentNode := b.head
	for currentNode != nil {
		if currentNode.key.AuctionID == auctionID {
			//fmt.Println("The auction is being accessed")
			return &currentNode.key
		}
		currentNode = currentNode.next
	}
	return &Auction{}
}

// A behavior of a hash table object used to delete a user within the table.
func (h *AuctionHashTable) AuctionHashAccessDelete(aid uint64) bool {
	index := HashAuction(aid)
	return h.array[index].deleteAuctionInLinkedList(aid)
}

func (b *AuctionLinkedList) deleteAuctionInLinkedList(aid uint64) bool {

	if b.head.key.AuctionID == aid {
		b.head = b.head.next
		return true
	}
	previousNode := b.head
	for previousNode.next != nil {
		if previousNode.next.key.AuctionID == aid {
			previousNode.next = previousNode.next.next
			return true
		}
		previousNode = previousNode.next
	}
	return false
}

// A behavior of a hash table object used to search of an auction object within the hash table using auction name of each auction.
func (h *AuctionHashTable) SearchAuctNameInHash(key Auction) bool {
	index := HashAuction(key.AuctionID)
	return h.array[index].searchAuctNameInLinkedList(key)
}

// Continuation of seachAuctNameHashTable() function to continue the search within the linked list at the hash index location.
func (b *AuctionLinkedList) searchAuctNameInLinkedList(k Auction) bool { //For checking when updated
	currentNode := b.head
	for currentNode != nil {
		if currentNode.key == k {
			return true
		}
		currentNode = currentNode.next
	}
	//fmt.Println("There is no auction with that name in the memory.")
	return false
}

func CreateAuction(auctioneer User, initBid uint64, bidStep uint64, id uint64) AuctionReport {

	var itemName string = "testItem"
	var duration time.Duration = 1

	auction := Auction{}
	auction = Auction{
		AuctionID:      id,
		AuctioneerID:   auctioneer.AccountID,
		ItemName:       itemName,
		CurrWinnerID:   auctioneer.AccountID,
		CurrWinnerName: auctioneer.Fullname,
		CurrMaxBid:     initBid,
		BidStep:        bidStep,
		LatestBidTime:  time.Now().Format(time.RFC3339Nano),
		StartTime:      time.Now().Format(time.RFC3339Nano),
		EndTime:        time.Now().Add(duration * time.Hour).Format(time.RFC3339Nano),
		ActionCount:    0,
	}
	result := AuctionReport{
		CreatedAuction: &auction,
		CreatedID:      id,
	}
	return result
}

func (a *Auction) UpdateAuctionWinner(b Bid) string {

	//fmt.Println("bid time ", b.bidTime)
	//fmt.Println("end time", a.endTime)

	if b.bidTime > a.EndTime {
		return "The auction has already ended"
	}

	if (b.bidPrice > a.CurrMaxBid) && (b.bidPrice-a.CurrMaxBid) >= a.BidStep {
		a.CurrMaxBid = b.bidPrice
		a.CurrWinnerID = b.bidderID
		a.LatestBidTime = b.bidTime
		a.CurrWinnerName = b.bidderUsername
	}

	time.Sleep(1 * time.Millisecond)
	report := fmt.Sprint(a.CurrWinnerID) + "is now the winner of auction" + fmt.Sprint(a.AuctionID)

	return report

	// where a is the updated auction.
}

// Create bidding to be used to update the auction.
func CreateBid(user User, price uint64, actionTime string) Bid {

	id := rand.Uint64()
	bid := Bid{}
	bid = Bid{
		biddingID:      id,
		bidderID:       user.AccountID,
		bidderUsername: user.Username,
		bidPrice:       price,
		bidTime:        actionTime,
	}
	return bid
}

// Bid is a datatype used to store bid interactions containing the bidding information.
type Bid struct {
	biddingID      uint64
	bidderID       uint64
	bidderUsername string
	bidPrice       uint64
	bidTime        string
}
