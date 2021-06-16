package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	rl "github.com/lachee/raylib-goplus/raylib"
)

var ( // MARK: var ███████████████████████████████

	// help
	helpon bool

	//endgame
	spacepause = 30
	endscore   int
	endfade    = float32(0.0)
	endnumber  = 1
	endscreen  bool
	endshapes  = make([]endshape, drawa)
	// start game
	startselect                     int
	startgame                       bool
	newgameoptions                  = make([]bool, 6)
	invincible, grayscale, supafast bool
	//options
	backgon      bool
	optionsonoff = make([]bool, 4)
	optionselect int
	// crate
	outlineblocks, circleblocks                             bool
	explodetimer                                            int
	exploderecs                                             = make([]exploderec, 8)
	crateactiv, crateon                                     bool
	cratetimer, cratetimercount, crateblock, cratetimernext int
	// backgrounds
	backtype    int
	currentback = make([]backg, 100)
	// imgs
	hart       = rl.NewRectangle(0, 40, 42, 36)
	coin       = rl.NewRectangle(0, 16, 16, 16)
	crate      = rl.NewRectangle(0, 0, 16, 16)
	rayliblogo = rl.NewRectangle(507, 2, 130, 130)
	gologo     = rl.NewRectangle(355, 0, 150, 140)
	// intro
	introon     bool
	introy      = -200
	introcolor1 rl.Color
	//snyk
	gametimecount                          int
	gametime                               = 120
	hppause                                bool
	hppausetimer                           int
	hp                                     = 3
	collectblok                            int
	multiplier                             = 1
	autosnyk, autosnykpause, autosnykcrate bool
	snykcount, autosnyktimer               int
	snyk                                   = make([]playerblok, drawa)
	//room
	blockw     = 15
	roomw      = 70
	roomh      = 40
	rooma      = roomw * roomh
	roomlayout = make([]blok, drawa)
	// core
	options, paused, scanlines, pixel_noise, ghosting          bool
	centerblok, drawblok, nextdrawblok, draww, drawh, drawa    int
	mouseblok                                                  int
	mousepos                                                   rl.Vector2
	gridon, debugon, fadeblinkon                               bool
	monw, monh                                                 int
	fps                                                        = 30
	framecount                                                 int
	imgs                                                       rl.Texture2D
	camera, cameraend                                          rl.Camera2D
	fadeblink                                                  = float32(0.2)
	onoff2, onoff3, onoff6, onoff10, onoff15, onoff30, onoff60 bool
)

type endshape struct {
	v2    rl.Vector2
	w     int
	color rl.Color
}
type exploderec struct {
	rec       rl.Rectangle
	color     rl.Color
	opac      float32
	direction int
}

type backg struct {
	rec                                     rl.Rectangle
	color                                   rl.Color
	sides, w, h, direction, speed, x, y     int
	rotating, poly, moving, resizing, lines bool
	opac, rotation                          float32
}
type blok struct {
	x, y                                        int
	special, crate, activ, solid, collect, snyk bool
	color, color2                               rl.Color
}
type playerblok struct {
	activ, bounce                               bool
	color                                       rl.Color
	x, y, blocknumber, previousblock, direction int
}

