(function () {
    'use strict';
    var reduced = window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    var wild = false;

    function mazeScene(canvas, splash) {
        if (!canvas || typeof THREE === 'undefined') return function () {};
        var renderer = new THREE.WebGLRenderer({canvas:canvas, antialias:true, alpha:true});
        renderer.setClearColor(0x03051d, splash ? 0 : 1);
        renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2));
        var scene = new THREE.Scene();
        var camera = new THREE.PerspectiveCamera(50, 1, .1, 100);
        camera.position.set(0, 0, splash ? 18 : 24);
        var objects = [], walls = [], dots = [], chasers = [], raf = 0, start = performance.now();
        var wallMat = new THREE.MeshBasicMaterial({color:0x168cff, transparent:true, opacity:splash ? .28 : .2});
        var dotMat = new THREE.MeshBasicMaterial({color:0xffd84a});

        function box(x, y, w, h) {
            var mesh = new THREE.Mesh(new THREE.BoxGeometry(w, h, .18), wallMat);
            mesh.position.set(x, y, -2);
            mesh.userData.homeX = x; mesh.userData.homeY = y;
            mesh.userData.phase = walls.length * .71;
            scene.add(mesh); walls.push(mesh); objects.push(mesh);
        }
        var paths = [[-8,7,8,.32],[4,7,7,.32],[-6,3,.32,8],[2,4,.32,6],[7,2,.32,10],[-2,0,8,.32],[-8,-4,9,.32],[4,-5,7,.32],[0,-8,15,.32]];
        paths.forEach(function (p) { box(p[0], p[1], p[2], p[3]); });
        for (var i = 0; i < 30; i++) {
            var d = new THREE.Mesh(new THREE.SphereGeometry(.09, 8, 8), dotMat);
            d.position.set(-9 + (i % 10) * 2, 6 - Math.floor(i / 10) * 5, -.5);
            d.userData.phase = i * .17; scene.add(d); dots.push(d); objects.push(d);
        }
        var orbColors = [0xff9fcf, 0x8ff0d2, 0xffad88, 0xc7a7ff];
        orbColors.forEach(function (color, i) {
            var group = new THREE.Group();
            var orb = new THREE.Mesh(new THREE.SphereGeometry(.48, 18, 14), new THREE.MeshBasicMaterial({color:color}));
            group.add(orb); objects.push(orb);
            [-.16,.16].forEach(function (x) {
                var eye = new THREE.Mesh(new THREE.SphereGeometry(.105, 10, 8), new THREE.MeshBasicMaterial({color:0xffffff}));
                eye.position.set(x,.1,.43); group.add(eye); objects.push(eye);
            });
            group.position.set(-7 + i * 4.5, -2 + (i % 2) * 2, 0); group.userData.phase = i * Math.PI / 2;
            scene.add(group); chasers.push(group);
        });
        var gold = new THREE.Mesh(new THREE.SphereGeometry(.62, 24, 18), dotMat);
        scene.add(gold); objects.push(gold);

        function resize() {
            var w = canvas.clientWidth || 2, h = canvas.clientHeight || 2;
            renderer.setSize(w, h, false); camera.aspect = w / h; camera.updateProjectionMatrix();
        }
        function draw(now) {
            var t = (now - start) / 1000;
            var fright = wild ? .5 + .5 * Math.sin(t * 7) : 0;
            wallMat.color.setHex(wild ? (fright > .5 ? 0x7367ff : 0x26e8ff) : 0x168cff);
            wallMat.opacity = splash ? .28 : (wild ? .34 + fright * .12 : .2);
            dotMat.color.setHex(wild ? (fright > .5 ? 0xffffff : 0xff4fd8) : 0xffd84a);
            renderer.setClearColor(wild ? (fright > .5 ? 0x17073b : 0x001c38) : 0x03051d, splash ? 0 : 1);

            walls.forEach(function (wall, i) {
                var surge = wild ? 1 + .18 * Math.sin(t * 5.2 + wall.userData.phase) : 1;
                wall.position.x = wall.userData.homeX + (wild ? Math.sin(t * 3.1 + i) * .75 : 0);
                wall.position.y = wall.userData.homeY + (wild ? Math.cos(t * 3.7 + i * .8) * .55 : 0);
                wall.rotation.z = wild ? Math.sin(t * 2.8 + i) * .18 : 0;
                wall.scale.set(surge, wild ? 1 + .3 * Math.cos(t * 4.4 + i) : 1, 1);
            });

            var goldSpeed = wild ? 5.8 : 1;
            gold.position.set(Math.sin(t * goldSpeed) * (wild ? 8.5 : 7), Math.sin(t * goldSpeed * 2) * (wild ? 3.8 : 3), .2);
            gold.scale.set(wild ? 1.1 + .35 * Math.sin(t * 12) : 1, wild ? .65 : 1, 1);
            gold.rotation.z = wild ? t * 8 : 0;
            dots.forEach(function (d, i) {
                var pulse = .65 + .35 * Math.sin(t * (wild ? 11 : 3) + d.userData.phase);
                d.scale.set(wild ? 1.2 + pulse * 3.8 : pulse, wild ? .45 + pulse * .45 : pulse, 1);
                d.rotation.z = wild ? Math.sin(t * 4 + i) * 1.2 : 0;
                d.material.color.setHex(wild ? (i % 3 === 0 ? 0xffffff : (i % 3 === 1 ? 0xff4fd8 : 0x65f7ff)) : 0xffd84a);
            });
            chasers.forEach(function (c, i) {
                var a = t * (wild ? 3.8 + i * .32 : .35) + c.userData.phase;
                var radiusX = wild ? 7.8 - i * .35 : 3.1 + i * .45;
                var radiusY = wild ? 4.6 - i * .25 : 1.7 + i * .22;
                c.position.set(Math.sin(a) * radiusX, Math.sin(a * 2 + i) * radiusY, wild ? Math.cos(a) * 1.4 : 0);
                c.rotation.z = wild ? -a + Math.cos(a * 2) * .4 : 0;
                c.scale.setScalar(wild ? .68 + .2 * Math.sin(t * 10 + i) : 1);
            });
            if (!splash) {
                camera.position.x = Math.sin(t * (wild ? 1.9 : .12)) * (wild ? 3.8 : 2);
                camera.position.y = wild ? Math.cos(t * 1.5) * 2.2 : 0;
                camera.position.z = wild ? 19 + Math.sin(t * 2.3) * 4 : 24;
                camera.lookAt(wild ? Math.sin(t * 2.1) * 1.6 : 0, wild ? Math.cos(t * 1.8) : 0, 0);
                camera.rotation.z = wild ? Math.sin(t * 1.7) * .12 : 0;
            }
            renderer.render(scene, camera);
            if (!reduced) raf = requestAnimationFrame(draw);
        }
        resize(); window.addEventListener('resize', resize); draw(start);
        return function (state) {
            if (typeof state === 'boolean') { wild = state; if (reduced) draw(performance.now()); return; }
            window.removeEventListener('resize', resize); if (raf) cancelAnimationFrame(raf);
            objects.forEach(function (o) { if (o.geometry) o.geometry.dispose(); });
            [wallMat, dotMat].forEach(function (m) { m.dispose(); });
            chasers.forEach(function (c) { c.children.forEach(function (o) { if (o.material) o.material.dispose(); }); });
            renderer.dispose();
        };
    }

    if (!document.documentElement.classList.contains('sno-splash-skip')) {
        window._snonuxSplashWebGLCleanup = mazeScene(document.getElementById('splash-gl-canvas'), true);
    }
    var setWild = mazeScene(document.getElementById('three-canvas'), false);

    function dotTrail(color) {
        if (reduced) return;
        for (var i = 0; i < 9; i++) {
            var dot = document.createElement('i'); dot.className = 'maze-dot';
            dot.style.left = (15 + i * 8) + 'vw'; dot.style.top = (48 + Math.sin(i) * 8) + 'vh';
            if (color) dot.style.background = dot.style.boxShadow = color;
            dot.style.animationDelay = (i * .035) + 's'; document.body.appendChild(dot);
            setTimeout((function (node) { return function () { node.remove(); }; })(dot), 850);
        }
    }
    window.snonuxOpenEffect = function () {
        var modal = document.getElementById('post-modal');
        if (modal) { modal.classList.add('sno-modal-zoom'); setTimeout(function () { modal.classList.remove('sno-modal-zoom'); }, 400); }
        dotTrail('#ffd84a');
    };
    window.snonuxCloseEffect = function () { dotTrail('#57b8ff'); };
    window.snonuxNavEffect = function () { dotTrail(wild ? '#ffffff' : '#ffd84a'); };
    window.snonuxPageEffect = function () { dotTrail('#ff9fcf'); };
    window.snonuxScrollEffect = function (dir) {
        var active = document.querySelector('.post-active'); if (!active) return;
        active.style.transform = 'translateX(' + (dir === 'down' ? '13px' : '-3px') + ')';
        setTimeout(function () { active.style.transform = ''; }, 170); dotTrail();
    };
    window.snonuxWildToggle = function () {
        wild = !wild; setWild(wild); document.body.classList.toggle('maze-frightened', wild);
        var badge = document.getElementById('sno-wild-badge'); if (badge) badge.classList.toggle('sno-wild-on', wild);
        document.documentElement.style.setProperty('--sno-wild-accent', wild ? '#777dff' : '#ffd84a');
    };
})();
