#Podman - a terminal based podcast client
Podman is a terminal based podcast client written in Go.

#Default Key Binds
∧ to scroll up 
∨ to scroll down
< to search
> to see downloaded
PgUp to skip forward
PgDown to skip backward
<spacebar> to pause and resume
<enter> to do actions (play, download, switch to detailed view, subscribe)

#Dependencies
Podman requires [gocui](https://github.com/jroimartin/gocui) for the interface, [go-sox](https://github.com/krig/go-sox) to play audio files, [sanitize](https://github.com/kennygrant/sanitize), [go-rss](https://github.com/ungerik/go-rss), and [pb](https://github.com/cheggaaa/pb). 
Searching is based on ITunes the largest source of podcasts
#Why?
I did not find any usable integrated podcast clients for the command line. There are a few that rely on RSS, but none that integrate searching and playback as well. I wanted a complete package like gPodder for the command line, but better.
#TODO
* Fix total length of download progress bar
* more advanced searching (categories, etc)
* strip HTML from descriptions again
