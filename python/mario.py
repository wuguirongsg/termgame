#!/usr/bin/env python3
"""
╔══════════════════════════════════════════════════════╗
║        FC TERMINAL MARIO  -  v1.0                   ║
║   NES风格超级玛丽 · 完整终端游戏                      ║
║   Controls: A/← D/→ Move  SPACE Jump  Q Quit        ║
╚══════════════════════════════════════════════════════╝
"""
import sys, time, threading, tty, termios, select, signal, os

# ═══════════════════════════════════════════════════════════
#  TERMINAL CONTROL
# ═══════════════════════════════════════════════════════════
ESC = '\033'
def _out(s): sys.stdout.write(s)
def _flush(): sys.stdout.flush()
def at(col, row): return f'{ESC}[{row};{col}H'
def fg(r,g,b):    return f'{ESC}[38;2;{r};{g};{b}m'
def bg(r,g,b):    return f'{ESC}[48;2;{r};{g};{b}m'
RST  = f'{ESC}[0m'
BOLD = f'{ESC}[1m'
HIDE = f'{ESC}[?25l'
SHOW = f'{ESC}[?25h'
ALT  = f'{ESC}[?1049h'   # alternate screen
NORM = f'{ESC}[?1049l'
CLER = f'{ESC}[2J{ESC}[H'

# ═══════════════════════════════════════════════════════════
#  SCREEN  (half-block rendering: 1 terminal row = 2 px rows)
# ═══════════════════════════════════════════════════════════
SW, SH = 80, 44          # pixel dimensions
HUD_H  = 2               # HUD terminal rows at top
T_ROWS = SH // 2         # game terminal rows = 22

# ═══════════════════════════════════════════════════════════
#  COLOR PALETTE  (NES-inspired 24-bit RGB)
# ═══════════════════════════════════════════════════════════
T    = None               # transparent

SKY  = (92,  148, 252)   # classic NES sky
SKY2 = (56,  100, 200)   # horizon sky

GT   = (0,   200, 0  )   # grass top
GM   = (180, 120, 60 )   # ground mid
GD   = (120,  75, 25 )   # ground dark

BK   = (204, 102, 40 )   # brick
BD   = (150,  65, 15 )   # brick dark
BH   = (235, 148, 85 )   # brick highlight
BM   = ( 85,  40,  8 )   # mortar

QB   = (255, 192,  0 )   # question gold
QD   = (196, 128,  0 )   # question dark
QQ   = (255, 255, 120)   # ? symbol
QK   = ( 70,  35,  0 )   # question border
UB   = (155,  95, 45 )   # used block

PL   = (  0, 228,  0 )   # pipe light
PD   = (  0, 148,  0 )   # pipe dark
PT_  = (  0, 255, 80 )   # pipe top
PDK  = (  0,  55,  0 )   # pipe deep dark

CC   = (255, 204,  0 )   # coin
CD   = (196, 140,  0 )   # coin dark

MR   = (204,  48,  48)   # mario red
MS   = (255, 188, 102)   # mario skin
MB   = ( 72,  72, 220)   # mario blue
MW   = (108,  60,  18)   # mario brown
ME   = ( 18,  18,  18)   # mario eye

GB_  = (192, 118,  40)   # goomba body
GK_  = ( 96,  56,   8)   # goomba dark
GF_  = ( 68,  32,   5)   # goomba foot
GW_  = (252, 252, 252)   # goomba eye white
GP_  = (  8,   8,   8)   # goomba pupil

CL   = (255, 255, 255)   # cloud white
CLS  = (200, 218, 255)   # cloud shadow
HG   = (  0, 168,  0 )   # hill green
HGD  = (  0, 120,  0 )   # hill dark

FP   = (188, 188, 188)   # flagpole
FC_  = (  0, 204,  0 )   # flag green
FCD  = (  0, 140,  0 )   # flag dark

COI  = (255, 255,  80)   # coin popup
SCP  = (255, 255, 255)   # score popup

# ═══════════════════════════════════════════════════════════
#  SPRITES  (pixel rows, each row = list of colors | T)
# ═══════════════════════════════════════════════════════════
def flip(spr): return [list(reversed(row)) for row in spr]

# ── Mario 8×12 ──────────────────────────────────────────────
MARIO_STAND_R = [
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
]
MARIO_WALK1_R = [
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
]
MARIO_WALK2_R = [
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
]
MARIO_JUMP_R = [
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
]
MARIO_DEAD_R = [
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
]
MARIO_STAND_L = flip(MARIO_STAND_R)
MARIO_WALK1_L = flip(MARIO_WALK1_R)
MARIO_WALK2_L = flip(MARIO_WALK2_R)
MARIO_JUMP_L  = flip(MARIO_JUMP_R)
MARIO_DEAD_L  = flip(MARIO_DEAD_R)

