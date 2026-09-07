package layout

type Packer func()

type Atlas struct {
	Imgs []Img
}

type Img struct {
	Name   string
	Width  int
	Height int
	X      int
	Y      int
}
