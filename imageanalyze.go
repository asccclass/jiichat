package main

import(
   "io"
   "os"
   "fmt"
   "strings"
   "net/http"
   "encoding/base64"
)

// 處理使用者發送消息的 form submit POST /video/analyze
func handleImageAnalyze(w http.ResponseWriter, r *http.Request) {
   if err := r.ParseMultipartForm(10 << 20); err != nil {  // 10 MB // 解析表單
      fmt.Println(err.Error())
		Response2User(w, ResponseChunks("無法獲取表單資料"))
      return
   }
   file, header, err := r.FormFile("image")  // 獲取 form 中的圖片文件
	if err != nil {
		Response2User(w, ResponseChunks("無法獲取攝像頭畫面"))
		return
	}
	defer file.Close()
   
	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {  // 檢查文件類型
		Response2User(w, ResponseChunks("畫面格式錯誤"))
		return
	}
	fileBytes, err := io.ReadAll(file)  // 讀取 image
	if err != nil {
		Response2User(w, ResponseChunks("無法讀取image"))
		return
	}
	// 將圖片編碼為base64
	base64Image := base64.StdEncoding.EncodeToString(fileBytes)
   if base64Image == "" {
      Response2User(w, ResponseChunks("無法獲取攝像頭畫面"))
      return
   }
   prompt := `請簡潔但詳細地描述附件圖檔中看到的內容，包括：
1. 主要物體、人物或活動
2. 場景位置和環境
3. 值得注意的變化或動作
4. 整體狀況評估

請用繁體中文回答，保持描述簡潔明確。`

   response := ""
   if os.Getenv("OllamaUrl") != "" {  // 如果有上傳圖片，則使用Ollama進行分析
      ola := AIs["Ollama"].(*OllamaClient)
      response, err = ola.Ask("llama3.2-vision:latest", prompt, base64Image, nil)
      if err != nil {
         fmt.Println("Ollama Ask error:", err.Error(), "base64Image length:", len(base64Image))
         Response2User(w, ResponseChunks("抱歉！無法處理您的請求。(" + err.Error() + ")"))
         return
      }
   }
   if response == "" {
      fmt.Println("抱歉！您的請求沒有回應，請稍後再試。")
      return   
   }
   chunks := ResponseChunks(response)
   Response2User(w, chunks)  // 模擬AI回應模式 
}