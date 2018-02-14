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
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "os/signal"
  "path"
  "strings"
  
  "github.com/gorilla/websocket"
  "github.com/julienschmidt/httprouter"
  "github.com/ranveerkunal/memfs"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "github.com/transparenciamx/tsv/bot"
  "github.com/transparenciamx/tsv/notifications"
  "github.com/transparenciamx/tsv/data"
  "encoding/json"
  "strconv"
)

// serverCmd represents the serve command
var serverCmd = &cobra.Command{
  Use:     "server",
  Short:   "Starts a TSV server instance",
  Aliases: []string{"serve", "start"},
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
    useCache           bool
    serverDocs         string
    serverSSLCert      string
    serverSSLKey       string
    storageHost        string
    storageDB          string
  )
  viper.SetDefault("server.port.http", 7788)
  viper.SetDefault("server.port.https", 7789)
  viper.SetDefault("server.redirect.http", true)
  viper.SetDefault("server.cache", false)
  viper.SetDefault("server.docs", "htdocs")
  viper.SetDefault("server.cert", "")
  viper.SetDefault("server.priv", "")
  viper.SetDefault("server.storage.host", "localhost:27017")
  viper.SetDefault("server.storage.db", "tsv")
  
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
  serverCmd.Flags().BoolVar(
    &useCache,
    "use-cache",
    false,
    "Use memory cache for statis files")
  serverCmd.Flags().StringVarP(
    &serverDocs,
    "htdocs",
    "c",
    "htdocs",
    "Full path for the content to serve")
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
  serverCmd.Flags().StringVar(
    &storageHost,
    "storage-host",
    "localhost:27017",
    "MongoDB instance used as storage component")
  serverCmd.Flags().StringVar(
    &storageDB,
    "storage-db",
    "tsv",
    "MongoDB database used")
  viper.BindPFlag("server.port.http", serverCmd.Flags().Lookup("http-port"))
  viper.BindPFlag("server.port.https", serverCmd.Flags().Lookup("https-port"))
  viper.BindPFlag("server.redirect.http", serverCmd.Flags().Lookup("redirect-http"))
  viper.BindPFlag("server.cache", serverCmd.Flags().Lookup("use-cache"))
  viper.BindPFlag("server.docs", serverCmd.Flags().Lookup("htdocs"))
  viper.BindPFlag("server.store", serverCmd.Flags().Lookup("store"))
  viper.BindPFlag("server.cert", serverCmd.Flags().Lookup("cert"))
  viper.BindPFlag("server.priv", serverCmd.Flags().Lookup("priv-key"))
  viper.BindPFlag("server.storage.host", serverCmd.Flags().Lookup("storage-host"))
  viper.BindPFlag("server.storage.db", serverCmd.Flags().Lookup("storage-db"))
  RootCmd.AddCommand(serverCmd)
}

// Utility method to properly return JSON content
func sendJSON(w http.ResponseWriter, data interface{}) {
  js, _ := json.Marshal(data)
  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, "%s", js)
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

