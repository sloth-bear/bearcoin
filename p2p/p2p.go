package p2p

import (
	"fmt"
	"net/http"

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
	initPeer(conn, "xx", "xx")
}

func AddPeer(address, port string) {
	url := fmt.Sprintf("ws://%s:%s/ws", address, port)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	utils.HandleErr(err)
	initPeer(conn, address, port)
}
