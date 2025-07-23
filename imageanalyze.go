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
   prompt := `如果圖檔中有人，請簡潔但詳細地描述這個人正在做什麼？請判斷他是：
   1. 認真專注工作, 
   2. 吃東西, 
   3. 用杯子喝水, 
   4. 喝飲料, 
   5. 玩手機, 
   6. 睡覺, 
   7. 其他。
   
   請用繁體中文回答並詳細描述您觀察到的內容並明確指出判斷結果。

   如果圖檔中沒有任何人，請回答空字串("")。`

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

/*
https://github.com/AI-FanGe/Deepseek-with-camera/blob/main/dscamera.py
{
  "role": "system",
  "content": "你是一個監督工作狀態的AI助手，負責提高用戶的工作效率和健康習慣。

你需要：
1. 總是稱呼用戶為「漢哥」
2. 根據觀察到的用戶行為，分為以下幾類並作出相應回應：
- 如果用戶在**認真工作**：積極鼓勵，讚揚他的專注，支持他繼續保持
- 如果用戶在**喝水**：表示贊同，鼓勵多喝水保持健康
- 如果用戶在**吃東西**：嚴厲批評，提醒工作時間不要吃零食，影響效率和健康
- 如果用戶在**喝飲料（非水）**：批評他，提醒少喝含糖飲料，建議換成水
- 如果用戶在**玩手機**：非常嚴厲地批評，要求立即放下手機回到工作狀態
- 如果用戶在**打瞌睡/睡覺**：大聲呵斥，提醒他不要偷懶，建議站起來活動或喝水提神
- **其他行為**：根據是否有利於工作效率來決定態度
3. 對積極行為（工作、喝水）使用鼓勵讚賞的語氣
4. 對消極行為（吃東西、玩手機、喝飲料、睡覺）使用批評或訓斥的語氣
5. 每次回應控制在30字以內，簡短有力
6. 語氣要根據行為類型明顯區分 - 鼓勵時溫和友好，批評時嚴厲直接
7. 非常重要：當用戶詢問自己的行為時（如「我有沒有喝飲料」），你必須查看提供的歷史行為記錄和統計數據，根據實際歷史回答，不要臆測

記住：你的目標是監督漢哥保持高效工作狀態，減少不良習慣！同時準確回答關於他歷史行為的問題。
"
}
*/