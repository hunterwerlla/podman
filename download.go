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
	if _, exists := downloading[ep.Link]; exists {
		// already downloading so bail
		return nil
	}
	// TODO do a real race condition fix here, although this is fine for human speed
	downloading[ep.Link] = &Download{TotalDownloaded: 0}
	folder := strings.Replace(podcast.CollectionName, " ", "", -1) //remove spaces
	basePath := configuration.StorageLocation + "/" + folder
	filePath := ""
	title := ""
	if len(ep.Title) > 30 {
		title = ep.Title[0:30]
	} else {
		title = ep.Title
	}
	// TODO support other types of files
	//check if title has extension, if not strip possible period and add extension
	if path.Ext(title) != "mp3" {
		title = strings.Replace(title, ".", "", -1)
		title += ".mp3"
	}
	//if empty, title invalid
	if title == ".mp3" {
		return errors.New("invalid path")
	}
	filePath = basePath + "/" + title
	err := os.MkdirAll(basePath, 0700)
	if err != nil {
		return err
	}
	file, err := os.Create(filePath)
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
	ep.StorageLocation = filePath
	//file download good so add it to downloaded
	configuration.Downloaded = append(configuration.Downloaded, ep)
	// remove from downloading map
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

func podcastIsDownloaded(configuration *Configuration, entry *PodcastEpisode) bool {
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
