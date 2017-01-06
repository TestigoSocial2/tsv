package storage

import (
	"fmt"
	"sync"
	"time"

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

func (storage *boltStorage) Cursor(bucket string, ch chan *Record) {
	storage.db.View(func(tx *bolt.Tx) error {
		tx.WriteFlag = txWriteFlag
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				ch <- &Record{
					Key:   k,
					Value: storage.compressor.Decompress(v),
				}
			}
		}
		return nil
	})
	close(ch)
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
