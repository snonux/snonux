package generator

// synthwaveTemplate is the 80s retrowave theme — dark purple sky, WebGL
// perspective grid with a glowing sun and scan-line rings, hot pink/orange
// accents, Russo One font. The CSS sky and grid-floor divs are replaced by WebGL.
const synthwaveTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>snonux.foo ⊕ SYNTHWAVE</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link href="https://fonts.googleapis.com/css2?family=Russo+One&family=Share+Tech+Mono&display=swap" rel="stylesheet">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/three.js/r134/three.min.js"></script>
    <style>
        :root { --pink:#ff2d78; --purple:#bf3fff; --orange:#ff6b2b; --bg:#0d0221; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:'Russo One','Arial Black',sans-serif; background:var(--bg);
               color:#fff; overflow:hidden; height:100vh; }
        #three-canvas { position:fixed; top:0; left:0; width:100%; height:100%; z-index:1; }
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
        .post-audio { width:100%; margin-top:10px; }
        .post-modal { display:none; position:fixed; inset:0; z-index:100;
                      background:rgba(13,2,33,0.96); overflow-y:auto; padding:40px 20px; }
        .post-modal.active { display:block; }
        .modal-inner { max-width:780px; margin:0 auto; background:rgba(20,5,50,0.98);
                       border:2px solid var(--pink); border-radius:6px;
                       box-shadow:0 0 60px rgba(255,45,120,0.35); padding:38px; }
        .modal-close { float:right; background:none; border:none; color:var(--orange);
                       font-family:'Russo One',sans-serif; font-size:0.9rem; cursor:pointer; letter-spacing:2px; }
        @media(max-width:640px) { .nav-hints{display:none;} header{padding:12px 18px;} }
    </style>
</head>
<body>
    <canvas id="three-canvas"></canvas>
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
    <script>
    // Synthwave WebGL: glowing sunset sphere with horizontal scan-line rings,
    // a receding grid floor, and pink star particles. Replaces CSS sky/grid.
    (function() {
        var scene, camera, renderer, clock;
        var sun, sunRings = [];

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x0d0221);
            scene.fog = new THREE.Fog(0x0d0221, 60, 180);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 300);
            camera.position.set(0, 10, 45);
            camera.lookAt(0, -5, -10);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();

            // Glowing orange sunset sphere
            sun = new THREE.Mesh(
                new THREE.SphereGeometry(12, 32, 16),
                new THREE.MeshBasicMaterial({ color: 0xff6b2b })
            );
            sun.position.set(0, -8, -55);
            scene.add(sun);

            // Horizontal scan-line rings — alternating pink/purple, stacked on the sun
            var ringColors = [0xff2d78, 0xbf3fff, 0xff2d78, 0xbf3fff, 0xff2d78, 0xbf3fff, 0xff2d78, 0xbf3fff];
            for (var i = 0; i < 8; i++) {
                var ring = new THREE.Mesh(
                    new THREE.TorusGeometry(13 + i * 1.2, 0.09, 8, 64),
                    new THREE.MeshBasicMaterial({ color: ringColors[i] })
                );
                ring.position.copy(sun.position);
                ring.position.y += -4 + i * 1.1;
                scene.add(ring);
                sunRings.push(ring);
            }

            // Receding grid floor
            var grid = new THREE.GridHelper(200, 40, 0xff2d78, 0x4a0060);
            grid.position.set(0, -18, -30);
            scene.add(grid);

            // 1200 star particles scattered in a sphere shell
            var starPos = new Float32Array(1200 * 3);
            for (var j = 0; j < 1200 * 3; j += 3) {
                var r = 80 + Math.random() * 40;
                var theta = Math.random() * Math.PI * 2;
                var phi = Math.acos(2 * Math.random() - 1);
                starPos[j]   = r * Math.sin(phi) * Math.cos(theta);
                starPos[j+1] = r * Math.sin(phi) * Math.sin(theta);
                starPos[j+2] = r * Math.cos(phi);
            }
            var starGeo = new THREE.BufferGeometry();
            starGeo.setAttribute('position', new THREE.BufferAttribute(starPos, 3));
            scene.add(new THREE.Points(starGeo, new THREE.PointsMaterial({
                color: 0xff88aa, size: 0.2, transparent: true, opacity: 0.7
            })));

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
            // Sun pulses subtly; camera drifts sideways for parallax
            var pulse = 1 + 0.015 * Math.sin(t * 1.5);
            sun.scale.setScalar(pulse);
            sunRings.forEach(function(r) { r.scale.setScalar(pulse); });
            camera.position.x = Math.sin(t * 0.08) * 4;
            renderer.render(scene, camera);
        }

        initThree();
    })();
    </script>
    {{template "navscript" .}}
</body>
</html>`
