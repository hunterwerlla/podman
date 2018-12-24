package main

import (
	"flag"
	"fmt"
	"github.com/hunterwerlla/podman/configuration"
	"github.com/hunterwerlla/podman/player"
	"github.com/jroimartin/gocui"
	"os/user"
	"time"
)

func main() {
	//get users home dir, the default storage
	// TODO fix this
	defaultStorage := "."
	usr, err := user.Current()
	//if no error, store in home directory
	if err == nil {
		defaultStorage = usr.HomeDir + "/.config/podman"
	}
	//read config file
	config := configuration.CreateDefault()
	config.StorageLocation = defaultStorage
	config = configuration.ReadConfig(config)
	//read command line flags
	noTui := flag.Bool("no-tui", false, "Select whether to use the GUI or not")
	flag.Parse()
	//make the channels used by player
	player.StartPlayer()
	//made a decision to use TUI or not
	if *noTui == true {
		runCui(&config)
	} else {
		runTui(&config)
	}
}

func runCui(config *configuration.Configuration) {
	end := false
	for end != true {
		end = CliCommand(config)
	}
}

func runTui(config *configuration.Configuration) {
	SetTuiConfiguration(config)
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
	configuration.WriteConfig(config) //update config on exit
	player.DisposePlayer()            //tell player to exit + wait
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
	if err := g.SetKeybinding("search", gocui.KeyEnter, gocui.ModNone, searchKeyword); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding("searchResults", gocui.KeyArrowRight, gocui.ModNone, switchListSubscribed); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding("downloads", gocui.KeyArrowLeft, gocui.ModNone, switchListSubscribed); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
}
