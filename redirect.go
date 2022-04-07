package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/appengine/v2/log"
)

type RedirectionLog struct {
	LogType     string      `json:"type"`
	HttpRequest HttpRequest `json:"httpRequest"`
	URL         string      `json:"url"`
	Org         string      `json:"org"`
	Channel     string      `json:"channel"`
}

type HttpRequest struct {
	Method     string `json:"method"`
	UserAgent  string `json:"userAgent"`
	RemoteAddr string `json:"remoteAddr"`
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	url := q.Get("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, http.StatusText(http.StatusBadRequest))
		return
	}

	isSlackbot := strings.Contains(r.UserAgent(), "Slackbot")
	if !isSlackbot {
		l := RedirectionLog{
			LogType: "redirect",
			HttpRequest: HttpRequest{
				Method:     r.Method,
				UserAgent:  r.UserAgent(),
				RemoteAddr: r.RemoteAddr,
			},
			URL:     url,
			Org:     q.Get("org"),
			Channel: q.Get("channel"),
		}
		j, _ := json.Marshal(l)
		log.Infof(r.Context(), "%s", j)
	}

	http.Redirect(w, r, url, http.StatusFound)
}
