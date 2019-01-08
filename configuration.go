package main

// Configuration Stores the application configuration
type Configuration struct {
	StorageLocation   string
	UpKeybind         string
	DownKeybind       string
	LeftKeybind       string
	RightKeybind      string
	PlayKeybind       string
	SearchKeybind     string
	ActionKeybind     string
	DeleteKeybind     string
	FastForward       string
	Rewind            string
	FastForwardLength int
	RewindLength      int
	Subscribed        []Podcast
	Downloaded        []PodcastEpisode
	Cached            []CachedPodcast
}

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

//CreateDefault creates default configuration
func CreateDefault() Configuration {
	return Configuration{
		StorageLocation:   "",
		UpKeybind:         "k",
		DownKeybind:       "j",
		LeftKeybind:       "h",
		RightKeybind:      "l",
		PlayKeybind:       "<Space>",
		SearchKeybind:     "/",
		ActionKeybind:     "s",
		DeleteKeybind:     "d",
		FastForward:       "<Previous>",
		Rewind:            "<Next>",
		FastForwardLength: 30,
		RewindLength:      10,
		Subscribed:        make([]Podcast, 0),
		Downloaded:        make([]PodcastEpisode, 0),
		Cached:            make([]CachedPodcast, 0),
	}
}
