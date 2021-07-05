package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/goLangCoin/blockchain"
	"github.com/goLangCoin/utils"
	"github.com/gorilla/mux"
)

var port string

type Message struct {
	Message string
}

type URL string

func (u URL) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

// List of URLs
type URLDescription struct {
	URL         URL    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

func (u URLDescription) String() string {
	return "Hello I am the URL Description"
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []URLDescription{
		{
			URL:         URL("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         URL("/blocks"),
			Method:      "GET",
			Description: "Get All Block",
		},
		{
			URL:         URL("/blocks"),
			Method:      "POST",
			Description: "Add a Block",
			Payload:     "data:string",
		},
		{
			URL:         URL("/blocks/{height}"),
			Method:      "POST",
			Description: "Get a Block",
		},
	}
	rw.Header().Add("Content-Type", "application/json")
	// simple version
	json.NewEncoder(rw).Encode(data)
	/*
		something hard version
		b, err := json.Marshal(data)
		utils.HandleErr(err)
		fmt.Fprintf(rw, "%s", b)
	*/
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.AllBlocks())
	case "POST":
		var message Message
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&message))
		blockchain.GetBlockchain().AddBlock(message.Message)
		rw.WriteHeader(http.StatusCreated)
	}
}

func getBlock(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["height"])
	utils.HandleErr(err)

	block, err := blockchain.GetBlockchain().GetBlock(id)
	if err == blockchain.ErrNotFound {
		json.NewEncoder(rw).Encode(fmt.Sprint(err))
	} else {
		json.NewEncoder(rw).Encode(block)
	}
}

func writeContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{height:[0-9]+}", getBlock).Methods("GET")
	fmt.Printf("Listening Server http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
