package blockchain

import (
	"errors"
	"sync"
	"time"

	"github.com/sloth-bear/bearcoin/utils"
	"github.com/sloth-bear/bearcoin/wallet"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs map[string]*Tx
	m   sync.Mutex
}

var m *mempool = &mempool{}
var memOnce sync.Once

func Mempool() *mempool {
	memOnce.Do(func() {
		m = &mempool{
			Txs: make(map[string]*Tx),
		}
	})
	return m
}

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct {
	TxID      string `json:"txID"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

func (t *Tx) generateId() {
	t.ID = utils.Hash(t)
}

func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(wallet.Wallet(), t.ID)
	}
}

func validate(t *Tx) bool {
	isValid := true

	for _, txIn := range t.TxIns {
		refferedTx := FindTx(Blockchain(), txIn.TxID)
		if refferedTx == nil {
			isValid = false
			break
		}

		address := refferedTx.TxOuts[txIn.Index].Address
		isValid = wallet.Verify(txIn.Signature, t.ID, address)
		if !isValid {
			break
		}
	}

	return isValid
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false

MempoolLoop:
	for _, tx := range Mempool().Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break MempoolLoop
			}
		}
	}

	return exists
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{Address: address, Amount: minerReward},
	}

	tx := Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.generateId()
	return &tx
}

var ErrorNoMoney = errors.New("Not enough money")
var ErrorNotValid = errors.New("Invalid Tx")

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(Blockchain(), from) < amount {
		return nil, ErrorNoMoney
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0

	uTxOuts := UTxOutsByAddress(Blockchain(), from)
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		// TODO from(owner)을 직접 넣는 것은 보안에 취약하다. 왜? UxOut에서 존재하는 코인만 가져올 텐데.
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}

	if change := total - amount; change > 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}

	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)

	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}

	tx.generateId()
	tx.sign()

	isValid := validate(tx)
	if !isValid {
		return nil, ErrorNotValid
	}
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) (*Tx, error) {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return nil, err
	}
	m.Txs[tx.ID] = tx
	return tx, nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbaseTx := makeCoinbaseTx(wallet.Wallet().Address)

	var txs []*Tx

	for _, mTx := range m.Txs {
		txs = append(txs, mTx)
	}

	txs = append(txs, coinbaseTx)
	m.Txs = make(map[string]*Tx)
	return txs
}

func (m *mempool) AddPeerTx(tx *Tx) {
	m.m.Lock()
	defer m.m.Unlock()

	m.Txs[tx.ID] = tx
}
