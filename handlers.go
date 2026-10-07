package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/abelgideon/blog-aggregator/internal/database"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username is required")
	}
	_, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("something went wrong: %v", err)
	}

	setUserErr := s.cfg.SetUser(cmd.args[0])
	if setUserErr != nil {
		return setUserErr
	}

	fmt.Printf("%v has been created and set as current user\n", cmd.args[0])
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username is required")
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("something went wrong: %v", err)
	}

	setUserErr := s.cfg.SetUser(cmd.args[0])
	if setUserErr != nil {
		return setUserErr
	}

	fmt.Printf("%v has been set as the current user\n", cmd.args[0])
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("reset does not take arguments")
	}

	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("something went wrong: %v", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("users does not take arguments")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("something went wrong: %v", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("agg does not take arguments")
	}

	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(rssFeed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("feeds does not take arguments")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feeds: %v", err)
	}

	for _, feed := range feeds {
		fmt.Printf("Name: %v\n", feed.Name)
		fmt.Printf("URL: %v\n", feed.Url)

		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error fetching user: %v", err)
		}
		fmt.Printf("Created by: %v\n", user.Name)

		fmt.Println()
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("unexpected amount of arguments")
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %v", err)
	}

	_, feedFollowErr := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if feedFollowErr != nil {
		return fmt.Errorf("failed to follow feed: %v", feedFollowErr)
	}

	fmt.Println("ID:", feed.ID)
	fmt.Println("Name:", feed.Name)
	fmt.Println("URL:", feed.Url)
	fmt.Println("User ID:", feed.UserID)
	fmt.Println("Created At:", feed.CreatedAt)
	fmt.Println("Updated At:", feed.UpdatedAt)

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("unexpected number of arguments")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error following feed: %v", err)
	}

	fmt.Printf("%v just followed %v\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("following takes no arguments")
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error fetching follows: %v", err)
	}

	for _, follow := range follows {
		fmt.Println("*", follow.FeedName)
	}

	return nil
}
