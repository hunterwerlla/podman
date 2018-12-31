package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	"strings"
)

func enterPressedHome(configuration *Configuration) {

}

func enterPressedSearch(configuration *Configuration) {
	// TODO refactor into two functions
	if currentMode == Normal {
		if currentSelected >= len(currentPodcastsInBuffer) || currentSelected < 0 {
			return
		}
		selectedPodcast := currentPodcastsInBuffer[currentSelected]
		// check if it's already part of the configuration
		for _, thing := range configuration.Subscribed {
			if selectedPodcast.ArtistName == thing.ArtistName && selectedPodcast.CollectionName == thing.CollectionName {
				//already subscribed so do nothing
				return
			}
		}
		configuration.Subscribed = append(configuration.Subscribed, selectedPodcast) //now subscribe by adding it to the subscribed list
		writeConfig(configuration)
	} else {
		// search TODO use go()
		currentSelected = 0
		var err error
		searchString := strings.Replace(userTextBuffer, "\n", "", -1)
		searchString = strings.Trim(searchString, "\n\t")
		searchString = strings.Replace(searchString, " ", "+", -1) //replace spaces with plus to not break everything
		currentPodcastsInBuffer, err = searchItunes(searchString)
		if err != nil {
			userTextBuffer = "error searching! " + err.Error()
		}
	}
}

func enterPressedDownloaded(configuration *Configuration) {

}

func escapePressedHome(configuration *Configuration) {

}

func escapePressedSearch(configuration *Configuration) {

}

func escapePressedDownloaded(configuration *Configuration) {

}

func upPressedHome(configuration *Configuration) {
	if currentSelected > 0 {
		currentSelected--
	}
}

func upPressedSearch(configuration *Configuration) {
	if currentSelected > 0 {
		currentSelected--
	} else {
		currentMode = Insert
	}
}

func upPressedDownloaded(configuration *Configuration) {

}

func downPressedHome(configuration *Configuration) {
	if currentSelected < currentListSize-1 {
		currentSelected++
	}
}

func downPressedSearch(configuration *Configuration) {
	if currentSelected < currentListSize-1 {
		currentSelected++
	}
}

func downPressedDownloaded(configuration *Configuration) {

}

func searchPressedHome(configuration *Configuration) {
	userTextBuffer = ""
	currentMode = Insert
}

func searchPressedSearch(configuration *Configuration) {
	userTextBuffer = ""
	currentMode = Insert
}

func searchPressedDownloaded(configuration *Configuration) {
	userTextBuffer = ""
	currentMode = Insert
}

func handleKeyboard(configuration *Configuration, event ui.Event) {
	if event.ID == "<Enter>" {
		enterPressed[currentScreen](configuration)
		// reset mode on enter after the action is done
		currentMode = Normal
	} else if event.ID == "<Escape>" {
		// reset mode on Escape as well as possibly do something else
		currentMode = Normal
		escapePressed[currentScreen](configuration)
	} else if (event.ID == configuration.LeftKeybind && currentMode == Normal) || event.ID == "<Left>" {
		transitionScreen(leftTransitions, currentScreen)
	} else if (event.ID == configuration.RightKeybind && currentMode == Normal) || event.ID == "<Right>" {
		transitionScreen(rightTransitions, currentScreen)
	} else if currentMode == Insert {
		if event.ID == "<Backspace>" || event.ID == "<Delete>" || event.ID == "C-8>" {
			if len(userTextBuffer) > 0 {
				userTextBuffer = userTextBuffer[:len(userTextBuffer)-1]
			}
		} else if event.ID == "<Space>" {
			userTextBuffer += " "
		} else if event.ID == "<Tab>" {
			userTextBuffer += "	"
		} else if len(event.ID) > 0 && string([]rune(event.ID)[0]) == "<" {
			// do not do anything if it's one of the other control keys
		} else {
			userTextBuffer += event.ID
		}
	} else if event.ID == configuration.SearchKeybind {
		searchPressed[currentScreen](configuration)
		// TODO refactor these out
	} else if event.ID == configuration.PlayKeybind {
		fmt.Println("play")
	} else if event.ID == configuration.UpKeybind || event.ID == "<Up>" {
		upPressed[currentScreen](configuration)
	} else if event.ID == configuration.DownKeybind || event.ID == "<Down>" {
		downPressed[currentScreen](configuration)
	}
}