// Start server
func runServer(_ *cobra.Command, _ []string) error {
  // Get storage handler
  db, err := connectStorage(viper.GetString("server.storage.host"), viper.GetString("server.storage.db"))
  if err != nil {
    log.Printf("Storage error: %s\n", err)
    return err
  }
  defer db.Close()
  
  // Subscribe to SIGINT signals
  stopChan := make(chan os.Signal)
  signal.Notify(stopChan, os.Interrupt)
  
  // Create in-memory cache filesystem for static files
  var httpFS http.FileSystem
  if viper.GetBool("server.cache") {
    log.Println("Using memory cache")
    httpFS, err = memfs.New(viper.GetString("server.docs"))
    if err != nil {
      log.Fatalf("Cache filesystem error: %+v", err)
    }
  } else {
    httpFS = http.Dir(viper.GetString("server.docs"))
  }
  
  // Load index.html contents
  indexFile, _ = ioutil.ReadFile(path.Join(viper.GetString("server.docs"), "index.html"))
  
  // Configure router
  router := httprouter.New()
  router.NotFound = http.FileServer(httpFS)
  router.POST("/profile", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
    up := data.UserProfile{}
    err := json.Unmarshal([]byte(r.FormValue("profile")), &up)
    if err != nil {
      log.Println("Decoding error:", err)
      fmt.Fprintf(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()))
      return
    }
    log.Printf("Store user profile: %s\n", up.User)
    db.Insert("user", up)
  
    // Dispatch SMS notification
    if up.EnableSMSNotifications {
      notifications.SendSMS(&notifications.SMSOptions{
        To:      fmt.Sprintf("+52%s", up.NotificationSMS),
        Sender:  "TS2",
        Message: "Bienvenido a Testigo Social 2.0. Ahora recibirás información sobre contrataciones públicas. El dinero público también es tu dinero.",
      })
    }
  
    // Dispatch email notification
    if up.EnableEmailNotifications {
      content, _ := notifications.PrepareContent(notifications.TSVEmailTemplate, map[string]interface{}{
        "Title":   "¡Gracias por tu interés!",
        "Content": "Bienvenido a Testigo Social 2.0. A partir de ahora recibirás información oportuna sobre los procedimientos de contratación pública que te interesan. El dinero público también es tu dinero.",
      })
      notifications.SendEmail(&notifications.EmailOptions{
        To:      up.NotificationEmail,
        From:    "notificaciones@testigosocial.mx",
        Subject: "TS 2.0",
        Body:    content,
      })
    }
    sendJSON(w, map[string]interface{}{"ok": true})
  })
  router.POST("/query", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
    result := make([]interface{}, 0)
    query, err := db.Query("contracts",r.FormValue("query"))
    if err != nil {
      log.Println("Decoding error:", err)
      fmt.Fprintf(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()))
      return
    }
    
    // Apply limit value
    if r.FormValue("limit") != "" {
      l, _ := strconv.Atoi(r.FormValue("limit"))
      query = query.Limit(l)
    }
  
    // Apply sort config
    if r.FormValue("sort") != "" {
      fields := []string{}
      err := json.Unmarshal([]byte(r.FormValue("sort")), &fields)
      if err == nil {
        query = query.Sort(fields...)
      }
    }
    
    query.All(&result)
    sendJSON(w, result)
  })
  router.POST("/indicators", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
    list := make([]interface{}, 0)
    query, err := db.Query("contracts",r.FormValue("query"))
    if err != nil {
      log.Println("Decoding error:", err)
      fmt.Fprintf(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()))
      return
    }
    query.All(&list)
    sendJSON(w, data.FormatIndicatorsResult(list))
  })
  router.GET("/stats/:bucket", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
    // Get result structure
    org := data.NewOrganization()
    org.Code = ps.ByName("bucket")
    switch org.Code {
    case "gacm":
      org.Description = "la construcción del Nuevo Aeropuerto Internacional de la Ciudad de México"
    case "cdmx":
      org.Description = "la contratación de servicios y obras públicas de la Ciudad de México"
    }
  
    // Get contract list
    list := make([]interface{}, 0)
    query, _ := db.Query("contracts", fmt.Sprintf("{\"project\":\"%s\"}", org.Code))
    query.All(&list)
    org.AddRecords(list)
    sendJSON(w, org)
  })
  router.GET("/contract/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
    res, err := db.GetByID("contracts", ps.ByName("id"))
    if err != nil {
      sendJSON(w, "")
      return
    }
    sendJSON(w,res)
  })
  
  // WebSocket
  router.GET("/ws", ws)
  
  // Facebook bot
  tsvBot := bot.New(&bot.Config{
    PageToken:         os.Getenv("TSV_FB_PAGE_TOKEN"),
    VerificationToken: os.Getenv("TSV_FB_VERIFY_TOKEN"),
  })
  router.GET("/fb", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    tsvBot.Verify(w, r)
  })
  router.POST("/fb", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    tsvBot.ReceiveMessage(w, r)
  })
  
  // Redirect ReactRouter paths to the 'index.html' file
  router.GET("/informacion", serveIndex)
  router.GET("/contratos", serveIndex)
  router.GET("/contratos/:id", serveIndex)
  router.GET("/indicadores", serveIndex)
  router.GET("/registro", serveIndex)
  
  // Handle incoming FB messages
  go func() {
    for {
      select {
      case msg := <- tsvBot.IncomingMessages:
        // Handle postbacks
        if pb := msg.Entry[0].Messaging[0].Postback.Payload; pb != "" {
          log.Printf("Handling postback: %s", pb)
          user := msg.Entry[0].Messaging[0].Sender.ID
          switch pb {
          case "NEW_QUERY":
            tsvBot.SendMessage(GetQueryMenu(user))
          case "NEW_THREAD":
            tsvBot.SendMessage(GetWelcomeMessage(user))
          }
        }
        
        // Handle quick replies
        if qr := msg.Entry[0].Messaging[0].Message.QuickReply.Payload; qr != "" {
          log.Printf("Handling quick reply: %s", qr)
          user := msg.Entry[0].Messaging[0].Sender.ID
          switch qr {
          case "RECENT":
            result := make([]interface{}, 0)
            query, _ := db.Query("contracts","{}")
            query = query.Limit(4).Sort("-releases.date")
            query.All(&result)
            tsvBot.SendMessage(GetContractListMessage(user, result))
          case "AMOUNT":
            result := make([]interface{}, 0)
            query, _ := db.Query("contracts","{}")
            query = query.Limit(4).Sort("-releases.planning.budget.amount.amount")
            query.All(&result)
            tsvBot.SendMessage(GetContractListMessage(user, result))
          case "GACM":
            result := make([]interface{}, 0)
            query, _ := db.Query("contracts","{\"project\":\"gacm\"}")
            query = query.Limit(4).Sort("-releases.planning.budget.amount.amount")
            query.All(&result)
            tsvBot.SendMessage(GetContractListMessage(user, result))
          case "CDMX":
            result := make([]interface{}, 0)
            query, _ := db.Query("contracts","{\"project\":\"cdmx\"}")
            query = query.Limit(4).Sort("-releases.planning.budget.amount.amount")
            query.All(&result)
            tsvBot.SendMessage(GetContractListMessage(user, result))
          }
        }
      case err := <- tsvBot.Errors:
        log.Printf("error: %+v", err)
      }
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
