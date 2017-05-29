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
  "github.com/transparenciamx/tsv/storage"
  "github.com/transparenciamx/tsv/bot"
  "encoding/json"
  "github.com/Jeffail/gabs"
  "fmt"
)

// Utility method to open a storage handler
func connectStorage(host, database string) (*storage.Handler, error ) {
  // Use linked storage if available
  linked := os.Getenv("STORAGE_PORT")
  if linked != "" {
    host = linked[6:]
  }
  
  // Connect to storage instance
  db, err := storage.NewHandler(host, database)
  if err != nil {
    return nil, err
  }
  log.Printf("Using storage: %s/%s", host, database)
  return db, nil
}

// GetWelcomeMessage returns a ready to use 'welcome message' for a new user
func GetWelcomeMessage(user string) *bot.Message {
  return &bot.Message{
    Recipient: bot.Recipient{
      ID: user,
    },
    Message: bot.MessageBody{
      Attachment: bot.Attachment{
        Type: "template",
        Payload: bot.Payload{
          TemplateType: "button",
          Text: "Elije una de las siguientes opciones:",
          Buttons: []bot.URLButton{
            {
              Type: "postback",
              Title: "Nueva Consulta",
              Payload: "NEW_QUERY",
            },
            {
              Type: "web_url",
              URL: "https://www.testigosocial.mx",
              Title: "Visitar Sitio Web",
            },
          },
        },
      },
    },
  }
}

// GetQueryMenu returns a menu with the available basic queries
func GetQueryMenu(user string) *bot.Message {
  return &bot.Message{
    Recipient: bot.Recipient{
      ID: user,
    },
    Message: bot.MessageBody{
      Text: "Seleccione el criterio para su consulta de los procesos de contratación registrados:",
      QuickReplies: []bot.QuickReply{
        {ContentType: "text", Title: "Recientes", Payload: "RECENT"},
        {ContentType: "text", Title: "Mayor Monto", Payload: "AMOUNT"},
        {ContentType: "text", Title: "Grupo Aeroportuario", Payload: "GACM"},
        {ContentType: "text", Title: "Ciudad de México", Payload: "CDMX"},
      },
    },
  }
}

// GetContractListMessage return a list message based on the provided query results
func GetContractListMessage(user string, list []interface{}) *bot.Message {
  els := []bot.ListElement{}
  for _, rec := range list {
    json, _ := json.Marshal(rec)
    r, _ := gabs.ParseJSON(json)
    releases, _ := r.Search("releases").Children()
    els = append(els, bot.ListElement{
      Title: releases[0].Path("ocid").String(),
      Subtitle: releases[0].Path("tender.title").String(),
      DefaultAction: bot.URLButton{
        Type: "web_url",
        URL: fmt.Sprintf("https://www.testigosocial.mx/contratos/%s", r.Path("_id").Data().(string)),
      },
    })
  }
  
  return &bot.Message{
    Recipient: bot.Recipient{
      ID: user,
    },
    Message: bot.MessageBody{
      Attachment: bot.Attachment{
        Type: "template",
        Payload: bot.Payload{
          TemplateType: "list",
          TopElementStyle: "compact",
          Elements: els,
        },
      },
    },
  }
}
