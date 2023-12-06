package santa

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type Controller interface {
	Control(e *Entity, g *Game)
}

type PlayerController struct {
	TouchID ebiten.TouchID
}

var inputs = map[ebiten.Key]Direction{
	ebiten.KeyW: Up,
	ebiten.KeyS: Down,
	ebiten.KeyA: Left,
	ebiten.KeyD: Right,

	ebiten.KeyUp:    Up,
	ebiten.KeyDown:  Down,
	ebiten.KeyLeft:  Left,
	ebiten.KeyRight: Right,

	ebiten.KeyNumpad8: Up,
	ebiten.KeyNumpad2: Down,
	ebiten.KeyNumpad5: Down,
	ebiten.KeyNumpad4: Left,
	ebiten.KeyNumpad6: Right,
}

func (c *PlayerController) Control(p *Entity, g *Game) {
	if (ebiten.IsKeyPressed(ebiten.Key7) && ebiten.IsKeyPressed(ebiten.Key8) && ebiten.IsKeyPressed(ebiten.Key9)) ||
		(ebiten.IsKeyPressed(ebiten.KeyNumpad1) && ebiten.IsKeyPressed(ebiten.KeyNumpad3) && ebiten.IsKeyPressed(ebiten.KeyNumpad7) && ebiten.IsKeyPressed(ebiten.KeyNumpad9)) {
		enemySpeed = 1.25
		//highScore = 917
	}

	if p.Moving != None {
		return
	}

	// keyboard
	for k, d := range inputs {
		if ebiten.IsKeyPressed(k) && g.MaybeMove(p, d) {
			return
		}
	}

	// mouse and touch
	tx, ty := int(g.ScreenWidth/2), int(g.ScreenHeight/2)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		tx, ty = ebiten.CursorPosition()
	}
	for _, id := range ebiten.AppendTouchIDs(nil) {
		tx, ty = ebiten.TouchPosition(id)
	}
	mx := (float64(tx) - (g.ScreenWidth / 2.0)) / (g.ScreenWidth / 2.0)
	my := (float64(ty) - (g.ScreenHeight / 2.0)) / (g.ScreenHeight / 2.0)
	amx, amy := math.Abs(mx), math.Abs(my)
	if amx < 0.1 && amy < 0.1 {
		return
	}
	d := None
	if amx > amy {
		if mx < 0 {
			d = Left
		} else {
			d = Right
		}
	} else {
		if my < 0 {
			d = Up
		} else {
			d = Down
		}
	}
	g.MaybeMove(p, d)
}

type EnemyHugController struct {
	Last          Direction
	AntiClockwise bool
}

func (c *EnemyHugController) Control(e *Entity, g *Game) {
	d := c.Last.Clockwise()
	if !c.AntiClockwise {
		d = d.Reverse()
	}
	for i := 0; i < 4; i++ {
		if g.MaybeMove(e, d) {
			c.Last = d
			return
		}
		d = d.Clockwise()
		if c.AntiClockwise {
			d = d.Reverse()
		}
	}
}

type EnemyBounceController struct {
	Last Direction
}

func (c *EnemyBounceController) Control(e *Entity, g *Game) {
	if g.MaybeMove(e, c.Last) {
		return
	}
	c.Last = c.Last.Reverse()
}

type EnemyDeflectController struct {
	Last          Direction
	AntiClockwise bool
}

func (c *EnemyDeflectController) Control(e *Entity, g *Game) {
	if g.MaybeMove(e, c.Last) {
		return
	}
	c.Last = c.Last.Clockwise()
	if c.AntiClockwise {
		c.Last = c.Last.Reverse()
	}
}

type SantaController struct {
	Chasing    bool
	StartSpeed float64
	Last       Direction
}

func (c *SantaController) Control(e *Entity, g *Game) {
	if e.Moving == None {
		if !c.Chasing {
			for y := e.Y; g.IsWalkable(e.Type, e.X, y); y++ {
				if g.Player.X == e.X && g.Player.Y == y {
					c.Chasing = true
					c.StartSpeed = e.Speed
					g.MaybeMove(e, Down)
					c.Last = Down
					return
				}
			}
		}
		if !g.MaybeMove(e, c.Last) {
			if c.Last == Up {
				c.Chasing = false
				e.Speed = c.StartSpeed
			} else {
				c.Last = Up
				e.Speed = e.Speed / 2.0
			}
		}
	}
}

const (
	walkableGreen = 11
	walkableEnd   = 17
	walkableRed   = 9
	walkableFall  = 74
)

var walkable = map[EntityType][]int{
	TypePlayer: {walkableGreen, walkableRed, walkableFall, walkableEnd},
	TypeEnemy:  {walkableGreen},
}

func (g *Game) IsWalkable(et EntityType, x, y int) bool {
	tid := g.Map.Layers[walkableIndex].Tiles[XyToIndex(x, y)].ID
	for _, valid := range walkable[et] {
		if int(tid) == valid {
			return true
		}
	}
	return false
}

func (g *Game) MaybeMove(e *Entity, d Direction) bool {
	nx, ny := e.MovesInto(d)
	//fmt.Printf("%d,%d > %d,%d = %d (d:%d)\n", p.X, p.Y, nx, ny, tid, d)
	if g.IsWalkable(e.Type, nx, ny) {
		e.Moving = d
		return true
	}
	return false
}
