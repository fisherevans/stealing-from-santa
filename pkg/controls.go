package santa

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Controller interface {
	Control(e *Entity, g *Game)
}

type PlayerController struct{}

func (c *PlayerController) Control(p *Entity, g *Game) {
	if p.Moving == None {
		if ebiten.IsKeyPressed(ebiten.KeyW) && g.MaybeMove(p, Up) {
			return
		} else if ebiten.IsKeyPressed(ebiten.KeyS) && g.MaybeMove(p, Down) {
			return
		} else if ebiten.IsKeyPressed(ebiten.KeyA) && g.MaybeMove(p, Left) {
			return
		} else if ebiten.IsKeyPressed(ebiten.KeyD) && g.MaybeMove(p, Right) {
			return
		}
	}
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
