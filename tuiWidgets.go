package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	"time"
)

const (
	playerHeight                   = 2
	searchBarHeight                = 3
	controlsHeight                 = 2
	podcastDetailDescriptionHeight = 3
)

func everyTwoSeconds() bool {
	t := time.Now().Second() % 4
	if t == 0 || t == 1 {
		return true
	}
	return false
}

func everyOtherSecond() bool {
	if time.Now().Second()%2 == 0 {
		return true
	}
	return false
}

func producePodcastListWidget(configuration *Configuration, width int, height int) *ui.List {
	podcastWidget := ui.NewList()
	podcastWidget.Width = width
	podcastWidget.Height = height - playerHeight - controlsHeight
	podcastWidget.Y = 0
	podcastWidget.Border = false
	podcastWidget.ItemFgColor = getForegroundColorForTheme(configuration)
	var listFormattedPodcasts []string
	podcasts := getCurrentPagePodcasts()
	cursor := getCurrentCursorPosition()
	currentListSize = len(podcasts)
	if currentListSize == 0 {
		tutorialString := fmt.Sprintf("No subscriptions, go left (<left>/%s) to search for podcasts!", configuration.LeftKeybind)
		listFormattedPodcasts = append(listFormattedPodcasts, tutorialString)
	}
	for currentNum, item := range podcasts {
		formattedPodcast := formatPodcast(item, width)
		if currentNum == cursor {
			formattedPodcast = termuiStyleText(formattedPodcast, getBackgroundColorForThemeString(configuration), getForegroundColorForThemeString(configuration))
		}
		listFormattedPodcasts = append(listFormattedPodcasts, formattedPodcast)
	}
	podcastWidget.Items = listFormattedPodcasts
	return podcastWidget
}

func produceNothingPlayingWidget(configuration *Configuration, width int, height int) ui.Bufferer {
	var widgetText string
	if GetPlayerState() == PlayerPause {
		widgetText = "Paused"
	} else {
		widgetText = "Nothing playing"
	}
	playerWidget := ui.NewParagraph(widgetText)
	playerWidget.TextFgColor = getForegroundColorForTheme(configuration)
	playerWidget.Width = width
	playerWidget.Height = playerHeight
	playerWidget.Y = height - playerHeight
	playerWidget.BorderLeft = false
	playerWidget.BorderRight = false
	playerWidget.BorderBottom = false
	return playerWidget
}

func fillPlayerGauge(playerWidget *ui.Gauge, configuration *Configuration, width int, height int) {
	var label string
	if (downloadInProgress() && GetPlayerState() != PlayerPlay) || (downloadInProgress() && everyTwoSeconds()) {
		label = "Downloading: "
		var (
			totalDownloadSize       int64
			totalDownloadCompleated int64
		)
		num := 0
		for key, value := range downloading {
			if num > 0 {
				label += " & "
			}
			label += key + " [" + byteCountDecimal(value.TotalDownloaded) + "/" + byteCountDecimal(value.FileSize) + "] (" + byteCountDecimal(value.Speed) + "/s)"

			totalDownloadSize += value.FileSize
			totalDownloadCompleated += value.TotalDownloaded
			num++
		}
		playerWidget.Percent = int((float64(totalDownloadCompleated) / float64(totalDownloadSize+1)) * 100)
	} else {
		lengthOfPlayingFile := GetLengthOfPlayingFile()
		currentPlayingPosition := GetPlayerPosition()
		label = fmt.Sprintf("%d/%d", int(currentPlayingPosition), lengthOfPlayingFile)
		playerWidget.Percent = int((float64(currentPlayingPosition) / float64(lengthOfPlayingFile)) * 100)
	}
	playerWidget.Label = label
}

func producePlayerWidget(configuration *Configuration, width int, height int) ui.Bufferer {
	// when nothing is happening, just display a generic message
	if !downloadInProgress() && GetPlayerState() != PlayerPlay {
		return produceNothingPlayingWidget(configuration, width, height)
	}
	playerWidget := ui.NewGauge()
	fillPlayerGauge(playerWidget, configuration, width, height)
	playerWidget.Width = width
	playerWidget.Height = playerHeight
	playerWidget.Y = height - playerHeight
	playerWidget.BorderLeft = false
	playerWidget.BorderRight = false
	playerWidget.BorderBottom = false
	playerWidget.BarColor = getForegroundColorForTheme(configuration)
	playerWidget.PercentColor = getForegroundColorForTheme(configuration)
	playerWidget.PercentColorHighlighted = getBackgroundColorForTheme(configuration)
	playerWidget.LabelAlign = ui.AlignLeft | ui.AlignCenterVertical
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
	if currentMode == Insert && everyOtherSecond() {
		text += "_"
	}
	searchWidget := ui.NewParagraph(text)
	searchWidget.Y = 0
	searchWidget.TextFgColor = getForegroundColorForTheme(configuration)
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
	searchResultsWidget.ItemFgColor = getForegroundColorForTheme(configuration)
	var listFormattedPodcasts []string
	podcasts := getCurrentPagePodcasts()
	currentListSize = len(podcasts)
	cursor := -1
	if currentListSize == 0 && searchFailed {
		listFormattedPodcasts = append(listFormattedPodcasts, "No results found")
	}
	if currentListSize != 0 && currentMode == Normal {
		// Only highlight when we are not searching and the list size is not zero
		cursor = getCurrentCursorPosition()
	}
	for currentNum, item := range podcasts {
		if currentNum < (cursor - (searchWidgetHeight / 2)) {
			continue
		}
		formattedPodcast := formatPodcast(item, width)
		if podcastIsSubscribed(configuration, &item) {
			formattedPodcast = "S - " + formattedPodcast
		} else {
			formattedPodcast = "    " + formattedPodcast
		}
		if currentNum == cursor {
			formattedPodcast = termuiStyleText(formattedPodcast, getBackgroundColorForThemeString(configuration), getForegroundColorForThemeString(configuration))
		}
		listFormattedPodcasts = append(listFormattedPodcasts, formattedPodcast)
	}
	searchResultsWidget.Items = listFormattedPodcasts
	return searchResultsWidget
}

