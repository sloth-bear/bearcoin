package p2p

import (
	"fmt"
	"net/http"
	"strings"

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
	result := strings.Split(r.RemoteAddr, ":")
	port := strings.TrimPrefix(r.URL.Query().Get("openPort"), ":")
	initPeer(conn, result[0], port)
}

func AddPeer(address, port, openPort string) {
	url := fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	utils.HandleErr(err)
	initPeer(conn, address, port)
}
