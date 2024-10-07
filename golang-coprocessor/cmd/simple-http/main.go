package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	coprocessor "github.com/apollographql/coprocessor/internal"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 3007
	}

	http.HandleFunc("/", coprocessor.RequestHandler)
	log.Printf("Starting on :%v", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
