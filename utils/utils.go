package utils

import (
	"bytes"
	"encoding/gob"
	"log"
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
