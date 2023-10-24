package main

import (
	"log"
)

func main() {
	loadEnv()
	store, err := NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	log.Println("database connected")
	server := NewAPIServer(":3000", store)
	server.Run()
}
