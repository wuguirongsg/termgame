// ╔══════════════════════════════════════════════════════╗
// ║        FC TERMINAL MARIO  —  Go Edition             ║
// ║   NES风格超级玛丽 · 完整终端游戏                      ║
// ║                                                     ║
// ║   Build:  go build -o mario . && ./mario            ║
// ║   Controls: A/← D/→ Move   SPACE Jump   Q Quit      ║
// ║   Requires: 80×25+ terminal, 24-bit color           ║
// ╚══════════════════════════════════════════════════════╝

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// ══════════════════════════════════════════════════════════
//  SCREEN CONFIG
// ══════════════════════════════════════════════════════════

const (
	SW         = 80 // screen pixel width
	SH         = 44 // screen pixel height  (half-block: each char = 2px rows)
	HUD_H      = 2  // HUD rows at top
	T_ROWS     = SH / 2
	TILE       = 8
	LEVEL_PX_W = 1600
	GROUND_Y   = 36
	FLAGPOLE_X = 1520
)

// ══════════════════════════════════════════════════════════
//  COLORS  (packed 0xRRGGBB, TRANSP = -1)
// ══════════════════════════════════════════════════════════

const TRANSP = -1

func rgb(r, g, b int) int { return (r << 16) | (g << 8) | b }
func cR(c int) int        { return (c >> 16) & 0xFF }
func cG(c int) int        { return (c >> 8) & 0xFF }
func cB(c int) int        { return c & 0xFF }

// short alias for transparent — used inline in sprite literals
const X = TRANSP

var (
	SKY = rgb(92, 148, 252)
	SK2 = rgb(56, 100, 200)
	GT  = rgb(0, 200, 0)
	GM  = rgb(180, 120, 60)
	GD  = rgb(120, 75, 25)
	BD  = rgb(150, 65, 15)
	BH  = rgb(235, 148, 85)
	BM  = rgb(85, 40, 8)
	QB  = rgb(255, 192, 0)
	QD  = rgb(196, 128, 0)
	QQ  = rgb(255, 255, 120)
	QK  = rgb(70, 35, 0)
	UB  = rgb(155, 95, 45)
	PL  = rgb(0, 228, 0)
	PD  = rgb(0, 148, 0)
	PT  = rgb(0, 255, 80)
	PDK = rgb(0, 55, 0)
	CC  = rgb(255, 204, 0)
	CD  = rgb(196, 140, 0)
	MR  = rgb(204, 48, 48)
	MS  = rgb(255, 188, 102)
	MB  = rgb(72, 72, 220)
	MW  = rgb(108, 60, 18)
	ME  = rgb(18, 18, 18)
	GBc = rgb(192, 118, 40)
	GKc = rgb(96, 56, 8)
	GFc = rgb(68, 32, 5)
	GWc = rgb(252, 252, 252)
	GPc = rgb(8, 8, 8)
	CL  = rgb(255, 255, 255)
	CLS = rgb(200, 218, 255)
	HG  = rgb(0, 168, 0)
	HGD = rgb(0, 120, 0)
	FP  = rgb(188, 188, 188)
	FG2 = rgb(150, 150, 150)
	FC  = rgb(0, 204, 0)
	FCD = rgb(0, 140, 0)
	COI = rgb(255, 255, 80)
)

// ══════════════════════════════════════════════════════════
//  SPRITES  ([][]int, X = transparent)
// ══════════════════════════════════════════════════════════

func flipSpr(s [][]int) [][]int {
	out := make([][]int, len(s))
	for i, row := range s {
		r := make([]int, len(row))
		for j, v := range row {
			r[len(row)-1-j] = v
		}
		out[i] = r
	}
	return out
}

