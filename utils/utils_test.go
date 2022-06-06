package utils

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	s := struct{ Test string }{Test: "test"}

	t.Run("Should be same hash if interface is same", func(t *testing.T) {
		hash := "e005c1d727f7776a57a661d61a182816d8953c0432780beeae35e337830b1746"
		x := Hash(s)

		if x != hash {
			t.Errorf("Expected %s, got %s", hash, x)
		}
	})

	t.Run("Should be hash is hex encoded", func(t *testing.T) {
		x := Hash(s)
		_, err := hex.DecodeString(x)
		if err != nil {
			t.Errorf("Hash should be hex encoded")
		}
	})
}

func ExampleHash() {
	s := struct{ Test string }{Test: "test"}
	x := Hash(s)
	fmt.Println(x)
	// Output: e005c1d727f7776a57a661d61a182816d8953c0432780beeae35e337830b1746
}
