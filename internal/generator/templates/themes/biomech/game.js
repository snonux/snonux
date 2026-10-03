/*
 * Spine Crawler — the biomech theme's game.
 *
 * A segmented worm winds down through a field of bone nodes: it crawls along
 * a row until a node or a wall blocks it, then drops a row and turns round.
 * The turret roams the bottom quarter and spits acid upward. A shot segment
 * calcifies into a new node and the worm splits in two at the wound, so every
 * hit reshapes the maze the rest of the worm has to crawl through.
 *
 * Around the worm: spore sacs fall and seed nodes, parasites zigzag through
 * the turret's zone eating nodes, and a fast eye poisons nodes so that a worm
 * touching one dives straight down. Level 10 adds the Queen, who breeds worms
 * and only takes damage in the head.
 */
(function () {
    'use strict';
    var G = window.SnoGame;
    var CELL = 24, COLS = 40, ROWS = 21, OY = 36;       // 40x21 cells below the HUD
    var ZONE = 16, ZONE_Y = OY + ZONE * CELL;           // turret territory: the bottom five rows
    var QROW = 5;                                       // first row below the Queen's body
    var NODE_HP = 3, MAX_NODES = 170, MAX_SHOTS = 5;
    var PLAYER_SPEED = 340, SHOT_SPEED = 820, FIRE_DELAY = 0.12, INVULN = 2.5;
    var BONE = '#d0c7bb', FLESH = '#803f5d', VEIN = '#f55b7d', ACID = '#93ffd8', STEEL = '#2d3642', BG = '#09070d';

    // ------------------------------------------------------------------
    // Level table
    // ------------------------------------------------------------------

    // One number per cell that is the same on every visit, so layouts can be
    // symmetric or patterned without consuming a shared random stream.
    function hash(c, r, seed) { return G.rng(seed * 7919 + c * 131 + r * 17)(); }

    // Layouts return 0 (empty), 1 (bone node) or 2 (poisoned node).
    function scatter(density) {
        return function (c, r, seed) { return hash(c, r, seed) < density ? 1 : 0; };
    }
    function mirror(density) {
        return function (c, r, seed) { return hash(Math.min(c, COLS - 1 - c), r, seed) < density ? 1 : 0; };
    }
    // Broken vertical ribs: lanes that hurry a worm downward, with gaps every
    // third row so it can still slip sideways into the next lane.
    function ribs(c, r, seed) {
        if (c % 8 === 4 && r >= 2 && r <= 13 && r % 3 !== 0) return 1;
        return hash(c, r, seed) < 0.02 ? 1 : 0;
    }
    // Diagonal seams of venom through an ordinary scatter.
    function bands(c, r, seed) {
        var h = hash(c, r, seed);
        if ((c + r * 2) % 13 === 0 && r > 1 && r < 13 && h < 0.6) return 2;
        return h < 0.06 ? 1 : 0;
    }

    // side -1 enters from the left edge, +1 from the right; `at` is seconds
    // into the level (the next wave comes early when the field is empty).
    function wave(len, side, row, at) { return { len: len, side: side, row: row, at: at }; }

    function swarm() {
        var out = [];
        // Three batches of six short worms, staggered over the top three rows.
        for (var i = 0; i < 18; i++) out.push(wave(i < 6 ? 4 : 3, i % 2 ? 1 : -1, i % 3, Math.floor(i / 6) * 15));
        return out;
    }

    // speed is cells per second; spore / para / eye are spawn intervals in
    // seconds (absent = that creature does not appear); regen is the seconds
    // between healing-and-sprouting pulses; queen is her head's hit points.
    var LEVELS = [
        { speed: 6, layout: scatter(0.045), waves: [wave(10, -1, 0, 0), wave(8, 1, 0, 16)] },
        { speed: 6.5, layout: scatter(0.05), spore: 5, waves: [wave(12, 1, 0, 0), wave(10, -1, 0, 18), wave(8, 1, 0, 34)] },
        { speed: 7, layout: scatter(0.055), spore: 8, para: 7, waves: [wave(12, -1, 0, 0), wave(11, 1, 0, 18), wave(9, -1, 0, 34)] },
        { speed: 7.5, layout: mirror(0.06), spore: 8, para: 8,
            waves: [wave(8, -1, 0, 0), wave(8, 1, 0, 0), wave(8, -1, 0, 18), wave(8, 1, 0, 18),
                wave(6, -1, 0, 36), wave(6, 1, 0, 36)] },
        { speed: 8, layout: scatter(0.06), spore: 9, para: 8, eye: 8,
            waves: [wave(14, -1, 0, 0), wave(1, 1, 0, 5), wave(1, -1, 2, 7), wave(12, 1, 0, 22), wave(10, -1, 0, 40)] },
        { speed: 8.5, layout: scatter(0.09), regen: 2.5, spore: 10, para: 8, eye: 12,
            waves: [wave(10, -1, 0, 0), wave(10, 1, 0, 0), wave(13, -1, 0, 22), wave(10, 1, 0, 40)] },
        { speed: 10.5, layout: scatter(0.05), spore: 8, para: 6, eye: 14, waves: swarm() },
        { speed: 10, layout: ribs, spore: 7, para: 6, eye: 9, waves: [wave(18, -1, 0, 0), wave(14, 1, 0, 22), wave(12, -1, 0, 42)] },
        { speed: 11, layout: bands, regen: 2.2, spore: 7, para: 6, eye: 6,
            waves: [wave(12, -1, 0, 0), wave(12, 1, 0, 0), wave(12, -1, 0, 20), wave(12, 1, 0, 20)] },
        { speed: 10, layout: scatter(0.05), top: QROW + 1, queen: 24, regen: 4, spore: 8, para: 8, eye: 11,
            waves: [wave(8, -1, QROW + 1, 0)] }
    ];

    // ------------------------------------------------------------------
    // Sounds: wet noise for flesh, short clicks for bone
    // ------------------------------------------------------------------

    var SND = {
        spit: function () { G.noise(0.06, { filter: 'bandpass', freq: 2200, slide: 500, q: 5, vol: 0.13 }); },
        squelch: function () {
            G.noise(0.16, { freq: 900, slide: 140, q: 3, vol: 0.3 });
            G.tone(170, 0.14, { type: 'sine', slide: 60, vol: 0.22 });
        },
        chip: function () { G.noise(0.03, { filter: 'bandpass', freq: 3200, q: 6, vol: 0.12 }); },
        crack: function () {
            G.noise(0.08, { filter: 'highpass', freq: 2400, vol: 0.16 });
            G.tone(820, 0.05, { type: 'triangle', slide: 300, vol: 0.12 });
        },
        hiss: function () { G.noise(0.5, { filter: 'bandpass', freq: 500, slide: 2600, q: 2, vol: 0.14, attack: 0.2 }); },
        thump: function () { G.tone(74, 0.07, { type: 'sine', slide: 46, vol: 0.2 }); },
        beat: function () {
            G.tone(58, 0.13, { type: 'sine', slide: 38, vol: 0.4 });
            G.tone(52, 0.13, { type: 'sine', slide: 34, vol: 0.3, delay: 0.16 });
        },
        spore: function () { G.tone(980, 0.5, { type: 'sine', slide: 180, vol: 0.09 }); },
        skitter: function () { G.noise(0.03, { filter: 'bandpass', freq: 4200, q: 8, vol: 0.1 }); },
        eye: function () { G.tone(1300, 0.45, { type: 'sawtooth', slide: 520, vol: 0.07 }); },
        poison: function () { G.tone(320, 0.12, { type: 'sine', slide: 900, vol: 0.12 }); },
        dive: function () {
            G.tone(700, 0.4, { type: 'sawtooth', slide: 90, vol: 0.14 });
            G.noise(0.3, { freq: 1600, slide: 200, vol: 0.14 });
        },
        gulp: function () { G.noise(0.1, { freq: 500, slide: 1400, q: 4, vol: 0.14 }); },
        sprout: function () { G.noise(0.12, { filter: 'bandpass', freq: 700, slide: 1500, q: 6, vol: 0.07 }); },
        clink: function () { G.tone(1700, 0.04, { type: 'triangle', vol: 0.09 }); },
        roar: function () {
            G.noise(0.7, { freq: 420, slide: 70, q: 2, vol: 0.4, attack: 0.08 });
            G.tone(96, 0.6, { type: 'sawtooth', slide: 48, vol: 0.18 });
        },
        queenHit: function () {
            G.noise(0.14, { freq: 1300, slide: 200, vol: 0.3 });
            G.tone(240, 0.16, { type: 'sawtooth', slide: 110, vol: 0.16 });
        }
    };

    // ------------------------------------------------------------------
    // Node field
    // ------------------------------------------------------------------

    function idx(c, r) { return r * COLS + c; }
    function cellX(c) { return c * CELL + CELL / 2; }
    function cellY(r) { return OY + r * CELL + CELL / 2; }

    function nodeAt(s, c, r) {
        return c < 0 || c >= COLS || r < 0 || r >= ROWS ? 0 : s.hp[idx(c, r)];
    }

    // The top rows stay open so worms can enter, the bottom row so the turret
    // always has a lane, and the turret's zone is thinned to a few stragglers.
    function buildNodes(s, level) {
        var top = s.cfg.top || 1;
        for (var r = 0; r < ROWS; r++) {
            for (var c = 0; c < COLS; c++) {
                var kind = r < top || r > ROWS - 2 ? 0 : s.cfg.layout(c, r, level);
                if (r >= ZONE && hash(c, r, level + 50) > 0.3) kind = 0;
                s.hp.push(kind ? NODE_HP : 0);
                s.poison.push(kind === 2);
            }
        }
    }

    function calcify(s, c, r, hp) {
        if (c < 0 || c >= COLS || r < 0 || r >= ROWS || s.hp[idx(c, r)]) return false;
        s.hp[idx(c, r)] = hp;
        s.poison[idx(c, r)] = false;
        return true;
    }

    function removeNode(s, c, r) {
        s.hp[idx(c, r)] = 0;
        s.poison[idx(c, r)] = false;
    }

    function damageNode(s, c, r) {
        var left = --s.hp[idx(c, r)];
        if (left > 0) { SND.chip(); return; }
        removeNode(s, c, r);
        SND.crack();
        G.addScore(1);
        G.burst(cellX(c), cellY(r), { n: 6, color: BONE, speed: 120, life: 0.35, gravity: 300 });
    }

    function countNodes(s) {
        var n = 0;
        for (var i = 0; i < s.hp.length; i++) if (s.hp[i]) n++;
        return n;
    }

    // Grows one bud next to an existing node, above the turret's zone.
    function sprout(s) {
        for (var tries = 0; tries < 24; tries++) {
            var c = Math.floor(G.rnd(0, COLS)), r = Math.floor(G.rnd(s.cfg.top || 1, ZONE));
            if (!nodeAt(s, c, r)) continue;
            var side = G.pick([[1, 0], [-1, 0], [0, 1], [0, -1]]);
            var nc = c + side[0], nr = r + side[1];
            if (nr < (s.cfg.top || 1) || nr >= ZONE || !calcify(s, nc, nr, 1)) continue;
            G.burst(cellX(nc), cellY(nr), { n: 5, color: FLESH, speed: 50, life: 0.5 });
            SND.sprout();
            return;
        }
    }

    // Regenerating levels: chipped nodes knit back together and the field
    // slowly spreads, so half-cleared lanes close again.
    function regenerate(s, dt) {
        if (!s.cfg.regen) return;
        s.regen -= dt;
        if (s.regen > 0) return;
        s.regen = s.cfg.regen;
        for (var i = 0; i < s.hp.length; i++) if (s.hp[i] > 0 && s.hp[i] < NODE_HP) s.hp[i]++;
        if (countNodes(s) < MAX_NODES) sprout(s);
    }

    // ------------------------------------------------------------------
    // Worms
    // ------------------------------------------------------------------
    //
    // Every segment carries the complete crawling state (target cell, the
    // cell it came from, heading, vertical direction, dive flag). Followers
    // copy the segment ahead on each step, so any segment can become a head
    // the moment the worm is cut in front of it.

    function segment(cx, cy, px, py, dir) {
        return { cx: cx, cy: cy, px: px, py: py, dir: dir, vd: 1, dive: false };
    }

    // A worm queued up beyond a side wall, crawling in along `row`.
    function sideWorm(len, side, row, speed) {
        var segs = [], dir = -side;
        for (var k = 0; k < len; k++) {
            var cx = side < 0 ? -1 - k : COLS + k;
            segs.push(segment(cx, row, cx - dir, row, dir));
        }
        return { segs: segs, t: 0, speed: speed, gate: 0 };
    }

    // A worm stacked up inside the Queen; segments above the gate row are
    // hidden by her body and crawl straight down until they emerge.
    function queenWorm(col, len, speed) {
        var segs = [], dir = Math.random() < 0.5 ? 1 : -1;
        for (var k = 0; k < len; k++) segs.push(segment(col, QROW - k, col, QROW - k - 1, dir));
        return { segs: segs, t: 0, speed: speed, gate: QROW };
    }

    function hidden(w, g) { return g.cx < 0 || g.cx >= COLS || g.cy < w.gate; }
    function segX(w, g) { return cellX(G.lerp(g.px, g.cx, w.t)); }
    function segY(w, g) { return cellY(G.lerp(g.py, g.cy, w.t)); }

    // Drops (or, once it has reached the floor, climbs) one row and reverses.
    // Inside the turret's zone the worm bounces between the floor and ZONE.
    function turn(h) {
        h.dir = -h.dir;
        if (h.vd > 0 && h.cy + 1 > ROWS - 1) h.vd = -1;
        else if (h.vd < 0 && h.cy - 1 < ZONE) h.vd = 1;
        h.cy += h.vd;
    }

    function advanceHead(s, w, h) {
        h.px = h.cx; h.py = h.cy;
        if (h.cy < w.gate) { h.cy++; return; }
        if (h.dive) { h.cy++; h.dive = h.cy < ZONE; return; }
        var nx = h.cx + h.dir;
        var wall = (nx < 0 && h.dir < 0) || (nx >= COLS && h.dir > 0);
        if (!wall && !nodeAt(s, nx, h.cy)) { h.cx = nx; return; }
        // Venom drives the worm mad: it plunges straight for the turret.
        if (!wall && s.poison[idx(nx, h.cy)] && h.cy < ZONE - 1) {
            h.dive = true; h.cy++;
            SND.dive();
            return;
        }
        turn(h);
        if (s.snd.thump <= 0) { SND.thump(); s.snd.thump = 0.14; }
    }

    function stepWorm(s, w) {
        for (var k = w.segs.length - 1; k > 0; k--) {
            var g = w.segs[k], lead = w.segs[k - 1];
            g.px = g.cx; g.py = g.cy;
            g.cx = lead.cx; g.cy = lead.cy; g.dir = lead.dir; g.vd = lead.vd; g.dive = lead.dive;
        }
        advanceHead(s, w, w.segs[0]);
    }

    function updateWorms(s, dt) {
        s.worms.forEach(function (w) {
            w.t += w.speed * dt;
            while (w.t >= 1) { w.t -= 1; stepWorm(s, w); }
        });
    }

    function liveSegments(s) {
        return s.worms.reduce(function (n, w) { return n + w.segs.length; }, 0);
    }

    function segmentsLeft(s) {
        return liveSegments(s) + s.queue.reduce(function (n, q) { return n + q.len; }, 0);
    }

    // The wound calcifies into a node; whatever was behind it becomes a worm
    // of its own, whose new head runs straight into that node and turns.
    function killSegment(s, w, k) {
        var g = w.segs[k], x = segX(w, g), y = segY(w, g);
        var rear = w.segs.splice(k).slice(1);
        calcify(s, g.cx, g.cy, NODE_HP);
        if (rear.length) s.worms.push({ segs: rear, t: w.t, speed: w.speed, gate: w.gate });
        var points = (k === 0 ? 100 : 10) * G.level;
        G.addScore(points);
        if (k === 0) G.popup(x, y - 10, points, ACID);
        G.burst(x, y, { n: 14, color: k === 0 ? VEIN : FLESH, speed: 190, life: 0.5, gravity: 260 });
        SND.squelch();
    }

    function spawnWaves(s, dt) {
        s.clock += dt;
        while (s.queue.length && (s.clock >= s.queue[0].at || !s.worms.length)) {
            var q = s.queue.shift();
            // Jumping the clock forward keeps waves that belong together together.
            s.clock = Math.max(s.clock, q.at);
            s.worms.push(sideWorm(q.len, q.side, q.row, s.cfg.speed));
            SND.hiss();
        }
    }

    // ------------------------------------------------------------------
    // Creatures: spore sacs, parasites and the eye
    // ------------------------------------------------------------------

    function spawnSpore(s) {
        var c = Math.floor(G.rnd(0, COLS));
        s.bugs.push({ kind: 'spore', x: cellX(c), y: OY - 12, vx: 0, vy: 220 + G.level * 8, r: 9, hp: 2, col: c, row: -1 });
        SND.spore();
    }

    function spawnPara(s) {
        var side = Math.random() < 0.5 ? -1 : 1;
        s.bugs.push({
            kind: 'para', x: side < 0 ? -14 : G.W + 14, y: ZONE_Y + G.rnd(0, 60),
            vx: -side * G.rnd(80, 130), vy: 190, r: 11, hp: 1, turn: 0.4
        });
    }

    function spawnEye(s) {
        var side = Math.random() < 0.5 ? -1 : 1, top = (s.cfg.top || 1) + 1;
        var row = Math.floor(G.rnd(top, ZONE - 2));
        s.bugs.push({ kind: 'eye', x: side < 0 ? -14 : G.W + 14, y: cellY(row), vx: -side * (300 + G.level * 10), vy: 0, r: 10, hp: 1, row: row });
        SND.eye();
    }

    function hasBug(s, kind) {
        return s.bugs.some(function (b) { return b.kind === kind; });
    }

    // Each creature has its own countdown; parasites and eyes come one at a
    // time so the turret's zone never becomes undodgeable.
    function spawnBugs(s, dt) {
        var spawners = { spore: spawnSpore, para: spawnPara, eye: spawnEye };
        for (var kind in spawners) {
            if (!s.cfg[kind]) continue;
            s.next[kind] -= dt;
            if (s.next[kind] > 0 || (kind !== 'spore' && hasBug(s, kind))) continue;
            s.next[kind] = s.cfg[kind] * G.rnd(0.7, 1.3);
            spawners[kind](s);
        }
    }

    // A sac seeds a node in roughly every third cell it falls through.
    function moveSpore(s, b, dt) {
        b.y += b.vy * dt;
        var row = Math.floor((b.y - OY) / CELL);
        if (row === b.row) return;
        b.row = row;
        if (row >= (s.cfg.top || 1) && row < ROWS - 2 && Math.random() < 0.3 && countNodes(s) < MAX_NODES) calcify(s, b.col, row, NODE_HP);
    }

    // Parasites change their vertical mind every fraction of a second and
    // chew through any node they cross, which keeps the zone from clogging.
    function movePara(s, b, dt) {
        b.turn -= dt;
        if (b.turn <= 0) { b.turn = G.rnd(0.25, 0.7); b.vy = G.pick([-1, 1]) * G.rnd(150, 230); }
        b.x += b.vx * dt; b.y += b.vy * dt;
        var top = ZONE_Y - CELL * 3 + b.r, bottom = G.H - b.r;
        if (b.y < top) { b.y = top; b.vy = Math.abs(b.vy); }
        if (b.y > bottom) { b.y = bottom; b.vy = -Math.abs(b.vy); }
        var c = Math.floor(b.x / CELL), r = Math.floor((b.y - OY) / CELL);
        if (nodeAt(s, c, r)) { removeNode(s, c, r); SND.gulp(); }
        if (s.snd.skit <= 0) { SND.skitter(); s.snd.skit = 0.17; }
    }

    function moveEye(s, b, dt) {
        b.x += b.vx * dt;
        var c = Math.floor(b.x / CELL);
        if (!nodeAt(s, c, b.row) || s.poison[idx(c, b.row)]) return;
        s.poison[idx(c, b.row)] = true;
        SND.poison();
        G.burst(cellX(c), cellY(b.row), { n: 6, color: ACID, speed: 70, life: 0.4 });
    }

    function updateBugs(s, dt) {
        var movers = { spore: moveSpore, para: movePara, eye: moveEye };
        s.bugs = s.bugs.filter(function (b) {
            movers[b.kind](s, b, dt);
            return b.hp > 0 && b.y < G.H + 30 && b.x > -40 && b.x < G.W + 40;
        });
    }

    // Parasites are worth more the closer the player lets them come.
    function bugScore(s, b) {
        if (b.kind === 'spore') return 200;
        if (b.kind === 'eye') return 1000;
        var d = G.dist(b.x, b.y, s.player.x, s.player.y);
        return d < 70 ? 900 : (d < 150 ? 600 : 300);
    }

    function hurtBug(s, b) {
        b.hp--;
        if (b.hp > 0) { b.vy *= 1.6; SND.chip(); return; }      // a wounded sac drops faster
        var points = bugScore(s, b);
        G.addScore(points);
        G.popup(b.x, b.y - 12, points, ACID);
        G.burst(b.x, b.y, { n: 16, color: b.kind === 'eye' ? BONE : VEIN, speed: 210, life: 0.5 });
        SND.squelch();
    }

    // ------------------------------------------------------------------
    // The Queen
    // ------------------------------------------------------------------

    function newQueen(hp) {
        return { x: G.W / 2, y: OY + 48, hp: hp, max: hp, hurt: 0, brood: 5 };
    }

    function headY(q) { return q.y + 46; }

    // She sweeps the whole width so her head cannot be camped under, breeds
    // faster once wounded, and never lets the field stay empty for long.
    function updateQueen(s, dt) {
        var q = s.queen;
        if (!q) return;
        q.x = G.W / 2 + Math.sin(s.time * 0.4) * 330;
        if (q.hurt > 0) q.hurt -= dt;
        q.brood -= dt;
        if (!s.worms.length) q.brood = Math.min(q.brood, 1.2);
        if (q.brood > 0 || liveSegments(s) >= 16) return;
        var angry = q.hp < q.max / 2;
        q.brood = angry ? 6 : 8.5;
        var col = G.clamp(Math.floor(q.x / CELL), 1, COLS - 2);
        s.worms.push(queenWorm(col, angry ? 7 : 5, s.cfg.speed));
        SND.roar();
        G.burst(q.x, headY(q), { n: 10, color: VEIN, speed: 120, angle: Math.PI / 2, spread: 1.6 });
    }

    function killQueen(s) {
        var q = s.queen;
        s.queen = null;
        G.addScore(5000);
        G.popup(q.x, q.y, 5000, ACID);
        G.burst(q.x, q.y, { n: 70, color: VEIN, speed: 340, life: 1.1, gravity: 200 });
        G.burst(q.x, q.y, { n: 40, color: BONE, speed: 260, life: 1.2, gravity: 300 });
        G.shake(14, 0.7); G.flash(VEIN, 0.3);
        G.sfx('bigboom'); SND.roar();
    }

    // Only the head below the carapace is soft; the armour swallows shots.
    function hitQueen(s, sh) {
        var q = s.queen;
        if (!q) return false;
        if (Math.abs(sh.x - q.x) < 17 && Math.abs(sh.y - headY(q)) < 17) {
            q.hp--; q.hurt = 0.12;
            G.addScore(50);
            G.burst(sh.x, sh.y, { n: 8, color: VEIN, speed: 160, angle: Math.PI / 2, spread: 2 });
            SND.queenHit();
            if (q.hp <= 0) killQueen(s);
            return true;
        }
        if (Math.abs(sh.x - q.x) > 92 || sh.y > q.y + 34) return false;
        SND.clink();
        G.burst(sh.x, sh.y, { n: 3, color: BONE, speed: 90, life: 0.25 });
        return true;
    }

    // ------------------------------------------------------------------
    // Turret and shots
    // ------------------------------------------------------------------

    // The turret is a 16 px box against whole node cells.
    function touchesNode(s, x, y) {
        for (var i = 0; i < 4; i++) {
            var c = Math.floor((x + (i % 2 ? 8 : -8)) / CELL), r = Math.floor((y + (i < 2 ? -8 : 8) - OY) / CELL);
            if (nodeAt(s, c, r)) return true;
        }
        return false;
    }

    // Moves along one axis and stops at a node. A turret that already overlaps
    // one (a segment calcified on top of it) may always move, so it can never
    // be sealed in.
    function slide(s, p, dx, dy) {
        var nx = G.clamp(p.x + dx, 12, G.W - 12), ny = G.clamp(p.y + dy, ZONE_Y + 12, G.H - 12);
        if (touchesNode(s, nx, ny) && !touchesNode(s, p.x, p.y)) {
            if (dx) p.vx = 0; else p.vy = 0;
            return;
        }
        p.x = nx; p.y = ny;
    }

    // Velocity eases toward the stick, so the turret has a little weight
    // without feeling sluggish.
    function movePlayer(s, dt) {
        var p = s.player;
        var ax = (G.key.right ? 1 : 0) - (G.key.left ? 1 : 0), ay = (G.key.down ? 1 : 0) - (G.key.up ? 1 : 0);
        var norm = ax && ay ? Math.SQRT1_2 : 1, k = Math.min(1, 14 * dt);
        p.vx += (ax * norm * PLAYER_SPEED - p.vx) * k;
        p.vy += (ay * norm * PLAYER_SPEED - p.vy) * k;
        slide(s, p, p.vx * dt, 0);
        slide(s, p, 0, p.vy * dt);
        if (p.inv > 0) p.inv -= dt;
    }

    function fire(s, dt) {
        if (s.cool > 0) s.cool -= dt;
        if (!G.key.a || s.cool > 0 || s.shots.length >= MAX_SHOTS) return;
        s.cool = FIRE_DELAY;
        s.shots.push({ x: s.player.x, y: s.player.y - 14 });
        s.recoil = 1;
        SND.spit();
    }

    function hitNode(s, sh) {
        var c = Math.floor(sh.x / CELL), r = Math.floor((sh.y - OY) / CELL);
        if (!nodeAt(s, c, r)) return false;
        damageNode(s, c, r);
        return true;
    }

    function hitWorm(s, sh) {
        for (var i = 0; i < s.worms.length; i++) {
            var w = s.worms[i];
            for (var k = 0; k < w.segs.length; k++) {
                var g = w.segs[k];
                if (hidden(w, g) || Math.abs(segX(w, g) - sh.x) > 11 || Math.abs(segY(w, g) - sh.y) > 12) continue;
                killSegment(s, w, k);
                return true;
            }
        }
        return false;
    }

    function hitBug(s, sh) {
        for (var i = 0; i < s.bugs.length; i++) {
            var b = s.bugs[i];
            if (b.hp <= 0 || Math.abs(b.x - sh.x) > b.r + 2 || Math.abs(b.y - sh.y) > b.r + 4) continue;
            hurtBug(s, b);
            return true;
        }
        return false;
    }

    function updateShots(s, dt) {
        s.shots = s.shots.filter(function (sh) {
            sh.y -= SHOT_SPEED * dt;
            return sh.y > OY && !(hitNode(s, sh) || hitWorm(s, sh) || hitBug(s, sh) || hitQueen(s, sh));
        });
        s.worms = s.worms.filter(function (w) { return w.segs.length; });
    }

    function wormTouches(s, p) {
        return s.worms.some(function (w) {
            return w.segs.some(function (g) {
                return !hidden(w, g) && G.circ(p.x, p.y, 8, segX(w, g), segY(w, g), 8);
            });
        });
    }

    // A hit costs a life but the fight goes on: the turret regrows at the
    // bottom centre, briefly untouchable, and the loose creatures scatter.
    function checkPlayerHit(s) {
        var p = s.player;
        if (p.inv > 0) return;
        var hit = wormTouches(s, p) || s.bugs.some(function (b) { return G.circ(p.x, p.y, 8, b.x, b.y, b.r - 2); });
        if (!hit) return;
        G.burst(p.x, p.y, { n: 26, color: ACID, speed: 260, life: 0.7 });
        G.noise(0.4, { freq: 700, slide: 90, vol: 0.35 });
        if (G.loseLife() <= 0) return;
        p.x = G.W / 2; p.y = G.H - 36; p.vx = 0; p.vy = 0; p.inv = INVULN;
        s.bugs = [];
    }

    // ------------------------------------------------------------------
    // Init and update
    // ------------------------------------------------------------------

    function init(level) {
        var cfg = LEVELS[level - 1];
        var s = {
            cfg: cfg, hp: [], poison: [], worms: [], queue: cfg.waves.slice(), shots: [], bugs: [],
            queen: cfg.queen ? newQueen(cfg.queen) : null,
            player: { x: G.W / 2, y: G.H - 36, vx: 0, vy: 0, inv: 0 },
            next: { spore: cfg.spore || 0, para: cfg.para || 0, eye: cfg.eye || 0 },
            snd: { thump: 0, beat: 0, skit: 0 },
            regen: cfg.regen || 0, clock: 0, time: 0, cool: 0, recoil: 0
        };
        buildNodes(s, level);
        return s;
    }

    function inZone(s) {
        return s.worms.some(function (w) {
            return w.segs.some(function (g) { return !hidden(w, g) && g.cy >= ZONE; });
        });
    }

    // Gated timers for the repeating sounds; a heartbeat warns while a worm
    // is loose in the turret's zone.
    function tickSounds(s, dt) {
        for (var k in s.snd) if (s.snd[k] > 0) s.snd[k] -= dt;
        if (s.snd.beat > 0 || !inZone(s)) return;
        SND.beat();
        s.snd.beat = 0.62;
    }

    function update(s, dt) {
        s.time += dt;
        s.recoil = Math.max(0, s.recoil - dt * 9);
        tickSounds(s, dt);
        movePlayer(s, dt);
        fire(s, dt);
        updateShots(s, dt);
        spawnWaves(s, dt);
        updateWorms(s, dt);
        spawnBugs(s, dt);
        updateBugs(s, dt);
        updateQueen(s, dt);
        regenerate(s, dt);
        if (!s.queue.length && !s.worms.length && !s.queen) { G.win(500 + G.lives * 250); return; }
        checkPlayerHit(s);
    }

    // ------------------------------------------------------------------
    // Drawing
    // ------------------------------------------------------------------

    function circle(ctx, x, y, r) { ctx.beginPath(); ctx.arc(x, y, r, 0, Math.PI * 2); }

    // Slow veins drifting behind the field keep the tissue visibly alive.
    function drawBackdrop(s, ctx) {
        ctx.fillStyle = BG;
        ctx.fillRect(0, 0, G.W, G.H);
        ctx.lineWidth = 2;
        for (var i = 0; i < 6; i++) {
            var base = 70 + i * 82, pulse = 0.3 + 0.14 * Math.sin(s.time * 2.2 + i);
            ctx.strokeStyle = 'rgba(128,63,93,' + pulse.toFixed(3) + ')';
            ctx.beginPath();
            for (var x = 0; x <= G.W; x += 40) {
                var y = base + Math.sin(x * 0.011 + s.time * 0.6 + i * 1.7) * 20 + Math.sin(x * 0.027 - s.time * 0.45 + i) * 9;
                if (x) ctx.lineTo(x, y); else ctx.moveTo(x, y);
            }
            ctx.stroke();
        }
        // The membrane marking how far up the turret may roam.
        ctx.fillStyle = 'rgba(147,255,216,0.035)';
        ctx.fillRect(0, ZONE_Y, G.W, G.H - ZONE_Y);
        ctx.fillStyle = 'rgba(147,255,216,' + (0.16 + 0.1 * Math.sin(s.time * 3)).toFixed(3) + ')';
        for (var d = (s.time * 20) % 24 - 24; d < G.W; d += 24) ctx.fillRect(d, ZONE_Y - 1, 12, 2);
    }

    // A node is a vertebra: a round body with two side processes. It shrinks
    // as it is chipped away, and venom turns it acid green.
    function drawNode(s, ctx, c, r) {
        var hp = s.hp[idx(c, r)], bad = s.poison[idx(c, r)];
        var x = cellX(c), y = cellY(r), rad = 4 + hp * 2;
        if (bad) {
            ctx.fillStyle = 'rgba(147,255,216,' + (0.14 + 0.1 * Math.sin(s.time * 6 + c)).toFixed(3) + ')';
            circle(ctx, x, y, rad + 5); ctx.fill();
        }
        ctx.fillStyle = bad ? ACID : BONE;
        ctx.fillRect(x - rad - 2, y - 2, rad * 2 + 4, 4);
        circle(ctx, x, y, rad); ctx.fill();
        ctx.fillStyle = bad ? '#1d5a48' : '#6f675f';
        circle(ctx, x, y, rad * 0.4); ctx.fill();
    }

    function drawNodes(s, ctx) {
        for (var r = 0; r < ROWS; r++) {
            for (var c = 0; c < COLS; c++) if (s.hp[idx(c, r)]) drawNode(s, ctx, c, r);
        }
    }

    function visiblePoints(w) {
        var pts = [];
        w.segs.forEach(function (g, k) {
            if (!hidden(w, g)) pts.push({ x: segX(w, g), y: segY(w, g), k: k, dive: g.dive });
        });
        return pts;
    }

    // The flesh between the plates, and a pair of legs per segment that
    // paddle out of phase so the crawl reads as a ripple.
    function drawWormTissue(s, ctx, pts) {
        ctx.strokeStyle = FLESH; ctx.lineWidth = 11; ctx.lineJoin = 'round'; ctx.lineCap = 'round';
        ctx.beginPath();
        pts.forEach(function (p, i) { if (i) ctx.lineTo(p.x, p.y); else ctx.moveTo(p.x, p.y); });
        ctx.stroke();
        ctx.strokeStyle = VEIN; ctx.lineWidth = 1.5;
        ctx.beginPath();
        pts.forEach(function (p) {
            var wig = Math.sin(s.time * 22 + p.k * 1.3) * 4;
            ctx.moveTo(p.x - 3, p.y - 7); ctx.lineTo(p.x - 3 + wig, p.y - 13);
            ctx.moveTo(p.x + 3, p.y + 7); ctx.lineTo(p.x + 3 - wig, p.y + 13);
        });
        ctx.stroke();
    }

    function drawWorm(s, ctx, w) {
        var pts = visiblePoints(w);
        if (!pts.length) return;
        drawWormTissue(s, ctx, pts);
        for (var i = pts.length - 1; i >= 0; i--) {
            var p = pts[i], head = p.k === 0;
            ctx.fillStyle = head ? VEIN : FLESH;
            ctx.strokeStyle = p.dive ? ACID : (head ? BONE : VEIN);
            ctx.lineWidth = 2;
            circle(ctx, p.x, p.y, head ? 11 : 9); ctx.fill(); ctx.stroke();
            ctx.fillStyle = head ? ACID : BONE;
            if (head) { ctx.fillRect(p.x - 6, p.y - 3, 4, 5); ctx.fillRect(p.x + 2, p.y - 3, 4, 5); }
            else { circle(ctx, p.x, p.y, 3); ctx.fill(); }
        }
    }

    function drawSpore(s, ctx, b) {
        ctx.strokeStyle = VEIN; ctx.lineWidth = 2;
        ctx.beginPath();
        for (var i = -1; i <= 1; i++) { ctx.moveTo(b.x + i * 4, b.y - 6); ctx.lineTo(b.x + i * 6 + Math.sin(s.time * 14 + i) * 3, b.y - 20); }
        ctx.stroke();
        ctx.fillStyle = b.hp > 1 ? VEIN : FLESH;
        circle(ctx, b.x, b.y, b.r); ctx.fill();
        ctx.fillStyle = ACID;
        circle(ctx, b.x, b.y + 2, 3.5); ctx.fill();
    }

    function drawPara(s, ctx, b) {
        ctx.strokeStyle = BONE; ctx.lineWidth = 2;
        ctx.beginPath();
        for (var i = 0; i < 6; i++) {
            var side = i < 3 ? -1 : 1, a = (i % 3 - 1) * 0.6 + Math.sin(s.time * 26 + i * 2) * 0.35;
            ctx.moveTo(b.x + side * 6, b.y);
            ctx.lineTo(b.x + side * (6 + Math.cos(a) * 13), b.y + Math.sin(a) * 13 + 4);
        }
        ctx.stroke();
        ctx.fillStyle = STEEL; ctx.strokeStyle = VEIN;
        circle(ctx, b.x, b.y, b.r - 2); ctx.fill(); ctx.stroke();
        ctx.fillStyle = ACID;
        ctx.fillRect(b.x - 5, b.y - 3, 3, 3); ctx.fillRect(b.x + 2, b.y - 3, 3, 3);
    }

    function drawEye(s, ctx, b) {
        var back = b.vx > 0 ? -1 : 1;
        ctx.strokeStyle = VEIN; ctx.lineWidth = 2;
        ctx.beginPath();
        for (var i = -1; i <= 1; i++) { ctx.moveTo(b.x + back * 8, b.y + i * 4); ctx.lineTo(b.x + back * (26 + Math.sin(s.time * 30 + i) * 5), b.y + i * 8); }
        ctx.stroke();
        ctx.fillStyle = BONE;
        circle(ctx, b.x, b.y, b.r); ctx.fill();
        ctx.fillStyle = ACID;
        circle(ctx, b.x - back * 3, b.y, 5.5); ctx.fill();
        ctx.fillStyle = BG;
        circle(ctx, b.x - back * 4, b.y, 2.5); ctx.fill();
    }

    function drawBugs(s, ctx) {
        var painters = { spore: drawSpore, para: drawPara, eye: drawEye };
        s.bugs.forEach(function (b) { painters[b.kind](s, ctx, b); });
    }

    function drawQueenHead(s, ctx, q) {
        var hy = headY(q), open = 4 + 3 * Math.sin(s.time * 5);
        ctx.fillStyle = q.hurt > 0 ? '#ffffff' : VEIN;
        ctx.strokeStyle = BONE; ctx.lineWidth = 2;
        circle(ctx, q.x, hy, 16); ctx.fill(); ctx.stroke();
        ctx.fillStyle = BG;
        ctx.beginPath(); ctx.ellipse(q.x, hy + 4, 9, open, 0, 0, Math.PI * 2); ctx.fill();
        ctx.fillStyle = ACID;
        ctx.fillRect(q.x - 10, hy - 9, 5, 4); ctx.fillRect(q.x + 5, hy - 9, 5, 4);
    }

    // An armoured carapace with breathing ribs; the bar above it shows how
    // much of the head is left.
    function drawQueen(s, ctx) {
        var q = s.queen;
        if (!q) return;
        var breathe = 1 + 0.04 * Math.sin(s.time * 2.4);
        ctx.fillStyle = STEEL; ctx.strokeStyle = FLESH; ctx.lineWidth = 4;
        ctx.beginPath(); ctx.ellipse(q.x, q.y, 92 * breathe, 36, 0, 0, Math.PI * 2); ctx.fill(); ctx.stroke();
        ctx.strokeStyle = BONE; ctx.lineWidth = 3;
        ctx.beginPath();
        for (var i = -3; i <= 3; i++) {
            if (!i) continue;
            ctx.moveTo(q.x + i * 22, q.y - 26); ctx.quadraticCurveTo(q.x + i * 30 * breathe, q.y, q.x + i * 20, q.y + 28);
        }
        ctx.stroke();
        drawQueenHead(s, ctx, q);
        ctx.fillStyle = 'rgba(208,199,187,0.25)';
        ctx.fillRect(q.x - 60, OY, 120, 4);
        ctx.fillStyle = VEIN;
        ctx.fillRect(q.x - 60, OY, 120 * q.hp / q.max, 4);
    }

    // Blinks while untouchable; the barrel kicks back on every shot.
    function drawPlayer(s, ctx) {
        var p = s.player;
        if (p.inv > 0 && Math.floor(s.time * 14) % 2) return;
        var kick = s.recoil * 4, glow = 4 + Math.sin(s.time * 8) * 1.5;
        ctx.fillStyle = STEEL; ctx.strokeStyle = BONE; ctx.lineWidth = 2;
        ctx.beginPath();
        ctx.moveTo(p.x - 13, p.y + 10); ctx.lineTo(p.x - 8, p.y - 4); ctx.lineTo(p.x + 8, p.y - 4); ctx.lineTo(p.x + 13, p.y + 10);
        ctx.closePath(); ctx.fill(); ctx.stroke();
        ctx.fillStyle = BONE;
        ctx.fillRect(p.x - 3, p.y - 15 + kick, 6, 13);
        ctx.fillStyle = ACID;
        circle(ctx, p.x, p.y + 3, glow); ctx.fill();
    }

    function drawShots(s, ctx) {
        ctx.fillStyle = 'rgba(147,255,216,0.35)';
        s.shots.forEach(function (sh) { ctx.fillRect(sh.x - 3, sh.y - 2, 6, 18); });
        ctx.fillStyle = ACID;
        s.shots.forEach(function (sh) { ctx.fillRect(sh.x - 1.5, sh.y - 4, 3, 14); });
    }

    function draw(s, ctx) {
        drawBackdrop(s, ctx);
        drawNodes(s, ctx);
        s.worms.forEach(function (w) { drawWorm(s, ctx, w); });
        drawQueen(s, ctx);
        drawBugs(s, ctx);
        drawShots(s, ctx);
        drawPlayer(s, ctx);
    }

    function hud(s) {
        return (s.queen ? 'QUEEN ' + s.queen.hp + '  ' : '') + 'SEGMENTS ' + segmentsLeft(s);
    }

    G.register('biomech', {
        title: 'SPINE CRAWLER',
        blurb: 'Shoot every segment of the worm before it reaches you.',
        controls: [
            '← → ↑ ↓ move the turret in the bottom zone',
            'SPACE (hold): rapid fire',
            'A shot segment turns to bone and splits the worm · green nodes make it dive',
            'Level 10: the Queen only bleeds from her head'
        ],
        levelNames: ['First Vertebra', 'Spore Fall', 'Parasites', 'Twin Spines', 'The Eye',
            'Marrow Bloom', 'Swarm', 'Rib Cage', 'Venom Tide', 'The Queen'],
        colors: { bg: BG, fg: BONE, accent: ACID, dim: '#8c7f86' },
        lives: 3,
        music: {
            bpm: 128, root: 40, scale: 'phrygian', prog: [0, 0, 1, 0, 0, 5, 1, 0],
            bass: 'x.x.x.xox.x.x.o5',
            lead: ['7...8.7.5...4...', '5.4.1...0---....', '7...8.7.a...8...', '7.5.4.1.0---....',
                '4.5.7...8.7.5...', 'b...a.8.7---....', '1.0.1.4.5.4.1...', '0---....7.8.7---'],
            arp: '0.1.2.1.',
            drums: { k: 'x..x..x.x..x..x.', s: '....x.......x..x', h: '..x...x...x...x.' },
            leadWave: 'sawtooth', bassWave: 'sawtooth', arpWave: 'triangle', leadOct: 2
        },
        init: init, update: update, draw: draw, hud: hud
    });
})();