# ── Goomba 8×10 ─────────────────────────────────────────────
GOOMBA_W1 = [
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
]
GOOMBA_W2 = [
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
]
GOOMBA_DEAD = [
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
]

# ═══════════════════════════════════════════════════════════
#  TILES  (8×8 pixel arrays)
# ═══════════════════════════════════════════════════════════
TILE_GROUND = [
  [GT, GT, GT, GT, GT, GT, GT, GT],
  [GT, GT, GT, GT, GT, GT, GT, GT],
  [GM, GM, GM, GM, GM, GM, GM, GM],
  [GM, GD, GM, GM, GD, GM, GM, GD],
  [GD, GD, GD, GD, GD, GD, GD, GD],
  [GD, GD, GD, GD, GD, GD, GD, GD],
  [GD, GD, GD, GD, GD, GD, GD, GD],
  [GD, GD, GD, GD, GD, GD, GD, GD],
]
TILE_BRICK = [
  [BM, BM, BM, BM, BM, BM, BM, BM],
  [BM, BH, BK, BK, BM, BH, BK, BK],
  [BM, BK, BK, BK, BM, BK, BK, BK],
  [BM, BK, BD, BK, BM, BK, BD, BK],
  [BM, BM, BM, BM, BM, BM, BM, BM],
  [BM, BK, BM, BH, BK, BK, BM, BH],
  [BM, BK, BM, BK, BK, BK, BM, BK],
  [BM, BD, BM, BK, BD, BK, BM, BK],
]
TILE_QBLOCK = [
  [QK, QB, QB, QB, QB, QB, QB, QK],
  [QB, QB, QD, QD, QB, QQ, QD, QB],
  [QB, QB, QQ, QD, QB, QQ, QD, QB],
  [QB, QB, QD, QQ, QB, QD, QD, QB],
  [QB, QB, QD, QQ, QB, QD, QD, QB],
  [QB, QB, QD, QD, QB, QQ, QD, QB],
  [QB, QD, QD, QD, QD, QD, QD, QB],
  [QK, QB, QB, QB, QB, QB, QB, QK],
]
TILE_QUSED = [
  [QK, UB, UB, UB, UB, UB, UB, QK],
  [UB, UB, BD, BD, UB, UB, BD, UB],
  [UB, UB, BD, BD, UB, UB, BD, UB],
  [UB, UB, BD, BD, UB, UB, BD, UB],
  [UB, UB, BD, BD, UB, UB, BD, UB],
  [UB, UB, BD, BD, UB, UB, BD, UB],
  [UB, BD, BD, BD, BD, BD, BD, UB],
  [QK, UB, UB, UB, UB, UB, UB, QK],
]
TILE_PIPE_TOP = [
  [T,  PL, PT_,PT_,PL, PD, PDK,T  ],
  [PL, PT_,PT_,PL, PL, PD, PDK,PD ],
  [PL, PT_,PL, PL, PD, PD, PDK,PD ],
  [PL, PL, PL, PD, PD, PDK,PDK,PD ],
  [PD, PD, PD, PD, PDK,PDK,PDK,PD ],
  [PDK,PDK,PDK,PDK,PDK,PDK,PDK,PDK],
  [T,  T,  T,  T,  T,  T,  T,  T  ],
  [T,  T,  T,  T,  T,  T,  T,  T  ],
]
TILE_PIPE_BODY = [
  [T,  PL, PL, PL, PD, PDK,PDK,T  ],
  [T,  PL, PL, PL, PD, PDK,PDK,T  ],
  [T,  PL, PL, PL, PD, PDK,PDK,T  ],
  [T,  PL, PL, PL, PD, PDK,PDK,T  ],
  [T,  PL, PL, PL, PD, PDK,PDK,T  ],
  [T,  PL, PL, PL, PD, PDK,PDK,T  ],
  [T,  PL, PL, PL, PD, PDK,PDK,T  ],
  [T,  PL, PL, PL, PD, PDK,PDK,T  ],
]

