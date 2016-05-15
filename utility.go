package main

//this holds things that really have no other place
import (
	"github.com/jroimartin/gocui"
)

func formatPodcastPrint(p Podcast, v *gocui.View) string {
	xMax, _ := v.Size()
	strin := p.CollectionName + " - " + p.ArtistName + " - " + p.Description
	if len(p.Description+p.CollectionName+p.ArtistName)+6 < xMax {
		//do nothing
	} else { //else truncate string
		strin = strin[0:xMax]
	}
	return strin
}

func setProperties(v *gocui.View) {
	//set properties
	v.BgColor = gocui.ColorWhite
	v.FgColor = gocui.ColorBlack
	v.Wrap = false
	v.Frame = false
}
