package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sloth-bear/bearcoin/blockchain"
	"github.com/sloth-bear/bearcoin/utils"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	port := r.URL.Query().Get("openPort")

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return ip != "" && port != ""
	}

	fmt.Printf("%s wants an upgrade\n", port)

	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)

	initPeer(conn, ip, port)
}

func AddPeer(address, port, openPort string, broadcast bool) {
	fmt.Printf("%s wants to connect to port %s\n", openPort, port)
	url := fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	utils.HandleErr(err)

	p := initPeer(conn, address, port)

	if broadcast {
		broadcastNewPeer(p)
		return
	}
	sendNewestBlock(p)
}

func BroadcastNewBlock(b *blockchain.Block) {
	for _, p := range Peers.v {
		notifyNewBlock(b, p)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	for _, p := range Peers.v {
		notifyNewTx(tx, p)
	}
}

func broadcastNewPeer(newPeer *peer) {
	for key, p := range Peers.v {
		if key != newPeer.key {
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			notifyNewPeer(payload, p)
		}
	}
}
