package main

import(
	"os"
   "fmt"
	"sync"
	"time"
	"bufio"
   // "strconv"
	"net/http"
	"crypto/tls"
	"encoding/json"
)
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCP Message 結構
type MCPMessage struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	ID      *string     `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCP Client 結構
type MCPClient struct {
	serverURL    string
	Path			string // SSE 端點路徑
	token        string
	httpClient   *http.Client
	sseReader    *bufio.Scanner
	sseResponse  *http.Response
	msgCounter   int
	msgMu        sync.Mutex
	notifications chan MCPMessage
	isConnected  bool
	connMu       sync.RWMutex
}

// 處理來自服務器的 SSE 消息
func(c *MCPClient) handleSSEMessages() {
	defer func() {
		c.connMu.Lock()
		c.isConnected = false
		c.connMu.Unlock()
		if c.sseResponse != nil {
			c.sseResponse.Body.Close()
		}
		fmt.Println("SSE connection closed")
	}()

	for c.sseReader.Scan() {
		line := c.sseReader.Text()  // SSE 格式: "data: {json}"
		if len(line) > 6 && line[:6] == "data: " {
			jsonData := line[6:]			
			var msg MCPMessage
			if err := json.Unmarshal([]byte(jsonData), &msg); err != nil {
				fmt.Printf("Failed to parse SSE message: %v", err)
				continue
			}

			fmt.Printf("Received SSE message: %+v", msg)
			
			// 處理通知消息
			if msg.Method != "" {
				fmt.Println("mcpclient.go", msg)
				//  c.handleNotification(&msg)
				select {  // 發送到通知通道
				case c.notifications <- msg:
				default:
					fmt.Printf("Notification channel full, dropping message")
				}
			}
		}
	}

	if err := c.sseReader.Err(); err != nil {
		fmt.Printf("SSE reader error: %v", err)
	}
}

// 連接到 MCP 服務器的 SSE 端點
func(c *MCPClient) Connect()(error) {
	if c.token == "" {
		return fmt.Errorf("authentication required before connecting")
	}
	// 建立 SSE 連接
	sseURL := fmt.Sprintf("%s%ssse?token=%s", c.serverURL, c.Path, c.token)
	fmt.Printf("mcpclient.go: Connecting to MCP server via SSE: %s\n", sseURL)
	c.connMu.Lock()
	if c.isConnected {
		c.connMu.Unlock()
		return fmt.Errorf("already connected")
	}
	c.connMu.Unlock()
	// 創建 HTTP 請求
	req, err := http.NewRequest("GET", sseURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create SSE request: %w", err)
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to SSE: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("SSE connection failed with status: %d", resp.StatusCode)
	}

	c.sseResponse = resp
	c.sseReader = bufio.NewScanner(resp.Body)
	
	c.connMu.Lock()
	c.isConnected = true
	c.connMu.Unlock()
	go c.handleSSEMessages()  // 啟動 SSE 消息處理協程
	return nil
}

func NewMCPClient()(*MCPClient) {  // 創建 HTTP 客戶端，跳過 TLS 驗證 (僅用於開發)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},  // true 僅用於演示，生產環境應設為 false
		// Proxy: http.ProxyFromEnvironment,
	}
	
	return &MCPClient{
		serverURL:     os.Getenv("MCPSrv"), // MCP 服務器 URL
		token:         os.Getenv("MCPToken"),     // MCP 認證令牌
		Path:			   os.Getenv("MCPSrvPath"),  // 默認 MCP 路徑
		httpClient:    &http.Client{Transport: tr, Timeout: 30 * time.Second},
		notifications: make(chan MCPMessage, 100),
	}
}

// 尋找工具
func SearchTool(srv *MCPServer, action string)(*Tool, error) {
	for _, tool := range srv.Capabilities.Tools {  // 遍歷所有工具
		if tool.Name == action {  // 如果找到匹配的工具
			return &tool, nil
		}
	}
	return nil, fmt.Errorf("未找到名為 %s 的工具", action)  // 如果沒有找到，返回錯誤
}

// 執行 Tools 工具
func RunTools(req GenerateRequest, prompt string)(string, error) {
	if len(McpHost.ConnectedServers) == 0 {  // 檢查是否有連接的 MCP Server
		return "", fmt.Errorf("no connected MCP servers")	
	}	
	for _, srv:= range McpHost.ConnectedServers {  // 遍歷所有MCP Server
      if srv.IsRelatedPrompt == "" {
	     continue  // 如果沒有相關提示，則跳過
	  }
      s, err := parseIntent(req, prompt, srv) // (map[string]interface{}, error)	
      if err != nil {
         continue  // 如果解析不相關，則跳過  fmt.Println("解析意圖不相關:", err.Error())
      }
      action, ok := s["action"].(string)
      if !ok || action == "general_chat" {
	     continue  // 如果沒有動作，則跳過
	  }
	  tool, err := SearchTool(srv, action)  // (string, error)
	  if err != nil {
	     continue  // 如果查找工具失敗，則跳過
	  }
      parameters, ok := s["parameters"].(map[string]interface{})
	  if !ok {
	     parameters = make(map[string]interface{})
	  }
	  return callMCPTool(tool.Name, parameters)  // 調用 MCP 工具
	}
	return "", fmt.Errorf("未找到相關的 MCP Server 或工具")
}
