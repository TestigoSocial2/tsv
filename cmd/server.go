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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/bcessa/tsv/storage"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/mxabierto/go-twilio"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type contracts struct {
	Budget    float64 `json:"budget"`
	Awarded   float64 `json:"awarded"`
	Total     int     `json:"total"`
	Active    int     `json:"active"`
	Completed int     `json:"completed"`
}

// Stats ...
type bucketStats struct {
	FirstDate    time.Time `json:"firstDate"`
	LastDate     time.Time `json:"lastDate"`
	Contracts    contracts `json:"contracts"`
	AssignMethod struct {
		Direct  contracts `json:"direct"`
		Limited contracts `json:"limited"`
		Public  contracts `json:"public"`
	} `json:"method"`
}

// serverCmd represents the serve command
var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "Starts a TSV server instance",
	Aliases: []string{"serve"},
	RunE:    runServer,
}

// WebSocket factory, adjust an incoming HTTP connection to include WS support
var wsu = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Setup
func init() {
	var (
		serverPort  int
		serverDocs  string
		serverStore string
	)
	viper.SetDefault("server.port", 7788)
	viper.SetDefault("server.docs", "htdocs")
	viper.SetDefault("server.store", os.TempDir())
	serverCmd.Flags().IntVarP(
		&serverPort,
		"port",
		"p",
		7788,
		"TCP port to use for server/client communications")
	serverCmd.Flags().StringVarP(
		&serverDocs,
		"htdocs",
		"c",
		"htdocs",
		"Full path for the content to serve")
	serverCmd.Flags().StringVarP(
		&serverStore,
		"store",
		"s",
		os.TempDir(),
		"Full path to use as data store location")
	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.docs", serverCmd.Flags().Lookup("htdocs"))
	viper.BindPFlag("server.store", serverCmd.Flags().Lookup("store"))
	RootCmd.AddCommand(serverCmd)
}

// Utility method to get a storage interface
func openStorage() (storage.Provider, error) {
	conf := storage.DefaultConfig()
	conf.Path = path.Join(viper.GetString("server.store"), "tsv.db")
	return storage.New(conf)
}

// Basic user profile structure
type userProfile struct {
	User                     string   `json:"user"`
	Password                 string   `json:"password"`
	UserType                 string   `json:"userType"`
	Age                      string   `json:"age"`
	PostalCode               string   `json:"postalCode"`
	SelectedAgencies         []string `json:"selectedAgencies"`
	SelectedProjects         []string `json:"selectedProjects"`
	NotificationEmail        string   `json:"notificationEmail"`
	EnableEmailNotifications bool     `json:"enableEmailNotifications"`
	NotificationSMS          string   `json:"notificationSMS"`
	EnableSMSNotifications   bool     `json:"enableSMSNotifications"`
}

// Create or update a user profile
func profile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Open storage
	store, err := openStorage()
	if err != nil {
		log.Println("Storage error:", err)
		fmt.Fprintf(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()))
		return
	}
	defer store.Close()

	up := userProfile{}
	err = json.Unmarshal([]byte(r.FormValue("profile")), &up)
	if err != nil {
		log.Println("Decoding error:", err)
		fmt.Fprintf(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()))
		return
	}
	log.Printf("Store user profile: %s\n", up.User)
	store.Write("profile", []byte(up.User), []byte(r.FormValue("profile")))

	tw := twilio.New("AC07ca38e0d96366ab706981b6b661a3ff", "fecb164a077922442f4f67cc797e0d4c")
	_, err = tw.SendMessage(&twilio.MessageOptions{
		To:   fmt.Sprintf("+521%s", up.NotificationSMS),
		From: "+14242964188",
		Body: "Bienvenido a Testigo Social Virtual 2.0 a partir de este momento comenzaras a recibir notificaciones relevantes sobre los procesos de contratación pública de tu interes.",
	})
	if err != nil {
		log.Println("SMS error:", err)
	}
	fmt.Fprintf(w, fmt.Sprintf("{\"ok\":true}"))
}

