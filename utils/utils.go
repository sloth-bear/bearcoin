package utils

import "log"

func HandleErr(error error) {
	if error != nil {
		log.Fatal(error)
	}
}
