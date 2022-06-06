// Package utils contains function to be used across the apllication.
package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func HandleErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// ToBytes takes an interface and then will return encoding the bytes of the interface.
func ToBytes(i interface{}) []byte {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(i)
	HandleErr(err)

	return buffer.Bytes()
}

// FromBytes takes an interface and data as bytes and then will encode the data to the interface.
func FromBytes(i interface{}, b []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(b))
	HandleErr(decoder.Decode(i))
}

// Hash takes an interface and hashes and returns the hex encoding of the hash.
func Hash(i interface{}) string {
	interfaceAsStr := fmt.Sprint(i)
	hash := sha256.Sum256([]byte(interfaceAsStr))
	return fmt.Sprintf("%x", hash)
}

// Splitter takes an source text as string, seperator for source text, and index for splitted text array, and then returns splitted text.
func Splitter(s string, sep string, i int) string {
	splitedStrs := strings.Split(s, sep)
	if len(splitedStrs)-1 < i {
		return ""
	}
	return splitedStrs[i]
}

// JsonToBytes takes an interface and then will return JSON encoding of interface.
func JsonToBytes(i interface{}) []byte {
	r, err := json.Marshal(i)
	HandleErr(err)
	return r
}
