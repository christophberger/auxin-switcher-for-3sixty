package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/christophberger/3sixty/internal/fsapi"
	"github.com/christophberger/3sixty/internal/hifiberry"
	"github.com/christophberger/3sixty/internal/systemctl"
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

// soundStatus reports the current status of the sound card
// and any status change since the last call
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

// setRadioFromSoundStatus checks the status of the sound card frequently and
// changes the radio's power status and mode accordingly.
func setRadioFromSoundStatus(a *app, fs *fsapi.Fsapi) error {
	status := soundStatus(a)
	switch status {
	case sndStatSwitchedOn:
		power, err := fs.GetPowerStatus()
		if err != nil {
			return fmt.Errorf("setRadioFromSoundStatus: cannot get power status: %w", err)
		}
		if power == fsapi.PowerOff {
			fs.SetPowerStatus(fsapi.PowerOn)
		}
		mode, err := fs.GetMode()
		if err != nil {
			return fmt.Errorf("setRadioFromSoundStatus: cannot get mode: %w", err)
		}
		if mode != fsapi.AuxIn {
			a.previousMode = mode
		}

		fs.SetMode(fsapi.AuxIn)
	case sndStatSwitchedOff:
		current, err := fs.GetMode()
		if err != nil {
			return fmt.Errorf("setRadioFromSoundStatus: cannot get mode: %w", err)
		}
		if current != fsapi.AuxIn {
			// Someone switched to another input while the Raspi player was playing
			// Leave the radio alone
			break
		}
		if a.previousMode != fsapi.AuxIn {
			err := fs.SetMode(a.previousMode)
			if err != nil {
				return fmt.Errorf("setRadioFromSoundStatus: cannot set mode: %w", err)
			}
		}
		err = fs.SetPowerStatus(fsapi.PowerOff)
		if err != nil {
			return fmt.Errorf("setRadioFromSoundStatus: cannot switch radio off: %w", err)
		}
	}
	return nil
}

// radioListens detects whether or not the radio is listening
// to aux in. If it is, the function returns true. If the radio is switched
// off or to another source, the function returns false.
// If querying the radio fails, the function returns false, assuming that
// the radio is not ready to play music from aux-in.
func radioListens(fs *fsapi.Fsapi) bool {
	power, _ := fs.GetPowerStatus()
	mode, _ := fs.GetMode()
	if power == fsapi.PowerOn && mode == fsapi.AuxIn {
		return true
	}
	return false
}

// stopSoundIfRadioStopsListening stops any sound output (from either
// Raspotify or shairport-sync) if the radio stops listening to aux-in.
func stopSoundIfRadioStopsListening(fs *fsapi.Fsapi) error {
	if !radioListens(fs) {
		isPlaying, err := hifiberry.IsPlaying()
		if err != nil {
			return fmt.Errorf("stopSoundIfRadioStopsListening: %w", err)
		}
		if isPlaying {
			err = systemctl.Restart("raspotify")
			if err != nil {
				return fmt.Errorf("stopSoundIfRadioStopsListening: %w", err)
			}
			err = systemctl.Restart("shairport-sync")
			if err != nil {
				return fmt.Errorf("stopSoundIfRadioStopsListening: %w", err)
			}
		}
	}
	return nil
}

func eventLoop(a *app, fs *fsapi.Fsapi) error {
	for {
		err := setRadioFromSoundStatus(a, fs)
		if err != nil {
			return fmt.Errorf("eventLoop: %w", err)
		}

		time.Sleep(1 * time.Second)

		err = stopSoundIfRadioStopsListening(fs)
		if err != nil {
			return fmt.Errorf("eventLoop: %w", err)
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime)
	url := flag.String("url", "http://CHANGE_ME/fsapi", "API URL to 3sixty")
	pin := flag.String("pin", "0000", "PIN of 3sixty")
	flag.Parse()
	fs := fsapi.New(*url, *pin)
	err := fs.CreateSession()
	if err != nil {
		log.Println(err)
		os.Exit(-1)
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
