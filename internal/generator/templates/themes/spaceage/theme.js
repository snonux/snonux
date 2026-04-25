
    // Splash WebGL: slowly rotating torus (space station ring) + star field.
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,g=new THREE.Group(),t0=performance.now();
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(50,1,0.1,60);ca.position.z=8;
        // Torus ring — space station
        var tor=new THREE.Mesh(new THREE.TorusGeometry(2.2,0.45,16,80),new THREE.MeshBasicMaterial({color:0x00e8e8,wireframe:true,transparent:true,opacity:0.85}));
        g.add(tor);
        // Inner hub
        var hub=new THREE.Mesh(new THREE.SphereGeometry(0.38,12,12),new THREE.MeshBasicMaterial({color:0x00e8e8,wireframe:true,transparent:true,opacity:0.6}));
        g.add(hub);
        sc.add(g);sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);var t=(now-t0)*0.001;g.rotation.x=t*0.28;g.rotation.y=t*0.45;hub.rotation.z=t*1.1;ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();


    // Space Age WebGL: toroidal space station ring + three satellite pods orbiting it
    // + a slowly rotating planet sphere + dense star field. Teal wireframe throughout.
    (function() {
        var _wild = false, _snoTOffset = 0, _snoLastT = 0;
        var scene, camera, renderer, clock;
        var station, pods = [], planet;

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x030a0f);
            scene.fog = new THREE.Fog(0x030a0f, 60, 160);

            camera = new THREE.PerspectiveCamera(58, window.innerWidth / window.innerHeight, 0.1, 300);
            camera.position.set(0, 14, 42);
            camera.lookAt(0, 0, 0);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();

            var tealMat  = new THREE.MeshBasicMaterial({ color: 0x00e8e8, wireframe: true });
            var dimMat   = new THREE.MeshBasicMaterial({ color: 0x1a4455, wireframe: true });
            var silverMat= new THREE.MeshBasicMaterial({ color: 0xc8d8e0, wireframe: true, transparent: true, opacity: 0.55 });

            // Main toroidal space station — the visual centrepiece
            station = new THREE.Mesh(new THREE.TorusGeometry(14, 2.8, 20, 100), tealMat.clone());
            scene.add(station);

            // Central hub sphere
            var hub = new THREE.Mesh(new THREE.SphereGeometry(2.2, 16, 16), dimMat.clone());
            scene.add(hub);

            // Three spoke arms from hub to ring
            for (var s = 0; s < 3; s++) {
                var angle = (s / 3) * Math.PI * 2;
                var spoke = new THREE.Mesh(
                    new THREE.CylinderGeometry(0.12, 0.12, 14, 6),
                    dimMat.clone()
                );
                spoke.rotation.z = angle + Math.PI / 2;
                spoke.position.set(Math.cos(angle) * 7, Math.sin(angle) * 7, 0);
                scene.add(spoke);
            }

            // Three satellite pods orbiting the station on a wider ring
            for (var p = 0; p < 3; p++) {
                var pod = new THREE.Mesh(new THREE.OctahedronGeometry(1.1, 0), silverMat.clone());
                pod._baseAngle = (p / 3) * Math.PI * 2;
                pods.push(pod);
                scene.add(pod);
            }

            // Distant planet — slowly rotating sphere
            planet = new THREE.Mesh(
                new THREE.SphereGeometry(9, 24, 24),
                new THREE.MeshBasicMaterial({ color: 0x0a2a38, wireframe: true, transparent: true, opacity: 0.6 })
            );
            planet.position.set(-55, -18, -80);
            scene.add(planet);

            // 1500 star particles spread through deep space
            var starPos = new Float32Array(1500 * 3);
            for (var i = 0; i < 1500 * 3; i++) {
                starPos[i] = (Math.random() - 0.5) * 220;
            }
            var starGeo = new THREE.BufferGeometry();
            starGeo.setAttribute('position', new THREE.BufferAttribute(starPos, 3));
            scene.add(new THREE.Points(starGeo, new THREE.PointsMaterial({
                color: 0xc8d8e0, size: 0.18, transparent: true, opacity: 0.65
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
            _snoTOffset += (realT - _snoLastT) * (_wild ? 10 : 0);
            _snoLastT = realT;
            var t = realT + _snoTOffset;

            // Station rotates; wild mode = warp-speed orbital mechanics
            station.rotation.y = t * 0.12;
            station.rotation.x = Math.sin(t * 0.07) * 0.3;
            if (_wild) { station.scale.setScalar(1 + 0.15 * Math.sin(realT * 4.5)); }
            else { station.scale.setScalar(1); }

            // Satellite pods orbit at expanded + oscillating radius in wild mode
            var podR = _wild ? 24 + Math.sin(realT * 0.7) * 10 : 24;
            var podY = _wild ? 12 : 5;
            for (var i = 0; i < pods.length; i++) {
                var angle = pods[i]._baseAngle + t * 0.38;
                pods[i].position.set(
                    Math.cos(angle) * podR,
                    Math.sin(angle * 0.6) * podY,
                    Math.sin(angle) * podR
                );
                pods[i].rotation.x = t * 0.7 + i;
                pods[i].rotation.z = t * 0.5 + i;
            }

            // Planet drifts very slowly
            planet.rotation.y = t * 0.04;

            // Camera: chaotic fly-by in wild, gentle drift otherwise
            if (_wild) {
                camera.position.x = Math.sin(realT * 0.34) * 18;
                camera.position.y = 14 + Math.sin(realT * 0.27) * 10;
                camera.position.z = 42 + Math.sin(realT * 0.21) * 16;
                camera.fov = 58 + Math.sin(realT * 0.44) * 18;
                camera.updateProjectionMatrix();
            } else {
                camera.position.x = Math.sin(t * 0.06) * 5;
                camera.position.y = 14 + Math.sin(t * 0.09) * 2;
                camera.position.z = 42;
                if (camera.fov !== 58) { camera.fov = 58; camera.updateProjectionMatrix(); }
            }

            renderer.render(scene, camera);
        }

        initThree();

        // Spaceage nav/wild effects — star streak on navigate, hyperspace jump on wild
        window.snonuxOpenEffect = function(post) {
            // Docking approach — expand from post position
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-expand'); setTimeout(function() { modal.classList.remove('sno-modal-expand'); }, 420); }
            var r = post ? post.getBoundingClientRect() : {left: window.innerWidth/2, top: window.innerHeight/2, width: 0, height: 0};
            var ring = document.createElement('div');
            ring.style.cssText = 'position:fixed;top:' + (r.top+r.height/2-6) + 'px;left:' + (r.left+r.width/2-6) + 'px;z-index:997;pointer-events:none;width:12px;height:12px;border-radius:50%;border:2px solid rgba(0,232,232,0.8);transition:all 0.38s ease,opacity 0.38s';
            document.body.appendChild(ring);
            setTimeout(function() { ring.style.transform='scale(28)'; ring.style.opacity='0'; setTimeout(function() { ring.remove(); }, 420); }, 15);
        };
        window.snonuxCloseEffect = function() {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(0,232,232,0.1);transition:opacity 0.2s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 230); }, 15);
        };
        window.snonuxScrollEffect = function(dir) {
            var isDown = dir === 'down';
            var thick = _wild ? '14px' : '5px';
            var d = document.createElement('div');
            // Space Age: teal energy pulse
            d.style.cssText = 'position:fixed;left:0;right:0;height:' + thick + ';z-index:9000;pointer-events:none;' +
                'background:linear-gradient(90deg,transparent,rgba(0,220,200,0.9),rgba(0,180,220,0.9),rgba(0,220,200,0.9),transparent);' +
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
            // Star streak flash — radial lines from center
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:radial-gradient(ellipse at center,rgba(0,232,232,0.2) 0%,transparent 55%);transform:scale(0.4);transition:transform 0.25s ease,opacity 0.25s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.transform='scale(1.6)'; d.style.opacity='0'; setTimeout(function() { d.remove(); }, 280); }, 15);
        };
        window.snonuxPageEffect = function() {
            // Hyperspace jump — white tunnel flash
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:radial-gradient(ellipse at center,rgba(255,255,255,0.3) 0%,rgba(0,232,232,0.15) 45%,transparent 70%);transition:opacity 0.2s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 230); }, 20);
        };
    })();
