package generator

// navDefs is appended to every theme template when parsing.
// It defines three named sub-templates shared across all themes:
//   - "navhints"  — keyboard shortcut hint bar HTML
//   - "navmodal"  — full-screen expanded-post modal HTML + image-sizing CSS
//   - "navscript" — keyboard navigation JavaScript with distinct sounds per action
//
// Each theme calls {{template "navhints" .}}, {{template "navmodal" .}}, and
// {{template "navscript" .}} at the appropriate points in its HTML.
// All theme-specific CSS lives in each theme file so themes stay self-contained.
const navDefs = `
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
/* Thumbnail sizing in list view; modal overrides to full width so images
   appear larger when a post is expanded with Enter. */
.post-image { max-height:220px; max-width:100%; object-fit:cover; cursor:pointer; }
#post-modal .post-image { max-height:none; width:100%; max-width:100%; object-fit:contain; cursor:default; }
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
