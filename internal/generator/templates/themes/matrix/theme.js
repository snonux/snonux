
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,pts,pos,t0=performance.now(),N=28,i,arr;
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren)ren.dispose();ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(55,1,0.1,80);ca.position.set(0,0.5,10);
        arr=new Float32Array(N*3*20);for(i=0;i<arr.length;i+=3){arr[i]=(Math.random()-0.5)*16;arr[i+1]=Math.random()*22;arr[i+2]=(Math.random()-0.5)*8;}
        var geo=new THREE.BufferGeometry();geo.setAttribute('position',new THREE.BufferAttribute(arr,3));
        pts=new THREE.Points(geo,new THREE.PointsMaterial({color:0x00ff41,size:0.14,transparent:true,opacity:0.85,blending:THREE.AdditiveBlending,depthWrite:false,sizeAttenuation:true}));
        sc.add(pts);pos=geo.attributes.position;sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);var t=(now-t0)*0.001,j,p;
            for(j=0;j<pos.count;j++){p=j*3;pos.array[p+1]-=0.045+Math.sin(t+j*0.1)*0.012;if(pos.array[p+1]<-2)pos.array[p+1]=20;}
            pos.needsUpdate=true;pts.rotation.y=t*0.15;ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();


    // Matrix WebGL scene: 80 columns of falling particles with per-vertex colour.
    // Each column has a "head" that falls at a random speed; particles near the head
    // are bright green and fade to near-black further behind, simulating digital rain.
    (function() {
        var _wild = false;
        var NUM_COLS   = 80;   // number of rain columns
        var COL_LEN    = 25;   // particles per column
        var SPACING    = 2.2;  // vertical gap between particles in a column
        var Y_TOP      = 30;   // world-space top of the rain field
        var Y_BOTTOM   = -30;  // world-space bottom

        var scene, camera, renderer;
        var points;
        var posArr, colArr;
        // Per-column state: x position, head y, and fall speed
        var colX = [], headY = [], speed = [];

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x000000);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth / window.innerHeight, 0.1, 200);
            camera.position.set(0, 0, 50);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: false });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

            var totalPts = NUM_COLS * COL_LEN;
            posArr = new Float32Array(totalPts * 3);
            colArr = new Float32Array(totalPts * 3);

            // Spread columns across x: -50..50; initialise heads at random y positions
            for (var c = 0; c < NUM_COLS; c++) {
                colX[c]  = -50 + (c / (NUM_COLS - 1)) * 100;
                headY[c] = Y_TOP - Math.random() * (Y_TOP - Y_BOTTOM);
                speed[c] = 0.08 + Math.random() * 0.07; // 0.08–0.15 units per frame
            }

            var geo = new THREE.BufferGeometry();
            geo.setAttribute('position', new THREE.BufferAttribute(posArr, 3));
            geo.setAttribute('color',    new THREE.BufferAttribute(colArr, 3));

            var mat = new THREE.PointsMaterial({ size: 0.35, vertexColors: true });
            points = new THREE.Points(geo, mat);
            scene.add(points);

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

            for (var c = 0; c < NUM_COLS; c++) {
                // Advance the head downward each frame; wild mode accelerates rain 10x
                headY[c] -= speed[c] * (_wild ? 10 : 1);
                // When the head exits the bottom, reset to a random point above the top
                if (headY[c] < Y_BOTTOM - COL_LEN * SPACING) {
                    headY[c] = Y_TOP + Math.random() * 20;
                }

                var base = c * COL_LEN;
                for (var p = 0; p < COL_LEN; p++) {
                    var i = base + p;
                    var y = headY[c] + p * SPACING; // particles trail upward from head
                    posArr[i * 3]     = colX[c];
                    posArr[i * 3 + 1] = y;
                    posArr[i * 3 + 2] = 0;

                    // Brightness falls off with distance behind the head:
                    // p=0 is the head (bright), p=COL_LEN-1 is the tail (dim)
                    var bright = Math.max(0, 1 - p / (COL_LEN * 0.7));
                    // Head particle: #00ff41, tail: #003b00
                    colArr[i * 3]     = 0;
                    colArr[i * 3 + 1] = bright * (p === 0 ? 1.0 : 0.88);
                    colArr[i * 3 + 2] = bright * (p === 0 ? 0.255 : 0.04);
                }
            }

            points.geometry.attributes.position.needsUpdate = true;
            points.geometry.attributes.color.needsUpdate    = true;
            // Wild: camera plunges into the rain like riding a digital waterfall
            if (_wild) {
                var wt = Date.now() * 0.001;
                camera.position.z = 50 + Math.sin(wt * 0.38) * 20;
                camera.position.x = Math.sin(wt * 0.31) * 14;
                camera.fov = 60 + Math.sin(wt * 0.51) * 18;
                camera.updateProjectionMatrix();
            } else {
                camera.position.z = 50;
                camera.position.x = 0;
                if (camera.fov !== 60) { camera.fov = 60; camera.updateProjectionMatrix(); }
            }
            renderer.render(scene, camera);
        }

        initThree();

        // Matrix nav/wild effects — shake + green flash on navigate, rain storm on wild
        window.snonuxOpenEffect = function() {
            // Slide in from left like terminal output decoding
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-slide'); setTimeout(function() { modal.classList.remove('sno-modal-slide'); }, 380); }
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:997;pointer-events:none;background:rgba(0,255,65,0.1);transition:opacity 0.2s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 230); }, 15);
        };
        window.snonuxCloseEffect = function() {
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-glitch'); setTimeout(function() { ov.classList.remove('sno-fx-glitch'); }, 280); }
        };
        window.snonuxScrollEffect = function(dir) {
            var isDown = dir === 'down';
            var thick = _wild ? '14px' : '5px';
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;left:0;right:0;height:' + thick + ';z-index:9000;pointer-events:none;' +
                'background:linear-gradient(90deg,transparent,rgba(0,255,70,0.88),rgba(0,200,50,0.88),transparent);' +
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
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-shake'); setTimeout(function() { ov.classList.remove('sno-fx-shake'); }, 380); }
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(0,255,65,0.15);transition:opacity 0.2s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 220); }, 30);
        };
        window.snonuxPageEffect = function() {
            // Glitch bars — horizontal displacement flicker
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-glitch'); setTimeout(function() { ov.classList.remove('sno-fx-glitch'); }, 320); }
        };
    })();
