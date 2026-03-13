package fc.terminal.mario;

import java.io.*;
import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.concurrent.*;

/**
 * ╔══════════════════════════════════════════════════════╗
 * ║        FC TERMINAL MARIO  — Java Edition            ║
 * ║   NES-style Super Mario in the terminal             ║
 * ║                                                     ║
 * ║   Build:  javac -d out src\/**\/*.java                ║
 * ║   Run:    java -cp out fc.terminal.mario.Main       ║
 * ║                                                     ║
 * ║   Requires: JDK 17+, ANSI-capable terminal          ║
 * ║   Controls: A/D ←/→ Move   SPACE Jump   Q Quit      ║
 * ╚══════════════════════════════════════════════════════╝
 *
 * Architecture:
 *   FrameBuffer  — half-block pixel renderer (▀)
 *   InputHandler — raw terminal keyboard reader
 *   Level        — tile map, entities, decorations
 *   Mario        — player physics + animation
 *   Goomba       — enemy AI
 *   Game         — main loop, collision, rendering
 */
public class FCTerminalMario {

    // ══════════════════════════════════════════════════════════
    //  TERMINAL HELPERS
    // ══════════════════════════════════════════════════════════
    static final String ESC   = "\033";
    static final String RESET = ESC + "[0m";
    static final String BOLD  = ESC + "[1m";
    static final String HIDE  = ESC + "[?25l";
    static final String SHOW  = ESC + "[?25h";
    static final String ALT   = ESC + "[?1049h";
    static final String NORM  = ESC + "[?1049l";
    static final String CLEAR = ESC + "[2J" + ESC + "[H";

    static String at(int col, int row)   { return ESC + "[" + row + ";" + col + "H"; }
    static String fg(int r, int g, int b){ return ESC + "[38;2;" + r + ";" + g + ";" + b + "m"; }
    static String bg(int r, int g, int b){ return ESC + "[48;2;" + r + ";" + g + ";" + b + "m"; }

    static PrintStream OUT;

    // ══════════════════════════════════════════════════════════
    //  SCREEN CONFIG
    // ══════════════════════════════════════════════════════════
    static final int SW      = 80;    // screen pixel width
    static final int SH      = 44;    // screen pixel height
    static final int HUD_H   = 2;     // HUD terminal rows
    static final int T_ROWS  = SH/2;  // game terminal rows = 22
    static final int TILE    = 8;     // tile size in pixels

    // ══════════════════════════════════════════════════════════
    //  COLOR PALETTE (stored as packed int 0xRRGGBB)
    // ══════════════════════════════════════════════════════════
    static final int TRANSP = -1;

    static final int SKY   = rgb(92, 148, 252);
    static final int SKY2  = rgb(56, 100, 200);
    static final int GT    = rgb(0, 200, 0);
    static final int GM    = rgb(180, 120, 60);
    static final int GD    = rgb(120, 75, 25);
    static final int BK    = rgb(204, 102, 40);
    static final int BD    = rgb(150, 65, 15);
    static final int BH    = rgb(235, 148, 85);
    static final int BM    = rgb(85, 40, 8);
    static final int QB    = rgb(255, 192, 0);
    static final int QD    = rgb(196, 128, 0);
    static final int QQ    = rgb(255, 255, 120);
    static final int QK    = rgb(70, 35, 0);
    static final int UB    = rgb(155, 95, 45);
    static final int PL    = rgb(0, 228, 0);
    static final int PD    = rgb(0, 148, 0);
    static final int PT    = rgb(0, 255, 80);
    static final int PDK   = rgb(0, 55, 0);
    static final int CC    = rgb(255, 204, 0);
    static final int CD    = rgb(196, 140, 0);
    static final int MR    = rgb(204, 48, 48);
    static final int MS    = rgb(255, 188, 102);
    static final int MB    = rgb(72, 72, 220);
    static final int MW    = rgb(108, 60, 18);
    static final int ME    = rgb(18, 18, 18);
    static final int GB    = rgb(192, 118, 40);
    static final int GK    = rgb(96, 56, 8);
    static final int GF    = rgb(68, 32, 5);
    static final int GW    = rgb(252, 252, 252);
    static final int GP    = rgb(8, 8, 8);
    static final int CL    = rgb(255, 255, 255);
    static final int CLS   = rgb(200, 218, 255);
    static final int HG    = rgb(0, 168, 0);
    static final int HGD   = rgb(0, 120, 0);
    static final int FP    = rgb(188, 188, 188);
    static final int FC    = rgb(0, 204, 0);
    static final int FCD   = rgb(0, 140, 0);
    static final int COI   = rgb(255, 255, 80);

    static int rgb(int r, int g, int b) { return (r<<16)|(g<<8)|b; }
    static int rr(int c){ return (c>>16)&0xFF; }
    static int rg(int c){ return (c>>8 )&0xFF; }
    static int rb(int c){ return  c     &0xFF; }

    // ══════════════════════════════════════════════════════════
    //  SPRITES  (int[][] — TRANSP = -1 = transparent)
    // ══════════════════════════════════════════════════════════
    static int T = TRANSP;

