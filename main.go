package main

import (
	"net/http"
	"text/template"

	"google.golang.org/appengine/v2"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.Must(template.ParseFiles("index.html"))
	t.Execute(w, map[string]string{
		"host": r.Host,
	})
}

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/p", ProxyFeed)
	http.HandleFunc("/r", Redirect)
	appengine.Main()
}
