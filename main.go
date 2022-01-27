package main

import (
	"github.com/sloth-bear/bearcoin/explorer"
	"github.com/sloth-bear/bearcoin/rest"
)

func main() {
	go explorer.Start(3000)
	rest.Start(4000)
}
