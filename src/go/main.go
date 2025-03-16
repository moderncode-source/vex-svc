package main

import (
	"io"
	"log"
	"net/http"
)

const PORT = ":8080"

func main() {
	rootHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello world!\n")
	}

	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
