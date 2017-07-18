package main

import "fmt"
import "Oon/brain"

func main() {
	fmt.Println("Starting Oon ...")
	b := brain.New("config.conf")
	b.Start()
	b.Destroy()
}
