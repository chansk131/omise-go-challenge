package summary

import (
	"testing"

	"github.com/chansk131/omise-go-challenge/donate"
)

func TestGetSummary_AllSuccess(t *testing.T) {
	donations := []*donate.Donation{
		{Name: "Alice", Amount: 100, Success: true},
		{Name: "Bob", Amount: 200, Success: true},
		{Name: "Charlie", Amount: 300, Success: true},
	}
	ch := make(chan *donate.Donation, len(donations))
	for _, d := range donations {
		ch <- d
	}
	close(ch)

	summary := GetSummary(ch)

	if summary.Total != 600 {
		t.Errorf("expected Total 600, got %d", summary.Total)
	}
	if summary.Success != 600 {
		t.Errorf("expected Success 600, got %d", summary.Success)
	}
	if summary.Faulty != 0 {
		t.Errorf("expected Faulty 0, got %d", summary.Faulty)
	}
	if summary.Count != 3 {
		t.Errorf("expected Count 3, got %d", summary.Count)
	}
	if summary.Average != 200 {
		t.Errorf("expected Average 200, got %f", summary.Average)
	}
	if summary.TopDonors[0].Name != "Charlie" {
		t.Errorf("expected top donor Charlie, got %s", summary.TopDonors[0].Name)
	}
	if summary.TopDonors[1].Name != "Bob" {
		t.Errorf("expected second donor Bob, got %s", summary.TopDonors[1].Name)
	}
	if summary.TopDonors[2].Name != "Alice" {
		t.Errorf("expected third donor Alice, got %s", summary.TopDonors[2].Name)
	}
}

func TestGetSummary_WithFaulty(t *testing.T) {
	donations := []*donate.Donation{
		{Name: "Alice", Amount: 100, Success: true},
		{Name: "Bob", Amount: 200, Success: false},
		{Name: "Charlie", Amount: 300, Success: true},
	}
	ch := make(chan *donate.Donation, len(donations))
	for _, d := range donations {
		ch <- d
	}
	close(ch)

	summary := GetSummary(ch)

	if summary.Total != 600 {
		t.Errorf("expected Total 600, got %d", summary.Total)
	}
	if summary.Success != 400 {
		t.Errorf("expected Success 400, got %d", summary.Success)
	}
	if summary.Faulty != 200 {
		t.Errorf("expected Faulty 200, got %d", summary.Faulty)
	}
	if summary.Count != 3 {
		t.Errorf("expected Count 3, got %d", summary.Count)
	}
	if summary.Average != 200 {
		t.Errorf("expected Average 200, got %f", summary.Average)
	}
	if summary.TopDonors[0].Name != "Charlie" {
		t.Errorf("expected top donor Charlie, got %s", summary.TopDonors[0].Name)
	}
	if summary.TopDonors[1].Name != "Alice" {
		t.Errorf("expected second donor Alice, got %s", summary.TopDonors[1].Name)
	}
	if summary.TopDonors[2].Name != "" {
		t.Errorf("expected third donor empty, got %s", summary.TopDonors[2].Name)
	}
}

func TestGetSummary_Empty(t *testing.T) {
	ch := make(chan *donate.Donation)
	close(ch)

	summary := GetSummary(ch)

	if summary.Total != 0 {
		t.Errorf("expected Total 0, got %d", summary.Total)
	}
	if summary.Success != 0 {
		t.Errorf("expected Success 0, got %d", summary.Success)
	}
	if summary.Faulty != 0 {
		t.Errorf("expected Faulty 0, got %d", summary.Faulty)
	}
	if summary.Count != 0 {
		t.Errorf("expected Count 0, got %d", summary.Count)
	}
	if summary.Average != 0 {
		t.Errorf("expected Average 0, got %f", summary.Average)
	}
	for i := 0; i < 3; i++ {
		if summary.TopDonors[i].Name != "" || summary.TopDonors[i].Amount != 0 {
			t.Errorf("expected empty top donor at %d, got %s %d", i, summary.TopDonors[i].Name, summary.TopDonors[i].Amount)
		}
	}
}

func TestSortTopDonors_Order(t *testing.T) {
	summary := &Summary{
		TopDonors: []TopDonor{
			{Name: "", Amount: 0},
			{Name: "", Amount: 0},
			{Name: "", Amount: 0},
		},
	}
	donations := []*donate.Donation{
		{Name: "A", Amount: 10, Success: true},
		{Name: "B", Amount: 20, Success: true},
		{Name: "C", Amount: 30, Success: true},
		{Name: "D", Amount: 25, Success: true},
		{Name: "F", Amount: 50, Success: false},
	}
	for _, d := range donations {
		sortTopDonors(d, summary)
	}
	if summary.TopDonors[0].Name != "C" || summary.TopDonors[0].Amount != 30 {
		t.Errorf("expected top donor C 30, got %s %d", summary.TopDonors[0].Name, summary.TopDonors[0].Amount)
	}
	if summary.TopDonors[1].Name != "D" || summary.TopDonors[1].Amount != 25 {
		t.Errorf("expected second donor D 25, got %s %d", summary.TopDonors[1].Name, summary.TopDonors[1].Amount)
	}
	if summary.TopDonors[2].Name != "B" || summary.TopDonors[2].Amount != 20 {
		t.Errorf("expected third donor B 20, got %s %d", summary.TopDonors[2].Name, summary.TopDonors[2].Amount)
	}
}
