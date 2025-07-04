// 滾動到底部
function scrollToBottom(idName) {
   const messagesContainer = document.getElementById(idName);
   messagesContainer.scrollTop = messagesContainer.scrollHeight;
}

// 格式化文本（Markdown-like）
function formatText(text) {
   return text
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
    .replace(/`([^`]+)`/g, '<code>$1</code>')
    .replace(/```(\w+)?\n([\s\S]*?)```/g, '<pre><code>$2</code></pre>')
    .replace(/^> (.+)$/gm, '<blockquote>$1</blockquote>')
    .replace(/^• (.+)$/gm, '<li>$1</li>')
    .replace(/(<li>.*<\/li>)/s, '<ul>$1</ul>')
    .replace(/^\d+\. (.+)$/gm, '<li>$1</li>')
    .replace(/^#{1,3} (.+)$/gm, function(match, text) {
        const level = match.indexOf(' ');
        return `<h${level}>${text}</h${level}>`;
    })
    .replace(/\n\n/g, '</p><p>')
    .replace(/^(.+)$/gm, function(match) {
        if (!match.startsWith('<') && match.trim()) {
            return `<p>${match}</p>`;
        }
        return match;
    });
}

 // HTML 轉義
function escapeHtml(text) {
   const div = document.createElement('div');
   div.textContent = text;
   return div.innerHTML;
}