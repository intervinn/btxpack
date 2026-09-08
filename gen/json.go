package gen

import (
	"encoding/json"
	"io"

	"github.com/intervinn/btxpack/layout"
)

type jsonEntry struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

func ExportJSON(f io.Writer, a *layout.Atlas, args []string) error {
	res := map[string]jsonEntry{}
	for _, v := range a.Recs {
		res[v.Name] = jsonEntry{v.X, v.Y, v.W, v.H}
	}

	return json.NewEncoder(f).Encode(res)
}
