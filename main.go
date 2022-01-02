package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/christophberger/3sixty/internal/fsapi"
	"github.com/christophberger/3sixty/internal/hifiberry"
)

// flags
var url string
var pin string

type sndStat int

const (
	sndStatUnknown sndStat = iota
	sndStatOff
	sndStatSwitchedOn
	sndStatOn
	sndStatSwitchedOff
)

type app struct {
	previousSndStat sndStat
	previousMode    string
}

// soundStatus continuously delivers the status of the sound card
// through the returned channel.
//
// The busy loop blocks until the reader fetches the next value.
// This way, the receiver can decide upon when status checks happen.
// Otherwise, the loop would have to use a timer for pausing, thus
// imposing a fixed interval of status updates.
func soundStatus(a *app) sndStat {
	current := sndStatOff

	// save current status on return for detecting status changes
	defer func() {
		a.previousSndStat = current
	}()

	// determine the current status
	playing, err := hifiberry.IsPlaying()
	if err != nil {
		return sndStatUnknown
	}
	if playing {
		current = sndStatOn
	}

	// if the status has changed, return it; else return the current status
	switch {
	case current == sndStatOn && a.previousSndStat == sndStatOff:
		return sndStatSwitchedOn
	case current == sndStatOff && a.previousSndStat == sndStatOn:
		return sndStatSwitchedOff
	default:
		return current
	}
}

// radioListenStatus detects whether or not the radio is listening
// to aux in. If it is, the function returns true. If the radio is switched
// off or to another source, the function returns false.
// If querying the radio fails, the function returns false, assuming that
// the radio is not ready to play music from aux-in.
func radioListenStatus() bool {
	fs := fsapi.New(url, pin)
	power, _ := fs.GetPowerStatus()
	mode, _ := fs.GetMode()
	if power == fsapi.PowerOn && mode == fsapi.AuxIn {
		return true
	}
	return false
}

// evenLoop checks the status of the sound card frequently and
// changes the radio's power status and mode accordingly.
func eventLoop(a *app, fs *fsapi.Fsapi) error {
	for {
		status := soundStatus(a)
		switch status {
		case sndStatSwitchedOn:
			power, err := fs.GetPowerStatus()
			if err != nil {
				return fmt.Errorf("eventLoop: cannot get power status: %w", err)
			}
			if power == fsapi.PowerOff {
				fs.SetPowerStatus(fsapi.PowerOn)
			}
			mode, err := fs.GetMode()
			if err != nil {
				return fmt.Errorf("eventLoop: cannot get mode: %w", err)
			}
			if mode != fsapi.AuxIn {
				a.previousMode = mode
			}

			fs.SetMode(fsapi.AuxIn)
		case sndStatSwitchedOff:
			if a.previousMode != fsapi.AuxIn {
				err := fs.SetMode(fsapi.AuxIn)
				if err != nil {
					return fmt.Errorf("eventLoop: cannot set mode: %w", err)
				}
			}
			err := fs.SetPowerStatus(fsapi.PowerOff)
			if err != nil {
				return fmt.Errorf("eventLoop: cannot switch radio off: %w", err)
			}
		}
	}
}

func main() {
	url := flag.String("url", "http://k--che.fritz.box/fsapi", "API URL to 3sixty")
	pin := flag.String("pin", "0000", "PIN of 3sixty")
	flag.Parse()
	fs := fsapi.New(*url, *pin)
	err := fs.CreateSession()
	if err != nil {
		log.Fatalln(err)
	}

	a := &app{}

	// start the event loop. In case of an error, log the error
	// and restart the loop.
	for {
		err = eventLoop(a, fs)
		if err != nil {
			log.Println(err)
		}
	}
}