var mStandR = [][]int{
	{X, X, MR, MR, MR, MR, X, X},
	{X, MR, MR, MR, MR, MR, MR, X},
	{X, MW, MS, MS, MW, MW, X, X},
	{X, MS, ME, MS, MS, MS, MS, X},
	{X, MS, MS, MS, MS, MS, MS, X},
	{X, X, MS, MS, MS, MS, X, X},
	{X, MR, MB, MB, MB, MB, MR, X},
	{X, MB, MB, MB, MB, MB, MB, X},
	{X, MB, MB, X, X, MB, MB, X},
	{X, MW, MB, X, X, MB, MW, X},
	{MW, MW, MW, X, X, MW, MW, MW},
	{X, X, X, X, X, X, X, X},
}
var mWalk1R = [][]int{
	{X, X, MR, MR, MR, MR, X, X},
	{X, MR, MR, MR, MR, MR, MR, X},
	{X, MW, MS, MS, MW, MW, X, X},
	{X, MS, ME, MS, MS, MS, MS, X},
	{X, MS, MS, MS, MS, MS, MS, X},
	{X, MR, MB, MB, MB, MB, MR, X},
	{X, MB, MB, MB, MB, MB, MB, X},
	{X, X, MB, X, MB, MB, X, X},
	{X, MB, MW, X, X, MW, X, X},
	{X, MW, MW, X, X, X, X, X},
	{X, X, X, X, X, X, X, X},
	{X, X, X, X, X, X, X, X},
}
var mWalk2R = [][]int{
	{X, X, MR, MR, MR, MR, X, X},
	{X, MR, MR, MR, MR, MR, MR, X},
	{X, MW, MS, MS, MW, MW, X, X},
	{X, MS, MS, MS, MS, ME, MS, X},
	{X, MS, MS, MS, MS, MS, MS, X},
	{X, MR, MB, MB, MB, MB, MR, X},
	{X, MB, MB, MB, MB, MB, MB, X},
	{X, X, MB, MB, X, MB, X, X},
	{X, X, MW, X, X, MW, MB, X},
	{X, X, X, X, X, MW, MW, X},
	{X, X, X, X, X, X, X, X},
	{X, X, X, X, X, X, X, X},
}
var mJumpR = [][]int{
	{X, X, MR, MR, MR, MR, X, X},
	{X, MR, MR, MR, MR, MR, MR, X},
	{X, MW, MS, MS, MW, MW, X, X},
	{X, MS, ME, MS, MS, MS, MS, X},
	{MS, MS, MS, MS, MS, MS, MS, MS},
	{MB, MR, MB, MB, MB, MB, MR, MB},
	{MB, MB, MB, MB, MB, MB, MB, MB},
	{X, MB, MB, X, X, MB, MB, X},
	{MW, MW, X, X, X, X, MW, MW},
	{X, X, X, X, X, X, X, X},
	{X, X, X, X, X, X, X, X},
	{X, X, X, X, X, X, X, X},
}
var mDeadR = [][]int{
	{X, X, X, MR, MR, X, X, X},
	{X, X, MR, MR, MR, MR, X, X},
	{X, MR, MR, MR, MR, MR, MR, X},
	{X, MW, MS, MS, MW, MW, X, X},
	{X, MS, ME, MS, MS, MS, MS, X},
	{MS, MS, MS, MS, MS, MS, MS, MS},
	{X, MB, MB, MB, MB, MB, MB, X},
	{X, MR, MB, MB, MB, MB, MR, X},
	{MW, MW, MW, X, X, MW, MW, MW},
	{X, X, X, X, X, X, X, X},
	{X, X, X, X, X, X, X, X},
	{X, X, X, X, X, X, X, X},
}

var (
	mStandL = flipSpr(mStandR)
	mWalk1L = flipSpr(mWalk1R)
	mWalk2L = flipSpr(mWalk2R)
	mJumpL  = flipSpr(mJumpR)
	mDeadL  = flipSpr(mDeadR)
)

var goombaW1 = [][]int{
	{X, X, GKc, GKc, GKc, GKc, X, X},
	{X, GKc, GBc, GBc, GBc, GBc, GKc, X},
	{GKc, GBc, GBc, GWc, GBc, GWc, GBc, GKc},
	{GKc, GBc, GKc, GPc, GBc, GPc, GBc, GKc},
	{GKc, GBc, GBc, GBc, GBc, GBc, GBc, GKc},
	{X, GKc, GKc, GKc, GKc, GKc, GKc, X},
	{X, GFc, GFc, GBc, GBc, GFc, GFc, X},
	{GFc, GFc, GBc, X, X, GBc, GFc, GFc},
	{GFc, GFc, GFc, X, X, GFc, GFc, GFc},
	{X, X, X, X, X, X, X, X},
}
var goombaW2 = [][]int{
	{X, X, GKc, GKc, GKc, GKc, X, X},
	{X, GKc, GBc, GBc, GBc, GBc, GKc, X},
	{GKc, GBc, GBc, GWc, GBc, GWc, GBc, GKc},
	{GKc, GBc, GKc, GPc, GBc, GPc, GBc, GKc},
	{GKc, GBc, GBc, GBc, GBc, GBc, GBc, GKc},
	{X, GKc, GKc, GKc, GKc, GKc, GKc, X},
	{X, GBc, GFc, GBc, GBc, GFc, GBc, X},
	{X, GFc, GFc, GBc, GBc, GFc, GFc, X},
	{X, GFc, GFc, X, X, GFc, GFc, X},
	{X, X, X, X, X, X, X, X},
}
var goombaDead = [][]int{
	{X, X, X, X, X, X, X, X},
	{X, X, X, X, X, X, X, X},
	{X, X, X, X, X, X, X, X},
	{X, X, X, X, X, X, X, X},
	{X, X, X, X, X, X, X, X},
	{X, X, X, X, X, X, X, X},
	{GKc, GKc, GKc, GKc, GKc, GKc, GKc, GKc},
	{GKc, GBc, GWc, GBc, GBc, GWc, GBc, GKc},
	{GKc, GFc, GFc, GFc, GFc, GFc, GFc, GKc},
	{X, X, X, X, X, X, X, X},
}

// ══════════════════════════════════════════════════════════
//  TILES (8×8)
// ══════════════════════════════════════════════════════════

