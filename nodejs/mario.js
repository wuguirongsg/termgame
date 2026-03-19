#!/usr/bin/env node
/**
 * ╔══════════════════════════════════════════════════════╗
 * ║        FC TERMINAL MARIO  -  v1.0  (Node.js)        ║
 * ║   NES风格超级玛丽 · 完整终端游戏                      ║
 * ║   Controls: A/← D/→ Move  SPACE Jump  Q Quit        ║
 * ╚══════════════════════════════════════════════════════╝
 */

'use strict';

// ═══════════════════════════════════════════════════════════
//  TERMINAL CONTROL
// ═══════════════════════════════════════════════════════════
const ESC = '\x1b';
const out  = (s) => process.stdout.write(s);
const flush = () => {}; // Node.js stdout is synchronous

const at   = (col, row) => `${ESC}[${row};${col}H`;
const fg   = (r, g, b)  => `${ESC}[38;2;${r};${g};${b}m`;
const bg   = (r, g, b)  => `${ESC}[48;2;${r};${g};${b}m`;
const RST  = `${ESC}[0m`;
const BOLD = `${ESC}[1m`;
const HIDE = `${ESC}[?25l`;
const SHOW = `${ESC}[?25h`;
const ALT  = `${ESC}[?1049h`;
const NORM = `${ESC}[?1049l`;
const CLER = `${ESC}[2J${ESC}[H`;

// ═══════════════════════════════════════════════════════════
//  SCREEN
// ═══════════════════════════════════════════════════════════
const SW    = 80;
const SH    = 44;
const HUD_H = 2;
const T_ROWS = SH >> 1;  // 22

// ═══════════════════════════════════════════════════════════
//  COLOR PALETTE (NES-inspired 24-bit RGB tuples)
// ═══════════════════════════════════════════════════════════
const T   = null;

const SKY  = [92,  148, 252];
const SKY2 = [56,  100, 200];

const GT   = [0,   200, 0  ];
const GM   = [180, 120, 60 ];
const GD   = [120,  75, 25 ];

const BK   = [204, 102, 40 ];
const BD   = [150,  65, 15 ];
const BH   = [235, 148, 85 ];
const BM   = [ 85,  40,  8 ];

const QB   = [255, 192,  0 ];
const QD   = [196, 128,  0 ];
const QQ   = [255, 255, 120];
const QK   = [ 70,  35,  0 ];
const UB   = [155,  95, 45 ];

const PL   = [  0, 228,  0 ];
const PD   = [  0, 148,  0 ];
const PT_  = [  0, 255, 80 ];
const PDK  = [  0,  55,  0 ];

const CC   = [255, 204,  0 ];
const CD   = [196, 140,  0 ];

const MR   = [204,  48,  48];
const MS   = [255, 188, 102];
const MB   = [ 72,  72, 220];
const MW   = [108,  60,  18];
const ME   = [ 18,  18,  18];

const GB_  = [192, 118,  40];
const GK_  = [ 96,  56,   8];
const GF_  = [ 68,  32,   5];
const GW_  = [252, 252, 252];
const GP_  = [  8,   8,   8];

const CL   = [255, 255, 255];
const CLS  = [200, 218, 255];
const HG   = [  0, 168,  0 ];
const HGD  = [  0, 120,  0 ];

const FP   = [188, 188, 188];
const FC_  = [  0, 204,  0 ];
const FCD  = [  0, 140,  0 ];

const COI  = [255, 255,  80];
const SCP  = [255, 255, 255];

// ═══════════════════════════════════════════════════════════
//  SPRITES
// ═══════════════════════════════════════════════════════════
function flip(spr) {
  return spr.map(row => [...row].reverse());
}