# ═══════════════════════════════════════════════════════════
#  FRAMEBUFFER  (half-block pixel renderer)
#  ▀ = upper half block: fg=top px, bg=bottom px
# ═══════════════════════════════════════════════════════════
class FrameBuffer:
    def __init__(self):
        blank = SKY
        dirty = (255, 0, 255)  # force initial redraw
        self.buf  = [[blank] * SW for _ in range(SH)]
        self.prev = [[dirty] * SW for _ in range(SH)]

    def put(self, x, y, color):
        if 0 <= x < SW and 0 <= y < SH and color is not None:
            self.buf[y][x] = color

    def fill(self, x, y, w, h, color):
        for dy in range(h):
            for dx in range(w):
                self.put(x+dx, y+dy, color)

    def sprite(self, spr, sx, sy):
        for ry, row in enumerate(spr):
            for rx, c in enumerate(row):
                if c is not None:
                    self.put(sx+rx, sy+ry, c)

    def clear(self, color=SKY):
        for row in self.buf:
            for i in range(SW):
                row[i] = color

    def flush(self):
        out = [RST]
        last_fg = last_bg = None
        for tr in range(T_ROWS):
            py0 = tr * 2
            py1 = py0 + 1
            row0 = self.buf[py0]
            row1 = self.buf[py1] if py1 < SH else row0
            pr0  = self.prev[py0]
            pr1  = self.prev[py1] if py1 < SH else pr0
            need_move = True
            for tx in range(SW):
                ct, cb = row0[tx], row1[tx]
                if ct == pr0[tx] and cb == pr1[tx]:
                    need_move = True
                    continue
                if need_move:
                    out.append(at(tx+1, tr + HUD_H + 1))
                    need_move = False
                if ct != last_fg:
                    out.append(fg(*ct)); last_fg = ct
                if cb != last_bg:
                    out.append(bg(*cb)); last_bg = cb
                out.append('▀')
                pr0[tx] = ct
                pr1[tx] = cb
        out.append(RST)
        _out(''.join(out)); _flush()

    def force_redraw(self):
        for row in self.prev:
            for i in range(SW):
                row[i] = (255, 0, 255)

# ═══════════════════════════════════════════════════════════
#  INPUT HANDLER  (background thread, raw terminal)
# ═══════════════════════════════════════════════════════════
class Input:
    def __init__(self):
        self._times = {}
        self._quit  = False
        self.fd     = sys.stdin.fileno()
        self._old   = termios.tcgetattr(self.fd)
        tty.setraw(self.fd)
        t = threading.Thread(target=self._loop, daemon=True)
        t.start()

    def _loop(self):
        while not self._quit:
            r, _, _ = select.select([sys.stdin], [], [], 0.04)
            if not r: continue
            ch = sys.stdin.read(1)
            if ch == '\x1b':
                try:
                    ch2 = sys.stdin.read(1)
                    if ch2 == '[':
                        ch3 = sys.stdin.read(1)
                        if   ch3 == 'A': self._press('up')
                        elif ch3 == 'B': self._press('down')
                        elif ch3 == 'C': self._press('right')
                        elif ch3 == 'D': self._press('left')
                except:
                    self._press('esc')
            else:
                self._press(ch.lower())

    def _press(self, k): self._times[k] = time.time()

    def held(self, k, window=0.3):
        return (time.time() - self._times.get(k, 0)) < window

    def tapped(self, k, window=0.08):
        return (time.time() - self._times.get(k, 0)) < window

    def consume(self, k): self._times.pop(k, None)

    def restore(self):
        self._quit = True
        termios.tcsetattr(self.fd, termios.TCSADRAIN, self._old)

# ═══════════════════════════════════════════════════════════
#  LEVEL DATA
# ═══════════════════════════════════════════════════════════
TILE_SIZE   = 8
GROUND_Y    = 36   # pixel y of ground top
MARIO_W, MARIO_H = 8, 12
LEVEL_PX_W  = 1600  # total level width in pixels

# Platforms: (world_px_x, world_px_y, width_px, tile_type)
# tile_type: 'ground'|'brick'|'qblock'|'qused'
PLATFORMS = [
    # ── Ground sections (with a gap/pit in between) ──────────
    (0,    GROUND_Y, 352, 'ground'),
    (400,  GROUND_Y, 160, 'ground'),
    (600,  GROUND_Y, 240, 'ground'),
    (896,  GROUND_Y, 352, 'ground'),
    (1296, GROUND_Y, 304, 'ground'),

    # ── Brick & Question block rows ──────────────────────────
    # y = 16 → about 3 tiles above ground  (game-feel spot)
    (88,  16, 8, 'qblock'),   # 1st ? block
    (112, 16, 8, 'brick'),
    (120, 16, 8, 'qblock'),   # 2nd ? block
    (128, 16, 8, 'brick'),
    (136, 16, 8, 'qblock'),   # 3rd ? block (triple)

    (232, 16, 8, 'qblock'),
    (240, 16, 8, 'brick'),
    (248, 16, 8, 'brick'),
    (256, 16, 8, 'qblock'),
    (264, 16, 8, 'brick'),

    # Upper row of bricks
    (288,  8, 40, 'brick'),
    (336,  8,  8, 'qblock'),
    (344,  8, 24, 'brick'),

    (440, 16, 8, 'qblock'),
    (448, 16, 8, 'brick'),

    (640, 16, 56, 'brick'),
    (696, 16,  8, 'qblock'),
    (704, 16, 40, 'brick'),

    # Staircase approach to flag
    (1232, GROUND_Y-8,  8, 'ground'),
    (1224, GROUND_Y-16, 8, 'ground'),
    (1216, GROUND_Y-24, 8, 'ground'),
    (1208, GROUND_Y-8,  8, 'ground'),
]

