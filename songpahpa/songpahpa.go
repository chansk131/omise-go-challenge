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

func InitialiseReader(filepath string) *csv.Reader {
	data, err := os.Open(filepath)
	if err != nil {
		log.Println(err)
	}

	reader, err := cipher.NewRot128Reader(data)
	if err != nil {
		log.Println(err)
	}

	csvReader := csv.NewReader(reader)

	return csvReader
}

func ReadCSV(csvReader *csv.Reader, songPahPaChannel chan<- *SongPahPa) {
	defer close(songPahPaChannel)

	// Skip the header
	_, err := csvReader.Read()
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

		songPahPa, err := ParseSongPahPa(row)
		if err != nil {
			continue
		}
		songPahPaChannel <- songPahPa
	}
}

func ParseSongPahPa(row []string) (*SongPahPa, error) {
	amount, err := strconv.ParseInt(row[1], 10, 64)
	if err != nil {
		log.Println("Error parsing Amount:", err)
		return nil, err
	}

	expMonthInt, err := strconv.Atoi(row[4])
	if err != nil {
		log.Println("Error parsing Month:", err)
		return nil, err
	}

	expYear, err := strconv.Atoi(row[5])
	if err != nil {
		log.Println("Error parsing year:", err)
		return nil, err
	}

	songPahPa := &SongPahPa{
		Name:     row[0],
		Amount:   amount,
		CCNumber: row[2],
		CVV:      row[3],
		ExpMonth: time.Month(expMonthInt),
		ExpYear:  expYear,
	}

	return songPahPa, nil
}
