package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sloth-bear/bearcoin/utils"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	port := r.URL.Query().Get("openPort")

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return ip != "" && port != ""
	}

	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)

	initPeer(conn, ip, port)
}

func AddPeer(address, port, openPort string) {
	url := fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port[1:], openPort)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	utils.HandleErr(err)
	initPeer(conn, address, port)
}
