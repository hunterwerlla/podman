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
		"|__|   |___||___||_|_|_||__,||_|_|\n"
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
)

var (
	currentSelected = 0
	currentListSize = 0
	currentScreen   = Home
)

func produceHeaderWidget(width int) *ui.Paragraph {
	headerWidget := ui.NewParagraph(podmanHeader)
	headerWidget.Height = headerHeight
	headerWidget.Width = width
	headerWidget.TextFgColor = ui.ColorBlack
	headerWidget.BorderTop = false
	headerWidget.BorderLeft = false
	headerWidget.BorderRight = false
	headerWidget.BorderBottom = true
	headerWidget.WrapLength = 0
	return headerWidget
}

func producePodcastListWidget(configruation *Configuration, width int, height int) *ui.List {
	podcastWidget := ui.NewList()
	podcastWidget.Width = width
	podcastWidget.Height = height - headerHeight - playerHeight
	podcastWidget.Y = headerHeight
	podcastWidget.Border = false
	var podcasts []string
	currentListSize = len(configruation.Subscribed)
	for ii, item := range configruation.Subscribed {
		formattedPodcast := formatPodcast(item, width)
		if ii == currentSelected {
			formattedPodcast = "[" + formattedPodcast + "](fg-white,bg-green)"
		}
		podcasts = append(podcasts, formattedPodcast)
	}
	podcastWidget.Items = podcasts
	return podcastWidget
}

func producePlayerWidget(width int, height int) *ui.Paragraph {
	headerWidget := ui.NewParagraph("play something")
	return headerWidget
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
	ui.Render(produceHeaderWidget(width))
	ui.Render(producePodcastListWidget(configuration, width, height))

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			if e.ID == "<C-c>" {
				break
			} else {
				handleKeyboard(configuration, e)
				ui.Render(producePodcastListWidget(configuration, width, height))
			}
		}
	}
}
