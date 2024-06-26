package main

import (
	"log"
	"net/http"
)

func main() {
	RegisterHandlers()
	log.Println("Listening server on: 8080")

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
