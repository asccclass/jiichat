package main

import(
   "io"
   "os"
   "fmt"
   "strings"
   "net/http"
   "encoding/base64"
)

// 處理使用者發送消息的 form submit
func handleImageAnalyze(w http.ResponseWriter, r *http.Request) {
   // 解析表單
   if err := r.ParseMultipartForm(10 << 20); err != nil {  // 10 MB
      fmt.Println(err.Error())
      http.Error(w, "無法解析表單", http.StatusBadRequest)
      return
   }

   // 獲取用戶消息
   file, header, err := r.FormFile("image")
	if err != nil {
		ResponseUser(w, ResponseChunks("無法獲取攝像頭畫面"))
		return
	}
	defer file.Close()
   // 檢查文件類型
	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
		ResponseUser(w, ResponseChunks("畫面格式錯誤"))
		return
	}
   // 讀取文件內容
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		ResponseUser(w, ResponseChunks("無法讀取畫面數據"))
		return
	}

   fmt.Println("Received image:", header.Filename, "Size:", len(fileBytes), "bytes")

	// 將圖片編碼為base64
	base64Image := base64.StdEncoding.EncodeToString(fileBytes)
   if base64Image == "" {
      ResponseUser(w, ResponseChunks("無法獲取攝像頭畫面"))
      return
   }
   response := ""   

      prompt = `請簡潔但詳細地描述當前畫面中看到的內容，包括：
1. 主要物體、人物或活動
2. 場景位置和環境
3. 值得注意的變化或動作
4. 整體狀況評估

請用繁體中文回答，保持描述簡潔明確，適合即時監控使用。`   
   if os.Getenv("OllamaUrl") != "" {  // 如果有上傳圖片，則使用Ollama進行分析
      ola := AIs["Ollama"].(*OllamaClient)
      response, err = ola.Ask("llama3.2-vision:latest", prompt, base64Image, nil)
      if err != nil {
         fmt.Println("Ollama Ask error:", err.Error(), "base64Image length:", len(base64Image))
         ResponseUser(w, ResponseChunks("抱歉！無法處理您的請求。(" + err.Error() + ")"))
         return
      }
   }
   chunks := ResponseChunks(response)
   ResponseUser(w, chunks)  // 模擬AI回應模式 
}