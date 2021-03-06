package main

import (
	"fmt"
	ui "github.com/gizak/termui"
)

type screen string
type mode string
type theme string

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

const (
	ThemeLight theme = "Light"
	ThemeDark  theme = "Dark"
)

var (
	currentListSize        = 0
	currentMode            = Normal
	currentScreen          = Home
	previousScreen         = None
	currentSelectedPodcast Podcast
	userTextBuffer         = ""
	searchFailed           = false
)

var (
	// TODO make spacing variable so when it's not wide enough it still works
	defaultControlsMap = map[screen]string{
		Home:          "[%s]elect/[<enter>]  [%s]/[<left>]search   [%s]/[<right>]downloaded   [%s]elete subscription    <Control-c> exit",
		Search:        "%s   [/]search   [esc]ape searching   [<enter>]%s   [j]down   [k]up   [l]/[<right>]home",
		Downloaded:    "[<enter>] Play   [%s]/[d]elete   [l]/[<left>]home",
		PodcastDetail: "[<enter>]/[%s] %s    [%s]elete downloaded    [%s]/[<left>]back",
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
		Home:          prepareDrawPageHome,
		Search:        prepareDrawPageSearch,
		Downloaded:    prepareDrawPageDownloaded,
		PodcastDetail: prepareDrawPagePodcastDetail,
	}

	drawPage = map[screen]func(configuration *Configuration, width int, height int) []ui.Bufferer{
		Home:          drawPageHome,
		Search:        drawPageSearch,
		Downloaded:    drawPageDownloaded,
		PodcastDetail: drawPagePodcastDetail,
	}

	refreshPage = map[screen]func(configuration *Configuration, width int, height int) []ui.Bufferer{
		Home:          nil,
		Search:        refreshPageSearch,
		Downloaded:    nil,
		PodcastDetail: nil,
	}

	actionPressed = map[screen]func(configuration *Configuration){
		Home:          enterPressedHome,
		Search:        actionPressedSearch,
		Downloaded:    deletePodcastSelectedByCursor,
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
		PodcastDetail: escapePressedPodcastDetail,
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
		Downloaded:    deletePodcastSelectedByCursor,
		PodcastDetail: deletePodcastSelectedByCursor,
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

func getCurrentPagePodcasts() []Podcast {
	return currentPodcastsInBuffers[currentScreen].([]Podcast)
}

func getCurrentPagePodcastEpisodes() []PodcastEpisode {
	return currentPodcastsInBuffers[currentScreen].([]PodcastEpisode)
}

func fillOutControlsMap(configuration *Configuration, controls map[screen]string) {
	switch currentScreen {
	case Home:
		controlsMap[Home] = fmt.Sprintf(controls[Home], configuration.ActionKeybind, configuration.LeftKeybind, configuration.RightKeybind, configuration.DeleteKeybind)
	case Search:
		var (
			searchText     string
			subscribedText string
		)
		if currentMode == Insert {
			searchText = "search         "
		} else {
			searchText = "look at podcast"
		}
		cursor := getCurrentCursorPosition()
		if len(currentPodcastsInBuffers[Search].([]Podcast)) > 0 &&
			podcastIsSubscribed(configuration, &currentPodcastsInBuffers[Search].([]Podcast)[cursor]) {
			subscribedText = fmt.Sprintf("un[%s]ubscribe", configuration.ActionKeybind)
		} else {
			subscribedText = fmt.Sprintf("[%s]ubscribe  ", configuration.ActionKeybind)
		}
		controlsMap[Search] = fmt.Sprintf(controls[Search], subscribedText, searchText)
	case Downloaded:
		controlsMap[Downloaded] = fmt.Sprintf(controls[Downloaded], configuration.ActionKeybind)
	case PodcastDetail:
		var actionText string
		cursor := getCurrentCursorPosition()
		if podcastIsDownloaded(configuration, &currentPodcastsInBuffers[PodcastDetail].([]PodcastEpisode)[cursor]) {
			actionText = "play episode    "
		} else {
			actionText = "download episode"
		}
		controlsMap[PodcastDetail] = fmt.Sprintf(controls[PodcastDetail], configuration.ActionKeybind, actionText, configuration.DeleteKeybind, configuration.LeftKeybind)
	}
	state := GetPlayerState()
	if state == PlayerPlay {
		playerText := "    <space> pause    <Pgup>ff %ds    <Pgdown>rw %ds"
		controlsMap[currentScreen] += fmt.Sprintf(playerText, configuration.FastForwardLength, configuration.RewindLength)
	} else if state == PlayerPause {
		controlsMap[currentScreen] += "    <space> resume"
	}
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
