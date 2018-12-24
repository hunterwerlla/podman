package configuration

import "time"

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
	Cached             []CachedPodcast
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
