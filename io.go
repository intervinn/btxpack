package btxpack

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Img struct {
	image.Rectangle
	Image image.Image
	Src   string
}

type Meta struct {
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	W    int    `json:"w"`
	H    int    `json:"h"`
}

func IsSupported(name string) bool {
	return strings.HasSuffix(name, ".png")
}

func ScanDir(root string) ([]Img, error) {
	res := []Img{}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !IsSupported(d.Name()) {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			return err
		}

		res = append(res, Img{
			Rectangle: img.Bounds(),
			Image:     img,
			Src:       path,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func getWidth(i []Rec) int {
	r := 0
	for _, v := range i {
		r += v.Bounds().Dx()
	}
	return r
}

func getHeight(i []Rec) int {
	r := 0
	for _, v := range i {
		if v.Bounds().Dy() > r {
			r = v.Bounds().Dy()
		}
	}
	return r
}

func WriteAtlasImg(recs []Rec, to string) error {
	w := getWidth(recs)
	h := getHeight(recs)

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for _, v := range recs {
		fmt.Printf("drawing %v\n", v.Src)
		draw.Draw(img, v.Bounds().Add(image.Pt(v.X, v.Y)), v.Image, image.Point{0, 0}, draw.Src)
	}

	os.MkdirAll(filepath.Dir(to), os.ModePerm)
	f, err := os.Create(to)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func WriteAtlasMeta(recs []Rec, to string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	metas := make([]Meta, len(recs))
	for i, v := range recs {
		name, err := filepath.Rel(wd, v.Src)
		if err != nil {
			return err
		}

		metas[i] = Meta{
			Name: name,
			X:    v.X,
			Y:    v.Y,
			W:    v.Bounds().Dx(),
			H:    v.Bounds().Dy(),
		}
	}

	fmt.Printf("writing atlas to %s\n", to)
	os.MkdirAll(filepath.Dir(to), os.ModePerm)
	f, err := os.Create(to)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", " ")
	e.Encode(metas)

	return nil
}
