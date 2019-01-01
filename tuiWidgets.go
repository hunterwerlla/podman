package main

import (
	ui "github.com/gizak/termui"
)

const (
	playerHeight    = 3
	searchBarHeight = 3
	controlsHeight  = 2
)

func producePodcastListWidget(configruation *Configuration, width int, height int) *ui.List {
	podcastWidget := ui.NewList()
	podcastWidget.Width = width
	podcastWidget.Height = height - playerHeight - controlsHeight
	podcastWidget.Y = 0
	podcastWidget.Border = false
	podcastWidget.ItemFgColor = ui.ColorBlack
	var podcasts []string
	currentListSize = len(currentPodcastsInBuffer)
	for ii, item := range currentPodcastsInBuffer {
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
	} else if len(userTextBuffer) == 0 {
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
	searchWidget.Y = 0
	searchWidget.TextFgColor = ui.ColorBlack
	searchWidget.Height = searchBarHeight
	searchWidget.Width = width
	searchWidget.Border = false
	return searchWidget
}

func produceSearchResultsWidget(configuration *Configuration, width int, height int) *ui.List {
	searchResultsWidget := ui.NewList()
	searchResultsWidget.Y = searchBarHeight
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
		for _, value := range configuration.Subscribed {
			if item.ArtistName == value.ArtistName && item.CollectionName == value.CollectionName {
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
	searchResultsWidget.Y = 0
	searchResultsWidget.Height = height - playerHeight - controlsHeight
	searchResultsWidget.Width = width
	searchResultsWidget.Border = false
	searchResultsWidget.ItemFgColor = ui.ColorBlack
	var podcasts []string
	var podcastList = currentPodcastsInBuffers[Downloaded].([]PodcastEpisode)
	currentListSize = len(podcastList)
	currentNum := 0
	for _, item := range podcastList {
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

func producePodcastDetailWidget(configuration *Configuration, width int, height int) ui.Bufferer {
	podcastDetailWidget := ui.NewParagraph("Stuff \n")
	podcastDetailWidget.TextFgColor = ui.ColorBlack
	podcastDetailWidget.Width = width
	podcastDetailWidget.Height = height
	podcastDetailWidget.Y = 0
	podcastDetailWidget.BorderLeft = false
	podcastDetailWidget.BorderRight = false
	podcastDetailWidget.BorderBottom = false
	return podcastDetailWidget
}
