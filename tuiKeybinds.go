package main

import (
	"fmt"
	ui "github.com/gizak/termui"
)

func enterPressedHome(configuration *Configuration) {

}

func enterPressedSearch(configuration *Configuration) {

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

}

func upPressedDownloaded(configuration *Configuration) {

}

func downPressedHome(configuration *Configuration) {
	if currentSelected < currentListSize-1 {
		currentSelected++
	}
}

func downPressedSearch(configuration *Configuration) {

}

func downPressedDownloaded(configuration *Configuration) {

}

func searchPressedHome(configuration *Configuration) {
	currentMode = Insert
}

func searchPressedSearch(configuration *Configuration) {
	currentMode = Insert
}

func searchPressedDownloaded(configuration *Configuration) {
	currentMode = Insert
}

func handleKeyboard(configuration *Configuration, event ui.Event) {
	if event.ID == "<Enter>" {
		enterPressed[currentScreen](configuration)
	}
	if event.ID == "<Escape>" {
		// reset mode on Escape as well as possibly do something else
		currentMode = Normal
		escapePressed[currentScreen](configuration)
	}
	if currentMode == Insert {
		if event.ID == "<Backspace>" || event.ID == "<Delete>" {
			userTextBuffer = userTextBuffer[:len(userTextBuffer)-1]
		} else {
			userTextBuffer += event.ID
		}
	}
	if event.ID == configuration.SearchKeybind {
		// Clear buffer and set insert every time
		userTextBuffer = ""
		currentMode = Insert
		searchPressed[currentScreen](configuration)
	}
	if event.ID == configuration.LeftKeybind || event.ID == "<Left>" {
		transitionScreen(leftTransitions, currentScreen)
	}
	if event.ID == configuration.RightKeybind || event.ID == "<Right>" {
		transitionScreen(rightTransitions, currentScreen)
	}
	// TODO refactor these out
	if event.ID == configuration.PlayKeybind {
		fmt.Println("play")
	}
	if event.ID == configuration.UpKeybind || event.ID == "<Up>" {
		upPressed[currentScreen](configuration)
	}
	if event.ID == configuration.DownKeybind || event.ID == "<Down>" {
		downPressed[currentScreen](configuration)
	}
}
