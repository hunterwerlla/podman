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
	// TODO make spacing variable so when it's not wide enough it still works
	defaultControlsMap = map[screen]string{
		Home:          "[%s]elect/[<enter>]  [h]/[<left>](search)   [l]/[<right>](downloaded)   [d]elete subscription",
		Search:        "[s]ubscribe/unsubscribe   [/]search   [esc]ape searching   [<enter>]%s   [j]down   [k]up   [l]/[<right>](home)",
		Downloaded:    "[<enter>] Play   [%s]/[d]elete   [l]/[<left>](home)",
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

	prepareDrawPage = map[screen]func(configuration *Configuration){
		Home:          prepareDrawPageMain,
		Search:        prepareDrawPageSearch,
		Downloaded:    prepareDrawPageDownloaded,
		PodcastDetail: prepareDrawPagePodcastDetail,
	}

	drawPage = map[screen]func(configuration *Configuration, width int, height int) []ui.Bufferer{
		Home:          drawPageMain,
		Search:        drawPageSearch,
		Downloaded:    drawPageDownloaded,
		PodcastDetail: drawPagePodcastDetail,
	}

	actionPressed = map[screen]func(configuration *Configuration){
		Home:          enterPressedHome,
		Search:        actionPressedSearch,
		Downloaded:    actionPressedDownloaded,
		PodcastDetail: doNothingWithInput,
	}

	enterPressed = map[screen]func(configuration *Configuration){
		Home:          enterPressedHome,
		Search:        enterPressedSearch,
		Downloaded:    enterPressedDownloaded,
		PodcastDetail: enterPressedPodcastDetail,
	}

	escapePressed = map[screen]func(configuration *Configuration){
		Home:          doNothingWithInput,
		Search:        doNothingWithInput,
		Downloaded:    doNothingWithInput,
		PodcastDetail: doNothingWithInput,
	}

	upPressed = map[screen]func(configuration *Configuration){
		Home:          upPressedGeneric,
		Search:        upPressedSearch,
		Downloaded:    upPressedGeneric,
		PodcastDetail: upPressedGeneric,
	}

	downPressed = map[screen]func(configuration *Configuration){
		Home:          downPressedGeneric,
		Search:        downPressedSearch,
		Downloaded:    downPressedGeneric,
		PodcastDetail: downPressedGeneric,
	}

	searchPressed = map[screen]func(configuration *Configuration){
		Home:          searchPressedHome,
		Search:        searchPressedSearch,
		Downloaded:    searchPressedDownloaded,
		PodcastDetail: doNothingWithInput,
	}

	deletePressed = map[screen]func(configuration *Configuration){
		Home:          deletePressedHome,
		Search:        doNothingWithInput,
		Downloaded:    actionPressedDownloaded,
		PodcastDetail: doNothingWithInput,
	}

	currentPodcastsInBuffers = map[screen]interface{}{
		Home:          make([]Podcast, 0),
		Search:        make([]Podcast, 0),
		Downloaded:    make([]PodcastEpisode, 0),
		PodcastDetail: make([]Podcast, 0),
	}

	currentssSelectedOnScreen = map[screen]int{
		Home:          0,
		Search:        0,
		Downloaded:    0,
		PodcastDetail: 0,
	}
)

var (
	currentListSize        = 0
	currentMode            = Normal
	currentScreen          = Home
	previousScreen         = None
	currentSelectedPodcast Podcast
	userTextBuffer         = ""
)

func getCurrentPagePodcasts() []Podcast {
	return currentPodcastsInBuffers[currentScreen].([]Podcast)
}

func getCurrentPagePodcastEpisodes() []PodcastEpisode {
	return currentPodcastsInBuffers[currentScreen].([]PodcastEpisode)
}

// TODO break into per screen functions
func fillOutControlsMap(configuration *Configuration, controls map[screen]string) {
	var searchText string
	controlsMap[Home] = fmt.Sprintf(controls[Home], configuration.ActionKeybind)
	if currentMode == Insert {
		searchText = "finish search  "
	} else {
		searchText = "look at podcast"
	}
	controlsMap[Search] = fmt.Sprintf(controls[Search], searchText)
	controlsMap[Downloaded] = fmt.Sprintf(controls[Downloaded], configuration.ActionKeybind)
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

func getCurrentCursorPosition() int {
	return currentssSelectedOnScreen[currentScreen]
}

func setCurrentCursorPosition(position int) {
	currentssSelectedOnScreen[currentScreen] = position
}

// StartTui starts the TUI with the Configuration passed in
func StartTui(configuration *Configuration) {
	fillOutControlsMap(configuration, defaultControlsMap)

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	tuiMainLoop(configuration)
}
