package generator

// auroraTemplate is a dark navy theme with a WebGL aurora borealis effect —
// six wide ribbon meshes whose vertices are animated with overlapping sine waves,
// rendered with additive blending to create the characteristic glow.
const auroraTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>snonux.foo ✦ AURORA</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/three.js/r134/three.min.js"></script>
    <style>
        :root { --green:#00ffb3; --teal:#00cfe8; --purple:#c084fc; --navy:#050d1a; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:'Segoe UI',system-ui,sans-serif; background:var(--navy);
               color:#e0f8f0; overflow:hidden; height:100vh; }
        #three-canvas { position:fixed; top:0; left:0; width:100%; height:100%; z-index:1; }
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
                        border-radius:20px; text-decoration:none; font-size:0.85rem; transition:all 0.2s; }
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
    <script>
    // Aurora WebGL: six wide ribbon meshes whose top-row vertices are animated
    // with overlapping sine waves, rendered with additive blending so they glow
    // against the dark navy sky like real aurora curtains.
    (function() {
        var RIBBON_COUNT = 6;
        var SEG_W = 60; // horizontal segments per ribbon
        var ribbonColors = [0x00ffb3, 0x00cfe8, 0xc084fc, 0x00ffb3, 0x48e8d0, 0xa855f7];
        var ribbonY     = [-10, -4, 2, 8, 14, 20];
        var ribbonZ     = [-40, -30, -22, -15, -10, -5];
        var ribbonFreq  = [0.6, 0.9, 0.7, 1.1, 0.5, 0.8];
        var ribbonPhase = [0.0, 1.2, 2.4, 0.8, 3.1, 1.7];
        var ribbonAmp   = [3.0, 2.5, 2.0, 3.5, 2.2, 2.8];

        var scene, camera, renderer, clock;
        var ribbons = [];

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x050d1a);
            scene.fog = new THREE.Fog(0x050d1a, 40, 120);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 200);
            camera.position.set(0, 5, 30);
            camera.lookAt(0, 5, 0);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

            clock = new THREE.Clock();

            for (var r = 0; r < RIBBON_COUNT; r++) {
                // Wide shallow plane; we animate the top row of vertices
                var geo = new THREE.PlaneGeometry(120, 8, SEG_W, 1);
                var mat = new THREE.MeshBasicMaterial({
                    color: ribbonColors[r], transparent: true, opacity: 0.32,
                    side: THREE.DoubleSide, blending: THREE.AdditiveBlending, depthWrite: false
                });
                var mesh = new THREE.Mesh(geo, mat);
                mesh.position.set(0, ribbonY[r], ribbonZ[r]);
                scene.add(mesh);
                ribbons.push({ mesh: mesh, geo: geo, freq: ribbonFreq[r],
                               phase: ribbonPhase[r], amp: ribbonAmp[r] });
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

            for (var r = 0; r < ribbons.length; r++) {
                var rb = ribbons[r];
                var pos = rb.geo.attributes.position;
                var count = pos.count;
                // PlaneGeometry vertices: (SEG_W+1)*2 total; top row is every other vertex
                for (var i = 0; i < count; i++) {
                    var x = pos.getX(i);
                    // Only animate top row (y > 0 in local space) for the waving top edge
                    if (pos.getY(i) > 0) {
                        pos.setY(i, rb.amp * Math.sin(t * rb.freq + x * 0.08 + rb.phase)
                                    + rb.amp * 0.4 * Math.cos(t * rb.freq * 0.7 + x * 0.05));
                    }
                }
                pos.needsUpdate = true;
            }
            renderer.render(scene, camera);
        }

        initThree();
    })();
    </script>
    {{template "navscript" .}}
</body>
</html>`
