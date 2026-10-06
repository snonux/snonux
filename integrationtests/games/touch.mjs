// Phone checks for the theme games: the touch half of e2e.mjs (--touch).
//
// Chrome is turned into a phone (mobile viewport, touch screen, coarse
// pointer) and every interaction is a real touch event. Per theme and per
// orientation it checks that the game launches from a tap on the splash
// button; that the touch pad is there, big enough, inside the viewport and
// nowhere on the canvas; that every pad button holds its key while touched
// and only then; that two fingers hold two buttons; that a full direction
// pad gives diagonals and follows a sliding thumb; that a tap on the canvas
// starts the level and, in play, moves SnoGame.mouse to the tap; and that
// the ✕ button quits.

import { join } from 'node:path';
import { sleep } from './lib.mjs';

// The two ways a phone is held (CSS pixels of a current mid-size phone).
export const PHONES = [
    { name: 'landscape', width: 844, height: 390 },
    { name: 'portrait', width: 390, height: 844 },
];

const STATE = 'SnoGame.debug.state()';
const MIN_BUTTON = 48;  // CSS px: the smallest target a thumb hits reliably
const MIN_CANVAS = 280; // CSS px wide: below this the playfield is unreadable

// The logical button each key code of the pad must also raise.
const LOGICAL = {
    ArrowLeft: 'left', ArrowRight: 'right', ArrowUp: 'up', ArrowDown: 'down',
    KeyA: 'left', KeyD: 'right', KeyW: 'up', KeyS: 'down', Space: 'a', KeyX: 'b',
};
const ARROWS = { left: 'ArrowLeft', right: 'ArrowRight', up: 'ArrowUp', down: 'ArrowDown' };
const WASD = { left: 'KeyA', right: 'KeyD', up: 'KeyW', down: 'KeyS' };

// Rectangles of everything the layout places, read from the live page.
const LAYOUT = `(function(){
  function r(e){var b=e.getBoundingClientRect();return {x:b.left,y:b.top,w:b.width,h:b.height};}
  var root=document.getElementById('sno-game'),pad=root&&root.querySelector('.sno-game-pad');
  if(!pad)return null;
  return {vw:innerWidth,vh:innerHeight,canvas:r(root.querySelector('canvas')),close:r(root.querySelector('.sno-game-close')),
    shown:getComputedStyle(pad).display!=='none'&&getComputedStyle(pad).visibility!=='hidden',
    buttons:[].map.call(pad.querySelectorAll('[data-code]'),function(e){var b=r(e);b.code=e.getAttribute('data-code');b.label=e.textContent;
      b.part=e.parentNode.className.split(' ')[0]+(/fire/.test(e.parentNode.className)?'-fire':'');return b;})};})()`;

function check(cond, message) { if (!cond) throw new Error(message); }

const centre = (r) => ({ x: r.x + r.w / 2, y: r.y + r.h / 2 });
const overlaps = (a, b) => a.x < b.x + b.w && a.x + a.w > b.x && a.y < b.y + b.h && a.y + a.h > b.y;
const inside = (r, l) => r.x >= 0 && r.y >= 0 && r.x + r.w <= l.vw && r.y + r.h <= l.vh;

// The key codes a game's def.touch asks the pad for.
function expectedCodes(touch) {
    const t = touch || {}, hide = new Set(t.hide || []), codes = [];
    for (const dir of ['left', 'right', 'up', 'down']) {
        if (hide.has(dir)) continue;
        codes.push(t.twin ? WASD[dir] : ARROWS[dir]);
        if (t.twin) codes.push(ARROWS[dir]);
    }
    if (!hide.has('a')) codes.push('Space');
    if (!hide.has('b')) codes.push('KeyX');
    return codes.sort();
}

