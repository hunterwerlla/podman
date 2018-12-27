package main

import (
	ui "github.com/gizak/termui"
)

func produceHeaderWidget(width int) *ui.Paragraph {
	headerWidget := ui.NewParagraph(podmanHeader)
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
	} else {
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
}

func produceSearchWidget(configuration *Configuration, width int, height int) *ui.Paragraph {
	return nil
}

func produceSearchResults(configuration *Configuration, width int, height int) *ui.List {
	return nil
}
