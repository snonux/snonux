
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,m,t0=performance.now();
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(48,1,0.1,60);ca.position.z=7;
        m=new THREE.Mesh(new THREE.IcosahedronGeometry(2.3,1),new THREE.MeshBasicMaterial({color:0x33ff33,wireframe:true,transparent:true,opacity:0.88}));
        sc.add(m);sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);var t=(now-t0)*0.001;m.rotation.x=t*0.62;m.rotation.y=t*0.88;ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();


    // Terminal WebGL scene: phosphor-green icosahedron wireframe + torus particle ring.
    // The scene sits behind the CRT scanline overlay (z-index:999) and the UI (z-index:10).
    (function() {
        var _wild = false, _snoTOffset = 0, _snoLastT = 0;
        var scene, camera, renderer, icosa, particles;
        var clock = new THREE.Clock();

        function initThree() {
            // Scene with pure-black background and distance fog
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x000000);
            scene.fog = new THREE.Fog(0x000000, 20, 80);

            // Perspective camera positioned in front of the orb
            camera = new THREE.PerspectiveCamera(60, window.innerWidth / window.innerHeight, 0.1, 200);
            camera.position.set(0, 0, 30);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

            // Large green phosphor wireframe icosahedron — the central CRT orb
            var icoGeo = new THREE.IcosahedronGeometry(8, 2);
            var icoMat = new THREE.MeshBasicMaterial({ color: 0x33ff33, wireframe: true });
            icosa = new THREE.Mesh(icoGeo, icoMat);
            scene.add(icosa);

            // 400 dim particles arranged on a torus path around the icosahedron
            var torusGeo = new THREE.TorusGeometry(14, 3, 16, 100);
            var positions = torusGeo.attributes.position;
            var ptGeo = new THREE.BufferGeometry();
            var pts = new Float32Array(400 * 3);
            for (var i = 0; i < 400; i++) {
                // Sample vertices from the torus geometry to place particles on its surface
                var idx = Math.floor(Math.random() * positions.count);
                pts[i * 3]     = positions.getX(idx);
                pts[i * 3 + 1] = positions.getY(idx);
                pts[i * 3 + 2] = positions.getZ(idx);
            }
            ptGeo.setAttribute('position', new THREE.BufferAttribute(pts, 3));
            var ptMat = new THREE.PointsMaterial({ color: 0x1a7a1a, size: 0.18 });
            particles = new THREE.Points(ptGeo, ptMat);
            scene.add(particles);

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
            _snoTOffset += (realT - _snoLastT) * (_wild ? 11 : 0);
            _snoLastT = realT;
            var t = realT + _snoTOffset;
            // Slow multi-axis rotation; wild mode overloads the phosphor orb
            icosa.rotation.x = t * 0.12;
            icosa.rotation.y = t * 0.18;
            icosa.rotation.z = t * 0.07;
            // Counter-rotate particles for visual contrast
            particles.rotation.y = -t * 0.08;
            particles.rotation.x =  t * 0.04;
            renderer.render(scene, camera);
        }

        initThree();

        // Terminal nav/wild effects — cursor glitch on navigate, buffer overflow on wild
        window.snonuxOpenEffect = function() {
            // Slide in like terminal output being printed
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-slide'); setTimeout(function() { modal.classList.remove('sno-modal-slide'); }, 360); }
            // Phosphor scan from top to bottom
            var scan = document.createElement('div');
            scan.style.cssText = 'position:fixed;top:0;left:0;right:0;height:2px;z-index:997;pointer-events:none;background:rgba(51,255,51,0.6);box-shadow:0 0 8px rgba(51,255,51,0.4);transition:top 0.3s linear,opacity 0.1s 0.3s';
            document.body.appendChild(scan);
            setTimeout(function() { scan.style.top='100vh'; setTimeout(function() { scan.style.opacity='0'; setTimeout(function() { scan.remove(); }, 120); }, 300); }, 15);
        };
        window.snonuxCloseEffect = function() {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(51,255,51,0.1);transition:opacity 0.18s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 200); }, 15);
            document.body.style.animationDuration = '9s';
        };
        window.snonuxScrollEffect = function(dir) {
            var isDown = dir === 'down';
            var thick = _wild ? '14px' : '5px';
            var d = document.createElement('div');
            // Terminal: phosphor green scan
            d.style.cssText = 'position:fixed;left:0;right:0;height:' + thick + ';z-index:9000;pointer-events:none;' +
                'background:linear-gradient(90deg,transparent,rgba(57,255,20,0.9),rgba(20,200,10,0.9),rgba(57,255,20,0.9),transparent);' +
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
            // Toggle intense scanline strobe in wild mode
            document.body.style.animationDuration = _wild ? '0.4s' : '9s';
        };
        window.snonuxNavEffect = function() {
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-glitch'); setTimeout(function() { ov.classList.remove('sno-fx-glitch'); }, 300); }
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(51,255,51,0.13);transition:opacity 0.18s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 210); }, 25);
        };
        window.snonuxPageEffect = function() {
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-glitch'); setTimeout(function() { ov.classList.remove('sno-fx-glitch'); setTimeout(function() { ov.classList.add('sno-fx-glitch'); setTimeout(function() { ov.classList.remove('sno-fx-glitch'); }, 280); }, 35); }, 300); }
        };
    })();
