package main

import (
   "os"
   //"io"
   "fmt"
   "time"
   "bytes"
   "strings"
   // "strconv"
   "net/http"
   "encoding/json"
   // "path/filepath"
   // "mime/multipart"
   // "github.com/joho/godotenv"
)

// Details represents the details of a model
type Details struct {
   ParentModel       string   `json:"parent_model"`
   Format            string   `json:"format"`
   Family            string   `json:"family"`
   Families          []string `json:"families"`
   ParameterSize     string   `json:"parameter_size"`
   QuantizationLevel string   `json:"quantization_level"`
}

// Model represents a single model's information 模型資訊結構
type Model struct {
   Name        string    `json:"name"`
   Model			string    `json:"model"`
   ModifiedAt  time.Time `json:"modified_at"`
   Size        int64     `json:"size"`
   Digest      string    `json:"digest"`
	Details		Details   `json:"details"`
   Description string    `json:"description,omitempty"`
	SizeGB		string	`json:"size_gb,omitempty"`	// 轉換後的大小
	ModiTime	string		`json:"modi_time,omitempty"`	// 轉換後的時間
}

// ModelsWrapper wraps the models array
type ModelsWrapper struct {
   Models []Model `json:"models"`
}

// 生成選項結構
type Options struct {
   Temperature      float64 `json:"temperature,omitempty"`
   TopP             float64 `json:"top_p,omitempty"`
   TopK             int     `json:"top_k,omitempty"`
   NumPredict       int     `json:"num_predict,omitempty"`
   NumKeep          int     `json:"num_keep,omitempty"`
   Seed             int     `json:"seed,omitempty"`
   FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
   PresencePenalty  float64 `json:"presence_penalty,omitempty"`
   Mirostat         int     `json:"mirostat,omitempty"`
   MirostatEta      float64 `json:"mirostat_eta,omitempty"`
   MirostatTau      float64 `json:"mirostat_tau,omitempty"`
   Stop             string  `json:"stop,omitempty"`
}

