package main

import (
	"fmt"
	ui "github.com/gizak/termui"
)

type screen string
type mode string

const (
	None          screen = "None"
	LastScreen    screen = "LastScreen"
	Home          screen = "Home"
	Search        screen = "Search"
	Downloaded    screen = "Downloaded"
	PodcastDetail screen = "PodcastDetail"
)

const (
	Normal mode = "Normal"
	Insert mode = "Insert"
)

var (
	defaultControlsMap = map[screen]string{
		Home:          "[%s]elect/[<enter>]  [h]left(search)   [l]right(downloaded)",
		Search:        "[s]ubscribe/unsubscribe   [/]search   [esc]ape searching   [<enter>]%s   [j]down   [k]up   [l]right(home)",
		Downloaded:    "[p]lay/<enter>",
		PodcastDetail: "[<enter>] download episode",
	}

	controlsMap = make(map[screen]string)

	leftTransitions = map[screen]screen{
		Home:          Search,
		Search:        None,
		Downloaded:    Home,
		PodcastDetail: LastScreen,
	}

	rightTransitions = map[screen]screen{
		Home:          Downloaded,
		Downloaded:    None,
		Search:        Home,
		PodcastDetail: LastScreen,
	}

	drawPage = map[screen]func(configuration *Configuration, width int, height int) []ui.Bufferer{
		Home:          drawPageMain,
		Search:        drawPageSearch,
		Downloaded:    drawPageDownloaded,
		PodcastDetail: drawPagePodcastDetail,
	}

	refreshPage = map[screen]func(configuration *Configuration, width int, height int) []ui.Bufferer{
		Home:       refreshPageMain,
		Search:     drawPageSearch,
		Downloaded: refreshPageDownloaded,
	}

	actionPressed = map[screen]func(configuration *Configuration){
		Home:          actionPressedHome,
		Search:        actionPressedSearch,
		Downloaded:    actionPressedDownloaded,
		PodcastDetail: actionPressedPodcastDetail,
	}

	enterPressed = map[screen]func(configuration *Configuration){
		Home:          enterPressedHome,
		Search:        enterPressedSearch,
		Downloaded:    enterPressedDownloaded,
		PodcastDetail: enterPressedPodcastDetail,
	}

	escapePressed = map[screen]func(configuration *Configuration){
		Home:       escapePressedHome,
		Search:     escapePressedSearch,
		Downloaded: escapePressedDownloaded,
	}

	upPressed = map[screen]func(configuration *Configuration){
		Home:       upPressedHome,
		Search:     upPressedSearch,
		Downloaded: upPressedDownloaded,
	}

	downPressed = map[screen]func(configuration *Configuration){
		Home:       downPressedHome,
		Search:     downPressedSearch,
		Downloaded: downPressedDownloaded,
	}

	searchPressed = map[screen]func(configuration *Configuration){
		Home:       searchPressedHome,
		Search:     searchPressedSearch,
		Downloaded: searchPressedDownloaded,
	}

	currentPodcastsInBuffers = map[screen]interface{}{
		Home:          nil,
		Search:        nil,
		Downloaded:    nil,
		PodcastDetail: nil,
	}
)

var (
	currentSelected         = 0
	currentListSize         = 0
	currentListOffset       = 0
	currentMode             = Normal
	currentScreen           = Home
	previousScreen          = None
	currentPodcastsInBuffer []Podcast
	currentSelectedPodcast  Podcast
	userTextBuffer          = ""
)

func fillOutControlsMap(configuration *Configuration, controls map[screen]string) {
	var searchText string
	controlsMap[Home] = fmt.Sprintf(controls[Home], configuration.ActionKeybind)
	if currentMode == Insert {
		searchText = "finish search  "
	} else {
		searchText = "look at podcast"
	}
	controlsMap[Search] = fmt.Sprintf(controls[Search], searchText)
	controlsMap[Downloaded] = controls[Downloaded]
}

func termuiStyleText(text string, fgcolor string, bgcolor string) string {
	text = "[" + text + "](fg-" + fgcolor + ",bg-" + string(bgcolor) + ")"
	return text
}

func transitionScreen(transitions map[screen]screen, screen screen) {
	if transitions[screen] == None {
		// do nothing
	} else if transitions[screen] == LastScreen {
		currentScreen = previousScreen
		previousScreen = None
	} else {
		currentScreen = transitions[screen]
	}
}

// StartTui starts the TUI with the Configuration passed in
func StartTui(configuration *Configuration) {
	fillOutControlsMap(configuration, defaultControlsMap)

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	width := ui.TermWidth()
	height := ui.TermHeight()

	ui.Render(drawPage[currentScreen](configuration, width, height)...)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			if e.ID == "<C-c>" {
				break
			} else {
				savedScreen := currentScreen
				handleKeyboard(configuration, e)
				// refresh screen after keyboard input or redraw screen entirely + reset state if we have changed screens
				if savedScreen == currentScreen {
					ui.Render(refreshPage[currentScreen](configuration, width, height)...)
				} else {
					// reset modes
					if currentScreen == Search {
						currentMode = Insert
					} else {
						currentMode = Normal
					}
					// reset text
					userTextBuffer = ""
					currentPodcastsInBuffer = nil
					currentSelected = 0
					// save last screen
					previousScreen = savedScreen
					// render new screen
					ui.Render(drawPage[currentScreen](configuration, width, height)...)
				}
			}
		} else if e.Type == ui.ResizeEvent {
			payload := e.Payload.(ui.Resize)
			width = payload.Width
			height = payload.Height
			ui.Render(drawPage[currentScreen](configuration, width, height)...)
		}
	}
}
