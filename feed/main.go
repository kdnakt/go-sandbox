package main

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func main() {
	fmt.Println("start!")
	p := gofeed.NewParser()
	// feed, err := p.ParseURL("https://kdnakt.hatenablog.com/rss")
	feed, err := p.ParseURL("https://crieit.net/users/dala00/feed")
	if err != nil {
		fmt.Printf("Error %v", err)
	}
	fmt.Println(feed.Title)
	fmt.Println("end!")
}
