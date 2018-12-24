package main

import (
	"encoding/json" //for reading itunes data
	"errors"
	"fmt"
	. "github.com/hunterwerlla/podman/configuration"
	"github.com/kennygrant/sanitize" //for stripping html tags
	"github.com/ungerik/go-rss"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//search itunes for a podcast with the string given, then returns an array of Podcast
func searchItunes(term string) ([]Podcast, error) {
	const itunesURL string = "https://itunes.apple.com/search?entity=podcast&term="
	searchURL := itunesURL + "\"" + term + "\""
	resp, err := http.Get(searchURL)
	if err != nil {
		return make([]Podcast, 0), errors.New("error cannot connect to itunes server")
	}
	defer resp.Body.Close()
	//empty body is also error
	if resp.ContentLength == 0 {
		return make([]Podcast, 0), errors.New("error cannot connect to itunes server")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return make([]Podcast, 0), errors.New("Error cannot read web page")
	}
	var results ItunesSearch
	err = json.Unmarshal(body, &results)
	if err != nil {
		return make([]Podcast, 0), errors.New("Bad data returned by Itunes")
	}
	return results.Results, nil
}

//this function will add additional data to the podcast beyond the itunes data (a description)
func podcastAddDescription(podcast *Podcast) error {
	feed, err := rss.Read(podcast.FeedURL)
	if err != nil {
		//fmt.Println("Unable to fetch RSS data, try again later")
		return err
	}
	podcast.Description = feed.Description
	return nil
}

//TODO strip HTML
func getPodcastEntries(podcast Podcast, input string, podcastCache *[]CachedPodcast) ([]PodcastEpisode, error) {
	var cacheEntry *CachedPodcast
	for _, thing := range *podcastCache {
		if podcast.CollectionName == thing.Type.CollectionName && podcast.ArtistName == thing.Type.ArtistName {
			cacheEntry = &thing
			break
		}
	}
	//first check if we need to update
	if cacheEntry != nil {
		//TODO set time to update
		if time.Since(cacheEntry.Checked).Hours() < 12 {
			return cacheEntry.Podcasts, nil
		}
	}
	feed, err := rss.Read(input)
	if err != nil {
		if cacheEntry != nil {
			//TODO return an error that isn't null
			return cacheEntry.Podcasts, nil
		}
		fmt.Println("Unable to fetch RSS data, try again later")
		return make([]PodcastEpisode, 0), nil
	}
	entries := make([]PodcastEpisode, 0)
	for _, item := range feed.Item {
		//change it from Item type from RSS to built in PodcastEpisode type, while also removing whitespace
		//it also strips HTML tags because a lot of podcasts include them in their RSS data
		content := sanitize.HTML(strings.Replace(item.Content, "\n", " ", -1))
		title := sanitize.HTML(strings.Replace(item.Title, "\n", " ", -1))
		description := sanitize.HTML(strings.Replace(item.Description, "\n", "", -1))
		url := ""
		for _, enc := range item.Enclosure {
			if enc.URL != "" {
				url = enc.URL
				break
			}
		}
		guid := ""
		if strings.TrimSpace(item.GUID) == "" {
			guid = content + title
		} else {
			guid = content + title + item.GUID
		}
		entries = append(entries, PodcastEpisode{feed.Title, title, description, url, content, guid, ""})
	}
	//if it's not nil we are updating
	if cacheEntry != nil {
		*cacheEntry = CachedPodcast{podcast, entries, time.Now()}
	} else { //otherwise we are creating
		*podcastCache = append(*podcastCache, CachedPodcast{podcast, entries, time.Now()})
	}
	return entries, nil
}
