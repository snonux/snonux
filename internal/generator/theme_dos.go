package generator

// dosTemplate is a classic DOS / IBM PC text-mode theme — blue background
// (#0000aa), white/yellow text, VT323 webfont for authentic VGA bitmap look,
// double-line box-drawing borders, and a BIOS-style layout.
// WebGL scene: falling green "rain" characters (BASIC-era) on the blue BG.
const dosTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SNONUX.FOO - DOS</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=VT323&display=swap" rel="stylesheet">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/three.js/r134/three.min.js"></script>
    <style>
        :root { --dos-blue:#0000aa; --dos-lblue:#5555ff; --dos-white:#aaaaaa;
                --dos-bwhite:#ffffff; --dos-yellow:#ffff55; --dos-cyan:#55ffff;
                --dos-red:#ff5555; --dos-bg:#000088; --dos-black:#000000; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:'VT323','Courier New',monospace; background:var(--dos-blue);
               color:var(--dos-bwhite); overflow:hidden; height:100vh; font-size:18px; }
        #three-canvas { position:fixed; top:0; left:0; width:100%; height:100%; z-index:1; }
        .overlay { position:relative; z-index:10; height:100vh; display:flex; flex-direction:column; }
        header { padding:8px 20px; background:var(--dos-white); color:var(--dos-blue);
                 display:flex; align-items:center; justify-content:space-between; }
        .logo { display:flex; align-items:center; gap:12px; }
        .logo-mark { font-size:1.8rem; color:var(--dos-blue); font-weight:bold; }
        .logo-title h1 { font-size:1.4rem; color:var(--dos-blue); font-weight:normal; letter-spacing:2px; }
        .logo-title .subtitle { font-size:0.85rem; color:var(--dos-blue); margin-top:1px; }
        .logo-title .subtitle a { color:var(--dos-blue); text-decoration:underline; }
        .logo-title .subtitle a:hover { color:var(--dos-black); }
        .transmit-btn { border:2px solid var(--dos-blue); color:var(--dos-blue); padding:4px 14px;
                        text-decoration:none; font-size:1rem; letter-spacing:1px;
                        transition:all 0.1s; }
        .transmit-btn:hover { background:var(--dos-blue); color:var(--dos-bwhite); }
        a.header-feed-link { color:var(--dos-blue); }
        a.header-feed-link:hover { color:var(--dos-black); }
        .nav-hints { background:var(--dos-blue); border-bottom:2px solid var(--dos-lblue);
                     color:var(--dos-cyan); padding:4px 20px; display:flex; gap:16px;
                     font-size:0.85rem; flex-wrap:wrap; }
        .nav-hints kbd { background:var(--dos-black); border:1px solid var(--dos-lblue);
                         color:var(--dos-yellow); padding:0 5px; margin:0 2px; }
        .content { flex:1; overflow-y:auto; padding:12px 20px;
                   scrollbar-width:thin; scrollbar-color:var(--dos-lblue) var(--dos-blue); }
        .page-nav { display:flex; justify-content:center; margin:10px 0; }
        .page-nav a { border:2px solid var(--dos-lblue); color:var(--dos-yellow); padding:6px 18px;
                      text-decoration:none; font-size:1rem; letter-spacing:1px; }
        .page-nav a:hover { background:var(--dos-lblue); color:var(--dos-bwhite); }
        .page-nav-footer { flex-shrink:0; padding:6px 20px; display:flex; justify-content:center;
            background:var(--dos-white); color:var(--dos-blue); }
        .page-nav-footer .page-nav a { border-color:var(--dos-blue); color:var(--dos-blue); }
        .page-nav-footer .page-nav a:hover { background:var(--dos-blue); color:var(--dos-bwhite); }
        .post { background:var(--dos-black); border:2px solid var(--dos-lblue);
                padding:12px 14px; margin-bottom:8px; cursor:pointer;
                transition:border-color 0.1s; }
        .post:hover { border-color:var(--dos-yellow);
                      box-shadow:0 0 0 1px var(--dos-yellow); }
        .post-active { border-color:var(--dos-yellow) !important;
                       background:rgba(0,0,170,0.3) !important;
                       box-shadow:0 0 0 2px var(--dos-yellow),inset 3px 0 0 var(--dos-yellow) !important; }
        .post-header { display:flex; justify-content:space-between; margin-bottom:8px; font-size:1rem; }
        .post-header strong { color:var(--dos-yellow); }
        .post-time { color:var(--dos-cyan); font-size:0.95rem; }
        .post-text { line-height:1.5; font-size:1.05rem; }
        .post-text a { color:var(--dos-cyan); text-decoration:underline; }
        .post-text a:hover { color:var(--dos-yellow); }
        .post-image { max-width:100%; margin-top:8px; border:2px solid var(--dos-lblue); }
        .post-audio { width:100%; margin-top:8px; }
        .post-modal { display:none; position:fixed; inset:0; z-index:100;
                      background:rgba(0,0,0,0.95); overflow-y:auto; padding:40px 20px; }
        .post-modal.active { display:block; }
        .modal-inner { max-width:740px; margin:0 auto; background:var(--dos-black);
                       border:2px solid var(--dos-yellow); padding:24px;
                       box-shadow:0 0 20px rgba(85,85,255,0.4); }
        .modal-close { float:right; background:var(--dos-white); border:2px outset var(--dos-bwhite);
                       color:var(--dos-blue); font-family:'VT323','Courier New',monospace;
                       font-size:1rem; cursor:pointer; padding:2px 8px; }
        .modal-close:hover { background:var(--dos-blue); color:var(--dos-bwhite);
                             border-style:inset; }
        @media(max-width:640px) { .nav-hints{display:none;} header{padding:6px 12px;} .content{padding:8px 12px;} }
        .splash-overlay.splash-dos { background:var(--dos-blue); font-family:'VT323','Courier New',monospace; }
        .splash-dos .splash-inner { position:relative; z-index:1; }
        .splash-dos .splash-title {
            font-size:clamp(1.4rem,4.5vw,2rem); color:var(--dos-bwhite);
            letter-spacing:0.15em;
            animation: splashDosBlink 1s step-end infinite;
        }
        @keyframes splashDosBlink { 0%,100%{border-right:0.6em solid var(--dos-bwhite)} 50%{border-right:0.6em solid transparent} }
        .splash-dos .splash-tag { color:var(--dos-yellow); letter-spacing:0.15em; }
        .splash-dos .splash-hint { color:var(--dos-cyan); }
        .splash-dos .splash-inner {
            background:var(--dos-black); border:2px solid var(--dos-lblue);
            text-shadow:none; box-shadow:4px 4px 0 rgba(0,0,0,0.5);
        }
{{template "navSharedCSSInner"}}
    </style>
