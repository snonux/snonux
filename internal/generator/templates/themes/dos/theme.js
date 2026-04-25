
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,drops=[],t0=performance.now();
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:false,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(1);
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(50,1,0.1,80);ca.position.z=20;
        var geo=new THREE.PlaneGeometry(0.22,0.32);
        for(var i=0;i<60;i++){
            var mat=new THREE.MeshBasicMaterial({color:0x55ff55,transparent:true,opacity:0.3+Math.random()*0.4});
            var m=new THREE.Mesh(geo,mat);
            m.position.set((Math.random()-0.5)*28, Math.random()*22-11, (Math.random()-0.5)*5);
            m.userData.speed=0.5+Math.random()*1.5;
            sc.add(m); drops.push(m);
        }
        sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);
            for(var i=0;i<drops.length;i++){
                drops[i].position.y-=drops[i].userData.speed*0.06;
                if(drops[i].position.y<-12) drops[i].position.y=12;
            }
            ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();


    (function() {
        var _wild = false, _snoTOffset = 0, _snoLastT = 0;
        var scene, camera, renderer, clock;
        var columns = [];

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x000088);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 200);
            camera.position.set(0, 0, 40);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: false });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(1);
            clock = new THREE.Clock();

            var geo = new THREE.PlaneGeometry(0.35, 0.5);

            for (var c = 0; c < 30; c++) {
                var col = [];
                var x = (c - 15) * 2.2;
                var speed = 1.5 + Math.random() * 3;
                var startY = Math.random() * 60 - 30;
                for (var r = 0; r < 8; r++) {
                    var brightness = 1.0 - (r / 8) * 0.7;
                    var color = new THREE.Color(brightness * 0.33, brightness, brightness * 0.33);
                    var mat = new THREE.MeshBasicMaterial({ color: color, transparent: true, opacity: brightness * 0.5 });
                    var mesh = new THREE.Mesh(geo, mat);
                    mesh.position.set(x, startY - r * 0.7, 0);
                    scene.add(mesh);
                    col.push({ mesh: mesh, offset: r * 0.7 });
                }
                columns.push({ chars: col, x: x, speed: speed, y: startY });
            }

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
            _snoTOffset += (realT - _snoLastT) * (_wild ? 9 : 0);
            _snoLastT = realT;
            var t = realT + _snoTOffset;
            for (var c = 0; c < columns.length; c++) {
                var col = columns[c];
                var y = col.y - t * col.speed;
                y = ((y % 60) + 60) % 60 - 30;
                for (var r = 0; r < col.chars.length; r++) {
                    col.chars[r].mesh.position.y = y - col.chars[r].offset;
                    // Wild: each char jitters horizontally like corrupted RAM
                    if (_wild) { col.chars[r].mesh.position.x = col.x + (Math.random() - 0.5) * 2.5; }
                    else { col.chars[r].mesh.position.x = col.x; }
                }
            }
            // Wild: camera lunges forward/back and sways like a CRT meltdown
            if (_wild) {
                camera.position.z = 40 + Math.sin(realT * 0.41) * 14;
                camera.position.x = Math.sin(realT * 0.37) * 8;
                camera.fov = 60 + Math.sin(realT * 0.53) * 16;
                camera.updateProjectionMatrix();
            } else {
                camera.position.z = 40;
                camera.position.x = 0;
                if (camera.fov !== 60) { camera.fov = 60; camera.updateProjectionMatrix(); }
            }
            renderer.render(scene, camera);
        }

        initThree();

        // DOS nav/wild effects — CRT glitch on navigate, system crash rain on wild
        window.snonuxOpenEffect = function() {
            // Slide in like a dialog box appearing on DOS screen
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-slide'); setTimeout(function() { modal.classList.remove('sno-modal-slide'); }, 360); }
            // CRT scan flash from top
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;top:0;left:0;right:0;height:4px;z-index:997;pointer-events:none;background:rgba(85,255,255,0.7);box-shadow:0 0 8px rgba(85,255,255,0.5);transition:top 0.28s linear,opacity 0.1s 0.28s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.top='100vh'; setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 120); }, 280); }, 15);
        };
        window.snonuxCloseEffect = function() {
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-glitch'); setTimeout(function() { ov.classList.remove('sno-fx-glitch'); }, 280); }
        };
        window.snonuxScrollEffect = function(dir) {
            var isDown = dir === 'down';
            var thick = _wild ? '14px' : '5px';
            var d = document.createElement('div');
            // DOS/CRT: grey-to-white scanline
            d.style.cssText = 'position:fixed;left:0;right:0;height:' + thick + ';z-index:9000;pointer-events:none;' +
                'background:linear-gradient(90deg,transparent,rgba(180,180,180,0.9),rgba(255,255,255,0.9),rgba(180,180,180,0.9),transparent);' +
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
            // CRT horizontal glitch
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-glitch'); setTimeout(function() { ov.classList.remove('sno-fx-glitch'); }, 300); }
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(85,255,85,0.12);transition:opacity 0.15s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 180); }, 25);
        };
        window.snonuxPageEffect = function() {
            // System crash — scanline strobe
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-glitch'); setTimeout(function() { ov.classList.remove('sno-fx-glitch'); setTimeout(function() { ov.classList.add('sno-fx-glitch'); setTimeout(function() { ov.classList.remove('sno-fx-glitch'); }, 280); }, 40); }, 310); }
        };
    })();
