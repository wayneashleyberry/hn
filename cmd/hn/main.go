package main

import (
	"context"
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

const (
	defaultExpiration = 1 * time.Hour
	cleanupInterval   = 1 * time.Minute
)

func main() {
	var gobtype []hackernews.Item

	gob.Register(gobtype)

	c := cache.New(defaultExpiration, cleanupInterval)

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
				cc := cache.NewFrom(defaultExpiration, cleanupInterval, cachedItems)

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

	bg := context.Background()

	ids, err := hackernews.TopStories(bg)
	if err != nil {
		panic(err)
	}

	items := []hackernews.Item{}

	for n, id := range ids {
		if n > 2 {
			continue
		}

		ctx, cancel := context.WithTimeout(bg, time.Second*2)
		defer cancel()

		i, err := hackernews.GetItem(ctx, id)
		if err != nil {
			panic(err)
		}

		items = append(items, i)
	}

	render(items)

	err = c.Add("items", items, cache.DefaultExpiration)
	if err != nil {
		panic(err)
	}

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
