package hackernews

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// TopStories returns a slice of story ID's
func TopStories(ctx context.Context) ([]int, error) {
	url := "https://hacker-news.firebaseio.com/v0/topstories.json"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []int{}, err
	}

	req = req.WithContext(ctx)

	client := http.DefaultClient

	res, err := client.Do(req)
	if err != nil {
		return []int{}, err
	}

	type ResponseBody []int

	var responseBody ResponseBody

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&responseBody)
	if err != nil {
		return []int{}, err
	}

	return responseBody, nil
}

// Item implementation
type Item struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Score       int    `json:"score"`
	Time        int64  `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

// GetItem will fetch an item by ID
func GetItem(ctx context.Context, id int) (Item, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Item{}, err
	}

	req = req.WithContext(ctx)

	client := http.DefaultClient

	res, err := client.Do(req)
	if err != nil {
		return Item{}, err
	}

	var responseBody Item

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&responseBody)
	if err != nil {
		return Item{}, err
	}

	return responseBody, nil
}
