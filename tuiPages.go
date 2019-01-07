package main

import (
	ui "github.com/gizak/termui"
)

func prepareDrawPageMain(configuration *Configuration) {
	currentPodcastsInBuffers[currentScreen] = configuration.Subscribed
}

func prepareDrawPageSearch(configuration *Configuration) {
}

func prepareDrawPageDownloaded(configuration *Configuration) {
	currentPodcastsInBuffers[currentScreen] = configuration.Downloaded
}

func prepareDrawPagePodcastDetail(configuration *Configuration) {
	entries, err := getPodcastEntries(currentSelectedPodcast, &configuration.Cached)
	if err == nil {
		currentPodcastsInBuffers[currentScreen] = entries
	} else {
		currentPodcastsInBuffers[currentScreen] = make([]PodcastEpisode, 0)
	}
}

func drawPageMain(configuration *Configuration, width int, height int) []ui.Bufferer {
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
	widgets := make([]ui.Bufferer, 3)
	widgets[0] = produceDownloadedWidget(configuration, width, height)
	widgets[1] = produceControlsWidget(configuration, width, height)
	widgets[2] = producePlayerWidget(configuration, width, height)
	return widgets
}

func drawPagePodcastDetail(configuration *Configuration, width int, height int) []ui.Bufferer {
	widgets := make([]ui.Bufferer, 2)
	widgets[0] = producePodcastDetailDescriptionWidget(configuration, width, height)
	widgets[1] = producePodcastDetailListWidget(configuration, width, height)
	return widgets
}
