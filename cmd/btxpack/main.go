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
	"strings"

	"github.com/intervinn/btxpack/gen"
	"github.com/intervinn/btxpack/layout"
)

var gens = map[string]gen.Generator{
	".c":    gen.ExportC,
	".h":    gen.ExportC,
	".json": gen.ExportJSON,
	"":      gen.ExportNoOp,
}

var packers = map[string]layout.Packer{
	"shelf": layout.Shelf,
}

func scanDir(root string) ([]layout.Img, error) {
	res := make([]layout.Img, 0)
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir(root, func(fpath string, d fs.DirEntry, err error) error {
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

		name, err := filepath.Rel(cwd, path.Join(root, d.Name()))
		if err != nil {
			return err
		}
		res = append(res, layout.Img{
			Name:   name,
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

func ignoreUndefinedFlags(args []string, fs *flag.FlagSet) []string {
	var clean []string
	for i := 0; i < len(args); i++ {
		arg := args[i]

		if strings.HasPrefix(arg, "-") {
			name := strings.TrimLeft(arg, "-")
			if idx := strings.Index(name, "="); idx != -1 {
				name = name[:idx]
			}

			if fs.Lookup(name) == nil {
				if !strings.Contains(arg, "=") && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					i++
				}
				continue
			}
		}
		clean = append(clean, arg)
	}
	return clean
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("more args")
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("failed to fetch wd:", err)
	}

	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	src := fs.String("src", "source", "source directory")
	alg := fs.String("alg", "shelf", "layout algorithm")
	gen := fs.String("gen", "", "metadata destination")
	out := fs.String("o", "", "(png) atlas destination file")
	fs.Parse(ignoreUndefinedFlags(os.Args[1:], fs))

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
