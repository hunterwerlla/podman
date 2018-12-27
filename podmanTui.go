//this holds Tui drawing functions
package main

import (
	"bytes"
	"fmt"
	"github.com/jroimartin/gocui"
	"strings"
	"time"
)

const (
	ScreenAll           string = ""
	ScreenHeader        string = "header"
	ScreenSubscribed    string = "subscribed"
	ScreenPodcast       string = "podcast"
	ScreenSearch        string = "search"
	ScreenSearchResults string = "searchResults"
	ScreenDownloads     string = "downloads"
)

const (
	ShowPlayer   = iota
	ShowDownload = iota
)

var (
	yCursorOffset          = 0
	selectedPodcast        Podcast
	selectedPodcastEntries []PodcastEpisode
	selectedPodcastSearch  []Podcast
	stateView              = ScreenSubscribed
	scrollingOffset        = 0
	playerOutputState      = ShowPlayer
	downloadProgressText   bytes.Buffer
	config                 *Configuration
)

func SetTuiConfiguration(configuration *Configuration) {
	config = configuration
}

func setTuiScreenProperties(v *gocui.View) {
	//set properties
	v.Clear()
	v.BgColor = gocui.ColorWhite
	v.FgColor = gocui.ColorBlack
	v.Wrap = false
	v.Frame = false
}

func formatPodcast(p Podcast, max int) string {
	strin := p.CollectionName + " - " + p.ArtistName + " - " + p.Description
	if len(p.Description+p.CollectionName+p.ArtistName)+6 < max {
		//do nothing
	} else { //else truncate string
		strin = strin[0:max]
	}
	return strin
}

func listSubscribed(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	g.Cursor = true
	v, err := g.SetView(ScreenHeader, -1, -1, maxX+1, 4)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	err = printHeader(v)
	v, err = g.SetView(ScreenSubscribed, -1, 3, maxX+1, maxY-1)
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
	if _, err = g.SetCurrentView(ScreenSubscribed); err != nil {
		return err
	}
	if err = v.SetCursor(0, yCursorOffset); err != nil {
		return err
	}
	return err
}

func listPodcast(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	//first 5 rows reserved for description
	v, err := g.SetView(ScreenPodcast, -1, 4, maxX+1, maxY-1)
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
	if _, err = g.SetCurrentView(ScreenPodcast); err != nil {
		return err
	}
	if err = v.SetCursor(0, 0+yCursorOffset); err != nil {
		return err
	}
	//first print the podcast description
	_ = printPodcastDescription(d)
	//print the list
	_ = printListPodcast(v)
	return err
}

func listSearch(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView(ScreenSearchResults, -1, 1, maxX+1, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView { //if not created yet cool we make it
			return err
		}
	}
	err = printSearch(v)
	if err != nil {
		return err
	}
	d, err := g.SetView(ScreenSearch, -1, -1, maxX+1, 1)
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
		_, err = g.SetCurrentView(ScreenSearch)
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
	setTuiScreenProperties(v)
	v.Highlight = true
	var err error
	//if nil then cache them
	if selectedPodcastEntries == nil {
		selectedPodcastEntries, err = getPodcastEntries(selectedPodcast, selectedPodcast.FeedURL, &config.Cached)
	}
	if err != nil {
		fmt.Fprintf(v, "Cannot download podcast list, check your connection")
		return nil
	}
	//now actually print
	for i, thing := range selectedPodcastEntries[scrollingOffset:] {
		fmt.Fprintf(v, "%d %s -  Dl:%t - %s\n", i+1+scrollingOffset, thing.Title, PodcastIsDownloaded(thing), thing.Summary)
	}
	return nil
}

//this function will print the podcast information when it goes to a podcast
func printPodcastDescription(v *gocui.View) error {
	setTuiScreenProperties(v)
	v.Wrap = true //turn wrap on
	//now actually print
	fmt.Fprintf(v, "Name: %s \nBy: %s", selectedPodcast.CollectionName, selectedPodcast.ArtistName)
	descString := selectedPodcast.Description
	fmt.Fprintf(v, "%s", descString)
	return nil
}

func printHeader(v *gocui.View) error {
	setTuiScreenProperties(v)
	fmt.Fprintln(v, " _____         _                  ")
	fmt.Fprintln(v, "|  _  | ___  _| | _____  ___  ___ ")
	fmt.Fprintln(v, "|   __|| . || . ||     || .'||   |")
	fmt.Fprintln(v, "|__|   |___||___||_|_|_||__,||_|_|")
	return nil
}

