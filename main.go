package main

import (
	rd "load_balance/redis"
)

//var S1_Usage, S2_Usage string

func main() {
	//var wg sync.WaitGroup
	//rd.KeyGen(rd.Client{Command: "create"})
	rd.GetAddressByID("1")
}
