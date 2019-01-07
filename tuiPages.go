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
	widgets := make([]ui.Bufferer, 0)
	widgets = append(widgets, producePodcastListWidget(configuration, width, height))
	widgets = append(widgets, produceControlsWidget(configuration, width, height))
	widgets = append(widgets, producePlayerWidget(configuration, width, height))
	return widgets
}

func drawPageSearch(configuration *Configuration, width int, height int) []ui.Bufferer {
	fillOutControlsMap(configuration, defaultControlsMap)
	widgets := make([]ui.Bufferer, 0)
	widgets = append(widgets, produceSearchWidget(configuration, width, height))
	widgets = append(widgets, produceSearchResultsWidget(configuration, width, height))
	widgets = append(widgets, produceControlsWidget(configuration, width, height))
	widgets = append(widgets, producePlayerWidget(configuration, width, height))
	return widgets
}

func drawPageDownloaded(configuration *Configuration, width int, height int) []ui.Bufferer {
	widgets := make([]ui.Bufferer, 0)
	widgets = append(widgets, produceDownloadedWidget(configuration, width, height))
	widgets = append(widgets, produceControlsWidget(configuration, width, height))
	widgets = append(widgets, producePlayerWidget(configuration, width, height))
	return widgets
}

func drawPagePodcastDetail(configuration *Configuration, width int, height int) []ui.Bufferer {
	widgets := make([]ui.Bufferer, 0)
	widgets = append(widgets, producePodcastDetailDescriptionWidget(configuration, width, height))
	widgets = append(widgets, producePodcastDetailListWidget(configuration, width, height))
	return widgets
}
