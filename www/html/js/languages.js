let currentLang = 'zh';

// 目前僅支援中英文切換
function switchLanguage(lang) {
   currentLang = lang;
   
   // Update button states
   document.querySelectorAll('.lang-btn').forEach(btn => {
         btn.classList.remove('active');
   });
   event.target.classList.add('active');

   // Update text content
   document.querySelectorAll('[data-zh][data-en]').forEach(element => {
         if (lang === 'zh') {
            element.innerHTML = element.getAttribute('data-zh');
         } else {
            element.innerHTML = element.getAttribute('data-en');
         }
   });
}