package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/ehafenmaier/boot-dev-gator/internal/database"
	"github.com/google/uuid"
	"html"
	"io"
	"net/http"
	"strings"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}

	// Set the request headers
	req.Header.Add("User-Agent", "gator")

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	feed := &RSSFeed{}
	err = xml.Unmarshal(body, feed)
	if err != nil {
		return nil, err
	}

	// Unescape the HTML entities in the feed
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	// Close the response body
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}

	// Parse the time between requests
	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("error parsing time between requests: %v", err)
	}

	// Create a ticker to scrape feeds
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	// Get the next feed to scrape
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("error getting next feed to fetch: %v\n", err)
		return
	}

	// Mark the feed as fetched
	dbParams := database.MarkFeedFetchedParams{
		ID:            feed.ID,
		UpdatedAt:     time.Now(),
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	feed, err = s.db.MarkFeedFetched(context.Background(), dbParams)
	if err != nil {
		fmt.Printf("error marking feed fetched: %v\n", err)
		return
	}

	// Fetch the RSS feed
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		fmt.Printf("error fetching feed: %v\n", err)
		return
	}

	// Iterate over the feed items and save them to the database
	for _, item := range rssFeed.Channel.Item {
		// Parse the published date
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Printf("error parsing published date: %v\n", err)
			continue
		}

		// Create a new post
		dbParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: item.Title, Valid: true},
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: publishedAt, Valid: true},
			FeedID:      feed.ID,
		}

		_, err = s.db.CreatePost(context.Background(), dbParams)
		if err != nil {
			// Skip duplicate key errors
			if strings.Contains(err.Error(), "duplicate key value") {
				continue
			}

			// Log other errors
			fmt.Printf("error creating post: %v\n", err)
		}
	}
}
