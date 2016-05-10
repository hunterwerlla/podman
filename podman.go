package main

import (
	"flag"
	"fmt"
	"github.com/jroimartin/gocui"
	"os/user"
)

//global state
var (
	globals GlobalState = GlobalState{"", nil, nil, nil}
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
	globals.playerControl = make(chan int)
	globals.playerFile = make(chan string)
	playerExit := make(chan bool)
	go play(playerExit)
	//set up annoying global variable
	globals.Config = &config
	//made a decision to use TUI or not
	if *noTui == true {
		end := false
		for end != true {
			config, end = CliInterface(config, globals.playerFile, globals.playerControl)
			globals.Config = &config
		}
	} else {
		g := gocui.NewGui()
		if err := g.Init(); err != nil {
			panic("Unable to start TUI, can atttempt to run --no-tui for minimal text based version")
		}
		defer g.Close()
		g.SetLayout(listSubscribed)
		//allow mouse
		g.Mouse = true
		//set keybinds
		if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitGui); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		//TODO fix keybinds
		if err := g.SetKeybinding("subscribed", gocui.KeyArrowDown, gocui.ModNone, cursorDownList); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("subscribed", gocui.KeyArrowUp, gocui.ModNone, cursorUpList); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		//enter on list goes to the list of episodes
		if err := g.SetKeybinding("subscribed", gocui.KeyEnter, gocui.ModNone, switchListPodcast); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		//Podcast view up down
		if err := g.SetKeybinding("podcast", gocui.KeyArrowDown, gocui.ModNone, cursorDownPodcast); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("podcast", gocui.KeyArrowUp, gocui.ModNone, cursorUpPodcast); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("podcast", gocui.KeyArrowLeft, gocui.ModNone, switchListSubscribed); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("podcast", gocui.KeyEnter, gocui.ModNone, switchPlayDownload); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		//main loop
		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}

	}
	globals.playerControl <- 5 //tell it to exit
	writeConfig(config)
	//wait for player to clean up
	<-playerExit
}
