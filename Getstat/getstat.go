package load_balance

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	//"reflect"
	"strconv"
	//rd "load_balance/redis"
)

var S1_Usage, S2_Usage float64

func GetStat(conn net.Conn) {
	//defer conn.Close()
	netData, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	if strings.TrimSpace(string(netData)) == "STOP" {
		fmt.Println("Exiting TCP server!")
	}

	usage := string(netData)
	fmt.Println(usage)
	if usage[12] == 116 {
		re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		submatchall := re.FindAllString(usage, -1)
		for _, element := range submatchall {
			S2temp := element
			S2_Usage_temp, err := strconv.ParseFloat(S2temp, 64)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(S2_Usage)
			S2_Usage = S2_Usage_temp
		}
		fmt.Println(S2_Usage)
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
			fmt.Println(S1_Usage)
			S1_Usage = S1_Usage_temp
		}
		fmt.Println(S1_Usage)
	}
	t := time.Now()
	myTime := t.Format(time.RFC3339) + "\n"
	conn.Write([]byte(myTime))

}
