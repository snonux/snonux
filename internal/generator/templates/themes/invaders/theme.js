(function () {
    var reduced = window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    var wild = false;
    var patterns = [
        [[1, 0], [3, 0], [0, 1], [1, 1], [2, 1], [3, 1], [4, 1], [0, 2], [2, 2], [4, 2], [1, 3], [3, 3]],
        [[0, 0], [4, 0], [1, 1], [2, 1], [3, 1], [0, 2], [1, 2], [2, 2], [3, 2], [4, 2], [1, 3], [3, 3]]
    ];

    function unit(pattern, size, material) {
        var group = new THREE.Group();
        var geometry = new THREE.BoxGeometry(size, size, size * .35);
        pattern.forEach(function (point) {
            var mesh = new THREE.Mesh(geometry, material);
            mesh.position.set((point[0] - 2) * size, (1.5 - point[1]) * size, 0);
            group.add(mesh);
        });
        return group;
    }

    function scene(canvas, splash) {
        if (!canvas || typeof THREE === 'undefined') return function () {};
        var renderer = new THREE.WebGLRenderer({canvas: canvas, antialias: false, alpha: true});
        var world = new THREE.Scene();
        var camera = new THREE.PerspectiveCamera(48, 1, .1, 100);
        var formation = new THREE.Group();
        var forms = [];
        var projectiles = [];
        var materials = [];
        var raf = 0;
        var dead = false;
        var last = performance.now();
        var elapsed = 0;
        var marchClock = 0;
        var marchX = 0;
        var marchY = 0;
        var direction = 1;
        var localWild = false;
        var rows = splash ? 3 : 5;
        var cols = splash ? 5 : 8;
        var spacingX = splash ? 1.45 : 2.4;
        var spacingY = splash ? 1.15 : 1.8;

        camera.position.z = splash ? 13 : 26;
        world.add(formation);
        world.add(new THREE.AmbientLight(0x79d9bd, .48));
        var pressureLight = new THREE.PointLight(0xff486d, 0, 45);
        pressureLight.position.set(0, 0, 8);
        world.add(pressureLight);

        for (var y = 0; y < rows; y++) {
            var material = new THREE.MeshStandardMaterial({
                color: 0x9fffd8,
                emissive: 0x123c31,
                emissiveIntensity: .18,
                roughness: .7
            });
            materials.push(material);
            for (var x = 0; x < cols; x++) {
                var form = unit(patterns[y % 2], splash ? .11 : .18, material);
                form.userData.homeX = (x - (cols - 1) / 2) * spacingX;
                form.userData.homeY = (rows / 2 - y) * spacingY;
                form.userData.row = y;
                form.userData.col = x;
                form.position.set(form.userData.homeX, form.userData.homeY, 0);
                forms.push(form);
                formation.add(form);
            }
        }

        /* A fixed pool makes repeated wild-mode toggles allocation-free. */
        if (!splash) {
            var shotGeometry = new THREE.BoxGeometry(.075, .7, .075);
            var shotMaterial = new THREE.MeshBasicMaterial({color: 0xffd477});
            materials.push(shotMaterial);
            for (var p = 0; p < 18; p++) {
                var shot = new THREE.Mesh(shotGeometry, shotMaterial);
                shot.userData.phase = p / 18;
                shot.visible = false;
                projectiles.push(shot);
                world.add(shot);
            }
        }

        function resize() {
            var width = canvas.clientWidth || 2;
            var height = canvas.clientHeight || 2;
            renderer.setSize(width, height, false);
            camera.aspect = width / height;
            camera.updateProjectionMatrix();
        }

        function draw(now) {
            if (dead) return;
            var dt = Math.min((now - last) / 1000, .05);
            last = now;
            elapsed += dt;

            if (localWild && !reduced) {
                marchClock += dt;
                if (marchClock >= .11) {
                    marchClock %= .11;
                    marchX += direction * .38;
                    if (Math.abs(marchX) > 2.2) {
                        direction *= -1;
                        marchX += direction * .38;
                        marchY -= .32;
                        if (marchY < -2.2) marchY = 1.25;
                    }
                }
                formation.position.x = marchX;
                formation.position.y = marchY + Math.sin(elapsed * 3.8) * .16;
                var scale = 1 + Math.sin(elapsed * 2.7) * .13;
                formation.scale.set(scale, scale, scale);
                formation.rotation.z = Math.sin(elapsed * 1.9) * .035;
                camera.position.x = Math.sin(elapsed * 1.35) * 1.5;
                camera.position.y = Math.cos(elapsed * 1.7) * .7;
                camera.position.z = 24 + Math.sin(elapsed * 1.1) * 2.2;
                camera.lookAt(0, 0, 0);
                pressureLight.intensity = 1.2 + Math.sin(elapsed * 8) * .65;

                forms.forEach(function (form) {
                    var row = form.userData.row;
                    var col = form.userData.col;
                    form.position.x = form.userData.homeX + Math.sin(elapsed * 5 + row * .8) * .28;
                    form.position.y = form.userData.homeY + Math.sin(elapsed * 6 + col * .58) * .23;
                    form.position.z = Math.sin(elapsed * 4.2 + col * .72 - row) * 2.25;
                    form.rotation.y = Math.sin(elapsed * 3.4 + col * .5) * .42;
                });
                materials.slice(0, rows).forEach(function (material, index) {
                    var pressure = (Math.sin(elapsed * 7 + index) + 1) / 2;
                    material.color.setHSL(.43 - pressure * .42, .9, .63);
                    material.emissive.setHSL(.43 - pressure * .4, 1, .18);
                    material.emissiveIntensity = .5 + pressure * 1.4;
                });
                projectiles.forEach(function (shot) {
                    var travel = (elapsed * 1.75 + shot.userData.phase) % 1;
                    var source = forms[Math.floor(shot.userData.phase * forms.length)];
                    shot.visible = true;
                    shot.position.x = source.position.x + formation.position.x;
                    shot.position.y = 8 - travel * 18;
                    shot.position.z = source.position.z + Math.sin(elapsed * 9 + shot.userData.phase * 20) * .7;
                    shot.scale.y = 1 + travel * 2.2;
                });
            } else {
                formation.position.x += (0 - formation.position.x) * Math.min(1, dt * 3);
                formation.position.y += (0 - formation.position.y) * Math.min(1, dt * 3);
                formation.scale.set(1, 1, 1);
                formation.rotation.z = 0;
                camera.position.x = 0;
                camera.position.y = 0;
                camera.position.z = splash ? 13 : 26;
                camera.lookAt(0, 0, 0);
                pressureLight.intensity = 0;
                forms.forEach(function (form) {
                    form.position.x = form.userData.homeX;
                    form.position.y = form.userData.homeY + (reduced ? 0 : Math.sin(elapsed * .7 + form.userData.col) * .035);
                    form.position.z = 0;
                    form.rotation.y = 0;
                });
                materials.slice(0, rows).forEach(function (material) {
                    material.color.setHex(0x9fffd8);
                    material.emissive.setHex(0x123c31);
                    material.emissiveIntensity = .18;
                });
                projectiles.forEach(function (shot) { shot.visible = false; });
            }
            renderer.render(world, camera);
            if (!reduced) raf = requestAnimationFrame(draw);
        }

        function cleanup() {
            dead = true;
            window.removeEventListener('resize', resize);
            if (raf) cancelAnimationFrame(raf);
            world.traverse(function (object) {
                if (object.geometry) object.geometry.dispose();
            });
            materials.forEach(function (material) { material.dispose(); });
            renderer.dispose();
            if (window._snonuxSplashWebGLCleanup === cleanup) window._snonuxSplashWebGLCleanup = null;
        }

        resize();
        window.addEventListener('resize', resize);
        draw(last);
        return function (state) {
            if (typeof state === 'boolean') {
                localWild = state && !splash;
                if (reduced) draw(performance.now());
            } else {
                cleanup();
            }
        };
    }

    if (!document.documentElement.classList.contains('sno-splash-skip')) {
        window._snonuxSplashWebGLCleanup = scene(document.getElementById('splash-gl-canvas'), true);
    }
    var control = scene(document.getElementById('three-canvas'), false);

    function pulse() {
        var pulseElement = document.createElement('div');
        pulseElement.style.cssText = 'position:fixed;z-index:9000;pointer-events:none;left:50%;bottom:0;width:3px;height:34px;background:#edfff8;box-shadow:0 0 10px #9fffd8;transition:transform .34s linear,opacity .34s';
        document.body.appendChild(pulseElement);
        requestAnimationFrame(function () {
            pulseElement.style.transform = 'translateY(-100vh)';
            pulseElement.style.opacity = '0';
        });
        setTimeout(function () { pulseElement.remove(); }, 390);
    }

    window.snonuxNavEffect = pulse;
    window.snonuxPageEffect = pulse;
    window.snonuxScrollEffect = pulse;
    window.snonuxOpenEffect = function () {
        pulse();
        var modal = document.getElementById('post-modal');
        if (modal) {
            modal.classList.add('sno-modal-zoom');
            setTimeout(function () { modal.classList.remove('sno-modal-zoom'); }, 400);
        }
    };
    window.snonuxCloseEffect = pulse;
    window.snonuxWildToggle = function () {
        wild = !wild;
        control(wild);
        var badge = document.getElementById('sno-wild-badge');
        if (badge) badge.classList.toggle('sno-wild-on', wild);
    };
}());