# Pipes: (world_px_x, height_in_tiles)  → 2 tiles wide (16px)
PIPES = [
    (176, 2),
    (264, 3),
    (384, 2),
    (464, 3),
    (768, 4),
    (896, 2),
]

# Coins: (world_px_x, world_px_y)
COINS = [
    (56, 22), (64, 22), (72, 22),
    (152, 22), (160, 22), (168, 22),
    (312, 22), (320, 22), (328, 22), (336, 22),
    (456, 22), (464, 22),
    (648, 22), (656, 22), (664, 22),
    (800, 14), (808, 14), (816, 14), (824, 14),
    (960, 22), (968, 22), (976, 22), (984, 22),
    (1072, 22), (1080, 22), (1088, 22),
    (1160, 22), (1168, 22),
]

# Goombas: (world_px_x, world_px_y, direction)  dir: +1=left, -1=right
ENEMIES = [
    (144, GROUND_Y-10, +1),
    (224, GROUND_Y-10, +1),
    (304, GROUND_Y-10, -1),
    (440, GROUND_Y-10, +1),
    (540, GROUND_Y-10, +1),
    (660, GROUND_Y-10, -1),
    (700, GROUND_Y-10, +1),
    (900, GROUND_Y-10, +1),
    (960, GROUND_Y-10, +1),
    (1020,GROUND_Y-10, -1),
    (1100,GROUND_Y-10, +1),
    (1180,GROUND_Y-10, -1),
]

# Cloud decorations: (world_px_x, world_px_y)
CLOUDS = [
    (40, 4), (120, 2), (200, 6), (300, 3), (400, 5),
    (520, 2), (620, 4), (720, 3), (820, 6), (920, 2),
    (1040,4), (1140,3), (1240,5), (1380,2),
]

# Flagpole x position
FLAGPOLE_X = 1520

