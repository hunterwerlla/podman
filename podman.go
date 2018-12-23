package main

import (
	"flag"
	"fmt"
	"github.com/jroimartin/gocui"
	"os/user"
	"time"
)

//global state
//TODO get rid of this whole awful thing
var (
	globals GlobalState = GlobalState{"", nil, nil, nil, -1, 0, nil, 0}
)

func main() {
	//get users home dir, the default storage
	usr, err := user.Current()
	defaultStorage := "."
	//if no error, sore in home directory
	if err == nil {
		defaultStorage = usr.HomeDir + "/" + "podman"
	}
	//make configuration struct that holds default settings
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
			end = CliCommand(globals.playerFile, globals.playerControl)
		}
		return
	}
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		fmt.Println(err)
		panic("Unable to start TUI, can atttempt to run --no-tui for minimal text based version")
	}
	defer g.Close()
	g.SetManagerFunc(guiHandler)
	setKeybinds(g)
	g.Mouse = true
	refreshGui(g)
	//main loop
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	writeConfig(*globals.Config)   //update config on exit
	globals.playerControl <- _exit //tell player to exit
	<-playerExit                   //wait for player to exit to finally exit
}

func refreshGui(g *gocui.Gui) {
	update := time.NewTicker(time.Millisecond * 500).C
	stopTick := make(chan bool)
	defer close(stopTick)
	go func() {
		for {
			select {
			case <-update:
				g.Update(guiHandler)
			case <-stopTick:
				return
			}
		}
	}()
}

func setKeybinds(g *gocui.Gui) {
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
	//player controls
	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, togglePlayerState); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, skipPlayerForward); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, skipPlayerBackward); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	//actions that are not switching views
	if err := g.SetKeybinding("", gocui.KeyDelete, gocui.ModNone, switchDeleteDownloaded); err != nil {
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
}
