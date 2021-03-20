

type User struct {
	username  string
	fullname  string
	accountID uint64
}

func createUser(username string, fullname string, id uint64) User {

	account := User{username: username}
	account.fullname = fullname
	account.accountID = id // need some algorithm to uniquely randomize username id
	// For first milestone, a counter is used to notate the number of users.
	return account
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

func hashUser(key User) uint64 {
	return key.accountID % ArraySize
}

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
