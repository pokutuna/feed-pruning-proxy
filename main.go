package main

import (
	"net/http"

	"google.golang.org/appengine/v2"
)

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/p", ProxyFeed)
	http.HandleFunc("/r", Redirect)

	appengine.Main()
}
