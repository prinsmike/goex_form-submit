package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/schema"
	"github.com/pressly/chi"
)

type Content struct {
	Text string
}

func (c Content) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", c.Form)      // GET /form - display the form.
	r.Post("/", c.PostForm) // POST /form - handle the form.

	return r
}

func (c Content) Form(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["Title"] = "Go Form Submission Test"

	templates := []string{
		"./templates/layout.html",
		"./templates/form.html",
	}

	t, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println("Could not parse templates:", err)
	}

	t.Execute(w, data)
}

type FormData struct {
	Title string `schema:"title"` // <input type="text" name="title" />
	Body  string `schema:"body"`  // <textarea name="body"></textarea>
}

func (c Content) PostForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Could not parse form data:", err)
	}

	data := new(FormData)

	decoder := schema.NewDecoder()
	err = decoder.Decode(data, r.PostForm)
	if err != nil {
		log.Println("Could not decode post form:", err)
	}

	// Write data to file.
	b, err := json.Marshal(data)
	if err != nil {
		log.Println("Could not marshal data to json:", err)
	}

	ioutil.WriteFile("./data.json", b, 0644)

	http.Redirect(w, r, "/", 303)
}
