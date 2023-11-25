package santa

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type Direction int

const (
	None Direction = iota
	Up
	Down
	Left
	Right
)

func (d Direction) Reverse() Direction {
	if d == Up {
		return Down
	} else if d == Down {
		return Up
	} else if d == Left {
		return Right
	} else if d == Right {
		return Left
	}
	return None
}

func (d Direction) Clockwise() Direction {
	if d == Up {
		return Right
	} else if d == Down {
		return Left
	} else if d == Left {
		return Up
	} else if d == Right {
		return Down
	}
	return None
}

func (d Direction) AntiCW() Direction {
	return d.Clockwise().Reverse()
}

type EntityType int

const (
	TypeEnemy = iota
	TypePlayer
)

type Entity struct {
	Speed      float64
	Controller Controller
	Type       EntityType
	Sprite     *ebiten.Image

	X              int
	Y              int
	Moving         Direction
	MoveTransition float64
}

func (e *Entity) EffectivePosition() (float64, float64) {
	x, y := float64(e.X), float64(e.Y)
	if e.Moving == Up {
		y -= math.Min(e.MoveTransition, 1.0)
	} else if e.Moving == Down {
		y += math.Min(e.MoveTransition, 1.0)
	} else if e.Moving == Left {
		x -= math.Min(e.MoveTransition, 1.0)
	} else if e.Moving == Right {
		x += math.Min(e.MoveTransition, 1.0)
	}
	return x, y
}

func (e *Entity) Update(delta float64) {
	if e.Moving != None {
		e.MoveTransition += delta * e.Speed
		if e.MoveTransition > 1.0 {
			if e.Moving == Up {
				e.Y -= 1
			} else if e.Moving == Down {
				e.Y += 1
			} else if e.Moving == Left {
				e.X -= 1
			} else if e.Moving == Right {
				e.X += 1
			}
			e.MoveTransition = 0
			e.Moving = None
		}
	}
}

func (e *Entity) Occupies(x, y int) bool {
	ex, ey := e.X, e.Y
	if ex == x && ey == y {
		return true
	}
	if e.Moving == Up {
		ey -= 1
	} else if e.Moving == Down {
		ey += 1
	} else if e.Moving == Left {
		ex -= 1
	} else if e.Moving == Right {
		ex += 1
	}
	return ex == x && ey == y
}

func (e *Entity) CollisionSpace() (int, int) {
	ex, ey := e.X, e.Y
	if e.MoveTransition < 0.5 {
		return ex, ey
	}
	if e.Moving == Up {
		ey -= 1
	} else if e.Moving == Down {
		ey += 1
	} else if e.Moving == Left {
		ex -= 1
	} else if e.Moving == Right {
		ex += 1
	}
	return ex, ey
}

func (e *Entity) MovesInto(d Direction) (x, y int) {
	if d == Up {
		return e.X, e.Y - 1
	} else if d == Down {
		return e.X, e.Y + 1
	} else if d == Left {
		return e.X - 1, e.Y
	} else if d == Right {
		return e.X + 1, e.Y
	}
	return e.X, e.Y

}
