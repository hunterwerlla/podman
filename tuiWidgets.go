package main

import (
	"fmt"
	ui "github.com/gizak/termui"
)

const (
	playerHeight                   = 2
	searchBarHeight                = 3
	controlsHeight                 = 2
	podcastDetailDescriptionHeight = 3
)

func producePodcastListWidget(configruation *Configuration, width int, height int) *ui.List {
	podcastWidget := ui.NewList()
	podcastWidget.Width = width
	podcastWidget.Height = height - playerHeight - controlsHeight
	podcastWidget.Y = 0
	podcastWidget.Border = false
	podcastWidget.ItemFgColor = ui.ColorBlack
	var listFormattedPodcasts []string
	podcasts := getCurrentPagePodcasts()
	cursor := getCurrentCursorPosition()
	currentListSize = len(podcasts)
	for currentNum, item := range podcasts {
		formattedPodcast := formatPodcast(item, width)
		if currentNum == cursor {
			formattedPodcast = termuiStyleText(formattedPodcast, "white", "black")
		}
		listFormattedPodcasts = append(listFormattedPodcasts, formattedPodcast)
	}
	podcastWidget.Items = listFormattedPodcasts
	return podcastWidget
}

// TODO figure out how to fix this mess
func producePlayerWidget(configuration *Configuration, width int, height int) ui.Bufferer {
	var widgetLabel string
	if downloadInProgress() {
		widgetLabel = "Downloading"
	} else {
		widgetLabel = "Nothing playing"
	}
	if GetPlayerState() != PlayerPlay || GetLengthOfPlayingFile() < 1 {
		playerWidget := ui.NewParagraph(widgetLabel)
		playerWidget.TextFgColor = ui.ColorBlack
		playerWidget.Width = width
		playerWidget.Height = playerHeight
		playerWidget.Y = height - playerHeight
		playerWidget.BorderLeft = false
		playerWidget.BorderRight = false
		playerWidget.BorderBottom = false
		return playerWidget
	}
	lengthOfPlayingFile := GetLengthOfPlayingFile()
	currentPlayingPosition := GetPlayerPosition()
	label := fmt.Sprintf("%d/%d", int(currentPlayingPosition), lengthOfPlayingFile)
	playerWidget := ui.NewGauge()
	playerWidget.Percent = int((float64(currentPlayingPosition) / float64(lengthOfPlayingFile)) * 100)
	playerWidget.Label = label
	playerWidget.Width = width
	playerWidget.Height = playerHeight
	playerWidget.Y = height - playerHeight
	playerWidget.BorderLeft = false
	playerWidget.BorderRight = false
	playerWidget.BorderBottom = false
	playerWidget.BarColor = ui.ColorBlack
	playerWidget.LabelAlign = ui.AlignLeft
	return playerWidget
}

