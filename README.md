# Gator

```
   ____      _             
  / ___| ___| |_ ___  _ __ 
 | |  _ / __| __/ _ \| '__|
 | |_| | (__| || (_) | |   
  \____|\___|\__\___/|_|    

```

Gator is a CLI tool to explore, follow, and fetch RSS feeds directly from your terminal.
It's lightweight, simple, and a great project to understand how real-world Go applications work.

---

## Introduction

Gator lets you:

- Register and log in users  
- Follow RSS feeds and store them in a Postgres database  
- Fetch and list new posts from your followed feeds  
- Read post summaries and jump to full posts  

The app is small, hackable, and easy to extend.

---

## Installation

Install using Go:

```bash
go install github.com/iashyam/gator@latest
```

Check if it works:

```bash
gator version
```

---

## Requirements

- Go
- PostgreSQL (running locally)

Install postgres using the following commands

ON linux

```bash
sudo apt update && sudo apt install postgresql postgresql-contrib
```

Set a password:

```bash
sudo passwd postgres
```

On macOS

```bash
brew install postgres
```

On Windows Users

If you're on Windows… consider installing Linux. You’re missing out.

### Verify installation

```bash
psql --version
```

---

## Configuration

Create a `.gatorconfig.json` file in your **home directory**:

```json
{"db_url": "postgres://username:password@host:port/database?sslmode=disable"}
```

This tells Gator where your Postgres database lives.

---

## Usage

Running `gator` without commands will show all available options.

Here is the full command list:

| Command    | Description                                         | Arguments        | Example                          |
|------------|-----------------------------------------------------|------------------|----------------------------------|
| register   | Creates a new user and logs them in                 | username         | gator register shyam             |
| reset      | Deletes all users from the database                 | —                | gator reset                      |
| users      | Prints all users from the database                  | —                | gator users                      |
| feeds      | Lists all the feeds in the database                 | —                | gator feeds                      |
| follow     | Follows a feed by URL                               | url              | gator follow https://abc.com/rss |
| following  | Lists all feeds followed by the current user        | —                | gator following                  |
| login      | Sets the current logged-in user                     | username         | gator login shyam                |
| agg        | Aggregates feeds from all users for a duration      | duration         | gator agg 1h                     |
| addfeed    | Adds a new feed to the database                     | name, url        | gator addfeed tech https://t.com |
| unfollow   | Unfollows a feed by URL                             | url              | gator unfollow https://t.com     |
| browse     | Lists posts from feeds followed by the user         | —                | gator browse                     |
| help       | Shows help information for all commands             | —                | gator help                       |

---

## Project Highlights

- Full CLI tool written in Go  
- Database interaction using `goose` + `sqlc`  
- RSS/XML parsing  
- Clean modular structure  
- Published Go module  

---

## Future Work

- Add a REST API  
- Add more tests (TDD)  
- Provide a Docker install method  
- GitHub CLI automation / publishing  
- More polished output & UX  

---

## Version

Check your installed version:

```bash
gator version
```
---

## Acknowledgments

Special thanks to:
- **[Boot.dev](https://boot.dev)** - For providing excellent Go learning resources and project ideas.
- **[Lane Wagner](https://twitter.com/wagslane)** - For mentoring, guidance, and creating the educational platform that inspired this project.

This project was built as part of the Boot.dev Go curriculum and serves as a practical demonstration of Go programming concepts.