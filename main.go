package main

import (
	"log"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err.Error())
	}
	store.Init()
	server := NewAPIServer(":4444", store)
	server.Run()
}