async function launchByTap(page, base, theme) {
    await page.goto(`${base}/index.html`);
    check(await page.eval(`matchMedia('(pointer: coarse)').matches`), 'phone emulation did not produce a coarse pointer');
    await page.eval(`snonuxSwitchTheme(${JSON.stringify(theme)})`);
    const title = await page.eval(`SnoGame.titles[${JSON.stringify(theme)}]`);
    const btn = '#splash-overlay .splash-game-btn';
    await page.waitFor(`(document.querySelector('${btn}')||{textContent:''}).textContent.indexOf(${JSON.stringify(title)})>=0`, 'splash "Play" button for this theme');
    const r = await page.eval(`(function(){var e=document.querySelector('${btn}');e.scrollIntoView({block:'center'});var b=e.getBoundingClientRect();return {x:b.left,y:b.top,w:b.width,h:b.height};})()`);
    await page.tap(centre(r).x, centre(r).y);
    await page.waitFor(`SnoGame.active && ${STATE} && ${STATE}.screen==='title'`, 'title screen after tapping the splash button');
    check((await page.eval(STATE)).theme === theme, `tap launched the wrong game`);
}

// Geometry: the pad exists, can be hit, and stays off the playfield.
function checkLayout(l, phone, touch) {
    check(l && l.shown, 'no visible touch pad');
    const want = expectedCodes(touch), got = l.buttons.map((b) => b.code).sort();
    check(want.length > 0, 'def.touch hides every pad button');
    check(JSON.stringify(got) === JSON.stringify(want), `pad buttons ${got} do not match def.touch (want ${want})`);
    check(l.canvas.w >= MIN_CANVAS && inside(l.canvas, l), `canvas ${Math.round(l.canvas.w)}x${Math.round(l.canvas.h)} is too small or off screen`);
    check(Math.abs(l.canvas.w / l.canvas.h - 16 / 9) < 0.02, 'canvas is not 16:9');
    check(inside(l.close, l) && l.close.w >= 40 && l.close.h >= 40, 'close button is off screen or too small to tap');
    check(!overlaps(l.close, l.canvas), 'close button covers the canvas');
    for (const b of l.buttons) {
        check(b.w >= MIN_BUTTON && b.h >= MIN_BUTTON, `pad button ${b.code} is ${Math.round(b.w)}x${Math.round(b.h)}px, below ${MIN_BUTTON}`);
        check(inside(b, l), `pad button ${b.code} is outside the viewport`);
        check(!overlaps(b, l.canvas), `pad button ${b.code} covers the canvas`);
        check(!overlaps(b, l.close), `pad button ${b.code} overlaps the close button`);
    }
    checkPlacement(l, phone);
}

// Landscape puts the pad in gutters beside the canvas, portrait below it.
function checkPlacement(l, phone) {
    const c = l.canvas;
    if (phone.name === 'portrait') {
        check(l.buttons.every((b) => b.y >= c.y + c.h), 'portrait: the pad must be below the canvas');
        check(c.y === 0, 'portrait: the canvas must be at the top');
        return;
    }
    const left = l.buttons.filter((b) => b.x + b.w <= c.x), right = l.buttons.filter((b) => b.x >= c.x + c.w);
    check(left.length + right.length === l.buttons.length, 'landscape: every pad button must be beside the canvas');
    check(left.length > 0, 'landscape: nothing in the left gutter');
    check(right.length > 0 || l.close.x >= c.x + c.w, 'landscape: nothing in the right gutter');
}

const keyState = (codes) => `(function(){var o={};${JSON.stringify(codes)}.forEach(function(c){o[c]=SnoGame.down(c);});
  o.key=JSON.parse(JSON.stringify(SnoGame.key));return o;})()`;

// Every button on its own: down while touched, up again after.
async function checkEachButton(page, l) {
    const codes = l.buttons.map((b) => b.code);
    for (const b of l.buttons) {
        const p = centre(b);
        await page.touch('touchStart', [p]);
        let st = await page.eval(keyState(codes));
        check(st[b.code], `touching ${b.code} did not set SnoGame.down('${b.code}')`);
        check(st.key[LOGICAL[b.code]], `touching ${b.code} did not set SnoGame.key.${LOGICAL[b.code]}`);
        check(codes.every((c) => c === b.code || !st[c]), `touching ${b.code} also pressed ${codes.filter((c) => c !== b.code && st[c])}`);
        await page.touch('touchEnd', []);
        st = await page.eval(keyState(codes));
        check(codes.every((c) => !st[c]), `releasing ${b.code} left ${codes.filter((c) => st[c])} down`);
        check(Object.values(st.key).every((v) => !v), `releasing ${b.code} left a logical button down`);
    }
}

