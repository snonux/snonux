package generator

// neonTemplate is the cyberpunk neon theme — dark background, Three.js 3D orb
// and rings, cyan/magenta/yellow palette, Orbitron font.
// Keyboard nav and modal HTML are injected via the shared navDefs sub-templates.
const neonTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <title>snonux.foo • NEON NEXUS</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/three.js/r134/three.min.js"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css">
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;500;700&display=swap');
        :root { --neon-cyan:#00f5ff; --neon-magenta:#ff00cc; --neon-yellow:#ffe700; }
        * { margin:0; padding:0; box-sizing:border-box; }
        body { font-family:'Orbitron',sans-serif; background:#0b001a; color:#e0f8ff; overflow:hidden; height:100vh; }
        #three-canvas { position:fixed; top:0; left:0; width:100%; height:100%; z-index:1; }
        .overlay { position:relative; z-index:10; height:100vh; display:flex; flex-direction:column; }
        header { padding:16px 30px; display:flex; align-items:center; justify-content:space-between;
                 background:rgba(11,0,26,0.8); backdrop-filter:blur(12px);
                 border-bottom:2px solid rgba(255,231,0,0.3); }
        .logo { display:flex; align-items:center; gap:12px; }
        #sn-logo { flex-shrink:0; }
        .logo-title h1 { font-size:2rem; font-weight:700; letter-spacing:-3px; text-shadow:0 0 25px var(--neon-cyan); }
        .logo-title .subtitle { font-size:0.68rem; opacity:0.6; letter-spacing:1px; margin-top:2px; }
        .logo-title .subtitle a { color:var(--neon-cyan); text-decoration:none; }
        .logo-title .subtitle a:hover { text-shadow:0 0 8px var(--neon-cyan); }
        .nav { display:flex; gap:16px; align-items:center; }
        .transmit-btn { background:transparent; border:3px solid var(--neon-yellow); color:var(--neon-yellow);
                        padding:12px 28px; border-radius:9999px; font-weight:600; letter-spacing:1px;
                        display:flex; align-items:center; gap:10px; box-shadow:0 0 30px var(--neon-yellow);
                        transition:all 0.3s; text-decoration:none; font-family:'Orbitron',sans-serif; font-size:0.9rem; }
        .transmit-btn:hover { background:var(--neon-yellow); color:#0b001a; transform:scale(1.08); }
        .content { flex:1; padding:30px; overflow-y:auto; scrollbar-width:thin; scrollbar-color:#ffe700 #1a0033; }
        .page-nav { display:flex; justify-content:center; margin:18px 0; }
        .page-nav a { background:transparent; border:2px solid var(--neon-cyan); color:var(--neon-cyan);
                      padding:10px 28px; border-radius:9999px; font-size:0.85rem; letter-spacing:2px;
                      text-decoration:none; transition:all 0.3s; }
        .page-nav a:hover { background:var(--neon-cyan); color:#0b001a; }
        .post { background:rgba(20,5,45,0.9); border:2px solid transparent;
                border-image:linear-gradient(45deg,var(--neon-cyan),var(--neon-magenta)) 1;
                border-radius:24px; padding:28px; margin-bottom:28px;
                box-shadow:0 0 35px rgba(0,245,255,0.5);
                transition:all 0.4s cubic-bezier(0.23,1,0.32,1); cursor:pointer; }
        .post:hover { transform:translateY(-8px) rotate(1deg); box-shadow:0 0 50px rgba(255,231,0,0.6); }
        .post-active { border-image:none !important; border-color:var(--neon-yellow) !important;
                       background:rgba(40,20,70,0.97) !important;
                       box-shadow:0 0 0 2px var(--neon-yellow),0 0 30px rgba(255,231,0,0.7),
                                  0 0 70px rgba(255,231,0,0.35),inset 4px 0 0 var(--neon-yellow) !important;
                       transform:translateY(-6px) scale(1.012); }
        .post-header { display:flex; justify-content:space-between; margin-bottom:18px; font-size:0.95rem; }
        .post-time { font-family:monospace; color:var(--neon-yellow); text-shadow:0 0 12px var(--neon-yellow); }
        .post-text { font-size:1.1rem; line-height:1.55; }
        .post-text a { color:var(--neon-cyan); text-decoration:none; }
        .post-text a:hover { text-shadow:0 0 8px var(--neon-cyan); }
        .post-image { max-width:100%; border-radius:12px; margin-top:12px; }
        .post-audio { width:100%; margin-top:12px; }
        .nav-hints { display:flex; gap:20px; justify-content:center; align-items:center;
                     padding:6px 20px; background:rgba(11,0,26,0.7);
                     border-bottom:1px solid rgba(0,245,255,0.15);
                     font-size:0.68rem; letter-spacing:1.5px; color:rgba(224,248,255,0.5); flex-wrap:wrap; }
        .nav-hints kbd { display:inline-block; background:rgba(0,245,255,0.1);
                         border:1px solid rgba(0,245,255,0.35); border-radius:4px; padding:1px 5px;
                         color:var(--neon-cyan); font-family:monospace; font-size:0.72rem; margin:0 2px; }
        .post-modal { display:none; position:fixed; inset:0; z-index:100;
                      background:rgba(11,0,26,0.95); backdrop-filter:blur(16px);
                      overflow-y:auto; padding:40px; }
        .post-modal.active { display:block; }
        .modal-inner { max-width:800px; margin:0 auto; background:rgba(20,5,45,0.98);
                       border:2px solid transparent;
                       border-image:linear-gradient(45deg,var(--neon-yellow),var(--neon-magenta)) 1;
                       border-radius:24px; padding:40px; box-shadow:0 0 80px rgba(255,231,0,0.4); }
        .modal-close { float:right; background:none; border:none; color:var(--neon-cyan);
                       font-size:1.4rem; cursor:pointer; font-family:'Orbitron',sans-serif; }
        @media(max-width:640px) {
            .logo-title h1 { font-size:1.6rem; } #sn-logo { width:44px; height:44px; }
            .post { padding:22px; margin-bottom:22px; } .content { padding:20px; }
            header { padding:14px 20px; } .transmit-btn { padding:9px 16px; font-size:0.8rem; }
            .nav-hints { display:none; } .modal-inner { padding:24px 16px; }
        }
    </style>
</head>
<body>
    <canvas id="three-canvas"></canvas>
    <div class="overlay">
        <header>
            <div class="logo">
                <svg id="sn-logo" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 56 56" width="56" height="56" aria-label="snonux logo">
                  <defs>
                    <linearGradient id="sn-grad" x1="0" y1="0" x2="1" y2="1">
                      <stop offset="0%" stop-color="#ffe700"/><stop offset="100%" stop-color="#ff00cc"/>
                    </linearGradient>
                    <radialGradient id="sn-bg" cx="40%" cy="35%" r="70%">
                      <stop offset="0%" stop-color="#2d0060"/><stop offset="100%" stop-color="#0b001a"/>
                    </radialGradient>
                    <filter id="sn-gc" x="-60%" y="-60%" width="220%" height="220%">
                      <feGaussianBlur stdDeviation="2.5" result="b"/>
                      <feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge>
                    </filter>
                    <filter id="sn-gm" x="-60%" y="-60%" width="220%" height="220%">
                      <feGaussianBlur stdDeviation="2.5" result="b"/>
                      <feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge>
                    </filter>
                    <filter id="sn-gh" x="-20%" y="-20%" width="140%" height="140%">
                      <feGaussianBlur stdDeviation="3" result="b"/>
                      <feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge>
                    </filter>
                  </defs>
                  <polygon points="55,28 41.5,51.4 14.5,51.4 1,28 14.5,4.6 41.5,4.6"
                    fill="none" stroke="#ffe700" stroke-width="5" opacity="0.18" filter="url(#sn-gh)"/>
                  <polygon points="55,28 41.5,51.4 14.5,51.4 1,28 14.5,4.6 41.5,4.6"
                    fill="url(#sn-bg)" stroke="url(#sn-grad)" stroke-width="1.8"/>
                  <line x1="34" y1="12" x2="22" y2="44" stroke="#ffe700" stroke-width="0.9" opacity="0.75"/>
                  <rect x="32.5" y="10.5" width="3" height="3" transform="rotate(45 34 12)" fill="#ffe700" opacity="0.8"/>
                  <rect x="20.5" y="42.5" width="3" height="3" transform="rotate(45 22 44)" fill="#ffe700" opacity="0.8"/>
                  <text x="9" y="37" font-family="Orbitron,monospace" font-weight="700" font-size="20"
                    fill="#00f5ff" filter="url(#sn-gc)">S</text>
                  <text x="28" y="37" font-family="Orbitron,monospace" font-weight="700" font-size="20"
                    fill="#ff00cc" filter="url(#sn-gm)">N</text>
                </svg>
                <div class="logo-title">
                    <h1>snonux.foo</h1>
                    <p class="subtitle">microblog &mdash; <a href="https://foo.zone">foo.zone</a> is the real blog</p>
                </div>
            </div>
            <div class="nav">
                <a href="https://foo.zone/about" class="transmit-btn">
                    <i class="fa-solid fa-feather-pointed"></i> TRANSMIT TO NEXUS
                </a>
            </div>
        </header>
        {{template "navhints" .}}
        <div class="content" id="post-content">
            {{if .PrevPage}}
            <div class="page-nav"><a href="{{.PrevPage}}">&larr; NEWER TRANSMISSIONS</a></div>
            {{end}}
            {{range $i, $post := .Posts}}
            <div class="post" data-index="{{$i}}" onclick="selectPost({{$i}})">
                <div class="post-header">
                    <div><strong>@snonux</strong></div>
                    <div class="post-time">{{$post.FormattedTime}}</div>
                </div>
                <div class="post-text">{{$post.ContentHTML}}</div>
            </div>
            {{end}}
            {{if .NextPage}}
            <div class="page-nav"><a href="{{.NextPage}}">OLDER TRANSMISSIONS &rarr;</a></div>
            {{end}}
        </div>
    </div>
    {{template "navmodal" .}}
    <script>
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
                particles.rotation.y+=0.0008;
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
    </script>
    {{template "navscript" .}}
</body>
</html>`
