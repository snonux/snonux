/*
 * Gravity Well — the cosmos theme's game.
 *
 * Asteroids on a wrapping screen where stars, black holes and pulsars pull on
 * everything: the ship, the rocks and every shot. Bullets curve, rocks
 * slingshot, and a ship or shot that touches a well's core is gone. Rocks are
 * too tough for that: they ricochet off a core in a flare, so gravity never
 * does the player's job. Destroy every rock to clear the level.
 */
(function () {
    'use strict';
    var G = window.SnoGame;
    // The playfield is the canvas below the HUD strip; it wraps on both axes.
    var TOP = G.HUD, PW = G.W, PH = G.H - G.HUD, TAU = Math.PI * 2;
    var THRUST = 270, TURN = 4.2, DRAG = 0.12, SHIP_R = 10, SHIP_MAX = 420;
    var SHOT_SPEED = 430, SHOT_LIFE = 1.15, MAX_SHOTS = 5, ROCK_MAX = 300;
    var ROCK_R = [0, 12, 22, 38], ROCK_PTS = [0, 100, 50, 20];
    var CORE_R = { star: 20, hole: 18, pulsar: 15, giant: 46 };
    var C = { bg: '#020214', gold: '#ffd166', purple: '#9b5de5', blue: '#4cc9f0', fg: '#d4e8ff', red: '#ff5d73', dim: '#6f7fa8' };
    var WELL_COLOR = { star: C.gold, hole: C.purple, pulsar: C.blue, giant: C.purple };

    // One entry per level. A well may orbit (cx, cy) at radius `or` with
    // angular speed `ow`, drift (dvx, dvy) or sweep a beam of length `beam`
    // at `bw` rad/s. `saucer` is the seconds between saucer visits. Spawn
    // points are deliberately off every axis of symmetry, so a ship left
    // alone always falls into something.
    var LEVELS = [
        { rocks: 3, spawns: [[150, 430], [810, 120]], hint: 'Stars pull on everything: thrust to stay out, and watch your shots bend',
            wells: [{ kind: 'star', cx: 480, cy: 285, gm: 0.9e6 }] },
        { rocks: 4, spawns: [[130, 440], [830, 110]], hint: 'Two suns: rocks slingshot from one to the other',
            wells: [{ kind: 'star', cx: 270, cy: 190, gm: 0.8e6 }, { kind: 'star', cx: 690, cy: 380, gm: 0.8e6 }] },
        { rocks: 4, spawns: [[120, 440], [840, 100]], hint: 'A binary pair of black holes waltzes around the centre',
            wells: [{ kind: 'hole', cx: 480, cy: 285, gm: 0.75e6, or: 105, ow: 0.5 },
                { kind: 'hole', cx: 480, cy: 285, gm: 0.75e6, or: 105, ow: 0.5, ph: Math.PI }] },
        { rocks: 5, saucer: 13, spawns: [[140, 440], [820, 110]], hint: 'Saucers ignore gravity. Their shots do not',
            wells: [{ kind: 'hole', cx: 480, cy: 285, gm: 1.3e6 }] },
        { rocks: 5, saucer: 12, spawns: [[100, 120], [860, 150]], hint: 'Three bodies, no stable orbit anywhere',
            wells: [{ kind: 'hole', cx: 480, cy: 160, gm: 1.0e6 }, { kind: 'star', cx: 250, cy: 400, gm: 0.7e6 },
                { kind: 'star', cx: 710, cy: 400, gm: 0.7e6 }] },
        { rocks: 5, iron: 2, saucer: 12, spawns: [[560, 90], [540, 480]], hint: 'Iron rocks (purple) take two hits each',
            wells: [{ kind: 'hole', cx: 300, cy: 285, gm: 0.6e6, or: 75, ow: 0.7 },
                { kind: 'hole', cx: 300, cy: 285, gm: 0.6e6, or: 75, ow: 0.7, ph: Math.PI },
                { kind: 'star', cx: 760, cy: 285, gm: 0.8e6 }] },
        { rocks: 5, iron: 1, saucer: 12, spawns: [[90, 440], [870, 110]], hint: 'The pulsar beam is deadly: stay out of its reach or slip in behind it',
            wells: [{ kind: 'pulsar', cx: 480, cy: 285, gm: 0.9e6, beam: 290, bw: 0.65 }] },
        { rocks: 6, iron: 2, saucer: 11, sniper: true, spawns: [[490, 80], [510, 490]], hint: 'Small saucers are snipers',
            wells: [{ kind: 'pulsar', cx: 250, cy: 285, gm: 0.8e6, beam: 250, bw: 0.7 },
                { kind: 'hole', cx: 720, cy: 285, gm: 0.55e6, or: 70, ow: -0.6 },
                { kind: 'hole', cx: 720, cy: 285, gm: 0.55e6, or: 70, ow: -0.6, ph: Math.PI }] },
        { rocks: 6, iron: 3, saucer: 10, maxFoes: 2, spawns: [[100, 450], [860, 100], [470, 60]], hint: 'A rogue black hole is drifting through',
            wells: [{ kind: 'pulsar', cx: 480, cy: 285, gm: 0.8e6, beam: 230, bw: -0.8 },
                { kind: 'hole', cx: 150, cy: 150, gm: 0.8e6, dvx: 34, dvy: 21 },
                { kind: 'star', cx: 800, cy: 440, gm: 0.6e6 }] },
        { rocks: 7, iron: 3, saucer: 10, sniper: true, maxFoes: 2, swirl: true,
            spawns: [[70, 80], [890, 80], [70, 490], [890, 490]], hint: 'Everything circles the giant. So should you',
            wells: [{ kind: 'giant', cx: 480, cy: 285, gm: 4.5e6 },
                { kind: 'pulsar', cx: 480, cy: 285, gm: 0.4e6, or: 250, ow: 0.2, beam: 160, bw: 1.0 }] }
    ];

    // ------------------------------------------------------------------
    // Wrapping space
    // ------------------------------------------------------------------

    function wrapX(x) { return ((x % PW) + PW) % PW; }
    function wrapY(y) { return TOP + ((((y - TOP) % PH) + PH) % PH); }
    // Shortest signed distance on the torus, so gravity and collisions work
    // across the screen edges exactly as they do in the middle.
    function dX(d) { return d > PW / 2 ? d - PW : (d < -PW / 2 ? d + PW : d); }
    function dY(d) { return d > PH / 2 ? d - PH : (d < -PH / 2 ? d + PH : d); }

    function gap(a, b) { return Math.hypot(dX(b.x - a.x), dY(b.y - a.y)); }
    function near(a, b, r) { return gap(a, b) < r; }

    function move(o, dt) {
        o.x = wrapX(o.x + o.vx * dt);
        o.y = wrapY(o.y + o.vy * dt);
    }

    function capSpeed(o, max) {
        var sp = Math.hypot(o.vx, o.vy);
        if (sp > max) { o.vx *= max / sp; o.vy *= max / sp; }
    }

    // Adds the pull of every well to o's velocity. The distance is floored so
    // the inverse square never blows up at a core (anything that close dies
    // on the same tick anyway).
    function pull(s, o, dt) {
        for (var i = 0; i < s.wells.length; i++) {
            var w = s.wells[i], dx = dX(w.x - o.x), dy = dY(w.y - o.y);
            var d2 = Math.max(dx * dx + dy * dy, 400), d = Math.sqrt(d2), a = w.gm / d2 * dt;
            o.vx += dx / d * a; o.vy += dy / d * a;
        }
    }

    // The well whose core (grown or shrunk by pad) contains o, if any.
    function coreAt(s, o, pad) {
        for (var i = 0; i < s.wells.length; i++) {
            if (near(o, s.wells[i], s.wells[i].r + pad)) return s.wells[i];
        }
        return null;
    }

    // Inside this radius gravity beats the ship's engine: only sideways speed
    // gets you out again. It is drawn as the dashed red ring.
    function dangerRadius(w) { return Math.sqrt(w.gm / THRUST); }

    // ------------------------------------------------------------------
    // Level setup
    // ------------------------------------------------------------------

    function placeWell(w) {
        w.x = w.cx + Math.cos(w.ph) * w.or;
        w.y = w.cy + Math.sin(w.ph) * w.or;
    }

    function makeWell(d) {
        var w = {
            kind: d.kind, gm: d.gm, r: CORE_R[d.kind],
            cx: d.cx, cy: d.cy, or: d.or || 0, ow: d.ow || 0, ph: d.ph || 0,
            dvx: d.dvx || 0, dvy: d.dvy || 0, beam: d.beam || 0, bw: d.bw || 0, ba: 0, x: 0, y: 0
        };
        placeWell(w);
        return w;
    }

    // The spawn point that is currently furthest from every well; wells move,
    // so the best one changes during a level.
    function bestSpawn(s) {
        var best = null, bestD = -1;
        s.spawns.forEach(function (p) {
            var spot = { x: p[0], y: p[1] }, d = Infinity;
            s.wells.forEach(function (w) { d = Math.min(d, gap(spot, w) - w.r - w.beam); });
            if (d > bestD) { bestD = d; best = spot; }
        });
        return best;
    }

    function newShip(s) {
        var spot = bestSpawn(s);
        return { x: spot.x, y: spot.y, vx: 0, vy: 0, a: -Math.PI / 2, inv: 2.5, wait: 0, flame: false };
    }

    function newRock(x, y, vx, vy, size, iron, rnd) {
        var shape = [];
        for (var i = 0; i < 10; i++) shape.push(0.74 + rnd() * 0.38);
        return {
            x: x, y: y, vx: vx, vy: vy, size: size, r: ROCK_R[size], iron: iron, hp: iron ? 2 : 1,
            a: rnd() * TAU, spin: (rnd() - 0.5) * 2.4, shape: shape
        };
    }

    // Rocks start on a rough orbit around the well that pulls hardest at
    // their position, so the level opens with slingshots instead of a rain
    // of rocks straight into the cores.
    function orbitVelocity(s, spot, rnd, swirl) {
        var best = s.wells[0], bestPull = 0;
        s.wells.forEach(function (w) {
            var d = Math.max(gap(spot, w), 1), p = w.gm / (d * d);
            if (p > bestPull) { bestPull = p; best = w; }
        });
        var dx = dX(best.x - spot.x), dy = dY(best.y - spot.y), d = Math.max(Math.hypot(dx, dy), 1);
        var sp = G.clamp(Math.sqrt(best.gm / d) * (0.8 + rnd() * 0.3), 40, 150);
        var dir = swirl || rnd() < 0.5 ? 1 : -1;
        return { vx: -dy / d * sp * dir, vy: dx / d * sp * dir };
    }

    function rockSpotOk(s, spot) {
        if (near(spot, { x: s.spawns[0][0], y: s.spawns[0][1] }, 230)) return false;
        return !s.wells.some(function (w) { return near(spot, w, w.r + 120); });
    }

    function placeRocks(s, L, rnd) {
        for (var i = 0; i < L.rocks; i++) {
            var spot = { x: 0, y: 0 };
            // A bounded search: the last candidate is used even if imperfect.
            for (var tries = 0; tries < 40; tries++) {
                spot.x = rnd() * PW; spot.y = TOP + rnd() * PH;
                if (rockSpotOk(s, spot)) break;
            }
            var v = orbitVelocity(s, spot, rnd, L.swirl);
            s.rocks.push(newRock(spot.x, spot.y, v.vx, v.vy, 3, i < (L.iron || 0), rnd));
        }
    }

    function makeStars(rnd) {
        var stars = [];
        for (var i = 0; i < 110; i++) {
            stars.push({ x: rnd() * PW, y: TOP + rnd() * PH, size: 1 + rnd() * 1.6, ph: rnd() * TAU, sp: 0.8 + rnd() * 2.4 });
        }
        return stars;
    }

    function init(level) {
        var L = LEVELS[level - 1], rnd = G.rng(level * 7919 + 13);
        var s = {
            time: 0, hint: L.hint, wells: L.wells.map(makeWell), spawns: L.spawns,
            rocks: [], shots: [], foes: [], foeShots: [], stars: makeStars(rnd),
            saucerEvery: L.saucer || 0, saucerT: (L.saucer || 0) * 0.6, sniper: !!L.sniper, maxFoes: L.maxFoes || 1,
            cool: 0, hyper: 0, tm: { thrust: 0, warn: 0, hum: 0, pulse: 0 }, ship: null
        };
        s.ship = newShip(s);
        placeRocks(s, L, rnd);
        return s;
    }

    // ------------------------------------------------------------------
    // Wells
    // ------------------------------------------------------------------

    function moveWells(s, dt) {
        s.wells.forEach(function (w) {
            if (w.dvx || w.dvy) { w.cx = wrapX(w.cx + w.dvx * dt); w.cy = wrapY(w.cy + w.dvy * dt); }
            w.ph += w.ow * dt;
            w.ba += w.bw * dt;
            placeWell(w);
        });
    }

    // Rocks survive a core: they ricochet off it elastically. If wells ate
    // rocks a level would clear itself while the player just kept away.
    function bounceRock(w, k) {
        var dx = dX(k.x - w.x), dy = dY(k.y - w.y), d = Math.max(Math.hypot(dx, dy), 1);
        var nx = dx / d, ny = dy / d, reach = w.r + k.r * 0.6;
        k.x = wrapX(w.x + nx * reach); k.y = wrapY(w.y + ny * reach);
        var vn = k.vx * nx + k.vy * ny;
        if (vn >= 0) return;
        k.vx -= 2 * vn * nx; k.vy -= 2 * vn * ny;
        G.burst(w.x + nx * w.r, w.y + ny * w.r, { n: 8, color: WELL_COLOR[w.kind], speed: 170, life: 0.4, size: 2 });
        G.tone(150 + k.size * 40, 0.16, { type: 'sine', slide: 420, vol: 0.12 });
    }

    function beamEnd(w) {
        return { x: w.x + Math.cos(w.ba) * w.beam, y: w.y + Math.sin(w.ba) * w.beam };
    }

    function beamHits(s, p) {
        return s.wells.some(function (w) {
            if (!w.beam) return false;
            var e = beamEnd(w), q = G.closestOnSeg(p.x, p.y, w.x, w.y, e.x, e.y);
            return G.dist(p.x, p.y, q.x, q.y) < SHIP_R;
        });
    }

    // ------------------------------------------------------------------
    // Ship
    // ------------------------------------------------------------------

    function thrust(s, p, dt) {
        p.vx += Math.cos(p.a) * THRUST * dt;
        p.vy += Math.sin(p.a) * THRUST * dt;
        if (s.tm.thrust > 0) return;
        // Gated: one puff of noise and exhaust every few ticks, not per tick.
        s.tm.thrust = 0.07;
        G.noise(0.1, { freq: 320, vol: 0.09 });
        G.burst(p.x - Math.cos(p.a) * 10, p.y - Math.sin(p.a) * 10, {
            n: 2, color: C.gold, speed: 150, life: 0.35, size: 2, angle: p.a + Math.PI, spread: 0.6
        });
    }

    function fire(s, p) {
        if (s.cool > 0 || s.shots.length >= MAX_SHOTS) return;
        s.cool = 0.17;
        var cx = Math.cos(p.a), cy = Math.sin(p.a);
        // Shots inherit the ship's velocity, as they would in free fall.
        s.shots.push({ x: p.x + cx * 14, y: p.y + cy * 14, vx: p.vx + cx * SHOT_SPEED, vy: p.vy + cy * SHOT_SPEED, life: SHOT_LIFE });
        G.tone(1040, 0.08, { slide: 320, vol: 0.09 });
    }

    // A random jump that keeps the ship's momentum. It never lands inside a
    // core, but it may well land next to one, or in front of a rock.
    function hyperspace(s, p) {
        if (s.hyper > 0) return;
        s.hyper = 3;
        G.burst(p.x, p.y, { n: 16, color: C.purple, speed: 200, life: 0.4 });
        var spot = { x: p.x, y: p.y };
        for (var tries = 0; tries < 20; tries++) {
            spot.x = G.rnd(20, PW - 20); spot.y = G.rnd(TOP + 20, G.H - 20);
            if (!coreAt(s, spot, 60)) break;
        }
        p.x = spot.x; p.y = spot.y;
        G.burst(p.x, p.y, { n: 16, color: C.blue, speed: 200, life: 0.4 });
        G.tone(200, 0.3, { type: 'sawtooth', slide: 2400, vol: 0.1 });
        G.noise(0.25, { filter: 'bandpass', freq: 400, slide: 3000, vol: 0.12 });
        G.flash(C.purple, 0.12);
    }

    function steerShip(s, dt) {
        var p = s.ship;
        if (p.wait > 0) {
            p.wait -= dt;
            if (p.wait <= 0) s.ship = newShip(s);
            return;
        }
        if (G.key.left) p.a -= TURN * dt;
        if (G.key.right) p.a += TURN * dt;
        p.flame = G.key.up;
        if (p.flame) thrust(s, p, dt);
        pull(s, p, dt);
        // A trace of drag keeps the ship controllable and means a drifting
        // ship always spirals inward in the end.
        p.vx -= p.vx * DRAG * dt; p.vy -= p.vy * DRAG * dt;
        capSpeed(p, SHIP_MAX);
        move(p, dt);
        if (p.inv > 0) p.inv -= dt;
        if (G.key.a) fire(s, p);
        if (G.hit.b) hyperspace(s, p);
    }

    function killShip(s) {
        var p = s.ship;
        G.burst(p.x, p.y, { n: 40, color: C.blue, speed: 260, life: 0.9 });
        G.burst(p.x, p.y, { n: 20, color: C.gold, speed: 140, life: 0.7 });
        G.sfx('bigboom');
        // The wreck is cleared away for a moment, then a new ship spawns with
        // a short shield; shots in flight are removed so it is a fair start.
        p.wait = 1.0; p.flame = false;
        s.foeShots = [];
        G.loseLife();
    }

    function warnNearWell(s, p) {
        if (s.tm.warn > 0) return;
        var close = s.wells.some(function (w) { return near(p, w, dangerRadius(w) * 1.25 + w.r); });
        if (!close) return;
        s.tm.warn = 0.32;
        G.tone(1320, 0.07, { type: 'triangle', vol: 0.1 });
    }

    function checkShip(s) {
        var p = s.ship;
        if (p.wait > 0) return;
        // A core kills even a freshly spawned, shielded ship.
        if (coreAt(s, p, SHIP_R * 0.5)) { killShip(s); return; }
        warnNearWell(s, p);
        if (p.inv > 0) return;
        var rock = s.rocks.some(function (k) { return near(p, k, k.r * 0.85 + SHIP_R * 0.7); });
        var foe = s.foes.some(function (f) { return near(p, f, f.r + SHIP_R * 0.7); });
        if (rock || foe || beamHits(s, p)) killShip(s);
    }

    // ------------------------------------------------------------------
    // Rocks and shots
    // ------------------------------------------------------------------

    function splitRock(s, k) {
        var sp = Math.hypot(k.vx, k.vy) || 1, nx = -k.vy / sp, ny = k.vx / sp, kick = G.rnd(45, 80);
        [1, -1].forEach(function (side) {
            s.rocks.push(newRock(k.x + nx * side * k.r * 0.4, k.y + ny * side * k.r * 0.4,
                k.vx + nx * side * kick, k.vy + ny * side * kick, k.size - 1, k.iron, Math.random));
        });
    }

    function breakRock(s, k) {
        s.rocks.splice(s.rocks.indexOf(k), 1);
        var pts = ROCK_PTS[k.size] * (k.iron ? 2 : 1);
        G.addScore(pts);
        G.popup(k.x, k.y, pts, C.gold);
        G.burst(k.x, k.y, { n: 6 + k.size * 5, color: k.iron ? C.purple : C.fg, speed: 150 + k.size * 30, life: 0.6 });
        // Bigger rocks rumble lower and longer.
        G.noise(0.14 + k.size * 0.1, { freq: 2200 - k.size * 550, slide: 80, vol: 0.2 + k.size * 0.06 });
        if (k.size === 3) G.shake(4, 0.15);
        if (k.size > 1) splitRock(s, k);
    }

    function hitRock(s, k, shot) {
        k.hp--;
        if (k.hp <= 0) { breakRock(s, k); return; }
        // Iron: the first hit only cracks it and shoves it a little.
        k.vx += shot.vx * 0.04; k.vy += shot.vy * 0.04;
        G.tone(1700, 0.09, { type: 'triangle', vol: 0.14 });
        G.burst(shot.x, shot.y, { n: 6, color: C.gold, speed: 160, life: 0.3, size: 2 });
    }

    function updateRocks(s, dt) {
        s.rocks.forEach(function (k) {
            pull(s, k, dt);
            capSpeed(k, ROCK_MAX);
            move(k, dt);
            k.a += k.spin * dt;
            var w = coreAt(s, k, k.r * 0.6);
            if (w) bounceRock(w, k);
        });
    }

    // True when the shot hit something and must be removed.
    function shotHits(s, b) {
        var i;
        for (i = 0; i < s.rocks.length; i++) {
            if (near(b, s.rocks[i], s.rocks[i].r + 2)) { hitRock(s, s.rocks[i], b); return true; }
        }
        for (i = 0; i < s.foes.length; i++) {
            if (near(b, s.foes[i], s.foes[i].r + 3)) { killFoe(s, s.foes[i]); return true; }
        }
        return false;
    }

    function fizzle(b, color) {
        G.burst(b.x, b.y, { n: 3, color: color, speed: 60, life: 0.25, size: 2 });
    }

    function updateShots(s, dt) {
        s.shots = s.shots.filter(function (b) {
            pull(s, b, dt);
            move(b, dt);
            b.life -= dt;
            if (coreAt(s, b, 0)) { fizzle(b, C.gold); return false; }
            return b.life > 0 && !shotHits(s, b);
        });
    }

    // ------------------------------------------------------------------
    // Saucers
    // ------------------------------------------------------------------

    function spawnFoe(s) {
        var dir = Math.random() < 0.5 ? 1 : -1, small = s.sniper && Math.random() < 0.6;
        s.foes.push({
            x: dir > 0 ? -20 : PW + 20, y: G.rnd(TOP + 70, G.H - 70), vx: dir * G.rnd(80, 120),
            ph: G.rnd(0, TAU), small: small, r: small ? 10 : 16, fireT: 1.3
        });
        G.sfx('alarm');
    }

    function foeFire(s, f) {
        var p = s.ship, err = f.small ? 0.1 : 0.32;
        var a = Math.atan2(dY(p.y - f.y), dX(p.x - f.x)) + G.rnd(-err, err);
        s.foeShots.push({ x: f.x, y: f.y, vx: Math.cos(a) * 230, vy: Math.sin(a) * 230, life: 2.4 });
        G.tone(520, 0.14, { type: 'sawtooth', slide: 140, vol: 0.1 });
    }

    // Saucers fly on anti-gravity: they cross the screen once, bobbing, and
    // swerve around any core in their lane. Returns false once off-screen.
    function moveFoe(s, f, dt) {
        var vy = Math.cos(f.ph * 1.6) * 50;
        s.wells.forEach(function (w) {
            if (near(f, w, w.r + 110)) vy += (f.y < w.y ? -1 : 1) * 110;
        });
        f.ph += dt;
        f.x += f.vx * dt;
        f.y = G.clamp(f.y + vy * dt, TOP + 20, G.H - 20);
        f.fireT -= dt;
        if (f.fireT <= 0 && s.ship.wait <= 0) { foeFire(s, f); f.fireT = f.small ? 1.5 : 2.1; }
        return f.x > -30 && f.x < PW + 30;
    }

    function killFoe(s, f) {
        s.foes.splice(s.foes.indexOf(f), 1);
        var pts = f.small ? 600 : 300;
        G.addScore(pts);
        G.popup(f.x, f.y, pts, C.red);
        G.burst(f.x, f.y, { n: 26, color: C.red, speed: 220, life: 0.7 });
        G.sfx('boom');
    }

    function updateFoes(s, dt) {
        if (s.saucerEvery) {
            s.saucerT -= dt;
            if (s.saucerT <= 0 && s.foes.length < s.maxFoes) { spawnFoe(s); s.saucerT = s.saucerEvery; }
        }
        s.foes = s.foes.filter(function (f) { return moveFoe(s, f, dt); });
        if (s.foes.length && s.tm.hum <= 0) {
            s.tm.hum = 0.28;
            G.tone(s.foes[0].small ? 640 : 380, 0.12, { type: 'sine', slide: s.foes[0].small ? 760 : 460, vol: 0.06 });
        }
    }

    function updateFoeShots(s, dt) {
        var p = s.ship, hit = false;
        s.foeShots = s.foeShots.filter(function (b) {
            pull(s, b, dt);
            move(b, dt);
            b.life -= dt;
            if (coreAt(s, b, 0)) { fizzle(b, C.red); return false; }
            if (p.wait <= 0 && p.inv <= 0 && near(b, p, SHIP_R)) { hit = true; return false; }
            return b.life > 0;
        });
        if (hit) killShip(s);
    }

    // ------------------------------------------------------------------
    // Update
    // ------------------------------------------------------------------

    function tickTimers(s, dt) {
        for (var k in s.tm) if (s.tm[k] > 0) s.tm[k] -= dt;
        if (s.cool > 0) s.cool -= dt;
        if (s.hyper > 0) s.hyper -= dt;
        // The pulsar ticks like a lighthouse so its beam can be heard coming.
        var pulsar = s.wells.some(function (w) { return w.beam > 0; });
        if (pulsar && s.tm.pulse <= 0) {
            s.tm.pulse = 0.9;
            G.tone(1760, 0.03, { type: 'sine', vol: 0.05 });
        }
    }

    function update(s, dt) {
        s.time += dt;
        tickTimers(s, dt);
        moveWells(s, dt);
        steerShip(s, dt);
        updateShots(s, dt);
        updateRocks(s, dt);
        updateFoes(s, dt);
        updateFoeShots(s, dt);
        checkShip(s);
        if (!s.rocks.length) G.win(500 + G.lives * 250);
    }

    // ------------------------------------------------------------------
    // Drawing
    // ------------------------------------------------------------------

    // Calls fn at (x, y) and at its mirror images across any edge it
    // overlaps, so things slide across the wrap instead of popping.
    function eachImage(x, y, r, fn) {
        var xs = [x], ys = [y];
        if (x < r) xs.push(x + PW); else if (x > PW - r) xs.push(x - PW);
        if (y < TOP + r) ys.push(y + PH); else if (y > G.H - r) ys.push(y - PH);
        xs.forEach(function (ix) { ys.forEach(function (iy) { fn(ix, iy); }); });
    }

    function circle(ctx, x, y, r) {
        ctx.beginPath();
        ctx.arc(x, y, Math.max(0.1, r), 0, TAU);
    }

    function drawSky(s, ctx) {
        ctx.fillStyle = C.bg;
        ctx.fillRect(0, 0, G.W, G.H);
        var neb = ctx.createRadialGradient(700, 150, 20, 700, 150, 520);
        neb.addColorStop(0, 'rgba(155,93,229,0.20)');
        neb.addColorStop(0.5, 'rgba(76,201,240,0.06)');
        neb.addColorStop(1, 'rgba(2,2,20,0)');
        ctx.fillStyle = neb;
        ctx.fillRect(0, TOP, G.W, PH);
        ctx.fillStyle = C.fg;
        s.stars.forEach(function (st) {
            ctx.globalAlpha = 0.3 + 0.3 * Math.sin(s.time * st.sp + st.ph);
            ctx.fillRect(st.x, st.y, st.size, st.size);
        });
        ctx.globalAlpha = 1;
    }

    // Rings that shrink toward the core show the inflow, and the dashed red
    // ring marks where gravity outpulls the engine.
    function drawField(s, ctx, w, x, y) {
        var reach = dangerRadius(w) * 1.7;
        ctx.strokeStyle = WELL_COLOR[w.kind];
        ctx.lineWidth = 1;
        for (var i = 0; i < 3; i++) {
            var f = 1 - ((s.time * 0.22 + i / 3) % 1);
            ctx.globalAlpha = 0.34 * (1 - f);
            circle(ctx, x, y, w.r + reach * f);
            ctx.stroke();
        }
        ctx.globalAlpha = 0.5;
        ctx.strokeStyle = C.red;
        ctx.setLineDash([4, 8]);
        ctx.lineDashOffset = s.time * 14;
        circle(ctx, x, y, dangerRadius(w) + w.r);
        ctx.stroke();
        ctx.setLineDash([]);
        ctx.globalAlpha = 1;
    }

    function drawGlow(ctx, x, y, r, color) {
        var g = ctx.createRadialGradient(x, y, r * 0.2, x, y, r);
        g.addColorStop(0, color);
        g.addColorStop(1, 'rgba(2,2,20,0)');
        ctx.fillStyle = g;
        circle(ctx, x, y, r);
        ctx.fill();
    }

    // A black disc inside a glowing rim, with tilted accretion arcs that
    // rotate so the hole never looks static.
    function drawHole(s, ctx, w, x, y) {
        ctx.globalAlpha = 0.55;
        drawGlow(ctx, x, y, w.r * 2.6, C.purple);
        ctx.globalAlpha = 1;
        ctx.fillStyle = '#000';
        circle(ctx, x, y, w.r);
        ctx.fill();
        ctx.lineWidth = 2;
        ctx.strokeStyle = C.purple;
        ctx.stroke();
        var rings = w.kind === 'giant' ? 3 : 1;
        for (var i = 0; i < rings; i++) {
            var a = s.time * (1.4 - i * 0.3) + i * 2;
            ctx.strokeStyle = i % 2 ? C.gold : C.blue;
            ctx.beginPath();
            ctx.ellipse(x, y, w.r * (1.55 + i * 0.4), w.r * (0.5 + i * 0.14), 0.35, a, a + 4.2);
            ctx.stroke();
        }
    }

    function drawSun(s, ctx, w, x, y) {
        var beat = 1 + 0.08 * Math.sin(s.time * 3 + w.cx);
        var pulsar = w.kind === 'pulsar';
        ctx.globalAlpha = 0.6;
        drawGlow(ctx, x, y, w.r * 2.8 * beat, pulsar ? C.blue : C.gold);
        ctx.globalAlpha = 1;
        ctx.fillStyle = pulsar ? '#eaffff' : '#fff3c4';
        circle(ctx, x, y, w.r * (pulsar ? 0.6 : 0.8) * beat);
        ctx.fill();
        ctx.strokeStyle = pulsar ? C.blue : C.gold;
        ctx.lineWidth = 2;
        circle(ctx, x, y, w.r);
        ctx.stroke();
    }

    function drawWell(s, ctx, w) {
        eachImage(w.x, w.y, dangerRadius(w) * 1.7 + w.r, function (x, y) {
            drawField(s, ctx, w, x, y);
            if (w.kind === 'hole' || w.kind === 'giant') drawHole(s, ctx, w, x, y);
            else drawSun(s, ctx, w, x, y);
        });
    }

    function drawBeam(ctx, w) {
        if (!w.beam) return;
        var e = beamEnd(w);
        ctx.lineCap = 'round';
        [[16, 0.16, C.blue], [7, 0.4, C.blue], [2.5, 1, '#ffffff']].forEach(function (layer) {
            ctx.globalAlpha = layer[1];
            ctx.strokeStyle = layer[2];
            ctx.lineWidth = layer[0];
            ctx.beginPath();
            ctx.moveTo(w.x, w.y);
            ctx.lineTo(e.x, e.y);
            ctx.stroke();
        });
        ctx.globalAlpha = 1;
    }

    function drawRock(ctx, k) {
        ctx.fillStyle = k.iron ? '#2a1746' : '#141a3c';
        // A cracked iron rock turns gold: one more hit breaks it.
        ctx.strokeStyle = k.iron ? (k.hp > 1 ? C.purple : C.gold) : C.fg;
        ctx.lineWidth = k.iron ? 3 : 2;
        eachImage(k.x, k.y, k.r * 1.2, function (x, y) {
            ctx.beginPath();
            for (var i = 0; i < k.shape.length; i++) {
                var a = k.a + i / k.shape.length * TAU, r = k.r * k.shape[i];
                ctx.lineTo(x + Math.cos(a) * r, y + Math.sin(a) * r);
            }
            ctx.closePath();
            ctx.fill();
            ctx.stroke();
        });
    }

    function drawShots(s, ctx) {
        ctx.fillStyle = C.gold;
        s.shots.forEach(function (b) { circle(ctx, b.x, b.y, 2.6); ctx.fill(); });
        ctx.fillStyle = C.red;
        s.foeShots.forEach(function (b) { circle(ctx, b.x, b.y, 3.4); ctx.fill(); });
    }

    function drawFoe(s, ctx, f) {
        var r = f.r;
        ctx.lineWidth = 2;
        ctx.strokeStyle = C.red;
        ctx.fillStyle = '#2b0d1c';
        ctx.beginPath();
        ctx.ellipse(f.x, f.y, r * 1.5, r * 0.55, 0, 0, TAU);
        ctx.fill();
        ctx.stroke();
        ctx.beginPath();
        ctx.arc(f.x, f.y - r * 0.3, r * 0.7, Math.PI, 0);
        ctx.stroke();
        // Running lights chase along the rim.
        for (var i = -1; i <= 1; i++) {
            ctx.fillStyle = Math.floor(s.time * 6 + i + 3) % 3 === 0 ? C.gold : C.red;
            ctx.fillRect(f.x + i * r * 0.8 - 1.5, f.y - 1.5, 3, 3);
        }
    }

    // Dots along the path a shot fired right now would take: the only way to
    // read how the wells will bend it.
    function drawAim(s, ctx, p) {
        var cx = Math.cos(p.a), cy = Math.sin(p.a);
        var o = { x: p.x + cx * 14, y: p.y + cy * 14, vx: p.vx + cx * SHOT_SPEED, vy: p.vy + cy * SHOT_SPEED };
        ctx.fillStyle = C.blue;
        for (var i = 1; i <= 44; i++) {
            pull(s, o, G.STEP);
            move(o, G.STEP);
            if (coreAt(s, o, 0)) break;
            if (i % 4) continue;
            ctx.globalAlpha = 0.85 * (1 - i / 52);
            ctx.fillRect(o.x - 1.5, o.y - 1.5, 3, 3);
        }
        ctx.globalAlpha = 1;
    }

    function shipPath(ctx, x, y, a) {
        var cx = Math.cos(a), cy = Math.sin(a);
        ctx.beginPath();
        ctx.moveTo(x + cx * 14, y + cy * 14);
        ctx.lineTo(x - cx * 9 - cy * 9, y - cy * 9 + cx * 9);
        ctx.lineTo(x - cx * 5, y - cy * 5);
        ctx.lineTo(x - cx * 9 + cy * 9, y - cy * 9 - cx * 9);
        ctx.closePath();
    }

    function drawShip(s, ctx) {
        var p = s.ship;
        if (p.wait > 0) return;
        drawAim(s, ctx, p);
        eachImage(p.x, p.y, 18, function (x, y) {
            if (p.flame) {
                var len = 16 + Math.random() * 10, cx = Math.cos(p.a), cy = Math.sin(p.a);
                ctx.strokeStyle = C.gold;
                ctx.lineWidth = 3;
                ctx.beginPath();
                ctx.moveTo(x - cx * 7, y - cy * 7);
                ctx.lineTo(x - cx * len, y - cy * len);
                ctx.stroke();
            }
            ctx.fillStyle = '#0b1638';
            ctx.strokeStyle = C.blue;
            ctx.lineWidth = 2;
            shipPath(ctx, x, y, p.a);
            ctx.fill();
            ctx.stroke();
            if (p.inv <= 0) return;
            // The spawn shield pulses and fades as it runs out.
            ctx.globalAlpha = Math.min(1, p.inv) * (0.5 + 0.4 * Math.sin(s.time * 18));
            ctx.strokeStyle = C.fg;
            circle(ctx, x, y, 19);
            ctx.stroke();
            ctx.globalAlpha = 1;
        });
    }

    function draw(s, ctx) {
        drawSky(s, ctx);
        // Mirror images near the top edge must not paint under the HUD.
        ctx.save();
        ctx.beginPath();
        ctx.rect(0, TOP, G.W, PH);
        ctx.clip();
        s.wells.forEach(function (w) { drawWell(s, ctx, w); });
        s.wells.forEach(function (w) { drawBeam(ctx, w); });
        s.rocks.forEach(function (k) { drawRock(ctx, k); });
        drawShots(s, ctx);
        s.foes.forEach(function (f) { drawFoe(s, ctx, f); });
        drawShip(s, ctx);
        ctx.restore();
        // s.time is 0 until the level is actually played, which keeps the hint
        // off the title screen where it would sit on top of the menu text.
        if (s.time > 0 && s.time < 6) {
            ctx.globalAlpha = Math.min(1, 6 - s.time);
            G.text(s.hint, G.W / 2, G.H - 16, { size: 15, color: C.dim, align: 'center', max: G.W - 40 });
            ctx.globalAlpha = 1;
        }
    }

    function hud(s) {
        return 'ROCKS ' + s.rocks.length + '  JUMP ' + (s.hyper > 0 ? Math.ceil(s.hyper) + 's' : 'OK');
    }

    G.register('cosmos', {
        title: 'GRAVITY WELL',
        blurb: 'Destroy every rock. Everything falls, including your shots.',
        controls: [
            '← → rotate · ↑ thrust · SPACE fire',
            'X hyperspace: a random jump that keeps your momentum',
            'Never touch a core: it swallows ships and shots, but rocks bounce off',
            'The dotted arc shows where your next shot will curve'
        ],
        levelNames: ['First Light', 'Twin Suns', 'Binary Waltz', 'Visitors', 'Three-Body Problem',
            'Iron Belt', 'Lighthouse', 'Crossfire', 'Rogue Hole', 'Event Horizon'],
        colors: { bg: C.bg, fg: C.fg, accent: C.blue, dim: C.dim },
        lives: 3,
        // Slow D dorian: long held lead notes stand in for pads over a
        // rolling sine arpeggio and a sparse, half-time beat.
        music: {
            bpm: 96, root: 38, scale: 'dorian', prog: [0, 3, 6, 4],
            bass: 'x.......x.....o.',
            lead: ['4-----7-9---..7-', '8-----6-5---....', 'b---9---7-8---4-', '7-----5-4-------',
                '2---4---7-----9-', 'a-----8-7---5---', '6---8---b-----9-', '7-------........'],
            arp: '0.12.3.21.0.23.1',
            drums: { k: 'x.........x.....', s: '........x.......', h: '..x...x...x...x.' },
            leadWave: 'triangle', bassWave: 'sine', arpWave: 'sine', leadOct: 2
        },
        init: init, update: update, draw: draw, hud: hud
    });
})();
