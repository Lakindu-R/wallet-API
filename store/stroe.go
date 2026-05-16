package store

import (
	"log"

	"github.com/cockroachdb/pebble"
)

var DB *pebble.DB

func InitDB() {
	var err error
	DB, err = pebble.Open("wallet.db", &pebble.Options{})

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Pebble DB Connected")
}
