
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,clock,rings=[];
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(50,1,0.1,60);ca.position.z=10;clock=new THREE.Clock();
        for(var i=0;i<3;i++){ var r=new THREE.Mesh(new THREE.TorusGeometry(1.4+i*0.5,0.04,8,48),new THREE.MeshBasicMaterial({color:0x63f3a8,transparent:true,opacity:0.68-i*0.1})); sc.add(r); rings.push(r);}
        var iris=new THREE.Mesh(new THREE.CircleGeometry(0.4,24),new THREE.MeshBasicMaterial({color:0xbcffd4,transparent:true,opacity:0.8})); sc.add(iris);
        sz();window.addEventListener('resize',sz);
        function loop(){ raf=requestAnimationFrame(loop); var t=clock.getElapsedTime(); for(var i=0;i<rings.length;i++){ rings[i].rotation.z=t*(0.4+i*0.3); rings[i].scale.setScalar(1+Math.sin(t*2+i)*0.03); } ren.render(sc,ca); }
        raf=requestAnimationFrame(loop);
    })();


    (function() {
        var _wild = false, _snoTOffset = 0, _snoLastT = 0;
        var scene, camera, renderer, clock, nodes = [], trackers = [], rain;

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x09100d);
            scene.fog = new THREE.Fog(0x09100d, 20, 120);
            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 200);
            camera.position.set(0, 6, 28);
            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            clock = new THREE.Clock();
            scene.add(new THREE.AmbientLight(0x35634f, 0.55));

            for (var i = 0; i < 12; i++) {
                var s = new THREE.Mesh(new THREE.PlaneGeometry(7, 4.2), new THREE.MeshBasicMaterial({ color:0x15221c, transparent:true, opacity:0.92, side:THREE.DoubleSide }));
                s.position.set((i % 4 - 1.5) * 11, 8 - Math.floor(i / 4) * 6, -10 - Math.floor(i / 4) * 8);
                scene.add(s); nodes.push(s);
                var box = new THREE.LineSegments(new THREE.EdgesGeometry(new THREE.PlaneGeometry(6.4, 3.6)), new THREE.LineBasicMaterial({ color:0x63f3a8, transparent:true, opacity:0.5 }));
                box.position.copy(s.position); box.position.z += 0.04; scene.add(box); trackers.push(box);
            }

            var rp = new Float32Array(1500 * 3);
            for (i = 0; i < rp.length; i += 3) { rp[i]=(Math.random()-0.5)*70; rp[i+1]=(Math.random()-0.5)*40; rp[i+2]=-80+Math.random()*90; }
            var rg = new THREE.BufferGeometry(); rg.setAttribute('position', new THREE.BufferAttribute(rp, 3));
            rain = new THREE.Points(rg, new THREE.PointsMaterial({ color:0xbcffd4, size:0.08, transparent:true, opacity:0.2 }));
            scene.add(rain);
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
            for (var i = 0; i < nodes.length; i++) {
                nodes[i].material.opacity = (_wild ? 0.58 : 0.9) - ((i % 4) * 0.06);
                trackers[i].rotation.z = Math.sin(t * 0.7 + i) * (_wild ? 0.12 : 0.03);
                trackers[i].material.opacity = _wild ? 0.84 : 0.5;
            }
            var pos = rain.geometry.attributes.position;
            for (i = 0; i < pos.count; i++) { var y = pos.getY(i) - (_wild ? 0.22 : 0.08); pos.setY(i, y < -20 ? 20 : y); }
            pos.needsUpdate = true;
            camera.position.x = Math.sin(realT * (_wild ? 1.6 : 0.3)) * (_wild ? 2.8 : 0.7);
            camera.position.y = 6 + Math.cos(realT * 0.4) * (_wild ? 1.1 : 0.3);
            camera.lookAt(0, 0, -20);
            renderer.render(scene, camera);
        }

        initThree();

        function overlay(css, ms) {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;' + css + ';transition:opacity ' + (ms || 200) + 'ms';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, ms || 200); }, 25);
        }
        window.snonuxOpenEffect = function(post) {
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-slide'); setTimeout(function() { modal.classList.remove('sno-modal-slide'); }, 340); }
            var r = post ? post.getBoundingClientRect() : { left: innerWidth*0.5, top: innerHeight*0.5, width: 0, height: 0 };
            var box = document.createElement('div');
            box.style.cssText = 'position:fixed;left:' + (r.left-6) + 'px;top:' + (r.top-6) + 'px;width:' + (r.width+12) + 'px;height:' + (r.height+12) + 'px;border:1px solid rgba(99,243,168,0.8);z-index:997;pointer-events:none;transition:transform 0.32s ease,opacity 0.32s ease;';
            document.body.appendChild(box);
            setTimeout(function() { box.style.transform='scale(1.18)'; box.style.opacity='0'; setTimeout(function() { box.remove(); }, 360); }, 18);
        };
        window.snonuxCloseEffect = function() { overlay('background:rgba(0,0,0,0.32)', 160); };
        window.snonuxNavEffect = function() { overlay('background:linear-gradient(90deg,transparent,rgba(99,243,168,0.08),transparent)', 160); };
        window.snonuxPageEffect = function() { overlay('background:radial-gradient(circle at center,rgba(255,77,92,0.12),transparent 68%)', 220); };
        window.snonuxScrollEffect = function(dir) {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;' + (dir === 'down' ? 'top:0;' : 'bottom:0;') + 'left:0;right:0;height:' + (_wild ? '16px' : '6px') + ';z-index:9000;pointer-events:none;background:linear-gradient(90deg,transparent,rgba(99,243,168,0.82),transparent);transition:transform 0.28s ease,opacity 0.28s ease;';
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
