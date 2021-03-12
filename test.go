// Golang program to illustrate the usage of
// time.Parse() function

// Including main package
package main

// Importing fmt and time
import (
	"fmt"
	"time"
)

// Calling main
func main() {

	init := time.Now().Format(time.RFC3339Nano)
	time.Sleep(1 * time.Second)
	end := time.Now().Format(time.RFC3339Nano)
	fmt.Println(end > init)
}
