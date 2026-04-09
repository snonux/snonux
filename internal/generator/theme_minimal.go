package generator

// minimalTemplate is a clean white theme — system font, subtle borders,
// no animations or decorations. Maximum readability.
const minimalTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>snonux.foo</title>
    <style>
        :root { --accent:#0066cc; --border:#e2e2e2; --muted:#666; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',system-ui,sans-serif;
               background:#fff; color:#111; overflow:hidden; height:100vh; }
        .overlay { height:100vh; display:flex; flex-direction:column; max-width:860px;
                   margin:0 auto; }
        header { padding:20px 32px; border-bottom:1px solid var(--border);
                 display:flex; align-items:center; justify-content:space-between; }
        .logo { display:flex; align-items:center; gap:14px; }
        .logo-mark { font-size:1.5rem; font-weight:800; color:#111; letter-spacing:-1px; }
        .logo-title h1 { font-size:1.35rem; font-weight:700; color:#111; letter-spacing:-0.5px; }
        .logo-title .subtitle { font-size:0.78rem; color:var(--muted); margin-top:1px; }
        .logo-title .subtitle a { color:var(--accent); text-decoration:none; }
        .logo-title .subtitle a:hover { text-decoration:underline; }
        .transmit-btn { border:1px solid var(--accent); color:var(--accent); padding:8px 18px;
                        border-radius:5px; text-decoration:none; font-size:0.88rem;
                        font-weight:500; transition:all 0.15s; }
        .transmit-btn:hover { background:var(--accent); color:#fff; }
        .nav-hints { padding:6px 32px; border-bottom:1px solid var(--border);
                     display:flex; gap:18px; font-size:0.72rem; color:var(--muted); flex-wrap:wrap; }
        .nav-hints kbd { background:#f5f5f5; border:1px solid #ccc; border-radius:3px;
                         padding:1px 5px; font-size:0.72rem; color:#333; margin:0 2px; }
        .content { flex:1; overflow-y:auto; padding:0 32px;
                   scrollbar-width:thin; scrollbar-color:#ccc #fff; }
        .page-nav { display:flex; justify-content:center; margin:16px 0; }
        .page-nav a { border:1px solid var(--border); color:var(--accent); padding:8px 20px;
                      border-radius:5px; text-decoration:none; font-size:0.88rem; }
        .page-nav a:hover { background:var(--accent); color:#fff; border-color:var(--accent); }
        .post { border-bottom:1px solid var(--border); padding:22px 0; cursor:pointer;
                transition:background 0.12s; }
        .post:hover { background:#f8f8f8; padding-left:8px; }
        .post-active { background:#eef5ff !important; border-left:3px solid var(--accent) !important;
                       padding-left:16px !important; }
        .post-header { display:flex; justify-content:space-between; margin-bottom:10px;
                       font-size:0.88rem; }
        .post-time { color:var(--muted); font-size:0.82rem; }
        .post-text { line-height:1.65; font-size:1rem; }
        .post-text a { color:var(--accent); text-decoration:none; }
        .post-text a:hover { text-decoration:underline; }
        .post-image { max-width:100%; border-radius:6px; margin-top:10px; border:1px solid var(--border); }
        .post-audio { width:100%; margin-top:10px; }
        .post-modal { display:none; position:fixed; inset:0; z-index:100;
                      background:rgba(255,255,255,0.97); overflow-y:auto; padding:40px 20px; }
        .post-modal.active { display:block; }
        .modal-inner { max-width:760px; margin:0 auto; background:#fff;
                       border:1px solid var(--border); border-radius:6px;
                       box-shadow:0 4px 24px rgba(0,0,0,0.1); padding:40px; }
        .modal-close { float:right; background:none; border:none; color:var(--muted);
                       font-size:0.9rem; cursor:pointer; }
        @media(max-width:640px) { .nav-hints{display:none;} .overlay{max-width:100%;} header{padding:16px 20px;} .content{padding:0 20px;} }
    </style>
</head>
<body>
    <div class="overlay">
        <header>
            <div class="logo">
                <span class="logo-mark">SN</span>
                <div class="logo-title">
                    <h1>snonux.foo</h1>
                    <p class="subtitle">microblog &mdash; <a href="https://foo.zone">foo.zone</a> is the real blog</p>
                </div>
            </div>
            <div class="nav">
                <a href="https://foo.zone/about" class="transmit-btn">Transmit to Nexus</a>
            </div>
        </header>
        {{template "navhints" .}}
        <div class="content" id="post-content">
            {{if .PrevPage}}<div class="page-nav"><a href="{{.PrevPage}}">&larr; Newer</a></div>{{end}}
            {{range $i, $post := .Posts}}
            <div class="post" data-index="{{$i}}" onclick="selectPost({{$i}})">
                <div class="post-header">
                    <div><strong>@snonux</strong></div>
                    <div class="post-time">{{$post.FormattedTime}}</div>
                </div>
                <div class="post-text">{{$post.ContentHTML}}</div>
            </div>
            {{end}}
            {{if .NextPage}}<div class="page-nav"><a href="{{.NextPage}}">Older &rarr;</a></div>{{end}}
        </div>
    </div>
    {{template "navmodal" .}}
    {{template "navscript" .}}
</body>
</html>`
