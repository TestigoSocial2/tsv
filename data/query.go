package data

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/Jeffail/gabs"
	"github.com/bcessa/tsv/storage"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Query contract information
type Query struct {
	Value  string `json:"value"`
	Filter string `json:"filter"`
	Bucket string `json:"bucket"`
	Limit  int    `json:"limit"`
}

// Remove UNICODE diacritics
func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
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
				dates := strings.Split(q.Value, "|")
				if len(dates) > 1 {
					startDate, _ := time.Parse("01-02-2006", dates[0])
					endDate, _ := time.Parse("01-02-2006", dates[1])
					date, _ := release.Path("date").Data().(string)
					rdate, _ := time.Parse(time.RFC3339, date)
					if rdate.After(startDate) && rdate.Before(endDate) {
						list = append(list, m)
					}
				} else {
					startDate, _ := time.Parse("01-02-2006", dates[0])
					date, _ := release.Path("date").Data().(string)
					rdate, _ := time.Parse(time.RFC3339, date)
					sy, sm, sd := startDate.Date()
					ry, rm, rd := rdate.Date()
					if sy == ry && sm == rm && sd == rd {
						list = append(list, m)
					}
				}
			case "amount":
				// releases[].planning.budget.amount.amount
				barrier := strings.Split(q.Value, "|")
				low, _ := strconv.Atoi(barrier[0])
				high, _ := strconv.Atoi(barrier[1])
				amount, _ := release.Path("planning.budget.amount.amount").Data().(float64)
				if amount >= float64(low) && amount <= float64(high) {
					list = append(list, m)
				}
			case "contractNumber":
				// releases[].ocid
				ocid, _ := release.Path("ocid").Data().(string)
				if ocid == q.Value {
					list = append(list, m)
				}
			case "procedureNumber":
				// releases[].tender.id
				tenderID, _ := release.Path("tender.id").Data().(string)
				if tenderID == q.Value {
					list = append(list, m)
				}
			case "procedureType":
				// releases[].tender.procurementMethod
				pType, _ := release.Path("tender.procurementMethod").Data().(string)
				if pType == q.Value {
					list = append(list, m)
				}
			case "buyer":
				// releases[].buyer.identifier.id
				// releases[].buyer.identifier.legalName
				buyerID, _ := release.Path("buyer.identifier.id").Data().(string)
				buyerName, _ := release.Path("buyer.identifier.legalName").Data().(string)

				// Normalize input and values
				t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
				qValue, _, _ := transform.String(t, q.Value)
				qValue = strings.ToLower(qValue)
				buyerID, _, _ = transform.String(t, buyerID)
				buyerID = strings.ToLower(buyerID)
				buyerName, _, _ = transform.String(t, buyerName)
				buyerName = strings.ToLower(buyerName)

				if strings.Contains(buyerID, qValue) || strings.Contains(buyerName, qValue) {
					list = append(list, m)
				}
			case "provider":
				// Normalize input value
				t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
				qValue, _, _ := transform.String(t, q.Value)
				qValue = strings.ToLower(qValue)

				awards, _ := release.Search("awards").Children()
				for _, award := range awards {
					// releases[].awards[].suppliers[].identifier.id
					providerID := []string{}
					json.Unmarshal([]byte(award.Path("suppliers.identifier.id").String()), &providerID)
					for _, p := range providerID {
						p, _, _ = transform.String(t, p)
						p = strings.ToLower(p)
						if strings.Contains(p, qValue) {
							list = append(list, m)
						}
					}

					// releases[].awards[].suppliers[].name
					providerName := []string{}
					json.Unmarshal([]byte(award.Path("suppliers.name").String()), &providerName)
					for _, p := range providerName {
						p, _, _ = transform.String(t, p)
						p = strings.ToLower(p)
						if strings.Contains(p, qValue) {
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
