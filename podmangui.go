//this holds the tui functions and information
package main

import (
	"bytes"
	"fmt"
	"github.com/jroimartin/gocui"
	"strings"
	"time"
)

var (
	yCursorOffset          int = 0
	selectedPodcast        Podcast
	selectedPodcastEntries []PodcastEntry
	selectedPodcastSearch  []Podcast
	stateView              int = _subscribed
	scrollingOffset        int = 0
	playerOutputState      int = _show_player
	downloadProgressText   bytes.Buffer
)

func listSubscribed(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	g.Cursor = true
	v, err := g.SetView("subscribed", -1, -1, maxX+1, maxY-1)
	//clear the view
	v.Clear()
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	//now print subscribed
	err = printSubscribed(v)
	if err != nil {
		return err
	}
	//now set current view to main view
	if _, err := g.SetCurrentView("subscribed"); err != nil {
		return err
	}
	if err := v.SetCursor(0, 1+yCursorOffset); err != nil {
		return err
	}
	return err
}

func listPodcast(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	//first 5 rows reserved for description
	v, err := g.SetView("podcast", -1, 4, maxX+1, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView { //if not created yet cool we make it
			return err
		}
	}
	//clear the view
	d, err := g.SetView("podcastDescription", -1, -1, maxX+1, 5)
	if err != nil {
		if err != gocui.ErrUnknownView { //if not created yet cool we make it
			return err
		}
	}
	//clear the view
	d.Clear()
	v.Clear()
	//set current view to podcast
	if _, err := g.SetCurrentView("podcast"); err != nil {
		return err
	}
	if err := v.SetCursor(0, 0+yCursorOffset); err != nil {
		return err
	}
	//first print the podcast description
	err = printPodcastDescription(d)
	//print the list
	err = printListPodcast(v)
	return err
}

func listSearch(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView("searchResults", -1, 1, maxX+1, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView { //if not created yet cool we make it
			return err
		}
	}
	err = printSearch(v)
	if err != nil {
		return err
	}
	d, err := g.SetView("search", -1, -1, maxX+1, 1)
	if err != nil {
		if err != gocui.ErrUnknownView { //if not created yet cool we make it
			return err
		}
	}
	err = printSearchBar(d)
	if err != nil {
		return err
	}
	//set view to search if selectedPodcasts are not null aka we have searched and have results
	if selectedPodcastSearch == nil {
		_, err = g.SetCurrentView("search")
		if err != nil {
			return err
		}
	}
	if err := v.SetCursor(0, 1+yCursorOffset); err != nil {
		return err
	}
	return nil
}

func listDownloaded(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView("downloads", -1, -1, maxX+1, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView { //if not created yet cool we make it
			return err
		}
	}
	err = printDownloaded(v)
	if err != nil {
		return err
	}
	_, err = g.SetCurrentView("downloads")
	if err != nil {
		return err
	}
	if err = v.SetCursor(0, 0+yCursorOffset); err != nil {
		return err
	}
	return nil
}

func printListPodcast(v *gocui.View) error {
	v.Clear()
	setProperties(v)
	v.Highlight = true
	var err error = nil
	//if nil then cache them
	if selectedPodcastEntries == nil {
		selectedPodcastEntries, err = getPodcastEntries(selectedPodcast, selectedPodcast.FeedURL)
	}
	if err != nil {
		fmt.Fprintln(v, "Cannot download podcast list, check your connection")
		return nil
	}
	//now actually print
	for i, thing := range selectedPodcastEntries[scrollingOffset:] {
		fmt.Fprintf(v, "%d %s -  Dl:%v - %s\n", i+1+scrollingOffset, thing.Title, isDownloaded(thing), thing.Summary)
	}
	return nil
}

//this function will print the podcast information when it goes to a podcast
func printPodcastDescription(v *gocui.View) error {
	setProperties(v)
	v.Wrap = true //turn wrap on
	//now actually print
	fmt.Fprintf(v, "Name: %s \nBy: %s\n", selectedPodcast.CollectionName, selectedPodcast.ArtistName)
	descString := selectedPodcast.Description
	fmt.Fprintf(v, "%s", descString)
	return nil
}

