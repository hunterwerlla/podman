package main

import (
	ui "github.com/gizak/termui"
	"strings"
	"time"
)

func switchToSelectedPodcastScreen(configuration *Configuration) {
	podcasts := getCurrentPagePodcasts()
	cursor := getCurrentCursorPosition()
	if cursor >= len(podcasts) || cursor < 0 {
		return
	}
	currentSelectedPodcast = podcasts[cursor]
	currentScreen = PodcastDetail
	// reset cursor
	setCurrentCursorPosition(0)
}

func searchPodcastsFromTui(configuration *Configuration) {
	// search TODO use go()
	setCurrentCursorPosition(0)
	var err error
	searchString := strings.Replace(userTextBuffer, "\n", "", -1)
	searchString = strings.Trim(searchString, "\n\t")
	searchString = strings.Replace(searchString, " ", "+", -1) //replace spaces with plus to not break everything
	currentPodcastsInBuffers[currentScreen], err = searchItunes(searchString)
	if err != nil {
		userTextBuffer = "error searching! " + err.Error()
	}
}

func doNothingWithInput(configuration *Configuration) {}

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
	cursor := getCurrentCursorPosition()
	if cursor >= len(podcasts) || cursor < 0 {
		return
	}
	SetPlaying(podcasts[cursor].StorageLocation)
	sendPlayerMessage(PlayerPlay)
}

func enterPressedPodcastDetail(configuration *Configuration) {
	podcasts := getCurrentPagePodcastEpisodes()
	cursor := getCurrentCursorPosition()
	if cursor >= len(podcasts) || cursor < 0 {
		return
	}
	if podcastIsDownloaded(configuration, podcasts[cursor]) {
		location := getPodcastLocation(configuration, podcasts[cursor])
		if location != "" {
			SetPlaying(location)
			sendPlayerMessage(PlayerPlay)
		}
		return
	}
	// TODO fix this race condition/bad configuration management.
	go func() {
		_ = downloadPodcast(configuration, currentSelectedPodcast, podcasts[cursor])
	}()
}

