package main

import (
	"log"

	"kaomojidb/src"
)

func main() {
	defer log.Println("Shutting down")
	src.Start()
}
