package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iashyam/gator/internal/database"
)

func FatchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {

	client := &http.Client{}

	request, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error making a request err: %v", err)
	}

	// this is just said to be a best practice
	request.Header.Set("User-Agent", "gator")

	response, err := client.Do(request)

	if err != nil {
		return nil, fmt.Errorf("error making a request err %v", err)
	}

	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("error making reading the response err %v", err)
	}

	var RssFeedObj RSSFeed

	err = xml.Unmarshal(bytes, &RssFeedObj)

	if err != nil {
		fmt.Printf("Reponse datas %s \n", string(bytes))
		return nil, fmt.Errorf("error unmarhalling the response %v", err)
	}

	/// there are some unconventional values to be taken care of
	RssFeedObj.Channel.Title = html.UnescapeString(RssFeedObj.Channel.Title)
	RssFeedObj.Channel.Description = html.UnescapeString(RssFeedObj.Channel.Description)

	return &RssFeedObj, nil
}

func ScrapeFeeds(state *State) error {
	feed, err := state.db.GetLastFetched(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feeds %v", err)
	}

	ctx := context.Background()
	rssFeed, err := FatchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed from url %v", err)
	}

	err = state.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		ID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("error fetching feed from url %v", err)
	}

	fmt.Printf("Fetched the feed: %s\n", rssFeed.Channel.Title)

	err = savePostsToDB(state, feed, rssFeed)
	if err != nil {
		return fmt.Errorf("error saving posts to database %v", err)
	}

	return nil
}

func savePostsToDB(state *State, feed database.Feed, rssFeed *RSSFeed) error {
	ctx := context.Background()
	for _, item := range rssFeed.Channel.Item {
		inLayout := "Mon, 02 Jan 2006 15:04:05 -0700"
		publishedAt, err := time.Parse(inLayout, item.PubDate)
		if err != nil {
			fmt.Printf("this error %v\n", err)
			continue
		}
		_, err = state.db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			FeedID:      feed.ID,
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: publishedAt,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			return fmt.Errorf("error saving post to database %v", err)
		}
	}
	return nil
}
