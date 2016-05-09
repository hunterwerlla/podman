//this holds the tui functions and information
package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

func mainLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	g.Cursor = true
	v, err := g.SetView("list", -1, -1, maxX+1, maxY-2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.BgColor = gocui.ColorBlack
	v.FgColor = gocui.ColorWhite
	v.Wrap = false
	v.Frame = false
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
	if err := v.SetCursor(0, 1); err != nil {
		return err
	}
	return err
}

func printSubscribed(v *gocui.View) error {
	v.Clear()
	fmt.Fprintf(v, "Podcast Name ------- Artist ------ Description -------\n")
	for _, item := range globals.Config.Subscribed {
		desc := item.Description
		if len(item.Description) > 30 {
			desc = item.Description[0:30]
		}
		fmt.Fprintf(v, "%s - %s - %s \n", item.CollectionName, item.ArtistName, desc)
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
	v.BgColor = gocui.ColorBlack
	v.FgColor = gocui.ColorWhite
	v.Frame = false
	fmt.Fprintln(v, "Play Something: [==================]")
	return nil
}

func scrollDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	return nil
}

func scrollUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}
func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		x, y := v.Cursor()
		/*if y == len(globals.Config.Subscribed) {
			return nil
		}*/
		if err := v.SetCursor(x, y+1); err != nil {
			return err
		}
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		x, y := v.Cursor()
		/*
			if y == 0 {
				return nil
			}*/
		if err := v.SetCursor(x, y-1); err != nil {
			return err
		}
	}
	return nil
}
func playSelected(g *gocui.Gui, v *gocui.View) error {
	return nil
}
func quitGui(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
