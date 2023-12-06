package santa

import (
	"fisherevans.com/stealingfromsanta/res"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/lafriks/go-tiled"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"log"
	"os"
)

const (
	entityTidHugLeft  = 1
	entityTidHugRight = 2
	entityTidLR       = 3
	entityTidUD       = 4
	entityTidDeflectL = 5
	entityTidDeflectR = 6
	entityTidSanta    = 7
)

var moneyTids = map[int]int{
	89:  1,
	94:  2,
	99:  3,
	104: 4,
	109: 7,

	//secrets
	92:  10,
	107: 40,
	115: 100,
}

var (
	playerSprite *ebiten.Image
	elfSprite    *ebiten.Image
	santaSprite  *ebiten.Image

	smallFont font.Face
	bigFont   font.Face

	entityIndex   int
	walkableIndex int
	propIndex     int
	mapIndex      int

	mapWidth int

	highScore  = 0
	enemySpeed = 3.0

	introTime = 5.0
)

func init() {
	// FONTS
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	smallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    8,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	bigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	// SPRITES
	playerSprite, _, err = ebitenutil.NewImageFromFileSystem(res.FS, "player.png")
	if err != nil {
		fmt.Printf("failed to load image: %s", err.Error())
		os.Exit(2)
	}

	elfSprite, _, err = ebitenutil.NewImageFromFileSystem(res.FS, "elf.png")
	if err != nil {
		fmt.Printf("failed to load image: %s", err.Error())
		os.Exit(2)
	}

	santaSprite, _, err = ebitenutil.NewImageFromFileSystem(res.FS, "santa.png")
	if err != nil {
		fmt.Printf("failed to load image: %s", err.Error())
		os.Exit(2)
	}

	// MAP DETAILS
	mapFile, err := tiled.LoadFile("map.tmx", tiled.WithFileSystem(res.FS))
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}

	for id, l := range mapFile.Layers {
		if l.Name == "map" {
			mapIndex = id
		} else if l.Name == "props" {
			propIndex = id
		} else if l.Name == "entities" {
			entityIndex = id
		} else if l.Name == "walkable" {
			walkableIndex = id
		}
	}

	highScore = 0
	for _, t := range mapFile.Layers[propIndex].Tiles {
		for tid, m := range moneyTids {
			if int(t.ID) == tid {
				highScore += m
			}
		}
	}

	mapWidth = mapFile.Width
}
