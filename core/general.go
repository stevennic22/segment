package core

import (
	"log"
)

func CheckError(err error, hardStop bool) bool {
	if err != nil {
		if hardStop {
			log.Fatalf("Fatal error within application: %s", err.Error())
		}
		log.Printf("Recoverable error within application: %s", err.Error())
		return (true)
	}

	return (false)
}
