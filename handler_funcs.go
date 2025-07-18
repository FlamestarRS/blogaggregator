package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/FlamestarRS/blogaggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {

		return fmt.Errorf("usage: %s <username>", cmd.name)
	}
	username := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Println("error: user does not exist")
		os.Exit(1)
	}
	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Println("Username has been set: " + username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {

		return fmt.Errorf("usage: %s <username>", cmd.name)
	}
	username := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		fmt.Println("error: user already exists")
		os.Exit(1)
	}
	params := database.CreateUserParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}
	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Println("New user created:" + username)
	fmt.Println(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Reset users table successfully")
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		return err
	}

	currentUser := s.cfg.CurrentUserName
	for _, user := range users {
		if user.Name == currentUser {
			fmt.Println(user.Name + " (current)")
			continue
		}
		fmt.Println(user.Name)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <time_between_requests> ex: 10s, 1m, 1h", cmd.name)
	}
	timeBetweenReqs := cmd.args[0]
	fmt.Println("Collecting feeds every " + timeBetweenReqs)
	parsedTime, err := time.ParseDuration(timeBetweenReqs)
	if err != nil {
		return fmt.Errorf("error parsing time: %v", err)
	}
	ticker := time.NewTicker(parsedTime)
	counter := 0
	for ; ; <-ticker.C {
		counter += 1
		fmt.Printf("scraping feed... (%v)\n", counter)
		scrapeFeeds(s)
	}

}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %s <feed_name> <feed_url>", cmd.name)
	}
	name := cmd.args[0]
	url := cmd.args[1]

	params := database.CreateFeedParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}
	newFeed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	handlerFollow(s, command{args: []string{url}}, user)
	fmt.Println(newFeed)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return err
	}
	if len(feeds) == 0 {
		fmt.Println("No feeds found")
	}
	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("user id not found for feed: %s", feed.Name)
		}
		fmt.Printf("Name: %s\nURL: %s\nUserID: %v\n", feed.Name, feed.Url, user)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.name)
	}
	url := cmd.args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	follow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error: could not follow feed %v, err: %v", feed.Name, err)
	}
	fmt.Println("Feed: " + follow.FeedName + " User: " + follow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, feed := range following {
		feedInfo, err := s.db.GetFeedByID(context.Background(), feed.FeedID)
		if err != nil {
			return err
		}
		fmt.Println(feedInfo.Name)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.name)
	}
	url := cmd.args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}
	params := database.DeleteFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error unfollowing feed")
	}
	return nil
}

func handlerBrowse(s *state, cmd command) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("usage: %s <int> (optional)", cmd.name)
	}
	limit := 2
	if len(cmd.args) == 1 {
		var err error
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("error setting limit: %v", err)
		}
	}
	postList, err := s.db.GetPostsForUser(context.Background(), int32(limit))
	if err != nil {
		return fmt.Errorf("error getting posts: %v", err)
	}
	for _, post := range postList {
		feed, err := s.db.GetFeedByID(context.Background(), post.FeedID)
		if err != nil {
			return fmt.Errorf("error getting feed from feedID")
		}
		if post.Description.Valid {
			fmt.Printf("Feed: %s\nTitle: %s\nDescription: %s\nLink: %s\nPubDate: %.10v\n", feed.Name, post.Title, post.Description.String, post.Url, post.PublishedAt)
			continue
		}
		fmt.Printf("Feed: %s\nTitle: %s\nLink: %s\nPubDate: %.10v\n", feed.Name, post.Title, post.Url, post.PublishedAt)
	}
	return nil
}
