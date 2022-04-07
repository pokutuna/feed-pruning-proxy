package main

var rssfeed = `
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
            <title>3/17 Born Digital Summit 2022 で発表します</title>
            <link>https://blog.pokutuna.com/entry/born-digital-summit-2022</link>
            <description>Long Long Content</description>
            <pubDate>Fri, 11 Mar 2022 11:00:00 +0900</pubDate>
            <guid isPermalink="false">hatenablog://entry/13574176438071665569</guid>
            <category>GCP</category>
            <category>宣伝</category>
            <enclosure url="https://cdn-ak.f.st-hatena.com/images/fotolife/p/pokutuna/20220311/20220311040604.png" type="image/png" length="0" />
        </item>
        <item>
            <title>Google Colaboratory でデータフローのドキュメントを書く試み</title>
            <link>https://blog.pokutuna.com/entry/dataflow-doc-on-colab</link>
            <description>Long Long Content</description>
            <pubDate>Thu, 02 Dec 2021 22:00:00 +0900</pubDate>
            <guid isPermalink="false">hatenablog://entry/13574176438038846201</guid>
            <category>データ</category>
            <category>ドキュメンテーション</category>
            <category>ツール</category>
            <enclosure url="https://cdn-ak.f.st-hatena.com/images/fotolife/p/pokutuna/20211202/20211202203702.png" type="image/png" length="0" />
        </item>
    </channel>
</rss>
`

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
        <title>3/17 Born Digital Summit 2022 で発表します</title>
        <link href="https://blog.pokutuna.com/entry/born-digital-summit-2022"/>
        <id>hatenablog://entry/13574176438071665569</id>
        <published>2022-03-11T11:00:00+09:00</published>
        <updated>2022-03-11T11:00:03+09:00</updated>
        <summary type="html">Long Long Summary</summary>
        <content type="html">Long Long Content</content>
        <category term="GCP" label="GCP" />
        <category term="宣伝" label="宣伝" />
        <link rel="enclosure" href="https://cdn-ak.f.st-hatena.com/images/fotolife/p/pokutuna/20220311/20220311040604.png" type="image/png" length="0" />
        <author>
            <name>pokutuna</name>
        </author>
    </entry>
    <entry>
        <title>Google Colaboratory でデータフローのドキュメントを書く試み</title>
        <link href="https://blog.pokutuna.com/entry/dataflow-doc-on-colab"/>
        <id>hatenablog://entry/13574176438038846201</id>
        <published>2021-12-02T22:00:00+09:00</published>
        <updated>2021-12-10T14:10:08+09:00</updated>
        <summary type="html">Long Long Summary</summary>
        <content type="html">Long Long Content</content>
        <category term="データ" label="データ" />
        <category term="ドキュメンテーション" label="ドキュメンテーション" />
        <category term="ツール" label="ツール" />
        <link rel="enclosure" href="https://cdn-ak.f.st-hatena.com/images/fotolife/p/pokutuna/20211202/20211202203702.png" type="image/png" length="0" />
        <author>
            <name>pokutuna</name>
        </author>
    </entry>
</feed>
`
