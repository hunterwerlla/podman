package configuration

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
		StorageLocation:    "",
		UpKeybind:          "k",
		DownKeybind:        "j",
		LeftKeybind:        "h",
		RightKeybind:       "l",
		PlayKeybind:        " ",
		SearchKeybind:      "/",
		forwardSkipLength:  30,
		backwardSkipLength: 10,
		Subscribed:         make([]Podcast, 0),
		Downloaded:         make(map[string]PodcastEpisode, 0),
		Cached:             make([]CachedPodcast, 0),
	}
}
