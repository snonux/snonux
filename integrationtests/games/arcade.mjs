// Checks of the game overview menu (shared/arcade.js), run by e2e.mjs in both
// modes: with mouse and keyboard on the desktop, and with taps on a phone.
import { join } from 'node:path';
import { sleep, themesWithGames } from './lib.mjs';

const STATE = 'SnoGame.debug.state()';
const SPLASH_UP = `!document.getElementById('splash-overlay').classList.contains('splash--dismissed')`;
// The game picked from the menu: deliberately not the test site's own theme
// (neon), since starting another theme's game is the point of the menu.
const PICK = 'dos';

function check(cond, message) { if (!cond) throw new Error(message); }

// Centre of an element in client pixels, scrolled into view first.
function centre(page, selector) {
    return page.eval(`(function(){var e=document.querySelector(${JSON.stringify(selector)});if(!e)return null;
      e.scrollIntoView({block:'center'});var r=e.getBoundingClientRect();return {x:r.left+r.width/2,y:r.top+r.height/2};})()`);
}

async function activate(page, selector, touch) {
    const p = await centre(page, selector);
    check(p, `no ${selector} on the page`);
    if (touch) await page.tap(p.x, p.y); else await page.click(selector);
}

// Every game is listed once, with its title, its theme, a description and a
// picture the server really has.
async function checkCards(page, base) {
    const cards = await page.eval(`Array.prototype.map.call(document.querySelectorAll('.sno-arcade-card'),function(c){
      return {theme:c.getAttribute('data-theme'),title:c.querySelector('.sno-arcade-title').textContent,
        label:c.querySelector('.sno-arcade-theme').textContent,blurb:c.querySelector('.sno-arcade-blurb').textContent,
        img:c.querySelector('img').getAttribute('src')};})`);
    const themes = themesWithGames();
    check(cards.length === themes.length, `the menu lists ${cards.length} games, want ${themes.length}`);
    check(cards[0].theme === 'neon' && /this theme/.test(cards[0].label), `the active theme's game must come first, got ${cards[0].theme}`);
    for (const theme of themes) {
        const c = cards.find((x) => x.theme === theme);
        check(c, `${theme} is missing from the menu`);
        check(c.title === await page.eval(`SnoGame.titles['${theme}']`), `${theme}: wrong title "${c.title}"`);
        check(c.blurb.length > 20, `${theme}: no description`);
        check(c.label.toLowerCase().includes(theme), `${theme}: the card does not name its theme`);
        const res = await fetch(`${base}/${c.img}`);
        check(res.ok && res.headers.get('content-type') === 'image/jpeg', `${theme}: picture ${c.img} not served (${res.status})`);
    }
}

// Arrow keys walk the cards, and no key gets through to the splash behind
// (Enter there would open the blog, "a" would start the theme's game).
async function checkKeys(page) {
    const focused = `document.activeElement.getAttribute('data-theme')`;
    const first = await page.eval(focused);
    await page.press('ArrowRight');
    const second = await page.eval(focused);
    check(second && second !== first, 'ArrowRight did not move to the next card');
    await page.press('ArrowDown');
    check(![first, second].includes(await page.eval(focused)), 'ArrowDown did not move a row down');
    await page.press('ArrowUp'); await page.press('ArrowLeft');
    check(await page.eval(focused) === first, 'ArrowUp and ArrowLeft did not lead back to the first card');
    await page.press('KeyA'); await page.press('KeyW');
    await sleep(150);
    check(await page.eval(SPLASH_UP), 'a key pressed in the menu dismissed the splash');
    check(!(await page.eval('SnoGame.active')), 'a key pressed in the menu started a game');
    check(await page.eval('SnoArcade.isOpen()'), 'a letter key closed the menu');
}

// A card starts its game; quitting the game leads back to the menu, on the
// same card.
async function checkPlayAndReturn(page, touch) {
    await activate(page, `.sno-arcade-card[data-theme="${PICK}"]`, touch);
    await page.waitFor(`${STATE} && ${STATE}.registered && ${STATE}.screen==='title'`, `the title screen of ${PICK}`);
    check((await page.eval(STATE)).theme === PICK, 'the wrong game started');
    check(!(await page.eval('SnoArcade.isOpen()')), 'the menu stayed open under the game');
    if (touch) check(await page.eval(`!!document.querySelector('.sno-game-pad')`), 'a game started by a tap in the menu has no touch pad');
    if (touch) await page.eval('SnoGame.quit()'); else await page.press('Escape');
    await page.waitFor('!SnoGame.active && SnoArcade.isOpen()', 'the menu to come back after the game');
    check(await page.eval(`document.activeElement.getAttribute('data-theme')`) === PICK, 'the menu did not return to the card that was played');
}

export async function testArcadeMenu(page, base, shots, touch) {
    if (touch) await page.emulatePhone(390, 844);
    await page.goto(`${base}/index.html`);
    await page.waitFor(`!!document.querySelector('.splash-arcade-btn')`, 'the "All games" button on the splash');
    await activate(page, '.splash-arcade-btn', touch);
    await page.waitFor('SnoArcade.isOpen()', 'the menu to open');
    check(await page.eval(SPLASH_UP), 'opening the menu dismissed the splash');
    await checkCards(page, base);
    check(await page.eval(`(function(){var r=document.getElementById('sno-arcade');return r.scrollWidth<=r.clientWidth;})()`), 'the menu is wider than the screen');
    if (shots) await page.screenshot(join(shots, `arcade-menu${touch ? '-touch' : ''}.jpg`));
    if (!touch) await checkKeys(page);
    await checkPlayAndReturn(page, touch);

    // Closing: Esc on the desktop, the close button by touch. A game started
    // the ordinary way afterwards must not bring the menu back when it ends.
    if (touch) await activate(page, '.sno-arcade-close', true); else await page.press('Escape');
    await page.waitFor('!SnoArcade.isOpen()', 'the menu to close');
    check(await page.eval(SPLASH_UP), 'closing the menu dismissed the splash');
    await page.eval(`SnoGame.launch()`);
    await page.waitFor(`${STATE} && ${STATE}.screen==='title'`, 'the title screen');
    await page.eval('SnoGame.quit()');
    await sleep(150);
    check(!(await page.eval('SnoArcade.isOpen()')), 'the menu opened after a game that was not started from it');

    // The header has the button too, once the splash is gone.
    await page.eval('window._snonuxDismissSplash && window._snonuxDismissSplash()');
    await page.waitFor(`!${SPLASH_UP}`, 'the splash to go');
    await sleep(touch ? 900 : 300);
    await activate(page, '.header-arcade-btn', touch);
    await page.waitFor('SnoArcade.isOpen()', 'the menu to open from the header');
    await page.eval('SnoArcade.close()');
}
