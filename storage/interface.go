package storage

import (
	"errors"
	"time"

	"github.com/bcessa/tsv/compressor"
)

// Record represents a basic record definition
type Record struct {
	Key   []byte
	Value []byte
}

type contracts struct {
	Budget    float64 `json:"budget"`
	Awarded   float64 `json:"awarded"`
	Total     int     `json:"total"`
	Active    int     `json:"active"`
	Completed int     `json:"completed"`
}

// Stats ...
type Stats struct {
	FirstDate    time.Time `json:firstDate`
	LastDate     time.Time `json:lastDate`
	Contracts    contracts `json:"contracts"`
	AssignMethod struct {
		Direct  contracts `json:"direct"`
		Limited contracts `json:"limited"`
		Public  contracts `json:"public"`
	} `json:"method"`
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

	// Return a BoltDB cursor; temporary fix
	Cursor(bucket string, s chan *Record)

	// GetStats ...
	GetStats(bucket string) (s *Stats)
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