func raylib() { // MARK: raylib
	rl.InitWindow(monw, monh, "GAME TITLE")
	rl.ToggleFullscreen()
	rl.SetExitKey(rl.KeyEnd) // key to end the game and close window
	// MARK: load images
	imgs = rl.LoadTexture("imgs.png") // load images
	paused = true
	introon = true
	rl.SetTargetFPS(fps)
	//rl.HideCursor()
	//	rl.ToggleFullscreen()
	for !rl.WindowShouldClose() {
		framecount++
		mousepos = rl.GetMousePosition()
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		drawnocameraback()
		rl.BeginMode2D(camera)
		if !paused {
			drawlayers()
		}
		if gridon {
			drawgrid()
		}

		rl.EndMode2D()
		drawnocamera()

		if debugon {
			drawdebug()
		}
		update()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}

func update() { // MARK: update

	if supafast {
		fps = 45
		rl.SetTargetFPS(fps)
	} else {
		fps = 30
		rl.SetTargetFPS(fps)
	}

	input()
	timers()
	if !paused {
		updateplayer()
		if hp <= 0 && !autosnyk {
			endscreen = true
		}
	}
	if crateon {
		createcrate()
	}
	if autosnyk {

		if !autosnykpause {

			if autosnykcrate {
				h, w := crateblock/draww, crateblock%draww
				hplayer, wplayer := snyk[0].blocknumber/draww, snyk[0].blocknumber%draww
				if flipcoin() {
					if h > hplayer {
						if snyk[0].direction != 1 {
							snyk[0].direction = 3
						} else {
							if flipcoin() {
								snyk[0].direction = 2
							} else {
								snyk[0].direction = 4
							}
						}
					} else if h < hplayer {
						if snyk[0].direction != 3 {
							snyk[0].direction = 1
						} else {
							if flipcoin() {
								snyk[0].direction = 2
							} else {
								snyk[0].direction = 4
							}
						}
					}
				} else {
					if w < wplayer {
						if snyk[0].direction != 2 {
							snyk[0].direction = 4
						} else {
							if flipcoin() {
								snyk[0].direction = 1
							} else {
								snyk[0].direction = 3
							}
						}
					} else if w > wplayer {
						if snyk[0].direction != 4 {
							snyk[0].direction = 2
						} else {
							if flipcoin() {
								snyk[0].direction = 1
							} else {
								snyk[0].direction = 3
							}
						}
					}
				}
			} else {
				h, w := collectblok/draww, collectblok%draww
				hplayer, wplayer := snyk[0].blocknumber/draww, snyk[0].blocknumber%draww
				if flipcoin() {
					if h > hplayer {
						if snyk[0].direction != 1 {
							snyk[0].direction = 3
						} else {
							if flipcoin() {
								snyk[0].direction = 2
							} else {
								snyk[0].direction = 4
							}
						}
					} else if h < hplayer {
						if snyk[0].direction != 3 {
							snyk[0].direction = 1
						} else {
							if flipcoin() {
								snyk[0].direction = 2
							} else {
								snyk[0].direction = 4
							}
						}
					}
				} else {
					if w < wplayer {
						if snyk[0].direction != 2 {
							snyk[0].direction = 4
						} else {
							if flipcoin() {
								snyk[0].direction = 1
							} else {
								snyk[0].direction = 3
							}
						}
					} else if w > wplayer {
						if snyk[0].direction != 4 {
							snyk[0].direction = 2
						} else {
							if flipcoin() {
								snyk[0].direction = 1
							} else {
								snyk[0].direction = 3
							}
						}
					}
				}
			}
			if rolldice()+rolldice() == 12 {
				snyk[0].direction = rInt(1, 5)
				autosnykpause = true
				autosnyktimer = 30
			}
		}

	}
}
func drawnocameraback() { // MARK: drawnocameraback
	if backgon {
		for a := 0; a < len(currentback); a++ {
			if currentback[a].poly {
				v2 := rl.NewVector2(float32(currentback[a].x), float32(currentback[a].y))
				rl.DrawPoly(v2, currentback[a].sides, float32(currentback[a].w/2), currentback[a].rotation, rl.Fade(currentback[a].color, currentback[a].opac))
			} else {
				rl.DrawRectangle(currentback[a].x, currentback[a].y, currentback[a].w, currentback[a].h, rl.Fade(currentback[a].color, currentback[a].opac))
			}
			if currentback[a].moving {
				switch currentback[a].direction {
				case 7:
					currentback[a].x -= rInt(8, 12)
					currentback[a].y -= rInt(8, 12)
				case 8:
					currentback[a].y -= rInt(8, 12)
				case 9:
					currentback[a].x += rInt(8, 12)
					currentback[a].y -= rInt(8, 12)
				case 4:
					currentback[a].x -= rInt(8, 12)
				case 5:
					currentback[a].x += rInt(8, 12)
				case 1:
					currentback[a].x -= rInt(8, 12)
					currentback[a].y += rInt(8, 12)
				case 2:
					currentback[a].y += rInt(8, 12)
				case 3:
					currentback[a].x += rInt(8, 12)
					currentback[a].y += rInt(8, 12)
				}
				if currentback[a].x > monw {
					currentback[a].x = 0
				}
				if currentback[a].x < 0 {
					currentback[a].x = monw
				}
				if currentback[a].y > monh {
					currentback[a].y = 0
				}
				if currentback[a].y < 0 {
					currentback[a].y = monh
				}

				if rolldice()+rolldice() > 10 {
					for {
						currentback[a].direction = rInt(1, 10)
						if currentback[a].direction != 5 {
							break
						}
					}
				}
			}
			if currentback[a].rotating {
				currentback[a].rotation += rFloat32(5, 11)
			}
			if currentback[a].resizing {
				currentback[a].w += rInt(1, 4)
				if currentback[a].w > 30 {
					currentback[a].w = rInt(5, 8)
				}
			}
		}
	}
	pausey := 0
	if rolldice() == 6 {
		pausey += rInt(-10, 21)
	}
	if paused && !options && !startgame && !endscreen {
		rl.DrawText("paused", monw/2-200, monh/2-40+pausey, 100, rl.Fade(rl.White, 0.3))
		if rolldice() == 6 {
			rl.DrawText("paused", monw/2-200, monh/2-100+pausey, 100, rl.Fade(randomcolor(), 0.5))
		}
	} else if !paused {
		// hp
		v2 := rl.NewVector2(float32(monw/2-70), float32(monh/2-100))
		for a := 0; a < hp; a++ {
			rl.DrawTextureRec(imgs, hart, v2, rl.Fade(darkred(), 0.8))
			if ghosting {
				ghostv2 := rl.NewVector2(v2.X+rFloat32(-4, 5), v2.Y+rFloat32(-4, 5))
				rl.DrawTextureRec(imgs, hart, ghostv2, rl.Fade(darkred(), 0.8))
			}
			v2.X += 50

		}
		// score
		snykcounttxt := strconv.Itoa(snykcount)
		multipliertxt := strconv.Itoa(multiplier)
		rl.DrawText("collected", monw/2-20, monh/2-40, 20, rl.Fade(rl.White, 0.3))
		if rolldice() == 6 {
			rl.DrawText("collected", monw/2-20, monh/2-40, 20, rl.Fade(randomcolor(), 0.5))
		}
		if snykcount < 20 {
			rl.DrawText(snykcounttxt, monw/2-70, monh/2-40, 40, rl.Fade(rl.White, 0.3))
			if rolldice() == 6 {
				rl.DrawText(snykcounttxt, monw/2-70, monh/2-40, 40, rl.Fade(randomcolor(), 0.5))
			}
		} else {
			rl.DrawText(snykcounttxt, monw/2-80, monh/2-40, 40, rl.Fade(rl.White, 0.3))
			if rolldice() == 6 {
				rl.DrawText(snykcounttxt, monw/2-80, monh/2-40, 40, rl.Fade(randomcolor(), 0.3))
			}
		}

		rl.DrawText("multiplier", monw/2-20, monh/2, 20, rl.Fade(rl.White, 0.3))
		if rolldice() == 6 {
			rl.DrawText("multiplier", monw/2-20, monh/2, 20, rl.Fade(randomcolor(), 0.5))
		}
		if multiplier < 20 {
			rl.DrawText(multipliertxt, monw/2-70, monh/2, 40, rl.Fade(rl.White, 0.3))
			if rolldice() == 6 {
				rl.DrawText(multipliertxt, monw/2-70, monh/2, 40, rl.Fade(randomcolor(), 0.5))
			}
		} else {
			rl.DrawText(multipliertxt, monw/2-80, monh/2, 40, rl.Fade(rl.White, 0.3))
			if rolldice() == 6 {
				rl.DrawText(multipliertxt, monw/2-80, monh/2, 40, rl.Fade(randomcolor(), 0.5))
			}

		}
		//time
		/*	gametimetxt := strconv.Itoa(gametime)
			rl.DrawText(gametimetxt, monw/2-20, monh/2+40, 60, rl.Fade(rl.White, 0.3))
			if rolldice() == 6 {
				rl.DrawText(gametimetxt, monw/2-20, monh/2+40, 60, rl.Fade(randomcolor(), 0.5))
			}
		*/
	}
}
func drawlayers() { // MARK: drawlayers

	// layer 1
	count := 0
	x, y := 0, 0
	drawblok = 0
	for a := 0; a < drawa; a++ {
		roomlayout[a].x = x
		roomlayout[a].y = y

		if roomlayout[a].activ && !roomlayout[a].crate && !roomlayout[a].collect {

			blokcolor := roomlayout[a].color
			if grayscale {
				blokcolor = roomlayout[a].color2
			}
			if outlineblocks && !circleblocks {
				rl.DrawRectangleLines(x, y, blockw, blockw, blokcolor)
				if ghosting {
					rl.DrawRectangleLines(x+rInt(-3, 4), y+rInt(-3, 4), blockw, blockw, rl.Fade(blokcolor, rF32(0.3, 1.1)))
				}
			} else if circleblocks && outlineblocks {
				rl.DrawCircleLines(x+8, y+8, float32(blockw/2), blokcolor)
				if ghosting {
					rl.DrawCircleLines((x+8)+rInt(-3, 4), (y+8)+rInt(-3, 4), float32(blockw/2), rl.Fade(blokcolor, rF32(0.3, 1.1)))
				}
			} else if circleblocks && !outlineblocks {
				rl.DrawCircle(x+8, y+8, float32(blockw/2), blokcolor)
				if ghosting {
					rl.DrawCircle((x+8)+rInt(-3, 4), (y+8)+rInt(-3, 4), float32(blockw/2), rl.Fade(blokcolor, rF32(0.3, 1.1)))
				}
			} else {
				rl.DrawRectangle(x, y, blockw, blockw, blokcolor)
				if ghosting {
					rl.DrawRectangle(x+rInt(-3, 4), y+rInt(-3, 4), blockw, blockw, rl.Fade(blokcolor, rF32(0.3, 1.1)))
				}
			}
		}
		if roomlayout[a].collect {
			destrec := rl.NewRectangle(float32(x-4), float32(y-4), 24, 24)
			origin := rl.NewVector2(0, 0)

			rl.DrawTexturePro(imgs, coin, destrec, origin, 0, brightyellow())
			if ghosting {
				destrec.X += rFloat32(-3, 4)
				destrec.Y += rFloat32(-3, 4)
				rl.DrawTexturePro(imgs, coin, destrec, origin, 0, rl.Fade(brightyellow(), rF32(0.3, 1.1)))
			}
			if snyk[0].blocknumber == a {
				collect(a)
			}
		}
		if roomlayout[a].crate {

			if snyk[0].blocknumber == a {
				createspecial()
				explode(x, y)
				autosnykcrate = false
				multiplier++
			}

			dest := rl.NewRectangle(float32(x+4), float32(y+4), 24, 24)
			org := rl.NewVector2(float32(8), float32(8))
			rotation := float32(0)
			if rolldice()+rolldice() > 10 {
				rotation += rFloat32(-15, 16)
				change := rInt(3, 6)
				dest.Height += float32(change)
				dest.Width += float32(change)
			}
			rl.DrawTexturePro(imgs, crate, dest, org, rotation, roomlayout[a].color)
			if ghosting {
				dest.X += rFloat32(-3, 4)
				dest.Y += rFloat32(-3, 4)
				rl.DrawTexturePro(imgs, crate, dest, org, rotation, rl.Fade(roomlayout[a].color, rF32(0.3, 1.1)))
			}
		}

		x += 16
		count++
		drawblok++
		if count == draww {
			count = 0
			x = 0
			y += 16

		}
	}

	if explodetimer != 0 {

		for b := 0; b < len(exploderecs); b++ {

			rl.DrawRectangleRec(exploderecs[b].rec, rl.Fade(exploderecs[b].color, exploderecs[b].opac))
			switch exploderecs[b].direction {
			case 7:
				exploderecs[b].rec.X -= rFloat32(1, 3)
				exploderecs[b].rec.Y -= rFloat32(1, 3)
			case 8:
				exploderecs[b].rec.Y -= rFloat32(1, 3)
			case 9:
				exploderecs[b].rec.X += rFloat32(1, 3)
				exploderecs[b].rec.Y -= rFloat32(1, 3)
			case 4:
				exploderecs[b].rec.X -= rFloat32(2, 4)
			case 6:
				exploderecs[b].rec.X += rFloat32(2, 4)
			case 1:
				exploderecs[b].rec.X -= rFloat32(1, 3)
				exploderecs[b].rec.Y += rFloat32(1, 3)
			case 2:
				exploderecs[b].rec.Y += rFloat32(1, 3)
			case 3:
				exploderecs[b].rec.X += rFloat32(1, 3)
				exploderecs[b].rec.Y += rFloat32(1, 3)
			}

			if exploderecs[b].rec.Width > 0 {
				exploderecs[b].rec.Width -= 0.1
			}
			if exploderecs[b].rec.Height > 0 {
				exploderecs[b].rec.Height -= 0.1
			}
			if exploderecs[b].opac > 0 {
				exploderecs[b].opac -= 0.01
			}

			explodetimer--
		}
	}

	// layer 2
	count = 0
	x, y = 0, 0
	drawblok = 0
	for a := 0; a < drawa; a++ {

		x += 16
		count++
		drawblok++
		if count == draww {
			count = 0
			x = 0
			y += 16

		}

	}

	// player layer

	for a := 0; a < len(snyk); a++ {
		if snyk[a].activ {
			snyk[a].x = (snyk[a].blocknumber % draww) * 16
			snyk[a].y = (snyk[a].blocknumber / draww) * 16
			if outlineblocks && !circleblocks {
				rl.DrawRectangleLines(snyk[a].x, snyk[a].y, blockw, blockw, snyk[a].color)
				if ghosting {
					rl.DrawRectangleLines(snyk[a].x+rInt(-3, 4), snyk[a].y+rInt(-3, 4), blockw, blockw, rl.Fade(snyk[a].color, rF32(0.3, 1.1)))
				}
			} else if outlineblocks && circleblocks {
				rl.DrawCircleLines(snyk[a].x+8, snyk[a].y+8, float32(blockw/2), snyk[a].color)
				if ghosting {
					rl.DrawCircleLines((snyk[a].x+8)+rInt(-3, 4), (snyk[a].y+8)+rInt(-3, 4), float32(blockw/2), rl.Fade(snyk[a].color, rF32(0.3, 1.1)))
				}
			} else if !outlineblocks && circleblocks {
				rl.DrawCircle(snyk[a].x+8, snyk[a].y+8, float32(blockw/2), snyk[a].color)
				if ghosting {
					rl.DrawCircle((snyk[a].x+8)+rInt(-3, 4), (snyk[a].y+8)+rInt(-3, 4), float32(blockw/2), rl.Fade(snyk[a].color, rF32(0.3, 1.1)))
				}
			} else {
				rl.DrawRectangle(snyk[a].x, snyk[a].y, blockw, blockw, snyk[a].color)
				if ghosting {
					rl.DrawRectangle(snyk[a].x+rInt(-3, 4), snyk[a].y+rInt(-3, 4), blockw, blockw, rl.Fade(snyk[a].color, rF32(0.3, 1.1)))
				}
			}
		}
	}

}
func drawnocamera() { // MARK: drawnocamera

	if introon {
		drawintro()
	}
	if options {

		optx := monw/2 - 250
		opty := 200
		rl.DrawRectangle(optx, 0, 500, monh, rl.Fade(rl.White, 0.8))
		if ghosting {
			rl.DrawRectangle(optx+rInt(-5, 6), 0+rInt(-5, 6), 500, monh, rl.Fade(rl.White, 0.3))
		}
		for a := 0; a < len(optionsonoff); a++ {
			if optionselect == a {
				rl.DrawRectangle(optx, opty-5, 500, 50, rl.Fade(randomcolor(), fadeblink))
			}
			rl.DrawRectangle(monw/2+150, opty, 40, 40, rl.Black)
			if optionsonoff[a] {
				rl.DrawRectangle(monw/2+155, opty+5, 30, 30, randomcolor())
			}

			opty += 60
		}

		opty = 100

		length := rl.MeasureText("options", 80)
		rl.DrawText("options", (monw/2)-(length/2)+rInt(-5, 6), opty+rInt(-5, 6), 80, randomcolor())
		rl.DrawText("options", (monw/2)-(length/2), opty, 80, rl.Black)
		opty += 100
		rl.DrawText("ghosting", optx+80, opty, 40, rl.Black)
		opty += 60
		rl.DrawText("pixel noise", optx+80, opty, 40, rl.Black)
		opty += 60
		rl.DrawText("scan lines", optx+80, opty, 40, rl.Black)
		opty += 60
		rl.DrawText("background", optx+80, opty, 40, rl.Black)
		opty += 60

	}
	if startgame {
		drawnewgamemenu()
	}
	if helpon {
		drawhelpscreen()
	}
	if endscreen {
		endgame()
	}
	if pixel_noise {
		for a := 0; a < 1000; a++ {
			w := rFloat32(1, 4)
			rec := rl.NewRectangle(rFloat32(0, monw), rFloat32(0, monh), w, w)
			rl.DrawRectangleRec(rec, rl.Black)
		}
	}
	if scanlines {
		y := 0
		for {
			rl.DrawLine(0, y, monw, y, rl.Fade(rl.Black, 0.3))
			y += 3
			if y >= monh {
				break
			}
		}
	}

}
func drawintro() { // MARK: drawintro
	rl.DrawRectangle(0, 0, monw, monh, rl.White)

	rl.DrawText("©", 10, 10, 50, rl.Fade(brightred(), fadeblink))
	rl.DrawText("2021 nicholasimon", 56, 24, 30, rl.Fade(brightred(), fadeblink))

	rl.DrawText("snyk", monw/2-200, introy, 200, rl.Black)
	rl.DrawText("snyk", monw/2-203, introy-3, 200, rl.White)
	rl.DrawText("snyk", monw/2-204, introy-4, 200, introcolor1)

	if introy < monh/2-100 {
		introy += 10
	}

	v2 := rl.NewVector2(float32(monw-380), float32(monh-180))
	rl.DrawTextureRec(imgs, gologo, v2, rl.Fade(rl.White, rF32(0.3, 1.1)))
	v2.X += 180
	v2.Y += 10
	rl.DrawTextureRec(imgs, rayliblogo, v2, rl.Fade(rl.White, rF32(0.3, 1.1)))
}
func collect(block int) { // MARK: collect

	roomlayout[block] = blok{}
	snyk[snykcount].activ = true
	snyk[snykcount].blocknumber = snyk[snykcount-1].previousblock
	snykcount++
	for {
		choose := rInt(0, drawa)
		if !roomlayout[choose].snyk {
			if !roomlayout[choose].solid {
				roomlayout[choose].activ = true
				roomlayout[choose].color = brightyellow()
				roomlayout[choose].collect = true
				collectblok = choose
				break
			}
		}
	}

}
func updateplayer() { // MARK: updateplayer

	// snyk room positions
	for a := 0; a < len(snyk); a++ {
		if snyk[a].activ {
			roomlayout[snyk[a].blocknumber].snyk = true
		} else if !snyk[a].activ {
			roomlayout[snyk[a].blocknumber].snyk = false
		}

	}

	// move
	switch snyk[0].direction {
	case 1:
		if !roomlayout[snyk[0].blocknumber-draww].solid {
			snyk[0].previousblock = snyk[0].blocknumber
			snyk[0].blocknumber -= draww
		} else {
			if !hppause {
				if !invincible {
					hp--

					hppausetimer = 30
					hppause = true
				}
			}

			if snyk[0].bounce {
				bounceplayer(0)
			}
		}
	case 2:
		if !roomlayout[snyk[0].blocknumber+1].solid {
			snyk[0].previousblock = snyk[0].blocknumber
			snyk[0].blocknumber++
		} else {
			if !hppause {
				if !invincible {
					hp--
					hppausetimer = 30
					hppause = true
				}
			}
			if snyk[0].bounce {
				bounceplayer(0)
			}
		}
	case 3:
		if !roomlayout[snyk[0].blocknumber+draww].solid {
			snyk[0].previousblock = snyk[0].blocknumber
			snyk[0].blocknumber += draww
		} else {
			if !hppause {
				if !invincible {
					hp--
					hppausetimer = 30
					hppause = true
				}
			}
			if snyk[0].bounce {
				bounceplayer(0)
			}
		}
	case 4:
		if !roomlayout[snyk[0].blocknumber-1].solid {
			snyk[0].previousblock = snyk[0].blocknumber
			snyk[0].blocknumber--
		} else {
			if !hppause {
				if !invincible {
					hp--
					hppausetimer = 30
					hppause = true
				}
			}
			if snyk[0].bounce {
				bounceplayer(0)
			}
		}
	}
	for a := 1; a < len(snyk); a++ {
		if snyk[a].activ {
			snyk[a].previousblock = snyk[a].blocknumber
			snyk[a].blocknumber = snyk[a-1].previousblock
		}
	}

}
func bounceplayer(snykblok int) { // MARK: drawnocamera

	switch snyk[snykblok].direction {
	case 1:
		snyk[snykblok].direction = rInt(2, 5)
	case 2:
		for {
			snyk[snykblok].direction = rInt(1, 5)
			if snyk[snykblok].direction != 2 {
				break
			}
		}
	case 3:
		for {
			snyk[snykblok].direction = rInt(1, 5)
			if snyk[snykblok].direction != 3 {
				break
			}
		}
	case 4:
		snyk[snykblok].direction = rInt(1, 4)
	}

}
func explode(x, y int) { // MARK: explode

	x += 8
	y += 8
	for a := 0; a < len(exploderecs); a++ {
		wid := float32(10)
		exploderecs[a].rec = rl.NewRectangle(float32(x+rInt(-2, 3)), float32(y+rInt(-2, 3)), wid, wid)
		exploderecs[a].color = brightyellow()
		exploderecs[a].opac = 0.4
		if a < 4 {
			exploderecs[a].direction = rInt(1, 5)
		} else {
			exploderecs[a].direction = rInt(5, 9)
		}
	}

	explodetimer = 60

}
func createspecial() { // MARK: createspecial

	roomlayout[crateblock] = blok{}
	crateactiv = false
	cratetimernext = rInt(30, 90)

	choose := rInt(0, 6)

	switch choose {
	case 5:
		if !invincible {
			hp++
		}
	case 4:
		if circleblocks {
			circleblocks = false
		} else {
			circleblocks = true
		}
	case 3:
		blockw += 2
	case 2:
		if outlineblocks {
			outlineblocks = false
		} else {
			outlineblocks = true
		}
	case 1:
		camera.Rotation += rFloat32(-5.0, 6.0)
	case 0:
		side := rInt(1, 5)
		switch side {
		case 4:
			startblok := centerblok
			startblok -= (roomw / 2)
			startblok += rInt(-(roomh/2), (roomh/3)) * draww
			width := rInt(8, 20)
			count := 0
			for {
				if roomlayout[startblok].snyk || roomlayout[startblok].collect || roomlayout[startblok].crate {
					break
				} else {
					roomlayout[startblok].activ = true
					roomlayout[startblok].solid = true
					startblok += draww
					count++
					if count == width {
						count = 0
						startblok -= width * draww
						startblok++
						startblok += rInt(1, 3) * draww
						width -= rInt(2, 4)
					}
					if width < 1 {
						break
					}
				}
			}
		case 2:
			startblok := centerblok
			startblok += (roomw / 2)
			startblok += rInt(-(roomh/2), (roomh/3)) * draww
			width := rInt(8, 20)
			count := 0
			for {
				if roomlayout[startblok].snyk || roomlayout[startblok].collect || roomlayout[startblok].crate {
					break
				} else {
					roomlayout[startblok].activ = true
					roomlayout[startblok].solid = true
					startblok += draww
					count++
					if count == width {
						count = 0
						startblok -= width * draww
						startblok--
						startblok += rInt(1, 3) * draww
						width -= rInt(2, 4)
					}
					if width < 1 {
						break
					}
				}
			}
		case 3:
			startblok := centerblok
			startblok += (roomh / 2) * draww
			startblok += rInt(-(roomw / 2), (roomw / 3))
			width := rInt(10, 25)
			count := 0
			for {
				if roomlayout[startblok].snyk || roomlayout[startblok].collect || roomlayout[startblok].crate {
					break
				} else {
					roomlayout[startblok].activ = true
					roomlayout[startblok].solid = true
					startblok++
					count++
					if count == width {
						count = 0
						startblok -= width
						startblok -= draww
						startblok += rInt(1, 3)
						width -= rInt(2, 4)
					}
					if width < 1 {
						break
					}
				}
			}
		case 1:
			startblok := centerblok
			startblok -= (roomh / 2) * draww
			startblok += rInt(-(roomw / 2), (roomw / 3))
			width := rInt(10, 25)
			count := 0
			for {
				if roomlayout[startblok].snyk || roomlayout[startblok].collect || roomlayout[startblok].crate {
					break
				} else {
					roomlayout[startblok].activ = true
					roomlayout[startblok].solid = true
					startblok++
					count++
					if count == width {
						count = 0
						startblok -= width
						startblok += draww
						startblok += rInt(1, 3)
						width -= rInt(2, 4)
					}
					if width < 1 {
						break
					}
				}
			}
		}

	}

}
func createcrate() { // MARK: createcrate

	if !crateactiv {
		for {
			choose := rInt(0, drawa)
			if !roomlayout[choose].snyk {
				if !roomlayout[choose].solid {
					roomlayout[choose].activ = true
					roomlayout[choose].crate = true
					roomlayout[choose].color = brightyellow()
					crateblock = choose
					break
				}
			}
		}
		crateactiv = true
		cratetimer = 10
		cratetimercount = 0
		if flipcoin() {
			autosnykcrate = true
		}
	}
	crateon = false
}
func createbackgrounds() { // MARK: createbackgrounds

	backorig := backg{}

	backorig.moving = true
	backorig.poly = true
	backorig.rotating = flipcoin()
	backorig.resizing = true
	backorig.lines = flipcoin()

	for a := 0; a < len(currentback); a++ {

		currentback[a].moving = backorig.moving
		currentback[a].poly = backorig.poly
		currentback[a].rotating = backorig.rotating
		currentback[a].resizing = backorig.resizing
		currentback[a].lines = backorig.lines

		currentback[a].color = randomcolor()
		currentback[a].h = rInt(1, 21)
		currentback[a].w = rInt(1, 21)
		if currentback[a].poly {
			if currentback[a].w < 6 {
				currentback[a].w += rInt(5, 10)
			}
		}
		currentback[a].opac = 0.1
		currentback[a].sides = rInt(3, 9)
		currentback[a].rotation = rFloat32(0, 360)
		currentback[a].x = rInt(0, monw)
		currentback[a].y = rInt(0, monh)

		for {
			currentback[a].direction = rInt(1, 9)
			if currentback[a].direction != 5 {
				break
			}
		}
	}

}
func drawnewgamemenu() { // MARK: drawnewgamemenu
	paused = true

	optx := monw/2 - 250
	opty := 200
	rl.DrawRectangle(optx, 0, 500, monh, rl.Fade(rl.White, 0.8))
	if ghosting {
		rl.DrawRectangle(optx+rInt(-5, 6), 0+rInt(-5, 6), 500, monh, rl.Fade(rl.White, 0.3))
	}
	for a := 0; a < len(newgameoptions); a++ {
		if a < 4 {
			if startselect == a {
				rl.DrawRectangle(optx, opty-5, 500, 50, rl.Fade(randomcolor(), fadeblink))
			}

			rl.DrawRectangle(monw/2+150, opty, 40, 40, rl.Black)
			if newgameoptions[a] {
				rl.DrawRectangle(monw/2+155, opty+5, 30, 30, randomcolor())
			}
			opty += 60
		} else {
			if startselect == a {
				rl.DrawRectangle(optx, opty-5, 500, 70, rl.Fade(randomcolor(), fadeblink))
			}
			opty += 60
		}

	}

	opty = 100

	length := rl.MeasureText("new game", 80)
	rl.DrawText("new game", (monw/2)-(length/2)+rInt(-5, 6), opty+rInt(-5, 6), 80, randomcolor())
	rl.DrawText("new game", (monw/2)-(length/2), opty, 80, rl.Black)
	opty += 100
	rl.DrawText("invincible", optx+80, opty, 40, rl.Black)
	opty += 60
	rl.DrawText("grayscale", optx+80, opty, 40, rl.Black)
	opty += 60
	rl.DrawText("supa fast", optx+80, opty, 40, rl.Black)
	opty += 60
	rl.DrawText("auto demo", optx+80, opty, 40, rl.Black)
	opty += 60
	rl.DrawText("start", optx+160+rInt(-5, 6), opty+rInt(-5, 6), 60, randomcolor())
	rl.DrawText("start", optx+160, opty, 60, rl.Black)
	opty += 60
	rl.DrawText("help", optx+175+rInt(-5, 6), opty+rInt(-5, 6), 60, randomcolor())
	rl.DrawText("help", optx+175, opty, 60, rl.Black)
}
func drawhelpscreen() { // MARK: drawhelpscreen

	rl.DrawRectangle(0, 0, monw, monh, rl.White)
	if rolldice() == 6 {
		rl.DrawRectangle(0, 0, monw, monh, randomcolor())
	}

	length := rl.MeasureText("help", 120)
	rl.DrawText("help", (monw/2)-(length/2)+rInt(-5, 6), 20+rInt(-5, 6), 120, randomcolor())
	rl.DrawText("help", (monw/2)-(length/2), 20, 120, rl.Black)

	length = rl.MeasureText("arrow keys move", 60)
	rl.DrawText("arrow keys move", (monw/2)-(length/2), 180, 60, rl.Black)
	length = rl.MeasureText("space key interact menus", 60)
	rl.DrawText("space key interact menus", (monw/2)-(length/2), 240, 60, rl.Black)
	length = rl.MeasureText("escape key options", 60)
	rl.DrawText("escape key options", (monw/2)-(length/2), 300, 60, rl.Black)
	length = rl.MeasureText("don't hit the walls", 60)
	rl.DrawText("don't hit the walls", (monw/2)-(length/2), 360, 60, rl.Black)
	length = rl.MeasureText("collect coins", 60)
	rl.DrawText("collect coins", (monw/2)-(length/2), 420, 60, rl.Black)
	length = rl.MeasureText("powerups change gameplay", 60)
	rl.DrawText("powerups change gameplay", (monw/2)-(length/2), 480, 60, rl.Black)
	length = rl.MeasureText("and increase score", 60)
	rl.DrawText("and increase score", (monw/2)-(length/2), 540, 60, rl.Black)
	length = rl.MeasureText("end key exits", 60)
	rl.DrawText("end key exits", (monw/2)-(length/2), 600, 60, rl.Black)

	if rl.IsKeyPressed(rl.KeySpace) || rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyEnter) {
		spacepause = 30
		helpon = false
	}

}
func endgame() { // MARK: endgame
	paused = true
	endfade += 0.05

	rl.BeginMode2D(cameraend)
	rl.DrawRectangle(0, 0, monw, monh, rl.Fade(rl.Black, endfade))
	if endnumber < len(endshapes)-50 {
		endnumber += 50
	}

	for a := 0; a < endnumber; a++ {

		rl.DrawRectangle(int(endshapes[a].v2.X), int(endshapes[a].v2.Y), endshapes[a].w, endshapes[a].w, endshapes[a].color)
	}
	rl.EndMode2D()
	if cameraend.Zoom < 4.0 {
		cameraend.Zoom += 0.05
	} else {
		rl.DrawText("the end", monw/2-204, monh/2-96, 120, rl.Black)
		rl.DrawText("the end", monw/2-200, monh/2-100, 120, rl.White)
		rl.DrawText("score:", monw/2-204, monh/2+4, 80, rl.Black)
		rl.DrawText("score:", monw/2-200, monh/2, 80, rl.White)
		endscore = (snykcount * 10) * multiplier
		endscoretxt := strconv.Itoa(endscore)
		rl.DrawText(endscoretxt, monw/2-204, monh/2+84, 120, rl.Black)
		rl.DrawText(endscoretxt, monw/2-200, monh/2+80, 120, rl.White)
	}
	if rl.IsKeyPressed(rl.KeySpace) || rl.IsKeyPressed(rl.KeyEnter) {
		endscreen = false
		endfade = float32(0.0)
		endnumber = 1
		spacepause = 30
		startgame = true
	}

}
func newlevel() { // MARK: newlevel

	for a := 0; a < len(roomlayout); a++ {
		roomlayout[a] = blok{}
	}

	multiplier = 1
	snykcount = 0
	hp = 3
	crateactiv = false
	cratetimernext = rInt(30, 90)
	createbackgrounds()
	for a := 0; a < len(roomlayout); a++ {
		roomlayout[a].color = randomcolor()
		roomlayout[a].color2 = randomgrey()
		roomlayout[a].activ = true
		roomlayout[a].solid = true
	}
	// clear unnecessary bloks
	area := 15 * drawh
	count := 0
	clearblok := 0
	for a := 0; a < area; a++ {
		roomlayout[clearblok].activ = false
		clearblok++
		count++
		if count == 15 {
			count = 0
			clearblok -= 15
			clearblok += draww
		}

	}
	clearblok = draww - 15
	for a := 0; a < area; a++ {
		roomlayout[clearblok].activ = false
		clearblok++
		count++
		if count == 15 {
			count = 0
			clearblok -= 15
			clearblok += draww
		}
	}
	area = draww * 5
	clearblok = 0
	for a := 0; a < area; a++ {
		roomlayout[clearblok].activ = false
		clearblok++
	}
	clearblok = draww * (drawh - 5)
	for a := 0; a < area; a++ {
		roomlayout[clearblok].activ = false
		clearblok++
	}

	roomblok := centerblok
	roomblok -= roomw / 2
	roomblok -= (roomh / 2) * draww
	count = 0
	for a := 0; a < rooma; a++ {
		roomlayout[roomblok].activ = false
		roomlayout[roomblok].solid = false
		roomblok++
		count++
		if count == roomw {
			count = 0
			roomblok -= roomw
			roomblok += draww
		}
	}

	snyk = make([]playerblok, drawa)
	snyk[0].activ = true
	snyk[0].bounce = true
	snyk[0].blocknumber = centerblok + rInt(-5, 6)
	snyk[0].blocknumber += rInt(-5, 6) * draww
	snyk[0].color = brightred()
	snyk[0].direction = rInt(1, 5)
	snykcount++

	for a := 1; a < len(snyk); a++ {
		snyk[a].color = brightred()
	}
	for {
		choose := rInt(0, drawa)
		if !roomlayout[choose].snyk {
			if !roomlayout[choose].solid {
				roomlayout[choose].activ = true
				roomlayout[choose].color = brightyellow()
				roomlayout[choose].collect = true
				collectblok = choose
				break
			}
		}
	}
	createcrate()

}

