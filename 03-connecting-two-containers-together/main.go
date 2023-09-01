package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var (
	ctx      = context.Background()
	dbClient *redis.Client
	key      = "pv"
)

func init() {
	dbClient = redis.NewClient(&redis.Options{
		Addr: "db:6379",
	})
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Ping from %s", r.RemoteAddr)

	// Get the visitor count from the Redis datastore.
	pageView, err := dbClient.Get(ctx, key).Int()
	if err != nil {
		if err == redis.Nil {
			pageView = 1
		} else {
			panic(err)
		}
	}

	// Increment the number of visitors.
	pageView++

	// Save the visitor count number to the Redis datastore.
	err = dbClient.Set(ctx, key, pageView, 0).Err()
	if err != nil {
		panic(err)
	}

	// Print the number of visitors to the ResponseWriter.
	fmt.Fprintf(w, "Welcome! You are visitor #%d.\n", pageView)
}
