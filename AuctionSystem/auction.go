package AuctionSystem

import (
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
func (h *AuctionHashTable) InsertAuctToHash(auction *Auction) bool {
	index := HashAuction(auction.AuctionID)
	return h.array[index].insertAuctToLinkedList(*auction)
}

// Continuation of hash function insertion to place it within a linked list as a node.
func (b *AuctionLinkedList) insertAuctToLinkedList(auction Auction) bool {
	if !b.searchAuctIDLinkedList(auction.AuctionID) {
		newNode := &AuctionNode{key: auction}
		newNode.next = b.head
		b.head = newNode
		return true
	} else {
		return false
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
	//fmt.Println("There is no auction with that ID in the memory.")
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
	return false
}

func CreateAuction(auctioneer User, initBid uint64, bidStep uint64, id uint64, duration time.Duration, itemName string) Auction {
	auction := Auction{
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
	}
	return auction
}

func (a *Auction) UpdateAuctionWinner(b Bid) bool {
	bidtime, err := time.Parse(time.RFC3339Nano, b.BidTime)
	endtime, err2 := time.Parse(time.RFC3339Nano, a.EndTime)
	if err != nil || err2 != nil {
		return false
	} else {
		if bidtime.After(endtime) {
			return false
		}
		if b.BidPrice > a.CurrMaxBid {
			if (b.BidPrice - a.CurrMaxBid) >= a.BidStep {
				a.CurrMaxBid = b.BidPrice
				a.CurrWinnerID = b.BidderID
				a.LatestBidTime = b.BidTime
				a.CurrWinnerName = b.BidderUsername
			}
		}
		return true
	}
}

// Create bidding to be used to update the auction.
func CreateBid(user User, price uint64, actionTime string) Bid {
	bid := Bid{}
	bid = Bid{
		BidderID:       user.AccountID,
		BidderUsername: user.Username,
		BidPrice:       price,
		BidTime:        actionTime,
	}
	return bid
}

// Bid is a datatype used to store bid interactions containing the bidding information.
type Bid struct {
	BidderID       uint64
	BidderUsername string
	BidPrice       uint64
	BidTime        string
}
