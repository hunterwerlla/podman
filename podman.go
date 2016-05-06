package main

import (
	"flag"
	"fmt"
	"strconv"
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
	//write config on sucessful exit
	defer writeConfig(config)
	//made a decision to use TUI or not
	if *noTui == true {
		end := false
		for end != true {
			config, end = CliInterface(config)
		}
	} else {
		//TUI
		panic("unimplemented")
	}
}

func CliInterface(config Configuration) (Configuration, bool) {
	command := ""
	fmt.Scanf("%s", &command)
	command = strings.ToLower(command)
	if command == "" {
		return config, false
	}
	if command == "exit" {
		return config, true
	} else if command[0] == '/' {
		fmt.Printf("%s is your search, sub to subscribe, exit to exit search, and \n", command[1:])
		results, err := searchItunes(string(command[1:]))
		if err != nil {
			fmt.Printf("Error searching Itunes, %s\n", err.Error())
		} else {
			for {
				fmt.Scanf("%s", &command)
				if command == "exit" {
					//exit this loop
					break
				} else if command == "sub" {
					//sub to podcast
					//load in the number
					fmt.Scanf("%s", &command)
					num, err := strconv.Atoi(command)
					if err != nil {
						fmt.Println("error converting to int")
						break
					}
					for i, _ := range results {
						if i == num {
							fmt.Println("appending the result to subscribed")
							config.Subscribed = append(config.Subscribed, results[i])
							goto searchEnd //yada yada considered harmful: this is much easier to jump out of these nested for loops
						}
					}
					fmt.Println("Number is in wrong format or too large, try again")
				} else {
					fmt.Println("Please input either exit or sub <number>")
				}
			}
		searchEnd:
		}
	} else if command == "list" {
		for _, entry := range config.Subscribed {
			//do nothing
			fmt.Println(entry.ArtistName)
		}
	} else if command == "help" {
		fmt.Println("Type list to list your subscriptions, /<string> to search, exit to exit, help to show this")
	} else if command == "settings" {
		fmt.Println(config)
	} else {
		fmt.Println("Type help for a list of commands")
	}
	return config, false
}
