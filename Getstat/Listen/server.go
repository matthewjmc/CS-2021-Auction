package main

// Locate on backend servers 1 and 2 to get cpu stat and send to the redis directly without
// passing the load balancer to reduce latency

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-redis/redis"
)

type Data struct {
	Usage float64
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "10.104.0.11: 80",
		Password: "",
		DB:       0,
	})
	defer client.Close()
	for {
		data := Usage()
		//fmt.Println(data)
		val, err := client.Get("1").Result()
		if err != nil {
			fmt.Println(err)
		}
		src := Data{}
		err = json.Unmarshal([]byte(val), &src)

		entry, err := json.Marshal(data)
		client.Set("1", entry, 0)
	}

}

func Usage() (data float64) {
	before := collectCPUStats()

	time.Sleep(time.Duration(1) * time.Second)
	after := collectCPUStats()

	total := float64(after.Total - before.Total)
	idle := float64(after.Idle-before.Idle) / total * 100

	//vs := strconv.FormatFloat(float64(idle), 'f', 2, 64)
	//fmt.Println(reflect.TypeOf(vs))
	fmt.Println("cpu idle:", idle)
	return idle
}

// Stats represents cpu statistics for linux
type Stats struct {
	User      uint64
	Nice      uint64
	System    uint64
	Idle      uint64
	Iowait    uint64
	Irq       uint64
	Softirq   uint64
	Steal     uint64
	Guest     uint64
	GuestNice uint64
	Total     uint64
	CPUCount  int
	StatCount int
}

type cpuStat struct {
	name string
	ptr  *uint64
}

func collectCPUStats() *Stats {
	file, err := os.Open("/proc/stat")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	out := io.Reader(file)
	scanner := bufio.NewScanner(out)
	var cpu Stats

	cpuStats := []cpuStat{
		{"user", &cpu.User},
		{"nice", &cpu.Nice},
		{"system", &cpu.System},
		{"idle", &cpu.Idle},
		{"iowait", &cpu.Iowait},
		{"irq", &cpu.Irq},
		{"softirq", &cpu.Softirq},
		{"steal", &cpu.Steal},
		{"guest", &cpu.Guest},
		{"guest_nice", &cpu.GuestNice},
	}

	if !scanner.Scan() {
		fmt.Println("failed to scan /proc/stat")
	}

	valStrs := strings.Fields(scanner.Text())[1:]
	cpu.StatCount = len(valStrs)
	for i, valStr := range valStrs {
		val, err := strconv.ParseUint(valStr, 10, 64)
		if err != nil {
			fmt.Println("failed to scan", cpuStats[i].name)
		}
		*cpuStats[i].ptr = val
		cpu.Total += val
	}

	// Since cpustat[CPUTIME_USER] includes cpustat[CPUTIME_GUEST], subtract the duplicated values from total.
	// https://github.com/torvalds/linux/blob/4ec9f7a18/kernel/sched/cputime.c#L151-L158
	cpu.Total -= cpu.Guest
	// cpustat[CPUTIME_NICE] includes cpustat[CPUTIME_GUEST_NICE]
	cpu.Total -= cpu.GuestNice

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") && unicode.IsDigit(rune(line[3])) {
			cpu.CPUCount++
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scan error for /proc/stat:", err)
	}

	return &cpu
}