var tileGround = [][]int{
	{GT, GT, GT, GT, GT, GT, GT, GT},
	{GT, GT, GT, GT, GT, GT, GT, GT},
	{GM, GM, GM, GM, GM, GM, GM, GM},
	{GM, GD, GM, GM, GD, GM, GM, GD},
	{GD, GD, GD, GD, GD, GD, GD, GD},
	{GD, GD, GD, GD, GD, GD, GD, GD},
	{GD, GD, GD, GD, GD, GD, GD, GD},
	{GD, GD, GD, GD, GD, GD, GD, GD},
}
var tileBrick = [][]int{
	{BM, BM, BM, BM, BM, BM, BM, BM},
	{BM, BH, BD, BD, BM, BH, BD, BD}, // BK replaced with BD for valid Go var
	{BM, BD, BD, BD, BM, BD, BD, BD},
	{BM, BD, BD, BD, BM, BD, BD, BD},
	{BM, BM, BM, BM, BM, BM, BM, BM},
	{BM, BD, BM, BH, BD, BD, BM, BH},
	{BM, BD, BM, BD, BD, BD, BM, BD},
	{BM, BD, BM, BD, BD, BD, BM, BD},
}
var tileQBlock = [][]int{
	{QK, QB, QB, QB, QB, QB, QB, QK},
	{QB, QB, QD, QD, QB, QQ, QD, QB},
	{QB, QB, QQ, QD, QB, QQ, QD, QB},
	{QB, QB, QD, QQ, QB, QD, QD, QB},
	{QB, QB, QD, QQ, QB, QD, QD, QB},
	{QB, QB, QD, QD, QB, QQ, QD, QB},
	{QB, QD, QD, QD, QD, QD, QD, QB},
	{QK, QB, QB, QB, QB, QB, QB, QK},
}
var tileQUsed = [][]int{
	{QK, UB, UB, UB, UB, UB, UB, QK},
	{UB, UB, BD, BD, UB, UB, BD, UB},
	{UB, UB, BD, BD, UB, UB, BD, UB},
	{UB, UB, BD, BD, UB, UB, BD, UB},
	{UB, UB, BD, BD, UB, UB, BD, UB},
	{UB, UB, BD, BD, UB, UB, BD, UB},
	{UB, BD, BD, BD, BD, BD, BD, UB},
	{QK, UB, UB, UB, UB, UB, UB, QK},
}
var tilePipeTop = [][]int{
	{X, PL, PT, PT, PL, PD, PDK, X},
	{PL, PT, PT, PL, PL, PD, PDK, PD},
	{PL, PT, PL, PL, PD, PD, PDK, PD},
	{PL, PL, PL, PD, PD, PDK, PDK, PD},
	{PD, PD, PD, PD, PDK, PDK, PDK, PD},
	{PDK, PDK, PDK, PDK, PDK, PDK, PDK, PDK},
	{X, X, X, X, X, X, X, X},
	{X, X, X, X, X, X, X, X},
}
var tilePipeBody = [][]int{
	{X, PL, PL, PL, PD, PDK, PDK, X},
	{X, PL, PL, PL, PD, PDK, PDK, X},
	{X, PL, PL, PL, PD, PDK, PDK, X},
	{X, PL, PL, PL, PD, PDK, PDK, X},
	{X, PL, PL, PL, PD, PDK, PDK, X},
	{X, PL, PL, PL, PD, PDK, PDK, X},
	{X, PL, PL, PL, PD, PDK, PDK, X},
	{X, PL, PL, PL, PD, PDK, PDK, X},
}

// ══════════════════════════════════════════════════════════
//  LEVEL DATA
// ══════════════════════════════════════════════════════════

type TileType int

const (
	TGround TileType = iota
	TBrick
	TQBlock
	TQUsed
)

type Platform struct {
	x, y, w int
	typ     TileType
}
type Pipe struct{ x, h int }
type Vec2i struct{ x, y int }
type GoombaInit struct{ x, y, dir int }

var initPlatforms = []Platform{
	{0, GROUND_Y, 352, TGround},
	{400, GROUND_Y, 160, TGround},
	{600, GROUND_Y, 240, TGround},
	{896, GROUND_Y, 352, TGround},
	{1296, GROUND_Y, 304, TGround},
	// Floating blocks
	{88, 16, 8, TQBlock},
	{112, 16, 8, TBrick},
	{120, 16, 8, TQBlock},
	{128, 16, 8, TBrick},
	{136, 16, 8, TQBlock},
	{232, 16, 8, TQBlock},
	{240, 16, 16, TBrick},
	{256, 16, 8, TQBlock},
	{264, 16, 8, TBrick},
	{288, 8, 40, TBrick},
	{336, 8, 8, TQBlock},
	{344, 8, 24, TBrick},
	{440, 16, 8, TQBlock},
	{448, 16, 8, TBrick},
	{640, 16, 56, TBrick},
	{696, 16, 8, TQBlock},
	{704, 16, 40, TBrick},
	// Staircase
	{1232, GROUND_Y - 8, 8, TGround},
	{1224, GROUND_Y - 16, 8, TGround},
	{1216, GROUND_Y - 24, 8, TGround},
}

var pipes = []Pipe{{176, 2}, {264, 3}, {384, 2}, {464, 3}, {768, 4}, {896, 2}}

var initCoins = []Vec2i{
	{56, 22}, {64, 22}, {72, 22}, {152, 22}, {160, 22}, {168, 22},
	{312, 22}, {320, 22}, {328, 22}, {336, 22}, {456, 22}, {464, 22},
	{648, 22}, {656, 22}, {664, 22}, {800, 14}, {808, 14}, {816, 14},
	{960, 22}, {968, 22}, {976, 22}, {984, 22},
	{1072, 22}, {1080, 22}, {1088, 22}, {1160, 22}, {1168, 22},
}

var initGoombas = []GoombaInit{
	{144, GROUND_Y - 10, 1}, {224, GROUND_Y - 10, 1}, {304, GROUND_Y - 10, -1},
	{440, GROUND_Y - 10, 1}, {540, GROUND_Y - 10, 1}, {660, GROUND_Y - 10, -1},
	{700, GROUND_Y - 10, 1}, {900, GROUND_Y - 10, 1}, {960, GROUND_Y - 10, 1},
	{1020, GROUND_Y - 10, -1}, {1100, GROUND_Y - 10, 1}, {1180, GROUND_Y - 10, -1},
}

