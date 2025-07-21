package main

import(
   "os"
   "fmt"
   "time"
   "strings"
   "net/http"
   "encoding/json"
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
		writeErrorResponse(w, "無法獲取攝像頭畫面")
		return
	}
	defer file.Close()
   // 檢查文件類型
	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
		writeErrorResponse(w, "畫面格式錯誤")
		return
	}
   // 讀取文件內容
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		writeErrorResponse(w, "無法讀取畫面數據")
		return
	}

	// 將圖片編碼為base64
	base64Image := base64.StdEncoding.EncodeToString(fileBytes)

   // 回傳用戶消息，實際處理會通過SSE進行
   w.Header().Set("Content-Type", "text/html")
   fmt.Fprintf(w, `<div class="message user-message"><div class="message-content">%s</div></div>`, message)
}
