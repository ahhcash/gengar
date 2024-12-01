package main

import (
	"flag"
	"log"
)

func main() {
	port := flag.Int("port", 50051, "the server port")
	flag.Parse()

	log.Printf("🏃‍♂️‍➡️ starting server on port %d", *port)
	if err := Run(*port); err != nil {
		log.Fatalf("☠️ server error: %v", err)
	}
}
