package load_balance

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	rv "RV/ReverseProxy"

	"github.com/go-redis/redis"
)

type Client struct {
	Command string
}

type Auction struct {
	AddressIP        string
	ConnectedClients int
}

type Data struct {
	Usage float64
}

//get command and track which function to call for client
func CommandFunction(cmd rv.Package) (string, rv.Package) {
	temp := rv.Package{}
	if cmd.Command == "create" {
		IP, Init := KeyGen(cmd)
		return IP, Init
	} else if cmd.Command == "join" {
		temp_key := strconv.FormatUint(cmd.AuctionID, 10)
		IP, Init := RequestConnection(temp_key, cmd)
		return IP, Init
	}
	return "NULL", temp
}

func getToken() uint64 {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	data := binary.BigEndian.Uint64(randomBytes) * uint64(time.Now().Unix())
	return data
}

// return new key to matthew's side (client) and lock the ip addr for that key
func KeyGen(init rv.Package) (key string, data rv.Package) {
	client := redis.NewClient(&redis.Options{
		Addr:     "10.104.0.11: 80",
		Password: "",
		DB:       0,
	})
	newkey := getToken()
	key1 := strconv.FormatUint(newkey, 10)

	S1, err := client.Get("1").Result()
	if err != nil {
		fmt.Println(err)
	}
	S1_Usage, _ := strconv.ParseFloat(S1, 64)

	S2, err := client.Get("2").Result()
	if err != nil {
		fmt.Println(err)
	}
	S2_Usage, _ := strconv.ParseFloat(S2, 64)

	var value Auction

	//set IP address to that key
	if S1_Usage > S2_Usage {
		value.AddressIP = "10.104.0.9:19530"
		value.ConnectedClients = 0
		entry, err := json.Marshal(value)
		if err != nil {
			fmt.Println(err)
		}
		err = client.Set(key1, entry, 0).Err()
		if err != nil {
			fmt.Println(err)
		}
	} else if S2_Usage > S1_Usage {
		value.AddressIP = "10.104.0.8:19530"
		value.ConnectedClients = 0
		entry, err := json.Marshal(value)
		if err != nil {
			fmt.Println(err)
		}
		err = client.Set(key1, entry, 0).Err()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		value.AddressIP = "10.104.0.9:19530"
		value.ConnectedClients = 0
		entry, err := json.Marshal(value)
		if err != nil {
			fmt.Println(err)
		}
		err = client.Set(key1, entry, 0).Err()
		if err != nil {
			fmt.Println(err)
		}
	}
	client.Close()
	init.Data.Value = newkey
	if err != nil {
		fmt.Println(err)
	}
	return value.AddressIP, init
}

// get addr to send to nonthicha reverse proxy
func RequestConnection(key string, init rv.Package) (string, rv.Package) { //(ip string)
	client := redis.NewClient(&redis.Options{
		Addr:     "10.104.0.11: 80",
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
	client.Close()
	if err != nil {
		fmt.Println(err)
	}
	return src.AddressIP, init
}

// update numbers of connected users and pass to reverse proxy
func UpdateConnections(id string) (user int) {
	client := redis.NewClient(&redis.Options{
		Addr:     "10.104.0.11: 80",
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
		AddressIP:        src.AddressIP,
		ConnectedClients: count,
	}
	entry, err := json.Marshal(newval)
	client.Set(id, entry, 0)
	client.Close()
	return count
}

// return all auctionID
func getAllAuctionID() (val []string) {
	client := redis.NewClient(&redis.Options{
		Addr:     "10.104.0.11: 80",
		Password: "",
		DB:       0,
	})
	val, err := client.Keys("*").Result()
	if err != nil {
		fmt.Println(err)
	}
	client.Close()
	return val
}

// delete auction from command and key
func deleteAuction(c Client, key string) bool {
	client := redis.NewClient(&redis.Options{
		Addr:     "10.104.0.11: 80",
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
	client.Close()
	return true
}
