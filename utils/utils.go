package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"strings"
)

func HandleErr(error error) {
	if error != nil {
		log.Fatal(error)
	}
}

func ToBytes(i interface{}) []byte {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(i)
	HandleErr(err)

	return buffer.Bytes()
}

func FromBytes(i interface{}, b []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(b))
	HandleErr(decoder.Decode(i))
}

func Hash(i interface{}) string {
	interfaceAsStr := fmt.Sprint(i)
	hash := sha256.Sum256([]byte(interfaceAsStr))
	return fmt.Sprintf("%x", hash)
}

func Splitter(s string, sep string, i int) string {
	splitedStrs := strings.Split(s, sep)
	if len(splitedStrs)-1 < i {
		return ""
	}
	return splitedStrs[i]
}
