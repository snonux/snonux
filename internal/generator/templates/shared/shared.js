    var SNONUX_SOUNDS = (typeof window !== "undefined" && window.SNONUX_SOUNDS) || {};
    // Inject wild-mode badge used by all themes
    (function() { var b=document.createElement('div'); b.id='sno-wild-badge'; b.textContent='WILD MODE'; document.body.appendChild(b); })();
    const SNONUX_WILD_PRESETS = {
        aurora: {
            banner: 'SOLAR STORM',
            ticker: ['FIELD INTERFERENCE', 'PLASMA DRIFT', 'PARTICLE BOMBARDMENT', 'CHROMATIC SPIKE'],
            scraps: ['MAGNETIC SHEAR', 'SOLAR WIND', 'AURORA NOISE', 'ION STORM', 'POLAR ARC'],
            flash: 'rgba(220,255,240,0.72)',
            emoji: ['\u2728','\u2604\uFE0F','\u{1F320}','\u{1F30C}','\u2B50']
        },
        brutalist: {
            banner: 'STRUCTURAL COLLAPSE',
            ticker: ['CONDEMNED', 'REBAR EXPOSED', 'FOUNDATION FAILURE', 'CRACK PROPAGATION'],
            scraps: ['CONDEMNED', 'RUST BLEED', 'SHEAR WALL LOST', 'LOAD PATH BROKEN', 'SPALLING'],
            flash: 'rgba(255,210,190,0.58)',
            emoji: ['\u{1F9F1}','\u2692\uFE0F','\u26A0\uFE0F','\u{1F6A7}','\u{1F4A5}']
        },
        cosmos: {
            banner: 'SUPERNOVA',
            ticker: ['SINGULARITY LENSING', 'GAMMA BURST', 'SHOCKWAVE EXPANDING', 'SPACETIME TEAR'],
            scraps: ['WHITEOUT', 'EVENT HORIZON', 'RADIATION FRONT', 'LENS LOCK', 'CORE BREACH'],
            flash: 'rgba(255,255,255,0.8)',
            emoji: ['\u{1F30C}','\u2604\uFE0F','\u{1F4AB}','\u2B50','\u{1FA90}']
        },
        dos: {
            banner: 'KERNEL PANIC',
            ticker: ['ABORT, RETRY, FAIL?', 'MEMORY CORRUPTION', 'STACK DUMP', 'SEGMENT FAULT'],
            scraps: ['DEAD BEEF', 'C0FFEE', 'BAD SECTOR', 'IRQ STORM', 'NULL PTR', 'HEX DUMP'],
            flash: 'rgba(255,255,255,0.88)',
            emoji: ['\u{1F4BE}','\u{1F4BB}','\u26A1','\u2620\uFE0F','\u{1F41B}']
        },
        matrix: {
            banner: 'CASCADE FAILURE',
            ticker: ['SENTINEL TRACE', 'GLYPH SATURATION', 'PHOSPHOR BURN-IN', 'RAIN AT TERMINAL VELOCITY'],
            scraps: ['SENTINEL', '0XDECODE', 'OVERRIDE', 'TRACE LOST', 'MACHINE DREAM'],
            flash: 'rgba(180,255,190,0.68)',
            emoji: ['\u{1F441}\uFE0F','\u{1F4A0}','\u{1F916}','\u{1F50D}','\u26D3\uFE0F']
        },
        neon: {
            banner: 'GAS DISCHARGE',
            ticker: ['TUBE ARC', 'ULTRAVIOLET BLEED', 'STROBE LOCK', 'SHORT CIRCUIT'],
            scraps: ['ARC OVERLOAD', 'PLASMA SIGN', 'NOBLE GAS', 'HARD STROBE', 'OVERDRIVE'],
            flash: 'rgba(255,245,170,0.72)',
            emoji: ['\u26A1','\u{1F4A5}','\u{1F52E}','\u2728','\u{1F388}']
        },
        ocean: {
            banner: 'HADAL DESCENT',
            ticker: ['PRESSURE SPIKE', 'BIOLUMINESCENT SWARM', 'ABYSSAL DRAG', 'TSUNAMI FRONT'],
            scraps: ['NO SURFACE', 'CRUSH DEPTH', 'TENTACLE DRIFT', 'SONAR LOST', 'DEEP CURRENT'],
            flash: 'rgba(200,255,255,0.54)',
            emoji: ['\u{1F419}','\u{1F420}','\u{1F30A}','\u{1F41A}','\u{1F9DC}']
        },
        plasma: {
            banner: 'FUSION BREACH',
            ticker: ['CONTAINMENT FAILURE', 'TOKAMAK DISTORTION', 'THERMAL RUNAWAY', 'WHITE-BLUE CORE'],
            scraps: ['ION SPRAY', 'FIELD LOSS', 'HEAT HAZE', 'QUENCH', 'ARC SHELL'],
            flash: 'rgba(230,250,255,0.78)',
            emoji: ['\u{1F300}','\u26A1','\u{1F4A0}','\u2728','\u{1F52C}']
        },
        retro: {
            banner: 'TAPE EAT',
            ticker: ['TRACKING LOSS', 'CHROMA SPLIT', 'MAGNETIC SNOW', 'CLICK-EJECT'],
            scraps: ['NO SIGNAL', 'HEAD DRAG', 'ROLL HOLD', 'SNOW PACK', 'EJECT CYCLE'],
            flash: 'rgba(255,226,178,0.6)',
            emoji: ['\u{1F4FC}','\u{1F4FA}','\u{1F3AE}','\u{1F579}\uFE0F','\u{1F4FB}']
        },
        retrofuture: {
            banner: 'ATOMIC TWILIGHT',
            ticker: ['GEIGER STATIC', 'FALLOUT DUST', 'RADIATION BURN', 'IRRADIATED SEPIA'],
            scraps: ['FALLOUT', 'BETA LEAK', 'ASH DRIFT', 'HALF-LIFE', 'GLOW CLOUD'],
            flash: 'rgba(255,240,180,0.62)',
            emoji: ['\u2622\uFE0F','\u{1F4A3}','\u{1F3ED}','\u2623\uFE0F','\u{1F9EA}']
        },
        spaceage: {
            banner: 'RE-ENTRY BURN',
            ticker: ['HEAT SHIELD LOSS', 'PLASMA BLACKOUT', 'COMMS STATIC', 'G-FORCE COMPRESSION'],
            scraps: ['BLACKOUT', 'SPARK SHOWER', 'PLASMA SHEATH', 'HULL GLOW', 'COMMS LOST'],
            flash: 'rgba(255,220,190,0.68)',
            emoji: ['\u{1F680}','\u{1F6F8}','\u{1FA90}','\u{1F30D}','\u2B50']
        },
        synthwave: {
            banner: 'GRID COLLAPSE',
            ticker: ['VOID PERSPECTIVE', 'MOLTEN SUN', 'CHROMA TEAR', 'OUT OF MEMORY'],
            scraps: ['VOID GRID', 'SUN DRIP', 'NEON PANIC', 'FRAME DROP', 'MEMORY STARVE'],
            flash: 'rgba(255,210,255,0.68)',
            emoji: ['\u{1F305}','\u{1F3B6}','\u{1F3B9}','\u{1F338}','\u{1F52E}']
        },
        terminal: {
            banner: 'FORK BOMB',
            ticker: ['PROCESS STORM', 'STACK TRACE WATERFALL', 'MEMORY GARBAGE', 'BSOD CREEP'],
            scraps: ['PID 65535', 'STACK OVERFLOW', 'OOM KILL', 'PANIC', '(:'],
            flash: 'rgba(180,255,180,0.7)',
            emoji: ['\u{1F4BB}','\u{1F41B}','\u2620\uFE0F','\u{1F5A5}\uFE0F','\u26A1']
        },
        tropicale: {
            banner: 'CATEGORY 5',
            ticker: ['HORIZONTAL RAIN', 'STORM SURGE', 'DEBRIS FIELD', 'WIND SHEAR'],
            scraps: ['PALM SNAP', 'SURGE LINE', 'SPRAY WALL', 'FLYING ROOF', 'LANDFALL'],
            flash: 'rgba(240,255,255,0.74)',
            emoji: ['\u{1F334}','\u{1F3D6}\uFE0F','\u{1F940}','\u{1F965}','\u{1F30A}']
        },
        noir: {
            banner: 'BLACKOUT DISTRICT',
            ticker: ['BLINDS SLAMMED SHUT', 'SIREN SWEEP', 'PROJECTOR BURN', 'MIDNIGHT DOWNPOUR'],
            scraps: ['NO WITNESSES', 'WET ASPHALT', 'RED CHANNEL', 'BLUE CHANNEL', 'SMOKE CURTAIN'],
            flash: 'rgba(255,245,225,0.66)',
            emoji: ['\u{1F576}\uFE0F','\u{1F52B}','\u{1F3A9}','\u{1F6AC}','\u{1F5DD}\uFE0F']
        },
        cathedral: {
            banner: 'LAST JUDGMENT',
            ticker: ['BELL SHOCKWAVE', 'INCENSE FIRESTORM', 'ROSE WINDOW FRACTURE', 'APSE IN FLAME'],
            scraps: ['REQUIEM', 'SHARD RAIN', 'VESPER BURN', 'GLORIA STATIC', 'NAVE COLLAPSE'],
            flash: 'rgba(255,239,202,0.72)',
            emoji: ['\u{1F54E}','\u{1F56F}\uFE0F','\u271D\uFE0F','\u{1F54A}\uFE0F','\u{1F3F0}']
        },
        surveillance: {
            banner: 'TOTAL COMPROMISE',
            ticker: ['CAMERA MESH BREACH', 'TRACKING LOSS', 'MULTIPLEX PANIC', 'ALERT CASCADE'],
            scraps: ['FLAGGED', 'OVERRIDDEN', 'TRACE LOOP', 'BOX LOST', 'ALERT 99'],
            flash: 'rgba(210,255,225,0.72)',
            emoji: ['\u{1F4F9}','\u{1F441}\uFE0F','\u{1F6A8}','\u{1F50D}','\u{1F4E1}']
        },
        biomech: {
            banner: 'CONTAINMENT RUPTURE',
            ticker: ['SYNAPSE STORM', 'TISSUE ARC', 'MEMBRANE TEAR', 'HYBRID OVERDRIVE'],
            scraps: ['VENTRICLE', 'MYCELIUM', 'RUPTURE', 'BIOFILM', 'NERVE GRID'],
            flash: 'rgba(255,205,220,0.7)',
            emoji: ['\u{1F9EC}','\u{1F9E0}','\u{1F9A0}','\u{1F52C}','\u{1FAC0}']
        },
        paper: {
            banner: 'PRESS JAM',
            ticker: ['TONER BLIZZARD', 'INK BLEED', 'COPY LAMP WHITEOUT', 'PAGE STORM'],
            scraps: ['MISPRINT', 'SKEWED FEED', 'RAG EDGE', 'CARBON DUST', 'REDACTION'],
            flash: 'rgba(255,250,236,0.82)',
            emoji: ['\u{1F4C4}','\u270F\uFE0F','\u{1F4CE}','\u2702\uFE0F','\u{1F5DE}\uFE0F']
        },
        volcano: {
            banner: 'PYROCLASTIC SURGE',
            ticker: ['ASH CASCADE', 'LAVA BOMB IMPACT', 'EARTHQUAKE SHAKE', 'SULFUR CLOUD'],
            scraps: ['ASHFALL', 'VENT BLAST', 'PYROCLAST', 'SEISMIC HIT', 'MAGMA SPRAY'],
            flash: 'rgba(255,220,150,0.72)',
            emoji: ['\u{1F30B}','\u{1F525}','\u{1F4A5}','\u2668\uFE0F','\u{1FAA8}']
        }
    };
    function snonuxDetectThemeName() {
        // The shell sets window.SNONUX_CURRENT_THEME synchronously in <head>.
        // Falls back to the splash class for resilience if something raced.
        if (typeof window !== 'undefined' && window.SNONUX_CURRENT_THEME) {
            return window.SNONUX_CURRENT_THEME;
        }
        var splash = document.getElementById('splash-overlay');
        if (splash) {
            for (var i = 0; i < splash.classList.length; i++) {
                var cls = splash.classList[i];
                if (cls.indexOf('splash-') === 0 && cls !== 'splash-overlay' && cls.indexOf('splash--') !== 0) {
                    return cls.slice(7);
                }
            }
        }
        return 'neon';
    }

    // snonuxSwitchTheme persists the user's choice and reloads.
    // The shell's boot script picks it up on the next load.
    function snonuxSwitchTheme(theme) {
        var all = (typeof window !== 'undefined' && window.SNONUX_ALL_THEMES) || [];
        if (all.indexOf(theme) < 0) return;
        try { localStorage.setItem('snonuxTheme', theme); } catch (_) {}
        location.reload();
    }

    function snonuxRandomTheme() {
        var all = (typeof window !== 'undefined' && window.SNONUX_ALL_THEMES) || [];
        if (all.length <= 1) return null;
        var current = snonuxDetectThemeName();
        var pool = all.filter(function (t) { return t !== current; });
        return pool[Math.floor(Math.random() * pool.length)];
    }
    function snonuxEnsureWildRoot() {
        var root = document.getElementById('sno-wild-root');
        if (root) return root;
        root = document.createElement('div');
        root.id = 'sno-wild-root';
        root.setAttribute('aria-hidden', 'true');
        root.innerHTML =
            '<div id="sno-wild-colorwash" class="sno-wild-layer"></div>' +
            '<div id="sno-wild-rain" class="sno-wild-layer"></div>' +
            '<div id="sno-wild-wave" class="sno-wild-layer"></div>' +
            '<div id="sno-wild-beacon" class="sno-wild-layer"></div>' +
            '<div id="sno-wild-noise" class="sno-wild-layer"></div>' +
            '<div id="sno-wild-banner"></div>' +
            '<div id="sno-wild-ticker"><span></span></div>' +
            '<div id="sno-wild-scraps"></div>';
        document.body.appendChild(root);
        return root;
    }
    function snonuxApplyWildPreset(theme) {
        var preset = SNONUX_WILD_PRESETS[theme] || SNONUX_WILD_PRESETS.neon;
        var html = document.documentElement;
        var body = document.body;
        html.setAttribute('data-sno-theme', theme);
        body.setAttribute('data-sno-theme', theme);
        var root = snonuxEnsureWildRoot();
        var banner = root.querySelector('#sno-wild-banner');
        var ticker = root.querySelector('#sno-wild-ticker span');
        var scraps = root.querySelector('#sno-wild-scraps');
        banner.textContent = preset.banner;
        ticker.textContent = ' ' + preset.ticker.join('   //   ') + '   //   ' + preset.ticker.join('   //   ') + ' ';
        scraps.innerHTML = '';
        var phrases = preset.scraps || [];
        var count = theme === 'dos' || theme === 'terminal' || theme === 'matrix' ? 22 : 16;
        for (var i = 0; i < count; i++) {
            var span = document.createElement('span');
            span.textContent = phrases[i % phrases.length];
            span.style.setProperty('--x', (8 + Math.random() * 84).toFixed(2) + '%');
            span.style.setProperty('--y', (12 + Math.random() * 72).toFixed(2) + '%');
            span.style.setProperty('--rot', ((Math.random() * 36) - 18).toFixed(1) + 'deg');
            span.style.setProperty('--dx', ((Math.random() * 120) - 60).toFixed(1) + 'px');
            span.style.setProperty('--dy', ((Math.random() * 100) - 50).toFixed(1) + 'px');
            span.style.setProperty('--dur', (4.5 + Math.random() * 5.5).toFixed(2) + 's');
            span.style.setProperty('--delay', (-Math.random() * 6).toFixed(2) + 's');
            scraps.appendChild(span);
        }
        window._snonuxWildTheme = theme;
        window._snonuxWildFlashColor = preset.flash || 'rgba(255,255,255,0.7)';
    }
    function snonuxPulseFlash(color, duration) {
        var ov = document.createElement('div');
        ov.style.cssText = 'position:fixed;inset:0;z-index:9998;pointer-events:none;opacity:0;';
        document.body.appendChild(ov);
        var tone = color || window._snonuxWildFlashColor || 'rgba(255,255,255,0.7)';
        [0, 70, 150, 250].forEach(function(d, i) {
            setTimeout(function() {
                ov.style.transition = 'opacity 0.08s linear';
                ov.style.background = tone;
                ov.style.opacity = (i % 2 === 0) ? '0.82' : '0';
            }, d);
        });
        setTimeout(function() {
            ov.style.transition = 'opacity 0.25s linear';
            ov.style.opacity = '0';
        }, Math.max(180, duration || 320));
        setTimeout(function() { ov.remove(); }, Math.max(520, duration || 320) + 260);
    }
    function snonuxScheduleWildBursts() {
        clearTimeout(window._snonuxWildBurstTimer);
        if (!window._snoWildActive) return;
        var delay = 1400 + Math.random() * 3600;
        window._snonuxWildBurstTimer = setTimeout(function() {
            if (!window._snoWildActive) return;
            snonuxPulseFlash(window._snonuxWildFlashColor, 260);
            snonuxScheduleWildBursts();
        }, delay);
    }
    function snonuxSetWildState(on) {
        var body = document.body;
        var badge = document.getElementById('sno-wild-badge');
        body.classList.toggle('sno-wild-active', !!on);
        if (badge) badge.classList.toggle('sno-wild-on', !!on);
        if (on) {
            snonuxApplyWildPreset(window._snonuxWildTheme || snonuxDetectThemeName());
            snonuxScheduleWildBursts();
            body.classList.add('sno-wild-hue');
            snonuxStartFlyingEmoji();
            snonuxStartRandomFlips();
        } else {
            clearTimeout(window._snonuxWildBurstTimer);
            body.classList.remove('sno-wild-hue');
            snonuxStopFlyingEmoji();
            snonuxStopRandomFlips();
        }
    }
    // === WILD FLYING EMOJI ===
    function snonuxStartFlyingEmoji() {
        snonuxStopFlyingEmoji();
        var zone = document.getElementById('sno-flyzone');
        if (!zone) {
            zone = document.createElement('div');
            zone.id = 'sno-flyzone';
            zone.setAttribute('aria-hidden', 'true');
            document.body.appendChild(zone);
        }
        function spawn() {
            if (!window._snoWildActive) return;
            var preset = SNONUX_WILD_PRESETS[window._snonuxWildTheme] || SNONUX_WILD_PRESETS.neon;
            var emojis = preset.emoji || ['\u2B50'];
            var s = document.createElement('span');
            s.textContent = emojis[Math.floor(Math.random() * emojis.length)];
            var top = (5 + Math.random() * 80).toFixed(1);
            var dur = (3 + Math.random() * 4).toFixed(2);
            var dir = Math.random() > 0.5 ? 'sno-fly-lr' : 'sno-fly-rl';
            var wobble = ((Math.random() - 0.5) * 60).toFixed(0);
            var rot = (180 + Math.random() * 540).toFixed(0);
            s.style.setProperty('--ftop', top + '%');
            s.style.setProperty('--fy', wobble + 'px');
            s.style.setProperty('--frot', rot + 'deg');
            s.style.animationName = dir;
            s.style.animationDuration = dur + 's';
            zone.appendChild(s);
            setTimeout(function() { s.remove(); }, parseFloat(dur) * 1000 + 200);
            window._snoFlyTimer = setTimeout(spawn, 800 + Math.random() * 2200);
        }
        spawn();
    }
    function snonuxStopFlyingEmoji() {
        clearTimeout(window._snoFlyTimer);
        var zone = document.getElementById('sno-flyzone');
        if (zone) zone.innerHTML = '';
    }
    // === WILD RANDOM FLIPS ===
    function snonuxStartRandomFlips() {
        snonuxStopRandomFlips();
        function flip() {
            if (!window._snoWildActive) return;
            var allPosts = document.querySelectorAll('.post:not(.post-active)');
            if (allPosts.length > 0) {
                var p = allPosts[Math.floor(Math.random() * allPosts.length)];
                p.classList.add('sno-fx-flip');
                setTimeout(function() { p.classList.remove('sno-fx-flip'); }, 1500);
            }
            window._snoFlipTimer = setTimeout(flip, 2500 + Math.random() * 4000);
        }
        window._snoFlipTimer = setTimeout(flip, 1500 + Math.random() * 2000);
    }
    function snonuxStopRandomFlips() {
        clearTimeout(window._snoFlipTimer);
        var flipped = document.querySelectorAll('.sno-fx-flip');
        flipped.forEach(function(el) { el.classList.remove('sno-fx-flip'); });
    }
    (function snonuxWildSetup() {
        window._snoWildActive = !!window._snoWildActive;
        snonuxApplyWildPreset(snonuxDetectThemeName());
        snonuxSetWildState(window._snoWildActive);
    })();
    // Dramatic lightning flash on wild mode activation/deactivation
    function snonuxWildFlash(on) {
        var ov = document.createElement('div');
        ov.style.cssText = 'position:fixed;inset:0;z-index:9998;pointer-events:none;';
        document.body.appendChild(ov);
        if (on) {
            // Three rapid lightning flashes on activation
            [0, 80, 180, 300, 440].forEach(function(d, i) {
                setTimeout(function() {
                    ov.style.background = (i % 2 === 0) ? 'rgba(255,255,200,0.72)' : 'transparent';
                }, d);
            });
            setTimeout(function() { ov.remove(); }, 550);
            // Persistent storm overlay with intermittent flicker
            var storm = document.createElement('div');
            storm.id = 'sno-wild-storm';
            storm.style.cssText = 'position:fixed;inset:0;z-index:4998;pointer-events:none;' +
                'background:radial-gradient(ellipse at 50% 0%,rgba(255,255,200,0.88) 0%,transparent 55%);' +
                'animation:sno-wild-flicker 3.7s ease-in-out infinite;';
            document.body.appendChild(storm);
        } else {
            // Brief dark veil on deactivation
            ov.style.background = 'rgba(0,0,0,0.45)';
            ov.style.transition = 'opacity 0.45s';
            setTimeout(function() { ov.style.opacity = '0'; }, 60);
            setTimeout(function() { ov.remove(); }, 550);
            var storm = document.getElementById('sno-wild-storm');
            if (storm) { storm.style.transition = 'opacity 0.5s'; storm.style.opacity = '0'; setTimeout(function() { storm.remove(); }, 600); }
        }
    }
    function snonuxWaveType(w) {
        if (w === 'square') return 'square';
        if (w === 'triangle') return 'triangle';
        if (w === 'sawtooth') return 'sawtooth';
        return 'sine';
    }

    // === SHARED AMBIENT ENGINE ===
    // One AudioContext + one master gain node. Drones are long-running
    // oscillators; pulses are short scheduled oscillators. Noise is a looped
    // buffer source (no AudioWorklet). Fade in/out respects preset attack/release.
    // Ambient never routes through the one-shot UI sound paths.
    (function ambientEngine() {
        var ctx = null;
        var masterGain = null;
        var droneNodes = [];
        var pulseTimer = null;
        var noiseSrc = null;
        var noiseGainNode = null;
        var isPlaying = false;
        var isWild = false;
        var currentPreset = null;
        var melodyTimer = null;
        var melodyIndex = 0;

        function wildifyPreset(base) {
            if (!base) return base;
            var w = {};
            for (var k in base) w[k] = base[k];
            // Denser pulses for wild mode
            if (w.bpm) w.bpm = w.bpm * 1.5;
            if (w.pulseInterval) w.pulseInterval = w.pulseInterval / 1.5;
            // More noise texture
            if (w.noiseGain != null) w.noiseGain = Math.min(w.noiseGain * 1.6, 0.08);
            // Slightly higher gain, capped for safety
            if (w.gain != null) w.gain = Math.min(w.gain * 1.3, 0.15);
            // Deeper filter sweep for pulses
            w.filterFreq = (w.filterFreq || 700) * 1.6;
            w.filterQ = (w.filterQ || 0.8) * 2.0;
            return w;
        }

        function getPreset() {
            var ambient = SNONUX_SOUNDS.ambient;
            if (!ambient) return null;
            var base = isWild ? (ambient.wild || ambient.normal) : ambient.normal;
            return isWild ? wildifyPreset(base) : base;
        }

        function ensureCtx() {
            if (!ctx) {
                ctx = new (window.AudioContext || window.webkitAudioContext)();
                masterGain = ctx.createGain();
                masterGain.gain.value = 0;
                masterGain.connect(ctx.destination);
            }
            if (ctx.state === 'suspended') {
                ctx.resume().catch(function() {});
            }
            return ctx;
        }

        function stopDrones() {
            droneNodes.forEach(function(node) {
                try { node.stop(); node.disconnect(); } catch (_) {}
            });
            droneNodes = [];
        }

        function startDrones(preset) {
            var c = ensureCtx();
            var freqs = preset.droneFreqs || [];
            if (freqs.length === 0) return;
            var wt = snonuxWaveType(preset.wave);
            var detune = preset.detuneCents || 0;
            // Drones must respect the theme's ambient gain so they don't
            // drown the melody.  Four drones at 0.25 each = 1.0 total,
            // which clips and buries everything else.
            var gBase = preset.gain != null ? preset.gain : 0.08;
            var perOscGain = gBase / Math.max(1, freqs.length);
            freqs.forEach(function(freq) {
                var osc = c.createOscillator();
                var g = c.createGain();
                g.gain.value = perOscGain;
                osc.type = wt;
                osc.frequency.value = freq;
                if (detune) osc.detune.value = detune;
                osc.connect(g);
                g.connect(masterGain);
                osc.start();
                droneNodes.push(osc);
            });
        }

        function stopNoise() {
            if (noiseSrc) { try { noiseSrc.stop(); noiseSrc.disconnect(); } catch (_) {} noiseSrc = null; }
            if (noiseGainNode) { try { noiseGainNode.disconnect(); } catch (_) {} noiseGainNode = null; }
        }

        function startNoise(preset) {
            if (!preset.noiseGain) return;
            var c = ensureCtx();
            var bufferSize = 2 * c.sampleRate;
            var buffer = c.createBuffer(1, bufferSize, c.sampleRate);
            var data = buffer.getChannelData(0);
            for (var i = 0; i < bufferSize; i++) {
                data[i] = Math.random() * 2 - 1;
            }
            var src = c.createBufferSource();
            src.buffer = buffer;
            src.loop = true;
            noiseGainNode = c.createGain();
            // Keep noise subtle so it doesn't bury the melody.
            noiseGainNode.gain.value = preset.noiseGain;
            src.connect(noiseGainNode);
            noiseGainNode.connect(masterGain);
            src.start();
            noiseSrc = src;
        }

        function stopMelody() {
            if (melodyTimer) { clearTimeout(melodyTimer); melodyTimer = null; }
            melodyIndex = 0;
        }

        function playMelodyNote(note, preset) {
            var c = ensureCtx();
            var freq = note.freq || 440;
            var dur = note.dur || 0.3;
            if (freq <= 0 || dur <= 0) return;
            var wt = snonuxWaveType(preset.wave);
            // masterGain is now a binary "gate" (fades to 1 when on), so per-note
            // volume carries the real preset.gain. Without the old masterGain
            // squaring, this makes music audible at normal levels.
            var g = note.gain != null ? note.gain : (preset.gain != null ? preset.gain : 0.08);
            var pulseGain = Math.min(g, 0.08);

            // Build a richer timbre with two slightly detuned oscillators so
            // melody lines sound like actual notes instead of sterile blips.
            var osc1 = c.createOscillator();
            var osc2 = c.createOscillator();
            var mix = c.createGain();
            var gain = c.createGain();

            osc1.type = wt;
            osc1.frequency.value = freq;
            osc2.type = wt;
            osc2.frequency.value = freq;
            osc2.detune.value = 5 + Math.random() * 3; // subtle chorus

            mix.gain.value = 0.5;
            osc1.connect(mix);
            osc2.connect(mix);

            var lastNode = mix;
            if (preset.filterFreq) {
                var filter = c.createBiquadFilter();
                filter.type = 'lowpass';
                filter.frequency.value = preset.filterFreq;
                filter.Q.value = preset.filterQ || 1;
                var now = c.currentTime;
                filter.frequency.setValueAtTime(preset.filterFreq, now);
                filter.frequency.exponentialRampToValueAtTime(Math.max(preset.filterFreq * 0.2, 100), now + dur);
                lastNode.connect(filter);
                lastNode = filter;
            }

            lastNode.connect(gain);
            gain.connect(masterGain);
            var now = c.currentTime;

            // Proper ADSR-ish envelope: attack up, sustain, then release.
            var attack = Math.min(0.04, dur * 0.15);
            var release = Math.min(0.12, dur * 0.35);
            var sustainTime = Math.max(0, dur - attack - release);

            gain.gain.setValueAtTime(0, now);
            gain.gain.linearRampToValueAtTime(pulseGain, now + attack);
            if (sustainTime > 0) {
                gain.gain.setValueAtTime(pulseGain, now + attack + sustainTime);
            }
            gain.gain.exponentialRampToValueAtTime(0.001, now + dur);
            osc1.start(now);
            osc1.stop(now + dur + 0.02);
            osc2.start(now);
            osc2.stop(now + dur + 0.02);
        }

        function scheduleMelodyNote(preset) {
            if (!isPlaying || !preset || !preset.melody || preset.melody.length === 0) return;
            var note = preset.melody[melodyIndex];
            playMelodyNote(note, preset);
            melodyIndex = (melodyIndex + 1) % preset.melody.length;
            // Timing to the next note: explicit 'step' on the current note
            // controls exactly when the sequencer advances.  If step is
            // omitted the scheduler falls back to dur + a tiny BPM-scaled
            // gap so the line stays intelligible.
            var step;
            if (note.step != null && note.step > 0) {
                step = note.step;
            } else {
                var beat = preset.bpm ? 60.0 / preset.bpm : 0.5;
                step = note.dur + beat * 0.15;
            }
            melodyTimer = setTimeout(function() {
                if (!isPlaying) return;
                scheduleMelodyNote(currentPreset);
            }, step * 1000);
        }

        function schedulePulse() {
            if (!isPlaying || !currentPreset) return;
            var preset = currentPreset;
            if (preset.melody && preset.melody.length > 0) {
                scheduleMelodyNote(preset);
                return;
            }
            var interval = preset.pulseInterval;
            if (!interval && preset.bpm) {
                interval = 60.0 / preset.bpm;
            }
            if (!interval) interval = 1.0;
            var jitter = interval * (0.8 + Math.random() * 0.4);
            pulseTimer = setTimeout(function() {
                if (!isPlaying) return;
                playPulse(preset);
                schedulePulse();
            }, jitter * 1000);
        }

        function playPulse(preset) {
            var c = ensureCtx();
            var freqs = preset.pulseFreqs || [];
            if (freqs.length === 0) return;
            var freq = freqs[Math.floor(Math.random() * freqs.length)];
            var wt = snonuxWaveType(preset.wave);
            var attack = preset.attack != null ? preset.attack : 0.05;
            var release = preset.release != null ? preset.release : 0.3;
            var g = preset.gain != null ? preset.gain : 0.08;
            var pulseGain = g * 0.6;

            var osc = c.createOscillator();
            var gain = c.createGain();
            osc.type = wt;
            osc.frequency.value = freq;

            var lastNode = osc;
            if (preset.filterFreq) {
                var filter = c.createBiquadFilter();
                filter.type = 'lowpass';
                filter.frequency.value = preset.filterFreq;
                filter.Q.value = preset.filterQ || 1;
                var now = c.currentTime;
                filter.frequency.setValueAtTime(preset.filterFreq, now);
                filter.frequency.exponentialRampToValueAtTime(Math.max(preset.filterFreq * 0.2, 100), now + attack + release);
                lastNode.connect(filter);
                lastNode = filter;
            }

            lastNode.connect(gain);
            gain.connect(masterGain);
            var now = c.currentTime;
            gain.gain.setValueAtTime(0, now);
            gain.gain.linearRampToValueAtTime(pulseGain, now + attack);
            gain.gain.exponentialRampToValueAtTime(0.001, now + attack + release);
            osc.start(now);
            osc.stop(now + attack + release + 0.02);
        }

        function fadeMasterTo(target, duration) {
            if (!masterGain || !ctx) return;
            var dur = duration != null ? duration : 0.5;
            var now = ctx.currentTime;
            masterGain.gain.cancelScheduledValues(now);
            masterGain.gain.setValueAtTime(masterGain.gain.value, now);
            masterGain.gain.linearRampToValueAtTime(target, now + dur);
        }

        function stopAll() {
            clearTimeout(pulseTimer);
            pulseTimer = null;
            stopMelody();
            stopDrones();
            stopNoise();
        }

        function startEngine() {
            var preset = getPreset();
            if (!preset) return;
            currentPreset = preset;
            ensureCtx();

            // Begin scheduling only once the AudioContext is running.
            // Modern browsers start AudioContext in 'suspended' state; firing
            // note events before it is running schedules them into the void.
            function begin() {
                if (!isPlaying) return; // user toggled off while awaiting resume
                // Audio must actually be running or all scheduled notes are lost
                if (ctx.state !== 'running') {
                    isPlaying = false;
                    return;
                }
                melodyIndex = 0;
                startDrones(preset);
                startNoise(preset);
                schedulePulse();
                // masterGain is a binary gate: open it quickly (50 ms) so the
                // first scheduled notes are not swallowed by a long ramp.
                // preset.attack controls per-note envelopes, not the gate.
                fadeMasterTo(1.0, 0.05);
            }

            if (ctx.state === 'suspended') {
                isPlaying = true; // flag the intent so begin() can run later
                ctx.resume().then(begin).catch(function() {
                    // Resume failed (no autoplay permission).  Reset state so
                    // the next keypress triggers a fresh start attempt.
                    isPlaying = false;
                });
            } else {
                isPlaying = true;
                begin();
            }
        }

        function pauseEngine() {
            isPlaying = false;
            var fadeOut = currentPreset && currentPreset.release != null ? currentPreset.release : 0.5;
            fadeMasterTo(0, fadeOut);
            setTimeout(function() {
                if (!isPlaying) stopAll();
            }, 1000);
        }

        window.snonuxAmbientStart = function(reason) {
            if (isPlaying) return;
            startEngine();
        };

        window.snonuxAmbientPause = function(reason) {
            if (!isPlaying) return;
            pauseEngine();
        };

        window.snonuxAmbientToggle = function() {
            if (isPlaying) snonuxAmbientPause('toggle');
            else snonuxAmbientStart('toggle');
        };

        window.snonuxAmbientSetWild = function(on) {
            var wasWild = isWild;
            isWild = !!on;
            if (isPlaying && wasWild !== isWild) {
                var newPreset = getPreset();
                if (!newPreset) {
                    snonuxAmbientPause();
                    return;
                }
                fadeMasterTo(0, 0.3);
                setTimeout(function() {
                    if (!isPlaying) return;
                    stopAll();
                    currentPreset = newPreset;
                    startDrones(newPreset);
                    startNoise(newPreset);
                    schedulePulse();
                    var targetGain = 1.0;
                    fadeMasterTo(targetGain, 0.5);
                }, 350);
            }
        };

        window.snonuxAmbientSyncPreset = function() {
            if (!isPlaying) return;
            var preset = getPreset();
            if (!preset) {
                snonuxAmbientPause();
                return;
            }
            stopAll();
            currentPreset = preset;
            startDrones(preset);
            startNoise(preset);
            schedulePulse();
            var targetGain = preset.gain != null ? preset.gain : 0.08;
            fadeMasterTo(targetGain, 0.5);
        };

        window.snonuxAmbientIsPlaying = function() {
            return isPlaying;
        };
    })();

    (function splashSetup() {
        var el = document.getElementById('splash-overlay');
        if (!el) return;
        var splashAudioCtx = null;
        var splashChimePlayed = false;
        function playSplashChime() {
            if (splashChimePlayed) return;
            try {
                if (!splashAudioCtx) {
                    splashAudioCtx = new (window.AudioContext || window.webkitAudioContext)();
                }
                var ctx = splashAudioCtx;
                function ring() {
                    splashChimePlayed = true;
                    // Theme override: if defined, plays its own sound and returns true to skip the default chime.
                    if (window.snonuxSplashSound && window.snonuxSplashSound(ctx)) return;
                    var now = ctx.currentTime;
                    var sp = SNONUX_SOUNDS.splash;
                    var freqs = sp.freqs;
                    var spacing = sp.spacing != null ? sp.spacing : 0.075;
                    var gainAm = sp.gain != null ? sp.gain : 0.1;
                    var wave = snonuxWaveType(sp.wave);
                    var i, osc, g, t0;
                    for (i = 0; i < freqs.length; i++) {
                        osc = ctx.createOscillator();
                        g = ctx.createGain();
                        osc.connect(g);
                        g.connect(ctx.destination);
                        osc.type = wave;
                        osc.frequency.value = freqs[i];
                        t0 = now + i * spacing;
                        g.gain.setValueAtTime(0, t0);
                        g.gain.linearRampToValueAtTime(gainAm, t0 + 0.028);
                        g.gain.exponentialRampToValueAtTime(0.001, t0 + 0.52);
                        osc.start(t0);
                        osc.stop(t0 + 0.55);
                    }
                }
                ctx.resume().then(ring).catch(function() {});
            } catch (_) {}
        }
        function dismiss() {
            if (typeof splashDrift !== 'undefined') splashDrift.stop();
            if (el.classList.contains('splash--dismissed')) return;
            el.classList.add('splash--dismissed');
            el.setAttribute('aria-hidden', 'true');
        }
        function show() {
            if (typeof splashDrift !== 'undefined') splashDrift.reset();
            document.documentElement.classList.remove('sno-splash-skip');
            el.classList.remove('splash--dismissed');
            el.removeAttribute('aria-hidden');
            el.focus({ preventScroll: true });
        }
        function openSplashFromHeader(e) {
            if (e.target.closest('a')) return;
            e.preventDefault();
            if (typeof modalDrift !== 'undefined') modalDrift.stop();
            var modal = document.getElementById('post-modal');
            if (modal) modal.classList.remove('active');
            show();
        }
        function bindHeaderTriggers() {
            var triggers = document.querySelectorAll('.logo-mark, .logo-title h1, #sn-logo');
            triggers.forEach(function(trigger) {
                trigger.addEventListener('click', openSplashFromHeader);
            });
        }
        bindHeaderTriggers();
        // Exposed so the runtime theme-meta apply can re-bind after replacing
        // <header> innerHTML — the freshly-injected nodes have no listeners.
        window._snonuxRebindHeader = bindHeaderTriggers;
        window._snonuxDismissSplash = dismiss;
        window._snonuxShowSplash = show;
        window._snonuxPlaySplashChime = playSplashChime;
        if (document.documentElement.classList.contains('sno-splash-skip')) {
            dismiss();
            return;
        }
        playSplashChime();
        el.addEventListener('pointerdown', function() { playSplashChime(); }, { passive: true });
        el.addEventListener('click', function(e) { e.preventDefault(); dismiss(); });
        el.focus({ preventScroll: true });
    })();

    // === SPECIAL EFFECTS HELPERS ===
    function playBassDrop() {
        try {
            var ctx = new (window.AudioContext || window.webkitAudioContext)();
            var osc = ctx.createOscillator();
            var gain = ctx.createGain();
            osc.connect(gain); gain.connect(ctx.destination);
            osc.type = 'sawtooth';
            osc.frequency.setValueAtTime(130, ctx.currentTime);
            osc.frequency.exponentialRampToValueAtTime(28, ctx.currentTime + 0.55);
            gain.gain.setValueAtTime(0.22, ctx.currentTime);
            gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.9);
            osc.start(); osc.stop(ctx.currentTime + 0.95);
        } catch(_) {}
    }
    function snonuxParticleBurst(count, color) {
        var burst = document.createElement('div');
        burst.id = 'sno-burst';
        burst.setAttribute('aria-hidden', 'true');
        document.body.appendChild(burst);
        var cx = window.innerWidth / 2, cy = window.innerHeight / 2;
        for (var i = 0; i < count; i++) {
            var s = document.createElement('span');
            var angle = (i / count) * Math.PI * 2 + (Math.random() - 0.5) * 0.6;
            var dist = 60 + Math.random() * 280;
            s.style.left = cx + 'px';
            s.style.top = cy + 'px';
            s.style.setProperty('--px', (Math.cos(angle) * dist).toFixed(1) + 'px');
            s.style.setProperty('--py', (Math.sin(angle) * dist).toFixed(1) + 'px');
            s.style.setProperty('--pdur', (0.25 + Math.random() * 0.35).toFixed(2) + 's');
            s.style.setProperty('--pdel', (Math.random() * 0.08).toFixed(2) + 's');
            s.style.width = (3 + Math.random() * 6) + 'px';
            s.style.height = s.style.width;
            s.style.background = color || 'currentColor';
            burst.appendChild(s);
        }
        setTimeout(function() { burst.remove(); }, 900);
    }
    function snonuxKonamiExplosion() {
        document.querySelectorAll('.post').forEach(function(p) {
            p.classList.add('sno-fx-flip');
            setTimeout(function() { p.classList.remove('sno-fx-flip'); }, 1500);
        });
        var ov = document.createElement('div');
        ov.style.cssText = 'position:fixed;inset:0;z-index:9997;pointer-events:none;background:rgba(255,255,255,0.92);mix-blend-mode:difference;transition:opacity 0.25s';
        document.body.appendChild(ov);
        setTimeout(function() { ov.style.opacity='0'; setTimeout(function(){ov.remove();}, 300); }, 220);
        snonuxParticleBurst(72, '#fff');
        playBassDrop();
        var banner = document.createElement('div');
        banner.id = 'sno-konami-banner';
        banner.textContent = 'UNLOCKED';
        document.body.appendChild(banner);
        requestAnimationFrame(function() { banner.style.opacity='1'; });
        setTimeout(function() { banner.style.opacity='0'; setTimeout(function(){banner.remove();}, 350); }, 2200);
    }
    function snonuxCriticalOverload() {
        var ov = document.createElement('div');
        ov.style.cssText = 'position:fixed;inset:0;z-index:9997;pointer-events:none;background:#fff;mix-blend-mode:difference;';
        document.body.appendChild(ov);
        [0, 70, 140, 210].forEach(function(d,i){
            setTimeout(function(){ ov.style.opacity = (i%2===0)?'1':'0.25'; }, d);
        });
        setTimeout(function(){ ov.style.transition='opacity 0.3s'; ov.style.opacity='0'; setTimeout(function(){ov.remove();}, 350); }, 320);
        document.querySelectorAll('.post').forEach(function(p){
            p.classList.add('sno-fx-flip');
            setTimeout(function(){ p.classList.remove('sno-fx-flip'); }, 1500);
        });
        playBassDrop();
        snonuxParticleBurst(96, window._snonuxWildFlashColor || '#fff');
    }
    // Konami code tracker
    (function konamiSetup(){
        var seq = ['ArrowUp','ArrowUp','ArrowDown','ArrowDown','ArrowLeft','ArrowRight','ArrowLeft','ArrowRight','b','a'];
        var idx = 0;
        document.addEventListener('keydown', function(e){
            if (e.key === seq[idx]) { idx++; if (idx >= seq.length) { idx = 0; snonuxKonamiExplosion(); } }
            else { idx = (e.key === seq[0]) ? 1 : 0; }
        });
    })();
    // CRT overlay injection
    (function crtSetup(){
        if (document.getElementById('sno-crt-overlay')) return;
        var crt = document.createElement('div');
        crt.id = 'sno-crt-overlay';
        crt.innerHTML = '<div class="crt-scanlines"></div><div class="crt-flicker"></div>';
        crt.setAttribute('aria-hidden','true');
        document.body.appendChild(crt);
        var svg = document.createElementNS('http://www.w3.org/2000/svg','svg');
        svg.id = 'sno-crt-svg';
        svg.setAttribute('style','position:absolute;width:0;height:0;');
        svg.innerHTML = '<defs><filter id="sno-crt-distort"><feTurbulence type="fractalNoise" baseFrequency="0.012 0.006" numOctaves="2" result="noise"/><feDisplacementMap in="SourceGraphic" in2="noise" scale="3" xChannelSelector="R" yChannelSelector="G"/></filter></defs>';
        document.body.appendChild(svg);
    })();
    // Seasonal effects
    (function seasonalEffects(){
        var month = new Date().getMonth();
        var season = (month >= 11 || month <= 1) ? 'winter' : (month <= 4) ? 'spring' : (month <= 7) ? 'summer' : 'autumn';
        document.body.classList.add('sno-season-' + season);
    })();
    // === VAPORWAVE SUNSET (synthwave / retro only) ===
    (function vaporwaveSetup(){
        var theme = snonuxDetectThemeName();
        if (theme !== 'synthwave' && theme !== 'retro') return;
        var sunset = document.createElement('div');
        sunset.id = 'sno-vaporwave-sunset';
        sunset.setAttribute('aria-hidden','true');
        sunset.innerHTML = '<svg style="width:100%;height:100%;" preserveAspectRatio="none" viewBox="0 0 100 100"><defs><linearGradient id="vwg" x1="0" y1="0" x2="0" y2="1"><stop offset="0%" stop-color="#ff00cc" stop-opacity="0.35"/><stop offset="40%" stop-color="#ff00cc" stop-opacity="0.1"/><stop offset="100%" stop-color="#0a0a1a" stop-opacity="0"/></linearGradient></defs><ellipse cx="50" cy="95" rx="55" ry="18" fill="url(#vwg)"/><rect x="0" y="0" width="100" height="100" fill="none"/><g opacity="0.35">' + Array.from({length:12},function(_,i){return '<rect x="0" y="' + (75 + i*2.1) + '" width="100" height="' + (1.2 + i*0.2) + '" fill="#ff00cc" opacity="' + (0.5 - i*0.04).toFixed(2) + '"/>';}).join('') + '</g></svg>';
        document.body.appendChild(sunset);
        document.body.classList.add('sno-vaporwave-on');
    })();
    // === PARALLAX TILT ===
    (function parallaxSetup(){
        var enabled = true; // on by default
        document.body.classList.add('sno-parallax-on');
        window.snonuxToggleParallax = function(on){ enabled = on; document.body.classList.toggle('sno-parallax-on', on); };
        document.addEventListener('mousemove', function(e){
            if (!enabled) return;
            var x = (e.clientX / window.innerWidth - 0.5) * 2;
            var y = (e.clientY / window.innerHeight - 0.5) * 2;
            var ov = document.querySelector('.overlay');
            if (ov) {
                var rotY = (x * 4.5).toFixed(2);
                var rotX = (-y * 4.0).toFixed(2);
                ov.style.transform = 'rotateY(' + rotY + 'deg) rotateX(' + rotX + 'deg) scale(1.04)';
            }
            // Subtle depth per post
            document.querySelectorAll('.post').forEach(function(p, i){
                var depth = ((i % 3) - 1) * 12;
                p.style.transform = 'translateZ(' + depth + 'px)';
            });
        }, { passive: true });
    })();
    // === HACKER HOVER SCRAMBLE ===
    (function scrambleSetup(){
        var glyphs = '!<>-_\\/[]{}—=+*^?#________';
        function scramble(el){
            var original = el.getAttribute('data-scramble') || el.textContent;
            if (!el.getAttribute('data-scramble')) el.setAttribute('data-scramble', original);
            var len = original.length;
            var iter = 0;
            var interval = setInterval(function(){
                el.textContent = original.split('').map(function(ch, i){
                    if (i < iter) return original[i];
                    return glyphs[Math.floor(Math.random() * glyphs.length)];
                }).join('');
                iter += 1/3;
                if (iter >= len) { clearInterval(interval); el.textContent = original; }
            }, 30);
        }
        function attach(){
            document.querySelectorAll('.post-header strong, .logo-title h1, .splash-title').forEach(function(el){
                if (el._scrambleAttached) return;
                el._scrambleAttached = true;
                el.addEventListener('mouseenter', function(){ scramble(el); });
            });
        }
        attach();
        // re-attach after any dynamic content changes (not common here but safe)
        window.snonuxAttachScramble = attach;
    })();
    // === SCREENSHOT FLASH ===
    function snonuxScreenshotFlash(){
        var flash = document.createElement('div');
        flash.id = 'sno-shutter-flash';
        document.body.appendChild(flash);
        flash.style.transition = 'none';
        flash.style.opacity = '1';
        // Shutter sound
        try {
            var ctx = new (window.AudioContext || window.webkitAudioContext)();
            var osc = ctx.createOscillator();
            var gain = ctx.createGain();
            var noise = ctx.createScriptProcessor(4096, 1, 1);
            noise.onaudioprocess = function(e){
                var out = e.outputBuffer.getChannelData(0);
                for (var i = 0; i < out.length; i++) out[i] = (Math.random() * 2 - 1) * 0.15;
            };
            noise.connect(gain); gain.connect(ctx.destination);
            gain.gain.setValueAtTime(0.15, ctx.currentTime);
            gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.12);
            noise.connect(gain);
            osc.connect(gain);
            osc.type = 'square';
            osc.frequency.setValueAtTime(600, ctx.currentTime);
            osc.frequency.exponentialRampToValueAtTime(60, ctx.currentTime + 0.1);
            osc.start(); osc.stop(ctx.currentTime + 0.12);
            setTimeout(function(){ noise.disconnect(); }, 200);
        } catch(_) {}
        setTimeout(function(){
            flash.style.transition = 'opacity 0.35s ease';
            flash.style.opacity = '0';
            setTimeout(function(){ flash.remove(); }, 400);
        }, 80);
    }
    // === POST SCATTER ===
    function snonuxPostScatter(){
        var all = document.querySelectorAll('.post');
        all.forEach(function(p){
            var sx = ((Math.random() - 0.5) * window.innerWidth * 0.6).toFixed(1) + 'px';
            var sy = ((Math.random() - 0.5) * window.innerHeight * 0.6).toFixed(1) + 'px';
            var sr = ((Math.random() - 0.5) * 25).toFixed(1) + 'deg';
            p.style.setProperty('--sx', sx);
            p.style.setProperty('--sy', sy);
            p.style.setProperty('--sr', sr);
            p.classList.add('sno-scatter');
            setTimeout(function(){ p.classList.remove('sno-scatter'); }, 950);
        });
        playBassDrop();
    }
    // === RAINBOW SPARKLE TRAIL ===
    (function rainbowSparkle(){
        var throttle = 0;
        var hue = 0;
        document.addEventListener('pointermove', function(e){
            var now = Date.now();
            if (now - throttle < 45) return;
            throttle = now;
            var d = document.createElement('div');
            d.className = 'sno-sparkle';
            var size = 3 + Math.random() * 5;
            d.style.width = size + 'px';
            d.style.height = size + 'px';
            d.style.left = (e.clientX - size / 2 + (Math.random() - 0.5) * 10) + 'px';
            d.style.top = (e.clientY - size / 2 + (Math.random() - 0.5) * 10) + 'px';
            d.style.background = 'hsl(' + hue + ',90%,65%)';
            d.style.opacity = '0.85';
            d.style.animation = 'sno-sparkle ' + (0.35 + Math.random() * 0.25).toFixed(2) + 's ease-out forwards';
            document.body.appendChild(d);
            setTimeout(function() { d.remove(); }, 650);
            hue = (hue + 22) % 360;
        }, { passive: true });
    })();
    // === KEYBOARD NAVIGATION ===
    // j / ArrowDown  → next post       k / ArrowUp    → previous post
    // h / ArrowLeft  → previous page   l / ArrowRight → next page
    // PageUp/PageDown → scroll the post list; re-highlight post at top of visible area
    // Enter / click post → expand modal    Esc → close modal
    const posts = document.querySelectorAll('.post');
    let currentIndex = posts.length > 0 ? 0 : -1;
    var prevPageURL = (typeof window !== "undefined") ? (window.snonuxPrevPageURL || null) : null;
    var nextPageURL = (typeof window !== "undefined") ? (window.snonuxNextPageURL || null) : null;

    if (currentIndex >= 0) selectPost(0);

    function setActiveHighlight(index, playSound, scrollIntoView) {
        if (posts.length === 0) return;
        var prevIdx = currentIndex;
        if (currentIndex >= 0) posts[currentIndex].classList.remove('post-active');
        currentIndex = Math.max(0, Math.min(index, posts.length - 1));
        posts[currentIndex].classList.add('post-active');
        if (prevIdx >= 0 && prevIdx !== currentIndex && posts[prevIdx]) {
            var ghost = posts[prevIdx].querySelector('.sno-afterimage');
            if (!ghost) {
                ghost = document.createElement('div');
                ghost.className = 'sno-afterimage';
                posts[prevIdx].appendChild(ghost);
            }
            ghost.classList.remove('sno-afterimage-active');
            void ghost.offsetWidth;
            ghost.classList.add('sno-afterimage-active');
        }
        if (scrollIntoView) {
            posts[currentIndex].scrollIntoView({ behavior: 'smooth', block: 'nearest' });
        }
        if (playSound) playNavSound();
    }

    function selectPost(index, direction) {
        setActiveHighlight(index, true, true);
        if (direction && posts[currentIndex]) {
            var post = posts[currentIndex];
            post.classList.remove('sno-enter-down', 'sno-enter-up');
            void post.offsetWidth;
            post.classList.add(direction === 'down' ? 'sno-enter-down' : 'sno-enter-up');
            setTimeout(function() { post.classList.remove('sno-enter-down', 'sno-enter-up'); }, 320);
        }
        if (window.snonuxNavEffect) window.snonuxNavEffect();
    }

    /** Pick the post that should be active for the current viewport (anchor near top of visible area). */
    function activeIndexForVisibleRegion(sc) {
        if (posts.length === 0) return -1;
        var scrTop, scrBot, anchorY;
        if (sc) {
            var scr = sc.getBoundingClientRect();
            scrTop = scr.top;
            scrBot = scr.bottom;
            anchorY = scr.top + Math.min(scr.height * 0.18, 100);
        } else {
            scrTop = 0;
            scrBot = window.innerHeight;
            anchorY = window.innerHeight * 0.15;
        }
        var i, pr;
        for (i = 0; i < posts.length; i++) {
            pr = posts[i].getBoundingClientRect();
            if (pr.top <= anchorY && anchorY < pr.bottom) return i;
        }
        for (i = 0; i < posts.length; i++) {
            pr = posts[i].getBoundingClientRect();
            if (pr.bottom > scrTop && pr.top < scrBot) return i;
        }
        return posts.length - 1;
    }

    function playNavSound() {
        try {
            var n = SNONUX_SOUNDS.nav;
            const ctx = new (window.AudioContext || window.webkitAudioContext)();
            const osc = ctx.createOscillator();
            const gain = ctx.createGain();
            osc.connect(gain); gain.connect(ctx.destination);
            osc.frequency.value = n.freq;
            osc.type = snonuxWaveType(n.wave);
            var dur = n.dur != null ? n.dur : 0.08;
            var g = n.gain != null ? n.gain : 0.12;
            gain.gain.setValueAtTime(g, ctx.currentTime);
            gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + dur);
            osc.start(ctx.currentTime); osc.stop(ctx.currentTime + dur + 0.02);
        } catch (_) {}
    }

    function playOpenSound() {
        try {
            var o = SNONUX_SOUNDS.open;
            const ctx = new (window.AudioContext || window.webkitAudioContext)();
            const osc = ctx.createOscillator();
            const gain = ctx.createGain();
            osc.connect(gain); gain.connect(ctx.destination);
            osc.type = snonuxWaveType(o.wave);
            var dur = o.dur != null ? o.dur : 0.14;
            var g = o.gain != null ? o.gain : 0.1;
            osc.frequency.setValueAtTime(o.start, ctx.currentTime);
            osc.frequency.exponentialRampToValueAtTime(o.end, ctx.currentTime + dur);
            gain.gain.setValueAtTime(g, ctx.currentTime);
            gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + dur + 0.06);
            osc.start(ctx.currentTime); osc.stop(ctx.currentTime + dur + 0.07);
        } catch (_) {}
    }

    function playCloseSound() {
        try {
            var c = SNONUX_SOUNDS.close;
            const ctx = new (window.AudioContext || window.webkitAudioContext)();
            const osc = ctx.createOscillator();
            const gain = ctx.createGain();
            osc.connect(gain); gain.connect(ctx.destination);
            osc.type = snonuxWaveType(c.wave);
            var dur = c.dur != null ? c.dur : 0.15;
            var g = c.gain != null ? c.gain : 0.1;
            osc.frequency.setValueAtTime(c.start, ctx.currentTime);
            osc.frequency.exponentialRampToValueAtTime(c.end, ctx.currentTime + dur);
            gain.gain.setValueAtTime(g, ctx.currentTime);
            gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + dur + 0.05);
            osc.start(ctx.currentTime); osc.stop(ctx.currentTime + dur + 0.06);
        } catch (_) {}
    }

    function playBounceSound() {
        try {
            var b = SNONUX_SOUNDS.bounce;
            const ctx = new (window.AudioContext || window.webkitAudioContext)();
            const osc = ctx.createOscillator();
            const gain = ctx.createGain();
            osc.connect(gain); gain.connect(ctx.destination);
            osc.type = snonuxWaveType(b.wave);
            var dur = b.dur != null ? b.dur : 0.12;
            var g = b.gain != null ? b.gain : 0.1;
            osc.frequency.setValueAtTime(b.start, ctx.currentTime);
            osc.frequency.exponentialRampToValueAtTime(b.end, ctx.currentTime + dur);
            gain.gain.setValueAtTime(g, ctx.currentTime);
            gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + dur + 0.05);
            osc.start(ctx.currentTime); osc.stop(ctx.currentTime + dur + 0.06);
        } catch (_) {}
    }

    var _snoBounceCls = ['sno-fx-bounce-left','sno-fx-bounce-right','sno-fx-bounce-left-wild','sno-fx-bounce-right-wild',
        'sno-fx-bounce-up','sno-fx-bounce-down','sno-fx-bounce-up-wild','sno-fx-bounce-down-wild'];

    function bounceEffect(dir) {
        var ov = document.querySelector('.overlay');
        if (!ov) return;
        var wild = !!window._snoWildActive;
        var map = { left: wild ? 'sno-fx-bounce-left-wild' : 'sno-fx-bounce-left',
                    right: wild ? 'sno-fx-bounce-right-wild' : 'sno-fx-bounce-right',
                    up: wild ? 'sno-fx-bounce-up-wild' : 'sno-fx-bounce-up',
                    down: wild ? 'sno-fx-bounce-down-wild' : 'sno-fx-bounce-down' };
        var cls = map[dir] || map.down;
        _snoBounceCls.forEach(function(c) { ov.classList.remove(c); });
        void ov.offsetWidth;
        ov.classList.add(cls);
        var dur = wild ? 540 : 380;
        setTimeout(function() { ov.classList.remove(cls); }, dur);
        playBounceSound();
        if (wild) snonuxPulseFlash(window._snonuxWildFlashColor, 200);
    }

    // === DRIFT PHYSICS — reusable controller for floating panels ===
    function makeDriftController(getEl, opts) {
        var x = 0, y = 0, vx = 0, vy = 0, raf = null;
        var PUSH = opts.push || 12;
        var FRICTION = opts.friction || 0.92;
        var BOUNCE_DAMP = opts.bounceDamp || 0.5;
        var STOP_THRESHOLD = opts.stopThreshold || 0.3;

        function clampAndBounce() {
            var el = getEl();
            if (!el) return;
            var w = el.offsetWidth, h = el.offsetHeight;
            var maxX = (window.innerWidth - w) / 2;
            var maxY = (window.innerHeight - h) / 2;
            if (maxX < 0) maxX = window.innerWidth * 0.3;
            if (maxY < 0) maxY = window.innerHeight * 0.3;
            var hit = false;
            if (x > maxX) { x = maxX; vx = -vx * BOUNCE_DAMP; hit = true; }
            if (x < -maxX) { x = -maxX; vx = -vx * BOUNCE_DAMP; hit = true; }
            if (y > maxY) { y = maxY; vy = -vy * BOUNCE_DAMP; hit = true; }
            if (y < -maxY) { y = -maxY; vy = -vy * BOUNCE_DAMP; hit = true; }
            if (hit && opts.onBounce) opts.onBounce(el, x, y, vx, vy);
        }

        function tick() {
            vx *= FRICTION;
            vy *= FRICTION;
            x += vx;
            y += vy;
            clampAndBounce();
            var el = getEl();
            if (el) {
                if (opts.applyTransform) opts.applyTransform(el, x, y, vx, vy);
                else el.style.transform = 'translate(' + x.toFixed(1) + 'px,' + y.toFixed(1) + 'px)';
            }
            if (Math.abs(vx) > STOP_THRESHOLD || Math.abs(vy) > STOP_THRESHOLD) {
                raf = requestAnimationFrame(tick);
            } else {
                raf = null;
            }
        }

        function ensureLoop() {
            if (!raf) raf = requestAnimationFrame(tick);
        }

        return {
            keyPush: function(e) {
                var dx = 0, dy = 0;
                switch (e.key) {
                    case 'h': case 'ArrowLeft':  dx = -PUSH; break;
                    case 'l': case 'ArrowRight': dx = PUSH;  break;
                    case 'k': case 'ArrowUp':    dy = -PUSH; break;
                    case 'j': case 'ArrowDown':  dy = PUSH;  break;
                    default: return false;
                }
                e.preventDefault();
                vx += dx;
                vy += dy;
                ensureLoop();
                return true;
            },
            kick: function(dx, dy) { vx += (dx || 0); vy += (dy || 0); ensureLoop(); },
            reset: function() {
                x = 0; y = 0; vx = 0; vy = 0;
                var el = getEl();
                if (el) el.style.transform = '';
                if (raf) { cancelAnimationFrame(raf); raf = null; }
            },
            stop: function() {
                if (raf) { cancelAnimationFrame(raf); raf = null; }
                var el = getEl();
                if (el) el.style.transform = '';
                x = 0; y = 0; vx = 0; vy = 0;
            }
        };
    }

    // === MODAL DRIFT — arrow/hjkl push the modal around with momentum ===
    var modalDrift = makeDriftController(
        function() { return document.querySelector('#post-modal .modal-inner'); },
        { push: 12, friction: 0.92, bounceDamp: 0.5, stopThreshold: 0.3 }
    );

    // === SPLASH DRIFT — same physics on the splash panel with velocity tilt ===
    var splashDrift = makeDriftController(
        function() { return document.querySelector('#splash-overlay .splash-inner'); },
        {
            push: 14,
            friction: 0.93,
            bounceDamp: 0.45,
            stopThreshold: 0.25,
            applyTransform: function(el, x, y, vx) {
                var rot = Math.max(-5, Math.min(5, vx * 0.15));
                el.style.transform = 'translate(' + x.toFixed(1) + 'px,' + y.toFixed(1) + 'px) rotate(' + rot.toFixed(2) + 'deg)';
            },
            onBounce: function() {
                playBounceSound();
                if (window._snoWildActive) snonuxPulseFlash(window._snonuxWildFlashColor, 180);
            }
        }
    );

    function getFxButton(name) {
        return document.querySelector('.nav-fx-button[data-sno-fx="' + name + '"]');
    }

    function pulseFxButton(name) {
        var button = getFxButton(name);
        if (!button) return;
        button.classList.remove('sno-fx-triggered');
        void button.offsetWidth;
        button.classList.add('sno-fx-triggered');
        setTimeout(function() { button.classList.remove('sno-fx-triggered'); }, 180);
    }

    function syncFxButtonStates() {
        var wildButton = getFxButton('wild');
        if (wildButton) wildButton.setAttribute('aria-pressed', window._snoWildActive ? 'true' : 'false');
        var crtButton = getFxButton('crt');
        if (crtButton) crtButton.setAttribute('aria-pressed', document.body.classList.contains('sno-crt-on') ? 'true' : 'false');
        var ghostButton = getFxButton('ghost');
        if (ghostButton) ghostButton.setAttribute('aria-pressed', document.body.classList.contains('sno-ghost-mode') ? 'true' : 'false');
        var ambientButton = getFxButton('ambient');
        if (ambientButton) ambientButton.setAttribute('aria-pressed', (window.snonuxAmbientIsPlaying && window.snonuxAmbientIsPlaying()) ? 'true' : 'false');
    }

    function setWildMode(on, opts) {
        var wasPlaying = window.snonuxAmbientIsPlaying && window.snonuxAmbientIsPlaying();
        window._snoWildActive = !!on;
        if (window.snonuxWildToggle) window.snonuxWildToggle();
        snonuxSetWildState(window._snoWildActive);
        if (opts && opts.splashMode) snonuxWildFlash(window._snoWildActive);
        else if (window._snoWildActive) snonuxCriticalOverload();
        else snonuxWildFlash(false);
        if (opts && opts.kickSplash && splashDrift) {
            splashDrift.kick((Math.random() - 0.5) * 24, -10 - Math.random() * 10);
        }
        pulseFxButton('wild');
        syncFxButtonStates();
        // Ambient follows wild state: switch preset and start wild ambient on activation.
        // If ambient was not playing before wild mode, mark it as wild-forced so we can
        // stop it on deactivation unless the user opted in during the session.
        if (window._snoWildActive) {
            if (window.snonuxAmbientSetWild) window.snonuxAmbientSetWild(window._snoWildActive);
            if (!wasPlaying && window.snonuxAmbientStart) {
                window._snonuxAmbientWildForced = true;
                window.snonuxAmbientStart('wild');
            }
        } else {
            // If ambient was wild-forced and the user never opted in, fade to silence.
            if (window._snonuxAmbientWildForced && !snonuxAmbientLoadPreference()) {
                if (window.snonuxAmbientPause) window.snonuxAmbientPause('wild-off');
            }
            window._snonuxAmbientWildForced = false;
            if (window.snonuxAmbientSetWild) window.snonuxAmbientSetWild(window._snoWildActive);
        }
        // Re-sync button states after ambient may have started or stopped.
        syncFxButtonStates();
    }

    function toggleWildMode(opts) {
        setWildMode(!window._snoWildActive, opts);
    }

    function toggleCrtMode() {
        document.body.classList.toggle('sno-crt-on');
        pulseFxButton('crt');
        syncFxButtonStates();
    }

    function toggleGhostMode() {
        document.body.classList.toggle('sno-ghost-mode');
        pulseFxButton('ghost');
        syncFxButtonStates();
    }

    function snonuxAmbientSavePreference(enabled) {
        try { localStorage.setItem('snonuxAmbientEnabled', enabled ? '1' : '0'); } catch (_) {}
    }
    function snonuxAmbientLoadPreference() {
        try { return localStorage.getItem('snonuxAmbientEnabled') === '1'; } catch (_) { return false; }
    }

    function toggleAmbientMode() {
        if (window.snonuxAmbientToggle) window.snonuxAmbientToggle();
        pulseFxButton('ambient');
        syncFxButtonStates();
        snonuxAmbientSavePreference(window.snonuxAmbientIsPlaying && window.snonuxAmbientIsPlaying());
    }

    function triggerFlashEffect() {
        snonuxScreenshotFlash();
        pulseFxButton('flash');
    }

    function triggerScatterEffect() {
        snonuxPostScatter();
        pulseFxButton('scatter');
    }

    (function setupNavFxButtons() {
        var fxHandlers = {
            wild: function() { toggleWildMode(); },
            crt: toggleCrtMode,
            ghost: toggleGhostMode,
            ambient: toggleAmbientMode,
            flash: triggerFlashEffect,
            scatter: triggerScatterEffect
        };
        document.querySelectorAll('.nav-fx-button').forEach(function(button) {
            button.addEventListener('click', function(e) {
                e.preventDefault();
                var fx = button.getAttribute('data-sno-fx');
                if (fxHandlers[fx]) fxHandlers[fx]();
            });
        });
        syncFxButtonStates();
    })();

    // Restore ambient preference on load (opt-in; default off).
    (function initAmbientPreference() {
        if (snonuxAmbientLoadPreference() && window.snonuxAmbientStart) {
            window.snonuxAmbientStart('restore');
        }
        syncFxButtonStates();
    })();

    // Inject keyboard controls hint into splash overlay (all themes)
    (function enhanceSplashHint() {
        var hint = document.querySelector('#splash-overlay .splash-hint');
        if (!hint || document.querySelector('#splash-overlay .splash-controls')) return;
        var extra = document.createElement('div');
        extra.className = 'splash-controls';
        extra.innerHTML = '<kbd>↑</kbd><kbd>↓</kbd><kbd>←</kbd><kbd>→</kbd> drift \u2022 <kbd>w</kbd> wild \u2022 <kbd>p</kbd> music \u2022 <kbd>Enter</kbd> open';
        hint.appendChild(extra);
    })();

    function openPostAt(index, scrollIntoView) {
        if (posts.length === 0) return;
        setActiveHighlight(index, false, !!scrollIntoView);
        var post = posts[currentIndex];
        var postText = post ? post.querySelector('.post-text') : null;
        if (!postText) return;
        var modal = document.getElementById('post-modal');
        var modalInner = modal ? modal.querySelector('.modal-inner') : null;
        document.getElementById('modal-content').innerHTML = postText.innerHTML;
        modal.classList.add('active');
        modalDrift.reset();
        if (window.snonuxOpenEffect) window.snonuxOpenEffect(post);
        modal.scrollTop = 0;
        if (modalInner) {
            modalInner.scrollTop = 0;
            requestAnimationFrame(function() {
                modalInner.scrollIntoView({ block: 'center', inline: 'nearest' });
            });
        }
        playOpenSound();
    }

    function closeModal() {
        modalDrift.stop();
        document.getElementById('post-modal').classList.remove('active');
        playCloseSound();
        if (window.snonuxCloseEffect) window.snonuxCloseEffect();
    }

    (function postClickOpen() {
        posts.forEach(function(post, idx) {
            post.addEventListener('click', function(e) {
            if (e.target.closest('a, button, audio, video, input, textarea, select, label')) return;
                openPostAt(idx, true);
            });
        });
    })();

    (function deepLinkFromHash() {
        var h = location.hash;
        if (!h || h.indexOf('#post-') !== 0) return;
        var id = decodeURIComponent(h.slice(6));
        var el = document.getElementById('post-' + id);
        if (!el) return;
        var idx = parseInt(el.getAttribute('data-index'), 10);
        if (isNaN(idx)) return;
        openPostAt(idx, true);
    })();

    document.addEventListener('keydown', function(e) {
        var tag = e.target.tagName;
        // Skip when typing into a form control. SELECT counts: native typeahead
        // (e.g. matching options by first letter) would otherwise be hijacked
        // by 'w', 'c', 't', etc. shortcuts.
        if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') return;
        var splash = document.getElementById('splash-overlay');
        if (splash && !splash.classList.contains('splash--dismissed')) {
            if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                if (window._snonuxPlaySplashChime) window._snonuxPlaySplashChime();
                if (window._snonuxDismissSplash) window._snonuxDismissSplash();
            } else if (e.key === 'Escape') {
                e.preventDefault();
                if (window._snonuxDismissSplash) window._snonuxDismissSplash();
            } else if (e.key === 'w' && !e.repeat) {
                e.preventDefault();
                toggleWildMode({ splashMode: true, kickSplash: true });
            } else if (e.key === 'p' && !e.repeat) {
                e.preventDefault();
                toggleAmbientMode();
            } else if (e.key === 'f' && !e.repeat) {
                e.preventDefault();
                triggerFlashEffect();
            } else if (e.key === 'c' && !e.repeat) {
                e.preventDefault();
                toggleCrtMode();
            } else if (e.key === 'g' && !e.repeat) {
                e.preventDefault();
                toggleGhostMode();
            } else if (e.key === 't' && !e.repeat) {
                e.preventDefault();
                var pick = snonuxRandomTheme();
                if (pick) snonuxSwitchTheme(pick);
            } else if (splashDrift.keyPush(e)) {
                playNavSound();
            }
            return;
        }
        if (document.getElementById('post-modal').classList.contains('active')) {
            if (e.key === 'Escape') { closeModal(); e.preventDefault(); }
            else if (e.key === 'p' && !e.repeat) {
                e.preventDefault();
                toggleAmbientMode();
            }
            else if (e.key === 'f' && !e.repeat) {
                e.preventDefault();
                triggerFlashEffect();
            }
            else if (e.key === 't' && !e.repeat) {
                e.preventDefault();
                var pick = snonuxRandomTheme();
                if (pick) snonuxSwitchTheme(pick);
            }
            else if (modalDrift.keyPush(e)) { playNavSound(); }
            return;
        }
        switch (e.key) {
            case 'PageUp':
            case 'PageDown': {
                var sc = document.getElementById('post-content');
                var step = (sc && sc.clientHeight) ? sc.clientHeight : window.innerHeight;
                var dy = (e.key === 'PageUp') ? -step : step;
                if (sc) {
                    sc.scrollTop += dy;
                } else {
                    window.scrollBy(0, dy);
                }
                var idx = activeIndexForVisibleRegion(sc);
                if (idx >= 0) setActiveHighlight(idx, true, false);
                if (window.snonuxScrollEffect) window.snonuxScrollEffect(e.key === 'PageUp' ? 'up' : 'down');
                e.preventDefault();
                break;
            }
            case 'j': case 'ArrowDown':
                if (currentIndex >= posts.length - 1) { bounceEffect('down'); }
                else { selectPost(currentIndex + 1, 'down'); }
                e.preventDefault(); break;
            case 'k': case 'ArrowUp':
                if (currentIndex <= 0) { bounceEffect('up'); }
                else { selectPost(currentIndex - 1, 'up'); }
                e.preventDefault(); break;
            case 'h': case 'ArrowLeft':
                if (prevPageURL) { playNavSound(); if (window.snonuxPageEffect) window.snonuxPageEffect(); window.location.href = prevPageURL; }
                else { bounceEffect('left'); }
                e.preventDefault(); break;
            case 'l': case 'ArrowRight':
                if (nextPageURL) { playNavSound(); if (window.snonuxPageEffect) window.snonuxPageEffect(); window.location.href = nextPageURL; }
                else { bounceEffect('right'); }
                e.preventDefault(); break;
            case 'Enter': openPostAt(currentIndex, true); e.preventDefault(); break;
            case 'w': {
                toggleWildMode();
                e.preventDefault(); break;
            }
            case 'p': {
                toggleAmbientMode();
                e.preventDefault(); break;
            }
            case 'c':
                toggleCrtMode();
                e.preventDefault(); break;
            case 'g':
                toggleGhostMode();
                e.preventDefault(); break;
            case 'f':
                triggerFlashEffect();
                e.preventDefault(); break;
            case 'x':
                triggerScatterEffect();
                e.preventDefault(); break;
            case 't': {
                var pick = snonuxRandomTheme();
                if (pick) snonuxSwitchTheme(pick);
                e.preventDefault(); break;
            }
        }
    });

    // === MODAL SCROLL-END INDICATOR ===
    (function modalScrollEnd() {
        var mi = document.querySelector('#post-modal .modal-inner');
        if (!mi) return;
        mi.addEventListener('scroll', function() {
            var atEnd = mi.scrollHeight - mi.scrollTop - mi.clientHeight < 4;
            var el = document.getElementById('sno-scroll-end');
            if (!el) return;
            if (atEnd) {
                el.classList.remove('sno-scroll-end-active');
                void el.offsetWidth;
                el.classList.add('sno-scroll-end-active');
            }
        }, { passive: true });
    })();

    // === IDLE BREATHING ===
    (function idleBreathe() {
        var timer = null;
        var IDLE_DELAY = 10000;
        function startBreathe() {
            stopBreathe();
            timer = setTimeout(function() {
                if (currentIndex >= 0 && posts[currentIndex]) {
                    posts[currentIndex].classList.add('sno-idle-breathe');
                }
            }, IDLE_DELAY);
        }
        function stopBreathe() {
            clearTimeout(timer);
            for (var i = 0; i < posts.length; i++) {
                posts[i].classList.remove('sno-idle-breathe');
            }
        }
        function resetIdle() { stopBreathe(); startBreathe(); }
        document.addEventListener('keydown', resetIdle);
        document.addEventListener('pointermove', resetIdle, { passive: true });
        document.addEventListener('pointerdown', resetIdle, { passive: true });
        startBreathe();
    })();

    // === FIRST-VISIT PARTICLE BURST ===
    (function firstVisitBurst() {
        var key = 'sno-visited';
        try { if (sessionStorage.getItem(key)) return; sessionStorage.setItem(key, '1'); } catch (_) { return; }
        if (document.documentElement.classList.contains('sno-splash-skip')) return;
        var origDismiss = window._snonuxDismissSplash;
        if (!origDismiss) return;
        window._snonuxDismissSplash = function() {
            origDismiss();
            var burst = document.createElement('div');
            burst.id = 'sno-burst';
            burst.setAttribute('aria-hidden', 'true');
            document.body.appendChild(burst);
            var cx = window.innerWidth / 2, cy = window.innerHeight / 2;
            for (var i = 0; i < 36; i++) {
                var s = document.createElement('span');
                var angle = (i / 36) * Math.PI * 2 + (Math.random() - 0.5) * 0.4;
                var dist = 80 + Math.random() * 180;
                s.style.left = cx + 'px';
                s.style.top = cy + 'px';
                s.style.setProperty('--px', (Math.cos(angle) * dist).toFixed(1) + 'px');
                s.style.setProperty('--py', (Math.sin(angle) * dist).toFixed(1) + 'px');
                s.style.setProperty('--pdur', (0.4 + Math.random() * 0.5).toFixed(2) + 's');
                s.style.setProperty('--pdel', (Math.random() * 0.12).toFixed(2) + 's');
                s.style.width = (4 + Math.random() * 5) + 'px';
                s.style.height = s.style.width;
                burst.appendChild(s);
            }
            setTimeout(function() { burst.remove(); }, 1200);
        };
    })();

    // === THEME APPLICATION (runtime) ===
    // The shell renders header/splash/title for SNONUX_DEFAULT_THEME. If the
    // user has saved a different theme, fetch its meta.json and swap the
    // theme-specific markup. Once the splash overlay is final we load theme.js
    // — its splash WebGL initialiser must bind to the canvas it will animate.
    (function snonuxApplyThemeMeta() {
        if (typeof window === 'undefined') return;
        var current = window.SNONUX_CURRENT_THEME;
        var def = window.SNONUX_DEFAULT_THEME;
        if (!current) return;

        var splashOverlay = document.getElementById('splash-overlay');
        if (splashOverlay) {
            // theme.css selectors target via [data-sno-theme] (set by the
            // boot script on <html>); the legacy splash-<theme> class on the
            // overlay is added defensively for any rule still keyed off it.
            splashOverlay.classList.add('splash-' + current);
        }

        function applyMeta(m) {
            if (m.title) document.title = m.title;
            var headerEl = document.querySelector('header');
            if (headerEl && m.header_html) headerEl.innerHTML = m.header_html;
            if (splashOverlay && m.splash_inner_html) splashOverlay.innerHTML = m.splash_inner_html;
            var prevA = document.getElementById('sno-prev-page');
            if (prevA && m.prev_page_text) prevA.innerHTML = m.prev_page_text;
            var nextA = document.getElementById('sno-next-page');
            if (nextA && m.next_page_text) nextA.innerHTML = m.next_page_text;
            if (typeof window._snonuxRebindHeader === 'function') window._snonuxRebindHeader();
        }

        function loadThemeJS() {
            var s = document.createElement('script');
            s.src = 'themes/' + current + '/theme.js';
            document.head.appendChild(s);
        }

        if (current === def) {
            // Baked markup is already correct. Load theme.js straight away.
            loadThemeJS();
            return;
        }

        // Switched theme: swap markup + sounds first, *then* load theme.js so
        // its splash WebGL attaches to the final canvas.
        var pending = 2;
        function done() { if (--pending === 0) loadThemeJS(); }
        fetch('themes/' + current + '/meta.json')
            .then(function (r) { return r.json(); })
            .then(applyMeta)
            .catch(function () {})
            .finally(done);
        fetch('themes/' + current + '/sounds.json')
            .then(function (r) { return r.json(); })
            .then(function (s) { window.SNONUX_SOUNDS = s; SNONUX_SOUNDS = s; if (window.snonuxAmbientSyncPreset) window.snonuxAmbientSyncPreset(); })
            .catch(function () {})
            .finally(done);
    })();

    // === THEME DROPDOWN (populates select#sno-theme-select) ===
    (function snonuxThemeDropdown() {
        var sel = document.getElementById('sno-theme-select');
        if (!sel) return;
        var all = (typeof window !== 'undefined' && window.SNONUX_ALL_THEMES) || [];
        var current = snonuxDetectThemeName();
        sel.innerHTML = '';
        for (var i = 0; i < all.length; i++) {
            var name = all[i];
            var opt = document.createElement('option');
            opt.value = name;
            opt.textContent = name;
            if (name === current) opt.selected = true;
            sel.appendChild(opt);
        }
        sel.addEventListener('change', function () {
            snonuxSwitchTheme(sel.value);
        });
    })();

