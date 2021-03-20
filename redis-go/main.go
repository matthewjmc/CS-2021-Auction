package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-redis/redis"
)

type Client struct {
	ClientObj net.Conn
	UserID    string
}

type Auction struct {
	Description      string `json: "description"`
	AddressIP        string `json: "address"`
	ConnectedClients int    `json: "numbers of connect users"`
}

func main() {
	getAuctionByID("1")
	//fmt.Println("start")

}

func createNewAuction() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost: 6379",
		Password: "",
		DB:       0,
	})
	d := Auction{
		Description:      "auction1",
		AddressIP:        "addr:1",
		ConnectedClients: 0,
	}
	json, err := json.Marshal(d)

	err = client.Set("2", json, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reflect.TypeOf(json))
}

func postNew() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost: 6379",
		Password: "",
		DB:       0,
	})
	d := &Auction{
		Description:      "auction1",
		AddressIP:        "addr:1",
		ConnectedClients: 0,
	}
	var m = make(map[string]interface{})

}

// return only id and description to the client
func getAuctionDescription(id string) (auctionid, val string) {
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
	returnval := strings.ReplaceAll(newval, "Description", "")
	return id, returnval
}

// return all column pass to nonthicha
func getAuctionByID(id string) (val string) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost: 6379",
		Password: "",
		DB:       0,
	})
	val, err := client.Get(id).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
	return val
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
