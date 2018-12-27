package main

import (
	"fmt"
	ui "github.com/gizak/termui"
)

type Screen string

const (
	podmanHeader = "" +
		" _____         _                  \n" +
		"|  _  | ___  _| | _____  ___  ___ \n" +
		"|   __|| . || . ||     || .'||   |\n" +
		"|__|   |___||___||_|_|_||__,||_|_|"
	headerHeight = 5
	playerHeight = 3

	None       Screen = "None"
	Home       Screen = "Home"
	Search     Screen = "Search"
	Downloaded Screen = "Downloaded"
)

var (
	leftTransitions = map[Screen]Screen{
		Home:       Search,
		Search:     None,
		Downloaded: Home,
	}

	rightTransitions = map[Screen]Screen{
		Home:       Downloaded,
		Downloaded: None,
		Search:     Home,
	}

	createPage = map[Screen]func(configruation *Configuration, width int, height int) []ui.Bufferer{
		Home:       drawMainPage,
		Search:     drawSearch,
		Downloaded: drawDownloaded,
	}

	refreshPage = map[Screen]func(configruation *Configuration, width int, height int) []ui.Bufferer{
		Home:       refreshMainPage,
		Search:     refreshSearch,
		Downloaded: refreshDownlaoded,
	}
)

var (
	currentSelected = 0
	currentListSize = 0
	currentScreen   = Home
)

func termuiStyleText(text string, fgcolor string, bgcolor string) string {
	text = "[" + text + "](fg-" + fgcolor + ",bg-" + string(bgcolor) + ")"
	return text
}

func drawMainPage(configruation *Configuration, width int, height int) []ui.Bufferer {
	widgets := make([]ui.Bufferer, 0)
	widgets = append(widgets, produceHeaderWidget(width))
	widgets = append(widgets, producePodcastListWidget(configruation, width, height))
	widgets = append(widgets, producePlayerWidget(configruation, width, height))
	return widgets
}

func refreshMainPage(configruation *Configuration, width int, height int) []ui.Bufferer {
	widgets := make([]ui.Bufferer, 0)
	widgets = append(widgets, producePodcastListWidget(configruation, width, height))
	widgets = append(widgets, producePlayerWidget(configruation, width, height))
	return widgets
}

func drawSearch(configruation *Configuration, width int, height int) []ui.Bufferer {
	widgets := make([]ui.Bufferer, 0)
	return widgets
}

func refreshSearch(configruation *Configuration, width int, height int) []ui.Bufferer {
	widgets := make([]ui.Bufferer, 0)
	return widgets
}

func drawDownloaded(configruation *Configuration, width int, height int) []ui.Bufferer {
	widgets := make([]ui.Bufferer, 0)
	return widgets
}

func refreshDownlaoded(configruation *Configuration, width int, height int) []ui.Bufferer {
	widgets := make([]ui.Bufferer, 0)
	return widgets
}

func refreshPlayer(configuration *Configuration, width int, height int) []ui.Bufferer {
	widgets := make([]ui.Bufferer, 0)
	widgets = append(widgets, producePlayerWidget(configuration, width, height))
	return widgets
}

func transitionScreen(transitions map[Screen]Screen, screen Screen) {
	if transitions[screen] == None {
		return
	}
	currentScreen = transitions[screen]
}

func handleKeyboard(configuration *Configuration, event ui.Event) {
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
		if currentSelected > 0 {
			currentSelected--
		}
	}
	if event.ID == configuration.DownKeybind || event.ID == "<Down>" {
		if currentSelected < currentListSize-1 {
			currentSelected++
		}
	}
}

func StartTui(configuration *Configuration) {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	width := ui.TermWidth()
	height := ui.TermHeight()

	ui.Render(createPage[currentScreen](configuration, width, height)...)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			if e.ID == "<C-c>" {
				break
			} else {
				handleKeyboard(configuration, e)
				ui.Render(refreshPage[currentScreen](configuration, width, height)...)
			}
		}
	}
}
