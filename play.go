package main

import (
	"fmt"
	"github.com/krig/go-sox"
	"time"
)

var (
	playerPosition int = 0
	startTime      time.Time
)

//this runs on its own thread to start/stop and select the media that is playing, it will also skip ahead in the future
//TODO make it skip ahead
//Control reference: 0 is play, 1 is pause, 2 is stop, 3 is skip ahead, 4 is reverse
func play(exit chan bool) {
	var (
		chain          *sox.EffectsChain = nil
		inFile         *sox.Format       = nil
		outFile        *sox.Format       = nil
		cachedFileName string            = ""
		fileName       string            = ""
		status         int               = -1
	)
	if !sox.Init() {
		panic("Unable to start the player")
	}
	defer sox.Quit()
	for {
		//wait for a signal or go if we already have a file name
		if fileName != "" {
			status = -1
		} else {
			status = -1
			fileName = ""
			select {
			case fileName = <-globals.playerFile:
			case status = <-globals.playerControl:
			}
		}
		//if filname is not empty, then new filename recieved
		if fileName != "" && status == -1 {
			//when switching stop before playing new file
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
			inFile = sox.OpenRead(fileName)
			//forward to the position
			//formula taken from example 2 of SoX
			if playerPosition != 0 {
				seek := uint64(float64(playerPosition)*float64(inFile.Signal().Rate())*float64((inFile.Signal().Channels())) + 0.5)
				seek -= seek % uint64(inFile.Signal().Channels())
				inFile.Seek(seek)
			}
			//TODO make it work on windows too
			outFile = sox.OpenWrite("default", inFile.Signal(), nil, "alsa")
			if outFile == nil {
				outFile = sox.OpenWrite("default", inFile.Signal(), nil, "pulseaudio")
				if outFile == nil {
					panic("Cannot open audio output devices")
				}
			}
			//Now actually play
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
			//process which also plays
			globals.playerState = 0                          //change state to play
			globals.Playing = fileName                       //set filename
			globals.LengthOfFile = getLengthOfFile(fileName) //set length
			go chain.Flow()
			//reset status and filename
			fileName = ""
			status = -1
		} else if status != -1 {
			switch status {
			case -1:
				//should not happen TODO error
			case 0: //case 0 play, only works after pause
				if playerPosition == -1 {
					fmt.Println("Have to select a file to play to resume playback")
				} else {
					//TODO fix this channel issue which makes no sense and makes the code worse
					//fileNameChannel <- cachedFileName
					fileName = cachedFileName
				}
			case 1: //case 1 pause
				//save time and file
				playerPosition += int(time.Since(startTime).Seconds())
				cachedFileName = inFile.Filename()
				globals.playerState = 1
				globals.Playing = ""
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
			case 2: //case 2 stop
				//reset position
				playerPosition = 0
				globals.playerState = -1
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
			case 3: //case 3 skip ahead
				//save time and file
				if playerPosition == -1 {
					fmt.Println("Have to select a file to play to resume playback")
				} else {
					playerPosition += int(time.Since(startTime).Seconds()) + globals.Config.forwardSkipLength
					fileName = inFile.Filename()
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
			case 4: //case 4 rewind
				//save time and file
				if playerPosition == -1 {
					fmt.Println("Have to select a file to play to resume playback")
				} else {
					playerPosition += int(time.Since(startTime).Seconds()) - globals.Config.backwardSkipLength
					fileName = inFile.Filename()
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
			case 5: //exit
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
