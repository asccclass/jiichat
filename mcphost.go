package main
/*
## 程式說明
這個Go程式模擬了MCP Host與多個MCP Server進行能力協商的過程，主要功能包括：

1. MCP Host初始化 ：創建一個MCP Host實例，準備連接到多個MCP Server。
2. 與MCP Server連接並進行能力協商 ：
   
   - 連接到多個MCP Server（示例中是file_server和database_server）
   - 從每個Server獲取其提供的工具列表和描述
   - 在實際應用中，這會通過HTTP請求實現
3. 處理用戶查詢 ：
   
   - 收集所有連接的Server提供的工具
   - 將用戶查詢和工具描述一起發送給LLM
   - LLM分析用戶需求，決定使用哪些工具
4. 執行LLM選擇的工具 ：
   
   - 解析LLM選擇的工具名稱，找到對應的Server
   - 向相應的Server發送請求，執行選定的工具
這個程式展示了MCP Host如何管理多個MCP工具的完整流程，包括初始化連接、能力協商、工具選擇和執行。在實際應用中，HTTP請求和響應處理會更加複雜，但基本流程是相同的。
*/

import (
	"os"
	"io"
   "fmt"
   "time"
   "strings"
   "net/http"
   "encoding/json"
)
// LLM客戶端
type LLMClient struct {
	Name	  string            // 客戶端名稱
	Endpoint string
}

// MCP Server提供的工具定義
type Tool struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Parameters  map[string]string `json:"parameters,omitempty"`
}

// MCP Server的能力描述
type ServerCapabilities struct {
	Version  string `json:"version"`
	ServerID string `json:"server_id"`
	Tools    []Tool `json:"tools"`
}

// MCP Server結構
type MCPServer struct {
	ID           string					`json:"id"` // Server ID
	Name         string					`json:"name"` // Server名稱			
	Capabilities ServerCapabilities	`json:"capabilities"` // Server能力描述
	Endpoint     string			`json:"endpoint",oomitempty` // Server的API端點
	IsRelatedPrompt string     `json:"isRelatedPrompt",omitempty` // 是否與ID服務事項相關
	ProcessPrompt string			`json:"processPrompt",omitempty` // 處理ID服務事項的提示，若是則需要做何處理
}

// MCP Host結構
type MCPHost struct {
	ConnectedServers map[string]*MCPServer
	LLMClient        *LLMClient
}


// 連接到MCP Server並進行能力協商，只執行一次
func(h *MCPHost) AddCapabilities(serviceName, endpoint string) (error) {
	// 建立HTTP客戶端，設定逾時時間
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	// 發送GET請求
	resp, err := client.Get(endpoint)
	if err != nil {
		return fmt.Errorf("HTTP請求失敗: %v", err)
	}
	defer resp.Body.Close()
	// 檢查HTTP狀態碼
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP狀態碼錯誤: %d", resp.StatusCode)
	}
	// 讀取回應內容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("讀取回應內容失敗: %v", err)
	}
	// 解析JSON
	var server MCPServer
	if err := json.Unmarshal(body, &server); err != nil {
		return fmt.Errorf("JSON解析失敗: %v", err)
	}
	h.ConnectedServers[serviceName] = &server
	fmt.Printf("成功連接到MCP Server: %s，獲取到%d個工具\n", serviceName, len(server.Capabilities.Tools))
	return nil
}

// 初始化MCP Host
func NewMCPHost()(*MCPHost) {
	url := ""
	name := ""
	ollama := os.Getenv("OllamaUrl") 
	if ollama != "" {
		 name = "ollama"
		 url = strings.TrimSuffix(ollama, "/") + "/api/chat" // 確保URL不以斜杠結尾
	}
	// TODO: 檢查是否使用其他LLM客戶端配置，並根據需要初始化
	return &MCPHost{
		ConnectedServers: make(map[string]*MCPServer),
		LLMClient: &LLMClient{
			Name:     name,   // LLM客戶端名稱
			Endpoint: url,
		},
	}
}


