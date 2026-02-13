package btxpack

type Rec struct {
	Img
	X int
	Y int
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
