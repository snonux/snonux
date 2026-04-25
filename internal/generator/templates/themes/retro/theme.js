
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,g=new THREE.Group(),t0=performance.now();
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(48,1,0.1,60);ca.position.z=7.5;
        var bx=new THREE.Mesh(new THREE.BoxGeometry(2.6,2.6,2.6),new THREE.MeshBasicMaterial({color:0xffb000,wireframe:true,transparent:true,opacity:0.9}));
        var oc=new THREE.Mesh(new THREE.OctahedronGeometry(1.35,0),new THREE.MeshBasicMaterial({color:0xffb000,wireframe:true,transparent:true,opacity:0.55}));
        g.add(bx);g.add(oc);sc.add(g);sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);var t=(now-t0)*0.001;g.rotation.x=t*0.44;g.rotation.y=t*0.71;oc.rotation.z=t*0.9;ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();


    // Retro WebGL scene: amber demo-scene cube + orbiting octahedrons + star particles.
    // Evokes classic 80s/90s PC demo aesthetics with amber phosphor colours.
    (function() {
        var _wild = false, _snoTOffset = 0, _snoLastT = 0;
        var scene, camera, renderer;
        var mainCube, orbiters = [];
        var clock = new THREE.Clock();

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x0a0800);
            scene.fog = new THREE.Fog(0x0a0800, 25, 90);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth / window.innerHeight, 0.1, 200);
            camera.position.set(0, 0, 35);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

            // Central amber wireframe cube — the demo-scene hero piece
            var boxGeo = new THREE.BoxGeometry(8, 8, 8);
            var boxMat = new THREE.MeshBasicMaterial({ color: 0xffb000, wireframe: true });
            mainCube = new THREE.Mesh(boxGeo, boxMat);
            scene.add(mainCube);

            // 6 small octahedron wireframes evenly spaced on an orbital ring
            var octoMat = new THREE.MeshBasicMaterial({ color: 0xffb000, wireframe: true });
            for (var i = 0; i < 6; i++) {
                var octoGeo = new THREE.OctahedronGeometry(1.5);
                var octo = new THREE.Mesh(octoGeo, octoMat.clone());
                var angle = (i / 6) * Math.PI * 2;
                octo.position.set(Math.cos(angle) * 18, Math.sin(angle) * 5, Math.sin(angle) * 18);
                orbiters.push({ mesh: octo, baseAngle: angle });
                scene.add(octo);
            }

            // 800 dim amber background star particles
            var starGeo = new THREE.BufferGeometry();
            var starPos = new Float32Array(800 * 3);
            for (var j = 0; j < 800 * 3; j++) {
                starPos[j] = (Math.random() - 0.5) * 120;
            }
            starGeo.setAttribute('position', new THREE.BufferAttribute(starPos, 3));
            var starMat = new THREE.PointsMaterial({ color: 0x7a5200, size: 0.15 });
            scene.add(new THREE.Points(starGeo, starMat));

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
            _snoTOffset += (realT - _snoLastT) * (_wild ? 8 : 0);
            _snoLastT = realT;
            var t = realT + _snoTOffset;
            // Main cube rotates; wild mode spins dramatically + pulses in scale
            mainCube.rotation.y = t * 0.35;
            mainCube.rotation.x = t * 0.18;
            if (_wild) { mainCube.scale.setScalar(1 + 0.2 * Math.sin(realT * 5.1)); }
            else { mainCube.scale.setScalar(1); }
            // Orbiters fly outward to larger radius in wild mode
            var orbR = _wild ? 18 + Math.sin(realT * 0.9) * 10 : 18;
            for (var i = 0; i < orbiters.length; i++) {
                var o = orbiters[i];
                var angle = o.baseAngle + t * 0.4;
                o.mesh.position.x = Math.cos(angle) * orbR;
                o.mesh.position.z = Math.sin(angle) * orbR;
                o.mesh.position.y = Math.sin(angle * 0.7) * (_wild ? 10 : 4);
                o.mesh.rotation.x = t * 0.9;
                o.mesh.rotation.z = t * 0.6;
            }
            // Camera: frenetic sway in wild, static otherwise
            if (_wild) {
                camera.position.x = Math.sin(realT * 0.39) * 14;
                camera.position.y = Math.sin(realT * 0.31) * 8;
                camera.position.z = 35 + Math.sin(realT * 0.23) * 14;
                camera.fov = 60 + Math.sin(realT * 0.48) * 16;
                camera.updateProjectionMatrix();
            } else {
                camera.position.set(0, 0, 35);
                if (camera.fov !== 60) { camera.fov = 60; camera.updateProjectionMatrix(); }
            }
            renderer.render(scene, camera);
        }

        initThree();

        // Retro nav/wild effects — amber ember burst on navigate, phosphor burnout on wild
        window.snonuxOpenEffect = function() {
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-zoom'); setTimeout(function() { modal.classList.remove('sno-modal-zoom'); }, 400); }
            // Amber phosphor glow flash
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:997;pointer-events:none;background:radial-gradient(ellipse at center,rgba(255,176,0,0.18) 0%,transparent 65%);transition:opacity 0.3s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 340); }, 15);
        };
        window.snonuxCloseEffect = function() {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(255,176,0,0.1);transition:opacity 0.2s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 230); }, 15);
            document.body.style.animationDuration = '11s';
        };
        window.snonuxScrollEffect = function(dir) {
            var isDown = dir === 'down';
            var thick = _wild ? '14px' : '5px';
            var d = document.createElement('div');
            // Retro: warm amber wave
            d.style.cssText = 'position:fixed;left:0;right:0;height:' + thick + ';z-index:9000;pointer-events:none;' +
                'background:linear-gradient(90deg,transparent,rgba(255,176,0,0.9),rgba(200,100,0,0.9),rgba(255,176,0,0.9),transparent);' +
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
            // Intensify the flicker in wild mode
            document.body.style.animationDuration = _wild ? '1.5s' : '11s';
        };
        window.snonuxNavEffect = function() {
            // Amber ember spray flash
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:radial-gradient(ellipse at 50% 50%,rgba(255,176,0,0.22) 0%,transparent 65%);transition:opacity 0.25s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 280); }, 25);
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-shake'); setTimeout(function() { ov.classList.remove('sno-fx-shake'); }, 380); }
        };
        window.snonuxPageEffect = function() {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(255,176,0,0.2);transition:opacity 0.18s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 200); }, 20);
        };
    })();
