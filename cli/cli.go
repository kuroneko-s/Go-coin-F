package cli

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
	fmt.Printf("-Rport : Set the PORT of the REST API server\n\n")
	fmt.Printf("-Eport : Set the PORT of the explorer server\n\n")
	fmt.Printf("-mode : Choose between 'html' and 'rest', 'dual' is run REST API(4000) and explorer(4001)\n\n")
	os.Exit(0)
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	Rport := flag.Int("Rport", 4000, "Set port of the REST API server")
	Eport := flag.Int("Eport", 4001, "Set port of the explorer server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest', 'dual' is run REST API(4000) and explorer(4001)")

	flag.Parse()

	switch *mode {
	case "html":
		explorer.Start(*Eport)
	case "rest":
		rest.Start(*Rport)
	case "dual":
		go explorer.Start(*Eport)
		rest.Start(*Rport)
	default:
		usage()
	}

	// rest := flag.NewFlagSet("rest", flag.ExitOnError)
	// portFlag := rest.Int("port", 4000, "Sets the port of the server")
}
