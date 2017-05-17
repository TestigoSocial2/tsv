// Copyright © 2016 Transparencia Mexicana AC. <ben@pixative.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package data

import (
  "encoding/json"
  "github.com/Jeffail/gabs"
  "strconv"
  "time"
)

// IndicatorsEntry represents a single statistic value item
type IndicatorsEntry struct {
  Count  int     `json:"count"`
  Amount float64 `json:"amount"`
}

// IndicatorsQueryResult defines the results returned on statistics queries
type IndicatorsQueryResult struct {
  Years     map[string]*IndicatorsEntry `json:"years"`
  Limited   IndicatorsEntry             `json:"limited"`
  Selective IndicatorsEntry             `json:"selective"`
  Open      IndicatorsEntry             `json:"open"`
}

// FormatIndicatorsResult return a list of contracts as valid indicators data
func FormatIndicatorsResult(list []interface{}) *IndicatorsQueryResult {
  res := &IndicatorsQueryResult{}
  res.Years = make(map[string]*IndicatorsEntry)
  
  for _, rec := range list {
    json, _ := json.Marshal(rec)
    r, _ := gabs.ParseJSON(json)
    
    // Inspect contract releases
    releases, _ := r.Search("releases").Children()
    for _, release := range releases {
      amount, _ := release.Path("planning.budget.amount.amount").Data().(float64)
      date, _ := release.Path("date").Data().(string)
      rdate, _ := time.Parse(time.RFC3339, date)
  
      _, ok := res.Years[strconv.Itoa(rdate.Year())]
      if ok {
        res.Years[strconv.Itoa(rdate.Year())].Count++
        res.Years[strconv.Itoa(rdate.Year())].Amount += amount
      } else {
        res.Years[strconv.Itoa(rdate.Year())] = &IndicatorsEntry{
          Count:  1,
          Amount: amount}
      }
  
      pType, _ := release.Path("tender.procurementMethod").Data().(string)
      switch pType {
      case "selective":
        res.Selective.Count++
        res.Selective.Amount += amount
      case "limited":
        res.Limited.Count++
        res.Limited.Amount += amount
      case "open":
        res.Open.Count++
        res.Open.Amount += amount
      }
    }
  }
  return res
}
