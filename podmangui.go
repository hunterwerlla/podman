//this holds the tui functions and information
package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

func mainLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	return nil
}

func quitGui(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
