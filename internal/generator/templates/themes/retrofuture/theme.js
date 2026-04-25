
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
        var atomic=new THREE.Mesh(new THREE.SphereGeometry(1.35,28,28),new THREE.MeshBasicMaterial({color:0x00d9c0,transparent:true,opacity:0.9}));
        atomic.position.y=0;g.add(atomic);
        var ringH=new THREE.Mesh(new THREE.TorusGeometry(2.1,0.06,8,64),new THREE.MeshBasicMaterial({color:0xff6b9d,transparent:true,opacity:0.8}));
        g.add(ringH);
        var ringV=new THREE.Mesh(new THREE.TorusGeometry(2.1,0.06,8,64),new THREE.MeshBasicMaterial({color:0xff6b9d,transparent:true,opacity:0.8}));
        ringV.rotation.x=Math.PI/2;g.add(ringV);
        var gr=new THREE.Mesh(new THREE.PlaneGeometry(28,28,20,20),new THREE.MeshBasicMaterial({color:0x00d9c0,wireframe:true,transparent:true,opacity:0.25}));
        gr.rotation.x=-Math.PI/2.4;gr.position.y=-2.8;g.add(gr);
        sc.add(g);sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);var t=(now-t0)*0.001;g.rotation.y=Math.sin(t*0.3)*0.12;atomic.position.y=Math.sin(t*1.1)*0.12;ringH.scale.setScalar(1+Math.sin(t*2.2)*0.06);ringV.scale.setScalar(1+Math.cos(t*2.2)*0.06);ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();


    // Retrofuture WebGL: atomic orb with crossed electron rings, receding teal grid floor,
    // chrome metallic star particles, and a slow drifting camera orbit.
    (function() {
        var _wild = false, _snoTOffset = 0, _snoLastT = 0;
        var scene, camera, renderer, clock;
        var atomic, ringH, ringV, stars;

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x0a0121);
            scene.fog = new THREE.Fog(0x0a0121, 60, 180);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 300);
            camera.position.set(0, 10, 45);
            camera.lookAt(0, -5, -10);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();

            // Glowing atomic teal orb
            atomic = new THREE.Mesh(
                new THREE.SphereGeometry(12, 32, 16),
                new THREE.MeshBasicMaterial({ color: 0x00d9c0 })
            );
            atomic.position.set(0, -8, -55);
            scene.add(atomic);

            // Horizontal electron ring — pink
            ringH = new THREE.Mesh(
                new THREE.TorusGeometry(15, 0.15, 8, 64),
                new THREE.MeshBasicMaterial({ color: 0xff6b9d })
            );
            ringH.position.copy(atomic.position);
            scene.add(ringH);

            // Vertical electron ring — coral
            ringV = new THREE.Mesh(
                new THREE.TorusGeometry(15, 0.15, 8, 64),
                new THREE.MeshBasicMaterial({ color: 0xff8c42 })
            );
            ringV.position.copy(atomic.position);
            ringV.rotation.x = Math.PI / 2;
            scene.add(ringV);

            // Receding grid floor — teal
            var grid = new THREE.GridHelper(200, 40, 0x00d9c0, 0x1a0840);
            grid.position.set(0, -18, -30);
            scene.add(grid);

            // 1200 chrome star particles scattered in a sphere shell
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
                color: 0xffd700, size: 0.22, transparent: true, opacity: 0.8
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
            // Atomic orb pulses; wild meltdown mode makes it throb intensely
            var pulse = _wild ? 1 + 0.1 * Math.sin(t * 12) : 1 + 0.015 * Math.sin(t * 1.5);
            atomic.scale.setScalar(pulse);
            ringH.rotation.z = t * 0.4;
            ringV.rotation.z = -t * 0.3;
            // Rings tilt chaotically in wild — full 3D tumble
            if (_wild) {
                ringH.rotation.x = Math.sin(realT * 2.3) * 0.8;
                ringV.rotation.y = Math.cos(realT * 1.9) * 0.9;
            } else {
                ringH.rotation.x = 0;
                ringV.rotation.y = 0;
            }
            // Camera: meltdown spiral in wild, slow orbit otherwise
            if (_wild) {
                camera.position.x = Math.sin(realT * 0.36) * 16;
                camera.position.y = 10 + Math.sin(realT * 0.29) * 10;
                camera.position.z = 45 + Math.sin(realT * 0.22) * 16;
                camera.fov = 60 + Math.sin(realT * 0.47) * 18;
                camera.updateProjectionMatrix();
            } else {
                camera.position.x = Math.sin(t * 0.07) * 5;
                camera.position.y = 10 + Math.sin(t * 0.05) * 2;
                camera.position.z = 45;
                if (camera.fov !== 60) { camera.fov = 60; camera.updateProjectionMatrix(); }
            }
            camera.lookAt(0, -5, -10);
            renderer.render(scene, camera);
        }

        initThree();

        // Retrofuture nav/wild effects — atomic pulse on navigate, meltdown on wild
        window.snonuxOpenEffect = function(post) {
            // Atomic rings expand outward from post — zoom into modal
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-zoom'); setTimeout(function() { modal.classList.remove('sno-modal-zoom'); }, 400); }
            var r = post ? post.getBoundingClientRect() : {left: window.innerWidth/2, top: window.innerHeight/2, width: 0, height: 0};
            [0, 80, 160].forEach(function(delay) {
                var ring = document.createElement('div');
                ring.style.cssText = 'position:fixed;top:' + (r.top+r.height/2-6) + 'px;left:' + (r.left+r.width/2-6) + 'px;z-index:997;pointer-events:none;width:12px;height:12px;border-radius:50%;border:2px solid rgba(0,217,192,0.7);transition:all 0.42s ease,opacity 0.42s';
                document.body.appendChild(ring);
                setTimeout(function() { ring.style.transform='scale(25)'; ring.style.opacity='0'; setTimeout(function() { ring.remove(); }, 460); }, delay + 15);
            });
        };
        window.snonuxCloseEffect = function() {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(0,217,192,0.1);transition:opacity 0.2s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 230); }, 15);
        };
        window.snonuxScrollEffect = function(dir) {
            var isDown = dir === 'down';
            var thick = _wild ? '14px' : '5px';
            var d = document.createElement('div');
            // Retrofuture: atomic orange-gold sweep
            d.style.cssText = 'position:fixed;left:0;right:0;height:' + thick + ';z-index:9000;pointer-events:none;' +
                'background:linear-gradient(90deg,transparent,rgba(255,140,0,0.9),rgba(255,80,0,0.9),rgba(255,140,0,0.9),transparent);' +
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
            // Atomic pulse — rings expand briefly as CSS overlay
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;top:50%;left:50%;transform:translate(-50%,-50%) scale(0.1);z-index:998;pointer-events:none;width:100vmax;height:100vmax;border-radius:50%;border:3px solid rgba(0,217,192,0.7);transition:transform 0.3s ease,opacity 0.3s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.transform='translate(-50%,-50%) scale(1.2)'; d.style.opacity='0'; setTimeout(function() { d.remove(); }, 330); }, 15);
        };
        window.snonuxPageEffect = function() {
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-zoom'); setTimeout(function() { ov.classList.remove('sno-fx-zoom'); }, 330); }
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(0,217,192,0.15);transition:opacity 0.2s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 230); }, 20);
        };
    })();
