package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/christophberger/3sixty/internal/fsapi"
)

// flags
var url string
var pin string

func main() {
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
}
