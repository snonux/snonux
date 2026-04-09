package generator

// navDefs is appended to every theme template when parsing.
// It defines three named sub-templates shared across all themes:
//   - "navhints"  — keyboard shortcut hint bar HTML
//   - "navmodal"  — full-screen expanded-post modal HTML
//   - "navscript" — keyboard navigation JavaScript
//
// Each theme calls {{template "navhints" .}}, {{template "navmodal" .}}, and
// {{template "navscript" .}} at the appropriate points in its HTML.
// All CSS for these elements (colours, borders, backdrop) lives in each theme
// so themes remain self-contained and independently styled.
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

    // playNavSound generates a short beep via the Web Audio API.
    // A fresh AudioContext per call avoids state issues across navigations.
    function playNavSound() {
        try {
            const ctx = new (window.AudioContext || window.webkitAudioContext)();
            const osc = ctx.createOscillator();
            const gain = ctx.createGain();
            osc.connect(gain); gain.connect(ctx.destination);
            osc.frequency.value = 220; osc.type = 'sine';
            gain.gain.setValueAtTime(0.15, ctx.currentTime);
            gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.08);
            osc.start(ctx.currentTime); osc.stop(ctx.currentTime + 0.08);
        } catch (_) {}
    }

    function openModal() {
        if (currentIndex < 0) return;
        document.getElementById('modal-content').innerHTML =
            posts[currentIndex].querySelector('.post-text').innerHTML;
        document.getElementById('post-modal').classList.add('active');
    }

    function closeModal() {
        document.getElementById('post-modal').classList.remove('active');
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
