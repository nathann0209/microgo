package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

// Constructor
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

// Method
func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello World")
	// Retrieve data in body of request in form of data slice
	data, err := io.ReadAll(r.Body)
	// Basic error handling: Type 2 using http.Error method.
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}

	// Write response back to user
	fmt.Fprintf(rw, "Hello %s", data)
}
