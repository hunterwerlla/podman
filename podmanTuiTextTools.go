package main

func formatPodcast(p Podcast, max int) string {
	formattedString := p.CollectionName + " - " + p.ArtistName + " - " + p.Description
	if len(p.Description+p.CollectionName+p.ArtistName)+6 < max {
		//do nothing
	} else { //else truncate string
		formattedString = formattedString[0:max]
	}
	return formattedString
}
