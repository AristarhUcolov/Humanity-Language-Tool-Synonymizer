// ============================================================
// UI translations for English, Russian, Moldovan
// ============================================================
const I18N = {
  en: {
    tagline: "Multi-language vocabulary elevation. Local, offline, no external APIs.",
    ui_lang_label: "UI:",
    docx_title: "Synonymize Word Document (.docx)",
    docx_desc: "Drag & drop a .docx file here or click to choose — fonts, tables, styles all preserved.",
    docx_choose: "Choose .docx…",
    original_text: "Original Text",
    elevated_text: "Elevated Text",
    paste_btn: "Paste",
    clear_btn: "Clear",
    copy_btn: "Copy",
    download_btn: "Download .txt",
    run_btn: "Synonymize",
    lang_label: "Synonym Language",
    lang_auto: "Auto-detect",
    hint_shortcut: "Tip: Ctrl + Enter",
    input_placeholder: "Paste your text here… (English, Russian, Moldovan)",
    output_placeholder: "Synonymized text will appear here…",
    diff_title: "Changes",
    diff_hide: "Hide",
    diff_show: "Show",
    footer_info: "Runs entirely offline on your computer. ~15,000 synonyms across English, Russian, and Moldovan.",
    footer_author: "Author:",
    donate_label: "Support the project:",
    msg_paste_first: "Paste some text first.",
    msg_processing: "Processing…",
    msg_done: "Done",
    msg_pasted: "Pasted",
    msg_copied: "Copied to clipboard",
    msg_failed: "Failed: ",
    msg_empty_file: "File is empty.",
    msg_not_docx: "Not a .docx file.",
    msg_doc_processing: "Processing ",
    msg_ready: "Ready: ",
    lang_detected: "Detected language: ",
  },
  ru: {
    tagline: "Многоязычное обогащение словаря. Локально, оффлайн, без внешних API.",
    ui_lang_label: "Интерфейс:",
    docx_title: "Синонимизировать документ Word (.docx)",
    docx_desc: "Перетащите .docx файл сюда или нажмите для выбора — шрифты, таблицы, стили сохраняются.",
    docx_choose: "Выбрать .docx…",
    original_text: "Исходный текст",
    elevated_text: "Обогащённый текст",
    paste_btn: "Вставить",
    clear_btn: "Очистить",
    copy_btn: "Копировать",
    download_btn: "Скачать .txt",
    run_btn: "Синонимизировать",
    lang_label: "Язык синонимов",
    lang_auto: "Автоопределение",
    hint_shortcut: "Совет: Ctrl + Enter",
    input_placeholder: "Вставьте текст здесь… (Английский, Русский, Молдавский)",
    output_placeholder: "Синонимизированный текст появится здесь…",
    diff_title: "Изменения",
    diff_hide: "Скрыть",
    diff_show: "Показать",
    footer_info: "Работает полностью оффлайн. ~15 000 синонимов на английском, русском и молдавском.",
    footer_author: "Автор:",
    donate_label: "Поддержать проект:",
    msg_paste_first: "Сначала вставьте текст.",
    msg_processing: "Обработка…",
    msg_done: "Готово",
    msg_pasted: "Вставлено",
    msg_copied: "Скопировано в буфер",
    msg_failed: "Ошибка: ",
    msg_empty_file: "Файл пустой.",
    msg_not_docx: "Не .docx файл.",
    msg_doc_processing: "Обработка ",
    msg_ready: "Готово: ",
    lang_detected: "Определён язык: ",
  },
  mo: {
    tagline: "Îmbogățirea vocabularului în mai multe limbi. Local, offline, fără API externe.",
    ui_lang_label: "Interfață:",
    docx_title: "Sinonimizează document Word (.docx)",
    docx_desc: "Trage fișierul .docx aici sau apasă pentru a alege — fonturi, tabele, stiluri păstrate.",
    docx_choose: "Alege .docx…",
    original_text: "Text original",
    elevated_text: "Text îmbogățit",
    paste_btn: "Lipește",
    clear_btn: "Șterge",
    copy_btn: "Copiază",
    download_btn: "Descarcă .txt",
    run_btn: "Sinonimizează",
    lang_label: "Limba sinonimelor",
    lang_auto: "Auto-detectare",
    hint_shortcut: "Sfat: Ctrl + Enter",
    input_placeholder: "Lipește textul aici… (Engleză, Rusă, Moldovenească)",
    output_placeholder: "Textul sinonimizat va apărea aici…",
    diff_title: "Schimbări",
    diff_hide: "Ascunde",
    diff_show: "Arată",
    footer_info: "Rulează complet offline. ~15.000 sinonime în engleză, rusă și moldovenească.",
    footer_author: "Autor:",
    donate_label: "Susține proiectul:",
    msg_paste_first: "Lipește mai întâi un text.",
    msg_processing: "Procesare…",
    msg_done: "Gata",
    msg_pasted: "Lipit",
    msg_copied: "Copiat în clipboard",
    msg_failed: "Eșuat: ",
    msg_empty_file: "Fișierul este gol.",
    msg_not_docx: "Nu este un fișier .docx.",
    msg_doc_processing: "Procesare ",
    msg_ready: "Gata: ",
    lang_detected: "Limbă detectată: ",
  },
};

