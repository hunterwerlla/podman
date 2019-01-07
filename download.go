package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync/atomic"
)

var (
	downloading          int32
	downloadProgressText bytes.Buffer
)

func downloadPodcast(configuration *Configuration, podcast Podcast, ep PodcastEpisode) error {
	atomic.AddInt32(&downloading, 1)
	defer func() { atomic.AddInt32(&downloading, -1) }()
	//get rid of all stdout data
	_, w, _ := os.Pipe()
	os.Stdout = w
	folder := strings.Replace(podcast.CollectionName, " ", "", -1) //remove spaces
	fullPath := configuration.StorageLocation + "/" + folder
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
		return errors.New("invalid path")
	}
	fullPathFile = fullPath + "/" + title
	err := os.MkdirAll(fullPath, 0700)
	if err != nil {
		return err
	}
	file, err := os.Create(fullPathFile)
	defer file.Close()
	if err != nil {
		return err
	}
	link, err := http.Get(ep.Link)
	if err != nil {
		return err
	}
	defer link.Body.Close()
	//actually download
	writeTo := io.Writer(file)
	_, err = io.Copy(writeTo, link.Body)
	downloadProgressText.Truncate(0)
	if err != nil {
		return err
	}
	//add location of file to structure
	ep.StorageLocation = fullPathFile
	//file download good so add it to downloaded
	configuration.Downloaded = append(configuration.Downloaded, ep)
	return nil
}

func deleteDownloadedPodcast(configuration *Configuration, entry PodcastEpisode) {
	for position, value := range configuration.Downloaded {
		if value.GUID == entry.GUID {
			configuration.Downloaded = append(configuration.Downloaded[:position], configuration.Downloaded[position+1:]...)
			break
		}
	}
	// TODO handle errors
	go os.Remove(entry.StorageLocation)
}

func podcastIsDownloaded(configuration *Configuration, entry PodcastEpisode) bool {
	for _, value := range configuration.Downloaded {
		if value.GUID == entry.GUID {
			return true
		}
	}
	return false
}

func downloadInProgress() bool {
	if downloading > 0 {
		return true
	}
	return false
}
