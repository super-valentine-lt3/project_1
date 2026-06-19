package main

import (
	"log"
	"image"

	"github.com/solarlune/goaseprite"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

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
var fade *Fade 

type Direction int

const (
    Up Direction = iota
    Down 
    Left 
    Right 
)

var CurrentDirection = Down 

type Game struct{
	Map CollisionMap  
	Chars  []*Character 
}

type CollisionMap struct {
    Width  int
    Height int
    Solid  [][]bool
}

func NewCollisionMap() CollisionMap {
	collision := make([][]bool, gameMap.Height)

	for y := range collision {
	    collision[y] = make([]bool, gameMap.Width)
	}	

	for _, layer := range gameMap.Layers {
	    for pos, tile := range layer.Tiles {
	        if tile.Nil {
	            continue
	        }

	        tt, err := tile.Tileset.GetTilesetTile(tile.ID)
	        if err != nil { continue }
	        if tt.Properties.GetBool("solid") {
	            x := pos % gameMap.Width
	            y := pos / gameMap.Width

	            collision[y][x] = true
	        }
	    }
	}

	return CollisionMap{
		Width: gameMap.Width, 
		Height: gameMap.Height, Solid: collision}
}

const tileWidth = 16
const tileHeight = 16
func (cm *CollisionMap) IsSolid(x, y float64) bool {
    tx := int(x) / tileWidth
    ty := int(y) / tileHeight

    return cm.Solid[ty][tx]
}

const PlayerWidth = 16
const PlayerHeight = 16

func (cm *CollisionMap) CanMove(newX, newY float64) bool {
    return !cm.IsSolid(newX, newY) &&
           !cm.IsSolid(newX+PlayerWidth-1, newY) &&
           !cm.IsSolid(newX, newY+PlayerHeight-1) &&
           !cm.IsSolid(newX+PlayerWidth-1, newY+PlayerHeight-1)
}

type Character struct {
	Sprite *SpriteAnim 
	PosX float64
	PosY float64	
	Actions map[string]ebiten.Key 
	Data map[string]string 
	CurrentDirection Direction 
}

func (c *Character) Get(name string) string {
	return c.Data[name]  
}

func (c *Character) SetData(name string, value string) {
	c.Data[name] = value 
}

func (c *Character) SetAction(name string, key ebiten.Key) {
	c.Actions[name] = key 
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
	character.Actions = make(map[string]ebiten.Key)
	character.Data = make(map[string]string)
	character.CurrentDirection = Down 

	return character
}

func (c *Character) Draw(screen *ebiten.Image) {
	c.Sprite.Draw(screen, c.PosX, c.PosY)
}

func (c *Character) Update(Map *CollisionMap) {
	//if ebiten.IsKeyPressed(ebiten.KeyUp) {
	if ebiten.IsKeyPressed(c.Actions["move_up"]) {	
		tempY :=  c.PosY - speed 
		tempX := c.PosX 
		if Map.CanMove(tempX, tempY) {
			c.Sprite.Play("walk_up")
			c.PosY = c.PosY - speed 
			c.CurrentDirection = Up 
		}
	} else if ebiten.IsKeyPressed(c.Actions["move_down"]) {
		tempY :=  c.PosY + speed 
		tempX := c.PosX 
		if Map.CanMove(tempX, tempY) {
			c.Sprite.Play("walk_down")
			c.PosY = c.PosY + speed 
			c.CurrentDirection = Down 
		}
	} else if ebiten.IsKeyPressed(c.Actions["move_left"]) {
		tempX := c.PosX - speed 
		tempY := c.PosY 
		if Map.CanMove(tempX, tempY) {
			c.Sprite.Play("walk_left")
			c.PosX = c.PosX - speed 
			c.CurrentDirection = Left 
		}
	} else if ebiten.IsKeyPressed(c.Actions["move_right"]) {
		tempX := c.PosX + speed 
		tempY := c.PosY 
		if Map.CanMove(tempX, tempY) {
			c.Sprite.Play("walk_right")
			c.PosX = c.PosX + speed 
			c.CurrentDirection = Right 
		}
	} else {
	    switch c.CurrentDirection {
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

	if inpututil.IsKeyJustPressed(c.Actions["shoot_projectile"]) {
		Projectiles = append(Projectiles, NewProjectile(c.PosX, c.PosY, c.CurrentDirection, 2, c.Get("projectile_color")))
	}
}

func (game *Game) Update() error {
	for _,  char := range game.Chars {
		char.Update(&game.Map)
	}

	if fade == nil &&  ebiten.IsKeyPressed(ebiten.KeyF) {
		f := NewFade(1, 0) 
		fade = &f 
	} else if fade != nil {
		if fade.Finished() {
			fade = nil 
		} 
	}
	for _, proj := range Projectiles {
		proj.Update()
	}
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")	
	DrawTiledLayer(screen, "base")
	DrawTiledLayer(screen, "over")

	for _, char := range game.Chars {
		char.Draw(screen) 
	}

	for _, proj := range Projectiles {
		proj.Draw(screen)
	}

	if fade != nil {
		fade.Draw(screen)
	} 
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	
	objectsLayer := GetObjectGroup("objects") 	// Search for Object Layer called "objects"
	object := GetObjectFromObjectLayer(objectsLayer, "PlayerStart")
	object2 := GetObjectFromObjectLayer(objectsLayer, "PlayerStart2")
	
	fmt.Println(object)

	characters := make([]*Character, 0)

	Char := NewCharacter(object, 
			CharacterSpriteFile, 
			CharacterSpriteDirectory, 
			CharacterSpriteStartAnim)
	Char.SetAction("move_up", ebiten.KeyUp)
	Char.SetAction("move_left", ebiten.KeyLeft)
	Char.SetAction("move_right", ebiten.KeyRight)
	Char.SetAction("move_down", ebiten.KeyDown)
	Char.SetAction("shoot_projectile", ebiten.KeySpace)	
	Char.SetData("projectile_color", "red")

	characters = append(characters, Char)

	Char2 := NewCharacter(object2, 
			CharacterSpriteFile, 
			CharacterSpriteDirectory, 
			CharacterSpriteStartAnim)
	Char2.SetAction("move_up", ebiten.KeyW)
	Char2.SetAction("move_left", ebiten.KeyA)
	Char2.SetAction("move_right", ebiten.KeyD)
	Char2.SetAction("move_down", ebiten.KeyS)
	Char2.SetAction("shoot_projectile", ebiten.KeyEnter)
	Char2.SetData("projectile_color", "green")

	characters = append(characters, Char2)


	game := &Game {
		Chars: characters, 
		Map: NewCollisionMap(), 
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}