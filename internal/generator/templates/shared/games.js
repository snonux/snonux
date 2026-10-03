/*
 * games.js — SnoGame, the small arcade engine behind the per-theme games.
 *
 * Every theme ships one game in themes/<name>/game.js. That file is loaded
 * lazily the first time the game is launched and calls
 * SnoGame.register('<name>', def). The engine owns everything the games have
 * in common, so each game only describes its own world:
 *
 *   - a fullscreen overlay with a 960x540 canvas, scaled to fit the viewport
 *   - a fixed 60 Hz update loop (dt is always SnoGame.STEP)
 *   - keyboard / mouse / touch input
 *   - Web Audio sound effects and a step-sequenced music track per game
 *   - particles, floating score popups, screen shake and flashes
 *   - small physics helpers (overlap tests, tile-grid movement)
 *   - title / level select / pause / level clear / game over / victory screens
 *   - ten levels per game, with progress saved in a cookie per game
 *
 * Esc always quits back to the blog. While a game is open, every key is
 * swallowed in the capture phase so shared.js shortcuts do not fire.
 *
 * The authoring guide, with the full API, lives in docs/games.md.
 */
(function () {
    'use strict';
    if (window.SnoGame) return;

    var W = 960, H = 540, STEP = 1 / 60, HUD_H = 30, LEVELS = 10;

    // Theme name → game title. The launch buttons need the title before the
    // game script itself has been fetched, so the list lives here.
    var TITLES = {
        aurora: 'Polar Rush', biomech: 'Spine Crawler', breakout: 'Brick Breaker',
        brutalist: 'Wrecking Ball', cathedral: 'Gargoyle', cosmos: 'Gravity Well',
        dos: 'DIGGER.EXE', invaders: 'Formation', matrix: 'Bullet Time',
        neon: 'Light Cycles', noir: 'Midnight Alley', nukem: 'Nukem Time',
        ocean: 'Deep Channel', pacmaze: 'Neon Maze', pinball: 'Electro Ball',
        plasma: 'Plasma Storm', retro: 'Carrier Lost', retrofuture: 'Atomic Defense',
        spaceage: 'Orbital Dock', surveillance: 'Blind Spot', synthwave: 'Outrun the Sun',
        terminal: 'kill -9', tetris: 'Block Blaster', tropicale: 'Island Hopper',
        volcano: 'Magma Rising'
    };

    var G = window.SnoGame = {
        W: W, H: H, STEP: STEP, HUD: HUD_H, LEVELS: LEVELS, titles: TITLES,
        active: false, level: 1, lives: 3, score: 0, t: 0,
        cam: { x: 0, y: 0 },
        key: { left: false, right: false, up: false, down: false, a: false, b: false },
        hit: {},
        mouse: { x: W / 2, y: H / 2, down: false, hit: false, rdown: false, rhit: false },
        fontFamily: 'monospace'
    };

    var defs = {};          // theme → registered game definition
    var cur = null;         // the running game: { theme, def, s, screen, ... }
    var dom = null;         // overlay nodes: { root, canvas, ctx }
    var scale = 1;          // backing-store pixels per logical pixel
    var raf = 0, lastTime = 0, acc = 0;
    var ambientWasPlaying = false;
    var errors = [];

    function noop() {}

    // ------------------------------------------------------------------
    // Math, random and overlap helpers
    // ------------------------------------------------------------------

    G.clamp = function (v, lo, hi) { return v < lo ? lo : (v > hi ? hi : v); };
    G.lerp = function (a, b, t) { return a + (b - a) * t; };
    G.dist = function (ax, ay, bx, by) { return Math.hypot(bx - ax, by - ay); };
    G.rnd = function (a, b) { return a + Math.random() * (b - a); };
    G.pick = function (arr) { return arr[Math.floor(Math.random() * arr.length)]; };

    // Deterministic generator (mulberry32) so level layouts built from a seed
    // are identical on every visit.
    G.rng = function (seed) {
        var a = seed >>> 0;
        return function () {
            a = (a + 0x6D2B79F5) >>> 0;
            var t = a;
            t = Math.imul(t ^ (t >>> 15), t | 1);
            t ^= t + Math.imul(t ^ (t >>> 7), t | 61);
            return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
        };
    };

    // Rectangles are {x, y, w, h} with x/y at the top-left corner.
    G.aabb = function (a, b) {
        return a.x < b.x + b.w && a.x + a.w > b.x && a.y < b.y + b.h && a.y + a.h > b.y;
    };
    G.circ = function (ax, ay, ar, bx, by, br) {
        var dx = bx - ax, dy = by - ay, r = ar + br;
        return dx * dx + dy * dy < r * r;
    };
    G.circRect = function (cx, cy, r, rx, ry, rw, rh) {
        var nx = G.clamp(cx, rx, rx + rw), ny = G.clamp(cy, ry, ry + rh);
        var dx = cx - nx, dy = cy - ny;
        return dx * dx + dy * dy < r * r;
    };
    // Closest point to (px,py) on the segment a→b; t is the 0..1 position.
    G.closestOnSeg = function (px, py, ax, ay, bx, by) {
        var dx = bx - ax, dy = by - ay, len = dx * dx + dy * dy;
        var t = len ? G.clamp(((px - ax) * dx + (py - ay) * dy) / len, 0, 1) : 0;
        return { x: ax + dx * t, y: ay + dy * t, t: t };
    };

    function hitsTiles(b, ts, solid) {
        var x0 = Math.floor(b.x / ts), x1 = Math.floor((b.x + b.w - 0.001) / ts);
        var y0 = Math.floor(b.y / ts), y1 = Math.floor((b.y + b.h - 0.001) / ts);
        for (var ty = y0; ty <= y1; ty++) {
            for (var tx = x0; tx <= x1; tx++) if (solid(tx, ty)) return true;
        }
        return false;
    }

    // Snaps the body back out of the tile it just entered along one axis and
    // records which side it touched.
    function resolveAxis(b, dx, dy, ts) {
        if (dx > 0) { b.x = Math.floor((b.x + b.w) / ts) * ts - b.w; b.wall = 1; }
        else if (dx < 0) { b.x = (Math.floor(b.x / ts) + 1) * ts; b.wall = -1; }
        if (dy > 0) { b.y = Math.floor((b.y + b.h) / ts) * ts - b.h; b.ground = true; }
        else if (dy < 0) { b.y = (Math.floor(b.y / ts) + 1) * ts; b.ceil = true; }
        if (dx) b.vx = 0;
        if (dy) b.vy = 0;
    }

    function moveAxis(b, dx, dy, ts, solid) {
        // Sub-step so a fast body cannot tunnel through a tile.
        var n = Math.ceil(Math.max(Math.abs(dx), Math.abs(dy)) / (ts * 0.5)) || 1;
        for (var i = 0; i < n; i++) {
            b.x += dx / n; b.y += dy / n;
            if (hitsTiles(b, ts, solid)) { resolveAxis(b, dx, dy, ts); return; }
        }
    }

    // Moves body b ({x, y, w, h, vx, vy}) by its velocity against a tile grid
    // of size ts. solid(tx, ty) says whether a tile blocks. Afterwards
    // b.ground / b.ceil are booleans and b.wall is -1, 0 or 1.
    // b.ground is also true for a body merely resting on a tile (vy = 0), so
    // a game need not push it down every tick to learn that. solid() is
    // called with negative and out-of-range indices and must cope.
    G.tileMove = function (b, dt, ts, solid) {
        b.ground = false; b.ceil = false; b.wall = 0;
        moveAxis(b, b.vx * dt, 0, ts, solid);
        moveAxis(b, 0, b.vy * dt, ts, solid);
        if (!b.ground && b.vy >= 0) {
            b.ground = hitsTiles({ x: b.x, y: b.y + b.h, w: b.w, h: 1 }, ts, solid);
        }
    };

    // ------------------------------------------------------------------
    // Saved progress: one small cookie per game
    // ------------------------------------------------------------------

    // document.cookie itself can throw (sandboxed frames), so every access
    // is guarded: a game must still run when nothing can be saved.
    function readCookie(name) {
        try {
            var m = document.cookie.match(new RegExp('(?:^|; )' + name + '=([^;]*)'));
            return m ? decodeURIComponent(m[1]) : null;
        } catch (_) { return null; }
    }

    function writeCookie(name, value) {
        document.cookie = name + '=' + encodeURIComponent(value) +
            '; max-age=31536000; path=/; SameSite=Lax';
    }

    // Cookies are unavailable on file:// pages; localStorage keeps local
    // previews working there. The cookie is the store everywhere else.
    function storeGet(name) {
        var v = readCookie(name);
        if (v !== null) return v;
        try { return localStorage.getItem(name); } catch (_) { return null; }
    }

    function storeSet(name, value) {
        try { writeCookie(name, value); } catch (_) {}
        if (readCookie(name) === value) return;
        try { localStorage.setItem(name, value); } catch (_) {}
    }

    var MAX_SCORE = 999999999;

    // Only plain digit strings count; "4abc" or "-5" are not numbers here.
    function intIn(v, lo, hi, fallback) {
        if (!/^\d{1,9}$/.test(String(v))) return fallback;
        var n = parseInt(v, 10);
        return (n < lo || n > hi) ? fallback : n;
    }

    // Save format: "unlocked.hiscore.lastLevel.won". Anything malformed or
    // out of range falls back to a fresh save rather than being trusted.
    function parseSave(raw) {
        var p = String(raw || '').split('.');
        var save = {
            unlocked: intIn(p[0], 1, LEVELS, 1),
            hi: intIn(p[1], 0, MAX_SCORE, 0),
            last: intIn(p[2], 1, LEVELS, 1),
            won: intIn(p[3], 0, 1, 0)
        };
        if (save.last > save.unlocked) save.last = save.unlocked;
        return save;
    }

    function loadSave(theme) { return parseSave(storeGet('snog_' + theme)); }

    function writeSave() {
        if (!cur) return;
        var s = cur.save;
        if (G.score > s.hi) s.hi = Math.min(MAX_SCORE, Math.floor(G.score));
        storeSet('snog_' + cur.theme, [s.unlocked, s.hi, s.last, s.won].join('.'));
    }

    // ------------------------------------------------------------------
    // Audio: sound effects
    // ------------------------------------------------------------------

    var actx = null, master = null, musicBus = null, sfxBus = null, noiseBuf = null;
    var muted = storeGet('snog_mute') === '1';
    // Counters for the e2e test: `sfx` counts sounds a game asked for itself,
    // as opposed to notes the music sequencer scheduled.
    var audioStats = { tones: 0, noises: 0, musicSteps: 0, sfx: 0, engineSfx: 0 };
    var audioHeld = false;   // true while paused, hidden or quit: nothing may resume the context
    var engineSound = false; // true while the engine itself plays a jingle

    function audio() {
        if (!actx) {
            try {
                actx = new (window.AudioContext || window.webkitAudioContext)();
                master = actx.createGain();
                master.gain.value = muted ? 0 : 1;
                master.connect(actx.destination);
                musicBus = actx.createGain(); musicBus.gain.value = 0.5; musicBus.connect(master);
                sfxBus = actx.createGain(); sfxBus.gain.value = 0.9; sfxBus.connect(master);
            } catch (_) { actx = null; return null; }
        }
        if (actx.state === 'suspended' && !audioHeld) actx.resume().catch(noop);
        return actx;
    }

    // Freezes or releases all sound. Pausing, a hidden tab and quitting hold
    // the context so neither the music pump nor a key press can restart it.
    function holdAudio(on) {
        audioHeld = on;
        if (!actx) return;
        if (on) actx.suspend().catch(noop); else actx.resume().catch(noop);
    }

    function countSfx(o) {
        if (o.bus) return;
        if (engineSound) audioStats.engineSfx++; else audioStats.sfx++;
    }

    // Engine jingles (hurt, win, over, menu blips) are kept out of the game's
    // own sfx count so the e2e test can tell a silent game from a noisy engine.
    function engineSfx(name) {
        engineSound = true;
        G.sfx(name);
        engineSound = false;
    }

    // Shared attack/decay envelope for both oscillators and noise.
    function envelope(c, t, dur, vol, attack) {
        var g = c.createGain();
        g.gain.setValueAtTime(0.0001, t);
        g.gain.linearRampToValueAtTime(vol, t + (attack || 0.005));
        g.gain.exponentialRampToValueAtTime(0.0001, t + dur);
        return g;
    }

    // G.tone(freq, dur, {type, vol, slide, delay, attack}) plays one
    // oscillator note; `slide` glides the pitch to that frequency.
    G.tone = function (freq, dur, o) {
        var c = audio();
        if (!c) return;
        o = o || {};
        audioStats.tones++;
        countSfx(o);
        var t = c.currentTime + (o.delay || 0);
        var osc = c.createOscillator();
        var g = envelope(c, t, dur, o.vol == null ? 0.2 : o.vol, o.attack);
        osc.type = o.type || 'square';
        osc.frequency.setValueAtTime(Math.max(1, freq), t);
        if (o.slide) osc.frequency.exponentialRampToValueAtTime(Math.max(1, o.slide), t + dur);
        osc.connect(g); g.connect(o.bus || sfxBus);
        osc.start(t); osc.stop(t + dur + 0.03);
    };

    function noiseBuffer(c) {
        if (noiseBuf) return noiseBuf;
        noiseBuf = c.createBuffer(1, c.sampleRate, c.sampleRate);
        var d = noiseBuf.getChannelData(0);
        for (var i = 0; i < d.length; i++) d[i] = Math.random() * 2 - 1;
        return noiseBuf;
    }

    // G.noise(dur, {vol, freq, slide, filter, q, delay, attack}) plays filtered
    // white noise: explosions, hats, wind, engines.
    G.noise = function (dur, o) {
        var c = audio();
        if (!c) return;
        o = o || {};
        audioStats.noises++;
        countSfx(o);
        var t = c.currentTime + (o.delay || 0);
        var src = c.createBufferSource(), f = c.createBiquadFilter();
        var g = envelope(c, t, dur, o.vol == null ? 0.25 : o.vol, o.attack);
        src.buffer = noiseBuffer(c); src.loop = true;
        f.type = o.filter || 'lowpass';
        f.Q.value = o.q || 1;
        f.frequency.setValueAtTime(o.freq || 1200, t);
        if (o.slide) f.frequency.exponentialRampToValueAtTime(Math.max(20, o.slide), t + dur);
        src.connect(f); f.connect(g); g.connect(o.bus || sfxBus);
        src.start(t); src.stop(t + dur + 0.03);
    };

    // Ready-made effects every game can use; games add their own flavour by
    // calling G.tone / G.noise directly.
    var SFX = {
        shoot: function () { G.tone(880, 0.09, { slide: 220, vol: 0.12 }); },
        laser: function () { G.tone(1400, 0.16, { type: 'sawtooth', slide: 180, vol: 0.1 }); },
        boom: function () { G.noise(0.35, { freq: 900, slide: 60, vol: 0.35 }); },
        bigboom: function () {
            G.noise(0.8, { freq: 1400, slide: 40, vol: 0.45 });
            G.tone(90, 0.6, { type: 'sine', slide: 30, vol: 0.4 });
        },
        jump: function () { G.tone(260, 0.16, { slide: 620, vol: 0.14 }); },
        coin: function () {
            G.tone(988, 0.07, { vol: 0.13 });
            G.tone(1319, 0.2, { vol: 0.13, delay: 0.07 });
        },
        power: function () {
            [523, 659, 784, 1047].forEach(function (f, i) {
                G.tone(f, 0.12, { type: 'triangle', vol: 0.16, delay: i * 0.06 });
            });
        },
        hit: function () { G.tone(200, 0.08, { slide: 90, vol: 0.18 }); },
        hurt: function () {
            G.tone(320, 0.3, { type: 'sawtooth', slide: 60, vol: 0.22 });
            G.noise(0.2, { freq: 600, vol: 0.2 });
        },
        bounce: function () { G.tone(440, 0.05, { type: 'triangle', vol: 0.15 }); },
        blip: function () { G.tone(660, 0.04, { vol: 0.1 }); },
        click: function () { G.tone(1200, 0.02, { vol: 0.08 }); },
        alarm: function () {
            G.tone(700, 0.18, { vol: 0.14 });
            G.tone(520, 0.18, { vol: 0.14, delay: 0.18 });
        },
        splash: function () { G.noise(0.3, { filter: 'bandpass', freq: 500, slide: 2500, vol: 0.2 }); },
        thrust: function () { G.noise(0.12, { freq: 400, vol: 0.12 }); },
        win: function () {
            [523, 659, 784, 1047, 1319].forEach(function (f, i) {
                G.tone(f, 0.22, { type: 'triangle', vol: 0.18, delay: i * 0.11 });
            });
        },
        over: function () {
            [392, 330, 262, 196].forEach(function (f, i) {
                G.tone(f, 0.32, { type: 'sawtooth', vol: 0.14, delay: i * 0.2 });
            });
        }
    };

    G.sfx = function (name) { if (SFX[name]) SFX[name](); };

    function setMuted(on) {
        muted = !!on;
        storeSet('snog_mute', muted ? '1' : '0');
        if (master && actx) master.gain.setTargetAtTime(muted ? 0 : 1, actx.currentTime, 0.02);
    }

    // ------------------------------------------------------------------
    // Audio: step-sequenced music
    // ------------------------------------------------------------------
    //
    // A game describes its tune in def.music; see docs/games.md for the
    // format. Sixteen steps make one bar, `prog` holds one chord root (a
    // scale degree) per bar, and the pattern strings loop independently.

    var SCALES = {
        major: [0, 2, 4, 5, 7, 9, 11], minor: [0, 2, 3, 5, 7, 8, 10],
        dorian: [0, 2, 3, 5, 7, 9, 10], phrygian: [0, 1, 3, 5, 7, 8, 10],
        lydian: [0, 2, 4, 6, 7, 9, 11], mixolydian: [0, 2, 4, 5, 7, 9, 10],
        harmonic: [0, 2, 3, 5, 7, 8, 11], pentatonic: [0, 3, 5, 7, 10],
        majorpenta: [0, 2, 4, 7, 9], blues: [0, 3, 5, 6, 7, 10],
        whole: [0, 2, 4, 6, 8, 10]
    };
    var music = { def: null, step: 0, next: 0, timer: 0, spb: 0.12 };

    function midiHz(n) { return 440 * Math.pow(2, (n - 69) / 12); }

    // semis shifts the result by plain semitones (used for a true fifth).
    function scaleHz(m, degree, octave, semis) {
        var sc = SCALES[m.scale] || SCALES.minor, n = sc.length;
        var oct = Math.floor(degree / n), idx = ((degree % n) + n) % n;
        return midiHz((m.root || 45) + 12 * (oct + octave) + sc[idx] + (semis || 0));
    }

    function scaleLen(m) { return (SCALES[m.scale] || SCALES.minor).length; }

    function patAt(p, i) {
        if (!p) return '.';
        if (Array.isArray(p)) p = p.join('');
        return p.length ? p.charAt(i % p.length) : '.';
    }

    // A note is held for as many steps as there are '-' characters after it.
    function heldSteps(p, i) {
        if (Array.isArray(p)) p = p.join('');
        var n = 1;
        while (n < 16 && p.charAt((i + n) % p.length) === '-') n++;
        return n;
    }

    function playBass(m, i, t, root) {
        var ch = patAt(m.bass, i);
        if (ch === '.' || ch === '-') return;
        // 'x' plays the chord root, 'o' the octave above it, '5' a perfect
        // fifth above it (in semitones, so it is right in every scale).
        G.tone(scaleHz(m, root, ch === 'o' ? 1 : 0, ch === '5' ? 7 : 0), music.spb * 1.8, {
            type: m.bassWave || 'triangle', vol: 0.3, bus: musicBus, delay: t
        });
    }

    function playLead(m, i, t) {
        var ch = patAt(m.lead, i);
        if (ch === '.' || ch === '-') return;
        var deg = parseInt(ch, 36);
        if (isNaN(deg)) return;
        G.tone(scaleHz(m, deg, m.leadOct == null ? 2 : m.leadOct), music.spb * heldSteps(m.lead, i) * 0.95, {
            type: m.leadWave || 'square', vol: 0.13, bus: musicBus, delay: t
        });
    }

    function playArp(m, i, t, root) {
        var ch = patAt(m.arp, i);
        if (ch === '.' || ch === '-') return;
        // Chord tones are stacked thirds in seven-note scales; five- and
        // six-note scales use neighbouring scale tones, which stay consonant.
        var n = scaleLen(m), steps = n === 7 ? [0, 2, 4, 7] : [0, 1, 2, n];
        var tone = steps[parseInt(ch, 10) % 4] || 0;
        G.tone(scaleHz(m, root + tone, 1), music.spb * 0.9, {
            type: m.arpWave || 'square', vol: 0.06, bus: musicBus, delay: t
        });
    }

    function playDrums(m, i, t) {
        var d = m.drums || {};
        if (patAt(d.k, i) === 'x') G.tone(150, 0.14, { type: 'sine', slide: 40, vol: 0.5, bus: musicBus, delay: t });
        if (patAt(d.s, i) === 'x') G.noise(0.12, { filter: 'bandpass', freq: 1800, vol: 0.22, bus: musicBus, delay: t });
        if (patAt(d.h, i) === 'x') G.noise(0.035, { filter: 'highpass', freq: 7000, vol: 0.09, bus: musicBus, delay: t });
    }

    function musicStep(i, when) {
        var m = music.def, t = Math.max(0, when - actx.currentTime);
        var prog = m.prog && m.prog.length ? m.prog : [0];
        var root = prog[Math.floor(i / 16) % prog.length];
        audioStats.musicSteps++;
        playBass(m, i, t, root);
        playLead(m, i, t);
        playArp(m, i, t, root);
        playDrums(m, i, t);
    }

    // Schedules a little ahead of the audio clock so timing stays tight even
    // when the main thread is busy drawing.
    function musicPump() {
        if (!music.def || !actx || audioHeld) return;
        // After a stall (busy page, throttled tab) drop the missed steps
        // rather than playing them all at once.
        if (music.next < actx.currentTime) {
            var missed = Math.ceil((actx.currentTime - music.next) / music.spb);
            music.step += missed; music.next += missed * music.spb;
        }
        while (music.next < actx.currentTime + 0.15) {
            musicStep(music.step, music.next);
            music.step++;
            music.next += music.spb;
        }
    }

    function musicStop() {
        if (music.timer) clearInterval(music.timer);
        music.timer = 0; music.def = null;
    }

    // Later levels play the same tune slightly faster to raise the pressure.
    function musicStart(def, level) {
        musicStop();
        var c = audio();
        if (!c || !def) return;
        music.def = def; music.step = 0;
        music.spb = 60 / ((def.bpm || 120) * (1 + 0.02 * (level - 1))) / 4;
        music.next = c.currentTime + 0.08;
        music.timer = setInterval(musicPump, 30);
    }

    // ------------------------------------------------------------------
    // Input
    // ------------------------------------------------------------------

    // Physical key (KeyboardEvent.code) → logical button. Games that need the
    // two halves separately (twin-stick) read G.down(code) instead.
    var LOGICAL = {
        ArrowLeft: 'left', KeyA: 'left', ArrowRight: 'right', KeyD: 'right',
        ArrowUp: 'up', KeyW: 'up', ArrowDown: 'down', KeyS: 'down',
        Space: 'a', KeyZ: 'a', KeyJ: 'a',
        KeyX: 'b', KeyK: 'b', ShiftLeft: 'b', ShiftRight: 'b'
    };
    var rawDown = {}, rawHit = {};

    G.down = function (code) { return !!rawDown[code]; };
    G.pressed = function (code) { return !!rawHit[code]; };

    function refreshLogical() {
        var k = G.key, name;
        for (name in k) k[name] = false;
        for (var code in LOGICAL) if (rawDown[code]) k[LOGICAL[code]] = true;
    }

    function setRaw(code, on) {
        if (on && !rawDown[code]) {
            rawHit[code] = true;
            if (LOGICAL[code]) G.hit[LOGICAL[code]] = true;
        }
        rawDown[code] = on;
        refreshLogical();
    }

    function clearInput() {
        rawDown = {}; rawHit = {}; G.hit = {};
        G.mouse.down = G.mouse.hit = G.mouse.rdown = G.mouse.rhit = false;
        refreshLogical();
    }

    // Edge-triggered state lives for exactly one update tick.
    function endTick() {
        rawHit = {}; G.hit = {};
        G.mouse.hit = false; G.mouse.rhit = false;
    }

    function isBrowserShortcut(e) {
        return e.ctrlKey || e.metaKey || e.altKey || /^F\d+$/.test(e.key);
    }

    // Keys still held when the game quit. Their auto-repeat would otherwise
    // reach the blog (a held arrow key would turn the page).
    var lingering = {};

    function onKeyDown(e) {
        if (!G.active) {
            if (lingering[e.code] && e.repeat) { e.preventDefault(); e.stopImmediatePropagation(); }
            else delete lingering[e.code];
            return;
        }
        // The blog must never see a key while a game is open, modifiers or
        // not; real browser shortcuts still keep their default action.
        e.stopImmediatePropagation();
        if (isBrowserShortcut(e)) return;
        e.preventDefault();
        if (e.code === 'Escape') { G.quit(); return; }
        if (e.repeat || !cur) return;
        audio();
        if (e.code === 'KeyM') { setMuted(!muted); return; }
        if (e.code === 'KeyP' && (cur.screen === 'play' || cur.screen === 'paused')) { togglePause(); return; }
        setRaw(e.code, true);
        if (cur.screen !== 'play') menuKey(e.code);
    }

    function onKeyUp(e) {
        delete lingering[e.code];
        if (!G.active) return;
        // preventDefault too: some browsers activate a focused button on the
        // Space keyup alone.
        e.preventDefault(); e.stopImmediatePropagation();
        setRaw(e.code, false);
    }

    function pointerPos(e) {
        var r = dom.canvas.getBoundingClientRect();
        G.mouse.x = G.clamp((e.clientX - r.left) / r.width * W, 0, W);
        G.mouse.y = G.clamp((e.clientY - r.top) / r.height * H, 0, H);
    }

    // Button state comes from e.buttons, not from which event fired: with two
    // buttons held, the second press and the first release arrive as
    // pointermove, so tracking down/up events alone leaves buttons stuck.
    function syncButtons(e) {
        var left = !!(e.buttons & 1), right = !!(e.buttons & 2), m = G.mouse;
        var leftHit = left && !m.down;
        if (leftHit) m.hit = true;
        if (right && !m.rdown) m.rhit = true;
        m.down = left; m.rdown = right;
        return leftHit;
    }

    function onPointerDown(e) {
        e.preventDefault();
        audio();
        pointerPos(e);
        if (syncButtons(e) && cur && cur.screen !== 'play') menuClick();
    }

    function onPointerMove(e) { pointerPos(e); syncButtons(e); }

    function onPointerUp(e) { syncButtons(e); }

    // On-screen pad for touch devices; each button simply holds a key.
    var PAD = [
        ['ArrowLeft', '◀', 'l'], ['ArrowRight', '▶', 'r'], ['ArrowUp', '▲', 'u'],
        ['ArrowDown', '▼', 'd'], ['Space', 'A', 'a'], ['KeyX', 'B', 'b']
    ];

    function buildTouchPad(root) {
        if (!window.matchMedia || !window.matchMedia('(pointer: coarse)').matches) return;
        var pad = document.createElement('div');
        pad.className = 'sno-game-pad';
        PAD.forEach(function (p) {
            var b = document.createElement('button');
            b.type = 'button'; b.className = 'sno-game-pad-' + p[2]; b.textContent = p[1];
            b.addEventListener('pointerdown', function (e) {
                e.preventDefault(); audio(); setRaw(p[0], true);
                // In menus the pad works like the keyboard: A confirms,
                // left/right pick the level.
                if (cur && cur.screen !== 'play') menuKey(p[0]);
            });
            ['pointerup', 'pointercancel', 'pointerleave'].forEach(function (ev) {
                b.addEventListener(ev, function () { setRaw(p[0], false); });
            });
            pad.appendChild(b);
        });
        root.appendChild(pad);
    }

    // ------------------------------------------------------------------
    // Effects: particles, popups, shake, flash
    // ------------------------------------------------------------------

    var particles = [], popups = [], shake = { mag: 0, t: 0 }, flash = { color: '#fff', t: 0, max: 0 };
    var MAX_PARTICLES = 600;

    // G.burst(x, y, {n, color, speed, life, size, gravity, angle, spread, drag})
    // sprays particles in world coordinates.
    G.burst = function (x, y, o) {
        o = o || {};
        var n = o.n || 12, spread = o.spread == null ? Math.PI * 2 : o.spread, base = o.angle || 0;
        for (var i = 0; i < n && particles.length < MAX_PARTICLES; i++) {
            var a = base + (Math.random() - 0.5) * spread, sp = (o.speed || 160) * (0.3 + Math.random() * 0.7);
            var life = (o.life || 0.6) * (0.5 + Math.random() * 0.5);
            particles.push({
                x: x, y: y, vx: Math.cos(a) * sp, vy: Math.sin(a) * sp, life: life, max: life,
                size: o.size || 3, color: o.color || '#fff', g: o.gravity || 0, drag: o.drag || 0
            });
        }
    };

    G.popup = function (x, y, text, color) {
        popups.push({ x: x, y: y, text: String(text), color: color || '#fff', life: 0.9 });
    };
    G.shake = function (mag, time) { shake.mag = Math.max(shake.mag, mag || 6); shake.t = Math.max(shake.t, time || 0.25); };
    G.flash = function (color, time) { flash.color = color || '#fff'; flash.t = flash.max = time || 0.15; };

    function clearFx() { particles.length = 0; popups.length = 0; shake.t = 0; flash.t = 0; }

    function updateFx(dt) {
        var i, p;
        for (i = particles.length - 1; i >= 0; i--) {
            p = particles[i];
            p.life -= dt;
            if (p.life <= 0) { particles.splice(i, 1); continue; }
            p.vy += p.g * dt;
            if (p.drag) { p.vx -= p.vx * p.drag * dt; p.vy -= p.vy * p.drag * dt; }
            p.x += p.vx * dt; p.y += p.vy * dt;
        }
        for (i = popups.length - 1; i >= 0; i--) {
            popups[i].life -= dt; popups[i].y -= 40 * dt;
            if (popups[i].life <= 0) popups.splice(i, 1);
        }
        if (shake.t > 0) shake.t -= dt;
        if (flash.t > 0) flash.t -= dt;
    }

    function drawFx(ctx) {
        var i, p;
        ctx.save();
        ctx.translate(-G.cam.x, -G.cam.y);
        for (i = 0; i < particles.length; i++) {
            p = particles[i];
            ctx.globalAlpha = Math.max(0, p.life / p.max);
            ctx.fillStyle = p.color;
            ctx.fillRect(p.x - p.size / 2, p.y - p.size / 2, p.size, p.size);
        }
        for (i = 0; i < popups.length; i++) {
            p = popups[i];
            ctx.globalAlpha = Math.min(1, p.life * 2);
            G.text(p.text, p.x, p.y, { size: 16, color: p.color, align: 'center', bold: true });
        }
        ctx.restore();
    }

    // ------------------------------------------------------------------
    // Text
    // ------------------------------------------------------------------

    G.font = function (px, bold) { return (bold ? '700 ' : '') + px + 'px ' + G.fontFamily; };

    // G.text(str, x, y, {size, color, align, bold, glow, max, baseline})
    G.text = function (str, x, y, o) {
        var ctx = dom && dom.ctx;
        if (!ctx) return;
        o = o || {};
        ctx.font = G.font(o.size || 16, o.bold);
        ctx.textAlign = o.align || 'left';
        ctx.textBaseline = o.baseline || 'alphabetic';
        ctx.fillStyle = o.color || '#fff';
        if (o.glow) { ctx.shadowColor = o.glow; ctx.shadowBlur = 14; }
        if (o.max) ctx.fillText(str, x, y, o.max); else ctx.fillText(str, x, y);
        ctx.shadowBlur = 0;
    };

    // ------------------------------------------------------------------
    // Game flow
    // ------------------------------------------------------------------

    function fail(err) {
        errors.push(String(err && err.stack || err));
        if (window.console) console.error('SnoGame:', err);
        if (cur) { cur.screen = 'error'; cur.timer = 0; }
        musicStop();
    }

    // Runs game-supplied code; one broken game must not take the page down.
    function guard(fn) {
        try { return fn(); } catch (e) { fail(e); return null; }
    }

    function setScreen(name, timer) {
        cur.screen = name; cur.timer = timer || 0; cur.age = 0;
        // The key or click that confirmed a menu must not also count as the
        // first press of the level (it would serve the ball, fire, bomb …).
        if (name === 'play') endTick();
    }

    function maxLives() { return cur.def.lives == null ? 3 : cur.def.lives; }

    function initState(level) {
        G.level = level; G.t = 0; G.cam.x = 0; G.cam.y = 0;
        cur.s = guard(function () { return cur.def.init(level, G); });
    }

    function beginLevel(level) {
        clearFx();
        cur.level = level;
        G.lives = maxLives();
        initState(level);
        if (cur.screen === 'error') return;
        cur.save.last = level;
        writeSave();
        musicStart(cur.def.music, level);
        setScreen('intro', 1.6);
        cur.banner = 'LEVEL ' + level;
    }

    function toTitle() {
        clearFx();
        cur.sel = cur.save.last;
        G.score = 0;
        G.lives = maxLives();
        initState(cur.sel);
        if (cur.screen !== 'error') setScreen('title');
    }

    // G.win(bonus) — the level's goal has been reached.
    G.win = function (bonus) {
        if (!cur || cur.screen !== 'play') return;
        G.score += bonus || 0;
        cur.bonus = bonus || 0;
        if (cur.level >= LEVELS) cur.save.won = 1;
        else if (cur.level + 1 > cur.save.unlocked) cur.save.unlocked = cur.level + 1;
        writeSave();
        engineSfx('win');
        setScreen(cur.level >= LEVELS ? 'victory' : 'clear', cur.level >= LEVELS ? 0 : 6);
    };

    function gameOver() {
        writeSave();
        musicStop();
        engineSfx('over');
        setScreen('over');
    }

    // G.loseLife() — costs one life but the game carries on from where it is
    // (the game handles its own respawn). Returns the lives left.
    G.loseLife = function () {
        if (!cur || cur.screen !== 'play') return G.lives;
        G.lives = Math.floor(G.lives) - 1;
        engineSfx('hurt'); G.shake(8, 0.3);
        if (G.lives <= 0) gameOver();
        return G.lives;
    };

    // G.die() — costs one life and restarts the current level from init().
    G.die = function () {
        if (!cur || cur.screen !== 'play') return;
        if (G.loseLife() > 0) setScreen('dead', 1.0);
    };

    // A bad value (undefined, NaN) is ignored so it cannot poison the score.
    G.addScore = function (n) { if (isFinite(n)) G.score += n; };

    // G.addLife(max) — one extra life, never above max (default 5).
    G.addLife = function (max) { G.lives = Math.min(max || 5, Math.floor(G.lives) + 1); return G.lives; };

    function respawn() {
        var lives = G.lives;
        initState(cur.level);
        G.lives = lives;
        if (cur.screen === 'error') return;
        cur.banner = 'READY';
        setScreen('intro', 0.9);
    }

    function onTimer() {
        if (cur.screen === 'intro') setScreen('play');
        else if (cur.screen === 'dead') respawn();
        else if (cur.screen === 'clear') beginLevel(cur.level + 1);
    }

    function togglePause() {
        if (cur.screen === 'play') { cur.screen = 'paused'; holdAudio(true); }
        else if (cur.screen === 'paused') { cur.screen = 'play'; holdAudio(false); endTick(); }
    }

    function selectLevel(n) {
        n = G.clamp(n, 1, cur.save.unlocked);
        if (n === cur.sel) return;
        cur.sel = n;
        engineSfx('blip');
        initState(n);
    }

    function menuConfirm() {
        var sc = cur.screen;
        if (sc === 'title') { G.score = 0; beginLevel(cur.sel); }
        else if (sc === 'intro') setScreen('play');
        else if (sc === 'paused') togglePause();
        else if (sc === 'error') G.quit();
        else if (cur.age < 0.5) return;
        else if (sc === 'clear') beginLevel(cur.level + 1);
        else if (sc === 'over') { G.score = 0; beginLevel(cur.level); }
        else if (sc === 'victory') toTitle();
    }

    function menuKey(code) {
        if (code === 'Enter' || code === 'Space' || code === 'NumpadEnter') { menuConfirm(); return; }
        if (cur.screen !== 'title') return;
        if (LOGICAL[code] === 'left') selectLevel(cur.sel - 1);
        else if (LOGICAL[code] === 'right') selectLevel(cur.sel + 1);
    }

    // Level boxes on the title screen: 10 squares in one centred row.
    var BOX = { size: 56, gap: 12, y: 340 };
    function boxX(i) { return (W - (LEVELS * BOX.size + (LEVELS - 1) * BOX.gap)) / 2 + i * (BOX.size + BOX.gap); }

    function menuClick() {
        if (cur.screen !== 'title') { menuConfirm(); return; }
        for (var i = 0; i < LEVELS; i++) {
            var hit = G.mouse.x >= boxX(i) && G.mouse.x <= boxX(i) + BOX.size &&
                G.mouse.y >= BOX.y && G.mouse.y <= BOX.y + BOX.size;
            if (!hit) continue;
            if (i + 1 > cur.save.unlocked) return;
            if (i + 1 === cur.sel) menuConfirm(); else selectLevel(i + 1);
            return;
        }
        menuConfirm();
    }

    function tick() {
        if (!cur || !cur.def) return;
        if (cur.screen === 'paused') { endTick(); return; }
        cur.age += STEP;
        if (cur.screen === 'play') {
            G.t += STEP;
            guard(function () { cur.def.update(cur.s, STEP, G); });
        } else if (cur.timer > 0) {
            cur.timer -= STEP;
            if (cur.timer <= 0) onTimer();
        }
        updateFx(STEP);
        endTick();
    }

    // ------------------------------------------------------------------
    // Rendering: HUD and menu screens
    // ------------------------------------------------------------------

    function colors() {
        var c = (cur && cur.def && cur.def.colors) || {};
        return { bg: c.bg || '#05060a', fg: c.fg || '#e8eef5', accent: c.accent || '#4cc9f0', dim: c.dim || '#7a8696' };
    }

    function levelName(n) {
        var names = cur.def.levelNames;
        return names && names[n - 1] ? names[n - 1] : '';
    }

    function pad6(n) { return ('000000' + Math.floor(n)).slice(-6); }

    function drawHud(ctx) {
        var c = colors(), name = levelName(cur.level);
        ctx.fillStyle = 'rgba(0,0,0,0.55)';
        ctx.fillRect(0, 0, W, HUD_H);
        G.text('LV ' + cur.level + '/' + LEVELS + (name ? '  ' + name : ''), 12, 21, { size: 15, color: c.fg, max: 300 });
        G.text(pad6(G.score), W / 2, 21, { size: 16, color: c.accent, align: 'center', bold: true });
        var extra = cur.def.hud ? guard(function () { return cur.def.hud(cur.s, G); }) : '';
        var lives = maxLives() > 0 ? new Array(Math.max(0, Math.floor(G.lives) || 0) + 1).join('♥') : '';
        // The right edge stays clear of the overlay's ESC button.
        G.text((extra ? extra + '   ' : '') + lives, W - 84, 21, { size: 15, color: c.fg, align: 'right', max: 380 });
    }

    function dimScreen(ctx, alpha) {
        ctx.fillStyle = 'rgba(0,0,0,' + alpha + ')';
        ctx.fillRect(0, 0, W, H);
    }

    function drawLevelBoxes(ctx, c) {
        for (var i = 0; i < LEVELS; i++) {
            var x = boxX(i), open = i + 1 <= cur.save.unlocked, sel = i + 1 === cur.sel;
            ctx.lineWidth = 2;
            ctx.strokeStyle = open ? c.accent : c.dim;
            ctx.fillStyle = sel ? c.accent : 'rgba(0,0,0,0.5)';
            ctx.fillRect(x, BOX.y, BOX.size, BOX.size);
            ctx.strokeRect(x, BOX.y, BOX.size, BOX.size);
            G.text(String(i + 1), x + BOX.size / 2, BOX.y + 38, {
                size: 26, bold: true, align: 'center', color: sel ? c.bg : (open ? c.fg : c.dim)
            });
        }
    }

    function drawTitle(ctx) {
        var c = colors(), d = cur.def, lines = d.controls || [], name = levelName(cur.sel);
        dimScreen(ctx, 0.74);
        G.text(d.title || TITLES[cur.theme], W / 2, 112, { size: 60, bold: true, color: c.accent, align: 'center', glow: c.accent, max: W - 60 });
        G.text(d.blurb || '', W / 2, 154, { size: 18, color: c.fg, align: 'center', max: W - 80 });
        for (var i = 0; i < lines.length; i++) {
            G.text(lines[i], W / 2, 200 + i * 24, { size: 15, color: c.dim, align: 'center', max: W - 80 });
        }
        drawLevelBoxes(ctx, c);
        G.text('LEVEL ' + cur.sel + (name ? ' — ' + name : ''), W / 2, 430, { size: 18, color: c.fg, align: 'center', max: W - 80 });
        G.text('ENTER start · ← → level · P pause · M ' + (muted ? 'unmute' : 'mute') + ' · ESC quit',
            W / 2, 478, { size: 14, color: c.dim, align: 'center' });
        G.text('HI ' + pad6(cur.save.hi) + (cur.save.won ? '  ★ COMPLETED' : ''), W / 2, 508, { size: 14, color: c.accent, align: 'center' });
    }

    function drawBanner(ctx, big, small, tint) {
        var c = colors();
        ctx.fillStyle = 'rgba(0,0,0,0.66)';
        ctx.fillRect(0, H / 2 - 74, W, 148);
        G.text(big, W / 2, H / 2 + 4, { size: 46, bold: true, color: tint || c.accent, align: 'center', glow: tint || c.accent, max: W - 60 });
        if (small) G.text(small, W / 2, H / 2 + 44, { size: 17, color: c.fg, align: 'center', max: W - 60 });
    }

    function drawScreen(ctx) {
        var sc = cur.screen;
        if (sc === 'title') drawTitle(ctx);
        else if (sc === 'intro') drawBanner(ctx, cur.banner, levelName(cur.level));
        else if (sc === 'paused') drawBanner(ctx, 'PAUSED', 'P resume · ESC quit');
        else if (sc === 'clear') drawBanner(ctx, 'LEVEL ' + cur.level + ' CLEAR', (cur.bonus ? 'BONUS ' + cur.bonus + ' · ' : '') + 'ENTER for level ' + (cur.level + 1));
        else if (sc === 'over') drawBanner(ctx, 'GAME OVER', 'SCORE ' + pad6(G.score) + ' · ENTER retry level ' + cur.level + ' · ESC quit', '#ff5566');
        else if (sc === 'victory') drawBanner(ctx, 'YOU WIN', 'All ' + LEVELS + ' levels cleared · SCORE ' + pad6(G.score) + ' · ENTER');
        else if (sc === 'error') drawBanner(ctx, 'GAME CRASHED', 'ESC to return to the blog', '#ff5566');
        else if (sc === 'loading') drawBanner(ctx, 'LOADING', TITLES[cur.theme] || '');
    }

    function drawWorld(ctx) {
        var sx = 0, sy = 0;
        if (shake.t > 0) { sx = (Math.random() - 0.5) * shake.mag * 2; sy = (Math.random() - 0.5) * shake.mag * 2; }
        ctx.save();
        ctx.translate(sx, sy);
        if (cur.s && cur.def && cur.screen !== 'error') {
            ctx.save();
            guard(function () { cur.def.draw(cur.s, ctx, G); });
            ctx.restore();
            ctx.globalAlpha = 1; ctx.shadowBlur = 0;
            drawFx(ctx);
        }
        ctx.restore();
    }

    function frame() {
        var ctx = dom.ctx, c = colors();
        ctx.setTransform(scale, 0, 0, scale, 0, 0);
        ctx.globalAlpha = 1;
        ctx.fillStyle = c.bg;
        ctx.fillRect(0, 0, W, H);
        if (!cur) return;
        drawWorld(ctx);
        if (flash.t > 0) {
            ctx.globalAlpha = Math.max(0, flash.t / flash.max) * 0.7;
            ctx.fillStyle = flash.color;
            ctx.fillRect(0, 0, W, H);
            ctx.globalAlpha = 1;
        }
        if (cur.def && cur.screen !== 'title' && cur.screen !== 'loading' && cur.screen !== 'error') drawHud(ctx);
        drawScreen(ctx);
    }

    // ------------------------------------------------------------------
    // Overlay, main loop, launch and quit
    // ------------------------------------------------------------------

    function resize() {
        if (!dom) return;
        var fit = Math.min(window.innerWidth / W, window.innerHeight / H);
        dom.canvas.style.width = Math.floor(W * fit) + 'px';
        dom.canvas.style.height = Math.floor(H * fit) + 'px';
        // Render at device resolution, but never below the logical size and
        // never so large that a slow GPU struggles.
        scale = G.clamp(fit * Math.min(window.devicePixelRatio || 1, 2), 1, 3);
        dom.canvas.width = Math.round(W * scale);
        dom.canvas.height = Math.round(H * scale);
    }

    function loop(now) {
        if (!G.active) return;
        raf = requestAnimationFrame(loop);
        acc += Math.min(0.1, (now - lastTime) / 1000);
        lastTime = now;
        while (acc >= STEP) { tick(); acc -= STEP; }
        frame();
    }

    function openOverlay() {
        var root = document.createElement('div');
        root.id = 'sno-game';
        root.setAttribute('role', 'dialog');
        root.setAttribute('aria-modal', 'true');
        root.setAttribute('aria-label', 'Game');
        var canvas = document.createElement('canvas');
        var close = document.createElement('button');
        close.type = 'button'; close.className = 'sno-game-close'; close.textContent = '✕ ESC';
        close.setAttribute('aria-label', 'Quit game');
        close.addEventListener('click', function () { G.quit(); });
        root.appendChild(canvas); root.appendChild(close);
        canvas.addEventListener('pointerdown', onPointerDown);
        canvas.addEventListener('pointermove', onPointerMove);
        canvas.addEventListener('pointercancel', onPointerUp);
        window.addEventListener('pointerup', onPointerUp);
        canvas.addEventListener('contextmenu', function (e) { e.preventDefault(); });
        buildTouchPad(root);
        root.tabIndex = -1;
        document.body.appendChild(root);
        document.body.classList.add('sno-game-on');
        // Move focus into the dialog so a focused blog button cannot be
        // activated from inside the game.
        root.focus({ preventScroll: true });
        dom = { root: root, canvas: canvas, ctx: canvas.getContext('2d') };
        G.fontFamily = getComputedStyle(document.body).fontFamily || 'monospace';
        resize();
    }

    function currentTheme() {
        return window.SNONUX_CURRENT_THEME || document.documentElement.getAttribute('data-sno-theme') || '';
    }

    function loadGameScript(theme) {
        var s = document.createElement('script');
        s.src = 'themes/' + theme + '/game.js?b=' + encodeURIComponent(window.SNONUX_BUILD || '');
        var stillWaiting = function () { return cur && cur.theme === theme && !cur.def && cur.screen === 'loading'; };
        s.onerror = function () { if (stillWaiting()) fail('could not load ' + s.src); };
        // A script with a syntax error loads "successfully" but never
        // registers; do not leave LOADING up forever.
        s.onload = function () { if (stillWaiting()) fail(s.src + ' did not register a game'); };
        document.head.appendChild(s);
    }

    function boot(theme) {
        cur.def = defs[theme];
        dom.canvas.style.cursor = cur.def.cursor || 'default';
        musicStart(cur.def.music, 1);
        toTitle();
    }

    // SnoGame.register(theme, def) — called by themes/<theme>/game.js.
    G.register = function (theme, def) {
        defs[theme] = def;
        if (cur && cur.theme === theme && !cur.def) boot(theme);
    };

    function pauseAmbient() {
        ambientWasPlaying = !!(window.snonuxAmbientIsPlaying && window.snonuxAmbientIsPlaying());
        if (ambientWasPlaying && window.snonuxAmbientPause) window.snonuxAmbientPause('game');
    }

    // SnoGame.launch(theme) — opens the overlay and starts that theme's game
    // (the active theme's game when no name is given).
    G.launch = function (theme) {
        theme = theme || currentTheme();
        if (G.active || !TITLES[theme]) return false;
        G.active = true;
        errors.length = 0;
        clearInput(); clearFx();
        pauseAmbient();
        openOverlay();
        holdAudio(false);
        audio();
        cur = { theme: theme, def: null, s: null, screen: 'loading', timer: 0, age: 0, level: 1, sel: 1, save: loadSave(theme) };
        if (defs[theme]) boot(theme); else loadGameScript(theme);
        lastTime = performance.now(); acc = 0;
        raf = requestAnimationFrame(loop);
        return true;
    };

    // SnoGame.quit() — closes the game and hands the page back to the blog.
    G.quit = function () {
        if (!G.active) return;
        G.active = false;
        cancelAnimationFrame(raf);
        if (cur && cur.def) writeSave();
        musicStop();
        // Silence anything already scheduled (jingles, music lookahead).
        holdAudio(true);
        lingering = rawDown;
        clearInput();
        window.removeEventListener('pointerup', onPointerUp);
        if (dom && dom.root.parentNode) dom.root.parentNode.removeChild(dom.root);
        document.body.classList.remove('sno-game-on');
        dom = null; cur = null;
        if (ambientWasPlaying && window.snonuxAmbientStart) window.snonuxAmbientStart('game');
        // Give keyboard focus back to the splash when the game was launched
        // from it, so Enter still opens the blog.
        var splash = document.getElementById('splash-overlay');
        if (splash && !splash.classList.contains('splash--dismissed')) splash.focus({ preventScroll: true });
    };

    window.addEventListener('keydown', onKeyDown, true);
    window.addEventListener('keyup', onKeyUp, true);
    window.addEventListener('resize', resize);
    window.addEventListener('blur', function () { if (G.active) clearInput(); });
    // A hidden tab pauses play and silences every other screen too, so a
    // throttled background tab never plays stuttering music.
    document.addEventListener('visibilitychange', function () {
        if (!G.active || !cur) return;
        if (document.hidden && cur.screen === 'play') togglePause();
        else if (cur.screen !== 'paused') holdAudio(document.hidden);
    });

    // ------------------------------------------------------------------
    // Launch buttons: splash, header and the fx row
    // ------------------------------------------------------------------

    var ICON = '<i class="fas fa-gamepad" aria-hidden="true"></i> ';

    function escapeHTML(s) {
        return String(s).replace(/[&<>"]/g, function (ch) {
            return { '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;' }[ch];
        });
    }

    // Creates the launch button inside `parent` on first use and keeps its
    // label in step with the active theme's game.
    function ensureButton(parent, cls, wrapCls, title, prepend) {
        var btn = parent.querySelector('.' + cls);
        if (!title) { if (btn) ((wrapCls && btn.closest('.' + wrapCls)) || btn).remove(); return; }
        if (!btn) {
            btn = document.createElement('button');
            btn.type = 'button';
            var node = btn;
            if (wrapCls) { node = document.createElement('div'); node.className = wrapCls; node.appendChild(btn); }
            if (prepend) parent.insertBefore(node, parent.firstChild); else parent.appendChild(node);
        }
        btn.className = cls + ' sno-game-launch' + (wrapCls ? ' splash-music-btn' : '');
        btn.setAttribute('aria-label', 'Play the game: ' + title);
        btn.innerHTML = ICON + (wrapCls ? 'Play ' : '') + escapeHTML(title);
    }

    // snonuxGameDecorate() is idempotent. shared.js calls it again whenever a
    // theme switch has replaced the splash or header markup.
    function decorate() {
        var title = TITLES[currentTheme()];
        var splash = document.getElementById('splash-overlay');
        if (splash) ensureButton(splash.querySelector('.splash-inner') || splash, 'splash-game-btn', 'splash-game-choice', title, false);
        var header = document.querySelector('header');
        if (header) ensureButton(header.querySelector('.nav') || header, 'header-game-btn', '', title, true);
        document.querySelectorAll('[data-sno-fx="game"]').forEach(function (b) {
            b.style.display = title ? '' : 'none';
            if (title) b.setAttribute('aria-label', 'Play the game: ' + title);
        });
    }

    // One delegated capture-phase listener: it runs before the splash's own
    // click-to-dismiss handler, so launching a game leaves the splash open.
    document.addEventListener('click', function (e) {
        var btn = e.target.closest && e.target.closest('.sno-game-launch');
        if (!btn) return;
        e.preventDefault(); e.stopPropagation();
        G.launch();
    }, true);

    window.snonuxGameDecorate = decorate;
    window.snonuxGameLaunch = function (theme) { return G.launch(theme); };
    decorate();

    // ------------------------------------------------------------------
    // Test hooks (used by integrationtests/games/e2e.mjs and by game authors)
    // ------------------------------------------------------------------

    function applyInput(inp, prev) {
        var k = G.key, name;
        for (name in k) {
            var on = !!inp[name];
            if (on && !prev[name]) G.hit[name] = true;
            k[name] = on;
        }
        var codes = inp.codes || {}, code;
        for (code in codes) if (codes[code] && !rawDown[code]) rawHit[code] = true;
        rawDown = {};
        for (code in codes) if (codes[code]) rawDown[code] = true;
        if (inp.mouse) {
            if (inp.mouse.x != null) { G.mouse.x = inp.mouse.x; G.mouse.y = inp.mouse.y; }
            if (inp.mouse.down && !G.mouse.down) G.mouse.hit = true;
            if (inp.mouse.rdown && !G.mouse.rdown) G.mouse.rhit = true;
            G.mouse.down = !!inp.mouse.down;
            G.mouse.rdown = !!inp.mouse.rdown;
        }
    }

    G.debug = {
        errors: errors,
        parseSave: parseSave,
        state: function () {
            if (!cur) return null;
            return {
                theme: cur.theme, screen: cur.screen, level: cur.level, sel: cur.sel, lives: G.lives,
                score: G.score, unlocked: cur.save.unlocked, hi: cur.save.hi, won: cur.save.won,
                registered: !!cur.def, t: G.t
            };
        },
        s: function () { return cur && cur.s; },
        // What the registered game declares, so tests can check it is complete.
        def: function () {
            var d = cur && cur.def;
            if (!d) return null;
            return {
                title: d.title || '', blurb: d.blurb || '', controls: (d.controls || []).length,
                levelNames: (d.levelNames || []).length, lives: maxLives(),
                // A tune needs a tempo and at least one voice that plays.
                music: !!(d.music && d.music.bpm && (d.music.lead || d.music.bass || d.music.arp) && /[^.\-]/.test([].concat(d.music.lead || '', d.music.bass || '', d.music.arp || '').join('')))
            };
        },
        audio: function () {
            return { state: actx ? actx.state : 'none', tones: audioStats.tones, noises: audioStats.noises, musicSteps: audioStats.musicSteps, sfx: audioStats.sfx, engineSfx: audioStats.engineSfx, muted: muted };
        },
        // Jumps straight into play on the given level (no intro banner).
        start: function (level) { G.score = 0; beginLevel(level); if (cur.screen !== 'error') setScreen('play'); },
        // Runs n update ticks synchronously. `input` is either a fixed
        // {left, right, up, down, a, b, codes, mouse} object or a function
        // (i, state, G) returning one per tick — handy for scripted bots.
        step: function (n, input) {
            var prev = {};
            for (var i = 0; i < n; i++) {
                var inp = (typeof input === 'function' ? input(i, cur.s, G) : input) || {};
                applyInput(inp, prev);
                prev = inp;
                tick();
            }
            clearInput();
            return G.debug.state();
        },
        win: function () { G.win(0); },
        lose: function () { while (cur && cur.screen === 'play') G.loseLife(); }
    };
})();
