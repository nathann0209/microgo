package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello World")
	tmplPath := filepath.Join("templates", "home.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(rw, "Error parsing template", http.StatusInternalServerError)
		return
	}
	data := struct{ Name string }{Name: "Temple Tester"}
	err = tmpl.Execute(rw, data)
	if err != nil {
		http.Error(rw, "Error executing template", http.StatusInternalServerError)
	}
}
