package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAccount(store Storage, fname, lname, pw string) *Account {
	acc, err := NewAccount(fname, lname, pw)

	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	fmt.Println("new account =>", acc.Number)

	return acc
}

func seedAccounts(s Storage) {
	seedAccount(s, "Nata", "Nael", "password280198")
}

func main() {
	// ./bin/gobank --seed to trigger the seeding function
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

	loadEnv()
	store, err := NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("seeding the database")
		seedAccounts(store)
		// seed account
	}

	log.Println("database connected")
	server := NewAPIServer(":3000", store)
	server.Run()
}
