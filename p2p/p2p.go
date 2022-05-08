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
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	port := r.URL.Query().Get("openPort")

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return ip != "" && port != ""
	}

	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)

	peer := initPeer(conn, ip, port)

	time.Sleep(20 * time.Second)
	peer.inbox <- []byte("Hello from Port 3000!")
}

func AddPeer(address, port, openPort string) {
	url := fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port[0:], openPort)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	utils.HandleErr(err)
	
	peer := initPeer(conn, address, port)

	time.Sleep(10 * time.Second)
	peer.inbox <- []byte("Hello from 4000!")
}