    static int[][] MARIO_STAND_R = {
        {T,  T,  MR, MR, MR, MR, T,  T },
        {T,  MR, MR, MR, MR, MR, MR, T },
        {T,  MW, MS, MS, MW, MW, T,  T },
        {T,  MS, ME, MS, MS, MS, MS, T },
        {T,  MS, MS, MS, MS, MS, MS, T },
        {T,  T,  MS, MS, MS, MS, T,  T },
        {T,  MR, MB, MB, MB, MB, MR, T },
        {T,  MB, MB, MB, MB, MB, MB, T },
        {T,  MB, MB, T,  T,  MB, MB, T },
        {T,  MW, MB, T,  T,  MB, MW, T },
        {MW, MW, MW, T,  T,  MW, MW, MW},
        {T,  T,  T,  T,  T,  T,  T,  T },
    };
    static int[][] MARIO_WALK1_R = {
        {T,  T,  MR, MR, MR, MR, T,  T },
        {T,  MR, MR, MR, MR, MR, MR, T },
        {T,  MW, MS, MS, MW, MW, T,  T },
        {T,  MS, ME, MS, MS, MS, MS, T },
        {T,  MS, MS, MS, MS, MS, MS, T },
        {T,  MR, MB, MB, MB, MB, MR, T },
        {T,  MB, MB, MB, MB, MB, MB, T },
        {T,  T,  MB, T,  MB, MB, T,  T },
        {T,  MB, MW, T,  T,  MW, T,  T },
        {T,  MW, MW, T,  T,  T,  T,  T },
        {T,  T,  T,  T,  T,  T,  T,  T },
        {T,  T,  T,  T,  T,  T,  T,  T },
    };
    static int[][] MARIO_WALK2_R = {
        {T,  T,  MR, MR, MR, MR, T,  T },
        {T,  MR, MR, MR, MR, MR, MR, T },
        {T,  MW, MS, MS, MW, MW, T,  T },
        {T,  MS, MS, MS, MS, ME, MS, T },
        {T,  MS, MS, MS, MS, MS, MS, T },
        {T,  MR, MB, MB, MB, MB, MR, T },
        {T,  MB, MB, MB, MB, MB, MB, T },
        {T,  T,  MB, MB, T,  MB, T,  T },
        {T,  T,  MW, T,  T,  MW, MB, T },
        {T,  T,  T,  T,  T,  MW, MW, T },
        {T,  T,  T,  T,  T,  T,  T,  T },
        {T,  T,  T,  T,  T,  T,  T,  T },
    };
    static int[][] MARIO_JUMP_R = {
        {T,  T,  MR, MR, MR, MR, T,  T },
        {T,  MR, MR, MR, MR, MR, MR, T },
        {T,  MW, MS, MS, MW, MW, T,  T },
        {T,  MS, ME, MS, MS, MS, MS, T },
        {MS, MS, MS, MS, MS, MS, MS, MS},
        {MB, MR, MB, MB, MB, MB, MR, MB},
        {MB, MB, MB, MB, MB, MB, MB, MB},
        {T,  MB, MB, T,  T,  MB, MB, T },
        {MW, MW, T,  T,  T,  T,  MW, MW},
        {T,  T,  T,  T,  T,  T,  T,  T },
        {T,  T,  T,  T,  T,  T,  T,  T },
        {T,  T,  T,  T,  T,  T,  T,  T },
    };
    static int[][] MARIO_DEAD_R = {
        {T,  T,  T,  MR, MR, T,  T,  T },
        {T,  T,  MR, MR, MR, MR, T,  T },
        {T,  MR, MR, MR, MR, MR, MR, T },
        {T,  MW, MS, MS, MW, MW, T,  T },
        {T,  MS, ME, MS, MS, MS, MS, T },
        {MS, MS, MS, MS, MS, MS, MS, MS},
        {T,  MB, MB, MB, MB, MB, MB, T },
        {T,  MR, MB, MB, MB, MB, MR, T },
        {MW, MW, MW, T,  T,  MW, MW, MW},
        {T,  T,  T,  T,  T,  T,  T,  T },
        {T,  T,  T,  T,  T,  T,  T,  T },
        {T,  T,  T,  T,  T,  T,  T,  T },
    };

    static int[][] GOOMBA_W1 = {
        {T,  T,  GK, GK, GK, GK, T,  T },
        {T,  GK, GB, GB, GB, GB, GK, T },
        {GK, GB, GB, GW, GB, GW, GB, GK},
        {GK, GB, GK, GP, GB, GP, GB, GK},
        {GK, GB, GB, GB, GB, GB, GB, GK},
        {T,  GK, GK, GK, GK, GK, GK, T },
        {T,  GF, GF, GB, GB, GF, GF, T },
        {GF, GF, GB, T,  T,  GB, GF, GF},
        {GF, GF, GF, T,  T,  GF, GF, GF},
        {T,  T,  T,  T,  T,  T,  T,  T },
    };
    static int[][] GOOMBA_W2 = {
        {T,  T,  GK, GK, GK, GK, T,  T },
        {T,  GK, GB, GB, GB, GB, GK, T },
        {GK, GB, GB, GW, GB, GW, GB, GK},
        {GK, GB, GK, GP, GB, GP, GB, GK},
        {GK, GB, GB, GB, GB, GB, GB, GK},
        {T,  GK, GK, GK, GK, GK, GK, T },
        {T,  GB, GF, GB, GB, GF, GB, T },
        {T,  GF, GF, GB, GB, GF, GF, T },
        {T,  GF, GF, T,  T,  GF, GF, T },
        {T,  T,  T,  T,  T,  T,  T,  T },
    };
    static int[][] GOOMBA_DEAD = {
        {T,T,T,T,T,T,T,T},{T,T,T,T,T,T,T,T},{T,T,T,T,T,T,T,T},{T,T,T,T,T,T,T,T},
        {T,T,T,T,T,T,T,T},{T,T,T,T,T,T,T,T},
        {GK,GK,GK,GK,GK,GK,GK,GK},
        {GK,GB,GW,GB,GB,GW,GB,GK},
        {GK,GF,GF,GF,GF,GF,GF,GK},
        {T,T,T,T,T,T,T,T},
    };

