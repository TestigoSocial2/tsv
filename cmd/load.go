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

package cmd

import (
  "os"
  
  "log"
  
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "io/ioutil"
  "strings"
  "fmt"
  "encoding/json"
  "path"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
  Use:     "load",
  Short:   "Load contract information",
  Aliases: []string{"add", "save", "store"},
  RunE:    runLoadCmd,
}

func init() {
  var (
    loadPath    string
    storageHost string
    storageDB   string
  )
  
  loadCmd.Flags().StringVar(
    &loadPath,
    "path",
    "",
    "Path to process")
  viper.SetDefault("load.path", "")
  viper.BindPFlag("load.path", loadCmd.Flags().Lookup("path"))
  
  loadCmd.Flags().StringVar(
    &storageHost,
    "storage-host",
    "localhost:27017",
    "MongoDB instance used as storage component")
  viper.SetDefault("load.storage.host", "localhost:27017")
  viper.BindPFlag("load.storage.host", loadCmd.Flags().Lookup("storage-host"))
  
  loadCmd.Flags().StringVar(&storageDB,
    "storage-db",
    "tsv",
    "MongoDB database used")
  viper.SetDefault("load.storage.db", "tsv")
  viper.BindPFlag("load.storage.db", loadCmd.Flags().Lookup("storage-db"))
  
  loadCmd.Flags().StringVar(&storageDB,
    "project-id",
    "",
    "Project identifier")
  viper.SetDefault("load.project.id", "")
  viper.BindPFlag("load.project.id", loadCmd.Flags().Lookup("project-id"))
  
  RootCmd.AddCommand(loadCmd)
}

// Command execution
func runLoadCmd(cmd *cobra.Command, args []string) error {
  // Get storage handler
  db, err := connectStorage(viper.GetString("load.storage.host"), viper.GetString("load.storage.db"))
  if err != nil {
    log.Printf("Storage error: %s\n", err)
    return err
  }
  defer db.Close()
  
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
    var item map[string]interface{}
    pid := viper.GetString("load.project.id")
    counter := 0
    for _, f := range files {
      name := f.Name()
      if strings.Contains(name, ".json") {
        b, err := ioutil.ReadFile(path.Join(viper.GetString("load.path"), name))
        if err == nil {
          err = json.Unmarshal(b, &item)
          if err == nil {
            if pid != "" {
              item["project"] = pid
            }
            db.Insert("contracts", item)
            counter++
          }
        }
      }
    }
    fmt.Printf("Files stored: %d\n", counter)
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
    var item map[string]interface{}
    err = json.Unmarshal(b, &item)
    if err != nil {
      return nil
    }
    
    // Attach project id, if provided
    if pid := viper.GetString("load.project.id"); pid != "" {
      item["project"] = pid
    }
    db.Insert("contracts", item)
  }
  return nil
}
