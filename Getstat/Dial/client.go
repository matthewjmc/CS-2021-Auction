package main

// Locate on the load balance

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

type Data struct {
	Usage float64
}

var S1_Usage, S2_Usage float64
var Usage float64

func main() {
	ln1, err := net.Dial("tcp4", "com1.mcmullin.org:20001")
	ln2, err := net.Dial("tcp4", "com2.mcmullin.org:20001")
	if err != nil {
		fmt.Println(err)
	}
	//for {
	//fmt.Println("top")
	go GetStat(ln1)
	fmt.Println("Start Getstat")
	GetStat(ln2)
	//fmt.Println("Bottom")
	//}
	//time.Sleep(1000 * time.Second)

}

// update numbers of connected users and pass to reverse proxy
// func UpdateUsage(id string, value float64) {
// 	client := redis.NewClient(&redis.Options{
// 		Addr:     "localhost: 6379",
// 		Password: "",
// 		DB:       0,
// 	})
// 	val, err := client.Get(id).Result()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	src := Data{}
// 	err = json.Unmarshal([]byte(val), &src)
// 	newval := value
// 	entry, err := json.Marshal(newval)
// 	client.Set(id, entry, 0)
// 	client.Close()
// }

func GetStat(conn net.Conn) {
	defer conn.Close()
	client := redis.NewClient(&redis.Options{
		Addr:     "10.104.0.11: 80",
		Password: "",
		DB:       0,
	})
	defer client.Close()
	for {
		//fmt.Println("Test")
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Exiting TCP server!")
		}

		usage := string(netData)
		//fmt.Println(usage)
		if usage[12] == 116 {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(usage, -1)
			for _, element := range submatchall {
				S2temp := element
				S2_Usage_temp, err := strconv.ParseFloat(S2temp, 64)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(S2_Usage_temp)
				Usage = S2_Usage_temp

				val, err := client.Get("2").Result()
				if err != nil {
					fmt.Println(err)
				}
				src := Data{}
				err = json.Unmarshal([]byte(val), &src)

				entry, err := json.Marshal(Usage)
				client.Set("2", entry, 0)
			}
			// fmt.Println(Usage)
		}
		if usage[12] == 111 {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(usage, -1)
			for _, element := range submatchall {
				S1temp := element
				S1_Usage_temp, err := strconv.ParseFloat(S1temp, 64)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(S1_Usage_temp)
				Usage = S1_Usage_temp

				val, err := client.Get("1").Result()
				if err != nil {
					fmt.Println(err)
				}
				src := Data{}
				err = json.Unmarshal([]byte(val), &src)

				entry, err := json.Marshal(Usage)
				client.Set("1", entry, 0)
			}
		}
		//fmt.Println(Usage)
		// t := time.Now()
		// myTime := t.Format(time.RFC3339) + "\n"
		// conn.Write([]byte(myTime))
	}

}
