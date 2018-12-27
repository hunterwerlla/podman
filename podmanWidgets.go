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
	headerWidget.WrapLength = 0
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

func producePlayerWidget(width int, height int) *ui.Paragraph {
	playerWidget := ui.NewParagraph("play something")
	playerWidget.Width = width
	playerWidget.Height = playerHeight
	playerWidget.Y = height - playerHeight
	playerWidget.BorderLeft = false
	playerWidget.BorderRight = false
	playerWidget.BorderBottom = false
	playerWidget.TextFgColor = ui.ColorBlack
	return playerWidget
}
