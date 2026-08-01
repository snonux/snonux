(function () {
    var reduced = window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    var colors = [0x36c5f0, 0x4361ee, 0xff922b, 0xffd43b, 0x51cf66, 0xb197fc, 0xff5c5c];
    var shapes = [
        [[-1,0],[0,0],[1,0],[2,0]], [[-1,0],[0,0],[1,0],[1,1]], [[-1,1],[-1,0],[0,0],[1,0]],
        [[0,0],[1,0],[0,1],[1,1]], [[-1,0],[0,0],[0,1],[1,1]], [[-1,0],[0,0],[1,0],[0,1]], [[-1,1],[0,1],[0,0],[1,0]]
    ];

    function tetromino(shape, color, size) {
        var group = new THREE.Group();
        var geo = new THREE.BoxGeometry(size * .92, size * .92, size * .32);
        var mat = new THREE.MeshPhongMaterial({color:color, shininess:55, flatShading:true});
        shape.forEach(function (cell) {
            var cube = new THREE.Mesh(geo, mat);
            cube.position.set(cell[0] * size, cell[1] * size, 0);
            group.add(cube);
        });
        return group;
    }

    function startScene(canvas, splash) {
        if (!canvas || typeof THREE === 'undefined') return function () {};
        var renderer = new THREE.WebGLRenderer({canvas:canvas, antialias:true, alpha:true});
        renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2));
        var scene = new THREE.Scene();
        var camera = new THREE.PerspectiveCamera(52, 1, .1, 100);
        camera.position.z = splash ? 12 : 25;
        scene.add(new THREE.AmbientLight(0xffffff, .85));
        var light = new THREE.DirectionalLight(0xffffff, 1.5); light.position.set(3, 8, 12); scene.add(light);
        var pieces = [], wild = false, raf = 0, last = performance.now();
        var count = splash ? 7 : 24;
        for (var i = 0; i < count; i++) {
            var p = tetromino(shapes[i % shapes.length], colors[i % colors.length], splash ? .42 : .72);
            p.position.set((Math.random() - .5) * (splash ? 13 : 34), (Math.random() - .5) * 28, (Math.random() - .5) * 10);
            p.rotation.z = Math.floor(Math.random() * 4) * Math.PI / 2;
            p.userData.speed = .8 + Math.random() * 1.1;
            pieces.push(p); scene.add(p);
        }
        function resize() {
            var w = canvas.clientWidth || 2, h = canvas.clientHeight || 2;
            renderer.setSize(w, h, false); camera.aspect = w / h; camera.updateProjectionMatrix();
        }
        function draw(now) {
            var dt = Math.min((now - last) / 1000, .05); last = now;
            pieces.forEach(function (p, i) {
                p.position.y -= p.userData.speed * dt * (wild ? 6 : 1);
                p.rotation.z += dt * (wild ? .75 : .04) * (i % 2 ? 1 : -1);
                if (p.position.y < -16) { p.position.y = 16; p.position.x = (Math.random() - .5) * (splash ? 13 : 34); }
            });
            if (!splash) { camera.position.x = Math.sin(now * .00012) * (wild ? 8 : 2); camera.lookAt(0, 0, 0); }
            renderer.render(scene, camera);
            if (!reduced) raf = requestAnimationFrame(draw);
        }
        resize(); window.addEventListener('resize', resize); draw(last);
        return function (state) {
            if (typeof state === 'boolean') { wild = state; return; }
            window.removeEventListener('resize', resize); if (raf) cancelAnimationFrame(raf); renderer.dispose();
        };
    }

    var cleanSplash = function () {};
    if (!document.documentElement.classList.contains('sno-splash-skip')) {
        cleanSplash = startScene(document.getElementById('splash-gl-canvas'), true);
        window._snonuxSplashWebGLCleanup = cleanSplash;
    }
    var setWild = startScene(document.getElementById('three-canvas'), false);
    var wild = false;

    function lineClear(color) {
        var line = document.createElement('div');
        line.style.cssText = 'position:fixed;left:0;right:0;top:50%;height:28px;z-index:9000;pointer-events:none;background:' + color + ';box-shadow:0 0 26px ' + color + ';transform:scaleX(0);transition:transform .16s steps(6,end),opacity .28s';
        document.body.appendChild(line);
        requestAnimationFrame(function () { line.style.transform = 'scaleX(1)'; });
        setTimeout(function () { line.style.opacity = '0'; }, 180);
        setTimeout(function () { line.remove(); }, 480);
    }
    window.snonuxOpenEffect = function () {
        var modal = document.getElementById('post-modal');
        if (modal) { modal.classList.add('sno-modal-zoom'); setTimeout(function () { modal.classList.remove('sno-modal-zoom'); }, 400); }
        lineClear('#ffd43b');
    };
    window.snonuxCloseEffect = function () { lineClear('#4361ee'); };
    window.snonuxNavEffect = function () { lineClear('#36c5f0'); };
    window.snonuxPageEffect = function () { lineClear('#b197fc'); };
    window.snonuxScrollEffect = function (dir) {
        var active = document.querySelector('.post-active');
        if (!active) return;
        active.style.transform = 'translateX(' + (dir === 'down' ? '14px' : '-2px') + ')';
        setTimeout(function () { active.style.transform = ''; }, 170);
    };
    window.snonuxWildToggle = function () {
        wild = !wild; setWild(wild);
        var badge = document.getElementById('sno-wild-badge'); if (badge) badge.classList.toggle('sno-wild-on', wild);
        document.documentElement.style.setProperty('--sno-wild-accent', wild ? '#ffd43b' : '#36c5f0');
    };
})();
