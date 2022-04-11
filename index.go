package main

import (
	"html/template"
	"net/http"
	"net/url"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		feedURL, err := generateFeedURLWithProxy(r)
		if err == nil {
			http.Redirect(w, r, feedURL, http.StatusFound)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.Must(template.ParseFiles("index.html"))
	t.Execute(w, map[string]string{
		"origin": ServerOrigin(r.Host),
	})
}

func generateFeedURLWithProxy(r *http.Request) (string, error) {
	u, _ := url.Parse(ServerOrigin(r.Host))
	u.Path = "/p"

	q := u.Query()
	setIfExist := func(k string) {
		if r.FormValue(k) != "" {
			q.Set(k, r.FormValue(k))
		}
	}
	setIfExist("feed")
	setIfExist("org")
	setIfExist("channel")

	// add a query param without value
	encoded := q.Encode()
	if r.FormValue("diet") != "" {
		encoded = encoded + "&diet"
	}

	u.RawQuery = encoded

	return u.String(), nil
}