var clouds = []Vec2i{
	{40, 4}, {120, 2}, {200, 6}, {300, 3}, {400, 5}, {520, 2},
	{620, 4}, {720, 3}, {820, 6}, {920, 2}, {1040, 4}, {1140, 3}, {1240, 5}, {1380, 2},
}

var hillXs = []int{0, 96, 256, 480, 640, 800, 960, 1200}

// ══════════════════════════════════════════════════════════
//  FRAME BUFFER  (half-block pixel renderer)
//  Each terminal row = 2 pixel rows.  Character '▀':
//    fg = upper pixel color,  bg = lower pixel color
// ══════════════════════════════════════════════════════════

type FrameBuffer struct {
	buf  [SH][SW]int
	prev [SH][SW]int
	w    *bufio.Writer
}

func newFB() *FrameBuffer {
	fb := &FrameBuffer{w: bufio.NewWriterSize(os.Stdout, 131072)}
	for y := range fb.buf {
		for x := range fb.buf[y] {
			fb.buf[y][x] = SKY
			fb.prev[y][x] = 0xFF00FF // force first redraw
		}
	}
	return fb
}

func (fb *FrameBuffer) put(x, y, color int) {
	if x >= 0 && x < SW && y >= 0 && y < SH && color != TRANSP {
		fb.buf[y][x] = color
	}
}

func (fb *FrameBuffer) sprite(spr [][]int, sx, sy int) {
	for ry, row := range spr {
		for rx, c := range row {
			if c != TRANSP {
				fb.put(sx+rx, sy+ry, c)
			}
		}
	}
}

func (fb *FrameBuffer) tile(t [][]int, sx, sy int) {
	for ry, row := range t {
		for rx, c := range row {
			if c != TRANSP {
				fb.put(sx+rx, sy+ry, c)
			}
		}
	}
}

func (fb *FrameBuffer) forceRedraw() {
	for y := range fb.prev {
		for x := range fb.prev[y] {
			fb.prev[y][x] = 0xFF00FF
		}
	}
}

func (fb *FrameBuffer) flush() {
	w := fb.w
	w.WriteString("\x1b[0m")
	lfR, lfG, lfB := -1, -1, -1
	lbR, lbG, lbB := -1, -1, -1
	needMove := true

	for tr := 0; tr < T_ROWS; tr++ {
		py0 := tr * 2
		py1 := py0 + 1
		if py1 >= SH {
			py1 = py0
		}
		needMove = true
		for tx := 0; tx < SW; tx++ {
			ct := fb.buf[py0][tx]
			cb := fb.buf[py1][tx]
			if ct == fb.prev[py0][tx] && cb == fb.prev[py1][tx] {
				needMove = true
				continue
			}
			if needMove {
				fmt.Fprintf(w, "\x1b[%d;%dH", tr+HUD_H+1, tx+1)
				needMove = false
			}
			fr, fg, fb2 := cR(ct), cG(ct), cB(ct)
			if fr != lfR || fg != lfG || fb2 != lfB {
				fmt.Fprintf(w, "\x1b[38;2;%d;%d;%dm", fr, fg, fb2)
				lfR, lfG, lfB = fr, fg, fb2
			}
			br, bg, bb := cR(cb), cG(cb), cB(cb)
			if br != lbR || bg != lbG || bb != lbB {
				fmt.Fprintf(w, "\x1b[48;2;%d;%d;%dm", br, bg, bb)
				lbR, lbG, lbB = br, bg, bb
			}
			w.WriteString("▀")
			fb.prev[py0][tx] = ct
			fb.prev[py1][tx] = cb
		}
	}
	w.WriteString("\x1b[0m")
	w.Flush()
}

// ══════════════════════════════════════════════════════════
//  INPUT  (goroutine reading keyboard)
// ══════════════════════════════════════════════════════════

type Input struct {
	mu    sync.Mutex
	times map[string]time.Time
	done  chan struct{}
}

func newInput() *Input {
	inp := &Input{times: make(map[string]time.Time), done: make(chan struct{})}
	go inp.loop()
	return inp
}

func (inp *Input) loop() {
	// openInputDevice() is implemented per-platform in terminal_unix.go / terminal_windows.go
	dev := openInputDevice()
	buf := make([]byte, 16)
	for {
		select {
		case <-inp.done:
			return
		default:
		}
		n, _ := dev.Read(buf)
		if n == 0 {
			continue
		}
		inp.parse(buf[:n])
	}
}

func (inp *Input) parse(b []byte) {
	if len(b) == 1 {
		ch := b[0]
		if ch == ' ' {
			inp.press(" ")
		} else {
			inp.press(strings.ToLower(string(ch)))
		}
	} else if len(b) >= 3 && b[0] == 0x1b && b[1] == '[' {
		switch b[2] {
		case 'A':
			inp.press("up")
		case 'B':
			inp.press("down")
		case 'C':
			inp.press("right")
		case 'D':
			inp.press("left")
		}
	}
}

func (inp *Input) press(k string) {
	inp.mu.Lock()
	inp.times[k] = time.Now()
	inp.mu.Unlock()
}

