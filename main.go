package main

import (
	"github.com/sloth-bear/bearcoin/cli"
	"github.com/sloth-bear/bearcoin/db"
)

func main() {
	defer db.Close()
	db.InitDB()
	cli.Start()
}