// MARK: core	core	core	core	core	core	core	core	core	core	core
func main() { // MARK: main
	rand.Seed(time.Now().UnixNano()) // random numbers
	rl.SetTraceLogLevel(rl.LogError) // hides info window
	rl.InitWindow(monw, monh, "setres")
	setres(0, 0)
	rl.CloseWindow()
	setinitialvalues()
	raylib()

}
func input() { // MARK: input

	if rl.IsKeyPressed(rl.KeyPause) {
		if paused {
			paused = false
		} else {
			paused = true
		}
	}

	// DEV KEYS DELETE
	if rl.IsKeyPressed(rl.KeyF5) {
		if endscreen {
			endscreen = false
		} else {
			endscreen = true
		}
	}
	if rl.IsKeyPressed(rl.KeyF4) {
		if startgame {
			startgame = false
		} else {
			startgame = true
		}
	}
	if rl.IsKeyPressed(rl.KeyF3) {
		if autosnyk {
			autosnyk = false
			newgameoptions[3] = true
		} else {
			autosnyk = true
			newgameoptions[3] = false
		}

	}
	if rl.IsKeyPressed(rl.KeyF2) {
		if crateon {
			crateon = false
		} else {
			crateon = true
		}

	}
	if rl.IsKeyPressed(rl.KeyF1) {
		if introon {
			introon = false
		} else {
			introon = true
		}

	}

	// DEV KEYS DELETE
	if rl.IsKeyPressed(rl.KeyEscape) {
		if !introon && !helpon && !startgame {
			if options {
				options = false
			} else {
				options = true
			}
		}
	}
	if introon {
		options = false
		startgame = false
		if rl.IsKeyPressed(rl.KeySpace) || rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyEscape) {
			introon = false
			startgame = true
		}
	} else if options {
		if rl.IsKeyPressed(rl.KeyUp) {
			optionselect--
			if optionselect < 0 {
				optionselect = len(optionsonoff) - 1
			}
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			optionselect++
			if optionselect > len(optionsonoff)-1 {
				optionselect = 0
			}
		}

		if rl.IsKeyPressed(rl.KeySpace) {

			switch optionselect {
			case 0:
				if ghosting {
					ghosting = false
					optionsonoff[0] = false
				} else {
					ghosting = true
					optionsonoff[0] = true
				}
			case 1:
				if pixel_noise {
					pixel_noise = false
					optionsonoff[1] = false
				} else {
					pixel_noise = true
					optionsonoff[1] = true
				}
			case 2:
				if scanlines {
					scanlines = false
					optionsonoff[2] = false
				} else {
					scanlines = true
					optionsonoff[2] = true
				}
			case 3:
				if backgon {
					backgon = false
					optionsonoff[3] = false
				} else {
					backgon = true
					optionsonoff[3] = true
				}
			}

		}

	} else if startgame {
		if rl.IsKeyPressed(rl.KeyUp) {
			startselect--
			if startselect < 0 {
				startselect = len(newgameoptions) - 1
			}
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			startselect++
			if startselect > len(newgameoptions)-1 {
				startselect = 0
			}
		}
		if spacepause == 0 {
			if rl.IsKeyPressed(rl.KeySpace) {

				switch startselect {
				case 0:
					if invincible {
						invincible = false
						newgameoptions[0] = false
					} else {
						invincible = true
						newgameoptions[0] = true
					}
				case 1:
					if grayscale {
						grayscale = false
						newgameoptions[1] = false
					} else {
						grayscale = true
						newgameoptions[1] = true
					}
				case 2:
					if supafast {
						supafast = false
						newgameoptions[2] = false
					} else {
						supafast = true
						newgameoptions[2] = true
					}
				case 3:
					if autosnyk {
						autosnyk = false
						newgameoptions[3] = false
					} else {
						autosnyk = true
						newgameoptions[3] = true
					}
				case 4:
					newlevel()
					startgame = false
					paused = false
					startselect = 0
				case 5:
					helpon = true

				}
			}
		}

	} else {
		if rl.IsKeyPressed(rl.KeyLeft) {
			if snyk[0].direction != 2 {
				snyk[0].direction = 4
			}
		}
		if rl.IsKeyPressed(rl.KeyRight) {
			if snyk[0].direction != 4 {
				snyk[0].direction = 2
			}
		}
		if rl.IsKeyPressed(rl.KeyUp) {
			if snyk[0].direction != 3 {
				snyk[0].direction = 1
			}
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			if snyk[0].direction != 1 {
				snyk[0].direction = 3
			}
		}
	}
	if rl.IsKeyPressed(rl.KeyKpAdd) {
		camera.Zoom += 0.2
	}
	if rl.IsKeyPressed(rl.KeyKpSubtract) {
		camera.Zoom -= 0.2
	}

	if rl.IsKeyPressed(rl.KeyKpDecimal) {
		if debugon {
			debugon = false
		} else {
			debugon = true
		}
	}

	if rl.IsKeyPressed(rl.KeyKp0) {
		if gridon {
			gridon = false
		} else {
			gridon = true
		}
	}
}
func drawdebug() { // MARK: DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG

	//centerlines
	rl.DrawLine(monw/2, 0, monw/2, monh, rl.Magenta)
	rl.DrawLine(0, monh/2, monw, monh/2, rl.Magenta)

	rl.DrawRectangle(monw-300, 0, 300, monh, rl.Fade(rl.Black, 0.5))
	textx := monw - 290
	textx2 := monw - 145
	texty := 10

	drawatext := strconv.Itoa(drawa)
	drawwtext := strconv.Itoa(draww)
	drawhtext := strconv.Itoa(drawh)
	roomatext := strconv.Itoa(rooma)
	snyklentext := strconv.Itoa(len(snyk))
	snyk0previoustext := strconv.Itoa(snyk[0].previousblock)
	snyk0currentbloktext := strconv.Itoa(snyk[0].blocknumber)
	camerazoomtext := fmt.Sprintf("%g", camera.Zoom)

	rl.DrawText("drawa", textx, texty, 10, rl.White)
	rl.DrawText(drawatext, textx2, texty, 10, rl.White)
	texty += 12
	rl.DrawText("draww", textx, texty, 10, rl.White)
	rl.DrawText(drawwtext, textx2, texty, 10, rl.White)
	texty += 12
	rl.DrawText("drawh", textx, texty, 10, rl.White)
	rl.DrawText(drawhtext, textx2, texty, 10, rl.White)
	texty += 12
	rl.DrawText("rooma", textx, texty, 10, rl.White)
	rl.DrawText(roomatext, textx2, texty, 10, rl.White)
	texty += 12
	rl.DrawText("snyklen", textx, texty, 10, rl.White)
	rl.DrawText(snyklentext, textx2, texty, 10, rl.White)
	texty += 12
	rl.DrawText("snyk0previoustext", textx, texty, 10, rl.White)
	rl.DrawText(snyk0previoustext, textx2, texty, 10, rl.White)
	texty += 12
	rl.DrawText("snyk0currentbloktext", textx, texty, 10, rl.White)
	rl.DrawText(snyk0currentbloktext, textx2, texty, 10, rl.White)
	texty += 12
	rl.DrawText("camerazoomtext", textx, texty, 10, rl.White)
	rl.DrawText(camerazoomtext, textx2, texty, 10, rl.White)
	texty += 12

	// fps
	rl.DrawRectangle(monw-110, monh-110, 100, 40, rl.Black)
	rl.DrawFPS(monw-100, monh-100)

}
func timers() { // MARK: timers

	if spacepause > 0 {
		spacepause--
	}

	gametimecount++
	if gametimecount%30 == 0 {
		gametime--
	}

	if hppausetimer != 0 {
		hppausetimer--
		if hppausetimer == 0 {
			hppause = false
		}
	}

	if onoff3 {
		coin.X += 16
		if coin.X > 70 {
			coin.X = 0
		}
	}

	if autosnyktimer > 0 {
		autosnyktimer--
		if autosnyktimer == 0 {
			autosnykpause = false
		}
	}

	if cratetimernext != 0 {
		cratetimernext--
		if cratetimernext == 0 {
			crateon = true
		}
	}

	if cratetimer != 0 {
		cratetimercount++
		if cratetimercount%30 == 0 {
			cratetimer--
		}
		if cratetimer == 0 {
			roomlayout[crateblock] = blok{}
			crateactiv = false
			cratetimernext = rInt(30, 90)
			cratetimercount = 0
		}
	}

	if introon {
		if onoff30 {
			introcolor1 = randomcolor()
		}
		if onoff15 {
			introy -= rInt(0, 16)
		}
	}

	if framecount%2 == 0 {
		if onoff2 {
			onoff2 = false
		} else {
			onoff2 = true
		}
	}
	if framecount%3 == 0 {
		if onoff3 {
			onoff3 = false
		} else {
			onoff3 = true
		}
	}
	if framecount%6 == 0 {
		if onoff6 {
			onoff6 = false
		} else {
			onoff6 = true
		}
	}
	if framecount%10 == 0 {
		if onoff10 {
			onoff10 = false
		} else {
			onoff10 = true
		}
	}
	if framecount%15 == 0 {
		if onoff15 {
			onoff15 = false
		} else {
			onoff15 = true
		}
	}
	if framecount%30 == 0 {
		if onoff30 {
			onoff30 = false
		} else {
			onoff30 = true
		}
	}
	if framecount%60 == 0 {
		if onoff60 {
			onoff60 = false
		} else {
			onoff60 = true
		}
	}
	if fadeblinkon {
		if fadeblink > 0.2 {
			fadeblink -= 0.05
		} else {
			fadeblinkon = false
		}
	} else {
		if fadeblink < 0.6 {
			fadeblink += 0.05
		} else {
			fadeblinkon = true
		}
	}
}

