package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/intervinn/btxpack"
	"github.com/intervinn/btxpack/layout"
)

func main() {
	out := flag.String("out", "atlas.png", "destination file (atlas.png)")
	src := flag.String("src", "assets", "assets directory")
	meta := flag.String("meta", "atlas.json", "atlas metadata path (json,c)")
	alg := flag.String("alg", "line", "algorithm (line,shelf)")
	raylib := flag.Bool("raylib", false, "use raylib types for c")

	flag.Parse()

	if out == nil || src == nil || *out == "" || *src == "" {
		fmt.Println("params `out` and `src` are required")
		flag.Usage()
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to fetch wd: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	img, err := btxpack.ScanDir(path.Join(cwd, *src))
	if err != nil {
		fmt.Printf("failed to fetch scan dir: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	var recs []layout.Rec
	switch *alg {
	case "shelf":
		{
			rrecs, side := layout.Shelf(img)
			recs = rrecs
			if err := btxpack.WriteAtlasImgSide(recs, *out, side); err != nil {
				fmt.Printf("failed to write atlas: %v\n", err)
				flag.Usage()
				os.Exit(1)
			}
		}
	case "line":
		{
			rrecs := layout.Layout(img)
			recs = rrecs
			if err := btxpack.WriteAtlasImg(recs, *out); err != nil {
				fmt.Printf("failed to write atlas: %v\n", err)
				flag.Usage()
				os.Exit(1)
			}
		}
	default:
		fmt.Printf("unknown alg: %v\n", *alg)
		flag.Usage()
		os.Exit(1)
	}

	if err := btxpack.WriteAtlasMeta(recs, *meta, *raylib); err != nil {
		fmt.Printf("failed to write atlas meta: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

}
