package main

import (
	ui "github.com/gizak/termui"
)

func prepareDrawPageHome(configuration *Configuration) {
	currentPodcastsInBuffers[currentScreen] = configuration.Subscribed
}

func prepareDrawPageSearch(configuration *Configuration) {
	searchFailed = false
}

func prepareDrawPageDownloaded(configuration *Configuration) {
	filteredList := make([]PodcastEpisode, 0, len(configuration.Downloaded))
	for _, v := range configuration.Downloaded {
		if podcastExistsOnDisk(v) {
			filteredList = append(filteredList, v)
		}
	}
	// If deleted outside of downloaded, we have to move the cursor down.
	if getCurrentCursorPosition() > len(filteredList)-1 {
		setCurrentCursorPosition(len(filteredList) - 1)
	} else if getCurrentCursorPosition() < 0 && len(filteredList) > 0 {
		setCurrentCursorPosition(0)
	}
	currentPodcastsInBuffers[currentScreen] = filteredList
}

func prepareDrawPagePodcastDetail(configuration *Configuration) {
	entries, err := getPodcastEntries(currentSelectedPodcast, &configuration.Cached)
	if err == nil {
		currentPodcastsInBuffers[currentScreen] = entries
	} else {
		currentPodcastsInBuffers[currentScreen] = make([]PodcastEpisode, 0)
	}
}

func drawPageHome(configuration *Configuration, width int, height int) []ui.Bufferer {
	fillOutControlsMap(configuration, defaultControlsMap)
	widgets := make([]ui.Bufferer, 3)
	widgets[0] = producePodcastListWidget(configuration, width, height)
	widgets[1] = produceControlsWidget(configuration, width, height)
	widgets[2] = producePlayerWidget(configuration, width, height)
	return widgets
}

func drawPageSearch(configuration *Configuration, width int, height int) []ui.Bufferer {
	fillOutControlsMap(configuration, defaultControlsMap)
	widgets := make([]ui.Bufferer, 4)
	widgets[0] = produceSearchWidget(configuration, width, height)
	widgets[1] = produceSearchResultsWidget(configuration, width, height)
	widgets[2] = produceControlsWidget(configuration, width, height)
	widgets[3] = producePlayerWidget(configuration, width, height)
	return widgets
}

func drawPageDownloaded(configuration *Configuration, width int, height int) []ui.Bufferer {
	fillOutControlsMap(configuration, defaultControlsMap)
	widgets := make([]ui.Bufferer, 3)
	widgets[0] = produceDownloadedWidget(configuration, width, height)
	widgets[1] = produceControlsWidget(configuration, width, height)
	widgets[2] = producePlayerWidget(configuration, width, height)
	return widgets
}

func drawPagePodcastDetail(configuration *Configuration, width int, height int) []ui.Bufferer {
	fillOutControlsMap(configuration, defaultControlsMap)
	widgets := make([]ui.Bufferer, 4)
	widgets[0] = producePodcastDetailDescriptionWidget(configuration, width, height)
	widgets[1] = producePodcastDetailListWidget(configuration, width, height)
	widgets[2] = produceControlsWidget(configuration, width, height)
	widgets[3] = producePlayerWidget(configuration, width, height)
	return widgets
}

func refreshPageSearch(configuration *Configuration, width int, height int) []ui.Bufferer {
	widgets := make([]ui.Bufferer, 1)
	widgets[0] = produceSearchWidget(configuration, width, height)
	return widgets
}
