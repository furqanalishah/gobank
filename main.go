package main

import "log"

func main() {
	store, err := NewPostgresStore()
	store.Init()
	if err != nil {
		log.Fatal(err)
	}
	server := NewAPIServer(":3000", store)
	server.Run()
}
