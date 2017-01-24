//
// Go Example - Form Submit
//
// HTML form submission example.
//
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// If no file exists redirect to form.
		_, err := os.Stat("./data.json")
		if err != nil {
			http.Redirect(w, r, "/form", 303)
		} else {
			// Open the data.json file.
			b, err := ioutil.ReadFile("./data.json")
			if err != nil {
				log.Println("Could not read file.")
			}

			// Marshal json into data struct.
			data := new(FormData)
			err = json.Unmarshal(b, data)
			if err != nil {
				log.Println("Could not unmarshal json.", err)
			}

			templates := []string{
				"./templates/layout.html",
				"./templates/home.html",
			}

			t, err := template.ParseFiles(templates...)
			if err != nil {
				log.Println("Could not parse templates.", err)
			}

			t.Execute(w, data)
		}
	})

	r.Mount("/form", Content{}.Routes())

	http.ListenAndServe(":3000", r)
}