// Mario 8×12
const MARIO_STAND_R = [
  [T,  T,  MR, MR, MR, MR, T,  T ],
  [T,  MR, MR, MR, MR, MR, MR, T ],
  [T,  MW, MS, MS, MW, MW, T,  T ],
  [T,  MS, ME, MS, MS, MS, MS, T ],
  [T,  MS, MS, MS, MS, MS, MS, T ],
  [T,  T,  MS, MS, MS, MS, T,  T ],
  [T,  MR, MB, MB, MB, MB, MR, T ],
  [T,  MB, MB, MB, MB, MB, MB, T ],
  [T,  MB, MB, T,  T,  MB, MB, T ],
  [T,  MW, MB, T,  T,  MB, MW, T ],
  [MW, MW, MW, T,  T,  MW, MW, MW],
  [T,  T,  T,  T,  T,  T,  T,  T ],
];
const MARIO_WALK1_R = [
  [T,  T,  MR, MR, MR, MR, T,  T ],
  [T,  MR, MR, MR, MR, MR, MR, T ],
  [T,  MW, MS, MS, MW, MW, T,  T ],
  [T,  MS, ME, MS, MS, MS, MS, T ],
  [T,  MS, MS, MS, MS, MS, MS, T ],
  [T,  MR, MB, MB, MB, MB, MR, T ],
  [T,  MB, MB, MB, MB, MB, MB, T ],
  [T,  T,  MB, T,  MB, MB, T,  T ],
  [T,  MB, MW, T,  T,  MW, T,  T ],
  [T,  MW, MW, T,  T,  T,  T,  T ],
  [T,  T,  T,  T,  T,  T,  T,  T ],
  [T,  T,  T,  T,  T,  T,  T,  T ],
];
const MARIO_WALK2_R = [
  [T,  T,  MR, MR, MR, MR, T,  T ],
  [T,  MR, MR, MR, MR, MR, MR, T ],
  [T,  MW, MS, MS, MW, MW, T,  T ],
  [T,  MS, MS, MS, MS, ME, MS, T ],
  [T,  MS, MS, MS, MS, MS, MS, T ],
  [T,  MR, MB, MB, MB, MB, MR, T ],
  [T,  MB, MB, MB, MB, MB, MB, T ],
  [T,  T,  MB, MB, T,  MB, T,  T ],
  [T,  T,  MW, T,  T,  MW, MB, T ],
  [T,  T,  T,  T,  T,  MW, MW, T ],
  [T,  T,  T,  T,  T,  T,  T,  T ],
  [T,  T,  T,  T,  T,  T,  T,  T ],
];
const MARIO_JUMP_R = [
  [T,  T,  MR, MR, MR, MR, T,  T ],
  [T,  MR, MR, MR, MR, MR, MR, T ],
  [T,  MW, MS, MS, MW, MW, T,  T ],
  [T,  MS, ME, MS, MS, MS, MS, T ],
  [MS, MS, MS, MS, MS, MS, MS, MS],
  [MB, MR, MB, MB, MB, MB, MR, MB],
  [MB, MB, MB, MB, MB, MB, MB, MB],
  [T,  MB, MB, T,  T,  MB, MB, T ],
  [MW, MW, T,  T,  T,  T,  MW, MW],
  [T,  T,  T,  T,  T,  T,  T,  T ],
  [T,  T,  T,  T,  T,  T,  T,  T ],
  [T,  T,  T,  T,  T,  T,  T,  T ],
];
const MARIO_DEAD_R = [
  [T,  T,  T,  MR, MR, T,  T,  T ],
  [T,  T,  MR, MR, MR, MR, T,  T ],
  [T,  MR, MR, MR, MR, MR, MR, T ],
  [T,  MW, MS, MS, MW, MW, T,  T ],
  [T,  MS, ME, MS, MS, MS, MS, T ],
  [MS, MS, MS, MS, MS, MS, MS, MS],
  [T,  MB, MB, MB, MB, MB, MB, T ],
  [T,  MR, MB, MB, MB, MB, MR, T ],
  [MW, MW, MW, T,  T,  MW, MW, MW],
  [T,  T,  T,  T,  T,  T,  T,  T ],
  [T,  T,  T,  T,  T,  T,  T,  T ],
  [T,  T,  T,  T,  T,  T,  T,  T ],
];
const MARIO_STAND_L = flip(MARIO_STAND_R);
const MARIO_WALK1_L = flip(MARIO_WALK1_R);
const MARIO_WALK2_L = flip(MARIO_WALK2_R);
const MARIO_JUMP_L  = flip(MARIO_JUMP_R);
const MARIO_DEAD_L  = flip(MARIO_DEAD_R);

// Goomba 8×10
const GOOMBA_W1 = [
  [T,   T,   GK_, GK_, GK_, GK_, T,   T  ],
  [T,   GK_, GB_, GB_, GB_, GB_, GK_, T  ],
  [GK_, GB_, GB_, GW_, GB_, GW_, GB_, GK_],
  [GK_, GB_, GK_, GP_, GB_, GP_, GB_, GK_],
  [GK_, GB_, GB_, GB_, GB_, GB_, GB_, GK_],
  [T,   GK_, GK_, GK_, GK_, GK_, GK_, T  ],
  [T,   GF_, GF_, GB_, GB_, GF_, GF_, T  ],
  [GF_, GF_, GB_, T,   T,   GB_, GF_, GF_],
  [GF_, GF_, GF_, T,   T,   GF_, GF_, GF_],
  [T,   T,   T,   T,   T,   T,   T,   T  ],
];
const GOOMBA_W2 = [
  [T,   T,   GK_, GK_, GK_, GK_, T,   T  ],
  [T,   GK_, GB_, GB_, GB_, GB_, GK_, T  ],
  [GK_, GB_, GB_, GW_, GB_, GW_, GB_, GK_],
  [GK_, GB_, GK_, GP_, GB_, GP_, GB_, GK_],
  [GK_, GB_, GB_, GB_, GB_, GB_, GB_, GK_],
  [T,   GK_, GK_, GK_, GK_, GK_, GK_, T  ],
  [T,   GB_, GF_, GB_, GB_, GF_, GB_, T  ],
  [T,   GF_, GF_, GB_, GB_, GF_, GF_, T  ],
  [T,   GF_, GF_, T,   T,   GF_, GF_, T  ],
  [T,   T,   T,   T,   T,   T,   T,   T  ],
];
const GOOMBA_DEAD = [
  [T,  T,  T,  T,  T,  T,  T,  T ],
  [T,  T,  T,  T,  T,  T,  T,  T ],
  [T,  T,  T,  T,  T,  T,  T,  T ],
  [T,  T,  T,  T,  T,  T,  T,  T ],
  [T,  T,  T,  T,  T,  T,  T,  T ],
  [T,  T,  T,  T,  T,  T,  T,  T ],
  [GK_,GK_,GK_,GK_,GK_,GK_,GK_,GK_],
  [GK_,GB_,GW_,GB_,GB_,GW_,GB_,GK_],
  [GK_,GF_,GF_,GF_,GF_,GF_,GF_,GK_],
  [T,  T,  T,  T,  T,  T,  T,  T ],
];

