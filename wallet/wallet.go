package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"io/fs"
	"math/big"
	"os"

	"github.com/sloth-bear/bearcoin/utils"
)

const (
	fileName string = "bearcoin.wallet"
)

type fileLayer interface {
	hasWalletFile() bool
	writeFile(name string, data []byte, perm fs.FileMode) error
	readFile(name string) ([]byte, error)
}

type layer struct{}

func (layer) hasWalletFile() bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func (layer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (layer) readFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

var files fileLayer = layer{}  // go는 interface를 암묵적으로 구현하기 때문에 layer struct에 명시해줄 필요가 없다.

type wallet struct {
	PrivateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

func createPrivateKey() *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return key
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = files.writeFile(fileName, bytes, 0644)
	utils.HandleErr(err)
}

func restoreKey() *ecdsa.PrivateKey {
	keyAsBytes, err := files.readFile(fileName)
	utils.HandleErr(err)
	key, err := x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return key
}

func encodeToHex(a, b []byte) string {
	bytes := append(a, b...)
	return fmt.Sprintf("%x", bytes)
}

func decodeFromHex(s string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return nil, nil, err
	}

	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]

	bigFirst, bigSecond := big.Int{}, big.Int{}
	bigFirst.SetBytes(firstHalfBytes)
	bigSecond.SetBytes(secondHalfBytes)

	return &bigFirst, &bigSecond, nil
}

func getAddress(key *ecdsa.PrivateKey) string {
	return encodeToHex(key.X.Bytes(), key.Y.Bytes())
}

func Sign(w *wallet, payload string) string {
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, payloadAsBytes)
	utils.HandleErr(err)

	return encodeToHex(r.Bytes(), s.Bytes())
}

func Verify(signature, payload, address string) bool {
	r, s, err := decodeFromHex(signature)
	utils.HandleErr(err)

	x, y, err := decodeFromHex(address)
	utils.HandleErr(err)

	publicKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	return ecdsa.Verify(&publicKey, payloadAsBytes, r, s)
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if files.hasWalletFile() {
			w.PrivateKey = restoreKey()
		} else {
			key := createPrivateKey()
			persistKey(key)
			w.PrivateKey = key
		}
		w.Address = getAddress(w.PrivateKey)

	}
	return w
}