    static int[][] flipSprite(int[][] s) {
        int[][] out = new int[s.length][];
        for (int i = 0; i < s.length; i++) {
            out[i] = new int[s[i].length];
            for (int j = 0; j < s[i].length; j++)
                out[i][j] = s[i][s[i].length - 1 - j];
        }
        return out;
    }

    static final int[][] MARIO_STAND_L = flipSprite(MARIO_STAND_R);
    static final int[][] MARIO_WALK1_L = flipSprite(MARIO_WALK1_R);
    static final int[][] MARIO_WALK2_L = flipSprite(MARIO_WALK2_R);
    static final int[][] MARIO_JUMP_L  = flipSprite(MARIO_JUMP_R);
    static final int[][] MARIO_DEAD_L  = flipSprite(MARIO_DEAD_R);

    // ══════════════════════════════════════════════════════════
    //  TILES
    // ══════════════════════════════════════════════════════════
    static final int[][] TILE_GROUND = {
        {GT,GT,GT,GT,GT,GT,GT,GT},{GT,GT,GT,GT,GT,GT,GT,GT},
        {GM,GM,GM,GM,GM,GM,GM,GM},{GM,GD,GM,GM,GD,GM,GM,GD},
        {GD,GD,GD,GD,GD,GD,GD,GD},{GD,GD,GD,GD,GD,GD,GD,GD},
        {GD,GD,GD,GD,GD,GD,GD,GD},{GD,GD,GD,GD,GD,GD,GD,GD},
    };
    static final int[][] TILE_BRICK = {
        {BM,BM,BM,BM,BM,BM,BM,BM},{BM,BH,BK,BK,BM,BH,BK,BK},
        {BM,BK,BK,BK,BM,BK,BK,BK},{BM,BK,BD,BK,BM,BK,BD,BK},
        {BM,BM,BM,BM,BM,BM,BM,BM},{BM,BK,BM,BH,BK,BK,BM,BH},
        {BM,BK,BM,BK,BK,BK,BM,BK},{BM,BD,BM,BK,BD,BK,BM,BK},
    };
    static final int[][] TILE_QBLOCK = {
        {QK,QB,QB,QB,QB,QB,QB,QK},{QB,QB,QD,QD,QB,QQ,QD,QB},
        {QB,QB,QQ,QD,QB,QQ,QD,QB},{QB,QB,QD,QQ,QB,QD,QD,QB},
        {QB,QB,QD,QQ,QB,QD,QD,QB},{QB,QB,QD,QD,QB,QQ,QD,QB},
        {QB,QD,QD,QD,QD,QD,QD,QB},{QK,QB,QB,QB,QB,QB,QB,QK},
    };
    static final int[][] TILE_QUSED = {
        {QK,UB,UB,UB,UB,UB,UB,QK},{UB,UB,BD,BD,UB,UB,BD,UB},
        {UB,UB,BD,BD,UB,UB,BD,UB},{UB,UB,BD,BD,UB,UB,BD,UB},
        {UB,UB,BD,BD,UB,UB,BD,UB},{UB,UB,BD,BD,UB,UB,BD,UB},
        {UB,BD,BD,BD,BD,BD,BD,UB},{QK,UB,UB,UB,UB,UB,UB,QK},
    };
    static final int[][] TILE_PIPE_TOP = {
        {T, PL,PT,PT,PL,PD,PDK,T },{PL,PT,PT,PL,PL,PD,PDK,PD},
        {PL,PT,PL,PL,PD,PD,PDK,PD},{PL,PL,PL,PD,PD,PDK,PDK,PD},
        {PD,PD,PD,PD,PDK,PDK,PDK,PD},{PDK,PDK,PDK,PDK,PDK,PDK,PDK,PDK},
        {T,T,T,T,T,T,T,T},{T,T,T,T,T,T,T,T},
    };
    static final int[][] TILE_PIPE_BODY = {
        {T,PL,PL,PL,PD,PDK,PDK,T},{T,PL,PL,PL,PD,PDK,PDK,T},
        {T,PL,PL,PL,PD,PDK,PDK,T},{T,PL,PL,PL,PD,PDK,PDK,T},
        {T,PL,PL,PL,PD,PDK,PDK,T},{T,PL,PL,PL,PD,PDK,PDK,T},
        {T,PL,PL,PL,PD,PDK,PDK,T},{T,PL,PL,PL,PD,PDK,PDK,T},
    };

    // ══════════════════════════════════════════════════════════
    //  FRAME BUFFER
    // ══════════════════════════════════════════════════════════
    static class FrameBuffer {
        int[][] buf  = new int[SH][SW];
        int[][] prev = new int[SH][SW];

        FrameBuffer() {
            for (int[] row : buf)  Arrays.fill(row, SKY);
            for (int[] row : prev) Arrays.fill(row, 0xFF00FF); // force redraw
        }

        void put(int x, int y, int color) {
            if (x >= 0 && x < SW && y >= 0 && y < SH && color != TRANSP)
                buf[y][x] = color;
        }

