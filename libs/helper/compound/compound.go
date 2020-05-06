package CHelperCompound

type Compound struct {
	width     float64
	height    float64
	fontPath  string
	fontSize  int
	lineSpace float64
	Y         float64
}

func NewCompound() *Compound {
	this := new(Compound)
	return this
}
func (this *Compound) Init(width float64, height float64, fontPath string, fontSize int, lineSpace float64) {
	this.width = width
	this.height = height
	this.fontPath = fontPath
	this.fontSize = fontSize
	this.lineSpace = lineSpace
}

func (this *Compound) SetTitle(title string) {

}
