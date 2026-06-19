package main 

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var Projectiles []*Projectile 
var ColorRed = color.RGBA{R: 255, G: 0, B: 0, A: 255}
var ColorGreen = color.RGBA{R: 0, G: 255, B: 0, A: 255}
var ColorWhite = color.RGBA{R: 255, G: 255, B: 255, A: 255}
var ColorRedShader = []float32{1.0, 0.0, 0.0, 1.0}
var ColorGreenShader = []float32{0.0, 1.0, 0.0, 1.0}


func init() {
	Projectiles = make([]*Projectile, 0)
}

func NewProjectile(
		BaseX float64, 
		BaseY float64, 
		Dir Direction, 
		Speed float64, ProjColor string) *Projectile{
	proj := &Projectile{
		X: BaseX, 
		Y: BaseY,
		Speed: Speed,  
	}
	if ProjColor == "red" {
		proj.Color = ColorRed
	} else if ProjColor == "green" {
		proj.Color = ColorGreen
	} else {
		proj.Color = ColorWhite 
	}

	switch Dir {
    case Up:
       proj.Y = proj.Y - 10 
    case Down:
       proj.Y = proj.Y + 10 
    case Left:
       proj.X = proj.X - 10 
    case Right:
       proj.X = proj.X + 10       
    default:
       proj.Y = proj.Y + 10 
    }
    proj.Dir = Dir 
    return proj 
}

type Projectile struct{
	X float64 
	Y float64 
	Dir Direction 
	Speed float64
	Color color.RGBA 
}

func (g *Projectile) Update()  { 
	switch g.Dir {
    case Up:
       g.Y = g.Y - g.Speed 
    case Down:
       g.Y = g.Y + g.Speed 
    case Left:
       g.X = g.X - g.Speed 
    case Right:
       g.X = g.X + g.Speed       
    default:
       g.Y = g.Y + g.Speed 
    }

}

func (g *Projectile) Draw(screen *ebiten.Image) {
	// Draw a solid circle at (x=320, y=240) with a radius of 50 in white
	vector.DrawFilledCircle(screen, float32(g.X), float32(g.Y), 2, g.Color, true)
    vector.StrokeCircle(
        screen,
		float32(g.X), float32(g.Y), 3, 2,
        color.Black,
        true,         // anti-alias
    )

}
