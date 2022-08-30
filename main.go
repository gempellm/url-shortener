package main

import (
	"fmt"
	"hash/adler32"
	"html/template"
	"log"
	"net/http"
	"os"
)

var templates = template.Must(template.ParseFiles("index.html", "saved.html"))

type Payload struct {
	URL string
}

func main() {
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/saved/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func hash(s string) string {
	h := adler32.New()
	h.Write([]byte(s))
	hash := h.Sum32()

	return fmt.Sprint(hash)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) != 1 {
		filename := r.URL.Path[1:] + ".txt"
		URL, err := os.ReadFile(filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, string(URL), http.StatusSeeOther)
		return
	}
	var pl Payload
	renderTemplate(w, "index", &pl)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	URL := r.FormValue("body")
	hash := hash(URL)
	filename := hash + ".txt"
	err := os.WriteFile(filename, []byte(URL), 0600)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "saved", &Payload{URL: "http://localhost:8080/" + hash})
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Payload) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
