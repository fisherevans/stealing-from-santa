package main

import (
	santa "fisherevans.com/stealingfromsanta/pkg"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
	"log"
)

func main() {
	g := &santa.Game{
		ScreenWidth:  320,
		ScreenHeight: 180,
		TileSize:     16,
	}
	g.Instance = santa.LoadNewInstance()
	ebiten.SetWindowSize(int(g.ScreenWidth*3), int(g.ScreenHeight*3))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Stealing From Santa")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
