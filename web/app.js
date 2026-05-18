const LANG_NAMES = {
  "auto": "Auto-detect",
  "en": "English",
  "ru": "Русский",
  "mo": "Moldovenească",
};

function countWords(s) {
  const m = s.trim().match(/\S+/g);
  return m ? m.length : 0;
}

function updateStats() {
  const inp = document.getElementById("input").value;
  const out = document.getElementById("output").value;
  document.getElementById("input-stats").textContent = `${inp.length} characters · ${countWords(inp)} words`;
  document.getElementById("output-stats").textContent = `${out.length} characters · ${countWords(out)} words`;
}

async function run() {
  const text = document.getElementById("input").value;
  const language = document.getElementById("language").value;
  const status = document.getElementById("status");
  const detected = document.getElementById("detected-lang");
  const btn = document.getElementById("run-btn");
  const changes = document.getElementById("changes");

  if (!text.trim()) {
    status.textContent = "Paste some text first.";
    status.className = "status err";
    return;
  }

  btn.disabled = true;
  status.textContent = "Processing…";
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
    status.textContent = "Done";
    status.className = "status ok";

    if (data.detectedLanguage) {
      const langDisplay = LANG_NAMES[data.detectedLanguage] || data.detectedLanguage;
      detected.textContent = `Detected: ${langDisplay}`;
    }

    if (Array.isArray(data.notes) && data.notes.length) {
      changes.textContent = data.notes.join(" · ");
    }
  } catch (e) {
    status.textContent = "Failed: " + e.message;
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
      status.textContent = "Pasted";
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
    status.textContent = "Copied to clipboard";
    status.className = "status ok";
  } catch (_) {
    const ta = document.getElementById("output");
    ta.select();
    document.execCommand("copy");
  }
}

async function uploadDocx(file) {
  const status = document.getElementById("docx-status");
  const language = document.getElementById("language").value;
  if (!file) return;
  if (file.size === 0) {
    status.textContent = "File is empty.";
    status.className = "status err";
    return;
  }
  if (!file.name.toLowerCase().endsWith(".docx")) {
    status.textContent = "Not a .docx file.";
    status.className = "status err";
    return;
  }

  status.textContent = `Processing "${file.name}"…`;
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

    status.textContent = `Ready: "${outName}" [${langDisplay}] (${changed}/${total} runs changed, ${synonyms} synonyms)`;
    status.className = "status ok";
  } catch (e) {
    status.textContent = "Failed: " + e.message;
    status.className = "status err";
  }
}

window.addEventListener("DOMContentLoaded", () => {
  document.getElementById("run-btn").addEventListener("click", run);
  document.getElementById("paste-btn").addEventListener("click", paste);
  document.getElementById("copy-btn").addEventListener("click", copyOut);
  document.getElementById("clear-btn").addEventListener("click", () => {
    document.getElementById("input").value = "";
    document.getElementById("output").value = "";
    document.getElementById("changes").textContent = "";
    document.getElementById("detected-lang").textContent = "";
    updateStats();
  });
  document.getElementById("input").addEventListener("input", updateStats);
  document.getElementById("docx-file").addEventListener("change", (e) => {
    const f = e.target.files && e.target.files[0];
    if (f) uploadDocx(f);
    e.target.value = ""; // allow re-upload of same file
  });
  document.addEventListener("keydown", (e) => {
    if ((e.ctrlKey || e.metaKey) && e.key === "Enter") run();
  });
  updateStats();
});
