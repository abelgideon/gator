# Gator Blog Aggregator

## Prerequisites
- go
- postgres
- config file at ~/.gatorconfig.json that stores the database URL and the current user in the format:

  `{"db_url": "your_database_connection_string", "current_user_name": "your_user_name"}`

## Installation
`go install github.com/abelgideon/gator`

## Commands
- `gator login <username>`: log in one of the registered users

- `gator register <username>`: register a new user

- `gator reset`: delete all users from database

- `gator users`: list all users

- `gator agg <time>`: aggregate posts from feeds you follow for a given period of time

- `gator addfeed <url>`: add a new feed

- `gator feeds`: list all feeds

- `gator follow <url>`: follow a feed

- `gator following`: list feeds you follow

- `gator unfollow <url>`: unfollow a feed

- `gator browse [limit]`: browse limit amount of posts from your feeds
