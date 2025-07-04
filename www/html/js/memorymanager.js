// 記憶管理類
class MemoryManager {
   constructor() {
      this.currentSession = {
         id: this.generateSessionId(),
         messages: [],
         startTime: new Date().toISOString(),
         preferences: this.loadPreferences()
      };
      this.init();
   }
   
   init() {
      this.loadSession();
      this.updateMemoryDisplay();
      this.loadPreferencesUI();
   }
   
   generateSessionId() {
      return 'session_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
   }
   
   // 消息管理
   addMessage(role, content) {
      const message = {
         id: Date.now(),
         role: role,
         content: content,
         timestamp: new Date().toISOString()
      };
      this.currentSession.messages.push(message);
      this.saveSession();
      this.updateMemoryDisplay();
      return message;
   }
   
   // 會話管理
   saveSession() {
         const sessionData = JSON.stringify(this.currentSession);  // 保存當前會話
         window.memoryStorage = window.memoryStorage || {};  // 使用變數代替 localStorage (因為在 Claude 環境中不可用)
         window.memoryStorage[this.currentSession.id] = sessionData;
         const sessions = this.getSessions();
         const existingIndex = sessions.findIndex(s => s.id === this.currentSession.id);
         
         if (existingIndex !== -1) {
            sessions[existingIndex] = {
               id: this.currentSession.id,
               startTime: this.currentSession.startTime,
               messageCount: this.currentSession.messages.length,
               lastMessage: this.currentSession.messages[this.currentSession.messages.length - 1]?.content || ''
            };
         } else {
            sessions.push({
               id: this.currentSession.id,
               startTime: this.currentSession.startTime,
               messageCount: this.currentSession.messages.length,
               lastMessage: this.currentSession.messages[this.currentSession.messages.length - 1]?.content || ''
            });
         }
         window.memoryStorage['sessions'] = JSON.stringify(sessions);
   }
   
   loadSession() {
      window.memoryStorage = window.memoryStorage || {};
      const sessionData = window.memoryStorage[this.currentSession.id];
      if (sessionData) {
         this.currentSession = JSON.parse(sessionData);
      }
   }
   
   getSessions() {
      window.memoryStorage = window.memoryStorage || {};
      const sessionsData = window.memoryStorage['sessions'];
      return sessionsData ? JSON.parse(sessionsData) : [];
   }
   
   // 偏好設定管理
   savePreferences(preferences) {
      this.currentSession.preferences = { ...this.currentSession.preferences, ...preferences };
      window.memoryStorage = window.memoryStorage || {};
      window.memoryStorage['preferences'] = JSON.stringify(this.currentSession.preferences);
      this.saveSession();
   }
   
   loadPreferences() {
      window.memoryStorage = window.memoryStorage || {};
      const preferencesData = window.memoryStorage['preferences'];
      return preferencesData ? JSON.parse(preferencesData) : {
         theme: 'blue',
         language: 'zh-TW',
         autoSave: true,
         notifications: true
      };
   }
   
