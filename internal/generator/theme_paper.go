package generator

// volcanoTemplate is a dark volcanic theme — ember and spark particles rise from
// below the screen, glowing orange/red/yellow with additive blending, set against
// a deep dark-rock background with a warm lava glow at the horizon.
const volcanoTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>snonux.foo ▲ VOLCANO</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/three.js/r134/three.min.js"></script>
    <style>
        :root { --lava:#ff4400; --ember:#ff8c00; --hot:#ffcc00; --bg:#0d0802; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:'Segoe UI',system-ui,sans-serif; background:var(--bg);
               color:#ffe8cc; overflow:hidden; height:100vh; }
        #three-canvas { position:fixed; top:0; left:0; width:100%; height:100%; z-index:1; }
        .overlay { position:relative; z-index:10; height:100vh; display:flex; flex-direction:column; }
        header { padding:16px 28px; background:rgba(13,8,2,0.82); backdrop-filter:blur(12px);
                 border-bottom:1px solid rgba(255,68,0,0.3); display:flex; align-items:center; justify-content:space-between; }
        .logo { display:flex; align-items:center; gap:14px; }
        .logo-mark { font-size:2rem; font-weight:800; color:var(--ember); text-shadow:0 0 16px var(--lava); }
        .logo-title h1 { font-size:1.5rem; font-weight:700; color:#ffe8cc; }
        .logo-title .subtitle { font-size:0.75rem; color:rgba(255,232,204,0.5); margin-top:2px; }
        .logo-title .subtitle a { color:var(--ember); text-decoration:none; }
        .logo-title .subtitle a:hover { text-shadow:0 0 8px var(--lava); }
        .transmit-btn { border:1px solid var(--lava); color:var(--lava); padding:9px 20px;
                        border-radius:4px; text-decoration:none; font-size:0.85rem; transition:all 0.2s; }
        .transmit-btn:hover { background:var(--lava); color:var(--bg); }
        .nav-hints { background:rgba(13,8,2,0.7); border-bottom:1px solid rgba(255,68,0,0.15);
                     color:rgba(255,232,204,0.4); padding:5px 28px; display:flex; gap:18px;
                     font-size:0.68rem; flex-wrap:wrap; }
        .nav-hints kbd { background:rgba(255,68,0,0.12); border:1px solid rgba(255,68,0,0.35);
                         color:var(--ember); border-radius:3px; padding:0 5px; margin:0 2px; }
        .content { flex:1; overflow-y:auto; padding:20px 28px;
                   scrollbar-width:thin; scrollbar-color:var(--lava) var(--bg); }
        .page-nav { display:flex; justify-content:center; margin:14px 0; }
        .page-nav a { border:1px solid var(--ember); color:var(--ember); padding:8px 20px;
                      border-radius:4px; text-decoration:none; font-size:0.82rem; }
        .page-nav a:hover { background:var(--lava); color:var(--bg); }
        .post { background:rgba(20,8,2,0.72); border:1px solid rgba(255,68,0,0.2); border-radius:8px;
                padding:20px; margin-bottom:14px; cursor:pointer;
                transition:all 0.25s; backdrop-filter:blur(4px); }
        .post:hover { border-color:var(--ember); box-shadow:0 0 20px rgba(255,68,0,0.25); transform:translateY(-2px); }
        .post-active { border-color:var(--hot) !important; background:rgba(30,8,2,0.9) !important;
                       box-shadow:0 0 24px rgba(255,140,0,0.4),inset 3px 0 0 var(--hot) !important; }
        .post-header { display:flex; justify-content:space-between; margin-bottom:12px; font-size:0.88rem; }
        .post-time { color:var(--ember); font-family:monospace; font-size:0.8rem; }
        .post-text { line-height:1.65; font-size:0.95rem; }
        .post-text a { color:var(--ember); text-decoration:none; }
        .post-text a:hover { text-shadow:0 0 8px var(--lava); }
        .post-audio { width:100%; margin-top:10px; }
        .post-modal { display:none; position:fixed; inset:0; z-index:100;
                      background:rgba(13,8,2,0.96); backdrop-filter:blur(20px);
                      overflow-y:auto; padding:40px 20px; }
        .post-modal.active { display:block; }
        .modal-inner { max-width:760px; margin:0 auto; background:rgba(20,8,2,0.98);
                       border:1px solid var(--lava); border-radius:10px;
                       box-shadow:0 0 60px rgba(255,68,0,0.3); padding:40px; }
        .modal-close { float:right; background:none; border:none; color:var(--ember);
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
    // Volcano WebGL: 2000 ember particles emitted from the bottom, rising with
    // drift and fade. Each particle has a randomised lifetime, speed, and hue
    // shifting from hot yellow through orange to red as it rises and cools.
    (function() {
        var N = 2000;
        var scene, camera, renderer, clock;
        var points;
        var posArr, colArr, alpArr;

        // Per-particle state
        var px = new Float32Array(N);
        var py = new Float32Array(N);
        var pz = new Float32Array(N);
        var vx = new Float32Array(N); // horizontal drift
        var vy = new Float32Array(N); // rise speed
        var life = new Float32Array(N);   // 0..1, resets at 0
        var maxLife = new Float32Array(N);

        function resetParticle(i) {
            // Spawn along a wide base strip at the bottom
            px[i] = (Math.random() - 0.5) * 60;
            py[i] = -25 + (Math.random() - 0.5) * 4;
            pz[i] = (Math.random() - 0.5) * 20 - 5;
            vx[i] = (Math.random() - 0.5) * 0.06;
            vy[i] = 0.06 + Math.random() * 0.12;
            maxLife[i] = 0.5 + Math.random() * 0.5;
            life[i] = Math.random(); // stagger initial phases
        }

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x0d0802);
            scene.fog = new THREE.Fog(0x0d0802, 30, 80);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 120);
            camera.position.set(0, 0, 45);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: false });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();

            posArr = new Float32Array(N * 3);
            colArr = new Float32Array(N * 3);

            for (var i = 0; i < N; i++) resetParticle(i);

            var geo = new THREE.BufferGeometry();
            geo.setAttribute('position', new THREE.BufferAttribute(posArr, 3));
            geo.setAttribute('color',    new THREE.BufferAttribute(colArr, 3));

            points = new THREE.Points(geo, new THREE.PointsMaterial({
                size: 0.25, vertexColors: true,
                transparent: true, opacity: 0.9,
                blending: THREE.AdditiveBlending, depthWrite: false
            }));
            scene.add(points);

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
            var dt = clock.getDelta();

            for (var i = 0; i < N; i++) {
                life[i] += dt / maxLife[i];
                if (life[i] > 1.0) resetParticle(i);

                py[i] += vy[i];
                px[i] += vx[i];

                var idx = i * 3;
                posArr[idx]   = px[i];
                posArr[idx+1] = py[i];
                posArr[idx+2] = pz[i];

                // Colour: young embers are yellow/hot, older ones shift to orange then red
                var t = life[i];
                var fade = Math.max(0, 1 - t * 1.4);
                colArr[idx]   = fade;                            // R: always full
                colArr[idx+1] = fade * Math.max(0, 1 - t * 2); // G: fades fast
                colArr[idx+2] = 0;                              // B: never
            }

            points.geometry.attributes.position.needsUpdate = true;
            points.geometry.attributes.color.needsUpdate    = true;
            renderer.render(scene, camera);
        }

        initThree();
    })();
    </script>
    {{template "navscript" .}}
</body>
</html>`
