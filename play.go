package main

import (
	"fmt"
	"github.com/krig/go-sox"
	"os"
	"time"
)

var (
	playerPosition int = -1
	startTime      time.Time
)

//this runs on its own thread to start/stop and select the media that is playing, it will also skip ahead in the future
//TODO make it skip ahead
//Control reference: 0 is play, 1 is pause, 2 is stop, 3 is skip ahead, 4 is reverse
func play(exit chan bool) {
	//get rid of all stderr data
	_, w, _ := os.Pipe()
	os.Stderr = w
	var (
		chain   *sox.EffectsChain = nil
		inFile  *sox.Format       = nil
		outFile *sox.Format       = nil
		status  int               = _nothing
	)
	if !sox.Init() {
		panic("Unable to start the player")
	}
	defer sox.Quit()
	for {
		status = _nothing
		toPlay := ""
		select {
		case toPlay = <-globals.playerFile:
		case status = <-globals.playerControl:
		}
		//if filname is not empty, then new filename recieved
		if toPlay != "" || globals.playerState == _play {
			if toPlay != "" {
				globals.Playing = toPlay
			}
			globals.playerState = _play
			inFile = sox.OpenRead(globals.Playing)
			if playerPosition == -1 {
				playerPosition = 0
			}
			//forward to the position
			//formula taken from example 2 of goSoX
			if playerPosition != 0 {
				seek := uint64(float64(playerPosition)*float64(inFile.Signal().Rate())*float64((inFile.Signal().Channels())) + 0.5)
				seek -= seek % uint64(inFile.Signal().Channels())
				inFile.Seek(seek)
			}
			outFile = sox.OpenWrite("default", inFile.Signal(), nil, "alsa")
			if outFile == nil {
				outFile = sox.OpenWrite("default", inFile.Signal(), nil, "pulseaudio")
				if outFile == nil {
					panic("Cannot open audio output devices")
				}
			}
			//Now actually play
			//play block
			chain = sox.CreateEffectsChain(inFile.Encoding(), outFile.Encoding())
			//make it output
			interm_signal := inFile.Signal().Copy()
			//set input
			e := sox.CreateEffect(sox.FindEffect("input"))
			e.Options(inFile)
			chain.Add(e, interm_signal, inFile.Signal())
			e.Release()
			//set output
			e = sox.CreateEffect(sox.FindEffect("output"))
			e.Options(outFile)
			chain.Add(e, interm_signal, inFile.Signal())
			e.Release()
			//start the timer which keeps track of position in the file
			startTime = time.Now()
			globals.LengthOfFile = getLengthOfFile(globals.Playing) //set length
			//process which also plays
			go chain.Flow()
		} else {
			switch status {
			case _nothing:
				//should not happen
			case _play: //case 0 play, only works after pause
				if playerPosition == -1 {
					fmt.Println("Have to select a file to play to resume playback")
				} else {
					globals.playerState = _play
					globals.playerFile <- globals.Playing
				}
			case _pause: //case 1 pause
				//save time and file
				playerPosition += int(time.Since(startTime).Seconds())
				globals.playerState = _pause
				//then stop and clear data
				if chain != nil {
					chain.Release()
					chain = nil
				}
				if inFile != nil {
					inFile.Release()
					inFile = nil
				}
				if outFile != nil {
					outFile.Release()
					outFile = nil
				}
			case _stop:
				//reset position
				playerPosition = -1
				globals.playerState = _nothing
				globals.Playing = ""
				globals.LengthOfFile = 0 //set length
				//then clean up
				if chain != nil {
					chain.Release()
					chain = nil
				}
				if inFile != nil {
					inFile.Release()
					inFile = nil
				}
				if outFile != nil {
					outFile.Release()
					outFile = nil
				}
			case _ff: //case 3 skip ahead
				//save time and file
				if playerPosition == -1 {
					fmt.Println("Have to select a file to play to resume playback")
				} else {
					playerPosition += int(time.Since(startTime).Seconds()) + globals.Config.forwardSkipLength
					//then stop and clear data
					if chain != nil {
						chain.Release()
						chain = nil
					}
					if inFile != nil {
						inFile.Release()
						inFile = nil
					}
					if outFile != nil {
						outFile.Release()
						outFile = nil
					}
				}
			case _rw: //case 4 rewind
				//save time and file
				if playerPosition == -1 {
					fmt.Println("Have to select a file to play to resume playback")
				} else {
					playerPosition += int(time.Since(startTime).Seconds()) - globals.Config.backwardSkipLength
					//then stop and clear data
					if chain != nil {
						chain.Release()
						chain = nil
					}
					if inFile != nil {
						inFile.Release()
						inFile = nil
					}
					if outFile != nil {
						outFile.Release()
						outFile = nil
					}
				}
			case _exit:
				goto exit //break out of loop for cleanup
			}
		}
	}
exit:
	if inFile != nil {
		inFile.Release()
	}
	if outFile != nil {
		outFile.Release()
	}
	if chain != nil {
		chain.Release()
	}
	exit <- true
}

func getLengthOfFile(fileName string) uint64 {
	inFile := sox.OpenRead(fileName)
	seek := uint64(float64(inFile.Signal().Length())/float64(inFile.Signal().Channels())/float64(inFile.Signal().Rate()) - 0.5)
	seek += seek % uint64(inFile.Signal().Channels())
	return seek
}
