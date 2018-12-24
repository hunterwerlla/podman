package main

import (
	"time"
)

//player output states
const (
	ShowPlayer   = iota
	ShowDownload = iota
)

// Configuration Stores the application configuration
type Configuration struct {
	StorageLocation    string
	UpKeybind          string
	DownKeybind        string
	LeftKeybind        string
	RightKeybind       string
	PlayKeybind        string
	SearchKeybind      string
	forwardSkipLength  int
	backwardSkipLength int
	Subscribed         []Podcast
	Downloaded         map[string]PodcastEpisode
	Cached             []cachedPodcast
}

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

// ItunesSearch returns an array of podcasts that are found by itunes
type ItunesSearch struct {
	Results []Podcast
}

type cachedPodcast struct {
	Type     Podcast
	Podcasts []PodcastEpisode
	Checked  time.Time
}

// PodcastEpisodeSlice takes a slice of a list of podcasts
type PodcastEpisodeSlice []PodcastEpisode

//TODO make this better
//now functions on slice of podcast entry
func (entries PodcastEpisodeSlice) Len() int {
	return len(entries)
}

func (entries PodcastEpisodeSlice) Less(i, j int) bool {
	return entries[i].Title < entries[j].Title
}

func (entries PodcastEpisodeSlice) Swap(i, j int) {
	entries[i], entries[j] = entries[j], entries[i]
}
