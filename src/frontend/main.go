package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const basePath = "./templates/"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	fmt.Println("starting frontend service on port 5021")
	err := http.ListenAndServe(":5021", nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, t string) {

	partials := []string{
		fmt.Sprintf("%sbase.layout.gohtml", basePath),
		fmt.Sprintf("%sheader.partial.gohtml", basePath),
		fmt.Sprintf("%sfooter.partial.gohtml", basePath),
	}

	var templates []string
	templates = append(templates, fmt.Sprintf("%s%s", basePath, t))

	for _, p := range partials {
		templates = append(templates, p)
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
