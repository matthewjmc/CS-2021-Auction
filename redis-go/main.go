package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"

	"github.com/go-redis/redis"
)

type Client struct {
	ClientObj net.Conn
	UserID    string
}

type Auction struct {
	Description      string
	AddressIP        string
	ConnectedClients int
}

func main() {
	GetAddressByID("9")
	// key1 := "9"
	// value1 := &Auction{Description: "someName", AddressIP: "addr:12345678", ConnectedClients: 1}
	// setKey(key1, value1)
}

func SetKey(key string, value interface{}) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost: 6379",
		Password: "",
		DB:       0,
	})
	entry, err := json.Marshal(value)
	if err != nil {
		fmt.Println(err)
	}
	err = client.Set(key, entry, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(entry)
}

// return only id and description to the client
func getAuctionDescription(id string) (auctionid, descr string) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost: 6379",
		Password: "",
		DB:       0,
	})
	val, err := client.Get(id).Result()
	if err != nil {
		fmt.Println(err)
	}
	newval := strings.Split(val, ",")[0]
	newval = strings.ReplaceAll(newval, "{", "")
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	newval = reg.ReplaceAllString(newval, "")
	desc := strings.ReplaceAll(newval, "Description", "")
	fmt.Println(desc)
	return id, desc
}

// return all column pass to nonthicha
func GetAddressByID(id string) (addr string) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost: 6379",
		Password: "",
		DB:       0,
	})
	val, err := client.Get(id).Result()
	if err != nil {
		fmt.Println(err)
	}
	newval := strings.Split(val, ",")[1]
	newval = strings.ReplaceAll(newval, "{", "")
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	newval = reg.ReplaceAllString(newval, "")
	addr = strings.ReplaceAll(newval, "AddressIP", "")
	fmt.Println(addr)
	return addr
}

func updateConnectedUsers(id string, count int64) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost: 6379",
		Password: "",
		DB:       0,
	})
	// TODO auto-increment value each time user connect
	val, err := client.IncrBy(id, count).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
}

// return all auctionID
func getAllAuctionID() (val []string) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost: 6379",
		Password: "",
		DB:       0,
	})
	val, err := client.Keys("*").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
	return val
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
