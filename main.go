package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Hello World")
		// Retrieve data in body of request in form of data slice
		data, err := io.ReadAll(r.Body)

		// // Basic error handling: Type 1
		// if err != nil {
		// 	rw.WriteHeader(http.StatusBadRequest)
		// 	rw.Write([]byte("Oops"))
		// }

		// Basic error handling: Type 2 using http.Error method.
		if err != nil {
			http.Error(rw, "Oops", http.StatusBadRequest)
			return
		}

		// Print the retrieved data
		log.Printf("Data: %s\n", data)

		// Write response back to user
		fmt.Fprintf(rw, "Hello %s", data)
	})

	http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
		log.Println("Goodbye World")
	})

	http.ListenAndServe(":9090", nil)
}