// ═══════════════════════════════════════════════════════════
//  TILES (8×8)
// ═══════════════════════════════════════════════════════════
const TILE_GROUND = [
  [GT, GT, GT, GT, GT, GT, GT, GT],
  [GT, GT, GT, GT, GT, GT, GT, GT],
  [GM, GM, GM, GM, GM, GM, GM, GM],
  [GM, GD, GM, GM, GD, GM, GM, GD],
  [GD, GD, GD, GD, GD, GD, GD, GD],
  [GD, GD, GD, GD, GD, GD, GD, GD],
  [GD, GD, GD, GD, GD, GD, GD, GD],
  [GD, GD, GD, GD, GD, GD, GD, GD],
];
const TILE_BRICK = [
  [BM, BM, BM, BM, BM, BM, BM, BM],
  [BM, BH, BK, BK, BM, BH, BK, BK],
  [BM, BK, BK, BK, BM, BK, BK, BK],
  [BM, BK, BD, BK, BM, BK, BD, BK],
  [BM, BM, BM, BM, BM, BM, BM, BM],
  [BM, BK, BM, BH, BK, BK, BM, BH],
  [BM, BK, BM, BK, BK, BK, BM, BK],
  [BM, BD, BM, BK, BD, BK, BM, BK],
];
const TILE_QBLOCK = [
  [QK, QB, QB, QB, QB, QB, QB, QK],
  [QB, QB, QD, QD, QB, QQ, QD, QB],
  [QB, QB, QQ, QD, QB, QQ, QD, QB],
  [QB, QB, QD, QQ, QB, QD, QD, QB],
  [QB, QB, QD, QQ, QB, QD, QD, QB],
  [QB, QB, QD, QD, QB, QQ, QD, QB],
  [QB, QD, QD, QD, QD, QD, QD, QB],
  [QK, QB, QB, QB, QB, QB, QB, QK],
];
const TILE_QUSED = [
  [QK, UB, UB, UB, UB, UB, UB, QK],
  [UB, UB, BD, BD, UB, UB, BD, UB],
  [UB, UB, BD, BD, UB, UB, BD, UB],
  [UB, UB, BD, BD, UB, UB, BD, UB],
  [UB, UB, BD, BD, UB, UB, BD, UB],
  [UB, UB, BD, BD, UB, UB, BD, UB],
  [UB, BD, BD, BD, BD, BD, BD, UB],
  [QK, UB, UB, UB, UB, UB, UB, QK],
];
const TILE_PIPE_TOP = [
  [T,  PL, PT_, PT_, PL, PD, PDK, T  ],
  [PL, PT_, PT_, PL, PL, PD, PDK, PD ],
  [PL, PT_, PL, PL, PD, PD, PDK, PD ],
  [PL, PL, PL, PD, PD, PDK, PDK, PD ],
  [PD, PD, PD, PD, PDK, PDK, PDK, PD ],
  [PDK,PDK,PDK,PDK,PDK,PDK,PDK,PDK],
  [T,  T,  T,  T,  T,  T,  T,  T  ],
  [T,  T,  T,  T,  T,  T,  T,  T  ],
];
const TILE_PIPE_BODY = [
  [T,  PL, PL, PL, PD, PDK, PDK, T  ],
  [T,  PL, PL, PL, PD, PDK, PDK, T  ],
  [T,  PL, PL, PL, PD, PDK, PDK, T  ],
  [T,  PL, PL, PL, PD, PDK, PDK, T  ],
  [T,  PL, PL, PL, PD, PDK, PDK, T  ],
  [T,  PL, PL, PL, PD, PDK, PDK, T  ],
  [T,  PL, PL, PL, PD, PDK, PDK, T  ],
  [T,  PL, PL, PL, PD, PDK, PDK, T  ],
];

// ═══════════════════════════════════════════════════════════
//  FRAMEBUFFER
// ═══════════════════════════════════════════════════════════
class FrameBuffer {
  constructor() {
    const blank = SKY;
    const dirty = [255, 0, 255];
    this.buf  = Array.from({ length: SH }, () => Array(SW).fill(blank));
    this.prev = Array.from({ length: SH }, () => Array(SW).fill(dirty));
  }

  put(x, y, color) {
    if (x >= 0 && x < SW && y >= 0 && y < SH && color !== null) {
      this.buf[y][x] = color;
    }
  }

  fill(x, y, w, h, color) {
    for (let dy = 0; dy < h; dy++)
      for (let dx = 0; dx < w; dx++)
        this.put(x + dx, y + dy, color);
  }

  sprite(spr, sx, sy) {
    for (let ry = 0; ry < spr.length; ry++) {
      const row = spr[ry];
      for (let rx = 0; rx < row.length; rx++) {
        const c = row[rx];
        if (c !== null) this.put(sx + rx, sy + ry, c);
      }
    }
  }

  clear(color = SKY) {
    for (let y = 0; y < SH; y++)
      for (let x = 0; x < SW; x++)
        this.buf[y][x] = color;
  }

  flush() {
    const parts = [RST];
    let lastFg = null, lastBg = null;
    for (let tr = 0; tr < T_ROWS; tr++) {
      const py0 = tr * 2;
      const py1 = py0 + 1;
      const row0 = this.buf[py0];
      const row1 = py1 < SH ? this.buf[py1] : row0;
      const pr0  = this.prev[py0];
      const pr1  = py1 < SH ? this.prev[py1] : pr0;
      let needMove = true;
      for (let tx = 0; tx < SW; tx++) {
        const ct = row0[tx];
        const cb = row1[tx];
        if (ct === pr0[tx] && cb === pr1[tx]) {
          needMove = true;
          continue;
        }
        if (needMove) {
          parts.push(at(tx + 1, tr + HUD_H + 1));
          needMove = false;
        }
        if (ct !== lastFg) {
          parts.push(fg(ct[0], ct[1], ct[2]));
          lastFg = ct;
        }
        if (cb !== lastBg) {
          parts.push(bg(cb[0], cb[1], cb[2]));
          lastBg = cb;
        }
        parts.push('\u2580'); // ▀
        pr0[tx] = ct;
        pr1[tx] = cb;
      }
    }
    parts.push(RST);
    out(parts.join(''));
  }

