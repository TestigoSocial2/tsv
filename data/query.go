package data

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/bcessa/tsv/storage"
)

// Query contract information
type Query struct {
	Value  string `json:"value"`
	Filter string `json:"filter"`
	Bucket string `json:"bucket"`
	Limit  int    `json:"limit"`
}

// Run the query and return results
func (q *Query) Run() []map[string]interface{} {
	// Open storage
	store, err := OpenStorage()
	if err != nil {
		log.Println("Storage error:", err)
		return nil
	}
	defer store.Close()

	// Final result is an array of map interfaces
	list := []map[string]interface{}{}

	// Iterate bucket
	cursor := make(chan *storage.Record)
	cancel := make(chan bool)
	go store.Cursor(q.Bucket, cursor, cancel)
	for rec := range cursor {
		// Skip faulty records
		r, err := gabs.ParseJSON(rec.Value)
		if err != nil {
			continue
		}

		// Unmarshal the contract document as a generic map interface
		m := make(map[string]interface{})
		json.Unmarshal(rec.Value, &m)

		// Inspect contract releases
		releases, _ := r.Search("releases").Children()
		for _, release := range releases {
			switch q.Filter {
			case "date":
				// releases[].date
				date, _ := release.Path("date").Data().(string)
				if strings.Contains(date, q.Value) {
					list = append(list, m)
				}
			case "amount":
				// releases[].planning.budget.amount.amount
				amount, _ := release.Path("planning.budget.amount.amount").Data().(float64)
				qval, _ := strconv.ParseFloat(q.Value, 64)
				if amount == qval {
					list = append(list, m)
				}
			case "contractNumber":
				// releases[].ocid
				ocid, _ := release.Path("ocid").Data().(string)
				if ocid == q.Value {
					list = append(list, m)
				}
			case "procedureNumber":
				// releases[].planning.budget.projectID
				projectID, _ := release.Path("planning.budget.projectID").Data().(string)
				if projectID == q.Value {
					list = append(list, m)
				}
			case "buyer":
				// releases[].buyer.identifier.id
				// releases[].buyer.identifier.legalName
				buyerID, _ := release.Path("buyer.identifier.id").Data().(string)
				buyerName, _ := release.Path("buyer.identifier.legalName").Data().(string)
				if strings.Contains(buyerID, q.Value) || strings.Contains(buyerName, q.Value) {
					list = append(list, m)
				}
			case "provider":
				awards, _ := release.Search("awards").Children()
				for _, award := range awards {
					// releases[].awards[].suppliers[].identifier.id
					providerID := []string{}
					json.Unmarshal([]byte(award.Path("suppliers.identifier.id").String()), &providerID)
					for _, p := range providerID {
						if strings.Contains(p, q.Value) {
							list = append(list, m)
						}
					}

					// releases[].awards[].suppliers[].name
					providerName := []string{}
					json.Unmarshal([]byte(award.Path("suppliers.name").String()), &providerName)
					for _, p := range providerName {
						if strings.Contains(p, q.Value) {
							list = append(list, m)
						}
					}
				}
			}
		}

		// Check query limit
		if q.Limit > 0 && len(list) == q.Limit {
			close(cancel)
			break
		}
	}

	return list
}
