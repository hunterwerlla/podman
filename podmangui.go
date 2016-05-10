//this holds the tui functions and information
package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"strings"
)

var (
	yCursorOffset   int      = 0
	selectedPodcast *Podcast = nil
	cachedPodcast   []PodcastEntry
)

func listLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	g.Cursor = true
	v, err := g.SetView("list", -1, -1, maxX+1, maxY-2)
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
	//now print the player
	err = printPlayer(g)
	//now set current view to main view
	if err := g.SetCurrentView("list"); err != nil {
		return err
	}
	if err := v.SetCursor(0, 1+yCursorOffset); err != nil {
		return err
	}
	return err
}

func printSubscribed(v *gocui.View) error {
	//first clear
	v.Clear()
	//then set properties
	//TODO add in colors from settings
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
	setProperties(v)
	fmt.Fprintf(v, "Play Something: [%s]", strings.Repeat("=", maxX-18))
	return nil
}

//cursor movement functions, should consolodate
//TODO reduce number
func cursorUpList(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		x, y := v.Cursor()
		//if Y is 1 at the top, so don't move up again
		if y == 1 {
			return nil
		}
		yCursorOffset--
		if err := v.SetCursor(x, y-1); err != nil {
			return err
		}
	}
	return nil
}

func cursorDownList(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		x, y := v.Cursor()
		//if y is equal to number subscribed+1 is at bottom
		if y == len(globals.Config.Subscribed) {
			return nil
		}
		yCursorOffset++
		if err := v.SetCursor(x, y+1); err != nil {
			return err
		}
	}
	return nil
}

//the second view a list podcast

func switchListPodcast(g *gocui.Gui, v *gocui.View) error {
	_, position := v.Cursor()                                  //get cursor position to select
	yCursorOffset = 0                                          //reset cursor
	selectedPodcast = &(globals.Config.Subscribed[position-1]) //minus 1 is needed
	cachedPodcast = nil                                        //now delete the cache from the last time
	g.SetLayout(listPodcast)
	return nil
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
	d, err := g.SetView("podcastDescription", -1, -1, maxX+1, 5)
	if err != nil {
		if err != gocui.ErrUnknownView { //if not created yet cool we make it
			return err
		}
	}
	//set current view to podcast
	if err := g.SetCurrentView("podcast"); err != nil {
		return err
	}
	if err := v.SetCursor(0, 1+yCursorOffset); err != nil {
		return err
	}
	//first print the podcast description
	err = printPodcastDescription(d)
	//first print the list
	err = printList(v)
	//now print the player
	err = printPlayer(g)
	return err
}

//this function will print the podcast information when it goes to a podcast
func printPodcastDescription(v *gocui.View) error {
	setProperties(v)
	//now actually print
	return nil
}
func printList(v *gocui.View) error {
	setProperties(v)
	//now actually print
	fmt.Fprintf(v, "BLEH! \n %v", *selectedPodcast)
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
func switchToDownload(g *gocui.Gui, v *gocui.View) error {
	return nil
}

func playSelected(g *gocui.Gui, v *gocui.View) error {
	return nil
}
func quitGui(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func cursorUpPodcast(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		x, y := v.Cursor()
		//if Y is 1 at the top, so don't move up again
		//TODO fix celing
		if y == 1 { //y=1 is the top because of the row that describes the colums
			return nil
		}
		yCursorOffset--
		if err := v.SetCursor(x, y-1); err != nil {
			return err
		}
	}
	return nil
}
func cursorDownPodcast(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		x, y := v.Cursor()
		if y == len(cachedPodcast) {
			return nil
		}
		yCursorOffset++
		if err := v.SetCursor(x, y+1); err != nil {
			return err
		}
	}
	return nil
}