  forceRedraw() {
    const dirty = [255, 0, 255];
    for (let y = 0; y < SH; y++)
      for (let x = 0; x < SW; x++)
        this.prev[y][x] = dirty;
  }
}

// ═══════════════════════════════════════════════════════════
//  INPUT HANDLER
// ═══════════════════════════════════════════════════════════
class Input {
  constructor() {
    this._times = {};
    this._quit  = false;
    this._setup();
  }

  _setup() {
    if (process.stdin.isTTY) {
      process.stdin.setRawMode(true);
    }
    process.stdin.resume();
    process.stdin.setEncoding('utf8');
    process.stdin.on('data', (data) => this._onData(data));
  }

  _onData(data) {
    // Arrow keys come as multi-byte sequences
    if (data === '\x1b[A') { this._press('up');    return; }
    if (data === '\x1b[B') { this._press('down');  return; }
    if (data === '\x1b[C') { this._press('right'); return; }
    if (data === '\x1b[D') { this._press('left');  return; }
    if (data === '\x03' || data === '\x1b') { this._press('q'); return; }
    // Normal characters
    for (const ch of data) {
      this._press(ch.toLowerCase());
    }
  }

  _press(k) {
    this._times[k] = Date.now();
  }

  held(k, windowMs = 300) {
    return (Date.now() - (this._times[k] || 0)) < windowMs;
  }

  tapped(k, windowMs = 80) {
    return (Date.now() - (this._times[k] || 0)) < windowMs;
  }

  consume(k) {
    delete this._times[k];
  }

  restore() {
    this._quit = true;
    if (process.stdin.isTTY) {
      process.stdin.setRawMode(false);
    }
    process.stdin.pause();
  }
}

// ═══════════════════════════════════════════════════════════
//  LEVEL DATA
// ═══════════════════════════════════════════════════════════
const TILE_SIZE  = 8;
const GROUND_Y   = 36;
const MARIO_W    = 8;
const MARIO_H    = 12;
const LEVEL_PX_W = 1600;

// Platforms: [world_px_x, world_px_y, width_px, tile_type]
const PLATFORMS_DEF = [
  [0,    GROUND_Y, 352, 'ground'],
  [400,  GROUND_Y, 160, 'ground'],
  [600,  GROUND_Y, 240, 'ground'],
  [896,  GROUND_Y, 352, 'ground'],
  [1296, GROUND_Y, 304, 'ground'],

  [88,  16, 8, 'qblock'],
  [112, 16, 8, 'brick' ],
  [120, 16, 8, 'qblock'],
  [128, 16, 8, 'brick' ],
  [136, 16, 8, 'qblock'],

  [232, 16, 8, 'qblock'],
  [240, 16, 8, 'brick' ],
  [248, 16, 8, 'brick' ],
  [256, 16, 8, 'qblock'],
  [264, 16, 8, 'brick' ],

  [288,  8, 40, 'brick' ],
  [336,  8,  8, 'qblock'],
  [344,  8, 24, 'brick' ],

  [440, 16, 8, 'qblock'],
  [448, 16, 8, 'brick' ],

  [640, 16, 56, 'brick' ],
  [696, 16,  8, 'qblock'],
  [704, 16, 40, 'brick' ],

  [1232, GROUND_Y-8,  8, 'ground'],
  [1224, GROUND_Y-16, 8, 'ground'],
  [1216, GROUND_Y-24, 8, 'ground'],
  [1208, GROUND_Y-8,  8, 'ground'],
];

// Pipes: [world_px_x, height_in_tiles]
const PIPES = [
  [176, 2],
  [264, 3],
  [384, 2],
  [464, 3],
  [768, 4],
  [896, 2],
];

// Coins: [world_px_x, world_px_y]
const COINS_DEF = [
  [56, 22], [64, 22], [72, 22],
  [152, 22], [160, 22], [168, 22],
  [312, 22], [320, 22], [328, 22], [336, 22],
  [456, 22], [464, 22],
  [648, 22], [656, 22], [664, 22],
  [800, 14], [808, 14], [816, 14], [824, 14],
  [960, 22], [968, 22], [976, 22], [984, 22],
  [1072, 22], [1080, 22], [1088, 22],
  [1160, 22], [1168, 22],
];

// Enemies: [world_px_x, world_px_y, direction]
const ENEMIES_DEF = [
  [144, GROUND_Y-10, +1],
  [224, GROUND_Y-10, +1],
  [304, GROUND_Y-10, -1],
  [440, GROUND_Y-10, +1],
  [540, GROUND_Y-10, +1],
  [660, GROUND_Y-10, -1],
  [700, GROUND_Y-10, +1],
  [900, GROUND_Y-10, +1],
  [960, GROUND_Y-10, +1],
  [1020,GROUND_Y-10, -1],
  [1100,GROUND_Y-10, +1],
  [1180,GROUND_Y-10, -1],
];

