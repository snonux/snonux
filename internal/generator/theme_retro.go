package generator

// retroTemplate is an amber DOS terminal theme — black background, amber
// phosphor (#ffb000) text, monospace throughout, no decorations.
// Distinct from terminal.go (green) — this one evokes vintage PC monitors.
// WebGL scene: spinning amber wireframe cube with orbiting octahedrons and
// dim amber star particles for a retro demo-scene feel.
const retroTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SNONUX.FOO // RETRO</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/three.js/r134/three.min.js"></script>
    <style>
        :root { --amber:#ffb000; --dim:#7a5200; --bg:#0a0800; --bg2:#050300; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:'Courier New',Courier,monospace; background:var(--bg); color:var(--amber);
               overflow:hidden; height:100vh; }
        /* Phosphor scanlines overlay — sits above WebGL */
        body::before { content:''; position:fixed; inset:0; z-index:999; pointer-events:none;
            background:repeating-linear-gradient(0deg,transparent,transparent 2px,
                rgba(0,0,0,0.15) 2px,rgba(0,0,0,0.15) 4px); }
        /* Subtle glow flicker */
        @keyframes amber-flicker { 0%,100%{opacity:1} 94%{opacity:0.98} 96%{opacity:0.93} }
        body { animation:amber-flicker 11s infinite; }
        /* WebGL background canvas */
        #three-canvas { position:fixed; top:0; left:0; width:100%; height:100%; z-index:1; }
        .overlay { position:relative; z-index:10; height:100vh; display:flex; flex-direction:column; }
        header { padding:12px 24px; background:var(--bg2); border-bottom:2px solid var(--amber);
                 display:flex; align-items:center; justify-content:space-between; }
        .logo { display:flex; align-items:center; gap:14px; }
        .logo-mark { font-size:1.6rem; color:var(--amber); text-shadow:0 0 14px var(--amber);
                     letter-spacing:2px; }
        .logo-title h1 { font-size:1.2rem; color:var(--amber); text-shadow:0 0 10px var(--amber);
                         letter-spacing:4px; font-weight:normal; }
        .logo-title .subtitle { font-size:0.72rem; color:var(--dim); margin-top:2px; }
        .logo-title .subtitle a { color:var(--amber); text-decoration:none; }
        .logo-title .subtitle a:hover { text-shadow:0 0 6px var(--amber); }
        .transmit-btn { border:1px solid var(--amber); color:var(--amber); padding:8px 18px;
                        text-decoration:none; font-size:0.82rem; letter-spacing:2px;
                        transition:all 0.1s; }
        .transmit-btn:hover { background:var(--amber); color:var(--bg); }
        a.header-feed-link { color:var(--dim); }
        a.header-feed-link:hover { color:var(--amber); text-shadow:0 0 6px var(--amber); }
        .nav-hints { background:var(--bg2); border-bottom:1px solid var(--dim); color:var(--dim);
                     padding:4px 24px; display:flex; gap:18px; font-size:0.68rem; flex-wrap:wrap; }
        .nav-hints kbd { background:transparent; border:1px solid var(--dim); color:var(--amber);
                         padding:0 5px; font-size:0.68rem; margin:0 2px; }
        .content { flex:1; overflow-y:auto; padding:14px 24px;
                   scrollbar-width:thin; scrollbar-color:var(--dim) var(--bg); }
        .page-nav { display:flex; justify-content:center; margin:12px 0; }
        .page-nav a { border:1px solid var(--dim); color:var(--amber); padding:7px 20px;
                      text-decoration:none; font-size:0.82rem; letter-spacing:2px; }
        .page-nav a:hover { background:var(--amber); color:var(--bg); border-color:var(--amber); }
        .page-nav-footer { flex-shrink:0; padding:6px 24px; display:flex; justify-content:center;
            background:var(--bg2); border-top:2px solid var(--amber); }
        .post { background:var(--bg); border:1px solid var(--dim); padding:16px 18px;
                margin-bottom:10px; cursor:pointer; transition:border-color 0.15s; }
        .post:hover { border-color:var(--amber); box-shadow:0 0 8px rgba(255,176,0,0.25); }
        .post-active { border-color:var(--amber) !important;
                       background:rgba(255,176,0,0.04) !important;
                       box-shadow:0 0 14px rgba(255,176,0,0.3),inset 3px 0 0 var(--amber) !important; }
        .post-header { display:flex; justify-content:space-between; margin-bottom:10px; font-size:0.85rem; }
        .post-time { color:var(--dim); font-size:0.78rem; }
        .post-text { line-height:1.6; font-size:0.88rem; }
        .post-text a { color:var(--amber); text-decoration:underline; }
        .post-image { max-width:100%; margin-top:10px; border:1px solid var(--dim);
                      filter:sepia(60%) hue-rotate(-10deg); }
        .post-audio { width:100%; margin-top:10px; }
        .post-modal { display:none; position:fixed; inset:0; z-index:100;
                      background:rgba(0,0,0,0.97); overflow-y:auto; padding:40px 20px; }
        .post-modal.active { display:block; }
        .modal-inner { max-width:740px; margin:0 auto; background:var(--bg);
                       border:1px solid var(--amber); padding:36px;
                       box-shadow:0 0 40px rgba(255,176,0,0.2); }
        .modal-close { float:right; background:none; border:none; color:var(--dim);
                       font-family:monospace; font-size:0.9rem; cursor:pointer; letter-spacing:2px; }
        @media(max-width:640px) { .nav-hints{display:none;} header{padding:10px 16px;} .content{padding:10px 16px;} }
        .splash-overlay.splash-retro { background: var(--bg); font-family:'Courier New',monospace; }
        .splash-retro::after {
            content:''; position:absolute; inset:0; pointer-events:none; opacity:0.35;
            background: repeating-linear-gradient(0deg, transparent, transparent 2px, rgba(0,0,0,0.2) 2px, rgba(0,0,0,0.2) 4px);
        }
        .splash-retro .splash-inner { position:relative; z-index:1; }
        .splash-retro .splash-title {
            font-size:clamp(1.15rem,3.8vw,1.55rem); color:var(--amber);
            text-shadow:0 0 14px var(--amber); letter-spacing:0.3em;
            animation: splashRetroFlicker 4s ease-in-out infinite;
        }
        @keyframes splashRetroFlicker { 0%,100%{opacity:1} 50%{opacity:0.92} }
        .splash-retro .splash-tag { color:#d4a020; }
        .splash-retro .splash-hint { color:#c99528; }
        .splash-retro .splash-inner { text-shadow: 0 0 10px #000, 0 2px 8px #000; }
    </style>
</head>
<body>
    {{template "splashGate"}}
    <div id="splash-overlay" class="splash-overlay splash-retro" tabindex="-1" aria-label="Open microblog">
        <canvas class="splash-gl-canvas" id="splash-gl-canvas" aria-hidden="true"></canvas>
        <div class="splash-inner">
            <div class="splash-title">*** SNONUX BBS ***</div>
            <div class="splash-tag">Amber phosphor mode</div>
            <div class="splash-hint">Press Enter or click to connect</div>
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
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(48,1,0.1,60);ca.position.z=7.5;
        var bx=new THREE.Mesh(new THREE.BoxGeometry(2.6,2.6,2.6),new THREE.MeshBasicMaterial({color:0xffb000,wireframe:true,transparent:true,opacity:0.9}));
        var oc=new THREE.Mesh(new THREE.OctahedronGeometry(1.35,0),new THREE.MeshBasicMaterial({color:0xffb000,wireframe:true,transparent:true,opacity:0.55}));
        g.add(bx);g.add(oc);sc.add(g);sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);var t=(now-t0)*0.001;g.rotation.x=t*0.44;g.rotation.y=t*0.71;oc.rotation.z=t*0.9;ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();
    </script>
    <canvas id="three-canvas"></canvas>
    <div class="overlay">
        <header>
            <div class="logo">
                <span class="logo-mark">[SN]</span>
                <div class="logo-title">
                    <h1>SNONUX.FOO</h1>
                    <p class="subtitle">MICROBLOG / <a href="https://foo.zone">FOO.ZONE</a> IS THE REAL BLOG</p>
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
                {{if .PrevPage}}<a href="{{.PrevPage}}">&lt;-- NEWER</a>{{end}}
                {{if .NextPage}}<a href="{{.NextPage}}">OLDER --&gt;</a>{{end}}
            </div>
        </footer>
        {{end}}
    </div>
    {{template "navmodal" .}}
    <script>
    // Retro WebGL scene: amber demo-scene cube + orbiting octahedrons + star particles.
    // Evokes classic 80s/90s PC demo aesthetics with amber phosphor colours.
    (function() {
        var scene, camera, renderer;
        var mainCube, orbiters = [];
        var clock = new THREE.Clock();

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x0a0800);
            scene.fog = new THREE.Fog(0x0a0800, 25, 90);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth / window.innerHeight, 0.1, 200);
            camera.position.set(0, 0, 35);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

            // Central amber wireframe cube — the demo-scene hero piece
            var boxGeo = new THREE.BoxGeometry(8, 8, 8);
            var boxMat = new THREE.MeshBasicMaterial({ color: 0xffb000, wireframe: true });
            mainCube = new THREE.Mesh(boxGeo, boxMat);
            scene.add(mainCube);

            // 6 small octahedron wireframes evenly spaced on an orbital ring
            var octoMat = new THREE.MeshBasicMaterial({ color: 0xffb000, wireframe: true });
            for (var i = 0; i < 6; i++) {
                var octoGeo = new THREE.OctahedronGeometry(1.5);
                var octo = new THREE.Mesh(octoGeo, octoMat.clone());
                var angle = (i / 6) * Math.PI * 2;
                octo.position.set(Math.cos(angle) * 18, Math.sin(angle) * 5, Math.sin(angle) * 18);
                orbiters.push({ mesh: octo, baseAngle: angle });
                scene.add(octo);
            }

            // 800 dim amber background star particles
            var starGeo = new THREE.BufferGeometry();
            var starPos = new Float32Array(800 * 3);
            for (var j = 0; j < 800 * 3; j++) {
                starPos[j] = (Math.random() - 0.5) * 120;
            }
            starGeo.setAttribute('position', new THREE.BufferAttribute(starPos, 3));
            var starMat = new THREE.PointsMaterial({ color: 0x7a5200, size: 0.15 });
            scene.add(new THREE.Points(starGeo, starMat));

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
            // Main cube rotates on Y and X axes for classic demo-scene look
            mainCube.rotation.y = t * 0.35;
            mainCube.rotation.x = t * 0.18;
            // Orbiters revolve around the central cube and spin individually
            for (var i = 0; i < orbiters.length; i++) {
                var o = orbiters[i];
                var angle = o.baseAngle + t * 0.4;
                o.mesh.position.x = Math.cos(angle) * 18;
                o.mesh.position.z = Math.sin(angle) * 18;
                o.mesh.position.y = Math.sin(angle * 0.7) * 4;
                o.mesh.rotation.x = t * 0.9;
                o.mesh.rotation.z = t * 0.6;
            }
            renderer.render(scene, camera);
        }

        initThree();
    })();
    </script>
    {{template "navscript" .}}
</body>
</html>`
