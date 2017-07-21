package main

import "fmt"
import "Oon/brain"
import "os"

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <config>\n", os.Args[0])
		return
	}

	fmt.Println("Starting Oon ...")
	b := brain.New(os.Args[1])
	b.Start()
	b.Destroy()
}
