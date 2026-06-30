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

	recs := layout.Layout(img)

	if err := btxpack.WriteAtlasImg(recs, *out); err != nil {
		fmt.Printf("failed to write atlas: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	if err := btxpack.WriteAtlasMeta(recs, *meta); err != nil {
		fmt.Printf("failed to write atlas meta: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

}
