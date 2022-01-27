package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sloth-bear/bearcoin/blockchain"
	"github.com/sloth-bear/bearcoin/utils"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type blockRequestBody struct {
	Message string
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{URL: url("/"), Method: "GET", Description: "See Documentation"},
		{URL: url("/blocks"), Method: "POST", Description: "Add A Block", Payload: "data:string"},
		{URL: url("/blocks"), Method: "GET", Description: "See All Blocks"},
		{URL: url("/blocks/{id}"), Method: "GET", Description: "See A Block"},
	}

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			rw.Header().Add("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())
		}
	case "POST":
		{
			var blockRequest blockRequestBody
			utils.HandleErr(json.NewDecoder(r.Body).Decode(&blockRequest))

			blockchain.GetBlockchain().AddBlock(blockRequest.Message)

			rw.WriteHeader(http.StatusCreated)
		}
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	height, err := strconv.Atoi(vars["height"])

	utils.HandleErr(err)

	json.NewEncoder(rw).Encode(blockchain.GetBlockchain().GetBlock(height))
}

func Start(aPort int) {
	handler := mux.NewRouter()
	handler.HandleFunc("/", documentation).Methods("GET")
	handler.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	handler.HandleFunc("/blocks/{height}", block).Methods("GET")

	port = fmt.Sprintf(":%d", aPort)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