func (inp *Input) held(k string, dur time.Duration) bool {
	inp.mu.Lock()
	t, ok := inp.times[k]
	inp.mu.Unlock()
	return ok && time.Since(t) < dur
}

func (inp *Input) heldD(k string) bool  { return inp.held(k, 140*time.Millisecond) }
func (inp *Input) tapped(k string) bool { return inp.held(k, 80*time.Millisecond) }

func (inp *Input) consume(k string) {
	inp.mu.Lock()
	delete(inp.times, k)
	inp.mu.Unlock()
}

func (inp *Input) stop() { close(inp.done) }

// ══════════════════════════════════════════════════════════
//  ENTITIES
// ══════════════════════════════════════════════════════════

type Mario struct {
	x, y, vx, vy float64
	onGround     bool
	facing       int // 1=right -1=left
	animT        int
	dead         bool
	win          bool
	deadVy       float64
	deadT        int
}

func newMario() *Mario {
	return &Mario{x: 24, y: GROUND_Y - 12, facing: 1}
}

func (m *Mario) spr() [][]int {
	if m.dead {
		if m.facing > 0 {
			return mDeadR
		}
		return mDeadL
	}
	if !m.onGround {
		if m.facing > 0 {
			return mJumpR
		}
		return mJumpL
	}
	if m.vx != 0 {
		f := (m.animT / 6) % 2
		if m.facing > 0 {
			if f == 0 {
				return mWalk1R
			}
			return mWalk2R
		}
		if f == 0 {
			return mWalk1L
		}
		return mWalk2L
	}
	if m.facing > 0 {
		return mStandR
	}
	return mStandL
}

const mW, mH = 8, 12
const mSpd = 2.0
const mJvy = -5.2
const mGrav = 0.45

type Goomba struct {
	x, y, vx, vy float64
	onGround     bool
	dead         bool
	deadT, animT int
}

func newGoomba(gi GoombaInit) *Goomba {
	return &Goomba{x: float64(gi.x), y: float64(gi.y), vx: -float64(gi.dir) * 0.8}
}

func (g *Goomba) spr() [][]int {
	if g.dead {
		return goombaDead
	}
	if (g.animT/8)%2 == 0 {
		return goombaW1
	}
	return goombaW2
}

const gW, gH = 8, 10
const gGrav = 0.45

type Coin struct {
	x, y      int
	collected bool
	animT     int
}

type Particle struct {
	x, y  float64
	vy    float64
	text  string
	color int
	life  int
	t     int
}

func newPart(x, y float64, text string, color int) *Particle {
	return &Particle{x: x, y: y, vy: -1.2, text: text, color: color, life: 28}
}

// ══════════════════════════════════════════════════════════
//  COLLISION RECT
// ══════════════════════════════════════════════════════════

type Rect struct{ x, y, w, h int }

func overlaps(ax, ay float64, aw, ah int, b Rect) bool {
	return ax < float64(b.x+b.w) && ax+float64(aw) > float64(b.x) &&
		ay < float64(b.y+b.h) && ay+float64(ah) > float64(b.y)
}

// ══════════════════════════════════════════════════════════
//  GAME
// ══════════════════════════════════════════════════════════

type Game struct {
	fb        *FrameBuffer
	inp       *Input
	mario     *Mario
	goombas   []*Goomba
	coins     []*Coin
	platforms []Platform
	usedQ     map[int]bool
	particles []*Particle

	score, coinCount, lives int
	timeLeft, timeTick      int
	cameraX                 float64
	state                   string // "play"|"over"|"win"
	frame                   int
}

func newGame() *Game {
	g := &Game{
		fb:        newFB(),
		inp:       newInput(),
		mario:     newMario(),
		platforms: make([]Platform, len(initPlatforms)),
		usedQ:     make(map[int]bool),
		lives:     3,
		timeLeft:  400,
		state:     "play",
	}
	copy(g.platforms, initPlatforms)
	for _, gi := range initGoombas {
		g.goombas = append(g.goombas, newGoomba(gi))
	}
	for _, c := range initCoins {
		g.coins = append(g.coins, &Coin{x: c.x, y: c.y})
	}
	return g
}

// solidRects returns all collision rectangles in world space.
func (g *Game) solidRects() []Rect {
	out := make([]Rect, 0, len(g.platforms)+len(pipes)*5)
	for _, p := range g.platforms {
		out = append(out, Rect{p.x, p.y, p.w, 8})
	}
	for _, p := range pipes {
		py := GROUND_Y - p.h*8
		out = append(out, Rect{p.x, py, 16, 8})
		for ti := 1; ti < p.h; ti++ {
			out = append(out, Rect{p.x + 2, py + ti*8, 12, 8})
		}
	}
	return out
}

