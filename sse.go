package main

import(
   "os"
   "fmt"
   "time"
   "strings"
   "net/http"
   "encoding/json"
)

// 消息結構
type ChatMessage struct {
   Type    string `json:"type"`
   Content string `json:"content"`
}

// 寫入響應
func ResponseChunks(response string) ([]ChatMessage) {
   chunks := []ChatMessage{}
   if os.Getenv("Stream") == "true" {    // 將回應分割成小塊以模擬流式輸出
      words := []rune(response)
      chunkSize := 5 // 每次發送5個字符

      for i := 0; i < len(words); i += chunkSize {
         end := i + chunkSize
         if end > len(words) {
            end = len(words)
         }
         chunk := string(words[i:end])
         chunks = append(chunks, ChatMessage{
            Type:    "chunk",
            Content: chunk,
         })
      }
   } else {  // 如果不需要流式輸出，則直接返回完整的回應
      /*
      s, err := MarkdownToHTML(response)
      if err != nil {
         s = response   
      }
         */

      chunks = append(chunks, ChatMessage{
         Type:    "chunk",
         Content: response,
      })
   }
   // fmt.Printf("模型: %s，使用者消息: %s，AI回應：%s\n", model, userMessage, response)
   return chunks
}

// AI回應（TODO: 需要增加記憶體）
func AIResponse(model, userMessage string)([]ChatMessage) {
   // 根據不同模型準備不同的回應
   response := ""
   var err error
   if os.Getenv("OllamaUrl") != "" {
      ola := AIs["Ollama"].(*OllamaClient)
      response, err = ola.Ask(model, userMessage, "", nil)
      if err != nil {
         response = "抱歉！無法處理您的請求。(" + err.Error() + ")"
      }
   }
   return ResponseChunks(response)
}

// 輸出 SSEChat 處理的聊天結果
func ResponseUser(w http.ResponseWriter, responses []ChatMessage) { 
   w.Header().Set("Content-Type", "text/event-stream")
   w.Header().Set("Cache-Control", "no-cache")
   w.Header().Set("Connection", "keep-alive")
   w.Header().Set("Access-Control-Allow-Origin", "*")
   
   flusher, ok := w.(http.Flusher)  // 創建SSE刷新器
   if !ok {
      http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
      return
   }
   // 逐步發送回應片段 SSE 格式要求每行以 \n 結尾，而 \r\n 會被視為額外的換行符。
   for _, chunk := range responses {      
      data, err := json.Marshal(chunk)  // 將消息轉換為JSON
      if err != nil {
         fmt.Println("JSON編碼錯誤:", err)
         continue
      }
      cleanContent := strings.ReplaceAll(string(data), "\r\n", "\\n")
      cleanContent = strings.ReplaceAll(cleanContent, "\r", "\\n")
      cleanContent = strings.ReplaceAll(cleanContent, "\n", "\\n")
      fmt.Fprintf(w, "data: %s\n\n", cleanContent)  // 發送SSE格式的消息
      flusher.Flush()
      time.Sleep(100 * time.Millisecond)  // 模擬打字延遲
   }

   completeMsg := ChatMessage{  // 發送完成信號
      Type:    "complete",
      Content: "",
   }
   completeData, _ := json.Marshal(completeMsg)
   fmt.Fprintf(w, "data: %s\n\n", completeData)
   flusher.Flush()
}

// 接收 SSE 請求並處理聊天
func SSEChat(w http.ResponseWriter, r *http.Request) {  
   // 獲取用戶消息和選擇的模型
   message := r.URL.Query().Get("message")
   model := r.URL.Query().Get("model")
   if message == "" {
      fmt.Println("使用者傳送空白資訊")
      return
   }
   ResponseUser(w, AIResponse(model, message))  // 模擬AI回應模式 
}

// 處理新對話請求
func handleNewChat(w http.ResponseWriter, r *http.Request) {
   // 返回空的對話區域以重置對話
   w.Header().Set("Content-Type", "text/html")
   fmt.Fprint(w, `<div class="message-container">
      <div class="message ai-message"><div class="message-content">您好！我已經為您開始了一個新對話。我能幫您什麼忙嗎？</div></div>
   </div>`)
}

// 處理使用者發送消息的 form submit
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
