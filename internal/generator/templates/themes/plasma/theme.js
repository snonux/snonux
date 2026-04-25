
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,g=new THREE.Group(),t0=performance.now(),spec=[[0x00f0ff,1.45,0,0,0],[0xff00e0,1.05,1.2,0.4,0],[0xffee00,0.75,-1.1,-0.3,0]];
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(50,1,0.1,60);ca.position.z=6.5;
        spec.forEach(function(s){var m=new THREE.Mesh(new THREE.SphereGeometry(s[1],20,20),new THREE.MeshBasicMaterial({color:s[0],transparent:true,opacity:0.42,blending:THREE.AdditiveBlending,depthWrite:false}));
            m.position.set(s[2],s[3],s[4]);m.userData.ph=s;g.add(m);});
        sc.add(g);sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);var t=(now-t0)*0.001;
            g.children.forEach(function(c,ix){var ph=c.userData.ph;c.position.x=ph[2]+Math.sin(t*0.9+ix)*0.35;c.position.y=ph[3]+Math.cos(t*0.7+ix*1.3)*0.28;c.scale.setScalar(1+Math.sin(t*1.5+ix)*0.06);});
            ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();


    // Plasma WebGL: 12 large translucent spheres drifting on independent sine
    // paths with additive blending — overlapping blobs mix colours and pulse
    // like a lava lamp or plasma ball. Dark bg, cyan/magenta/yellow palette.
    (function() {
        var _wild = false, _snoTOffset = 0, _snoLastT = 0;
        var scene, camera, renderer, clock;
        var blobs = [];

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x050008);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 200);
            camera.position.set(0, 0, 40);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();

            // Each blob: [color, radius, baseX, baseY, baseZ, ampX, ampY, freqX, freqY, phaseX, phaseY]
            var cfg = [
                [0x00f0ff, 10, -5,  0, -15, 8,  6, 0.30, 0.40, 0.0, 1.0],
                [0xff00e0,  9,  8, -4, -18, 7,  8, 0.40, 0.30, 1.5, 0.5],
                [0xffee00,  8, -8,  6, -20, 6,  7, 0.50, 0.20, 3.0, 2.0],
                [0x00f0ff,  7,  4,  8, -12, 9,  5, 0.20, 0.50, 0.8, 3.5],
                [0xff00e0,  9, -6, -8, -16, 8,  6, 0.35, 0.45, 2.2, 1.2],
                [0xffee00, 11,  2,  2, -22, 7,  9, 0.25, 0.35, 4.0, 0.3],
                [0x8800ff,  8,-12,  4, -14, 6,  7, 0.45, 0.25, 1.0, 2.5],
                [0x00ff88,  7, 10, -6, -19, 8,  5, 0.30, 0.40, 3.5, 1.8],
                [0xff4400,  9,  0, 10, -17, 7,  8, 0.40, 0.30, 0.5, 4.0],
                [0x00f0ff,  6, -4, -4, -11, 5,  6, 0.55, 0.35, 2.8, 0.9],
                [0xff00e0, 10,  6,  0, -25, 9,  5, 0.20, 0.50, 1.3, 3.2],
                [0xffee00,  7,-10, -2, -13, 6,  8, 0.40, 0.30, 4.5, 1.5],
            ];

            cfg.forEach(function(c) {
                var geo = new THREE.SphereGeometry(c[1], 24, 24);
                var mat = new THREE.MeshBasicMaterial({
                    color: c[0], transparent: true, opacity: 0.18,
                    blending: THREE.AdditiveBlending, depthWrite: false
                });
                var mesh = new THREE.Mesh(geo, mat);
                mesh.position.set(c[2], c[3], c[4]);
                blobs.push({ mesh: mesh,
                    bx: c[2], by: c[3],
                    ax: c[5], ay: c[6],
                    fx: c[7], fy: c[8],
                    px: c[9], py: c[10] });
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
            var realT = clock.getElapsedTime();
            _snoTOffset += (realT - _snoLastT) * (_wild ? 7 : 0);
            _snoLastT = realT;
            var t = realT + _snoTOffset;
            var ampMult = _wild ? 3.5 : 1;
            blobs.forEach(function(b) {
                b.mesh.position.x = b.bx + b.ax * ampMult * Math.sin(t * b.fx + b.px);
                b.mesh.position.y = b.by + b.ay * ampMult * Math.cos(t * b.fy + b.py);
                // Blobs pulse in size when wild
                if (_wild) { b.mesh.scale.setScalar(1 + 0.3 * Math.sin(t * b.fy * 1.3)); }
                else { b.mesh.scale.setScalar(1); }
            });
            // Wild: camera sucked into plasma vortex
            if (_wild) {
                camera.position.x = Math.sin(realT * 0.32) * 14;
                camera.position.y = Math.sin(realT * 0.24) * 10;
                camera.position.z = 40 + Math.sin(realT * 0.19) * 14;
                camera.fov = 60 + Math.sin(realT * 0.41) * 16;
                camera.updateProjectionMatrix();
            } else {
                camera.position.set(0, 0, 40);
                if (camera.fov !== 60) { camera.fov = 60; camera.updateProjectionMatrix(); }
            }
            renderer.render(scene, camera);
        }

        initThree();

        // Plasma nav/wild effects — blob pulse burst on navigate, supernova on wild
        window.snonuxOpenEffect = function(post) {
            // Blobs merge to form the modal — expand with plasma splash
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-expand'); setTimeout(function() { modal.classList.remove('sno-modal-expand'); }, 420); }
            var r = post ? post.getBoundingClientRect() : {left: window.innerWidth/2, top: window.innerHeight/2, width: 0, height: 0};
            var colors = ['rgba(0,240,255,0.35)', 'rgba(255,0,224,0.3)', 'rgba(255,238,0,0.25)'];
            colors.forEach(function(col, i) {
                var b = document.createElement('div');
                var angle = (i/3)*Math.PI*2; var dist = 30+i*15;
                b.style.cssText = 'position:fixed;top:' + (r.top+r.height/2) + 'px;left:' + (r.left+r.width/2) + 'px;z-index:997;pointer-events:none;width:20px;height:20px;border-radius:50%;background:' + col + ';filter:blur(6px);transition:all 0.4s ease,opacity 0.4s';
                document.body.appendChild(b);
                setTimeout(function() { b.style.transform='translate(' + (Math.cos(angle)*dist) + 'px,' + (Math.sin(angle)*dist) + 'px) scale(8)'; b.style.opacity='0'; setTimeout(function() { b.remove(); }, 450); }, 15+i*30);
            });
        };
        window.snonuxCloseEffect = function() {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:radial-gradient(ellipse at center,rgba(255,0,224,0.12) 0%,transparent 60%);transition:opacity 0.2s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 230); }, 15);
        };
        window.snonuxScrollEffect = function(dir) {
            var isDown = dir === 'down';
            var thick = _wild ? '14px' : '5px';
            var d = document.createElement('div');
            // Plasma: shifting magenta-to-cyan rainbow
            d.style.cssText = 'position:fixed;left:0;right:0;height:' + thick + ';z-index:9000;pointer-events:none;' +
                'background:linear-gradient(90deg,transparent,rgba(255,0,128,0.9),rgba(128,0,255,0.9),rgba(0,200,255,0.9),transparent);' +
                (isDown ? 'top:0;' : 'bottom:0;') +
                'transition:transform 0.3s ease,opacity 0.3s ease;';
            document.body.appendChild(d);
            setTimeout(function() { d.style.transform = isDown ? 'translateY(100vh)' : 'translateY(-100vh)'; d.style.opacity='0'; }, 16);
            setTimeout(function() { d.remove(); }, 380);
        };
        window.snonuxWildToggle = function() {
            _wild = !_wild;
            var b = document.getElementById('sno-wild-badge');
            if (b) b.classList.toggle('sno-wild-on', _wild);
        };
        window.snonuxNavEffect = function() {
            // Blob pulse burst — rainbow flash radiating from center
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:radial-gradient(ellipse at center,rgba(0,240,255,0.25) 0%,rgba(255,0,224,0.15) 45%,transparent 70%);transform:scale(0.5);transition:transform 0.28s ease,opacity 0.28s ease';
            document.body.appendChild(d);
            setTimeout(function() { d.style.transform='scale(1.3)'; d.style.opacity='0'; setTimeout(function() { d.remove(); }, 310); }, 20);
        };
        window.snonuxPageEffect = function() {
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-zoom'); setTimeout(function() { ov.classList.remove('sno-fx-zoom'); }, 330); }
        };
    })();