func produceSearchWidget(configuration *Configuration, width int, height int) *ui.Paragraph {
	text := ""
	podcasts := getCurrentPagePodcasts()
	if len(podcasts) > 0 {
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
	if len(userTextBuffer) > 0 {
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
	searchWidgetHeight := height - searchBarHeight - playerHeight - controlsHeight
	searchResultsWidget := ui.NewList()
	searchResultsWidget.Y = searchBarHeight
	searchResultsWidget.Height = searchWidgetHeight
	searchResultsWidget.Width = width
	searchResultsWidget.Border = false
	searchResultsWidget.ItemFgColor = ui.ColorBlack
	var formattedPodcastList []string
	podcasts := getCurrentPagePodcasts()
	currentListSize = len(podcasts)
	cursor := getCurrentCursorPosition()
	for currentNum, item := range podcasts {
		if currentNum < (cursor - (searchWidgetHeight / 2)) {
			continue
		}
		formattedPodcast := formatPodcast(item, width)
		// TODO make an isSubscribed function
		subscribed := false
		for _, value := range configuration.Subscribed {
			if item.ArtistName == value.ArtistName && item.CollectionName == value.CollectionName {
				//already subscribed, add S -
				formattedPodcast = "S - " + formattedPodcast
				subscribed = true
			}
		}
		if subscribed != true {
			formattedPodcast = "    " + formattedPodcast
		}
		if currentNum == cursor {
			formattedPodcast = termuiStyleText(formattedPodcast, "white", "black")
		}
		formattedPodcastList = append(formattedPodcastList, formattedPodcast)
	}
	searchResultsWidget.Items = formattedPodcastList
	return searchResultsWidget
}

func produceDownloadedWidget(configuration *Configuration, width int, height int) *ui.List {
	searchResultsWidgetHeight := height - playerHeight - controlsHeight
	searchResultsWidget := ui.NewList()
	searchResultsWidget.Y = 0
	searchResultsWidget.Height = searchResultsWidgetHeight
	searchResultsWidget.Width = width
	searchResultsWidget.Border = false
	searchResultsWidget.ItemFgColor = ui.ColorBlack
	searchResultsWidget.Overflow = "wrap"
	var listFormattedPodcasts []string
	var podcast = getCurrentPagePodcastEpisodes()
	currentListSize = len(podcast)
	cursor := getCurrentCursorPosition()
	for currentNum, item := range podcast {
		if currentNum < (cursor - (searchResultsWidgetHeight / 2)) {
			continue
		}
		// TODO add function for this
		formattedPodcast := item.PodcastTitle + " " + item.Title + " " + item.Summary
		if currentNum == cursor {
			formattedPodcast = wrapString(formattedPodcast, width)
			formattedPodcast = termuiStyleText(formattedPodcast, "white", "black")
		} else if len(formattedPodcast) > width {
			formattedPodcast = formattedPodcast[0 : width-3]
			formattedPodcast += "..."
		}
		listFormattedPodcasts = append(listFormattedPodcasts, formattedPodcast)
	}
	searchResultsWidget.Items = listFormattedPodcasts
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

func producePodcastDetailDescriptionWidget(configuration *Configuration, width int, height int) ui.Bufferer {
	blurb := currentSelectedPodcast.CollectionName + ", " + currentSelectedPodcast.ArtistName + "\n" + currentSelectedPodcast.Description
	podcastDetailWidget := ui.NewParagraph(blurb)
	podcastDetailWidget.TextFgColor = ui.ColorBlack
	podcastDetailWidget.Width = width
	podcastDetailWidget.Height = podcastDetailDescriptionHeight
	podcastDetailWidget.Y = 0
	podcastDetailWidget.BorderLeft = false
	podcastDetailWidget.BorderRight = false
	podcastDetailWidget.BorderTop = false
	return podcastDetailWidget
}

func producePodcastDetailListWidget(configuration *Configuration, width int, height int) ui.Bufferer {
	podcastDetailListWidgetHeight := height - podcastDetailDescriptionHeight
	var listFormattedPodcasts []string
	podcasts := currentPodcastsInBuffers[currentScreen].([]PodcastEpisode)
	podcastDetailListWidget := ui.NewList()
	podcastDetailListWidget.Width = width
	podcastDetailListWidget.Height = podcastDetailListWidgetHeight
	podcastDetailListWidget.Y = podcastDetailDescriptionHeight
	podcastDetailListWidget.Border = false
	podcastDetailListWidget.ItemFgColor = ui.ColorBlack
	podcastDetailListWidget.Overflow = "wrap"
	currentListSize = len(podcasts)
	cursor := getCurrentCursorPosition()
	for currentNum, item := range podcasts {
		if currentNum < (cursor - (podcastDetailListWidgetHeight / 2)) {
			continue
		}
		formattedPodcast := formatPodcastEpisode(item)
		if podcastIsDownloaded(configuration, item) {
			formattedPodcast = "D - " + formattedPodcast
		} else {
			formattedPodcast = "    " + formattedPodcast
		}
		if currentNum == getCurrentCursorPosition() {
			formattedPodcast = wrapString(formattedPodcast, width)
			formattedPodcast = termuiStyleText(formattedPodcast, "white", "black")
		} else if len(formattedPodcast) > width {
			formattedPodcast = formattedPodcast[0 : width-3]
			formattedPodcast += "..."
		}
		listFormattedPodcasts = append(listFormattedPodcasts, formattedPodcast)
		currentNum++
	}
	podcastDetailListWidget.Items = listFormattedPodcasts
	return podcastDetailListWidget
}
