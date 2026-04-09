package generator

// terminalTemplate is the green phosphor CRT terminal theme.
// Monospace throughout, scanline overlay via CSS, no external dependencies.
const terminalTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>snonux.foo // TERMINAL</title>
    <style>
        :root { --p:#33ff33; --dim:#1a7a1a; --bg:#0a0a0a; --bg2:#050505; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:'Courier New',Courier,monospace; background:var(--bg); color:var(--p);
               overflow:hidden; height:100vh; position:relative; }
        /* CRT scanlines */
        body::before { content:''; position:fixed; inset:0; z-index:999; pointer-events:none;
            background:repeating-linear-gradient(0deg,transparent,transparent 2px,
                rgba(0,0,0,0.12) 2px,rgba(0,0,0,0.12) 4px); }
        /* Subtle screen flicker */
        @keyframes flicker { 0%,100%{opacity:1} 93%{opacity:0.97} 95%{opacity:0.91} 97%{opacity:0.98} }
        body { animation:flicker 9s infinite; }
        .overlay { position:relative; z-index:10; height:100vh; display:flex; flex-direction:column; }
        header { padding:12px 24px; background:var(--bg2); border-bottom:2px solid var(--p);
                 display:flex; align-items:center; justify-content:space-between; }
        .logo { display:flex; align-items:center; gap:14px; }
        .logo-mark { font-size:1.6rem; color:var(--p); text-shadow:0 0 14px var(--p); letter-spacing:2px; }
        .logo-title h1 { font-size:1.3rem; color:var(--p); text-shadow:0 0 10px var(--p);
                         letter-spacing:3px; font-weight:normal; }
        .logo-title .subtitle { font-size:0.72rem; color:var(--dim); margin-top:2px; }
        .logo-title .subtitle a { color:var(--p); text-decoration:none; }
        .logo-title .subtitle a:hover { text-shadow:0 0 6px var(--p); }
        .nav a.transmit-btn { border:1px solid var(--p); color:var(--p); padding:8px 18px;
                              border-radius:0; text-decoration:none; letter-spacing:2px; font-size:0.85rem;
                              transition:all 0.2s; }
        .nav a.transmit-btn:hover { background:var(--p); color:var(--bg); }
        .nav-hints { background:var(--bg2); border-bottom:1px solid var(--dim); color:var(--dim);
                     padding:5px 24px; display:flex; gap:18px; font-size:0.68rem; flex-wrap:wrap; }
        .nav-hints kbd { background:transparent; border:1px solid var(--dim); color:var(--p);
                         border-radius:0; padding:0 5px; font-size:0.7rem; margin:0 2px; }
        .content { flex:1; overflow-y:auto; padding:16px 24px;
                   scrollbar-width:thin; scrollbar-color:var(--dim) var(--bg); }
        .page-nav { display:flex; justify-content:center; margin:14px 0; }
        .page-nav a { border:1px solid var(--dim); color:var(--p); padding:7px 20px;
                      border-radius:0; text-decoration:none; letter-spacing:2px; font-size:0.82rem; }
        .page-nav a:hover { background:var(--p); color:var(--bg); border-color:var(--p); }
        .post { background:var(--bg); border:1px solid var(--dim); border-radius:0;
                padding:18px 20px; margin-bottom:12px; cursor:pointer; transition:border-color 0.15s; }
        .post:hover { border-color:var(--p); box-shadow:0 0 8px rgba(51,255,51,0.3); }
        .post-active { border-color:var(--p) !important; background:rgba(51,255,51,0.04) !important;
                       box-shadow:0 0 14px rgba(51,255,51,0.3),inset 3px 0 0 var(--p) !important; }
        .post-header { display:flex; justify-content:space-between; margin-bottom:12px; font-size:0.88rem; }
        .post-time { color:var(--dim); font-size:0.82rem; }
        .post-text { line-height:1.6; font-size:0.92rem; }
        .post-text a { color:var(--p); text-decoration:underline; }
        .post-image { max-width:100%; margin-top:10px; border:1px solid var(--dim); }
        .post-audio { width:100%; margin-top:10px; }
        .post-modal { display:none; position:fixed; inset:0; z-index:100;
                      background:rgba(0,0,0,0.97); overflow-y:auto; padding:40px 20px; }
        .post-modal.active { display:block; }
        .modal-inner { max-width:760px; margin:0 auto; background:var(--bg);
                       border:1px solid var(--p); border-radius:0;
                       box-shadow:0 0 40px rgba(51,255,51,0.25); padding:36px; }
        .modal-close { float:right; background:none; border:none; color:var(--p);
                       font-family:monospace; font-size:0.9rem; cursor:pointer; letter-spacing:2px; }
        @media(max-width:640px) { .nav-hints{display:none;} header{padding:10px 16px;} .content{padding:12px 16px;} }
    </style>
</head>
<body>
    <div class="overlay">
        <header>
            <div class="logo">
                <span class="logo-mark">[SN]</span>
                <div class="logo-title">
                    <h1>snonux.foo</h1>
                    <p class="subtitle">microblog / <a href="https://foo.zone">foo.zone</a> is the real blog</p>
                </div>
            </div>
            <div class="nav">
                <a href="https://foo.zone/about" class="transmit-btn">&gt; TRANSMIT</a>
            </div>
        </header>
        {{template "navhints" .}}
        <div class="content" id="post-content">
            {{if .PrevPage}}<div class="page-nav"><a href="{{.PrevPage}}">&lt;-- NEWER</a></div>{{end}}
            {{range $i, $post := .Posts}}
            <div class="post" data-index="{{$i}}" onclick="selectPost({{$i}})">
                <div class="post-header">
                    <div><strong>@snonux</strong></div>
                    <div class="post-time">{{$post.FormattedTime}}</div>
                </div>
                <div class="post-text">{{$post.ContentHTML}}</div>
            </div>
            {{end}}
            {{if .NextPage}}<div class="page-nav"><a href="{{.NextPage}}">OLDER --&gt;</a></div>{{end}}
        </div>
    </div>
    {{template "navmodal" .}}
    {{template "navscript" .}}
</body>
</html>`
