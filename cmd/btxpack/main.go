package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/intervinn/btxpack/gen"
	"github.com/intervinn/btxpack/layout"
)

var gens = map[string]gen.Generator{
	".c":    gen.ExportC,
	".json": gen.ExportJSON,
	"":      gen.ExportNoOp,
}

var packers = map[string]layout.Packer{
	"shelf": layout.Shelf,
}

func scanDir(root string) ([]layout.Img, error) {
	res := make([]layout.Img, 0)
	err := filepath.WalkDir(root, func(fpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		f, err := os.Open(fpath)
		if err != nil {
			return err
		}

		img, _, err := image.Decode(f)
		if err != nil {
			return err
		}

		res = append(res, layout.Img{
			Name:   path.Join(root, d.Name()),
			Width:  img.Bounds().Dx(),
			Height: img.Bounds().Dy(),
		})

		f.Close()
		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, err
}

func writePacked(to string, recs []layout.Rec, w int, h int) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for _, v := range recs {
		fi, err := os.Open(v.Name)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer fi.Close()

		i, _, err := image.Decode(fi)
		if err != nil {
			return fmt.Errorf("failed to decode image: %w", err)
		}

		draw.Draw(img, i.Bounds().Add(image.Pt(v.X, v.Y)), i, image.Point{0, 0}, draw.Src)
	}

	os.MkdirAll(filepath.Dir(to), os.ModePerm)
	f, err := os.Create(to)
	if err != nil {
		return err
	}

	defer f.Close()

	return png.Encode(f, img)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("more args")
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("failed to fetch wd:", err)
	}

	src := flag.String("src", "source", "source directory")
	alg := flag.String("alg", "shelf", "layout algorithm")
	gen := flag.String("gen", "", "metadata destination")
	out := flag.String("o", "", "(png) atlas destination file")
	flag.Parse()

	a := os.Args[2:]

	if *out == "" {
		log.Fatalln("out file is required")
	}

	packer, ok := packers[*alg]
	if !ok {
		log.Fatalln("invalid algorithm:", *alg)
	}

	generator, ok := gens[filepath.Ext(*gen)]
	if !ok {
		log.Fatalln("invalid generator:", *gen)
	}

	imgs, err := scanDir(path.Join(cwd, *src))
	if err != nil {
		log.Fatalln("failed to scan dir:", err)
	}

	atlas, w, h := packer(imgs)
	if err = writePacked(*out, atlas.Recs, w, h); err != nil {
		log.Fatalln("failed to write packed:", err)
	}

	if err := os.MkdirAll(filepath.Dir(*gen), os.ModePerm); err != nil {
		log.Fatalln("failed to mkdir for meta:", err)
	}

	genf, err := os.Create(path.Join(cwd, *gen))
	if err != nil {
		log.Fatalln("failed to create meta file:", err)
	}
	defer genf.Close()

	if err = generator(genf, atlas, a); err != nil {
		log.Fatalln("failed to write meta:", err)
	}
}
