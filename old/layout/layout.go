package layout

import (
	"fmt"
	"image"
)

type Img struct {
	image.Rectangle
	Image image.Image
	Src   string
}

type Rec struct {
	Img
	X int
	Y int
}

func (i Img) String() string {
	return fmt.Sprintf("Img{%dx%d}", i.Bounds().Dx(), i.Bounds().Dy())
}

func (r Rec) String() string {
	return fmt.Sprintf("Rec{%d,%d, %dx%d}", r.X, r.Y, r.Bounds().Dx(), r.Bounds().Dy())
}

// rn it just puts all sprites in line with eachother
// need to make some more efficient algorithm in future
func Layout(imgs []Img) []Rec {
	out := make([]Rec, len(imgs))
	ox := 0
	for i, v := range imgs {
		out[i] = Rec{
			Img: v,
			X:   ox,
			Y:   0,
		}
		ox += v.Bounds().Dx()
	}

	return out
}

type Meta struct {
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	W    int    `json:"w"`
	H    int    `json:"h"`
}
