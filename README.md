Build: [![Build Status](https://travis-ci.org/hunterwerlla/podman.svg?branch=master)](https://travis-ci.org/hunterwerlla/podman)
### Podman - a terminal based podcast client
Podman is a terminal based podcast client written in Go with a fully features TUI and somewhat functional CUI

### Default Key Bindings
**K/&lt;Up&gt;** to scroll up<br/>
**J/&lt;Down&gt;** to scroll down<br/>
**H/&lt;**- to move left (search)<br/>
**L/-&gt;** to move right (downloaded)<br/>
**PgUp** to skip forward<br/>
**PgDown** to skip backward<br/>
**D** to delete downloads/subscriptions<br/>
**&lt;spacebar&gt;** to pause/resume when playing<br/>
**&lt;enter&gt;** to do actions (play, download, view all of the episodes, subscribe)<br/>
**/** to search podcast list/downloads/new podcasts to subscribe to

Keybinds are editable by changing ~/.config/podman/config.json

### Dependencies
Podman requires:

[termui](https://github.com/gizak/termui) for the interface,  
[faiface/beep](https://github.com/faiface/beep) to play audio files,  
[sanitize](https://github.com/kennygrant/sanitize) to clean up podcast descriptions,  
[go-rss](https://github.com/ungerik/go-rss) to grab podcast feeds  
[kr/text](https://github.com/kr/text) for text wrapping

Searching utilizes ITunes.

# TODO
* make search UX more obvious
* Skipping too much will crash or stop playing. It seems like an issue in the mp3 lib, but I need to look into it.
* re-add download bar
* add verification of downloads
* store user data in a database
* make it work on dark theme terminals
* add a settings menu
* make it work on Windows
* look at optimization