</head>
<body>
    {{template "splashGate"}}
    <div id="splash-overlay" class="splash-overlay splash-dos" role="dialog" aria-modal="true" aria-label="Open microblog" tabindex="-1">
        <canvas class="splash-gl-canvas" id="splash-gl-canvas" aria-hidden="true"></canvas>
        <div class="splash-inner">
            <div class="splash-title">C:\SNONUX&gt;</div>
            <div class="splash-tag">MS-DOS v6.22</div>
            <div class="splash-hint">Press any key to continue...</div>
        </div>
    </div>
    <script>
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,drops=[],t0=performance.now();
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:false,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(1);
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(50,1,0.1,80);ca.position.z=20;
        var geo=new THREE.PlaneGeometry(0.22,0.32);
        for(var i=0;i<60;i++){
            var mat=new THREE.MeshBasicMaterial({color:0x55ff55,transparent:true,opacity:0.3+Math.random()*0.4});
            var m=new THREE.Mesh(geo,mat);
            m.position.set((Math.random()-0.5)*28, Math.random()*22-11, (Math.random()-0.5)*5);
            m.userData.speed=0.5+Math.random()*1.5;
            sc.add(m); drops.push(m);
        }
        sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);
            for(var i=0;i<drops.length;i++){
                drops[i].position.y-=drops[i].userData.speed*0.06;
                if(drops[i].position.y<-12) drops[i].position.y=12;
            }
            ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();
    </script>
    <canvas id="three-canvas"></canvas>
    <div class="overlay">
        <header>
            <div class="logo">
                <span class="logo-mark">C:\&gt;</span>
                <div class="logo-title">
                    <h1>SNONUX.FOO</h1>
                    <p class="subtitle">MICROBLOG &mdash; <a href="https://foo.zone">FOO.ZONE</a> IS THE REAL BLOG</p>
                    <p class="logo-host">Served by NetBSD on a Raspberry Pi 3</p>
                </div>
            </div>
            <div class="nav">
                <a href="atom.xml" class="header-feed-link" rel="alternate" title="Atom feed" type="application/atom+xml">Atom feed</a>
                <a href="https://foo.zone/about" class="transmit-btn">ABOUT</a>
            </div>
        </header>
        {{template "navhints" .}}
        <div class="content" id="post-content">
            {{range $i, $post := .Posts}}
            <div class="post" id="post-{{$post.ID}}" data-index="{{$i}}">
                <div class="post-header">
                    <div><strong>@SNONUX</strong></div>
                    <div class="post-time">{{$post.FormattedTime}}</div>
                </div>
                <div class="post-text">{{$post.ContentHTML}}</div>
            </div>
            {{end}}
        </div>
        {{if or .PrevPage .NextPage}}
        <footer class="page-nav-footer" aria-label="Pagination">
            <div class="page-nav page-nav-dual">
                {{if .PrevPage}}<a href="{{.PrevPage}}">&lt;-- NEWER</a>{{end}}
                {{if .NextPage}}<a href="{{.NextPage}}">OLDER --&gt;</a>{{end}}
            </div>
        </footer>
        {{end}}
    </div>
    {{template "navmodal" .}}
    <script>
    (function() {
        var scene, camera, renderer, clock;
        var columns = [];

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x000088);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 200);
            camera.position.set(0, 0, 40);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: false });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(1);
            clock = new THREE.Clock();

            var geo = new THREE.PlaneGeometry(0.35, 0.5);

            for (var c = 0; c < 30; c++) {
                var col = [];
                var x = (c - 15) * 2.2;
                var speed = 1.5 + Math.random() * 3;
                var startY = Math.random() * 60 - 30;
                for (var r = 0; r < 8; r++) {
                    var brightness = 1.0 - (r / 8) * 0.7;
                    var color = new THREE.Color(brightness * 0.33, brightness, brightness * 0.33);
                    var mat = new THREE.MeshBasicMaterial({ color: color, transparent: true, opacity: brightness * 0.5 });
                    var mesh = new THREE.Mesh(geo, mat);
                    mesh.position.set(x, startY - r * 0.7, 0);
                    scene.add(mesh);
                    col.push({ mesh: mesh, offset: r * 0.7 });
                }
                columns.push({ chars: col, x: x, speed: speed, y: startY });
            }

            window.addEventListener('resize', onResize);
            animate();
        }

        function onResize() {
            camera.aspect = window.innerWidth / window.innerHeight;
            camera.updateProjectionMatrix();
            renderer.setSize(window.innerWidth, window.innerHeight);
        }

        function animate() {
            requestAnimationFrame(animate);
            var t = clock.getElapsedTime();
            for (var c = 0; c < columns.length; c++) {
                var col = columns[c];
                var y = col.y - t * col.speed;
                y = ((y % 60) + 60) % 60 - 30;
                for (var r = 0; r < col.chars.length; r++) {
                    col.chars[r].mesh.position.y = y - col.chars[r].offset;
                }
            }
            renderer.render(scene, camera);
        }

        initThree();
    })();
    </script>
    {{template "navscript" .}}
</body>
</html>`
