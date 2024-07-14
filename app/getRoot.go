package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func getRoot(tmpl map[string]*template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s\n", r.Method, r.URL.Path)

		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		err := tmpl["index.html"].ExecuteTemplate(w, "base.html", nil)
		if err != nil {
			fmt.Printf("error executing template: %s\n", err)
		}
	}
}
