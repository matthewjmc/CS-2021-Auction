package AuctionSystem

import (
	"fmt"
	"math/rand"
	"time"
)

const ArraySize = 1000

// User contains a user's information for every other implementation.
type User struct {
	Username  string
	Fullname  string
	AccountID uint64
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

func HashUser(Key User) uint64 {
	return Key.AccountID % ArraySize
}

// A behavior of a hash table object used to insert a user into a hash function to properly placed it at the correct index.
func (h *UserHashTable) InsertUserToHash(user User) {
	index := HashUser(user)
	h.Array[index].insertUserToLinkedList(user)
}

// Continuation of hash function insertion to place it within a linked list as a node.
func (b *UserLinkedList) insertUserToLinkedList(User User) {
	if !b.searchUserIDLinkedList(User) {
		newNode := &UserNode{Key: User}
		newNode.Next = b.Head
		b.Head = newNode
		//fmt.Println(k)
	} else {
		//fmt.Println(k, "already exists")
	}
}

// A behavior of a hash table object used to search of an user object within the hash table using account ID of each auction.
func (h *UserHashTable) SearchUserIDHashTable(user User) bool {
	index := HashUser(user)
	return h.Array[index].searchUserIDLinkedList(user)
}

// Continuation of seachUserIDHashTable() function to continue the search within the linked list at the hash index location.
func (b *UserLinkedList) searchUserIDLinkedList(user User) bool { //For search the user by using accouintID
	currentNode := b.Head
	temp := user.AccountID
	for currentNode != nil {
		if currentNode.Key.AccountID == temp {
			return true
		}
		currentNode = currentNode.Next
	}
	fmt.Println("There is no auction with that ID in the memory.")
	return false
}

func (h *UserHashTable) UserHashAccessDelete(key User) {
	index := HashUser(key)
	h.Array[index].DeleteUserInLinkedList(key)
}

func (b *UserLinkedList) DeleteUserInLinkedList(k User) {

	if b.Head.Key.AccountID == k.AccountID {
		b.Head = b.Head.Next
		return
	}
	previousNode := b.Head
	for previousNode.Next != nil {
		if previousNode.Next.Key.AccountID == k.AccountID {
			previousNode.Next = previousNode.Next.Next
			return
		}
		previousNode = previousNode.Next
	}
}

func (h *UserHashTable) UserHashAccessUpdate(Key User) {
	index := HashUser(Key)
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
