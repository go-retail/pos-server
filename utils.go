package main

import (
	"fmt"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", err, msg)
		panic(fmt.Sprintf("%s: %s", err, msg))
	}
}
