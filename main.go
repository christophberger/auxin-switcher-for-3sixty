package main

import (
	"flag"
	"log"

	"github.com/christophberger/3sixty/internal/fsapi"
	"github.com/christophberger/3sixty/internal/hifiberry"
)

// flags
var url string
var pin string

type sndStat int

const (
	sndStatOff sndStat = iota
	sndStatSwitchedOn
	sndStatOn
	sndStatSwitchedOff
	sndStatUnknown
)

type app struct {
	previous sndStat
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
		a.previous = current
	}()

	// determine the current status
	status, err := hifiberry.GetSoundStatus()
	if err != nil {
		return sndStatUnknown
	}
	if status {
		current = sndStatOn
	}

	// if the status has changed, return it; else return the current status
	switch {
	case current == sndStatOn && a.previous == sndStatOff:
		return sndStatSwitchedOn
	case current == sndStatOff && a.previous == sndStatOn:
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
	if power == fsapi.PowerOn && mode == fsapi.AuxInId {
		return true
	}
	return false
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
}
