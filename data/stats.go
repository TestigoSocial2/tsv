package data

import (
	"time"

	"github.com/Jeffail/gabs"
)

// Relevant contracts details
type contracts struct {
	Budget    float64 `json:"budget"`
	Awarded   float64 `json:"awarded"`
	Total     int     `json:"total"`
	Active    int     `json:"active"`
	Completed int     `json:"completed"`
}

// Organization represent the information summary of a specific organization
type Organization struct {
	Code         string    `json:"code"`
	Description  string    `json:"description"`
	FirstDate    time.Time `json:"firstDate"`
	LastDate     time.Time `json:"lastDate"`
	Contracts    contracts `json:"contracts"`
	AssignMethod struct {
		Limited   contracts `json:"limited"`
		Selective contracts `json:"selective"`
		Open      contracts `json:"open"`
	} `json:"method"`
}

// NewOrganization initialize a organization entry
func NewOrganization() *Organization {
	org := &Organization{}
	org.FirstDate = time.Now()
	org.LastDate = time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	return org
}

// AddRecord aggregate the information of a specific contract to the organization entry
func (org *Organization) AddRecord(rec []byte) error {
	r, err := gabs.ParseJSON(rec)
	if err != nil {
		return err
	}

	activeState := "active"
	org.Contracts.Total++
	releases, _ := r.Search("releases").Children()
	for _, release := range releases {
		// date
		date, _ := release.Path("date").Data().(string)
		t, err := time.Parse(time.RFC3339, date)
		if err == nil {
			if t.Before(org.FirstDate) {
				org.FirstDate = t
			}
			if t.After(org.LastDate) {
				org.LastDate = t
			}
		}

		// planning.budget.amount.amount
		amount, ok := release.Path("planning.budget.amount.amount").Data().(float64)
		if ok {
			org.Contracts.Budget += amount
		}

		// tender.status
		status, ok := release.Path("tender.status").Data().(string)
		if ok {
			switch status {
			case activeState:
				org.Contracts.Active++
			case "complete":
				org.Contracts.Completed++
			}
		}

		// contracts.value.amount
		contracts, _ := release.Search("contracts").Children()
		for _, contract := range contracts {
			award, ok := contract.Path("value.amount").Data().(float64)
			if ok {
				org.Contracts.Awarded += award
			}
		}

		// tender.numberOfTenderers
		if release.ExistsP("tender.procurementMethod") {
			procurementMethod, _ := release.Path("tender.procurementMethod").Data().(string)
			switch procurementMethod {
			case "limited":
				org.AssignMethod.Limited.Total++
				org.AssignMethod.Limited.Budget += amount
				if status != activeState {
					org.AssignMethod.Limited.Active++
				} else {
					org.AssignMethod.Limited.Completed++
				}
				break
			case "selective":
				org.AssignMethod.Selective.Total++
				org.AssignMethod.Selective.Budget += amount
				if status != activeState {
					org.AssignMethod.Selective.Active++
				} else {
					org.AssignMethod.Selective.Completed++
				}
				break
			case "open":
				org.AssignMethod.Open.Total++
				org.AssignMethod.Open.Budget += amount
				if status != activeState {
					org.AssignMethod.Open.Active++
				} else {
					org.AssignMethod.Open.Completed++
				}
			}
		}
	}
	return nil
}
