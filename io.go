package btxpack

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/intervinn/btxpack/codegen"
	"github.com/intervinn/btxpack/layout"
)

func IsSupported(name string) bool {
	return strings.HasSuffix(name, ".png")
}

func ScanDir(root string) ([]layout.Img, error) {
	res := []layout.Img{}

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

		res = append(res, layout.Img{
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

func getWidth(i []layout.Rec) int {
	r := 0
	for _, v := range i {
		r += v.Bounds().Dx()
	}
	return r
}

func getHeight(i []layout.Rec) int {
	r := 0
	for _, v := range i {
		if v.Bounds().Dy() > r {
			r = v.Bounds().Dy()
		}
	}
	return r
}

func WriteAtlasImg(recs []layout.Rec, to string) error {
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

	fmt.Printf("writing image to %s\n", to)

	return png.Encode(f, img)
}

func WriteAtlasImgSide(recs []layout.Rec, to string, side int) error {
	w := side
	h := side

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

	fmt.Printf("writing image to %s\n", to)

	return png.Encode(f, img)
}

func IsExtSupported(s string) bool {
	return s == ".c" || s == ".h" || s == ".json"
}

func WriteAtlasMeta(recs []layout.Rec, to string, raylib bool) error {
	ext := path.Ext(to)
	if !IsExtSupported(ext) {
		return fmt.Errorf("unsupported format")
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	metas := make([]layout.Meta, len(recs))
	for i, v := range recs {
		name, err := filepath.Rel(wd, v.Src)
		if err != nil {
			return err
		}

		metas[i] = layout.Meta{
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

	switch ext {
	case ".json":
		return codegen.WriteAtlasJson(metas, f)
	case ".c", ".h":
		if raylib {
			return codegen.WriteAtlasRaylibC(metas, f)
		}
		return codegen.WriteAtlasC(metas, f)
	}

	return fmt.Errorf("unknown extension: %s", ext)
}