        void fill(int x, int y, int w, int h, int color) {
            for (int dy = 0; dy < h; dy++)
                for (int dx = 0; dx < w; dx++)
                    put(x+dx, y+dy, color);
        }

        void drawSprite(int[][] spr, int sx, int sy) {
            for (int ry = 0; ry < spr.length; ry++)
                for (int rx = 0; rx < spr[ry].length; rx++)
                    if (spr[ry][rx] != TRANSP)
                        put(sx+rx, sy+ry, spr[ry][rx]);
        }

        void clear(int color) {
            for (int[] row : buf) Arrays.fill(row, color);
        }

        void forceRedraw() {
            for (int[] row : prev) Arrays.fill(row, 0xFF00FF);
        }

        void flush() {
            StringBuilder sb = new StringBuilder(65536);
            sb.append(RESET);
            int lastFgR=-1,lastFgG=-1,lastFgB=-1;
            int lastBgR=-1,lastBgG=-1,lastBgB=-1;
            boolean needMove = true;

            for (int tr = 0; tr < T_ROWS; tr++) {
                int py0 = tr * 2, py1 = py0 + 1;
                int[] row0 = buf[py0];
                int[] row1 = py1 < SH ? buf[py1] : row0;
                int[] pr0  = prev[py0];
                int[] pr1  = py1 < SH ? prev[py1] : pr0;
                needMove = true;

                for (int tx = 0; tx < SW; tx++) {
                    int ct = row0[tx], cb = row1[tx];
                    if (ct == pr0[tx] && cb == pr1[tx]) { needMove=true; continue; }
                    if (needMove) {
                        sb.append(at(tx+1, tr+HUD_H+1));
                        needMove = false;
                    }
                    int fr=rr(ct),fg=rg(ct),fb2=rb(ct);
                    if (fr!=lastFgR||fg!=lastFgG||fb2!=lastFgB) {
                        sb.append(fg(fr,fg,fb2)); lastFgR=fr;lastFgG=fg;lastFgB=fb2;
                    }
                    int br=rr(cb),bgv=rg(cb),bb=rb(cb);
                    if (br!=lastBgR||bgv!=lastBgG||bb!=lastBgB) {
                        sb.append(bg(br,bgv,bb)); lastBgR=br;lastBgG=bgv;lastBgB=bb;
                    }
                    sb.append('▀');
                    pr0[tx]=ct; pr1[tx]=cb;
                }
            }
            sb.append(RESET);
            OUT.print(sb);
            OUT.flush();
        }
    }

    // ══════════════════════════════════════════════════════════
    //  INPUT HANDLER (raw terminal, non-blocking)
    // ══════════════════════════════════════════════════════════
    static class InputHandler {
        final ConcurrentHashMap<String,Long> times = new ConcurrentHashMap<>();
        volatile boolean quit = false;
        Process sttyProc;
        static final boolean IS_WINDOWS = System.getProperty("os.name").toLowerCase().contains("win");
        InputStream terminalIn;

        InputHandler() throws IOException {
            if (IS_WINDOWS) {
                terminalIn = System.in;
            } else {
                new ProcessBuilder("sh", "-c", "stty raw -echo </dev/tty").inheritIO().start();
                terminalIn = new FileInputStream("/dev/tty");
            }
            Thread t = new Thread(this::loop, "input");
            t.setDaemon(true);
            t.start();
        }

        void loop() {
            try {
                byte[] buf = new byte[8];
                while (!quit) {
                    if (terminalIn.available() > 0) {
                        int n = terminalIn.read(buf);
                        if (n > 0) parseInput(buf, n);
                    } else {
                        Thread.sleep(16);
                    }
                }
            } catch (Exception e) { /* ignore */ }
        }

        void parseInput(byte[] buf, int n) {
            if (n == 1) {
                char c = (char)(buf[0] & 0xFF);
                if (c == 27) {
                    press("esc");
                    return;
                }
                press(String.valueOf(c).toLowerCase());
            } else if (n >= 3 && buf[0]==0x1B && buf[1]=='[') {
                switch(buf[2]) {
                    case 'A': press("up");    break;
                    case 'B': press("down");  break;
                    case 'C': press("right"); break;
                    case 'D': press("left");  break;
                }
            }
        }

        void press(String k) { 
            times.put(k, System.currentTimeMillis()); 
        }

        boolean held(String k, int ms) {
            Long t = times.get(k);
            return t != null && (System.currentTimeMillis()-t) < ms;
        }
        boolean held(String k)    { return held(k, 300); }
        boolean tapped(String k)  { return held(k, 80);  }
        void consume(String k)    { times.remove(k); }

        void restore() throws IOException {
            quit = true;
            if (!IS_WINDOWS) {
                new ProcessBuilder("sh", "-c", "stty sane </dev/tty").inheritIO().start();
            }
        }
    }

    // ══════════════════════════════════════════════════════════
    //  LEVEL DATA
    // ══════════════════════════════════════════════════════════
    static final int GROUND_Y   = 36;
    static final int LEVEL_PX_W = 1600;
    static final int FLAGPOLE_X = 1520;

