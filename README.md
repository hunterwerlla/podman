#Podman - a terminal based podcast client
Podman is a terminal based podcast client written in Go.

#Default Key Binds
  K 
  ^
H< >L
  ∨
  J

/ to search

#Dependencies
Podman requires [gocui](https://github.com/jroimartin/gocui) for the interface and [go-sox](https://github.com/krig/go-sox) to play audio files. The search relies on ITunes beause it is the largest source of podcasts on the internet (but podcasts are downloaded from their respective RSS feeds).
#Why?
I did not find any usable integrated podcast clients for the command line. There are a few that rely on RSS, but none that integrate searching and playback as well. I wanted a complete package like gPodder for the command line, but better.
#Why Go?
Initially I was going to write it in Rust since I really enjoy the language, but I saw this amazing HackerNews reader written in go, [go-hn](https://gitlab.com/shank/go-hn) which uses gocui which I found to be way better than any other command line interface library. 
