package generator

// terminalTemplate is the green phosphor CRT terminal theme.
// Monospace throughout, scanline overlay via CSS, no external dependencies.
// WebGL scene: large IcosahedronGeometry wireframe with orbiting torus particles,
// giving a rotating phosphor-green 3D orb behind the terminal interface.
const terminalTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>snonux.foo // TERMINAL</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/three.js/r134/three.min.js"></script>
    <style>
        :root { --p:#33ff33; --dim:#1a7a1a; --bg:#0a0a0a; --bg2:#050505; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:'Courier New',Courier,monospace; background:var(--bg); color:var(--p);
               overflow:hidden; height:100vh; position:relative; }
        /* CRT scanlines sit above the WebGL canvas */
        body::before { content:''; position:fixed; inset:0; z-index:999; pointer-events:none;
            background:repeating-linear-gradient(0deg,transparent,transparent 2px,
                rgba(0,0,0,0.12) 2px,rgba(0,0,0,0.12) 4px); }
        /* Subtle screen flicker */
        @keyframes flicker { 0%,100%{opacity:1} 93%{opacity:0.97} 95%{opacity:0.91} 97%{opacity:0.98} }
        body { animation:flicker 9s infinite; }
        /* WebGL background canvas — fills the viewport behind everything */
        #three-canvas { position:fixed; top:0; left:0; width:100%; height:100%; z-index:1; }
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
        a.header-feed-link { color:var(--dim); }
        a.header-feed-link:hover { color:var(--p); }
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
        .page-nav-footer { flex-shrink:0; padding:6px 24px; display:flex; justify-content:center;
            background:var(--bg2); border-top:2px solid var(--p); }
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
        .splash-overlay.splash-terminal { background: var(--bg); font-family:'Courier New',monospace; }
        .splash-terminal .splash-prompt { text-align:left; font-size:0.9rem; color:rgba(51,255,51,0.78); margin-bottom:0.5rem; }
        .splash-terminal .splash-title { font-size:clamp(1.2rem,4vw,1.65rem); color:var(--p);
            text-shadow:0 0 12px var(--p); letter-spacing:0.15em; }
        .splash-terminal .splash-cursor::after { content:'█'; animation: splashTermBlink 1s step-end infinite; color:var(--p); }
        @keyframes splashTermBlink { 0%,100%{opacity:1} 50%{opacity:0} }
        .splash-terminal .splash-tag { color:rgba(51,255,51,0.85); letter-spacing:0.25em; }
        .splash-terminal .splash-hint { color:rgba(51,255,51,0.8); }
        .splash-terminal .splash-inner { text-shadow: 0 0 8px #000, 0 2px 12px #000; }
{{template "navSharedCSSInner"}}
    </style>
</head>
<body>
    {{template "splashGate"}}
    <div id="splash-overlay" class="splash-overlay splash-terminal" role="dialog" aria-modal="true" aria-label="Open microblog" tabindex="-1">
        <canvas class="splash-gl-canvas" id="splash-gl-canvas" aria-hidden="true"></canvas>
        <div class="splash-inner">
            <div class="splash-prompt">&gt; ./snonux --boot</div>
            <div class="splash-title splash-cursor">LINK ESTABLISHED</div>
            <div class="splash-tag">TERMINAL SESSION</div>
            <div class="splash-hint">[ click / enter to continue ]</div>
        </div>
    </div>
    <script>
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,m,t0=performance.now();
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(48,1,0.1,60);ca.position.z=7;
        m=new THREE.Mesh(new THREE.IcosahedronGeometry(2.3,1),new THREE.MeshBasicMaterial({color:0x33ff33,wireframe:true,transparent:true,opacity:0.88}));
        sc.add(m);sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);var t=(now-t0)*0.001;m.rotation.x=t*0.62;m.rotation.y=t*0.88;ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();
    </script>
    <canvas id="three-canvas"></canvas>
    <div class="overlay">
        <header>
            <div class="logo">
                <span class="logo-mark">[SN]</span>
                <div class="logo-title">
                    <h1>snonux.foo</h1>
                    <p class="subtitle">microblog / <a href="https://foo.zone">foo.zone</a> is the real blog</p>
                    <p class="logo-host">Served by NetBSD on a Raspberry Pi 3</p>
                </div>
            </div>
            <div class="nav">
                <a href="atom.xml" class="header-feed-link" rel="alternate" title="Atom feed" type="application/atom+xml">atom.xml</a>
                <a href="https://foo.zone/about" class="transmit-btn">&gt; TRANSMIT</a>
            </div>
        </header>
        {{template "navhints" .}}
        <div class="content" id="post-content">
            {{range $i, $post := .Posts}}
            <div class="post" id="post-{{$post.ID}}" data-index="{{$i}}">
                <div class="post-header">
                    <div><strong>@snonux</strong></div>
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
    // Terminal WebGL scene: phosphor-green icosahedron wireframe + torus particle ring.
    // The scene sits behind the CRT scanline overlay (z-index:999) and the UI (z-index:10).
    (function() {
        var scene, camera, renderer, icosa, particles;
        var clock = new THREE.Clock();

        function initThree() {
            // Scene with pure-black background and distance fog
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x000000);
            scene.fog = new THREE.Fog(0x000000, 20, 80);

            // Perspective camera positioned in front of the orb
            camera = new THREE.PerspectiveCamera(60, window.innerWidth / window.innerHeight, 0.1, 200);
            camera.position.set(0, 0, 30);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

            // Large green phosphor wireframe icosahedron — the central CRT orb
            var icoGeo = new THREE.IcosahedronGeometry(8, 2);
            var icoMat = new THREE.MeshBasicMaterial({ color: 0x33ff33, wireframe: true });
            icosa = new THREE.Mesh(icoGeo, icoMat);
            scene.add(icosa);

            // 400 dim particles arranged on a torus path around the icosahedron
            var torusGeo = new THREE.TorusGeometry(14, 3, 16, 100);
            var positions = torusGeo.attributes.position;
            var ptGeo = new THREE.BufferGeometry();
            var pts = new Float32Array(400 * 3);
            for (var i = 0; i < 400; i++) {
                // Sample vertices from the torus geometry to place particles on its surface
                var idx = Math.floor(Math.random() * positions.count);
                pts[i * 3]     = positions.getX(idx);
                pts[i * 3 + 1] = positions.getY(idx);
                pts[i * 3 + 2] = positions.getZ(idx);
            }
            ptGeo.setAttribute('position', new THREE.BufferAttribute(pts, 3));
            var ptMat = new THREE.PointsMaterial({ color: 0x1a7a1a, size: 0.18 });
            particles = new THREE.Points(ptGeo, ptMat);
            scene.add(particles);

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
            // Slow multi-axis rotation for the phosphor orb
            icosa.rotation.x = t * 0.12;
            icosa.rotation.y = t * 0.18;
            icosa.rotation.z = t * 0.07;
            // Counter-rotate particles for visual contrast
            particles.rotation.y = -t * 0.08;
            particles.rotation.x =  t * 0.04;
            renderer.render(scene, camera);
        }

        initThree();
    })();
    </script>
    {{template "navscript" .}}
</body>
</html>`
