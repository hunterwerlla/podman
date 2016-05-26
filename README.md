#Podman - a terminal based podcast client
Podman is a terminal based podcast client written in Go.

#Default Key Binds
  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;K   
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;∧  
H&nbsp;<&nbsp;&nbsp;&nbsp;>&nbsp;L  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;∨  
&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;J  
 
/ to search

#Dependencies
Podman requires [gocui](https://github.com/jroimartin/gocui) for the interface, [go-sox](https://github.com/krig/go-sox) to play audio files, [sanitize](https://github.com/kennygrant/sanitize), [go-rss](https://github.com/ungerik/go-rss), and [pb](https://github.com/cheggaaa/pb). 
Searching is based on ITunes the largest source of podcasts
#Why?
I did not find any usable integrated podcast clients for the command line. There are a few that rely on RSS, but none that integrate searching and playback as well. I wanted a complete package like gPodder for the command line, but better.
#TODO
* Cache results more intelligently
* Download progress bar
* Properly pad playing bar
* Delete downloaded podcasts
* A list of downloaded podcasts
* Switch to using maps for subscribed podcasts and downloads
* more advanced searching (categories, etc)
* strip HTML from descriptions again
