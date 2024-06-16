package handlers

import (
	"fmt"
	"log"
	"net/http"
)

type GoodBye struct {
	l *log.Logger
}

func NewGoodbye(l *log.Logger) *GoodBye {
	return &GoodBye{l}
}

func (g *GoodBye) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	html := `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Class Notes</title>
    </head>
    <body>
        <h1>Class Notes</h1>
        <ul>
            <li><a href="https://drive.google.com/file/d/17aT3sAagh3VW9ByShlsPBk1GI3_Nejp5/view?usp=sharing" target="_blank">Note: Example Note</a></li>
        </ul>
    </body>
    </html>`

	fmt.Fprint(rw, html)

	rw.Write([]byte("Bye!"))
}
