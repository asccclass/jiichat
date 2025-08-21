package main

import(
   "os"
   "fmt"
   // "time"
   // "strings"
   "net/http"
   // "encoding/json"
)

// 消息結構
type ChatMessage struct {
   Type    string `json:"type"`
   Content string `json:"content"`
}

// AI回應（TODO: 需要增加記憶體）
func AIResponse(model, userMessage string)(string) /*([]ChatMessage)*/ {
   response := ""
   var err error
   if os.Getenv("OllamaUrl") != "" {   // 根據不同模型準備不同的回應
      ola := AIs["Ollama"].(*OllamaClient)
      response, err = ola.Ask(model, userMessage, "", nil) // response 為純文字或json字串
      if err != nil {
         response = "抱歉！無法處理您的請求。(" + err.Error() + ")"
      }
   }
   // return ResponseChunks(response)
   return response
}

// 接收 SSE 請求並處理聊天  /sse
func SSEChat(w http.ResponseWriter, r *http.Request) { 
   message := r.URL.Query().Get("message")
   model := r.URL.Query().Get("model")
   if message == "" {
      fmt.Println("使用者傳送空白資訊")
      return
   }
   res :=  AIResponse(model, message)
   if len(res) > 0 {
      ResponseChunks(w, res)
      fmt.Println("send finish")
      // Response2User(w, res)  // 模擬AI回應模式 
   } else {
      ResponseChunks(w, "沒收到資料")
   }
}

// 處理新對話請求
func handleNewChat(w http.ResponseWriter, r *http.Request) {
   // 返回空的對話區域以重置對話
   w.Header().Set("Content-Type", "text/html")
   fmt.Fprint(w, `<div class="message-container">
      <div class="message ai-message"><div class="message-content">您好！我已經為您開始了一個新對話。我能幫您什麼忙嗎？</div></div>
   </div>`)
}

// 處理使用者發送消息的 form submit /send-message
func handleSendMessage(w http.ResponseWriter, r *http.Request) {
   // 解析表單
   if err := r.ParseForm(); err != nil {
      fmt.Println(err.Error())
      http.Error(w, "無法解析表單", http.StatusBadRequest)
      return
   }

   // 獲取用戶消息
   message := r.FormValue("message-input")
   if message == "" {
      fmt.Println("empty message @ Send Message")
      http.Error(w, "消息不能為空", http.StatusBadRequest)
      return
   }

   // 回傳用戶消息，實際處理會通過SSE進行
   w.Header().Set("Content-Type", "text/html")
   fmt.Fprintf(w, `<div class="message user-message"><div class="message-content">%s</div></div>`, message)
}
