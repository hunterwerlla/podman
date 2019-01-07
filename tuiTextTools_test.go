package main

import (
	"testing"
)

func TestPodcastFormat(t *testing.T) {
	p := Podcast{ArtistName: "?", CollectionName: ">", FeedURL: "NULL", Description: "a"}
	formattedString := formatPodcast(p, 1)
	if formattedString != ">" {
		t.Errorf("Truncating podcast name not working size 1: %s ", formattedString)
	}
	formattedString = formatPodcast(p, 10)
	if formattedString != "> - ? - a" {
		t.Errorf("Truncating podcast name not working size 10: %s", formattedString)
	}
}
