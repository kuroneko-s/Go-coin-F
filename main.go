package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/goLangCoin/blockchain"
)

const PORT string = ":4000"

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	// template.Must => err를 체크해줌
	tmpl := template.Must(template.("templates/pages/home.gohtml"))
	data := homeData{"Home", blockchain.AllBlocks()}
	tmpl.Execute(rw, data)
}

func main() {
	http.HandleFunc("/", home)
	fmt.Printf("Listening on http://localhost%s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
