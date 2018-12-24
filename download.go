package main

import (
	"errors"
	"github.com/cheggaaa/pb"
	"github.com/jroimartin/gocui"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync/atomic"
)

var (
	downloading      int32
	downloadProgress *pb.ProgressBar
)

func download(config Configuration, podcast Podcast, ep PodcastEntry, g *gocui.Gui) (Configuration, error) {
	atomic.AddInt32(&downloading, 1)
	defer func() { atomic.AddInt32(&downloading, -1) }()
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
	if downloadProgress != nil {
		//make a progress bar length of the content
		downloadProgress.Add(int(link.ContentLength))
	} else {
		downloadProgress = pb.New(int(link.ContentLength))
		downloadProgress.SetUnits(pb.U_BYTES)
		downloadProgress.Format("[=-]")
		downloadProgress.Start()
		defer downloadProgress.Finish()
		downloadProgress.Output = &downloadProgressText
	}
	writeTo := io.MultiWriter(file, downloadProgress)
	_, err = io.Copy(writeTo, link.Body)
	//stop download progress bar
	downloadProgress = nil
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

func downloadInProgress() bool {
	if downloading > 0 {
		return true
	}
	return false
}
