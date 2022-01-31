package main

import (
	"flag"
	"fmt"
	"os"
)

func usages() {
	fmt.Printf("Welcome to bearcoin!\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-mode=rest: Choose between 'html' and 'rest'\n")
	fmt.Printf("-port=4000 Set port of the server\n\n")

	os.Exit(0)
}

var mode string
var port int

func main() {
	flag.StringVar(&mode, "mode", "rest", "Choose between 'html' and 'rest'")
	flag.IntVar(&port, "port", 4000, "Sets port of the server")
	flag.Parse()

	if !flag.Parsed() {
		os.Exit(0)
	}

	switch mode {
	case "html":
		fmt.Println("Start explorer")
	case "rest":
		fmt.Println("Start REST API")
	default:
		usages()
	}
}
