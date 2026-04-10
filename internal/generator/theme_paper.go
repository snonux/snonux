package generator

// volcanoTemplate is a dark volcanic theme — rising ember particles, a glowing
// lava plane at the base, molten rock boulders, smoke plumes, and a deep
// underground furnace glow on the horizon.
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
                      overflow-y:auto; padding:40px 20px; }
        .post-modal.active { display:block; }
        .modal-inner { max-width:760px; margin:0 auto; background:rgba(20,8,2,0.92);
                       border:1px solid var(--lava); border-radius:10px; backdrop-filter:blur(16px);
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
    // Volcano WebGL: glowing lava floor, molten rock boulders, smoke plumes,
    // underground furnace glow sphere, and 3000 rising ember particles.
    (function() {
        var N_EMBER = 3000;
        var N_SMOKE = 800;
        var scene, camera, renderer, clock;
        var emberPoints, smokePoints;
        var ePosArr, eColArr, sPosArr;
        var ePX, ePY, ePZ, eVX, eVY, eLife, eMaxLife;
        var sPX, sPY, sPZ, sSVY, sSLife, sSMaxLife;
        var lavaGeo, lavaFloor;

        function resetEmber(i) {
            ePX[i] = (Math.random() - 0.5) * 70;
            ePY[i] = -22 + (Math.random() - 0.5) * 4;
            ePZ[i] = (Math.random() - 0.5) * 30 - 5;
            eVX[i] = (Math.random() - 0.5) * 0.08;
            eVY[i] = 0.07 + Math.random() * 0.14;
            eMaxLife[i] = 0.4 + Math.random() * 0.6;
            eLife[i] = Math.random();
        }

        function resetSmoke(i) {
            sPX[i] = (Math.random() - 0.5) * 40;
            sPY[i] = -18 + Math.random() * 5;
            sPZ[i] = (Math.random() - 0.5) * 20 - 5;
            sSVY[i] = 0.015 + Math.random() * 0.025;
            sSMaxLife[i] = 1.5 + Math.random() * 2.0;
            sSLife[i] = Math.random();
        }

        function buildLavaFloor() {
            lavaGeo = new THREE.PlaneGeometry(200, 200, 60, 60);
            lavaFloor = new THREE.Mesh(lavaGeo, new THREE.MeshPhongMaterial({
                color: 0x8b1000, emissive: 0xff2200, emissiveIntensity: 0.6,
                shininess: 120
            }));
            lavaFloor.rotation.x = -Math.PI / 2;
            lavaFloor.position.y = -22;
            scene.add(lavaFloor);
        }

        function buildBoulders() {
            // Molten rock boulders with glowing emissive cores
            var boulderData = [
                [-18,-16,-15, 5], [20,-15,-20, 7], [-8,-14,-30, 4],
                [30,-16,-12, 6], [-28,-15,-25, 5]
            ];
            boulderData.forEach(function(b) {
                var mesh = new THREE.Mesh(
                    new THREE.IcosahedronGeometry(b[3], 1),
                    new THREE.MeshPhongMaterial({ color: 0x1a0500, emissive: 0xff4400, emissiveIntensity: 0.7, shininess: 20 })
                );
                mesh.position.set(b[0], b[1], b[2]);
                mesh.rotation.set(Math.random(), Math.random(), Math.random());
                scene.add(mesh);
            });
        }

        function buildFurnaceGlow() {
            // Underground furnace: massive low-opacity emissive sphere below the lava
            var glow = new THREE.Mesh(
                new THREE.SphereGeometry(45, 16, 16),
                new THREE.MeshBasicMaterial({ color: 0xff3300, transparent: true, opacity: 0.22, blending: THREE.AdditiveBlending, depthWrite: false })
            );
            glow.position.set(0, -60, -30);
            scene.add(glow);
        }

        function buildParticles() {
            ePX = new Float32Array(N_EMBER); ePY = new Float32Array(N_EMBER);
            ePZ = new Float32Array(N_EMBER); eVX = new Float32Array(N_EMBER);
            eVY = new Float32Array(N_EMBER); eLife = new Float32Array(N_EMBER);
            eMaxLife = new Float32Array(N_EMBER);
            ePosArr = new Float32Array(N_EMBER * 3);
            eColArr = new Float32Array(N_EMBER * 3);
            for (var i = 0; i < N_EMBER; i++) resetEmber(i);
            var eGeo = new THREE.BufferGeometry();
            eGeo.setAttribute('position', new THREE.BufferAttribute(ePosArr, 3));
            eGeo.setAttribute('color',    new THREE.BufferAttribute(eColArr, 3));
            emberPoints = new THREE.Points(eGeo, new THREE.PointsMaterial({
                size: 0.3, vertexColors: true,
                transparent: true, opacity: 0.95, blending: THREE.AdditiveBlending, depthWrite: false
            }));
            scene.add(emberPoints);

            sPX = new Float32Array(N_SMOKE); sPY = new Float32Array(N_SMOKE);
            sPZ = new Float32Array(N_SMOKE); sSVY = new Float32Array(N_SMOKE);
            sSLife = new Float32Array(N_SMOKE); sSMaxLife = new Float32Array(N_SMOKE);
            sPosArr = new Float32Array(N_SMOKE * 3);
            for (var j = 0; j < N_SMOKE; j++) resetSmoke(j);
            var sGeo = new THREE.BufferGeometry();
            sGeo.setAttribute('position', new THREE.BufferAttribute(sPosArr, 3));
            smokePoints = new THREE.Points(sGeo, new THREE.PointsMaterial({
                color: 0x444444, size: 1.8, transparent: true, opacity: 0.15, depthWrite: false
            }));
            scene.add(smokePoints);
        }

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x0d0802);
            scene.fog = new THREE.Fog(0x0d0802, 35, 100);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 150);
            camera.position.set(0, 8, 50);
            camera.lookAt(0, -5, 0);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: false });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();

            scene.add(new THREE.AmbientLight(0x220800, 1.0));
            var lavaLight = new THREE.PointLight(0xff4400, 4, 80);
            lavaLight.position.set(0, -15, 0);
            scene.add(lavaLight);

            buildLavaFloor();
            buildBoulders();
            buildFurnaceGlow();
            buildParticles();

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
            var t  = clock.getElapsedTime();

            // Pulse the lava floor emissive intensity
            lavaFloor.material.emissiveIntensity = 0.4 + 0.25 * Math.sin(t * 1.8);

            // Animate lava floor vertices
            var lp = lavaGeo.attributes.position;
            for (var i = 0; i < lp.count; i++) {
                var lx = lp.getX(i), lz = lp.getZ(i);
                lp.setY(i, Math.sin(lx * 0.08 + t * 0.7) * 0.8 + Math.cos(lz * 0.1 + t * 0.5) * 0.6);
            }
            lp.needsUpdate = true;

            // Embers
            for (var ei = 0; ei < N_EMBER; ei++) {
                eLife[ei] += dt / eMaxLife[ei];
                if (eLife[ei] > 1.0) resetEmber(ei);
                ePY[ei] += eVY[ei];
                ePX[ei] += eVX[ei];
                var idx = ei * 3, te = eLife[ei];
                ePosArr[idx] = ePX[ei]; ePosArr[idx+1] = ePY[ei]; ePosArr[idx+2] = ePZ[ei];
                var fade = Math.max(0, 1 - te * 1.3);
                eColArr[idx] = fade; eColArr[idx+1] = fade * Math.max(0, 1 - te * 2.2); eColArr[idx+2] = 0;
            }
            emberPoints.geometry.attributes.position.needsUpdate = true;
            emberPoints.geometry.attributes.color.needsUpdate    = true;

            // Smoke
            for (var si = 0; si < N_SMOKE; si++) {
                sSLife[si] += dt / sSMaxLife[si];
                if (sSLife[si] > 1.0) resetSmoke(si);
                sPY[si] += sSVY[si];
                sPX[si] += (Math.random() - 0.5) * 0.04;
                var si3 = si * 3;
                sPosArr[si3] = sPX[si]; sPosArr[si3+1] = sPY[si]; sPosArr[si3+2] = sPZ[si];
            }
            smokePoints.geometry.attributes.position.needsUpdate = true;

            renderer.render(scene, camera);
        }

        initThree();
    })();
    </script>
    {{template "navscript" .}}
</body>
</html>`
