package main

import (
	"testing"
	"time"
)

func Testinsert(t *testing.T) {
	hashTable := Init()
	auction1 := Auction{
		auctionID:     1,
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
		auctionID:     1,
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

	hashTable.Insert(auction1)
	hashTable.Update(auction2)
	hashTable.Delete(auction2)

}
