package main

import (
	"log"

	"Kaomoji-DB/src"
)

func main() {
	defer log.Println("Shutting down")
	src.Start()
}
