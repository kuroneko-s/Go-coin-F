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

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Place       string `json:"place,omitempty"`
}

// get json data step 1
type addBlockBody struct {
	Message string
}
type url string
type errorResponse struct {
	ErrMessage string `json:"errMessage"`
}

var port string

func (u url) MarshalText() (text []byte, err error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "show description doc",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "See All blocks",
			Place:       "body:string",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "add blocks",
			Place:       "body:string",
		},
		{
			URL:         url("/blocks/{height}"),
			Method:      "GET",
			Description: "See A Block",
		},
	}
	json.NewEncoder(rw).Encode(data)

	// b, err := json.Marshal(data)
	// rw.Write(b)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.AllBlocks())
	case "POST":
		// Get the json type data for client
		// step 2
		var addBlockBody addBlockBody
		utils.CheckErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.GetBlockchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	height, err := strconv.Atoi(vars["height"])
	utils.CheckErr(err)
	block, err := blockchain.GetBlockchain().GetBlock(height)
	encoder := json.NewEncoder(rw)

	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func Start(aPort int) {
	// NewServeMux url과 handler를 매핑해줌 (dispatcherServlet)
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{height:[0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
