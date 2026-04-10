package generator

// oceanTemplate is a deep-ocean theme — dark navy/midnight blue background,
// WebGL animated wave surface with per-vertex sine displacement, teal/aqua accents.
const oceanTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>snonux.foo ~ OCEAN</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/three.js/r134/three.min.js"></script>
    <style>
        :root { --teal:#00b4d8; --aqua:#48cae4; --deep:#023e8a; --navy:#03045e; --foam:#caf0f8; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:'Segoe UI',system-ui,sans-serif; background:var(--navy);
               color:var(--foam); overflow:hidden; height:100vh; }
        #three-canvas { position:fixed; top:0; left:0; width:100%; height:100%; z-index:1; }
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
                        border-radius:20px; text-decoration:none; font-size:0.85rem; transition:all 0.2s; }
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
    // Ocean WebGL: a large PlaneGeometry wave surface whose vertices are displaced
    // each frame by two overlapping sine functions, lit by a moving teal point light.
    (function() {
        var scene, camera, renderer, clock;
        var waveMesh, waveGeo, pointLight;

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x03045e);
            scene.fog = new THREE.Fog(0x03045e, 30, 120);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 200);
            camera.position.set(0, 25, 50);
            camera.lookAt(0, 0, 0);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();

            // Wave surface — high segment count so vertex displacement looks smooth
            waveGeo = new THREE.PlaneGeometry(200, 200, 80, 80);
            waveMesh = new THREE.Mesh(waveGeo, new THREE.MeshPhongMaterial({
                color: 0x0077b6, emissive: 0x023e8a, emissiveIntensity: 0.3,
                transparent: true, opacity: 0.85, side: THREE.DoubleSide
            }));
            waveMesh.rotation.x = -Math.PI / 2;
            waveMesh.position.y = -5;
            scene.add(waveMesh);

            // Moving teal light circling above the wave
            pointLight = new THREE.PointLight(0x48cae4, 2, 80);
            pointLight.position.set(0, 20, 10);
            scene.add(pointLight);
            scene.add(new THREE.AmbientLight(0x023e8a, 0.6));

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
            var pos = waveGeo.attributes.position;

            // Two overlapping sine waves produce realistic ocean surface chop
            for (var i = 0; i < pos.count; i++) {
                var x = pos.getX(i);
                var z = pos.getZ(i);
                pos.setY(i,
                    Math.sin(x * 0.05 + t * 1.2) * 1.8 +
                    Math.cos(z * 0.07 + t * 0.9) * 1.4
                );
            }
            pos.needsUpdate = true;
            waveGeo.computeVertexNormals();

            // Light orbits lazily
            pointLight.position.x = Math.cos(t * 0.3) * 30;
            pointLight.position.z = Math.sin(t * 0.3) * 30;

            renderer.render(scene, camera);
        }

        initThree();
    })();
    </script>
    {{template "navscript" .}}
</body>
</html>`
