package blockchain

import (
	"errors"
	"time"

	"github.com/sloth-bear/bearcoin/utils"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

type TxIn struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{Owner: "COINBASE", Amount: minerReward},
	}
	txOuts := []*TxOut{
		{Owner: address, Amount: minerReward},
	}

	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

// 트랜잭션에는 입력값과 출력값이 필요하다.
// 누군가 얼마나 가지고 있는지 알고 싶다면 트랜잭션 출력값을 참고한다.
// 트랜잭션을 시작하고 싶다면 입력값을 생성해야 한다.
// 그리고 입력값은 blockchain에서 가지고 있는 돈이어야 한다.
func makeTx(from, to string, amount int) (*Tx, error) {
	if (Blockchain().BalanceByAddress(from)) < amount {
		return nil, errors.New("not enough money")
	}
	var txIns []*TxIn
	var txOuts []*TxOut

	total := 0
	oldTxOuts := Blockchain().TxOutsByAddress(from)
	for _, oldTxOut := range oldTxOuts {
		if total > amount {
			break
		}
		txIn := &TxIn{oldTxOut.Owner, oldTxOut.Amount}
		txIns = append(txIns, txIn)
		total += txIn.Amount
	}

	change := total - amount
	if change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}

	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("slothbear", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbaseTx := makeCoinbaseTx("slothbear")
	txs := m.Txs
	txs = append(txs, coinbaseTx)
	m.Txs = nil
	return txs
}