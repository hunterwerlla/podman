package main

import (
	"fmt"
	"github.com/bbrks/wrap"
	"strings"
	"unicode/utf8"
)

func formatPodcast(p Podcast, max int) string {
	var formatBuilder strings.Builder
	formatBuilder.WriteString(p.CollectionName)
	if p.ArtistName != "" {
		formatBuilder.WriteString(" - ")
		formatBuilder.WriteString(p.ArtistName)
	}
	if p.Description != "" {
		formatBuilder.WriteString(" - ")
		formatBuilder.WriteString(p.Description)
	}
	formattedString := formatBuilder.String()
	if len(formattedString) > max {
		formattedString = formattedString[0:max]
	}
	return formattedString
}

func wrapOrBreakText(configuration *Configuration, formattedPodcast string, width int, selected bool) string {
	if selected {
		formattedPodcast = wrapString(formattedPodcast, width)
		formattedPodcast = termuiStyleText(formattedPodcast, "white", "black")
	} else if utf8.RuneCountInString(formattedPodcast) > width {
		formattedPodcast = substringUTF8(formattedPodcast, 0, width-3)
		formattedPodcast += "..."
	}
	return formattedPodcast
}

func formatPodcastEpisode(p PodcastEpisode) string {
	var formatBuilder strings.Builder
	formatBuilder.WriteString(p.Title)
	if p.Summary != "" {
		formatBuilder.WriteString(" - ")
		formatBuilder.WriteString(p.Summary)
	}
	if p.Content != "" {
		formatBuilder.WriteString(" - ")
		formatBuilder.WriteString(p.Content)
	}
	return formatBuilder.String()
}

func substringUTF8(input string, begin int, end int) string {
	stringStart := 0
	i := 0
	for j := range input {
		if i == begin {
			stringStart = j
		}
		if i == end {
			return input[stringStart:j]
		}
		i++
	}
	return input[stringStart:]
}

func wrapString(input string, max int) string {
	w := wrap.NewWrapper()
	w.StripTrailingNewline = true
	output := w.Wrap(input, max)
	return output
}

// lovingly stolen from: https://programming.guide/go/formatting-byte-size-to-human-readable-format.html
func byteCountDecimal(number int64) string {
	const unit = 1000
	if number < unit {
		return fmt.Sprintf("%d B", number)
	}
	div, exp := int64(unit), 0
	for n := number / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(number)/float64(div), "kMGTPE"[exp])
}
