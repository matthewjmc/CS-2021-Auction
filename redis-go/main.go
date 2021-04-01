package CS-2021-Auction

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"
	"strconv"

	gs "CS-2021-Auction/getstat"

	"github.com/go-redis/redis"
)

type Client struct {
	Command string
}

type Auction struct {
	Description      string
	AddressIP        string
	ConnectedClients int
}

// func main() {
// 	updateConnectedUsers("1")
// 	key1 := "10"
// 	value1 := &Auction{Description: "someName1", AddressIP: "addr.1.23456781", ConnectedClients: +1}
// 	SetKey("1", value1)
// 	val := Client{Command: "create", Description: "asdd"}
// 	keygen(Client{Command: "create", Description: "asdd"})
// 	deleteAuction(Client{Command: "stop", Description: "asdd"}, "1")
// }

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

// return new key to matthew's side (client) and lock the ip addr for that key
func KeyGen(cli Client) (key string) {
	if cli.Command != "create" {
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
	//fmt.Println(temp)
	var count = len(temp)
	fmt.Println(count)
	var newkey = count + 1
	key1 := strconv.Itoa(newkey)
	fmt.Println(key1, reflect.TypeOf(key1))

	//get CPU stat
	// ln, err := net.Listen("tcp4", ":19530")
	// 		if err != nil {
	// 				fmt.Println(err)
	// 		}
	// 	fmt.Println("in1")
	// for{
	// 	conn, err := ln.Accept()
	// 	fmt.Println("in2")
	// 	if err != nil {
	// 			fmt.Println(err)
	// 			return
	// 	}
	// 	fmt.Println("in3")
	// 	gs.GetStat(conn)
	// 	fmt.Println("in4")
	// 	}

	ln1, err := net.Dial("tcp4", "com1.mcmullin.org:19530")
	ln2, err := net.Dial("tcp4", "com2.mcmullin.org:19530")
	if err != nil {
		fmt.Println(err)
	}
	gs.GetStat(ln1)
	gs.GetStat(ln2)
	fmt.Println(gs.S1_Usage)
	fmt.Println(gs.S2_Usage)

	//set IP address to that key
	if gs.S1_Usage > gs.S2_Usage {
		value := &Auction{Description: "", AddressIP: "com1.mcmullin.org", ConnectedClients: 0}
		entry, err := json.Marshal(value)
		if err != nil {
			fmt.Println(err)
		}
		err = client.Set(key1, entry, 0).Err()
		if err != nil {
			fmt.Println(err)
		}
	} else if gs.S2_Usage > gs.S1_Usage {
		value := &Auction{Description: "", AddressIP: "com2.mcmullin.org", ConnectedClients: 0}
		entry, err := json.Marshal(value)
		if err != nil {
			fmt.Println(err)
		}
		err = client.Set(key1, entry, 0).Err()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		value := &Auction{Description: "", AddressIP: "com1.mcmullin.org", ConnectedClients: 0}
		entry, err := json.Marshal(value)
		if err != nil {
			fmt.Println(err)
		}
		err = client.Set(key1, entry, 0).Err()
		if err != nil {
			fmt.Println(err)
		}
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

// update numbers of connected users
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
	if c.Command == "delete" {
		val, err := client.Del(key).Result()
		if err != nil {
			fmt.Println(err)
			return false
		}
		fmt.Println(val)
	}
	return true
}
