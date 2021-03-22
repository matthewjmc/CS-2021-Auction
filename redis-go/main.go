package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/go-redis/redis"
)

type Client struct {
	Command string
}

type Address struct {
	addr1: "com1.mcmullin.org:19530"
	addr2: "com2.mcmullin.org:19530"
}

type Auction struct {
	Description      string
	AddressIP        string
	ConnectedClients int
}

func main() {
	//updateConnectedUsers("1")
	// key1 := "10"
	// value1 := &Auction{Description: "someName1", AddressIP: "addr.1.23456781", ConnectedClients: +1}
	// SetKey("1", value1)
	//val := Client{Command: "create", Description: "asdd"}
	//keygen(Client{Command: "create", Description: "asdd"})
	//deleteAuction(Client{Command: "stop", Description: "asdd"}, "1")
}

// set auctionID or key to auction struct input
func SetKey(key string, value interface{}) (key1 string) {
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
	return key
}

// return to matthew client
func keygen(c Client) (key string) {
	if c.Command != "create" {
		fmt.Println("not create command")
	}
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost: 6379",
		Password: "",
		DB:       0,
	})
	temp, err := client.Keys("*").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(temp)
	var count = len(temp)
	fmt.Println(count)
	var newkey = count + 1
	key1 := strconv.Itoa(newkey)
	fmt.Println(key1, reflect.TypeOf(key1))
	//set IP address to that key
	if S1_Usage > S2_Usage{
		src := Auction{}
		err = json.Unmarshal([]byte(val), &src)
		var newval = Auction{
		Description:      src.Description,
		AddressIP:        Address.addr1,
		ConnectedClients: src.ConnectedClients,
	}
	entry, err := json.Marshal(newval)
	client.Set(id, entry, 0)
	}else if S2_Usage > S1_Usage{
		src := Auction{}
		err = json.Unmarshal([]byte(val), &src)
		var newval = Auction{
		Description:      src.Description,
		AddressIP:        Address.addr2,
		ConnectedClients: src.ConnectedClients,
	}
	entry, err := json.Marshal(newval)
	client.Set(id, entry, 0)
	}else {
		src := Auction{}
		err = json.Unmarshal([]byte(val), &src)
		var newval = Auction{
		Description:      src.Description,
		AddressIP:        Address.addr1,
		ConnectedClients: src.ConnectedClients,
	}
	entry, err := json.Marshal(newval)
	client.Set(id, entry, 0)
	}
	return key1
}

// get addr to send to nonthicha
func GetAddressByID(key string) (ip string) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost: 6379",
		Password: "",
		DB:       0,
	})
	val, err := client.Get(key).Result()
	if err == redis.Nil || err != nil {
		log.Fatal(err)
	}
	src := Auction{}
	err = json.Unmarshal([]byte(val), &src)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(src.AddressIP)
	return src.AddressIP
}

// get auction description from ID input
func GetDescriptionByID(key string) (ip string) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost: 6379",
		Password: "",
		DB:       0,
	})
	val, err := client.Get(key).Result()
	if err == redis.Nil || err != nil {
		log.Fatal(err)
	}
	src := Auction{}
	err = json.Unmarshal([]byte(val), &src)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(src.Description)
	return src.Description
}

func updateConnectedUsers(id string) (user int) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost: 6379",
		Password: "",
		DB:       0,
	})
	val, err := client.Get(id).Result()
	if err != nil {
		fmt.Println(err)
	}
	src := Auction{}
	err = json.Unmarshal([]byte(val), &src)
	var count = src.ConnectedClients
	count = count + 1
	var newval = Auction{
		Description:      src.Description,
		AddressIP:        src.AddressIP,
		ConnectedClients: count,
	}
	entry, err := json.Marshal(newval)
	client.Set(id, entry, 0)
	return count
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
	return val
}

// delete auction from command and key
func deleteAuction(c Client, key string) bool {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost: 6379",
		Password: "",
		DB:       0,
	})
	if c.Command == "stop" {
		val, err := client.Del(key).Result()
		if err != nil {
			fmt.Println(err)
			return false
		}
		fmt.Println(val)
	}
	return true
}
