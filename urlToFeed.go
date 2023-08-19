package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RssFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RssItem `xml:"item"`
	} `xml:"channel"`
}

type RssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubdate"`
}

func urlToRssFeed(url string) (RssFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return RssFeed{}, err
	}
	defer resp.Body.Close()

	feed, err := io.ReadAll(resp.Body)
	if err != nil {
		return RssFeed{}, err
	}

	rssFeed := RssFeed{}
	err = xml.Unmarshal(feed, &rssFeed)
	if err != nil {
		return RssFeed{}, err
	}
	return rssFeed, nil
}
