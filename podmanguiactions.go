//this holds the tui functions and information
package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"os"
	"sort"
	"strings"
	"time"
)

func guiHandler(g *gocui.Gui) error {
	if stateView == 0 {
		listSubscribed(g)
	} else if stateView == 1 {
		listPodcast(g)
	} else if stateView == 2 {
		listSearch(g)
	} else if stateView == 3 {
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
		if stateView == 0 {
			//starts at 1
			if y >= len(globals.Config.Subscribed[scrollingOffset:]) {
				return nil
			}
		} else if stateView == 1 || stateView == 3 {
			//starts at 0
			if y >= len(selectedPodcastEntries[scrollingOffset:])-1 {
				return nil
			}
		} else if stateView == 2 {
			//never allow scroll down on search, only allow transitioning view
			if g.CurrentView().Name() == "search" {
				if len(selectedPodcastSearch) > 0 {
					g.SetCurrentView("searchResults")
					yCursorOffset = 0
					return nil
				} else { //else don't allow scrolling down
					return nil
				}
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
		if stateView == 0 {
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
		} else if stateView == 1 || stateView == 3 {
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
		} else if stateView == 2 {
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
					g.SetCurrentView("search")
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
	if len(globals.Config.Subscribed) == 0 {
		return nil //TODO return an actual error
	}
	if position > len(globals.Config.Subscribed) {
		return nil
	}
	selectedPodcast = globals.Config.Subscribed[position-1] //select the podcast put in memory
	selectedPodcastEntries = nil                            //now delete the cache
	selectedPodcastSearch = nil
	scrollingOffset = 0
	//change layout
	stateView = 1
	//delete old views
	g.DeleteView("subscribed")
	return nil
}
func switchListSubscribed(g *gocui.Gui, v *gocui.View) error {
	yCursorOffset = 0 //reset cursor
	scrollingOffset = 0
	//change layout
	stateView = 0
	//delete other views
	g.DeleteView("subscribed")
	g.DeleteView("podcast")
	g.DeleteView("downloads")
	g.DeleteView("podcastDescription")
	listSubscribed(g)
	return nil
}
func switchListSearch(g *gocui.Gui, v *gocui.View) error {
	yCursorOffset = 0 //rest cursor
	scrollingOffset = 0
	stateView = 2 //2 is search
	g.DeleteView("subscribed")
	g.DeleteView("podcast")
	g.DeleteView("downloads")
	g.DeleteView("podcastDescription")
	listSearch(g)
	g.SetCurrentView("search")
	return nil
}

func switchListDownloads(g *gocui.Gui, v *gocui.View) error {
	yCursorOffset = 0 //rest cursor
	scrollingOffset = 0
	stateView = 3 //3 is downloads
	g.DeleteView("subscribed")
	g.DeleteView("podcast")
	g.DeleteView("downloads")
	g.DeleteView("podcastDescription")
	listDownloaded(g)
	g.SetCurrentView("downloads")
	//now sort the map and put in selectedPodcastEntries
	var tmp []PodcastEntry
	for _, thing := range globals.Config.Downloaded {
		tmp = append(tmp, thing)
	}
	sort.Sort(PodcastEntrySlice(tmp))
	selectedPodcastEntries = tmp
	return nil
}
func switchKeyword(g *gocui.Gui, v *gocui.View) error {
	queue := v.ViewBuffer()
	queue = strings.Replace(queue, "\n", "", -1)
	queue = strings.Trim(queue, "\n\t ")
	queue = strings.Replace(queue, " ", "+", -1) //replace spaces with plus to not break everything
	podcasts, err := searchItunes(queue)
	if err != nil {
		fmt.Fprintln(v, "error searching! %s", err.Error())
		return nil
	}
	//clear the buffer
	v.Clear()
	selectedPodcastSearch = podcasts
	g.SetCurrentView("searchResults")
	return nil
}

func switchSubscribe(g *gocui.Gui, v *gocui.View) error {
	_, position := v.Cursor() //get cursor position to select
	if len(selectedPodcastSearch) == 0 {
		return nil //TODO return an actual error
	}
	selectedPodcast = selectedPodcastSearch[position-1] //select the podcast put in memory
	//now check if already added
	for _, thing := range globals.Config.Subscribed {
		if selectedPodcast.ArtistName == thing.ArtistName && selectedPodcast.CollectionName == thing.CollectionName {
			//already subscribed
			return nil
		}
	}
	globals.Config.Subscribed = append(globals.Config.Subscribed, selectedPodcast) //now subscribe by adding it to the subscribed list
	writeConfig(*globals.Config)
	return nil
}
func playDownload(g *gocui.Gui, v *gocui.View) error {
	_, position := v.Cursor() //get cursor position to select
	var toPlay PodcastEntry
	guid := selectedPodcastEntries[position].GUID
	if isDownloaded(selectedPodcastEntries[position]) == false {
		download(*globals.Config, selectedPodcast, selectedPodcastEntries[position])
	} else {
		if thing := globals.Config.Downloaded[guid]; thing != (PodcastEntry{}) {
			toPlay = thing
		} else {
			return nil //TODO real error
		}
		//now play
		if toPlay := toPlay.StorageLocation; toPlay != "" {
			globals.playerFile <- toPlay
		}
	}
	return nil
}

func switchRemoveSubscription(g *gocui.Gui, v *gocui.View) error {
	_, position := v.Cursor() //get cursor position to select
	if len(globals.Config.Subscribed) == 0 {
		return nil //TODO return an actual error
	}
	if position > len(globals.Config.Subscribed) {
		return nil
	}
	globals.Config.Subscribed = append(globals.Config.Subscribed[0:position-1], globals.Config.Subscribed[position:]...)
	return nil
}
func switchDeleteDownloaded(g *gocui.Gui, v *gocui.View) error {
	_, position := v.Cursor()             //get cursor position to select
	if stateView == 1 || stateView == 3 { //in subscribed is very different from in download list
		if isDownloaded(selectedPodcastEntries[position]) {
			//remove entry in list, then remove entry on disk
			toDelete, ok := globals.Config.Downloaded[selectedPodcastEntries[position].GUID]
			if ok {
				os.Remove(toDelete.StorageLocation)
				delete(globals.Config.Downloaded, toDelete.GUID)
			}
		}
		//update if stateview is downloads, update due to custom sort
		if stateView == 3 {
			switchListDownloads(g, v)
		}
	}
	return nil
}

func quitGui(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func togglePlayerState(g *gocui.Gui, v *gocui.View) error {
	//pause so will not enter invalid state
	time.Sleep(time.Millisecond * 50)
	if globals.playerState == 1 {
		globals.playerControl <- 0
		globals.playerState = 0
	} else if globals.playerState == 0 {
		globals.playerControl <- 1
		globals.playerState = 1
	}
	return nil
}

func skipPlayerForward(g *gocui.Gui, v *gocui.View) error {
	//pause so will not enter invalid state
	time.Sleep(time.Millisecond * 50)
	globals.playerControl <- 3
	return nil
}

func skipPlayerBackward(g *gocui.Gui, v *gocui.View) error {
	//pause so will not enter invalid state
	time.Sleep(time.Millisecond * 50)
	globals.playerControl <- 4
	return nil
}