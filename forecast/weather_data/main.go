package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/mediocregopher/radix.v2/redis"
)

var (
	location = flag.String("loc", "", "location to load data from")
	hist     = flag.Bool("hist", false, "wether to load historic data")
	key      = flag.String("key", "", "api-key from weather online")
	limit    = flag.Int("limit", 250, "limit of request (based on key)")
)

func main() {
	flag.Parse()

	client, err := redis.DialTimeout("tcp", "localhost:6379", 10*time.Second)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	requests, err := client.Cmd("INCR", *key).Int()
	if err != nil {
		panic(err)
	}
	log.Printf("Start %d. requests with given key\n", requests)
	if requests > *limit {
		client.Cmd("DECR", *key)
		log.Fatalf("Too many requests for your key")
	}

	fmt.Println("Do something")
	fmt.Println(url.QueryEscape("Köln"))

	r, err := GetWeather("Köln")
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}
