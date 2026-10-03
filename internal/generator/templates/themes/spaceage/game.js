/*
 * Orbital Dock — the spaceage theme's game: a thrust lander.
 *
 * Rotate the capsule, burn the engine against gravity and set it down on the
 * pad slowly and upright. Rock, flak, mines and an empty tank all end the
 * flight. The ten levels add canyons, wind, a moving ferry deck, caves with
 * fuel canisters, gun turrets, cargo that needs two landings, a low-gravity
 * moon, a heavy world and finally a rotating starbase to dock with.
 *
 * Every pad (ground, hovering or orbiting) is described by a centre, a normal
 * angle and a velocity, so one touchdown test serves all of them.
 */
(function () {
    'use strict';
    var G = window.SnoGame;
    var SEG = 24;                    // horizontal spacing of the terrain points
    var R = 12;                      // ship radius: centre to the feet
    var THRUST = 170, TURN = 3.2, BURN = 8, TANK = 100, CAN_FUEL = 35;
    var HAUL = 0.88;                 // thrust left over while carrying cargo
    var LIMIT = { down: 60, side: 40, tilt: 0.35 };   // a softer touch than this lands
    var TURRET = { range: 420, period: 2.2, shell: 150, lead: 0.6 };
    var C = { bg: '#030a0f', sky: '#06202b', teal: '#00e8e8', dim: '#1a4455', red: '#ff3320', silver: '#c8d8e0', rock: '#082230', far: '#05171f' };

    // One entry per level. floor/ceil are skeletons of [x, y] points that get
    // seeded jaggedness; ground pads are {i, n}: first terrain point and
    // width in segments, so their edges always sit on the terrain grid.
    var LEVELS = [
        { w: 960, g: 55, fuel: 100, start: [220, 90, 30], amp: 10, hint: 'LAND GENTLY AND UPRIGHT ON THE PAD',
            floor: [[0, 430], [120, 470], [260, 440], [380, 490], [520, 455], [624, 470], [744, 470], [840, 420], [960, 400]],
            pads: [{ i: 26, n: 5, kind: 'goal' }] },
        { w: 960, g: 60, fuel: 90, start: [140, 60, 20], amp: 8, hint: 'THE PAD LIES AT THE BOTTOM OF THE CANYON',
            floor: [[0, 300], [150, 260], [260, 330], [330, 220], [400, 250], [456, 500], [528, 500], [590, 260], [680, 200], [780, 330], [860, 280], [960, 330]],
            pads: [{ i: 19, n: 3, kind: 'goal' }] },
        { w: 960, g: 58, fuel: 100, start: [760, 100, -20], amp: 10, wind: [26, 10], hint: 'CROSSWIND — LEAN INTO THE GUSTS',
            floor: [[0, 470], [100, 430], [200, 480], [300, 380], [384, 340], [456, 340], [540, 420], [650, 470], [760, 400], [860, 450], [960, 410]],
            pads: [{ i: 16, n: 3, kind: 'goal' }] },
        { w: 960, g: 58, fuel: 100, start: [130, 110, 20], amp: 14, hint: 'MATCH THE FERRY DECK, THEN SET DOWN',
            floor: [[0, 380], [120, 420], [200, 500], [330, 470], [420, 515], [560, 480], [680, 515], [800, 460], [880, 400], [960, 370]],
            pads: [{ x: 500, y: 330, w: 96, sway: [150, 0.3], kind: 'goal' }] },
        { w: 1920, g: 55, fuel: 40, start: [100, 150, 30], amp: 10, hint: 'LOW TANK — GRAB THE FUEL CANISTERS',
            ceil: [[0, 60], [200, 80], [400, 150], [520, 90], [700, 60], [860, 200], [1000, 230], [1150, 120], [1300, 70], [1500, 160], [1640, 100], [1800, 60], [1920, 60]],
            floor: [[0, 330], [200, 360], [400, 400], [520, 330], [700, 300], [860, 440], [1000, 470], [1150, 400], [1300, 300], [1500, 400], [1640, 430], [1800, 430], [1920, 380]],
            cans: [[480, 210], [980, 350], [1380, 220]], pads: [{ i: 70, n: 4, kind: 'goal' }] },
        { w: 1440, g: 58, fuel: 100, start: [120, 100, 30], amp: 10, hint: 'FLAK TURRETS — NEVER FLY A STRAIGHT LINE',
            floor: [[0, 420], [150, 450], [300, 380], [450, 470], [600, 440], [720, 490], [850, 400], [1000, 460], [1130, 430], [1250, 470], [1344, 470], [1440, 400]],
            turrets: [14, 31, 43], cans: [[700, 200]], pads: [{ i: 52, n: 4, kind: 'goal' }] },
        { w: 1440, g: 58, fuel: 100, start: [720, 100, 0], amp: 10, hint: 'FETCH THE CARGO IN THE WEST, DELIVER IT EAST',
            floor: [[0, 300], [96, 300], [192, 300], [300, 460], [450, 420], [600, 480], [760, 430], [900, 500], [1050, 440], [1200, 360], [1296, 360], [1440, 300]],
            cans: [[860, 300]], pads: [{ i: 4, n: 4, kind: 'cargo' }, { i: 50, n: 4, kind: 'goal' }] },
        { w: 1920, g: 24, fuel: 60, start: [90, 180, 20], amp: 7, hint: 'LOW GRAVITY — SMALL BURNS, MIND THE MINES',
            ceil: [[0, 80], [250, 120], [400, 260], [560, 300], [720, 170], [900, 90], [1080, 200], [1250, 300], [1420, 220], [1600, 120], [1760, 200], [1920, 220]],
            floor: [[0, 300], [250, 300], [400, 420], [560, 460], [720, 340], [900, 250], [1080, 360], [1250, 470], [1420, 390], [1600, 290], [1760, 400], [1832, 400], [1920, 380]],
            mines: [[640, 317, 45, 1.2], [1160, 330, 45, 1.5], [1510, 255, 45, 1.7]],
            cans: [[900, 170], [1420, 305]], pads: [{ i: 73, n: 3, kind: 'goal' }] },
        { w: 1440, g: 96, fuel: 100, start: [700, 90, 0], amp: 9, wind: [14, 6], hint: 'HEAVY WORLD — CARGO WEST, DELIVERY EAST',
            floor: [[0, 400], [144, 400], [240, 400], [350, 480], [500, 440], [650, 490], [800, 450], [950, 500], [1080, 460], [1176, 420], [1272, 420], [1440, 380]],
            turrets: [27], cans: [[780, 250], [1020, 280]], pads: [{ i: 6, n: 4, kind: 'cargo' }, { i: 49, n: 4, kind: 'goal' }] },
        { w: 960, g: 0, fuel: 100, start: [130, 120, 25], amp: 0, hint: 'DOCK WITH THE ROTATING ARM — FEET TOWARDS THE STATION',
            station: { x: 480, y: 290, r: 58, d: 122, spin: 0.3, g: 40, debris: [[205, -0.45, 0], [205, -0.45, Math.PI], [235, 0.35, 1]] },
            pads: [{ w: 64, orbit: true, kind: 'goal' }] }
    ];

    // ---------------------------------------------------------------- setup

    function skeletonY(points, x) {
        for (var k = 1; k < points.length; k++) {
            if (x <= points[k][0]) {
                var a = points[k - 1], b = points[k];
                return G.lerp(a[1], b[1], (x - a[0]) / (b[0] - a[0]));
            }
        }
        return points[points.length - 1][1];
    }

    // Turns a skeleton into one height per SEG pixels, roughened by the
    // level's seeded generator so the rock is jagged but always the same.
    function profile(points, cfg, rnd) {
        var line = [];
        for (var i = 0; i <= cfg.w / SEG; i++) line.push(skeletonY(points, i * SEG) + (rnd() * 2 - 1) * cfg.amp);
        return line;
    }

    // A ground pad flattens the terrain points beneath it; the others hover
    // (sway) or ride the station arm (orbit) and are placed by placePads().
    function buildPad(s, def) {
        var pad = { kind: def.kind, w: def.w, ground: false, cx: 0, cy: 0, ang: 0, vx: 0, vy: 0, crate: def.kind === 'cargo' };
        if (def.n) {
            pad.w = def.n * SEG; pad.ground = true;
            pad.cx = def.i * SEG + pad.w / 2; pad.cy = skeletonY(s.cfg.floor, pad.cx);
            for (var k = def.i; k <= def.i + def.n; k++) s.floor[k] = pad.cy;
        } else if (def.sway) {
            pad.sway = def.sway; pad.x0 = def.x; pad.cy = def.y;
        } else {
            pad.orbit = true;
        }
        return pad;
    }

    function placePads(s) {
        var st = s.cfg.station;
        s.pads.forEach(function (pad) {
            if (pad.sway) {
                pad.cx = pad.x0 + pad.sway[0] * Math.sin(s.t * pad.sway[1]);
                pad.vx = pad.sway[0] * pad.sway[1] * Math.cos(s.t * pad.sway[1]);
            } else if (pad.orbit) {
                var th = s.t * st.spin, c = Math.cos(th), sn = Math.sin(th);
                pad.cx = st.x + c * st.d; pad.cy = st.y + sn * st.d;
                pad.vx = -sn * st.d * st.spin; pad.vy = c * st.d * st.spin;
                pad.ang = th + Math.PI / 2;      // the pad's normal points away from the hub
            }
        });
    }

    function buildHazards(s, cfg) {
        s.turrets = (cfg.turrets || []).map(function (i, k) {
            return { x: i * SEG, y: s.floor[i], cool: 1.2 + k * 0.7, warned: false, aim: -Math.PI / 2 };
        });
        s.mines = (cfg.mines || []).map(function (m) { return { bob: m }; });
        if (cfg.station) cfg.station.debris.forEach(function (d) { s.mines.push({ orbit: d }); });
        s.cans = (cfg.cans || []).map(function (c) { return { x: c[0], y: c[1], taken: false }; });
    }

    function buildSky(rnd) {
        var stars = [];
        for (var i = 0; i < 70; i++) stars.push({ x: rnd() * G.W, y: G.HUD + rnd() * 400, r: 0.6 + rnd() * 1.4, ph: rnd() * 6.3 });
        return { stars: stars, planet: { x: 560 + rnd() * 300, y: 90 + rnd() * 90, r: 26 + rnd() * 30 } };
    }

    function camTarget(s) { return G.clamp(s.ship.x - G.W / 2, 0, s.cfg.w - G.W); }

    function init(level) {
        var cfg = LEVELS[level - 1], rnd = G.rng(level * 733);
        var s = {
            cfg: cfg, t: 0, dead: false, windX: 0, stranded: 0, shells: [],
            floor: cfg.floor ? profile(cfg.floor, cfg, rnd) : null,
            ceil: cfg.ceil ? profile(cfg.ceil, cfg, rnd) : null,
            ship: { x: cfg.start[0], y: cfg.start[1], vx: cfg.start[2], vy: 0, ang: 0, fuel: cfg.fuel, landed: -1, off: 0, cargo: false, burning: false },
            msg: { text: cfg.hint, t: 5 },
            timers: { engine: 0, beep: 0, ping: 0, wind: 0 }
        };
        s.pads = cfg.pads.map(function (def) { return buildPad(s, def); });
        s.needCargo = s.pads.some(function (p) { return p.kind === 'cargo'; });
        placePads(s);
        buildHazards(s, cfg);
        s.sky = buildSky(rnd);
        s.camX = camTarget(s);
        return s;
    }

    // -------------------------------------------------------------- helpers

    function angDiff(a, b) { return Math.atan2(Math.sin(a - b), Math.cos(a - b)); }

    function heightAt(line, x) {
        var i = G.clamp(Math.floor(x / SEG), 0, line.length - 2);
        return G.lerp(line[i], line[i + 1], x / SEG - i);
    }

    function windAt(s) {
        var w = s.cfg.wind;
        return w ? w[0] * Math.sin(s.t * 0.4) + w[1] * Math.sin(s.t * 1.7 + 1) : 0;
    }

    function minePos(s, m) {
        if (m.bob) return { x: m.bob[0], y: m.bob[1] + m.bob[2] * Math.sin(s.t * m.bob[3]) };
        var st = s.cfg.station, a = m.orbit[2] + s.t * m.orbit[1];
        return { x: st.x + Math.cos(a) * m.orbit[0], y: st.y + Math.sin(a) * m.orbit[0] };
    }

    // The pad the pilot has to reach next: the crate first, then the beacon.
    function target(s) {
        var want = s.needCargo && !s.ship.cargo ? 'cargo' : 'goal';
        for (var i = 0; i < s.pads.length; i++) if (s.pads[i].kind === want) return s.pads[i];
        return s.pads[0];
    }

    // The ship as the pad sees it: offset along the deck, height above it,
    // closing speed, sideways slip and tilt against the pad's normal.
    function padFrame(ship, pad) {
        var nx = Math.sin(pad.ang), ny = -Math.cos(pad.ang);
        var dx = ship.x - pad.cx, dy = ship.y - pad.cy, rvx = ship.vx - pad.vx, rvy = ship.vy - pad.vy;
        return {
            along: -dx * ny + dy * nx, h: dx * nx + dy * ny,
            down: -(rvx * nx + rvy * ny), side: Math.abs(-rvx * ny + rvy * nx),
            tilt: Math.abs(angDiff(ship.ang, pad.ang))
        };
    }

    function say(s, text) { s.msg = { text: text, t: 3 }; }

    // Two-tone mission-control beeps.
    function radio(notes) {
        notes.forEach(function (f, i) { G.tone(f, 0.09, { type: 'sine', vol: 0.12, delay: i * 0.11 }); });
    }

    // ------------------------------------------------------------ collisions

    // side is 1 for the floor (rock below the line) and -1 for the ceiling.
    function lineHit(line, x, y, r, side) {
        if (!line) return false;
        if ((y - heightAt(line, x)) * side > 0) return true;       // already inside the rock
        var i0 = Math.max(0, Math.floor((x - r) / SEG)), i1 = Math.min(line.length - 2, Math.floor((x + r) / SEG));
        for (var i = i0; i <= i1; i++) {
            var p = G.closestOnSeg(x, y, i * SEG, line[i], (i + 1) * SEG, line[i + 1]);
            if (G.dist(x, y, p.x, p.y) < r) return true;
        }
        return false;
    }

    function hitsRock(s, x, y, r) { return lineHit(s.floor, x, y, r, 1) || lineHit(s.ceil, x, y, r, -1); }

    // The hull ring and the arm are deadly; the arm stops short of the deck
    // so a ship standing on the pad does not touch it.
    function hitsStation(s) {
        var st = s.cfg.station, ship = s.ship;
        if (!st) return false;
        if (G.dist(ship.x, ship.y, st.x, st.y) < st.r + R - 1) return true;
        var th = s.t * st.spin, c = Math.cos(th), sn = Math.sin(th);
        var p = G.closestOnSeg(ship.x, ship.y, st.x + c * st.r, st.y + sn * st.r, st.x + c * (st.d - 9), st.y + sn * (st.d - 9));
        return G.dist(ship.x, ship.y, p.x, p.y) < R + 2;
    }

    function hitsMine(s) {
        return s.mines.some(function (m) {
            var p = minePos(s, m);
            return G.dist(s.ship.x, s.ship.y, p.x, p.y) < R + 7;
        });
    }

    // '' = clear of the pad, 'land' = a good touchdown, 'crash' = too fast,
    // too tilted, off the edge or from underneath. A ground pad's edges are
    // left to the terrain test; a hovering deck is solid all the way round.
    function padContact(ship, pad) {
        var f = padFrame(ship, pad), reach = pad.w / 2 + (pad.ground ? -5 : R * 0.6);
        if (Math.abs(f.along) > reach || Math.abs(f.h) > R) return '';
        var soft = f.down < LIMIT.down && f.side < LIMIT.side && f.tilt < LIMIT.tilt;
        return soft && f.h > 0 && Math.abs(f.along) <= pad.w / 2 - 5 ? 'land' : 'crash';
    }

    function crash(s) {
        var ship = s.ship;
        if (s.dead) return;
        s.dead = true;
        G.burst(ship.x, ship.y, { n: 36, color: C.red, speed: 260, life: 0.9, gravity: 180 });
        G.burst(ship.x, ship.y, { n: 16, color: C.silver, speed: 180, life: 1.2, gravity: 180, size: 4 });
        G.sfx('bigboom');
        G.noise(0.5, { filter: 'bandpass', freq: 900, slide: 120, vol: 0.3 });
        G.shake(12, 0.5); G.flash(C.red, 0.2);
        G.die();
    }

    // -------------------------------------------------------------- landing

    // Pins the ship to its pad. With lift it is released a hair above the
    // deck and pushed off, so the touchdown test does not catch it again.
    function seat(s, lift) {
        var ship = s.ship, pad = s.pads[ship.landed], nx = Math.sin(pad.ang), ny = -Math.cos(pad.ang);
        var h = R + (lift ? 1.5 : 0), kick = lift ? 30 : 0;
        ship.x = pad.cx - ny * ship.off + nx * h; ship.y = pad.cy + nx * ship.off + ny * h;
        ship.vx = pad.vx + nx * kick; ship.vy = pad.vy + ny * kick;
        ship.ang = pad.ang;
    }

    function touchDown(s, i) {
        var ship = s.ship, pad = s.pads[i];
        ship.off = padFrame(ship, pad).along; ship.landed = i; ship.burning = false;
        seat(s, false);
        G.noise(0.12, { freq: 300, vol: 0.3 });
        G.burst(ship.x, ship.y + R, { n: 10, color: C.silver, speed: 90, life: 0.5, angle: -Math.PI / 2, spread: Math.PI });
        radio([1320, 1760]);
        if (pad.crate) {
            pad.crate = false; ship.cargo = true;
            G.sfx('power'); G.addScore(200);
            G.popup(ship.x, ship.y - 26, 'CARGO +200', C.teal);
            say(s, 'CARGO ABOARD — DELIVER IT TO THE BEACON');
        } else if (pad.kind === 'goal' && s.needCargo && !ship.cargo) {
            G.tone(140, 0.3, { type: 'square', vol: 0.12 });
            say(s, 'FETCH THE CARGO FIRST');
        } else if (pad.kind === 'goal') {
            radio([1320, 1760, 2093, 2637]);
            G.win(200 + Math.round(ship.fuel) * 5 + G.lives * 100);
        }
    }

    // While landed the ship rides its pad. Thrust lifts off; an empty tank
    // on a pad that is not the goal would strand the pilot for ever, so it
    // costs a life after a moment.
    function ride(s, dt) {
        var ship = s.ship, lift = (G.key.up || G.key.a) && ship.fuel > 0;
        seat(s, lift);
        if (lift) {
            ship.landed = -1;
            G.tone(180, 0.2, { type: 'sawtooth', slide: 420, vol: 0.1 });
            return;
        }
        if (ship.fuel > 0) return;
        if (s.stranded === 0) say(s, 'STRANDED — TANK DRY');
        s.stranded += dt;
        if (s.stranded > 1.5) G.die();
    }

    // --------------------------------------------------------------- flight

    function steer(s, dt) {
        var ship = s.ship;
        ship.ang += ((G.key.right ? 1 : 0) - (G.key.left ? 1 : 0)) * TURN * dt;
        if (G.hit.left || G.hit.right) G.noise(0.05, { filter: 'highpass', freq: 3000, vol: 0.06 });   // RCS puff
        ship.burning = (G.key.up || G.key.a) && ship.fuel > 0;
    }

    function burn(s, dt) {
        var ship = s.ship;
        ship.fuel = Math.max(0, ship.fuel - BURN * dt);
        if (ship.fuel > 0) return;
        say(s, 'TANK DRY');
        G.tone(520, 0.5, { type: 'sawtooth', slide: 90, vol: 0.14 });
    }

    function keepInBounds(s) {
        var ship = s.ship, w = s.cfg.w;
        if (ship.x < R) { ship.x = R; ship.vx = Math.max(0, ship.vx); }
        if (ship.x > w - R) { ship.x = w - R; ship.vx = Math.min(0, ship.vx); }
        if (ship.y < G.HUD + R) { ship.y = G.HUD + R; ship.vy = Math.max(0, ship.vy); }
        if (ship.y > G.H - R) { ship.y = G.H - R; ship.vy = Math.min(0, ship.vy); }
    }

    // Gravity pulls down, or towards the hub on the station level; wind is a
    // sideways acceleration; the engine pushes along the ship's nose.
    function fly(s, dt) {
        var ship = s.ship, st = s.cfg.station, ax = windAt(s), ay = s.cfg.g;
        if (st) {
            var d = Math.max(1, G.dist(ship.x, ship.y, st.x, st.y));
            ax += (st.x - ship.x) / d * st.g; ay += (st.y - ship.y) / d * st.g;
        }
        if (ship.burning) {
            var push = THRUST * (ship.cargo ? HAUL : 1);
            ax += Math.sin(ship.ang) * push; ay -= Math.cos(ship.ang) * push;
            burn(s, dt);
        }
        ship.vx += ax * dt; ship.vy += ay * dt;
        ship.x += ship.vx * dt; ship.y += ship.vy * dt;
        keepInBounds(s);
    }

    function checkContacts(s) {
        var ship = s.ship;
        for (var i = 0; i < s.pads.length; i++) {
            var hit = padContact(ship, s.pads[i]);
            if (hit === 'land') { touchDown(s, i); return; }
            if (hit === 'crash') { crash(s); return; }
        }
        if (hitsRock(s, ship.x, ship.y, R - 1) || hitsStation(s)) crash(s);
    }

    function collectCans(s) {
        var ship = s.ship;
        s.cans.forEach(function (can) {
            if (can.taken || G.dist(ship.x, ship.y, can.x, can.y) > R + 14) return;
            can.taken = true;
            ship.fuel = Math.min(TANK, ship.fuel + CAN_FUEL);
            G.addScore(50);
            G.sfx('coin'); G.tone(330, 0.25, { type: 'triangle', slide: 660, vol: 0.12 });
            G.popup(can.x, can.y - 14, 'FUEL +' + CAN_FUEL, C.teal);
        });
    }

    // -------------------------------------------------------------- hazards

    // A turret tracks the ship, whines half a second before it fires and
    // holds its fire while the ship is out of range. It leads its target by
    // part of the shell's flight time: a full lead would be unfair, none
    // would never hit a moving ship.
    function updateTurrets(s, dt) {
        var ship = s.ship;
        s.turrets.forEach(function (tu) {
            var my = tu.y - 14, lead = G.dist(tu.x, my, ship.x, ship.y) / TURRET.shell * TURRET.lead;
            tu.aim = Math.atan2(ship.y + ship.vy * lead - my, ship.x + ship.vx * lead - tu.x);
            if (G.dist(tu.x, my, ship.x, ship.y) > TURRET.range) { tu.cool = Math.max(tu.cool, 0.8); tu.warned = false; return; }
            tu.cool -= dt;
            if (tu.cool < 0.5 && !tu.warned) { tu.warned = true; G.tone(420, 0.4, { type: 'sawtooth', slide: 900, vol: 0.07 }); }
            if (tu.cool > 0) return;
            tu.cool = TURRET.period; tu.warned = false;
            s.shells.push({ x: tu.x + Math.cos(tu.aim) * 16, y: my + Math.sin(tu.aim) * 16, vx: Math.cos(tu.aim) * TURRET.shell, vy: Math.sin(tu.aim) * TURRET.shell, life: 4 });
            G.sfx('shoot'); G.noise(0.1, { freq: 700, vol: 0.15 });
        });
    }

    function updateShells(s, dt) {
        var ship = s.ship;
        s.shells = s.shells.filter(function (sh) {
            sh.x += sh.vx * dt; sh.y += sh.vy * dt; sh.life -= dt;
            if (G.dist(sh.x, sh.y, ship.x, ship.y) < R + 3) { crash(s); return false; }
            if (hitsRock(s, sh.x, sh.y, 2)) {
                G.burst(sh.x, sh.y, { n: 5, color: C.red, speed: 70, life: 0.3 });
                return false;
            }
            return sh.life > 0;
        });
    }

    // ------------------------------------------------------------- ambience

    function engineSound(s) {
        var ship = s.ship, back = Math.atan2(Math.cos(ship.ang), -Math.sin(ship.ang));
        s.timers.engine = 0.09;
        G.noise(0.14, { freq: 260, vol: 0.14 });
        G.tone(46 + Math.random() * 8, 0.12, { type: 'sawtooth', vol: 0.06 });
        G.burst(ship.x - Math.sin(ship.ang) * R, ship.y + Math.cos(ship.ang) * R, { n: 2, color: C.red, speed: 150, life: 0.35, angle: back, spread: 0.5 });
    }

    // Continuous sounds are gated by timers: engine rumble, the low-fuel
    // warning, the radar altimeter that pings faster near the pad, and gusts.
    function ambience(s, dt) {
        var ship = s.ship, tm = s.timers, wind = Math.abs(windAt(s)), f = padFrame(ship, target(s));
        for (var k in tm) if (tm[k] > 0) tm[k] -= dt;
        if (s.msg.t > 0) s.msg.t -= dt;
        s.windX += windAt(s) * dt * 6;
        if (ship.burning && tm.engine <= 0) engineSound(s);
        if (ship.fuel > 0 && ship.fuel < 20 && tm.beep <= 0) {
            tm.beep = 0.7;
            G.tone(1046, 0.07, { vol: 0.07 });
        }
        if (ship.landed < 0 && tm.ping <= 0 && f.h > 0 && f.h < 170 && Math.abs(f.along) < target(s).w / 2 + 40) {
            tm.ping = G.clamp(f.h / 170, 0.12, 0.8);
            G.tone(1500, 0.03, { type: 'sine', vol: 0.09 });
        }
        if (wind > 10 && tm.wind <= 0) {
            tm.wind = 1.4;
            G.noise(1.2, { filter: 'bandpass', freq: 400, slide: 900, vol: 0.003 * wind, attack: 0.4 });
        }
    }

    function update(s, dt) {
        var ship = s.ship;
        s.t += dt;
        placePads(s);
        if (ship.landed >= 0) ride(s, dt);
        else { steer(s, dt); fly(s, dt); }
        collectCans(s);
        updateTurrets(s, dt);
        updateShells(s, dt);
        if (ship.landed < 0) checkContacts(s);
        if (hitsMine(s)) crash(s);
        s.camX += (camTarget(s) - s.camX) * Math.min(1, dt * 6);
        ambience(s, dt);
    }

    // -------------------------------------------------------------- drawing

    function drawSky(ctx, s) {
        var grad = ctx.createLinearGradient(0, 0, 0, G.H), p = s.sky.planet, px = p.x - s.camX * 0.08;
        grad.addColorStop(0, C.bg); grad.addColorStop(1, C.sky);
        ctx.fillStyle = grad; ctx.fillRect(0, 0, G.W, G.H);
        s.sky.stars.forEach(function (st) {
            var x = (((st.x - s.camX * 0.2) % G.W) + G.W) % G.W;
            ctx.globalAlpha = 0.45 + 0.4 * Math.sin(s.t * 2 + st.ph);
            ctx.fillStyle = C.silver; ctx.fillRect(x, st.y, st.r, st.r);
        });
        ctx.globalAlpha = 1;
        if (s.ceil) return;          // no sky to see from inside a cave
        // A ringed planet far behind everything, in the space-age poster style.
        ctx.globalAlpha = 0.5;
        ctx.fillStyle = C.dim;
        ctx.beginPath(); ctx.arc(px, p.y, p.r, 0, Math.PI * 2); ctx.fill();
        ctx.strokeStyle = C.teal; ctx.lineWidth = 1.5;
        ctx.beginPath(); ctx.ellipse(px, p.y, p.r * 1.8, p.r * 0.35, -0.35, 0, Math.PI * 2); ctx.stroke();
        ctx.globalAlpha = 1;
    }

    // Distant ridges for the open-sky levels; they scroll at half speed.
    function drawFarHills(ctx, s) {
        if (!s.floor || s.ceil) return;
        ctx.fillStyle = C.far;
        ctx.beginPath(); ctx.moveTo(0, G.H);
        for (var x = 0; x <= G.W; x += SEG) {
            var wx = x + s.camX * 0.5;
            ctx.lineTo(x, 380 + 45 * Math.sin(wx * 0.011 + s.cfg.g) + 22 * Math.sin(wx * 0.031));
        }
        ctx.lineTo(G.W, G.H); ctx.fill();
    }

    function trace(ctx, line, i0, i1, dy) {
        ctx.beginPath();
        for (var i = i0; i <= i1; i++) ctx.lineTo(i * SEG, line[i] + dy);
    }

    // edge is the screen border the rock grows from: G.H (floor) or 0 (ceiling).
    function drawRock(ctx, s, line, edge) {
        if (!line) return;
        var i0 = Math.max(0, Math.floor(s.camX / SEG)), i1 = Math.min(line.length - 1, Math.ceil((s.camX + G.W) / SEG));
        var inward = edge ? 9 : -9;
        trace(ctx, line, i0, i1, 0);
        ctx.lineTo(i1 * SEG, edge); ctx.lineTo(i0 * SEG, edge); ctx.closePath();
        ctx.fillStyle = C.rock; ctx.fill();
        trace(ctx, line, i0, i1, inward);                 // a dim inner contour gives the rock depth
        ctx.strokeStyle = C.dim; ctx.lineWidth = 2; ctx.stroke();
        trace(ctx, line, i0, i1, 0);
        ctx.strokeStyle = C.teal; ctx.lineWidth = 2;
        ctx.shadowColor = C.teal; ctx.shadowBlur = 8; ctx.stroke();
        ctx.shadowBlur = 0;
    }

    function drawStation(ctx, s) {
        var st = s.cfg.station;
        if (!st) return;
        ctx.save();
        ctx.translate(st.x, st.y); ctx.rotate(s.t * st.spin);
        ctx.fillStyle = C.dim; ctx.fillRect(st.r - 2, -3, st.d - st.r - 4, 6);        // docking arm
        ctx.strokeStyle = C.silver; ctx.lineWidth = 8;
        ctx.beginPath(); ctx.arc(0, 0, st.r - 4, 0, Math.PI * 2); ctx.stroke();
        ctx.lineWidth = 3;
        for (var k = 0; k < 4; k++) {
            ctx.rotate(Math.PI / 2);
            ctx.beginPath(); ctx.moveTo(16, 0); ctx.lineTo(st.r - 8, 0); ctx.stroke();
            ctx.fillStyle = Math.sin(s.t * 3 + k) > 0 ? C.teal : C.dim;
            ctx.fillRect(st.r * 0.68, st.r * 0.68, 4, 4);                              // ring windows
        }
        ctx.fillStyle = C.dim; ctx.strokeStyle = C.teal; ctx.lineWidth = 2;
        ctx.beginPath(); ctx.arc(0, 0, 16, 0, Math.PI * 2); ctx.fill(); ctx.stroke();
        ctx.restore();
    }

    function drawCrate(ctx, y) {
        ctx.fillStyle = C.bg; ctx.strokeStyle = C.teal; ctx.lineWidth = 2;
        ctx.fillRect(-8, y - 8, 16, 16); ctx.strokeRect(-8, y - 8, 16, 16);
        ctx.beginPath(); ctx.moveTo(-8, y - 8); ctx.lineTo(8, y + 8); ctx.moveTo(8, y - 8); ctx.lineTo(-8, y + 8); ctx.stroke();
    }

    // Drawn in the pad's own frame (x along the deck, -y along its normal),
    // so ground, ferry and station pads share the code.
    function drawPad(ctx, s, pad) {
        var next = pad === target(s), blink = Math.sin(s.t * 6) > 0, bob = 4 * Math.sin(s.t * 4), half = pad.w / 2;
        ctx.save();
        ctx.translate(pad.cx, pad.cy); ctx.rotate(pad.ang);
        ctx.fillStyle = C.silver; ctx.fillRect(-half, 0, pad.w, 5);
        ctx.fillStyle = C.dim; ctx.fillRect(-half + 6, 5, pad.w - 12, 4);
        ctx.fillStyle = blink ? (next ? C.teal : C.red) : C.dim;
        ctx.fillRect(-half, -4, 5, 4); ctx.fillRect(half - 5, -4, 5, 4);
        if (pad.sway) {                                    // hover jets keep the ferry deck up
            ctx.fillStyle = C.red;
            ctx.fillRect(-half + 10, 9, 6, 6 + Math.random() * 8); ctx.fillRect(half - 16, 9, 6, 6 + Math.random() * 8);
        }
        if (pad.crate) drawCrate(ctx, -34 + bob);
        else if (next) {                                   // a bobbing chevron marks the pad to reach
            ctx.fillStyle = C.teal;
            ctx.beginPath(); ctx.moveTo(-9, -52 + bob); ctx.lineTo(9, -52 + bob); ctx.lineTo(0, -40 + bob); ctx.fill();
        }
        ctx.restore();
    }

    function drawShip(ctx, s) {
        var ship = s.ship;
        ctx.save();
        ctx.translate(ship.x, ship.y); ctx.rotate(ship.ang);
        if (ship.burning) {
            ctx.fillStyle = C.red;
            ctx.beginPath(); ctx.moveTo(-4, 6); ctx.lineTo(4, 6); ctx.lineTo(0, 17 + Math.random() * 12); ctx.fill();
            ctx.fillStyle = C.silver;
            ctx.beginPath(); ctx.moveTo(-2, 6); ctx.lineTo(2, 6); ctx.lineTo(0, 11 + Math.random() * 6); ctx.fill();
        }
        ctx.strokeStyle = C.silver; ctx.lineWidth = 2;
        ctx.beginPath();
        ctx.moveTo(-6, 4); ctx.lineTo(-11, 12); ctx.moveTo(6, 4); ctx.lineTo(11, 12);          // legs
        ctx.moveTo(-14, 12); ctx.lineTo(-8, 12); ctx.moveTo(8, 12); ctx.lineTo(14, 12);        // feet
        ctx.stroke();
        ctx.fillStyle = C.silver;
        ctx.beginPath(); ctx.moveTo(0, -15); ctx.quadraticCurveTo(9, -8, 7, 5); ctx.lineTo(-7, 5); ctx.quadraticCurveTo(-9, -8, 0, -15); ctx.fill();
        ctx.fillStyle = C.red; ctx.fillRect(-7, 0, 14, 2);
        ctx.fillStyle = C.teal;
        ctx.beginPath(); ctx.arc(0, -5, 3, 0, Math.PI * 2); ctx.fill();
        if (ship.cargo) { ctx.fillRect(-4, 6, 8, 6); ctx.strokeRect(-4, 6, 8, 6); }             // crate slung between the legs
        ctx.restore();
    }

    function drawMine(ctx, s, m) {
        var p = minePos(s, m);
        ctx.save();
        ctx.translate(p.x, p.y); ctx.rotate(s.t * 1.5);
        ctx.strokeStyle = C.red; ctx.lineWidth = 2;
        ctx.beginPath(); ctx.moveTo(-11, 0); ctx.lineTo(11, 0); ctx.moveTo(0, -11); ctx.lineTo(0, 11); ctx.stroke();
        ctx.fillStyle = C.bg;
        ctx.beginPath(); ctx.arc(0, 0, 7, 0, Math.PI * 2); ctx.fill(); ctx.stroke();
        ctx.fillStyle = Math.sin(s.t * 8) > 0 ? C.red : C.dim; ctx.fillRect(-2, -2, 4, 4);
        ctx.restore();
    }

    function drawTurret(ctx, tu) {
        var hot = tu.cool < 0.5;
        ctx.fillStyle = C.dim;
        ctx.beginPath(); ctx.moveTo(tu.x - 12, tu.y + 2); ctx.lineTo(tu.x - 7, tu.y - 14); ctx.lineTo(tu.x + 7, tu.y - 14); ctx.lineTo(tu.x + 12, tu.y + 2); ctx.fill();
        ctx.strokeStyle = hot ? C.red : C.silver; ctx.lineWidth = 4;
        ctx.beginPath(); ctx.moveTo(tu.x, tu.y - 14); ctx.lineTo(tu.x + Math.cos(tu.aim) * 16, tu.y - 14 + Math.sin(tu.aim) * 16); ctx.stroke();
        ctx.fillStyle = hot ? C.red : C.silver;
        ctx.beginPath(); ctx.arc(tu.x, tu.y - 14, 5, 0, Math.PI * 2); ctx.fill();
    }

    function drawThings(ctx, s) {
        s.cans.forEach(function (can, k) {
            if (can.taken) return;
            var y = can.y + 3 * Math.sin(s.t * 3 + k);
            ctx.fillStyle = C.teal; ctx.fillRect(can.x - 7, y - 9, 14, 18);
            ctx.fillStyle = C.silver; ctx.fillRect(can.x - 3, y - 12, 6, 3);
            G.text('F', can.x, y + 5, { size: 12, bold: true, color: C.bg, align: 'center' });
        });
        s.turrets.forEach(function (tu) { drawTurret(ctx, tu); });
        s.mines.forEach(function (m) { drawMine(ctx, s, m); });
        ctx.fillStyle = C.red; ctx.shadowColor = C.red; ctx.shadowBlur = 10;
        s.shells.forEach(function (sh) { ctx.beginPath(); ctx.arc(sh.x, sh.y, 3, 0, Math.PI * 2); ctx.fill(); });
        ctx.shadowBlur = 0;
    }

    // Streaks drift with the wind so its direction and strength can be read.
    function drawWind(ctx, s) {
        var wind = windAt(s), span = G.W + 80;
        if (!s.cfg.wind) return;
        ctx.strokeStyle = C.silver; ctx.lineWidth = 1; ctx.globalAlpha = 0.3;
        ctx.beginPath();
        for (var k = 0; k < 16; k++) {
            var x = (((k * 137 + s.windX) % span) + span) % span - 40, y = 50 + (k * 89) % 400;
            ctx.moveTo(x, y); ctx.lineTo(x - wind * 1.4, y);
        }
        ctx.stroke(); ctx.globalAlpha = 1;
    }

    // The instrument panel: closing speed, side slip and tilt against the
    // target pad, each with its limit mark, and the tank. Red means that
    // touching down now would wreck the ship.
    function drawPanel(ctx, s) {
        var ship = s.ship, f = padFrame(ship, target(s)), landed = ship.landed >= 0;
        var rows = [['DESC', f.down / LIMIT.down], ['SIDE', f.side / LIMIT.side], ['TILT', f.tilt / LIMIT.tilt]];
        ctx.fillStyle = 'rgba(2,6,8,0.72)'; ctx.fillRect(10, 38, 150, 68);
        ctx.strokeStyle = C.dim; ctx.lineWidth = 1; ctx.strokeRect(10.5, 38.5, 150, 68);
        rows.forEach(function (row, k) {
            var y = 46 + k * 15, v = landed ? 0 : Math.max(0, row[1]);
            G.text(row[0], 16, y + 8, { size: 10, color: C.silver });
            ctx.fillStyle = C.dim; ctx.fillRect(52, y, 100, 8);
            ctx.fillStyle = v < 1 ? C.teal : C.red; ctx.fillRect(52, y, Math.min(100, v * 50), 8);
            ctx.fillStyle = C.silver; ctx.fillRect(101, y - 2, 2, 12);                 // the limit
        });
        G.text('FUEL', 16, 99, { size: 10, color: C.silver });
        ctx.fillStyle = C.dim; ctx.fillRect(52, 91, 100, 8);
        ctx.fillStyle = ship.fuel < 20 ? C.red : C.teal; ctx.fillRect(52, 91, ship.fuel / TANK * 100, 8);
    }

    // The mission line at the top, and an arrow at the screen edge while the
    // target pad is scrolled out of sight.
    function drawGuide(ctx, s) {
        var pad = target(s), sx = pad.cx - s.camX, off = sx < 0 ? -1 : (sx > G.W ? 1 : 0);
        if (s.msg.t > 0) {
            ctx.globalAlpha = Math.min(1, s.msg.t);
            G.text(s.msg.text, G.W / 2, 132, { size: 15, color: C.teal, align: 'center', glow: C.teal, max: G.W - 60 });
            ctx.globalAlpha = 1;
        }
        if (!off || Math.sin(s.t * 6) < -0.4) return;
        var x = off < 0 ? 16 : G.W - 16, y = G.clamp(pad.cy - 40, 150, G.H - 40);
        ctx.fillStyle = C.teal;
        ctx.beginPath(); ctx.moveTo(x + off * 10, y); ctx.lineTo(x - off * 6, y - 10); ctx.lineTo(x - off * 6, y + 10); ctx.fill();
        G.text(Math.round(Math.abs(pad.cx - s.ship.x) / 10) + 'm', x - off * 14, y + 5, { size: 12, color: C.teal, align: off < 0 ? 'left' : 'right' });
    }

    function draw(s, ctx) {
        drawSky(ctx, s);
        drawFarHills(ctx, s);
        G.cam.x = s.camX;            // engine particles and popups scroll with the world
        ctx.save();
        ctx.translate(-s.camX, 0);
        drawRock(ctx, s, s.ceil, 0);
        drawRock(ctx, s, s.floor, G.H);
        drawStation(ctx, s);
        s.pads.forEach(function (pad) { drawPad(ctx, s, pad); });
        drawThings(ctx, s);
        if (!s.dead) drawShip(ctx, s);
        ctx.restore();
        drawWind(ctx, s);
        drawPanel(ctx, s);
        drawGuide(ctx, s);
    }

    function hud(s) {
        return 'FUEL ' + Math.ceil(s.ship.fuel) + (s.ship.cargo ? ' +CARGO' : '');
    }

    G.register('spaceage', {
        title: 'ORBITAL DOCK',
        blurb: 'Fly the capsule to the pad and touch down slowly and upright.',
        controls: [
            '← → rotate · ↑ or SPACE fire the engine',
            'Land with DESC, SIDE and TILT all below their marks',
            'Fuel is limited: canisters refill the tank',
            'Cargo levels need two landings: crate first, then the beacon'
        ],
        levelNames: ['First Contact', 'Razor Canyon', 'Crosswind', 'Ferry Deck', 'Fuel Run', 'Flak Alley', 'Cargo Haul', 'Selene Caverns', 'Iron Giant', 'Starbase Halo'],
        colors: { bg: C.bg, fg: C.silver, accent: C.teal, dim: '#5d8a9a' },
        lives: 3,
        // A slow lydian tune: long sine notes over a sparse pulse, calm
        // enough to concentrate on a landing.
        music: {
            bpm: 84, root: 45, scale: 'lydian', prog: [0, 4, 1, 0, 5, 4, 1, 0],
            bass: 'x.......5.....x.',
            lead: [
                '4-----7---6-4---', '6-----8---6---4-', '5-----3---5-8---', '7-------4---2---',
                '9-----7---5-7---', '8-----b---8-6---', 'a-----8---5-3---', '4-----------....'
            ],
            arp: '0...1...2...1.3.',
            drums: { k: 'x.......x.......', h: '....x.......x.x.' },
            leadWave: 'sine', bassWave: 'triangle', arpWave: 'triangle', leadOct: 2
        },
        init: init, update: update, draw: draw, hud: hud
    });
})();
