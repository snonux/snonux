package generator

// glassTemplate is a glassmorphism theme — semi-transparent frosted panels over
// a WebGL background of slowly drifting crystal icosahedron shards.
// CSS gradient blobs are replaced by the WebGL scene.
const glassTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>snonux.foo · glass</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/three.js/r134/three.min.js"></script>
    <style>
        :root { --blue:#6366f1; --purple:#a855f7; --pink:#ec4899; --text:#1e1b4b; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:'Segoe UI',system-ui,sans-serif; overflow:hidden; height:100vh;
               background:#f0f4ff; color:var(--text); }
        #three-canvas { position:fixed; top:0; left:0; width:100%; height:100%; z-index:1; }
        .overlay { position:relative; z-index:10; height:100vh; display:flex; flex-direction:column; }
        header { padding:16px 28px; background:rgba(255,255,255,0.55); backdrop-filter:blur(20px);
                 border-bottom:1px solid rgba(255,255,255,0.6); display:flex; align-items:center; justify-content:space-between;
                 box-shadow:0 2px 12px rgba(99,102,241,0.08); }
        .logo { display:flex; align-items:center; gap:14px; }
        .logo-mark { font-size:2rem; font-weight:800;
                     background:linear-gradient(135deg,var(--blue),var(--purple));
                     -webkit-background-clip:text; -webkit-text-fill-color:transparent; }
        .logo-title h1 { font-size:1.5rem; font-weight:700; color:var(--text); }
        .logo-title .subtitle { font-size:0.75rem; color:#6b7280; margin-top:1px; }
        .logo-title .subtitle a { color:var(--blue); text-decoration:none; }
        .logo-title .subtitle a:hover { text-decoration:underline; }
        .transmit-btn { border:1px solid rgba(99,102,241,0.4); color:var(--blue); padding:9px 20px;
                        border-radius:20px; text-decoration:none; font-size:0.85rem;
                        background:rgba(255,255,255,0.5); backdrop-filter:blur(8px); transition:all 0.2s; }
        .transmit-btn:hover { background:var(--blue); color:#fff; border-color:var(--blue); }
        .nav-hints { background:rgba(255,255,255,0.35); backdrop-filter:blur(10px);
                     border-bottom:1px solid rgba(255,255,255,0.5); color:#6b7280;
                     padding:4px 28px; display:flex; gap:18px; font-size:0.68rem; flex-wrap:wrap; }
        .nav-hints kbd { background:rgba(255,255,255,0.7); border:1px solid rgba(99,102,241,0.25);
                         color:var(--blue); border-radius:4px; padding:0 5px; margin:0 2px; font-size:0.68rem; }
        .content { flex:1; overflow-y:auto; padding:20px 28px;
                   scrollbar-width:thin; scrollbar-color:rgba(99,102,241,0.4) transparent; }
        .page-nav { display:flex; justify-content:center; margin:14px 0; }
        .page-nav a { border:1px solid rgba(99,102,241,0.35); color:var(--blue); padding:8px 20px;
                      border-radius:20px; text-decoration:none; font-size:0.82rem;
                      background:rgba(255,255,255,0.45); backdrop-filter:blur(8px); }
        .page-nav a:hover { background:var(--blue); color:#fff; }
        .post { background:rgba(255,255,255,0.45); backdrop-filter:blur(18px);
                border:1px solid rgba(255,255,255,0.6); border-radius:14px;
                padding:22px; margin-bottom:14px; cursor:pointer;
                box-shadow:0 4px 20px rgba(99,102,241,0.08); transition:all 0.25s; }
        .post:hover { background:rgba(255,255,255,0.6); box-shadow:0 8px 30px rgba(99,102,241,0.18); transform:translateY(-2px); }
        .post-active { border-color:var(--blue) !important; background:rgba(238,240,255,0.75) !important;
                       box-shadow:0 0 0 2px rgba(99,102,241,0.3),0 8px 30px rgba(99,102,241,0.2),
                                  inset 3px 0 0 var(--blue) !important; }
        .post-header { display:flex; justify-content:space-between; margin-bottom:12px; font-size:0.88rem; }
        .post-time { color:#9ca3af; font-family:monospace; font-size:0.8rem; }
        .post-text { line-height:1.65; font-size:0.95rem; }
        .post-text a { color:var(--blue); text-decoration:none; }
        .post-text a:hover { text-decoration:underline; }
        .post-audio { width:100%; margin-top:10px; }
        .post-modal { display:none; position:fixed; inset:0; z-index:100;
                      background:rgba(240,244,255,0.85); backdrop-filter:blur(28px);
                      overflow-y:auto; padding:40px 20px; }
        .post-modal.active { display:block; }
        .modal-inner { max-width:760px; margin:0 auto; background:rgba(255,255,255,0.7);
                       backdrop-filter:blur(24px); border:1px solid rgba(255,255,255,0.75);
                       border-radius:16px; box-shadow:0 20px 60px rgba(99,102,241,0.18); padding:40px; }
        .modal-close { float:right; background:none; border:none; color:#9ca3af;
                       font-size:0.9rem; cursor:pointer; }
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
    // Glass WebGL: 8 crystal icosahedron shards, each rendered as a semi-transparent
    // solid mesh with a wireframe overlay, drifting and rotating slowly in the light bg.
    // alpha:true so the light body background (#f0f4ff) shows through the canvas.
    (function() {
        var scene, camera, renderer, clock;
        var shards = [];

        function initThree() {
            scene = new THREE.Scene();

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 200);
            camera.position.set(0, 0, 40);

            renderer = new THREE.WebGLRenderer({
                canvas: document.getElementById('three-canvas'),
                antialias: true, alpha: true
            });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            renderer.setClearColor(0xf0f4ff, 1);
            clock = new THREE.Clock();

            scene.add(new THREE.AmbientLight(0xffffff, 0.8));
            var pl = new THREE.PointLight(0x7c3aed, 2, 80);
            pl.position.set(10, 10, 10);
            scene.add(pl);

            var colors  = [0x6366f1, 0xa855f7, 0xec4899, 0x6366f1, 0x818cf8, 0xc084fc, 0xf472b6, 0x6366f1];
            var sizes   = [5, 3.5, 2.5, 4, 3, 2, 4.5, 2.2];
            var details = [2, 1, 0, 2, 1, 0, 1, 2];
            var positions = [
                [-12, 6,-8], [14,-4,-12], [0,12,-6], [-8,-10,-15],
                [16, 8,-4], [-14,2,-10], [6,-8,-5], [-4,10,-14]
            ];

            for (var i = 0; i < 8; i++) {
                var geo = new THREE.IcosahedronGeometry(sizes[i], details[i]);
                var solid = new THREE.Mesh(geo, new THREE.MeshPhongMaterial({
                    color: colors[i], transparent: true, opacity: 0.12, side: THREE.DoubleSide
                }));
                var wire = new THREE.Mesh(geo, new THREE.MeshBasicMaterial({
                    color: colors[i], wireframe: true, transparent: true, opacity: 0.28
                }));
                solid.position.set(positions[i][0], positions[i][1], positions[i][2]);
                wire.position.copy(solid.position);

                var rx = (Math.random() - 0.5) * 0.004;
                var ry = (Math.random() - 0.5) * 0.006;
                var rz = (Math.random() - 0.5) * 0.003;
                shards.push({ solid: solid, wire: wire, rx: rx, ry: ry, rz: rz });
                scene.add(solid);
                scene.add(wire);
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
            shards.forEach(function(s) {
                s.solid.rotation.x += s.rx; s.solid.rotation.y += s.ry; s.solid.rotation.z += s.rz;
                s.wire.rotation.copy(s.solid.rotation);
            });
            renderer.render(scene, camera);
        }

        initThree();
    })();
    </script>
    {{template "navscript" .}}
</body>
</html>`
