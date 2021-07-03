package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const PORT string = ":4000"

type URL string

func (u URL) MarshalText() (text []byte, err error) {
	url := fmt.Sprintf("http://localhost%s%s", PORT, u)
	return []byte(url), nil
}

type URLDescription struct {
	URL         URL    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Place       string `json:"place,omitempty"`
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []URLDescription{
		{
			URL:         URL("/"),
			Method:      "GET",
			Description: "show description doc",
		},
		{
			URL:         URL("/add"),
			Method:      "POST",
			Description: "add blocks",
			Place:       "body:string",
		},
	}
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)

	// b, err := json.Marshal(data)
	// rw.Write(b)
}

func main() {
	http.HandleFunc("/", documentation)
	fmt.Printf("Listening on http://localhost%s", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
