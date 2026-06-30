package layout

import (
	"fmt"
	"math"
	"slices"
)

// shelf algo
func Shelf(imgs []Img) ([]Rec, int) {
	// 0. find a power two for max width
	// 1. sort all imgs by height
	// 2. place imgs left to right
	// 3. when max width is exceeded, start a new row

	s := 0
	for _, i := range imgs {
		s += i.Bounds().Dx() * i.Bounds().Dy()
	}

	iside := math.Sqrt(float64(s))
	bside := iside * 1.1
	var side int

	for i := 0; ; i++ {
		n := math.Pow(2, float64(i))
		if n >= bside {
			side = int(n)
			break
		}
	}
	side = min(side, 4096)

	fmt.Printf("side: %d\n", side)

	imgsh := make([]Img, len(imgs))
	copy(imgsh, imgs)
	slices.SortStableFunc(imgsh, func(a Img, b Img) int {
		if a.Bounds().Dy() > b.Bounds().Dy() {
			return -1
		}
		if a.Bounds().Dy() < b.Bounds().Dy() {
			return 1
		}
		return 0
	})

	fmt.Println("SORTED")
	for _, i := range imgsh {
		fmt.Printf("%v\n", i.String())
	}

	recs := make([]Rec, 0)

	mx := 0 // width counter
	tl := 0 // current row tallest rec
	x := 0
	y := 0
	for _, i := range imgsh {
		w := i.Bounds().Dx()
		h := i.Bounds().Dy()

		if mx+w > side {
			y += tl
			mx = 0
			x = 0
			tl = 0
		}

		tl = max(tl, h)

		recs = append(recs, Rec{
			Img: i,
			X:   x,
			Y:   y,
		})
		mx += w
		x += w
	}

	fmt.Println(recs)

	return recs, side
}
