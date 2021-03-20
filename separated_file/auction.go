
var wg sync.WaitGroup

const ArraySize = 1000

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

type auctionReport struct {
	createdAuction     *Auction
	created_auction_id uint64
}

func hashAuction(targetID uint64) uint64 {
	return targetID % ArraySize
}

func auctionAllocate() *auctionHashTable {
	result := &auctionHashTable{}
	for i := range result.array {
		result.array[i] = &auctionLinkedList{}
	}
	return result
}

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

func createAuctionMain(A *auctionHashTable, report chan uint64, report_log chan string) {

	count := randomize(1, 1000000)
	newUser := createUser("testUsername"+fmt.Sprint(count), "test"+fmt.Sprint(count), randomize(100000, 999999))
	newAuction := createAuction(newUser, randomize(100, 10000), randomize(100, 1000), 992129)

	A.insertAuctToHash(newAuction.createdAuction)

	report <- newAuction.created_auction_id // This line is used to notate new user created.
	report_log <- "auction has been created completely"

}

// This function will be called by createAuctionMain() which is located within the main().
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