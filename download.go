package main

import (
	"errors"
	"github.com/jroimartin/gocui"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func download(config Configuration, podcast Podcast, ep PodcastEntry, g *gocui.Gui) (Configuration, error) {
	//get rid of all stdout data
	_, w, _ := os.Pipe()
	os.Stdout = w
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
	//actually download
	if globals.downloadProgress != nil {
		//make a progress bar length of the content
		globals.downloadProgress.Add(int(link.ContentLength))
	} else {
		globals.downloadProgress = pb.New(int(link.ContentLength))
		globals.downloadProgress.SetUnits(pb.U_BYTES)
		globals.downloadProgress.Format("[=-]")
		globals.downloadProgress.Start()
		defer globals.downloadProgress.Finish()
		globals.downloadProgress.Output = &downloadProgressText
	}
	writeTo := io.MultiWriter(file, globals.downloadProgress)
	_, err = io.Copy(writeTo, link.Body)
	//stop download progress bar
	globals.downloadProgress = nil
	downloadProgressText.Truncate(0)
	if err != nil {
		return config, err
	}
	//add location of file to structure
	ep.StorageLocation = fullPathFile
	//file download good so add it to downloaded
	config.Downloaded[ep.GUID] = ep
	globals.Config = &config //update configuration
	writeConfig(*globals.Config)
	return config, nil
}

func isDownloaded(entry PodcastEntry) bool {
	_, ok := globals.Config.Downloaded[entry.GUID]
	if ok {
		return true
	}
	return false
}

func isDownloadedPath(path string) bool {
	for _, item := range globals.Config.Downloaded {
		if strings.Contains(item.Link, path) {
			return true
		}
	}
	return false
}