func printSubscribed(v *gocui.View) error {
	//first clear
	v.Clear()
	//then set properties
	setProperties(v)
	v.Highlight = true
	//if none print message and return
	if len(globals.Config.Subscribed) == 0 {
		fmt.Fprintln(v, "Scroll left to search for podcasts to subscribe to.")
		return nil
	}
	fmt.Fprintf(v, "Podcast Name - Artist - Description\n")
	for _, item := range globals.Config.Subscribed[scrollingOffset:] {
		_, xWidth := v.Size()
		fmt.Fprintf(v, "%s\n", formatPodcast(item, xWidth))
	}
	return nil
}

func printSearch(v *gocui.View) error {
	setProperties(v)
	v.Clear()
	if selectedPodcastSearch != nil && len(selectedPodcastSearch) > 0 {
		fmt.Fprintf(v, "Search Results: \n")
	} else if selectedPodcastSearch == nil {
		fmt.Fprintf(v, "Type to search \n")
	} else {
		fmt.Fprintf(v, "No results \n")
	}
	for _, thing := range selectedPodcastSearch[scrollingOffset:] {
		_, xWidth := v.Size()
		fmt.Fprintf(v, "%s\n", formatPodcast(thing, xWidth))
	}
	return nil
}
func printSearchBar(v *gocui.View) error {
	setProperties(v)
	v.Autoscroll = true //to hide subsequent entries
	v.Editable = true
	return nil
}

func printDownloaded(v *gocui.View) error {
	//first clear
	v.Clear()
	//then set properties
	setProperties(v)
	v.Highlight = true
	if len(selectedPodcastEntries) == 0 {
		fmt.Fprintln(v, "Subscribe to some podcasts and download episodes")
		return nil
	}
	for i, thing := range selectedPodcastEntries[scrollingOffset:] {
		fmt.Fprintf(v, "%d %s - %s - %s\n", i+1+scrollingOffset, thing.PodcastTitle, thing.Title, thing.Summary)
	}
	return nil
}

//TODO fix time the frame after resume (jumps quite a lot for no reason)
//the audio player
func printPlayer(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView("player", -1, maxY-2, maxX+1, maxY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		err = nil
	}
	maxX, maxY = v.Size()
	setProperties(v)
	v.Clear()
	playingPlayerPosition := 0
	playingMessage := ""
	if globals.playerState == _play {
		playingPlayerPosition = playerPosition + int(time.Since(startTime).Seconds())
	}
	count := globals.LengthOfFile
	percent := float64(playingPlayerPosition) / float64(count)
	//10 is width of numbers, 2 is width of ends
	numFilled := int(percent * float64(maxX-10.0-2.0))
	if numFilled == 0 {
		numFilled = 1 //have to have at least 1
	}
	//actually print player
	if (downloadProgressText.Len() == 0) || (playerOutputState == _show_player && downloadProgressText.Len() != 0) {
		//if playing and valid length of file
		if globals.playerState == _play && globals.LengthOfFile != 0 {
			playingMessage = fmt.Sprintf("%d/%d", playingPlayerPosition, count)
			numEmpty := int((1.0 - float64(percent)) * float64(maxX-10.0-2.0))
			playingMessage = fmt.Sprintf("%s%s%s%s%s%s\n", playingMessage, "[", strings.Repeat("=", numFilled-1), ">", strings.Repeat("-", numEmpty), "]")
		} else if globals.playerState == _pause {
			playingMessage = "paused"
		} else if globals.playerState == _stop {
			playingMessage = "stopped"
		} else {
			playingMessage = "Nothing playing"
		}
		fmt.Fprintf(v, playingMessage)
		if globals.downloadProgress != nil {
			playerOutputState = _show_download //alternate state
		}
	} else { //else print progress bar
		fmt.Fprintf(v, "%s", downloadProgressText.String())
		//only alternate if playing
		if globals.playerState == _play {
			playerOutputState = _show_player
		}
	}
	return nil
}