    // Platform: [worldX, worldY, widthPx, type] (type: 0=ground,1=brick,2=qblock,3=qused)
    static final int GROUND=0, BRICK=1, QBLOCK=2, QUSED=3;
    static List<int[]> platforms = new ArrayList<>(Arrays.asList(
        new int[]{0,   GROUND_Y, 352, GROUND},
        new int[]{400, GROUND_Y, 160, GROUND},
        new int[]{600, GROUND_Y, 240, GROUND},
        new int[]{896, GROUND_Y, 352, GROUND},
        new int[]{1296,GROUND_Y, 304, GROUND},
        new int[]{88,  16, 8,  QBLOCK}, new int[]{112,16, 8, BRICK},
        new int[]{120, 16, 8,  QBLOCK}, new int[]{128,16, 8, BRICK},
        new int[]{136, 16, 8,  QBLOCK},
        new int[]{232, 16, 8,  QBLOCK}, new int[]{240,16,16, BRICK},
        new int[]{256, 16, 8,  QBLOCK}, new int[]{264,16, 8, BRICK},
        new int[]{288,  8, 40, BRICK},  new int[]{336, 8, 8, QBLOCK},
        new int[]{344,  8, 24, BRICK},
        new int[]{440, 16, 8,  QBLOCK}, new int[]{448,16, 8, BRICK},
        new int[]{640, 16, 56, BRICK},  new int[]{696,16, 8, QBLOCK},
        new int[]{704, 16, 40, BRICK},
        new int[]{1232, GROUND_Y-8,  8, GROUND},
        new int[]{1224, GROUND_Y-16, 8, GROUND},
        new int[]{1216, GROUND_Y-24, 8, GROUND}
    ));

    // Pipes [worldX, heightTiles]
    static final int[][] PIPES = {{176,2},{264,3},{384,2},{464,3},{768,4},{896,2}};
    // Coins [worldX, worldY]
    static final int[][] COIN_POS = {
        {56,22},{64,22},{72,22},{152,22},{160,22},{168,22},
        {312,22},{320,22},{328,22},{336,22},{456,22},{464,22},
        {648,22},{656,22},{664,22},{800,14},{808,14},{816,14},
        {960,22},{968,22},{976,22},{984,22},{1072,22},{1080,22},{1088,22},
        {1160,22},{1168,22}
    };
    // Goombas [worldX, worldY, dir]
    static final int[][] GOOMBA_POS = {
        {144,GROUND_Y-10,1},{224,GROUND_Y-10,1},{304,GROUND_Y-10,-1},
        {440,GROUND_Y-10,1},{540,GROUND_Y-10,1},{660,GROUND_Y-10,-1},
        {700,GROUND_Y-10,1},{900,GROUND_Y-10,1},{960,GROUND_Y-10,1},
        {1020,GROUND_Y-10,-1},{1100,GROUND_Y-10,1},{1180,GROUND_Y-10,-1}
    };
    // Clouds [worldX, worldY]
    static final int[][] CLOUD_POS = {
        {40,4},{120,2},{200,6},{300,3},{400,5},{520,2},{620,4},
        {720,3},{820,6},{920,2},{1040,4},{1140,3},{1240,5},{1380,2}
    };

    // ══════════════════════════════════════════════════════════
    //  ENTITY CLASSES
    // ══════════════════════════════════════════════════════════
    static class Entity {
        double x, y, vx, vy;
        boolean onGround;
        Entity(double x, double y){ this.x=x; this.y=y; }
    }

    static class MarioEnt extends Entity {
        static final int W=8, H=12;
        static final double MAX_SPD=3.5, ACCEL=0.15, DECEL=0.2, JUMP_VY=-5.5, GRAV=0.45, AIR_CTRL=0.12;
        int facing=1, animT=0;
        boolean dead=false, win=false;
        double deadVy=0;
        int deadT=0;

        MarioEnt(double x, double y){ super(x,y); }

        void applyPhysics(boolean left, boolean right, boolean onGround) {
            double targetVx = 0;
            if (left && !right) targetVx = -MAX_SPD;
            else if (right && !left) targetVx = MAX_SPD;

            double accel = onGround ? ACCEL : AIR_CTRL;
            if (targetVx == 0) {
                if (vx > 0) vx = Math.max(0, vx - DECEL);
                else if (vx < 0) vx = Math.min(0, vx + DECEL);
            } else {
                if (targetVx > vx) vx = Math.min(targetVx, vx + accel);
                else vx = Math.max(targetVx, vx - accel);
            }

            if (vx != 0) facing = vx > 0 ? 1 : -1;
        }

        int[][] getSprite() {
            if (dead) return facing>0 ? MARIO_DEAD_R : MARIO_DEAD_L;
            if (!onGround) return facing>0 ? MARIO_JUMP_R : MARIO_JUMP_L;
            if (Math.abs(vx) > 0.1) {
                int f = (animT/6)%2;
                if (facing>0) return f==0?MARIO_WALK1_R:MARIO_WALK2_R;
                else          return f==0?MARIO_WALK1_L:MARIO_WALK2_L;
            }
            return facing>0 ? MARIO_STAND_R : MARIO_STAND_L;
        }
    }

    static class GoombaEnt extends Entity {
        static final int W=8, H=10;
        static final double SPD=0.8, GRAV=0.45;
        boolean dead=false;
        int deadT=0, animT=0;
        GoombaEnt(double x, double y, int dir){ super(x,y); vx=-dir*SPD; }
        int[][] getSprite() {
            if (dead) return GOOMBA_DEAD;
            return (animT/8)%2==0 ? GOOMBA_W1 : GOOMBA_W2;
        }
    }

    static class CoinEnt {
        int x, y, animT=0;
        boolean collected=false;
        CoinEnt(int x, int y){ this.x=x; this.y=y; }
    }