# ═══════════════════════════════════════════════════════════
#  GAME ENTITIES
# ═══════════════════════════════════════════════════════════
class Mario:
    W, H = 8, 12
    MAX_SPD = 3.5
    ACCEL = 0.15
    DECEL = 0.2
    JUMP_VY  = -5.5
    GRAVITY  = 0.45
    AIR_CTRL = 0.12

    def __init__(self, x, y):
        self.x = float(x)
        self.y = float(y)
        self.vx = 0.0
        self.vy = 0.0
        self.on_ground = False
        self.facing = 1
        self.anim_t  = 0
        self.dead    = False
        self.dead_vy = 0.0
        self.dead_t  = 0
        self.win     = False
        self.walk_frame = 0

    def apply_physics(self, left, right, on_ground):
        target_vx = 0
        if left and not right:
            target_vx = -self.MAX_SPD
        elif right and not left:
            target_vx = self.MAX_SPD

        accel = self.ACCEL if on_ground else self.AIR_CTRL
        if target_vx == 0:
            if self.vx > 0:
                self.vx = max(0, self.vx - self.DECEL)
            elif self.vx < 0:
                self.vx = min(0, self.vx + self.DECEL)
        else:
            if target_vx > self.vx:
                self.vx = min(target_vx, self.vx + accel)
            else:
                self.vx = max(target_vx, self.vx - accel)

        if self.vx != 0:
            self.facing = 1 if self.vx > 0 else -1

    def sprite(self):
        if self.dead:
            return MARIO_DEAD_R if self.facing > 0 else MARIO_DEAD_L
        if not self.on_ground:
            return MARIO_JUMP_R if self.facing > 0 else MARIO_JUMP_L
        if abs(self.vx) > 0.1:
            f = (self.anim_t // 6) % 2
            if self.facing > 0:
                return MARIO_WALK1_R if f == 0 else MARIO_WALK2_R
            else:
                return MARIO_WALK1_L if f == 0 else MARIO_WALK2_L
        return MARIO_STAND_R if self.facing > 0 else MARIO_STAND_L

class Goomba:
    W, H = 8, 10
    SPD = 0.8

    def __init__(self, x, y, direction):
        self.x = float(x)
        self.y = float(y)
        self.vx = -direction * self.SPD
        self.vy = 0.0
        self.on_ground = False
        self.dead = False
        self.dead_t = 0
        self.anim_t = 0

    def sprite(self):
        if self.dead: return GOOMBA_DEAD
        f = (self.anim_t // 8) % 2
        return GOOMBA_W1 if f == 0 else GOOMBA_W2

class Coin:
    def __init__(self, x, y):
        self.x = x
        self.y = y
        self.collected = False
        self.anim_t = 0
        self.pop_t = 0    # popup animation timer
        self.pop_y = 0.0

class Particle:
    def __init__(self, x, y, text, color=(255,255,80)):
        self.x = float(x)
        self.y = float(y)
        self.text = text
        self.color = color
        self.vy = -1.2
        self.life = 28
        self.t = 0

# ═══════════════════════════════════════════════════════════
#  GAME
# ═══════════════════════════════════════════════════════════
class Game:
    def __init__(self):
        self.fb    = FrameBuffer()
        self.inp   = Input()
        self.mario = Mario(24, GROUND_Y - Mario.H)
        self.goombas  = [Goomba(x, y, d) for x,y,d in ENEMIES]
        self.coins    = [Coin(x, y) for x,y in COINS]
        self.platforms = list(PLATFORMS)
        self.used_qblocks = set()  # indices of hit question blocks
        self.particles = []
        self.camera_x  = 0
        self.score     = 0
        self.coin_count= 0
        self.lives     = 3
        self.state     = 'play'   # 'play'|'dead'|'win'|'over'
        self.frame     = 0
        self.time_left = 400
        self.time_tick = 0

    # ── Physics helpers ──────────────────────────────────────
    def rect_overlap(self, ax,ay,aw,ah, bx,by,bw,bh):
        return ax < bx+bw and ax+aw > bx and ay < by+bh and ay+ah > by

    def get_solid_rects(self):
        """Returns list of (x,y,w,h) for all solid blocks in world space."""
        rects = []
        for px, py, pw, ptype in self.platforms:
            rects.append((px, py, pw, 8))
        for px, th in PIPES:
            pipe_py = GROUND_Y - th * 8
            # top tile (16px wide)
            rects.append((px, pipe_py, 16, 8))
            # body tiles
            for ti in range(1, th):
                rects.append((px+2, pipe_py + ti*8, 12, 8))
        return rects

    def resolve_entity(self, ent, w, h, gravity=0.45):
        """Apply gravity + collide entity with solid blocks."""
        ent.vy += gravity
        ent.y  += ent.vy
        ent.x  += ent.vx

        # Clamp to level
        ent.x = max(0, min(ent.x, LEVEL_PX_W - w))

        on_ground = False
        bump_block = -1   # index of qblock hit from below

        for i, (px, py, pw, ph) in enumerate(self.get_solid_rects()):
            if not self.rect_overlap(ent.x, ent.y, w, h, px, py, pw, ph):
                continue
            # Compute overlap amounts
            ol = (ent.x + w) - px     # overlap from left
            or_ = (px + pw) - ent.x   # overlap from right
            ot = (ent.y + h) - py     # overlap from top (entity falling into)
            ob = (py + ph) - ent.y    # overlap from bottom (entity going up)

            min_ov = min(ol, or_, ot, ob)
            if min_ov == ot and ent.vy >= 0:
                ent.y = py - h
                ent.vy = 0
                on_ground = True
            elif min_ov == ob and ent.vy < 0:
                ent.y = py + ph
                ent.vy = 0
                # Check if this is a qblock
                for pi, (px2,py2,pw2,ptype) in enumerate(self.platforms):
                    if px==px2 and py==py2 and ptype=='qblock' and pi not in self.used_qblocks:
                        bump_block = pi
            elif min_ov == ol:
                ent.x = px - w
                ent.vx = 0
            elif min_ov == or_:
                ent.x = px + pw
                ent.vx = 0

        ent.on_ground = on_ground
        return bump_block

    # ── Update ───────────────────────────────────────────────
    def update(self):
        if self.state not in ('play',): return
        inp = self.inp
        mario = self.mario

        if mario.dead:
            mario.dead_vy += 0.5
            mario.y += mario.dead_vy
            mario.dead_t += 1
            if mario.dead_t > 80:
                self.lives -= 1
                if self.lives <= 0:
                    self.state = 'over'
                else:
                    self._respawn()
            return

        if mario.win:
            mario.x += 1.0
            return

        # ── Input ──────────────────────────────────────────
        left = inp.held('a') or inp.held('left')
        right = inp.held('d') or inp.held('right')
        mario.apply_physics(left, right, mario.on_ground)

        can_jump = mario.on_ground
        if (inp.tapped(' ') or inp.tapped('w') or inp.tapped('up')) and can_jump:
            mario.vy = Mario.JUMP_VY
            mario.on_ground = False
            inp.consume(' '); inp.consume('w'); inp.consume('up')

        # ── Physics ────────────────────────────────────────
        bump = self.resolve_entity(mario, Mario.W, Mario.H, Mario.GRAVITY)
        if bump >= 0:
            self.used_qblocks.add(bump)
            self.platforms[bump] = (self.platforms[bump][0], self.platforms[bump][1], self.platforms[bump][2], 'qused')
            px, py, *_ = self.platforms[bump]
            self.particles.append(Particle(px + self.camera_x, py - 10, '+100', COI))
            self.score += 100
            self.coin_count += 1

        if mario.y > SH + 10:
            self._kill_mario()

        if abs(mario.vx) > 0.1 or not mario.on_ground:
            mario.anim_t += 1

        # ── Goombas ────────────────────────────────────────
        for g in self.goombas:
            if g.dead:
                g.dead_t += 1
                continue
            g.anim_t += 1
            self.resolve_entity(g, Goomba.W, Goomba.H, 0.45)

            # Reverse at platform edges
            standing_x = int(g.x)
            standing_y = int(g.y) + Goomba.H + 1
            on_edge = True
            for px,py,pw,ph in self.get_solid_rects():
                if px <= standing_x < px+pw and py <= standing_y < py+ph+1:
                    on_edge = False; break
            if on_edge or g.x <= 0 or g.x >= LEVEL_PX_W - Goomba.W:
                g.vx = -g.vx

            # Goomba vs Mario
            if not mario.dead and self.rect_overlap(mario.x,mario.y,Mario.W,Mario.H,
                                                     g.x, g.y, Goomba.W, Goomba.H):
                if mario.vy > 0 and mario.y + Mario.H < g.y + Goomba.H - 2:
                    # Jump on goomba
                    g.dead = True
                    mario.vy = -3.5
                    self.score += 200
                    self.particles.append(Particle(
                        g.x - self.camera_x + Goomba.W//2, g.y - 8, '+200', SCP))
                else:
                    self._kill_mario()

        # ── Coins ──────────────────────────────────────────
        for c in self.coins:
            if c.collected: continue
            c.anim_t += 1
            if self.rect_overlap(mario.x,mario.y,Mario.W,Mario.H, c.x,c.y,4,8):
                c.collected = True
                self.coin_count += 1
                self.score += 50
                self.particles.append(Particle(c.x - self.camera_x, c.y - 6, '+50', COI))

        # ── Flagpole ───────────────────────────────────────
        if not mario.win and mario.x + Mario.W >= FLAGPOLE_X:
            mario.win = True
            self.state = 'win'

        # ── Particles ─────────────────────────────────────
        for p in self.particles:
            p.y += p.vy
            p.t += 1
        self.particles = [p for p in self.particles if p.t < p.life]

        # ── Camera ────────────────────────────────────────
        target_cam = mario.x - SW // 3
        self.camera_x = max(0, min(target_cam, LEVEL_PX_W - SW))

        # ── Timer ─────────────────────────────────────────
        self.time_tick += 1
        if self.time_tick >= 60:
            self.time_tick = 0
            self.time_left = max(0, self.time_left - 1)
            if self.time_left == 0:
                self._kill_mario()

        self.frame += 1

    def _kill_mario(self):
        m = self.mario
        if not m.dead:
            m.dead = True
            m.vy = 0
            m.dead_vy = -6.5
            m.vx = 0

    def _respawn(self):
        self.mario = Mario(24, GROUND_Y - Mario.H)
        self.camera_x = 0
        self.state = 'play'
        self.time_left = 400
        self.time_tick = 0
        self.fb.force_redraw()

    # ── Draw helpers ─────────────────────────────────────────
    def _wx(self, world_x): return int(world_x - self.camera_x)

    def draw_tile(self, tile, sx, sy):
        for ry, row in enumerate(tile):
            for rx, c in enumerate(row):
                if c is not None:
                    self.fb.put(sx+rx, sy+ry, c)

    def draw_cloud(self, sx, sy):
        # Simple blocky cloud
        data = [
            (2,0,4), (1,1,6), (0,2,8), (0,3,8), (1,4,6)
        ]
        for ox, oy, w in data:
            for dx in range(w):
                c = CLS if (dx + oy) % 3 == 0 else CL
                self.fb.put(sx+ox+dx, sy+oy, c)

    def draw_hill(self, sx, sy, w, color, dark):
        # Draw a simple triangle hill
        half = w // 2
        for row in range(half):
            x0 = half - row - 1
            x1 = half + row + 1
            for dx in range(x0, x1):
                c = dark if row == 0 else color
                self.fb.put(sx+dx, sy+row, c)

    def draw_flagpole(self, sx):
        if sx < -2 or sx >= SW: return
        for y in range(4, GROUND_Y):
            c = FP if y % 2 == 0 else (150,150,150)
            self.fb.put(sx, y, c)
        # Flag
        for fy in range(4, 12):
            for fx in range(1, 9):
                c = FC_ if fx < 6 else FCD
                self.fb.put(sx + fx, fy, c)

    # ── Render ───────────────────────────────────────────────
    def render(self):
        fb = self.fb
        cam = self.camera_x

        # ── Sky ────────────────────────────────────────────
        for y in range(SH):
            c = SKY2 if y < 6 else SKY
            for x in range(SW):
                fb.buf[y][x] = c

        # ── Clouds (parallax at 0.6x) ──────────────────────
        for cx, cy in CLOUDS:
            sx = self._wx(cx * 0.6 + cam * 0.4 - cam * 0.4 + cx - cam * 0.35)
            # Simplified: clouds scroll at 0.5 speed
            sx = int(cx - cam * 0.5)
            self.draw_cloud(sx, cy)

        # ── Hills / bushes ─────────────────────────────────
        for hx in [0, 96, 256, 480, 640, 800, 960, 1200]:
            sx = int(hx - cam)
            self.draw_hill(sx, GROUND_Y-8, 20, HG, HGD)

        # ── Ground / Platforms ─────────────────────────────
        for px, py, pw, ptype in self.platforms:
            sx = self._wx(px)
            if sx + pw < 0 or sx >= SW: continue
            tile = {
                'ground': TILE_GROUND, 'brick': TILE_BRICK,
                'qblock': TILE_QBLOCK, 'qused': TILE_QUSED
            }.get(ptype, TILE_GROUND)
            for tx in range(pw // TILE_SIZE):
                self.draw_tile(tile, sx + tx * TILE_SIZE, py)

        # ── Pipes ──────────────────────────────────────────
        for px, th in PIPES:
            sx = self._wx(px)
            if sx + 16 < 0 or sx >= SW: continue
            pipe_py = GROUND_Y - th * 8
            # Top tile (spans 2 tiles width)
            self.draw_tile(TILE_PIPE_TOP, sx,   pipe_py)
            self.draw_tile(TILE_PIPE_TOP, sx+8, pipe_py)
            # Body
            for ti in range(1, th):
                self.draw_tile(TILE_PIPE_BODY, sx,   pipe_py + ti*8)
                self.draw_tile(TILE_PIPE_BODY, sx+8, pipe_py + ti*8)

        # ── Flagpole ───────────────────────────────────────
        self.draw_flagpole(self._wx(FLAGPOLE_X))

        # ── Coins ──────────────────────────────────────────
        for c in self.coins:
            if c.collected: continue
            sx = self._wx(c.x)
            if sx < -4 or sx >= SW: continue
            # Coin wobble
            wobble_x = 1 if (c.anim_t // 8) % 4 < 2 else 0
            for ry in range(8):
                for rx in range(4):
                    cc_color = CC if rx < 2 else CD
                    if ry == 0 or ry == 7: cc_color = None if rx == 0 or rx == 3 else cc_color
                    if cc_color:
                        fb.put(sx + rx + wobble_x, c.y + ry, cc_color)

        # ── Goombas ────────────────────────────────────────
        for g in self.goombas:
            if g.dead and g.dead_t > 25: continue
            sx = self._wx(g.x)
            if sx < -8 or sx >= SW: continue
            fb.sprite(g.sprite(), sx, int(g.y))

        # ── Mario ──────────────────────────────────────────
        mario = self.mario
        if not (mario.dead and mario.y > SH + 8):
            sx = self._wx(mario.x)
            fb.sprite(mario.sprite(), sx, int(mario.y))

        # ── Particles ──────────────────────────────────────
        for p in self.particles:
            fade = 1.0 - (p.t / p.life)
            if fade < 0.3: continue
            r,g_,b_ = p.color
            col = (int(r*fade), int(g_*fade), int(b_*fade))
            px = int(p.x)
            py = int(p.y)
            for i, ch in enumerate(p.text):
                cx = px + i
                if 0 <= cx < SW and 0 <= py < SH:
                    fb.buf[py][cx] = col

        fb.flush()

    def draw_hud(self):
        # Line 1
        score_str = f"SCORE:{self.score:06d}"
        coins_str = f"⭐{self.coin_count:02d}"
        world_str = "WORLD 1-1"
        time_str  = f"TIME:{self.time_left:03d}"
        lives_str = f"♥x{self.lives}"

        line = (f"  {BOLD}{fg(255,255,100)}{score_str}{RST}  "
                f"{fg(255,220,0)}{coins_str}{RST}  "
                f"{BOLD}{fg(255,255,255)}{world_str}{RST}  "
                f"{fg(255,100,100)}{time_str}{RST}  "
                f"{fg(255,80,80)}{lives_str}{RST}")

        _out(at(1, 1))
        _out(bg(*SKY2) + ' ' * SW + RST)
        _out(at(1, 2))
        _out(bg(*SKY2) + ' ' * SW + RST)
        _out(at(1, 1))
        _out(bg(*SKY2) + line + RST)
        _flush()

    def draw_overlay(self, title, subtitle, color=(255,255,100)):
        cx = SW // 2
        cy = T_ROWS // 2 + HUD_H
        w = max(len(title), len(subtitle)) + 6
        # Box
        box_x = cx - w//2
        box_y = cy - 2
        _out(at(box_x, box_y))
        _out(bg(20,20,80) + fg(100,100,255) + '╔' + '═'*(w) + '╗' + RST)
        _out(at(box_x, box_y+1))
        _out(bg(20,20,80) + fg(100,100,255) + '║' + RST)
        pad1 = (w - len(title)) // 2
        _out(at(box_x+1+pad1, box_y+1))
        _out(bg(20,20,80) + BOLD + fg(*color) + title + RST)
        _out(at(box_x+w, box_y+1))
        _out(bg(20,20,80) + fg(100,100,255) + '║' + RST)
        _out(at(box_x, box_y+2))
        pad2 = (w - len(subtitle)) // 2
        _out(bg(20,20,80) + fg(100,100,255) + '║' + RST)
        _out(at(box_x+1+pad2, box_y+2))
        _out(bg(20,20,80) + fg(200,200,200) + subtitle + RST)
        _out(at(box_x+w, box_y+2))
        _out(bg(20,20,80) + fg(100,100,255) + '║' + RST)
        _out(at(box_x, box_y+3))
        _out(bg(20,20,80) + fg(100,100,255) + '╚' + '═'*(w) + '╝' + RST)
        _flush()

    # ── Main Loop ────────────────────────────────────────────
    def run(self):
        _out(ALT + CLER + HIDE)
        _flush()
        # Draw border / initial clear
        self.fb.force_redraw()

        FPS = 30
        frame_time = 1.0 / FPS

        try:
            while True:
                t0 = time.time()

                # Quit check
                if self.inp.held('q') or self.inp.held('esc'):
                    break

                self.update()
                self.render()
                self.draw_hud()

                if self.state == 'over':
                    self.draw_overlay('GAME OVER', 'Press Q to quit', (255, 60, 60))
                elif self.state == 'win':
                    self.draw_overlay(' STAGE CLEAR! ', f'Score: {self.score}  Coins: {self.coin_count}', (255,255,80))

                # Frame rate cap
                elapsed = time.time() - t0
                sleep_t = frame_time - elapsed
                if sleep_t > 0:
                    time.sleep(sleep_t)

        finally:
            self.inp.restore()
            _out(NORM + SHOW + RST + CLER)
            _flush()
            print(f"\n👾 Thanks for playing FC Terminal Mario!")
            print(f"   Final Score: {self.score}  Coins: {self.coin_count}  Lives left: {self.lives}")

# ═══════════════════════════════════════════════════════════
#  ENTRY POINT
# ═══════════════════════════════════════════════════════════
def main():
    if not sys.stdin.isatty():
        print("Error: must run in an interactive terminal")
        sys.exit(1)

    # Check terminal size
    try:
        cols, rows = os.get_terminal_size()
    except:
        cols, rows = 80, 26

    if cols < SW or rows < (T_ROWS + HUD_H + 1):
        print(f"⚠️  Terminal too small! Need at least {SW}×{T_ROWS+HUD_H+1}")
        print(f"   Current: {cols}×{rows}")
        print(f"   Please resize your terminal and try again.")
        sys.exit(1)

    game = Game()
    game.run()

if __name__ == '__main__':
    main()