package generator

// synthwaveTemplate is the 80s retrowave theme — dark purple sky, CSS perspective
// grid floor, hot pink/orange accents, Russo One font.
const synthwaveTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>snonux.foo ⊕ SYNTHWAVE</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link href="https://fonts.googleapis.com/css2?family=Russo+One&family=Share+Tech+Mono&display=swap" rel="stylesheet">
    <style>
        :root { --pink:#ff2d78; --purple:#bf3fff; --orange:#ff6b2b; --bg:#0d0221; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:'Russo One','Arial Black',sans-serif; background:var(--bg);
               color:#fff; overflow:hidden; height:100vh; }
        /* Sunset sky gradient */
        .sky { position:fixed; inset:0; z-index:0;
               background:linear-gradient(180deg,#0d0221 0%,#1a0533 45%,#4a0080 72%,#8b1070 88%,#c8365a 100%); }
        /* Perspective grid floor */
        .grid-floor { position:fixed; bottom:0; left:0; width:100%; height:46vh; z-index:1;
            background-image:linear-gradient(rgba(255,45,120,0.35) 1px,transparent 1px),
                             linear-gradient(90deg,rgba(255,45,120,0.35) 1px,transparent 1px);
            background-size:44px 44px;
            transform:perspective(380px) rotateX(76deg); transform-origin:bottom;
            mask-image:linear-gradient(to top,rgba(0,0,0,0.85) 0%,transparent 100%); }
        .overlay { position:relative; z-index:10; height:100vh; display:flex; flex-direction:column; }
        header { padding:16px 28px; background:rgba(13,2,33,0.82); backdrop-filter:blur(10px);
                 border-bottom:2px solid var(--pink); display:flex; align-items:center; justify-content:space-between; }
        .logo { display:flex; align-items:center; gap:14px; }
        .logo-mark { font-size:1.8rem; background:linear-gradient(90deg,var(--pink),var(--purple));
                     -webkit-background-clip:text; -webkit-text-fill-color:transparent; }
        .logo-title h1 { font-size:1.7rem; background:linear-gradient(90deg,var(--pink),var(--orange));
                         -webkit-background-clip:text; -webkit-text-fill-color:transparent; letter-spacing:2px; }
        .logo-title .subtitle { font-size:0.7rem; color:rgba(255,255,255,0.55); margin-top:2px;
                                 font-family:'Share Tech Mono',monospace; }
        .logo-title .subtitle a { color:var(--pink); text-decoration:none; }
        .logo-title .subtitle a:hover { text-shadow:0 0 8px var(--pink); }
        .transmit-btn { border:2px solid var(--orange); color:var(--orange); padding:10px 22px;
                        border-radius:4px; text-decoration:none; letter-spacing:1px;
                        font-size:0.88rem; transition:all 0.2s; }
        .transmit-btn:hover { background:var(--orange); color:var(--bg); }
        .nav-hints { background:rgba(13,2,33,0.75); border-bottom:1px solid rgba(255,45,120,0.3);
                     color:rgba(255,255,255,0.45); padding:5px 20px; display:flex; gap:18px;
                     font-size:0.68rem; font-family:'Share Tech Mono',monospace; flex-wrap:wrap; }
        .nav-hints kbd { background:rgba(255,45,120,0.15); border:1px solid rgba(255,45,120,0.4);
                         color:var(--pink); border-radius:3px; padding:0 5px; margin:0 2px; font-size:0.7rem; }
        .content { flex:1; overflow-y:auto; padding:22px 28px;
                   scrollbar-width:thin; scrollbar-color:var(--pink) var(--bg); }
        .page-nav { display:flex; justify-content:center; margin:14px 0; }
        .page-nav a { border:2px solid var(--purple); color:var(--purple); padding:8px 22px;
                      border-radius:4px; text-decoration:none; letter-spacing:2px; font-size:0.82rem; }
        .page-nav a:hover { background:var(--purple); color:#fff; }
        .post { background:rgba(20,5,50,0.85); border:1px solid var(--purple); border-radius:6px;
                padding:22px; margin-bottom:18px; cursor:pointer; transition:all 0.25s; }
        .post:hover { border-color:var(--pink); box-shadow:0 0 22px rgba(255,45,120,0.35); transform:translateY(-3px); }
        .post-active { border-color:var(--orange) !important; background:rgba(30,8,60,0.96) !important;
                       box-shadow:0 0 22px rgba(255,107,43,0.45),inset 3px 0 0 var(--orange) !important; }
        .post-header { display:flex; justify-content:space-between; margin-bottom:14px; }
        .post-time { color:var(--orange); font-family:'Share Tech Mono',monospace; font-size:0.85rem; }
        .post-text { line-height:1.6; font-size:0.95rem; font-family:'Share Tech Mono',monospace; }
        .post-text a { color:var(--pink); text-decoration:none; }
        .post-image { max-width:100%; border-radius:6px; margin-top:10px; }
        .post-audio { width:100%; margin-top:10px; }
        .post-modal { display:none; position:fixed; inset:0; z-index:100;
                      background:rgba(13,2,33,0.96); overflow-y:auto; padding:40px 20px; }
        .post-modal.active { display:block; }
        .modal-inner { max-width:780px; margin:0 auto; background:rgba(20,5,50,0.98);
                       border:2px solid var(--pink); border-radius:6px;
                       box-shadow:0 0 60px rgba(255,45,120,0.35); padding:38px; }
        .modal-close { float:right; background:none; border:none; color:var(--orange);
                       font-family:'Russo One',sans-serif; font-size:0.9rem; cursor:pointer; letter-spacing:2px; }
        @media(max-width:640px) { .nav-hints{display:none;} .grid-floor{height:30vh;} header{padding:12px 18px;} }
    </style>
</head>
<body>
    <div class="sky"></div>
    <div class="grid-floor"></div>
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
                <a href="https://foo.zone/about" class="transmit-btn">TRANSMIT TO NEXUS</a>
            </div>
        </header>
        {{template "navhints" .}}
        <div class="content" id="post-content">
            {{if .PrevPage}}<div class="page-nav"><a href="{{.PrevPage}}">&larr; NEWER</a></div>{{end}}
            {{range $i, $post := .Posts}}
            <div class="post" data-index="{{$i}}" onclick="selectPost({{$i}})">
                <div class="post-header">
                    <div><strong>@snonux</strong></div>
                    <div class="post-time">{{$post.FormattedTime}}</div>
                </div>
                <div class="post-text">{{$post.ContentHTML}}</div>
            </div>
            {{end}}
            {{if .NextPage}}<div class="page-nav"><a href="{{.NextPage}}">OLDER &rarr;</a></div>{{end}}
        </div>
    </div>
    {{template "navmodal" .}}
    {{template "navscript" .}}
</body>
</html>`
