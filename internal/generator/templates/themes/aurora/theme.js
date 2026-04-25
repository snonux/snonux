
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,clock,ribbons=[],SEG=48;
        var cols=[0x00ffb3,0x00cfe8,0xc084fc,0x48e8d0,0xa855f7];
        var yPos=[2,5.5,9,12.5,16], zPos=[-18,-14,-10,-7,-4];
        function cleanup(){
            window.removeEventListener('resize',sz);
            if(raf)cancelAnimationFrame(raf);raf=null;
            ribbons.forEach(function(rb){rb.geo.dispose();rb.mesh.material.dispose();});
            ribbons=[];
            if(ren){ren.dispose();ren=null;}
            window._snonuxSplashWebGLCleanup=null;
        }
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){
            var w=cv.clientWidth||2,h=cv.clientHeight||2;
            if(ren)ren.setSize(w,h,false);
            if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}
        }
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});
        ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();
        ca=new THREE.PerspectiveCamera(52,1,0.1,120);
        ca.position.set(0,10,26);ca.lookAt(0,8,-6);
        clock=new THREE.Clock();
        for(var r=0;r<5;r++){
            var geo=new THREE.PlaneGeometry(100,7,SEG,1);
            var mat=new THREE.MeshBasicMaterial({
                color:cols[r],transparent:true,opacity:0.26+r*0.02,
                side:THREE.DoubleSide,blending:THREE.AdditiveBlending,depthWrite:false
            });
            var mesh=new THREE.Mesh(geo,mat);
            mesh.position.set(0,yPos[r],zPos[r]);
            sc.add(mesh);
            ribbons.push({mesh:mesh,geo:geo,freq:0.55+0.12*r,phase:r*1.15,amp:2.4+0.2*r});
        }
        function loop(){
            raf=requestAnimationFrame(loop);
            var t=clock.getElapsedTime();
            for(var i=0;i<ribbons.length;i++){
                var rb=ribbons[i],pos=rb.geo.attributes.position;
                for(var v=0;v<pos.count;v++){
                    if(pos.getY(v)>0){
                        var x=pos.getX(v);
                        pos.setY(v,rb.amp*Math.sin(t*rb.freq+x*0.065+rb.phase)
                            +rb.amp*0.38*Math.cos(t*rb.freq*0.72+x*0.042));
                    }
                }
                pos.needsUpdate=true;
            }
            ren.render(sc,ca);
        }
        sz();window.addEventListener('resize',sz);
        raf=requestAnimationFrame(loop);
    })();


    // Aurora WebGL: six wide ribbon meshes whose top-row vertices are animated
    // with overlapping sine waves, rendered with additive blending so they glow
    // against the dark navy sky like real aurora curtains.
    (function() {
        var _wild = false, _snoTOffset = 0, _snoLastT = 0;
        var RIBBON_COUNT = 6;
        var SEG_W = 60; // horizontal segments per ribbon
        var ribbonColors = [0x00ffb3, 0x00cfe8, 0xc084fc, 0x00ffb3, 0x48e8d0, 0xa855f7];
        var ribbonY     = [-10, -4, 2, 8, 14, 20];
        var ribbonZ     = [-40, -30, -22, -15, -10, -5];
        var ribbonFreq  = [0.6, 0.9, 0.7, 1.1, 0.5, 0.8];
        var ribbonPhase = [0.0, 1.2, 2.4, 0.8, 3.1, 1.7];
        var ribbonAmp   = [3.0, 2.5, 2.0, 3.5, 2.2, 2.8];

        var scene, camera, renderer, clock;
        var ribbons = [];

        function initThree() {
            scene = new THREE.Scene();
            scene.background = new THREE.Color(0x050d1a);
            scene.fog = new THREE.Fog(0x050d1a, 40, 120);

            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 200);
            camera.position.set(0, 5, 30);
            camera.lookAt(0, 5, 0);

            renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('three-canvas'), antialias: true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

            clock = new THREE.Clock();

            for (var r = 0; r < RIBBON_COUNT; r++) {
                // Wide shallow plane; we animate the top row of vertices
                var geo = new THREE.PlaneGeometry(120, 8, SEG_W, 1);
                var mat = new THREE.MeshBasicMaterial({
                    color: ribbonColors[r], transparent: true, opacity: 0.32,
                    side: THREE.DoubleSide, blending: THREE.AdditiveBlending, depthWrite: false
                });
                var mesh = new THREE.Mesh(geo, mat);
                mesh.position.set(0, ribbonY[r], ribbonZ[r]);
                scene.add(mesh);
                ribbons.push({ mesh: mesh, geo: geo, freq: ribbonFreq[r],
                               phase: ribbonPhase[r], amp: ribbonAmp[r] });
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
            // Wild mode: accelerate time 18× so waves churn much faster
            _snoTOffset += (realT - _snoLastT) * (_wild ? 18 : 0);
            _snoLastT = realT;
            var t = realT + _snoTOffset;

            var ampMult  = _wild ? 3.2 : 1;
            // Camera sways left/right and bobs up/down in wild mode
            camera.position.x = _wild ? Math.sin(t * 0.28) * 10 : 0;
            camera.position.y = _wild ? 5 + Math.cos(t * 0.19) * 4 : 5;

            for (var r = 0; r < ribbons.length; r++) {
                var rb = ribbons[r];
                var pos = rb.geo.attributes.position;
                var count = pos.count;
                // In wild mode ribbons also drift vertically so they cross and tangle
                var yDrift = _wild ? Math.sin(t * rb.freq * 0.4 + r * 1.1) * 6 : 0;
                rb.mesh.position.y = ribbonY[r] + yDrift;
                // PlaneGeometry vertices: (SEG_W+1)*2 total; top row is every other vertex
                for (var i = 0; i < count; i++) {
                    var x = pos.getX(i);
                    // Only animate top row (y > 0 in local space) for the waving top edge
                    if (pos.getY(i) > 0) {
                        pos.setY(i, rb.amp * ampMult * Math.sin(t * rb.freq + x * 0.08 + rb.phase)
                                    + rb.amp * ampMult * 0.4 * Math.cos(t * rb.freq * 0.7 + x * 0.05));
                    }
                }
                pos.needsUpdate = true;
            }
            renderer.render(scene, camera);
        }

        initThree();

        // Aurora nav/wild effects — snow burst on navigate, blizzard storm on wild
        window.snonuxOpenEffect = function() {
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-zoom'); setTimeout(function() { modal.classList.remove('sno-modal-zoom'); }, 400); }
            // Frost shimmer — aurora-colored radial
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:997;pointer-events:none;background:radial-gradient(ellipse at center,rgba(0,255,179,0.14) 0%,rgba(192,132,252,0.1) 55%,transparent 80%);transition:opacity 0.3s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 340); }, 15);
        };
        window.snonuxCloseEffect = function() {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(0,207,232,0.12);transition:opacity 0.2s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 230); }, 15);
        };
        window.snonuxScrollEffect = function(dir) {
            var isDown = dir === 'down';
            var thick = _wild ? '14px' : '5px';
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;left:0;right:0;height:' + thick + ';z-index:9000;pointer-events:none;' +
                'background:linear-gradient(90deg,transparent,rgba(0,207,232,0.9),rgba(120,200,100,0.9),rgba(0,207,232,0.9),transparent);' +
                (isDown ? 'top:0;' : 'bottom:0;') +
                'transition:transform 0.32s ease,opacity 0.32s ease;';
            document.body.appendChild(d);
            setTimeout(function() { d.style.transform = isDown ? 'translateY(100vh)' : 'translateY(-100vh)'; d.style.opacity='0'; }, 16);
            setTimeout(function() { d.remove(); }, 400);
        };
        window.snonuxWildToggle = function() {
            _wild = !_wild;
            var b = document.getElementById('sno-wild-badge');
            if (b) b.classList.toggle('sno-wild-on', _wild);
        };
        window.snonuxNavEffect = function() {
            // Snow burst — CSS snowflakes scatter from cursor
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-shake'); setTimeout(function() { ov.classList.remove('sno-fx-shake'); }, 380); }
            // Frost flash
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:radial-gradient(ellipse at center,rgba(0,255,179,0.18) 0%,rgba(192,132,252,0.1) 60%,transparent 100%);transition:opacity 0.22s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 250); }, 30);
        };
        window.snonuxPageEffect = function() {
            var ov = document.querySelector('.overlay');
            if (ov) { ov.classList.add('sno-fx-zoom'); setTimeout(function() { ov.classList.remove('sno-fx-zoom'); }, 330); }
        };
    })();