// resolve moves entity by velocity, applies gravity, collides with world.
// Returns index of bumped QBlock, or -1.
func (g *Game) resolve(
	ex, ey, evx, evy *float64,
	onGround *bool,
	ew, eh int,
	grav float64,
) int {
	*evy += grav
	*ey += *evy
	*ex += *evx

	maxX := float64(LEVEL_PX_W - ew)
	if *ex < 0 {
		*ex = 0
	} else if *ex > maxX {
		*ex = maxX
	}

	*onGround = false
	bump := -1

	for _, r := range g.solidRects() {
		if !overlaps(*ex, *ey, ew, eh, r) {
			continue
		}
		ol := (*ex + float64(ew)) - float64(r.x)
		or_ := float64(r.x+r.w) - *ex
		ot := (*ey + float64(eh)) - float64(r.y)
		ob := float64(r.y+r.h) - *ey

		mn := min(min(ol, or_), min(ot, ob))
		switch {
		case mn == ot && *evy >= 0: // landing on top
			*ey = float64(r.y - eh)
			*evy = 0
			*onGround = true
		case mn == ob && *evy < 0: // bumping from below
			*ey = float64(r.y + r.h)
			*evy = 0
			for pi, p := range g.platforms {
				if p.x == r.x && p.y == r.y && p.typ == TQBlock && !g.usedQ[pi] {
					bump = pi
				}
			}
		case mn == ol: // hitting left wall
			*ex = float64(r.x - ew)
			*evx = 0
		default: // hitting right wall
			*ex = float64(r.x + r.w)
			*evx = 0
		}
	}
	return bump
}

// ── Update ──────────────────────────────────────────────────

func (g *Game) update() {
	if g.state != "play" {
		return
	}
	m := g.mario
	inp := g.inp

	// ── Dead Mario arc ────────────────────────────────────
	if m.dead {
		m.deadVy += 0.5
		m.y += m.deadVy
		m.deadT++
		if m.deadT > 80 {
			g.lives--
			if g.lives <= 0 {
				g.state = "over"
			} else {
				g.respawn()
			}
		}
		return
	}
	if m.win {
		m.x += 1.0
		return
	}

	// ── Player input ──────────────────────────────────────
	moving := false
	if inp.heldD("a") || inp.heldD("left") {
		m.vx = -mSpd
		m.facing = -1
		moving = true
	} else if inp.heldD("d") || inp.heldD("right") {
		m.vx = mSpd
		m.facing = 1
		moving = true
	} else {
		m.vx = 0
	}

	if (inp.tapped(" ") || inp.tapped("w") || inp.tapped("up")) && m.onGround {
		m.vy = mJvy
		inp.consume(" ")
		inp.consume("w")
		inp.consume("up")
	}

	// ── Physics ───────────────────────────────────────────
	bump := g.resolve(&m.x, &m.y, &m.vx, &m.vy, &m.onGround, mW, mH, mGrav)
	if bump >= 0 {
		g.usedQ[bump] = true
		g.platforms[bump].typ = TQUsed
		g.score += 100
		g.coinCount++
		px := float64(g.platforms[bump].x) - g.cameraX
		py := float64(g.platforms[bump].y) - 10
		g.particles = append(g.particles, newPart(px, py, "+100", COI))
	}
	if m.y > SH+10 {
		g.killMario()
	}
	if moving || !m.onGround {
		m.animT++
	}

	// ── Goombas ───────────────────────────────────────────
	rects := g.solidRects()
	for _, gb := range g.goombas {
		if gb.dead {
			gb.deadT++
			continue
		}
		gb.animT++
		g.resolve(&gb.x, &gb.y, &gb.vx, &gb.vy, &gb.onGround, gW, gH, gGrav)

		// Edge / boundary reversal
		sx, sy := int(gb.x), int(gb.y)+gH+1
		onEdge := true
		for _, r := range rects {
			if r.x <= sx && sx < r.x+r.w && r.y <= sy && sy < r.y+r.h+1 {
				onEdge = false
				break
			}
		}
		if onEdge || gb.x <= 0 || gb.x >= float64(LEVEL_PX_W-gW) {
			gb.vx = -gb.vx
		}

		// Collision with Mario
		if !m.dead && overlaps(m.x, m.y, mW, mH, Rect{int(gb.x), int(gb.y), gW, gH}) {
			if m.vy > 0 && m.y+float64(mH) < gb.y+float64(gH)-2 {
				// Mario stomps goomba
				gb.dead = true
				m.vy = -3.5
				g.score += 200
				g.particles = append(g.particles,
					newPart(gb.x-g.cameraX+4, gb.y-8, "+200", rgb(255, 255, 255)))
			} else {
				g.killMario()
			}
		}
	}

	// ── Coins ─────────────────────────────────────────────
	for _, c := range g.coins {
		if c.collected {
			continue
		}
		c.animT++
		if overlaps(m.x, m.y, mW, mH, Rect{c.x, c.y, 4, 8}) {
			c.collected = true
			g.coinCount++
			g.score += 50
			g.particles = append(g.particles,
				newPart(float64(c.x)-g.cameraX, float64(c.y)-6, "+50", COI))
		}
	}

	// ── Flagpole ──────────────────────────────────────────
	if !m.win && m.x+float64(mW) >= FLAGPOLE_X {
		m.win = true
		g.state = "win"
	}

	// ── Particles ─────────────────────────────────────────
	live := g.particles[:0]
	for _, p := range g.particles {
		p.y += p.vy
		p.t++
		if p.t < p.life {
			live = append(live, p)
		}
	}
	g.particles = live

	// ── Camera ────────────────────────────────────────────
	target := m.x - SW/3.0
	if target < 0 {
		target = 0
	}
	if target > LEVEL_PX_W-SW {
		target = LEVEL_PX_W - SW
	}
	g.cameraX = target

	// ── Timer ─────────────────────────────────────────────
	g.timeTick++
	if g.timeTick >= 60 {
		g.timeTick = 0
		if g.timeLeft > 0 {
			g.timeLeft--
		}
		if g.timeLeft == 0 {
			g.killMario()
		}
	}
	g.frame++
}

