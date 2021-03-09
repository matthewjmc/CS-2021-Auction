package main

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println("<---------- Commiting Code to Github Now ---------->")
	data,err := exec.Command("git","commit","-a","-m","'Pushing Everything into Git Time to Sleep'").Output()
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(string(data))
	fmt.Println("<---------- Pushing Code to Github Now ---------->")
	data,err = exec.Command("git","push","origin","Matthew").Output()
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(string(data))
	data,err = exec.Command("sl").Output()
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(string(data))
}
