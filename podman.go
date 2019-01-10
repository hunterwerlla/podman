package main

import (
	"flag"
	"os/user"
	"time"
)

// Podcast Holds one podcast
type Podcast struct {
	ArtistName     string
	CollectionName string
	FeedURL        string
	Description    string
}

// PodcastEpisode Holds one episode of a podcast
type PodcastEpisode struct {
	// The title of the podcast according to ITunes
	PodcastTitle    string
	Title           string
	Summary         string
	Link            string
	Content         string
	GUID            string
	StorageLocation string
}

// CachedPodcast a struct that holds podcasts cached in memory/on disk
type CachedPodcast struct {
	Type     Podcast
	Podcasts []PodcastEpisode
	Checked  time.Time
}

// ItunesSearch returns an array of podcasts that are found by itunes
type ItunesSearch struct {
	Results []Podcast
}

// PodcastEpisodeSlice takes a slice of a list of podcasts
type PodcastEpisodeSlice []PodcastEpisode

func main() {
	//get users home dir, the default storage
	// TODO fix this to work on Windows + to be not a mess
	defaultStorage := "."
	usr, err := user.Current()
	//if no error, store in home directory
	if err == nil {
		defaultStorage = usr.HomeDir + "/.config/podman"
	}
	//read config file
	configuration := CreateDefault()
	configuration.StorageLocation = defaultStorage
	readConfig(&configuration)
	//read command line flags
	noTui := flag.Bool("no-tui", false, "Select whether to use the GUI or not")
	flag.Parse()
	// Start up the global player
	StartPlayer(&configuration)
	//made a decision to use TUI or not
	if *noTui == true {
		RunCui(&configuration)
	} else {
		StartTui(&configuration)
	}
}
