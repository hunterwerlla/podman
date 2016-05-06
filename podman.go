package main

import (
	"encoding/json" //for decoding search data from itunes
	"flag"
	"fmt"
	"os"
	"strings"
	//"github.com/krig/go-sox" //for playing podcasts
)

type Configuration struct {
	storageLocation string
	upKeybind       string
	downKeybind     string
	leftKeybind     string
	rightKeybind    string
	playKeybind     string
	searchKeybind   string
	subscribed      []PodcastEntry
}

type PodcastEntry struct {
	name   string
	author string
	uri    string
}

func main() {
	//make configurationg struct that holds default settings
	config := Configuration{"~/podman", "<up>", "<down>", "<left>", "<right>", "<space>", "/", make([]PodcastEntry, 0)}
	//read command line flags first
	noTui := flag.Bool("no-tui", false, "Select whether to use the TUI or not")
	flag.Parse()
	//read config file
	config = readConfig(config)
	//made a decision to use TUI or not
	if *noTui == true {
		end := false
		for end != true {
			end = CliInterface()
		}
	} else {
		//TUI
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
	return make([]PodcastEntry, 1)
}

//read config in
func readConfig(c Configuration) Configuration {
	//check if there is a config file
	config, err := os.Open("./config.json")
	if err != nil {
		//config does not exist so build one out of the defult settings
		//first check if the storage location is ok
		if _, err := os.Stat(c.storageLocation); os.IsNotExist(err) {
			//path does not exist try to make
			err := os.Mkdir(c.storageLocation, 0666)
			if err != nil {
				//failed to create folder to store, store files in same directory as program
				c.storageLocation = "."
			}
		}
		writeConfig(c)
		return c
	}
	defer config.Close()
	//now read in the settings and write it to the configuration object
	return c
}

//save current config to file
func writeConfig(c Configuration) {
	config, err := os.Create("./config.json")
	if err != nil {
		//using default settings because cannot write settings
		return
	}
	defer config.Close()
	encoder := json.NewEncoder(config)
	encoder.Encode(c)
}
