# Humanity Language Tool — Synonymizer

<div align="center">

**Multi-language vocabulary elevation tool with .docx support**

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Languages](https://img.shields.io/badge/Languages-EN%20%7C%20RU%20%7C%20MO-blue)]()
[![Offline](https://img.shields.io/badge/Offline-100%25-success)]()

[**🇬🇧 English**](#-english) · [**🇷🇺 Русский**](#-русский)

</div>

---

## 🇬🇧 English

### About <a name="english"></a>

**Humanity Language Tool — Synonymizer** is a powerful offline tool that elevates text vocabulary by replacing common words with sophisticated synonyms. Supports **English**, **Russian**, and **Moldovan/Romanian** with automatic language detection, and processes both plain text and `.docx` files while preserving all formatting (fonts, sizes, colors, tables, lists, etc.).

### ✨ Features

- 🌍 **Multi-language support**: English, Russian, Moldovan/Romanian with auto-detection
- 📚 **Huge synonym dictionaries** (~18,800 total):
  - English: **7,900+** entries
  - Russian: **5,700+** entries
  - Moldovan: **5,100+** entries
- 📄 **Full DOCX support**: preserves all formatting (cross-run synonym matching)
  - All fonts (Times New Roman, Arial, Calibri, etc.), sizes (6pt–72pt+)
  - Bold, italic, underline, strikethrough, colors (RGB)
  - Tables, lists, headers, footers, hyperlinks, comments, tracked changes
- 📑 **Batch DOCX**: process many .docx at once → one ZIP
- 📕 **PDF text extraction**: pull text from PDFs, synonymize, get .txt
- ✨ **Click-to-revert**: changed words are highlighted; click one to undo it
- 🌗 **Light / dark theme**
- 🔒 **100% offline** — no data leaves your computer
- 🌐 **Network access** — share with devices on local network
- 🎨 **Smart punctuation** — handles dashes, quotes, ellipses across languages
- ⚡ **Fast & lightweight** — single 7 MB executable

### 🚀 Quick Start

#### Download & Run
```bash
# Just download humanity.exe and run it
humanity.exe
```

#### Options
```bash
humanity.exe                        # Listen on all network interfaces (default)
humanity.exe -addr 127.0.0.1:8080   # Localhost only
humanity.exe -no-open               # Don't auto-open browser
```

When you start the server, it shows:
```
━━━ Synonymaizer Server ━━━
🌐 Network Access (from other devices):
   http://192.168.0.57:8080
Local only (localhost):   http://127.0.0.1:8080
```

Open one of those URLs in your browser.

### 📖 Usage

1. **Select language** (or use Auto-detect)
2. **Paste your text** in the left textarea
3. Click **Synonymize**
4. Copy the elevated text from the right

#### DOCX Processing
1. Click **Choose .docx…** at the top
2. Select your Word document
3. The synonymized file downloads automatically
4. All formatting is preserved exactly

### 🔧 API

#### POST /api/humanize
```json
{
  "text": "The good man works hard.",
  "language": "auto"
}
```

Response:
```json
{
  "output": "The excellent gentleman labors arduous.",
  "detectedLanguage": "en",
  "synonymsApplied": 5,
  "notes": ["Language: English", "Synonyms applied: 5"]
}
```

#### POST /api/humanize-docx
Multipart form with `file` (DOCX) and optional `language`. Returns the processed DOCX.

#### POST /api/humanize-docx-batch
Multipart form with several `files` (DOCX) and optional `language`. Returns a ZIP of all processed documents.

#### POST /api/humanize-pdf
Multipart form with `file` (PDF) and optional `language`. Extracts the PDF text, synonymizes it, returns a `.txt`.

#### GET /api/languages
Lists all supported languages.

### 🛠 Build from Source

Requirements: Go 1.21+

```bash
git clone https://github.com/aristarh-ucolov/Humanity-Language-Tool-Synonymizer
cd Humanity-Language-Tool-Synonymizer
go build -ldflags="-s -w" -o humanity.exe
```

### 📁 Project Structure

```
.
├── main.go                          # HTTP server
├── internal/
│   ├── docx/docx.go                 # DOCX processing
│   ├── docx/pdf.go                  # PDF text extraction
│   └── humanize/
│       ├── humanize.go              # Main pipeline
│       ├── analyze.go               # Text analysis
│       ├── language.go              # Language detection
│       ├── replace.go               # Word replacement
│       ├── punctuation.go           # Punctuation rules
│       ├── dicts.go                 # Dictionary structures
│       ├── dict_en.go               # English (7900+)
│       ├── dict_ru.go               # Russian (5700+)
│       └── dict_mo.go               # Moldovan (5100+)
└── web/
    ├── index.html                   # Web UI
    ├── app.js                       # Frontend logic
    └── style.css                    # Styles
```

### 💖 Support the Project

If you find this tool useful, please consider supporting development:

- ☕ [Buy Me a Coffee](https://buymeacoffee.com/aristarh.ucolov)
- 💜 [DonationAlerts](https://www.donationalerts.com/r/aristarh_ucolov)

### 👤 Author

**Aristarh Ucolov** ([aristarh.ucolov@gmail.com](mailto:aristarh.ucolov@gmail.com))

### 📜 License

MIT License

---

## 🇷🇺 Русский

### О проекте <a name="русский"></a>

**Humanity Language Tool — Synonymizer** — мощный оффлайн-инструмент, который заменяет обычные слова на более изысканные синонимы. Поддерживает **английский**, **русский** и **молдавский/румынский** языки с автоматическим определением, обрабатывает простой текст и `.docx` документы с полным сохранением форматирования (шрифты, размеры, цвета, таблицы, списки и т.д.).

### ✨ Возможности

- 🌍 **Многоязычность**: Английский, русский, молдавский/румынский с автоопределением
- 📚 **Огромные словари синонимов** (~18 800 в сумме):
  - Английский: **7 900+** записей
  - Русский: **5 700+** записей
  - Молдавский: **5 100+** записей
- 📄 **Полная поддержка DOCX**: сохраняет всё форматирование (cross-run матчинг)
  - Все шрифты (Times New Roman, Arial, Calibri и др.), размеры (6pt–72pt+)
  - Жирный, курсив, подчёркивание, зачёркивание, цвета (RGB)
  - Таблицы, списки, колонтитулы, гиперссылки, комментарии, правки
- 📑 **Пакетная обработка DOCX**: много .docx за раз → один ZIP
- 📕 **Извлечение текста из PDF**: текст из PDF синонимизируется → .txt
- ✨ **Откат по клику**: заменённые слова подсвечены, клик отменяет замену
- 🌗 **Светлая / тёмная тема**
- 🔒 **100% оффлайн** — данные никуда не передаются
- 🌐 **Сетевой доступ** — работа с любых устройств в локальной сети
- 🎨 **Умная пунктуация** — корректные тире, кавычки, многоточие для всех языков
- ⚡ **Быстро и легко** — один файл .exe весом 7 МБ

### 🚀 Быстрый старт

#### Скачать и запустить
```bash
# Просто скачайте humanity.exe и запустите
humanity.exe
```

#### Параметры запуска
```bash
humanity.exe                        # Сетевой доступ (по умолчанию)
humanity.exe -addr 127.0.0.1:8080   # Только localhost
humanity.exe -no-open               # Не открывать браузер автоматически
```

При запуске сервер показывает:
```
━━━ Synonymaizer Server ━━━
🌐 Network Access (from other devices):
   http://192.168.0.57:8080
Local only (localhost):   http://127.0.0.1:8080
```

Откройте любой из этих адресов в браузере.

### 📖 Использование

1. **Выберите язык** (или используйте автоопределение)
2. **Вставьте текст** в левое окно
3. Нажмите **Synonymize**
4. Скопируйте улучшенный текст из правого окна

#### Обработка DOCX
1. Нажмите **Choose .docx…** сверху
2. Выберите ваш Word-документ
3. Обработанный файл автоматически скачается
4. Всё форматирование сохраняется в точности

### 🔧 API

#### POST /api/humanize
```json
{
  "text": "Хороший человек работает.",
  "language": "auto"
}
```

Ответ:
```json
{
  "output": "Превосходный индивид трудится.",
  "detectedLanguage": "ru",
  "synonymsApplied": 3,
  "notes": ["Language: Русский", "Synonyms applied: 3"]
}
```

#### POST /api/humanize-docx
Multipart форма с полем `file` (DOCX) и опциональным `language`.

Возвращает обработанный DOCX как бинарную загрузку.

#### GET /api/languages
Список всех поддерживаемых языков.

### 🛠 Сборка из исходного кода

Требования: Go 1.21+

```bash
git clone https://github.com/aristarh-ucolov/Humanity-Language-Tool-Synonymizer
cd Humanity-Language-Tool-Synonymizer
go build -ldflags="-s -w" -o humanity.exe
```

### 📁 Структура проекта

```
.
├── main.go                          # HTTP-сервер
├── internal/
│   ├── docx/docx.go                 # Обработка DOCX
│   ├── docx/pdf.go                  # Извлечение текста из PDF
│   └── humanize/
│       ├── humanize.go              # Главный pipeline
│       ├── analyze.go               # Анализ текста
│       ├── language.go              # Определение языка
│       ├── replace.go               # Замена слов
│       ├── punctuation.go           # Правила пунктуации
│       ├── dicts.go                 # Структуры словарей
│       ├── dict_en.go               # Английский (7900+)
│       ├── dict_ru.go               # Русский (5700+)
│       └── dict_mo.go               # Молдавский (5100+)
└── web/
    ├── index.html                   # Веб-интерфейс
    ├── app.js                       # Логика фронтенда
    └── style.css                    # Стили
```

### 💖 Поддержать проект

Если инструмент оказался вам полезен, пожалуйста, поддержите развитие:

- ☕ [Buy Me a Coffee](https://buymeacoffee.com/aristarh.ucolov)
- 💜 [DonationAlerts](https://www.donationalerts.com/r/aristarh_ucolov)

### 👤 Автор

**Аристарх Уколов** ([aristarh.ucolov@gmail.com](mailto:aristarh.ucolov@gmail.com))

### 📜 Лицензия

MIT License

---

<div align="center">

Made with ❤️ by **Aristarh Ucolov**

</div>
