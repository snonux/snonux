
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,clock,core;
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(48,1,0.1,60);ca.position.z=9;clock=new THREE.Clock();
        core=new THREE.Mesh(new THREE.SphereGeometry(1.2,24,24),new THREE.MeshBasicMaterial({color:0xf55b7d,transparent:true,opacity:0.76})); sc.add(core);
        var shell=new THREE.Mesh(new THREE.TorusKnotGeometry(2.4,0.36,80,14),new THREE.MeshBasicMaterial({color:0x93ffd8,wireframe:true,transparent:true,opacity:0.42})); sc.add(shell); shell.userData.rot=0.006;
        sz();window.addEventListener('resize',sz);
        function loop(){ raf=requestAnimationFrame(loop); var t=clock.getElapsedTime(); shell.rotation.x=t*0.2; shell.rotation.y=t*0.3; core.scale.setScalar(1+Math.sin(t*3.2)*0.08); ren.render(sc,ca); }
        raf=requestAnimationFrame(loop);
    })();


    (function() {
        var _wild = false, _snoTOffset = 0, _snoLastT = 0;
        var scene, camera, renderer, clock, core, shellA, shellB, orbiters = [];

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x09070d);
            scene.fog = new THREE.Fog(0x09070d, 18, 120);
            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 220);
            camera.position.set(0, 6, 26);
            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();
            scene.add(new THREE.AmbientLight(0x553a47, 0.45));
            var coreLight = new THREE.PointLight(0xf55b7d, 1.6, 80); coreLight.position.set(0,0,0); scene.add(coreLight);

            core = new THREE.Mesh(new THREE.SphereGeometry(4.2, 36, 36), new THREE.MeshPhongMaterial({ color:0x803f5d, emissive:0xf55b7d, emissiveIntensity:0.52, shininess:90 }));
            shellA = new THREE.Mesh(new THREE.TorusKnotGeometry(7.4, 0.45, 180, 24, 2, 5), new THREE.MeshBasicMaterial({ color:0x93ffd8, wireframe:true, transparent:true, opacity:0.34 }));
            shellB = new THREE.Mesh(new THREE.TorusKnotGeometry(5.9, 0.28, 160, 16, 3, 7), new THREE.MeshBasicMaterial({ color:0xd0c7bb, wireframe:true, transparent:true, opacity:0.18 }));
            scene.add(core); scene.add(shellA); scene.add(shellB);
            for (var i = 0; i < 9; i++) {
                var orb = new THREE.Mesh(new THREE.SphereGeometry(0.55 + Math.random() * 0.45, 14, 14), new THREE.MeshPhongMaterial({ color: i % 2 === 0 ? 0x93ffd8 : 0xf55b7d, emissive: i % 2 === 0 ? 0x24473b : 0x5b1f32, emissiveIntensity:0.45 }));
                orb.userData.radius = 11 + Math.random() * 8;
                orb.userData.speed = 0.2 + Math.random() * 0.5;
                orb.userData.phase = Math.random() * Math.PI * 2;
                orbiters.push(orb); scene.add(orb);
            }
            window.addEventListener('resize', onResize);
            animate();
        }

        function onResize() { camera.aspect = window.innerWidth / window.innerHeight; camera.updateProjectionMatrix(); renderer.setSize(window.innerWidth, window.innerHeight); }

        function animate() {
            requestAnimationFrame(animate);
            var realT = clock.getElapsedTime();
            _snoTOffset += (realT - _snoLastT) * (_wild ? 11 : 0);
            _snoLastT = realT;
            var t = realT + _snoTOffset;
            core.scale.setScalar(1 + Math.sin(t * (_wild ? 6 : 1.8)) * (_wild ? 0.22 : 0.08));
            shellA.rotation.x = t * (_wild ? 0.9 : 0.25); shellA.rotation.y = t * (_wild ? 1.2 : 0.32);
            shellB.rotation.y = -t * (_wild ? 1.1 : 0.22); shellB.rotation.z = t * (_wild ? 0.8 : 0.18);
            shellA.material.opacity = _wild ? 0.54 : 0.34;
            for (var i = 0; i < orbiters.length; i++) {
                var o = orbiters[i], a = t * o.userData.speed + o.userData.phase;
                o.position.set(Math.cos(a) * o.userData.radius, Math.sin(a * 1.4) * 4, Math.sin(a) * o.userData.radius * 0.7);
            }
            camera.position.x = Math.sin(realT * (_wild ? 1.8 : 0.35)) * (_wild ? 3.2 : 1.1);
            camera.position.y = 6 + Math.sin(realT * (_wild ? 1.2 : 0.28)) * (_wild ? 1.8 : 0.4);
            camera.lookAt(0, 0, 0);
            renderer.render(scene, camera);
        }

        initThree();

        function flash(css, ms) {
            var d=document.createElement('div');
            d.style.cssText='position:fixed;inset:0;z-index:998;pointer-events:none;'+css+';transition:opacity '+(ms||220)+'ms';
            document.body.appendChild(d);
            setTimeout(function(){d.style.opacity='0';setTimeout(function(){d.remove();},ms||220);},25);
        }
        window.snonuxOpenEffect = function() {
            var modal=document.getElementById('post-modal');
            if(modal){modal.classList.add('sno-modal-expand');setTimeout(function(){modal.classList.remove('sno-modal-expand');},400);}
            flash('background:radial-gradient(circle at center,rgba(245,91,125,0.16),transparent 70%)',240);
        };
        window.snonuxCloseEffect = function(){ flash('background:rgba(0,0,0,0.3)',160); };
        window.snonuxNavEffect = function(){ flash('background:linear-gradient(90deg,transparent,rgba(147,255,216,0.1),transparent)',160); };
        window.snonuxPageEffect = function(){ flash('background:radial-gradient(circle at center,rgba(147,255,216,0.12),transparent 72%)',220); };
        window.snonuxScrollEffect = function(dir){
            var d=document.createElement('div');
            d.style.cssText='position:fixed;'+(dir==='down'?'top:0;':'bottom:0;')+'left:0;right:0;height:'+(_wild?'16px':'6px')+';z-index:9000;pointer-events:none;background:linear-gradient(90deg,transparent,rgba(245,91,125,0.8),rgba(147,255,216,0.7),transparent);transition:transform 0.32s ease,opacity 0.32s ease;';
            document.body.appendChild(d);
            setTimeout(function(){d.style.transform=dir==='down'?'translateY(100vh)':'translateY(-100vh)';d.style.opacity='0';},16);
            setTimeout(function(){d.remove();},380);
        };
        window.snonuxWildToggle = function(){ _wild=!_wild; var b=document.getElementById('sno-wild-badge'); if(b)b.classList.toggle('sno-wild-on',_wild); };
    })();
