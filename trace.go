package main

import (
	"fmt"
	"net/http"
	"regexp"
)

var (
	traceCtxRe = regexp.MustCompile(`^(\w+)/(\d+)(?:;o=.+)?$`)
)

type Trace struct {
	TraceID string `json:"logging.googleapis.com/trace,omitempty"`
	SpanID  string `json:"logging.googleapis.com/spanId,omitempty"`
}

func GetTrace(req *http.Request, projectID string) Trace {
	if projectID == "" {
		return Trace{"", ""}
	}

	traceHeader := req.Header.Get("X-Cloud-Trace-Context")
	if traceHeader == "" {
		return Trace{"", ""}
	}
	matches := traceCtxRe.FindAllStringSubmatch(traceHeader, 2)
	if len(matches) < 1 || len(matches[0]) < 3 {
		return Trace{"", ""}
	}

	t := matches[0][1]
	s := matches[0][2]
	return Trace{
		TraceID: fmt.Sprintf("projects/%s/traces/%s", projectID, t),
		SpanID:  s,
	}
}
