package main

import (
   "fmt"
   "time"
   "strings"
   "net/http"
   "encoding/json"
)

// 輸出 SSEChat 處理的聊天結果
func Response2User(w http.ResponseWriter, responses []ChatMessage) { 
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
      // flusher.Flush()
      time.Sleep(50 * time.Millisecond)  // 模擬打字延遲
   }

   completeMsg := ChatMessage{  // 發送完成信號
      Type:    "complete",
      Content: "",
   }
   completeData, _ := json.Marshal(completeMsg)
   fmt.Fprintf(w, "data: %s\n\n", completeData)
   flusher.Flush()
}

// 寫入響應
func ResponseChunksOld(w http.ResponseWriter, response string)/* ([]ChatMessage) */{
   // chunks := []ChatMessage{}
   // fmt.Println("response.go", "Ollama 回應:", response)
   w.Header().Set("Content-Type", "text/event-stream")
   w.Header().Set("Cache-Control", "no-cache")
   w.Header().Set("Connection", "keep-alive")
   w.Header().Set("Access-Control-Allow-Origin", "*")
   w.WriteHeader(http.StatusOK)  // 設定狀態碼為200 OK
   // 確保響應立即發送
	flusher, ok := w.(http.Flusher)
   if !ok {
      fmt.Println("response.go", "服務器不支持串流回應")
		http.Error(w, "服務器不支持串流回應", http.StatusInternalServerError)
      // TODO: 這邊可以改成設定參數後，採用一次回復
      return
	}
   // 設定超時和錯誤恢復
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Stream handler panic: %v", r)
		}
	}()
   // 模擬流式輸出  

   // 修正: 確保 JSON 字串內的換行符被正確轉義
   cleanContent := strings.ReplaceAll(string(response), "\r\n", "\\n")
   cleanContent = strings.ReplaceAll(cleanContent, "\r", "\\n")
   // cleanContent = strings.ReplaceAll(cleanContent, "\n", "\\n")

   words := strings.Fields(cleanContent)
   for i, word := range words {  // 創建響應結構
      restxt := ChatMessage{
         Type:    "chunk",
         Content: word,
      }
      if i > 0 {
         restxt.Content = " " + restxt.Content  // 在單詞間添加空格
      }
      // 轉換為JSON
      jsonData, err := json.Marshal(restxt)
      if err != nil {
         fmt.Println("JSON編碼錯誤:", err.Error())
         continue
      }

      fmt.Fprintf(w, "data: %s\n\n", string(jsonData))  // 寫入SSE格式的數據
      fmt.Println("response.go", i, "發送 SSE 數據:", "data:", string(cleanContent))
		flusher.Flush()
      time.Sleep(time.Millisecond * 50)  // 模擬打字速度
   }
   // 發送結束信號
   completeMsg := ChatMessage{
      Type:    "complete",
      Content: "",
   }
   completeData, _ := json.Marshal(completeMsg)
   fmt.Fprintf(w, "data: %s\n\n", completeData)
   flusher.Flush()
/*
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
      chunks = append(chunks, ChatMessage{
         Type:    "chunk",
         Content: response,
      })
   }
   return chunks
*/
}


// 寫入響應
func ResponseChunks(w http.ResponseWriter, response string)/* ([]ChatMessage) */{
   // chunks := []ChatMessage{}
   // fmt.Println("response.go", "Ollama 回應:", response)
   w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
   w.Header().Set("Cache-Control", "no-cache")
   w.Header().Set("Connection", "keep-alive")
   w.Header().Set("Access-Control-Allow-Origin", "*")
   w.Header().Set("X-Accel-Buffering", "no")  // 避免代理伺服器緩衝,適用於 Nginx/Apache
   w.WriteHeader(http.StatusOK)  // 設定狀態碼為200 OK
   // 確保響應立即發送
	flusher, ok := w.(http.Flusher)
   if !ok {
      fmt.Println("response.go", "服務器不支持串流回應")
		http.Error(w, "服務器不支持串流回應", http.StatusInternalServerError)
      // TODO: 這邊可以改成設定參數後，採用一次回復
      return
	}
   // 設定超時和錯誤恢復
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Stream handler panic: %v", r)
		}
	}()
   // 模擬流式輸出  

   // 修正: 確保 JSON 字串內的換行符被正確轉義
   cleanContent := strings.Split(string(response), "\n")
   for _, word := range cleanContent {  // 創建響應結構
      restxt := ChatMessage{
         Type:    "chunk",
         Content: word,
      }
      // 轉換為JSON
      jsonData, err := json.Marshal(restxt)
      if err != nil {
         fmt.Println("JSON編碼錯誤:", err.Error())
         return
      }
      fmt.Fprintf(w, "data: %s\n\n", string(jsonData))  // 寫入SSE格式的數據
      flusher.Flush()
      time.Sleep(time.Millisecond * 50)  // 模擬打字速度
   }

   // 發送結束信號
   completeMsg := ChatMessage{
      Type:    "complete",
      Content: "",
   }
   completeData, _ := json.Marshal(completeMsg)
   fmt.Fprintf(w, "data: %s\n\n", completeData)
   flusher.Flush()
   time.Sleep(time.Millisecond * 5000)  // 模擬打字速度
}