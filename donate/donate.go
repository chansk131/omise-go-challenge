package donate

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/chansk131/omise-go-challenge/songpahpa"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"golang.org/x/time/rate"
)

type Donation struct {
	Name    string
	Amount  int64
	Success bool
}

type Donator struct {
	client  *omise.Client
	limiter *rate.Limiter
}

func Initialise(publicKey string, secretKey string) *Donator {
	if len(publicKey) == 0 || len(secretKey) == 0 {
		panic(errors.New("your publicKey or secretKey is empty"))
	}

	client, err := omise.NewClient(publicKey, secretKey)
	if err != nil {
		log.Println(err)
	}

	limiter := rate.NewLimiter(1, 1) // 1 req/s, burst of 1
	return &Donator{client, limiter}
}

func (d *Donator) Donate(
	songPahPaChannel <-chan *songpahpa.SongPahPa,
	summaryChannel chan<- *Donation) {

	for songPahPa := range songPahPaChannel {

		isSuccess := d.createCharge(songPahPa)
		summaryChannel <- &Donation{
			Name:    songPahPa.Name,
			Amount:  songPahPa.Amount,
			Success: isSuccess,
		}
	}
}

func (d *Donator) createCharge(songPahPa *songpahpa.SongPahPa) bool {
	if songPahPa.ExpYear < time.Now().Year() {
		return false
	}

	err := d.limiter.Wait(context.Background())
	if err != nil {
		fmt.Println("Rate limiter error:", err)
		return false
	}

	token, createToken := &omise.Token{}, &operations.CreateToken{
		Name:            songPahPa.Name,
		Number:          songPahPa.CCNumber,
		ExpirationMonth: songPahPa.ExpMonth,
		ExpirationYear:  songPahPa.ExpYear,
		SecurityCode:    songPahPa.CVV,
	}
	if e := d.client.Do(token, createToken); e != nil {
		log.Println(e)
		return false
	}

	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   songPahPa.Amount,
		Currency: "thb",
		Card:     token.ID,
	}
	if e := d.client.Do(charge, createCharge); e != nil {
		log.Println(e)
		return false
	}

	return charge.Paid
}
