package main

import (
	"github.com/goLangCoin/cli"
	"github.com/goLangCoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}
