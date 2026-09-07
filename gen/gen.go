package gen

import (
	"io"

	"github.com/intervinn/btxpack/layout"
)

type Generator func(f io.Writer, a *layout.Atlas) error
