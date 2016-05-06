package main

type Configuration struct {
	StorageLocation string
	UpKeybind       string
	DownKeybind     string
	LeftKeybind     string
	RightKeybind    string
	PlayKeybind     string
	SearchKeybind   string
	Subscribed      []Podcast
}

type Podcast struct {
	ArtistName     string
	CollectionName string
	FeedURL        string
	Description    string
}

type PodcastEntry struct {
	title       string
	number      int
	description string
}

type ItunesSearch struct {
	Results []Podcast
}
