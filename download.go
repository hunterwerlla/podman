package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var (
	downloading          = make(map[string]*Download, 0)
	downloadProgressText bytes.Buffer
)

const downloadWindowSize = 10

type Download struct {
	TotalDownloaded int64
	FileSize        int64
	Speed           int64
	downloadStart   time.Time
}

func (download *Download) Write(p []byte) (int, error) {
	numberBytes := len(p)
	if download.downloadStart.IsZero() {
		download.downloadStart = time.Now()
	}
	download.TotalDownloaded += int64(numberBytes)
	// do a very bad but working average
	// TODO improve to sliding window
	timeBetweenSamples := download.TotalDownloaded * 1000000000 / (time.Since(download.downloadStart).Nanoseconds() + 1)
	download.Speed = timeBetweenSamples
	return numberBytes, nil
}

func downloadPodcast(configuration *Configuration, podcast Podcast, ep PodcastEpisode) error {
	//get rid of all stdout data
	//_, w, _ := os.Pipe()
	//os.Stdout = w
	if downloading[ep.Link] != nil {
		// already downloading so bail
		return nil
	}
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
	httpResponse, err := http.Get(ep.Link)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != 200 {
		return errors.New("Download failed")
	}
	//actually download
	size := httpResponse.ContentLength
	downloadProgressTracker := &Download{FileSize: size}
	downloading[ep.Link] = downloadProgressTracker
	writeTo := io.Writer(file)
	_, err = io.Copy(writeTo, io.TeeReader(httpResponse.Body, downloadProgressTracker))
	downloadProgressText.Truncate(0)
	if err != nil {
		return err
	}
	//add location of file to structure
	ep.StorageLocation = fullPathFile
	//file download good so add it to downloaded
	configuration.Downloaded = append(configuration.Downloaded, ep)
	// remove from map
	delete(downloading, ep.Link)
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

func podcastExistsOnDisk(entry PodcastEpisode) bool {
	if _, err := os.Stat(entry.StorageLocation); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func podcastIsDownloaded(configuration *Configuration, entry PodcastEpisode) bool {
	for _, value := range configuration.Downloaded {
		if value.GUID == entry.GUID {
			if podcastExistsOnDisk(value) {
				return true
			}
			return false
		}
	}
	return false
}

func getPodcastLocation(configuration *Configuration, entry PodcastEpisode) string {
	for _, value := range configuration.Downloaded {
		if value.GUID == entry.GUID {
			return value.StorageLocation
		}
	}
	return ""
}

func downloadInProgress() bool {
	if len(downloading) > 0 {
		return true
	}
	return false
}