func setres(w, h int) { // MARK: setres

	monw = rl.GetMonitorWidth(0)
	monh = rl.GetMonitorHeight(0)

	if monw > 1500 && monw < 1800 {
		camera.Offset.X = float32(50 - (monw / 4))
		camera.Offset.Y = float32(0 - (monh / 5))
		camera.Zoom = 1.2
	} else if monw <= 1500 && monw > 1400 {
		camera.Offset.X = float32(0 - (monw / 6))
		camera.Offset.Y = float32(0 - (monh / 9))
		camera.Zoom = 1.0
	} else if monw <= 1300 {
		camera.Offset.X = float32(0 - (monw / 6))
		camera.Offset.Y = float32(0 - (monh / 6))
		camera.Zoom = 0.9
	} else if monw > 1300 && monw <= 1400 {
		camera.Offset.X = -280
		camera.Offset.Y = float32(0 - (monh / 5))
		camera.Zoom = 1.0
	} else if monw >= 1800 && monw < 2000 {
		camera.Offset.X = -500
		camera.Offset.Y = -250
		camera.Zoom = 1.5
	} else if monw >= 2000 && monw < 3000 {
		camera.Offset.X = -250
		camera.Offset.Y = -100
		camera.Zoom = 1.6
	} else if monw > 3000 {
		camera.Offset.X = float32(0 - (monw / 5))
		camera.Offset.Y = float32(10 - (monh / 6))
		camera.Zoom = 2.8
	}
}

