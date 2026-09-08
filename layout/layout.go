package layout

type Packer func([]Img) (*Atlas, int, int)

type Atlas struct {
	Recs []Rec
}

type Img struct {
	Name   string
	Width  int
	Height int
}

type Rec struct {
	Name string
	W    int
	H    int
	X    int
	Y    int
}
