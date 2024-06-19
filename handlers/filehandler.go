package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/russross/blackfriday/v2"
)

type FileHandler struct {
	l *log.Logger
}

func NewFileHandler(l *log.Logger) *FileHandler {
	return &FileHandler{l}
}

func (fh *FileHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var filepath string
	switch r.URL.Path {
	case "/pointers":
		filepath = "static/pointers.md"
	default:
		http.NotFound(rw, r)
		return
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		http.Error(rw, "File not found", 404)
		return
	}

	htmlContent := blackfriday.Run(content)

	rw.Header().Set("Content-Type", "text/html")
	rw.Write(htmlContent)
}
