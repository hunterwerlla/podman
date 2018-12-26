//this file holds the tui action functions
package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"os"
	"sort"
	"strings"
)

func SetTuiKeybinds(g *gocui.Gui) {
	//global keybinds
	if err := g.SetKeybinding(ScreenAll, gocui.KeyCtrlC, gocui.ModNone, quitGui); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding(ScreenAll, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding(ScreenAll, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	//player controls
	if err := g.SetKeybinding(ScreenAll, gocui.KeySpace, gocui.ModNone, togglePlayerState); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding(ScreenAll, gocui.KeyPgup, gocui.ModNone, skipPlayerForward); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding(ScreenAll, gocui.KeyPgdn, gocui.ModNone, skipPlayerBackward); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	//actions taken
	if err := g.SetKeybinding(ScreenAll, gocui.KeyDelete, gocui.ModNone, switchDeleteDownloaded); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding(ScreenPodcast, gocui.KeyEnter, gocui.ModNone, playDownload); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding(ScreenDownloads, gocui.KeyEnter, gocui.ModNone, playDownload); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding(ScreenSearchResults, gocui.KeyEnter, gocui.ModNone, subscribeToPodcast); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding(ScreenSearch, gocui.KeyEnter, gocui.ModNone, searchKeyword); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	//switching views
	if err := g.SetKeybinding(ScreenSubscribed, gocui.KeyArrowLeft, gocui.ModNone, switchListSearch); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding(ScreenSubscribed, gocui.KeyArrowRight, gocui.ModNone, switchListDownloads); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding(ScreenSubscribed, gocui.KeyEnter, gocui.ModNone, switchListPodcast); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding(ScreenSubscribed, gocui.KeyDelete, gocui.ModNone, switchRemoveSubscription); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding(ScreenPodcast, gocui.KeyArrowLeft, gocui.ModNone, switchListSubscribed); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding(ScreenSearch, gocui.KeyArrowRight, gocui.ModNone, switchListSubscribed); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding(ScreenSearchResults, gocui.KeyArrowRight, gocui.ModNone, switchListSubscribed); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
	if err := g.SetKeybinding(ScreenDownloads, gocui.KeyArrowLeft, gocui.ModNone, switchListSubscribed); err != nil {
		panic(fmt.Sprintf("Error in GUI, have to exit %s", err.Error()))
	}
}

func TuiHandler(g *gocui.Gui) error {
	if stateView == ScreenSubscribed {
		listSubscribed(g)
	} else if stateView == ScreenPodcast {
		listPodcast(g)
	} else if stateView == ScreenSearch {
		listSearch(g)
	} else if stateView == ScreenDownloads {
		listDownloaded(g)
	}
	printPlayer(g)
	return nil
}

//TODO add scrolling beyond screen and not crash
func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, maxY := v.Size()
		x, y := v.Cursor()
		if stateView == ScreenSubscribed {
			//starts at 1
			if y >= len(config.Subscribed[scrollingOffset:]) {
				return nil
			}
		} else if stateView == ScreenPodcast || stateView == ScreenDownloads {
			//starts at 0
			if y >= len(selectedPodcastEntries[scrollingOffset:])-1 {
				return nil
			}
		} else if stateView == ScreenSearch {
			//never allow scroll down on search, only allow transitioning view
			if g.CurrentView().Name() == ScreenSearch {
				if len(selectedPodcastSearch) > 0 {
					g.SetCurrentView(ScreenSearchResults)
					yCursorOffset = 0
					return nil
				}
				// don't allow scrolling down
				return nil
			}
			//starts at 1
			if y >= len(selectedPodcastSearch[scrollingOffset:]) {
				return nil
			}
		} else { //unknown state TODO return error
			return nil
		}
		//go to another page
		if y == maxY-1 {
			scrollingOffset += maxY //add height
			yCursorOffset = 0
			if err := v.SetCursor(x, 0); err != nil {
				return err
			}
			return nil
		}
		yCursorOffset++
		if err := v.SetCursor(x, y+1); err != nil {
			return err
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, maxY := v.Size()
		x, y := v.Cursor()
		if stateView == ScreenSubscribed {
			//if Y is 1 at the top, so don't move up again
			if y == 1 {
				if scrollingOffset != 0 {
					//this seems like some magic numbers but it's due to the way the cursor is updated
					yCursorOffset = maxY - 2
					//NOTE if this ever breaks, it's because the library changed something with cursor updating
					if err := v.SetCursor(x, yCursorOffset); err != nil {
						return err
					}
					scrollingOffset -= maxY //subtract height
					//make sure it's not negative
					if scrollingOffset < 0 {
						scrollingOffset = 0
					}
				}
				return nil
			}
		} else if stateView == ScreenPodcast || stateView == ScreenDownloads {
			if y == 0 {
				if scrollingOffset != 0 {
					yCursorOffset = maxY - 1
					if err := v.SetCursor(x, 1); err != nil {
						return err
					}
					scrollingOffset -= maxY //subtract height
					//make sure it's not negative
					if scrollingOffset < 0 {
						scrollingOffset = 0
					}
				}
				return nil
			}
		} else if stateView == ScreenSearch {
			if y < 2 { //y==0 included because search bar
				if scrollingOffset != 0 {
					//NOTE this is the same magic as list podcast
					yCursorOffset = maxY - 2
					if err := v.SetCursor(x, yCursorOffset); err != nil {
						return err
					}
					scrollingOffset -= maxY //subtract height
					//make sure it's not negative
					if scrollingOffset < 0 {
						scrollingOffset = 0
					}
					return nil
				}
				if y == 1 || y == 0 { //if y is 0 set active view to search bar
					g.SetCurrentView(ScreenSearch)
				}
				return nil
			}
		} else {
			return nil
		}
		if y <= 0 || yCursorOffset <= 0 {
			return nil
		}
		yCursorOffset--
		if err := v.SetCursor(x, y-1); err != nil {
			return err
		}
	}
	return nil
}

//the second view a list podcast
//TODO add real error
func switchListPodcast(g *gocui.Gui, v *gocui.View) error {
	_, position := v.Cursor() //get cursor position to select
	yCursorOffset = 0         //reset cursor
	if len(config.Subscribed) == 0 {
		return nil //TODO return an actual error
	}
	if position > len(config.Subscribed) {
		return nil
	}
	selectedPodcast = config.Subscribed[position-1] //select the podcast put in memory
	selectedPodcastEntries = nil                    //now delete the cache
	selectedPodcastSearch = nil
	scrollingOffset = 0
	//change layout
	stateView = ScreenPodcast
	//delete old views
	g.DeleteView(ScreenSubscribed)
	return nil
}

func switchListSubscribed(g *gocui.Gui, v *gocui.View) error {
	yCursorOffset = 0 //reset cursor
	scrollingOffset = 0
	//change layout
	stateView = ScreenSubscribed
	//delete other views
	g.DeleteView(ScreenSubscribed)
	g.DeleteView(ScreenPodcast)
	g.DeleteView("downloads")
	g.DeleteView("podcastDescription")
	listSubscribed(g)
	return nil
}

func switchListSearch(g *gocui.Gui, v *gocui.View) error {
	yCursorOffset = 0 //rest cursor
	scrollingOffset = 0
	stateView = ScreenSearch
	g.DeleteView(ScreenSubscribed)
	g.DeleteView(ScreenPodcast)
	g.DeleteView("downloads")
	g.DeleteView("podcastDescription")
	listSearch(g)
	g.SetCurrentView(ScreenSearch)
	return nil
}

func switchListDownloads(g *gocui.Gui, v *gocui.View) error {
	yCursorOffset = 0 //rest cursor
	scrollingOffset = 0
	stateView = ScreenDownloads
	g.DeleteView(ScreenSubscribed)
	g.DeleteView(ScreenPodcast)
	g.DeleteView("downloads")
	g.DeleteView("podcastDescription")
	listDownloaded(g)
	g.SetCurrentView("downloads")
	//now sort the map and put in selectedPodcastEntries
	var tmp []PodcastEpisode
	for _, thing := range config.Downloaded {
		tmp = append(tmp, thing)
	}
	sort.Sort(PodcastEpisodeSlice(tmp))
	selectedPodcastEntries = tmp
	return nil
}
func searchKeyword(g *gocui.Gui, v *gocui.View) error {
	searchQuery := v.ViewBuffer()
	searchQuery = strings.Replace(searchQuery, "\n", "", -1)
	searchQuery = strings.Trim(searchQuery, "\n\t ")
	searchQuery = strings.Replace(searchQuery, " ", "+", -1) //replace spaces with plus to not break everything
	podcasts, err := searchItunes(searchQuery)
	if err != nil {
		fmt.Fprintf(v, "error searching! %s", err.Error())
		return nil
	}
	//clear the buffer
	v.Clear()
	selectedPodcastSearch = podcasts
	g.SetCurrentView(ScreenSearchResults)
	return nil
}

func subscribeToPodcast(g *gocui.Gui, v *gocui.View) error {
	_, position := v.Cursor() //get cursor position to select
	if len(selectedPodcastSearch) == 0 {
		return nil //TODO return an actual error
	}
	selectedPodcast = selectedPodcastSearch[position-1] //select the podcast put in memory
	//now check if already added
	for _, thing := range config.Subscribed {
		if selectedPodcast.ArtistName == thing.ArtistName && selectedPodcast.CollectionName == thing.CollectionName {
			//already subscribed
			return nil
		}
	}
	config.Subscribed = append(config.Subscribed, selectedPodcast) //now subscribe by adding it to the subscribed list
	WriteConfig(config)
	return nil
}

func switchRemoveSubscription(g *gocui.Gui, v *gocui.View) error {
	_, position := v.Cursor() //get cursor position to select
	if len(config.Subscribed) == 0 {
		return nil //TODO return an actual error
	}
	if position > len(config.Subscribed) {
		return nil
	}
	item := config.Subscribed[position-1]
	//now remove from cache
	for i, thing := range config.Cached {
		if thing.Type.ArtistName == item.ArtistName && thing.Type.CollectionName == item.CollectionName {
			config.Cached = append(config.Cached[0:i], config.Cached[i+1:]...)
			break
		}
	}
	config.Subscribed = append(config.Subscribed[0:position-1], config.Subscribed[position:]...)
	return nil
}
func switchDeleteDownloaded(g *gocui.Gui, v *gocui.View) error {
	_, position := v.Cursor()                                       //get cursor position to select
	if stateView == ScreenSubscribed || stateView == ScreenDownloads { //in subscribed is very different from in download list
		if PodcastIsDownloaded(selectedPodcastEntries[position]) {
			//remove entry in list, then remove entry on disk
			toDelete, ok := config.Downloaded[selectedPodcastEntries[position].GUID]
			if ok {
				os.Remove(toDelete.StorageLocation)
				delete(config.Downloaded, toDelete.GUID)
			}
		}
		//update if stateview is downloads, update due to custom sort
		if stateView == ScreenDownloads {
			switchListDownloads(g, v)
		}
	}
	return nil
}

func playDownload(g *gocui.Gui, v *gocui.View) error {
	_, yPosition := v.Cursor() //get cursor position to select
	if len(selectedPodcastEntries) <= yPosition {
		return nil
	}
	guid := selectedPodcastEntries[yPosition].GUID
	if PodcastIsDownloaded(selectedPodcastEntries[yPosition]) == false {
		go func() {
			config, err := DownloadPodcast(config, selectedPodcast, selectedPodcastEntries[yPosition], g)
			if err != nil {
				// TODO rasise a UI error
				return
			}
			WriteConfig(config)
		}() //download async
	} else {
		if podcast := config.Downloaded[guid]; podcast != (PodcastEpisode{}) { //if it is not empty
			SetPlaying(podcast.StorageLocation)
			SetPlayerState(Play)
		}
	}
	return nil
}

func togglePlayerState(g *gocui.Gui, v *gocui.View) error {
	if GetPlayerState() == Play {
		SetPlayerState(Pause)
	} else if GetPlayerState() == Pause {
		SetPlayerState(Play)
	}
	// needed for gocui
	return nil
}

func skipPlayerForward(g *gocui.Gui, v *gocui.View) error {
	SetPlayerState(FastForward)
	return nil
}

func skipPlayerBackward(g *gocui.Gui, v *gocui.View) error {
	SetPlayerState(Rewind)
	return nil
}

func quitGui(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
