package summary

import (
	"fmt"

	"github.com/chansk131/omise-go-challenge/donate"
)

type TopDonor struct {
	Name   string
	Amount int64
}

type Summary struct {
	Total     int64
	Success   int64
	Faulty    int64
	Count     uint
	Average   float64
	TopDonors []TopDonor
}

func GetSummary(donationChannel <-chan *donate.Donation) *Summary {
	summary := &Summary{
		Total:   0,
		Success: 0,
		Faulty:  0,
		Count:   0,
		Average: 0,
		TopDonors: []TopDonor{
			{Name: "", Amount: 0},
			{Name: "", Amount: 0},
			{Name: "", Amount: 0},
		},
	}

	for donation := range donationChannel {
		summary.Count++
		summary.Total += donation.Amount

		if donation.Success {
			summary.Success += donation.Amount
		} else {
			summary.Faulty += donation.Amount
		}
		sortTopDonors(donation, summary)
	}
	averageAmount := float64(0)
	if summary.Count > 0 {
		averageAmount = float64(summary.Total) / float64(summary.Count)
	}
	summary.Average = averageAmount

	return summary
}

func sortTopDonors(donation *donate.Donation, summary *Summary) {
	if !donation.Success {
		return
	}

	if donation.Amount > summary.TopDonors[0].Amount {
		summary.TopDonors[2] = summary.TopDonors[1]
		summary.TopDonors[1] = summary.TopDonors[0]
		summary.TopDonors[0] = TopDonor{
			Name:   donation.Name,
			Amount: donation.Amount,
		}
	} else if donation.Amount > summary.TopDonors[1].Amount {
		summary.TopDonors[2] = summary.TopDonors[1]
		summary.TopDonors[1] = TopDonor{
			Name:   donation.Name,
			Amount: donation.Amount,
		}
	} else if donation.Amount > summary.TopDonors[2].Amount {
		summary.TopDonors[2] = TopDonor{
			Name:   donation.Name,
			Amount: donation.Amount,
		}
	}

}

func (s *Summary) Print() {
	fmt.Printf("%25s THB %14d.00\n", "total received:", s.Total)
	fmt.Printf("%25s THB %14d.00\n", "successfully donated:", s.Success)
	fmt.Printf("%25s THB %14d.00\n", "faulty donation:", s.Faulty)
	fmt.Println()
	fmt.Printf("%25s THB %17.2f\n", "average per person:", s.Average)
	fmt.Printf("%25s %s\n", "top donors:", s.TopDonors[0].Name)
	fmt.Printf("%25s %s\n", "", s.TopDonors[1].Name)
	fmt.Printf("%25s %s\n", "", s.TopDonors[2].Name)
}
