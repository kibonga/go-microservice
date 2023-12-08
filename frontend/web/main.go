package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	fmt.Println("Starting front end service on port 8092")
	err := http.ListenAndServe(":8092", nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, t string) {

	partials := []string{
		"./frontend/web/templates/base.layout.gohtml",
		"./frontend/web/templates/header.partial.gohtml",
		"./frontend/web/templates/footer.partial.gohtml",
	}
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	fmt.Println("current dir", dir)

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("./frontend/web/templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
