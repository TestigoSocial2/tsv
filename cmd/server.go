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
  "crypto/tls"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "os/signal"
  "path"
  "path/filepath"
  "strings"
  
  "github.com/gorilla/websocket"
  "github.com/julienschmidt/httprouter"
  "github.com/ranveerkunal/memfs"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "github.com/transparenciamx/tsv/bot"
  "github.com/transparenciamx/tsv/data"
  "github.com/transparenciamx/tsv/notifications"
  "github.com/transparenciamx/tsv/storage"
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

// Hold index file content
var indexFile []byte

// Setup
func init() {
  var (
    serverPortHTTP     int
    serverPortHTTPS    int
    serverRedirectHTTP bool
    serverDocs         string
    serverStore        string
    serverSSLCert      string
    serverSSLKey       string
  )
  viper.SetDefault("server.port.http", 7788)
  viper.SetDefault("server.port.https", 7789)
  viper.SetDefault("server.redirect.http", true)
  viper.SetDefault("server.docs", "htdocs")
  viper.SetDefault("server.store", "htdocs")
  viper.SetDefault("server.cert", "")
  viper.SetDefault("server.priv", "")
  
  serverCmd.Flags().IntVar(
    &serverPortHTTP,
    "http-port",
    7788,
    "HTTP port to use for server/client communications")
  serverCmd.Flags().IntVar(
    &serverPortHTTPS,
    "https-port",
    7789,
    "HTTPS port to use for secure server/client communications")
  serverCmd.Flags().BoolVar(
    &serverRedirectHTTP,
    "redirect-http",
    true,
    "Redirect all HTTP traffic if a secure channel is available")
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
    "htdocs",
    "Full path to use as data store location")
  serverCmd.Flags().StringVar(
    &serverSSLCert,
    "cert",
    "",
    "SSL certificate to enable HTTPS")
  serverCmd.Flags().StringVar(
    &serverSSLKey,
    "priv-key",
    "",
    "SSL certificate's private key")
  viper.BindPFlag("server.port.http", serverCmd.Flags().Lookup("http-port"))
  viper.BindPFlag("server.port.https", serverCmd.Flags().Lookup("https-port"))
  viper.BindPFlag("server.redirect.http", serverCmd.Flags().Lookup("redirect-http"))
  viper.BindPFlag("server.docs", serverCmd.Flags().Lookup("htdocs"))
  viper.BindPFlag("server.store", serverCmd.Flags().Lookup("store"))
  viper.BindPFlag("server.cert", serverCmd.Flags().Lookup("cert"))
  viper.BindPFlag("server.priv", serverCmd.Flags().Lookup("priv-key"))
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
      Message: "Bienvenido a Testigo Social 2.0. Ahora recibirás información sobre contrataciones públicas. El dinero público también es tu dinero.",
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
      "Content": "Bienvenido a Testigo Social 2.0. A partir de ahora recibirás información oportuna sobre los procedimientos de contratación pública que te interesan. El dinero público también es tu dinero.",
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
    org.Description = "la construcción del Nuevo Aeropuerto Internacional de la Ciudad de México"
  case "cdmx":
    org.Description = "la contratación de servicios y obras públicas de la Ciudad de México"
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

// Handle websockets
func ws(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  c, err := wsu.Upgrade(w, r, nil)
  if err != nil {
    log.Println("ws:", err)
    return
  }
  
  go func(c *websocket.Conn) {
    for {
      mType, data, err := c.ReadMessage()
      if err != nil {
        // Close connection upon request or error
        c.Close()
        return
      }
      
      // TODO - Handle messages
      log.Printf("message: %+v - %+v", mType, data)
      c.WriteMessage(mType, data)
    }
  }(c)
}

// Redirect HTTP to HTTPS
func redirect(w http.ResponseWriter, req *http.Request) {
  host := strings.Split(req.Host, ":")[0]
  target := "https://" + host + req.URL.Path
  if len(req.URL.RawQuery) > 0 {
    target += "?" + req.URL.RawQuery
  }
  log.Printf("redirect to: %s", target)
  http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

// Redirect ReactRouter paths to the 'index.hmtl' file
func serveIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, fmt.Sprintf("%s", indexFile))
}

func runServer(_ *cobra.Command, _ []string) error {
  // Set storage config variable
  sp, err := filepath.Abs(path.Join(viper.GetString("server.store"), "tsv.db"))
  if err != nil {
    return err
  }
  os.Setenv("TSV_STORAGE", sp)
  
  // Subscribe to SIGINT signals
  stopChan := make(chan os.Signal)
  signal.Notify(stopChan, os.Interrupt)
  
  // Create in-memory cache filesystem for static files
  cacheFS, err := memfs.New(viper.GetString("server.docs"))
  if err != nil {
    log.Fatalf("Cache filesystem error: %+v", err)
  }
  
  // Start facebook bot
  tsvBot := bot.New(&bot.Config{
    VerificationToken: os.Getenv("TSV_FB_VERIFY_TOKEN"),
    PageToken:         os.Getenv("TSV_FB_PAGE_TOKEN"),
  })
  
  // Load index.html contents
  indexFile, _ = ioutil.ReadFile(path.Join(viper.GetString("server.docs"), "index.html"))
  
  // Configure router
  router := httprouter.New()
  router.NotFound = http.FileServer(cacheFS)
  router.POST("/profile", profile)
  router.POST("/query/:bucket", query)
  router.POST("/indicators", indicators)
  router.GET("/stats/:bucket", stats)
  router.GET("/ws", ws)
  router.GET("/fb", tsvBot.Verify)
  router.POST("/fb", tsvBot.ReceiveMessage)
  
  // Redirect ReactRouter paths to the 'index.hmtl' file
  router.GET("/informacion", serveIndex)
  router.GET("/contratos", serveIndex)
  router.GET("/indicadores", serveIndex)
  router.GET("/registro", serveIndex)
  
  // Handle incoming FB messages
  go func() {
    for msg := range tsvBot.IncomingMessages {
      log.Printf("%+v", msg)
      tsvBot.DispatchMessage(msg.User, "Tu mensaje es muy importante para nosotros, en breve nos comunicaremos contigo!")
    }
  }()
  
  // Start server
  go func() {
    cert := viper.GetString("server.cert")
    priv := viper.GetString("server.priv")
    addr := fmt.Sprintf(":%d", viper.GetInt("server.port.http"))
    if cert != "" && priv != "" {
      // Secure HTTPS server
      secureAddr := fmt.Sprintf(":%d", viper.GetInt("server.port.https"))
      srv := &http.Server{
        Addr:    secureAddr,
        Handler: router,
        TLSConfig: &tls.Config{
          MinVersion:               tls.VersionTLS12,
          CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
          PreferServerCipherSuites: true,
          CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
            tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
          },
        },
        TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
      }
      
      // Start a basic HTTP server that will redirect all requests to the HTTPS location,
      // if required
      if viper.GetBool("server.redirect.http") {
        go http.ListenAndServe(addr, http.HandlerFunc(redirect))
      }
      
      log.Println("Handling secure requests on:", secureAddr)
      if err := srv.ListenAndServeTLS(cert, priv); err != nil {
        log.Printf("HTTPS server error: %s\n", err)
      }
    } else {
      // Standard HTTP server
      log.Println("Handling requests on:", addr)
      if err := http.ListenAndServe(addr, router); err != nil {
        log.Printf("HTTP server error: %s\n", err)
      }
    }
  }()
  
  // Wait for SIGINT
  <-stopChan
  log.Println("Shutting down server")
  return nil
}
