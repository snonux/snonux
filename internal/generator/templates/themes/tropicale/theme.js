
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,t0=performance.now(),i;
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(50,1,0.1,80);ca.position.set(0,1.5,10);
        // Golden sun sphere
        var sun=new THREE.Mesh(new THREE.SphereGeometry(1.4,18,18),new THREE.MeshBasicMaterial({color:0xfbbf24,transparent:true,opacity:0.9}));
        sun.position.set(3,2.5,-8);sc.add(sun);
        // Soft corona halo around sun
        var halo=new THREE.Mesh(new THREE.SphereGeometry(2.2,16,16),new THREE.MeshBasicMaterial({color:0xf97316,transparent:true,opacity:0.25,blending:THREE.AdditiveBlending,depthWrite:false}));
        halo.position.copy(sun.position);sc.add(halo);
        // Sandy beach: wide flat plane at the bottom of the scene
        var beach=new THREE.Mesh(new THREE.PlaneGeometry(26,7),new THREE.MeshBasicMaterial({color:0xe8c97a,transparent:true,opacity:0.88,side:THREE.DoubleSide}));
        beach.rotation.x=-Math.PI/2;beach.position.set(0,-1.5,-2);sc.add(beach);
        // Shallow water strip at the shore edge — lighter turquoise
        var shore=new THREE.Mesh(new THREE.PlaneGeometry(26,1.4),new THREE.MeshBasicMaterial({color:0x38c9d8,transparent:true,opacity:0.55,side:THREE.DoubleSide}));
        shore.rotation.x=-Math.PI/2;shore.position.set(0,-1.48,-5.5);sc.add(shore);
        // Palm trunk: leaning cylinder
        var trunk=new THREE.Mesh(new THREE.CylinderGeometry(0.09,0.16,4.5,8),new THREE.MeshBasicMaterial({color:0x6b4226}));
        trunk.position.set(-4.2,-0.2,-3);trunk.rotation.z=0.2;sc.add(trunk);
        // Palm fronds: five half-ellipses fanning from the crown
        var frondMat=new THREE.MeshBasicMaterial({color:0x2d8a4e,transparent:true,opacity:0.92,side:THREE.DoubleSide});
        var frondAngles=[0,1.15,2.3,3.55,4.8];
        frondAngles.forEach(function(a){
            var f=new THREE.Mesh(new THREE.PlaneGeometry(2.6,0.55,6,1),frondMat);
            f.position.set(-4.2+Math.cos(a)*1.35,2.1+Math.sin(a)*0.5,-3+Math.sin(a)*1.0);
            f.rotation.set(-0.1,a,Math.PI*0.06);sc.add(f);
        });
        // Seagulls — simple V arcs made of thin tori
        for(i=0;i<5;i++){
            var b=new THREE.Mesh(new THREE.TorusGeometry(0.2+Math.random()*0.1,0.03,6,10,Math.PI),new THREE.MeshBasicMaterial({color:0xfef9e7,transparent:true,opacity:0.8}));
            b.position.set((Math.random()-0.5)*7,1.2+Math.random()*2.5,(Math.random()-0.5)*3-3);
            b.userData.sx=(Math.random()-0.5)*0.011;b.userData.y0=b.position.y;b.userData.phase=Math.random()*Math.PI*2;
            sc.add(b);
        }
        sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);var t=(now-t0)*0.001;
            var birds=sc.children.filter(function(c){return c.userData.sx!==undefined;});
            birds.forEach(function(b){b.position.x+=b.userData.sx;b.position.y=b.userData.y0+Math.sin(t*1.4+b.userData.phase)*0.12;
                if(b.position.x>7)b.position.x=-7;if(b.position.x<-7)b.position.x=7;});
            // Shore shimmer: opacity pulses like sunlight on shallow water
            shore.material.opacity=0.45+Math.sin(t*1.8)*0.12;
            sun.scale.setScalar(1+Math.sin(t*0.8)*0.04);halo.scale.setScalar(1+Math.sin(t*0.6)*0.07);
            ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();


    // Tropicale WebGL: tropical sunset beach — rolling ocean waves, a glowing sun
    // sinking toward the horizon, drifting seagulls, a palm tree silhouette on
    // the shore, and golden sparkle particles on the water surface.
    // Wild mode: cyclone swells with crashing foam and a darkened storm sky.
    (function() {
        var _wild = false, _snoTOffset = 0, _snoLastT = 0;
        var scene, camera, renderer, clock;
        var waveGeo, waveMesh, sunMesh, sunHalo, skyLight, sunLight;
        var seagulls = [];
        var SPARKLE_COUNT = 500;
        var sparklePos, sparkleGeo;

        function buildWaves() {
            // High-density plane for smooth shoreline swells
            waveGeo = new THREE.PlaneGeometry(280, 180, 90, 60);
            waveMesh = new THREE.Mesh(waveGeo, new THREE.MeshPhongMaterial({
                color: 0x0e7490, emissive: 0x0a3a4a, emissiveIntensity: 0.3,
                transparent: true, opacity: 0.9, side: THREE.DoubleSide, shininess: 120
            }));
            waveMesh.rotation.x = -Math.PI / 2;
            waveMesh.position.y = 0;
            scene.add(waveMesh);
        }

        function buildSun() {
            // Large warm sphere near the horizon simulating the setting sun
            sunMesh = new THREE.Mesh(
                new THREE.SphereGeometry(7, 24, 24),
                new THREE.MeshBasicMaterial({ color: 0xfbbf24, transparent: true, opacity: 0.95 })
            );
            sunMesh.position.set(30, 10, -70);
            scene.add(sunMesh);

            // Wide soft halo with additive blending for the glow corona
            sunHalo = new THREE.Mesh(
                new THREE.SphereGeometry(13, 16, 16),
                new THREE.MeshBasicMaterial({
                    color: 0xf97316, transparent: true, opacity: 0.28,
                    blending: THREE.AdditiveBlending, depthWrite: false
                })
            );
            sunHalo.position.copy(sunMesh.position);
            scene.add(sunHalo);
        }

        function buildPalm() {
            // Trunk: tall narrow cylinder leaning right
            var trunkGeo = new THREE.CylinderGeometry(0.3, 0.6, 18, 8);
            var trunkMat = new THREE.MeshPhongMaterial({ color: 0x6b4226, emissive: 0x2a1a0a, emissiveIntensity: 0.2 });
            var trunk = new THREE.Mesh(trunkGeo, trunkMat);
            trunk.position.set(-28, 5, 10);
            trunk.rotation.z = 0.18; // lean toward the sea
            scene.add(trunk);

            // Fronds: five flat ellipses radiating from the crown
            var frondMat = new THREE.MeshPhongMaterial({ color: 0x2d8a4e, emissive: 0x0a2a18, emissiveIntensity: 0.2, side: THREE.DoubleSide });
            var frondAngles = [0, 1.2, 2.5, 3.9, 5.2];
            frondAngles.forEach(function(angle) {
                var frond = new THREE.Mesh(new THREE.PlaneGeometry(7, 1.4, 8, 1), frondMat);
                frond.position.set(-28 + Math.cos(angle) * 4.5, 14.5 + Math.sin(angle * 0.5) * 1.2, 10 + Math.sin(angle) * 3);
                frond.rotation.set(0.15, angle, Math.PI * 0.08);
                scene.add(frond);
            });
        }

        function buildSeagulls() {
            // V-shaped torus arcs flying above the horizon
            var positions = [[-18,14,-30],[8,16,-24],[-5,18,-40],[22,13,-35],[-30,12,-20],[12,20,-28]];
            positions.forEach(function(p) {
                var body = new THREE.Mesh(
                    new THREE.TorusGeometry(0.8, 0.12, 6, 14, Math.PI),
                    new THREE.MeshBasicMaterial({ color: 0xfef9e7, transparent: true, opacity: 0.85 })
                );
                body.position.set(p[0], p[1], p[2]);
                body.userData.baseX = p[0];
                body.userData.baseY = p[1];
                body.userData.phase = Math.random() * Math.PI * 2;
                body.userData.speed = 0.012 + Math.random() * 0.01;
                seagulls.push(body);
                scene.add(body);
            });
        }

        function buildSparkles() {
            // Golden sunlight glitter on the wave surface
            sparklePos = new Float32Array(SPARKLE_COUNT * 3);
            for (var i = 0; i < SPARKLE_COUNT; i++) {
                sparklePos[i*3]   = (Math.random() - 0.5) * 200;
                sparklePos[i*3+1] = 0.5 + Math.random() * 0.5;
                sparklePos[i*3+2] = (Math.random() - 0.5) * 100 - 20;
            }
            sparkleGeo = new THREE.BufferGeometry();
            sparkleGeo.setAttribute('position', new THREE.BufferAttribute(sparklePos, 3));
            scene.add(new THREE.Points(sparkleGeo, new THREE.PointsMaterial({
                color: 0xfbbf24, size: 0.22, transparent: true, opacity: 0.65
            })));
        }

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x0a1e2e);
            scene.fog = new THREE.FogExp2(0x0a1e2e, 0.006);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 250);
            camera.position.set(0, 22, 60);
            camera.lookAt(0, 2, 0);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();

            // Warm ambient light imitating diffuse tropical sky
            scene.add(new THREE.AmbientLight(0x1a4a5a, 0.6));
            // Sunlight: warm golden point light from the horizon direction
            sunLight = new THREE.PointLight(0xfbbf24, 2.8, 200);
            sunLight.position.set(30, 20, -60);
            scene.add(sunLight);
            // Soft fill from below the horizon (scattered sea light)
            var seaFill = new THREE.PointLight(0x38c9d8, 1.2, 80);
            seaFill.position.set(0, -5, 0);
            scene.add(seaFill);

            buildWaves();
            buildSun();
            buildPalm();
            buildSeagulls();
            buildSparkles();

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
            // Wild mode accelerates time 7× to simulate a tropical storm
            _snoTOffset += (realT - _snoLastT) * (_wild ? 7 : 0);
            _snoLastT = realT;
            var t = realT + _snoTOffset;
            var waveAmp = _wild ? 3.2 : 1;

            // Rolling shoreline swells with three overlapping wave components
            var pos = waveGeo.attributes.position;
            for (var i = 0; i < pos.count; i++) {
                var x = pos.getX(i), z = pos.getZ(i);
                pos.setY(i,
                    Math.sin(x * 0.035 + t * 1.0) * 2.8 * waveAmp +
                    Math.cos(z * 0.05  + t * 0.75) * 2.0 * waveAmp +
                    Math.sin((x - z)   * 0.022 + t * 0.5) * 1.2 * waveAmp
                );
            }
            pos.needsUpdate = true;
            waveGeo.computeVertexNormals();

            // Sun gently bobs toward the horizon and pulses its halo
            sunMesh.position.y = 10 + Math.sin(t * 0.12) * 1.5;
            sunHalo.position.y = sunMesh.position.y;
            sunHalo.scale.setScalar(1 + Math.sin(t * 0.45) * 0.06);
            // In wild mode the sky darkens as the storm rolls in
            sunLight.intensity = _wild ? 1.0 + Math.sin(t * 2) * 0.5 : 2.8;

            // Seagulls drift lazily across the horizon, dipping on thermals
            seagulls.forEach(function(b) {
                b.position.x += b.userData.speed * (_wild ? 3.5 : 1);
                b.position.y = b.userData.baseY + Math.sin(t * 1.1 + b.userData.phase) * 0.8;
                b.rotation.z = Math.sin(t * 0.7 + b.userData.phase) * 0.15;
                if (b.position.x > 50) b.position.x = -50;
            });

            // Sparkles ride the wave surface, flickering in sunlight
            var sp = sparkleGeo.attributes.position;
            var sMult = _wild ? 0 : 1; // hide sparkles during storm
            for (var si = 0; si < SPARKLE_COUNT; si++) {
                var sx = sparklePos[si*3], sz = sparklePos[si*3+2];
                sparklePos[si*3+1] = (0.5 + Math.random() * 0.3) * sMult +
                    Math.sin(sx * 0.035 + t * 1.0) * 2.8 * waveAmp +
                    Math.cos(sz * 0.05  + t * 0.75) * 2.0 * waveAmp;
            }
            sp.needsUpdate = true;

            renderer.render(scene, camera);
        }

        initThree();

        // --- Tropical audio helpers ---

        // Monkey call: three rising "oo-oo" chirps with vibrato, then a falling "aah" screech.
        // Registered as window.snonuxSplashSound so navscript calls it instead of the default chime.
        window.snonuxSplashSound = function(ctx) {
            var now = ctx.currentTime;
            // Three staccato rising chirps — "oo oo oo"
            [[380,820,0.00],[480,980,0.19],[600,1180,0.39]].forEach(function(p) {
                var osc = ctx.createOscillator();
                var lfo = ctx.createOscillator(); // vibrato LFO
                var lfoGain = ctx.createGain();
                var g = ctx.createGain();
                lfo.connect(lfoGain); lfoGain.connect(osc.frequency);
                osc.connect(g); g.connect(ctx.destination);
                osc.type = 'sine';
                lfo.type = 'sine';
                lfo.frequency.value = 14; // fast wobble for monkey texture
                lfoGain.gain.value = 55;  // ±55 Hz wobble depth
                var t = now + p[2];
                osc.frequency.setValueAtTime(p[0], t);
                osc.frequency.linearRampToValueAtTime(p[1], t + 0.11);
                g.gain.setValueAtTime(0.09, t);
                g.gain.exponentialRampToValueAtTime(0.001, t + 0.17);
                lfo.start(t); lfo.stop(t + 0.18);
                osc.start(t); osc.stop(t + 0.18);
            });
            // Long falling "aah" screech with heavier vibrato
            var osc2 = ctx.createOscillator();
            var lfo2 = ctx.createOscillator();
            var lg2  = ctx.createGain();
            var g2   = ctx.createGain();
            lfo2.connect(lg2); lg2.connect(osc2.frequency);
            osc2.connect(g2); g2.connect(ctx.destination);
            osc2.type = 'sine';
            lfo2.type = 'sine';
            lfo2.frequency.value = 10;
            lg2.gain.value = 80;
            var t2 = now + 0.62;
            osc2.frequency.setValueAtTime(1050, t2);
            osc2.frequency.linearRampToValueAtTime(420, t2 + 0.32);
            g2.gain.setValueAtTime(0.10, t2);
            g2.gain.exponentialRampToValueAtTime(0.001, t2 + 0.36);
            lfo2.start(t2); lfo2.stop(t2 + 0.38);
            osc2.start(t2); osc2.stop(t2 + 0.38);
            return true; // skip default chime
        };

        // Seagull chirp: rapid sine glide 1600→2400→1900 Hz — recognisable bird cry.
        function snoPlayBirdChirp(delay) {
            try {
                var ctx = new (window.AudioContext || window.webkitAudioContext)();
                var osc = ctx.createOscillator();
                var g   = ctx.createGain();
                osc.connect(g); g.connect(ctx.destination);
                osc.type = 'sine';
                var t = ctx.currentTime + (delay || 0);
                osc.frequency.setValueAtTime(1600, t);
                osc.frequency.linearRampToValueAtTime(2400, t + 0.06);
                osc.frequency.linearRampToValueAtTime(1900, t + 0.14);
                g.gain.setValueAtTime(0, t);
                g.gain.linearRampToValueAtTime(0.07, t + 0.02);
                g.gain.exponentialRampToValueAtTime(0.001, t + 0.22);
                osc.start(t); osc.stop(t + 0.24);
            } catch(_) {}
        }

        // Wave crash: shaped white-noise burst with a low-pass sweep — shore sound.
        function snoPlayWaveCrash(gainMult) {
            try {
                var ctx = new (window.AudioContext || window.webkitAudioContext)();
                var dur = 0.42;
                var buf = ctx.createBuffer(1, Math.ceil(ctx.sampleRate * dur), ctx.sampleRate);
                var data = buf.getChannelData(0);
                for (var i = 0; i < data.length; i++) data[i] = Math.random() * 2 - 1;
                var src = ctx.createBufferSource();
                src.buffer = buf;
                var flt = ctx.createBiquadFilter();
                flt.type = 'lowpass';
                flt.frequency.setValueAtTime(1200, ctx.currentTime);
                flt.frequency.exponentialRampToValueAtTime(180, ctx.currentTime + dur);
                var gn = ctx.createGain();
                var gv = (gainMult || 1) * 0.18;
                gn.gain.setValueAtTime(0, ctx.currentTime);
                gn.gain.linearRampToValueAtTime(gv, ctx.currentTime + 0.04);
                gn.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + dur);
                src.connect(flt); flt.connect(gn); gn.connect(ctx.destination);
                src.start(); src.stop(ctx.currentTime + dur + 0.02);
            } catch(_) {}
        }

        // Gentle wave lap: softer, shorter version of crash for scroll feedback.
        function snoPlayWaveLap() { snoPlayWaveCrash(0.38); }

        // --- Tropical nav / wild effects ---

        // Open post: foam spray radiates outward + wave crash sound.
        window.snonuxOpenEffect = function(post) {
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-fly'); setTimeout(function() { modal.classList.remove('sno-modal-fly'); }, 390); }
            snoPlayWaveCrash(1);
            // Sandy foam dots scatter outward from the post card
            var r = post ? post.getBoundingClientRect() : {left: window.innerWidth/2, top: window.innerHeight/2, width: 0, height: 0};
            for (var i = 0; i < 12; i++) {
                (function(i) {
                    var b  = document.createElement('div');
                    var sz = 4 + Math.random() * 7;
                    var dx = (Math.random() - 0.5) * 180;
                    var dy = -(40 + Math.random() * 80);
                    // Alternate between sandy foam and turquoise spray
                    var col = (i % 3 === 0) ? 'rgba(232,201,122,0.75)' : 'rgba(56,201,216,0.5)';
                    b.style.cssText = 'position:fixed;top:' + (r.top + r.height * 0.6) + 'px;left:' + (r.left + r.width * 0.5) + 'px;' +
                        'z-index:997;pointer-events:none;width:' + sz + 'px;height:' + sz + 'px;border-radius:50%;' +
                        'background:' + col + ';transition:transform 0.5s ease,opacity 0.45s;';
                    document.body.appendChild(b);
                    setTimeout(function() {
                        b.style.transform = 'translate(' + dx + 'px,' + dy + 'px) scale(0.4)';
                        b.style.opacity = '0';
                        setTimeout(function() { b.remove(); }, 520);
                    }, 18 + i * 22);
                })(i);
            }
        };

        // Close post: tide-wash — teal shimmer that retreats downward + brief chirp.
        window.snonuxCloseEffect = function() {
            snoPlayBirdChirp(0);
            var d = document.createElement('div');
            // Gradient simulates the thin sheen of receding water on wet sand
            d.style.cssText = 'position:fixed;bottom:0;left:0;right:0;height:38px;z-index:998;pointer-events:none;' +
                'background:linear-gradient(180deg,transparent,rgba(56,201,216,0.35));transition:transform 0.38s ease,opacity 0.32s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.transform = 'translateY(40px)'; d.style.opacity = '0'; setTimeout(function() { d.remove(); }, 420); }, 20);
        };

        // Scroll: tide-mark bar sweeps across the viewport + gentle wave lap.
        window.snonuxScrollEffect = function(dir) {
            snoPlayWaveLap();
            var isDown = dir === 'down';
            var thick = _wild ? '18px' : '6px';
            var d = document.createElement('div');
            // Sandy-gold centre fading to lagoon teal at the edges — tide-mark stripe
            d.style.cssText = 'position:fixed;left:0;right:0;height:' + thick + ';z-index:9000;pointer-events:none;' +
                'background:linear-gradient(90deg,transparent,rgba(14,116,144,0.85),rgba(232,201,122,0.9),rgba(14,116,144,0.85),transparent);' +
                (isDown ? 'top:0;' : 'bottom:0;') +
                'transition:transform 0.38s ease,opacity 0.38s ease;';
            document.body.appendChild(d);
            setTimeout(function() { d.style.transform = isDown ? 'translateY(100vh)' : 'translateY(-100vh)'; d.style.opacity = '0'; }, 16);
            setTimeout(function() { d.remove(); }, 440);
        };

        window.snonuxWildToggle = function() {
            _wild = !_wild;
            var b = document.getElementById('sno-wild-badge');
            if (b) b.classList.toggle('sno-wild-on', _wild);
        };

        // Nav (j/k): post bobs gently like driftwood + seagull chirp.
        window.snonuxNavEffect = function() {
            snoPlayBirdChirp(0);
            var ov = document.querySelector('.overlay');
            if (!ov) return;
            // Soft vertical bob — no skew — differentiates from ocean/other themes
            ov.style.transition = 'transform 0.14s ease-out';
            ov.style.transform = 'translateY(-5px)';
            setTimeout(function() {
                ov.style.transform = 'translateY(2px)';
                setTimeout(function() { ov.style.transform = ''; ov.style.transition = ''; }, 120);
            }, 140);
            // Warm golden-sun shimmer flash — contrasts with ocean's cold blue flash
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;' +
                'background:radial-gradient(ellipse 70% 55% at 72% 28%,rgba(251,191,36,0.14) 0%,transparent 70%);transition:opacity 0.28s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity = '0'; setTimeout(function() { d.remove(); }, 300); }, 30);
        };

        // Page nav (h/l): full wave-crash surge across the screen.
        window.snonuxPageEffect = function() {
            snoPlayWaveCrash(0.7);
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-zoom'); setTimeout(function() { ov.classList.remove('sno-fx-zoom'); }, 330); }
        };
    })();
