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

type errorResponse struct {
	Message string `json:"message"`
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{URL: url("/"), Method: "GET", Description: "See Documentation"},
		{URL: url("/blocks"), Method: "POST", Description: "Add A Block", Payload: "data:string"},
		{URL: url("/blocks"), Method: "GET", Description: "See All Blocks"},
		{URL: url("/blocks/{id}"), Method: "GET", Description: "See A Block"},
	}

	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
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
	vars := mux.Vars(r)
	height, err := strconv.Atoi(vars["height"])
	utils.HandleErr(err)

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
	handler := mux.NewRouter()
	handler.HandleFunc("/", documentation).Methods("GET")
	handler.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	handler.HandleFunc("/blocks/{height}", block).Methods("GET")
	handler.Use(jsonContentTypeMiddleware)

	port = fmt.Sprintf(":%d", aPort)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
