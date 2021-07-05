package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/goLangCoin/explorer"
	"github.com/goLangCoin/rest"
)

func usage() {
	fmt.Printf("Welcome golang coin project\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port : Set the PORT of the server\n\n")
	fmt.Printf("-mode : Choose between 'html' and 'rest'\n\n")
	os.Exit(0)
}

func main() {
	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	flag.Parse()

	switch *mode {
	case "html":
		explorer.Start(*port)
	case "rest":
		rest.Start(*port)
	default:
		usage()
	}

	// rest := flag.NewFlagSet("rest", flag.ExitOnError)
	// portFlag := rest.Int("port", 4000, "Sets the port of the server")
}
