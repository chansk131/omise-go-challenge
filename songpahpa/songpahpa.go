package songpahpa

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/chansk131/omise-go-challenge/cipher"
)

type SongPahPa struct {
	Name     string
	Amount   int64
	CCNumber string
	CVV      string
	ExpMonth time.Month
	ExpYear  int
}

func ReadCSV(filepath string, songPahPaChannel chan<- *SongPahPa) {
	defer close(songPahPaChannel)

	data, err := os.Open(filepath)
	if err != nil {
		log.Println(err)
	}

	reader, err := cipher.NewRot128Reader(data)
	if err != nil {
		log.Println(err)
	}

	csvReader := csv.NewReader(reader)

	// Skip the header
	_, err = csvReader.Read()
	if err != nil {
		log.Println("Error read csv:", err)
		return
	}

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Error read csv:", err)
			break
		}

		amount, err := strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			log.Println("Error parsing Amount:", err)
			break
		}

		expMonthInt, err := strconv.Atoi(row[4])
		if err != nil {
			log.Println("Error parsing Month:", err)
			break
		}

		expYear, err := strconv.Atoi(row[5])
		if err != nil {
			log.Println("Error parsing year:", err)
			break
		}

		songPahPaChannel <- &SongPahPa{
			Name:     row[0],
			Amount:   amount,
			CCNumber: row[2],
			CVV:      row[3],
			ExpMonth: time.Month(expMonthInt),
			ExpYear:  expYear,
		}
	}
}
