package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/bcessa/tsv/compressor"
	"github.com/boltdb/bolt"
)

// boltStorage data store implementation using BoltDB as provider
type boltStorage struct {
	db         *bolt.DB
	state      sync.Mutex
	compressor compressor.Provider
}

func (storage *boltStorage) Open(path string) error {
	// Open database
	db, err := bolt.Open(path, 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("Error opening storage file")
	}
	storage.db = db

	// Create main bucket
	storage.state.Lock()
	defer storage.state.Unlock()
	return storage.db.Update(func(tx *bolt.Tx) error {
		tx.WriteFlag = txWriteFlag
		tx.CreateBucketIfNotExists([]byte("main"))
		return nil
	})
}

func (storage *boltStorage) Close() error {
	storage.state.Lock()
	defer storage.state.Unlock()
	return storage.db.Close()
}

func (storage *boltStorage) Write(bucket string, key, val []byte) error {
	storage.state.Lock()
	defer storage.state.Unlock()
	return storage.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		tx.WriteFlag = txWriteFlag
		return tx.Bucket([]byte(bucket)).Put(key, storage.compressor.Compress(val))
	})
}

func (storage *boltStorage) Read(bucket string, key []byte) (val []byte) {
	storage.db.View(func(tx *bolt.Tx) error {
		tx.WriteFlag = txWriteFlag
		val = storage.compressor.Decompress(tx.Bucket([]byte(bucket)).Get(key))
		return nil
	})
	return val
}

func (storage *boltStorage) Cursor(bucket string, s chan *Record) {
	storage.db.View(func(tx *bolt.Tx) error {
		tx.WriteFlag = txWriteFlag
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			s <- &Record{
				Key:   k,
				Value: storage.compressor.Decompress(v),
			}
		}
		return nil
	})
	return
}

func (storage *boltStorage) GetStats(bucket string) (s *Stats) {
	s = &Stats{}
	s.FirstDate = time.Now()
	s.LastDate = time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	storage.db.View(func(tx *bolt.Tx) error {
		tx.WriteFlag = txWriteFlag
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			val := storage.compressor.Decompress(v)
			r, err := gabs.ParseJSON(val)
			if err == nil {
				s.Contracts.Total++
				releases, _ := r.Search("releases").Children()
				for _, child := range releases {
					// date
					date, _ := child.Path("date").Data().(string)
					t, err := time.Parse("2006-01-02T15:04:05.000Z", date)
					if err == nil {
						if t.Before(s.FirstDate) {
							s.FirstDate = t
						}
						if t.After(s.LastDate) {
							s.LastDate = t
						}
					}

					// planning.budget.amount.amount
					amount, ok := child.Path("planning.budget.amount.amount").Data().(float64)
					if ok {
						s.Contracts.Budget += amount
					}

					// tender.status
					status, ok := child.Path("tender.status").Data().(string)
					if ok {
						switch status {
						case "active":
							s.Contracts.Active++
						case "complete":
							s.Contracts.Completed++
						}
					}

					// contracts.value.amount
					contracts, _ := child.Search("contracts").Children()
					for _, contract := range contracts {
						award, ok := contract.Path("value.amount").Data().(float64)
						if ok {
							s.Contracts.Awarded += award
						}
					}

					// tender.numberOfTenderers
					if child.ExistsP("tender.numberOfTenderers") {
						participants, _ := child.Path("tender.numberOfTenderers").Data().(float64)
						switch {
						case (participants == 1):
							s.AssignMethod.Direct.Total++
							s.AssignMethod.Direct.Budget += amount
							if status != "active" {
								s.AssignMethod.Direct.Active++
							} else {
								s.AssignMethod.Direct.Completed++
							}
							break
						case (participants >= 1 && participants <= 3):
							s.AssignMethod.Limited.Total++
							s.AssignMethod.Limited.Budget += amount
							if status != "active" {
								s.AssignMethod.Limited.Active++
							} else {
								s.AssignMethod.Limited.Completed++
							}
							break
						default:
							s.AssignMethod.Public.Total++
							s.AssignMethod.Public.Budget += amount
							if status != "active" {
								s.AssignMethod.Public.Active++
							} else {
								s.AssignMethod.Public.Completed++
							}
						}
					}
				}
			}
		}
		return nil
	})
	return s
}

func (storage *boltStorage) Count(bucket string) (count uint64) {
	storage.state.Lock()
	defer storage.state.Unlock()

	storage.db.View(func(tx *bolt.Tx) error {
		tx.WriteFlag = txWriteFlag
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			count = uint64(b.Stats().KeyN)
		}
		return nil
	})
	return count
}

func (storage *boltStorage) NextIndex(bucket string) uint64 {
	return storage.Count(bucket) + 1
}

func (storage *boltStorage) GetHead(bucket string) (key, val []byte) {
	storage.state.Lock()
	defer storage.state.Unlock()

	storage.db.View(func(tx *bolt.Tx) error {
		tx.WriteFlag = txWriteFlag
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()
		k, v := c.Last()

		// Decompress head value
		key = k
		val = storage.compressor.Decompress(v)
		return nil
	})
	return key, val
}
