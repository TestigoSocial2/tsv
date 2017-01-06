package data

import (
	"os"

	"github.com/bcessa/tsv/storage"
)

// OpenStorage is an utility method to get a storage interface
func OpenStorage() (storage.Provider, error) {
	conf := storage.DefaultConfig()
	if os.Getenv("TSV_STORAGE") != "" {
		conf.Path = os.Getenv("TSV_STORAGE")
	}
	return storage.New(conf)
}
