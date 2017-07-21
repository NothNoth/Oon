package main

import "fmt"
import "Oon/brain"
import "os"

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <motors config> <camera config>\n", os.Args[0])
		return
	}

	fmt.Println("Starting Oon ...")
	b := brain.New(os.Args[1], os.Args[2])
	b.Start()
	b.Destroy()
}
