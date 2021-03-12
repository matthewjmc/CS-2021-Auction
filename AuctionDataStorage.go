package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

const ArraySize = 5

type HashTable struct {
	array [ArraySize]*linkedList
}

type linkedList struct {
	head *linkedListNode
}

type linkedListNode struct {
	key  Auction
	next *linkedListNode
}

type Auction struct {
	auctionID     uint32
	auctioneerID  uint32
	itemName      string
	currWinnerID  uint32
	currMaxBid    uint32
	bidStep       uint32
	latestBidTime uint32
	startTime     uint32
	endTime       uint32
	actionCount   uint32
}

// Function used to find the hash index of an object.,
func hash(key Auction) uint32 {
	return key.auctionID % ArraySize
}

// Functions to insert the auction data into a hash table and into the linked list nodes.

// Insert an auction into a hash function to properly placed it at the correct index.
func (h *HashTable) insertAuctToHash(auction Auction) {
	index := hash(auction)
	h.array[index].insertAuctToLinkedList(auction)
}

// Continuation of hash function insertion to place it within a linked list as a node.
func (b *linkedList) insertAuctToLinkedList(auction Auction) {
	if !b.searchAuctIDLinkedList(auction) {
		newNode := &linkedListNode{key: auction}
		newNode.next = b.head
		b.head = newNode
		//fmt.Println(k)
	} else {
		//fmt.Println(k, "already exists")
	}
}

// Search of an auction object within the hash table using auction ID of each auction.
func (h *HashTable) searchAuctIDHashTable(auction Auction) bool {
	index := hash(auction)
	return h.array[index].searchAuctIDLinkedList(auction)
}

// Continuation of seachAuctIDHashTable() function to continue the search within the linked list at the hash index location.
func (b *linkedList) searchAuctIDLinkedList(auction Auction) bool { //For search the auction by using auctionID
	currentNode := b.head
	temp := auction.auctionID
	for currentNode != nil {
		if currentNode.key.auctionID == temp {
			return true
		}
		currentNode = currentNode.next
	}
	return false
}

// Search of an auction object within the hash table using auction name of each auction.
func (h *HashTable) searchAuctNameInHash(key Auction) bool {
	index := hash(key)
	return h.array[index].searchAuctNameInLinkedList(key)
}

// Continuation of seachAuctNameHashTable() function to continue the search within the linked list at the hash index location.
func (b *linkedList) searchAuctNameInLinkedList(k Auction) bool { //For checking when updated
	currentNode := b.head
	for currentNode != nil {
		if currentNode.key == k {
			return true
		}
		currentNode = currentNode.next
	}
	return false
}

// Hashing
func (h *HashTable) hashAccessDelete(key Auction) {
	index := hash(key)
	h.array[index].deleteAuctionInLinkedList(key)
}

func (b *linkedList) deleteAuctionInLinkedList(k Auction) {

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

func (h *HashTable) hashAccessUpdate(key Auction) {
	index := hash(key)
	h.array[index].updateAuctionInLinkedList(key)
}

func (b *linkedList) updateAuctionInLinkedList(k Auction) { //update auction
	currentNode := b.head
	temp := k.auctionID
	for currentNode != nil {
		if currentNode.key.auctionID == temp {
			currentNode.key = k
			fmt.Println(currentNode.key)
			return
		}
		currentNode = currentNode.next
	}
}

// Hash block allocation.
func Init() *HashTable {
	result := &HashTable{}
	for i := range result.array {
		result.array[i] = &linkedList{}
	}
	return result
}

func main() { //Testing the corrupt then trying access the same data at same time.
	hashTable := Init()
	auction1 := Auction{
		auctionID:     0,
		auctioneerID:  123,
		itemName:      "au1",
		currWinnerID:  123,
		currMaxBid:    123,
		bidStep:       123,
		latestBidTime: 0,
		startTime:     0,
		endTime:       0,
		actionCount:   0,
	}
	for i := 0; i < 50; i++ {
		auction1.auctionID = uint32(i)
		hashTable.insertAuctToHash(auction1)
	}

	ch := make(chan int, 2)
	s := time.Now()
	ch <- 1
	ch <- 2
	go test1(*hashTable, ch)
	go test2(*hashTable, ch)
	check1 := Auction{
		auctionID:     1,
		auctioneerID:  1,
		itemName:      "au1",
		currWinnerID:  1,
		currMaxBid:    1,
		bidStep:       1,
		latestBidTime: 1,
		startTime:     1,
		endTime:       1,
		actionCount:   0,
	}
	check2 := Auction{
		auctionID:     1,
		auctioneerID:  2,
		itemName:      "au2",
		currWinnerID:  2,
		currMaxBid:    2,
		bidStep:       2,
		latestBidTime: 2,
		startTime:     2,
		endTime:       2,
		actionCount:   0,
	}
	fmt.Println(time.Since(s))

	fmt.Println("ID1c", hashTable.searchAuctIDHashTable(check1))
	fmt.Println("ID2c", hashTable.searchAuctIDHashTable(check2))
	fmt.Println("auction1c", hashTable.searchAuctNameInHash(check1))
	fmt.Println("auction2c", hashTable.searchAuctNameInHash(check2))

}

func test1(h HashTable, c chan int) {
	auction2 := Auction{
		auctionID:     0,
		auctioneerID:  1,
		itemName:      "au1",
		currWinnerID:  1,
		currMaxBid:    1,
		bidStep:       1,
		latestBidTime: 1,
		startTime:     1,
		endTime:       1,
		actionCount:   0,
	}
	for j := 0; j < 50; j++ {
		auction2.auctionID = uint32(j)
		h.hashAccessUpdate(auction2)
	}

}

func test2(h HashTable, c chan int) {
	auction2 := Auction{
		auctionID:     0,
		auctioneerID:  2,
		itemName:      "au2",
		currWinnerID:  2,
		currMaxBid:    2,
		bidStep:       2,
		latestBidTime: 2,
		startTime:     2,
		endTime:       2,
		actionCount:   0,
	}

	for j := 0; j < 50; j++ {
		auction2.auctionID = uint32(j)
		h.hashAccessUpdate(auction2)
	}

}
