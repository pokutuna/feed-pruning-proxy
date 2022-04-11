package main

import (
	"fmt"
	"net/url"

	"github.com/beevik/etree"
	"github.com/microcosm-cc/bluemonday"
)

type FeedEditor interface {
	// Add "(with slack-feed-proxy)" notice to feed title
	UpdateFeedTitle(doc *etree.Document)

	// Remove entry contents
	// Slack expands not only the page link but also links contained in these. Noisy!
	RemoveEntryContent(doc *etree.Document)

	// Remove links in entry contents
	DietEntryContent(doc *etree.Document)

	// Replace entry links with redirector
	// The flood of feeds no one reading flushes out human conversations. Stop Robots Empire!
	TapRedirector(doc *etree.Document, conf TransformConfig)
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

func (e RSSFeedEditor) DietEntryContent(doc *etree.Document) {
	for _, i := range doc.FindElements("/rss/channel/item") {
		if d := i.SelectElement("description"); d != nil {
			d.SetText(stripTags(d.Text()))
		}
	}
}

func (e RSSFeedEditor) TapRedirector(doc *etree.Document, conf TransformConfig) {
	for _, i := range doc.FindElements("/rss/channel/item") {
		if d := i.SelectElement("link"); d != nil {
			d.SetText(tapRedirector(d.Text(), conf))
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

func (e AtomFeedEditor) DietEntryContent(doc *etree.Document) {
	for _, i := range doc.FindElements("/feed/entry") {
		if d := i.SelectElement("summary"); d != nil {
			d.SetText(stripTags(d.Text()))
		}
		if d := i.SelectElement("content"); d != nil {
			d.SetText(stripTags(d.Text()))
		}
	}
}

func (e AtomFeedEditor) TapRedirector(doc *etree.Document, conf TransformConfig) {
	for _, i := range doc.FindElements("/feed/entry") {
		if d := i.SelectElement("link"); d != nil {
			u := d.SelectAttrValue("href", "")
			d.RemoveAttr("href")
			d.CreateAttr("href", tapRedirector(u, conf))
		}
	}
}

func addTitleNotice(title string) string {
	return fmt.Sprintf("%s (with slack-feed-proxy)", title)
}

func tapRedirector(link string, conf TransformConfig) string {
	u, _ := url.Parse(conf.ProxyOrigin)
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

func stripTags(text string) string {
	p := bluemonday.StrictPolicy()
	return p.Sanitize(text)
}