func setinitialvalues() { // MARK: setinitialvalues

	invincible = false
	newgameoptions[0] = false
	grayscale = false
	newgameoptions[1] = false
	supafast = false
	newgameoptions[2] = false
	crateactiv = false

	createbackgrounds()
	introcolor1 = randomcolor()
	backgon = true
	optionsonoff[3] = true
	pixel_noise = true
	optionsonoff[1] = true
	scanlines = true
	optionsonoff[2] = true
	ghosting = true
	optionsonoff[0] = true
	//	gridon = true

	draww = (1920 / 16) + 1
	drawh = (1080 / 16) + 1
	drawa = draww * drawh

	centerblok = draww / 2
	centerblok += ((drawh / 2) - 1) * draww

	roomlayout = make([]blok, drawa)
	endshapes = make([]endshape, drawa)
	count := 0
	x := float32(0)
	y := float32(0)
	choosecolor := randomcolor()
	switch rolldice() {
	case 1:
		choosecolor = randomcolor()
	case 2:
		choosecolor = randombluedark()
	case 3:
		choosecolor = randomorange()
	case 4:
		choosecolor = randomyellow()
	case 5:
		choosecolor = randomgreen()
	case 6:
		choosecolor = randomred()
	}
	for a := 0; a < len(endshapes); a++ {
		endshapes[a].v2.X = x
		endshapes[a].v2.Y = y
		endshapes[a].color = choosecolor
		endshapes[a].w = rInt(12, 21)
		x += 16
		count++
		if count == draww {
			count = 0
			x = 0
			y += 16
		}
	}

	for a := 0; a < len(roomlayout); a++ {
		roomlayout[a].color = randomcolor()
		roomlayout[a].color2 = randomgrey()
		roomlayout[a].activ = true
		roomlayout[a].solid = true
	}
	// clear unnecessary bloks
	area := 15 * drawh
	count = 0
	clearblok := 0
	for a := 0; a < area; a++ {
		roomlayout[clearblok].activ = false
		clearblok++
		count++
		if count == 15 {
			count = 0
			clearblok -= 15
			clearblok += draww
		}

	}
	clearblok = draww - 15
	for a := 0; a < area; a++ {
		roomlayout[clearblok].activ = false
		clearblok++
		count++
		if count == 15 {
			count = 0
			clearblok -= 15
			clearblok += draww
		}
	}
	area = draww * 5
	clearblok = 0
	for a := 0; a < area; a++ {
		roomlayout[clearblok].activ = false
		clearblok++
	}
	clearblok = draww * (drawh - 5)
	for a := 0; a < area; a++ {
		roomlayout[clearblok].activ = false
		clearblok++
	}

	roomblok := centerblok
	roomblok -= roomw / 2
	roomblok -= (roomh / 2) * draww
	count = 0
	for a := 0; a < rooma; a++ {
		roomlayout[roomblok].activ = false
		roomlayout[roomblok].solid = false
		roomblok++
		count++
		if count == roomw {
			count = 0
			roomblok -= roomw
			roomblok += draww
		}
	}

	snyk = make([]playerblok, drawa)

}

