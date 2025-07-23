package main

import(
	"os"
	"io"
	"fmt"
	"time"
	"bytes"
	"net/http"
	"encoding/json"
)

// CallToolRequestParams MCP 工具調用請求參數
type CallToolRequestParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// CallToolRequest MCP 工具調用請求
type CallToolRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	ID      string      `json:"id"`
	Params CallToolRequestParams `json:"params"`
}

// CallToolResultContent MCP 工具調用結果內容
type CallToolResultContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type CallToolResult struct {	
	Content []CallToolResultContent `json:"content"`
	IsError *bool                   `json:"isError,omitempty"`
}

// CallToolResult MCP 工具調用結果
type CallToolResults struct {
	JSONRPC string                  `json:"jsonrpc"`
	ID      string                  `json:"id"`
	Result CallToolResult			  `json:"result,omitempty"`
}

// callMCPTool 調用 MCP Server 的工具
func callMCPTool(toolName string, args map[string]interface{}) (string, error) {
	request := CallToolRequest{
		JSONRPC: os.Getenv("JSONRPCVersion"), //JSONRPC Version 版本
		Method: "tools/call",	// tools/list)列出所有工具名稱  tools/call)調用工具
		Params: CallToolRequestParams{
			Name:      toolName,   // 工具名稱
			Arguments: args,
		},
	}
	serverURL := os.Getenv("MCPSrv") // MCP Server URL
	if serverURL == "" {
		return "", fmt.Errorf("MCPSrv environment variable not set")
	}
	serverPath := os.Getenv("MCPSrvPath") // MCP Server Path
	if serverPath == "" {
		serverPath = "/" // Default path if not set
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("marshal request: %s", err.Error())
	}
   if os.Getenv("Debug") == "true" {
		fmt.Println("MCPServer 請求內容:", string(jsonData))  // MCPServer 請求內容
	}
	hClient := &http.Client {
	   Timeout: 60 * time.Second,
	}
	resp, err := hClient.Post(serverURL + serverPath + "request", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("make request: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server error (status %d): %s", resp.StatusCode, string(body))
	}
	var msg CallToolResults
	if err := json.Unmarshal(body, &msg); err != nil {
		return "", fmt.Errorf("unmarshal response: %s", err.Error())
	}
	result := msg.Result
	if result.IsError != nil && *result.IsError {
		return "", fmt.Errorf("tool error: %s", result.Content[0].Text)
	}
	if len(result.Content) > 0 {
		return result.Content[0].Text, nil
	}
	return "操作完成", nil
}