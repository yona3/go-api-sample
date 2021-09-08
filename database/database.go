package database

import (
	"log"

	"github.com/yona3/go-api-sample/ent"

	_ "github.com/lib/pq"
)

var client *ent.Client

func Init() {
	c, err := ent.Open("postgres", "host=localhost port=5433 user=root dbname=go-api-sample password=password sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	client = c
}

func GetClient() *ent.Client {
	return client
}

func CloseClient() error {
	log.Println("Closing client")
	return client.Close()
}
