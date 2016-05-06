package main

import (
	"errors"
	"os"
	"strings"
)

func download(config Configuration, podcast Podcast, ep PodcastEntry) (Configuration, error) {
	folder := strings.Trim(podcast.CollectionName, " \t\n")
	fullPath := config.StorageLocation + "/" + folder
	err := os.MkdirAll(fullPath, 0700)
	if err != nil {
		return config, err
	}
	if true {
		return config, errors.New("")
	}
	return config, nil
}