// Clouds: [world_px_x, world_px_y]
const CLOUDS = [
  [40, 4], [120, 2], [200, 6], [300, 3], [400, 5],
  [520, 2], [620, 4], [720, 3], [820, 6], [920, 2],
  [1040, 4], [1140, 3], [1240, 5], [1380, 2],
];

const FLAGPOLE_X = 1520;

// ═══════════════════════════════════════════════════════════
//  GAME ENTITIES
// ═══════════════════════════════════════════════════════════
class MarioEntity {
  constructor(x, y) {
    this.x = x;
    this.y = y;
    this.vx = 0;
    this.vy = 0;
    this.onGround = false;
    this.facing = 1;
    this.animT = 0;
    this.dead = false;
    this.deadVy = 0;
    this.deadT = 0;
    this.win = false;
    this.walkFrame = 0;
  }

  static MAX_SPD = 3.5;
  static ACCEL   = 0.15;
  static DECEL   = 0.2;
  static JUMP_VY = -5.5;
  static GRAVITY = 0.45;
  static AIR_CTRL= 0.12;

  applyPhysics(left, right, onGround) {
    let targetVx = 0;
    if (left && !right)  targetVx = -MarioEntity.MAX_SPD;
    if (right && !left) targetVx = MarioEntity.MAX_SPD;

    const accel = onGround ? MarioEntity.ACCEL : MarioEntity.AIR_CTRL;
    if (targetVx === 0) {
      if (this.vx > 0) this.vx = Math.max(0, this.vx - MarioEntity.DECEL);
      else if (this.vx < 0) this.vx = Math.min(0, this.vx + MarioEntity.DECEL);
    } else {
      if (targetVx > this.vx) this.vx = Math.min(targetVx, this.vx + accel);
      else this.vx = Math.max(targetVx, this.vx - accel);
    }
    if (this.vx !== 0) this.facing = this.vx > 0 ? 1 : -1;
  }

  sprite() {
    if (this.dead)
      return this.facing > 0 ? MARIO_DEAD_R : MARIO_DEAD_L;
    if (!this.onGround)
      return this.facing > 0 ? MARIO_JUMP_R : MARIO_JUMP_L;
    if (Math.abs(this.vx) > 0.1) {
      const f = Math.floor(this.animT / 6) % 2;
      if (this.facing > 0) return f === 0 ? MARIO_WALK1_R : MARIO_WALK2_R;
      else                 return f === 0 ? MARIO_WALK1_L : MARIO_WALK2_L;
    }
    return this.facing > 0 ? MARIO_STAND_R : MARIO_STAND_L;
  }
}

class Goomba {
  static W   = 8;
  static H   = 10;
  static SPD = 0.8;

  constructor(x, y, direction) {
    this.x = x;
    this.y = y;
    this.vx = -direction * Goomba.SPD;
    this.vy = 0;
    this.onGround = false;
    this.dead = false;
    this.deadT = 0;
    this.animT = 0;
  }

  sprite() {
    if (this.dead) return GOOMBA_DEAD;
    return (Math.floor(this.animT / 8) % 2 === 0) ? GOOMBA_W1 : GOOMBA_W2;
  }
}

class Coin {
  constructor(x, y) {
    this.x = x;
    this.y = y;
    this.collected = false;
    this.animT = 0;
    this.popT  = 0;
    this.popY  = 0;
  }
}

class Particle {
  constructor(x, y, text, color = [255, 255, 80]) {
    this.x = x;
    this.y = y;
    this.text  = text;
    this.color = color;
    this.vy   = -1.2;
    this.life = 28;
    this.t    = 0;
  }
}

// ═══════════════════════════════════════════════════════════
//  GAME
// ═══════════════════════════════════════════════════════════
class Game {
  constructor() {
    this.fb        = new FrameBuffer();
    this.inp       = new Input();
    this.mario     = new MarioEntity(24, GROUND_Y - MARIO_H);
    this.goombas   = ENEMIES_DEF.map(([x, y, d]) => new Goomba(x, y, d));
    this.coins     = COINS_DEF.map(([x, y]) => new Coin(x, y));
    this.platforms = PLATFORMS_DEF.map(p => [...p]);
    this.usedQBlocks = new Set();
    this.particles = [];
    this.cameraX   = 0;
    this.score     = 0;
    this.coinCount = 0;
    this.lives     = 3;
    this.state     = 'play';
    this.frame     = 0;
    this.timeLeft  = 400;
    this.timeTick  = 0;
    this._timer    = null;
  }

  // ── Physics helpers ──────────────────────────────────────
  rectOverlap(ax, ay, aw, ah, bx, by, bw, bh) {
    return ax < bx+bw && ax+aw > bx && ay < by+bh && ay+ah > by;
  }

  getSolidRects() {
    const rects = [];
    for (const [px, py, pw] of this.platforms) {
      rects.push([px, py, pw, 8]);
    }
    for (const [px, th] of PIPES) {
      const pipePy = GROUND_Y - th * 8;
      rects.push([px, pipePy, 16, 8]);
      for (let ti = 1; ti < th; ti++) {
        rects.push([px + 2, pipePy + ti * 8, 12, 8]);
      }
    }
    return rects;
  }

