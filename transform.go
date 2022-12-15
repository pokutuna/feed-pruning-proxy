package main

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/beevik/etree"
)

type ErrXMLParseFailed error

func NewErrXMLParseFailed(err error) ErrXMLParseFailed {
	return ErrXMLParseFailed(fmt.Errorf("failed to parse: %w", err))
}

type ErrUnExpectedFormat string

func (e ErrUnExpectedFormat) Error() string {
	return fmt.Sprintf("unexpected feed format: %s", string(e))
}

type TransformConfig struct {
	ProxyOrigin   string
	Org           string
	Channel       string
	UseRedirector bool
	DietMode      bool
}

func Transform(feed io.Reader, conf TransformConfig) (io.WriterTo, error) {
	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(feed); err != nil {
		return nil, NewErrXMLParseFailed(err)
	}
	if doc.Root() == nil {
		return nil, NewErrXMLParseFailed(errors.New("document has no root element"))
	}

	var editor FeedEditor
	format := strings.ToLower(doc.Root().Tag)
	if format == "rss" {
		editor = RSSFeedEditor{}
	} else if format == "feed" {
		editor = AtomFeedEditor{}
	} else if format == "rdf" {
		editor = RDFFeedEditor{}
	} else {
		return nil, ErrUnExpectedFormat(format)
	}

	editor.UpdateFeedTitle(doc)

	if conf.DietMode {
		editor.DietEntryContent(doc)
	} else {
		editor.RemoveEntryContent(doc)
	}

	if conf.UseRedirector {
		editor.TapRedirector(doc, conf)
	}

	return doc, nil
}
