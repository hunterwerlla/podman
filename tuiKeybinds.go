package main

import (
	ui "github.com/gizak/termui"
	"strings"
)

func switchToSelectedPodcastScreen(configuration *Configuration) {
	podcasts := getCurrentPagePodcasts()
	if currentSelected >= len(podcasts) || currentSelected < 0 {
		return
	}
	currentSelectedPodcast = podcasts[currentSelected]
	currentScreen = PodcastDetail
}

func searchPodcastsFromTui(configuration *Configuration) {
	// search TODO use go()
	currentSelected = 0
	var err error
	searchString := strings.Replace(userTextBuffer, "\n", "", -1)
	searchString = strings.Trim(searchString, "\n\t")
	searchString = strings.Replace(searchString, " ", "+", -1) //replace spaces with plus to not break everything
	currentPodcastsInBuffers[currentScreen], err = searchItunes(searchString)
	if err != nil {
		userTextBuffer = "error searching! " + err.Error()
	}
}

func enterPressedHome(configuration *Configuration) {
	// TODO refactor into two functions
	if currentMode == Normal {
		switchToSelectedPodcastScreen(configuration)
	}
}

func enterPressedSearch(configuration *Configuration) {
	if currentMode == Normal {
		switchToSelectedPodcastScreen(configuration)
	} else {
		searchPodcastsFromTui(configuration)
	}
}

func enterPressedDownloaded(configuration *Configuration) {
	podcasts := getCurrentPagePodcastEpisodes()
	if currentSelected >= len(podcasts) || currentSelected < 0 {
		return
	}
	SetPlaying(podcasts[currentSelected].StorageLocation)
	SetPlayerState(Play)
}

func enterPressedPodcastDetail(configuration *Configuration) {

}

func escapePressedHome(configuration *Configuration) {

}

func escapePressedSearch(configuration *Configuration) {

}

func escapePressedDownloaded(configuration *Configuration) {

}

func actionPressedHome(configuration *Configuration) {

}

func actionPressedSearch(configuration *Configuration) {
	subscribedKey := -1
	podcasts := getCurrentPagePodcasts()
	if currentSelected >= len(podcasts) || currentSelected < 0 {
		return
	}
	selectedPodcast := podcasts[currentSelected]
	// check if it's already part of the configuration
	for key, value := range configuration.Subscribed {
		if selectedPodcast.ArtistName == value.ArtistName && selectedPodcast.CollectionName == value.CollectionName {
			subscribedKey = key
		}
	}
	if subscribedKey != -1 {
		configuration.Subscribed = append(configuration.Subscribed[:subscribedKey], configuration.Subscribed[subscribedKey+1:]...)
	} else {
		configuration.Subscribed = append(configuration.Subscribed, selectedPodcast) //now subscribe by adding it to the subscribed list
	}
	writeConfig(configuration)
}

func actionPressedDownloaded(configuration *Configuration) {

}

func actionPressedPodcastDetail(configuration *Configuration) {

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
	if currentMode == Insert && len(userTextBuffer) > 0 {
		searchPodcastsFromTui(configuration)
	} else if currentSelected < currentListSize-1 {
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

func handleEventsGlobal(configuration *Configuration, event ui.Event) bool {
	if event.ID == "<Enter>" {
		enterPressed[currentScreen](configuration)
		// reset mode on enter after the action is done
		currentMode = Normal
	} else if event.ID == "<Escape>" {
		// reset mode on Escape as well as possibly do an action
		currentMode = Normal
		escapePressed[currentScreen](configuration)
	} else if (event.ID == configuration.LeftKeybind && currentMode == Normal) || event.ID == "<Left>" {
		transitionScreen(leftTransitions, currentScreen)
	} else if (event.ID == configuration.RightKeybind && currentMode == Normal) || event.ID == "<Right>" {
		transitionScreen(rightTransitions, currentScreen)
	} else if (event.ID == configuration.UpKeybind && currentMode == Normal) || event.ID == "<Up>" {
		upPressed[currentScreen](configuration)
	} else if (event.ID == configuration.DownKeybind && currentMode == Normal) || event.ID == "<Down>" {
		downPressed[currentScreen](configuration)
		currentMode = Normal
	} else if event.ID == configuration.PlayKeybind {
		TogglePlayerState()
	} else {
		// nothing matches, return false
		return false
	}
	// matches, return true
	return true
}

func handleKeyboard(configuration *Configuration, event ui.Event) {
	if handleEventsGlobal(configuration, event) {
		// handled
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
	} else if event.ID == configuration.ActionKeybind {
		actionPressed[currentScreen](configuration)
	}
}
