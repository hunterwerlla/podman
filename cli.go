package main

import (
	"fmt"
	"strconv"
	"strings"
)

func CliInterface(playerFile chan string, playerControl chan int) bool {
	command := ""
	fmt.Scanf("%s", &command)
	command = strings.ToLower(command)
	if command == "" {
		return false
	}
	if command == "exit" {
		return true
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
							globals.Config.Subscribed = append(globals.Config.Subscribed, result)
							writeConfig(*globals.Config) //update config on disk
							fmt.Println("subscribed and written to disk")
							goto searchEnd //considered harmful
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
		for i, entry := range globals.Config.Subscribed {
			//do nothing
			fmt.Printf("%2d\t%2s\t%15s\n", i, entry.CollectionName, entry.ArtistName)
		}
	} else if command == "rm" {
		fmt.Scanf("%s", &command)
		num, err := strconv.Atoi(command)
		if err != nil {
			fmt.Println("please use in the form of \"rm <number>\"")
			return false
		}
		for i, _ := range globals.Config.Subscribed {
			if i == num {
				//then remove this one
				fmt.Printf("Removing %s\n", globals.Config.Subscribed[i].CollectionName)
				globals.Config.Subscribed = append(globals.Config.Subscribed[:i], globals.Config.Subscribed[i+1:]...)
				writeConfig(*globals.Config)
			}
		}
	} else if command == "show" {
		fmt.Scanf("%s", &command)
		num, err := strconv.Atoi(command)
		if err != nil {
			fmt.Println("please use in the form of \"show <number>\"")
			return false
		}
		for i, pc := range globals.Config.Subscribed {
			if i == num {
				entries, err := getPodcastEntries(pc, pc.FeedURL)
				if err != nil {
					fmt.Printf("%d when attempting to parse RSS\n", err.Error())
					break
				}
				for i, entry := range entries {
					fmt.Printf("%d Title: %s\n Summary: %s\n Content: %s\n", i, entry.Title, entry.Summary, entry.Content)
				}
			}
		}
		return false
	} else if command == "download" {
		fmt.Scanf("%s", &command)
		pcNum, err := strconv.Atoi(command)
		if err != nil {
			fmt.Println("please use in the form of \"download <podcast number> <episode number>\"")
			return false
		}
		fmt.Scanf("%s", &command)
		epNum, err := strconv.Atoi(command)
		if err != nil {
			fmt.Println("please use in the form of \"download <podcast number> <episode number>\"")
			return false
		}
		for ii, pc := range globals.Config.Subscribed {
			if ii == pcNum {
				entries, err := getPodcastEntries(pc, pc.FeedURL)
				if err != nil {
					fmt.Printf("%d when attempting to parse RSS\n", err.Error())
					break
				}
				for i, entry := range entries {
					if i == epNum {
						_, err := download(*globals.Config, pc, entry)
						if err != nil {
							fmt.Printf("Error when downloading: %s\n", err.Error())
						} else {
							fmt.Println("Finished downloading")
						}
						return false
					}
				}
				fmt.Printf("Invalid episode number %d\n", epNum)
				return false
			}
		}
		fmt.Println("Invalid subscription number")
		return false
	} else if command == "play" {
		//TODO make it send a message to a goroutine that runs all the time instead
		fmt.Scanf("%s", &command)
		pcNum, err := strconv.Atoi(command)
		if err != nil {
			fmt.Println("please use in the form of \"play <downloaded episode number>\"")
			return false
		}
		for i, item := range globals.Config.Downloaded {
			if i == pcNum {
				//send storage location to player
				playerFile <- item.StorageLocation
				return false
			}
		}
		fmt.Println("episode not found")
		return false
	} else if command == "stop" {
		playerControl <- 2
	} else if command == "pause" {
		playerControl <- 1
	} else if command == "resume" {
		playerControl <- 0
	} else if command == "ff" {
		playerControl <- 3
	} else if command == "rewind" {
		playerControl <- 4
	} else if command == "ls-download" {
		for i, podcast := range globals.Config.Downloaded {
			fmt.Printf("%d %s %s\n", i, podcast.PodcastTitle, podcast.Title)
		}
	} else if command == "help" {
		fmt.Println("Type ls to list your subscriptions, ls-download to list downloads, start <num> to play, stop to stop, resume to resume, /<string> to search, exit to exit, help to show this")
	} else if command == "settings" {
		fmt.Println(*globals.Config)
	} else {
		fmt.Println("Type help for a list of commands")
	}
	return false
}
