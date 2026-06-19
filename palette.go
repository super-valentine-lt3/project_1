package main 
import (
	"github.com/hajimehoshi/ebiten/v2"
)

var paletteShader *ebiten.Shader

func init() {
    var err error

    paletteShader, err = ebiten.NewShader([]byte(`
package main

var PlayerColor vec4

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
    pixel := imageSrc0At(texCoord)

    skin := vec3(
        246.0/255.0,
        187.0/255.0,
        148.0/255.0,
    )

    diff := distance(pixel.rgb, skin)

    if diff < 0.08 {
        return vec4(
            PlayerColor.rgb,
            pixel.a,
        )
    }

    return pixel
}
`))

    if err != nil {
        panic(err)
    }
}