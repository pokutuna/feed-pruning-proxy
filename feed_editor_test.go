package main

import (
	"strings"
	"testing"

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
)

func TestRSSFeedEditor(t *testing.T) {
	rssfeed := `
    <?xml version="1.0"?>
    <rss version="2.0">
        <channel>
            <title>ぽ靴な缶</title>
            <link>https://blog.pokutuna.com/</link>
            <description></description>
            <lastBuildDate>Fri, 11 Mar 2022 11:00:00 +0900</lastBuildDate>
            <docs>http://blogs.law.harvard.edu/tech/rss</docs>
            <generator>Hatena::Blog</generator>
            <item>
                <title>Title 1</title>
                <link>https://blog.pokutuna.com/entry/1</link>
                <description>Hello This is a &lt;a href=&quot;https://example.test&quot;&gt;link&lt;/a&gt;</description>
                <pubDate>Fri, 11 Mar 2022 11:00:00 +0900</pubDate>
                <guid isPermalink="false">hatenablog://entry/13574176438071665569</guid>
            </item>
            <item>
                <title>Title 2</title>
                <link>https://blog.pokutuna.com/entry/2</link>
                <description><![CDATA[Hello This is a <a href="https://example.test">link</a>]]></description>
                <pubDate>Thu, 02 Dec 2021 22:00:00 +0900</pubDate>
                <guid isPermalink="false">hatenablog://entry/13574176438038846201</guid>
            </item>
            <item>
                <title>Title 3</title>
                <link>https://blog.pokutuna.com/entry/3</link>
                <content:encoded><![CDATA[Hello This is a <a href="https://example.test">link</a>]]></content:encoded>
                <pubDate>Thu, 03 Dec 2021 03:00:00 +0900</pubDate>
                <guid isPermalink="false">hatenablog://entry/3</guid>
            </item>
        </channel>
    </rss>
    `

	editor := RSSFeedEditor{}
	var doc *etree.Document

	reload := func() {
		doc = etree.NewDocument()
		_, err := doc.ReadFrom(strings.NewReader(rssfeed))
		assert.NoError(t, err)
	}

	t.Run("UpdateFeedTitle", func(t *testing.T) {
		reload()

		editor.UpdateFeedTitle(doc)

		title := doc.FindElement("/rss/channel/title").Text()
		assert.Equal(t, "ぽ靴な缶 (with feed-pruning-proxy)", title)
	})

	t.Run("RemoveEntryContent", func(t *testing.T) {
		t.Run("description", func(t *testing.T) {
			reload()
			item := doc.FindElement("/rss/channel/item")
			assert.NotNil(t, item.SelectElement("description"))

			editor.RemoveEntryContent(doc)
			assert.Nil(t, item.SelectElement("description"))
		})

		t.Run("content:encoded", func(t *testing.T) {
			reload()
			item := doc.FindElements("/rss/channel/item")[2]
			assert.NotNil(t, item.SelectElement("content:encoded"))

			editor.RemoveEntryContent(doc)
			assert.Nil(t, item.SelectElement("content:encoded"))
		})
	})

	t.Run("DietEntryContent", func(t *testing.T) {
		reload()

		items := doc.FindElements("/rss/channel/item")
		item1 := items[0] // description with Escaped
		item2 := items[1] // description with CDATA
		item3 := items[2] // content:encoded with CDATA

		editor.DietEntryContent(doc)

		assert.Equal(t, "Hello This is a link", item1.SelectElement("description").Text())
		assert.Equal(t, "Hello This is a link", item2.SelectElement("description").Text())
		assert.Equal(t, "Hello This is a link", item3.SelectElement("content:encoded").Text())
	})

	t.Run("TapRedirector", func(t *testing.T) {
	})
}

func TestAtomFeedEditor(t *testing.T) {
	var atomfeed = `
    <feed xmlns="http://www.w3.org/2005/Atom" xml:lang="ja">
        <title>ぽ靴な缶</title>
        <link href="https://blog.pokutuna.com/"/>
        <updated>2022-03-11T11:00:00+09:00</updated>
        <author>
            <name>pokutuna</name>
        </author>
        <generator uri="https://blog.hatena.ne.jp/" version="7f2239af4151fedf9100f66511debb">Hatena::Blog</generator>
        <id>hatenablog://blog/12704591929886254776</id>
        <entry>
            <title>Title 1</title>
            <link href="https://blog.pokutuna.com/entry/1"/>
            <id>hatenablog://entry/13574176438071665569</id>
            <published>2022-03-11T11:00:00+09:00</published>
            <updated>2022-03-11T11:00:03+09:00</updated>
            <summary type="html">Hello This is a &lt;a href=&quot;https://example.test&quot;&gt;link&lt;/a&gt; in summary</summary>
            <content type="html">Hello This is a &lt;a href=&quot;https://example.test&quot;&gt;link&lt;/a&gt; in content</content>
            <author>
                <name>pokutuna</name>
            </author>
        </entry>
        <entry>
            <title>Title 2</title>
            <link href="https://blog.pokutuna.com/entry/2"/>
            <id>hatenablog://entry/13574176438038846201</id>
            <published>2021-12-02T22:00:00+09:00</published>
            <updated>2021-12-10T14:10:08+09:00</updated>
            <summary type="html"><![CDATA[Hello This is a <a href="https://example.test">link</a> in summary]]></summary>
            <content type="html"><![CDATA[Hello This is a <a href="https://example.test">link</a> in content]]></content>
            <author>
                <name>pokutuna</name>
            </author>
        </entry>
    </feed>
    `

	editor := AtomFeedEditor{}
	var doc *etree.Document

	reload := func() {
		doc = etree.NewDocument()
		_, err := doc.ReadFrom(strings.NewReader(atomfeed))
		assert.NoError(t, err)
	}

	t.Run("UpdateFeedTitle", func(t *testing.T) {
		reload()

		editor.UpdateFeedTitle(doc)

		title := doc.FindElement("/feed/title").Text()
		assert.Equal(t, "ぽ靴な缶 (with feed-pruning-proxy)", title)
	})

	t.Run("RemoveEntryContent", func(t *testing.T) {
		reload()

		item := doc.FindElement("/feed/entry")
		assert.NotNil(t, item.SelectElement("summary"))
		assert.NotNil(t, item.SelectElement("content"))

		editor.RemoveEntryContent(doc)
		assert.Nil(t, item.SelectElement("summary"))
		assert.Nil(t, item.SelectElement("content"))
	})

	t.Run("DietEntryContent", func(t *testing.T) {
		reload()

		items := doc.FindElements("/feed/entry")
		item1 := items[0] // Escaped
		item2 := items[1] // CDATA

		editor.DietEntryContent(doc)

		assert.Equal(t, "Hello This is a link in summary", item1.SelectElement("summary").Text())
		assert.Equal(t, "Hello This is a link in content", item1.SelectElement("content").Text())
		assert.Equal(t, "Hello This is a link in summary", item2.SelectElement("summary").Text())
		assert.Equal(t, "Hello This is a link in content", item2.SelectElement("content").Text())
	})

	t.Run("TapRedirector", func(t *testing.T) {
	})
}