func printSubscribed(v *gocui.View) error {
	// set properties
	setTuiScreenProperties(v)
	v.Highlight = true
	//if none print message and return
	if len(config.Subscribed) == 0 {
		fmt.Fprintln(v, "Scroll left to search for podcasts to subscribe to.")
		return nil
	}
	for _, item := range config.Subscribed[scrollingOffset:] {
		xWidth, _ := v.Size()
		fmt.Fprintf(v, "%s\n", formatPodcast(item, xWidth))
	}
	return nil
}

func printSearch(v *gocui.View) error {
	setTuiScreenProperties(v)
	v.Clear()
	if selectedPodcastSearch != nil && len(selectedPodcastSearch) > 0 {
		fmt.Fprintf(v, "Search Results: \n")
	} else if selectedPodcastSearch == nil {
		fmt.Fprintf(v, "Type to search \n")
	} else {
		fmt.Fprintf(v, "No results \n")
	}
	for _, thing := range selectedPodcastSearch[scrollingOffset:] {
		xWidth, _ := v.Size()
		fmt.Fprintf(v, "%s\n", formatPodcast(thing, xWidth))
	}
	return nil
}

func printSearchBar(v *gocui.View) error {
	setTuiScreenProperties(v)
	v.Autoscroll = true //to hide subsequent entries
	v.Editable = true
	return nil
}

func printDownloaded(v *gocui.View) error {
	//first clear
	v.Clear()
	//then set properties
	setTuiScreenProperties(v)
	v.Highlight = true
	if len(selectedPodcastEntries) == 0 {
		fmt.Fprintf(v, "Subscribe to some podcasts and download episodes")
		return nil
	}
	for i, thing := range selectedPodcastEntries[scrollingOffset:] {
		fmt.Fprintf(v, "%d %s - %s - %s", i+1+scrollingOffset, thing.PodcastTitle, thing.Title, thing.Summary)
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
	}
	maxX, _ = v.Size()
	setTuiScreenProperties(v)
	v.Clear()
	playingPlayerPosition := 0
	playingMessage := ""
	if GetPlayerState() == Play {
		// TODO fix this logic
		playingPlayerPosition = GetPlayerPosition()
	}
	playingFileLength := GetLengthOfPlayingFile()
	percent := float64(playingPlayerPosition) / float64(playingFileLength)
	//10 is width of numbers, 2 is width of ends
	numFilled := int(percent * float64(maxX-10.0-2.0))
	if numFilled == 0 {
		numFilled = 1 //have to have at least 1
	}
	//actually print player
	if (downloadProgressText.Len() == 0) || (playerOutputState == ShowPlayer && downloadProgressText.Len() != 0) {
		//if playing and valid length of file
		if GetPlayerState() == Play && playingFileLength != 0 {
			playingMessage = fmt.Sprintf("%d/%d", playingPlayerPosition, playingFileLength)
			numEmpty := int((1.0 - float64(percent)) * float64(maxX-10.0-2.0))
			playingMessage = fmt.Sprintf("%s%s%s%s%s%s\n", playingMessage, "[", strings.Repeat("=", numFilled-1), ">", strings.Repeat("-", numEmpty), "]")
		} else if GetPlayerState() == Pause {
			playingMessage = "paused"
		} else if GetPlayerState() == Stop {
			playingMessage = "stopped"
		} else {
			playingMessage = "Nothing playing"
		}
		fmt.Fprintf(v, playingMessage)
		if downloadInProgress() {
			playerOutputState = ShowDownload //alternate state
		}
	} else { //else print progress bar
		fmt.Fprintf(v, "%s", downloadProgressText.String())
		//only alternate if playing
		if GetPlayerState() == Play {
			playerOutputState = ShowPlayer
		}
	}
	return nil
}

// TODO figure out how to re-enable this in a reasonable way
func refreshTui(g *gocui.Gui) {
	update := time.NewTicker(time.Millisecond * 500).C
	stopTick := make(chan bool)
	defer close(stopTick)
	go func() {
		for {
			select {
			case <-update:
				g.Update(TuiHandler)
			case <-stopTick:
				return
			}
		}
	}()
}
