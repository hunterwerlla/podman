package main

func podcastIsSubscribed(configuration *Configuration, podcast *Podcast) bool {
	for _, value := range configuration.Subscribed {
		if podcast.ArtistName == value.ArtistName && podcast.CollectionName == value.CollectionName {
			return true
		}
	}
	return false
}
