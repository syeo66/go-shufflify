package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"text/template"
)

func main() {
	http.HandleFunc("/", getRoot)

	fmt.Printf("starting server port 3333\n")
	fmt.Printf("open http://localhost:3333/\n")
	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Printf("error parsing template: %s\n", err)
	}

	err = t.Execute(w, "")
	if err != nil {
		fmt.Printf("error executing template: %s\n", err)
	}
}
