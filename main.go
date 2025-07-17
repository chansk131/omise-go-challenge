package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/chansk131/omise-go-challenge/donate"
	"github.com/chansk131/omise-go-challenge/songpahpa"
	"github.com/chansk131/omise-go-challenge/summary"
	"github.com/joho/godotenv"
)

func main() {
	filepath := os.Args[1]

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	publicKey := os.Getenv("OMISE_PUBLIC_KEY")
	secretKey := os.Getenv("OMISE_SECRET_KEY")
	donator := donate.Initialise(publicKey, secretKey)

	fmt.Println("performing donations...")

	songPahPaChannel := make(chan *songpahpa.SongPahPa)
	reader := songpahpa.InitialiseReader(filepath)
	go songpahpa.ReadCSV(reader, songPahPaChannel)

	donationChannel := make(chan *donate.Donation)

	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup

	for i := range numWorkers {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			donator.Donate(songPahPaChannel, donationChannel)
		}(i)
	}

	go func() {
		wg.Wait()
		close(donationChannel)
	}()

	summary := summary.GetSummary(donationChannel)

	fmt.Println("done.")
	summary.Print()
}
