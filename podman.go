package main

import (
	"flag"
	"fmt"
	"strings"
	//"github.com/krig/go-sox" //for playing podcasts
)

func main() {
	//make configurationg struct that holds default settings
	config := Configuration{"~/podman", "k", "j", "h", "l", " ", "/", make([]Podcast, 0)}
	//read command line flags first
	noTui := flag.Bool("no-tui", false, "Select whether to use the TUI or not")
	flag.Parse()
	//read config file
	config = readConfig(config)
	//made a decision to use TUI or not
	if *noTui == true {
		end := false
		for end != true {
			end = CliInterface(config)
		}
	} else {
		//TUI
		panic("unimplemented")
	}
}

func CliInterface(config Configuration) bool {
	command := ""
	fmt.Scanf("%s", &command)
	command = strings.ToLower(command)
	//handle empty string
	if command == "" {
		return false
	} else if command == "exit" {
		return true
	} else if command[0] == '/' {
		fmt.Printf("%s is your search\n", command[1:])
	} else if command == "list" {
		for _, entry := range config.Subscribed {
			//do nothing
			fmt.Println(entry.author)
		}
	} else if command == "help" {
		fmt.Println("Type list to list your subscriptions, /<string> to search, exit to exit, help to show this")
	} else if command == "settings" {
		fmt.Println(config)
	} else {
		fmt.Println("Type help for a list of commands")
	}
	return false
}
