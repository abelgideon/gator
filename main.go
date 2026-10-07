package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/abelgideon/blog-aggregator/internal/config"
	"github.com/abelgideon/blog-aggregator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
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
	myCommands.register("agg", handlerAgg)

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
