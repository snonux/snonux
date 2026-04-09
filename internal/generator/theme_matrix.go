package generator

// matrixTemplate is a hacker-style theme inspired by The Matrix — black
// background, bright matrix-green (#00ff41) text, monospace throughout,
// no decorations beyond a faint scanline overlay.
const matrixTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>snonux.foo // MATRIX</title>
    <style>
        :root { --g:#00ff41; --g2:#008f11; --g3:#003b00; --bg:#000; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:'Courier New',Courier,monospace; background:var(--bg); color:var(--g);
               overflow:hidden; height:100vh; }
        /* scanline overlay */
        body::before { content:''; position:fixed; inset:0; z-index:999; pointer-events:none;
            background:repeating-linear-gradient(0deg,transparent,transparent 3px,
                rgba(0,0,0,0.08) 3px,rgba(0,0,0,0.08) 4px); }
        @keyframes blink { 0%,100%{opacity:1} 50%{opacity:0} }
        .overlay { position:relative; z-index:10; height:100vh; display:flex; flex-direction:column; }
        header { padding:12px 24px; background:#000; border-bottom:1px solid var(--g2);
                 display:flex; align-items:center; justify-content:space-between; }
        .logo { display:flex; align-items:center; gap:14px; }
        .logo-mark { font-size:1.8rem; color:var(--g); text-shadow:0 0 18px var(--g); letter-spacing:3px; }
        /* blinking cursor after logo mark */
        .logo-mark::after { content:'_'; animation:blink 1.2s step-start infinite; }
        .logo-title h1 { font-size:1.2rem; color:var(--g); text-shadow:0 0 10px var(--g);
                         letter-spacing:4px; font-weight:normal; }
        .logo-title .subtitle { font-size:0.72rem; color:var(--g2); margin-top:2px; letter-spacing:1px; }
        .logo-title .subtitle a { color:var(--g); text-decoration:none; }
        .logo-title .subtitle a:hover { text-shadow:0 0 6px var(--g); }
        .transmit-btn { border:1px solid var(--g2); color:var(--g); padding:8px 18px;
                        text-decoration:none; font-size:0.82rem; letter-spacing:2px;
                        transition:all 0.1s; }
        .transmit-btn:hover { background:var(--g); color:var(--bg); }
        .nav-hints { background:#000; border-bottom:1px solid var(--g3); color:var(--g2);
                     padding:4px 24px; display:flex; gap:18px; font-size:0.68rem; flex-wrap:wrap; }
        .nav-hints kbd { background:transparent; border:1px solid var(--g3); color:var(--g);
                         padding:0 5px; font-size:0.68rem; margin:0 2px; }
        .content { flex:1; overflow-y:auto; padding:14px 24px;
                   scrollbar-width:thin; scrollbar-color:var(--g2) var(--bg); }
        .page-nav { display:flex; justify-content:center; margin:12px 0; }
        .page-nav a { border:1px solid var(--g2); color:var(--g); padding:7px 20px;
                      text-decoration:none; font-size:0.82rem; letter-spacing:2px; }
        .page-nav a:hover { background:var(--g); color:var(--bg); }
        .post { background:#000; border:1px solid var(--g3); padding:16px 18px;
                margin-bottom:10px; cursor:pointer; transition:border-color 0.15s; }
        .post:hover { border-color:var(--g2); box-shadow:0 0 8px rgba(0,255,65,0.2); }
        .post-active { border-color:var(--g) !important; background:rgba(0,255,65,0.03) !important;
                       box-shadow:0 0 14px rgba(0,255,65,0.3),inset 3px 0 0 var(--g) !important; }
        .post-header { display:flex; justify-content:space-between; margin-bottom:10px; font-size:0.85rem; }
        .post-time { color:var(--g2); font-size:0.78rem; }
        .post-text { line-height:1.6; font-size:0.88rem; }
        .post-text a { color:var(--g); text-decoration:underline; }
        .post-image { max-width:100%; margin-top:10px; border:1px solid var(--g3); }
        .post-audio { width:100%; margin-top:10px; }
        .post-modal { display:none; position:fixed; inset:0; z-index:100;
                      background:rgba(0,0,0,0.98); overflow-y:auto; padding:40px 20px; }
        .post-modal.active { display:block; }
        .modal-inner { max-width:740px; margin:0 auto; background:#000;
                       border:1px solid var(--g); padding:36px;
                       box-shadow:0 0 40px rgba(0,255,65,0.25); }
        .modal-close { float:right; background:none; border:none; color:var(--g2);
                       font-family:monospace; font-size:0.9rem; cursor:pointer; letter-spacing:2px; }
        @media(max-width:640px) { .nav-hints{display:none;} header{padding:10px 16px;} .content{padding:10px 16px;} }
    </style>
</head>
<body>
    <div class="overlay">
        <header>
            <div class="logo">
                <span class="logo-mark">SN</span>
                <div class="logo-title">
                    <h1>SNONUX.FOO</h1>
                    <p class="subtitle">MICROBLOG / <a href="https://foo.zone">FOO.ZONE</a> IS THE REAL BLOG</p>
                </div>
            </div>
            <div class="nav">
                <a href="https://foo.zone/about" class="transmit-btn">TRANSMIT</a>
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
