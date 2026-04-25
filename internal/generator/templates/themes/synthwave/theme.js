
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,g=new THREE.Group(),t0=performance.now();
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(58,1,0.1,120);ca.position.set(0,1.2,7);
        var sun=new THREE.Mesh(new THREE.SphereGeometry(1.35,28,28),new THREE.MeshBasicMaterial({color:0xff6b2b,transparent:true,opacity:0.95}));
        sun.position.y=2.1;g.add(sun);
        var gr=new THREE.Mesh(new THREE.PlaneGeometry(28,28,20,20),new THREE.MeshBasicMaterial({color:0xbf3fff,wireframe:true,transparent:true,opacity:0.4}));
        gr.rotation.x=-Math.PI/2.15;gr.position.y=-2.4;g.add(gr);
        sc.add(g);sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);var t=(now-t0)*0.001;g.rotation.y=Math.sin(t*0.35)*0.08;sun.position.y=2.1+Math.sin(t*1.2)*0.08;sun.scale.setScalar(1+Math.sin(t*2)*0.04);ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();


    // Synthwave WebGL: glowing sunset sphere with horizontal scan-line rings,
    // a receding grid floor, and pink star particles. Replaces CSS sky/grid.
    (function() {
        var _wild = false, _snoTOffset = 0, _snoLastT = 0;
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
            var realT = clock.getElapsedTime();
            _snoTOffset += (realT - _snoLastT) * (_wild ? 9 : 0);
            _snoLastT = realT;
            var t = realT + _snoTOffset;
            // Sun pulses subtly; wild mode makes it throb dramatically
            var pulse = _wild ? 1 + 0.12 * Math.sin(t * 8) : 1 + 0.015 * Math.sin(t * 1.5);
            sun.scale.setScalar(pulse);
            sunRings.forEach(function(r) { r.scale.setScalar(pulse); });
            // Camera: warp-drive rush in wild, gentle sway otherwise
            if (_wild) {
                camera.position.x = Math.sin(realT * 0.35) * 16;
                camera.position.y = 10 + Math.sin(realT * 0.28) * 9;
                camera.position.z = 45 + Math.sin(realT * 0.23) * 18;
                camera.fov = 60 + Math.sin(realT * 0.46) * 16;
                camera.updateProjectionMatrix();
            } else {
                camera.position.x = Math.sin(t * 0.08) * 4;
                camera.position.y = 10;
                camera.position.z = 45;
                if (camera.fov !== 60) { camera.fov = 60; camera.updateProjectionMatrix(); }
            }
            renderer.render(scene, camera);
        }

        initThree();

        // Synthwave nav/wild effects — grid zoom on navigate, turbo drive on wild
        window.snonuxOpenEffect = function() {
            // Fly up from the grid floor + neon scan line
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-fly'); setTimeout(function() { modal.classList.remove('sno-modal-fly'); }, 390); }
            // Pink scan line rises from bottom
            var line = document.createElement('div');
            line.style.cssText = 'position:fixed;bottom:0;left:0;right:0;height:3px;z-index:997;pointer-events:none;background:linear-gradient(90deg,transparent,rgba(255,45,120,0.9),rgba(191,63,255,0.9),transparent);box-shadow:0 0 12px rgba(255,45,120,0.6);transition:bottom 0.32s ease,opacity 0.12s 0.32s';
            document.body.appendChild(line);
            setTimeout(function() { line.style.bottom='100vh'; setTimeout(function() { line.style.opacity='0'; setTimeout(function() { line.remove(); }, 140); }, 320); }, 15);
        };
        window.snonuxCloseEffect = function() {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:linear-gradient(180deg,transparent 50%,rgba(255,45,120,0.12) 100%);transition:opacity 0.2s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 230); }, 15);
        };
        window.snonuxScrollEffect = function(dir) {
            var isDown = dir === 'down';
            var thick = _wild ? '14px' : '5px';
            var d = document.createElement('div');
            // Synthwave: hot pink retrowave beam
            d.style.cssText = 'position:fixed;left:0;right:0;height:' + thick + ';z-index:9000;pointer-events:none;' +
                'background:linear-gradient(90deg,transparent,rgba(255,0,128,0.95),rgba(255,0,200,0.95),rgba(255,0,128,0.95),transparent);' +
                (isDown ? 'top:0;' : 'bottom:0;') +
                'transition:transform 0.28s ease,opacity 0.28s ease;';
            document.body.appendChild(d);
            setTimeout(function() { d.style.transform = isDown ? 'translateY(100vh)' : 'translateY(-100vh)'; d.style.opacity='0'; }, 16);
            setTimeout(function() { d.remove(); }, 360);
        };
        window.snonuxWildToggle = function() {
            _wild = !_wild;
            var b = document.getElementById('sno-wild-badge');
            if (b) b.classList.toggle('sno-wild-on', _wild);
        };
        window.snonuxNavEffect = function() {
            // Grid zoom forward
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-zoom'); setTimeout(function() { ov.classList.remove('sno-fx-zoom'); }, 330); }
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:linear-gradient(180deg,transparent 40%,rgba(255,45,120,0.2) 100%);transition:opacity 0.2s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 230); }, 20);
        };
        window.snonuxPageEffect = function() {
            // Warp-speed tunnel flash
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:radial-gradient(ellipse at 50% 100%,rgba(255,107,43,0.35) 0%,transparent 65%);transform:scaleY(0.2);transition:transform 0.22s ease,opacity 0.22s ease';
            document.body.appendChild(d);
            setTimeout(function() { d.style.transform='scaleY(1.4)'; d.style.opacity='0'; setTimeout(function() { d.remove(); }, 250); }, 20);
        };
    })();
