package main

import (
	ui "github.com/gizak/termui"
)

func drawPageMain(configuration *Configuration, width int, height int) []ui.Bufferer {
	widgets := make([]ui.Bufferer, 0)
	currentPodcastsInBuffers[currentScreen] = configuration.Subscribed
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
	currentPodcastsInBuffers[currentScreen] = configuration.Downloaded
	widgets = append(widgets, produceDownloadedWidget(configuration, width, height))
	widgets = append(widgets, produceControlsWidget(configuration, width, height))
	widgets = append(widgets, producePlayerWidget(configuration, width, height))
	return widgets
}

func drawPagePodcastDetail(configuration *Configuration, width int, height int) []ui.Bufferer {
	widgets := make([]ui.Bufferer, 0)
	widgets = append(widgets, producePodcastDetailWidget(configuration, width, height))
	return widgets
}