package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Amarnaath548/rssagg/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetwenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetwenRequest)
	ticker := time.NewTicker(timeBetwenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFeatch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching foods:", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error making feed as fetched:", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description:= pgtype.Text{}
		if item.Description!=""{
			description.String=item.Description
			description.Valid=true
		}

		pubAt, err := time.Parse(time.RFC1123Z,item.PubDate)
		if err != nil {
			log.Printf("couldn't parse date %v with err %v", item.PubDate, err)
			continue
		}

		_, err = db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:        uuid.New(),
				CreatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
				UpdatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
				Title:     item.Title,
				Description: description,
				PublishedAt: pgtype.Timestamp{Time: pubAt, Valid: true},
				Url: item.Link,
				FeedID: feed.ID,
			})
			if err != nil {
				if strings.Contains(err.Error(),"duplicate key") {
					continue
				}
				log.Println("failed to create post:",err)
			}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