func (g *Game) killMario() {
	m := g.mario
	if !m.dead {
		m.dead = true
		m.vy = 0
		m.deadVy = -6.5
		m.vx = 0
	}
}

func (g *Game) respawn() {
	g.mario = newMario()
	g.cameraX = 0
	g.state = "play"
	g.timeLeft = 400
	g.timeTick = 0
	g.fb.forceRedraw()
}

// ── Render helpers ──────────────────────────────────────────

func (g *Game) wx(worldX float64) int { return int(worldX - g.cameraX) }
func (g *Game) wxi(worldX int) int    { return worldX - int(g.cameraX) }

func (g *Game) drawCloud(sx, sy int) {
	type rowDef struct{ ox, w int }
	rows := [5]rowDef{{2, 4}, {1, 6}, {0, 8}, {0, 8}, {1, 6}}
	for oy, r := range rows {
		for dx := 0; dx < r.w; dx++ {
			c := CL
			if (dx+oy)%3 == 0 {
				c = CLS
			}
			g.fb.put(sx+r.ox+dx, sy+oy, c)
		}
	}
}

func (g *Game) drawHill(sx, sy, w, color, dark int) {
	half := w / 2
	for row := 0; row < half; row++ {
		c := color
		if row == 0 {
			c = dark
		}
		for dx := half - row - 1; dx < half+row+1; dx++ {
			g.fb.put(sx+dx, sy+row, c)
		}
	}
}

func (g *Game) drawFlagpole(sx int) {
	if sx < -2 || sx >= SW {
		return
	}
	for y := 4; y < GROUND_Y; y++ {
		c := FP
		if y%2 != 0 {
			c = FG2
		}
		g.fb.put(sx, y, c)
	}
	for fy := 4; fy < 12; fy++ {
		for fx := 1; fx < 9; fx++ {
			c := FC
			if fx >= 6 {
				c = FCD
			}
			g.fb.put(sx+fx, fy, c)
		}
	}
}

// ── Main render ─────────────────────────────────────────────

func (g *Game) render() {
	fb := g.fb

	// Sky gradient
	for y := 0; y < SH; y++ {
		c := SKY
		if y < 6 {
			c = SK2
		}
		for x := 0; x < SW; x++ {
			fb.buf[y][x] = c
		}
	}

	// Clouds — parallax at 0.5x
	for _, cl := range clouds {
		sx := int(float64(cl.x) - g.cameraX*0.5)
		g.drawCloud(sx, cl.y)
	}

	// Hills
	for _, hx := range hillXs {
		g.drawHill(g.wxi(hx), GROUND_Y-8, 20, HG, HGD)
	}

	// Platforms
	for _, p := range g.platforms {
		sx := g.wxi(p.x)
		if sx+p.w < 0 || sx >= SW {
			continue
		}
		var t [][]int
		switch p.typ {
		case TGround:
			t = tileGround
		case TBrick:
			t = tileBrick
		case TQBlock:
			t = tileQBlock
		case TQUsed:
			t = tileQUsed
		}
		for tx := 0; tx < p.w/TILE; tx++ {
			fb.tile(t, sx+tx*TILE, p.y)
		}
	}

	// Pipes
	for _, p := range pipes {
		sx := g.wxi(p.x)
		if sx+16 < 0 || sx >= SW {
			continue
		}
		py := GROUND_Y - p.h*8
		fb.tile(tilePipeTop, sx, py)
		fb.tile(tilePipeTop, sx+8, py)
		for ti := 1; ti < p.h; ti++ {
			fb.tile(tilePipeBody, sx, py+ti*8)
			fb.tile(tilePipeBody, sx+8, py+ti*8)
		}
	}

	// Flagpole
	g.drawFlagpole(g.wx(FLAGPOLE_X))

	// Coins with wobble animation
	for _, c := range g.coins {
		if c.collected {
			continue
		}
		sx := g.wxi(c.x)
		if sx < -4 || sx >= SW {
			continue
		}
		wobble := 0
		if (c.animT/8)%4 < 2 {
			wobble = 1
		}
		for ry := 0; ry < 8; ry++ {
			for rx := 0; rx < 4; rx++ {
				if (ry == 0 || ry == 7) && (rx == 0 || rx == 3) {
					continue
				}
				col := CC
				if rx >= 2 {
					col = CD
				}
				fb.put(sx+rx+wobble, c.y+ry, col)
			}
		}
	}

	// Goombas
	for _, gb := range g.goombas {
		if gb.dead && gb.deadT > 25 {
			continue
		}
		sx := g.wx(gb.x)
		if sx < -8 || sx >= SW {
			continue
		}
		fb.sprite(gb.spr(), sx, int(gb.y))
	}

	// Mario
	m := g.mario
	if !(m.dead && m.y > float64(SH)+8) {
		fb.sprite(m.spr(), g.wx(m.x), int(m.y))
	}

	// Particles (fade-out text)
	for _, p := range g.particles {
		fade := 1.0 - float64(p.t)/float64(p.life)
		if fade < 0.3 {
			continue
		}
		col := rgb(
			int(float64(cR(p.color))*fade),
			int(float64(cG(p.color))*fade),
			int(float64(cB(p.color))*fade),
		)
		px, py := int(p.x), int(p.y)
		for i := 0; i < len(p.text); i++ {
			cx := px + i
			if cx >= 0 && cx < SW && py >= 0 && py < SH {
				fb.buf[py][cx] = col
			}
		}
	}

	fb.flush()
}

