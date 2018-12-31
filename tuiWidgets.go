package main

import (
	ui "github.com/gizak/termui"
)

const (
	podmanHeader = "" +
		" _____         _                  \n" +
		"|  _  | ___  _| | _____  ___  ___ \n" +
		"|   __|| . || . ||     || .'||   |\n" +
		"|__|   |___||___||_|_|_||__,||_|_|"
	downloadedHeader = "" +
		"____                     _              _         _\n" +
		"|    \\  ___  _ _ _  ___ | | ___  ___  _| | ___  _| |\n" +
		"|  |  || . || | | ||   || || . || .'|| . || -_|| . |\n" +
		"|____/ |___||_____||_|_||_||___||__,||___||___||___|"
	searchHeader = "" +
		" _____                      _\n" +
		"|   __| ___  ___  ___  ___ | |_\n" +
		"|__   || -_|| .'||  _||  _||   |\n" +
		"|_____||___||__,||_|  |___||_|_|"
	headerHeight    = 5
	playerHeight    = 3
	searchBarHeight = 3
	controlsHeight  = 2
)

func produceHeaderWidget(width int, header string) *ui.Paragraph {
	headerWidget := ui.NewParagraph(header)
	headerWidget.Height = headerHeight
	headerWidget.Width = width
	headerWidget.TextFgColor = ui.ColorBlack
	headerWidget.TextBgColor = ui.ColorRGB(4, 4, 4)
	headerWidget.BorderTop = false
	headerWidget.BorderLeft = false
	headerWidget.BorderRight = false
	headerWidget.BorderBottom = true
	return headerWidget
}

func producePodcastListWidget(configruation *Configuration, width int, height int) *ui.List {
	podcastWidget := ui.NewList()
	podcastWidget.Width = width
	podcastWidget.Height = height - headerHeight - playerHeight
	podcastWidget.Y = headerHeight
	podcastWidget.Border = false
	podcastWidget.ItemFgColor = ui.ColorBlack
	var podcasts []string
	currentListSize = len(configruation.Subscribed)
	for ii, item := range configruation.Subscribed {
		formattedPodcast := formatPodcast(item, width)
		if ii == currentSelected {
			formattedPodcast = termuiStyleText(formattedPodcast, "white", "black")
		}
		podcasts = append(podcasts, formattedPodcast)
	}
	podcastWidget.Items = podcasts
	return podcastWidget
}

// TODO figure out how to fix this mess
func producePlayerWidget(configuration *Configuration, width int, height int) ui.Bufferer {
	if GetPlayerState() != Play {
		playerWidget := ui.NewParagraph("Nothing playing")
		playerWidget.TextFgColor = ui.ColorBlack
		playerWidget.Width = width
		playerWidget.Height = playerHeight
		playerWidget.Y = height - playerHeight
		playerWidget.BorderLeft = false
		playerWidget.BorderRight = false
		playerWidget.BorderBottom = false
		return playerWidget
	}
	playerWidget := ui.NewGauge()
	playerWidget.Label = "whatever {{percent}}%"
	playerWidget.Width = width
	playerWidget.Height = playerHeight
	playerWidget.Y = height - playerHeight
	playerWidget.BorderLeft = false
	playerWidget.BorderRight = false
	playerWidget.BorderBottom = false
	return playerWidget
}

func produceSearchWidget(configuration *Configuration, width int, height int) *ui.Paragraph {
	text := ""
	if len(currentPodcastsInBuffer) > 0 {
		text = "    Results:\n"
	} else {
		text = "    Search for podcasts:\n"
	}
	if currentMode == Insert {
		text += "   >"
	} else {
		text += "    "
	}
	if userTextBuffer != "" {
		text += userTextBuffer
	}
	if currentMode == Insert {
		text += "_"
	}
	searchWidget := ui.NewParagraph(text)
	searchWidget.Y = headerHeight
	searchWidget.TextFgColor = ui.ColorBlack
	searchWidget.Height = searchBarHeight
	searchWidget.Width = width
	searchWidget.Border = false
	return searchWidget
}

func produceSearchResultsWidget(configuration *Configuration, width int, height int) *ui.List {
	searchResultsWidget := ui.NewList()
	searchResultsWidget.Y = searchBarHeight + headerHeight
	searchResultsWidget.Height = height - searchBarHeight - playerHeight - controlsHeight
	searchResultsWidget.Width = width
	searchResultsWidget.Border = false
	searchResultsWidget.ItemFgColor = ui.ColorBlack
	var podcasts []string
	currentListSize = len(currentPodcastsInBuffer)
	for ii, item := range currentPodcastsInBuffer {
		formattedPodcast := formatPodcast(item, width)
		// TODO make an isSubscribed function
		subscribed := false
		for _, thing := range configuration.Subscribed {
			if item.ArtistName == thing.ArtistName && item.CollectionName == thing.CollectionName {
				//already subscribed so do nothing
				formattedPodcast = "S - " + formattedPodcast
				subscribed = true
			}
		}
		if subscribed != true {
			formattedPodcast = "    " + formattedPodcast
		}
		if ii == currentSelected {
			formattedPodcast = termuiStyleText(formattedPodcast, "white", "black")
		}
		podcasts = append(podcasts, formattedPodcast)
	}
	searchResultsWidget.Items = podcasts
	return searchResultsWidget
}

func produceDownloadedWidget(configuration *Configuration, width int, height int) *ui.List {
	searchResultsWidget := ui.NewList()
	searchResultsWidget.Y = headerHeight
	searchResultsWidget.Height = height - playerHeight - headerHeight
	searchResultsWidget.Width = width
	searchResultsWidget.Border = false
	searchResultsWidget.ItemFgColor = ui.ColorBlack
	var podcasts []string
	currentListSize = len(configuration.Downloaded)
	currentNum := 0
	for _, item := range configuration.Downloaded {
		// TODO add function for this
		formattedPodcast := item.PodcastTitle + " " + item.Title + " " + item.Summary
		if currentNum == currentSelected {
			formattedPodcast = termuiStyleText(formattedPodcast, "white", "black")
		}
		podcasts = append(podcasts, formattedPodcast)
		currentNum++
	}
	searchResultsWidget.Items = podcasts
	return searchResultsWidget
}

func produceControlsWidget(configuration *Configuration, width int, height int) ui.Bufferer {
	controlsWidget := ui.NewParagraph(controlsMap[currentScreen])
	controlsWidget.TextFgColor = ui.ColorBlack
	controlsWidget.Width = width
	controlsWidget.Height = controlsHeight
	controlsWidget.Y = height - playerHeight - controlsHeight
	controlsWidget.BorderLeft = false
	controlsWidget.BorderRight = false
	controlsWidget.BorderBottom = false
	return controlsWidget
}
