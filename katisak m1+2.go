

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