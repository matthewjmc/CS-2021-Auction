package main

import (
	"fmt"
	"time"
	"sync"
)

var wg sync.WaitGroup

const ArraySize = 100

type HashTable struct {
	array [ArraySize]*bucket
}


type bucket struct {
	head *bucketNode
}

type bucketNode struct {
	key  Auction
	next *bucketNode
}

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

func (b *bucket) insert(k Auction) {
	if  !b.search(k) {
		newNode := &bucketNode{key: k}
		newNode.next = b.head
		b.head = newNode
		//fmt.Println(k)
	} else {
		//fmt.Println(k, "already exists")
	}
}

func (b *bucket) search(k Auction) bool {
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

func (b *bucket) searchAuction(k Auction) bool {
	currentNode := b.head
	
	for currentNode != nil {
		if currentNode.key == k {
			return true
		}
		currentNode = currentNode.next
	}
	return false
}

func (b *bucket) update(k Auction) {
	currentNode := b.head
	AID := k.auctionID
	for currentNode != nil {
		if currentNode.key.auctionID == AID {
			currentNode.key = k
			//fmt.Println(currentNode.key)
			return
		}
		currentNode = currentNode.next
	}
	
}

func (b *bucket) delete(k Auction) {

	if b.head.key.auctionID == k.auctionID {
		b.head = b.head.next
		return
	}

	previousNode := b.head
	for previousNode.next != nil {
		if previousNode.next.key.auctionID == k.auctionID {
			//delete
			previousNode.next = previousNode.next.next
			return
		}
		previousNode = previousNode.next
	}
}

// hash
func hash(key Auction) uint32 {
	return key.auctionID % ArraySize
}


func Init() *HashTable {
	result := &HashTable{}
	for i := range result.array {
		result.array[i] = &bucket{}
	}
	return result
}

func test1(c chan int) {
	hashTable := Init()
	auction1 := Auction{
		auctionID:     123,
		auctioneerID:  123,
		itemName:      "TTT",
		currWinnerID:  123,
		currMaxBid:    123,
		bidStep:       123,
		latestBidTime: time.Now(),
		startTime:     time.Now(),
		endTime:       time.Now(),
		actionCount:   0,
	}

	auction2 := Auction{
		auctionID:     123,
		auctioneerID:  222,
		itemName:      "YYY",
		currWinnerID:  222,
		currMaxBid:    222,
		bidStep:       222,
		latestBidTime: time.Now(),
		startTime:     time.Now(),
		endTime:       time.Now(),
		actionCount:   0,
	}
	for i:=0 ; i<10000; i++{
		hashTable.Insert(auction1)
		hashTable.Update(auction2)
		hashTable.Delete(auction2)
	}
	fmt.Println("done")
	wg.Done()
	
}

func test2(c chan int) {
	hashTable := Init()
	
	auction1 := Auction{
		auctionID:     130,
		auctioneerID:  123,
		itemName:      "TTT",
		currWinnerID:  123,
		currMaxBid:    123,
		bidStep:       123,
		latestBidTime: time.Now(),
		startTime:     time.Now(),
		endTime:       time.Now(),
		actionCount:   0,
	}

	auction2 := Auction{
		auctionID:     130,
		auctioneerID:  222,
		itemName:      "YYY",
		currWinnerID:  222,
		currMaxBid:    222,
		bidStep:       222,
		latestBidTime: time.Now(),
		startTime:     time.Now(),
		endTime:       time.Now(),
		actionCount:   0,
	}
	for i:=0 ; i<10000; i++{
		hashTable.Insert(auction1)
		hashTable.Update(auction2)
		hashTable.Delete(auction2)
	}
	fmt.Println("done")
	wg.Done()

}

func main(){
	ch := make(chan int)
	s := time.Now()
	wg.Add(1)
	go test1(ch)
	wg.Add(1)
	go test1(ch)
	wg.Wait()
	fmt.Println(time.Since(s))
	
}