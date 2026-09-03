/* TrieOH deck · navegação mínima (~70 linhas, sem dependências)
   Setas/space/PageUp/PageDown/Home/End · F = fullscreen · hash sync
   Qualquer erro inesperado aparece num toast no canto p/ depurar. */
(function () {
  var slides = Array.prototype.slice.call(document.querySelectorAll(".slide"));
  if (!slides.length) return;
  var cur = 0;
  var curEl = document.getElementById("cur");
  var totEl = document.getElementById("tot");
  var errBox = null;
  if (totEl) totEl.textContent = String(slides.length).padStart(2, "0");

  function toast(msg) {
    console.error("[deck]", msg);
    if (!errBox) {
      errBox = document.createElement("div");
      errBox.style.cssText =
        "position:fixed;left:16px;bottom:16px;z-index:99;max-width:70vw;" +
        "background:#3a1216;color:#ffb3b8;border:1px solid #ff5f57;border-radius:10px;" +
        "padding:10px 14px;font:12px/1.4 ui-monospace,monospace;white-space:pre-wrap";
      document.body.appendChild(errBox);
    }
    errBox.textContent = "[erro] " + msg;
  }
  window.addEventListener("error", function (e) { toast(e.message + " @ " + (e.filename || "").split("/").pop() + ":" + e.lineno); });
  window.addEventListener("unhandledrejection", function (e) { toast("promise: " + e.reason); });

  function pad(n) { return String(n + 1).padStart(2, "0"); }

  function titles() {
    var t = document.querySelector(".slide.active h1, .slide.active h2");
    return t ? t.textContent.replace(/\s+/g, " ").trim() : "";
  }

  function go(i, push) {
    try {
      cur = Math.max(0, Math.min(slides.length - 1, i));
      slides.forEach(function (s, k) { s.classList.toggle("active", k === cur); });
      if (curEl) curEl.textContent = pad(cur);
      var act = slides[cur];
      if (act) { act.scrollTop = 0; void act.getBoundingClientRect(); } /* força repaint agendado */
      document.title = (titles() ? titles() + " · " : "") + "O Poder da Ambição";
      if (push !== false) { try { history.replaceState(null, "", "#" + pad(cur)); } catch (e) {} }
    } catch (e) { toast(e.message); }
  }

  function next() { go(cur + 1); }
  function prev() { go(cur - 1); }

  document.addEventListener("keydown", function (e) {
    switch (e.key) {
      case "ArrowRight": case "ArrowDown": case "PageDown":
        e.preventDefault(); next(); break;
      case "ArrowLeft": case "ArrowUp": case "PageUp":
        e.preventDefault(); prev(); break;
      case "Home": e.preventDefault(); go(0); break;
      case "End": e.preventDefault(); go(slides.length - 1); break;
      case " ":
        /* se o foco está num botão, deixa o clique nativo agir (evita avanço duplo) */
        if (e.target && e.target.closest && e.target.closest("button")) return;
        e.preventDefault(); next(); break;
      case "Enter":
        if (e.target && e.target.closest && e.target.closest("button")) return;
        e.preventDefault(); next(); break;
      case "f": case "F":
        if (document.fullscreenElement) document.exitFullscreen();
        else if (document.documentElement.requestFullscreen) document.documentElement.requestFullscreen();
        break;
    }
  });

  var n = document.getElementById("next"), p = document.getElementById("prev");
  function bind(btn, fn) {
    if (!btn) return;
    btn.addEventListener("click", function () { fn(); btn.blur(); });
  }
  bind(n, next);
  bind(p, prev);

  /* toque simples (swipe) */
  var x0 = null;
  document.addEventListener("touchstart", function (e) { x0 = e.touches[0].clientX; }, { passive: true });
  document.addEventListener("touchend", function (e) {
    if (x0 === null) return;
    var dx = e.changedTouches[0].clientX - x0;
    if (Math.abs(dx) > 48) (dx < 0 ? next() : prev());
    x0 = null;
  }, { passive: true });

  /* back/forward no histórico */
  window.addEventListener("hashchange", function () {
    var m = location.hash.match(/^#(\d+)/);
    if (m) go(Math.min(parseInt(m[1], 10) - 1, slides.length - 1), false);
  });

  var m = location.hash.match(/^#(\d+)/);
  go(m ? Math.min(parseInt(m[1], 10) - 1, slides.length - 1) : 0, false);
})();
