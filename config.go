package main

import (
   "encoding/json"
   "log"
   "os"
   "sync"
)

// Config 結構包含應用程序配置
type Config struct {
   Server    ServerConfig    `json:"server"`
   OpenAI    OpenAIConfig    `json:"openai"`
   Anthropic AnthropicConfig `json:"anthropic"`
}

// ServerConfig 包含服務器相關配置
type ServerConfig struct {
   Port         string `json:"port"`
   Host         string `json:"host"`
   TemplatePath string `json:"template_path"`
}

// OpenAIConfig 包含OpenAI API相關配置
type OpenAIConfig struct {
   APIKey string `json:"api_key"`
}

// AnthropicConfig 包含Anthropic API相關配置
type AnthropicConfig struct {
   APIKey string `json:"api_key"`
}

var (
   config     *Config
   configOnce sync.Once
)

// LoadConfig 從文件加載配置
func LoadConfig(filePath string) *Config {
   configOnce.Do(func() {
      // 預設配置
      config = &Config{
         Server: ServerConfig{
            Port:         "8080",
            Host:         "localhost",
            TemplatePath: "./templates",
         },
      }

      // 嘗試從文件讀取配置
      data, err := os.ReadFile(filePath)
      if err != nil {
         log.Printf("警告: 無法讀取配置文件 '%s': %v", filePath, err)
         log.Println("使用預設配置")

         // 創建範例配置文件
         createExampleConfig(filePath)
         return
      }

      // 解析JSON配置
      err = json.Unmarshal(data, config)
      if err != nil {
         log.Printf("警告: 解析配置文件失敗: %v", err)
         log.Println("使用預設配置")
      }
   })

   return config
}

// 創建範例配置文件
func createExampleConfig(filePath string) {
   exampleConfig := Config{
      Server: ServerConfig{
         Port:         "8080",
         Host:         "localhost",
         TemplatePath: "./templates",
      },
      OpenAI: OpenAIConfig{
         APIKey: "your-openai-api-key-here",
      },
      Anthropic: AnthropicConfig{
         APIKey: "your-anthropic-api-key-here",
      },
   }

   data, err := json.MarshalIndent(exampleConfig, "", "  ")
   if err != nil {
      log.Printf("無法創建範例配置: %v", err)
      return
   }

   err = os.WriteFile(filePath+".example", data, 0644)
   if err != nil {
      log.Printf("無法寫入範例配置文件: %v", err)
   } else {
      log.Printf("已創建範例配置文件: %s.example", filePath)
   }
}

// GetConfig 取得配置單例
func GetConfig() *Config {
   if config == nil {
      return LoadConfig("config.json")
   }
   return config
}
