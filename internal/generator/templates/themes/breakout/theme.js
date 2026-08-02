(function () {
    'use strict';
    var reduced = window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    var palette = [0xc75c5c, 0xd98255, 0xd5aa58, 0x77a66e, 0x559ca0, 0x657eae, 0x8c6f9f];

    function buildScene(canvas, splash) {
        if (!canvas || typeof THREE === 'undefined') return { wild: function () {}, hit: function () {}, destroy: function () {} };
        var renderer = new THREE.WebGLRenderer({canvas: canvas, alpha: true, antialias: true});
        renderer.setClearColor(0, 0); renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2));
        var scene = new THREE.Scene();
        var camera = new THREE.OrthographicCamera(-10, 10, 7, -7, .1, 50); camera.position.z = 12;
        scene.add(new THREE.AmbientLight(0xffffff, 1.4));
        var bricks = [], balls = [], rows = splash ? 3 : 6, cols = splash ? 7 : 11;
        var brickGeo = new THREE.BoxGeometry(1.45, .42, .3);
        for (var r = 0; r < rows; r++) for (var c = 0; c < cols; c++) {
            var brick = new THREE.Mesh(brickGeo, new THREE.MeshLambertMaterial({color: palette[r % 7], transparent: true}));
            brick.position.set((c - (cols - 1) / 2) * 1.65, 5.3 - r * .62, 0);
            brick.userData.row = r; brick.userData.col = c;
            brick.userData.homeX = brick.position.x; brick.userData.homeY = brick.position.y;
            bricks.push(brick); scene.add(brick);
        }
        var paddle = new THREE.Mesh(new THREE.BoxGeometry(3.2, .28, .36), new THREE.MeshLambertMaterial({color: 0xf2eddf}));
        paddle.position.y = -5.5; scene.add(paddle);
        var ballGeo = new THREE.SphereGeometry(.22, 16, 12);
        var ballMaterial = new THREE.MeshLambertMaterial({color: 0xf2eddf});
        function addBall(seed) {
            var ball = new THREE.Mesh(ballGeo, ballMaterial);
            ball.position.set((seed - 1) * .6, -3 + seed * .25, .35);
            ball.userData.vx = (seed % 2 ? 1 : -1) * (2.8 + (seed % 6) * .35);
            ball.userData.vy = 3.4 + (seed % 7) * .25;
            balls.push(ball); scene.add(ball);
        }
        addBall(0);
        var maxBalls = splash ? 10 : 18, activeBalls = 1;
        var isWild = false, raf = 0, last = performance.now(), nextRow = 0, dead = false;
        function resize() {
            var w = canvas.clientWidth || 2, h = canvas.clientHeight || 2, aspect = w / h;
            camera.left = -7 * aspect; camera.right = 7 * aspect; camera.top = 7; camera.bottom = -7; camera.updateProjectionMatrix(); renderer.setSize(w, h, false);
        }
        function frame(now) {
            if (dead) return;
            var dt = Math.min((now - last) / 1000, .04) * (isWild ? 2.8 : 1); last = now;
            balls.slice(0, activeBalls).forEach(function (b) {
                b.position.x += b.userData.vx * dt; b.position.y += b.userData.vy * dt;
                var edge = 6.8 * camera.right / 10;
                if (b.position.x > edge || b.position.x < -edge) { b.position.x = Math.max(-edge, Math.min(edge, b.position.x)); b.userData.vx *= -1; }
                if (b.position.y > 6.4) b.userData.vy = -Math.abs(b.userData.vy);
                if (b.position.y < -5.15 && b.position.y > -5.8 && Math.abs(b.position.x - paddle.position.x) < 1.9) b.userData.vy = Math.abs(b.userData.vy);
                if (b.position.y < -6.7) { b.position.set(0, -3, .35); b.userData.vy = Math.abs(b.userData.vy); }
            });
            paddle.position.x = isWild
                ? (Math.sin(now * .0063) * 3.4 + Math.sin(now * .014) * 1.25)
                : Math.sin(now * .00055) * 3.6;
            if (isWild) {
                bricks.forEach(function (b) {
                    var wave = now * .0048 + b.userData.col * .52 + b.userData.row * .76;
                    b.position.x = b.userData.homeX + Math.sin(now * .0017 + b.userData.row) * .55;
                    b.position.y = b.userData.homeY + Math.sin(wave) * .25;
                    b.scale.y = .65 + .55 * (1 + Math.sin(wave * 1.35)) / 2;
                    b.material.opacity = .22 + .78 * (1 + Math.sin(wave - b.userData.row)) / 2;
                });
                camera.position.x = Math.sin(now * .0021) * .7;
                camera.position.y = Math.cos(now * .0016) * .38;
                camera.rotation.z = Math.sin(now * .00135) * .022;
            } else {
                bricks.forEach(function (b) {
                    b.position.x = b.userData.homeX; b.position.y = b.userData.homeY; b.scale.y = 1;
                });
                camera.position.x = 0; camera.position.y = 0; camera.rotation.z = 0;
            }
            renderer.render(scene, camera); if (!reduced) raf = requestAnimationFrame(frame);
        }
        function hit() {
            var row = nextRow++ % rows;
            bricks.forEach(function (b) { if (b.userData.row === row) b.material.opacity = .08; });
            window.setTimeout(function () { if (!dead) bricks.forEach(function (b) { if (b.userData.row === row) b.material.opacity = 1; }); }, 480);
        }
        function wild(on) {
            isWild = on && !reduced;
            if (isWild) {
                while (balls.length < maxBalls) addBall(balls.length);
                balls.forEach(function (ball) { if (ball.parent !== scene) scene.add(ball); });
                activeBalls = maxBalls;
            } else {
                activeBalls = 1;
                balls.forEach(function (ball, index) { if (index && ball.parent === scene) scene.remove(ball); });
                bricks.forEach(function (b) { b.material.opacity = 1; });
            }
        }
        function destroy() {
            dead = true; window.removeEventListener('resize', resize); if (raf) cancelAnimationFrame(raf);
            bricks.forEach(function (b) { b.material.dispose(); });
            brickGeo.dispose(); paddle.geometry.dispose(); paddle.material.dispose(); ballGeo.dispose(); ballMaterial.dispose();
            renderer.dispose(); if (window._snonuxSplashWebGLCleanup === destroy) window._snonuxSplashWebGLCleanup = null;
        }
        resize(); window.addEventListener('resize', resize); frame(last);
        return {wild: wild, hit: hit, destroy: destroy};
    }

    if (!document.documentElement.classList.contains('sno-splash-skip')) {
        var splashScene = buildScene(document.getElementById('splash-gl-canvas'), true);
        window._snonuxSplashWebGLCleanup = splashScene.destroy;
    }
    var game = buildScene(document.getElementById('three-canvas'), false), wild = false;
    function flash(color, vertical) {
        var el = document.createElement('div');
        el.style.cssText = 'position:fixed;z-index:9000;pointer-events:none;background:' + color + ';opacity:.72;transition:transform .2s ease,opacity .28s ease;' + (vertical ? 'top:0;bottom:0;width:5px;left:50%;transform:scaleY(0)' : 'left:0;right:0;height:7px;top:50%;transform:scaleX(0)');
        document.body.appendChild(el); requestAnimationFrame(function () { el.style.transform = vertical ? 'scaleY(1)' : 'scaleX(1)'; });
        setTimeout(function () { el.style.opacity = '0'; }, 180); setTimeout(function () { el.remove(); }, 500);
    }
    window.snonuxOpenEffect = function () { game.hit(); flash('#d5aa58', false); };
    window.snonuxCloseEffect = function () { flash('#657eae', true); };
    window.snonuxNavEffect = function () { game.hit(); flash('#559ca0', false); };
    window.snonuxPageEffect = function () { game.hit(); flash('#8c6f9f', false); };
    window.snonuxScrollEffect = function (direction) {
        var post = document.querySelector('.post-active'); if (!post) return;
        post.style.transform = 'translateX(' + (direction === 'down' ? '9px' : '-3px') + ')'; setTimeout(function () { post.style.transform = ''; }, 150);
    };
    window.snonuxWildToggle = function () {
        wild = !wild; game.wild(wild);
        var badge = document.getElementById('sno-wild-badge'); if (badge) badge.classList.toggle('sno-wild-on', wild);
        document.documentElement.style.setProperty('--sno-wild-accent', wild ? '#c75c5c' : '#d5aa58');
    };
}());
