package main

import (
	"flag"
	"fmt"
	"github.com/jroimartin/gocui"
	"os/user"
)

//global state
var (
	globals GlobalState = GlobalState{"", nil}
)

func main() {
	//get users home dir, the default storage
	usr, err := user.Current()
	defaultStorage := "."
	//if no error, sore in home directory
	if err == nil {
		defaultStorage = usr.HomeDir + "/" + "podman"
	}
	//make configurationg struct that holds default settings
	config := Configuration{defaultStorage, "k", "j", "h", "l", " ", "/", 30, 10, make([]Podcast, 0), make([]PodcastEntry, 0)}
	//read command line flags
	noTui := flag.Bool("no-gui", false, "Select whether to use the GUI or not")
	flag.Parse()
	//read config file
	config = readConfig(config)
	//make the channels used by player
	playerControl := make(chan int)
	playerFile := make(chan string)
	playerExit := make(chan bool)
	go play(config, playerFile, playerControl, playerExit)
	//set up annoying global variable
	globals.Config = &config
	//made a decision to use TUI or not
	if *noTui == true {
		end := false
		for end != true {
			config, end = CliInterface(config, playerFile, playerControl)
			globals.Config = &config
		}
	} else {
		g := gocui.NewGui()
		if err := g.Init(); err != nil {
			panic("Unable to start TUI, can atttempt to run --no-tui for minimal text based version")
		}
		defer g.Close()
		g.SetLayout(mainLayout)
		//allow mouse
		g.Mouse = true
		if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitGui); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("player", gocui.KeySpace, gocui.ModNone, playSelected); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		//TODO fix keybinds
		if err := g.SetKeybinding("list", gocui.KeyCtrlA, gocui.ModNone, cursorDown); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("list", gocui.KeyCtrlB, gocui.ModNone, cursorUp); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
	}
	playerControl <- 5 //tell it to exit
	writeConfig(config)
	//wait for player to clean up
	<-playerExit
}
