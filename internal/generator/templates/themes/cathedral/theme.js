
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,clock,embers,rose;
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(46,1,0.1,80);ca.position.set(0,6,22);clock=new THREE.Clock();
        var floor=new THREE.Mesh(new THREE.PlaneGeometry(28,42),new THREE.MeshBasicMaterial({color:0x120f16,transparent:true,opacity:0.95})); floor.rotation.x=-Math.PI/2; floor.position.y=-4; sc.add(floor);
        for(var i=0;i<8;i++){ var c=new THREE.Mesh(new THREE.CylinderGeometry(0.38,0.52,11,8),new THREE.MeshBasicMaterial({color:0x302737}));
            c.position.set(i<4?-6.2:6.2,1,-16+(i%4)*8); sc.add(c);}
        rose=new THREE.Mesh(new THREE.CircleGeometry(3.4,34),new THREE.MeshBasicMaterial({color:0xd9bf78,transparent:true,opacity:0.26})); rose.position.set(0,8,-16); sc.add(rose);
        var beam=new THREE.Mesh(new THREE.ConeGeometry(3.8,12,20,1,true),new THREE.MeshBasicMaterial({color:0xd9bf78,transparent:true,opacity:0.08,side:THREE.DoubleSide})); beam.position.set(0,4,-6); beam.rotation.x=Math.PI; sc.add(beam);
        var ep=new Float32Array(240*3); for(i=0;i<ep.length;i+=3){ ep[i]=(Math.random()-0.5)*14; ep[i+1]=Math.random()*8; ep[i+2]=-14+Math.random()*18; }
        var eg=new THREE.BufferGeometry(); eg.setAttribute('position',new THREE.BufferAttribute(ep,3));
        embers=new THREE.Points(eg,new THREE.PointsMaterial({color:0xffb35a,size:0.12,transparent:true,opacity:0.24})); sc.add(embers);
        sz();window.addEventListener('resize',sz);
        function loop(){ raf=requestAnimationFrame(loop); var t=clock.getElapsedTime(),pos=embers.geometry.attributes.position;
            for(var i=0;i<pos.count;i++){ var y=pos.getY(i)+0.05; pos.setY(i,y>10?0:y); }
            pos.needsUpdate=true; rose.rotation.z=t*0.08; beam.material.opacity=0.08+Math.sin(t*1.5)*0.02; ren.render(sc,ca); }
        raf=requestAnimationFrame(loop);
    })();


    (function() {
        var _wild = false, _snoTOffset = 0, _snoLastT = 0;
        var scene, camera, renderer, clock, dust, embers, beams = [], candles = [], rose, halo, chandelier, pipes = [];

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x0f0d14);
            scene.fog = new THREE.Fog(0x0f0d14, 16, 140);
            camera = new THREE.PerspectiveCamera(56, window.innerWidth/window.innerHeight, 0.1, 260);
            camera.position.set(0, 10, 40);
            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();

            scene.add(new THREE.AmbientLight(0x50394b, 0.38));
            var floor = new THREE.Mesh(new THREE.PlaneGeometry(90, 220, 1, 1), new THREE.MeshPhongMaterial({ color:0x19141c, shininess:12 }));
            floor.rotation.x = -Math.PI/2; floor.position.y = -3; scene.add(floor);
            var aisle = new THREE.Mesh(new THREE.PlaneGeometry(10, 210), new THREE.MeshBasicMaterial({ color:0x3d2830, transparent:true, opacity:0.24 }));
            aisle.rotation.x = -Math.PI/2; aisle.position.set(0,-2.98,-38); scene.add(aisle);
            rose = new THREE.Mesh(new THREE.CircleGeometry(7.2, 48), new THREE.MeshBasicMaterial({ color:0xd9bf78, transparent:true, opacity:0.18 }));
            rose.position.set(0, 16, -102); scene.add(rose);
            halo = new THREE.Mesh(new THREE.CircleGeometry(10.8, 48), new THREE.MeshBasicMaterial({ color:0x71233d, transparent:true, opacity:0.08 }));
            halo.position.set(0, 16, -103); scene.add(halo);
            chandelier = new THREE.Mesh(new THREE.TorusGeometry(4.5, 0.18, 12, 44), new THREE.MeshBasicMaterial({ color:0xe0c47f, transparent:true, opacity:0.5 }));
            chandelier.position.set(0, 18, -16); scene.add(chandelier);

            for (var i = 0; i < 18; i++) {
                var side = i < 9 ? -1 : 1;
                var z = -95 + (i % 9) * 12;
                var col = new THREE.Mesh(new THREE.CylinderGeometry(0.9, 1.05, 20, 10), new THREE.MeshPhongMaterial({ color:0x2d2632 }));
                col.position.set(side * 13, 7, z); scene.add(col);
                var pipe = new THREE.Mesh(new THREE.BoxGeometry(1.4, 12 + (i % 5) * 2.8, 1.4), new THREE.MeshPhongMaterial({ color:0x55474b, shininess:28 }));
                pipe.position.set(side * 18, pipe.geometry.parameters.height * 0.5 - 1, z - 4); scene.add(pipe); pipes.push(pipe);
                var beam = new THREE.Mesh(new THREE.ConeGeometry(4.4, 26, 22, 1, true), new THREE.MeshBasicMaterial({ color:i % 2 === 0 ? 0x7bc2ff : 0x8e2f49, transparent:true, opacity:0.08, side:THREE.DoubleSide }));
                beam.position.set(side * 9, 9, z); beam.rotation.x = Math.PI; scene.add(beam); beams.push(beam);
            }

            for (i = 0; i < 14; i++) {
                var flame = new THREE.PointLight(i % 3 === 0 ? 0xffd58a : 0xffb35a, 0.52, 18);
                flame.position.set((i % 2 === 0 ? -4.5 : 4.5) + (Math.random() - 0.5), 1.4, -10 - i * 7.6);
                scene.add(flame); candles.push(flame);
            }

            var dp = new Float32Array(1400 * 3);
            for (i = 0; i < dp.length; i += 3) { dp[i]=(Math.random()-0.5)*42; dp[i+1]=Math.random()*28; dp[i+2]=-120+Math.random()*120; }
            var dg = new THREE.BufferGeometry(); dg.setAttribute('position', new THREE.BufferAttribute(dp, 3));
            dust = new THREE.Points(dg, new THREE.PointsMaterial({ color:0xf0e8d9, size:0.12, transparent:true, opacity:0.34 }));
            scene.add(dust);
            var ep = new Float32Array(460 * 3);
            for (i = 0; i < ep.length; i += 3) { ep[i]=(Math.random()-0.5)*24; ep[i+1]=Math.random()*16; ep[i+2]=-100+Math.random()*100; }
            var eg = new THREE.BufferGeometry(); eg.setAttribute('position', new THREE.BufferAttribute(ep, 3));
            embers = new THREE.Points(eg, new THREE.PointsMaterial({ color:0xffb35a, size:0.18, transparent:true, opacity:0.0 }));
            scene.add(embers);
            window.addEventListener('resize', onResize);
            animate();
        }

        function onResize() { camera.aspect = window.innerWidth / window.innerHeight; camera.updateProjectionMatrix(); renderer.setSize(window.innerWidth, window.innerHeight); }

        function animate() {
            requestAnimationFrame(animate);
            var realT = clock.getElapsedTime();
            _snoTOffset += (realT - _snoLastT) * (_wild ? 9 : 0);
            _snoLastT = realT;
            var t = realT + _snoTOffset;
            var pos = dust.geometry.attributes.position;
            for (var i = 0; i < pos.count; i++) {
                var y = pos.getY(i) + (_wild ? 0.12 : 0.02);
                var x = pos.getX(i) + Math.sin(t * 0.3 + i) * (_wild ? 0.03 : 0.005);
                pos.setX(i, x); pos.setY(i, y > 26 ? -2 : y);
            }
            pos.needsUpdate = true;

            var emberPos = embers.geometry.attributes.position;
            for (i = 0; i < emberPos.count; i++) {
                var ey = emberPos.getY(i) + (_wild ? 0.2 : 0.01);
                emberPos.setY(i, ey > 22 ? 0 : ey);
                if (_wild) emberPos.setX(i, emberPos.getX(i) + Math.sin(realT * 0.8 + i) * 0.04);
            }
            emberPos.needsUpdate = true;
            embers.material.opacity = _wild ? 0.74 : 0.0;

            rose.rotation.z = realT * (_wild ? 0.9 : 0.08);
            halo.scale.setScalar(1 + Math.sin(realT * (_wild ? 3.8 : 1.2)) * (_wild ? 0.18 : 0.03));
            chandelier.rotation.z = Math.sin(realT * (_wild ? 1.6 : 0.24)) * (_wild ? 0.2 : 0.03);
            chandelier.rotation.x = Math.sin(realT * (_wild ? 1.2 : 0.18)) * (_wild ? 0.12 : 0.02);

            for (i = 0; i < beams.length; i++) beams[i].material.opacity = (_wild ? 0.18 : 0.08) + Math.sin(t * 0.8 + i) * 0.02;
            for (i = 0; i < candles.length; i++) candles[i].intensity = (_wild ? 1.18 : 0.52) + Math.sin(realT * 5 + i) * 0.12;
            for (i = 0; i < pipes.length; i++) pipes[i].scale.y = 1 + Math.sin(realT * (_wild ? 4.2 : 0.4) + i) * (_wild ? 0.18 : 0.015);

            camera.position.x = Math.sin(realT * (_wild ? 1.4 : 0.16)) * (_wild ? 3.4 : 0.6);
            camera.position.y = 10 + Math.sin(realT * 0.3) * (_wild ? 1.8 : 0.4);
            camera.position.z = _wild ? 35 + Math.sin(realT * 0.6) * 4 : 40;
            camera.lookAt(0, 5, -48);
            renderer.render(scene, camera);
        }

        initThree();

        function veil(css, ms) {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;' + css + ';transition:opacity ' + (ms || 240) + 'ms';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, ms || 240); }, 25);
        }
        window.snonuxOpenEffect = function() {
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-zoom'); setTimeout(function() { modal.classList.remove('sno-modal-zoom'); }, 390); }
            veil('background:radial-gradient(ellipse at center,rgba(224,196,127,0.18),rgba(123,194,255,0.08),transparent 70%)', 260);
        };
        window.snonuxCloseEffect = function() { veil('background:rgba(0,0,0,0.3)', 180); };
        window.snonuxNavEffect = function() { veil('background:linear-gradient(90deg,transparent,rgba(123,194,255,0.12),rgba(224,196,127,0.08),transparent)', 190); };
        window.snonuxPageEffect = function() { veil('background:radial-gradient(ellipse at center,rgba(142,47,73,0.22),rgba(224,196,127,0.12),transparent 72%)', 260); };
        window.snonuxScrollEffect = function(dir) {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;' + (dir === 'down' ? 'top:0;' : 'bottom:0;') + 'left:0;right:0;height:' + (_wild ? '16px' : '6px') + ';z-index:9000;pointer-events:none;background:linear-gradient(90deg,transparent,rgba(224,196,127,0.84),rgba(123,194,255,0.62),rgba(142,47,73,0.54),transparent);transition:transform 0.34s ease,opacity 0.34s ease;';
            document.body.appendChild(d);
            setTimeout(function() { d.style.transform = dir === 'down' ? 'translateY(100vh)' : 'translateY(-100vh)'; d.style.opacity='0'; }, 16);
            setTimeout(function() { d.remove(); }, 400);
        };
        window.snonuxWildToggle = function() {
            _wild = !_wild;
            var b = document.getElementById('sno-wild-badge');
            if (b) b.classList.toggle('sno-wild-on', _wild);
        };
    })();
