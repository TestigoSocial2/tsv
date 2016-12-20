// Copyright © 2016 Transparencia Mexicana AC. <ben@datos.mx>
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

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/bcessa/tsv/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:     "load",
	Short:   "Load contract information",
	Aliases: []string{"store"},
	RunE:    runLoadCmd,
}

func init() {
	var (
		loadPath   string
		loadBucket string
		loadStore  string
	)
	viper.SetDefault("load.path", ".")
	viper.SetDefault("load.bucket", "main")
	viper.SetDefault("load.store", os.TempDir())
	loadCmd.Flags().StringVarP(
		&loadPath,
		"path",
		"p",
		".",
		"Full path for the contents to load, can be a file or a directory to scan")
	loadCmd.Flags().StringVarP(
		&loadBucket,
		"bucket",
		"b",
		"main",
		"Bucket to use for information storage")
	loadCmd.Flags().StringVarP(
		&loadStore,
		"store",
		"s",
		os.TempDir(),
		"Full path to use as data store location")
	viper.BindPFlag("load.path", loadCmd.Flags().Lookup("path"))
	viper.BindPFlag("load.bucket", loadCmd.Flags().Lookup("bucket"))
	viper.BindPFlag("load.store", loadCmd.Flags().Lookup("store"))
	RootCmd.AddCommand(loadCmd)
}

// Command execution
func runLoadCmd(cmd *cobra.Command, args []string) error {
	// Prepare storage provider
	conf := storage.DefaultConfig()
	conf.Path = path.Join(viper.GetString("load.store"), "tsv.db")
	store, err := storage.New(conf)
	if err != nil {
		return err
	}
	defer store.Close()

	src, err := os.Stat(viper.GetString("load.path"))
	if err != nil {
		return err
	}
	mode := src.Mode()

	// Iterate directory
	if mode.IsDir() {
		files, err := ioutil.ReadDir(viper.GetString("load.path"))
		if err != nil {
			return err
		}

		// Iterate files and store
		counter := 0
		bucket := viper.GetString("load.bucket")
		for _, f := range files {
			name := f.Name()
			if strings.Contains(name, ".json") {
				b, err := ioutil.ReadFile(path.Join(viper.GetString("load.path"), name))
				if err == nil {
					store.Write(bucket, []byte(strings.Split(name, ".json")[0]), b)
					counter++
				}
			}
		}
		fmt.Printf("Files stored: %d in bucket: %s\n", counter, bucket)
		return nil
	}

	// Use regular file
	if mode.IsRegular() {
		name := src.Name()
		if !strings.Contains(name, ".json") {
			return fmt.Errorf("Invalid file format, expecting a '.json' file")
		}
		b, err := ioutil.ReadFile(viper.GetString("load.path"))
		if err != nil {
			return err
		}
		bucket := viper.GetString("load.bucket")
		store.Write(bucket, []byte(strings.Split(name, ".json")[0]), b)
		fmt.Printf("File stored: %s in bucket: %s\n", name, bucket)
	}
	return nil
}
