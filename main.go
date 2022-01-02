package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/christophberger/3sixty/internal/fsapi"
	"github.com/christophberger/3sixty/internal/hifiberry"
)

// flags
var url string
var pin string

func main() {

	// test3sixty()
	testSoundStatus()
}

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
//
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

func testSoundStatus() {
	ctx := context.Background()
	statCh := monitorSoundStatus(ctx)
	for {
		select {
		case stat := <-statCh:
			fmt.Println(stat)
		case <-ctx.Done():
			return
		}
		<-time.After(1 * time.Second)
	}
}

func test3sixty() {
	url := flag.String("url", "http://k--che.fritz.box/fsapi", "API URL to 3sixty")
	pin := flag.String("pin", "0000", "PIN of 3sixty")
	flag.Parse()
	fs := fsapi.New(*url, *pin)
	err := fs.CreateSession()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(fs.Sid())

	err = fs.SetMode("7")
	if err != nil {
		log.Fatalln(err)
	}

	<-time.After(5 * time.Second)
	err = fs.SetMode("4")
	if err != nil {
		log.Fatalln(err)
	}

	power, err := fs.GetPowerStatus()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(power)

	err = fs.SetPowerStatus("0")
	if err != nil {
		log.Fatalln(err)
	}
}
