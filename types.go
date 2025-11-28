package main

import (
	"fmt"

	"github.com/iashyam/gator/internal/config"
	"github.com/iashyam/gator/internal/database"
)

type Commands struct {
	Commandlist map[string]func(*State, Command) error
	CommandMap  map[string]Command
}

type Command struct {
	Name        string
	Description string
	Arguments  []string
	Args        []string
	c     Commands
}

type State struct {
	db     *database.Queries
	Config *config.Config
}


func (commands Commands) Run(state *State, command Command) error {

	handler, exists := commands.Commandlist[command.Name]
	if !exists {
		return fmt.Errorf("unregisterd command")
	}
	err := handler(state, command)
	if err != nil {
		return fmt.Errorf("error while executing command %s: %v", command.Name, err)
	}

	return nil
}

func (commands Commands) Register(name string, command Command, f func(*State, Command) error) error {

	_, ok := commands.Commandlist[name]
	if ok {
		return fmt.Errorf("command %s already exist", name)
	}

	commands.Commandlist[name] = f
	commands.CommandMap[name] = command

	return nil
}


type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}