package main 

import (
	"time"
	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

const ScreenWidth = 320
const ScreenHeight = 240

type Fade struct {
    FromAlpha float64
    ToAlpha   float64

    Duration time.Duration
    Start    time.Time

    fadeImage *ebiten.Image
}

func (f *Fade) Finished() bool {
	return f.Alpha() == f.ToAlpha 
}

func (f *Fade) Alpha() float64 {
    t := float64(time.Since(f.Start)) / float64(f.Duration)
    if t > 1 {
        t = 1
    }

    return f.FromAlpha + (f.ToAlpha-f.FromAlpha)*t
}

func (f *Fade) Draw(screen *ebiten.Image) {

    op := &ebiten.DrawImageOptions{}
    op.ColorScale.Scale(0, 0, 0, float32(f.Alpha()))

    screen.DrawImage(f.fadeImage, op)
}

func NewFade(FromAlpha, ToAlpha float64) Fade {
	fade := ebiten.NewImage(ScreenWidth, ScreenHeight)
	fade.Fill(color.White)

	return Fade{
	    Duration: time.Second,
	    Start:    time.Now(),
	    fadeImage: fade,
	    FromAlpha: FromAlpha, 
	    ToAlpha: ToAlpha,  
	}
} 