    static class Particle {
        double x, y, vy=-1.2;
        String text; int color; int life=28, t=0;
        Particle(double x, double y, String text, int color){
            this.x=x; this.y=y; this.text=text; this.color=color;
        }
    }

    // ══════════════════════════════════════════════════════════
    //  GAME
    // ══════════════════════════════════════════════════════════
    static class Game {
        FrameBuffer fb = new FrameBuffer();
        InputHandler inp;
        MarioEnt mario = new MarioEnt(24, GROUND_Y - MarioEnt.H);
        List<GoombaEnt> goombas = new ArrayList<>();
        List<CoinEnt>   coins   = new ArrayList<>();
        List<Particle>  parts   = new ArrayList<>();
        Set<Integer>    usedQ   = new HashSet<>();

        int score=0, coinCount=0, lives=3, frame=0, timeTick=0, timeLeft=400;
        double cameraX=0;
        String state = "play";

        Game() throws Exception {
            inp = new InputHandler();
            for (int[] g : GOOMBA_POS) goombas.add(new GoombaEnt(g[0],g[1],g[2]));
            for (int[] c : COIN_POS)   coins.add(new CoinEnt(c[0],c[1]));
        }

        List<int[]> getSolidRects() {
            List<int[]> r = new ArrayList<>();
            for (int[] p : platforms) r.add(new int[]{p[0],p[1],p[2],8});
            for (int[] p : PIPES) {
                int pipePy = GROUND_Y - p[1]*8;
                r.add(new int[]{p[0], pipePy, 16, 8});
                for (int ti=1; ti<p[1]; ti++)
                    r.add(new int[]{p[0]+2, pipePy+ti*8, 12, 8});
            }
            return r;
        }

        boolean overlap(double ax,double ay,int aw,int ah, double bx,double by,int bw,int bh){
            return ax<bx+bw && ax+aw>bx && ay<by+bh && ay+ah>by;
        }

        int resolveEntity(Entity e, int w, int h, double grav) {
            e.vy += grav;
            e.y  += e.vy;
            e.x  += e.vx;
            e.x   = Math.max(0, Math.min(e.x, LEVEL_PX_W - w));
            e.onGround = false;
            int bumpBlock = -1;

            List<int[]> rects = getSolidRects();
            for (int ri=0; ri<rects.size(); ri++) {
                int[] r = rects.get(ri);
                if (!overlap(e.x,e.y,w,h, r[0],r[1],r[2],r[3])) continue;
                double ol  = (e.x+w) - r[0];
                double or_ = (r[0]+r[2]) - e.x;
                double ot  = (e.y+h) - r[1];
                double ob  = (r[1]+r[3]) - e.y;
                double min = Math.min(Math.min(ol,or_),Math.min(ot,ob));
                if (min==ot && e.vy>=0)       { e.y=r[1]-h; e.vy=0; e.onGround=true; }
                else if (min==ob && e.vy<0)   {
                    e.y=r[1]+r[3]; e.vy=0;
                    // Check if qblock
                    for (int pi=0; pi<platforms.size(); pi++) {
                        int[] p=platforms.get(pi);
                        if (p[0]==r[0] && p[1]==r[1] && p[3]==QBLOCK && !usedQ.contains(pi))
                            bumpBlock = pi;
                    }
                }
                else if (min==ol) { e.x=r[0]-w;     e.vx=0; }
                else               { e.x=r[0]+r[2];  e.vx=0; }
            }
            return bumpBlock;
        }

        void update() {
            if (!"play".equals(state)) return;
            MarioEnt m = mario;

            if (m.dead) {
                m.deadVy += 0.5; m.y += m.deadVy; m.deadT++;
                if (m.deadT > 80) {
                    lives--;
                    if (lives <= 0) state="over"; else respawn();
                }
                return;
            }
            if (m.win) { m.x += 1.0; return; }

            // Input
            boolean left = inp.held("a") || inp.held("left");
            boolean right = inp.held("d") || inp.held("right");
            m.applyPhysics(left, right, m.onGround);

            if ((inp.tapped(" ")||inp.tapped("w")||inp.tapped("up")) && m.onGround) {
                m.vy=MarioEnt.JUMP_VY;
                inp.consume(" "); inp.consume("w"); inp.consume("up");
            }

            int bump = resolveEntity(m, MarioEnt.W, MarioEnt.H, MarioEnt.GRAV);
            if (bump >= 0) {
                usedQ.add(bump);
                int[] p = platforms.get(bump);
                p[3] = QUSED;
                score += 100; coinCount++;
                parts.add(new Particle(p[0]-cameraX, p[1]-10, "+100", COI));
            }
            if (m.y > SH+10) killMario();

            if (Math.abs(m.vx) > 0.1 || !m.onGround) m.animT++;

            // Goombas
            for (GoombaEnt g : goombas) {
                if (g.dead) { g.deadT++; continue; }
                g.animT++;
                resolveEntity(g, GoombaEnt.W, GoombaEnt.H, GoombaEnt.GRAV);

                // Edge detection
                boolean onEdge = true;
                int sx=(int)g.x, sy=(int)g.y+GoombaEnt.H+1;
                for (int[] r : getSolidRects())
                    if (r[0]<=sx && sx<r[0]+r[2] && r[1]<=sy && sy<r[1]+r[3]+1) { onEdge=false; break; }
                if (onEdge || g.x<=0 || g.x>=LEVEL_PX_W-GoombaEnt.W) g.vx=-g.vx;

                if (!m.dead && overlap(m.x,m.y,MarioEnt.W,MarioEnt.H, g.x,g.y,GoombaEnt.W,GoombaEnt.H)) {
                    if (m.vy>0 && m.y+MarioEnt.H < g.y+GoombaEnt.H-2) {
                        g.dead=true; m.vy=-3.5; score+=200;
                        parts.add(new Particle(g.x-cameraX+4, g.y-8, "+200", 0xFFFFFF));
                    } else killMario();
                }
            }

            // Coins
            for (CoinEnt c : coins) {
                if (c.collected) continue;
                c.animT++;
                if (overlap(m.x,m.y,MarioEnt.W,MarioEnt.H, c.x,c.y,4,8)) {
                    c.collected=true; coinCount++; score+=50;
                    parts.add(new Particle(c.x-cameraX, c.y-6, "+50", COI));
                }
            }

            if (!m.win && m.x+MarioEnt.W >= FLAGPOLE_X) { m.win=true; state="win"; }

            parts.removeIf(p -> p.t >= p.life);
            for (Particle p : parts) { p.y+=p.vy; p.t++; }

            // Camera
            double target = m.x - SW/3.0;
            cameraX = Math.max(0, Math.min(target, LEVEL_PX_W-SW));

            timeTick++;
            if (timeTick >= 60) { timeTick=0; timeLeft=Math.max(0,timeLeft-1); if(timeLeft==0)killMario(); }

            frame++;
        }

