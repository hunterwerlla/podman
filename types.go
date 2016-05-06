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
	Downloaded      []PodcastEntry
}

type Podcast struct {
	ArtistName     string
	CollectionName string
	FeedURL        string
	Description    string
}

type PodcastEntry struct {
	title   string
	Summary string
	Link    string
	Content string
}

type ItunesSearch struct {
	Results []Podcast
}
