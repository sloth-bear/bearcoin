package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"testing"
)

const (
	testKey       string = "307702010104206e5bc5e8b8ade34ac456038001b71d79df5f86ab01619afb33951ba110479044a00a06082a8648ce3d030107a14403420004214490f5193ef33115a8a0e387f53f4b608b7d20b5b6f7ad31f9dc1cdd22c0f3f273d188e954148130b0887d0b05e90dd8054dd808cb15580cfae358feecb4a6"
	testPayload   string = "0039f51fbc5fcb0bd1cbd154ae9f25a9261a3e478b160e78c08b813ee054fbc5"
	testSigniture string = "7d320d82942cf0ca1bef7dfa7b8325d643283788da6e7c0157b5ded63d3aeedf504fd83a546aff659a330b7b0d6961996c8f3614d788eee5ac00a4f34109f417"
)

func makeTestWallet() *wallet {
	b, _ := hex.DecodeString(testKey)
	key, _ := x509.ParseECPrivateKey(b)

	w := *Wallet()
	w.PrivateKey = key
	w.Address = getAddress(key)

	return &w
}

func TestSign(t *testing.T) {
	s := Sign(makeTestWallet(), testPayload)
	_, err := hex.DecodeString(s)
	if err != nil {
		t.Errorf("Sign() should return a hex encoded string, but got %s", s)
	}
}

func TestVerify(t *testing.T) {
	type test struct {
		input string
		ok    bool
	}

	tests := []test{
		{input: testPayload, ok: true},
		{input: "1039f51fbc5fcb0bd1cbd154ae9f25a9261a3e478b160e78c08b813ee054fbc5", ok: false},
	}

	for _, tc := range tests {
		w := makeTestWallet()
		ok := Verify(testSigniture, tc.input, w.Address)

		if ok != tc.ok {
			t.Errorf("Verify() could not verify test signature and payload")
		}
	}
}

func TestDecodeFromHex(t *testing.T) {
	_, _, err := decodeFromHex("InvalidHex")
	if err == nil {
		t.Error("decodeFromHex should return error when payload is not hex.")
	}
}
