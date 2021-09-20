package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/goLangCoin/blockchain"
	"github.com/goLangCoin/utils"
	"github.com/goLangCoin/wallet"
	"github.com/gorilla/mux"
)

var port string

type addTxPayload struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type Message struct {
	Message string
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type myWalletResponse struct {
	Address string `json:"address"`
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
			URL:         URL("/status"),
			Method:      "GET",
			Description: "See the Status of the Blockchain",
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
			URL:         URL("/blocks/{hash}"),
			Method:      "POST",
			Description: "Get a Block",
		},
		{
			URL:         URL("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an Address",
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
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain())))
	case "POST":
		blockchain.Blockchain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
	}
}

func getBlock(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	if err == blockchain.ErrNotFound {
		json.NewEncoder(rw).Encode(errorResponse{fmt.Sprint(err)})
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

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.Blockchain())
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	// queyr URL 받는법
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(address, blockchain.Blockchain())
		json.NewEncoder(rw).Encode(balanceResponse{address, amount})
	default:
		// fmt.Println(blockchain.UTxOutsByAddress(address, blockchain.Blockchain()))
		utils.HandleErr(json.NewEncoder(rw).
			Encode(blockchain.UTxOutsByAddress(address, blockchain.Blockchain())))
	}

}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	// json.NewEncoder(rw).Encode(myWalletResponse{address})
	json.NewEncoder(rw).Encode(struct {
		Address string `json:"address"`
	}{address}) // anonymous struct
}

func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	router.Use(writeContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[0-9a-f]+}", getBlock).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	fmt.Printf("Listening Server http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
