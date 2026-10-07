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
