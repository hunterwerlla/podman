package main

import (
	"time"
)

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
	Downloaded         map[string]PodcastEntry
	Cached             []cachedPodcast
}

type Podcast struct {
	ArtistName     string
	CollectionName string
	FeedURL        string
	Description    string
}

type PodcastEntry struct {
	PodcastTitle    string
	Title           string
	Summary         string
	Link            string
	Content         string
	GUID            string
	StorageLocation string
}

type ItunesSearch struct {
	Results []Podcast
}

type GlobalState struct {
	Playing       string
	Config        *Configuration
	playerFile    chan string
	playerControl chan int
	playerState   int
	LengthOfFile  uint64
}

type cachedPodcast struct {
	Type     Podcast
	Podcasts []PodcastEntry
	Checked  time.Time
}
type PodcastEntrySlice []PodcastEntry

//TODO make this better
//now functions on slice of podcast entry
func (entries PodcastEntrySlice) Len() int {
	return len(entries)
}

func (entries PodcastEntrySlice) Less(i, j int) bool {
	return entries[i].Title < entries[j].Title
}

func (entries PodcastEntrySlice) Swap(i, j int) {
	entries[i], entries[j] = entries[j], entries[i]
}
