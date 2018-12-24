Build: [![Build Status](https://travis-ci.org/hunterwerlla/podman.svg?branch=master)](https://travis-ci.org/hunterwerlla/podman)
### Podman - a terminal based podcast client
Podman is a terminal based podcast client written in Go.

### Default Key Binds
∧ to scroll up 

∨ to scroll down

&lt; to search

&gt; to see downloaded

PgUp to skip forward

PgDown to skip backward

&lt;spacebar&gt; to pause and resume

&lt;enter&gt; to do actions (play, download, switch to detailed view, subscribe)


### Dependencies
Podman requires:

[gocui](https://github.com/jroimartin/gocui) for the interface, [go-sox](https://github.com/krig/go-sox) to play audio files, [sanitize](https://github.com/kennygrant/sanitize), [go-rss](https://github.com/ungerik/go-rss), and [pb](https://github.com/cheggaaa/pb).

Podman utilizes go modules to simplify this.

Searching utilizes ITunes.

# TODO
* Fix total length of download progress bar
* fix graphical glitches
* more advanced searching (categories, etc)
* store data in a database
