package main

import (
	"fmt"
	"io"
	"net/url"

	"github.com/beevik/etree"
)

type ErrUnExpectedFormat string

func (e ErrUnExpectedFormat) Error() string {
	return fmt.Sprintf("unexpected format: %s", string(e))
}

type TransformConfig struct {
	ProxyHost     string
	Org           string
	Channel       string
	UseRedirector bool
}

func Transform(feed io.Reader, conf TransformConfig) (string, io.WriterTo, error) {

	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(feed); err != nil {
		return "", nil, err
	}

	var editor FeedEditor
	var contentType string
	format := doc.Root().Tag
	if format == "rss" {
		editor = RSSFeedEditor{}
		contentType = "application/rss+xml; charset=utf-8"
	} else if format == "feed" {
		editor = AtomFeedEditor{}
		contentType = "application/atom+xml; charset=utf-8"
	} else {
		return "", nil, ErrUnExpectedFormat(format)
	}

	editor.UpdateFeedTitle(doc)
	editor.RemoveEntryContent(doc)

	if conf.UseRedirector {
		editor.TapRedirector(doc, conf)
	}

	return contentType, doc, nil
}

type FeedEditor interface {
	UpdateFeedTitle(doc *etree.Document)
	RemoveEntryContent(doc *etree.Document)
	TapRedirector(doc *etree.Document, conf TransformConfig)
}

func addTitleNotice(title string) string {
	return title + " (with slack-feed-proxy)"
}

func tapRedirector(link string, conf TransformConfig) string {
	u, _ := url.Parse(fmt.Sprintf("https://%s", conf.ProxyHost))
	u.Path = "/r"

	q := u.Query()
	q.Set("url", link)
	if conf.Org != "" {
		q.Set("org", conf.Org)
	}
	if conf.Channel != "" {
		q.Set("channel", conf.Channel)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

type RSSFeedEditor struct{}

func (e RSSFeedEditor) UpdateFeedTitle(doc *etree.Document) {
	if i := doc.FindElement("/rss/channel/title"); i != nil {
		i.SetText(addTitleNotice(i.Text()))
	}
}

func (e RSSFeedEditor) RemoveEntryContent(doc *etree.Document) {
	for _, i := range doc.FindElements("/rss/channel/item") {
		if d := i.SelectElement("description"); d != nil {
			i.RemoveChild(d)
		}
	}
}

func (e RSSFeedEditor) TapRedirector(doc *etree.Document, conf TransformConfig) {
	for _, i := range doc.FindElements("/rss/channel/item") {
		if d := i.SelectElement("link"); d != nil {
			i.SetText(tapRedirector(d.Text(), conf))
		}
	}
}

type AtomFeedEditor struct{}

func (e AtomFeedEditor) UpdateFeedTitle(doc *etree.Document) {
	if i := doc.FindElement("/feed/title"); i != nil {
		i.SetText(addTitleNotice(i.Text()))
	}
}

func (e AtomFeedEditor) RemoveEntryContent(doc *etree.Document) {
	for _, i := range doc.FindElements("/feed/entry") {
		if d := i.SelectElement("summary"); d != nil {
			i.RemoveChild(d)
		}
		if d := i.SelectElement("content"); d != nil {
			i.RemoveChild(d)
		}
	}
}

func (e AtomFeedEditor) TapRedirector(doc *etree.Document, conf TransformConfig) {
	for _, i := range doc.FindElements("/feed/entry") {
		if d := i.SelectElement("link"); d != nil {
			i.SetText(tapRedirector(d.SelectAttrValue("href", ""), conf))
		}
	}
}
