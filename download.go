package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func download(config Configuration, podcast Podcast, ep PodcastEntry) (Configuration, error) {
	folder := strings.Replace(podcast.CollectionName, " ", "", -1) //remove spaces
	fullPath := config.StorageLocation + "/" + folder
	fullPathFile := ""
	title := ""
	if len(ep.Title) > 30 {
		title = ep.Title[0:30]
	} else {
		title = ep.Title
	}
	//check if title has extension, if not strip possible period and add extension
	if path.Ext(title) != "mp3" {
		title = strings.Replace(title, ".", "", -1)
		title += ".mp3"
	}
	//if empty, title invalid
	if title == ".mp3" {
		return config, errors.New("invalid path")
	}
	fullPathFile = fullPath + "/" + title
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
	globals.Config = &config //update configuration
	return config, nil
}

//TODO updat to use hashmap
func isDownloaded(entry PodcastEntry) bool {
	for _, item := range globals.Config.Downloaded {
		if entry.Link == item.Link {
			return true
		}
	}
	return false
}
