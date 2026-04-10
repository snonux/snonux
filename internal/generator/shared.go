package generator

// navDefs is appended to every theme template when parsing.
// It defines named sub-templates shared across all themes:
//   - "splashGate" — synchronous script: first child of <body>; sets html.sno-splash-skip when
//     splash should not run (?splash=0, not index.html, or Referer from same-site index/pageN).
//   - "navhints"  — keyboard shortcut hint bar HTML
//   - "navmodal"  — full-screen expanded-post modal HTML + image-sizing CSS
//   - "navscript" — keyboard navigation JavaScript with distinct sounds per action
//
// Each theme calls {{template "splashGate"}}, {{template "navhints" .}}, {{template "navmodal" .}},
// and {{template "navscript" .}} at the appropriate points in its HTML.
// All theme-specific CSS lives in each theme file so themes stay self-contained.
const navDefs = `
{{define "splashGate"}}
<script>
(function(){
  try {
    var sp = new URLSearchParams(location.search);
    if (sp.get('splash') === '0') {
      document.documentElement.classList.add('sno-splash-skip');
      return;
    }
  } catch (_) {}
  var parts = location.pathname.split('/').filter(function(s) { return s.length; });
  var seg = (parts.length ? parts[parts.length - 1] : '').toLowerCase();
  var onIndex = (!seg || seg === 'index.html');
  var ref = document.referrer;
  function refIsSameSiteBlogPage(url) {
    if (!url) return false;
    try {
      var ru = new URL(url), cu = new URL(location.href);
      if (ru.origin !== cu.origin) return false;
      var rp = ru.pathname.split('/').filter(function(s) { return s.length; });
      var rs = (rp.length ? rp[rp.length - 1] : '').toLowerCase();
      if (rs === 'index.html' || rs === '') return true;
      if (/^page\d+\.html$/.test(rs)) return true;
      return false;
    } catch (_) { return false; }
  }
  if (!onIndex || refIsSameSiteBlogPage(ref)) document.documentElement.classList.add('sno-splash-skip');
})();
</script>
{{end}}

{{define "navhints"}}
<div class="nav-hints" aria-label="keyboard shortcuts">
    <span><kbd>j</kbd><kbd>k</kbd> or <kbd>↑</kbd><kbd>↓</kbd> select post</span>
    <span><kbd>Enter</kbd> expand</span>
    <span><kbd>Esc</kbd> close</span>
    <span><kbd>h</kbd><kbd>l</kbd> or <kbd>←</kbd><kbd>→</kbd> change page</span>
</div>
{{end}}

{{define "navmodal"}}
<style>
/* Thumbnail sizing in list view; modal overrides to full width. */
.post-image { max-height:220px; max-width:100%; object-fit:cover; cursor:pointer; }
#post-modal .post-image { max-height:none; width:100%; max-width:100%; object-fit:contain; cursor:default; }
/* Semi-transparent modal backdrop so the WebGL scene stays visible behind
   the expanded post. Theme-specific modal-inner keeps its own background. */
.post-modal { background:rgba(0,0,0,0.55) !important; backdrop-filter:blur(6px) !important; }
/* Content area max-width across all themes */
.overlay { max-width:1200px; margin-left:auto; margin-right:auto; }
/* Pagination: newer + older side by side at the bottom of the feed */
.page-nav-dual { display:flex; justify-content:center; align-items:center; flex-wrap:wrap;
                gap:clamp(16px,4vw,48px); }
/* Host note under the site subtitle (all themes) */
.logo-host { font-size:0.65rem; opacity:0.55; margin-top:4px; letter-spacing:0.3px; line-height:1.3; }
/* Atom feed link in header (paired with transmit in .nav) */
.nav { display:flex; align-items:center; gap:clamp(10px,2.2vw,20px); flex-wrap:wrap; justify-content:flex-end; }
a.header-feed-link { font-size:0.8rem; text-decoration:none; opacity:0.82; letter-spacing:0.04em; white-space:nowrap; }
a.header-feed-link:hover { opacity:1; text-decoration:underline; }
/* Full-viewport splash (theme-specific colours/animation on each .splash-THEMENAME) */
#splash-overlay { position:fixed; inset:0; z-index:2000; display:flex; flex-direction:column; align-items:center;
  justify-content:center; text-align:center; padding:max(16px,4vw); box-sizing:border-box; cursor:pointer;
  transition:opacity .55s ease, visibility .55s ease, transform .55s ease; }
#splash-overlay.splash--dismissed { opacity:0 !important; visibility:hidden !important;
  pointer-events:none !important; transform:scale(1.02); }
#splash-overlay:focus { outline:2px solid rgba(255,255,255,0.35); outline-offset:4px; }
/* Vignette over WebGL so 3D motion does not overpower the edges */
#splash-overlay::before { content:""; position:absolute; inset:0; z-index:1; pointer-events:none;
  background: radial-gradient(ellipse 92% 82% at 50% 42%, rgba(0,0,0,0) 32%, rgba(0,0,0,0.26) 68%, rgba(0,0,0,0.48) 100%); }
.splash-title { font-weight:700; letter-spacing:0.06em; line-height:1.15; }
.splash-tag { margin-top:0.35rem; font-size:0.76rem; letter-spacing:0.2em; text-transform:uppercase; }
.splash-hint { margin-top:1.25rem; font-size:0.72rem; letter-spacing:0.12em; }
#splash-overlay .splash-gl-canvas { position:absolute; inset:0; width:100%; height:100%; display:block; z-index:0; pointer-events:none; }
/* Frosted panel so title/tag/hint stay readable over busy shaders */
#splash-overlay .splash-inner { position:relative; z-index:2; max-width:min(520px,92vw);
  padding: clamp(1.15rem, 3.2vw, 1.75rem) clamp(1.3rem, 3.8vw, 1.95rem); border-radius:14px;
  background: rgba(0, 0, 0, 0.58); backdrop-filter: blur(14px); -webkit-backdrop-filter: blur(14px);
  box-shadow: 0 14px 44px rgba(0, 0, 0, 0.58), inset 0 1px 0 rgba(255, 255, 255, 0.07); }
#splash-overlay.splash-brutalist .splash-inner.splash-frame {
  padding: clamp(1.4rem, 4.5vw, 2.25rem) clamp(1.1rem, 3.5vw, 1.9rem); background: rgba(0, 0, 0, 0.78); }
html.sno-splash-skip #splash-overlay { display:none !important; visibility:hidden !important; pointer-events:none !important; }
</style>
<div class="post-modal" id="post-modal">
    <div class="modal-inner">
        <button class="modal-close" onclick="closeModal()">[ ESC ] CLOSE</button>
        <div id="modal-content"></div>
    </div>
</div>
{{end}}

{{define "navscript"}}
<script>
    (function splashSetup() {
        var el = document.getElementById('splash-overlay');
        if (!el) return;
        if (document.documentElement.classList.contains('sno-splash-skip')) {
            if (typeof window._snonuxSplashWebGLCleanup === 'function') {
                try { window._snonuxSplashWebGLCleanup(); } catch (_) {}
                window._snonuxSplashWebGLCleanup = null;
            }
            el.remove();
            return;
        }
        var splashAudioCtx = null;
        var splashChimePlayed = false;
        // Soft major arpeggio (G4 → C5 → E5 → G5); works once autopolicy allows audio.
        function playSplashChime() {
            if (splashChimePlayed) return;
            try {
                if (!splashAudioCtx) {
                    splashAudioCtx = new (window.AudioContext || window.webkitAudioContext)();
                }
                var ctx = splashAudioCtx;
                function ring() {
                    splashChimePlayed = true;
                    var now = ctx.currentTime;
                    var freqs = [392, 523.25, 659.25, 783.99];
                    var i, osc, g, t0;
                    for (i = 0; i < freqs.length; i++) {
                        osc = ctx.createOscillator();
                        g = ctx.createGain();
                        osc.connect(g);
                        g.connect(ctx.destination);
                        osc.type = 'sine';
                        osc.frequency.value = freqs[i];
                        t0 = now + i * 0.075;
                        g.gain.setValueAtTime(0, t0);
                        g.gain.linearRampToValueAtTime(0.1, t0 + 0.028);
                        g.gain.exponentialRampToValueAtTime(0.001, t0 + 0.52);
                        osc.start(t0);
                        osc.stop(t0 + 0.55);
                    }
                }
                ctx.resume().then(ring).catch(function() {});
            } catch (_) {}
        }
        playSplashChime();
        el.addEventListener('pointerdown', function() { playSplashChime(); }, { passive: true });
        function dismiss() {
            if (el.classList.contains('splash--dismissed')) return;
            playSplashChime();
            if (typeof window._snonuxSplashWebGLCleanup === 'function') {
                try { window._snonuxSplashWebGLCleanup(); } catch (_) {}
                window._snonuxSplashWebGLCleanup = null;
            }
            el.classList.add('splash--dismissed');
            setTimeout(function() { if (el.parentNode) el.parentNode.removeChild(el); }, 600);
        }
        el.addEventListener('click', function(e) { e.preventDefault(); dismiss(); });
        window._snonuxDismissSplash = dismiss;
        el.focus({ preventScroll: true });
    })();

    // === KEYBOARD NAVIGATION ===
    // j / ArrowDown  → next post       k / ArrowUp    → previous post
    // h / ArrowLeft  → previous page   l / ArrowRight → next page
    // Enter          → expand modal    Esc            → close modal
    const posts = document.querySelectorAll('.post');
    let currentIndex = posts.length > 0 ? 0 : -1;
    const prevPageURL = {{.PrevPageJSON}};
    const nextPageURL = {{.NextPageJSON}};

    if (currentIndex >= 0) selectPost(0);

    function selectPost(index) {
        if (posts.length === 0) return;
        if (currentIndex >= 0) posts[currentIndex].classList.remove('post-active');
        currentIndex = Math.max(0, Math.min(index, posts.length - 1));
        posts[currentIndex].classList.add('post-active');
        posts[currentIndex].scrollIntoView({ behavior: 'smooth', block: 'nearest' });
        playNavSound();
    }

    // playNavSound: short low beep for post selection (j/k navigation).
    function playNavSound() {
        try {
            const ctx = new (window.AudioContext || window.webkitAudioContext)();
            const osc = ctx.createOscillator();
            const gain = ctx.createGain();
            osc.connect(gain); gain.connect(ctx.destination);
            osc.frequency.value = 220; osc.type = 'sine';
            gain.gain.setValueAtTime(0.12, ctx.currentTime);
            gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.08);
            osc.start(ctx.currentTime); osc.stop(ctx.currentTime + 0.08);
        } catch (_) {}
    }

    // playOpenSound: bright ascending chime when modal opens (Enter key).
    function playOpenSound() {
        try {
            const ctx = new (window.AudioContext || window.webkitAudioContext)();
            const osc = ctx.createOscillator();
            const gain = ctx.createGain();
            osc.connect(gain); gain.connect(ctx.destination);
            osc.type = 'triangle';
            osc.frequency.setValueAtTime(440, ctx.currentTime);
            osc.frequency.exponentialRampToValueAtTime(880, ctx.currentTime + 0.14);
            gain.gain.setValueAtTime(0.10, ctx.currentTime);
            gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.20);
            osc.start(ctx.currentTime); osc.stop(ctx.currentTime + 0.20);
        } catch (_) {}
    }

    // playCloseSound: descending sweep when modal closes (Esc key).
    function playCloseSound() {
        try {
            const ctx = new (window.AudioContext || window.webkitAudioContext)();
            const osc = ctx.createOscillator();
            const gain = ctx.createGain();
            osc.connect(gain); gain.connect(ctx.destination);
            osc.type = 'sine';
            osc.frequency.setValueAtTime(440, ctx.currentTime);
            osc.frequency.exponentialRampToValueAtTime(110, ctx.currentTime + 0.15);
            gain.gain.setValueAtTime(0.10, ctx.currentTime);
            gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.18);
            osc.start(ctx.currentTime); osc.stop(ctx.currentTime + 0.18);
        } catch (_) {}
    }

    function openModal() {
        if (currentIndex < 0) return;
        document.getElementById('modal-content').innerHTML =
            posts[currentIndex].querySelector('.post-text').innerHTML;
        document.getElementById('post-modal').classList.add('active');
        playOpenSound();
    }

    function closeModal() {
        document.getElementById('post-modal').classList.remove('active');
        playCloseSound();
    }

    document.addEventListener('keydown', function(e) {
        if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') return;
        var splash = document.getElementById('splash-overlay');
        if (splash && !splash.classList.contains('splash--dismissed')) {
            if (e.key === 'Enter' || e.key === ' ' || e.key === 'Escape') {
                e.preventDefault();
                if (window._snonuxDismissSplash) window._snonuxDismissSplash();
            }
            return;
        }
        if (document.getElementById('post-modal').classList.contains('active')) {
            if (e.key === 'Escape') { closeModal(); e.preventDefault(); }
            return;
        }
        switch (e.key) {
            case 'j': case 'ArrowDown':  selectPost(currentIndex + 1); e.preventDefault(); break;
            case 'k': case 'ArrowUp':    selectPost(currentIndex - 1); e.preventDefault(); break;
            case 'h': case 'ArrowLeft':
                if (prevPageURL) { playNavSound(); window.location.href = prevPageURL; }
                e.preventDefault(); break;
            case 'l': case 'ArrowRight':
                if (nextPageURL) { playNavSound(); window.location.href = nextPageURL; }
                e.preventDefault(); break;
            case 'Enter': openModal(); e.preventDefault(); break;
        }
    });
</script>
{{end}}
`