  resolveEntity(ent, w, h, gravity = 0.45) {
    ent.vy += gravity;
    ent.y  += ent.vy;
    ent.x  += ent.vx;

    ent.x = Math.max(0, Math.min(ent.x, LEVEL_PX_W - w));

    let onGround   = false;
    let bumpBlock  = -1;

    const rects = this.getSolidRects();
    for (let i = 0; i < rects.length; i++) {
      const [px, py, pw, ph] = rects[i];
      if (!this.rectOverlap(ent.x, ent.y, w, h, px, py, pw, ph)) continue;

      const ol  = (ent.x + w) - px;
      const orr = (px + pw) - ent.x;
      const ot  = (ent.y + h) - py;
      const ob  = (py + ph) - ent.y;

      const minOv = Math.min(ol, orr, ot, ob);
      if (minOv === ot && ent.vy >= 0) {
        ent.y = py - h;
        ent.vy = 0;
        onGround = true;
      } else if (minOv === ob && ent.vy < 0) {
        ent.y = py + ph;
        ent.vy = 0;
        // Check if qblock
        for (let pi = 0; pi < this.platforms.length; pi++) {
          const [px2, py2, , ptype] = this.platforms[pi];
          if (px === px2 && py === py2 && ptype === 'qblock' && !this.usedQBlocks.has(pi)) {
            bumpBlock = pi;
          }
        }
      } else if (minOv === ol) {
        ent.x = px - w;
        ent.vx = 0;
      } else {
        ent.x = px + pw;
        ent.vx = 0;
      }
    }

    ent.onGround = onGround;
    return bumpBlock;
  }

  // ── Update ───────────────────────────────────────────────
  update() {
    if (this.state !== 'play') return;
    const inp   = this.inp;
    const mario = this.mario;

    if (mario.dead) {
      mario.deadVy += 0.5;
      mario.y += mario.deadVy;
      mario.deadT += 1;
      if (mario.deadT > 80) {
        this.lives--;
        if (this.lives <= 0) this.state = 'over';
        else this._respawn();
      }
      return;
    }

    if (mario.win) {
      mario.x += 1.0;
      return;
    }

    // Input
    const left  = inp.held('a') || inp.held('left');
    const right = inp.held('d') || inp.held('right');
    mario.applyPhysics(left, right, mario.onGround);

    const canJump = mario.onGround;
    if ((inp.tapped(' ') || inp.tapped('w') || inp.tapped('up')) && canJump) {
      mario.vy = MarioEntity.JUMP_VY;
      mario.onGround = false;
      inp.consume(' '); inp.consume('w'); inp.consume('up');
    }

    // Physics
    const bump = this.resolveEntity(mario, MARIO_W, MARIO_H, MarioEntity.GRAVITY);
    if (bump >= 0) {
      this.usedQBlocks.add(bump);
      this.platforms[bump][3] = 'qused';
      const [px, py] = this.platforms[bump];
      this.particles.push(new Particle(px + this.cameraX, py - 10, '+100', COI));
      this.score += 100;
      this.coinCount += 1;
    }

    if (mario.y > SH + 10) this._killMario();

    if (Math.abs(mario.vx) > 0.1 || !mario.onGround) mario.animT++;

    // Goombas
    for (const g of this.goombas) {
      if (g.dead) { g.deadT++; continue; }
      g.animT++;
      this.resolveEntity(g, Goomba.W, Goomba.H, 0.45);

      const standingX = Math.floor(g.x);
      const standingY = Math.floor(g.y) + Goomba.H + 1;
      let onEdge = true;
      for (const [px, py, pw, ph] of this.getSolidRects()) {
        if (px <= standingX && standingX < px + pw && py <= standingY && standingY < py + ph + 1) {
          onEdge = false; break;
        }
      }
      if (onEdge || g.x <= 0 || g.x >= LEVEL_PX_W - Goomba.W) g.vx = -g.vx;

      if (!mario.dead && this.rectOverlap(mario.x, mario.y, MARIO_W, MARIO_H, g.x, g.y, Goomba.W, Goomba.H)) {
        if (mario.vy > 0 && mario.y + MARIO_H < g.y + Goomba.H - 2) {
          g.dead = true;
          mario.vy = -3.5;
          this.score += 200;
          this.particles.push(new Particle(
            g.x - this.cameraX + Goomba.W / 2, g.y - 8, '+200', SCP));
        } else {
          this._killMario();
        }
      }
    }

    // Coins
    for (const c of this.coins) {
      if (c.collected) continue;
      c.animT++;
      if (this.rectOverlap(mario.x, mario.y, MARIO_W, MARIO_H, c.x, c.y, 4, 8)) {
        c.collected = true;
        this.coinCount++;
        this.score += 50;
        this.particles.push(new Particle(c.x - this.cameraX, c.y - 6, '+50', COI));
      }
    }

    // Flagpole
    if (!mario.win && mario.x + MARIO_W >= FLAGPOLE_X) {
      mario.win = true;
      this.state = 'win';
    }

    // Particles
    for (const p of this.particles) {
      p.y += p.vy;
      p.t++;
    }
    this.particles = this.particles.filter(p => p.t < p.life);

    // Camera
    const targetCam = mario.x - SW / 3;
    this.cameraX = Math.max(0, Math.min(targetCam, LEVEL_PX_W - SW));

    // Timer
    this.timeTick++;
    if (this.timeTick >= 60) {
      this.timeTick = 0;
      this.timeLeft = Math.max(0, this.timeLeft - 1);
      if (this.timeLeft === 0) this._killMario();
    }

    this.frame++;
  }

  _killMario() {
    const m = this.mario;
    if (!m.dead) {
      m.dead   = true;
      m.vy     = 0;
      m.deadVy = -6.5;
      m.vx     = 0;
    }
  }

  _respawn() {
    this.mario     = new MarioEntity(24, GROUND_Y - MARIO_H);
    this.cameraX   = 0;
    this.state     = 'play';
    this.timeLeft  = 400;
    this.timeTick  = 0;
    this.fb.forceRedraw();
  }