   // UI 更新
   updateMemoryDisplay() {
/*      
      const currentSessionEl = document.getElementById('currentSession');   // 更新當前會話顯示
      currentSessionEl.innerHTML = `
         <div class="memory-item">
            <strong>會話ID:</strong> ${this.currentSession.id}
         </div>
         <div class="memory-item">
            <strong>開始時間:</strong> ${new Date(this.currentSession.startTime).toLocaleString()}
         </div>
         <div class="memory-item">
            <strong>消息數量:</strong> ${this.currentSession.messages.length}
         </div>
      `;
*/      
      // 更新歷史記錄顯示
      // const historyListEl = document.getElementById('historyList');
      const historyListEl = document.getElementById('talkNumber');
      const sessions = this.getSessions();
      
      if(sessions.length === 0) {
         historyListEl.innerHTML = '0';  // '<div class="memory-item">暫無歷史記錄</div>';
         return;
      } else {
         historyListEl.innerHTML = sessions.slice(-5).map(session => `${session.messageCount}`).join('');;  // 更新對話數量顯示
         /*
         historyListEl.innerHTML = sessions.slice(-5).map(session => `
            <div class="memory-item">
                  <strong>會話:</strong> ${new Date(session.startTime).toLocaleString()}<br>
                  <strong>消息數:</strong> ${session.messageCount}<br>
                  <strong>最後消息:</strong> ${session.lastMessage.substring(0, 50)}...
            </div>
         `).join('');
         */
      }
/*      
      // 更新用戶偏好顯示
      const userPreferencesEl = document.getElementById('userPreferences');
      userPreferencesEl.innerHTML = `
         <div class="memory-item">
            <strong>主題:</strong> ${this.currentSession.preferences.theme}
         </div>
         <div class="memory-item">
            <strong>語言:</strong> ${this.currentSession.preferences.language}
         </div>
         <div class="memory-item">
            <strong>自動保存:</strong> ${this.currentSession.preferences.autoSave ? '是' : '否'}
         </div>
         <div class="memory-item">
            <strong>通知:</strong> ${this.currentSession.preferences.notifications ? '是' : '否'}
         </div>
      `;
*/      
   }
   
   loadPreferencesUI() {
/*      
      const prefs = this.currentSession.preferences;
      document.getElementById('themeColor').value = prefs.theme;
      document.getElementById('language').value = prefs.language;
      document.getElementById('autoSave').checked = prefs.autoSave;
      document.getElementById('notifications').checked = prefs.notifications;
*/      
   }
   
   // 記憶操作
   clearMemory() {
      window.memoryStorage = {};
      this.currentSession = {
         id: this.generateSessionId(),
         messages: [],
         startTime: new Date().toISOString(),
         preferences: this.loadPreferences()
      };
      // document.getElementById('chatMessages').innerHTML = '';
      this.updateMemoryDisplay();
   }
   
   exportMemory() {
         const exportData = {
            version: '1.0',
            exportTime: new Date().toISOString(),
            sessions: this.getSessions(),
            currentSession: this.currentSession,
            preferences: this.currentSession.preferences
         };
         
         const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' });
         const url = URL.createObjectURL(blob);
         const a = document.createElement('a');
         a.href = url;
         a.download = `ai_memory_${new Date().toISOString().split('T')[0]}.json`;
         a.click();
         URL.revokeObjectURL(url);
   }
   
   importMemory() {
         const input = document.createElement('input');
         input.type = 'file';
         input.accept = '.json';
         input.onchange = (e) => {
            const file = e.target.files[0];
            if (file) {
               const reader = new FileReader();
               reader.onload = (e) => {
                     try {
                        const importData = JSON.parse(e.target.result);
                        
                        // 恢復數據
                        window.memoryStorage = window.memoryStorage || {};
                        
                        if (importData.sessions) {
                           window.memoryStorage['sessions'] = JSON.stringify(importData.sessions);
                        }
                        if (importData.preferences) {
                           window.memoryStorage['preferences'] = JSON.stringify(importData.preferences);
                        }
                        if (importData.currentSession) {
                           this.currentSession = importData.currentSession;
                           this.loadSession();
                           this.displayMessages();
                        }
                        this.updateMemoryDisplay();
                        this.loadPreferencesUI();
                        
                        alert('記憶匯入成功！');
                     } catch (error) {
                        alert('匯入失敗：' + error.message);
                     }
               };
               reader.readAsText(file);
            }
         };
         input.click();
   }
   
   displayMessages() {
      /*
      const chatMessages = document.getElementById('chatMessages');
      chatMessages.innerHTML = '';
      
      this.currentSession.messages.forEach(message => {
         const messageEl = document.createElement('div');
         messageEl.className = `message ${message.role}`;
         messageEl.innerHTML = `
            <div>${message.content}</div>
            <div class="message-time">${new Date(message.timestamp).toLocaleString()}</div>
         `;
         chatMessages.appendChild(messageEl);
      });
      chatMessages.scrollTop = chatMessages.scrollHeight; 
   */
   }
}