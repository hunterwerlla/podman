package main

import (
	"fmt"
	"strconv"
	"strings"
)

func RunCui(config *Configuration) {
	end := false
	for end != true {
		end = CliCommand(config)
	}
}

// CliCommand runs the main loop for a CLI based session
func CliCommand(config *Configuration) bool {
	command := ""
	_, err := fmt.Scanf("%s", &command)
	if err != nil {
		return false
	}
	command = strings.ToLower(command)
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
							config.Subscribed = append(config.Subscribed, result)
							WriteConfig(config) //update config on disk
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
		for i, entry := range config.Subscribed {
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
		for i := range config.Subscribed {
			if i == num {
				//then remove this one
				fmt.Printf("Removing %s\n", config.Subscribed[i].CollectionName)
				config.Subscribed = append(config.Subscribed[:i], config.Subscribed[i+1:]...)
				WriteConfig(config)
			}
		}
	} else if command == "show" {
		fmt.Scanf("%s", &command)
		num, err := strconv.Atoi(command)
		if err != nil {
			fmt.Println("please use in the form of \"show <number>\"")
			return false
		}
		found := false
		for i, pc := range config.Subscribed {
			if i == num {
				found = true
				entries, err := getPodcastEntries(pc, pc.FeedURL, &config.Cached)
				if err != nil {
					fmt.Printf("%s when attempting to parse RSS\n", err.Error())
					break
				}
				for i, entry := range entries {
					fmt.Printf("%d Title: %s\n    Summary: %s \n    Content: %s\n\n", i, entry.Title, entry.Summary, entry.Content)
				}
			}
		}
		if found == false {
			fmt.Printf("Num %d not found \n", num)
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
		for ii, pc := range config.Subscribed {
			if ii == pcNum {
				entries, err := getPodcastEntries(pc, pc.FeedURL, &config.Cached)
				if err != nil {
					fmt.Printf("%s when attempting to parse RSS\n", err.Error())
					break
				}
				for i, entry := range entries {
					if i == epNum {
						config, err := DownloadPodcast(config, pc, entry)
						WriteConfig(config)
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
		i := 0
		for _, value := range config.Downloaded {
			if i == pcNum {
				//send storage location to player
				SetPlaying(value.StorageLocation)
				SetPlayerState(Play)
				break
			}
			i++
		}
		fmt.Println("episode not found")
		return false
	} else if command == "stop" {
		SetPlayerState(Stop)
	} else if command == "pause" {
		SetPlayerState(Pause)
	} else if command == "resume" {
		SetPlayerState(Resume)
	} else if command == "ff" {
		SetPlayerState(FastForward)
	} else if command == "rewind" {
		SetPlayerState(Rewind)
	} else if command == "ls-downloaded" {
		for i, podcast := range config.Downloaded {
			fmt.Printf("%s %s %s\n", i, podcast.PodcastTitle, podcast.Title)
		}
	} else if command == "help" {
		fmt.Println("Type ls to list your subscriptions, ls-downloaded to list downloads, start <num> to play, stop to stop, resume to resume, show <num> to show subscribed, /<string> to search, exit to exit, help to show this")
	} else if command == "settings" {
		fmt.Println(*config)
	} else {
		fmt.Println("Type help for a list of commands")
	}
	return false
}
