package main

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func download(config Configuration, podcast Podcast, ep PodcastEntry) (Configuration, error) {
	folder := strings.Replace(podcast.CollectionName, " ", "", -1) //remove spaces
	fullPath := config.StorageLocation + "/" + folder
	fullPathFile := ""
	if len(ep.title) > 20 {
		fullPathFile = fullPath + "/" + ep.title[0:20]
	} else {
		fullPathFile = fullPath + "/" + ep.title
	}
	err := os.MkdirAll(fullPath, 0700)
	if err != nil {
		return config, err
	}
	file, err := os.Create(fullPathFile)
	defer file.Close()
	if err != nil {
		return config, err
	}
	link, err := http.Get(ep.Link)
	defer link.Body.Close()
	if err != nil {
		return config, err
	}
	_, err = io.Copy(file, link.Body)
	if err != nil {
		return config, err
	}
	//add location of file to structure
	ep.StorageLocation = fullPathFile
	//file download good so add it to downloaded
	config.Downloaded = append(config.Downloaded, ep)
	writeConfig(config)
	return config, nil
}
