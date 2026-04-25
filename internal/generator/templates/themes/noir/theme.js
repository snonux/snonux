
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,clock,rain;
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(50,1,0.1,80);ca.position.set(0,4,20);clock=new THREE.Clock();
        var street=new THREE.Mesh(new THREE.PlaneGeometry(28,24),new THREE.MeshBasicMaterial({color:0x0d0d0d,transparent:true,opacity:0.92}));
        street.rotation.x=-Math.PI/2;street.position.y=-3;sc.add(street);
        for(var i=0;i<5;i++){ var b=new THREE.Mesh(new THREE.BoxGeometry(3+Math.random()*2,8+Math.random()*9,3+Math.random()*2),new THREE.MeshBasicMaterial({color:0x111111}));
            b.position.set(-10+i*5,1.5+b.geometry.parameters.height*0.5,-10-Math.random()*8);sc.add(b);}
        var lamp=new THREE.Mesh(new THREE.CylinderGeometry(0.08,0.12,8,6),new THREE.MeshBasicMaterial({color:0x3d3d3d})); lamp.position.set(0,1,-3); sc.add(lamp);
        var glow=new THREE.Mesh(new THREE.SphereGeometry(0.5,12,12),new THREE.MeshBasicMaterial({color:0xf0ead6,transparent:true,opacity:0.85})); glow.position.set(0,5,-3); sc.add(glow);
        var cone=new THREE.Mesh(new THREE.ConeGeometry(4.5,10,18,1,true),new THREE.MeshBasicMaterial({color:0xf0ead6,transparent:true,opacity:0.12,side:THREE.DoubleSide}));
        cone.position.set(0,0,-3); cone.rotation.x=Math.PI; sc.add(cone);
        var rp=new Float32Array(700*3); for(i=0;i<rp.length;i+=3){ rp[i]=(Math.random()-0.5)*24; rp[i+1]=Math.random()*20-2; rp[i+2]=(Math.random()-0.5)*20; }
        var rg=new THREE.BufferGeometry(); rg.setAttribute('position',new THREE.BufferAttribute(rp,3));
        rain=new THREE.Points(rg,new THREE.PointsMaterial({color:0xd8d1c4,size:0.08,transparent:true,opacity:0.38})); sc.add(rain);
        sz(); window.addEventListener('resize',sz);
        function loop(){ raf=requestAnimationFrame(loop); var t=clock.getElapsedTime(),pos=rain.geometry.attributes.position;
            for(var i=0;i<pos.count;i++){ var y=pos.getY(i)-0.32; pos.setY(i,y<-3?18:y); }
            pos.needsUpdate=true; glow.scale.setScalar(1+Math.sin(t*2.3)*0.05); cone.material.opacity=0.1+Math.sin(t*1.8)*0.03; ren.render(sc,ca); }
        raf=requestAnimationFrame(loop);
    })();


    (function() {
        var _wild = false, _snoTOffset = 0, _snoLastT = 0;
        var scene, camera, renderer, clock, rain, leftSweep, rightSweep, street, buildings = [], fogPlanes = [], signPlane, lampHalo;

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x050505);
            scene.fog = new THREE.Fog(0x050505, 18, 120);
            camera = new THREE.PerspectiveCamera(58, window.innerWidth/window.innerHeight, 0.1, 200);
            camera.position.set(0, 10, 34);
            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();

            street = new THREE.Mesh(new THREE.PlaneGeometry(90, 180, 1, 1), new THREE.MeshPhongMaterial({ color:0x111111, shininess:8 }));
            street.rotation.x = -Math.PI/2; street.position.y = -2; scene.add(street);
            var stripe = new THREE.Mesh(new THREE.PlaneGeometry(1.2, 120), new THREE.MeshBasicMaterial({ color:0xf0ead6, transparent:true, opacity:0.08 }));
            stripe.rotation.x = -Math.PI/2; stripe.position.set(0,-1.98,-20); scene.add(stripe);
            scene.add(new THREE.AmbientLight(0x404040, 0.45));

            var lampLight = new THREE.PointLight(0xf0ead6, 1.2, 70); lampLight.position.set(0, 18, -18); scene.add(lampLight);
            lampHalo = new THREE.Mesh(new THREE.SphereGeometry(1.4, 16, 16), new THREE.MeshBasicMaterial({ color:0xf0ead6, transparent:true, opacity:0.12 }));
            lampHalo.position.set(0,18,-18); scene.add(lampHalo);
            leftSweep = new THREE.PointLight(0x223b88, 0.0, 60); leftSweep.position.set(-25, 8, -10); scene.add(leftSweep);
            rightSweep = new THREE.PointLight(0x882222, 0.0, 60); rightSweep.position.set(25, 8, -10); scene.add(rightSweep);

            for (var i = 0; i < 18; i++) {
                var h = 10 + Math.random() * 28, w = 4 + Math.random() * 4, d = 4 + Math.random() * 6;
                var b = new THREE.Mesh(new THREE.BoxGeometry(w, h, d), new THREE.MeshPhongMaterial({ color:0x131313 }));
                var side = i < 9 ? -1 : 1;
                b.position.set(side * (14 + Math.random() * 22), h * 0.5 - 2, -80 + (i % 9) * 10);
                scene.add(b); buildings.push(b);
            }
            signPlane = new THREE.Mesh(new THREE.PlaneGeometry(8, 2.4), new THREE.MeshBasicMaterial({ color:0xa9372b, transparent:true, opacity:0.34, side:THREE.DoubleSide }));
            signPlane.position.set(18, 10, -32); signPlane.rotation.y = -0.42; scene.add(signPlane);
            for (i = 0; i < 4; i++) {
                var fog = new THREE.Mesh(new THREE.PlaneGeometry(40, 10), new THREE.MeshBasicMaterial({ color:0xffffff, transparent:true, opacity:0.03, side:THREE.DoubleSide, depthWrite:false }));
                fog.position.set((Math.random()-0.5)*20, 1 + Math.random()*6, -50 + i * 16);
                scene.add(fog); fogPlanes.push(fog);
            }

            var rp = new Float32Array(2600 * 3);
            for (i = 0; i < rp.length; i += 3) {
                rp[i] = (Math.random() - 0.5) * 80;
                rp[i + 1] = Math.random() * 60;
                rp[i + 2] = -90 + Math.random() * 120;
            }
            var rg = new THREE.BufferGeometry();
            rg.setAttribute('position', new THREE.BufferAttribute(rp, 3));
            rain = new THREE.Points(rg, new THREE.PointsMaterial({ color:0xd8d1c4, size:0.1, transparent:true, opacity:0.46 }));
            scene.add(rain);

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

            var pos = rain.geometry.attributes.position;
            var speed = _wild ? 1.35 : 0.42;
            for (var i = 0; i < pos.count; i++) {
                var y = pos.getY(i) - speed;
                if (y < -4) { y = 56; pos.setX(i, (Math.random() - 0.5) * 80); }
                pos.setY(i, y);
                if (_wild) pos.setX(i, pos.getX(i) + Math.sin(t * 0.9 + i) * 0.02);
            }
            pos.needsUpdate = true;
            for (i = 0; i < fogPlanes.length; i++) {
                fogPlanes[i].position.x = Math.sin(realT * (0.18 + i * 0.05)) * (_wild ? 12 : 4);
                fogPlanes[i].material.opacity = (_wild ? 0.08 : 0.03) + Math.sin(realT * 0.7 + i) * 0.01;
            }
            signPlane.material.opacity = _wild ? (0.18 + Math.abs(Math.sin(realT * 8.2)) * 0.48) : (0.28 + Math.sin(realT * 1.4) * 0.06);
            lampHalo.scale.setScalar(1 + Math.sin(realT * (_wild ? 6 : 2.2)) * (_wild ? 0.18 : 0.05));

            leftSweep.intensity = _wild ? 1.4 + Math.sin(realT * 3.5) * 0.4 : 0;
            rightSweep.intensity = _wild ? 1.2 + Math.cos(realT * 3.2) * 0.4 : 0;
            leftSweep.position.x = -28 + Math.sin(realT * 1.7) * 10;
            rightSweep.position.x = 28 + Math.cos(realT * 1.6) * 10;

            camera.position.x = _wild ? Math.sin(realT * 1.9) * 2.2 : Math.sin(realT * 0.22) * 1.4;
            camera.position.y = 10 + Math.sin(realT * 0.3) * (_wild ? 1.4 : 0.5);
            camera.position.z = _wild ? 32 + Math.sin(realT * 0.7) * 4 : 34;
            camera.lookAt(0, 6, -35);
            renderer.render(scene, camera);
        }

        initThree();

        function flash(css, ms) {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;' + css + ';transition:opacity ' + (ms || 220) + 'ms';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity = '0'; setTimeout(function() { d.remove(); }, ms || 220); }, 25);
        }

        window.snonuxOpenEffect = function(post) {
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-fly'); setTimeout(function() { modal.classList.remove('sno-modal-fly'); }, 360); }
            var r = post ? post.getBoundingClientRect() : { left: innerWidth * 0.5, top: innerHeight * 0.5, width: 0, height: 0 };
            var s = document.createElement('div');
            s.style.cssText = 'position:fixed;left:' + (r.left + r.width * 0.5 - 12) + 'px;top:' + (r.top + r.height * 0.5 - 12) + 'px;width:24px;height:24px;border-radius:50%;z-index:997;pointer-events:none;background:radial-gradient(circle,rgba(240,234,214,0.88),rgba(240,234,214,0.18) 55%,transparent 72%);transition:transform 0.4s ease,opacity 0.4s ease;';
            document.body.appendChild(s);
            setTimeout(function() { s.style.transform='scale(18)'; s.style.opacity='0'; setTimeout(function() { s.remove(); }, 420); }, 18);
        };
        window.snonuxCloseEffect = function() { flash('background:rgba(0,0,0,0.45)', 180); };
        window.snonuxNavEffect = function() { flash('background:repeating-linear-gradient(90deg,rgba(0,0,0,0.8) 0 12%,rgba(240,234,214,0.06) 12% 14%,rgba(0,0,0,0.8) 14% 24%)', 170); };
        window.snonuxPageEffect = function() { flash('background:linear-gradient(90deg,rgba(36,65,130,0.14),transparent 38%,rgba(169,55,43,0.18) 62%,transparent)', 240); };
        window.snonuxScrollEffect = function(dir) {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;' + (dir === 'down' ? 'top:0;' : 'bottom:0;') + 'left:0;right:0;height:' + (_wild ? '18px' : '6px') + ';z-index:9000;pointer-events:none;background:linear-gradient(90deg,transparent,rgba(240,234,214,0.8),transparent);transition:transform 0.28s ease,opacity 0.28s ease;';
            document.body.appendChild(d);
            setTimeout(function() { d.style.transform = dir === 'down' ? 'translateY(100vh)' : 'translateY(-100vh)'; d.style.opacity='0'; }, 16);
            setTimeout(function() { d.remove(); }, 340);
        };
        window.snonuxWildToggle = function() {
            _wild = !_wild;
            var b = document.getElementById('sno-wild-badge');
            if (b) b.classList.toggle('sno-wild-on', _wild);
        };
    })();
