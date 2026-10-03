/*
 * Formation — the invaders theme's game.
 *
 * A fixed shooter: a formation marches side to side, drops a step at every
 * wall and speeds up as it thins. Destroy all of it (and the mothership, on
 * the boss levels) before it lands. Bunkers erode cell by cell where anything
 * hits them. Later levels add divers that peel off and swoop, shielded
 * invaders, splitters that burst into loose minis, cloaking, lancers that
 * charge a beam, and falling capsules (rapid, double, shield).
 */
(function () {
    'use strict';
    var G = window.SnoGame;
    var C = { bg: '#030706', mint: '#9fffd8', white: '#edfff8', dim: '#78988b', line: '#24473b', alert: '#ffcc70', deep: '#3b6e5c', soft: '#62c9a4' };
    var PX = 3, CW = 46, RH = 36, HALF_W = 16, HALF_H = 12;      // sprite pixel, slot spacing, invader half size
    var PY = 496, LAND_Y = 470, GROUND_Y = 516, MARGIN = 26, DROP = 24;
    var SHOT_V = 720, CHARGE = 1.0, BEAM = 0.4, FADE = 0.4;
    var BUNK = { cols: 14, rows: 7, cell: 6, y: 408 };
    var MARCH = [146.83, 130.81, 116.54, 110];                   // D C Bb A: the descending four-note march
    var CAPS = ['rapid', 'double', 'shield'];

    // hp, score and look per invader kind. 's' is shielded (the first hit only
    // pops the shield), 'x' splits into two minis ('m'), 'l' is a lancer.
    var KIND = {
        a: { hp: 1, pts: 10, color: C.soft, spr: 'a' }, b: { hp: 1, pts: 20, color: C.mint, spr: 'b' },
        c: { hp: 1, pts: 30, color: C.white, spr: 'c' }, s: { hp: 2, pts: 40, color: C.mint, spr: 'b' },
        x: { hp: 1, pts: 30, color: C.white, spr: 'x' }, l: { hp: 2, pts: 60, color: C.alert, spr: 'l' },
        m: { hp: 1, pts: 15, color: C.white, spr: 'm' }
    };

    // Sprite bitmaps: a shared body plus two leg variants, so every invader
    // steps in time with the march.
    var ART = {
        a: [['..x.....x..', '...x...x...', '..xxxxxxx..', '.xx.xxx.xx.', 'xxxxxxxxxxx', 'x.xxxxxxx.x'], ['x.x.....x.x', '...xx.xx...'], ['..x.....x..', '.x.......x.']],
        b: [['...xxxxx...', '.xxxxxxxxx.', 'xxxxxxxxxxx', 'xxx..x..xxx', 'xxxxxxxxxxx'], ['..xxx.xxx..', '.xx..x..xx.', 'xx.......xx'], ['...xx.xx...', '..xx.x.xx..', '..xx...xx..']],
        c: [['....xxx....', '...xxxxx...', '..xxxxxxx..', '.xx.xxx.xx.', '.xxxxxxxxx.'], ['...x.x.x...', '..x.....x..', '.x.x...x.x.'], ['..x.x.x.x..', '.x.......x.', '..x.....x..']],
        x: [['.xxx...xxx.', 'xxxxx.xxxxx', 'x.xxxxxxx.x', 'xxxxxxxxxxx', 'xxxxx.xxxxx', '.xxx...xxx.'], ['.x.x...x.x.', 'x...x.x...x'], ['x.x.....x.x', '.x..x.x..x.']],
        l: [['xxxxxxxxxxx', 'x.xxxxxxx.x', '.xxx.x.xxx.', '..xxxxxxx..', '..xx.x.xx..', '...xxxxx...'], ['....x.x....', '.....x.....'], ['....xxx....', '....x.x....']],
        m: [['.x...x.', 'xxxxxxx', 'x.xxx.x', 'xxxxxxx'], ['.x.x.x.'], ['x.x.x.x']]
    };

    // Turns bitmap rows into horizontal runs so a sprite costs a dozen
    // fillRects instead of one per pixel.
    function compile(rows) {
        var runs = [];
        rows.forEach(function (row, y) {
            for (var x = 0; x < row.length; x++) {
                if (row.charAt(x) !== 'x') continue;
                var x0 = x;
                while (row.charAt(x + 1) === 'x') x++;
                runs.push([x0, y, x - x0 + 1]);
            }
        });
        return { w: rows[0].length, h: rows.length, runs: runs };
    }

    var SPR = {};
    Object.keys(ART).forEach(function (k) {
        SPR[k] = [compile(ART[k][0].concat(ART[k][1])), compile(ART[k][0].concat(ART[k][2]))];
    });

    // One entry per level:
    //   map      formation shape, one KIND letter per slot ('.' = empty)
    //   top      y of the first row (default 84; lower on boss levels)
    //   split    first row of a second group that marches the other way
    //   lo, hi   march speed in px/s with the formation full / almost gone;
    //            lo is high enough that a player who never shoots is landed on
    //   bomb     mean seconds between formation bombs
    //   bombV    bomb fall speed in px/s
    //   maxBombs bombs in the air above which the formation holds fire
    //   bunkers  number of bunkers
    //   caps     true when kills and the saucer may drop capsules
    //   dive     mean seconds between divers; divers: how many may be out
    //   zig      share of bombs that weave
    //   cloak    true when a wave of invisibility sweeps the rows
    //   lance    mean seconds between lancer beams
    //   boss     mothership: hp, fan (seconds between bolt fans), n (bolts per
    //            fan), minis (seconds between mini launches), lance (seconds
    //            between its own beams once below half hp)
    var LV = [
        { map: ['cccccccc', 'bbbbbbbb', 'bbbbbbbb', 'aaaaaaaa', 'aaaaaaaa'], lo: 36, hi: 120, bomb: 1.5, bombV: 165, maxBombs: 2, bunkers: 4 },
        { map: ['cc.......cc', 'bbbb...bbbb', 'abbbb.bbbba', '.aaabbbaaa.', '..aaaaaaa..', '....aaa....'], lo: 38, hi: 130, bomb: 1.3, bombV: 175, maxBombs: 3, bunkers: 4, caps: true },
        { map: ['cccc...cccc', 'bbbbb.bbbbb', 'aaaaa.aaaaa', 'aaaaa.aaaaa', 'aaa.....aaa'], lo: 38, hi: 135, bomb: 1.25, bombV: 180, maxBombs: 3, bunkers: 4, caps: true, dive: 4.5, divers: 1 },
        { map: ['...ccccc...', '.bbbbbbbbb.', 'aaaaaaaaaaa', '.asasasasa.', '...sssss...'], lo: 40, hi: 140, bomb: 1.15, bombV: 190, maxBombs: 3, bunkers: 4, caps: true, dive: 4.2, divers: 1, zig: 0.25 },
        { map: ['sbbbbbbbs', 'bbbbbbbbb', 'aaaaaaaaa', 'a.a.a.a.a'], top: 144, lo: 36, hi: 120, bomb: 1.4, bombV: 190, maxBombs: 3, bunkers: 3, caps: true, dive: 5, divers: 1, zig: 0.2, boss: { hp: 60, fan: 2.6, n: 3 } },
        { map: ['..c.c.c.c..', '.xbxbxbxbx.', 'bbbbbbbbbbb', 'aaaaaaaaaaa', '.xaaxaxaax.'], lo: 40, hi: 145, bomb: 1.1, bombV: 195, maxBombs: 4, bunkers: 4, caps: true, dive: 4, divers: 2, zig: 0.25 },
        { map: ['ccccccccc', 'bsbsbsbsb', 'bbbbbbbbb', 'axaaxaaxa', 'aaaaaaaaa', 'a.a.a.a.a'], split: 3, lo: 36, hi: 140, bomb: 1.0, bombV: 200, maxBombs: 4, bunkers: 4, caps: true, dive: 3.6, divers: 2, zig: 0.3 },
        { map: ['ccccccccccc', 'bbbbbbbbbbb', '.sbsbsbsbs.', 'aaaaaaaaaaa', 'x.a.x.a.x.a'], cloak: true, lo: 40, hi: 150, bomb: 1.0, bombV: 205, maxBombs: 4, bunkers: 2, caps: true, dive: 3.4, divers: 2, zig: 0.3 },
        { map: ['l.c.l.l.c.l', 'bbsbbbbbsbb', 'aaaaaaaaaaa', '.xaasasaax.', '.aa.a.a.aa.'], lo: 42, hi: 155, bomb: 0.95, bombV: 210, maxBombs: 4, bunkers: 3, caps: true, dive: 3.2, divers: 2, zig: 0.35, lance: 4.5 },
        { map: ['l.s.l.s.l', 'bxbbbbbxb', 'aaaaaaaaa', '.a.a.a.a.'], top: 144, lo: 38, hi: 140, bomb: 1.1, bombV: 210, maxBombs: 4, bunkers: 2, caps: true, dive: 3.4, divers: 2, zig: 0.3, lance: 6, boss: { hp: 80, fan: 2.2, n: 5, minis: 7, lance: 5 } }
    ];

    // ------------------------------------------------------------------
    // Level setup
    // ------------------------------------------------------------------

    function bunkerCell(c, r) {
        var corner = (r === 0 && (c < 2 || c > 11)) || (r === 1 && (c < 1 || c > 12));
        var arch = (r >= 5 && c >= 5 && c <= 8) || (r === 4 && c >= 6 && c <= 7);
        return corner || arch ? 0 : 1;
    }

    function buildBunkers(n) {
        var out = [], w = BUNK.cols * BUNK.cell;
        for (var i = 0; i < n; i++) {
            var cells = [];
            for (var r = 0; r < BUNK.rows; r++) for (var c = 0; c < BUNK.cols; c++) cells.push(bunkerCell(c, r));
            out.push({ x: Math.round(G.W * (i + 1) / (n + 1) - w / 2), cells: cells });
        }
        return out;
    }

    function buildStars(rnd) {
        var stars = [];
        for (var i = 0; i < 70; i++) stars.push({ x: rnd() * G.W, y: rnd() * 600, v: 6 + rnd() * 26, tw: 1 + rnd() * 4 });
        return stars;
    }

    function newInvader(kind, col, row, g) {
        return { kind: kind, col: col, row: row, g: g, hp: KIND[kind].hp, state: 'form', x: 0, y: 0, vx: 0, vy: 0, flash: 0, fade: 0, bombed: false, crunch: false, bounty: true };
    }

    function slotX(s, e) { return s.groups[e.g].x + e.col * CW; }
    function slotY(s, e) { return s.groups[e.g].y + e.row * RH; }

    function buildFormation(s) {
        var cfg = s.cfg, cols = cfg.map[0].length, x0 = (G.W - (cols - 1) * CW) / 2, y0 = cfg.top || 84;
        s.groups = [{ x: x0, y: y0, dir: 1 }];
        if (cfg.split) s.groups.push({ x: x0, y: y0, dir: -1 });
        cfg.map.forEach(function (line, row) {
            for (var col = 0; col < cols; col++) {
                var kind = line.charAt(col);
                if (!KIND[kind]) continue;
                var e = newInvader(kind, col, row, cfg.split && row >= cfg.split ? 1 : 0);
                e.x = slotX(s, e); e.y = slotY(s, e);
                s.inv.push(e);
            }
        });
        s.total = s.inv.length;
    }

    function newBoss(b) {
        return { x: G.W / 2, y: 82, vx: 0, tx: G.W / 2, hp: b.hp, max: b.hp, t: 0, moveT: 0, fanT: 3, miniT: b.minis || 0, lanceT: b.lance || 0, flash: 0, capAt: b.hp - 9 };
    }

    function init(level) {
        var cfg = LV[level - 1], rnd = G.rng(level * 7919 + 13);
        var s = {
            cfg: cfg, rnd: rnd, player: { x: G.W / 2, vx: 0 }, inv: [], groups: [], total: 0, drop: 0,
            shots: [], bombs: [], caps: [], beams: [], bunkers: buildBunkers(cfg.bunkers), stars: buildStars(rnd),
            boss: cfg.boss ? newBoss(cfg.boss) : null, saucer: null, saucerT: 9 + rnd() * 6, warble: 0,
            power: { rapid: 0, double: 0 }, shield: false, invuln: 0, cooldown: 0, wipe: false,
            beat: 0.4, note: 0, frame: 0, alarmT: 0, bombT: 2, diveT: cfg.dive || 0, lanceT: cfg.lance || 0
        };
        buildFormation(s);
        return s;
    }

    // ------------------------------------------------------------------
    // Bunkers
    // ------------------------------------------------------------------

    function cellAt(s, x, y) {
        var r = Math.floor((y - BUNK.y) / BUNK.cell);
        if (r < 0 || r >= BUNK.rows) return false;
        return s.bunkers.some(function (b) {
            var c = Math.floor((x - b.x) / BUNK.cell);
            return c >= 0 && c < BUNK.cols && b.cells[r * BUNK.cols + c] === 1;
        });
    }

    // Clears every cell whose box touches the rectangle, each with the given
    // chance, so hits leave ragged craters. Returns how many cells went.
    function clearCells(s, x0, y0, x1, y1, chance) {
        var n = 0, cs = BUNK.cell;
        var r0 = Math.max(0, Math.floor((y0 - BUNK.y) / cs)), r1 = Math.min(BUNK.rows - 1, Math.floor((y1 - BUNK.y) / cs));
        s.bunkers.forEach(function (b) {
            var c0 = Math.max(0, Math.floor((x0 - b.x) / cs)), c1 = Math.min(BUNK.cols - 1, Math.floor((x1 - b.x) / cs));
            for (var r = r0; r <= r1; r++) {
                for (var c = c0; c <= c1; c++) {
                    var i = r * BUNK.cols + c;
                    if (b.cells[i] && s.rnd() < chance) { b.cells[i] = 0; n++; }
                }
            }
        });
        return n;
    }

    // A shot or bomb at (x, y) bites a crater out of a bunker; false when
    // there was nothing solid at that point.
    function crater(s, x, y) {
        if (!cellAt(s, x, y)) return false;
        clearCells(s, x, y, x, y, 1);
        clearCells(s, x - 7, y - 8, x + 7, y + 8, 0.7);
        G.burst(x, y, { n: 5, color: C.deep, speed: 90, life: 0.35, gravity: 300 });
        G.noise(0.06, { freq: 700, vol: 0.12 });
        return true;
    }

    // ------------------------------------------------------------------
    // Player
    // ------------------------------------------------------------------

    function tickTimers(s, dt) {
        if (s.invuln > 0) s.invuln -= dt;
        if (s.cooldown > 0) s.cooldown -= dt;
        if (s.alarmT > 0) s.alarmT -= dt;
        for (var k in s.power) if (s.power[k] > 0) s.power[k] -= dt;
    }

    // The cannon has a little inertia: it accelerates hard and brakes by
    // friction, so quick taps nudge it and a held key slides it.
    function movePlayer(s, dt) {
        var p = s.player, dir = (G.key.right ? 1 : 0) - (G.key.left ? 1 : 0);
        if (dir) p.vx = G.clamp(p.vx + dir * 2800 * dt, -340, 340);
        else p.vx -= p.vx * Math.min(1, 12 * dt);
        p.x += p.vx * dt;
        if (p.x < 24 || p.x > G.W - 24) { p.x = G.clamp(p.x, 24, G.W - 24); p.vx = 0; }
    }

    // Two volleys may be in the air at once (four with the rapid capsule);
    // the double capsule makes every volley a pair.
    function fire(s) {
        var twin = s.power.double > 0, rapid = s.power.rapid > 0;
        var limit = (rapid ? 4 : 2) * (twin ? 2 : 1), x = s.player.x;
        if (!G.key.a || s.cooldown > 0 || s.shots.length >= limit) return;
        s.cooldown = rapid ? 0.13 : 0.3;
        if (twin) s.shots.push({ x: x - 11, y: PY - 14 }, { x: x + 11, y: PY - 14 });
        else s.shots.push({ x: x, y: PY - 16 });
        G.sfx('shoot');
    }

    function hitPlayer(s) {
        var p = s.player;
        if (s.invuln > 0) return;
        if (s.shield) {
            s.shield = false; s.invuln = 1;
            G.burst(p.x, PY, { n: 16, color: C.alert, speed: 220, life: 0.4 });
            G.tone(700, 0.25, { type: 'triangle', slide: 180, vol: 0.2 });
            return;
        }
        G.burst(p.x, PY, { n: 28, color: C.mint, speed: 260, life: 0.8, gravity: 300 });
        G.noise(0.5, { freq: 1100, slide: 70, vol: 0.35 });
        // A hit costs the capsule powers and wipes the bombs in the air (at
        // the end of the tick, as this may run inside their loop). Divers,
        // minis and charging beams stay: the invulnerability below is what
        // gives the player time to get clear of those.
        s.wipe = true; s.power.rapid = 0; s.power.double = 0;
        s.invuln = 2.2; p.vx = 0;
        G.loseLife();
    }

    function playerBox(s) { return { x: s.player.x - 16, y: PY - 8, w: 32, h: 18 }; }

    // ------------------------------------------------------------------
    // Formation march
    // ------------------------------------------------------------------

    // The fewer invaders are left, the faster the rest march; squaring keeps
    // the rush for the last handful, like the original.
    // Minis are loose and do not count: splitting an invader must not slow
    // the march down again.
    function marchSpeed(s) {
        var left = s.inv.filter(function (e) { return e.kind !== 'm'; }).length;
        var thin = 1 - Math.min(1, left / s.total);
        return s.cfg.lo + (s.cfg.hi - s.cfg.lo) * thin * thin;
    }

    // Marches one group sideways. Returns true when it met a wall and turned.
    function marchGroup(s, g, gi, v, dt) {
        var min = 99, max = -1;
        s.inv.forEach(function (e) {
            if (e.g !== gi) return;
            if (e.col < min) min = e.col;
            if (e.col > max) max = e.col;
        });
        if (max < 0) return false;
        g.x += g.dir * v * dt;
        var hitRight = g.dir > 0 && g.x + max * CW + HALF_W > G.W - MARGIN;
        var hitLeft = g.dir < 0 && g.x + min * CW - HALF_W < MARGIN;
        if (hitRight || hitLeft) g.dir = -g.dir;
        return hitRight || hitLeft;
    }

    // Whichever group meets a wall, every group drops the same step, so
    // counter-marching groups keep their rows aligned and never interleave.
    function marchFormation(s, v, dt) {
        if (s.drop > 0) {
            var d = Math.min(s.drop, 80 * dt);
            s.groups.forEach(function (g) { g.y += d; });
            s.drop -= d;
            return;
        }
        s.groups.forEach(function (g, gi) { if (marchGroup(s, g, gi, v, dt)) s.drop = DROP; });
    }

    // The march: four descending notes whose tempo follows the formation's
    // speed. The same beat flips the sprites' legs.
    function beat(s, v, dt) {
        s.beat -= dt;
        if (s.beat > 0) return;
        s.beat = G.clamp(14 / v, 0.1, 0.7);
        s.frame = 1 - s.frame;
        G.tone(MARCH[s.note % 4], 0.11, { type: 'square', vol: 0.17 });
        G.tone(MARCH[s.note % 4] / 2, 0.14, { type: 'triangle', vol: 0.22 });
        s.note++;
    }

    function wrapFlier(s, e) {
        // Re-enter just below the HUD strip, fading in, rather than flying
        // down through the score line.
        e.y = G.HUD + 12; e.vy = 60; e.bombed = false; e.crunch = false; e.fade = FADE;
        // A diver flies home to its slot; a mini has no slot and dives again.
        if (e.kind === 'm') { e.vx = 0; e.x = 60 + s.rnd() * (G.W - 120); }
        else e.state = 'return';
    }

    function returnToSlot(s, e, dt) {
        var tx = slotX(s, e), ty = slotY(s, e), d = G.dist(e.x, e.y, tx, ty), step = 300 * dt;
        if (d <= step) { e.state = 'form'; e.x = tx; e.y = ty; return; }
        e.x += (tx - e.x) / d * step; e.y += (ty - e.y) / d * step;
    }

    // A flier ploughs through bunker cells just as shots and bombs do, so a
    // bunker never shelters a diver from the cannon underneath it.
    function crunchBunker(s, e) {
        if (e.y + 9 < BUNK.y || !clearCells(s, e.x - 10, e.y - 8, e.x + 10, e.y + 8, 1)) return;
        G.burst(e.x, e.y + 8, { n: 3, color: C.deep, speed: 90, life: 0.3, gravity: 300 });
        if (e.crunch) return;
        e.crunch = true;
        G.noise(0.18, { freq: 500, slide: 150, vol: 0.16 });
    }

    // Divers and minis fall with gravity and steer toward the cannon until
    // they are low; then they are committed and can be sidestepped.
    function updateFlier(s, e, dt) {
        if (e.state === 'return') { returnToSlot(s, e, dt); return; }
        var mini = e.kind === 'm';
        e.vy = Math.min(mini ? 150 : 215, e.vy + 300 * dt);
        if (e.y < 400) e.vx = G.clamp(e.vx + (s.player.x > e.x ? 1 : -1) * 340 * dt, -190, 190);
        e.x = G.clamp(e.x + e.vx * dt, 20, G.W - 20);
        e.y += e.vy * dt;
        if (!e.bombed && e.y > 250 && e.y < 380) { e.bombed = true; dropBomb(s, e.x, e.y + 10, 0, s.cfg.bombV, 'bolt'); }
        if (e.y > G.H + 24) { wrapFlier(s, e); return; }
        crunchBunker(s, e);
        if (G.aabb({ x: e.x - 12, y: e.y - 9, w: 24, h: 18 }, playerBox(s)) && s.invuln <= 0) {
            e.hp = 0;
            G.burst(e.x, e.y, { n: 12, color: KIND[e.kind].color, speed: 200 });
            hitPlayer(s);
        }
    }

    // Moves the groups and every invader. Returns true when the formation
    // has reached the landing line.
    function updateInvaders(s, dt) {
        var v = marchSpeed(s), landed = false;
        marchFormation(s, v, dt);
        if (s.inv.length) beat(s, v, dt);
        s.inv.forEach(function (e) {
            if (e.flash > 0) e.flash -= dt;
            if (e.fade > 0) e.fade -= dt;
            if (e.state !== 'form') { updateFlier(s, e, dt); return; }
            e.x = slotX(s, e); e.y = slotY(s, e);
            if (e.y + HALF_H > BUNK.y) clearCells(s, e.x - HALF_W, e.y - HALF_H, e.x + HALF_W, e.y + HALF_H, 1);
            if (e.y > 380 && s.alarmT <= 0) { s.alarmT = 1.2; G.sfx('alarm'); }
            if (e.y + HALF_H >= LAND_Y) landed = true;
        });
        return landed;
    }

    // ------------------------------------------------------------------
    // Enemy attacks: bombs, divers, lancer beams
    // ------------------------------------------------------------------

    function dropBomb(s, x, y, vx, vy, kind) {
        s.bombs.push({ x: x, x0: x, y: y, vx: vx, vy: vy, kind: kind, t: 0 });
    }

    // Picks who bombs next: usually a random column, sometimes the one above
    // the cannon, and always the lowest invader of that column.
    function pickShooter(s) {
        var form = s.inv.filter(function (e) { return e.state === 'form'; }), px = s.player.x;
        if (!form.length) return null;
        var pick = form[Math.floor(s.rnd() * form.length)];
        if (s.rnd() < 0.35) form.forEach(function (e) { if (Math.abs(e.x - px) < Math.abs(pick.x - px)) pick = e; });
        form.forEach(function (e) { if (Math.abs(e.x - pick.x) < 4 && e.y > pick.y) pick = e; });
        return pick;
    }

    function scheduleBombs(s, dt) {
        s.bombT -= dt;
        if (s.bombT > 0 || s.bombs.length >= s.cfg.maxBombs) return;
        var e = pickShooter(s);
        if (!e) return;
        s.bombT = s.cfg.bomb * (0.6 + s.rnd() * 0.8);
        var zig = s.rnd() < (s.cfg.zig || 0);
        dropBomb(s, e.x, e.y + HALF_H, 0, s.cfg.bombV * (zig ? 0.8 : 1), zig ? 'zig' : 'bolt');
        G.tone(330, 0.12, { type: 'triangle', slide: 120, vol: 0.07 });
    }

    function scheduleDives(s, dt) {
        if (!s.cfg.dive) return;
        s.diveT -= dt;
        if (s.diveT > 0) return;
        var out = s.inv.filter(function (e) { return e.state !== 'form' && e.kind !== 'm'; }).length;
        var form = s.inv.filter(function (e) { return e.state === 'form' && e.kind !== 'l'; });
        if (out >= s.cfg.divers || !form.length) return;
        s.diveT = s.cfg.dive * (0.7 + s.rnd() * 0.6);
        var e = form[Math.floor(s.rnd() * form.length)];
        // It starts by pulling up and away, which reads as peeling off. The
        // climb is cut short for top rows so it stays out of the saucer lane
        // and clear of the mothership's hull.
        var room = Math.max(0, e.y - (s.boss ? 122 : 76));
        e.state = 'dive'; e.bombed = false; e.crunch = false;
        e.vy = -Math.min(140, Math.sqrt(600 * room)); e.vx = (s.rnd() < 0.5 ? -1 : 1) * 120;
        G.tone(950, 0.4, { type: 'sawtooth', slide: 210, vol: 0.09 });
    }

    function startBeam(s, src) {
        s.beams.push({ src: src, t: 0 });
        G.tone(180, CHARGE, { type: 'sawtooth', slide: 1300, vol: 0.07 });
    }

    function scheduleLances(s, dt) {
        if (!s.cfg.lance) return;
        s.lanceT -= dt;
        if (s.lanceT > 0) return;
        var free = s.inv.filter(function (e) {
            return e.kind === 'l' && e.state === 'form' && !s.beams.some(function (b) { return b.src === e; });
        });
        if (!free.length) return;
        s.lanceT = s.cfg.lance * (0.7 + s.rnd() * 0.6);
        startBeam(s, free[Math.floor(s.rnd() * free.length)]);
    }

    // A beam is a thin warning line while it charges, then a column of light
    // that burns through bunkers and the cannon alike.
    function updateBeams(s, dt) {
        s.beams = s.beams.filter(function (b) {
            var src = b.src, before = b.t;
            if (src.hp <= 0 || (src.state && src.state !== 'form')) return false;
            b.t += dt;
            if (b.t < CHARGE) return true;
            if (before < CHARGE) { G.sfx('laser'); G.shake(4, 0.25); }
            clearCells(s, src.x - 8, BUNK.y, src.x + 8, BUNK.y + BUNK.rows * BUNK.cell, 1);
            if (Math.abs(s.player.x - src.x) < 22) hitPlayer(s);
            return b.t < CHARGE + BEAM;
        });
    }

    function updateBombs(s, dt) {
        var box = playerBox(s);
        s.bombs = s.bombs.filter(function (b) {
            b.t += dt; b.y += b.vy * dt; b.x0 += b.vx * dt;
            // Zig bombs weave around their path, so they cannot be dodged by
            // a single small step.
            b.x = b.x0 + (b.kind === 'zig' ? Math.sin(b.t * 8) * 20 : 0);
            if (crater(s, b.x, b.y + 6)) return false;
            if (b.x > box.x - 3 && b.x < box.x + box.w + 3 && b.y + 6 > box.y && b.y - 6 < box.y + box.h && s.invuln <= 0) {
                hitPlayer(s);
                return false;
            }
            if (b.y < GROUND_Y) return b.x > -20 && b.x < G.W + 20;
            G.burst(b.x, GROUND_Y, { n: 4, color: C.alert, speed: 80, life: 0.3, angle: -Math.PI / 2, spread: 2 });
            return false;
        });
    }

    // ------------------------------------------------------------------
    // Kills, capsules, saucer
    // ------------------------------------------------------------------

    function dropCapsule(s, x, y, force) {
        if (!s.cfg.caps || s.caps.length >= 2 || (!force && s.rnd() > 0.09)) return;
        s.caps.push({ x: x, y: y, kind: CAPS[Math.floor(s.rnd() * CAPS.length)] });
    }

    // bounty is false for minis the mothership launches: it launches them
    // for as long as it lives, so paying for them would be a score farm.
    function splitInvader(s, e, bounty) {
        [-1, 1].forEach(function (side) {
            var m = newInvader('m', -1, -1, -1);
            m.bounty = bounty;
            m.state = 'dive'; m.x = e.x + side * 10; m.y = e.y; m.vx = side * 130; m.vy = -90;
            s.inv.push(m);
        });
        G.tone(520, 0.2, { type: 'triangle', slide: 1040, vol: 0.14 });
    }

    function damageInvader(s, e) {
        var k = KIND[e.kind];
        e.hp--;
        if (e.hp > 0) {
            e.flash = 0.15;
            G.burst(e.x, e.y, { n: 8, color: C.alert, speed: 150, life: 0.3 });
            G.tone(1500, 0.08, { type: 'triangle', slide: 700, vol: 0.14 });
            return;
        }
        // Shooting an invader out of its dive is harder, so it pays double.
        var pts = e.bounty ? k.pts * (e.state === 'form' ? 1 : 2) : 0;
        G.addScore(pts);
        if (pts) G.popup(e.x, e.y - 10, pts, k.color);
        G.burst(e.x, e.y, { n: 12, color: k.color, speed: 190, life: 0.5 });
        G.noise(0.14, { freq: 1600, slide: 200, vol: 0.2 });
        G.tone(700 - e.row * 60, 0.09, { slide: 160, vol: 0.1 });
        if (e.kind === 'x') splitInvader(s, e, true);
        if (e.bounty) dropCapsule(s, e.x, e.y, false);
    }

    function applyCapsule(s, kind) {
        if (kind === 'shield') s.shield = true;
        else s.power[kind] = kind === 'rapid' ? 10 : 12;
        G.addScore(50);
        G.popup(s.player.x, PY - 26, kind.toUpperCase(), C.alert);
        G.sfx('power');
    }

    function updateCapsules(s, dt) {
        var px = s.player.x;
        s.caps = s.caps.filter(function (c) {
            c.y += 120 * dt;
            var caught = Math.abs(c.x - px) < 28 && c.y > PY - 16 && c.y < PY + 16;
            if (caught) applyCapsule(s, c.kind);
            return !caught && c.y < GROUND_Y;
        });
    }

    // The bonus saucer crosses the top now and then; it always carries a
    // capsule on the levels that have them.
    function updateSaucer(s, dt) {
        var u = s.saucer;
        if (!u) {
            if (s.boss) return;
            s.saucerT -= dt;
            if (s.saucerT > 0) return;
            var dir = s.rnd() < 0.5 ? 1 : -1;
            s.saucer = { x: dir > 0 ? -30 : G.W + 30, y: 48, dir: dir };
            return;
        }
        u.x += u.dir * 130 * dt;
        s.warble -= dt;
        if (s.warble <= 0) { s.warble = 0.2; G.tone(880, 0.16, { type: 'sine', slide: 1320, vol: 0.06 }); }
        if (u.x < -40 || u.x > G.W + 40) { s.saucer = null; s.saucerT = 14 + s.rnd() * 8; }
    }

    function killSaucer(s) {
        var u = s.saucer, pts = [100, 150, 200, 300][Math.floor(s.rnd() * 4)];
        G.addScore(pts);
        G.popup(u.x, u.y + 22, pts, C.alert);
        G.burst(u.x, u.y, { n: 22, color: C.alert, speed: 220, life: 0.6 });
        G.sfx('boom'); G.sfx('coin');
        dropCapsule(s, u.x, u.y, true);
        s.saucer = null; s.saucerT = 14 + s.rnd() * 8;
    }

    // ------------------------------------------------------------------
    // Mothership
    // ------------------------------------------------------------------

    // The core hatch is open for part of every five seconds; a shot into the
    // open core counts triple.
    function coreOpen(b) { return b.t % 5 > 2.8; }

    function bossFan(s, b) {
        var n = s.cfg.boss.n, v = s.cfg.bombV + 30;
        var aim = Math.atan2(PY - b.y, s.player.x - b.x);
        for (var i = 0; i < n; i++) {
            var a = aim + (i - (n - 1) / 2) * 0.22;
            dropBomb(s, b.x, b.y + 20, Math.cos(a) * v, Math.sin(a) * v, 'bolt');
        }
        G.tone(240, 0.3, { type: 'sawtooth', slide: 90, vol: 0.14 });
    }

    function bossMinis(s, b) {
        var loose = s.inv.filter(function (e) { return e.kind === 'm'; }).length;
        if (loose >= 4) return;
        splitInvader(s, { x: b.x, y: b.y + 20 }, false);
    }

    // The mothership is pulled toward a wandering target by a damped spring,
    // which gives it a heavy, drifting motion.
    function updateBoss(s, dt) {
        var b = s.boss, cfg = s.cfg.boss;
        b.t += dt;
        if (b.flash > 0) b.flash -= dt;
        b.moveT -= dt;
        if (b.moveT <= 0) { b.moveT = 2 + s.rnd() * 1.5; b.tx = 110 + s.rnd() * (G.W - 220); }
        // While its own beam is charging or firing it brakes to a halt, so
        // the warning line marks where the beam will really be.
        if (s.beams.some(function (bm) { return bm.src === b; })) b.vx -= b.vx * Math.min(1, 6 * dt);
        else b.vx += ((b.tx - b.x) * 1.6 - b.vx * 1.8) * dt;
        b.x += b.vx * dt;
        b.fanT -= dt;
        if (b.fanT <= 0) { b.fanT = cfg.fan; bossFan(s, b); }
        if (cfg.minis) {
            b.miniT -= dt;
            if (b.miniT <= 0) { b.miniT = cfg.minis; bossMinis(s, b); }
        }
        // Wounded, it starts using the lancers' beam itself.
        if (cfg.lance && b.hp < b.max / 2) {
            b.lanceT -= dt;
            if (b.lanceT <= 0) { b.lanceT = cfg.lance; startBeam(s, b); }
        }
    }

    function killBoss(s) {
        var b = s.boss;
        G.addScore(200 * G.level);
        G.popup(b.x, b.y, 200 * G.level, C.alert);
        G.burst(b.x, b.y, { n: 70, color: C.alert, speed: 340, life: 1.1 });
        G.burst(b.x, b.y, { n: 40, color: C.white, speed: 200, life: 0.8 });
        G.sfx('bigboom'); G.shake(12, 0.6); G.flash(C.white, 0.25);
        s.bombs = [];
        s.boss = null;
    }

    function damageBoss(s, sh) {
        var b = s.boss, core = coreOpen(b) && Math.abs(sh.x - b.x) < 18;
        b.hp -= core ? 3 : 1;
        b.flash = 0.1;
        G.addScore(core ? 30 : 10);
        G.burst(sh.x, b.y + 14, { n: core ? 14 : 5, color: core ? C.alert : C.mint, speed: 170, life: 0.35 });
        G.tone(core ? 220 : 1300, 0.1, { type: core ? 'sawtooth' : 'triangle', slide: core ? 660 : 500, vol: 0.16 });
        if (b.hp <= 0) { killBoss(s); return; }
        if (b.hp <= b.capAt) { b.capAt -= 9; dropCapsule(s, b.x, b.y + 20, true); }
    }

    // ------------------------------------------------------------------
    // Player shots
    // ------------------------------------------------------------------

    function shotHitsBomb(s, sh) {
        for (var i = 0; i < s.bombs.length; i++) {
            var b = s.bombs[i];
            if (Math.abs(b.x - sh.x) > 7 || Math.abs(b.y - sh.y) > 12) continue;
            s.bombs.splice(i, 1);
            G.addScore(5);
            G.burst(b.x, b.y, { n: 6, color: C.alert, speed: 120, life: 0.25 });
            G.sfx('click');
            return true;
        }
        return false;
    }

    function shotHitsInvader(s, sh) {
        for (var i = 0; i < s.inv.length; i++) {
            var e = s.inv[i], mini = e.kind === 'm';
            if (e.hp <= 0 || Math.abs(e.x - sh.x) > (mini ? 12 : HALF_W + 2) || Math.abs(e.y - sh.y) > (mini ? 12 : HALF_H + 6)) continue;
            damageInvader(s, e);
            return true;
        }
        return false;
    }

    // True when the shot was used up. Checked nearest-first from the cannon's
    // point of view: bunker, bombs, invaders, saucer, mothership.
    function shotHits(s, sh) {
        if (crater(s, sh.x, sh.y) || crater(s, sh.x, sh.y + 7)) return true;
        if (shotHitsBomb(s, sh) || shotHitsInvader(s, sh)) return true;
        var u = s.saucer, b = s.boss;
        if (u && Math.abs(u.x - sh.x) < 24 && Math.abs(u.y - sh.y) < 12) { killSaucer(s); return true; }
        if (b && Math.abs(b.x - sh.x) < 78 && Math.abs(b.y - sh.y) < 22) { damageBoss(s, sh); return true; }
        return false;
    }

    function updateShots(s, dt) {
        s.shots = s.shots.filter(function (sh) {
            sh.y -= SHOT_V * dt;
            return sh.y > G.HUD + 4 && !shotHits(s, sh);
        });
    }

    function update(s, dt) {
        tickTimers(s, dt);
        movePlayer(s, dt);
        fire(s);
        if (updateInvaders(s, dt)) {
            // Touchdown: the level is lost and starts over.
            G.flash(C.alert, 0.3); G.sfx('bigboom');
            G.die();
            return;
        }
        scheduleBombs(s, dt); scheduleDives(s, dt); scheduleLances(s, dt);
        if (s.boss) updateBoss(s, dt);
        updateShots(s, dt); updateBombs(s, dt); updateBeams(s, dt);
        updateCapsules(s, dt); updateSaucer(s, dt);
        s.inv = s.inv.filter(function (e) { return e.hp > 0; });
        if (s.wipe) { s.bombs = []; s.wipe = false; }
        if (!s.inv.length && !s.boss) G.win(500 + G.lives * 250 + G.level * 100);
    }

    // ------------------------------------------------------------------
    // Drawing
    // ------------------------------------------------------------------

    function drawSprite(ctx, spr, cx, cy, color) {
        var ox = Math.round(cx - spr.w * PX / 2), oy = Math.round(cy - spr.h * PX / 2);
        ctx.fillStyle = color;
        for (var i = 0; i < spr.runs.length; i++) {
            var r = spr.runs[i];
            ctx.fillRect(ox + r[0] * PX, oy + r[1] * PX, r[2] * PX, PX);
        }
    }

    function drawBackdrop(s, ctx) {
        var span = GROUND_Y - G.HUD;
        ctx.fillStyle = C.bg;
        ctx.fillRect(0, 0, G.W, G.H);
        s.stars.forEach(function (st, i) {
            ctx.globalAlpha = 0.35 + 0.3 * Math.sin(G.t * st.tw + i);
            ctx.fillStyle = i % 9 ? C.dim : C.mint;
            ctx.fillRect(st.x, G.HUD + (st.y + G.t * st.v) % span, 2, 2);
        });
        // The dashed line is where the formation lands; it crawls so the
        // threat reads even when the player is not looking at the invaders.
        ctx.globalAlpha = 0.6;
        ctx.fillStyle = C.deep;
        for (var x = -24 + (G.t * 20) % 24; x < G.W; x += 24) ctx.fillRect(x, LAND_Y, 12, 2);
        ctx.globalAlpha = 1;
        ctx.fillStyle = C.mint;
        ctx.fillRect(0, GROUND_Y, G.W, 2);
        ctx.fillStyle = C.line;
        ctx.fillRect(0, GROUND_Y + 6, G.W, 1);
    }

    function drawBunkers(s, ctx) {
        var cs = BUNK.cell;
        s.bunkers.forEach(function (b) {
            for (var i = 0; i < b.cells.length; i++) {
                if (!b.cells[i]) continue;
                var c = i % BUNK.cols, r = Math.floor(i / BUNK.cols);
                ctx.fillStyle = r < 2 ? C.soft : C.deep;
                ctx.fillRect(b.x + c * cs, BUNK.y + r * cs, cs, cs);
            }
        });
    }

    // On the cloak level a wave of invisibility sweeps through the rows.
    function invaderAlpha(s, e) {
        if (e.fade > 0) return 1 - e.fade / FADE;
        if (!s.cfg.cloak || e.state !== 'form') return 1;
        return G.clamp(0.55 + Math.sin(G.t * 1.3 + e.row * 0.9 + e.col * 0.25) * 1.1, 0.07, 1);
    }

    function drawInvaders(s, ctx) {
        s.inv.forEach(function (e) {
            var k = KIND[e.kind];
            ctx.globalAlpha = invaderAlpha(s, e);
            drawSprite(ctx, SPR[k.spr][s.frame], e.x, e.y, e.flash > 0 ? C.white : k.color);
            if (e.kind === 's' && e.hp > 1) {
                ctx.strokeStyle = C.alert; ctx.lineWidth = 2;
                ctx.beginPath(); ctx.arc(e.x, e.y + 2, 20, 0.15 * Math.PI, 0.85 * Math.PI); ctx.stroke();
            }
            if (e.kind === 'l' && e.hp > 1) { ctx.fillStyle = C.bg; ctx.fillRect(e.x - 2, e.y - 6, 4, 4); }
        });
        ctx.globalAlpha = 1;
    }

    function drawBeams(s, ctx) {
        s.beams.forEach(function (b) {
            var x = b.src.x, y = b.src.y + 14, live = b.t >= CHARGE;
            ctx.fillStyle = C.alert;
            if (!live) {
                // Charging: a flickering hairline that thickens as it fills.
                ctx.globalAlpha = 0.25 + 0.5 * (b.t / CHARGE) * (Math.sin(G.t * 40) > 0 ? 1 : 0.5);
                ctx.fillRect(x - 1, y, 2, GROUND_Y - y);
            } else {
                ctx.globalAlpha = 0.35; ctx.fillRect(x - 16, y, 32, GROUND_Y - y);
                ctx.globalAlpha = 1; ctx.fillRect(x - 7, y, 14, GROUND_Y - y);
                ctx.fillStyle = C.white; ctx.fillRect(x - 2, y, 4, GROUND_Y - y);
            }
        });
        ctx.globalAlpha = 1;
    }

    function drawBoss(s, ctx) {
        var b = s.boss, open = coreOpen(b), hit = b.flash > 0;
        ctx.fillStyle = hit ? C.white : C.line;
        ctx.strokeStyle = C.mint; ctx.lineWidth = 2;
        ctx.beginPath(); ctx.ellipse(b.x, b.y - 6, 42, 18, 0, Math.PI, 0); ctx.fill(); ctx.stroke();
        ctx.beginPath(); ctx.ellipse(b.x, b.y + 4, 78, 17, 0, 0, Math.PI * 2); ctx.fill(); ctx.stroke();
        for (var i = -3; i <= 3; i++) {
            // Running lights chase along the hull.
            ctx.fillStyle = (Math.floor(b.t * 6) + i + 30) % 4 === 0 ? C.white : C.deep;
            if (i) ctx.fillRect(b.x + i * 20 - 3, b.y + 1, 6, 5);
        }
        ctx.fillStyle = open ? C.alert : C.bg;
        ctx.globalAlpha = open ? 0.7 + 0.3 * Math.sin(b.t * 14) : 1;
        ctx.beginPath(); ctx.arc(b.x, b.y + 6, open ? 12 : 8, 0, Math.PI * 2); ctx.fill();
        ctx.globalAlpha = 1;
        ctx.strokeStyle = open ? C.alert : C.mint;
        ctx.stroke();
        // Hull strength, just under the HUD.
        ctx.fillStyle = C.line; ctx.fillRect(G.W / 2 - 150, 36, 300, 5);
        ctx.fillStyle = C.alert; ctx.fillRect(G.W / 2 - 150, 36, 300 * Math.max(0, b.hp) / b.max, 5);
    }

    function drawSaucer(s, ctx) {
        var u = s.saucer;
        ctx.fillStyle = C.alert;
        ctx.beginPath(); ctx.ellipse(u.x, u.y + 2, 22, 7, 0, 0, Math.PI * 2); ctx.fill();
        ctx.beginPath(); ctx.ellipse(u.x, u.y - 3, 10, 7, 0, Math.PI, 0); ctx.fill();
        ctx.fillStyle = C.bg;
        for (var i = -1; i <= 1; i++) if ((Math.floor(G.t * 8) + i + 9) % 3) ctx.fillRect(u.x + i * 11 - 2, u.y, 4, 3);
    }

    function drawPlayer(s, ctx) {
        var x = Math.round(s.player.x), twin = s.power.double > 0;
        // Blink while invulnerable after a hit.
        if (s.invuln > 0 && Math.floor(s.invuln * 12) % 2) return;
        ctx.fillStyle = s.power.rapid > 0 ? C.white : C.mint;
        ctx.fillRect(x - 18, PY - 2, 36, 12);
        ctx.fillRect(x - 12, PY - 7, 24, 5);
        if (twin) { ctx.fillRect(x - 13, PY - 14, 5, 8); ctx.fillRect(x + 8, PY - 14, 5, 8); }
        else ctx.fillRect(x - 3, PY - 15, 6, 9);
        ctx.fillStyle = C.bg;
        ctx.fillRect(x - 12, PY + 3, 6, 3); ctx.fillRect(x + 6, PY + 3, 6, 3);
        if (!s.shield) return;
        ctx.strokeStyle = C.alert; ctx.lineWidth = 2;
        ctx.globalAlpha = 0.6 + 0.4 * Math.sin(G.t * 9);
        ctx.beginPath(); ctx.arc(x, PY + 8, 30, Math.PI * 1.08, Math.PI * 1.92); ctx.stroke();
        ctx.globalAlpha = 1;
    }

    function drawProjectiles(s, ctx) {
        ctx.fillStyle = C.white;
        s.shots.forEach(function (sh) { ctx.fillRect(sh.x - 1.5, sh.y, 3, 12); });
        ctx.fillStyle = C.alert;
        s.bombs.forEach(function (b) {
            // Bolts are a straight bar; zig bombs are drawn as a wriggle.
            var k = b.kind === 'zig' ? (Math.floor(b.t * 14) % 2 ? 3 : -3) : 0;
            ctx.fillRect(b.x - 2 + k, b.y - 6, 4, 6);
            ctx.fillRect(b.x - 2 - k, b.y, 4, 6);
        });
        s.caps.forEach(function (c) {
            ctx.fillStyle = C.alert;
            ctx.fillRect(c.x - 13, c.y - 8, 26, 16);
            ctx.fillStyle = C.bg;
            ctx.fillRect(c.x - 11, c.y - 6, 22, 12);
            G.text(c.kind.charAt(0).toUpperCase(), c.x, c.y + 5, { size: 13, bold: true, color: C.alert, align: 'center' });
        });
    }

    function draw(s, ctx) {
        drawBackdrop(s, ctx);
        drawBunkers(s, ctx);
        drawBeams(s, ctx);
        drawInvaders(s, ctx);
        if (s.boss) drawBoss(s, ctx);
        if (s.saucer) drawSaucer(s, ctx);
        drawProjectiles(s, ctx);
        drawPlayer(s, ctx);
    }

    // HOSTILES is everything that still has to be destroyed to clear the
    // level, so it includes loose minis and the mothership (and goes up by
    // one when a splitter bursts into two).
    function hud(s) {
        var on = [];
        if (s.power.rapid > 0) on.push('RAPID');
        if (s.power.double > 0) on.push('DOUBLE');
        if (s.shield) on.push('SHIELD');
        return (on.length ? on.join('+') + '  ' : '') + 'HOSTILES ' + (s.inv.length + (s.boss ? 1 : 0));
    }

    G.register('invaders', {
        title: 'FORMATION',
        blurb: 'Destroy the formation before it lands. It marches faster as it thins.',
        controls: [
            '← → move the cannon · SPACE fire (hold)',
            'Shots and bombs chew through the bunkers; shots can pick bombs out of the air',
            'Capsules: R rapid fire · D double barrel · S shield',
            'Mothership: aim for the core while it glows'
        ],
        levelNames: ['First Contact', 'Arrowhead', 'Peel Off', 'Phalanx', 'Mothership', 'Mitosis', 'Crossfire', 'Cloak', 'Lancers', 'Armada'],
        colors: { bg: C.bg, fg: C.white, accent: C.mint, dim: C.dim },
        lives: 3,
        // Sparse on purpose: the march played by the formation is the real
        // bass line, the sequencer only adds an eerie minor lead over it.
        music: {
            bpm: 104, root: 38, scale: 'minor', prog: [0, 5, 3, 4],
            bass: 'x.......x.....o.',
            lead: [
                '7--.9--.a---..9.', 'c--.a--.9---....', '7--.9--.a--.c--.', 'b---9---8---....',
                'e.e.d.c.b---a...', 'c.c.b.a.9---7...', '7.9.a.b.c-b-a-9.', 'b-------........'
            ],
            arp: '0.2.1.3.',
            drums: { k: 'x.....x.x.......', s: '....x.......x...', h: '..x...x...x...x.' },
            leadWave: 'square', bassWave: 'triangle', arpWave: 'triangle', leadOct: 2
        },
        init: init, update: update, draw: draw, hud: hud
    });
})();
