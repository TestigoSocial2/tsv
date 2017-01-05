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
	FirstDate    time.Time `json:"firstDate"`
	LastDate     time.Time `json:"lastDate"`
	Contracts    contracts `json:"contracts"`
	AssignMethod struct {
		Direct  contracts `json:"direct"`
		Limited contracts `json:"limited"`
		Public  contracts `json:"public"`
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
	for _, child := range releases {
		// date
		date, _ := child.Path("date").Data().(string)
		t, err := time.Parse("2006-01-02T15:04:05.000Z", date)
		if err == nil {
			if t.Before(org.FirstDate) {
				org.FirstDate = t
			}
			if t.After(org.LastDate) {
				org.LastDate = t
			}
		}

		// planning.budget.amount.amount
		amount, ok := child.Path("planning.budget.amount.amount").Data().(float64)
		if ok {
			org.Contracts.Budget += amount
		}

		// tender.status
		status, ok := child.Path("tender.status").Data().(string)
		if ok {
			switch status {
			case activeState:
				org.Contracts.Active++
			case "complete":
				org.Contracts.Completed++
			}
		}

		// contracts.value.amount
		contracts, _ := child.Search("contracts").Children()
		for _, contract := range contracts {
			award, ok := contract.Path("value.amount").Data().(float64)
			if ok {
				org.Contracts.Awarded += award
			}
		}

		// tender.numberOfTenderers
		if child.ExistsP("tender.numberOfTenderers") {
			participants, _ := child.Path("tender.numberOfTenderers").Data().(float64)
			switch {
			case (participants == 1):
				org.AssignMethod.Direct.Total++
				org.AssignMethod.Direct.Budget += amount
				if status != activeState {
					org.AssignMethod.Direct.Active++
				} else {
					org.AssignMethod.Direct.Completed++
				}
				break
			case (participants >= 1 && participants <= 3):
				org.AssignMethod.Limited.Total++
				org.AssignMethod.Limited.Budget += amount
				if status != activeState {
					org.AssignMethod.Limited.Active++
				} else {
					org.AssignMethod.Limited.Completed++
				}
				break
			default:
				org.AssignMethod.Public.Total++
				org.AssignMethod.Public.Budget += amount
				if status != activeState {
					org.AssignMethod.Public.Active++
				} else {
					org.AssignMethod.Public.Completed++
				}
			}
		}
	}
	return nil
}
