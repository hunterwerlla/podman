//this holds the tui functions and information
package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"strings"
	"time"
)

var (
	yCursorOffset          int = 0
	selectedPodcast        Podcast
	selectedPodcastEntries []PodcastEntry
	stateView              int = 0 //0 is listSubscribed, 1 is listPodcast, 2 is listSearch, 3 is listDownload
)

func guiHandler(g *gocui.Gui) error {
	if stateView == 0 {
		listSubscribed(g)
	} else if stateView == 1 {
		listPodcast(g)
	}
	printPlayer(g)
	return nil
}

func listSubscribed(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	g.Cursor = true
	v, err := g.SetView("subscribed", -1, -1, maxX+1, maxY-2)
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
	if err := g.SetCurrentView("subscribed"); err != nil {
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
	v, err := g.SetView("podcast", -1, 5, maxX+1, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView { //if not created yet cool we make it
			return err
		}
	}
	//clear the view
	v.Clear()
	d, err := g.SetView("podcastDescription", -1, -1, maxX+1, 5)
	if err != nil {
		if err != gocui.ErrUnknownView { //if not created yet cool we make it
			return err
		}
	}
	//clear the view
	d.Clear()
	//set current view to podcast
	if err := g.SetCurrentView("podcast"); err != nil {
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

func printSubscribed(v *gocui.View) error {
	//first clear
	v.Clear()
	//then set properties
	setProperties(v)
	xMax, _ := v.Size()
	spacing := (xMax - 34) / 3 //43 chracters
	space := strings.Repeat("-", spacing)
	fmt.Fprintf(v, "Podcast Name %s Artist %s Description %s\n", space, space, space)
	for _, item := range globals.Config.Subscribed {
		strin := item.CollectionName + " - " + item.ArtistName + " - " + item.Description
		if len(item.Description+item.CollectionName+item.ArtistName)+6 < xMax {
			//do nothing
		} else { //else truncate string
			strin = strin[0:xMax]
		}
		fmt.Fprintf(v, "%s\n", strin)
	}
	return nil
}

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
	v.Clear()
	setProperties(v)
	if globals.playerState == -1 {
		fmt.Fprintf(v, "Play Something!")
	} else {
		playingPlayerPosition := 0
		playingMessage := ""
		if globals.playerState == 0 {
			playingPlayerPosition = playerPosition + int(time.Since(startTime).Seconds())
		} else {
			startTime.Add(time.Second) //else keep updating start time
		}
		count := globals.LengthOfFile
		percent := float64(playingPlayerPosition) / float64(count)
		maxX, _ := v.Size()
		//10 is width of numbers, 2 is width of ends
		numFilled := int(percent * float64(maxX-10.0-2.0))
		if numFilled == 0 {
			numFilled++
		}
		if globals.playerState == 0 {
			playingMessage = fmt.Sprintf("%d/%d", playingPlayerPosition, count)
		} else if globals.playerState == 1 {
			playingMessage = "paused"
		} else if globals.playerState == 2 {
			playingMessage = "stopped"
		} else {
			playingMessage = "Play Something"
		}
		numEmpty := int((1.0 - float64(percent)) * float64(maxX-10.0-2.0))
		fmt.Fprintf(v, "%s%s%s%s%s%s\n", playingMessage, "[", strings.Repeat("=", numFilled-1), ">", strings.Repeat("-", numEmpty), "]")
	}
	return nil
}

//TODO add scrolling beyond screen and not crash
func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		x, y := v.Cursor()
		if stateView == 0 {
			//if y is equal to number subscribed+1 is at bottom
			if y == len(globals.Config.Subscribed) {
				return nil
			}
		} else if stateView == 1 {
			if y == len(selectedPodcastEntries)-1 {
				return nil
			}
		} else {
			return nil
		}
		yCursorOffset++
		if err := v.SetCursor(x, y+1); err != nil {
			return err
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		x, y := v.Cursor()
		//if Y is 1 at the top, so don't move up again
		//TODO fix celing
		if stateView == 0 {
			if y == 1 {
				return nil
			}
		} else if stateView == 1 {
			if y == 0 {
				return nil
			}
		} else {
			return nil
		}
		yCursorOffset--
		if err := v.SetCursor(x, y-1); err != nil {
			return err
		}
	}
	return nil
}

//the second view a list podcast

func switchListPodcast(g *gocui.Gui, v *gocui.View) error {
	_, position := v.Cursor()                               //get cursor position to select
	yCursorOffset = 0                                       //reset cursor
	selectedPodcast = globals.Config.Subscribed[position-1] //select the podcast put in memory
	selectedPodcastEntries = nil                            //now delete the cache from the last time
	//change layout
	stateView = 1
	//delete old views
	g.DeleteView("subscribed")
	return nil
}
func switchListSubscribed(g *gocui.Gui, v *gocui.View) error {
	yCursorOffset = 0 //reset cursor
	//change layout
	stateView = 0
	//delete other views
	g.DeleteView("podcast")
	g.DeleteView("podcastDescription")
	return nil
}

//TODO null check
func switchPlayDownload(g *gocui.Gui, v *gocui.View) error {
	_, position := v.Cursor() //get cursor position to select
	var toPlay PodcastEntry
	guid := selectedPodcastEntries[position].GUID
	if isDownloaded(selectedPodcastEntries[position]) == false {
		download(*globals.Config, selectedPodcast, selectedPodcastEntries[position])
		//point it at the new podcast
		toPlay = globals.Config.Downloaded[len(globals.Config.Downloaded)-1]
	}
	//TODO fix this awful code
	for _, thing := range globals.Config.Downloaded {
		if thing.GUID == guid {
			toPlay = thing
			break
		}
	}
	//now play
	if toPlay := toPlay.StorageLocation; toPlay != "" {
		globals.playerFile <- toPlay
	}
	return nil
}

//this function will print the podcast information when it goes to a podcast
func printPodcastDescription(v *gocui.View) error {
	setProperties(v)
	//now actually print
	fmt.Fprintf(v, "Name: %s By: %s\n", selectedPodcast.CollectionName, selectedPodcast.ArtistName)
	maxX, _ := v.Size()
	descString := ""
	n := 0
	for ii := len(selectedPodcast.Description); ii > maxX; ii -= maxX {
		descString += selectedPodcast.Description[n*maxX:n*maxX+maxX] + "\n"
		n++
	}
	descString += selectedPodcast.Description[maxX*n:]
	fmt.Fprintf(v, "%s", descString)
	return nil
}
func printListPodcast(v *gocui.View) error {
	v.Clear()
	setProperties(v)
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
	for i, thing := range selectedPodcastEntries {
		//TODO make this efficent by adding a map
		fmt.Fprintf(v, "%d %s - %s - Downloaded: %v\n", i+1, thing.Title, thing.Content, isDownloaded(thing))
	}
	return nil
}

func setProperties(v *gocui.View) {
	//First clear
	v.Clear()
	//set properties
	v.BgColor = gocui.ColorWhite
	v.FgColor = gocui.ColorBlack
	v.Wrap = false
	v.Frame = false
}

func quitGui(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func togglePlayerState(g *gocui.Gui, v *gocui.View) error {
	//pause so will not enter invalid state
	time.Sleep(time.Millisecond * 100)
	if globals.playerState == 1 {
		globals.playerControl <- 0
		globals.playerState = 0
	} else if globals.playerState == 0 {
		globals.playerControl <- 1
		globals.playerState = 1
	}
	return nil
}
