package main

import (
	//"encoding/json" //for decoding search data from itunes
	"flag"
	"fmt"
	"strings"
	//"github.com/krig/go-sox" //for playing podcasts
)

type PodcastEntry struct {
	author string
}

func main() {
	//read command line flags first
	guiSelection := flag.Bool("no-tui", false, "Select whether to use the TUI or not")
	flag.Parse()
	//made a decision to use TUI or not
	if *guiSelection == true {
		end := false
		for end != true {
			end = CliInterface()
		}
	} else {
		panic("unimplemented")
	}
}

func CliInterface() bool {
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
		subs := getSubscribed()
		for _, entry := range subs {
			//do nothing
			fmt.Println(entry.author)
		}
	} else if command == "help" {
		fmt.Println("Type list to list your subscriptions, /<string> to search, exit to exit, help to show this")
	} else {
		fmt.Println("Type help for a list of commands")
	}
	return false
}

func getSubscribed() []PodcastEntry {
	return make([]PodcastEntry, 5)
}
