package storage

import (
	"errors"

	"github.com/bcessa/tsv/compressor"
)

// Record represents a basic record definition
type Record struct {
	Key   []byte
	Value []byte
}

// Provider defines the common interface for the underlaying
// data store component
type Provider interface {
	// Open will start a connection to the data store
	Open(path string) error

	// Close terminate the underlaying connection to the data store
	Close() error

	// Write will store a new entry in the data store
	Write(bucket string, key, val []byte) error

	// Read attempt to retrieve an existing record from the data store
	Read(bucket string, key []byte) []byte

	// Count return the number of elements in the chain
	Count(bucket string) uint64

	// NextIndex return a valid auto increment index for elements in the chain
	NextIndex(bucket string) uint64

	// GetHead return the top-most key/value pair on the chain
	GetHead(bucket string) (key, val []byte)

	// Iterate the contents of a given bucket
	Cursor(bucket string, out chan<- *Record, cancel chan bool)
}

// New is a constructor/setup method to easily create a storage instance
func New(opts *Config) (storage Provider, err error) {
	switch opts.Type {
	case "bolt":
		storage = &boltStorage{
			compressor: compressor.New(opts.Compressor),
		}
		err := storage.Open(opts.Path)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Invalid storage type")
	}
	return storage, nil
}