        void killMario() {
            if (!mario.dead) { mario.dead=true; mario.vy=0; mario.deadVy=-6.5; mario.vx=0; }
        }

        void respawn() {
            mario=new MarioEnt(24, GROUND_Y-MarioEnt.H);
            cameraX=0; state="play"; timeLeft=400; timeTick=0;
            fb.forceRedraw();
        }

        int wx(double worldX){ return (int)(worldX - cameraX); }

        void drawTile(int[][] tile, int sx, int sy) {
            for (int ry=0; ry<tile.length; ry++)
                for (int rx=0; rx<tile[ry].length; rx++)
                    if (tile[ry][rx]!=TRANSP) fb.put(sx+rx, sy+ry, tile[ry][rx]);
        }

        void drawCloud(int sx, int sy) {
            int[] widths={4,6,8,8,6}; int[] offX={2,1,0,0,1};
            for (int row=0; row<5; row++)
                for (int dx=0; dx<widths[row]; dx++)
                    fb.put(sx+offX[row]+dx, sy+row, (dx+row)%3==0?CLS:CL);
        }

        void drawHill(int sx, int sy, int w, int color, int dark) {
            int half=w/2;
            for (int row=0; row<half; row++) {
                int x0=half-row-1, x1=half+row+1;
                for (int dx=x0; dx<x1; dx++) fb.put(sx+dx, sy+row, row==0?dark:color);
            }
        }

        void drawFlagpole(int sx) {
            if (sx<-2 || sx>=SW) return;
            for (int y=4; y<GROUND_Y; y++) fb.put(sx, y, y%2==0?FP:rgb(150,150,150));
            for (int fy=4; fy<12; fy++)
                for (int fx=1; fx<9; fx++) fb.put(sx+fx, fy, fx<6?FC:FCD);
        }

        void render() {
            // Sky
            for (int y=0; y<SH; y++) { int c=y<6?SKY2:SKY; Arrays.fill(fb.buf[y],c); }

            // Clouds
            for (int[] cl : CLOUD_POS) drawCloud((int)(cl[0]-cameraX*0.5), cl[1]);

            // Hills
            for (int hx : new int[]{0,96,256,480,640,800,960,1200})
                drawHill(wx(hx), GROUND_Y-8, 20, HG, HGD);

            // Platforms
            for (int[] p : platforms) {
                int sx=wx(p[0]);
                if (sx+p[2]<0||sx>=SW) continue;
                int[][] tile = p[3]==GROUND?TILE_GROUND : p[3]==BRICK?TILE_BRICK :
                               p[3]==QBLOCK?TILE_QBLOCK : TILE_QUSED;
                for (int tx=0; tx<p[2]/TILE; tx++) drawTile(tile, sx+tx*TILE, p[1]);
            }

            // Pipes
            for (int[] p : PIPES) {
                int sx=wx(p[0]);
                if (sx+16<0||sx>=SW) continue;
                int pipePy=GROUND_Y-p[1]*8;
                drawTile(TILE_PIPE_TOP, sx, pipePy);
                drawTile(TILE_PIPE_TOP, sx+8, pipePy);
                for (int ti=1; ti<p[1]; ti++) {
                    drawTile(TILE_PIPE_BODY, sx, pipePy+ti*8);
                    drawTile(TILE_PIPE_BODY, sx+8, pipePy+ti*8);
                }
            }

            // Flagpole
            drawFlagpole(wx(FLAGPOLE_X));

            // Coins
            for (CoinEnt c : coins) {
                if (c.collected) continue;
                int sx=wx(c.x);
                if (sx<-4||sx>=SW) continue;
                int wobble=(c.animT/8)%4<2?1:0;
                for (int ry=0; ry<8; ry++)
                    for (int rx=0; rx<4; rx++) {
                        if ((ry==0||ry==7)&&(rx==0||rx==3)) continue;
                        fb.put(sx+rx+wobble, c.y+ry, rx<2?CC:CD);
                    }
            }

            // Goombas
            for (GoombaEnt g : goombas) {
                if (g.dead && g.deadT>25) continue;
                int sx=wx(g.x);
                if (sx<-8||sx>=SW) continue;
                fb.drawSprite(g.getSprite(), sx, (int)g.y);
            }

            // Mario
            MarioEnt m=mario;
            if (!(m.dead && m.y>SH+8)) fb.drawSprite(m.getSprite(), wx(m.x), (int)m.y);

            // Particles
            for (Particle p : parts) {
                double fade=1.0-(double)p.t/p.life;
                if (fade<0.3) continue;
                int c=p.color;
                int col=rgb((int)(rr(c)*fade),(int)(rg(c)*fade),(int)(rb(c)*fade));
                int px=(int)p.x, py=(int)p.y;
                for (int i=0; i<p.text.length(); i++)
                    if (px+i>=0&&px+i<SW&&py>=0&&py<SH) fb.buf[py][px+i]=col;
            }

            fb.flush();
        }

