package hackernews

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TopStories returns a slice of story ID's
func TopStories() ([]int, error) {
	url := "https://hacker-news.firebaseio.com/v0/topstories.json"
	client := &http.Client{
		Timeout: time.Second * 1,
	}
	r, err := client.Get(url)
	if err != nil {
		return []int{}, err
	}

	type ResponseBody []int

	var responseBody ResponseBody

	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&responseBody)
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
func GetItem(id int) (Item, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	client := &http.Client{
		Timeout: time.Second * 1,
	}
	r, err := client.Get(url)
	if err != nil {
		return Item{}, err
	}

	var responseBody Item

	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&responseBody)
	if err != nil {
		return Item{}, err
	}

	return responseBody, nil
}