// Two fingers on two different parts of the pad: both hold, and lifting one
// leaves the other held.
async function checkTwoFingers(page, l) {
    const first = l.buttons[0], second = l.buttons.find((b) => b.part !== first.part);
    if (!second) return false; // a pad with a single part: nothing to combine
    const p1 = { ...centre(first), id: 1 }, p2 = { ...centre(second), id: 2 };
    await page.touch('touchStart', [p1]);
    await page.touch('touchStart', [p1, p2]);
    let st = await page.eval(keyState([first.code, second.code]));
    check(st[first.code] && st[second.code], `two fingers: ${first.code}=${st[first.code]} ${second.code}=${st[second.code]}, want both`);
    await page.touch('touchEnd', [p2]);
    st = await page.eval(keyState([first.code, second.code]));
    check(st[first.code] && !st[second.code], `lifting the second finger: ${first.code}=${st[first.code]} ${second.code}=${st[second.code]}`);
    await page.touch('touchEnd', []);
    st = await page.eval(keyState([first.code, second.code]));
    check(!st[first.code] && !st[second.code], 'two fingers: a button stayed down after both lifted');
    return true;
}

// A full direction pad is an 8-way stick: its corner is a diagonal, and a
// thumb sliding across it changes direction without lifting.
async function checkDiagonal(page, l, part) {
    const btn = (code) => l.buttons.find((b) => b.part === part && LOGICAL[b.code] === code);
    const up = btn('up'), left = btn('left'), right = btn('right');
    if (!up || !left || !right) return false;
    const codes = [up.code, left.code, right.code];
    await page.touch('touchStart', [{ x: centre(left).x, y: centre(up).y }]);
    let st = await page.eval(keyState(codes));
    check(st[up.code] && st[left.code] && !st[right.code], `${part}: the up-left corner gave ${codes.filter((c) => st[c])}`);
    await page.touch('touchMove', [centre(right)]);
    st = await page.eval(keyState(codes));
    check(st[right.code] && !st[left.code] && !st[up.code], `${part}: sliding to the right gave ${codes.filter((c) => st[c])}`);
    await page.touch('touchEnd', []);
    st = await page.eval(keyState(codes));
    check(codes.every((c) => !st[c]), `${part}: a direction stayed down after the slide`);
    return true;
}

// Client coordinates of a logical canvas point (960x540 space).
const onCanvas = (l, x, y) => ({ x: l.canvas.x + x / 960 * l.canvas.w, y: l.canvas.y + y / 540 * l.canvas.h });

// A tap above the level boxes starts the selected level from the title.
async function checkTapStarts(page, l) {
    check((await page.eval(STATE)).screen === 'title', 'expected the title screen');
    const p = onCanvas(l, 480, 60);
    await page.tap(p.x, p.y);
    await page.waitFor(`${STATE}.screen==='intro' || ${STATE}.screen==='play'`, 'the level to start after a tap on the canvas');
}

// In play a finger on the canvas is the mouse: position and held button.
async function checkTapAims(page, l) {
    for (const [x, y] of [[300, 200], [820, 430]]) {
        await page.eval('SnoGame.debug.start(1)');
        const p = onCanvas(l, x, y);
        await page.touch('touchStart', [p]);
        const m = await page.eval(`({x:SnoGame.mouse.x,y:SnoGame.mouse.y,down:SnoGame.mouse.down,screen:${STATE}.screen})`);
        await page.touch('touchEnd', []);
        check(Math.abs(m.x - x) <= 3 && Math.abs(m.y - y) <= 3, `tap at ${x},${y} put SnoGame.mouse at ${m.x.toFixed(1)},${m.y.toFixed(1)}`);
        check(m.down, 'a finger on the canvas did not set SnoGame.mouse.down');
        check(!(await page.eval('SnoGame.mouse.down')), 'lifting the finger left SnoGame.mouse.down set');
    }
}

async function checkCloseQuits(page, l) {
    await page.tap(centre(l.close).x, centre(l.close).y);
    await page.waitFor(`!SnoGame.active && !document.getElementById('sno-game')`, 'the ✕ button to quit');
    check(!(await page.eval('document.fullscreenElement')), 'quitting left the page in fullscreen');
}

