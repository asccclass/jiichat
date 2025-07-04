<!DOCTYPE html>
<html lang="zh-TW">
<head>
   <meta charset="UTF-8">
   <meta name="viewport" content="width=device-width, initial-scale=1.0">
   <title>AI 聊天界面</title>
   <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
   <script src="https://cdnjs.cloudflare.com/ajax/libs/htmx/1.9.5/htmx.min.js"></script>
   <script src="https://cdnjs.cloudflare.com/ajax/libs/htmx/1.9.5/ext/sse.min.js"></script>
   <style>
        :root {
            --primary-color: #10a37f;
            --bg-color: #ffffff;
            --text-color: #343541;
            --secondary-bg: #f7f7f8;
            --border-color: #e5e5e6;
            --hover-color: #f1f1f2;
            --ai-msg-bg: transparlent;
            --user-msg-bg: #ffffff;
        }

        body {
            font-family: 'Söhne', ui-sans-serif, system-ui, -apple-system, 'Segoe UI', Roboto, Ubuntu, Cantarell, sans-serif;
            margin: 0;
            padding: 0;
            color: var(--text-color);
            background-color: var(--bg-color);
            display: flex;
            height: 100vh;
        }

        .sidebar {
            width: 260px;
            background-color: #202123;
            color: white;
            display: flex;
            flex-direction: column;
            height: 100%;
        }

        .new-chat-btn {
            margin: 10px;
            padding: 12px;
            border: 1px solid rgba(255, 255, 255, 0.2);
            border-radius: 5px;
            display: flex;
            align-items: center;
            cursor: pointer;
            transition: background-color 0.3s;
        }

        .new-chat-btn:hover {
           background-color: rgba(255, 255, 255, 0.1);
        }

        .new-chat-btn i {
           margin-right: 10px;
        }

        .sidebar-footer {
           margin-top: auto;
           padding: 15px;
           border-top: 1px solid rgba(255, 255, 255, 0.2);
        }

        .main-content {
           flex-grow: 1;
           display: flex;
           flex-direction: column;
           height: 100%;
           position: relative;
        }

        .chat-history {
           flex-grow: 1;
           overflow-y: auto;
           padding: 20px 0;
        }

        .message_OLD {
           margin: 0 auto;
           line-height: 1.5;
        }
	.message {
           padding: 10px;
           max-width: 800px;
           margin-bottom: 15px;
           width: 100%;
           clear: both;
           display: flex;
           box-sizing: border-box;
        }
	.message-content {
           padding: 10px 15px;
           border-radius: 20px;
           word-wrap: break-word;
        }

	.ai-message .message-content {
            background-color: #7030a0; /* 紫色背景 */
            color: white;
            border-top-left-radius: 5px;
        }

        .user-message {
           background-color: var(--user-msg-bg);
        }

        .message-avatar {
           width: 30px;
           height: 30px;
           border-radius: 3px;
           background-color: #10a37f;
           color: white;
           display: flex;
           align-items: center;
           justify-content: center;
           margin-right: 15px;
           flex-shrink: 0;
        }

        .user-avatar {
           background-color: #5436DA;
        }

        .message-content {
           flex-grow: 1;
           overflow-wrap: break-word;
        }

        .message-content pre {
           background-color: #282c34;
           color: #abb2bf;
           padding: 15px;
           border-radius: 6px;
           overflow-x: auto;
        }

        .message-content code {
           font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, Courier, monospace;
        }

        .input-area {
           border-top: 1px solid var(--border-color);
           padding: 20px;
           position: sticky;
           bottom: 0;
           background-color: var(--bg-color);
        }

        .model-select {
           display: flex;
           justify-content: center;
           margin-top: 15px;
        }

        .model-select select {
           padding: 8px 12px;
           border-radius: 6px;
           border: 1px solid var(--border-color);
           background-color: var(--bg-color);
           font-size: 14px;
        }

        .input-container {
            display: flex;
            align-items: center;
            position: relative;
            max-width: 800px;
            margin: 0 auto;
        }

        .input-box {
            flex-grow: 1;
            padding: 12px 45px 12px 12px;
            border: 1px solid var(--border-color);
            border-radius: 6px;
            font-size: 16px;
            line-height: 1.5;
            resize: none;
            height: 24px;
            max-height: 200px;
            overflow-y: auto;
            background-color: var(--bg-color);
        }

        .input-box:focus {
            outline: none;
            border-color: var(--primary-color);
            box-shadow: 0 0 0 2px rgba(16, 163, 127, 0.2);
        }

        .send-button {
            position: absolute;
            right: 12px;
            bottom: 12px;
            background: none;
            border: none;
            color: var(--primary-color);
            cursor: pointer;
            font-size: 18px;
            display: flex;
            align-items: center;
            justify-content: center;
            width: 32px;
            height: 32px;
            border-radius: 6px;
        }

        .send-button:hover {
            background-color: var(--hover-color);
        }

        .send-button:disabled {
            color: var(--border-color);
            cursor: not-allowed;
        }

        .welcome-container {
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            height: 100%;
            padding: 20px;
            text-align: center;
        }

        .welcome-title {
            font-size: 32px;
            font-weight: bold;
            margin-bottom: 20px;
        }

        .examples-container {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 16px;
            width: 100%;
            max-width: 900px;
        }

        .example-card {
            background-color: var(--secondary-bg);
            border-radius: 8px;
            padding: 16px;
            cursor: pointer;
            transition: background-color 0.3s;
        }

        .example-card:hover {
            background-color: var(--hover-color);
        }

        .typing-indicator {
            display: inline-flex;
            align-items: center;
            gap: 5px;
        }

        .typing-dot {
            width: 8px;
            height: 8px;
            background-color: var(--text-color);
            border-radius: 50%;
            opacity: 0.6;
            animation: typing-dot 1.4s infinite ease-in-out;
        }

        .typing-dot:nth-child(1) {
            animation-delay: 0s;
        }

        .typing-dot:nth-child(2) {
            animation-delay: 0.2s;
        }

        .typing-dot:nth-child(3) {
            animation-delay: 0.4s;
        }
	/* 保留原始樣式，只添加必要的修改 */
        .model-dropdown-container {
            margin-bottom: 15px;
        }

        .typing-indicator {
            display: none;
            padding: 8px;
            background-color: #f1f1f1;
            border-radius: 5px;
            margin-bottom: 10px;
            font-style: italic;
            color: #666;
        }

        @keyframes typing-dot {
            0%, 100% {
                transform: scale(0.7);
            }
            50% {
                transform: scale(1);
            }
        }

        /* Responsive design */
        @media (max-width: 768px) {
            .sidebar {
                display: none;
            }
            
            .examples-container {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
   <div class="sidebar">
      <!-- 新增：將下拉選單移到新對話框上方 -->
      <div class="model-select">
         <select id="model-dropdown">
            <option value="gpt-4">GPT-4</option>
            <option value="gpt-3.5-turbo">GPT-3.5 Turbo</option>
            <option value="claude-3-opus">Claude 3 Opus</option>
            <option value="claude-3-sonnet">Claude 3 Sonnet</option>
            <option value="llama-3">Llama 3</option>
         </select>
      </div>

      <div class="new-chat-btn" hx-post="/aiui/new-chat" hx-trigger="click" hx-target=".chat-area" hx-swap="innerHTML">
         <i class="fas fa-plus"></i>
         <span>新對話</span>
      </div>
        <div class="history-list">
            <!-- 聊天歷史會在這裡顯示 -->
        </div>
        <div class="sidebar-footer">
            <span>AI 聊天系統</span>
        </div>
    </div>
    
    <div class="main-content">
        <div class="chat-area">
            <div class="message-container">
                <!-- 保留原始歡迎訊息 -->
                <div class="message ai-message">
		   <div class="message-content">您好！我是AI助手，有什麼我可以幫您的嗎？</div>
               </div>
            </div>
            <!-- 新增：打字指示器 -->
            <div class="typing-indicator">AI正在思考回應...</div>
        </div>
        <div id="chat-container" class="chat-history">
            <div class="welcome-container">
                <h1 class="welcome-title">AI 聊天系統</h1>
                <div class="examples-container">
                    <div class="example-card" hx-post="/api/chat/example" hx-vals='{"prompt": "用簡單的語言解釋量子計算"}'>
                        <h3>用簡單的語言解釋量子計算</h3>
                        <p>解釋複雜的科學概念</p>
                    </div>
                    <div class="example-card" hx-post="/api/chat/example" hx-vals='{"prompt": "給我寫一個Python的網絡爬蟲程式"}'>
                        <h3>給我寫一個Python的網絡爬蟲程式</h3>
                        <p>幫助編寫程式碼</p>
                    </div>
                    <div class="example-card" hx-post="/api/chat/example" hx-vals='{"prompt": "創作一個短篇科幻故事"}'>
                        <h3>創作一個短篇科幻故事</h3>
                        <p>生成創意內容</p>
                    </div>
                </div>
            </div>
        </div>
        
        <div class="input-area">
           <form class="input-container" hx-post="/aiui/send-message" hx-trigger="submit" hx-ext="sse-events" hx-target="#chat-container" hx-swap="beforeend">
                <input type="hidden" name="model" id="model-input">
                <textarea id="user-input"
                    class="input-box" 
                    name="message" 
                    placeholder="傳送訊息..." 
                    required
                    onkeydown="if(event.key === 'Enter' && !event.shiftKey) { event.preventDefault(); this.form.dispatchEvent(new Event('submit')); }"
                ></textarea>
                <button type="submit" class="send-button">
                    <i class="fas fa-paper-plane"></i>
                </button>
            </form>
        </div>
    </div>

    <script>
        // 調整textarea高度
        document.querySelector('.input-box').addEventListener('input', function() {
            this.style.height = 'auto';
            this.style.height = (this.scrollHeight) + 'px';
        });

        // 將選擇的模型值傳遞到隱藏輸入字段
        document.getElementById('model-dropdown').addEventListener('change', function() {
            document.getElementById('model-input').value = this.value;
        });

        // 初始化隱藏字段值
        document.getElementById('model-input').value = document.getElementById('model-dropdown').value;

        // SSE連接處理
        function setupSSEConnection(messageId) {
            const eventSource = new EventSource(`/aiui/${messageId}`);
            const messageContent = document.querySelector(`#message-${messageId} .message-content`);
            
            eventSource.onmessage = function(event) {
                const data = JSON.parse(event.data);
                
                if (data.content) {
                    messageContent.innerHTML = data.content;
                }
                
                if (data.done) {
                    eventSource.close();
                }
            };
            
            eventSource.onerror = function() {
                eventSource.close();
            };
        }

        // HTMX事件處理器
        document.body.addEventListener('htmx:afterSwap', function(event) {
            // 如果新增了訊息元素並且有特定的ID
            const newElements = event.detail.target.querySelectorAll('[id^="message-"]');
            newElements.forEach(element => {
                const messageId = element.id.replace('message-', '');
                if(element.classList.contains('ai-message')) {
                    setupSSEConnection(messageId);
                }
            });
        });
       // 添加使用者訊息
       function addUserMessage(text) {
          const messageDiv = document.createElement('div');
          messageDiv.className = 'message user-message';
                
          const contentDiv = document.createElement('div');
          contentDiv.className = 'message-content';
          contentDiv.textContent = text;
                
          messageDiv.appendChild(contentDiv);
          chatContainer.appendChild(messageDiv);
                
          // 滾動到底部
          chatContainer.scrollTop = chatContainer.scrollHeight;
       }
       // 初始化HTMX的SSE擴展
        document.addEventListener('DOMContentLoaded', function() {
            htmx.defineExtension('sse-events', {
                onEvent: function(name, event) {
				    console.log(name);
                    // 處理SSE事件
                }
            });
        });
       // 提交表單時處理用戶消息並連接SSE
        document.querySelector('form').addEventListener('submit', function(e) {
            e.preventDefault();
            
            const userInput = document.getElementById('user-input');
            const messageContainer = document.querySelector('.message-container');
            const typingIndicator = document.querySelector('.typing-indicator');
            
            // 顯示用戶消息
            const userMessage = document.createElement('div');
            userMessage.className = 'message user-message';
            userMessage.textContent = userInput.value;
            messageContainer.appendChild(userMessage);
            
            // 顯示打字指示器
            typingIndicator.style.display = 'block';
            
            // 連接到SSE端點
            const selectedModel = document.getElementById('model-dropdown').value;
            const eventSource = new EventSource(`/aiui/sse?message=${encodeURIComponent(userInput.value)}&model=${encodeURIComponent(selectedModel)}`);
            
            let currentAiMessage = null;
            
            eventSource.onmessage = function(event) {
                const data = JSON.parse(event.data);
                if (data.type === 'complete') {  // 完成時關閉連接並隱藏指示器
                    eventSource.close();
                    typingIndicator.style.display = 'none';
                    if(currentAiMessage) {
                       currentAiMessage.dataset.complete = 'true';
                    }
                } else if(data.type === 'chunk') {  // 第一個塊時創建新消息元素
                    if (!currentAiMessage) {
                        currentAiMessage = document.createElement('div');
                        currentAiMessage.className = 'message ai-message';
                        currentAiMessage.dataset.complete = 'false';
                        messageContainer.appendChild(currentAiMessage);
                    }
                    
                    // 添加新的文本
                    currentAiMessage.textContent += data.content;
                    
                    // 滾動到最新消息
                    messageContainer.scrollIntoView({ behavior: 'smooth', block: 'end' });
                }
            };
            
            // 錯誤處理
            eventSource.onerror = function() {
                console.error('SSE連接錯誤');
                eventSource.close();
                typingIndicator.style.display = 'none';
                
                // 顯示錯誤消息
                const errorMessage = document.createElement('div');
                errorMessage.className = 'message ai-message error';
                errorMessage.textContent = '連接出現問題，請稍後再試。';
                messageContainer.appendChild(errorMessage);
            };
            
            // 清空輸入框
            userInput.value = '';
        });
    </script>
</body>
</html>
