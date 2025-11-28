package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/iashyam/gator/internal/config"
	"github.com/iashyam/gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	cfg, _ := config.ReadConfig()
	dbUrl := cfg.Db_url
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Can't open the postgress databse check url %s", dbUrl)
	}
	dbQuiries := database.New(db)
	state := State{dbQuiries, &cfg}

	commands := Commands{
		Commandlist: make(map[string]func(*State, Command) error),
		CommandMap:  make(map[string]Command),
	}
	commandLogin := Command{
		Name:        "login",
		Description: "sets the users to given user",
		Args:        os.Args,
	}
	commandRegister := Command{
		Name:        "register",
		Description: "creates a new user and logs it in",
		Args:        os.Args,
	}
	commandReset := Command{
		Name:        "reset",
		Description: "deletes all the users from the database",
		Args:        os.Args,
	}
	commandListUser := Command{
		Name:        "users",
		Description: "prints all the users from the database",
		Args:        os.Args,
	}
	commandAgg := Command{
		Name:        "agg",
		Description: "fatches the rss feed of a url",
		Args:        os.Args,
	}
	commandAddFeed := Command{
		Name:        "addfeed",
		Description: "takes a name and url and adds it to the databse",
		Args:        os.Args,
	}
	commandListFeeds := Command{
		Name:        "feeds",
		Description: "lists all the feeds in the database",
		Args:        os.Args,
	}
	commandFollow := Command{
		Name:        "follow",
		Description: "follow a feed by url",
		Args:        os.Args,
	}
	commandListFollowing := Command{
		Name:        "following",
		Description: "lists all the feeds followed by the user",
		Args:        os.Args,
	}
	commandUnfollowFeed := Command{
		Name:        "unfollow",
		Description: "unfllows feed given in url",
		Args:        os.Args,
	}
	commandBrowse := Command{
		Name:        "browse",
		Description: "lists the posts from feeds followed by the user",
		Args:        os.Args,
	}
	// fmt.Println(os.Args)
	command := os.Args[1]

	commands.Register("login", commandLogin, HandlerLogin)
	commands.Register("register", commandRegister, HandlerRegistger)
	commands.Register("reset", commandReset, HandlerReset)
	commands.Register("users", commandListUser, HandlerListUsers)
	commands.Register("agg", commandAgg, HandlerAgg)
	commands.Register("addfeed", commandAddFeed, middleWareLoggedIn(HandlerAddFeed))
	commands.Register("feeds", commandListFeeds, middleWareLoggedIn(HandlerListFeeds))
	commands.Register("follow", commandFollow, middleWareLoggedIn(HandlerFollow))
	commands.Register("following", commandListFollowing, middleWareLoggedIn(HandlerListFollowing))
	commands.Register("unfollow", commandUnfollowFeed, middleWareLoggedIn(HandlerUnfollowFeed))
	commands.Register("browse", commandBrowse, middleWareLoggedIn(HandlerBrowse))

	commandToRun, ok := commands.CommandMap[command]
	if !ok {
		fmt.Printf("Unknown command %s\n", command)
		os.Exit(1)
	}
	err = commands.Run(&state, commandToRun)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

}
