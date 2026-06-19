package main

import (
	"log"
	"image"

	"github.com/solarlune/goaseprite"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

    "fmt"
    "os"

    "github.com/lafriks/go-tiled"
)
const mapPath = "assets/map_1.tmx" // Path to your Tiled Map.
var img *ebiten.Image
var gameMap *tiled.Map 
var speed = 2.0

const CharacterSpriteFile = "character_base_16x16.json"
const CharacterSpriteDirectory = "./assets"
const CharacterSpriteStartAnim = "idle_down "

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

func GetObjectGroup(name string) *tiled.ObjectGroup {
	for _, objectGroup := range gameMap.ObjectGroups {
		if objectGroup.Name == name {
			return objectGroup 
		}
	}
	return nil 
}

func GetObjectFromObjectLayer(objectGroup *tiled.ObjectGroup, name string) *tiled.Object{
	for _, object := range objectGroup.Objects {
		if object.Name == name {
			return object 
		}
	}	
	return nil 
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
type Direction int

const (
    Up Direction = iota
    Down 
    Left 
    Right 
)

var CurrentDirection = Down 

type Game struct{
	Char *Character 
}

type Character struct {
	Sprite *SpriteAnim 
	PosX float64
	PosY float64	
}

type SpriteAnim struct {
	Sprite    *goaseprite.File
	AsePlayer *goaseprite.Player
	Img       *ebiten.Image	
}

func (sprite *SpriteAnim) Play(animation string) {
	sprite.AsePlayer.Play(animation)
}
func (sprite *SpriteAnim) Update(delta float32) {
	sprite.AsePlayer.Update(delta)
}

func NewSpriteAnim(
	file string, directory string, start_anim string) *SpriteAnim {
	//sprite, err := goaseprite.Open("character_base_16x16.json", os.DirFS("./assets"))
	
	sprite, err := goaseprite.Open(file, os.DirFS(directory))

	if err != nil {
		panic(err)
	}
	spriteAnim := &SpriteAnim {
		Sprite: sprite, 
	}

	spriteAnim.AsePlayer = spriteAnim.Sprite.CreatePlayer()
	img, _, err := ebitenutil.NewImageFromFile(spriteAnim.Sprite.ImagePath)
	if err != nil {
		panic(err)
	}

	// game.Sprite.PlaySpeed = 2

	spriteAnim.Img = img
	spriteAnim.AsePlayer.Play(start_anim)
	return spriteAnim
}

func (sprite *SpriteAnim) Draw(Screen *ebiten.Image, PosX float64,  PosY float64) {
	opts := &ebiten.DrawImageOptions{}
    opts.GeoM.Translate(
        float64(PosX),
        float64(PosY),
    )
	sub := sprite.Img.SubImage(image.Rect(sprite.AsePlayer.CurrentFrameCoords()))

	Screen.DrawImage(sub.(*ebiten.Image), opts)
}

func NewCharacter(obj *tiled.Object, 
			sprite_file string, 
			sprite_directory string, 
			default_anim string ) *Character{
	character := &Character{}
	// sprite, err := goaseprite.Open("character_base_16x16.json", os.DirFS("./assets"))
	// if err != nil {
	// 	panic(err)
	// }
	// character := &Character{
	// 	Sprite: sprite,
	// }

	// character.AsePlayer = character.Sprite.CreatePlayer()
	// img, _, err := ebitenutil.NewImageFromFile(character.Sprite.ImagePath)
	// if err != nil {
	// 	panic(err)
	// }

	// // game.Sprite.PlaySpeed = 2

	// character.Img = img
	// character.AsePlayer.Play("idle_down")
	// character.PosX = obj.X 
	// character.PosY = obj.Y 
	spriteAnim := NewSpriteAnim(sprite_file, sprite_directory, default_anim)
	character.Sprite = spriteAnim
	character.PosX = obj.X 
	character.PosY = obj.Y 
	return character
}

func (c *Character) Draw(screen *ebiten.Image) {
	c.Sprite.Draw(screen, c.PosX, c.PosY)
}

func (c *Character) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		c.Sprite.Play("walk_up")
		c.PosY = c.PosY - speed 
		CurrentDirection = Up 
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		c.Sprite.Play("walk_down")
		c.PosY = c.PosY + speed 
		CurrentDirection = Down 
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		c.Sprite.Play("walk_left")
		c.PosX = c.PosX - speed 
		CurrentDirection = Left 
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		c.Sprite.Play("walk_right")
		c.PosX = c.PosX + speed 
		CurrentDirection = Right 
	} else {
	    switch CurrentDirection {
	    case Up:
	       c.Sprite.Play("idle_up")
	    case Down:
	       c.Sprite.Play("idle_down")
	    case Left:
	       c.Sprite.Play("idle_left")
	    case Right:
	       c.Sprite.Play("idle_right")	       
	    default:
	       c.Sprite.Play("idle_down")
	    }
	}

	c.Sprite.Update(float32(1.0 / 60.0))

}

func (game *Game) Update() error {
	game.Char.Update()
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")	
	DrawTiledLayer(screen, "base")
	DrawTiledLayer(screen, "over")
	game.Char.Draw(screen) 
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	objectsLayer := GetObjectGroup("objects") 	// Search for Object Layer called "objects"
	object := GetObjectFromObjectLayer(objectsLayer, "PlayerStart")
	fmt.Println(object)


	Char := NewCharacter(object, 
			CharacterSpriteFile, 
			CharacterSpriteDirectory, 
			CharacterSpriteStartAnim)

	game := &Game {
		Char: Char, 
	}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}