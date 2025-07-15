package donate

import (
	"errors"
	"log"
	"time"

	"github.com/chansk131/omise-go-challenge/songpahpa"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type Donation struct {
	Name    string
	Amount  int64
	Success bool
}

func InitialiseClient(publicKey string, secretKey string) *omise.Client {
	if len(publicKey) == 0 || len(secretKey) == 0 {
		panic(errors.New("your publicKey or secretKey is empty"))
	}

	client, err := omise.NewClient(publicKey, secretKey)
	if err != nil {
		log.Println(err)
	}
	return client
}

func Donate(client *omise.Client,
	songPahPaChannel <-chan *songpahpa.SongPahPa,
	summaryChannel chan<- *Donation) {
	for songPahPa := range songPahPaChannel {
		isSuccess := createCharge(client, songPahPa)
		summaryChannel <- &Donation{
			Name:    songPahPa.Name,
			Amount:  songPahPa.Amount,
			Success: isSuccess,
		}
	}
}

func createCharge(client *omise.Client, songPahPa *songpahpa.SongPahPa) bool {
	if songPahPa.ExpYear < time.Now().Year() {
		return false
	}

	time.Sleep(200 * time.Millisecond)

	token, createToken := &omise.Token{}, &operations.CreateToken{
		Name:            songPahPa.Name,
		Number:          songPahPa.CCNumber,
		ExpirationMonth: songPahPa.ExpMonth,
		ExpirationYear:  songPahPa.ExpYear,
		SecurityCode:    songPahPa.CVV,
	}
	log.Println(createToken)
	if e := client.Do(token, createToken); e != nil {
		log.Println(e)
		return false
	}

	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   songPahPa.Amount,
		Currency: "thb",
		Card:     token.ID,
	}
	if e := client.Do(charge, createCharge); e != nil {
		log.Println(e)
		return false
	}

	return charge.Paid
}