// Everything about the pad itself, for the def.touch the game declared:
// where it is, and that its buttons, two fingers and diagonals work.
async function checkPad(page, phone, touch) {
    const l = await page.eval(LAYOUT);
    checkLayout(l, phone, touch);
    await checkEachButton(page, l);
    const combined = await checkTwoFingers(page, l);
    const parts = [...new Set(l.buttons.filter((b) => b.part.startsWith('sno-game-dpad')).map((b) => b.part))];
    let diagonals = 0;
    for (const part of parts) if (await checkDiagonal(page, l, part)) diagonals++;
    // A pad is allowed to be one part, or to have no full stick, but only
    // because its def.touch says so: never by accident.
    const hidden = (touch && touch.hide) || [];
    check(combined || hidden.includes('a') && hidden.includes('b'), 'no second pad part to hold together with the first');
    check(diagonals === (touch && touch.twin ? 2 : 1) || hidden.some((h) => ARROWS[h]), `${diagonals} direction pad(s) with diagonals`);
    return l;
}

async function testOrientation(page, base, theme, phone, shots) {
    await page.emulatePhone(phone.width, phone.height);
    await launchByTap(page, base, theme);
    const touch = (await page.eval('SnoGame.debug.def()')).touch;
    if (shots) { await sleep(120); await page.screenshot(join(shots, `${theme}-touch-${phone.name}-title.jpg`)); }
    await checkTapStarts(page, await page.eval(LAYOUT));
    await page.eval('SnoGame.debug.start(1)');
    const l = await checkPad(page, phone, touch);
    await checkTapAims(page, l);
    if (shots) { await sleep(120); await page.screenshot(join(shots, `${theme}-touch-${phone.name}-play.jpg`)); }
    const errs = await page.eval('SnoGame.debug.errors.slice()');
    check(!errs.length, `engine error: ${errs[0]}`);
    await checkCloseQuits(page, l);
}

// def.touch variants no shipped game may happen to use; each is tried on a
// stand-in game so the pad options stay covered whatever the games declare.
const PAD_CONFIGS = [
    { twin: true },
    { twin: true, hide: ['a', 'b'] },
    { a: 'JUMP', b: 'BOMB', hide: ['up', 'down'] },
    { a: 'FIRE', hide: ['b', 'down'] },
    { hide: ['a', 'b'] },
];

const standIn = (touch) => `SnoGame.register('breakout', { touch: ${JSON.stringify(touch)},
  init: function () { return {}; }, update: function () {}, draw: function () {} })`;

// Engine-level check of def.touch: buttons, labels and layout per variant.
export async function testPadConfigs(page, base, shots) {
    for (const phone of PHONES) {
        await page.emulatePhone(phone.width, phone.height);
        for (const [i, touch] of PAD_CONFIGS.entries()) {
            const what = `def.touch ${JSON.stringify(touch)} in ${phone.name}`;
            await page.goto(`${base}/index.html`);
            await page.eval(`${standIn(touch)}; SnoGame.launch('breakout'); SnoGame.debug.start(1)`);
            let l;
            try { l = await checkPad(page, phone, touch); } catch (err) { err.message = `${what}: ${err.message}`; throw err; }
            const label = (code) => (l.buttons.find((b) => b.code === code) || {}).label;
            check(label('Space') === ((touch.hide || []).includes('a') ? undefined : touch.a || 'A'), `${what}: A is labelled ${label('Space')}`);
            check(label('KeyX') === ((touch.hide || []).includes('b') ? undefined : touch.b || 'B'), `${what}: B is labelled ${label('KeyX')}`);
            if (shots) await page.screenshot(join(shots, `pad-config-${i}-${phone.name}.jpg`));
            await checkCloseQuits(page, l);
        }
    }
}

// Runs every phone orientation for one theme.
export async function testThemeTouch(page, base, theme, shots) {
    page.problems.length = 0;
    await page.send('Network.clearBrowserCookies');
    for (const phone of PHONES) {
        try {
            await testOrientation(page, base, theme, phone, shots);
        } catch (err) {
            err.message = `${phone.name} ${phone.width}x${phone.height}: ${err.message}`;
            throw err;
        }
    }
    check(!page.problems.length, `page errors: ${page.problems.slice(0, 3).join(' | ')}`);
}
