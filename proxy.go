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
	isSelf := r.Host == parsed.Host // prevent redirection loop
	if feed == "" || err != nil || isSelf {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", feed, nil)
	req.Header.Set("User-Agent", fmt.Sprintf("feed-pruning-proxy (%s)", r.Host))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, context.DeadlineExceeded) {
			code = http.StatusRequestTimeout
		}
		http.Error(w, http.StatusText(code), code)
		log.Errorf(r.Context(), "failed to request to fetch a feed: url=%s; err=%v", feed, err)
		return
	}
	defer resp.Body.Close()

	if 400 <= resp.StatusCode {
		http.Error(w, http.StatusText(resp.StatusCode), resp.StatusCode)
		log.Warningf(r.Context(), "failed to fetch a feed: url=%s; status=%d", feed, resp.StatusCode)
		return
	}

	_, isDietMode := q["diet"]
	conf := TransformConfig{
		ProxyOrigin:   ServerOrigin(r.Host),
		Org:           q.Get("org"),
		Channel:       q.Get("channel"),
		UseRedirector: q.Get("org") != "" || q.Get("channel") != "",
		DietMode:      isDietMode,
	}

	wt, err := Transform(resp.Body, conf)
	if err != nil {
		var code int
		var msg string

		var errParse ErrXMLParseFailed
		var errFormat ErrUnExpectedFormat
		if errors.As(err, &errParse) || errors.As(err, &errFormat) {
			code = http.StatusBadRequest
			msg = err.Error()
		} else {
			code = http.StatusInternalServerError
		}

		http.Error(w, fmt.Sprintf("%s\n%s\n", http.StatusText(code), msg), code)
		log.Warningf(r.Context(), "Failed to transfrom a feed: url=%s; err=%v", feed, err)
		return
	}

	w.Header().Set("Content-Type", "application/xml") // or respect the original content-type?
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.WriteHeader(http.StatusOK)
	_, err = wt.WriteTo(w)
	if err != nil {
		log.Errorf(r.Context(), err.Error())
	}
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
