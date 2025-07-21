/*
需要搭配下列HTML片段使用：
<button class="send-button" id="mic-button"><i class="fa-solid fa-microphone"></i></button>
*/

// 語音識別相關變數
var recognition;
var isRecognizing = false;
var finalTranscript = '';

// 檢查瀏覽器是否支援語音識別
function InitializeSpeechRecognition() {
    if (!('webkitSpeechRecognition' in window)) {
        // document.getElementById('status').innerHTML = '您的瀏覽器不支援語音識別，請使用 Chrome 瀏覽器';
        document.getElementById('mic-button').style.display = 'none';
    } else {
        // 初始化語音識別
        recognition = new webkitSpeechRecognition();
        recognition.continuous = true;
        recognition.interimResults = true;
        recognition.lang = 'cmn-Hant-TW';
        
        // 語音識別開始
        recognition.onstart = function() {
            isRecognizing = true;
            // document.getElementById('mic-button').style.display = 'none';
            // document.getElementById('status').innerHTML = '正在聆聽中...';
            // document.getElementById('mic-image').src = 'https://www.google.com/intl/en/chrome/assets/common/images/content/mic-animate.gif';
        };
        
        // 語音識別結束
        recognition.onend = function() {
            isRecognizing = false;
            // document.getElementById('status').innerHTML = '點擊麥克風開始語音輸入';
            // document.getElementById('mic-image').src = 'https://www.google.com/intl/z/none/images/logo-mic.png';
            if(finalTranscript) {
                // document.getElementById('status').innerHTML = '語音輸入完成！';
                // document.getElementById('mic-button').classList.remove('active');
            }
        };
        
        // 語音識別結果
        recognition.onresult = function(event) {
            let interimTranscript = '';
            
            for(let i = event.resultIndex; i < event.results.length; i++)  {
                if(event.results[i].isFinal) {
                    finalTranscript += event.results[i][0].transcript;
                } else {
                    interimTranscript += event.results[i][0].transcript;
                }
            }            
            // 顯示結果 id="message-input"
            document.getElementById('message-input').value = finalTranscript;
            // document.getElementById('final-text').innerHTML = finalTranscript;
            // document.getElementById('interim-text').innerHTML = interimTranscript;
        };
        
        // 語音識別錯誤處理
        recognition.onerror = function(event) {
            isRecognizing = false;
            // document.getElementById('mic-button').classList.remove('active');
            // document.getElementById('mic-image').src = 'https://www.google.com/intl/en/chrome/assets/common/images/content/mic.gif';
            switch(event.error) {
                case 'no-speech':
                    document.getElementById('message-input').value = "";
                    // document.getElementById('status').innerHTML = '沒有檢測到語音';
                    break;
                case 'audio-capture':
                    // document.getElementById('status').innerHTML = '沒有找到麥克風';
                    console.error('沒有找到麥克風');
                    break;
                case 'not-allowed':
                    console.error('沒有麥克風使用權限');
                    document.getElementById('status').innerHTML = '沒有麥克風使用權限';
                    break;
                default:
                    console.error('語音識別發生錯誤: ' + event.error);
                    document.getElementById('status').innerHTML = '語音識別發生錯誤';
            }
        };
    }

    // 切換語音識別狀態
    recognition.toggleRecognition = function() {
        if(isRecognizing) {
            recognition.stop();
        } else {
            finalTranscript = '';
            document.getElementById('message-input').value = "";
            // document.getElementById('final-text').innerHTML = '';
            // document.getElementById('interim-text').innerHTML = '';
            // 設定選擇的語言
            recognition.lang = "cmn-Hant-TW";  // document.getElementById('language-select').value;
            // document.getElementById('mic-button').style.backgroundColor = '#4facfe'; // 按下時的藍色             
            recognition.start();
        }
    }
}