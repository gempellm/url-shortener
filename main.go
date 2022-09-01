package main

import (
	"fmt"
	"hash/adler32"
	"html/template"
	"io/ioutil"
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

	return fmt.Sprintf("%x", hash)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) > 1 {
		curdir, _ := os.Getwd()
		filename := curdir + "/data/" + r.URL.Path[1:] + ".txt"
		URL, err := ioutil.ReadFile(filename)
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
	filename := "data/" + hash + ".txt"
	err := ioutil.WriteFile(filename, []byte(URL), 0600)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "saved", &Payload{URL: "http://91.203.192.110:8080/" + hash})
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Payload) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
