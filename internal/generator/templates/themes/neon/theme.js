
    (function(){
        if(document.documentElement.classList.contains('sno-splash-skip'))return;
        var cv=document.getElementById('splash-gl-canvas');
        if(!cv||typeof THREE==='undefined')return;
        var raf,ren,sc,ca,g=new THREE.Group(),t0=performance.now();
        function cleanup(){window.removeEventListener('resize',sz);if(raf)cancelAnimationFrame(raf);raf=null;if(ren){ren.dispose();}ren=null;window._snonuxSplashWebGLCleanup=null;}
        window._snonuxSplashWebGLCleanup=cleanup;
        function sz(){var w=cv.clientWidth||2,h=cv.clientHeight||2;if(ren)ren.setSize(w,h,false);if(ca){ca.aspect=w/h;ca.updateProjectionMatrix();}}
        ren=new THREE.WebGLRenderer({canvas:cv,antialias:true,alpha:true});
        ren.setClearColor(0,0);ren.setPixelRatio(Math.min(window.devicePixelRatio||1,2));
        sc=new THREE.Scene();ca=new THREE.PerspectiveCamera(52,1,0.1,100);ca.position.set(0,0.4,9);
        var cols=[0x00f5ff,0xff00cc,0xffe700],i,m;
        for(i=0;i<3;i++){m=new THREE.Mesh(new THREE.TorusGeometry(1.55+i*0.48,0.055,8,48),new THREE.MeshBasicMaterial({color:cols[i],transparent:true,opacity:0.92}));m.rotation.x=Math.PI/2;m.userData.sp=0.01+i*0.004;g.add(m);}
        g.add(new THREE.Mesh(new THREE.SphereGeometry(0.52,20,20),new THREE.MeshBasicMaterial({color:0xffe700,transparent:true,opacity:0.95})));
        sc.add(g);sz();window.addEventListener('resize',sz);
        function loop(now){raf=requestAnimationFrame(loop);var t=(now-t0)*0.001;g.rotation.y=t*0.42;g.rotation.x=Math.sin(t*0.65)*0.12;g.children.forEach(function(c){if(c.userData.sp)c.rotation.z+=c.userData.sp;});ren.render(sc,ca);}
        raf=requestAnimationFrame(loop);
    })();


        // Three.js neon nexus scene — central orb, orbiting rings, particle field.
        let scene, camera, renderer, centralSphere, rings = [], particles;
        function initThree() {
            const canvas = document.getElementById('three-canvas');
            renderer = new THREE.WebGLRenderer({ canvas, antialias:true, alpha:true });
            renderer.setSize(window.innerWidth, window.innerHeight);
            renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
            scene = new THREE.Scene();
            scene.fog = new THREE.Fog(0x0b001a, 15, 80);
            camera = new THREE.PerspectiveCamera(60, window.innerWidth/window.innerHeight, 0.1, 200);
            camera.position.set(0, 12, 35);
            scene.add(new THREE.AmbientLight(0x00f5ff, 0.8));
            const coreLight = new THREE.PointLight(0xff00cc, 4, 100);
            coreLight.position.set(0,0,0); scene.add(coreLight);
            centralSphere = new THREE.Mesh(new THREE.SphereGeometry(6,64,64),
                new THREE.MeshPhongMaterial({color:0x00f5ff,emissive:0xff00cc,emissiveIntensity:1.8,
                    shininess:100,transparent:true,opacity:0.95}));
            scene.add(centralSphere);
            scene.add(new THREE.Mesh(new THREE.SphereGeometry(4.5,64,64),
                new THREE.MeshBasicMaterial({color:0x00f5ff,transparent:true,opacity:0.4,blending:THREE.AdditiveBlending})));
            const rc=[0x00f5ff,0xff00cc,0x00f5ff,0xffe700];
            for(let i=0;i<14;i++){
                const ring=new THREE.Mesh(new THREE.TorusGeometry(12+i*2.2,0.35,32,128),
                    new THREE.MeshPhongMaterial({color:rc[i%4],emissive:rc[i%4],emissiveIntensity:2.5,
                        shininess:80,transparent:true,opacity:0.9,side:THREE.DoubleSide}));
                ring.rotation.x=Math.random()*Math.PI;
                ring.userData={speed:0.008+i*0.003,axisTilt:Math.random()*0.6};
                scene.add(ring); rings.push(ring);
            }
            const pCount=2200,pos=new Float32Array(pCount*3),col=new Float32Array(pCount*3);
            for(let i=0;i<pCount*3;i+=3){
                const r=30+Math.random()*40,t=Math.random()*Math.PI*2,p=Math.acos(2*Math.random()-1);
                pos[i]=r*Math.sin(p)*Math.cos(t);pos[i+1]=r*Math.sin(p)*Math.sin(t);pos[i+2]=r*Math.cos(p);
                const c=new THREE.Color().setHSL(Math.random()>0.5?0.55:0.8,1,1);
                col[i]=c.r;col[i+1]=c.g;col[i+2]=c.b;
            }
            const pg=new THREE.BufferGeometry();
            pg.setAttribute('position',new THREE.BufferAttribute(pos,3));
            pg.setAttribute('color',new THREE.BufferAttribute(col,3));
            particles=new THREE.Points(pg,new THREE.PointsMaterial(
                {size:0.22,vertexColors:true,transparent:true,opacity:0.9,blending:THREE.AdditiveBlending}));
            scene.add(particles);
            let mouseX=0;
            window.addEventListener('mousemove',e=>{mouseX=(e.clientX/window.innerWidth)*2-1;});
            (function animate(){
                requestAnimationFrame(animate);
                const time=Date.now()*0.0004;
                camera.position.x=Math.sin(time)*35+mouseX*6;
                camera.position.z=Math.cos(time)*35+10;
                camera.lookAt(0,4,0);
                centralSphere.rotation.y+=0.003;
                rings.forEach((ring,i)=>{
                    ring.rotation.y+=ring.userData.speed;
                    ring.rotation.x=Math.sin(time*1.5+i)*ring.userData.axisTilt;
                });
                particles.rotation.y+=window._snoNeonWild ? 0.012 : 0.0008;
                renderer.render(scene,camera);
            })();
        }
        window.addEventListener('resize',()=>{
            if(!camera||!renderer) return;
            camera.aspect=window.innerWidth/window.innerHeight;
            camera.updateProjectionMatrix();
            renderer.setSize(window.innerWidth,window.innerHeight);
        });
        window.onload=initThree;


    // Neon nav/wild effects — lightning flash on navigate, ring frenzy on wild
    (function() {
        function flash(color, ms) {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:' + color + ';transition:opacity ' + (ms||180) + 'ms';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity = '0'; setTimeout(function() { d.remove(); }, ms || 180); }, 30);
        }
        function fxOverlay(cls, ms) {
            var ov = document.querySelector('.overlay');
            if (!ov) return;
            ov.classList.add(cls);
            setTimeout(function() { ov.classList.remove(cls); }, ms || 380);
        }
        var _wild = false;
        window.snonuxOpenEffect = function(post) {
            // Modal burst from center with lightning ring
            var modal = document.getElementById('post-modal');
            if (modal) { modal.classList.add('sno-modal-expand'); setTimeout(function() { modal.classList.remove('sno-modal-expand'); }, 420); }
            // Cyan ring pulse radiating outward
            var ring = document.createElement('div');
            var r = post ? post.getBoundingClientRect() : {left: window.innerWidth/2, top: window.innerHeight/2};
            ring.style.cssText = 'position:fixed;top:' + (r.top+20) + 'px;left:' + (r.left+20) + 'px;z-index:997;pointer-events:none;width:10px;height:10px;border-radius:50%;border:3px solid rgba(0,245,255,0.9);transition:all 0.38s ease,opacity 0.38s';
            document.body.appendChild(ring);
            setTimeout(function() { ring.style.transform='scale(35)'; ring.style.opacity='0'; setTimeout(function() { ring.remove(); }, 420); }, 15);
        };
        window.snonuxCloseEffect = function() {
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;inset:0;z-index:998;pointer-events:none;background:rgba(255,0,204,0.12);transition:opacity 0.18s';
            document.body.appendChild(d);
            setTimeout(function() { d.style.opacity='0'; setTimeout(function() { d.remove(); }, 200); }, 15);
        };
        window.snonuxNavEffect = function() {
            flash('rgba(0,245,255,0.22)', 160);
            fxOverlay('sno-fx-shake', 380);
        };
        window.snonuxPageEffect = function() {
            flash('rgba(255,231,0,0.18)', 140);
            fxOverlay('sno-fx-zoom', 320);
        };
        window.snonuxScrollEffect = function(dir) {
            var isDown = dir === 'down';
            var thick = _wild ? '14px' : '5px';
            var d = document.createElement('div');
            d.style.cssText = 'position:fixed;left:0;right:0;height:' + thick + ';z-index:9000;pointer-events:none;' +
                'background:linear-gradient(90deg,transparent,rgba(0,245,255,0.9),rgba(255,0,204,0.9),transparent);' +
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
            // Speed up all rings and particles when wild
            if (rings && rings.length) {
                rings.forEach(function(r, i) {
                    r.userData.speed = _wild ? (0.008 + i * 0.003) * 14 : 0.008 + i * 0.003;
                });
            }
            // Store wild state for particle rotation boost in animate loop
            window._snoNeonWild = _wild;
        };
    })();