        void drawHud() {
            String line = String.format("  %sSCORE:%06d%s  %s⭐%02d%s  %sWORLD 1-1%s  %sTIME:%03d%s  %s♥x%d%s",
                BOLD+fg(255,255,100), score, RESET,
                fg(255,220,0), coinCount, RESET,
                BOLD+fg(255,255,255), RESET,
                fg(255,100,100), timeLeft, RESET,
                fg(255,80,80), lives, RESET);
            StringBuilder sb = new StringBuilder();
            sb.append(at(1,1)).append(bg(rr(SKY2),rg(SKY2),rb(SKY2)));
            for (int i=0; i<SW; i++) sb.append(' ');
            sb.append(RESET).append(at(1,2)).append(bg(rr(SKY2),rg(SKY2),rb(SKY2)));
            for (int i=0; i<SW; i++) sb.append(' ');
            sb.append(RESET).append(at(1,1)).append(bg(rr(SKY2),rg(SKY2),rb(SKY2))).append(line).append(RESET);
            OUT.print(sb);
            OUT.flush();
        }

        void drawOverlay(String title, String sub, int[] color) {
            int cx=SW/2, cy=T_ROWS/2+HUD_H;
            int w=Math.max(title.length(),sub.length())+6;
            int bx=cx-w/2, by=cy-2;
            StringBuilder sb=new StringBuilder();
            sb.append(at(bx,by))  .append(bg(20,20,80)).append(fg(100,100,255)).append('╔').append("═".repeat(w)).append('╗').append(RESET);
            sb.append(at(bx,by+1)).append(bg(20,20,80)).append(fg(100,100,255)).append('║').append(RESET);
            int p1=(w-title.length())/2;
            sb.append(at(bx+1+p1,by+1)).append(bg(20,20,80)).append(BOLD).append(fg(color[0],color[1],color[2])).append(title).append(RESET);
            sb.append(at(bx+w,by+1)).append(bg(20,20,80)).append(fg(100,100,255)).append('║').append(RESET);
            sb.append(at(bx,by+2)).append(bg(20,20,80)).append(fg(100,100,255)).append('║').append(RESET);
            int p2=(w-sub.length())/2;
            sb.append(at(bx+1+p2,by+2)).append(bg(20,20,80)).append(fg(200,200,200)).append(sub).append(RESET);
            sb.append(at(bx+w,by+2)).append(bg(20,20,80)).append(fg(100,100,255)).append('║').append(RESET);
            sb.append(at(bx,by+3)).append(bg(20,20,80)).append(fg(100,100,255)).append('╚').append("═".repeat(w)).append('╝').append(RESET);
            OUT.print(sb); OUT.flush();
        }

        void run() throws Exception {
            OUT.print(ALT + CLEAR + HIDE); OUT.flush();
            fb.forceRedraw();
            final int FPS=30;
            final long FRAME_NS=1_000_000_000L/FPS;

            while (true) {
                long t0=System.nanoTime();
                if (inp.held("q") || inp.held("\u001b")) break;
                update(); render(); drawHud();
                if ("over".equals(state)) drawOverlay("GAME OVER","Press Q to quit",new int[]{255,60,60});
                else if ("win".equals(state)) drawOverlay(" STAGE CLEAR! ","Score:"+score+" Coins:"+coinCount,new int[]{255,255,80});
                long elapsed=System.nanoTime()-t0;
                long sleep=(FRAME_NS-elapsed)/1_000_000;
                if (sleep>0) Thread.sleep(sleep);
            }

            inp.restore();
            OUT.print(NORM+SHOW+RESET+CLEAR); OUT.flush();
            System.out.println("\n👾 Thanks for playing FC Terminal Mario!");
            System.out.printf("   Final Score: %d  Coins: %d  Lives: %d%n", score, coinCount, lives);
        }
    }

    // ══════════════════════════════════════════════════════════
    //  ENTRY POINT
    // ══════════════════════════════════════════════════════════
    public static void main(String[] args) throws Exception {
        System.setOut(new PrintStream(System.out, true, StandardCharsets.UTF_8));
        OUT = new PrintStream(
            new BufferedOutputStream(System.out, 131072),
            false, StandardCharsets.UTF_8
        );
        OUT.print(ESC + "[?1049h" + ESC + "[?25l" + CLEAR);
        System.out.println("Starting FC Terminal Mario...");
        Thread.sleep(500);
        new Game().run();
    }
}