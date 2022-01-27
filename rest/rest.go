package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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

func getBlock(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/blocks/")
	fmt.Println(path)

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(blockchain.GetBlockchain().GetBlock(path))
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)

	handler := http.NewServeMux()
	handler.HandleFunc("/", documentation)
	handler.HandleFunc("/blocks", blocks)
	handler.HandleFunc("/blocks/", getBlock)

	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
