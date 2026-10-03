/*
 * Outrun the Sun — the synthwave theme's game.
 *
 * A pseudo-3D racer: the road is a list of short segments, each with a curve
 * and a height, projected toward a striped sun. The sun sinks as the clock
 * runs down; every checkpoint lifts it again. Curves push the car outward,
 * hills change the top speed, the verge and everything standing on it cost
 * speed, and squeezing past traffic refills the nitro tank.
 */
(function () {
    'use strict';
    var G = window.SnoGame;
    var W = G.W, H = G.H, HORIZON = H / 2;
    // World units: a segment is 200 long, the road is 2000 from centre to
    // edge, and lateral positions ("x", "off") are fractions of that half-width.
    var SEG = 200, ROAD = 2000, CAM_H = 1000, DEPTH = 0.84, PZ = 1050, DRAW = 180;
    var MAX = 12000, ACCEL = MAX * 0.4, BRAKE = MAX * 0.9, COAST = MAX * 0.16, GRAV = MAX * 0.35;
    var STEER = 2.2, CENT = 0.6, OFF_LIMIT = MAX * 0.3, OFF_DECEL = MAX * 0.9;
    var LANES = [-0.75, -0.25, 0.25, 0.75], CAR_HALF = 0.11, PLAYER_Y = 486;
    var GRID_X = [1.7, 2.7, 4, 5.6, 7.6];
    var CAR_COLORS = ['#35e0ff', '#ffd24a', '#ff6b2b', '#7dff9a', '#bf3fff', '#f4f0ff'];

    var PALS = {
        sunset: { skyTop: '#2a0845', skyBot: '#ff6b2b', sun: ['#ffe45c', '#ff2d78'], ground: ['#1a0a3a', '#220d4a'], grid: '#ff2d78',
            road: ['#1b1030', '#21143a'], rumble: ['#ff2d78', '#35e0ff'], haze: '255,107,43', ridge: '#3b0f5e', edge: '#bf3fff', leaf: '#ff2d78', stars: 0.35 },
        dusk: { skyTop: '#0d0221', skyBot: '#bf3fff', sun: ['#ff9a3c', '#ff2d78'], ground: ['#12062e', '#190a3c'], grid: '#bf3fff',
            road: ['#160b2a', '#1c0f34'], rumble: ['#bf3fff', '#ff6b2b'], haze: '191,63,255', ridge: '#220a48', edge: '#ff2d78', leaf: '#35e0ff', stars: 0.7 },
        night: { skyTop: '#030014', skyBot: '#3a1078', sun: ['#ff5fa2', '#7a2bff'], ground: ['#07031a', '#0c0526'], grid: '#35e0ff',
            road: ['#0d0a1e', '#121028'], rumble: ['#35e0ff', '#ff2d78'], haze: '58,16,120', ridge: '#0a0524', edge: '#35e0ff', leaf: '#bf3fff', stars: 1 },
        canyon: { skyTop: '#3a0a2a', skyBot: '#ff8a3c', sun: ['#fff07a', '#ff6b2b'], ground: ['#2a0a1e', '#361026'], grid: '#ff6b2b',
            road: ['#1f0d1c', '#281224'], rumble: ['#ff6b2b', '#ffd24a'], haze: '255,138,60', ridge: '#5a1530', edge: '#ff6b2b', leaf: '#ffd24a', stars: 0.2 },
        storm: { skyTop: '#05060f', skyBot: '#2b2f55', sun: ['#c8d0ff', '#5a4fb0'], ground: ['#080a18', '#0c1022'], grid: '#5f7bff',
            road: ['#0c0e1c', '#111426'], rumble: ['#5f7bff', '#ff2d78'], haze: '43,47,85', ridge: '#0a0c1c', edge: '#5f7bff', leaf: '#5f7bff', stars: 0 },
        dawn: { skyTop: '#1a0b4a', skyBot: '#ffb347', sun: ['#fff3a0', '#ff2d78'], ground: ['#200a40', '#2b0e52'], grid: '#ffd24a',
            road: ['#1c1034', '#24153f'], rumble: ['#ffd24a', '#ff2d78'], haze: '255,179,71', ridge: '#4a1466', edge: '#ffd24a', leaf: '#ff2d78', stars: 0.4 }
    };

    // len: segments to the finish. curve/hold/alt/straight shape the bends,
    // hill the height changes. traffic, oil and works are counts per 1000
    // segments; oncoming is the size of the pool of cars driving toward you.
    // pace is the share of top speed the clock expects on average.
    var LEVELS = [
        { pal: 'sunset', len: 2800, cps: 3, curve: 2, hold: 1, straight: 0.4, hill: 0, traffic: 9, side: ['palm', 'pylon'], gap: 6, sideOff: 1.5, spread: 0.8, pace: 0.58 },
        { pal: 'sunset', len: 3000, cps: 3, curve: 3, hold: 0.5, alt: true, straight: 0.1, hill: 0, traffic: 11, side: ['palm'], gap: 3, sideOff: 1.25, spread: 0.25, pace: 0.62 },
        { pal: 'dusk', len: 3200, cps: 3, curve: 2.5, hold: 1, straight: 0.3, hill: 30, traffic: 12, side: ['pyramid', 'palm'], gap: 6, sideOff: 1.4, spread: 1.2, pace: 0.68 },
        { pal: 'night', city: true, len: 3400, cps: 3, curve: 2.5, hold: 1, straight: 0.5, hill: 0, traffic: 24, trucks: 0.3, side: ['tower', 'lamp'], gap: 4, sideOff: 1.3, spread: 0.5, pace: 0.705 },
        { pal: 'dawn', len: 3600, cps: 3, curve: 3.2, hold: 1, straight: 0.3, hill: 12, traffic: 14, oil: 14, side: ['palm', 'rock'], gap: 5, sideOff: 1.4, spread: 1, pace: 0.73 },
        { pal: 'canyon', len: 3600, cps: 4, curve: 5.5, hold: 1.4, straight: 0.25, hill: 30, traffic: 10, side: ['rock'], gap: 2, sideOff: 1.2, spread: 0.3, pace: 0.68 },
        { pal: 'night', len: 3800, cps: 4, curve: 3, hold: 1, straight: 0.35, hill: 15, traffic: 16, oncoming: 3, side: ['lamp', 'pylon'], gap: 5, sideOff: 1.3, spread: 0.6, pace: 0.7 },
        { pal: 'dusk', city: true, len: 4000, cps: 4, curve: 3.5, hold: 1, straight: 0.4, hill: 15, traffic: 16, trucks: 0.35, works: 6, side: ['tower', 'pylon'], gap: 5, sideOff: 1.3, spread: 0.6, pace: 0.74 },
        { pal: 'storm', len: 4200, cps: 4, curve: 4, hold: 1, straight: 0.3, hill: 25, traffic: 14, oncoming: 3, oil: 7, storm: true, wind: true, side: ['lamp', 'palm'], gap: 5, sideOff: 1.35, spread: 0.8, pace: 0.757 },
        { pal: 'dawn', len: 4800, cps: 4, curve: 5, hold: 1.1, straight: 0.3, hill: 35, traffic: 18, trucks: 0.25, oncoming: 4, oil: 8, works: 3, wind: true, side: ['rock', 'palm', 'pylon'], gap: 4, sideOff: 1.3, spread: 0.7, pace: 0.75 }
    ];

    // w/h: drawn size in world units. col: width that collides. hit: what
    // touching it does — 'side' objects stand on the verge and cost most of
    // the speed, 'break' objects shatter (cost = share of speed kept), 'oil'
    // takes the grip away.
    var KINDS = {
        palm: { w: 1900, h: 2500, col: 260, hit: 'side', draw: drawPalm },
        pylon: { w: 320, h: 1500, col: 200, hit: 'side', draw: drawPylon },
        pyramid: { w: 1000, h: 800, col: 800, hit: 'side', draw: drawPyramid },
        tower: { w: 1500, h: 4600, col: 1500, hit: 'side', draw: drawTower },
        lamp: { w: 900, h: 1900, col: 140, hit: 'side', draw: drawLamp },
        rock: { w: 1300, h: 1300, col: 1000, hit: 'side', draw: drawRock },
        chevron: { w: 800, h: 1000, col: 700, hit: 'side', draw: drawChevron },
        barrier: { w: 860, h: 520, col: 860, hit: 'break', cost: 0.4, draw: drawBarrier },
        cone: { w: 260, h: 340, col: 240, hit: 'break', cost: 0.9, draw: drawCone },
        oil: { w: 820, h: 140, col: 700, hit: 'oil', draw: drawOil }
    };
    var CARS = {
        car: { half: 0.11, h: 300, draw: drawCarBody },
        truck: { half: 0.15, h: 760, draw: drawTruck }
    };

    // ------------------------------------------------------------------
    // Track building
    // ------------------------------------------------------------------

    function easeIn(a, b, p) { return a + (b - a) * p * p; }
    function easeInOut(a, b, p) { return a + (b - a) * (0.5 - Math.cos(p * Math.PI) / 2); }
    function lastY(segs) { return segs.length ? segs[segs.length - 1].y2 : 0; }

    // The projected fields (sc, sx, …) are scratch space for draw(); keeping
    // them on the segment avoids allocating per frame.
    function addSeg(segs, curve, y) {
        segs.push({
            i: segs.length, y1: lastY(segs), y2: y, curve: curve, sprites: [], cars: [], gate: '',
            cx: 0, cx2: 0, sc: 0, sc2: 0, sx: 0, sx2: 0, sy: 0, sy2: 0, sw: 0, sw2: 0, clip: 0, vis: false
        });
    }

    // A stretch that eases into a curve, holds it and eases out, while the
    // height changes smoothly by `hill` segment lengths over the whole stretch.
    function addRoad(segs, enter, hold, leave, curve, hill) {
        var y0 = lastY(segs), y1 = y0 + hill * SEG, total = enter + hold + leave, n;
        for (n = 0; n < enter; n++) addSeg(segs, easeIn(0, curve, n / enter), easeInOut(y0, y1, (n + 1) / total));
        for (n = 0; n < hold; n++) addSeg(segs, curve, easeInOut(y0, y1, (enter + n + 1) / total));
        for (n = 0; n < leave; n++) addSeg(segs, easeInOut(curve, 0, n / leave), easeInOut(y0, y1, (enter + hold + n + 1) / total));
    }

    function addSection(segs, L, rnd, k) {
        var straight = rnd() < L.straight;
        // `alt` forces left-right-left S-bends instead of random directions.
        var dir = L.alt ? (k % 2 ? 1 : -1) : (rnd() < 0.5 ? -1 : 1);
        var curve = straight ? 0 : dir * L.curve * (0.5 + 0.5 * rnd());
        var hill = L.hill ? (rnd() * 2 - 1) * L.hill : 0;
        var len = 25 + Math.floor(rnd() * 30);
        addRoad(segs, len, Math.floor(len * L.hold * (0.6 + rnd())), len, curve, hill);
    }

    function buildTrack(L, rnd) {
        var segs = [], k = 0, finish;
        addRoad(segs, 20, 40, 20, 0, 0);
        while (segs.length < L.len) addSection(segs, L, rnd, k++);
        addRoad(segs, 20, 30, 20, 0, 0);
        finish = segs.length - 20;
        // Run-off, so there is still road on the horizon at the finish line.
        addRoad(segs, 60, 180, 60, 0, 0);
        for (k = 1; k <= L.cps; k++) segs[Math.round(finish * k / (L.cps + 1))].gate = 'cp';
        segs[finish].gate = 'finish';
        return { segs: segs, finish: finish };
    }

    // Closest a verge object may stand so that a car with all four wheels on
    // the rumble strip can never touch it.
    function minOff(kind) { return 1.08 + CAR_HALF + KINDS[kind].col / ROAD / 2; }

    // Lines the verge with the level's scenery. On the outside of a real
    // bend it puts chevron boards instead, which is how the player reads the
    // road ahead.
    function decorate(segs, L, rnd) {
        var i, side, seg, kind, bend;
        for (i = 12; i < segs.length; i += L.gap) {
            seg = segs[i];
            bend = Math.abs(seg.curve) >= 2 && Math.floor(i / L.gap) % 3 === 0;
            for (side = -1; side <= 1; side += 2) {
                if (bend && side === (seg.curve > 0 ? -1 : 1)) {
                    seg.sprites.push({ kind: 'chevron', off: side * minOff('chevron'), dir: -side, c: 0, dead: false });
                    continue;
                }
                kind = L.side[Math.floor(rnd() * L.side.length)];
                seg.sprites.push({
                    kind: kind, off: side * (Math.max(L.sideOff, minOff(kind)) + rnd() * L.spread),
                    dir: side, c: Math.floor(rnd() * 4), dead: false
                });
            }
        }
    }

    // Lanes that carry traffic in the player's direction.
    function ownLanes(L) { return L.oncoming ? [0.25, 0.75] : LANES; }

    function placeOil(segs, L, rnd, finish) {
        var n = Math.round((L.oil || 0) * finish / 1000), lanes = ownLanes(L), i;
        for (i = 0; i < n; i++) {
            segs[150 + Math.floor(rnd() * (finish - 200))].sprites.push({
                kind: 'oil', off: lanes[Math.floor(rnd() * lanes.length)], dir: 1, c: 0, dead: false
            });
        }
    }

    // A roadworks site: two rows of barriers across the blocked lanes with a
    // run of cones before them as the warning.
    function placeWorks(segs, L, rnd, finish) {
        var n = Math.round((L.works || 0) * finish / 1000), i, at, blocked, b, k;
        var sets = L.oncoming ? [[2], [3]] : [[0, 1], [1, 2], [2, 3]];
        for (i = 0; i < n; i++) {
            at = 200 + Math.floor(rnd() * (finish - 300));
            blocked = sets[Math.floor(rnd() * sets.length)];
            for (b = 0; b < blocked.length; b++) {
                for (k = 0; k <= 10; k += 10) segs[at + k].sprites.push({ kind: 'barrier', off: LANES[blocked[b]], dir: 1, c: 0, dead: false });
                for (k = 10; k <= 30; k += 10) segs[at - k].sprites.push({ kind: 'cone', off: LANES[blocked[b]], dir: 1, c: 0, dead: false });
            }
        }
    }

    function newCar(rnd, L, z, off, dir) {
        var truck = dir > 0 && rnd() < (L.trucks || 0);
        return {
            z: z, off: off, goal: off, dir: dir, kind: truck ? 'truck' : 'car', seg: -1,
            speed: MAX * (truck ? 0.22 + 0.1 * rnd() : 0.27 + 0.28 * rnd()),
            color: CAR_COLORS[Math.floor(rnd() * CAR_COLORS.length)], passed: false, honked: false, gone: false
        };
    }

    function placeTraffic(s, rnd) {
        var L = s.L, n = Math.round(L.traffic * s.finish / 1000), lanes = ownLanes(L), i, car;
        for (i = 0; i < n; i++) {
            car = newCar(rnd, L, (90 + rnd() * (s.finish - 100)) * SEG, lanes[Math.floor(rnd() * lanes.length)], 1);
            s.cars.push(car);
            reseat(s, car);
        }
        // Oncoming cars are a small pool that is recycled ahead of the
        // player, so the left lanes stay busy for the whole run.
        for (i = 0; i < (L.oncoming || 0); i++) {
            car = newCar(rnd, L, 0, LANES[i % 2], -1);
            car.speed = MAX * (0.3 + 0.15 * rnd());
            s.cars.push(car);
            respawnOncoming(s, car, 60 + i * 70);
        }
    }

    function backdrop(rnd) {
        var stars = [], ridge = [], i;
        for (i = 0; i < 70; i++) stars.push({ x: rnd() * W, y: 34 + rnd() * (HORIZON - 70), p: rnd() * 6.28 });
        for (i = 0; i < 48; i++) ridge.push(0.25 + 0.75 * rnd());
        return { stars: stars, ridge: ridge };
    }

    function init(level) {
        var L = LEVELS[level - 1], rnd = G.rng(level * 7919 + 13), track = buildTrack(L, rnd), sky = backdrop(rnd);
        // Each leg of the course buys this much time at the checkpoint; the
        // start gets a few seconds extra for pulling away from standstill.
        var leg = Math.round(track.finish / (L.cps + 1) * SEG / (MAX * L.pace));
        var s = {
            L: L, pal: PALS[L.pal], rnd: rnd, segs: track.segs, finish: track.finish, cars: [],
            pos: 0, x: 0.25, vx: 0, speed: 0, nitro: 1, boost: false, slip: 0, off: false,
            lastSeg: Math.floor(PZ / SEG), time: leg + 6, time0: leg + 6, cpTime: leg,
            wind: 0, gust: 0, gustDir: 1, gustIn: 6, bolt: 0, boltIn: 5,
            cp: 0, clock: 0, skyX: 0, sunX: 0, msg: '', msgT: 0, engT: 0, skidT: 0, beepT: 0, crashT: 0, dustT: 0,
            stars: sky.stars, ridge: sky.ridge
        };
        decorate(s.segs, L, rnd);
        placeOil(s.segs, L, rnd, s.finish);
        placeWorks(s.segs, L, rnd, s.finish);
        placeTraffic(s, rnd);
        return s;
    }

    // ------------------------------------------------------------------
    // Driving
    // ------------------------------------------------------------------

    function segAt(s, z) { return s.segs[G.clamp(Math.floor(z / SEG), 0, s.segs.length - 1)]; }
    function playerSeg(s) { return segAt(s, s.pos + PZ); }

    function say(s, text) { s.msg = text; s.msgT = 1.6; }

    // Thrust fades toward the top speed (air drag), and the slope adds or
    // takes its share, so hills really change how fast the car can go.
    function drive(s, dt) {
        var seg = playerSeg(s), slope = (seg.y2 - seg.y1) / SEG, pct = s.speed / MAX, a;
        var boost = G.key.a && s.nitro > 0;
        if (boost && !s.boost) nitroSound();
        s.boost = boost;
        var top = boost ? 1.3 : 1;
        if (G.key.down) a = -BRAKE;
        else if (G.key.up || boost) a = ACCEL * (boost ? 2.2 : 1) * (1 - (pct / top) * (pct / top));
        else a = -COAST;
        // A parked car on a slope stays put (handbrake); a moving one rolls.
        if (s.speed > 0 || G.key.up || boost) a -= slope * GRAV;
        s.speed = G.clamp(s.speed + a * dt, 0, MAX * 1.4);
        if (boost) {
            s.nitro = Math.max(0, s.nitro - dt / 3.2);
            G.burst(W / 2 + (Math.random() < 0.5 ? -38 : 38), PLAYER_Y - 8, { n: 1, color: '#35e0ff', speed: 140, life: 0.3, angle: Math.PI / 2, spread: 0.6 });
        }
    }

    // Steering has a little inertia; the curve pushes the car outward with
    // the square of the speed, so a hard bend must be taken slower.
    function steer(s, dt) {
        var pct = Math.min(1.2, s.speed / MAX), seg = playerSeg(s);
        var input = (G.key.right ? 1 : 0) - (G.key.left ? 1 : 0);
        var target = input * STEER * Math.min(1, pct * 3);
        // On oil the wheels barely answer, so the car keeps sliding the way it went.
        s.vx += (target - s.vx) * Math.min(1, dt * (s.slip > 0 ? 1 : 9));
        var drift = CENT * pct * pct * seg.curve * (s.L.storm ? 1.25 : 1) * (s.slip > 0 ? 1.5 : 1);
        s.x = G.clamp(s.x + (s.vx - drift + s.wind * pct) * dt, -2.3, 2.3);
        s.skyX += seg.curve * pct * dt * 40;
        s.sunX += (-seg.curve * 40 - s.sunX) * Math.min(1, dt * 1.5);
        if (Math.abs(drift) > 1.3 && s.skidT <= 0) {
            s.skidT = 0.16;
            G.noise(0.16, { filter: 'bandpass', freq: 1900 + Math.abs(drift) * 200, q: 6, vol: 0.05 });
        }
    }

    function offroad(s, dt) {
        s.off = Math.abs(s.x) > 1.05;
        if (!s.off || s.speed < 600) return;
        if (s.speed > OFF_LIMIT) s.speed = Math.max(OFF_LIMIT, s.speed - OFF_DECEL * dt);
        if (s.dustT > 0) return;
        s.dustT = 0.1;
        G.burst(W / 2 + G.rnd(-70, 70), PLAYER_Y - 6, { n: 2, color: s.pal.grid, speed: 150, life: 0.45, angle: -Math.PI / 2, spread: 1.6, gravity: 400 });
        G.noise(0.09, { freq: 300, vol: 0.1 });
    }

    function crashFx(s, big) {
        G.shake(big ? 14 : 7, big ? 0.5 : 0.25);
        G.burst(W / 2, PLAYER_Y - 40, { n: big ? 40 : 16, color: '#ffd24a', speed: 380, life: 0.6, gravity: 500 });
        if (s.crashT > 0) return;
        s.crashT = 0.3;
        G.sfx(big ? 'bigboom' : 'boom');
        G.tone(big ? 90 : 160, 0.25, { type: 'sawtooth', slide: 40, vol: 0.2 });
    }

    function touchSprite(s, spr) {
        var k = KINDS[spr.kind];
        if (spr.dead || Math.abs(s.x - spr.off) > CAR_HALF + k.col / ROAD / 2) return;
        if (k.hit === 'oil') {
            s.slip = 1.3;
            say(s, 'OIL!');
            G.tone(700, 0.5, { type: 'sine', slide: 180, vol: 0.14 });
            G.noise(0.4, { filter: 'bandpass', freq: 900, slide: 300, q: 3, vol: 0.12 });
        } else if (k.hit === 'break') {
            spr.dead = true;
            s.speed *= k.cost;
            G.burst(W / 2, PLAYER_Y - 30, { n: 14, color: '#ff6b2b', speed: 320, life: 0.7, gravity: 600 });
            G.noise(0.18, { freq: 2200, slide: 300, vol: 0.2 });
            G.tone(k.cost < 0.5 ? 130 : 520, 0.09, { type: 'square', slide: 60, vol: 0.15 });
            if (k.cost < 0.5) G.shake(8, 0.25);
        } else {
            // Verge objects do not stop the car dead (that would trap it);
            // they take most of its speed and let it plough on.
            s.speed = Math.max(MAX * 0.08, s.speed * 0.3);
            crashFx(s, false);
        }
    }

    function passGate(s, seg) {
        if (seg.gate === 'finish') {
            G.win(1000 * G.level + Math.ceil(s.time) * 100);
            return;
        }
        s.time += s.cpTime;
        s.cp++;
        s.nitro = Math.min(1, s.nitro + 0.34);
        say(s, 'CHECKPOINT  +' + s.cpTime + 's');
        G.addScore(500);
        [523, 784, 1047, 1568].forEach(function (f, i) { G.tone(f, 0.14, { type: 'triangle', vol: 0.16, delay: i * 0.07 }); });
    }

    // Handles every segment the car entered this tick (more than one at
    // nitro speed), so nothing on the road can be skipped over.
    function crossSegments(s) {
        var now = playerSeg(s).i, i, seg, j;
        for (i = s.lastSeg + 1; i <= now; i++) {
            seg = s.segs[i];
            G.addScore(1);
            for (j = 0; j < seg.sprites.length; j++) touchSprite(s, seg.sprites[j]);
            if (seg.gate) passGate(s, seg);
        }
        // Never backwards: a rear-end shunt drops the car back a little, and
        // the segments it re-enters must not count twice.
        s.lastSeg = Math.max(s.lastSeg, now);
    }

    // ------------------------------------------------------------------
    // Traffic
    // ------------------------------------------------------------------

    function reseat(s, car) {
        var idx = Math.floor(car.z / SEG), old = s.segs[car.seg];
        if (idx === car.seg) return;
        if (old) old.cars.splice(old.cars.indexOf(car), 1);
        car.seg = idx;
        s.segs[idx].cars.push(car);
        if (car.dir > 0) dodgeWorks(s, car);
    }

    function laneBlocked(seg, lane) {
        for (var j = 0; j < seg.sprites.length; j++) {
            if (seg.sprites[j].kind === 'barrier' && !seg.sprites[j].dead && Math.abs(seg.sprites[j].off - lane) < 0.2) return true;
        }
        return false;
    }

    // Traffic pulls out of a coned-off lane well before it reaches the
    // barriers, into the nearest lane that is open.
    function dodgeWorks(s, car) {
        var ahead = s.segs[car.seg + 30], lanes = ownLanes(s.L), best = car.goal, d = 9, i;
        if (!ahead || !laneBlocked(ahead, car.goal)) return;
        for (i = 0; i < lanes.length; i++) {
            if (!laneBlocked(ahead, lanes[i]) && Math.abs(lanes[i] - car.goal) < d) { d = Math.abs(lanes[i] - car.goal); best = lanes[i]; }
        }
        car.goal = best;
    }

    function respawnOncoming(s, car, ahead) {
        var old = s.segs[car.seg], idx = Math.floor(s.pos / SEG) + Math.floor(ahead);
        if (old) old.cars.splice(old.cars.indexOf(car), 1);
        car.seg = -1;
        // Nothing drives out of the void beyond the finish line.
        car.gone = idx > s.finish + 150;
        if (car.gone) return;
        car.z = idx * SEG + 1;
        car.off = car.goal = LANES[s.rnd() < 0.5 ? 0 : 1];
        car.passed = false; car.honked = false;
        reseat(s, car);
    }

    // The speed of whatever blocks this car: a slower car just ahead in its
    // lane, or the player's own car, so that traffic never rams from behind.
    function leaderSpeed(s, car) {
        var v = car.speed, n, seg, i, o, gap = s.pos + PZ - car.z;
        for (n = 0; n < 3; n++) {
            seg = s.segs[car.seg + n];
            if (!seg) break;
            for (i = 0; i < seg.cars.length; i++) {
                o = seg.cars[i];
                if (o !== car && o.dir > 0 && o.z > car.z && Math.abs(o.off - car.off) < 0.2) v = Math.min(v, o.speed);
            }
        }
        if (gap > 0 && gap < 3 * SEG && Math.abs(s.x - car.off) < CAR_HALF + CARS[car.kind].half) v = Math.min(v, s.speed);
        return v;
    }

    function moveCar(s, car, dt) {
        if (car.gone) return;
        if (car.dir < 0) {
            car.z -= car.speed * dt;
            // Recycled once it is behind the camera (and before it leaves the track).
            if (car.z < Math.max(SEG, s.pos - SEG)) { respawnOncoming(s, car, DRAW + 5 + s.rnd() * 200); return; }
        } else {
            car.z = Math.min(car.z + leaderSpeed(s, car) * dt, (s.segs.length - 3) * SEG);
            car.off += G.clamp(car.goal - car.off, -0.8 * dt, 0.8 * dt);
        }
        reseat(s, car);
    }

    function carCrash(s, car) {
        if (car.dir < 0) {
            // Head-on: the worst thing on the road. All speed is gone.
            car.passed = true;
            s.speed = MAX * 0.04;
            say(s, 'HEAD-ON!');
            G.flash('#ff2d78', 0.25);
            crashFx(s, true);
            return;
        }
        // Rear-ending a car drops the player behind it at less than its speed.
        s.speed = car.speed * 0.7;
        s.pos = car.z - PZ - 30;
        crashFx(s, false);
    }

    // Squeezing past a car is rewarded: that is where the nitro comes from.
    function carPassed(s, car, dx) {
        var near = dx < 0.45;
        car.passed = true;
        G.addScore(near ? 100 * (car.dir < 0 ? 2 : 1) : 25);
        if (!near) return;
        s.nitro = Math.min(1, s.nitro + (car.dir < 0 ? 0.15 : 0.1));
        G.popup(W / 2 + (car.off - s.x) * 300, PLAYER_Y - 90, 'NEAR MISS', '#35e0ff');
        G.noise(0.22, { filter: 'bandpass', freq: 500, slide: 2600, q: 2, vol: 0.14 });
        G.tone(1320, 0.07, { type: 'triangle', vol: 0.1 });
    }

    function honk(car) {
        car.honked = true;
        G.tone(392, 0.28, { type: 'square', vol: 0.09 });
        G.tone(494, 0.28, { type: 'square', vol: 0.09 });
    }

    // A car counts as met on the tick it goes from ahead to behind, which
    // works at any closing speed.
    function meetTraffic(s) {
        var pz = s.pos + PZ, i, car, dz, dx;
        for (i = 0; i < s.cars.length; i++) {
            car = s.cars[i];
            if (car.gone) continue;
            dz = car.z - pz;
            dx = Math.abs(car.off - s.x);
            if (car.dir < 0 && !car.honked && dz > 0 && dz < 60 * SEG && dx < 0.3) honk(car);
            if (dz >= 0) { car.passed = false; continue; }
            if (car.passed) continue;
            if (dx < CAR_HALF + CARS[car.kind].half) carCrash(s, car);
            else carPassed(s, car, dx);
        }
    }

    // ------------------------------------------------------------------
    // Weather, clock and sound
    // ------------------------------------------------------------------

    // Crosswind comes in announced gusts that swell and fade, pushing the car
    // sideways the faster it goes.
    function weather(s, dt) {
        if (s.L.wind) {
            s.gustIn -= dt;
            if (s.gustIn <= 0) {
                s.gustIn = 7 + s.rnd() * 5; s.gust = 2.6; s.gustDir = s.rnd() < 0.5 ? -1 : 1;
                G.noise(2.4, { filter: 'bandpass', freq: 400, slide: 1400, q: 1.5, vol: 0.16, attack: 0.8 });
            }
            s.gust = Math.max(0, s.gust - dt);
            s.wind = s.gustDir * 0.85 * Math.sin(Math.PI * s.gust / 2.6);
        }
        if (!s.L.storm) return;
        s.boltIn -= dt;
        s.bolt = Math.max(0, s.bolt - dt);
        if (s.boltIn <= 0) {
            s.boltIn = 4 + s.rnd() * 6; s.bolt = 0.35;
            G.noise(1.6, { freq: 220, slide: 50, vol: 0.3, delay: 0.25 });
        }
    }

    function nitroSound() {
        G.tone(180, 0.5, { type: 'sawtooth', slide: 1300, vol: 0.14 });
        G.noise(0.6, { filter: 'highpass', freq: 800, slide: 5000, vol: 0.14 });
    }

    // The engine is a short note retriggered about eleven times a second; its
    // pitch climbs through four "gears" so speed can be heard.
    function engineSound(s) {
        if (s.engT > 0) return;
        s.engT = 0.09;
        var p = Math.min(1.3, s.speed / MAX) * 4, gear = Math.min(3, Math.floor(p));
        var f = 52 + gear * 20 + (p - gear) * 62;
        G.tone(f, 0.11, { type: 'sawtooth', vol: 0.05 });
        G.tone(f / 2, 0.11, { type: 'triangle', vol: 0.07 });
    }

    function tickTimers(s, dt) {
        s.clock += dt;
        s.time -= dt;
        ['msgT', 'engT', 'skidT', 'crashT', 'dustT', 'slip', 'beepT'].forEach(function (k) { if (s[k] > 0) s[k] -= dt; });
        // The last ten seconds beep, faster in the last five.
        if (s.time < 10 && s.time > 0 && s.beepT <= 0) {
            s.beepT = s.time < 5 ? 0.5 : 1;
            G.tone(s.time < 5 ? 1175 : 880, 0.09, { type: 'square', vol: 0.12 });
        }
    }

    function update(s, dt) {
        tickTimers(s, dt);
        drive(s, dt);
        steer(s, dt);
        s.pos += s.speed * dt;
        offroad(s, dt);
        crossSegments(s);
        for (var i = 0; i < s.cars.length; i++) moveCar(s, s.cars[i], dt);
        meetTraffic(s);
        weather(s, dt);
        engineSound(s);
        // The sun has set: out of time.
        if (s.time <= 0) G.die();
    }

    // ------------------------------------------------------------------
    // Drawing: sky
    // ------------------------------------------------------------------

    function drawSun(s, ctx) {
        var pal = s.pal, R = 100, left = G.clamp(s.time / s.time0, 0, 1), cx = W / 2 + s.sunX;
        // The sun's height is the clock: it touches the horizon as time runs out.
        var cy = HORIZON - 70 + (1 - left) * 150, k, yy, g;
        ctx.save();
        ctx.beginPath(); ctx.rect(0, 0, W, HORIZON); ctx.clip();
        g = ctx.createRadialGradient(cx, cy, R * 0.6, cx, cy, R * 2.1);
        g.addColorStop(0, 'rgba(' + pal.haze + ',0.5)'); g.addColorStop(1, 'rgba(' + pal.haze + ',0)');
        ctx.fillStyle = g; ctx.fillRect(cx - R * 2.2, cy - R * 2.2, R * 4.4, R * 4.4);
        ctx.beginPath(); ctx.arc(cx, cy, R, 0, Math.PI * 2); ctx.clip();
        g = ctx.createLinearGradient(0, cy - R, 0, cy + R);
        g.addColorStop(0, pal.sun[0]); g.addColorStop(1, pal.sun[1]);
        ctx.fillStyle = g; ctx.fillRect(cx - R, cy - R, R * 2, R * 2);
        // Dark bands slide down the lower half and thicken, the classic look.
        ctx.fillStyle = pal.skyTop;
        for (k = 0; k < 7; k++) {
            yy = ((k + (s.clock * 0.3) % 1) / 7);
            ctx.fillRect(cx - R, cy - R * 0.1 + yy * R * 1.1, R * 2, 1.5 + yy * 9);
        }
        ctx.restore();
    }

    // A skyline that scrolls sideways with the bends: jagged peaks, or the
    // stepped blocks of a city.
    function drawRidge(s, ctx, parallax, height, fill, edge) {
        var step = 40, n = s.ridge.length, shift = ((s.skyX * parallax) % (n * step) + n * step) % (n * step);
        var first = Math.floor(shift / step), x = -(shift % step), i, h;
        ctx.beginPath();
        ctx.moveTo(x, HORIZON);
        for (i = 0; x <= W + step; i++, x += step) {
            h = HORIZON - s.ridge[(first + i) % n] * height;
            if (s.L.city) { ctx.lineTo(x, h); ctx.lineTo(x + step, h); } else ctx.lineTo(x + step / 2, h);
        }
        ctx.lineTo(W + step, HORIZON);
        ctx.fillStyle = fill; ctx.fill();
        if (!edge) return;
        ctx.strokeStyle = edge; ctx.lineWidth = 1.5; ctx.globalAlpha = 0.7; ctx.stroke(); ctx.globalAlpha = 1;
    }

    function drawSky(s, ctx) {
        var pal = s.pal, g = ctx.createLinearGradient(0, 0, 0, HORIZON), i, st;
        g.addColorStop(0, pal.skyTop); g.addColorStop(1, pal.skyBot);
        ctx.fillStyle = g; ctx.fillRect(0, 0, W, HORIZON);
        ctx.fillStyle = '#fff';
        for (i = 0; i < s.stars.length && pal.stars > 0; i++) {
            st = s.stars[i];
            ctx.globalAlpha = pal.stars * (0.45 + 0.55 * Math.sin(s.clock * 2 + st.p));
            ctx.fillRect(st.x, st.y, 2, 2);
        }
        ctx.globalAlpha = 1;
        drawSun(s, ctx);
        if (s.bolt > 0) { ctx.fillStyle = 'rgba(220,225,255,' + (s.bolt * 1.6).toFixed(2) + ')'; ctx.fillRect(0, 0, W, HORIZON); }
        drawRidge(s, ctx, 0.5, 46, pal.skyTop, null);
        drawRidge(s, ctx, 1, s.L.city ? 78 : 60, pal.ridge, pal.edge);
        ctx.fillStyle = pal.ground[0]; ctx.fillRect(0, HORIZON, W, H - HORIZON);
    }

    // ------------------------------------------------------------------
    // Drawing: road
    // ------------------------------------------------------------------

    function camera(s) {
        var seg = playerSeg(s), p = ((s.pos + PZ) % SEG) / SEG;
        return { x: s.x * ROAD, y: G.lerp(seg.y1, seg.y2, p) + CAM_H, z: s.pos };
    }

    // Projects both ends of a segment. x and dx are the sideways offsets the
    // curves so far have piled up at its near and far end.
    function projectSeg(seg, cam, x, dx) {
        var z1 = seg.i * SEG - cam.z, z2 = z1 + SEG;
        seg.cx = x; seg.cx2 = x + dx;
        seg.sc = z1 > DEPTH ? DEPTH / z1 : 0;
        seg.sc2 = z2 > DEPTH ? DEPTH / z2 : 0;
        seg.sx = W / 2 + seg.sc * (x - cam.x) * W / 2;
        seg.sx2 = W / 2 + seg.sc2 * (x + dx - cam.x) * W / 2;
        seg.sy = H / 2 - seg.sc * (seg.y1 - cam.y) * H / 2;
        seg.sy2 = H / 2 - seg.sc2 * (seg.y2 - cam.y) * H / 2;
        seg.sw = seg.sc * ROAD * W / 2;
        seg.sw2 = seg.sc2 * ROAD * W / 2;
    }

    function quad(ctx, x1, y1, w1, x2, y2, w2, color) {
        ctx.fillStyle = color;
        ctx.beginPath();
        ctx.moveTo(x1 - w1, y1); ctx.lineTo(x1 + w1, y1); ctx.lineTo(x2 + w2, y2); ctx.lineTo(x2 - w2, y2);
        ctx.fill();
    }

    function drawSegment(s, ctx, seg, y1) {
        var pal = s.pal, band = Math.floor(seg.i / 3) % 2, x1 = seg.sx, x2 = seg.sx2, w1 = seg.sw, w2 = seg.sw2, y2 = seg.sy2, k;
        ctx.fillStyle = pal.ground[band];
        ctx.fillRect(0, y2, W, y1 - y2 + 1);
        quad(ctx, x1, y1, w1 * 1.12, x2, y2, w2 * 1.12, pal.rumble[band]);
        quad(ctx, x1, y1, w1, x2, y2, w2, seg.gate ? (seg.gate === 'finish' ? '#ffffff' : pal.grid) : pal.road[band]);
        if (s.L.oncoming) {
            // A solid double line: the left half belongs to oncoming traffic.
            quad(ctx, x1 - w1 * 0.02, y1, w1 * 0.008, x2 - w2 * 0.02, y2, w2 * 0.008, '#ff6b2b');
            quad(ctx, x1 + w1 * 0.02, y1, w1 * 0.008, x2 + w2 * 0.02, y2, w2 * 0.008, '#ff6b2b');
        }
        if (band) return;
        for (k = -1; k <= 1; k++) {
            if (k === 0 && s.L.oncoming) continue;
            quad(ctx, x1 + w1 * k * 0.5, y1, w1 * 0.012, x2 + w2 * k * 0.5, y2, w2 * 0.012, '#ffd9f0');
        }
    }

    // Draws the road front to back. A segment is skipped when a nearer hill
    // already covers it; seg.clip remembers that height for the objects pass.
    function drawRoad(s, ctx, cam) {
        var base = Math.floor(cam.z / SEG), x = 0, dx = -s.segs[base].curve * ((cam.z % SEG) / SEG), maxy = H, n, seg;
        for (n = 0; n < DRAW; n++) {
            seg = s.segs[base + n];
            if (!seg) break;
            projectSeg(seg, cam, x, dx);
            x += dx; dx += seg.curve;
            seg.clip = maxy;
            seg.vis = seg.sc > 0 && seg.sy2 < seg.sy && seg.sy2 < maxy;
            if (!seg.vis) continue;
            drawSegment(s, ctx, seg, Math.min(seg.sy, maxy));
            maxy = seg.sy2;
        }
    }

    function gridLine(ctx, seg, k) {
        var x1 = seg.sx + k * seg.sw, x2 = seg.sx2 + k * seg.sw2;
        if ((x1 < 0 && x2 < 0) || (x1 > W && x2 > W)) return;
        ctx.moveTo(x1, Math.min(seg.sy, seg.clip)); ctx.lineTo(x2, seg.sy2);
    }

    // The neon grid on the plain beside the road, stroked as one path.
    function drawGrid(s, ctx, base) {
        var n, seg, k, y;
        ctx.strokeStyle = s.pal.grid; ctx.globalAlpha = 0.4; ctx.lineWidth = 1;
        ctx.beginPath();
        for (n = 1; n < DRAW; n++) {
            seg = s.segs[base + n];
            if (!seg || !seg.vis) continue;
            for (k = 0; k < GRID_X.length; k++) { gridLine(ctx, seg, GRID_X[k]); gridLine(ctx, seg, -GRID_X[k]); }
            if (seg.i % 4 || n > 110) continue;
            y = seg.sy2;
            ctx.moveTo(0, y); ctx.lineTo(seg.sx2 - seg.sw2 * 1.15, y);
            ctx.moveTo(seg.sx2 + seg.sw2 * 1.15, y); ctx.lineTo(W, y);
        }
        ctx.stroke();
        ctx.globalAlpha = 1;
    }

    // ------------------------------------------------------------------
    // Drawing: roadside objects and traffic. Every draw function gets the
    // ground point (x, y) and u, the pixels per world unit at that distance.
    // ------------------------------------------------------------------

    // Returns 0 when a hill hides the object completely, 2 when it had to
    // set a clip (the caller must restore), 1 otherwise.
    function openClip(ctx, clipY, top, bottom) {
        if (top >= clipY || bottom < 0) return 0;
        if (bottom <= clipY) return 1;
        ctx.save();
        ctx.beginPath(); ctx.rect(0, 0, W, clipY); ctx.clip();
        return 2;
    }

    function drawPalm(ctx, x, y, u, spr, s) {
        var lean = spr.dir * 160 * u, tx = x + lean, ty = y - 1950 * u, i, a;
        ctx.strokeStyle = '#2a0f4a'; ctx.lineWidth = Math.max(1, 110 * u);
        ctx.beginPath(); ctx.moveTo(x, y); ctx.quadraticCurveTo(x - lean * 0.6, y - 900 * u, tx, ty); ctx.stroke();
        ctx.strokeStyle = s.pal.leaf; ctx.lineWidth = Math.max(1, 60 * u);
        ctx.beginPath();
        // Seven fronds fanned over the crown, each drooping at its tip.
        for (i = 0; i < 7; i++) {
            a = -2.95 + i * 0.46;
            ctx.moveTo(tx, ty);
            ctx.quadraticCurveTo(tx + Math.cos(a) * 520 * u, ty + Math.sin(a) * 520 * u - 200 * u, tx + Math.cos(a) * 900 * u, ty + Math.sin(a) * 450 * u + 260 * u);
        }
        ctx.stroke();
    }

    function drawPylon(ctx, x, y, u, spr, s) {
        var col = spr.c % 2 ? '#35e0ff' : '#ff2d78', w = 150 * u, h = 1500 * u;
        ctx.globalAlpha = 0.22 + 0.1 * Math.sin(s.clock * 6 + spr.c);
        ctx.fillStyle = col; ctx.fillRect(x - w * 1.1, y - h, w * 2.2, h);
        ctx.globalAlpha = 1;
        ctx.fillRect(x - w / 2, y - h, w, h);
        ctx.fillStyle = '#fff'; ctx.fillRect(x - w / 2, y - h, w, Math.max(1, 70 * u));
    }

    function drawPyramid(ctx, x, y, u, spr, s) {
        var w = 500 * u, h = 800 * u, k;
        ctx.beginPath(); ctx.moveTo(x - w, y); ctx.lineTo(x, y - h); ctx.lineTo(x + w, y); ctx.closePath();
        ctx.fillStyle = '#0d0221'; ctx.fill();
        ctx.strokeStyle = s.pal.edge; ctx.lineWidth = Math.max(1, 26 * u);
        for (k = 1; k < 4; k++) { ctx.moveTo(x - w * (1 - k / 4), y - h * k / 4); ctx.lineTo(x + w * (1 - k / 4), y - h * k / 4); }
        ctx.moveTo(x, y - h); ctx.lineTo(x, y);
        ctx.stroke();
    }

    function drawTower(ctx, x, y, u, spr, s) {
        var w = 750 * u, h = (2600 + spr.c * 650) * u, rows = 9, r;
        ctx.fillStyle = '#0b0420'; ctx.fillRect(x - w, y - h, w * 2, h);
        // Lit floors: which ones glow depends on the tower, so no two match.
        for (r = 0; r < rows; r++) {
            if ((r * 7 + spr.c * 3) % 4 === 0) continue;
            ctx.fillStyle = (r + spr.c) % 3 ? '#ff2d78' : '#35e0ff';
            ctx.globalAlpha = 0.55;
            ctx.fillRect(x - w * 0.8, y - h + (r + 0.5) * h / (rows + 1), w * 1.6, h / (rows + 1) * 0.35);
        }
        ctx.globalAlpha = 1;
        ctx.strokeStyle = s.pal.edge; ctx.lineWidth = 1; ctx.strokeRect(x - w, y - h, w * 2, h);
    }

    function drawLamp(ctx, x, y, u, spr) {
        var h = 1900 * u, arm = -spr.dir * 420 * u;
        ctx.strokeStyle = '#4a3a78'; ctx.lineWidth = Math.max(1, 50 * u);
        ctx.beginPath(); ctx.moveTo(x, y); ctx.lineTo(x, y - h); ctx.lineTo(x + arm, y - h); ctx.stroke();
        ctx.fillStyle = 'rgba(255,214,120,0.25)';
        ctx.beginPath(); ctx.arc(x + arm, y - h, 230 * u, 0, Math.PI * 2); ctx.fill();
        ctx.fillStyle = '#ffe9a8';
        ctx.beginPath(); ctx.arc(x + arm, y - h, Math.max(1, 80 * u), 0, Math.PI * 2); ctx.fill();
    }

    function drawRock(ctx, x, y, u, spr, s) {
        var w = 650 * u, h = (900 + spr.c * 130) * u;
        ctx.beginPath();
        ctx.moveTo(x - w, y); ctx.lineTo(x - w * 0.7, y - h * 0.6); ctx.lineTo(x - w * 0.2, y - h);
        ctx.lineTo(x + w * 0.3, y - h * 0.75); ctx.lineTo(x + w * 0.75, y - h * 0.9); ctx.lineTo(x + w, y);
        ctx.closePath();
        ctx.fillStyle = s.pal.ridge; ctx.fill();
        ctx.strokeStyle = s.pal.edge; ctx.lineWidth = Math.max(1, 22 * u); ctx.stroke();
    }

    function drawChevron(ctx, x, y, u, spr, s) {
        var w = 400 * u, h = 520 * u, top = y - 1000 * u, k, ax, on = Math.floor(s.clock * 5) % 2;
        ctx.fillStyle = '#2a0f4a'; ctx.fillRect(x - 30 * u, top + h, 60 * u, 1000 * u - h);
        ctx.fillStyle = '#12041f'; ctx.fillRect(x - w, top, w * 2, h);
        ctx.strokeStyle = on ? '#ffd24a' : '#ff6b2b'; ctx.lineWidth = Math.max(1, 70 * u);
        ctx.beginPath();
        // Two arrow heads pointing the way the road turns.
        for (k = -1; k <= 1; k += 2) {
            ax = x + k * 170 * u;
            ctx.moveTo(ax - spr.dir * 110 * u, top + h * 0.2); ctx.lineTo(ax + spr.dir * 110 * u, top + h * 0.5); ctx.lineTo(ax - spr.dir * 110 * u, top + h * 0.8);
        }
        ctx.stroke();
    }

    function drawBarrier(ctx, x, y, u, spr, s) {
        var w = 430 * u, h = 300 * u, top = y - 520 * u, k;
        ctx.fillStyle = '#3a2a55';
        ctx.fillRect(x - w * 0.8, top + h, 50 * u, 520 * u - h); ctx.fillRect(x + w * 0.8 - 50 * u, top + h, 50 * u, 520 * u - h);
        for (k = 0; k < 6; k++) {
            ctx.fillStyle = k % 2 ? '#fff' : '#ff2d55';
            ctx.fillRect(x - w + k * w / 3, top, w / 3 + 0.5, h);
        }
        if (Math.floor(s.clock * 4) % 2) return;
        ctx.fillStyle = '#ffd24a';
        ctx.beginPath(); ctx.arc(x, top - 60 * u, Math.max(1.5, 80 * u), 0, Math.PI * 2); ctx.fill();
    }

    function drawCone(ctx, x, y, u) {
        var w = 130 * u, h = 340 * u;
        ctx.fillStyle = '#ff6b2b';
        ctx.beginPath(); ctx.moveTo(x - w, y); ctx.lineTo(x, y - h); ctx.lineTo(x + w, y); ctx.fill();
        ctx.fillStyle = '#fff'; ctx.fillRect(x - w * 0.55, y - h * 0.5, w * 1.1, h * 0.14);
    }

    function drawOil(ctx, x, y, u, spr, s) {
        var w = 410 * u, h = 70 * u;
        ctx.fillStyle = '#05010c';
        ctx.beginPath(); ctx.ellipse(x, y - h, w, h, 0, 0, Math.PI * 2); ctx.fill();
        // A shimmering rim so the slick reads against the dark road.
        ctx.strokeStyle = Math.floor(s.clock * 3) % 2 ? '#35e0ff' : '#bf3fff'; ctx.lineWidth = Math.max(1, 30 * u);
        ctx.beginPath(); ctx.ellipse(x, y - h, w * 0.7, h * 0.6, 0, 0, Math.PI * 2); ctx.stroke();
    }

    // Seen from behind, or from the front (white headlights) when oncoming.
    function drawCarBody(ctx, x, y, u, car) {
        var w = 220 * u, h = 300 * u, front = car.dir < 0;
        ctx.fillStyle = 'rgba(0,0,0,0.5)'; ctx.fillRect(x - w, y - h * 0.08, w * 2, h * 0.1);
        ctx.fillStyle = '#05010c';
        ctx.fillRect(x - w, y - h * 0.3, w * 0.3, h * 0.3); ctx.fillRect(x + w * 0.7, y - h * 0.3, w * 0.3, h * 0.3);
        ctx.fillStyle = car.color; ctx.fillRect(x - w, y - h * 0.62, w * 2, h * 0.48);
        ctx.beginPath();
        ctx.moveTo(x - w * 0.9, y - h * 0.62); ctx.lineTo(x - w * 0.68, y - h); ctx.lineTo(x + w * 0.68, y - h); ctx.lineTo(x + w * 0.9, y - h * 0.62);
        ctx.fill();
        ctx.fillStyle = '#160629';
        ctx.fillRect(x - w * 0.62, y - h * 0.93, w * 1.24, h * 0.26);
        ctx.fillStyle = front ? '#fff6c0' : '#ff2d55';
        ctx.fillRect(x - w * 0.92, y - h * 0.55, w * 0.42, h * 0.16); ctx.fillRect(x + w * 0.5, y - h * 0.55, w * 0.42, h * 0.16);
        if (!front) return;
        ctx.fillStyle = 'rgba(255,246,192,0.28)';
        ctx.beginPath(); ctx.arc(x - w * 0.7, y - h * 0.47, w * 0.6, 0, Math.PI * 2); ctx.arc(x + w * 0.7, y - h * 0.47, w * 0.6, 0, Math.PI * 2); ctx.fill();
    }

    function drawTruck(ctx, x, y, u, car) {
        var w = 300 * u, h = 760 * u;
        ctx.fillStyle = '#05010c';
        ctx.fillRect(x - w, y - h * 0.14, w * 0.34, h * 0.14); ctx.fillRect(x + w * 0.66, y - h * 0.14, w * 0.34, h * 0.14);
        ctx.fillStyle = '#241040'; ctx.fillRect(x - w, y - h, w * 2, h * 0.9);
        ctx.fillStyle = car.color;
        ctx.fillRect(x - w, y - h, w * 2, h * 0.07); ctx.fillRect(x - w * 0.04, y - h, w * 0.08, h * 0.9);
        ctx.fillStyle = '#ff2d55';
        ctx.fillRect(x - w * 0.95, y - h * 0.22, w * 0.3, h * 0.07); ctx.fillRect(x + w * 0.65, y - h * 0.22, w * 0.3, h * 0.07);
    }

    function drawSprite(s, ctx, seg, spr) {
        var k = KINDS[spr.kind], u = seg.sc * W / 2, x = seg.sx + spr.off * seg.sw, y = seg.sy, clip;
        // Too close to the lens, too small to see, or off to the side.
        if (spr.dead || seg.sc <= 0 || seg.sc > 0.004 || k.w * u < 1 || x + k.w * u < 0 || x - k.w * u > W) return;
        clip = openClip(ctx, seg.clip, y - k.h * u, y);
        if (!clip) return;
        k.draw(ctx, x, y, u, spr, s);
        if (clip === 2) ctx.restore();
    }

    function drawCar(s, ctx, cam, seg, car) {
        var cz = car.z - cam.z, k = CARS[car.kind], p, sc, u, x, y, clip;
        if (car.gone || cz < 300) return;
        p = (car.z - seg.i * SEG) / SEG; sc = DEPTH / cz; u = sc * W / 2;
        x = W / 2 + u * (car.off * ROAD + G.lerp(seg.cx, seg.cx2, p) - cam.x);
        y = H / 2 - sc * (G.lerp(seg.y1, seg.y2, p) - cam.y) * H / 2;
        clip = openClip(ctx, seg.clip, y - k.h * u, y);
        if (!clip) return;
        k.draw(ctx, x, y, u, car);
        if (clip === 2) ctx.restore();
    }

    // The arch over a checkpoint or the finish line.
    function drawGate(s, ctx, seg) {
        var u = seg.sc * W / 2, x = seg.sx, y = seg.sy, w = seg.sw * 1.25, h = 1700 * u, clip;
        if (seg.sc <= 0 || seg.sc > 0.004) return;
        clip = openClip(ctx, seg.clip, y - h, y);
        if (!clip) return;
        ctx.fillStyle = seg.gate === 'finish' ? '#ffd24a' : '#35e0ff';
        ctx.fillRect(x - w, y - h, 90 * u, h); ctx.fillRect(x + w - 90 * u, y - h, 90 * u, h);
        ctx.fillStyle = '#12041f'; ctx.fillRect(x - w, y - h, w * 2, 420 * u);
        ctx.strokeStyle = seg.gate === 'finish' ? '#ffd24a' : '#35e0ff'; ctx.lineWidth = Math.max(1, 30 * u);
        ctx.strokeRect(x - w, y - h, w * 2, 420 * u);
        if (u > 0.012) {
            G.text(seg.gate === 'finish' ? 'F I N I S H' : 'CHECKPOINT', x, y - h + 300 * u, { size: G.clamp(300 * u, 6, 64), color: '#fff', align: 'center', bold: true });
        }
        if (clip === 2) ctx.restore();
    }

    // Far to near, so nearer things cover farther ones.
    function drawObjects(s, ctx, cam, base) {
        var n, seg, i;
        for (n = DRAW - 1; n >= 0; n--) {
            seg = s.segs[base + n];
            if (!seg) continue;
            if (seg.gate) drawGate(s, ctx, seg);
            for (i = 0; i < seg.sprites.length; i++) drawSprite(s, ctx, seg, seg.sprites[i]);
            for (i = 0; i < seg.cars.length; i++) drawCar(s, ctx, cam, seg, seg.cars[i]);
        }
    }

    // ------------------------------------------------------------------
    // Drawing: weather, the player's car and the dashboard
    // ------------------------------------------------------------------

    function drawHaze(s, ctx) {
        var g = ctx.createLinearGradient(0, HORIZON - 30, 0, HORIZON + 46);
        g.addColorStop(0, 'rgba(' + s.pal.haze + ',0)'); g.addColorStop(0.4, 'rgba(' + s.pal.haze + ',0.55)'); g.addColorStop(1, 'rgba(' + s.pal.haze + ',0)');
        ctx.fillStyle = g; ctx.fillRect(0, HORIZON - 30, W, 76);
    }

    function drawRain(s, ctx) {
        var i, x, y, slant = s.wind * 14 - 5, t = s.clock;
        ctx.strokeStyle = '#9fb4ff'; ctx.globalAlpha = 0.4; ctx.lineWidth = 1.5;
        ctx.beginPath();
        for (i = 0; i < 90; i++) {
            x = (i * 97.3 + t * (120 + s.wind * 400)) % W; if (x < 0) x += W;
            y = G.HUD + (i * 53.7 + t * (1100 + i * 9)) % (H - G.HUD);
            ctx.moveTo(x, y); ctx.lineTo(x + slant, y + 24);
        }
        ctx.stroke();
        ctx.globalAlpha = 1;
    }

    function drawFlames(s, ctx, x, y) {
        var len = 26 + 10 * Math.sin(s.clock * 60), k;
        for (k = -1; k <= 1; k += 2) {
            ctx.fillStyle = '#35e0ff';
            ctx.beginPath(); ctx.moveTo(x + k * 38 - 9, y - 16); ctx.lineTo(x + k * 38 + 9, y - 16); ctx.lineTo(x + k * 38, y - 16 + len); ctx.fill();
            ctx.fillStyle = '#fff';
            ctx.beginPath(); ctx.moveTo(x + k * 38 - 4, y - 16); ctx.lineTo(x + k * 38 + 4, y - 16); ctx.lineTo(x + k * 38, y - 16 + len * 0.55); ctx.fill();
        }
    }

    function poly(ctx, pts, color) {
        ctx.fillStyle = color;
        ctx.beginPath();
        ctx.moveTo(pts[0], pts[1]);
        for (var i = 2; i < pts.length; i += 2) ctx.lineTo(pts[i], pts[i + 1]);
        ctx.fill();
    }

    function drawPlayer(s, ctx) {
        var pct = s.speed / MAX, x = W / 2, braking = G.key.down && s.speed > 0;
        // The cabin leans into the steering; on oil the whole car wobbles.
        var lean = G.clamp(s.vx / STEER, -1, 1) * 14 + (s.slip > 0 ? Math.sin(s.clock * 22) * 12 : 0);
        var y = PLAYER_Y + Math.sin(s.clock * 38) * (0.6 + pct * 1.2) + (s.off ? Math.sin(s.clock * 71) * 3 * Math.min(1, pct * 3) : 0);
        ctx.fillStyle = 'rgba(0,0,0,0.45)';
        ctx.beginPath(); ctx.ellipse(x, PLAYER_Y + 3, 96, 9, 0, 0, Math.PI * 2); ctx.fill();
        if (s.boost) drawFlames(s, ctx, x, y);
        ctx.fillStyle = '#05010c';
        ctx.fillRect(x - 86, y - 26, 24, 26); ctx.fillRect(x + 62, y - 26, 24, 26);
        poly(ctx, [x - 62 + lean, y - 86, x + 62 + lean, y - 86, x + 76, y - 50, x - 76, y - 50], '#bf3fff');
        poly(ctx, [x - 52 + lean, y - 79, x + 52 + lean, y - 79, x + 62, y - 54, x - 62, y - 54], '#1a0636');
        poly(ctx, [x - 86, y - 14, x + 86, y - 14, x + 80, y - 52, x - 80, y - 52], '#ff2d78');
        ctx.fillStyle = '#7a1fd0'; ctx.fillRect(x - 90 + lean * 0.5, y - 64, 180, 6);
        ctx.fillStyle = braking ? '#fff2f2' : '#ff6b2b'; ctx.fillRect(x - 74, y - 44, 148, 9);
        ctx.fillStyle = '#3a0a4a'; ctx.fillRect(x - 80, y - 26, 160, 9);
        ctx.fillStyle = '#35e0ff'; ctx.fillRect(x - 82, y - 15, 164, 2);
    }

    function bar(ctx, x, y, w, h, fill, color) {
        ctx.fillStyle = 'rgba(13,2,33,0.7)'; ctx.fillRect(x, y, w, h);
        ctx.fillStyle = color; ctx.fillRect(x, y, w * G.clamp(fill, 0, 1), h);
        ctx.strokeStyle = 'rgba(255,255,255,0.5)'; ctx.lineWidth = 1; ctx.strokeRect(x + 0.5, y + 0.5, w, h);
    }

    function drawDash(s, ctx) {
        var kmh = Math.round(s.speed / MAX * 290), k, x;
        G.text(String(kmh), W - 96, H - 34, { size: 38, bold: true, color: '#fff', align: 'right', glow: '#ff2d78' });
        G.text('KM/H', W - 24, H - 34, { size: 14, color: '#ffd9f0', align: 'right' });
        G.text('NITRO', W - 196, H - 13, { size: 12, color: '#35e0ff', align: 'right' });
        bar(ctx, W - 190, H - 23, 166, 10, s.nitro, s.boost ? '#fff' : '#35e0ff');
        // Course map: how far to the finish, with the checkpoints marked.
        bar(ctx, 24, H - 23, 220, 10, (s.pos + PZ) / (s.finish * SEG), '#ff2d78');
        ctx.fillStyle = '#ffd24a';
        for (k = 1; k <= s.L.cps; k++) { x = 24 + 220 * k / (s.L.cps + 1); ctx.fillRect(x - 1, H - 27, 2, 18); }
    }

    function drawClock(s, ctx) {
        var low = s.time < 10, blink = low && Math.floor(s.clock * 4) % 2;
        G.text('TIME', W / 2, 50, { size: 13, color: '#ffd9f0', align: 'center' });
        G.text(String(Math.max(0, Math.ceil(s.time))), W / 2, 92, { size: 44, bold: true, color: blink ? '#fff' : (low ? '#ff2d55' : '#ffd24a'), align: 'center', glow: low ? '#ff2d55' : '#ff6b2b' });
        if (s.msgT > 0) {
            ctx.globalAlpha = Math.min(1, s.msgT * 2);
            G.text(s.msg, W / 2, 140, { size: 26, bold: true, color: '#35e0ff', align: 'center', glow: '#35e0ff' });
            ctx.globalAlpha = 1;
        }
        if (s.gust > 0) {
            G.text(s.gustDir > 0 ? 'CROSSWIND  > > >' : '< < <  CROSSWIND', W / 2, 172, { size: 18, bold: true, color: '#ffd24a', align: 'center' });
        } else if (s.speed < 200 && s.pos < SEG * 4) {
            // Pulses, so the start line is never a frozen picture.
            ctx.globalAlpha = 0.55 + 0.45 * Math.sin(s.clock * 6);
            G.text('HOLD  ↑  TO DRIVE', W / 2, 172, { size: 18, color: '#fff', align: 'center' });
            ctx.globalAlpha = 1;
        }
    }

    function draw(s, ctx) {
        var cam = camera(s), base = Math.floor(cam.z / SEG);
        drawSky(s, ctx);
        drawRoad(s, ctx, cam);
        drawGrid(s, ctx, base);
        drawObjects(s, ctx, cam, base);
        drawHaze(s, ctx);
        if (s.L.storm) drawRain(s, ctx);
        drawPlayer(s, ctx);
        drawDash(s, ctx);
        drawClock(s, ctx);
    }

    function hud(s) {
        return 'CHECKPOINT ' + s.cp + '/' + s.L.cps;
    }

    G.register('synthwave', {
        title: 'OUTRUN THE SUN',
        blurb: 'Reach every checkpoint before the sun sets. Lift off for the hard bends.',
        controls: [
            '← → steer · ↑ accelerate · ↓ brake',
            'SPACE nitro — near misses and checkpoints refill the tank',
            'Leaving the road, hitting traffic or roadside objects costs speed'
        ],
        levelNames: ['Sunset Strip', 'Palm Slalom', 'Rolling Hills', 'Neon City', 'Oil Coast',
            'Canyon Hairpins', 'Two-Way Highway', 'Roadworks', 'Midnight Storm', 'Outrun the Sun'],
        colors: { bg: '#0d0221', fg: '#ffffff', accent: '#ff2d78', dim: '#bf3fff' },
        // One life: when the clock runs out the race is over, as in the arcade.
        lives: 1,
        music: {
            bpm: 118, root: 45, scale: 'minor', prog: [0, 5, 2, 6],
            bass: 'x.xox.xox.xox.xo',
            lead: ['4---..2.4---7...', '5---..7.9---7.5.', '6---..4.6---9...', '8---..6.8-a-8-6.',
                '7-4-7-9-b---9.7.', '9-7-5-7-9---c---', 'b-9-6-9-b---d-b-', 'a---8---6-8-a...'],
            arp: '0213',
            // Four on the floor with the big snare on two and four.
            drums: { k: 'x...x...x...x...', s: '....x.......x...', h: '..x...x...x...xx' },
            leadWave: 'sawtooth', bassWave: 'sawtooth', arpWave: 'square', leadOct: 2
        },
        init: init, update: update, draw: draw, hud: hud
    });
})();
