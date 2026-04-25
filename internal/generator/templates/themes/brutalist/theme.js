
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,g=new THREE.Group(),t0=performance.now();
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(50,1,0.1,80);ca.position.z=8;
        var b1=new THREE.LineSegments(new THREE.EdgesGeometry(new THREE.BoxGeometry(3.4,2.4,2.4)),new THREE.LineBasicMaterial({color:0xffffff}));
        var b2=new THREE.LineSegments(new THREE.EdgesGeometry(new THREE.BoxGeometry(2.2,1.6,1.6)),new THREE.LineBasicMaterial({color:0xff2200}));
        b2.position.set(0.3,0.2,0.5);g.add(b1);g.add(b2);sc.add(g);sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);var t=(now-t0)*0.001;g.rotation.x=t*0.51;g.rotation.y=t*0.73;ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();


    // Brutalist WebGL: harsh slowly-rotating boxes — solid white and wireframe red.
    // No fog, no softness. Pure geometric violence against the black void.
    (function() {
        var _wild = false;
        var scene, camera, renderer, clock;
        var boxes = [];

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x000000);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 200);
            camera.position.set(0, 0, 40);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: false });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();

            // Box configurations: [size, posX, posY, posZ, rotSpeedX, rotSpeedY, wireframe, color]
            var configs = [
                [10, 0,   0,  0,   0.002, 0.005, false, 0xffffff],
                [6,  18, -6,  -8,  0.004, 0.003, true,  0xff2200],
                [7,  -16, 5, -10,  0.003, 0.006, true,  0xff2200],
                [5,  8,  12, -5,   0.006, 0.002, false, 0xff2200],
                [4,  -10,-10, -3,  0.005, 0.004, false, 0xffffff],
            ];

            configs.forEach(function(c) {
                var geo = new THREE.BoxGeometry(c[0], c[0], c[0]);
                var mat = new THREE.MeshBasicMaterial({ color: c[7], wireframe: c[6] });
                var mesh = new THREE.Mesh(geo, mat);
                mesh.position.set(c[1], c[2], c[3]);
                boxes.push({ mesh: mesh, rx: c[4], ry: c[5] });
                scene.add(mesh);
            });

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
            var sm = _wild ? 15 : 1;
            boxes.forEach(function(b) {
                b.mesh.rotation.x += b.rx * sm;
                b.mesh.rotation.y += b.ry * sm;
                // Wild mode: random jitter on positions
                if (_wild) { b.mesh.position.x += (Math.random()-0.5)*0.4; b.mesh.position.y += (Math.random()-0.5)*0.4; }
            });
            renderer.render(scene, camera);
        }

        initThree();

        // Brutalist nav/wild effects — violent shake on navigate, geometric chaos on wild
        window.snonuxOpenEffect = function() {
            // Expand violently from nothing — pure brutalist impact
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-expand'); setTimeout(function() { modal.classList.remove('sno-modal-expand'); }, 420); }
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:997;pointer-events:none;background:rgba(255,34,0,0.22);transition:opacity 0.14s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 170); }, 15);
        };
        window.snonuxCloseEffect = function() {
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-shake'); setTimeout(function() { ov.classList.remove('sno-fx-shake'); }, 360); }
        };
        window.snonuxScrollEffect = function(dir) {
            var isDown = dir === 'down';
            var thick = _wild ? '14px' : '5px';
            var d = document.createElement('div');
            // Brutalist: harsh black-and-white hard edge
            d.style.cssText = 'position:fixed;left:0;right:0;height:' + thick + ';z-index:9000;pointer-events:none;' +
                'background:linear-gradient(90deg,rgba(0,0,0,0.95),rgba(255,255,255,0.95),rgba(0,0,0,0.95));' +
                (isDown ? 'top:0;' : 'bottom:0;') +
                'transition:transform 0.25s ease,opacity 0.25s ease;';
            document.body.appendChild(d);
            setTimeout(function() { d.style.transform = isDown ? 'translateY(100vh)' : 'translateY(-100vh)'; d.style.opacity='0'; }, 16);
            setTimeout(function() { d.remove(); }, 320);
        };
        window.snonuxWildToggle = function() {
            _wild = !_wild;
            var b = document.getElementById('sno-wild-badge');
            if (b) b.classList.toggle('sno-wild-on', _wild);
        };
        window.snonuxNavEffect = function() {
            // Violent double shake + red flash
            var ov = document.querySelector('.overlay');
            if (ov) {
                ov.classList.add('sno-fx-shake');
                setTimeout(function() { ov.classList.remove('sno-fx-shake'); setTimeout(function() { ov.classList.add('sno-fx-shake'); setTimeout(function() { ov.classList.remove('sno-fx-shake'); }, 380); }, 50); }, 400);
            }
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(255,34,0,0.28);transition:opacity 0.18s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 200); }, 25);
        };
        window.snonuxPageEffect = function() {
            // Color inversion flash
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:#fff;mix-blend-mode:difference;transition:opacity 0.15s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 180); }, 20);
        };
    })();
