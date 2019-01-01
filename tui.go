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
		Home:          upPressedHome,
		Search:        upPressedSearch,
		Downloaded:    upPressedDownloaded,
		PodcastDetail: upPressedPodcastDetail,
	}

	downPressed = map[screen]func(configuration *Configuration){
		Home:          downPressedHome,
		Search:        downPressedSearch,
		Downloaded:    downPressedDownloaded,
		PodcastDetail: downPressedPodcastDetail,
	}

	searchPressed = map[screen]func(configuration *Configuration){
		Home:       searchPressedHome,
		Search:     searchPressedSearch,
		Downloaded: searchPressedDownloaded,
	}

	currentPodcastsInBuffers = map[screen]interface{}{
		Home:          make([]Podcast, 0),
		Search:        make([]Podcast, 0),
		Downloaded:    make([]PodcastEpisode, 0),
		PodcastDetail: make([]Podcast, 0),
	}
)

var (
	currentSelected        = 0
	currentListSize        = 0
	currentListOffset      = 0
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

	prepareDrawPage[currentScreen](configuration)
	ui.Render(drawPage[currentScreen](configuration, width, height)...)

	for e := range ui.PollEvents() {
		savedScreen := currentScreen
		if e.Type == ui.KeyboardEvent {
			if e.ID == "<C-c>" {
				break
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
			// reset text and selected if not transitioning between detail screen
			if (currentScreen != PodcastDetail) && (savedScreen != PodcastDetail) {
				userTextBuffer = ""
				currentSelected = 0
			}
			// save last screen
			previousScreen = savedScreen
		}
		ui.Render(drawPage[currentScreen](configuration, width, height)...)
	}
}
