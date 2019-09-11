package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/wayneashleyberry/hn/pkg/hackernews"
	"github.com/wayneashleyberry/hn/pkg/hyperlink"
	"github.com/wayneashleyberry/truecolor/pkg/color"
)

func main() {
	stories, err := hackernews.TopStories()
	if err != nil {
		return
	}

	for n, id := range stories {
		if n > 2 {
			return
		}

		s, err := hackernews.GetItem(id)
		if err != nil {
			return
		}

		tm := time.Unix(s.Time, 0)

		fmt.Printf("%d. ", n+1)
		hyperlink.Write(os.Stdout, s.URL, s.Title)
		fmt.Print("\n")
		color.White().Dim().Printf("   %d points by %s %s\n", s.Score, s.By, humanize.Time(tm))
	}
}
