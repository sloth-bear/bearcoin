package main

import (
	"fmt"
	"os"

	"github.com/sloth-bear/bearcoin/explorer"
)

func usages() {
	fmt.Printf("Welcome to bearcoin!\n\n")
	fmt.Printf("Please use the following commands:\n\n")
	fmt.Printf("explorer:   Start the HTML Explorer\n")
	fmt.Printf("rest:	    Start the REST API (recommended)\n\n")

	os.Exit(0)
}

func main() {
	args := os.Args

	if len(args) < 2 {
		usages()
	}

	switch args[1] {
	case "explorer":
		fmt.Println("Start Explorer")
		explorer.Start(3000)

	case "rest":
		fmt.Println("Start REST API")
		explorer.Start(4000)

	default:
		usages()
	}
}
