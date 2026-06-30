package codegen

import (
	"encoding/json"
	"io"

	"github.com/intervinn/btxpack/layout"
)

func WriteAtlasJson(metas []layout.Meta, to io.Writer) error {
	e := json.NewEncoder(to)
	e.SetIndent("", " ")
	e.Encode(metas)

	return nil
}
