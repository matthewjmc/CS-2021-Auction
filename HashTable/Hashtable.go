package main

import (
	"fmt"
	"time"
	"sync"
)

var wg sync.WaitGroup

const ArraySize = 5

type HashTable struct {
	array [ArraySize]*linklist
}


type linklist struct {
	head *linklistNode
}

type linklistNode struct {
	key  Auction
	next *linklistNode
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

func (h *HashTable) Insert(key Auction) {
	index := hash(key)
	h.array[index].insert(key)
	
}

func (h *HashTable) Search(key Auction) bool {
	index := hash(key)
	return h.array[index].search(key)
}

func (h *HashTable) SearchAuction(key Auction) bool {
	index := hash(key)
	return h.array[index].searchAuction(key)
}

func (h *HashTable) Delete(key Auction) {
	index := hash(key)
	h.array[index].delete(key)
}

func (h *HashTable) Update(key Auction) {
	index := hash(key)
	h.array[index].update(key)
}



func Init() *HashTable {//Allocate the hash block
	result := &HashTable{}
	for i := range result.array {
		result.array[i] = &linklist{}
	}
	return result
}

func hash(key Auction) uint32 {//Find the hash index
	return key.auctionID % ArraySize
}

func (b *linklist) insert(k Auction) {//Create auction
	if  !b.search(k) {
		newNode := &linklistNode{key: k}
		newNode.next = b.head
		b.head = newNode
		//fmt.Println(k)
	} else {
		//fmt.Println(k, "already exists")
	}
}

func (b *linklist) search(k Auction) bool {//For search the auction by using auctionID
	currentNode := b.head
	AID := k.auctionID
	for currentNode != nil {
		if currentNode.key.auctionID == AID {
			return true
		}
		currentNode = currentNode.next
	}
	return false
}

func (b *linklist) searchAuction(k Auction) bool {//For checking when updated
	currentNode := b.head
	for currentNode != nil {
		if currentNode.key == k {
			return true
		}
		currentNode = currentNode.next
	}
	return false
}

func (b *linklist) update(k Auction) {//update auction
	currentNode := b.head
	AID := k.auctionID
	for currentNode != nil {
		if currentNode.key.auctionID == AID {
			currentNode.key = k
			fmt.Println(currentNode.key)
			return
		}
		currentNode = currentNode.next
	}
}

func (b *linklist) delete(k Auction) {//delete auction

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





func main(){ //Testing the corrupt then trying access the same data at same time.
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
	for i:=0 ; i<50; i++{
		auction1.auctionID = uint32(i)
		hashTable.Insert(auction1)
	}

	ch := make(chan int,2)
	s := time.Now()
	ch<-1
	ch<-2
	go test1(*hashTable,ch)
	go test2(*hashTable,ch)
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

	
	fmt.Println("ID1c",hashTable.Search(check1))
	fmt.Println("ID2c",hashTable.Search(check2))
	fmt.Println("auction1c",hashTable.SearchAuction(check1))
	fmt.Println("auction2c",hashTable.SearchAuction(check2))
	

}

func test1(h HashTable, c chan int){
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
	for j:=0 ; j<50;j++{
		auction2.auctionID = uint32(j)
		h.Update(auction2)
	}
	
}

func test2(h HashTable, c chan int){
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
	
	for j:=0 ; j<50;j++{
		auction2.auctionID = uint32(j)
		h.Update(auction2)
	}
	
}