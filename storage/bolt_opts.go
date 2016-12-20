// +build !linux

package storage

// txWriteFlag specify the value for Tx.WriteFlag, useful when
// working with larger-than-RAM datasets
// https://godoc.org/github.com/boltdb/bolt#Tx
const txWriteFlag = 0
