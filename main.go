package main

import (
	"log"
)

func main() {
	store, err := newSqliteStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%+v\n", store)
	server := newApiServer(":3000", store)
	server.Run()
}
