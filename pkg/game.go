package santa

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"image/color"
	"math"
	"time"
)

type Game struct {
	ScreenWidth  float64
	ScreenHeight float64
	TileSize     float64

	Instance
}

type Instance struct {
	IntroTime float64

	Map         *tiled.Map
	Renderer    *render.Renderer
	RenderedMap *ebiten.Image

	lastTime time.Time

	Speed float64

	Player *Entity
	Mobs   []*Entity

	Money int

	GameOverMessage string
	Escaped         bool
	RestartAt       time.Time
}

func (g *Game) Update() error {
	if g.lastTime.IsZero() {
		g.lastTime = time.Now()
	}
	now := time.Now()
	if !g.RestartAt.IsZero() && g.RestartAt.Before(now) {
		g.Instance = LoadNewInstance()
		return nil
	}
	dt := now.Sub(g.lastTime).Seconds() * g.Speed
	g.lastTime = now

	g.IntroTime = math.Max(0, g.IntroTime-dt)
	if g.IntroTime > 0 {
		return nil
	}

	for _, e := range append(g.Mobs, g.Player) {
		e.Update(dt)
	}

	if !g.RestartAt.IsZero() {
		g.Speed *= 0.9
		return nil
	}

	for _, e := range append(g.Mobs, g.Player) {
		if e.Moving == None {
			e.Controller.Control(e, g)
		}
	}

	// money
	pt := g.Map.Layers[propIndex].Tiles[XyToIndex(g.Player.X, g.Player.Y)]
	for tid, m := range moneyTids {
		if int(pt.ID) == tid {
			pt.ID = 0
			pt.Nil = true
			g.Money += m
			g.RenderedMap = RenderMap(g.Renderer)
		}
	}

	// tile effects
	wt := g.Map.Layers[walkableIndex].Tiles[XyToIndex(g.Player.X, g.Player.Y)]
	if int(wt.ID) == walkableFall {
		g.Lost("You fell to your death.")
		return nil
	}
	if int(wt.ID) == walkableEnd {
		g.Won()
		return nil
	}

	// mobs
	px, py := g.Player.CollisionSpace()
	for _, m := range g.Mobs {
		ex, ey := m.CollisionSpace()
		if ex == px && ey == py {
			if m.Sprite == santaSprite {
				g.Lost("Santa devoured you.")
			} else {
				g.Lost("An elf arrested you.")
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(g.ScreenWidth/2-g.TileSize/2, g.ScreenHeight/2-g.TileSize/2)
		px, py := g.Player.EffectivePosition()
		op.GeoM.Translate(-px*g.TileSize, -py*g.TileSize)
		screen.DrawImage(g.RenderedMap, op)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Snow Gold: %d / %d", g.Money, highScore))
	}
	for _, m := range g.Mobs {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(g.ScreenWidth/2-g.TileSize/2, g.ScreenHeight/2-g.TileSize/4.0-g.TileSize/2)
		px, py := g.Player.EffectivePosition()
		mx, my := m.EffectivePosition()
		dx, dy := mx-px, my-py
		op.GeoM.Translate(dx*g.TileSize, dy*g.TileSize)
		screen.DrawImage(m.Sprite, op)

	}
	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(g.ScreenWidth/2-g.TileSize/2, g.ScreenHeight/2-g.TileSize/4.0-g.TileSize/2)
		screen.DrawImage(playerSprite, op)
	}

	if g.IntroTime > 0 {
		s := math.Min(1, g.IntroTime)
		bgClr := color.RGBA{
			A: uint8(s * 255),
		}
		ebitenutil.DrawRect(screen, 0, 0, g.ScreenWidth, g.ScreenHeight, bgClr)

		x1 := uint8(math.Max(0, math.Min(1, introTime-g.IntroTime)) * s * 255)
		txtClr1 := color.RGBA{R: x1, G: x1, B: x1, A: x1}
		DrawCenteredText(screen, smallFont, "Steal Snow Gold.", g.ScreenWidth/2, g.ScreenHeight/2-20, txtClr1)

		x2 := uint8(math.Max(0, math.Min(1, (introTime-1)-(g.IntroTime))) * s * 255)
		txtClr2 := color.RGBA{R: x2, G: x2, B: x2, A: x2}
		DrawCenteredText(screen, smallFont, "Escape.", g.ScreenWidth/2, g.ScreenHeight/2, txtClr2)

		x3 := uint8(math.Max(0, math.Min(1, (introTime-2)-(g.IntroTime))) * s * 255)
		txtClr3 := color.RGBA{R: x3, G: x3, B: x3, A: x3}
		DrawCenteredText(screen, smallFont, "See the high score.", g.ScreenWidth/2, g.ScreenHeight/2+20, txtClr3)
	}

	if g.Speed < 1 {
		clr := color.RGBA{
			A: uint8((1.0 - g.Speed) * 255),
		}
		ebitenutil.DrawRect(screen, 0, 0, g.ScreenWidth, g.ScreenHeight, clr)
	}
	if g.GameOverMessage != "" {
		DrawCenteredText(screen, bigFont, "You Lost!", g.ScreenWidth/2, g.ScreenHeight/2-20, colornames.White)
		DrawCenteredText(screen, smallFont, g.GameOverMessage, g.ScreenWidth/2, g.ScreenHeight/2+10, colornames.White)
	}
	if g.Escaped {
		DrawCenteredText(screen, bigFont, "You Escaped!", g.ScreenWidth/2, g.ScreenHeight/2-50, colornames.White)
		DrawCenteredText(screen, smallFont, fmt.Sprintf("High Score: %d", highScore), g.ScreenWidth/2, g.ScreenHeight/2-20, colornames.Red)
		DrawCenteredText(screen, smallFont, fmt.Sprintf("Your Score: %d", g.Money), g.ScreenWidth/2, g.ScreenHeight/2, colornames.White)
		DrawCenteredText(screen, smallFont, "What a Loser!", g.ScreenWidth/2, g.ScreenHeight/2+20, colornames.White)
		DrawCenteredText(screen, smallFont, "A game by Laundering'R'Us", g.ScreenWidth/2, g.ScreenHeight/2+40, colornames.Forestgreen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(g.ScreenWidth), int(g.ScreenHeight)
}

func (g *Game) Lost(msg string) {
	g.GameOverMessage = msg
	g.RestartAt = time.Now().Add(time.Second * 5)
}

func (g *Game) Won() {
	g.Escaped = true
	g.RestartAt = time.Now().Add(time.Second * 15)
}

func DrawCenteredText(screen *ebiten.Image, font font.Face, s string, cx, cy float64, clr color.RGBA) {
	bounds := text.BoundString(font, s)
	x, y := int(cx)-bounds.Min.X-bounds.Dx()/2, int(cy)-bounds.Min.Y-bounds.Dy()/2
	text.Draw(screen, s, font, x, y, clr)
}
