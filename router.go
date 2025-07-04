// router.go
package main

import(
   "os"
   "net/http"
   "github.com/asccclass/sherryserver"
)

func NewRouter(srv *SherryServer.Server, documentRoot string)(*http.ServeMux) {
   router := http.NewServeMux()

   // Static File server
   staticfileserver := SherryServer.StaticFileServer{documentRoot, "index.html"}
   staticfileserver.AddRouter(router)

   if os.Getenv("OllamaUrl") != "" { // AI Chat
      ollama := NewOllamaClient()
      if ollama == nil { 
         return nil
      }
      AIs["Ollama"] = ollama
      ollama.AddRouter(router)
   }

   // Input App router
   router.Handle("/new-chat", http.HandlerFunc(handleNewChat))
   router.Handle("/sse", http.HandlerFunc(SSEChat))
   router.Handle("/send-message", http.HandlerFunc(handleSendMessage))
   return router
}
