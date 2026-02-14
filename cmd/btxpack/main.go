package main

import (
	"os"
	"path"

	"github.com/intervinn/btxpack"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "btxpack",
	Short: "Automatic texture packer and atlas generator",

	Args: cobra.ExactArgs(3),

	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		arg := args[0]

		img, err := btxpack.ScanDir(path.Join(cwd, arg))
		if err != nil {
			return err
		}

		recs := btxpack.Layout(img)

		if err := btxpack.WriteAtlasImg(recs, args[1]); err != nil {
			return err
		}

		if err := btxpack.WriteAtlasMeta(recs, args[2]); err != nil {
			return err
		}

		return nil
	},
}

func main() {
	root.Execute()
}