const LANG_NAMES = {
  "auto": "Auto-detect",
  "en": "English",
  "ru": "Русский",
  "mo": "Moldovenească",
};

let currentUiLang = localStorage.getItem('uiLang') || 'en';

// ============================================================
// Theme (light / dark)
// ============================================================
function applyTheme(theme) {
  document.documentElement.setAttribute('data-theme', theme);
  const btn = document.getElementById('theme-toggle');
  if (btn) btn.textContent = theme === 'light' ? '☀️' : '🌙';
  try { localStorage.setItem('theme', theme); } catch (_) {}
}

function toggleTheme() {
  const cur = document.documentElement.getAttribute('data-theme') || 'dark';
  applyTheme(cur === 'dark' ? 'light' : 'dark');
}

function t(key) {
  return (I18N[currentUiLang] && I18N[currentUiLang][key]) || I18N.en[key] || key;
}

function applyTranslations() {
  document.documentElement.lang = currentUiLang === 'mo' ? 'ro' : currentUiLang;
  document.querySelectorAll('[data-i18n]').forEach(el => {
    const key = el.getAttribute('data-i18n');
    el.textContent = t(key);
  });
  document.querySelectorAll('[data-i18n-placeholder]').forEach(el => {
    const key = el.getAttribute('data-i18n-placeholder');
    el.placeholder = t(key);
  });
}

// ============================================================
// Text utilities
// ============================================================
function countWords(s) {
  const m = s.trim().match(/\S+/g);
  return m ? m.length : 0;
}

function updateStats() {
  const inp = document.getElementById("input").value;
  const out = document.getElementById("output").value;
  document.getElementById("input-stats").textContent = `${inp.length} ${currentUiLang === 'ru' ? 'символов' : currentUiLang === 'mo' ? 'caractere' : 'characters'} · ${countWords(inp)} ${currentUiLang === 'ru' ? 'слов' : currentUiLang === 'mo' ? 'cuvinte' : 'words'}`;
  document.getElementById("output-stats").textContent = `${out.length} ${currentUiLang === 'ru' ? 'символов' : currentUiLang === 'mo' ? 'caractere' : 'characters'} · ${countWords(out)} ${currentUiLang === 'ru' ? 'слов' : currentUiLang === 'mo' ? 'cuvinte' : 'words'}`;
}

