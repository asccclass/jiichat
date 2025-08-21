package main

import (
   "fmt"
   "regexp"
   "strings"
   "encoding/json"
)

// parseIntentWithOllama 使用 Ollama 解析使用者意圖
func parseIntent(req GenerateRequest, userInput string, mcpsrv *MCPServer) (map[string]interface{}, error) {
   prompt := fmt.Sprintf("%s\n\n使用者輸入：`%s`", mcpsrv.IsRelatedPrompt, userInput)
   response := ""
   var err error
   if AIs["Ollama"].(*OllamaClient) != nil {
      jData, err := AIs["Ollama"].(*OllamaClient).Prompt2String(req, "user", prompt, "")
      if err != nil {
         return nil, fmt.Errorf("prepare prompt for ollama: %s", err.Error())
      }
      res, err := AIs["Ollama"].(*OllamaClient).Send2LLM(jData, false)  // (string, error) 
      if err != nil {
         return nil, fmt.Errorf("query ollama for intent: %s", err.Error())
      }
      response = res
   } else {
      fmt.Println("No Ollama client initialized, cannot parse intent")
      return nil, fmt.Errorf("Not any LLM is initialized")
   }
   // 清除 <think>...</think> 標籤   
   re := regexp.MustCompile(`(?s)^.*</think>`)
	response = re.ReplaceAllString(response, "")
   // 清除 ```
   response = strings.ReplaceAll(response, "```", "")
   // 清理回應，只保留 JSON 部分
   response = strings.TrimSpace(response)
   start := strings.Index(response, "{")
   end := strings.LastIndex(response, "}") + 1
   if start >= 0 && end > start {
      response = response[start:end]
   }
   var intent map[string]interface{}
   if err = json.Unmarshal([]byte(response), &intent); err != nil { // 如果 JSON 解析失敗，預設為一般對話
      return map[string]interface{}{
         "is_related": false,
         "action":          "general_chat",
         "parameters":      map[string]interface{}{},
      }, nil 

      return nil, fmt.Errorf("Parse intent failed（Mcpsrv 回傳格式錯誤）: %s", err.Error())      
      /* return map[string]interface{}{
         "is_related": false,
         "action":          "general_chat",
         "parameters":      map[string]interface{}{},
      }, nil */
   }   
   return intent, nil
}