package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/FlamestarRS/blogaggregator/internal/database"
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
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(body, &rssFeed)
	if err != nil {
		return nil, err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}
	return &rssFeed, nil
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feed: %v", err)
	}
	params := database.MarkFeedFetchedParams{
		ID:        nextFeed.ID,
		UpdatedAt: time.Now(),
	}
	err = s.db.MarkFeedFetched(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error marking feed: %v", err)
	}
	fetchedFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}
	fmt.Printf("Feed: %s\nFetched %v posts:\n", fetchedFeed.Channel.Title, len(fetchedFeed.Channel.Item))
	newPostCounter := 0
	newPosts := []string{}
	for _, item := range fetchedFeed.Channel.Item {
		const timeFormat = "Mon, 02 Jan 2006 15:04:05 -0700"
		pubDate, _ := time.Parse(timeFormat, item.PubDate)

		hasDesc := true
		if strings.HasPrefix(item.Description, "<") {
			hasDesc = false
		}

		paramsCreatePost := database.CreatePostParams{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  hasDesc,
			},
			PublishedAt: pubDate,
			FeedID:      nextFeed.ID,
		}
		_, err = s.db.CreatePost(context.Background(), paramsCreatePost)
		if err != nil {
			continue
		}
		newPostCounter += 1
		newPosts = append(newPosts, item.Title)

	}
	fmt.Printf("Found %v new posts:\n\n", newPostCounter)
	for _, postTitle := range newPosts {
		fmt.Println(postTitle)
	}

	return nil
}
