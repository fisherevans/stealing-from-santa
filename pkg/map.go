package santa

import (
	"fisherevans.com/stealingfromsanta/res"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	"os"
)

func LoadNewInstance() Instance {
	var err error
	i := Instance{}
	i.Player = &Entity{
		Type:       TypePlayer,
		Sprite:     playerSprite,
		Controller: &PlayerController{},
		Speed:      4,
	}
	i.Speed = 1
	i.IntroTime = introTime

	i.Map, err = tiled.LoadFile("map.tmx", tiled.WithFileSystem(res.FS))
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}

	i.Renderer, err = render.NewRendererWithFileSystem(i.Map, res.FS)
	if err != nil {
		fmt.Printf("map unsupported for rendering: %s", err.Error())
		os.Exit(2)
	}

	for x := 0; x < i.Map.Width; x++ {
		for y := 0; y < i.Map.Height; y++ {
			t := i.Map.Layers[entityIndex].Tiles[XyToIndex(x, y)]
			if t.Nil {
				continue
			}
			switch t.ID {
			case 0:
				i.Player.X = x
				i.Player.Y = y
			case entityTidHugLeft:
				i.Mobs = append(i.Mobs, &Entity{
					Sprite:     elfSprite,
					Speed:      enemySpeed,
					Controller: &EnemyHugController{Last: Left},
					X:          x,
					Y:          y,
				})
			case entityTidHugRight:
				i.Mobs = append(i.Mobs, &Entity{
					Sprite:     elfSprite,
					Speed:      enemySpeed,
					Controller: &EnemyHugController{Last: Left, AntiClockwise: true},
					X:          x,
					Y:          y,
				})
			case entityTidLR:
				i.Mobs = append(i.Mobs, &Entity{
					Sprite:     elfSprite,
					Speed:      enemySpeed,
					Controller: &EnemyBounceController{Last: Left},
					X:          x,
					Y:          y,
				})
			case entityTidUD:
				i.Mobs = append(i.Mobs, &Entity{
					Sprite:     elfSprite,
					Speed:      enemySpeed,
					Controller: &EnemyBounceController{Last: Up},
					X:          x,
					Y:          y,
				})
			case entityTidDeflectR:
				i.Mobs = append(i.Mobs, &Entity{
					Sprite:     elfSprite,
					Speed:      enemySpeed,
					Controller: &EnemyDeflectController{Last: Up},
					X:          x,
					Y:          y,
				})
			case entityTidDeflectL:
				i.Mobs = append(i.Mobs, &Entity{
					Sprite:     elfSprite,
					Speed:      enemySpeed,
					Controller: &EnemyDeflectController{Last: Down, AntiClockwise: true},
					X:          x,
					Y:          y,
				})
			case entityTidSanta:
				i.Mobs = append(i.Mobs, &Entity{
					Sprite:     santaSprite,
					Speed:      enemySpeed * 3,
					Controller: &SantaController{Last: None},
					X:          x,
					Y:          y,
				})
			}
		}
	}

	i.RenderedMap = RenderMap(i.Renderer)

	return i
}

func RenderMap(r *render.Renderer) *ebiten.Image {
	r.Clear()
	r.RenderLayer(mapIndex)
	r.RenderLayer(propIndex)
	return ebiten.NewImageFromImage(r.Result)
}

func XyToIndex(x, y int) int {
	return x + y*mapWidth
}

func IndexToXY(i int) (int, int) {
	return i / mapWidth, i % mapWidth
}
