package main

import (
	"encoding/json" //for reading itunes data
	"errors"
	"fmt"
	//"github.com/SlyMarbo/rss"
	"github.com/kennygrant/sanitize" //for stripping html tags
	"github.com/ungerik/go-rss"      //test differnt kind
	"io/ioutil"
	"net/http"
	"strings"
)

//search itunes for a podcast with the string given, then returns an array of Podcast
func searchItunes(term string) ([]Podcast, error) {
	const itunesUrl string = "https://itunes.apple.com/search?entity=podcast&term="
	searchUrl := itunesUrl + "\"" + term + "\""
	resp, err := http.Get(searchUrl)
	if err != nil {
		return make([]Podcast, 0), errors.New("error cannot connect to itunes server")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return make([]Podcast, 0), errors.New("Error cannot read web page")
	}
	var results ItunesSearch
	err = json.Unmarshal(body, &results)
	if err != nil {
		return make([]Podcast, 0), errors.New("Bad data returned by Itunes")
	}
	for i, n := range results.Results {
		fmt.Printf("%d \n \t artist: %s\n\t collection: %s\n\t url: %s\n", i, n.ArtistName, n.CollectionName, n.FeedURL)
	}
	return results.Results, nil
}

/*
func podcastAddDescription(podcast *Podcast) error {
	feed, err := rss.Fetch(podcast.FeedURL)
	if err != nil {
		fmt.Println("Unable to fetch RSS data, try again later")
		return err
	}
	podcast.Description = feed.Description
	return nil
}

//takes an RSS url and returns the data in the form of an array of podcast episode entries
//The Item type comes from RSS

func parseRss(input string) ([]PodcastEntry, error) {
	feed, err := rss.Fetch(input)
	if err != nil {
		fmt.Println("Unable to fetch RSS data, try again later")
		return make([]PodcastEntry, 0), nil
	}
	entries := make([]PodcastEntry, 0)
	for _, item := range feed.Items {
		//change it from Item type from RSS to built in PodcastEntry type, while also removing whitespace
		//it also strips HTML tags because a lot of podcasts include them in their RSS data
		content := sanitize.HTML(strings.Replace(item.Content, "\n", " ", -1))
		content = strings.Replace(content, "\n", "", -1)
		title := sanitize.HTML(strings.Replace(item.Title, "\n", " ", -1))
		title = strings.Replace(content, "\n", "", -1)
		fmt.Println(item.Enclosures[0])
		entries = append(entries, PodcastEntry{title, item.Summary, item.Link, content, ""})
	}
	return entries, nil
}
*/

//this function will add additional data to the podcast beyond the itunes data (a description)
func podcastAddDescription(podcast *Podcast) error {
	feed, err := rss.Read(podcast.FeedURL)
	if err != nil {
		fmt.Println("Unable to fetch RSS data, try again later")
		return err
	}
	podcast.Description = feed.Description
	return nil
}
func parseRss(input string) ([]PodcastEntry, error) {
	feed, err := rss.Read(input)
	if err != nil {
		fmt.Println("Unable to fetch RSS data, try again later")
		return make([]PodcastEntry, 0), nil
	}
	entries := make([]PodcastEntry, 0)
	for _, item := range feed.Item {
		//change it from Item type from RSS to built in PodcastEntry type, while also removing whitespace
		//it also strips HTML tags because a lot of podcasts include them in their RSS data
		//content := sanitize.HTML(strings.Replace(item.Content, "\n", " ", -1))
		content := item.Content
		//title := sanitize.HTML(strings.Replace(item.Title, "\n", " ", -1))
		//title = strings.Replace(content, "\n", "", -1)
		title := sanitize.HTML(strings.Replace(item.Title, "\n", " ", -1))
		description := strings.Replace(item.Description, "\n", "", -1)
		url := ""
		for _, enc := range item.Enclosure {
			if enc.URL != "" {
				url = enc.URL
				break
			}
		}
		entries = append(entries, PodcastEntry{title, description, url, content, ""})
	}
	return entries, nil
}
