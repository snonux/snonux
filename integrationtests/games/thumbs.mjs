// Regenerates the pictures of the game overview menu (shared/arcade.js):
//
//   node integrationtests/games/thumbs.mjs            every game
//   node integrationtests/games/thumbs.mjs neon dos   only these
//
// Each game is started on LEVEL, played for a few seconds by a seeded random
// bot (so a rerun gives the same picture), frozen, and its canvas is written
// as internal/generator/templates/themes/<theme>/game-thumb.jpg. Run it again
// after changing how a game looks.
import { writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { sleep, startSession, themesWithGames, THEMES_DIR } from './lib.mjs';

const LEVEL = 3, THUMB_W = 320, THUMB_H = 180;
// Ticks to play before the picture is taken; a game the bot loses that
// quickly is retried with the next, shorter run.
const RUNS = [300, 150, 60, 10];

// Wraps SnoGame.register so update() can be frozen once the scene is right,
// and makes Math.random repeatable.
const PATCH = `(function(){var R=SnoGame.register;SnoGame.register=function(t,d){var u=d.update;
 d.update=function(){if(window.__freeze)return;return u.apply(this,arguments);};return R.call(SnoGame,t,d);};
 window.__seed=function(){var a=12345;Math.random=function(){a|=0;a=a+0x6D2B79F5|0;var t=Math.imul(a^a>>>15,1|a);
 t=t+Math.imul(t^t>>>7,61|t)^t;return((t^t>>>14)>>>0)/4294967296;};};return 1;})()`;

const BOT = `(function(){var r=SnoGame.rng(7),c={},u=0;return function(i){
 if(i>=u){u=i+5+Math.floor(r()*40);c={left:r()<.35,right:r()<.35,up:r()<.5,down:r()<.1,a:r()<.6,b:r()<.2,
 mouse:{x:r()*960,y:r()*540,down:r()<.5}};}return c;};})()`;

// Plays `steps` ticks and reports the screen the game ended on.
function play(page, steps) {
    return page.eval(`(function(){window.__freeze=0;__seed();SnoGame.debug.start(${LEVEL});var bot=${BOT},n=0,st=SnoGame.debug.state();
      for(;n<${steps};n+=5){st=SnoGame.debug.step(5,function(i){return bot(n+i);});if(st.screen!=='play')break;}
      window.__freeze=1;return st.screen;})()`);
}

const SHRINK = `(function(){var src=document.querySelector('#sno-game canvas'),c=document.createElement('canvas');
 c.width=${THUMB_W};c.height=${THUMB_H};var x=c.getContext('2d');x.imageSmoothingQuality='high';
 x.drawImage(src,0,0,${THUMB_W},${THUMB_H});return c.toDataURL('image/jpeg',0.82);})()`;

async function shoot(page, base, theme) {
    await page.goto(`${base}/index.html`);
    await page.eval(PATCH);
    await page.eval(`SnoGame.launch('${theme}')`);
    await page.waitFor(`SnoGame.debug.state() && SnoGame.debug.state().registered`, 'register');
    let steps = 0, screen = '';
    for (steps of RUNS) {
        screen = await play(page, steps);
        if (screen === 'play') break;
    }
    await sleep(300);    // let the frozen scene be drawn
    const url = await page.eval(SHRINK);
    const file = join(THEMES_DIR, theme, 'game-thumb.jpg');
    writeFileSync(file, Buffer.from(url.split(',')[1], 'base64'));
    console.log(`${screen === 'play' ? 'ok  ' : 'WARN'} ${theme} (${steps} ticks, screen ${screen})`);
}

const wanted = process.argv.slice(2);
const themes = wanted.length ? wanted : themesWithGames();
const { page, base, close } = await startSession();
try {
    for (const theme of themes) await shoot(page, base, theme);
} finally { await close(); }
