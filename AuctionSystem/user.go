package AuctionSystem

import (
	"fmt"
	"math/rand"
	"time"
)

const ArraySize = 1000

// User contains a user's information for every other implementation.
type User struct {
	AccountID uint64
	Username  string
	Fullname  string
}

type UserHashTable struct {
	Array [ArraySize]*UserLinkedList
}

type UserLinkedList struct {
	Head *UserNode
}

type UserNode struct {
	Key  User
	Next *UserNode
}

func HashUser(uid uint64) uint64 {
	return uid % ArraySize
}

// A behavior of a hash table object used to insert a user into a hash function to properly placed it at the correct index.
func (h *UserHashTable) InsertUserToHash(user User) bool {
	index := HashUser(user.AccountID)
	return h.Array[index].insertUserToLinkedList(user)
}

// Continuation of hash function insertion to place it within a linked list as a node.
func (b *UserLinkedList) insertUserToLinkedList(User User) bool {
	if !b.searchUserIDLinkedList(User.AccountID) {
		newNode := &UserNode{Key: User}
		newNode.Next = b.Head
		b.Head = newNode
		//fmt.Println(k)
		return true
	} else {
		return false
	}
}

// A behavior of a hash table object used to search of an user object within the hash table using account ID of each auction.
func (h *UserHashTable) SearchUserIDHashTable(uid uint64) bool {
	index := HashUser(uid)
	return h.Array[index].searchUserIDLinkedList(uid)
}

// Continuation of seachUserIDHashTable() function to continue the search within the linked list at the hash index location.
func (b *UserLinkedList) searchUserIDLinkedList(uid uint64) bool { //For search the user by using accouintID
	currentNode := b.Head
	for currentNode != nil {
		if currentNode.Key.AccountID == uid {
			return true
		}
		currentNode = currentNode.Next
	}
	//fmt.Println("There is no auction with that ID in the memory.")
	return false
}

// A behavior of a hash table object used to search of an user object within the hash table using account ID of each auction.
func (h *UserHashTable) AccessUserHash(uid uint64) *User {
	index := HashUser(uid)
	return h.Array[index].accessUserLinkedList(uid)
}

// Continuation of seachUserIDHashTable() function to continue the search within the linked list at the hash index location.
func (b *UserLinkedList) accessUserLinkedList(uid uint64) *User { //For search the user by using accouintID
	currentNode := b.Head
	for currentNode != nil {
		if currentNode.Key.AccountID == uid {
			return &currentNode.Key
		}
		currentNode = currentNode.Next
	}
	//fmt.Println("There is no user with that ID in the memory.")
	return &User{}
}

func (h *UserHashTable) UserHashAccessDelete(uid uint64) bool {
	index := HashUser(uid)
	return h.Array[index].DeleteUserInLinkedList(uid)
}

func (b *UserLinkedList) DeleteUserInLinkedList(uid uint64) bool {

	if b.Head.Key.AccountID == uid {
		b.Head = b.Head.Next
		return true
	}
	previousNode := b.Head
	for previousNode.Next != nil {
		if previousNode.Next.Key.AccountID == uid {
			previousNode.Next = previousNode.Next.Next
			return true
		}
		previousNode = previousNode.Next
	}
	return false
}

func (h *UserHashTable) UserHashAccessUpdate(Key User) {
	index := HashUser(Key.AccountID)
	h.Array[index].UpdateUserInLinkedList(Key)
}

func (b *UserLinkedList) UpdateUserInLinkedList(k User) { //update user
	currentNode := b.Head
	temp := k.AccountID
	for currentNode != nil {
		if currentNode.Key.AccountID == temp {
			currentNode.Key = k
			fmt.Println(currentNode.Key)
			return
		}
		currentNode = currentNode.Next
	}
}

// Hash block allocation.
func UserAllocate() *UserHashTable {
	result := &UserHashTable{}
	for i := range result.Array {
		result.Array[i] = &UserLinkedList{}
	}
	return result
}

// Create new users into the system
func CreateUser(Username string, Fullname string, id uint64) User {

	account := User{Username: Username}
	account.Fullname = Fullname
	account.AccountID = id // need some algorithm to uniquely randomize username id
	// For first milestone, a counter is used to notate the number of users.
	return account
}

// Functions to insert the data into a hash table and into the linked list nodes of their corresponding datatypes.

// used to randomize integers for different test cases.
func Randomize(min int, max int) uint64 {
	rand.Seed(time.Now().UnixNano())
	var check int = rand.Intn(max-min+1) + min
	//fmt.Println(check)
	random := uint64(check)

	return random
}