func actionPressedSearch(configuration *Configuration) {
	subscribedKey := -1
	podcasts := getCurrentPagePodcasts()
	cursor := getCurrentCursorPosition()
	if cursor >= len(podcasts) || cursor < 0 {
		return
	}
	selectedPodcast := podcasts[cursor]
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

func deletePodcastSelectedByCursor(configuration *Configuration) {
	cursor := getCurrentCursorPosition()
	podcasts := getCurrentPagePodcastEpisodes()
	if cursor >= len(podcasts) || cursor < 0 {
		return
	}
	deleteDownloadedPodcast(configuration, podcasts[cursor])
	// reset cursor if needed
	if cursor == len(podcasts)-1 {
		setCurrentCursorPosition(len(podcasts) - 2)
	}
	writeConfig(configuration)
}

func upPressedGeneric(configuration *Configuration) {
	cursor := getCurrentCursorPosition()
	if cursor > 0 {
		setCurrentCursorPosition(cursor - 1)
	}
}

func upPressedSearch(configuration *Configuration) {
	cursor := getCurrentCursorPosition()
	if cursor > 0 {
		setCurrentCursorPosition(cursor - 1)
	} else {
		currentMode = Insert
	}
}

func downPressedGeneric(configuration *Configuration) {
	cursor := getCurrentCursorPosition()
	if cursor < currentListSize-1 {
		setCurrentCursorPosition(cursor + 1)
	}
}

func downPressedSearch(configuration *Configuration) {
	cursor := getCurrentCursorPosition()
	if currentMode == Insert && len(userTextBuffer) > 0 {
		searchPodcastsFromTui(configuration)
	} else if cursor < currentListSize-1 {
		setCurrentCursorPosition(cursor + 1)
	}
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

func deletePressedHome(configuration *Configuration) {
	subscribedKey := -1
	podcasts := getCurrentPagePodcasts()
	cursor := getCurrentCursorPosition()
	if cursor >= len(podcasts) || cursor < 0 {
		return
	}
	for key, value := range configuration.Subscribed {
		if podcasts[cursor].ArtistName == value.ArtistName && podcasts[cursor].CollectionName == value.CollectionName {
			subscribedKey = key
		}
	}
	if subscribedKey == -1 {
		return
	}
	// reset cursor if needed
	if cursor == len(configuration.Subscribed)-1 {
		setCurrentCursorPosition(len(configuration.Subscribed) - 2)
	}
	configuration.Subscribed = append(configuration.Subscribed[:subscribedKey], configuration.Subscribed[subscribedKey+1:]...)
	writeConfig(configuration)
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
	} else if (event.ID == configuration.FastForward && currentMode == Normal) || event.ID == "<Previous>" {
		sendPlayerMessage(PlayerFastForward)
	} else if (event.ID == configuration.Rewind && currentMode == Normal) || event.ID == "<Next>" {
		sendPlayerMessage(PlayerRewind)
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
	} else if event.ID == configuration.PlayKeybind {
		TogglePlayerState()
	} else if event.ID == configuration.SearchKeybind {
		searchPressed[currentScreen](configuration)
	} else if event.ID == configuration.ActionKeybind {
		actionPressed[currentScreen](configuration)
	} else if event.ID == configuration.DeleteKeybind {
		deletePressed[currentScreen](configuration)
		prepareDrawPage[currentScreen](configuration)
	}
}

func handleMouse(configuration *Configuration, event ui.Event) {
	if event.ID == "<MouseWheelUp>" {
		upPressed[currentScreen](configuration)
	} else if event.ID == "<MouseWheelDown>" {
		downPressed[currentScreen](configuration)
	} else if event.ID == "<MouseLeft>" {
		enterPressed[currentScreen](configuration)
	} else if event.ID == "<MouseRight>" && currentScreen == PodcastDetail {
		transitionScreen(leftTransitions, currentScreen)
	}
}

func tuiMainLoop(configuration *Configuration) {
	width := ui.TermWidth()
	height := ui.TermHeight()

	prepareDrawPage[currentScreen](configuration)
	ui.Render(drawPage[currentScreen](configuration, width, height)...)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C

	for {
		select {
		case e := <-uiEvents:
			{
				savedScreen := currentScreen
				if e.Type == ui.KeyboardEvent {
					if e.ID == "<C-c>" {
						goto exitMainLoop
					} else {
						handleKeyboard(configuration, e)
					}
				} else if e.Type == ui.MouseEvent {
					handleMouse(configuration, e)
					ui.Render(drawPage[currentScreen](configuration, width, height)...)
				} else if e.Type == ui.ResizeEvent {
					payload := e.Payload.(ui.Resize)
					width = payload.Width
					height = payload.Height
					ui.Render(drawPage[currentScreen](configuration, width, height)...)
				}
				// refresh screen after keyboard input or redraw screen entirely + reset state if we have changed screens
				if savedScreen != currentScreen {
					prepareDrawPage[currentScreen](configuration)
					// reset modes
					if currentScreen == Search && (savedScreen != PodcastDetail) {
						currentMode = Insert
					} else {
						currentMode = Normal
					}
					// save last screen
					previousScreen = savedScreen
				}
				ui.Render(drawPage[currentScreen](configuration, width, height)...)
			}
		case <-ticker:
			// refresh player or reset if needed
			if GetPlayerState() == PlayerPlay || downloadInProgress() {
				ui.Render(producePlayerWidget(configuration, width, height))
			} else if GetPlayerState() == PlayerStop {
				sendPlayerMessage(PlayerNothingPlaying)
				ui.Render(producePlayerWidget(configuration, width, height))
			}
			if refreshPage[currentScreen] != nil {
				ui.Render(refreshPage[currentScreen](configuration, width, height)...)
			}
		}
	}
exitMainLoop:
}
