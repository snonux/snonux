
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,sparks=[],t0=performance.now(),N=120;
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:false,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(1);
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(50,1,0.1,80);ca.position.z=18;
        var geo=new THREE.SphereGeometry(0.12,4,4);
        for(var i=0;i<N;i++){
            var c=Math.random()>0.5?0xff0000:0xffd700;
            var mat=new THREE.MeshBasicMaterial({color:c,transparent:true,opacity:0.4+Math.random()*0.5});
            var m=new THREE.Mesh(geo,mat);
            m.position.set((Math.random()-0.5)*24,(Math.random()-0.5)*16,(Math.random()-0.5)*6);
            m.userData.vx=(Math.random()-0.5)*0.08;
            m.userData.vy=0.03+Math.random()*0.06;
            m.userData.vz=(Math.random()-0.5)*0.04;
            sc.add(m);sparks.push(m);
        }
        sz();window.addEventListener('resize',sz);
        function loop(){raf=requestAnimationFrame(loop);
            for(var i=0;i<sparks.length;i++){
                var s=sparks[i];
                s.position.x+=s.userData.vx;s.position.y+=s.userData.vy;s.position.z+=s.userData.vz;
                if(s.position.y>10){s.position.y=-10;s.position.x=(Math.random()-0.5)*24;}
            }
            ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();


    (function() {
        var _wild = false, _snoTOffset = 0, _snoLastT = 0;
        var scene, camera, renderer, clock;
        var N_SPARKS = 2000;
        var N_DEBRIS = 400;
        var sparkPoints, debrisPoints;
        var spArr, spColArr, spVX, spVY, spVZ, spLife, spMaxLife;
        var dbArr, dbVX, dbVY, dbVZ, dbLife, dbMaxLife;

        function resetSpark(i) {
            var cx = (Math.random() - 0.5) * 80;
            var cy = -25 + Math.random() * 5;
            var cz = (Math.random() - 0.5) * 40 - 10;
            spArr[i*3] = cx; spArr[i*3+1] = cy; spArr[i*3+2] = cz;
            spVX[i] = (Math.random() - 0.5) * 0.15;
            spVY[i] = 0.08 + Math.random() * 0.18;
            spVZ[i] = (Math.random() - 0.5) * 0.08;
            spMaxLife[i] = 0.5 + Math.random() * 0.8;
            spLife[i] = Math.random();
        }

        function resetDebris(i) {
            dbArr[i*3] = (Math.random() - 0.5) * 60;
            dbArr[i*3+1] = 30 + Math.random() * 10;
            dbArr[i*3+2] = (Math.random() - 0.5) * 30 - 10;
            dbVX[i] = (Math.random() - 0.5) * 0.06;
            dbVY[i] = -(0.04 + Math.random() * 0.1);
            dbVZ[i] = (Math.random() - 0.5) * 0.04;
            dbMaxLife[i] = 1.5 + Math.random() * 2;
            dbLife[i] = Math.random();
        }

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x0a0a0a);
            scene.fog = new THREE.Fog(0x0a0a0a, 40, 110);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 150);
            camera.position.set(0, 5, 45);
            camera.lookAt(0, -5, 0);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: false });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();

            scene.add(new THREE.AmbientLight(0x220000, 0.8));
            var redLight = new THREE.PointLight(0xff0000, 3, 80);
            redLight.position.set(0, -10, 0);
            scene.add(redLight);
            var goldLight = new THREE.PointLight(0xffd700, 2, 60);
            goldLight.position.set(10, 15, 10);
            scene.add(goldLight);

            // Ground plane — dark concrete
            var ground = new THREE.Mesh(
                new THREE.PlaneGeometry(200, 200),
                new THREE.MeshPhongMaterial({ color: 0x1a1a1a, emissive: 0x220000, emissiveIntensity: 0.3 })
            );
            ground.rotation.x = -Math.PI / 2;
            ground.position.y = -25;
            scene.add(ground);

            // Rubble blocks
            var rubbleData = [
                [-20,-22,-18,4], [15,-21,-25,5], [-10,-23,-30,3],
                [28,-22,-15,4], [-25,-21,-22,6], [5,-23,-35,3]
            ];
            rubbleData.forEach(function(b) {
                var mesh = new THREE.Mesh(
                    new THREE.BoxGeometry(b[3], b[3]*0.6, b[3]*0.8),
                    new THREE.MeshPhongMaterial({ color: 0x2a2a2a, emissive: 0x330000, emissiveIntensity: 0.2 })
                );
                mesh.position.set(b[0], b[1], b[2]);
                mesh.rotation.set(Math.random()*0.5, Math.random(), Math.random()*0.3);
                scene.add(mesh);
            });

            // Explosion glow sphere
            var glow = new THREE.Mesh(
                new THREE.SphereGeometry(35, 12, 12),
                new THREE.MeshBasicMaterial({ color: 0xff2200, transparent: true, opacity: 0.12, blending: THREE.AdditiveBlending, depthWrite: false })
            );
            glow.position.set(0, -40, -20);
            scene.add(glow);

            // Rising sparks
            spArr = new Float32Array(N_SPARKS * 3);
            spColArr = new Float32Array(N_SPARKS * 3);
            spVX = new Float32Array(N_SPARKS); spVY = new Float32Array(N_SPARKS); spVZ = new Float32Array(N_SPARKS);
            spLife = new Float32Array(N_SPARKS); spMaxLife = new Float32Array(N_SPARKS);
            for (var i = 0; i < N_SPARKS; i++) resetSpark(i);
            var spGeo = new THREE.BufferGeometry();
            spGeo.setAttribute('position', new THREE.BufferAttribute(spArr, 3));
            spGeo.setAttribute('color', new THREE.BufferAttribute(spColArr, 3));
            sparkPoints = new THREE.Points(spGeo, new THREE.PointsMaterial({
                size: 0.25, vertexColors: true, transparent: true, opacity: 0.9,
                blending: THREE.AdditiveBlending, depthWrite: false
            }));
            scene.add(sparkPoints);

            // Falling debris
            dbArr = new Float32Array(N_DEBRIS * 3);
            dbVX = new Float32Array(N_DEBRIS); dbVY = new Float32Array(N_DEBRIS); dbVZ = new Float32Array(N_DEBRIS);
            dbLife = new Float32Array(N_DEBRIS); dbMaxLife = new Float32Array(N_DEBRIS);
            for (var j = 0; j < N_DEBRIS; j++) resetDebris(j);
            var dbGeo = new THREE.BufferGeometry();
            dbGeo.setAttribute('position', new THREE.BufferAttribute(dbArr, 3));
            debrisPoints = new THREE.Points(dbGeo, new THREE.PointsMaterial({
                color: 0x888888, size: 0.6, transparent: true, opacity: 0.4, depthWrite: false
            }));
            scene.add(debrisPoints);

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
            var dt = clock.getDelta();
            var realT = clock.getElapsedTime();
            _snoTOffset += (realT - _snoLastT) * (_wild ? 6 : 0);
            _snoLastT = realT;

            var sparkMult = _wild ? 6 : 1;
            for (var i = 0; i < N_SPARKS; i++) {
                spLife[i] += dt / spMaxLife[i] * sparkMult;
                if (spLife[i] > 1.0) resetSpark(i);
                var idx = i * 3;
                spArr[idx] += spVX[i] * sparkMult;
                spArr[idx+1] += spVY[i] * sparkMult;
                spArr[idx+2] += spVZ[i] * sparkMult;
                var fade = Math.max(0, 1 - spLife[i] * 1.2);
                // Red to gold gradient
                spColArr[idx] = fade;
                spColArr[idx+1] = fade * Math.max(0, 0.85 - spLife[i]);
                spColArr[idx+2] = 0;
            }
            sparkPoints.geometry.attributes.position.needsUpdate = true;
            sparkPoints.geometry.attributes.color.needsUpdate = true;

            var debrisMult = _wild ? 4 : 1;
            for (var j = 0; j < N_DEBRIS; j++) {
                dbLife[j] += dt / dbMaxLife[j] * debrisMult;
                if (dbLife[j] > 1.0) resetDebris(j);
                var di = j * 3;
                dbArr[di] += dbVX[j] * debrisMult;
                dbArr[di+1] += dbVY[j] * debrisMult;
                dbArr[di+2] += dbVZ[j] * debrisMult;
            }
            debrisPoints.geometry.attributes.position.needsUpdate = true;

            // Camera shake in wild mode
            if (_wild) {
                camera.position.x = Math.sin(realT * 2.7) * 5;
                camera.position.z = 45 + Math.sin(realT * 1.3) * 10;
                camera.fov = 60 + Math.sin(realT * 3.1) * 12;
                camera.updateProjectionMatrix();
            } else {
                camera.position.x = Math.sin(realT * 0.15) * 2;
                camera.position.z = 45;
                if (camera.fov !== 60) { camera.fov = 60; camera.updateProjectionMatrix(); }
            }

            renderer.render(scene, camera);
        }

        initThree();

        // Duke Nukem effects — explosive flash on open, screen shake on nav
        window.snonuxOpenEffect = function(post) {
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-fly'); setTimeout(function() { modal.classList.remove('sno-modal-fly'); }, 380); }
            // Explosion burst from post
            var r = post ? post.getBoundingClientRect() : {left: window.innerWidth/2, top: window.innerHeight/2, width:0, height:0};
            for (var i = 0; i < 8; i++) {
                (function(i) {
                    var e = document.createElement('div');
                    var ang = (i / 8) * Math.PI * 2; var dist = 50 + Math.random() * 70;
                    var c = i % 2 === 0 ? '255,0,0' : '255,215,0';
                    e.style.cssText = 'position:fixed;top:' + (r.top+(r.height||0)/2) + 'px;left:' + (r.left+(r.width||0)/2) + 'px;z-index:997;pointer-events:none;width:6px;height:6px;border-radius:50%;background:rgba(' + c + ',0.9);transition:all 0.4s ease,opacity 0.4s';
                    document.body.appendChild(e);
                    setTimeout(function() { e.style.transform='translate(' + (Math.cos(ang)*dist) + 'px,' + (Math.sin(ang)*dist-40) + 'px) scale(0.2)'; e.style.opacity='0'; setTimeout(function() { e.remove(); }, 450); }, 16);
                })(i);
            }
        };
        window.snonuxCloseEffect = function() {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(255,0,0,0.15);transition:opacity 0.2s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 220); }, 15);
        };
        window.snonuxScrollEffect = function(dir) {
            var isDown = dir === 'down';
            var thick = _wild ? '14px' : '5px';
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;left:0;right:0;height:' + thick + ';z-index:9000;pointer-events:none;' +
                'background:linear-gradient(90deg,transparent,rgba(255,0,0,0.9),rgba(255,215,0,0.9),rgba(255,0,0,0.9),transparent);' +
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
            // Heavy screen shake
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-shake'); setTimeout(function() { ov.classList.remove('sno-fx-shake'); }, 400); }
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(255,215,0,0.15);transition:opacity 0.2s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 220); }, 25);
        };
        window.snonuxPageEffect = function() {
            // Double explosion flash
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(255,0,0,0.25);transition:opacity 0.18s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 200); }, 15);
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-shake'); setTimeout(function() { ov.classList.remove('sno-fx-shake'); setTimeout(function() { ov.classList.add('sno-fx-shake'); setTimeout(function() { ov.classList.remove('sno-fx-shake'); }, 300); }, 30); }, 320); }
        };
    })();
