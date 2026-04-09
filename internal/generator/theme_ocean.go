package generator

// oceanTemplate is a deep-ocean theme — dark navy/midnight blue background,
// teal/aqua/seafoam accents, subtle wave gradient at the bottom.
const oceanTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>snonux.foo ~ OCEAN</title>
    <style>
        :root { --teal:#00b4d8; --aqua:#48cae4; --deep:#023e8a; --navy:#03045e; --foam:#caf0f8; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:'Segoe UI',system-ui,sans-serif; background:var(--navy);
               color:var(--foam); overflow:hidden; height:100vh; }
        /* Deep ocean gradient base */
        body { background:linear-gradient(180deg,#03045e 0%,#023e8a 60%,#0077b6 100%); }
        /* Animated wave shimmer at bottom */
        @keyframes wave { 0%,100%{transform:translateX(0)} 50%{transform:translateX(-30px)} }
        .wave-bottom { position:fixed; bottom:0; left:-5%; width:110%; height:120px; z-index:1;
            background:radial-gradient(ellipse 80% 60% at 50% 100%,rgba(0,180,216,0.22),transparent);
            animation:wave 8s ease-in-out infinite; }
        .overlay { position:relative; z-index:10; height:100vh; display:flex; flex-direction:column; }
        header { padding:16px 28px; background:rgba(3,4,94,0.82); backdrop-filter:blur(12px);
                 border-bottom:1px solid rgba(0,180,216,0.3); display:flex; align-items:center; justify-content:space-between; }
        .logo { display:flex; align-items:center; gap:14px; }
        .logo-mark { font-size:2rem; font-weight:800; color:var(--aqua); text-shadow:0 0 16px var(--teal); }
        .logo-title h1 { font-size:1.5rem; font-weight:700; color:var(--foam); letter-spacing:1px; }
        .logo-title .subtitle { font-size:0.75rem; color:rgba(202,240,248,0.55); margin-top:2px; }
        .logo-title .subtitle a { color:var(--aqua); text-decoration:none; }
        .logo-title .subtitle a:hover { text-shadow:0 0 8px var(--teal); }
        .transmit-btn { border:1px solid var(--teal); color:var(--teal); padding:9px 20px;
                        border-radius:20px; text-decoration:none; font-size:0.85rem;
                        transition:all 0.2s; }
        .transmit-btn:hover { background:var(--teal); color:var(--navy); }
        .nav-hints { background:rgba(3,4,94,0.65); border-bottom:1px solid rgba(0,180,216,0.18);
                     color:rgba(202,240,248,0.45); padding:5px 28px; display:flex; gap:18px;
                     font-size:0.68rem; flex-wrap:wrap; }
        .nav-hints kbd { background:rgba(0,180,216,0.12); border:1px solid rgba(0,180,216,0.35);
                         color:var(--aqua); border-radius:3px; padding:0 5px; margin:0 2px; }
        .content { flex:1; overflow-y:auto; padding:20px 28px;
                   scrollbar-width:thin; scrollbar-color:var(--teal) var(--navy); }
        .page-nav { display:flex; justify-content:center; margin:14px 0; }
        .page-nav a { border:1px solid var(--deep); color:var(--aqua); padding:8px 20px;
                      border-radius:20px; text-decoration:none; font-size:0.82rem; }
        .page-nav a:hover { background:var(--teal); color:var(--navy); }
        .post { background:rgba(3,4,94,0.55); border:1px solid rgba(0,180,216,0.22); border-radius:10px;
                padding:20px; margin-bottom:14px; cursor:pointer;
                transition:all 0.25s; backdrop-filter:blur(6px); }
        .post:hover { border-color:var(--teal); box-shadow:0 4px 24px rgba(0,180,216,0.22); transform:translateY(-2px); }
        .post-active { border-color:var(--aqua) !important; background:rgba(0,100,150,0.55) !important;
                       box-shadow:0 0 22px rgba(72,202,228,0.35),inset 3px 0 0 var(--aqua) !important; }
        .post-header { display:flex; justify-content:space-between; margin-bottom:12px; font-size:0.88rem; }
        .post-time { color:var(--teal); font-family:monospace; font-size:0.8rem; }
        .post-text { line-height:1.65; font-size:0.95rem; }
        .post-text a { color:var(--aqua); text-decoration:none; }
        .post-text a:hover { text-shadow:0 0 8px var(--teal); }
        .post-image { max-width:100%; border-radius:8px; margin-top:10px; }
        .post-audio { width:100%; margin-top:10px; }
        .post-modal { display:none; position:fixed; inset:0; z-index:100;
                      background:rgba(3,4,94,0.96); backdrop-filter:blur(20px);
                      overflow-y:auto; padding:40px 20px; }
        .post-modal.active { display:block; }
        .modal-inner { max-width:760px; margin:0 auto; background:rgba(2,30,80,0.98);
                       border:1px solid var(--teal); border-radius:12px;
                       box-shadow:0 0 60px rgba(0,180,216,0.3); padding:40px; }
        .modal-close { float:right; background:none; border:none; color:var(--teal);
                       font-size:0.9rem; cursor:pointer; letter-spacing:1px; }
        @media(max-width:640px) { .nav-hints{display:none;} header{padding:12px 18px;} .content{padding:14px 18px;} }
    </style>
</head>
<body>
    <div class="wave-bottom"></div>
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
