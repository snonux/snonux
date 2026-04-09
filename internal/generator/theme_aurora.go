package generator

// auroraTemplate is a dark navy theme with a CSS-animated aurora borealis
// effect — shifting green/purple/teal gradients across the background sky.
const auroraTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>snonux.foo ✦ AURORA</title>
    <style>
        :root { --green:#00ffb3; --teal:#00cfe8; --purple:#c084fc; --navy:#050d1a; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:'Segoe UI',system-ui,sans-serif; background:var(--navy);
               color:#e0f8f0; overflow:hidden; height:100vh; }
        /* Animated aurora bands */
        @keyframes aurora1 { 0%,100%{opacity:0.18;transform:scaleX(1) translateY(0)} 50%{opacity:0.28;transform:scaleX(1.15) translateY(-14px)} }
        @keyframes aurora2 { 0%,100%{opacity:0.12;transform:scaleX(1) translateY(0)} 50%{opacity:0.22;transform:scaleX(0.88) translateY(10px)} }
        @keyframes aurora3 { 0%,100%{opacity:0.10;transform:scaleX(1) skewY(0deg)} 50%{opacity:0.18;transform:scaleX(1.08) skewY(2deg)} }
        .aurora-bg { position:fixed; inset:0; z-index:0; overflow:hidden; }
        .aurora-bg::before { content:''; position:absolute; left:-20%; top:5%; width:140%; height:45%;
            background:radial-gradient(ellipse,rgba(0,255,179,0.38) 0%,rgba(0,207,232,0.22) 40%,transparent 70%);
            filter:blur(40px); animation:aurora1 12s ease-in-out infinite; }
        .aurora-bg::after { content:''; position:absolute; left:10%; top:20%; width:120%; height:55%;
            background:radial-gradient(ellipse,rgba(192,132,252,0.28) 0%,rgba(0,255,179,0.18) 45%,transparent 70%);
            filter:blur(50px); animation:aurora2 16s ease-in-out infinite; }
        .aurora-band3 { position:fixed; left:-10%; top:35%; width:130%; height:40%; z-index:0;
            background:radial-gradient(ellipse,rgba(0,207,232,0.22) 0%,rgba(192,132,252,0.14) 50%,transparent 75%);
            filter:blur(45px); animation:aurora3 20s ease-in-out infinite; }
        .overlay { position:relative; z-index:10; height:100vh; display:flex; flex-direction:column; }
        header { padding:16px 28px; background:rgba(5,13,26,0.78); backdrop-filter:blur(14px);
                 border-bottom:1px solid rgba(0,255,179,0.25); display:flex; align-items:center; justify-content:space-between; }
        .logo { display:flex; align-items:center; gap:14px; }
        .logo-mark { font-size:2rem; font-weight:800; background:linear-gradient(90deg,var(--green),var(--teal));
                     -webkit-background-clip:text; -webkit-text-fill-color:transparent; }
        .logo-title h1 { font-size:1.5rem; font-weight:700; color:#e0f8f0; letter-spacing:1px; }
        .logo-title .subtitle { font-size:0.75rem; color:rgba(224,248,240,0.55); margin-top:2px; }
        .logo-title .subtitle a { color:var(--green); text-decoration:none; }
        .logo-title .subtitle a:hover { text-shadow:0 0 8px var(--green); }
        .transmit-btn { border:1px solid var(--teal); color:var(--teal); padding:9px 20px;
                        border-radius:20px; text-decoration:none; font-size:0.85rem;
                        transition:all 0.2s; }
        .transmit-btn:hover { background:var(--teal); color:var(--navy); }
        .nav-hints { background:rgba(5,13,26,0.6); border-bottom:1px solid rgba(0,255,179,0.15);
                     color:rgba(224,248,240,0.45); padding:5px 28px; display:flex; gap:18px;
                     font-size:0.68rem; flex-wrap:wrap; }
        .nav-hints kbd { background:rgba(0,255,179,0.1); border:1px solid rgba(0,255,179,0.35);
                         color:var(--green); border-radius:3px; padding:0 5px; margin:0 2px; }
        .content { flex:1; overflow-y:auto; padding:20px 28px;
                   scrollbar-width:thin; scrollbar-color:var(--green) var(--navy); }
        .page-nav { display:flex; justify-content:center; margin:14px 0; }
        .page-nav a { border:1px solid var(--teal); color:var(--teal); padding:8px 20px;
                      border-radius:20px; text-decoration:none; font-size:0.82rem; letter-spacing:1px; }
        .page-nav a:hover { background:var(--teal); color:var(--navy); }
        .post { background:rgba(5,20,35,0.72); border:1px solid rgba(0,255,179,0.2); border-radius:10px;
                padding:20px; margin-bottom:14px; cursor:pointer;
                transition:all 0.25s; backdrop-filter:blur(6px); }
        .post:hover { border-color:var(--green); box-shadow:0 0 20px rgba(0,255,179,0.2); transform:translateY(-2px); }
        .post-active { border-color:var(--purple) !important; background:rgba(15,5,40,0.9) !important;
                       box-shadow:0 0 24px rgba(192,132,252,0.35),inset 3px 0 0 var(--purple) !important; }
        .post-header { display:flex; justify-content:space-between; margin-bottom:12px; font-size:0.88rem; }
        .post-time { color:var(--teal); font-family:monospace; font-size:0.8rem; }
        .post-text { line-height:1.65; font-size:0.95rem; }
        .post-text a { color:var(--green); text-decoration:none; }
        .post-text a:hover { text-shadow:0 0 8px var(--green); }
        .post-image { max-width:100%; border-radius:8px; margin-top:10px; }
        .post-audio { width:100%; margin-top:10px; }
        .post-modal { display:none; position:fixed; inset:0; z-index:100;
                      background:rgba(5,13,26,0.95); backdrop-filter:blur(20px);
                      overflow-y:auto; padding:40px 20px; }
        .post-modal.active { display:block; }
        .modal-inner { max-width:760px; margin:0 auto; background:rgba(5,20,40,0.97);
                       border:1px solid var(--green); border-radius:12px;
                       box-shadow:0 0 60px rgba(0,255,179,0.25); padding:40px; }
        .modal-close { float:right; background:none; border:none; color:var(--teal);
                       font-size:0.9rem; cursor:pointer; letter-spacing:1px; }
        @media(max-width:640px) { .nav-hints{display:none;} header{padding:12px 18px;} .content{padding:14px 18px;} }
    </style>
</head>
<body>
    <div class="aurora-bg"></div>
    <div class="aurora-band3"></div>
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
