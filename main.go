package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/abelgideon/blog-aggregator/internal/config"
	"github.com/abelgideon/blog-aggregator/internal/database"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	getCommand map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	getComm, exists := c.getCommand[cmd.name]
	if !exists {
		return fmt.Errorf("command not found")
	}

	err := getComm(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.getCommand[name] = f
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

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Printf("failed to read file: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", conf.DBURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	myState := state{
		db:  dbQueries,
		cfg: &conf,
	}

	myCommands := commands{
		getCommand: make(map[string]func(*state, command) error),
	}

	myCommands.register("login", handlerLogin)
	myCommands.register("register", handlerRegister)
	myCommands.register("reset", handlerReset)
	myCommands.register("users", handlerUsers)

	userArgs := os.Args
	if len(userArgs) < 2 {
		fmt.Println("not enough arguments were provided")
		os.Exit(1)
	}

	cmd := command{
		name: userArgs[1],
		args: userArgs[2:],
	}

	runErr := myCommands.run(&myState, cmd)
	if runErr != nil {
		fmt.Println(runErr)
		os.Exit(1)
	}
}
