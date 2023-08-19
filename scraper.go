package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Ayushlm10/rssAgg/internal/database"
	"github.com/google/uuid"
)

func scraper(concurrency int, timeBetweenRequests time.Duration, db *database.Queries) {
	log.Printf("Scraping %d feeds, at interval %v", concurrency, timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println(err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(feed, db, wg)
		}
		wg.Wait()
	}
}

func scrapeFeed(feed database.Feed, db *database.Queries, wg *sync.WaitGroup) {
	defer wg.Done()

	//fetch feed
	rssFeed, err := urlToRssFeed(feed.Url)
	if err != nil {
		log.Println("Error fetchign feed:", err)
		return
	}

	_, err = db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error Fetching feed:", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {

		//description is a null string.
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		// parse the date as it is a string
		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)

		if err != nil {
			log.Printf("Cant't parse date: %v", err)
			continue
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			FeedID:      feed.ID,
			PublishedAt: pubDate,
			Description: description,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("failed to create post: %v\n", err)
		}
	}
	log.Printf("Found %s collected , %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
