package main

import (
	"flag"
	"fmt"
	"github.com/jroimartin/gocui"
	"os/user"
	"time"
)

// Podcast Holds one podcast
type Podcast struct {
	ArtistName     string
	CollectionName string
	FeedURL        string
	Description    string
}

// PodcastEpisode Holds one episode of a podcast
type PodcastEpisode struct {
	// The title of the podcast according to ITunes
	PodcastTitle    string
	Title           string
	Summary         string
	Link            string
	Content         string
	GUID            string
	StorageLocation string
}

// CachedPodcast a struct that holds podcasts cached in memory/on disk
type CachedPodcast struct {
	Type     Podcast
	Podcasts []PodcastEpisode
	Checked  time.Time
}

// ItunesSearch returns an array of podcasts that are found by itunes
type ItunesSearch struct {
	Results []Podcast
}

// PodcastEpisodeSlice takes a slice of a list of podcasts
type PodcastEpisodeSlice []PodcastEpisode

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
	config := CreateDefault()
	config.StorageLocation = defaultStorage
	config = ReadConfig(config)
	//read command line flags
	noTui := flag.Bool("no-tui", false, "Select whether to use the GUI or not")
	flag.Parse()
	//make the channels used by player
	StartPlayer()
	//made a decision to use TUI or not
	if *noTui == true {
		runCui(&config)
	} else {
		runTui(&config)
	}
}

func runCui(config *Configuration) {
	end := false
	for end != true {
		end = CliCommand(config)
	}
}

func runTui(config *Configuration) {
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
	WriteConfig(config) //update config on exit
	DisposePlayer()            //tell player to exit + wait
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
	if err := g.SetKeybinding("searchResults", gocui.KeyEnter, gocui.ModNone, actionSubscribe); err != nil {
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