// ============================================================
// Diff highlighting
// ============================================================
function computeDiff(original, modified) {
  // Simple word-level diff. For each word in modified that differs from
  // the corresponding word in original, show as a change pair.
  const origWords = original.match(/\S+/g) || [];
  const modWords = modified.match(/\S+/g) || [];
  const pairs = [];
  const len = Math.min(origWords.length, modWords.length);

  for (let i = 0; i < len; i++) {
    const a = origWords[i].replace(/[.,!?;:()«»"']/g, '');
    const b = modWords[i].replace(/[.,!?;:()«»"']/g, '');
    if (a.toLowerCase() !== b.toLowerCase() && a && b) {
      // Avoid duplicates
      if (!pairs.some(p => p.from === a && p.to === b)) {
        pairs.push({ from: a, to: b });
      }
    }
  }
  return pairs;
}

function renderDiff(pairs) {
  const section = document.getElementById('diff-section');
  const content = document.getElementById('diff-content');
  content.innerHTML = '';
  if (pairs.length === 0) {
    section.style.display = 'none';
    return;
  }
  section.style.display = 'block';
  for (const pair of pairs) {
    const div = document.createElement('div');
    div.className = 'diff-pair';
    div.innerHTML = `<span class="from">${escapeHtml(pair.from)}</span><span class="arrow">→</span><span class="to">${escapeHtml(pair.to)}</span>`;
    content.appendChild(div);
  }
}

function escapeHtml(s) {
  const div = document.createElement('div');
  div.textContent = s;
  return div.innerHTML;
}

// ============================================================
// Main actions
// ============================================================
async function run() {
  const text = document.getElementById("input").value;
  const language = document.getElementById("language").value;
  const status = document.getElementById("status");
  const detected = document.getElementById("detected-lang");
  const btn = document.getElementById("run-btn");
  const changes = document.getElementById("changes");

  if (!text.trim()) {
    status.textContent = t('msg_paste_first');
    status.className = "status err";
    return;
  }

  btn.disabled = true;
  status.textContent = t('msg_processing');
  status.className = "status";
  changes.textContent = "";
  detected.textContent = "";

  try {
    const resp = await fetch("/api/humanize", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ text, language }),
    });
    if (!resp.ok) throw new Error(await resp.text());
    const data = await resp.json();
    document.getElementById("output").value = data.output;
    updateStats();
    status.textContent = t('msg_done');
    status.className = "status ok";

    if (data.detectedLanguage) {
      const langDisplay = LANG_NAMES[data.detectedLanguage] || data.detectedLanguage;
      detected.textContent = `${t('lang_detected')}${langDisplay}`;
    }

    if (Array.isArray(data.notes) && data.notes.length) {
      changes.textContent = data.notes.join(" · ");
    }

    // Compute and render diff
    const pairs = computeDiff(text, data.output);
    renderDiff(pairs);

    // Save to localStorage
    try {
      localStorage.setItem('lastInput', text);
      localStorage.setItem('lastOutput', data.output);
    } catch (_) {}
  } catch (e) {
    status.textContent = t('msg_failed') + e.message;
    status.className = "status err";
  } finally {
    btn.disabled = false;
  }
}

async function paste() {
  try {
    const text = await navigator.clipboard.readText();
    if (text) {
      document.getElementById("input").value = text;
      updateStats();
      const status = document.getElementById("status");
      status.textContent = t('msg_pasted');
      status.className = "status ok";
    }
  } catch (_) {
    document.getElementById("input").focus();
  }
}

async function copyOut() {
  const text = document.getElementById("output").value;
  if (!text) return;
  try {
    await navigator.clipboard.writeText(text);
    const status = document.getElementById("status");
    status.textContent = t('msg_copied');
    status.className = "status ok";
  } catch (_) {
    const ta = document.getElementById("output");
    ta.select();
    document.execCommand("copy");
  }
}

function downloadTxt() {
  const text = document.getElementById("output").value;
  if (!text) return;
  const blob = new Blob([text], { type: 'text/plain;charset=utf-8' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `synonymized_${Date.now()}.txt`;
  document.body.appendChild(a);
  a.click();
  a.remove();
  setTimeout(() => URL.revokeObjectURL(url), 1000);
}

async function uploadDocx(file) {
  const status = document.getElementById("docx-status");
  const language = document.getElementById("language").value;
  if (!file) return;
  if (file.size === 0) {
    status.textContent = t('msg_empty_file');
    status.className = "status err";
    return;
  }
  if (!file.name.toLowerCase().endsWith(".docx")) {
    status.textContent = t('msg_not_docx');
    status.className = "status err";
    return;
  }

  status.textContent = `${t('msg_doc_processing')}"${file.name}"…`;
  status.className = "status";

  const fd = new FormData();
  fd.append("file", file, file.name);
  fd.append("language", language);

  try {
    const resp = await fetch("/api/humanize-docx", { method: "POST", body: fd });
    if (!resp.ok) {
      const txt = await resp.text();
      throw new Error(txt || resp.statusText);
    }
    const total = parseInt(resp.headers.get("X-Humanity-Runs-Total") || "0", 10);
    const changed = parseInt(resp.headers.get("X-Humanity-Runs-Changed") || "0", 10);
    const synonyms = parseInt(resp.headers.get("X-Humanity-Synonyms") || "0", 10);
    const lang = resp.headers.get("X-Humanity-Language") || "";
    const langDisplay = LANG_NAMES[lang] || lang;

    const blob = await resp.blob();
    const cd = resp.headers.get("Content-Disposition") || "";
    const m = /filename="?([^";]+)"?/i.exec(cd);
    const outName = (m && m[1]) || file.name.replace(/\.docx$/i, ".synonymized.docx");

    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = outName;
    document.body.appendChild(a);
    a.click();
    a.remove();
    setTimeout(() => URL.revokeObjectURL(url), 5000);

    status.textContent = `${t('msg_ready')}"${outName}" [${langDisplay}] (${changed}/${total} runs, ${synonyms} synonyms)`;
    status.className = "status ok";
  } catch (e) {
    status.textContent = t('msg_failed') + e.message;
    status.className = "status err";
  }
}

// ============================================================
// Drag and drop for DOCX
// ============================================================
function setupDragDrop() {
  const dropzone = document.getElementById('docx-dropzone');
  ['dragenter', 'dragover'].forEach(ev => {
    dropzone.addEventListener(ev, (e) => {
      e.preventDefault();
      e.stopPropagation();
      dropzone.classList.add('drag-over');
    });
  });
  ['dragleave', 'drop'].forEach(ev => {
    dropzone.addEventListener(ev, (e) => {
      e.preventDefault();
      e.stopPropagation();
      dropzone.classList.remove('drag-over');
    });
  });
  dropzone.addEventListener('drop', (e) => {
    const files = e.dataTransfer.files;
    if (files && files.length > 0) {
      uploadDocx(files[0]);
    }
  });
}

// ============================================================
// Init
// ============================================================
window.addEventListener("DOMContentLoaded", () => {
  // Apply saved theme (default dark)
  applyTheme(localStorage.getItem('theme') || 'dark');
  document.getElementById("theme-toggle").addEventListener("click", toggleTheme);

  // Apply initial translations
  document.getElementById("ui-lang-select").value = currentUiLang;
  applyTranslations();

  document.getElementById("ui-lang-select").addEventListener("change", (e) => {
    currentUiLang = e.target.value;
    localStorage.setItem('uiLang', currentUiLang);
    applyTranslations();
    updateStats();
  });

  document.getElementById("run-btn").addEventListener("click", run);
  document.getElementById("paste-btn").addEventListener("click", paste);
  document.getElementById("copy-btn").addEventListener("click", copyOut);
  document.getElementById("download-btn").addEventListener("click", downloadTxt);
  document.getElementById("clear-btn").addEventListener("click", () => {
    document.getElementById("input").value = "";
    document.getElementById("output").value = "";
    document.getElementById("changes").textContent = "";
    document.getElementById("detected-lang").textContent = "";
    document.getElementById("diff-section").style.display = 'none';
    updateStats();
    try { localStorage.removeItem('lastInput'); localStorage.removeItem('lastOutput'); } catch (_) {}
  });

  document.getElementById("input").addEventListener("input", updateStats);

  document.getElementById("docx-file").addEventListener("change", (e) => {
    const f = e.target.files && e.target.files[0];
    if (f) uploadDocx(f);
    e.target.value = ""; // allow re-upload of same file
  });

  // Diff toggle
  document.getElementById("diff-toggle").addEventListener("click", () => {
    const content = document.getElementById("diff-content");
    const btn = document.getElementById("diff-toggle");
    if (content.style.display === 'none') {
      content.style.display = 'flex';
      btn.textContent = t('diff_hide');
    } else {
      content.style.display = 'none';
      btn.textContent = t('diff_show');
    }
  });

  // Keyboard shortcut: Ctrl+Enter to run
  document.addEventListener("keydown", (e) => {
    if ((e.ctrlKey || e.metaKey) && e.key === "Enter") run();
  });

  // Drag and drop
  setupDragDrop();

  // Restore from localStorage
  try {
    const lastInput = localStorage.getItem('lastInput');
    if (lastInput) document.getElementById('input').value = lastInput;
    const lastOutput = localStorage.getItem('lastOutput');
    if (lastOutput) document.getElementById('output').value = lastOutput;
  } catch (_) {}

  updateStats();
});
