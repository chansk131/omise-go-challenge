package main

import (
	"fmt"
	"os"

	"github.com/chansk131/omise-go-challenge/songpahpa"
)

func main() {
	filepath := os.Args[1]

	fmt.Println("performing donations...")

	songPahPaChannel := make(chan *songpahpa.SongPahPa)
	go songpahpa.ReadCSV(filepath, songPahPaChannel)
	for songpahpa := range songPahPaChannel {
		fmt.Println(songpahpa)
	}
}
