package main

import (
	"log"

	"GO-API-template/src"
)

func main() {
	defer log.Println("Shutting down")
	src.Start()
}
