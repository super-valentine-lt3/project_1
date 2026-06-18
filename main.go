package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

    "fmt"
    "os"

    "github.com/lafriks/go-tiled"
)
const mapPath = "assets/map_1.tmx" // Path to your Tiled Map.
var img *ebiten.Image
var gameMap *tiled.Map 

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("assets/16oga.png")
	if err != nil {
		log.Fatal(err)
	}

    // Parse .tmx file.
    gameMap, err = tiled.LoadFile(mapPath)
    if err != nil {
        fmt.Printf("error parsing map: %s", err.Error())
        os.Exit(2)
    }


}

func DrawTiledLayer(screen *ebiten.Image, name string) {
	for _, layer := range gameMap.Layers {
		if layer.Name != name { continue }

	    for pos, tile := range layer.Tiles {
	    	if tile.Nil { continue }
	    	
	        rect := tile.GetTileRect()

	        sub := img.SubImage(rect).(*ebiten.Image)
	        tileX := pos % gameMap.Width
	        tileY := pos / gameMap.Width

	        op := &ebiten.DrawImageOptions{}
	        op.GeoM.Translate(
	            float64(tileX*gameMap.TileWidth),
	            float64(tileY*gameMap.TileHeight),
	        )

	        screen.DrawImage(sub, op)
	    }
	}	
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")	
	DrawTiledLayer(screen, "base")
	DrawTiledLayer(screen, "over")

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}