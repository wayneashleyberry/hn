package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/wayneashleyberry/hn/pkg/hyperlink"
	"github.com/wayneashleyberry/truecolor/pkg/color"
)

func main() {
	url := "https://hacker-news.firebaseio.com/v0/topstories.json"
	client := &http.Client{
		Timeout: time.Second * 1,
	}
	r, err := client.Get(url)
	if err != nil {
		return
	}

	type ResponseBody []int

	var responseBody ResponseBody

	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&responseBody)
	if err != nil {
		return
	}

	for n, id := range responseBody {
		if n > 2 {
			return
		}

		s, err := getStory(id)
		if err != nil {
			return
		}

		tm := time.Unix(s.Time, 0)

		hyperlink.Write(os.Stdout, s.URL, s.Title)
		fmt.Print("\n")
		color.White().Dim().Printf("%d points by %s %s\n", s.Score, s.By, humanize.Time(tm))
	}
}

type story struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Score       int    `json:"score"`
	Time        int64  `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

func getStory(id int) (story, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	client := &http.Client{
		Timeout: time.Second * 1,
	}
	r, err := client.Get(url)
	if err != nil {
		return story{}, err
	}

	var responseBody story

	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&responseBody)
	if err != nil {
		return story{}, err
	}

	return responseBody, nil
}
