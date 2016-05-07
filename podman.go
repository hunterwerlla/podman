package main

import (
	"flag"
	"fmt"
	"github.com/krig/go-sox" //for playing podcasts
	"os/user"
	"strconv"
	"strings"
	"time"
)

func main() {
	//get users home dir, the default storage
	usr, err := user.Current()
	defaultStorage := "."
	//if no error, sore in home directory
	if err == nil {
		defaultStorage = usr.HomeDir + "/" + "podman"
	}
	//make configurationg struct that holds default settings
	config := Configuration{defaultStorage, "k", "j", "h", "l", " ", "/", make([]Podcast, 0), make([]PodcastEntry, 0)}
	//read command line flags first
	noTui := flag.Bool("no-tui", false, "Select whether to use the TUI or not")
	flag.Parse()
	//read config file
	config = readConfig(config)
	//write config on sucessful exit
	//defer writeConfig(config)
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
	writeConfig(config)
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
					for i, result := range results {
						if i == num {
							fmt.Println("appending the result to subscribed")
							//add description to it
							podcastAddDescription(&result)
							//then add
							config.Subscribed = append(config.Subscribed, result)
							writeConfig(config) //update config on disk
							goto searchEnd      //considered harmful
						}
					}
					fmt.Println("Number is in wrong format or too large, try again")
				} else {
					fmt.Println("Please input either exit or sub <number>")
				}
			}
		searchEnd:
		}
	} else if command == "ls" {
		for i, entry := range config.Subscribed {
			//do nothing
			fmt.Printf("%2d\t%2s\t%15s\n", i, entry.CollectionName, entry.ArtistName)
		}
	} else if command == "rm" {
		fmt.Scanf("%s", &command)
		num, err := strconv.Atoi(command)
		if err != nil {
			fmt.Println("please use in the form of \"rm <number>\"")
			return config, false
		}
		for i, _ := range config.Subscribed {
			if i == num {
				//then remove this one
				fmt.Printf("Removing %s\n", config.Subscribed[i].CollectionName)
				config.Subscribed = append(config.Subscribed[:i], config.Subscribed[i+1:]...)
				writeConfig(config)
			}
		}
	} else if command == "show" {
		fmt.Scanf("%s", &command)
		num, err := strconv.Atoi(command)
		if err != nil {
			fmt.Println("please use in the form of \"show <number>\"")
			return config, false
		}
		for i, pc := range config.Subscribed {
			if i == num {
				entries, err := parseRss(pc.FeedURL)
				if err != nil {
					fmt.Printf("%d when attempting to parse RSS\n", err.Error())
					break
				}
				for i, _ := range entries {
					fmt.Printf("%d Title: %s\n Summary: %s\n Content: %s\n", i, entries[i].title, entries[i].Summary, entries[i].Content)
				}
			}
		}
		return config, false
	} else if command == "download" {
		fmt.Scanf("%s", &command)
		pcNum, err := strconv.Atoi(command)
		if err != nil {
			fmt.Println("please use in the form of \"download <podcast number> <episode number>\"")
			return config, false
		}
		fmt.Scanf("%s", &command)
		epNum, err := strconv.Atoi(command)
		if err != nil {
			fmt.Println("please use in the form of \"download <podcast number> <episode number>\"")
			return config, false
		}
		for ii, pc := range config.Subscribed {
			if ii == pcNum {
				entries, err := parseRss(pc.FeedURL)
				if err != nil {
					fmt.Printf("%d when attempting to parse RSS\n", err.Error())
					break
				}
				for i, entry := range entries {
					if i == epNum {
						config, err := download(config, pc, entry)
						if err != nil {
							fmt.Printf("Error when downloading: %s\n", err.Error())
						} else {
							fmt.Println("Finished downloading")
						}
						return config, false
					}
				}
				fmt.Printf("Invalid episode number %d\n", epNum)
				return config, false
			}
		}
		fmt.Println("Invalid subscription number")
		return config, false
	} else if command == "play" {
		//TODO make it send a message to a goroutine instead
		fmt.Scanf("%s", &command)
		pcNum, err := strconv.Atoi(command)
		if err != nil {
			fmt.Println("please use in the form of \"play <downloaded episode number>\"")
			return config, false
		}
		for i, item := range config.Downloaded {
			if i == pcNum {
				if !sox.Init() {
					fmt.Println("Unable to start the player")
					return config, false
				}
				defer sox.Quit()
				in := sox.OpenRead(item.StorageLocation)
				defer in.Release()
				out := sox.OpenWrite("default", in.Signal(), nil, "alsa")
				if out == nil {
					out = sox.OpenWrite("default", in.Signal(), nil, "pulseaudio")
					if out == nil {
						fmt.Println("Cannot open audio output devices")
						return config, false
					}
				}
				time.Sleep(time.Second * 10)
				return config, false
			}
		}
		fmt.Println("episode not found")
		return config, false
	} else if command == "help" {
		fmt.Println("Type ls to list your subscriptions, /<string> to search, exit to exit, help to show this")
	} else if command == "settings" {
		fmt.Println(config)
	} else {
		fmt.Println("Type help for a list of commands")
	}
	return config, false
}
