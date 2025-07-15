package songpahpa

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/chansk131/omise-go-challenge/cipher"
)

type SongPahPa struct {
	Name     string
	Amount   int64
	CCNumber string
	CVV      string
	ExpMonth string
	ExpYear  string
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

		songPahPaChannel <- &SongPahPa{
			Name:     row[0],
			Amount:   amount,
			CCNumber: row[2],
			CVV:      row[3],
			ExpMonth: row[4],
			ExpYear:  row[5],
		}
	}
}
