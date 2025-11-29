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
	cfg, cfgError := config.ReadConfig()

	if cfgError != nil {
		fmt.Println("Error reading config file. Make sure you have created a config file using the instructions in the README.md")
		log.Fatalf("Error reading config file: %v", cfgError)
	}

	dbUrl := cfg.Db_url
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println("Make sure that postgres is installed and running and the url is correct!")
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
		Arguments:   []string{"username"},
		Description: "sets the users to given user",
		Args:        os.Args,
	}
	commandRegister := Command{
		Name:        "register",
		Arguments:   []string{"username"},
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
		Arguments:   []string{"duration"},
		Description: "aggregates feeds from all users for a given duration",
		Args:        os.Args,
	}
	commandAddFeed := Command{
		Name:        "addfeed",
		Arguments:   []string{"name", "url"},
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
		Arguments:   []string{"url"},
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
		Arguments:   []string{"url"},
		Description: "unfllows feed given in url",
		Args:        os.Args,
	}
	commandBrowse := Command{
		Name:        "browse",
		Description: "lists the posts from feeds followed by the user",
		Args:        os.Args,
	}
	commandVersion := Command{
		Name:        "version",
		Description: "Provides help information for commands",
	}
	commandHelp := Command{
		Name:        "help",
		Description: "Provides help information for commands",
	}

	// fmt.Println(os.Args)

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
	commands.Register("version", commandVersion, HandlerVersion)
	commands.Register("help", commandHelp, HandlerHelp)

	if len(os.Args) < 2 || os.Args[1] == "help" {
		fmt.Println("Welcomt to Gator RSS Feed Reader!")
		commandHelp.c = commands
		err = commands.Run(&state, commandHelp)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		os.Exit(0)
	}
	command := os.Args[1]

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
