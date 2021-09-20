package main

import (
	"github.com/goLangCoin/cli"
	"github.com/goLangCoin/db"
	"github.com/goLangCoin/wallet"
)

func main() {
	defer db.Close()
	db.DB()
	cli.Start()
	wallet.Wallet()
}
