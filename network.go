package main

import (
	"encoding/json" //for reading itunes data
	"errors"
	"github.com/kennygrant/sanitize"
	"github.com/ungerik/go-rss"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func sanitizePodcast(entry string) string {
	entry = sanitize.HTML(entry)
	entry = strings.Replace(entry, "\n", " ", -1)
	entry = strings.Replace(entry, "\r", " ", -1)
	entry = strings.TrimSpace(entry)
	return entry
}

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
	for ii := range results.Results {
		results.Results[ii].Description = sanitizePodcast(results.Results[ii].Description)
		results.Results[ii].ArtistName = sanitizePodcast(results.Results[ii].ArtistName)
		results.Results[ii].CollectionName = sanitizePodcast(results.Results[ii].CollectionName)
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
func getPodcastEntries(podcast Podcast, podcastCache *[]CachedPodcast) ([]PodcastEpisode, error) {
	var cacheEntry *CachedPodcast
	input := podcast.FeedURL
	for _, value := range *podcastCache {
		if podcast.CollectionName == value.Type.CollectionName && podcast.ArtistName == value.Type.ArtistName {
			cacheEntry = &value
			break
		}
	}
	// TODO create a lower timeout client to load from filesystem faster?
	feed, err := rss.Read(input)
	if err != nil {
		if cacheEntry != nil {
			return cacheEntry.Podcasts, nil
		}
		return []PodcastEpisode{{Title: "Unable to fetch RSS data, try again later"}}, nil
	}
	entries := make([]PodcastEpisode, 0)
	for _, item := range feed.Item {
		// change it from Item type from RSS to built in PodcastEpisode type,
		// while also removing whitespace and stripping HTML tags
		url := ""
		content := sanitizePodcast(item.Content)
		title := sanitizePodcast(item.Title)
		description := sanitizePodcast(item.Description)
		for _, enc := range item.Enclosure {
			if len(enc.URL) > 0 {
				url = enc.URL
				break
			}
		}
		guid := ""
		// If the GUID is empty, make one up
		if strings.TrimSpace(item.GUID) == "" {
			guid = content + title
		} else {
			guid = content + title + item.GUID
		}
		entries = append(entries, PodcastEpisode{
			PodcastTitle:    feed.Title,
			Title:           title,
			Summary:         description,
			Link:            url,
			Content:         content,
			GUID:            guid,
			StorageLocation: "",
		})
	}
	//if it's not nil we are updating
	if cacheEntry != nil {
		*cacheEntry = CachedPodcast{podcast, entries, time.Now()}
	} else { //otherwise we are creating
		*podcastCache = append(*podcastCache, CachedPodcast{podcast, entries, time.Now()})
	}
	return entries, nil
}
