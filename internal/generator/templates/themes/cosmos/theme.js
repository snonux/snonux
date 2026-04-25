
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,g=new THREE.Group(),t0=performance.now();
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(48,1,0.1,90);ca.position.set(0,0.5,10);
        var planet=new THREE.Mesh(new THREE.SphereGeometry(1.35,24,24),new THREE.MeshBasicMaterial({color:0xc8853a,transparent:true,opacity:0.92}));
        var ring=new THREE.Mesh(new THREE.TorusGeometry(2.1,0.05,8,64),new THREE.MeshBasicMaterial({color:0xffd166,transparent:true,opacity:0.85}));
        ring.rotation.x=Math.PI/2.25;var moon=new THREE.Mesh(new THREE.SphereGeometry(0.35,12,12),new THREE.MeshBasicMaterial({color:0x9b5de5,transparent:true,opacity:0.8}));
        moon.position.set(2.8,0.6,0.5);g.add(planet);g.add(ring);g.add(moon);sc.add(g);sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);var t=(now-t0)*0.001;g.rotation.y=t*0.22;planet.rotation.y=t*0.35;
            moon.position.x=2.6*Math.cos(t*0.7);moon.position.z=2.6*Math.sin(t*0.7);ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();


    // Cosmos WebGL: ringed planet, swirling nebula blobs, asteroid belt, and stars.
    // The planet sits at lower-right and slowly rotates; asteroids orbit it;
    // nebula clouds drift with additive blending for a deep-space glow.
    (function() {
        var _wild = false;
        var scene, camera, renderer, clock;
        var planet, planetRings = [];
        var asteroids = [];
        var ASTEROID_COUNT = 300;
        var asteroidAngles, asteroidRadii, asteroidSpeeds, asteroidY;

        function buildPlanet() {
            // Planet body — warm golden tone
            planet = new THREE.Mesh(
                new THREE.SphereGeometry(14, 48, 48),
                new THREE.MeshPhongMaterial({
                    color: 0xc8853a, emissive: 0x3a1800, emissiveIntensity: 0.4, shininess: 60
                })
            );
            planet.position.set(28, -18, -55);
            scene.add(planet);

            // Ring system — 5 tilted torus rings in gold/purple
            var ringCols = [0xffd166, 0xc07c30, 0xffd166, 0x9b5de5, 0xffd166];
            for (var i = 0; i < 5; i++) {
                var ring = new THREE.Mesh(
                    new THREE.TorusGeometry(18 + i * 2.5, 0.5 - i * 0.06, 8, 128),
                    new THREE.MeshBasicMaterial({ color: ringCols[i], transparent: true, opacity: 0.55 - i * 0.08, side: THREE.DoubleSide })
                );
                ring.position.copy(planet.position);
                ring.rotation.x = Math.PI / 2.4;
                ring.rotation.z = 0.2;
                scene.add(ring);
                planetRings.push(ring);
            }
        }

        function buildNebula() {
            // Large translucent additive blobs for the nebula cloud
            var nCols = [0x9b5de5, 0x4cc9f0, 0x7b2fff, 0x4cc9f0, 0x9b5de5];
            var nPos  = [[-30,20,-80],[-10,-10,-90],[20,30,-70],[-20,-25,-95],[10,15,-85]];
            nCols.forEach(function(c, i) {
                var mesh = new THREE.Mesh(
                    new THREE.SphereGeometry(22 + i * 4, 16, 16),
                    new THREE.MeshBasicMaterial({
                        color: c, transparent: true, opacity: 0.09,
                        blending: THREE.AdditiveBlending, depthWrite: false
                    })
                );
                mesh.position.set(nPos[i][0], nPos[i][1], nPos[i][2]);
                scene.add(mesh);
            });
        }

        function buildAsteroids() {
            asteroidAngles = new Float32Array(ASTEROID_COUNT);
            asteroidRadii  = new Float32Array(ASTEROID_COUNT);
            asteroidSpeeds = new Float32Array(ASTEROID_COUNT);
            asteroidY      = new Float32Array(ASTEROID_COUNT);

            var geo = new THREE.BufferGeometry();
            var pos = new Float32Array(ASTEROID_COUNT * 3);
            for (var i = 0; i < ASTEROID_COUNT; i++) {
                asteroidAngles[i] = Math.random() * Math.PI * 2;
                asteroidRadii[i]  = 20 + Math.random() * 12;
                asteroidSpeeds[i] = 0.003 + Math.random() * 0.004;
                asteroidY[i]      = (Math.random() - 0.5) * 3;
                pos[i*3] = pos[i*3+1] = pos[i*3+2] = 0;
            }
            geo.setAttribute('position', new THREE.BufferAttribute(pos, 3));
            asteroids = new THREE.Points(geo, new THREE.PointsMaterial({
                color: 0xaaaaaa, size: 0.3, transparent: true, opacity: 0.7
            }));
            // Asteroids orbit the planet — they live in planet-relative coords
            planet.add(asteroids);
        }

        function buildStars() {
            var pos = new Float32Array(2500 * 3);
            for (var i = 0; i < 2500 * 3; i += 3) {
                var r = 100 + Math.random() * 80, t = Math.random() * Math.PI * 2, p = Math.acos(2 * Math.random() - 1);
                pos[i]   = r * Math.sin(p) * Math.cos(t);
                pos[i+1] = r * Math.sin(p) * Math.sin(t);
                pos[i+2] = r * Math.cos(p);
            }
            var geo = new THREE.BufferGeometry();
            geo.setAttribute('position', new THREE.BufferAttribute(pos, 3));
            scene.add(new THREE.Points(geo, new THREE.PointsMaterial({ color: 0xffffff, size: 0.18, transparent: true, opacity: 0.85 })));
        }

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x020214);
            scene.fog = new THREE.Fog(0x020214, 80, 200);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 300);
            camera.position.set(0, 6, 38);
            camera.lookAt(0, 0, 0);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();

            scene.add(new THREE.AmbientLight(0x4cc9f0, 0.4));
            var sun = new THREE.PointLight(0xffd166, 3, 300);
            sun.position.set(-60, 40, 30);
            scene.add(sun);

            buildPlanet();
            buildNebula();
            buildAsteroids();
            buildStars();

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
            var t = clock.getElapsedTime();

            var sm = _wild ? 22 : 1;
            planet.rotation.y += 0.0015 * sm;

            // Planet rings tilt and counter-spin wildly
            for (var ri = 0; ri < planetRings.length; ri++) {
                planetRings[ri].rotation.y += 0.002 * sm;
                planetRings[ri].rotation.z += 0.001 * sm * (ri % 2 === 0 ? 1 : -1);
            }

            // Asteroid belt: warp-speed orbit + radial breathing + vertical scatter in wild
            var pos = asteroids.geometry.attributes.position;
            for (var i = 0; i < ASTEROID_COUNT; i++) {
                asteroidAngles[i] += asteroidSpeeds[i] * sm;
                var a = asteroidAngles[i];
                var r = asteroidRadii[i] + (_wild ? Math.sin(t * 2.1 + i * 0.3) * 9 : 0);
                var y = asteroidY[i]      + (_wild ? Math.sin(t * 3.0 + i * 0.7) * 2.5 : 0);
                pos.setXYZ(i, Math.cos(a) * r, y, Math.sin(a) * r);
            }
            pos.needsUpdate = true;

            // Camera: chaotic deep-space flight in wild, gentle orbit otherwise
            if (_wild) {
                camera.position.x = Math.sin(t * 0.38) * 20;
                camera.position.y = 6 + Math.sin(t * 0.29) * 10;
                camera.position.z = 38 + Math.sin(t * 0.21) * 14;
                camera.fov = 60 + Math.sin(t * 0.47) * 18;
            } else {
                camera.position.x = Math.sin(t * 0.06) * 6;
                camera.position.y = 6 + Math.sin(t * 0.04) * 2;
                camera.position.z = 38;
                camera.fov = 60;
            }
            camera.updateProjectionMatrix();

            renderer.render(scene, camera);
        }

        initThree();

        // Cosmos nav/wild effects — shooting star on navigate, warp speed on wild
        window.snonuxOpenEffect = function(post) {
            // Materialise from deep space — zoom + gold shimmer + orbiting ring
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-zoom'); setTimeout(function() { modal.classList.remove('sno-modal-zoom'); }, 400); }
            var r = post ? post.getBoundingClientRect() : {left: window.innerWidth/2, top: window.innerHeight/2, width: 0, height: 0};
            var ring = document.createElement('div');
            ring.style.cssText = 'position:fixed;top:' + (r.top+r.height/2-5) + 'px;left:' + (r.left+r.width/2-5) + 'px;z-index:997;pointer-events:none;width:10px;height:10px;border-radius:50%;border:2px solid rgba(255,209,102,0.85);box-shadow:0 0 12px rgba(155,93,229,0.5);transition:all 0.4s ease,opacity 0.4s';
            document.body.appendChild(ring);
            setTimeout(function() { ring.style.transform='scale(32)'; ring.style.opacity='0'; setTimeout(function() { ring.remove(); }, 440); }, 15);
        };
        window.snonuxCloseEffect = function() {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:radial-gradient(ellipse at center,rgba(155,93,229,0.12) 0%,transparent 65%);transition:opacity 0.25s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 280); }, 15);
        };
        window.snonuxScrollEffect = function(dir) {
            var isDown = dir === 'down';
            var thick = _wild ? '14px' : '5px';
            var d = document.createElement('div');
            // Cosmos: purple-to-gold stardust streak
            d.style.cssText = 'position:fixed;left:0;right:0;height:' + thick + ';z-index:9000;pointer-events:none;' +
                'background:linear-gradient(90deg,transparent,rgba(155,93,229,0.9),rgba(255,200,0,0.9),rgba(155,93,229,0.9),transparent);' +
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
            // Shooting star CSS line flashes across screen
            var d = document.createElement('div');
            var angle = 30 + Math.random() * 30;
            d.style.cssText = 'position:fixed;top:' + (10+Math.random()*40) + '%;left:-10%;z-index:998;pointer-events:none;width:35%;height:2px;background:linear-gradient(90deg,transparent,rgba(255,209,102,0.9),transparent);transform:rotate(-' + angle + 'deg);transition:left 0.28s ease,opacity 0.28s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.left='110%'; d.style.opacity='0'; setTimeout(function() { d.remove(); }, 320); }, 20);
        };
        window.snonuxPageEffect = function() {
            // Warp flash — stars streak white
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:radial-gradient(ellipse at center,rgba(255,255,255,0.25) 0%,transparent 70%);transition:opacity 0.2s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 230); }, 20);
        };
    })();
