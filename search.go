package main

import (
	"encoding/json" //for reading itunes data
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//itunes search url

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
		fmt.Printf("%d \n \t artist: %s\n\t collection: %s\n\t url: %s\n\n", i, n.ArtistName, n.CollectionName, n.FeedURL)
	}
	return results.Results, nil
}

//takes an RSS url and returns the data in the form of an array of podcast episode entries
func parseRss(input string) ([]PodcastEntry, error) {
	return make([]PodcastEntry, 0), nil
}
