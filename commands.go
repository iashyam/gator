package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iashyam/gator/internal/database"
)

func middleWareLoggedIn(handler func(state *State, cmd Command, user database.User) error) func(*State, Command) error {

	return func(state *State, command Command) error {

		ctx := context.Background()
		username := state.Config.Current_user_name
		user, err := state.db.GetUser(ctx, username)

		if err != nil {
			return fmt.Errorf("error fetching user from database %v", err)
		}

		return handler(state, command, user)

	}
}

func HandlerLogin(state *State, command Command) error {
	if len(command.Args) < 2 {
		return fmt.Errorf("not enough arguments provided")
	}
	if len(command.Args) != 3 {
		return fmt.Errorf("a username is required for login")
	}
	username := command.Args[2]
	ctx := context.Background()
	_, err := state.db.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("unregistered user %s", username)
	}
	err = state.Config.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting the user %v", err)
	}
	fmt.Printf("User set to %s successfully\n", username)
	return nil
}

func HandlerRegistger(state *State, command Command) error {
	if len(command.Args) < 2 {
		return fmt.Errorf("not enough arguments provided")
	}
	if len(command.Args) != 3 {
		return fmt.Errorf("a username is required for register")
	}
	username := command.Args[2]
	ctx := context.Background()
	newUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	user, err := state.db.CreateUser(ctx, newUserParams)
	if err != nil {
		return fmt.Errorf("user %s already exists", username)
	}

	state.Config.SetUser(user.Name)
	fmt.Printf("user %s was created, and logged in\n", user.Name)

	return nil

}

func HandlerReset(state *State, command Command) error {

	ctx := context.Background()
	err := state.db.Reset(ctx)
	if err != nil {
		return fmt.Errorf("can't reset the databse")
	}
	fmt.Println("user table reset successfully")
	return nil
}

func HandlerListUsers(state *State, command Command) error {

	ctx := context.Background()
	users, err := state.db.ListUsers(ctx)
	if err != nil {
		return fmt.Errorf("can't list the users")
	}

	for _, username := range users {
		if username == state.Config.Current_user_name {
			fmt.Printf("* %s (current)\n", username)
			continue
		}
		fmt.Printf("* %s\n", username)
	}

	return nil
}

func HandlerAgg(state *State, command Command) error {

	if len(command.Args) < 2 {
		return fmt.Errorf("not enough arguments provided")
	}
	if len(command.Args) != 3 {
		return fmt.Errorf("a time duartion is needed to run this command")
	}

	t := command.Args[2]
	fmt.Println("Scrapping feeds every", t)
	time_duration, err := time.ParseDuration(t)

	if err != nil {
		return fmt.Errorf("can't parse duarion %v", err)
	}

	ticker := time.NewTicker(time_duration)

	for ; ; <-ticker.C {
		err := ScrapeFeeds(state)
		if err != nil {
			return fmt.Errorf("error appread while scraping feed %v", err)
		}
		time.Sleep(time.Second)
	}

	return nil
}

func HandlerAddFeed(state *State, command Command, user database.User) error {
	if len(command.Args) < 3 {
		return fmt.Errorf("not enough arguments provided")
	}
	if len(command.Args) != 4 {
		return fmt.Errorf("a name and url is required for register")
	}

	ctx := context.Background()
	feedName := command.Args[2]
	url := command.Args[3]

	feed, err := state.db.AddFeed(ctx, database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       url,
		UserID:    user.ID,
	})

	if err != nil {
		return fmt.Errorf("error adding feed to database %v", err)
	}
	_, err = state.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})

	if err != nil {
		return fmt.Errorf("error adding feedfollow to database %v", err)
	}

	fmt.Printf("Feed %s added with url %s\n", feed.Name, feed.Url)
	return nil
}

func HandlerListFeeds(state *State, command Command, user database.User) error {

	ctx := context.Background()
	feeds, err := state.db.ListFeeds(ctx)

	if err != nil {
		return fmt.Errorf("error fetching the feeds from database %v", err)
	}

	for _, feed := range feeds {
		fmt.Printf("* Name: %s\n", feed.Name)
		fmt.Printf("	- Url: %s\n", feed.Url)

		creator, _ := state.db.GetUserNameByID(ctx, feed.UserID)
		fmt.Printf("	- Created by: %s\n", creator.Name)
	}

	return nil
}

func HandlerFollow(state *State, command Command, user database.User) error {
	ctx := context.Background()
	if len(command.Args) < 3 {
		return fmt.Errorf("not enough arguments provided")
	}
	if len(command.Args) != 3 {
		return fmt.Errorf("a is required for register")
	}

	feedURL := command.Args[2]
	feed, err := state.db.GetFeedByURL(ctx, feedURL)

	if err != nil {
		return fmt.Errorf("error fetching feed from database %v", err)
	}

	followfeed, err := state.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})

	if err != nil {
		return fmt.Errorf("error following feed %v", err)
	}

	fmt.Printf("User %s is now following feed %s \n", followfeed.UserName, followfeed.FeedName)

	return nil
}

func HandlerListFollowing(state *State, command Command, user database.User) error {

	ctx := context.Background()

	followFeeds, err := state.db.GetFeedFollowsForUser(ctx, user.ID)

	if err != nil {
		return err
	}

	fmt.Printf("Feed followed by %s\n", user)
	for _, feed := range followFeeds {
		thisFeed, _ := state.db.GetFeedByID(ctx, feed.FeedID)
		creator, _ := state.db.GetUserNameByID(ctx, thisFeed.UserID)
		fmt.Printf("	- %s by %s\n", thisFeed.Name, creator.Name)
	}

	return nil
}

func HandlerUnfollowFeed(state *State, command Command, user database.User) error {

	ctx := context.Background()

	if len(command.Args) < 2 {
		return fmt.Errorf("not enough arguments provided")
	}
	if len(command.Args) < 3 {
		return fmt.Errorf("specity a feed url to unfollow")
	}
	if len(command.Args) > 3 {
		return fmt.Errorf("what are you trying to do mate?")
	}

	feedURL := command.Args[2]
	feed, err := state.db.GetFeedByURL(ctx, feedURL)

	if err != nil {
		return fmt.Errorf("error getting feed from database %v", err)
	}

	err = state.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})

	if err != nil {
		return fmt.Errorf("error deleting connection from database %v", err)
	}

	fmt.Printf("Unfollowed %s\n", feed.Name)

	return nil
}

func HandlerBrowse(state *State, command Command, user database.User) error {

	if len(command.Args) < 2 {
		return fmt.Errorf("not enough arguments provided")
	}

	posts, err := state.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  5, //command.Args[2],
	})
	if err != nil {
		return fmt.Errorf("error fetching posts for user %v", err)
	}

	fmt.Printf("Posts for user %s\n", user.Name)
	for _, post := range posts {
		fmt.Printf(" - %s (%s)\n", post.Title, post.Url)
	}

	return nil
}

func HandlerHelp(state *State, command Command) error {
	fmt.Println("Available commands:")

	for _, cmd := range command.c.CommandMap {
		fmt.Printf(" - %s: %s\n", cmd.Name, cmd.Description)
		if len(cmd.Arguments) > 0 {
			fmt.Printf("   Arguments: %v\n", cmd.Arguments)
		}
	}
	return nil
}
