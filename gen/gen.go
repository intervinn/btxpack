package gen

import (
	"flag"
	"io"
	"strings"

	"github.com/intervinn/btxpack/layout"
)

type Generator func(f io.Writer, a *layout.Atlas, args []string) error

func ExportNoOp(f io.Writer, a *layout.Atlas, args []string) error {
	return nil
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
