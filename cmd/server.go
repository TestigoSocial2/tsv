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
	"os/signal"
	"path"

	"github.com/bcessa/tsv/data"
	"github.com/bcessa/tsv/notifications"
	"github.com/bcessa/tsv/storage"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

// Create or update a user profile
func profile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Open storage
	store, err := data.OpenStorage()
	if err != nil {
		log.Println("Storage error:", err)
		fmt.Fprintf(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()))
		return
	}
	defer store.Close()

	up := data.UserProfile{}
	err = json.Unmarshal([]byte(r.FormValue("profile")), &up)
	if err != nil {
		log.Println("Decoding error:", err)
		fmt.Fprintf(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()))
		return
	}
	log.Printf("Store user profile: %s\n", up.User)
	store.Write("profile", []byte(up.User), []byte(r.FormValue("profile")))

	// Dispatch SMS notification
	if up.EnableSMSNotifications {
		smsID, err := notifications.SendSMS(&notifications.SMSOptions{
			To:      fmt.Sprintf("+52%s", up.NotificationSMS),
			Sender:  "TS2",
			Message: "Bienvenido a Testigo Social Virtual 2.0 a partir de este momento comenzaras a recibir notificaciones relevantes",
		})
		if err != nil {
			log.Println("SMS error:", err)
		}
		log.Println("SMS notifications dipatched with ID:", smsID)
	}

	// Dispatch email notification
	if up.EnableEmailNotifications {
		content, _ := notifications.PrepareContent(notifications.TSVEmailTemplate, map[string]interface{}{
			"Title":   "¡Gracias por tu interés!",
			"Content": "Bienvenido a Testigo Social Virtual 2.0 a partir de este momento comenzaras a recibir notificaciones relevantes sobre los procesos de contratación pública de tu interes.",
		})
		emailID, err := notifications.SendEmail(&notifications.EmailOptions{
			To:      up.NotificationEmail,
			From:    "notificaciones@testigosocial.mx",
			Subject: "TS 2.0",
			Body:    content,
		})
		if err != nil {
			log.Println("Email error:", err)
		}
		log.Println("Email notifications dipatched with ID:", emailID)
	}

	fmt.Fprintf(w, fmt.Sprintf("{\"ok\":true}"))
}

// Calculate and return stats about the contracts on a specific storage bucket
func stats(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Open storage
	store, err := data.OpenStorage()
	if err != nil {
		log.Println("Storage error:", err)
		fmt.Fprintf(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()))
		return
	}
	defer store.Close()

	log.Printf("Calculating stats for: %s\n", ps.ByName("bucket"))
	org := data.NewOrganization()
	org.Code = ps.ByName("bucket")
	switch org.Code {
	case "gacm":
		org.Description = "Nuevo Aeropuerto Internacional de la Ciudad de México"
	case "cdmx":
		org.Description = "Contratación de Servicios y Obras de la Ciudad de México"
	}

	cursor := make(chan *storage.Record)
	go store.Cursor(ps.ByName("bucket"), cursor, nil)
	for rec := range cursor {
		org.AddRecord(rec.Value)
	}

	res, _ := json.Marshal(org)
	fmt.Fprintf(w, string(res))
}

// Contracts list query
func query(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	query := data.Query{}
	err := json.Unmarshal([]byte(r.FormValue("query")), &query)
	if err != nil {
		log.Println("Decoding error:", err)
		fmt.Fprintf(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()))
		return
	}
	query.Bucket = ps.ByName("bucket")

	log.Printf("Running query: %+v\n", query)
	res, _ := json.Marshal(query.Run())
	fmt.Fprintf(w, fmt.Sprintf("%s", res))
}

// Statistics query
func indicators(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	query := data.IndicatorsQuery{}
	err := json.Unmarshal([]byte(r.FormValue("query")), &query)
	if err != nil {
		log.Println("Decoding error:", err)
		fmt.Fprintf(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()))
		return
	}

	log.Printf("Running indicators query: %+v\n", query)
	res, _ := json.Marshal(query.Run())
	fmt.Fprintf(w, fmt.Sprintf("%s", res))
}

// Store pre-register email
func preregister(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Printf("Pre-register for: %+v\n", r.FormValue("email"))
	f, err := os.OpenFile("/data/pre-register", os.O_APPEND|os.O_WRONLY, 0600)
	if err == nil {
		f.WriteString(fmt.Sprintf("%s\n", r.FormValue("email")))
	}
	defer f.Close()
	fmt.Fprintf(w, fmt.Sprintf("%s", "ok"))
}

// Handle websockets
func ws(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c, err := wsu.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ws:", err)
		return
	}
	defer c.Close()
	messageType, p, err := c.ReadMessage()

	if err != nil {
		// TODO - Inspect and handle message type and content
		c.WriteMessage(messageType, p)
	}
}

func runServer(cmd *cobra.Command, args []string) error {
	// Set storage config variable
	os.Setenv("TSV_STORAGE", path.Join(viper.GetString("server.store"), "tsv.db"))

	// Subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	// Configure router
	router := httprouter.New()
	router.NotFound = http.FileServer(http.Dir(viper.GetString("server.docs")))
	router.POST("/profile", profile)
	router.POST("/query/:bucket", query)
	router.POST("/indicators", indicators)
	router.POST("/preregister", preregister)
	router.GET("/stats/:bucket", stats)
	router.GET("/ws", ws)

	// Start server
	log.Println("Handling requests on port:", viper.GetInt("server.port"))
	go func() {
		addr := fmt.Sprintf(":%d", viper.GetInt("server.port"))
		if err := http.ListenAndServe(addr, router); err != nil {
			log.Printf("Connection error: %s\n", err)
		}
	}()

	// Wait for SIGINT
	<-stopChan
	log.Println("Shutting down server")
	return nil
}