// Calculate and return stats about the contracts
func stats(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Open storage
	store, err := openStorage()
	if err != nil {
		log.Println("Storage error:", err)
		fmt.Fprintf(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()))
		return
	}
	defer store.Close()

	log.Printf("Calculating stats for: %s\n", ps.ByName("bucket"))
	// stats := store.GetStats(ps.ByName("bucket"))
	// res, _ := json.Marshal(stats)
	// fmt.Fprintf(w, string(res))

	data := &bucketStats{}
	data.FirstDate = time.Now()
	data.LastDate = time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	c := make(chan *storage.Record)
	go store.Cursor(ps.ByName("bucket"), c)
loop:
	for {
		select {
		case rec, ok := <-c:
			if !ok {
				break loop
			}
			val := rec.Value
			r, err := gabs.ParseJSON(val)
			if err == nil {
				data.Contracts.Total++
				releases, _ := r.Search("releases").Children()
				for _, child := range releases {
					// date
					date, _ := child.Path("date").Data().(string)
					t, err := time.Parse("2006-01-02T15:04:05.000Z", date)
					if err == nil {
						if t.Before(data.FirstDate) {
							data.FirstDate = t
						}
						if t.After(data.LastDate) {
							data.LastDate = t
						}
					}

					// planning.budget.amount.amount
					amount, ok := child.Path("planning.budget.amount.amount").Data().(float64)
					if ok {
						data.Contracts.Budget += amount
					}

					// tender.status
					status, ok := child.Path("tender.status").Data().(string)
					if ok {
						switch status {
						case "active":
							data.Contracts.Active++
						case "complete":
							data.Contracts.Completed++
						}
					}

					// contracts.value.amount
					contracts, _ := child.Search("contracts").Children()
					for _, contract := range contracts {
						award, ok := contract.Path("value.amount").Data().(float64)
						if ok {
							data.Contracts.Awarded += award
						}
					}

					// tender.numberOfTenderers
					if child.ExistsP("tender.numberOfTenderers") {
						participants, _ := child.Path("tender.numberOfTenderers").Data().(float64)
						switch {
						case (participants == 1):
							data.AssignMethod.Direct.Total++
							data.AssignMethod.Direct.Budget += amount
							if status != "active" {
								data.AssignMethod.Direct.Active++
							} else {
								data.AssignMethod.Direct.Completed++
							}
							break
						case (participants >= 1 && participants <= 3):
							data.AssignMethod.Limited.Total++
							data.AssignMethod.Limited.Budget += amount
							if status != "active" {
								data.AssignMethod.Limited.Active++
							} else {
								data.AssignMethod.Limited.Completed++
							}
							break
						default:
							data.AssignMethod.Public.Total++
							data.AssignMethod.Public.Budget += amount
							if status != "active" {
								data.AssignMethod.Public.Active++
							} else {
								data.AssignMethod.Public.Completed++
							}
						}
					}
				}
			}
		}
	}

	res, _ := json.Marshal(data)
	fmt.Fprintf(w, string(res))
}

// Receives a query request and runs it
// func query(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

// Handle websockets
// func ws(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	c, err := wsu.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("ws:", err)
// 		return
// 	}
// 	defer c.Close()
// 	messageType, p, err := c.ReadMessage()
// 	c.WriteMessage(messageType, p)
// }

func runServer(cmd *cobra.Command, args []string) error {
	// Configure router
	router := httprouter.New()
	router.NotFound = http.FileServer(http.Dir(viper.GetString("server.docs")))
	router.POST("/profile", profile)
	router.GET("/stats/:bucket", stats)
	// router.GET("/ws", ws)

	log.Println("Handling requests on port:", viper.GetInt("server.port"))
	err := http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("server.port")), router)
	if err != nil {
		return err
	}
	return nil
}