// 生成回應結構
type GenerateResponse struct {
   Model              string    `json:"model"`
   CreatedAt          time.Time `json:"created_at"`
   Message            Message    `json:"message"`
	DoneReason         string    `json:"done_reason"`
   Done               bool      `json:"done"`
   TotalDuration      int64     `json:"total_duration,omitempty"`
   LoadDuration       int64     `json:"load_duration,omitempty"`
   PromptEvalCount    int64     `json:"prompt_eval_count,omitempty"`
   PromptEvalDuration int64     `json:"prompt_eval_duration,omitempty"`
   EvalCount          int64     `json:"eval_count,omitempty"`
   EvalDuration       int64     `json:"eval_duration,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 生成請求結構
type GenerateRequest struct {
   Model    string	`json:"model"`
   Messages   []Message `json:"messages"`
   System   string	`json:"system,omitempty"`
   Format   string	`json:"format,omitempty"`
   Stream   bool	`json:"stream"`
   Options  Options	`json:"options,omitempty"`
   Images   []string	`json:"images,omitempty"`
   Context  []int	`json:"context,omitempty"`
   Template string	`json:"template,omitempty"`
}

// 定義請求結構
type RequestPayload struct {
   Message string `json:"message"` // 客戶端發送的訊息
   UserID  string `json:"user_id"` // 可選：用於標識用戶
}

type OllamaClient struct {
   URL		string		`json:"url"`		// Ollama 服務的 URL
   Models	[]Model		`json:"models"`   // 模型列表回應結構
}

// 獲取所有可用模型
func(app *OllamaClient) getModels() ([]Model, error) {
   // 這樣可以保留 DefaultTransport 內建的撥號逾時、TLS 握手逾時等合理設置
	transport := http.DefaultTransport.(*http.Transport).Clone()
   // 如果你需要設置代理或自定義 TLS 配置，在這裡添加
	// proxyURL, _ := url.Parse("http://your-proxy-server:8080")
	// transport.Proxy = http.ProxyURL(proxyURL)
	// transport.TLSHandshakeTimeout = 10 * time.Second // 已經在 DefaultTransport 中預設
   // 可以根據需要設定其他 Client 參數，例如 CheckRedirect、Timeout 等
   client := &http.Client{
      Transport: transport,
      Timeout: 60 * time.Second, // 整個請求的逾時時間
   }
   resp, err := client.Get(app.URL + "/api/tags")
   if err != nil {
      return nil, fmt.Errorf("連接到 Ollama 服務失敗: %v", err)
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("獲取模型列表失敗，狀態碼: %d", resp.StatusCode)
   }
   var listResp ModelsWrapper
   if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
      return nil, fmt.Errorf("解析回應失敗: %v", err)
   }
   return listResp.Models, nil
}

// 顯示所有可用模型 <option value="gpt-4">GPT-4</option>
func(app *OllamaClient) ListModelsFromWeb(w http.ResponseWriter, r *http.Request) {
   s := []string{}
   for _, model := range app.Models {
      val := strings.ReplaceAll(strings.ToLower(model.Name), " ", "-")
      s = append(s, fmt.Sprintf("<option value=\"%s\">%s</option>", val, model.Name))
   }
   jsonData, err := json.Marshal(s)
   if err != nil {
      fmt.Println("序列化 JSON 失敗:", err)
      return
   }
   // 回傳用戶消息，實際處理會通過SSE進行
   w.Header().Set("Content-Type", "application/json")
   w.WriteHeader(http.StatusOK)
   w.Write(jsonData)	// 將 JSON 數據寫入響應
}

// 送出給Ollams
func(app *OllamaClient) Send2LLM(jsonData string)(string, error) {
   resp, err := http.Post(app.URL+"/api/chat", "application/json", bytes.NewBuffer([]byte(jsonData))) // 發送請求給 Ollama   
   if err != nil {
      return "", fmt.Errorf("發送請求失敗: %v", err)
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK { // 如果狀態碼不是 200 OK，則返回錯誤
      return "", fmt.Errorf("%s生成回應失敗，狀態碼: %d", app.URL,resp.StatusCode)
   }
   // 解析回應
   var genResp GenerateResponse
   if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
      return "", fmt.Errorf("解析回應失敗: %v", err)
   }
   return genResp.Message.Content, nil
}

// 轉換 GenerateRequest 為字串格式
func(app *OllamaClient) Prompt2String(req GenerateRequest, role, prompt string)(string, error) {
   req.Messages = append(req.Messages, Message{Role: role, Content: prompt})  // 如果沒有工具套用，則使用原始提示
   jData, err := json.Marshal(req)  // 將請求轉為 JSON
   if err != nil {
      return "", fmt.Errorf("Marshal request failed, 序列化請求失敗: %v", err)
   }
   return string(jData), nil
}

// 產生回應（非串流模式） ola.Ask(model, userMessage, nil)
func(app *OllamaClient) Ask(modelName, userinput string, base64Image string, files []string) (string, error) {
   prompt := strings.TrimSpace(userinput)
   if prompt == "" {
      return "", fmt.Errorf("No data")
   }
   if modelName == "" {
      modelName = "phi4:latest"  // 預設模型名稱
   }
   reqBody := GenerateRequest {  // 初始化
      Model:  modelName,
      Messages: []Message{},   // role, content
      Stream: false,          // 不使用串流模式
   }
   // 如果有上傳圖片，將圖片編碼為 base64 並添加到請求中
   if base64Image != "" {
      reqBody.Images = []string{base64Image}
   } else {
      toolsResponse, err := RunTools(reqBody, prompt)  // (map[string]interface, error)    // MCP 工具套用
      if err == nil {  // 如果有工具套用，則使用工具回應
         // jData, err := app.Prompt2String(reqBody, "user", "把下列內容，用人類的語氣重新改寫，請使用繁體中文回答，並去除掉簡體字：" + toolsResponse)  // 如果沒有工具套用，則使用原始提示
         // if err != nil {   
            return toolsResponse, nil
         // }
         if os.Getenv("Debug") == "true" {
            fmt.Println("Tools response:", toolsResponse)
         }
         // return app.Send2LLM(string(jData))
      } else {
         fmt.Println("toolsResponse RunTools error:", err.Error())
      }
   }
   jData, err := app.Prompt2String(reqBody, "user", prompt)  // 如果沒有工具套用，則使用原始提示
   if os.Getenv("Debug") == "true" {
      fmt.Println("Prompt2String:", jData)
   }
   if err != nil {   
      fmt.Print("Prepare prompt failed: ", err.Error())
      return "", fmt.Errorf("prepare prompt for ollama: %s", err.Error())
   }
/*
   // 如果有上傳文件，將文件內容添加到提示
   if len(files) > 0 {
   	var fileContents []string
   	for _, file := range files {
   		content, err := os.ReadFile(file)
   		if err != nil {
   			return "", fmt.Errorf("讀取文件 %s 失敗: %v", file, err)
   		}
   		fileContents = append(fileContents, fmt.Sprintf("\n文件 '%s' 內容:\n%s", filepath.Base(file), string(content)))
   	}

   	// 將文件內容附加到提示
   	if len(fileContents) > 0 {
   		reqBody.Prompt += "\n\n以下是相關文件內容，請參考這些內容回答我的問題:" + strings.Join(fileContents, "\n\n")
   	}
   }
*/
   return app.Send2LLM(string(jData))
}

func(app *OllamaClient) AddRouter(router *http.ServeMux) {
   router.HandleFunc("GET /models", app.ListModelsFromWeb)                // 列出工具列表
}

func NewOllamaClient()(*OllamaClient) {
   app := &OllamaClient{
      URL: os.Getenv("OllamaUrl"),
   }
   if app.URL == "" {
      fmt.Printf("未設定 envfile 中的 OllamaUrl 參數")
      return nil
   }
	var err error
   app.Models, err = app.getModels()
   if err != nil {
      fmt.Printf("獲取模型列表失敗: %s", err.Error())
      return nil
   }

   for i, model := range app.Models {
      size := float64(model.Size) / (1024 * 1024 * 1024)
      app.Models[i].SizeGB = fmt.Sprintf("%.2f GB", size)
      app.Models[i].ModiTime = model.ModifiedAt.Format("2006-01-02 15:04:05")
   }
   return app
}

/*

// 上傳檔案處理
func handleFileUpload() []string {
   files := []string{}
   for {
   	fmt.Println("\n請選擇操作:")
   	fmt.Println("1. 上傳檔案")
   	fmt.Println("2. 完成上傳並繼續")
   	
   	choice := getUserInput("請輸入選項 (1-2): ")
   	
   	switch choice {
   	case "1":
   		path := getUserInput("請輸入檔案路徑: ")
   		if _, err := os.Stat(path); os.IsNotExist(err) {
   			fmt.Printf("錯誤: 檔案 '%s' 不存在\n", path)
   			continue
   		}
   		files = append(files, path)
   		fmt.Printf("已添加檔案: %s\n", path)
   	case "2":
   		return files
   	default:
   		fmt.Println("無效選項，請重新選擇")
   	}
   }
}
*/
