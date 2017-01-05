package data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
}

// OpenStorage is an utility method to get a storage interface
func OpenStorage() (storage.Provider, error) {
	conf := storage.DefaultConfig()
	if os.Getenv("TSV_STORAGE") != "" {
		conf.Path = os.Getenv("TSV_STORAGE")
	}
	return storage.New(conf)
}

// Run the query and return results
func (q *Query) Run() string {
	// Open storage
	store, err := OpenStorage()
	if err != nil {
		log.Println("Storage error:", err)
		return ""
	}
	defer store.Close()

	// Final result is an array of map interfaces
	res := []map[string]interface{}{}

	// Iterate bucket
	c := make(chan *storage.Record)
	go store.Cursor(q.Bucket, c)
CURSOR:
	for {
		select {
		case rec, ok := <-c:
			if !ok {
				break CURSOR
			}
			r, err := gabs.ParseJSON(rec.Value)
			if err == nil {
				// Unmarshal the contract document as a generic map interface
				m := make(map[string]interface{})
				json.Unmarshal(rec.Value, &m)

				// Inspect contract releases
				releases, _ := r.Search("releases").Children()
				for _, child := range releases {
					switch q.Filter {
					case "contractNumber":
						// releases[].ocid
						ocid, _ := child.Path("ocid").Data().(string)
						if ocid == q.Value {
							res = append(res, m)
						}
					case "procedureNumber":
						// releases[].planning.budget.projectID
						projectID, _ := child.Path("planning.budget.projectID").Data().(string)
						if projectID == q.Value {
							res = append(res, m)
						}
					case "buyer":
						// releases[].buyer.identifier.id
						// releases[].buyer.identifier.legalName
						buyerID, _ := child.Path("buyer.identifier.id").Data().(string)
						buyerName, _ := child.Path("buyer.identifier.legalName").Data().(string)
						if strings.Contains(buyerID, q.Value) || strings.Contains(buyerName, q.Value) {
							res = append(res, m)
						}
					case "date":
						// releases[].date
						date, _ := child.Path("date").Data().(string)
						if strings.Contains(date, q.Value) {
							res = append(res, m)
						}
					case "amount":
						// releases[].planning.budget.amount.amount
						amount, _ := child.Path("planning.budget.amount.amount").Data().(float64)
						qval, _ := strconv.ParseFloat(q.Value, 64)
						if amount == qval {
							res = append(res, m)
						}
					case "provider":
						awards, _ := child.Search("awards").Children()
						for _, award := range awards {
							// releases[].awards[].suppliers[].identifier.id
							providerID := []string{}
							json.Unmarshal([]byte(award.Path("suppliers.identifier.id").String()), &providerID)
							if len(providerID) > 0 {
								for _, p := range providerID {
									if strings.Contains(p, q.Value) {
										res = append(res, m)
									}
								}
							}

							// releases[].awards[].suppliers[].name
							providerName := []string{}
							json.Unmarshal([]byte(award.Path("suppliers.name").String()), &providerName)
							if len(providerName) > 0 {
								for _, p := range providerName {
									if strings.Contains(p, q.Value) {
										res = append(res, m)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	k, _ := json.Marshal(res)
	return fmt.Sprintf("%s", k)
}
