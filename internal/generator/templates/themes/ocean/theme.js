
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,g=new THREE.Group(),t0=performance.now(),i;
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(50,1,0.1,70);ca.position.set(0,0.2,9);
        for(i=0;i<5;i++){var b=new THREE.Mesh(new THREE.SphereGeometry(0.25+Math.random()*0.35,12,12),new THREE.MeshBasicMaterial({color:0x48cae4,transparent:true,opacity:0.65}));
            b.position.set((Math.random()-0.5)*7,(Math.random()-0.5)*4,(Math.random()-0.5)*3);b.userData.dy=0.02+Math.random()*0.03;b.userData.x=b.position.x;b.userData.y0=b.position.y;g.add(b);}
        var jelly=new THREE.Mesh(new THREE.SphereGeometry(1.1,16,16),new THREE.MeshBasicMaterial({color:0x00b4d8,transparent:true,opacity:0.35,wireframe:true}));
        g.add(jelly);sc.add(g);sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);var t=(now-t0)*0.001;jelly.rotation.y=t*0.4;
            g.children.forEach(function(c){if(c.userData.dy){c.position.y+=Math.sin(t*2+c.userData.x)*0.008;c.position.x=c.userData.x+Math.sin(t+c.userData.y0)*0.15;}});ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();


    // Ocean WebGL: dramatic wave surface + sea rock spires + bioluminescent
    // jellyfish + rising bubbles + a slow whale cruising the deep.
    (function() {
        var _wild = false, _snoTOffset = 0, _snoLastT = 0;
        var scene, camera, renderer, clock;
        var waveGeo, waveMesh, sunLight;
        var whale, jellyfish = [];
        var BUBBLE_COUNT = 600;
        var bubblePos, bubbleVY;

        function buildWaves() {
            // High-density plane for smooth vertex displacement
            waveGeo = new THREE.PlaneGeometry(300, 300, 100, 100);
            waveMesh = new THREE.Mesh(waveGeo, new THREE.MeshPhongMaterial({
                color: 0x0077b6, emissive: 0x023e8a, emissiveIntensity: 0.25,
                transparent: true, opacity: 0.88, side: THREE.DoubleSide, shininess: 80
            }));
            waveMesh.rotation.x = -Math.PI / 2;
            waveMesh.position.y = 0;
            scene.add(waveMesh);
        }

        function buildRocks() {
            // 5 jagged sea rock spires poking above the wave baseline
            var rockPositions = [[-30,0,-30],[20,-2,-20],[-15,2,-45],[35,-1,-35],[-45,1,-25]];
            rockPositions.forEach(function(p) {
                var h = 8 + Math.random() * 10;
                var rock = new THREE.Mesh(
                    new THREE.ConeGeometry(2 + Math.random(), h, 6),
                    new THREE.MeshPhongMaterial({ color: 0x023e8a, emissive: 0x00b4d8, emissiveIntensity: 0.15 })
                );
                rock.position.set(p[0], p[1] + h / 2 - 3, p[2]);
                scene.add(rock);
            });
        }

        function buildJellyfish() {
            // Bioluminescent jellyfish: torus body + cone cap, additive blending
            var jPos = [[-12, 6,-15],[18,10,-22],[-25,4,-18],[8,8,-30]];
            jPos.forEach(function(p) {
                var body = new THREE.Mesh(
                    new THREE.TorusGeometry(2.2, 0.5, 12, 24),
                    new THREE.MeshBasicMaterial({ color: 0x48cae4, transparent: true, opacity: 0.5, blending: THREE.AdditiveBlending, depthWrite: false })
                );
                var cap = new THREE.Mesh(
                    new THREE.SphereGeometry(2.2, 12, 8, 0, Math.PI * 2, 0, Math.PI / 2),
                    new THREE.MeshBasicMaterial({ color: 0x00b4d8, transparent: true, opacity: 0.35, blending: THREE.AdditiveBlending, depthWrite: false, side: THREE.DoubleSide })
                );
                cap.position.y = 0.5;
                body.add(cap);
                body.position.set(p[0], p[1], p[2]);
                jellyfish.push({ mesh: body, baseY: p[1], phase: Math.random() * Math.PI * 2 });
                scene.add(body);
            });
        }

        function buildWhale() {
            // Dark elongated flattened sphere — whale silhouette in the deep
            var geo = new THREE.SphereGeometry(1, 16, 8);
            whale = new THREE.Mesh(geo, new THREE.MeshBasicMaterial({ color: 0x011f40, transparent: true, opacity: 0.7 }));
            whale.scale.set(12, 3, 5);
            whale.position.set(-60, -8, -20);
            scene.add(whale);
        }

        function buildBubbles() {
            bubblePos = new Float32Array(BUBBLE_COUNT * 3);
            bubbleVY  = new Float32Array(BUBBLE_COUNT);
            for (var i = 0; i < BUBBLE_COUNT; i++) {
                bubblePos[i*3]   = (Math.random() - 0.5) * 100;
                bubblePos[i*3+1] = -15 - Math.random() * 15;
                bubblePos[i*3+2] = (Math.random() - 0.5) * 60 - 10;
                bubbleVY[i] = 0.04 + Math.random() * 0.06;
            }
            var geo = new THREE.BufferGeometry();
            geo.setAttribute('position', new THREE.BufferAttribute(bubblePos, 3));
            scene.add(new THREE.Points(geo, new THREE.PointsMaterial({
                color: 0xcaf0f8, size: 0.18, transparent: true, opacity: 0.6
            })));
            return geo;
        }

        var bubbleGeo;

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x03045e);
            scene.fog = new THREE.Fog(0x03045e, 40, 130);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 220);
            camera.position.set(0, 20, 55);
            camera.lookAt(0, 0, 0);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();

            scene.add(new THREE.AmbientLight(0x023e8a, 0.5));
            sunLight = new THREE.PointLight(0x48cae4, 2.5, 100);
            sunLight.position.set(0, 30, 10);
            scene.add(sunLight);
            var deepLight = new THREE.PointLight(0x0077b6, 1.5, 60);
            deepLight.position.set(0, -10, 0);
            scene.add(deepLight);

            buildWaves();
            buildRocks();
            buildJellyfish();
            buildWhale();
            bubbleGeo = buildBubbles();

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
            _snoTOffset += (realT - _snoLastT) * (_wild ? 6 : 0);
            _snoLastT = realT;
            var t = realT + _snoTOffset;
            var waveAmp = _wild ? 2.8 : 1; // tsunami amplitude in wild mode
            var pos = waveGeo.attributes.position;

            // Dramatic overlapping waves — tsunami amplitude in wild mode
            for (var i = 0; i < pos.count; i++) {
                var x = pos.getX(i), z = pos.getZ(i);
                pos.setY(i,
                    Math.sin(x * 0.04 + t * 1.1) * 3.2 * waveAmp +
                    Math.cos(z * 0.06 + t * 0.85) * 2.4 * waveAmp +
                    Math.sin((x + z) * 0.025 + t * 0.6) * 1.5 * waveAmp
                );
            }
            pos.needsUpdate = true;
            waveGeo.computeVertexNormals();

            // Jellyfish bob and slowly drift horizontally
            jellyfish.forEach(function(j) {
                j.mesh.position.y = j.baseY + Math.sin(t * 0.8 + j.phase) * 1.2;
                j.mesh.position.x += 0.005 * Math.sin(t * 0.3 + j.phase);
                j.mesh.rotation.y += 0.006;
            });

            // Whale cruises across at depth, wraps around
            whale.position.x += 0.04;
            if (whale.position.x > 80) whale.position.x = -80;
            whale.position.y = -8 + Math.sin(t * 0.15) * 2;

            // Rising bubbles — explode upward in wild tsunami mode
            var bp = bubbleGeo.attributes.position;
            var bMult = _wild ? 8 : 1;
            for (var bi = 0; bi < BUBBLE_COUNT; bi++) {
                bubblePos[bi*3+1] += bubbleVY[bi] * bMult;
                if (bubblePos[bi*3+1] > 8) {
                    bubblePos[bi*3]   = (Math.random() - 0.5) * 100;
                    bubblePos[bi*3+1] = -15 - Math.random() * 10;
                }
            }
            bp.needsUpdate = true;

            // Sunlight orbits above
            sunLight.position.x = Math.cos(t * 0.2) * 35;
            sunLight.position.z = Math.sin(t * 0.2) * 35;

            renderer.render(scene, camera);
        }

        initThree();

        // Ocean nav/wild effects — wave surge on navigate, tsunami on wild
        window.snonuxOpenEffect = function(post) {
            // Fly up from depth — like surfacing from the ocean
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-fly'); setTimeout(function() { modal.classList.remove('sno-modal-fly'); }, 390); }
            // Bubble burst from post
            var r = post ? post.getBoundingClientRect() : {left: window.innerWidth/2, top: window.innerHeight/2, width: 0, height: 0};
            for (var i = 0; i < 8; i++) {
                (function(i) {
                    var b = document.createElement('div');
                    var sz = 4 + Math.random() * 8;
                    b.style.cssText = 'position:fixed;top:' + (r.top + r.height * 0.7) + 'px;left:' + (r.left + r.width * 0.3 + Math.random() * r.width * 0.4) + 'px;z-index:997;pointer-events:none;width:' + sz + 'px;height:' + sz + 'px;border-radius:50%;border:1px solid rgba(72,202,228,0.7);background:rgba(0,180,216,0.2);transition:all 0.5s ease,opacity 0.5s';
                    document.body.appendChild(b);
                    setTimeout(function() { b.style.transform='translateY(-' + (60+Math.random()*60) + 'px) scale(1.5)'; b.style.opacity='0'; setTimeout(function() { b.remove(); }, 560); }, 20 + i*40);
                })(i);
            }
        };
        window.snonuxCloseEffect = function() {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(0,180,216,0.12);transition:opacity 0.2s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 230); }, 15);
        };
        window.snonuxScrollEffect = function(dir) {
            var isDown = dir === 'down';
            var thick = _wild ? '14px' : '5px';
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;left:0;right:0;height:' + thick + ';z-index:9000;pointer-events:none;' +
                'background:linear-gradient(90deg,transparent,rgba(0,119,182,0.9),rgba(0,180,216,0.9),rgba(0,119,182,0.9),transparent);' +
                (isDown ? 'top:0;' : 'bottom:0;') +
                'transition:transform 0.32s ease,opacity 0.32s ease;';
            document.body.appendChild(d);
            setTimeout(function() { d.style.transform = isDown ? 'translateY(100vh)' : 'translateY(-100vh)'; d.style.opacity='0'; }, 16);
            setTimeout(function() { d.remove(); }, 400);
        };
        window.snonuxWildToggle = function() {
            _wild = !_wild;
            var b = document.getElementById('sno-wild-badge');
            if (b) b.classList.toggle('sno-wild-on', _wild);
        };
        window.snonuxNavEffect = function() {
            // Wave surge: content skews briefly
            var ov = document.querySelector('.overlay');
            if (!ov) return;
            ov.style.transition = 'transform 0.1s';
            ov.style.transform = 'skewX(-1.5deg) translateY(-3px)';
            setTimeout(function() { ov.style.transform = ''; setTimeout(function() { ov.style.transition = ''; }, 180); }, 110);
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(0,180,216,0.18);transition:opacity 0.22s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 250); }, 30);
        };
        window.snonuxPageEffect = function() {
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-zoom'); setTimeout(function() { ov.classList.remove('sno-fx-zoom'); }, 330); }
        };
    })();
