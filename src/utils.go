package main

import (
	"fmt"
	"log"
)

// HandleError is a common helper function that checks
// if error has occurred, if it does then check if the
// error context was fatal, then decide if the program
// should be halted or not. Returns boolean if error
// has occurred.
func HandleError(err error, fatal bool) bool {
	if err != nil {
		fmt.Println(err.Error())
		if fatal {
			log.Fatalln(err)
		} else {
			log.Println(err)
		}
		return true
	}
	return false
}
