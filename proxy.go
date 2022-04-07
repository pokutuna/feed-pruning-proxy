package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"google.golang.org/appengine/v2/log"
)

func ProxyFeed(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	feed := q.Get("feed")
	parsed, err := url.Parse(feed)
	if feed == "" || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, http.StatusText(http.StatusBadRequest))
		return
	}

	// prevent redirection loop
	if r.Host == parsed.Host {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, http.StatusText(http.StatusBadRequest))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", feed, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, http.StatusText(http.StatusInternalServerError))
		log.Errorf(r.Context(), "Failed to create request to fetch a feed: %v", err.Error())
		return
	}
	req.Header.Set("User-Agent", fmt.Sprintf("slack-feed-proxy (%s)", r.Host))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, http.StatusText(http.StatusInternalServerError))
		log.Errorf(r.Context(), "Failed to request to fetch a feed: %v", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		w.WriteHeader(resp.StatusCode)
		fmt.Fprintln(w, http.StatusText(resp.StatusCode))
		log.Warningf(r.Context(), "Failed to fetch a feed: status=%d; url=%s", resp.StatusCode, feed)
		return
	}

	conf := TransformConfig{
		ProxyOrigin:   ServerOrigin(r.Host),
		Org:           q.Get("org"),
		Channel:       q.Get("channel"),
		UseRedirector: q.Get("org") != "" || q.Get("channel") != "",
	}

	ct, wt, err := Transform(resp.Body, conf)
	if err != nil {
		var code int
		var errFormat ErrUnExpectedFormat
		if errors.As(err, &errFormat) {
			code = http.StatusBadRequest
		} else {
			code = http.StatusInternalServerError
		}
		w.WriteHeader(code)
		fmt.Fprintln(w, http.StatusText(code))
		log.Warningf(r.Context(), "Failed to transfrom a feed: %v", err)
		return
	}

	w.Header().Set("Content-Type", ct)
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.WriteHeader(http.StatusOK)
	wt.WriteTo(w)
}

func ServerOrigin(host string) string {
	var scheme string
	if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") {
		scheme = "http"
	} else {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, host)
}
