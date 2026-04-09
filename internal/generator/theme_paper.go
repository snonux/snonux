package generator

// paperTemplate is a warm vintage newspaper theme — Georgia serif, sepia tones,
// subtle texture simulation via CSS, no animations.
const paperTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>snonux.foo — the microblog</title>
    <style>
        :root { --ink:#2c1810; --sepia:#c8a96e; --paper:#f5f0e8; --muted:#8b6f47; --accent:#7b2d00; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:Georgia,'Times New Roman',serif; background:var(--paper); color:var(--ink);
               overflow:hidden; height:100vh; }
        /* subtle paper grain simulation */
        body::after { content:''; position:fixed; inset:0; pointer-events:none; z-index:999;
            background-image:url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='200' height='200'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.85' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='200' height='200' filter='url(%23n)' opacity='0.03'/%3E%3C/svg%3E");
            opacity:0.4; }
        .overlay { height:100vh; display:flex; flex-direction:column; max-width:800px; margin:0 auto; }
        header { padding:18px 28px 14px; border-bottom:3px double var(--ink);
                 display:flex; align-items:center; justify-content:space-between; }
        .logo { display:flex; align-items:center; gap:14px; }
        .logo-mark { font-size:2.6rem; font-weight:700; color:var(--accent); line-height:1; font-style:italic; }
        .logo-title h1 { font-size:1.6rem; font-weight:700; color:var(--ink); letter-spacing:1px; font-variant:small-caps; }
        .logo-title .subtitle { font-size:0.78rem; color:var(--muted); margin-top:2px; font-style:italic; }
        .logo-title .subtitle a { color:var(--accent); text-decoration:none; }
        .logo-title .subtitle a:hover { text-decoration:underline; }
        .transmit-btn { border:2px solid var(--ink); color:var(--ink); padding:8px 16px;
                        text-decoration:none; font-size:0.82rem; font-variant:small-caps; letter-spacing:1px;
                        transition:all 0.15s; }
        .transmit-btn:hover { background:var(--ink); color:var(--paper); }
        .nav-hints { background:transparent; border-bottom:1px solid var(--sepia); color:var(--muted);
                     padding:4px 28px; display:flex; gap:18px; font-size:0.68rem; font-family:monospace; flex-wrap:wrap; }
        .nav-hints kbd { background:transparent; border:1px solid var(--muted); color:var(--ink);
                         padding:0 4px; font-size:0.68rem; margin:0 1px; }
        .content { flex:1; overflow-y:auto; padding:16px 28px;
                   scrollbar-width:thin; scrollbar-color:var(--sepia) var(--paper); }
        .page-nav { display:flex; justify-content:center; margin:14px 0; }
        .page-nav a { border:1px solid var(--ink); color:var(--ink); padding:7px 20px;
                      text-decoration:none; font-size:0.82rem; font-variant:small-caps; letter-spacing:1px; }
        .page-nav a:hover { background:var(--ink); color:var(--paper); }
        .post { border-bottom:1px solid var(--sepia); padding:18px 0; cursor:pointer;
                transition:background 0.12s; }
        .post:hover { background:#ede8dc; padding-left:8px; }
        .post-active { background:#e8e0cc !important; border-left:4px solid var(--accent) !important;
                       padding-left:12px !important; }
        .post-header { display:flex; justify-content:space-between; margin-bottom:10px; font-size:0.85rem; }
        .post-time { color:var(--muted); font-family:monospace; font-size:0.78rem; }
        .post-text { line-height:1.7; font-size:1rem; }
        .post-text a { color:var(--accent); text-decoration:none; }
        .post-text a:hover { text-decoration:underline; }
        .post-image { max-width:100%; margin-top:10px; border:1px solid var(--sepia); }
        .post-audio { width:100%; margin-top:10px; }
        .post-modal { display:none; position:fixed; inset:0; z-index:100;
                      background:rgba(245,240,232,0.97); overflow-y:auto; padding:40px 20px; }
        .post-modal.active { display:block; }
        .modal-inner { max-width:720px; margin:0 auto; background:var(--paper);
                       border:2px solid var(--ink); padding:40px;
                       box-shadow:4px 4px 0 var(--sepia); }
        .modal-close { float:right; background:none; border:none; color:var(--muted);
                       font-size:0.9rem; cursor:pointer; font-variant:small-caps; letter-spacing:1px; }
        @media(max-width:640px) { .nav-hints{display:none;} .overlay{max-width:100%;} header{padding:14px 16px;} .content{padding:12px 16px;} }
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
                <a href="https://foo.zone/about" class="transmit-btn">Transmit</a>
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
