package gen

import (
	"io"

	"github.com/intervinn/btxpack/layout"
)

type Generator func(f io.Writer, a *layout.Atlas, args []string) error

func ExportNoOp(f io.Writer, a *layout.Atlas, args []string) error {
	return nil
}
