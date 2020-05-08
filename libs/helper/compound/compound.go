package CHelperCompound

import (
	"errors"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"unicode"
)

var (
	fontPath  string
	fontBytes []byte
	fontType  *truetype.Font
)

type Compound struct {
	width      float64
	height     float64
	dpi        int
	fontPath   string
	fontSize   float64
	lineSpace  float64
	obj        *image.RGBA
	fontBytes  []byte
	fontType   *truetype.Font
	once       sync.Once
	fontDrawer *font.Drawer
	Point      *fixed.Point26_6
}

func NewCompound() *Compound {
	this := new(Compound)
	return this
}

func SetCompound(fontpath string) error {
	var err error
	fontPath = fontpath
	fontBytes, err = ioutil.ReadFile(fontPath)
	if err != nil {
		return err
	}
	fontType, err = truetype.Parse(fontBytes)
	if err != nil {
		return err
	}
	return nil
}

//1200,800,"./src/simsun.ttf",18,72,1.2,1100
func (this *Compound) Init(width float64, height float64, fontpath string, fontSize float64, dpi int, lineSpace float64) error {
	this.width = width
	this.height = height
	this.fontSize = fontSize
	this.lineSpace = lineSpace
	this.dpi = dpi

	if fontPath == "" {
		if fontPath == "" {
			return errors.New("FontPath is Null")
		}
		this.fontPath = fontPath
		this.fontBytes = fontBytes
		this.fontType = fontType
	} else {
		this.fontPath = fontpath
		var err error
		this.fontBytes, err = ioutil.ReadFile(this.fontPath)
		if err != nil {
			return err
		}
		this.fontType, err = truetype.Parse(this.fontBytes)
		if err != nil {
			return err
		}
	}

	this.fontDrawer = &font.Drawer{
		Face: truetype.NewFace(this.fontType, &truetype.Options{
			Size: float64(this.fontSize),
		}),
	}

	this.obj = image.NewRGBA(image.Rect(0, 0, int(this.width), int(this.height)))
	this.Point = new(fixed.Point26_6)
	*(this.Point) = freetype.Pt(10, 10)

	draw.Draw(this.obj, this.obj.Bounds(), image.White, image.Point{}, draw.Src)
	return nil
}

func (this *Compound) splitOnSpace(x string) []string {
	var result []string
	pi := 0
	ps := false
	for i, c := range x {
		s := unicode.IsSpace(c)
		if s != ps && i > 0 {
			result = append(result, x[pi:i])
			pi = i
		}
		ps = s
	}
	result = append(result, x[pi:])
	return result
}

func (this *Compound) measureString(s string) (w, h float64) {
	a := this.fontDrawer.MeasureString(s)
	return float64(a >> 6), this.lineSpace
}

func (this *Compound) WordWrap(s string, width float64) []string {
	var result []string
	for _, line := range strings.Split(s, "\n") {
		if len(line) == 0 {
			result = append(result, "\n")
			continue
		}
		fields := this.splitOnSpace(line)

		if len(fields)%2 == 1 {
			fields = append(fields, "")
		}
		x := ""
		for i := 0; i < len(fields); i++ {
			runes := []rune(fields[i])
			for k, v := range runes {
				w, _ := this.measureString(x + string(runes[k]))
				if w > width {
					result = append(result, x)
					x = string(v)
				} else {
					x += string(v)
				}
			}
		}
		if x != "" {
			result = append(result, x)
		}
	}
	for i, line := range result {
		result[i] = strings.TrimSpace(line)
	}
	return result
}

func (this *Compound) AddTitle(title string, size float64) error {
	fg := image.Black

	c := freetype.NewContext()
	c.SetFont(this.fontType)
	c.SetFontSize(float64(size))
	c.SetClip(this.obj.Bounds())
	c.SetDst(this.obj)
	c.SetSrc(fg)

	w, _ := this.measureString(title)
	x := (this.width - w) / 3
	if x < 0 {
		x = 0
	}

	*this.Point = freetype.Pt(int(x), 10+int(c.PointToFixed(size)>>6))
	_, err := c.DrawString(title, freetype.Pt(int(x), 10+int(c.PointToFixed(size)>>6)))
	if err != nil {
		return err
	}

	this.Point.Y += c.PointToFixed(size * this.lineSpace)
	this.Point.Y += c.PointToFixed(size * this.lineSpace)
	return nil
}

func (this *Compound) AddBody(body string, x int, width float64) error {
	fg := image.Black

	c := freetype.NewContext()
	c.SetFont(this.fontType)
	c.SetFontSize(this.fontSize)
	c.SetClip(this.obj.Bounds())
	c.SetDst(this.obj)
	c.SetSrc(fg)
	Y := this.Point.Y
	*this.Point = freetype.Pt(x, 0)
	this.Point.Y = Y

	text := this.WordWrap(body, width)
	for _, v := range text {
		_, err := c.DrawString(v, *this.Point)
		if err != nil {
			return err
		}
		this.Point.Y += c.PointToFixed(this.fontSize * this.lineSpace)
	}
	return nil
}

// absolutePosition ture 绝对位置   false 相对位置
func (this *Compound) AddImage(imagePath string, imageWidth uint, imageHeight uint, x int, y int, absolutePosition bool, addHeight bool) error {
	var img image.Image
	if strings.Index(imagePath, "http") == 0 {
		resp, err := http.Get(imagePath)
		if err != nil {
			return err
		}
		img, _, err = image.Decode(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

	} else {
		f, err := os.Open(imagePath)
		if err != nil {
			return err
		}
		defer f.Close()
		img, _, err = image.Decode(f)
		if err != nil {
			return err
		}
	}

	img = resize.Resize(imageWidth, imageHeight, img, resize.Lanczos3)

	if absolutePosition {
		draw.Draw(this.obj, img.Bounds().Add(image.Pt(x, y)), img, img.Bounds().Min, draw.Src)
	} else {
		draw.Draw(this.obj, img.Bounds().Add(image.Pt(x, this.Point.Y.Floor()+y)), img, img.Bounds().Min, draw.Src)
	}

	if addHeight {
		c := freetype.NewContext()
		this.Point.Y += c.PointToFixed(float64(imageHeight))
	}

	return nil
}

func (this *Compound) AddMarginSpace(Y float64) {
	c := freetype.NewContext()
	this.Point.Y += c.PointToFixed(Y)
}

func (this *Compound) Save(writer io.Writer) error {
	return png.Encode(writer, this.obj)
}
