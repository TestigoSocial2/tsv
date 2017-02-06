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
	"fmt"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/transparenciamx/tsv/compressor"
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

func (storage *boltStorage) Cursor(bucket string, out chan<- *Record, cancel chan bool) {
	storage.db.View(func(tx *bolt.Tx) error {
		tx.WriteFlag = txWriteFlag
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			c := b.Cursor()
			k, v := c.First()
			for {
				select {
				case <-cancel:
					return nil
				default:
					if k == nil {
						return nil
					}
					out <- &Record{
						Key:   k,
						Value: storage.compressor.Decompress(v),
					}
					k, v = c.Next()
				}
			}
		}
		return nil
	})
	close(out)
	return
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
