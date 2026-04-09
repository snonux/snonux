package generator

// brutalistTemplate is a raw brutalist theme — pure black, thick white borders,
// Impact font, red as the only accent. No rounded corners anywhere.
const brutalistTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SNONUX.FOO</title>
    <style>
        :root { --red:#ff2200; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:Impact,'Arial Narrow',Arial,sans-serif;
               background:#000; color:#fff; overflow:hidden; height:100vh; }
        .overlay { height:100vh; display:flex; flex-direction:column; }
        header { padding:14px 24px; background:#000; border-bottom:4px solid #fff;
                 display:flex; align-items:center; justify-content:space-between; }
        .logo { display:flex; align-items:center; gap:16px; }
        .logo-mark { font-size:2.8rem; color:var(--red); line-height:1; }
        .logo-title h1 { font-size:2rem; color:#fff; letter-spacing:0; line-height:1; }
        .logo-title .subtitle { font-size:0.78rem; color:#888; margin-top:3px;
                                 font-family:'Courier New',monospace; }
        .logo-title .subtitle a { color:var(--red); text-decoration:none; }
        .logo-title .subtitle a:hover { text-decoration:underline; }
        .transmit-btn { border:3px solid var(--red); color:var(--red); padding:10px 20px;
                        border-radius:0; text-decoration:none; font-family:Impact; font-size:1.05rem;
                        letter-spacing:2px; transition:all 0.1s; }
        .transmit-btn:hover { background:var(--red); color:#000; }
        .nav-hints { background:#111; border-bottom:2px solid #333; color:#888;
                     padding:5px 24px; display:flex; gap:18px; font-family:'Courier New',monospace;
                     font-size:0.7rem; flex-wrap:wrap; }
        .nav-hints kbd { background:#000; border:1px solid #555; color:#fff;
                         border-radius:0; padding:0 5px; margin:0 2px; font-size:0.7rem; }
        .content { flex:1; overflow-y:auto; padding:20px 24px;
                   scrollbar-width:thin; scrollbar-color:#fff #000; }
        .page-nav { display:flex; justify-content:center; margin:14px 0; }
        .page-nav a { border:3px solid #fff; color:#fff; padding:9px 22px;
                      border-radius:0; text-decoration:none; font-family:Impact;
                      font-size:1rem; letter-spacing:2px; }
        .page-nav a:hover { background:#fff; color:#000; }
        .post { background:#000; border:3px solid #fff; border-radius:0;
                padding:20px 22px; margin-bottom:14px; cursor:pointer;
                transition:border-color 0.1s,background 0.1s; }
        .post:hover { border-color:var(--red); }
        .post-active { border-color:var(--red) !important; background:#0d0000 !important;
                       border-left-width:8px !important; box-shadow:none !important; }
        .post-header { display:flex; justify-content:space-between; margin-bottom:12px; }
        .post-time { color:#aaa; font-family:'Courier New',monospace; font-size:0.82rem; }
        .post-text { font-family:'Arial',sans-serif; font-size:1rem; line-height:1.5; }
        .post-text a { color:var(--red); text-decoration:underline; }
        .post-image { max-width:100%; margin-top:10px; border:3px solid #fff; }
        .post-audio { width:100%; margin-top:10px; }
        .post-modal { display:none; position:fixed; inset:0; z-index:100;
                      background:rgba(0,0,0,0.98); overflow-y:auto; padding:40px 20px; }
        .post-modal.active { display:block; }
        .modal-inner { max-width:780px; margin:0 auto; background:#000;
                       border:4px solid #fff; border-radius:0; padding:38px;
                       box-shadow:8px 8px 0 var(--red); }
        .modal-close { float:right; background:none; border:none; color:var(--red);
                       font-family:Impact; font-size:1.3rem; cursor:pointer; letter-spacing:2px; }
        @media(max-width:640px) { .nav-hints{display:none;} header{padding:10px 16px;} .logo-mark{font-size:2rem;} }
    </style>
</head>
<body>
    <div class="overlay">
        <header>
            <div class="logo">
                <span class="logo-mark">SN</span>
                <div class="logo-title">
                    <h1>SNONUX.FOO</h1>
                    <p class="subtitle">MICROBLOG &mdash; <a href="https://foo.zone">FOO.ZONE</a> IS THE REAL BLOG</p>
                </div>
            </div>
            <div class="nav">
                <a href="https://foo.zone/about" class="transmit-btn">TRANSMIT</a>
            </div>
        </header>
        {{template "navhints" .}}
        <div class="content" id="post-content">
            {{if .PrevPage}}<div class="page-nav"><a href="{{.PrevPage}}">&larr; NEWER</a></div>{{end}}
            {{range $i, $post := .Posts}}
            <div class="post" data-index="{{$i}}" onclick="selectPost({{$i}})">
                <div class="post-header">
                    <div><strong>@SNONUX</strong></div>
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