/*
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// MCP Server提供的工具定義
type Tool struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Parameters  map[string]string `json:"parameters,omitempty"`
}

// MCP Server的能力描述
type ServerCapabilities struct {
	Version  string `json:"version"`
	ServerID string `json:"server_id"`
	Tools    []Tool `json:"tools"`
}

// 用戶查詢結構
type UserQuery struct {
	Query string `json:"query"`
}

// 發送給LLM的請求
type LLMRequest struct {
	UserQuery UserQuery         `json:"user_query"`
	Tools     []Tool            `json:"tools"`
	Context   map[string]string `json:"context,omitempty"`
}

// LLM的響應
type LLMResponse struct {
	Response     string   `json:"response"`
	SelectedTools []string `json:"selected_tools,omitempty"`
}

// MCP Host結構
type MCPHost struct {
	ConnectedServers map[string]*MCPServer
	LLMClient        *LLMClient
}

// MCP Server結構
type MCPServer struct {
	ID           string
	Capabilities ServerCapabilities
	Endpoint     string
}

// LLM客戶端
type LLMClient struct {
	Endpoint string
}

// 初始化MCP Host
func NewMCPHost() *MCPHost {
	return &MCPHost{
		ConnectedServers: make(map[string]*MCPServer),
		LLMClient: &LLMClient{
			Endpoint: "https://api.llm-service.example/v1/chat",
		},
	}
}

// 連接到MCP Server並進行能力協商
func (h *MCPHost) ConnectToServer(serverID, endpoint string) error {
	fmt.Printf("正在連接到MCP Server: %s (%s)\n", serverID, endpoint)
	
	// 模擬HTTP請求獲取Server能力
	// 在實際應用中，這裡會發送HTTP請求到endpoint
	capabilities := ServerCapabilities{
		Version:  "1.0",
		ServerID: serverID,
		Tools: []Tool{
			{
				Name:        "file_reader",
				Description: "讀取本地文件內容",
				Parameters: map[string]string{
					"file_path": "文件路徑",
				},
			},
			{
				Name:        "database_query",
				Description: "執行SQL查詢",
				Parameters: map[string]string{
					"query":      "SQL查詢語句",
					"database":   "數據庫名稱",
					"connection": "連接字符串",
				},
			},
		},
	}

	// 創建並存儲Server連接
	server := &MCPServer{
		ID:           serverID,
		Capabilities: capabilities,
		Endpoint:     endpoint,
	}

	h.ConnectedServers[serverID] = server
	fmt.Printf("成功連接到MCP Server: %s，獲取到%d個工具\n", serverID, len(capabilities.Tools))
	return nil
}

// 從所有連接的Server獲取工具列表
func (h *MCPHost) GetAllTools() []Tool {
	var allTools []Tool

	for serverID, server := range h.ConnectedServers {
		fmt.Printf("從Server %s獲取工具列表\n", serverID)
		for _, tool := range server.Capabilities.Tools {
			// 在工具名稱前加上serverID前綴，以區分不同Server的同名工具
			tool.Name = fmt.Sprintf("%s.%s", serverID, tool.Name)
			allTools = append(allTools, tool)
		}
	}

	return allTools
}

// 處理用戶查詢
func (h *MCPHost) ProcessUserQuery(query string) (*LLMResponse, error) {
	fmt.Printf("處理用戶查詢: %s\n", query)
	
	// 1. 獲取所有可用工具
	allTools := h.GetAllTools()
	fmt.Printf("獲取到總共%d個工具\n", len(allTools))

	// 2. 構建發送給LLM的請求
	llmRequest := LLMRequest{
		UserQuery: UserQuery{Query: query},
		Tools:     allTools,
		Context:   map[string]string{"session_id": "12345"},
	}

	// 3. 發送請求給LLM
	fmt.Println("發送請求給LLM，包含用戶查詢和工具描述")
	// 在實際應用中，這裡會發送HTTP請求到LLM API
	
	// 模擬LLM的響應
	llmResponse := &LLMResponse{
		Response: "我將幫助您查詢數據庫中的用戶信息",
		SelectedTools: []string{"database_server.database_query"},
	}

	fmt.Printf("LLM決定使用以下工具: %s\n", strings.Join(llmResponse.SelectedTools, ", "))
	return llmResponse, nil
}

// 執行LLM選擇的工具
func (h *MCPHost) ExecuteTools(llmResponse *LLMResponse, query string) error {
	for _, toolName := range llmResponse.SelectedTools {
		// 解析工具名稱，獲取serverID和實際的工具名
		parts := strings.SplitN(toolName, ".", 2)
		if len(parts) != 2 {
			return fmt.Errorf("無效的工具名稱格式: %s", toolName)
		}

		serverID, actualToolName := parts[0], parts[1]
		server, exists := h.ConnectedServers[serverID]
		if !exists {
			return fmt.Errorf("找不到Server: %s", serverID)
		}

		// 檢查工具是否存在
		var toolExists bool
		for _, tool := range server.Capabilities.Tools {
			if tool.Name == actualToolName {
				toolExists = true
				break
			}
		}

		if !toolExists {
			return fmt.Errorf("Server %s上找不到工具: %s", serverID, actualToolName)
		}

		// 模擬執行工具
		fmt.Printf("執行工具: %s.%s\n", serverID, actualToolName)
		// 在實際應用中，這裡會發送請求到Server執行工具
	}

	return nil
}

func main() {
	// 創建MCP Host
	host := NewMCPHost()

	// 連接到多個MCP Server
	err := host.ConnectToServer("file_server", "https://file-server.example/mcp")
	if err != nil {
		log.Fatalf("連接到file_server失敗: %v", err)
	}

	err = host.ConnectToServer("database_server", "https://db-server.example/mcp")
	if err != nil {
		log.Fatalf("連接到database_server失敗: %v", err)
	}

	// 處理用戶查詢
	userQuery := "查詢所有用戶的信息並保存到文件中"
	llmResponse, err := host.ProcessUserQuery(userQuery)
	if err != nil {
		log.Fatalf("處理用戶查詢失敗: %v", err)
	}

	// 執行LLM選擇的工具
	err = host.ExecuteTools(llmResponse, userQuery)
	if err != nil {
		log.Fatalf("執行工具失敗: %v", err)
	}

	fmt.Println("任務完成!")
}
*/