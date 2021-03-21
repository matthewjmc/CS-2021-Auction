package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func main() {
	before := collectCPUStats()

	time.Sleep(time.Duration(1) * time.Second)
	after := collectCPUStats()

	total := float64(after.Total - before.Total)
	idle := float64(after.Idle-before.Idle) / total * 100
	fmt.Println("cpu idle: ", idle)
}

// cpu statistics for linux
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

	// Since cpustat[CPUTIME_USER] includes cpustat[CPUTIME_GUEST]
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

// https://github.com/torvalds/linux/blob/4ec9f7a18/kernel/sched/cputime.c#L151-L158
