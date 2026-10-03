/*
 * Wrecking Ball — the brutalist theme's game.
 *
 * A crane drives along the ground with a ball on a chain. The ball is a real
 * pendulum: it only moves because the crane's pivot drags it, so the player
 * builds momentum by driving, braking hard (the ball whips forward) and
 * reeling the chain. Concrete structures stand on a grid with gravity: a
 * block stays up while something holds it from below, or while its row is
 * anchored on both sides (a lintel). Knock out the support and everything
 * above comes down and breaks. Demolish the target share before time is up,
 * and keep the cab out from under what falls.
 */
(function () {
    'use strict';
    var G = window.SnoGame;
    var GROUND = 470, CW = 40, CH = 30, COLS = 24, ROWS = 13, PIVOT_Y = 62, BALL_R = 15;
    var GRAV = 900, L_MIN = 70, L_MAX = GROUND - BALL_R - PIVOT_Y - 2, REEL = 170, SUB = 4;
    var CAB_ACC = 620, CAB_MAX = 250, CAB_HALF = 30, CAB_H = 36, BLAST = 78;
    var CONC = 0, REBAR = 1, GLASS = 2, CHARGE = 3;
    var RED = '#ff2200';

    // Per block kind: hit points, the slowest ball impact that does anything
    // (minV), the damage every such impact adds on top (base), how much a
    // fall hurts (fall) and the score.
    var KINDS = [
        { hp: 1, minV: 80, base: 0, fall: 1, score: 10 },
        { hp: 1.5, minV: 320, base: 0.5, fall: 0, score: 40 },
        { hp: 0.2, minV: 20, base: 1, fall: 10, score: 60 },
        { hp: 1, minV: 60, base: 0, fall: 1, score: 25 }
    ];
    var LEVELS = [
        { time: 75, target: 0.7 }, { time: 90, target: 0.75 }, { time: 100, target: 0.75 },
        { time: 110, target: 0.75 }, { time: 120, target: 0.75 }, { time: 120, target: 0.7 },
        { time: 130, target: 0.65 }, { time: 140, target: 0.7, wind: true },
        { time: 150, target: 0.7, tremor: true }, { time: 180, target: 0.75, wind: true, tremor: true }
    ];

    // ------------------------------------------------------------------
    // Structures: each builder writes block kinds into a ROWS x COLS grid
    // (row 0 stands on the ground). Everything must be held up at the start.
    // ------------------------------------------------------------------

    function fill(g, c0, r0, w, h, kind) {
        for (var r = r0; r < r0 + h; r++) for (var c = c0; c < c0 + w; c++) g[r][c] = kind;
    }

    function buildWall(g) { fill(g, 9, 0, 6, 4, CONC); fill(g, 17, 0, 2, 2, CONC); }

    function buildStacks(g) { fill(g, 6, 0, 2, 8, CONC); fill(g, 10, 0, 2, 3, CONC); fill(g, 15, 0, 3, 9, CONC); }

    // Concrete frame with glass infill on every other floor.
    function buildGlassHouse(g) {
        for (var r = 0; r < 8; r++) {
            for (var c = 7; c < 17; c++) g[r][c] = (r % 2 === 1 && (c - 7) % 3 !== 0) ? GLASS : CONC;
        }
    }

    function buildRebarCore(g) {
        fill(g, 9, 0, 7, 9, CONC); fill(g, 12, 0, 1, 9, REBAR);
        fill(g, 5, 0, 2, 5, CONC); fill(g, 18, 0, 2, 6, CONC);
    }

    // Pillars with doorways between them: the slab above only holds while
    // both of its ends are still standing.
    function buildArcade(g) {
        for (var c = 5; c <= 19; c += 2) fill(g, c, 0, 1, 2, (c === 5 || c === 19) ? REBAR : CONC);
        fill(g, 5, 2, 15, 1, CONC);
        for (c = 5; c <= 19; c++) fill(g, c, 3, 1, 2, (c - 5) % 2 === 1 ? GLASS : CONC);
        fill(g, 5, 5, 15, 1, CONC); fill(g, 9, 6, 7, 2, CONC); fill(g, 12, 6, 1, 2, REBAR);
    }

    // Rebar footings make the usual trick (take out the ground floor)
    // useless; the charges on the faces are the way in.
    function buildShortFuse(g) {
        fill(g, 5, 0, 4, 8, CONC); fill(g, 5, 0, 4, 1, REBAR);
        g[1][5] = CHARGE; g[5][8] = CHARGE;
        fill(g, 13, 0, 6, 10, CONC); fill(g, 13, 0, 6, 1, REBAR);
        fill(g, 13, 1, 1, 3, REBAR); fill(g, 18, 1, 1, 3, REBAR);
        g[4][13] = CHARGE; g[6][18] = CHARGE; g[9][15] = CHARGE;
    }

    function buildSilos(g) {
        [[5, 12], [11, 11], [17, 12]].forEach(function (t) {
            fill(g, t[0], 0, 2, t[1], CONC);
            for (var r = 0; r < t[1]; r += 4) fill(g, t[0], r, 2, 1, REBAR);
            fill(g, t[0], t[1] - 1, 2, 1, GLASS);
        });
        g[2][11] = CHARGE; g[6][18] = CHARGE; g[9][5] = CHARGE;
    }

    function buildZiggurat(g) {
        for (var r = 0; r < 8; r++) {
            for (var c = 4 + r; c <= 19 - r; c++) {
                var edge = c === 4 + r || c === 19 - r;
                g[r][c] = (c === 11 || c === 12) ? REBAR : (r % 2 === 1 && edge ? GLASS : CONC);
            }
        }
        g[2][6] = CHARGE; g[4][15] = CHARGE;
    }

    // A slab block riddled with window holes, on a rebar skeleton.
    function buildCondemned(g) {
        fill(g, 4, 0, 16, 10, CONC);
        fill(g, 4, 0, 1, 5, REBAR); fill(g, 19, 0, 1, 5, REBAR); fill(g, 11, 0, 2, 10, REBAR);
        for (var r = 1; r < 10; r += 2) {
            for (var c = 5; c < 19; c += 3) g[r][c] = (r === 3 || r === 7) ? GLASS : null;
        }
        g[2][4] = CHARGE; g[6][19] = CHARGE; g[9][9] = CHARGE;
    }

    // Podium on pillars, two towers with rebar cores and a glazed skybridge
    // that hangs between them.
    function buildMega(g) {
        for (var c = 5; c <= 18; c++) {
            if ((c - 5) % 2 === 0 || c === 18) fill(g, c, 0, 1, 2, (c === 5 || c === 18 || c === 11) ? REBAR : CONC);
        }
        fill(g, 5, 2, 14, 1, CONC);
        fill(g, 5, 3, 4, 9, CONC); fill(g, 7, 3, 1, 9, REBAR);
        fill(g, 15, 3, 4, 8, CONC); fill(g, 16, 3, 1, 8, REBAR);
        fill(g, 10, 3, 4, 2, CONC); fill(g, 9, 7, 6, 1, CONC); fill(g, 9, 8, 6, 1, GLASS);
        for (var r = 4; r < 11; r += 2) { g[r][6] = GLASS; g[r][17] = GLASS; }
        g[5][5] = CHARGE; g[6][18] = CHARGE; g[3][11] = CHARGE; g[11][8] = CHARGE;
    }

    var BUILD = [buildWall, buildStacks, buildGlassHouse, buildRebarCore, buildArcade,
        buildShortFuse, buildSilos, buildZiggurat, buildCondemned, buildMega];

    function cellY(r) { return GROUND - (r + 1) * CH; }

    function buildGrid(level, rnd) {
        var kinds = [], grid = [], total = 0;
        for (var r = 0; r < ROWS; r++) { kinds.push(new Array(COLS).fill(null)); grid.push(new Array(COLS).fill(null)); }
        BUILD[level - 1](kinds);
        for (r = 0; r < ROWS; r++) {
            for (var c = 0; c < COLS; c++) {
                var k = kinds[r][c];
                if (k === null) continue;
                grid[r][c] = { c: c, r: r, kind: k, hp: KINDS[k].hp, y: cellY(r), vy: 0, loose: 0, fuse: 0, falling: false, dead: false, shade: rnd() };
                total++;
            }
        }
        return { grid: grid, total: total };
    }

    function buildSky(rnd) {
        var sky = [];
        for (var x = -20; x < G.W; x += 50 + rnd() * 70) {
            sky.push({ x: x, w: 50 + rnd() * 80, h: 60 + rnd() * 230, tone: rnd() < 0.5 ? '#151515' : '#1c1c1c' });
        }
        return sky;
    }

    function buildMotes(rnd) {
        var motes = [];
        for (var i = 0; i < 28; i++) motes.push({ x: rnd() * G.W, y: 40 + rnd() * 420, v: 8 + rnd() * 26, size: 2 + rnd() * 3 });
        return motes;
    }

    function init(level) {
        var rnd = G.rng(level * 7919 + 13), built = buildGrid(level, rnd), cfg = LEVELS[level - 1];
        return {
            cfg: cfg, grid: built.grid, total: built.total, gone: 0, need: Math.ceil(built.total * cfg.target),
            time: cfg.time, done: 0, dirty: false, combo: 0, comboT: 0,
            cab: { x: 70, vx: 0, inv: 0 }, L: 300, Ldot: 0, ball: { x: 70, y: PIVOT_Y + 300, vx: 0, vy: 0 },
            falling: [], loose: [], lit: [], chunks: [],
            wind: { f: 0, t: 0, next: 4 }, tremor: { next: 8, warn: false },
            cool: { tink: 0, thud: 0, engine: 0, reel: 0 },
            sky: buildSky(rnd), motes: buildMotes(rnd)
        };
    }

    // ------------------------------------------------------------------
    // Sounds
    // ------------------------------------------------------------------

    // Landings and impacts come in bursts; the gate keeps a collapse from
    // turning into one long wall of noise.
    function thud(s, power) {
        if (s.cool.thud > 0) return;
        s.cool.thud = 0.05;
        G.noise(0.16 + power * 0.2, { freq: 520, slide: 60, vol: Math.min(0.5, 0.14 + power * 0.3) });
        G.tone(96, 0.2, { type: 'sine', slide: 38, vol: Math.min(0.45, 0.15 + power * 0.3) });
    }

    function clang(vol) {
        G.tone(1240, 0.22, { type: 'square', vol: vol * 0.6 });
        G.tone(1873, 0.3, { type: 'triangle', vol: vol });
    }

    function shatter() {
        G.noise(0.25, { filter: 'highpass', freq: 5200, vol: 0.2 });
        G.tone(2600, 0.14, { type: 'triangle', slide: 3900, vol: 0.08 });
    }

    // ------------------------------------------------------------------
    // Blocks: damage, destruction, explosions
    // ------------------------------------------------------------------

    function blockX(b) { return b.c * CW + CW / 2; }

    function eachBlock(s, fn) {
        for (var r = 0; r < ROWS; r++) for (var c = 0; c < COLS; c++) if (s.grid[r][c]) fn(s.grid[r][c]);
        s.falling.forEach(fn);
    }

    // Lists (falling, loose, lit) drop dead blocks on their next pass, so a
    // block can be destroyed from anywhere without touching those arrays.
    function removeBlock(s, b) {
        b.dead = true;
        if (b.falling) return;
        s.grid[b.r][b.c] = null;
        s.dirty = true;
    }

    function rubble(b) {
        var x = blockX(b), y = b.y + CH / 2;
        if (b.kind === GLASS) {
            shatter();
            G.burst(x, y, { n: 12, color: '#ffffff', speed: 240, gravity: 700, size: 3, life: 0.6 });
            return;
        }
        G.burst(x, y, { n: 9, color: b.kind === REBAR ? '#5a5a5a' : '#a0a0a0', speed: 190, gravity: 600, size: 5, life: 0.7 });
        G.burst(x, y, { n: 5, color: b.kind === REBAR ? RED : '#606060', speed: 80, life: 0.9, size: 7, drag: 2 });
        if (b.kind === REBAR) clang(0.1);
    }

    // Blocks that go down in quick succession multiply their score, which is
    // what makes a clean collapse worth more than chipping away.
    function destroy(s, b) {
        if (b.dead) return;
        removeBlock(s, b);
        s.gone++;
        s.combo = s.comboT > 0 ? s.combo + 1 : 1;
        s.comboT = 1.2;
        G.addScore(KINDS[b.kind].score * Math.min(s.combo, 8));
        if (s.combo % 5 === 0) G.popup(blockX(b), b.y, 'x' + s.combo, RED);
        else if (b.kind === GLASS) G.popup(blockX(b), b.y, '+' + KINDS[GLASS].score, '#ffffff');
        if (b.kind === CHARGE) explode(s, b); else rubble(b);
    }

    function lightFuse(s, b, t) {
        if (b.fuse > 0) { b.fuse = Math.min(b.fuse, t); return; }
        b.fuse = t;
        s.lit.push(b);
        G.sfx('alarm');
        G.popup(blockX(b), b.y - 6, 'FUSE LIT', RED);
    }

    // Returns true when the block was destroyed. A charge never breaks: any
    // damage lights its fuse instead.
    function damage(s, b, amount) {
        if (b.dead || amount <= 0) return false;
        if (b.kind === CHARGE) { lightFuse(s, b, 1.6); return false; }
        b.hp -= amount;
        if (b.hp > 0) return false;
        destroy(s, b);
        return true;
    }

    function flightTime(y, vy) {
        return (-vy + Math.sqrt(vy * vy + 2 * GRAV * Math.max(0, GROUND - y))) / GRAV;
    }

    // Debris that can hit the cab. tx is where it will land, so the ground
    // can show a warning mark while it is in the air.
    function addChunk(s, x, y, vx, vy) {
        s.chunks.push({ x: x, y: y, vx: vx, vy: vy, tx: x + vx * flightTime(y, vy), rot: 0, spin: G.rnd(-9, 9), size: 22 });
    }

    function blastBall(s, x, y) {
        var b = s.ball, d = G.dist(x, y, b.x, b.y);
        if (d > 170 || d < 1) return;
        var kick = 460 * (1 - d / 170);
        b.vx += (b.x - x) / d * kick; b.vy += (b.y - y) / d * kick;
    }

    function explode(s, b) {
        var x = blockX(b), y = b.y + CH / 2;
        G.sfx('bigboom'); G.shake(14, 0.5); G.flash('#ffffff', 0.12);
        G.burst(x, y, { n: 26, color: RED, speed: 420, life: 0.7, size: 5, gravity: 300 });
        G.burst(x, y, { n: 14, color: '#ffffff', speed: 260, life: 0.5, size: 4 });
        eachBlock(s, function (o) {
            if (o.dead || G.dist(x, y, blockX(o), o.y + CH / 2) > BLAST) return;
            // A neighbouring charge goes off a moment later: a chain reaction.
            if (o.kind === CHARGE) lightFuse(s, o, 0.2); else damage(s, o, 2);
        });
        blastBall(s, x, y);
        if (Math.abs(s.cab.x - x) < 85 && y > GROUND - 150) hurtCab(s);
        for (var i = 0; i < 3; i++) addChunk(s, x, y, G.rnd(-260, 260), G.rnd(-520, -260));
    }

    // ------------------------------------------------------------------
    // Structural gravity
    // ------------------------------------------------------------------

    function firm(b) { return !!b && b.loose <= 0; }

    // Higher blocks let go a little later, so a column comes down as a
    // ripple and every block lands (and breaks) on its own.
    function loosen(s, b) {
        b.loose = 0.1 + 0.03 * b.r;
        s.loose.push(b);
    }

    // Within a run of touching blocks, everything between the first and the
    // last block that stands on something is held (a beam on two supports).
    // Whatever sticks out past them is a cantilever and lets go.
    function settleRow(s, r, below) {
        var row = s.grid[r], held = [], c = 0;
        while (c < COLS) {
            if (!firm(row[c])) { c++; continue; }
            var start = c, first = -1, last = -1;
            for (; c < COLS && firm(row[c]); c++) {
                if (below && !below[c]) continue;
                if (first < 0) first = c;
                last = c;
            }
            for (var i = start; i < c; i++) {
                if (first >= 0 && i >= first && i <= last) held[i] = true; else loosen(s, row[i]);
            }
        }
        return held;
    }

    function settle(s) {
        var below = null;
        s.dirty = false;
        for (var r = 0; r < ROWS; r++) below = settleRow(s, r, below);
    }

    function updateLoose(s, dt) {
        s.loose = s.loose.filter(function (b) {
            if (b.dead) return false;
            b.loose -= dt;
            if (b.loose > 0) return true;
            s.grid[b.r][b.c] = null;
            b.loose = 0; b.falling = true; b.vy = 0;
            s.falling.push(b);
            return false;
        });
    }

    // The row a falling block will come to rest in: one above the highest
    // block that is still below its feet.
    function landingRow(s, b) {
        var feet = b.y + CH;
        for (var r = ROWS - 1; r >= 0; r--) {
            if (s.grid[r][b.c] && cellY(r) >= feet - CH / 2) return r + 1;
        }
        return 0;
    }

    function land(s, b, row) {
        var v = b.vy, dmg = v / 520, under = row > 0 ? s.grid[row - 1][b.c] : null;
        b.falling = false; b.r = row; b.y = cellY(row); b.vy = 0;
        s.grid[row][b.c] = b;
        s.dirty = true;
        if (v < 60) return;
        thud(s, v / 600);
        G.burst(blockX(b), b.y + CH, { n: 5, color: '#707070', speed: 110, angle: -Math.PI / 2, spread: Math.PI, life: 0.5, size: 4 });
        // The block that is landed on takes part of the blow as well.
        if (under) damage(s, under, dmg * 0.4 * KINDS[under.kind].fall);
        damage(s, b, dmg * KINDS[b.kind].fall);
    }

    function hitsCab(s, x, y, w, h) {
        return x < s.cab.x + CAB_HALF && x + w > s.cab.x - CAB_HALF && y + h > GROUND - CAB_H && y < GROUND;
    }

    // After a hit the cab is safe for two seconds, and it cannot be hurt
    // once the target is reached, so the final collapse can be enjoyed.
    function hurtCab(s) {
        if (s.cab.inv > 0 || s.done > 0) return;
        s.cab.inv = 2;
        G.flash(RED, 0.2);
        G.burst(s.cab.x, GROUND - 24, { n: 18, color: '#ffffff', speed: 240, gravity: 500 });
        G.loseLife();
    }

    function updateFalling(s, dt) {
        s.falling = s.falling.filter(function (b) {
            if (b.dead) return false;
            b.vy += GRAV * dt; b.y += b.vy * dt;
            if (hitsCab(s, b.c * CW, b.y, CW, CH)) { destroy(s, b); hurtCab(s); return false; }
            var row = landingRow(s, b);
            if (b.y < cellY(row)) return true;
            land(s, b, row);
            return false;
        });
    }

    // An explosion can light further fuses, so the list is swapped out
    // before it is walked and survivors are pushed onto the new one.
    function updateFuses(s, dt) {
        var lit = s.lit;
        s.lit = [];
        lit.forEach(function (b) {
            if (b.dead) return;
            var before = b.fuse;
            b.fuse -= dt;
            if (Math.floor(before * 6) !== Math.floor(b.fuse * 6)) G.tone(1500, 0.03, { vol: 0.08 });
            if (b.fuse > 0) s.lit.push(b); else destroy(s, b);
        });
    }

    function updateChunks(s, dt) {
        s.chunks = s.chunks.filter(function (k) {
            k.vy += GRAV * dt; k.x += k.vx * dt; k.y += k.vy * dt; k.rot += k.spin * dt;
            var hit = hitsCab(s, k.x - k.size / 2, k.y - k.size / 2, k.size, k.size);
            if (!hit && k.y < GROUND - k.size / 2) return true;
            if (hit) hurtCab(s);
            thud(s, 0.3);
            G.burst(k.x, GROUND - 6, { n: 8, color: '#8a8a8a', speed: 160, angle: -Math.PI / 2, spread: Math.PI, gravity: 500, size: 4 });
            return false;
        });
    }

    // ------------------------------------------------------------------
    // Crane and ball
    // ------------------------------------------------------------------

    function approach(v, target, step) {
        return v < target ? Math.min(target, v + step) : Math.max(target, v - step);
    }

    // SPACE is a hard brake: the cab stops dead and the ball keeps going,
    // which is the quickest way to throw it at a wall.
    function driveCab(s, dt) {
        var cab = s.cab, dir = (G.key.right ? 1 : 0) - (G.key.left ? 1 : 0);
        if (G.key.a) {
            if (G.hit.a && Math.abs(cab.vx) > 90) G.noise(0.22, { filter: 'bandpass', freq: 2400, slide: 900, q: 6, vol: 0.12 });
            cab.vx = approach(cab.vx, 0, 1800 * dt);
        } else if (dir) {
            cab.vx = G.clamp(cab.vx + dir * CAB_ACC * dt, -CAB_MAX, CAB_MAX);
        } else {
            cab.vx = approach(cab.vx, 0, 420 * dt);
        }
        if (Math.abs(cab.vx) > 20 && s.cool.engine <= 0) {
            s.cool.engine = 0.13;
            G.noise(0.14, { freq: 140 + Math.abs(cab.vx) * 0.8, vol: 0.07 });
        }
    }

    function reelChain(s, dt) {
        var dir = (G.key.down ? 1 : 0) - (G.key.up ? 1 : 0), old = s.L;
        s.L = G.clamp(old + dir * REEL * dt, L_MIN, L_MAX);
        s.Ldot = (s.L - old) / dt;
        if (s.Ldot && s.cool.reel <= 0) {
            s.cool.reel = 0.07;
            G.tone(dir > 0 ? 420 : 560, 0.02, { vol: 0.05 });
        }
    }

    function moveCab(s, h) {
        var cab = s.cab;
        cab.x += cab.vx * h;
        if (cab.x < 40) { cab.x = 40; cab.vx = Math.max(0, cab.vx); }
        if (cab.x > G.W - 40) { cab.x = G.W - 40; cab.vx = Math.min(0, cab.vx); }
    }

    // The chain is a rope, not a rod: it only acts when taut. It then takes
    // away the ball's velocity away from the pivot (relative to the moving
    // crane), which is what drags the ball along and makes it swing.
    function rope(s, h) {
        var b = s.ball, cab = s.cab, dx = b.x - cab.x, dy = b.y - PIVOT_Y, d = Math.hypot(dx, dy);
        if (d <= s.L || d < 1) return;
        var nx = dx / d, ny = dy / d, rvx = b.vx - cab.vx, rvy = b.vy;
        // Pull back at a limited rate so a snagged ball is not teleported
        // through the wall that holds it.
        var pull = Math.min(d - s.L, 5);
        b.x -= nx * pull; b.y -= ny * pull;
        var vn = rvx * nx + rvy * ny;
        if (vn > s.Ldot) { rvx -= nx * (vn - s.Ldot); rvy -= ny * (vn - s.Ldot); vn = s.Ldot; }
        // Reeling in a swinging ball speeds it up (angular momentum).
        var k = 1 - s.Ldot * h / s.L;
        b.vx = cab.vx + nx * vn + (rvx - nx * vn) * k;
        b.vy = ny * vn + (rvy - ny * vn) * k;
    }

    // A ball caught on a structure holds the crane back: the chain can
    // stretch 30 px and no further. If the ball hangs too low for that, the
    // winch slips instead.
    function snag(s) {
        var b = s.ball, cab = s.cab, dx = b.x - cab.x, dy = b.y - PIVOT_Y, max = s.L + 30;
        if (dx * dx + dy * dy <= max * max) return;
        if (Math.abs(dy) >= max) { s.L = G.clamp(Math.hypot(dx, dy) - 30, L_MIN, L_MAX); return; }
        var reach = Math.sqrt(max * max - dy * dy);
        cab.x = G.clamp(b.x - (dx > 0 ? reach : -reach), 40, G.W - 40);
        if (cab.vx * dx < 0) cab.vx = 0;
    }

    // What the ball does to a block it hits at speed v. Returns true when
    // the block broke.
    function strike(s, blk, v) {
        var K = KINDS[blk.kind];
        if (v < K.minV) {
            if (v > 40 && s.cool.tink <= 0) {
                s.cool.tink = 0.25;
                if (blk.kind === REBAR && v > 150) { clang(0.07); G.popup(blockX(blk), blk.y, 'HIT HARDER', '#c8c8c8'); }
                else G.tone(620, 0.04, { type: 'triangle', vol: 0.07 });
            }
            return false;
        }
        thud(s, v / 500);
        G.shake(Math.min(7, v / 90), 0.12);
        if (blk.kind === REBAR) {
            clang(0.12);
            G.burst(s.ball.x, s.ball.y, { n: 8, color: RED, speed: 260, life: 0.3, size: 3 });
        }
        return damage(s, blk, (v - K.minV) / 140 + K.base);
    }

    // Unit vector from the block's surface to the ball's centre, and how far
    // apart they are. A ball whose centre is inside the block is pushed out
    // through the nearer pair of faces.
    function contact(b, rx, ry) {
        var dx = b.x - G.clamp(b.x, rx, rx + CW), dy = b.y - G.clamp(b.y, ry, ry + CH), d = Math.hypot(dx, dy);
        if (d > 0.001) return { x: dx / d, y: dy / d, d: d };
        var ox = b.x - (rx + CW / 2), oy = b.y - (ry + CH / 2);
        if (Math.abs(ox) / CW > Math.abs(oy) / CH) return { x: ox < 0 ? -1 : 1, y: 0, d: 0 };
        return { x: 0, y: oy < 0 ? -1 : 1, d: 0 };
    }

    function hitBlock(s, blk, rx, ry) {
        var b = s.ball;
        if (blk.dead || !G.circRect(b.x, b.y, BALL_R, rx, ry, CW, CH)) return;
        var n = contact(b, rx, ry), vn = b.vx * n.x + b.vy * n.y;
        b.x += n.x * (BALL_R - n.d); b.y += n.y * (BALL_R - n.d);
        if (vn >= 0) return;
        // A block that breaks only slows the ball, so one good swing can
        // plough through several; a block that holds bounces it back.
        if (strike(s, blk, -vn)) { b.vx *= 0.8; b.vy *= 0.8; return; }
        b.vx -= 1.35 * vn * n.x; b.vy -= 1.35 * vn * n.y;
    }

    function collideBlocks(s) {
        var b = s.ball;
        var c0 = Math.max(0, Math.floor((b.x - BALL_R) / CW)), c1 = Math.min(COLS - 1, Math.floor((b.x + BALL_R) / CW));
        var r0 = Math.max(0, Math.floor((GROUND - b.y - BALL_R) / CH)), r1 = Math.min(ROWS - 1, Math.floor((GROUND - b.y + BALL_R) / CH));
        for (var r = r0; r <= r1; r++) {
            for (var c = c0; c <= c1; c++) if (s.grid[r][c]) hitBlock(s, s.grid[r][c], c * CW, cellY(r));
        }
        for (var i = 0; i < s.falling.length; i++) hitBlock(s, s.falling[i], s.falling[i].c * CW, s.falling[i].y);
    }

    function stepBall(s, h) {
        var b = s.ball;
        moveCab(s, h);
        b.vy += GRAV * h; b.vx += s.wind.f * h;
        b.vx -= b.vx * 0.12 * h; b.vy -= b.vy * 0.12 * h;
        b.x += b.vx * h; b.y += b.vy * h;
        rope(s, h);
        collideBlocks(s);
        snag(s);
        if (b.y < G.HUD + BALL_R) { b.y = G.HUD + BALL_R; b.vy = Math.abs(b.vy) * 0.4; }
        if (b.y > GROUND - BALL_R) { b.y = GROUND - BALL_R; b.vy = -Math.abs(b.vy) * 0.3; }
    }

    // ------------------------------------------------------------------
    // Hazards of the late levels
    // ------------------------------------------------------------------

    // Gusts push the ball sideways for a few seconds; a gust can ruin a
    // swing or, timed well, feed it.
    function updateWind(s, dt) {
        var w = s.wind;
        if (w.t > 0) {
            w.t -= dt;
            if (w.t <= 0) w.f = 0;
            return;
        }
        w.next -= dt;
        if (w.next > 0) return;
        w.f = G.pick([-1, 1]) * 230; w.t = 2.6; w.next = G.rnd(4, 7);
        G.noise(2.4, { filter: 'bandpass', freq: 500, slide: 1400, q: 2, vol: 0.14, attack: 0.6 });
    }

    function topBlocks(s) {
        var tops = [];
        for (var c = 0; c < COLS; c++) {
            for (var r = ROWS - 1; r >= 0; r--) {
                var b = s.grid[r][c];
                if (!b) continue;
                if (firm(b) && b.kind !== CHARGE) tops.push(b);
                break;
            }
        }
        return tops;
    }

    // A tremor shakes roof blocks loose. The first is thrown at where the
    // cab stands right now, so standing still after the warning is what
    // gets punished; the others land at random.
    function quake(s) {
        var tops = topBlocks(s), n = Math.min(tops.length, G.level >= 10 ? 3 : 2);
        G.shake(12, 0.5); G.sfx('boom');
        for (var i = 0; i < n; i++) {
            var b = tops.splice(Math.floor(Math.random() * tops.length), 1)[0];
            var x = blockX(b), y = b.y + CH / 2, tx = i === 0 ? s.cab.x : G.rnd(60, G.W - 60);
            removeBlock(s, b);
            s.gone++;
            addChunk(s, x, y, (tx - x) / flightTime(y, -240), -240);
        }
    }

    function updateTremor(s, dt) {
        var t = s.tremor;
        t.next -= dt;
        if (t.next < 1.4 && !t.warn) {
            t.warn = true;
            G.noise(1.4, { freq: 90, slide: 240, vol: 0.3, attack: 0.5 });
        }
        if (t.warn) G.shake(2.5, 0.1);
        if (t.next > 0) return;
        quake(s);
        t.next = G.rnd(8, 11); t.warn = false;
    }

    function updateMotes(s, dt) {
        s.motes.forEach(function (m) {
            m.x += (m.v + s.wind.f * 0.6) * dt;
            if (m.x > G.W + 5) m.x = -5;
            if (m.x < -5) m.x = G.W + 5;
        });
    }

    // ------------------------------------------------------------------
    // Update
    // ------------------------------------------------------------------

    function checkEnd(s, dt) {
        if (s.done > 0) {
            s.done -= dt;
            if (s.done <= 0) G.win(Math.ceil(s.time) * 10 + G.lives * 200);
            return;
        }
        if (s.gone >= s.need) {
            // A short pause lets the last collapse play out before the banner.
            s.done = 1.2;
            G.flash('#ffffff', 0.15);
            G.popup(G.W / 2, 200, 'TARGET DEMOLISHED', '#ffffff');
            return;
        }
        var before = s.time;
        s.time -= dt;
        if (s.time <= 10 && Math.ceil(before) !== Math.ceil(s.time)) G.sfx('alarm');
        if (s.time <= 0) G.die();
    }

    function tickTimers(s, dt) {
        for (var k in s.cool) if (s.cool[k] > 0) s.cool[k] -= dt;
        if (s.comboT > 0) s.comboT -= dt;
        if (s.cab.inv > 0) s.cab.inv -= dt;
    }

    function update(s, dt) {
        tickTimers(s, dt);
        driveCab(s, dt);
        reelChain(s, dt);
        // Sub-steps keep a fast ball from skipping through a 30 px block.
        for (var i = 0; i < SUB; i++) stepBall(s, dt / SUB);
        if (s.dirty) settle(s);
        updateLoose(s, dt);
        updateFalling(s, dt);
        updateFuses(s, dt);
        updateChunks(s, dt);
        if (s.cfg.wind) updateWind(s, dt);
        if (s.cfg.tremor) updateTremor(s, dt);
        updateMotes(s, dt);
        checkEnd(s, dt);
    }

    // ------------------------------------------------------------------
    // Drawing
    // ------------------------------------------------------------------

    function drawSky(s, ctx) {
        ctx.fillStyle = '#0a0a0a';
        ctx.fillRect(0, 0, G.W, G.H);
        s.sky.forEach(function (b) {
            ctx.fillStyle = b.tone;
            ctx.fillRect(b.x, GROUND - b.h, b.w, b.h);
        });
        // Smog that slowly breathes: also keeps the picture alive while the
        // ball hangs still.
        ctx.fillStyle = 'rgba(255,255,255,' + (0.045 + 0.03 * Math.sin(G.t * 2.2)) + ')';
        ctx.fillRect(0, 230, G.W, GROUND - 230);
        ctx.fillStyle = '#5a5a5a';
        s.motes.forEach(function (m) { ctx.fillRect(m.x, m.y + Math.sin(G.t + m.x * 0.02) * 6, m.size, m.size); });
    }

    function drawBar(ctx, x, w, part, color, label) {
        ctx.fillStyle = '#000000';
        ctx.fillRect(x, 508, w, 16);
        ctx.fillStyle = color;
        ctx.fillRect(x, 508, w * G.clamp(part, 0, 1), 16);
        ctx.strokeStyle = '#ffffff'; ctx.lineWidth = 2;
        ctx.strokeRect(x, 508, w, 16);
        G.text(label, x, 502, { size: 13, color: '#c8c8c8' });
    }

    function drawGround(s, ctx) {
        ctx.fillStyle = '#111111';
        ctx.fillRect(0, GROUND, G.W, G.H - GROUND);
        ctx.fillStyle = '#ffffff';
        ctx.fillRect(0, GROUND, G.W, 4);
        var low = s.time <= 10 && Math.floor(G.t * 4) % 2 === 0;
        drawBar(ctx, 20, 600, s.gone / s.total, RED, 'DEMOLISHED ' + Math.floor(s.gone / s.total * 100) + '%  —  TARGET ' + Math.round(s.cfg.target * 100) + '%');
        drawBar(ctx, 660, 280, s.time / s.cfg.time, low ? RED : '#ffffff', 'TIME ' + Math.max(0, Math.ceil(s.time)));
        // The mark the red bar has to reach.
        ctx.fillStyle = '#ffffff';
        ctx.fillRect(20 + 600 * s.need / s.total - 2, 503, 4, 26);
    }

    function drawCracks(ctx, b, x, y) {
        var part = b.hp / KINDS[b.kind].hp, n = part < 0.4 ? 3 : (part < 0.75 ? 2 : 1), sx = x + 7 + b.shade * 14;
        ctx.strokeStyle = '#000000'; ctx.lineWidth = 1.5;
        ctx.beginPath();
        for (var i = 0; i < n; i++) {
            ctx.moveTo(sx + i * 9, y + 2); ctx.lineTo(sx + i * 9 - 5, y + CH * 0.45); ctx.lineTo(sx + i * 9 + 4, y + CH - 2);
        }
        ctx.stroke();
    }

    function drawConcrete(ctx, b, x, y) {
        var rebar = b.kind === REBAR, g = rebar ? 70 + Math.floor(b.shade * 20) : 128 + Math.floor(b.shade * 56);
        ctx.fillStyle = 'rgb(' + g + ',' + g + ',' + g + ')';
        ctx.fillRect(x, y, CW, CH);
        // Board-marked concrete: a light top edge and a shadowed bottom.
        ctx.fillStyle = 'rgba(255,255,255,0.18)';
        ctx.fillRect(x, y, CW, 3);
        ctx.fillStyle = 'rgba(0,0,0,0.28)';
        ctx.fillRect(x, y + CH - 5, CW, 5);
        if (rebar) {
            ctx.fillStyle = RED;
            for (var i = 0; i < 3; i++) ctx.fillRect(x + 8 + i * 11, y + 3, 3, CH - 6);
        }
        if (b.hp < KINDS[b.kind].hp) drawCracks(ctx, b, x, y);
        ctx.strokeStyle = '#000000'; ctx.lineWidth = 2;
        ctx.strokeRect(x + 1, y + 1, CW - 2, CH - 2);
    }

    function drawGlass(ctx, x, y) {
        ctx.fillStyle = 'rgba(255,255,255,0.16)';
        ctx.fillRect(x + 2, y + 2, CW - 4, CH - 4);
        ctx.strokeStyle = '#ffffff'; ctx.lineWidth = 2;
        ctx.strokeRect(x + 2, y + 2, CW - 4, CH - 4);
        ctx.beginPath();
        ctx.moveTo(x + 8, y + CH - 8); ctx.lineTo(x + 20, y + 7);
        ctx.moveTo(x + 17, y + CH - 8); ctx.lineTo(x + 25, y + 12);
        ctx.stroke();
    }

    function drawCharge(ctx, b, x, y) {
        var lit = b.fuse > 0, blink = lit && Math.floor(b.fuse * 10) % 2 === 0;
        ctx.fillStyle = blink ? '#ffffff' : RED;
        ctx.fillRect(x, y, CW, CH);
        ctx.fillStyle = '#000000';
        ctx.fillRect(x, y + 4, CW, 3); ctx.fillRect(x, y + CH - 7, CW, 3);
        G.text(lit ? b.fuse.toFixed(1) : 'TNT', x + CW / 2, y + 20, { size: 13, bold: true, color: '#000000', align: 'center' });
        ctx.strokeStyle = '#ffffff'; ctx.lineWidth = 2;
        ctx.strokeRect(x + 1, y + 1, CW - 2, CH - 2);
    }

    function drawBlock(ctx, b) {
        // A block about to fall shivers first: the warning to get clear.
        var x = b.c * CW + (b.loose > 0 ? G.rnd(-1.5, 1.5) : 0), y = b.y;
        if (b.dead) return;
        if (b.kind === GLASS) drawGlass(ctx, x, y);
        else if (b.kind === CHARGE) drawCharge(ctx, b, x, y);
        else drawConcrete(ctx, b, x, y);
    }

    function drawChunks(s, ctx) {
        s.chunks.forEach(function (k) {
            // Landing mark on the ground.
            ctx.fillStyle = RED;
            ctx.fillRect(k.tx - 16, GROUND + 6, 32, 5);
            ctx.save();
            ctx.translate(k.x, k.y); ctx.rotate(k.rot);
            ctx.fillStyle = '#9a9a9a';
            ctx.fillRect(-k.size / 2, -k.size / 2, k.size, k.size);
            ctx.strokeStyle = '#000000'; ctx.lineWidth = 2;
            ctx.strokeRect(-k.size / 2, -k.size / 2, k.size, k.size);
            ctx.restore();
        });
    }

    function drawMast(s, ctx) {
        var x = s.cab.x, top = PIVOT_Y - 8, foot = GROUND - CAB_H;
        ctx.strokeStyle = '#d8d8d8'; ctx.lineWidth = 2;
        ctx.beginPath();
        ctx.moveTo(x - 7, foot); ctx.lineTo(x - 7, top);
        ctx.moveTo(x + 7, foot); ctx.lineTo(x + 7, top);
        for (var y = foot, side = 1; y > top + 26; y -= 26, side = -side) {
            ctx.moveTo(x - 7 * side, y); ctx.lineTo(x + 7 * side, y - 26);
        }
        ctx.stroke();
        ctx.fillStyle = '#ffffff';
        ctx.fillRect(x - 16, top - 4, 32, 8);
        // Warning beacon on the mast head.
        ctx.fillStyle = Math.floor(G.t * 3) % 2 === 0 ? RED : '#551100';
        ctx.fillRect(x - 4, top - 12, 8, 8);
    }

    function drawCab(s, ctx) {
        var x = s.cab.x;
        // Blinks while it is invulnerable after a hit.
        if (s.cab.inv > 0 && Math.floor(s.cab.inv * 12) % 2 === 0) return;
        ctx.fillStyle = '#1e1e1e';
        ctx.fillRect(x - 34, GROUND - 13, 68, 15);
        ctx.strokeStyle = '#ffffff'; ctx.lineWidth = 2;
        ctx.strokeRect(x - 34, GROUND - 13, 68, 15);
        ctx.fillStyle = '#ffffff';
        for (var i = 0; i < 6; i++) {
            // Track links crawl with the cab so driving reads at a glance.
            var lx = ((i * 12 + x * 0.6) % 66 + 66) % 66;
            ctx.fillRect(x - 33 + lx, GROUND - 2, 4, 3);
        }
        ctx.fillRect(x - 26, GROUND - CAB_H, 52, 22);
        ctx.fillStyle = '#000000';
        ctx.fillRect(x + 4, GROUND - CAB_H + 4, 18, 11);
        ctx.fillStyle = RED;
        ctx.fillRect(x - 26, GROUND - 19, 52, 5);
    }

    function drawBall(s, ctx) {
        var b = s.ball, speed = Math.hypot(b.vx, b.vy);
        ctx.strokeStyle = '#ffffff'; ctx.lineWidth = 3;
        ctx.setLineDash([7, 4]);
        ctx.beginPath(); ctx.moveTo(s.cab.x, PIVOT_Y); ctx.lineTo(b.x, b.y); ctx.stroke();
        ctx.setLineDash([]);
        ctx.fillStyle = '#2c2c2c';
        ctx.beginPath(); ctx.arc(b.x, b.y, BALL_R, 0, Math.PI * 2); ctx.fill();
        // The rim turns red once the ball is fast enough to crack rebar.
        ctx.strokeStyle = speed >= KINDS[REBAR].minV ? RED : '#ffffff';
        ctx.stroke();
        ctx.strokeStyle = 'rgba(255,255,255,0.5)'; ctx.lineWidth = 2;
        ctx.beginPath(); ctx.arc(b.x, b.y, BALL_R - 6, Math.PI * 1.1, Math.PI * 1.6); ctx.stroke();
    }

    function drawWind(s, ctx) {
        var f = s.wind.f;
        if (!f) return;
        ctx.strokeStyle = 'rgba(255,255,255,0.4)'; ctx.lineWidth = 2;
        ctx.beginPath();
        for (var i = 0; i < 12; i++) {
            var x = ((i * 137 + G.t * f * 2.4) % G.W + G.W) % G.W, y = 70 + i * 31;
            ctx.moveTo(x, y); ctx.lineTo(x - (f > 0 ? 46 : -46), y);
        }
        ctx.stroke();
        G.text(f > 0 ? 'WIND >>>' : '<<< WIND', G.W / 2, 54, { size: 18, bold: true, color: '#ffffff', align: 'center' });
    }

    function drawNotes(s, ctx) {
        if (s.tremor.warn && Math.floor(G.t * 8) % 2 === 0) {
            G.text('TREMOR — KEEP MOVING', G.W / 2, 78, { size: 20, bold: true, color: RED, align: 'center' });
        }
        if (G.level === 1 && G.t < 9) {
            G.text('DRIVE TO SWING THE BALL · BRAKE TO WHIP IT · HIT THE BOTTOM ROW', G.W / 2, 54, { size: 16, color: '#c8c8c8', align: 'center' });
        }
    }

    function draw(s, ctx) {
        drawSky(s, ctx);
        drawGround(s, ctx);
        eachBlock(s, function (b) { drawBlock(ctx, b); });
        drawChunks(s, ctx);
        drawMast(s, ctx);
        drawCab(s, ctx);
        drawBall(s, ctx);
        drawWind(s, ctx);
        drawNotes(s, ctx);
    }

    function hud(s) {
        return 'DEMO ' + Math.floor(s.gone / s.total * 100) + '/' + Math.round(s.cfg.target * 100) + '%  T ' + Math.max(0, Math.ceil(s.time));
    }

    G.register('brutalist', {
        title: 'WRECKING BALL',
        blurb: 'Swing the ball, knock out the support, bring the concrete down before time runs out.',
        controls: [
            '← → drive the crane: the ball swings by its own momentum',
            '↑ ↓ reel the chain in / out · SPACE (touch: A) hard brake, whips the ball forward',
            'Red-barred rebar needs a fast hit · TNT blows 1.6 s after it is struck',
            'Stay out from under falling slabs and flying debris'
        ],
        levelNames: ['Garden Wall', 'Twin Stacks', 'Glass House', 'Rebar Core', 'Arcade',
            'Short Fuse', 'Silos', 'Crosswind', 'Aftershock', 'Megastructure'],
        colors: { bg: '#0a0a0a', fg: '#ffffff', accent: RED, dim: '#9a9a9a' },
        lives: 3,
        // A slow, heavy dorian stomp: sawtooth bass, a clanking off-beat
        // snare and an eight-bar lead that climbs and falls back to the root.
        music: {
            bpm: 104, root: 38, scale: 'dorian', prog: [0, 0, 3, 0, 6, 3, 4, 0],
            bass: 'x..x..x...x.x.o.',
            lead: ['0...2.3...2.0...', '4---....3.2.0...', '5...4.3...2.3.4.', '7---....4-..2-..',
                '6...5.4...2.4...', '5---3---2---0...', '4.4.3.2.4.3.2.0.', '0-------........'],
            arp: '0..2..1.',
            drums: { k: 'x.....x.x.....x.', s: '....x.......x..x', h: 'x.xxx.x.x.xxx.x.' },
            leadWave: 'sawtooth', bassWave: 'sawtooth', arpWave: 'triangle', leadOct: 2
        },
        init: init, update: update, draw: draw, hud: hud
    });
})();
