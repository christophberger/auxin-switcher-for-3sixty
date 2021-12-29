package main

import (
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
	soundCardStatus()
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

func soundCardStatus() {
	for {
		fmt.Println(hifiberry.GetSoundStatus())
		<-time.After(1 * time.Second)
	}
}
