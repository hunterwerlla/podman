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
Podman requires [gocui](https://github.com/jroimartin/gocui) for the interface, [go-sox](https://github.com/krig/go-sox) to play audio files, [sanitize](https://github.com/kennygrant/sanitize), and [go-rss](https://github.com/ungerik/go-rss).  
Searching is based on ITunes the largest source of podcasts
#Why?
I did not find any usable integrated podcast clients for the command line. There are a few that rely on RSS, but none that integrate searching and playback as well. I wanted a complete package like gPodder for the command line, but better.
#Why Go?
Initially I was going to write it in Rust since I really enjoy the language, but I saw this amazing HackerNews reader written in go, [go-hn](https://gitlab.com/shank/go-hn) which uses gocui which I found to be way better than any other command line interface library. 
#TODO
* Cache results more intelligently
* Download progress bar
* Properly pad playing bar
* Delete downloaded podcasts
* A list of downloaded podcasts
* Switch to using maps for subscribed podcasts and downloads
* more advanced searching
