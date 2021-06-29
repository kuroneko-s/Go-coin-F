package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

const PORT string = ":4000"

type homeData struct {
	PageTitle string
}

func home(rw http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		log.Fatal(err)
	}
	data := homeData{"Home"}
	tmpl.Execute(rw, data)
}

func main() {
	http.HandleFunc("/", home)
	fmt.Printf("Listening on http://localhost%s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
