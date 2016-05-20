package main

import (
	"flag"
	"fmt"
	"github.com/jroimartin/gocui"
	"os/user"
	"time"
)

//global state
var (
	globals GlobalState = GlobalState{"", nil, nil, nil, -1, 0}
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
	config := Configuration{defaultStorage, "k", "j", "h", "l", " ", "/", 30, 10, make([]Podcast, 0), make(map[string]PodcastEntry, 0), make([]cachedPodcast, 0)}
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
			end = CliInterface(globals.playerFile, globals.playerControl)
		}
	} else {
		g := gocui.NewGui()
		if err := g.Init(); err != nil {
			panic("Unable to start TUI, can atttempt to run --no-tui for minimal text based version")
		}
		defer g.Close()
		//set main window
		g.SetLayout(guiHandler)
		//now a goroutine that updates every second
		update := time.NewTicker(time.Millisecond * 1000).C
		stopTick := make(chan bool)
		go func() {
			for {
				select {
				case <-update:
					g.Execute(guiHandler)
				case <-stopTick:
					return
				}
			}
		}()
		//allow mouse
		g.Mouse = true
		//global keybinds
		if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitGui); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		//actions that are not switching views
		//player controls
		if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, togglePlayerState); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("podcast", gocui.KeyEnter, gocui.ModNone, playDownload); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("downloads", gocui.KeyEnter, gocui.ModNone, playDownload); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("searchResults", gocui.KeyEnter, gocui.ModNone, switchSubscribe); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		//switching views
		if err := g.SetKeybinding("subscribed", gocui.KeyArrowLeft, gocui.ModNone, switchListSearch); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("subscribed", gocui.KeyArrowRight, gocui.ModNone, switchListDownloads); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("subscribed", gocui.KeyEnter, gocui.ModNone, switchListPodcast); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("subscribed", gocui.KeyDelete, gocui.ModNone, switchRemoveSubscription); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("podcast", gocui.KeyArrowLeft, gocui.ModNone, switchListSubscribed); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("search", gocui.KeyArrowRight, gocui.ModNone, switchListSubscribed); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("search", gocui.KeyEnter, gocui.ModNone, switchKeyword); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("searchResults", gocui.KeyArrowRight, gocui.ModNone, switchListSubscribed); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		if err := g.SetKeybinding("downloads", gocui.KeyArrowLeft, gocui.ModNone, switchListSubscribed); err != nil {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		//main loop
		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
		}
		//clean up on exit
		close(stopTick)
	}
	globals.playerControl <- 5   //tell it to exit
	writeConfig(*globals.Config) //update config
	//wait for player to clean up
	<-playerExit
}