func produceDownloadedWidget(configuration *Configuration, width int, height int) *ui.List {
	searchResultsWidgetHeight := height - playerHeight - controlsHeight
	searchResultsWidget := ui.NewList()
	searchResultsWidget.Y = 0
	searchResultsWidget.Height = searchResultsWidgetHeight
	searchResultsWidget.Width = width
	searchResultsWidget.Border = false
	searchResultsWidget.ItemFgColor = getForegroundColorForTheme(configuration)
	searchResultsWidget.Overflow = "wrap"
	var listFormattedPodcasts []string
	var podcast = getCurrentPagePodcastEpisodes()
	currentListSize = len(podcast)
	cursor := -1
	if currentListSize == 0 {
		listFormattedPodcasts = append(listFormattedPodcasts, "No podcasts downloaded yet")
	} else {
		cursor = getCurrentCursorPosition()
	}
	for currentNum, item := range podcast {
		if currentNum < (cursor - (searchResultsWidgetHeight / 2)) {
			continue
		}
		formattedPodcast := item.PodcastTitle + " " + item.Title + " " + item.Summary
		formattedPodcast = wrapOrBreakText(configuration, formattedPodcast, width, currentNum == getCurrentCursorPosition())
		listFormattedPodcasts = append(listFormattedPodcasts, formattedPodcast)
	}
	searchResultsWidget.Items = listFormattedPodcasts
	return searchResultsWidget
}

func produceControlsWidget(configuration *Configuration, width int, height int) ui.Bufferer {
	controlsWidget := ui.NewParagraph(controlsMap[currentScreen])
	controlsWidget.TextFgColor = getForegroundColorForTheme(configuration)
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
	podcastDetailWidget.TextFgColor = getForegroundColorForTheme(configuration)
	podcastDetailWidget.Width = width
	podcastDetailWidget.Height = podcastDetailDescriptionHeight
	podcastDetailWidget.Y = 0
	podcastDetailWidget.BorderLeft = false
	podcastDetailWidget.BorderRight = false
	podcastDetailWidget.BorderTop = false
	return podcastDetailWidget
}

func producePodcastDetailListWidget(configuration *Configuration, width int, height int) ui.Bufferer {
	podcastDetailListWidgetHeight := height - podcastDetailDescriptionHeight - playerHeight - controlsHeight
	var listFormattedPodcasts []string
	podcasts := currentPodcastsInBuffers[currentScreen].([]PodcastEpisode)
	podcastDetailListWidget := ui.NewList()
	podcastDetailListWidget.Width = width
	podcastDetailListWidget.Height = podcastDetailListWidgetHeight
	podcastDetailListWidget.Y = podcastDetailDescriptionHeight
	podcastDetailListWidget.Border = false
	podcastDetailListWidget.ItemFgColor = getForegroundColorForTheme(configuration)
	podcastDetailListWidget.Overflow = "wrap"
	currentListSize = len(podcasts)
	cursor := -1
	if currentListSize == 0 {
		listFormattedPodcasts = append(listFormattedPodcasts, "Could not load podcast")
	} else {
		cursor = getCurrentCursorPosition()
	}
	for currentNum, item := range podcasts {
		if currentNum < (cursor - (podcastDetailListWidgetHeight / 2)) {
			continue
		}
		formattedPodcast := formatPodcastEpisode(item)
		if podcastIsDownloaded(configuration, &item) {
			formattedPodcast = "D - " + formattedPodcast
		} else {
			formattedPodcast = "    " + formattedPodcast
		}
		formattedPodcast = wrapOrBreakText(configuration, formattedPodcast, width, currentNum == getCurrentCursorPosition())
		listFormattedPodcasts = append(listFormattedPodcasts, formattedPodcast)
	}
	podcastDetailListWidget.Items = listFormattedPodcasts
	return podcastDetailListWidget
}
