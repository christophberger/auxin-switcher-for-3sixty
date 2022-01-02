package main

import (
	"context"
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

// soundStatus continuously delivers the status of the sound card
// through the returned channel.
//
// The busy loop blocks until the reader fetches the next value.
// This way, the receiver can decide upon when status checks happen.
// Otherwise, the loop would have to use a timer for pausing, thus
// imposing a fixed interval of status updates.
func monitorSoundStatus(ctx context.Context) chan sndStat {
	previous, current := sndStatUnknown, sndStatUnknown
	statCh := make(chan sndStat)
	go func() {
		for {
			status, err := hifiberry.GetSoundStatus()
			if err != nil {
				statCh <- sndStatUnknown
			}
			current = sndStatOff
			if status {
				current = sndStatOn
			}
			switch {
			case current == sndStatOn && previous == sndStatOff:
				statCh <- sndStatSwitchedOn
			case current == sndStatOff && previous == sndStatOn:
				statCh <- sndStatSwitchedOff
			default:
				statCh <- current
			}
			previous = current
		}
	}()
	return statCh
}

// monitorRadioListenStatus detects whether or not the radio is listening
// to aux in. If it is, the function returns true. If the radio is switched
// off or to another source, the function returns false.
// If querying the radio fails, the function returns false, assuming that
// the radio is not ready to play music.
func monitorRadioListenStatus(ctx context.Context) chan bool {
	statCh := make(chan bool)
	fs := fsapi.New(url, pin)
	go func() {
		for {
			// No error checking for the following two calls.
			// If the calls fail, the radio is probably not
			// ready to listen.
			power, _ := fs.GetPowerStatus()
			mode, _ := fs.GetMode()
			if power == fsapi.PowerOn && mode == fsapi.AuxInId {
				statCh <- true
				continue
			}
		}
	}()
	return statCh
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