// ── HUD ─────────────────────────────────────────────────────

func (g *Game) drawHUD() {
	w := g.fb.w
	sky := fmt.Sprintf("\x1b[48;2;%d;%d;%dm", cR(SK2), cG(SK2), cB(SK2))

	line := fmt.Sprintf(
		"  \x1b[1m\x1b[38;2;255;255;100mSCORE:%06d\x1b[0m"+
			"  \x1b[38;2;255;220;0m\u2b50%02d\x1b[0m"+
			"  \x1b[1m\x1b[38;2;255;255;255mWORLD 1-1\x1b[0m"+
			"  \x1b[38;2;255;100;100mTIME:%03d\x1b[0m"+
			"  \x1b[38;2;255;80;80m\u2665x%d\x1b[0m",
		g.score, g.coinCount, g.timeLeft, g.lives,
	)

	// Row 1: fill
	w.WriteString("\x1b[1;1H")
	w.WriteString(sky)
	w.WriteString(strings.Repeat(" ", SW))
	w.WriteString("\x1b[0m")
	// Row 2: fill
	w.WriteString("\x1b[2;1H")
	w.WriteString(sky)
	w.WriteString(strings.Repeat(" ", SW))
	w.WriteString("\x1b[0m")
	// Row 1: content
	w.WriteString("\x1b[1;1H")
	w.WriteString(sky)
	w.WriteString(line)
	w.WriteString("\x1b[0m")
	w.Flush()
}

// ── Overlay (game over / win) ────────────────────────────────

func (g *Game) drawOverlay(title, sub string, r, gv, b int) {
	w := g.fb.w
	boxW := max(len([]rune(title)), len([]rune(sub))) + 6
	bx := SW/2 - boxW/2
	by := T_ROWS/2 + HUD_H - 2

	border := "\x1b[48;2;20;20;80m\x1b[38;2;100;100;255m"
	titleFg := fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, gv, b)

	// Top border
	fmt.Fprintf(w, "\x1b[%d;%dH%s╔%s╗\x1b[0m",
		by, bx, border, strings.Repeat("═", boxW))
	// Title row
	p1 := (boxW - len([]rune(title))) / 2
	fmt.Fprintf(w, "\x1b[%d;%dH%s║\x1b[0m", by+1, bx, border)
	fmt.Fprintf(w, "\x1b[%d;%dH\x1b[48;2;20;20;80m\x1b[1m%s%s\x1b[0m",
		by+1, bx+1+p1, titleFg, title)
	fmt.Fprintf(w, "\x1b[%d;%dH%s║\x1b[0m", by+1, bx+boxW, border)
	// Subtitle row
	p2 := (boxW - len([]rune(sub))) / 2
	fmt.Fprintf(w, "\x1b[%d;%dH%s║\x1b[0m", by+2, bx, border)
	fmt.Fprintf(w, "\x1b[%d;%dH\x1b[48;2;20;20;80m\x1b[38;2;200;200;200m%s\x1b[0m",
		by+2, bx+1+p2, sub)
	fmt.Fprintf(w, "\x1b[%d;%dH%s║\x1b[0m", by+2, bx+boxW, border)
	// Bottom border
	fmt.Fprintf(w, "\x1b[%d;%dH%s╚%s╝\x1b[0m",
		by+3, bx, border, strings.Repeat("═", boxW))
	w.Flush()
}

// ══════════════════════════════════════════════════════════
//  MAIN LOOP
// ══════════════════════════════════════════════════════════

func (g *Game) run() {
	// Enter alternate screen, clear, hide cursor
	fmt.Print("\x1b[?1049h\x1b[2J\x1b[H\x1b[?25l")
	g.fb.forceRedraw()

	defer func() {
		g.inp.stop()
		restoreTerminal()
		fmt.Print("\x1b[?1049l\x1b[?25h\x1b[0m\x1b[2J\x1b[H")
		fmt.Println("\n👾 Thanks for playing FC Terminal Mario!  (Go Edition)")
		fmt.Printf("   Final Score: %d  Coins: %d  Lives: %d\n",
			g.score, g.coinCount, g.lives)
	}()

	const fps = 30
	tick := time.NewTicker(time.Second / fps)
	defer tick.Stop()

	for range tick.C {
		// Quit
		if g.inp.heldD("q") || g.inp.heldD("\x1b") {
			break
		}

		g.update()
		g.render()
		g.drawHUD()

		switch g.state {
		case "over":
			g.drawOverlay("GAME OVER", "Press Q to quit", 255, 60, 60)
		case "win":
			g.drawOverlay(" STAGE CLEAR! ",
				fmt.Sprintf("Score:%d  Coins:%d", g.score, g.coinCount),
				255, 255, 80)
		}
	}
}

// ══════════════════════════════════════════════════════════
//  ENTRY POINT
// ══════════════════════════════════════════════════════════

// termSize and enableRaw/restoreTerminal are in terminal_unix.go / terminal_windows.go

func main() {
	cols, rows := termSize()
	if cols > 0 && (cols < SW || rows < T_ROWS+HUD_H+1) {
		fmt.Printf("⚠️  Terminal too small! Need at least %dx%d\n", SW, T_ROWS+HUD_H+1)
		fmt.Printf("   Current size: %dx%d\n", cols, rows)
		os.Exit(1)
	}
	enableRaw()
	newGame().run()
}
