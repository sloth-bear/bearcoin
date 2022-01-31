package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/sloth-bear/bearcoin/explorer"
	"github.com/sloth-bear/bearcoin/rest"
)

func usages() {
	fmt.Printf("Welcome to bearcoin!\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-mode:    Choose between 'html' and 'rest' (default: rest)\n")
	fmt.Printf("-port:    Set port of the server (default: 4000)\n\n")

	os.Exit(0)
}

func Start() {
	if len(os.Args) == 1 {
		usages()
	}

	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")
	port := flag.Int("port", 4000, "Sets port of the server")

	flag.Parse()

	switch *mode {
	case "html":
		explorer.Start(*port)
	case "rest":
		rest.Start(*port)
	default:
		usages()
	}
}
