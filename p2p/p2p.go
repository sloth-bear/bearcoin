package p2p

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sloth-bear/bearcoin/utils"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)

	for {
		_, p, err := conn.ReadMessage() // blocking operation
		if err != nil {
			conn.Close()
			break
		}

		fmt.Println("Message arrived! ")

		pStr := fmt.Sprintf("%s", p)
		fmt.Printf("%s", pStr+"\n\n")

		time.Sleep(2 * time.Second)

		resMsg := "Client: " + pStr
		conn.WriteMessage(1, []byte(resMsg))
	}
}
