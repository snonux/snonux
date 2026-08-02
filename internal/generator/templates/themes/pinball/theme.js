(function(){
    'use strict';
    var reduced=window.matchMedia&&window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    function scene(canvas,splash){
        if(!canvas||typeof THREE==='undefined')return function(){};
        var renderer=new THREE.WebGLRenderer({canvas:canvas,antialias:true,alpha:true}),world=new THREE.Scene();
        var camera=new THREE.PerspectiveCamera(52,1,.1,100),raf=0,wild=false,last=performance.now(),balls=[],bumpers=[],lights=[];
        camera.position.set(0,splash?1:8,splash?13:25); renderer.setPixelRatio(Math.min(devicePixelRatio||1,2));
        world.add(new THREE.AmbientLight(0xffffff,.75)); var light=new THREE.PointLight(0x20e7f2,3,60);light.position.set(-6,8,12);world.add(light);
        var railMat=new THREE.MeshPhongMaterial({color:0xc9dade,shininess:120});
        [-8,8].forEach(function(x){var rail=new THREE.Mesh(new THREE.TorusGeometry(10,.13,10,50,Math.PI),railMat);rail.position.x=x;rail.rotation.z=x<0?-.65:.65;world.add(rail);});
        [[-5,3,0xff3d4e],[0,-1,0xffd43b],[5,4,0x20e7f2]].forEach(function(b,i){var bumper=new THREE.Mesh(new THREE.CylinderGeometry(1.45,1.8,.8,24),new THREE.MeshPhongMaterial({color:b[2],emissive:b[2],emissiveIntensity:.35,shininess:70}));bumper.position.set(b[0],b[1],-2);bumper.rotation.x=Math.PI/2;bumper.userData={x:b[0],y:b[1],phase:i*2.1};bumpers.push(bumper);world.add(bumper);var glow=new THREE.PointLight(b[2],0,18);glow.position.set(b[0],b[1],3);glow.userData.phase=i*2.1;lights.push(glow);world.add(glow);});
        function addBall(i){var ball=new THREE.Mesh(new THREE.SphereGeometry(splash ? .42 : .55,20,16),railMat);ball.position.set((Math.random()-.5)*15,(Math.random()-.5)*18,i%3-1);ball.userData={vx:(Math.random()-.5)*5,vy:2+Math.random()*4};balls.push(ball);world.add(ball);}
        var normalBalls=splash?3:6,maxBalls=splash?12:22;
        for(var i=0;i<maxBalls;i++){addBall(i);balls[i].visible=i<normalBalls;}
        function resize(){var w=canvas.clientWidth||2,h=canvas.clientHeight||2;renderer.setSize(w,h,false);camera.aspect=w/h;camera.updateProjectionMatrix();}
        function draw(now){var dt=Math.min((now-last)/1000,.05),t=now/1000;last=now;balls.forEach(function(ball,i){if(!ball.visible)return;var energy=wild?3.8:1;ball.position.x+=ball.userData.vx*dt*energy;ball.position.y+=ball.userData.vy*dt*energy;if(Math.abs(ball.position.x)>10){ball.position.x=Math.sign(ball.position.x)*10;ball.userData.vx*=-1;}if(Math.abs(ball.position.y)>11){ball.position.y=Math.sign(ball.position.y)*11;ball.userData.vy*=-1;}if(wild)ball.position.z=-1+Math.sin(t*7+i)*2.2;ball.rotation.x+=dt*(wild?14:4);});if(wild&&!reduced){bumpers.forEach(function(b){var pulse=1+Math.sin(t*8+b.userData.phase)*.22;b.scale.set(pulse,pulse,1);b.position.x=b.userData.x+Math.sin(t*3+b.userData.phase)*.55;b.position.y=b.userData.y+Math.cos(t*4+b.userData.phase)*.4;b.material.emissiveIntensity=.55+Math.sin(t*10+b.userData.phase)*.35;});lights.forEach(function(l){l.intensity=2.5+Math.sin(t*11+l.userData.phase)*2;});camera.position.x=Math.sin(t*2.7)*1.4;camera.position.y=(splash?1:8)+Math.cos(t*3.3)*.8;camera.rotation.z=Math.sin(t*2.2)*.035;}renderer.render(world,camera);if(!reduced)raf=requestAnimationFrame(draw);}
        resize();window.addEventListener('resize',resize);draw(last);
        return function(state){if(typeof state==='boolean'){wild=state;balls.forEach(function(ball,i){ball.visible=i<(wild?maxBalls:normalBalls);if(wild&&i>=normalBalls){ball.position.set((Math.random()-.5)*15,(Math.random()-.5)*18,i%4-1);}});if(!wild){bumpers.forEach(function(b){b.scale.set(1,1,1);b.position.set(b.userData.x,b.userData.y,-2);b.material.emissiveIntensity=.35;});lights.forEach(function(l){l.intensity=0;});camera.position.set(0,splash?1:8,splash?13:25);camera.rotation.z=0;}if(reduced)renderer.render(world,camera);return;}window.removeEventListener('resize',resize);if(raf)cancelAnimationFrame(raf);renderer.dispose();};
    }
    if(!document.documentElement.classList.contains('sno-splash-skip'))window._snonuxSplashWebGLCleanup=scene(document.getElementById('splash-gl-canvas'),true);
    var setWild=scene(document.getElementById('three-canvas'),false),wild=false;
    function hit(color,x,y){if(reduced)return;var d=document.createElement('i');d.setAttribute('aria-hidden','true');d.style.cssText='position:fixed;z-index:9000;pointer-events:none;left:'+(x||innerWidth/2)+'px;top:'+(y||innerHeight/2)+'px;width:18px;height:18px;border:4px solid '+color+';border-radius:50%;box-shadow:0 0 18px '+color+';transform:translate(-50%,-50%);transition:transform .28s,opacity .28s';document.body.appendChild(d);requestAnimationFrame(function(){d.style.transform='translate(-50%,-50%) scale(8)';d.style.opacity='0';});setTimeout(function(){d.remove();},330);}
    window.snonuxOpenEffect=function(post){var r=post&&post.getBoundingClientRect();hit('#ffd43b',r?r.left+r.width/2:0,r?r.top+30:0);};
    window.snonuxCloseEffect=function(){hit('#ff3d4e');};
    window.snonuxNavEffect=function(){hit('#20e7f2',innerWidth*(.25+Math.random()*.5),innerHeight*(.2+Math.random()*.6));};
    window.snonuxPageEffect=function(){hit('#ffd43b');};
    window.snonuxScrollEffect=function(){var p=document.querySelector('.post-active');if(p){p.classList.add('pinball-hit');setTimeout(function(){p.classList.remove('pinball-hit');},180);}hit('#20e7f2',innerWidth*.12,innerHeight*.5);};
    window.snonuxWildToggle=function(){wild=!wild;setWild(wild);document.documentElement.classList.toggle('sno-pinball-multiball',wild);var b=document.getElementById('sno-wild-badge');if(b)b.classList.toggle('sno-wild-on',wild);};
})();