// MARK:  █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █
func drawgrid() { // MARK: drawgrid

	x := 16
	for {
		rl.DrawLine(x, 0, x, monh, rl.Fade(rl.Magenta, 0.1))
		x += 16
		if x > monw {
			break
		}
	}
	y := 16
	for {
		rl.DrawLine(0, y, monw, y, rl.Fade(rl.Magenta, 0.1))
		y += 16
		if y > monh {
			break
		}
	}

}

// MARK: colors
// https://www.rapidtables.com/web/color/RGB_Color.html
func darkred() rl.Color {
	color := rl.NewColor(55, 0, 0, 255)
	return color
}
func semidarkred() rl.Color {
	color := rl.NewColor(70, 0, 0, 255)
	return color
}
func brightred() rl.Color {
	color := rl.NewColor(230, 0, 0, 255)
	return color
}
func randomgrey() rl.Color {
	color := rl.NewColor(uint8(rInt(160, 193)), uint8(rInt(160, 193)), uint8(rInt(160, 193)), uint8(rInt(0, 255)))
	return color
}
func randombluelight() rl.Color {
	color := rl.NewColor(uint8(rInt(0, 180)), uint8(rInt(120, 256)), uint8(rInt(120, 256)), 255)
	return color
}
func randombluedark() rl.Color {
	color := rl.NewColor(0, 0, uint8(rInt(120, 250)), 255)
	return color
}
func randomyellow() rl.Color {
	color := rl.NewColor(255, uint8(rInt(150, 256)), 0, 255)
	return color
}
func randomorange() rl.Color {
	color := rl.NewColor(uint8(rInt(250, 256)), uint8(rInt(60, 210)), 0, 255)
	return color
}
func randomred() rl.Color {
	color := rl.NewColor(uint8(rInt(128, 256)), uint8(rInt(0, 129)), uint8(rInt(0, 129)), 255)
	return color
}
func randomgreen() rl.Color {
	color := rl.NewColor(uint8(rInt(0, 170)), uint8(rInt(100, 256)), uint8(rInt(0, 50)), 255)
	return color
}
func randomcolor() rl.Color {
	color := rl.NewColor(uint8(rInt(0, 256)), uint8(rInt(0, 256)), uint8(rInt(0, 256)), 255)
	return color
}
func brightyellow() rl.Color {
	color := rl.NewColor(uint8(255), uint8(255), uint8(0), 255)
	return color
}
func brightbrown() rl.Color {
	color := rl.NewColor(uint8(218), uint8(165), uint8(32), 255)
	return color
}
func brightgrey() rl.Color {
	color := rl.NewColor(uint8(212), uint8(212), uint8(213), 255)
	return color
}

// random numbers
func rF32(min, max float32) float32 {
	return (rand.Float32() * (max - min)) + min
}
func rInt(min, max int) int {
	return rand.Intn(max-min) + min
}
func rInt32(min, max int) int32 {
	a := int32(rand.Intn(max-min) + min)
	return a
}
func rFloat32(min, max int) float32 {
	a := float32(rand.Intn(max-min) + min)
	return a
}
func flipcoin() bool {
	var b bool
	a := rInt(0, 10001)
	if a < 5000 {
		b = true
	}
	return b
}
func rolldice() int {
	a := rInt(1, 7)
	return a
}
