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

package storage

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "fmt"
)

// Handler provides a storage usage point
type Handler struct {
  db *mgo.Database
}

// IndexOptions defines the available configuration values when creating a new index
type IndexOptions struct {
  Key            []string
  Unique         bool
  DropDuplicates bool
  Background     bool
}

// NewHandler creates a new handler interface on the provider
func NewHandler(host, db string) (*Handler, error) {
  session, err := mgo.Dial(host)
  session.SetSafe(&mgo.Safe{WMode: "majority"})
  if err != nil {
    return nil, err
  }
  
  return &Handler{db:session.DB(db)}, nil
}

// Close end any sessions/connections on the provider
func (h *Handler) Close() {
  h.db.Session.Close()
  return
}

// Count returns the number of items in a collection
func (h *Handler) Count(collection string) int {
  res := 0
  res, _ = h.db.C(collection).Count()
  return res
}

// Insert save entries to the requested collection
func (h *Handler) Insert(collection string, docs ...interface{}) error {
  return h.db.C(collection).Insert(docs...)
}

// BulkInsert will store several docs at once for faster writes
func (h *Handler) BulkInsert(collection string, docs []interface{}) error {
  b := h.db.C(collection).Bulk()
  b.Insert(docs...)
  _, err := b.Run()
  return err
}

// Update an existing record with the provided JSON encoded data
func (h *Handler) Update(collection string, uuid string, data string) error {
  // Turn provided data to BSON
  var d interface{}
  err := bson.UnmarshalJSON([]byte(data), &d)
  if err != nil {
    return err
  }
  
  // Perform update
  err = h.db.C(collection).Update(bson.M{"uuid":uuid}, d)
  if err != nil {
    return err
  }
  return nil
}

// Remove will permanently delete a record based on it's "uuid" field
func(h *Handler) Remove(collection string, uuid string) error {
  return h.db.C(collection).Remove(bson.M{"uuid":uuid})
}

// Retrieve an entry based on it's "uuid" field
func(h *Handler) Get(collection string, uuid string) (interface{}, error) {
  var r interface{}
  err := h.db.C(collection).Find(bson.M{"uuid":uuid}).One(&r)
  if err != nil {
    return nil, err
  }
  return r, nil
}

// Retrieve an entry based on it's "_id" field
func(h *Handler) GetByID(collection string, id string) (interface{}, error) {
  var r interface{}
  if ! bson.IsObjectIdHex(id) {
    return nil, fmt.Errorf("invalid id: %s", id)
  }
  err := h.db.C(collection).FindId(bson.ObjectIdHex(id)).One(&r)
  if err != nil {
    return nil, err
  }
  return r, nil
}

// Query run a JSON encoded query in the selected collection
func(h *Handler) Query(collection string, query string) (*mgo.Query, error) {
  var q interface{}
  err := bson.UnmarshalJSON([]byte(query), &q)
  if err != nil {
    return nil, err
  }
  return h.db.C(collection).Find(q), nil
}

// Index will create a new index in the collection;
// keys can be of the format: [$<kind>:][-]<field name>
func(h *Handler) Index(collection string, opts *IndexOptions) error {
  index := mgo.Index{
    Key: opts.Key,
    Unique: opts.Unique,
    DropDups: opts.DropDuplicates,
    Background: opts.Background,
  }
  return h.db.C(collection).EnsureIndex(index)
}