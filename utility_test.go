package main

import (
	. "github.com/hunterwerlla/podman/configuration"
	"testing"
)

func TestFormat(t *testing.T) {
	p := Podcast{ArtistName: "?", CollectionName: ">", FeedURL: "NULL", Description: "a"}
	formattedString := formatPodcast(p, 1)
	if formattedString != ">" {
		t.Errorf("Truncating podcast name not working: %s ", formattedString)
	}
	formattedString = formatPodcast(p, 10)
	if formattedString != "> - ? - a" {
		t.Errorf("Truncating podcast name not working: %s", formattedString)
	}
}
