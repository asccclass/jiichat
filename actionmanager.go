package main

import(
   "sync"
	"time"
)
type ActionState string
const (
	StatePending    ActionState = "pending"
	StateCompleted  ActionState = "completed"
	StateExecuting  ActionState = "executing"
	StateFailed     ActionState = "failed"
)

// FieldInfo 表示欄位的資訊和狀態
type FieldInfo struct {
	Name        string      `json:"name"`
	Value       interface{} `json:"value,omitempty"`
	Required    bool        `json:"required"`
	Collected   bool        `json:"collected"`
	Description string      `json:"description"`
	Type        string      `json:"type"` // string, int, bool, date, etc.
}

// PendingAction 表示一個待完成的操作
type PendingAction struct {
	ID          string                 `json:"id"`
	ActionType  string                 `json:"action_type"`
	UserID      string                 `json:"user_id"`
	Fields      map[string]*FieldInfo  `json:"fields"`
	State       ActionState            `json:"state"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Context     map[string]interface{} `json:"context,omitempty"`
	NextPrompt  string                 `json:"next_prompt,omitempty"`
}

// ActionTemplate 定義操作的模板
type ActionTemplate struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Fields      map[string]*FieldInfo `json:"fields"`
	MCPTool     string                `json:"mcp_tool"` // 對應的 MCP Server 工具名稱
}

// MemoryManager 管理待完成的操作
type MemoryManager struct {
	mu              sync.RWMutex
	pendingActions  map[string]*PendingAction
	actionTemplates map[string]*ActionTemplate
	mcpClient       MCPClient
}
