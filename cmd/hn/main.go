package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/patrickmn/go-cache"
	"github.com/wayneashleyberry/hn/pkg/hackernews"
	"github.com/wayneashleyberry/hn/pkg/hyperlink"
	"github.com/wayneashleyberry/truecolor/pkg/color"
)

func main() {
	var gobtype []hackernews.Item

	gob.Register(gobtype)

	c := cache.New(30*time.Minute, 1*time.Hour)

	usercache, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	appcachedir := path.Join(usercache, "hn-cli")

	if _, err := os.Stat(appcachedir); os.IsNotExist(err) {
		err = os.Mkdir(appcachedir, 0700)
		if err != nil {
			panic(err)
		}
	}

	cachefile := path.Join(appcachedir, "cache.gob")

	if _, err := os.Stat(cachefile); err == nil {
		f, err := os.Open(cachefile)
		if err == nil {
			dec := gob.NewDecoder(f)

			var cachedItems map[string]cache.Item
			err := dec.Decode(&cachedItems)
			if err == nil {
				cc := cache.NewFrom(30*time.Minute, 1*time.Hour, cachedItems)

				values, present := cc.Get("items")
				if present {
					items, ok := values.([]hackernews.Item)
					if ok && len(items) > 0 {
						render(items)
						return
					}
				}
			}
		}
	}

	ids, err := hackernews.TopStories()
	if err != nil {
		return
	}

	items := []hackernews.Item{}

	for n, id := range ids {
		if n > 2 {
			continue
		}

		i, err := hackernews.GetItem(id)
		if err != nil {
			return
		}

		items = append(items, i)
	}

	render(items)

	c.Add("items", items, cache.DefaultExpiration)

	f, err := os.Create(cachefile)
	if err != nil {
		panic(err)
	}

	enc := gob.NewEncoder(f)

	err = enc.Encode(c.Items())
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}
}

func render(items []hackernews.Item) {
	for n, item := range items {
		tm := time.Unix(item.Time, 0)

		fmt.Printf("%d. ", n+1)
		err := hyperlink.Write(os.Stdout, item.URL, item.Title)
		if err != nil {
			continue
		}

		fmt.Print("\n")
		color.White().Dim().Printf("   %d points by %s %s\n", item.Score, item.By, humanize.Time(tm))
	}
}
