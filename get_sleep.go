package main

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println("<---------- Pushing Code to Github Now ---------->")
	cmd := exec.Command("git commit -m", "'Finished for the Day pushing to Github now'")
	if err := cmd.Run(); err != nil {
		fmt.Println(err)

	}
	cmd = exec.Command("git push origin Matthew")
	if err := cmd.Run(); err != nil {
		fmt.Println(err)

	}
	fmt.Println("<---------- Bye Bye ---------->")
	cmd = exec.Command("sl")
	if err := cmd.Run(); err != nil {
		fmt.Println(err)

	}
}
