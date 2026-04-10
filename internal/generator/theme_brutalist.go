package generator

// brutalistTemplate is a raw brutalist theme — pure black, thick white borders,
// Impact font, red as the only accent. WebGL scene: harsh slowly-rotating boxes,
// wireframe red and solid white, no fog — brutal clarity.
const brutalistTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SNONUX.FOO</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/three.js/r134/three.min.js"></script>
    <style>
        :root { --red:#ff2200; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:Impact,'Arial Narrow',Arial,sans-serif;
               background:#000; color:#fff; overflow:hidden; height:100vh; }
        #three-canvas { position:fixed; top:0; left:0; width:100%; height:100%; z-index:1; }
        .overlay { height:100vh; display:flex; flex-direction:column; position:relative; z-index:10; }
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
        a.header-feed-link { color:#aaa; font-family:'Courier New',monospace; font-size:0.78rem; }
        a.header-feed-link:hover { color:var(--red); }
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
        .page-nav-footer { flex-shrink:0; padding:6px 24px; display:flex; justify-content:center;
            background:#000; border-top:4px solid #fff; }
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
        .splash-overlay.splash-brutalist { background:#000; }
        .splash-brutalist .splash-frame {
            border:6px solid #fff; padding:clamp(1.5rem,5vw,2.5rem) clamp(1.25rem,4vw,2rem);
            box-shadow: 12px 12px 0 var(--red); animation: splashBrutalJolt 3s steps(2,end) infinite;
        }
        @keyframes splashBrutalJolt { 0%,100% { transform: translate(0,0); } 50% { transform: translate(2px,-2px); } }
        .splash-brutalist .splash-title { font-family:Impact,sans-serif; font-size:clamp(1.8rem,6vw,2.8rem); color:#fff; }
        .splash-brutalist .splash-tag { font-family:'Courier New',monospace; color:var(--red); }
        .splash-brutalist .splash-hint { font-family:'Courier New',monospace; color:#c8c8c8; }
        .splash-brutalist .splash-inner { text-shadow: 0 0 12px #000, 0 2px 8px #000; }
{{template "navSharedCSSInner"}}
    </style>
</head>
<body>
    {{template "splashGate"}}
    <div id="splash-overlay" class="splash-overlay splash-brutalist" role="dialog" aria-modal="true" aria-label="Open microblog" tabindex="-1">
        <canvas class="splash-gl-canvas" id="splash-gl-canvas" aria-hidden="true"></canvas>
        <div class="splash-inner splash-frame">
            <div class="splash-title">SNONUX.FOO</div>
            <div class="splash-tag">Brutalist theme</div>
            <div class="splash-hint">[ CLICK OR ENTER TO TRANSMIT ]</div>
        </div>
    </div>
    <script>
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,g=new THREE.Group(),t0=performance.now();
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(50,1,0.1,80);ca.position.z=8;
        var b1=new THREE.LineSegments(new THREE.EdgesGeometry(new THREE.BoxGeometry(3.4,2.4,2.4)),new THREE.LineBasicMaterial({color:0xffffff}));
        var b2=new THREE.LineSegments(new THREE.EdgesGeometry(new THREE.BoxGeometry(2.2,1.6,1.6)),new THREE.LineBasicMaterial({color:0xff2200}));
        b2.position.set(0.3,0.2,0.5);g.add(b1);g.add(b2);sc.add(g);sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);var t=(now-t0)*0.001;g.rotation.x=t*0.51;g.rotation.y=t*0.73;ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();
    </script>
    <canvas id="three-canvas"></canvas>
    <div class="overlay">
        <header>
            <div class="logo">
                <span class="logo-mark">SN</span>
                <div class="logo-title">
                    <h1>SNONUX.FOO</h1>
                    <p class="subtitle">MICROBLOG &mdash; <a href="https://foo.zone">FOO.ZONE</a> IS THE REAL BLOG</p>
                    <p class="logo-host">Site served by a Raspberry Pi 3</p>
                </div>
            </div>
            <div class="nav">
                <a href="atom.xml" class="header-feed-link" rel="alternate" title="Atom feed" type="application/atom+xml">Atom feed</a>
                <a href="https://foo.zone/about" class="transmit-btn">TRANSMIT</a>
            </div>
        </header>
        {{template "navhints" .}}
        <div class="content" id="post-content">
            {{range $i, $post := .Posts}}
            <div class="post" data-index="{{$i}}" onclick="selectPost({{$i}})">
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
                {{if .PrevPage}}<a href="{{.PrevPage}}">&larr; NEWER</a>{{end}}
                {{if .NextPage}}<a href="{{.NextPage}}">OLDER &rarr;</a>{{end}}
            </div>
        </footer>
        {{end}}
    </div>
    {{template "navmodal" .}}
    <script>
    // Brutalist WebGL: harsh slowly-rotating boxes — solid white and wireframe red.
    // No fog, no softness. Pure geometric violence against the black void.
    (function() {
        var scene, camera, renderer, clock;
        var boxes = [];

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x000000);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 200);
            camera.position.set(0, 0, 40);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: false });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();

            // Box configurations: [size, posX, posY, posZ, rotSpeedX, rotSpeedY, wireframe, color]
            var configs = [
                [10, 0,   0,  0,   0.002, 0.005, false, 0xffffff],
                [6,  18, -6,  -8,  0.004, 0.003, true,  0xff2200],
                [7,  -16, 5, -10,  0.003, 0.006, true,  0xff2200],
                [5,  8,  12, -5,   0.006, 0.002, false, 0xff2200],
                [4,  -10,-10, -3,  0.005, 0.004, false, 0xffffff],
            ];

            configs.forEach(function(c) {
                var geo = new THREE.BoxGeometry(c[0], c[0], c[0]);
                var mat = new THREE.MeshBasicMaterial({ color: c[7], wireframe: c[6] });
                var mesh = new THREE.Mesh(geo, mat);
                mesh.position.set(c[1], c[2], c[3]);
                boxes.push({ mesh: mesh, rx: c[4], ry: c[5] });
                scene.add(mesh);
            });

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
            boxes.forEach(function(b) {
                b.mesh.rotation.x += b.rx;
                b.mesh.rotation.y += b.ry;
            });
            renderer.render(scene, camera);
        }

        initThree();
    })();
    </script>
    {{template "navscript" .}}
</body>
</html>`