  // ── Draw helpers ─────────────────────────────────────────
  _wx(worldX) { return Math.floor(worldX - this.cameraX); }

  drawTile(tile, sx, sy) {
    for (let ry = 0; ry < tile.length; ry++) {
      const row = tile[ry];
      for (let rx = 0; rx < row.length; rx++) {
        const c = row[rx];
        if (c !== null) this.fb.put(sx + rx, sy + ry, c);
      }
    }
  }

  drawCloud(sx, sy) {
    const data = [[2,0,4],[1,1,6],[0,2,8],[0,3,8],[1,4,6]];
    for (const [ox, oy, w] of data) {
      for (let dx = 0; dx < w; dx++) {
        const c = ((dx + oy) % 3 === 0) ? CLS : CL;
        this.fb.put(sx + ox + dx, sy + oy, c);
      }
    }
  }

  drawHill(sx, sy, w, color, dark) {
    const half = Math.floor(w / 2);
    for (let row = 0; row < half; row++) {
      const x0 = half - row - 1;
      const x1 = half + row + 1;
      for (let dx = x0; dx < x1; dx++) {
        const c = (row === 0) ? dark : color;
        this.fb.put(sx + dx, sy + row, c);
      }
    }
  }

  drawFlagpole(sx) {
    if (sx < -2 || sx >= SW) return;
    for (let y = 4; y < GROUND_Y; y++) {
      const c = (y % 2 === 0) ? FP : [150, 150, 150];
      this.fb.put(sx, y, c);
    }
    for (let fy = 4; fy < 12; fy++) {
      for (let fx = 1; fx < 9; fx++) {
        const c = (fx < 6) ? FC_ : FCD;
        this.fb.put(sx + fx, fy, c);
      }
    }
  }

  // ── Render ───────────────────────────────────────────────
  render() {
    const fb  = this.fb;
    const cam = this.cameraX;

    // Sky
    for (let y = 0; y < SH; y++) {
      const c = (y < 6) ? SKY2 : SKY;
      for (let x = 0; x < SW; x++) fb.buf[y][x] = c;
    }

    // Clouds (parallax 0.5x)
    for (const [cx, cy] of CLOUDS) {
      const sx = Math.floor(cx - cam * 0.5);
      this.drawCloud(sx, cy);
    }

    // Hills
    for (const hx of [0, 96, 256, 480, 640, 800, 960, 1200]) {
      this.drawHill(Math.floor(hx - cam), GROUND_Y - 8, 20, HG, HGD);
    }

    // Ground / Platforms
    const tileMap = { ground: TILE_GROUND, brick: TILE_BRICK, qblock: TILE_QBLOCK, qused: TILE_QUSED };
    for (const [px, py, pw, ptype] of this.platforms) {
      const sx = this._wx(px);
      if (sx + pw < 0 || sx >= SW) continue;
      const tile = tileMap[ptype] || TILE_GROUND;
      for (let tx = 0; tx < Math.floor(pw / TILE_SIZE); tx++) {
        this.drawTile(tile, sx + tx * TILE_SIZE, py);
      }
    }

    // Pipes
    for (const [px, th] of PIPES) {
      const sx    = this._wx(px);
      if (sx + 16 < 0 || sx >= SW) continue;
      const pipePy = GROUND_Y - th * 8;
      this.drawTile(TILE_PIPE_TOP,  sx,   pipePy);
      this.drawTile(TILE_PIPE_TOP,  sx+8, pipePy);
      for (let ti = 1; ti < th; ti++) {
        this.drawTile(TILE_PIPE_BODY, sx,   pipePy + ti * 8);
        this.drawTile(TILE_PIPE_BODY, sx+8, pipePy + ti * 8);
      }
    }

    // Flagpole
    this.drawFlagpole(this._wx(FLAGPOLE_X));

    // Coins
    for (const c of this.coins) {
      if (c.collected) continue;
      const sx = this._wx(c.x);
      if (sx < -4 || sx >= SW) continue;
      const wobbleX = (Math.floor(c.animT / 8) % 4 < 2) ? 1 : 0;
      for (let ry = 0; ry < 8; ry++) {
        for (let rx = 0; rx < 4; rx++) {
          let ccColor = (rx < 2) ? CC : CD;
          if (ry === 0 || ry === 7) { if (rx === 0 || rx === 3) ccColor = null; }
          if (ccColor) fb.put(sx + rx + wobbleX, c.y + ry, ccColor);
        }
      }
    }

    // Goombas
    for (const g of this.goombas) {
      if (g.dead && g.deadT > 25) continue;
      const sx = this._wx(g.x);
      if (sx < -8 || sx >= SW) continue;
      fb.sprite(g.sprite(), sx, Math.floor(g.y));
    }

    // Mario
    const mario = this.mario;
    if (!(mario.dead && mario.y > SH + 8)) {
      const sx = this._wx(mario.x);
      fb.sprite(mario.sprite(), sx, Math.floor(mario.y));
    }

    // Particles
    for (const p of this.particles) {
      const fade = 1.0 - (p.t / p.life);
      if (fade < 0.3) continue;
      const r  = Math.floor(p.color[0] * fade);
      const g_ = Math.floor(p.color[1] * fade);
      const b_ = Math.floor(p.color[2] * fade);
      const col = [r, g_, b_];
      const px = Math.floor(p.x);
      const py = Math.floor(p.y);
      for (let i = 0; i < p.text.length; i++) {
        const cx = px + i;
        if (cx >= 0 && cx < SW && py >= 0 && py < SH) {
          fb.buf[py][cx] = col;
        }
      }
    }

    fb.flush();
  }

