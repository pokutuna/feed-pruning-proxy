package main

import (
	"fmt"
	"os"

	"github.com/beevik/etree"
)

func main() {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile("feed.rss"); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", doc.Root().Tag)

	items := doc.FindElements("//rss/channel/item")
	for _, i := range items {
		if d := i.SelectElement("description"); d != nil {
			i.RemoveChild(d)
		}
	}

	doc.WriteTo(os.Stdout)
}