  drawHud() {
    const scoreStr = `SCORE:${String(this.score).padStart(6, '0')}`;
    const coinsStr = `\u2B50${String(this.coinCount).padStart(2, '0')}`;
    const worldStr = 'WORLD 1-1';
    const timeStr  = `TIME:${String(this.timeLeft).padStart(3, '0')}`;
    const livesStr = `\u2665x${this.lives}`;

    const line = `  ${BOLD}${fg(255,255,100)}${scoreStr}${RST}  ` +
                 `${fg(255,220,0)}${coinsStr}${RST}  ` +
                 `${BOLD}${fg(255,255,255)}${worldStr}${RST}  ` +
                 `${fg(255,100,100)}${timeStr}${RST}  ` +
                 `${fg(255,80,80)}${livesStr}${RST}`;

    out(at(1, 1));
    out(bg(SKY2[0],SKY2[1],SKY2[2]) + ' '.repeat(SW) + RST);
    out(at(1, 2));
    out(bg(SKY2[0],SKY2[1],SKY2[2]) + ' '.repeat(SW) + RST);
    out(at(1, 1));
    out(bg(SKY2[0],SKY2[1],SKY2[2]) + line + RST);
  }

  drawOverlay(title, subtitle, color = [255, 255, 100]) {
    const cx = Math.floor(SW / 2);
    const cy = Math.floor(T_ROWS / 2) + HUD_H;
    const w  = Math.max(title.length, subtitle.length) + 6;
    const boxX = cx - Math.floor(w / 2);
    const boxY = cy - 2;

    out(at(boxX, boxY));
    out(bg(20,20,80) + fg(100,100,255) + '\u2554' + '\u2550'.repeat(w) + '\u2557' + RST);
    out(at(boxX, boxY + 1));
    out(bg(20,20,80) + fg(100,100,255) + '\u2551' + RST);
    const pad1 = Math.floor((w - title.length) / 2);
    out(at(boxX + 1 + pad1, boxY + 1));
    out(bg(20,20,80) + BOLD + fg(color[0],color[1],color[2]) + title + RST);
    out(at(boxX + w, boxY + 1));
    out(bg(20,20,80) + fg(100,100,255) + '\u2551' + RST);
    out(at(boxX, boxY + 2));
    const pad2 = Math.floor((w - subtitle.length) / 2);
    out(bg(20,20,80) + fg(100,100,255) + '\u2551' + RST);
    out(at(boxX + 1 + pad2, boxY + 2));
    out(bg(20,20,80) + fg(200,200,200) + subtitle + RST);
    out(at(boxX + w, boxY + 2));
    out(bg(20,20,80) + fg(100,100,255) + '\u2551' + RST);
    out(at(boxX, boxY + 3));
    out(bg(20,20,80) + fg(100,100,255) + '\u255a' + '\u2550'.repeat(w) + '\u255d' + RST);
  }

  // ── Main Loop ────────────────────────────────────────────
  run() {
    out(ALT + CLER + HIDE);
    this.fb.forceRedraw();

    const FPS       = 30;
    const frameTime = Math.floor(1000 / FPS);

    // Quit handler
    const onQuit = () => {
      if (this.inp.held('q') || this.inp.held('esc') || this.inp.held('\x03')) {
        this._cleanup();
      }
    };

    this._timer = setInterval(() => {
      // Quit check
      if (this.inp.held('q') || this.inp.tapped('\x03')) {
        this._cleanup();
        return;
      }

      this.update();
      this.render();
      this.drawHud();

      if (this.state === 'over') {
        this.drawOverlay('GAME OVER', 'Press Q to quit', [255, 60, 60]);
      } else if (this.state === 'win') {
        this.drawOverlay(' STAGE CLEAR! ', `Score: ${this.score}  Coins: ${this.coinCount}`, [255, 255, 80]);
      }
    }, frameTime);
  }

  _cleanup() {
    if (this._timer) {
      clearInterval(this._timer);
      this._timer = null;
    }
    this.inp.restore();
    out(NORM + SHOW + RST + CLER);
    process.stdout.write(`\n\u{1F47E} Thanks for playing FC Terminal Mario!\n`);
    process.stdout.write(`   Final Score: ${this.score}  Coins: ${this.coinCount}  Lives left: ${this.lives}\n`);
    process.exit(0);
  }
}

// ═══════════════════════════════════════════════════════════
//  ENTRY POINT
// ═══════════════════════════════════════════════════════════
function main() {
  if (!process.stdin.isTTY) {
    console.error('Error: must run in an interactive terminal');
    process.exit(1);
  }

  let cols = 80, rows = 26;
  try {
    cols = process.stdout.columns || 80;
    rows = process.stdout.rows    || 26;
  } catch (e) {}

  const needCols = SW;
  const needRows = T_ROWS + HUD_H + 1;
  if (cols < needCols || rows < needRows) {
    console.error(`\u26A0\uFE0F  Terminal too small! Need at least ${needCols}\u00D7${needRows}`);
    console.error(`   Current: ${cols}\u00D7${rows}`);
    console.error(`   Please resize your terminal and try again.`);
    process.exit(1);
  }

  // Handle Ctrl+C gracefully
  process.on('SIGINT', () => {
    out(NORM + SHOW + RST + CLER);
    process.exit(0);
  });

  const game = new Game();
  game.run();
}

main();
