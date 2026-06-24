package logo

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"runtime"

	"github.com/dolmen-go/kittyimg"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

//go:embed alpha.png
var logo []byte

func Image() image.Image {
	img, _, _ := image.Decode(bytes.NewReader(logo))
	return img
}

func Render() string {
	img, err := RenderImage(bytes.NewReader(logo))
	if err != nil {
		panic(err)
	}
	return img
}

func RenderWithSize(width, height int) string {
	img, err := RenderImageResized(bytes.NewReader(logo), width, height)
	if err != nil {
		panic(err)
	}
	return img
}

func RenderImage(r io.Reader) (string, error) {
	src, _, err := image.Decode(r)
	if err != nil {
		return "", err
	}
	return encode(fillTransparentBlack(src))
}

func RenderImageResized(r io.Reader, width, height int) (string, error) {
	src, _, err := image.Decode(r)
	if err != nil {
		return "", err
	}
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return encode(fillTransparentBlack(dst))
}

func RenderFetch(img image.Image, width, height int) string {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	logoImg := fillTransparentBlack(dst)

	textImg := renderText(logoImg.Bounds().Dy())
	combined := stitchSideBySide(logoImg, textImg)

	var imgBuf bytes.Buffer
	png.Encode(&imgBuf, combined)

	var buf bytes.Buffer
	kittyimg.Transcode(&buf, &imgBuf)
	return buf.String()
}

func encode(img image.Image) (string, error) {
	var imgBuf bytes.Buffer
	if err := png.Encode(&imgBuf, img); err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := kittyimg.Transcode(&buf, &imgBuf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

const (
	cellW = 8
	cellH = 16
	padX  = 12
	padY  = 8
)

func renderText(height int) image.Image {
	lines := runtimeInfo()

	lineHeight := cellH
	numLines := len(lines)
	textH := numLines*lineHeight + padY*2
	if textH < height {
		textH = height
	}

	maxTextW := 0
	for _, l := range lines {
		w := 0
		for _, seg := range l.segments {
			w += len(seg.text) * cellW
		}
		if w > maxTextW {
			maxTextW = w
		}
	}
	textW := maxTextW + padX*2

	img := image.NewRGBA(image.Rect(0, 0, textW, textH))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	face := basicfont.Face7x13
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.White),
		Face: face,
	}

	y := padY + face.Metrics().Ascent.Ceil()
	for _, l := range lines {
		totalW := 0
		for _, seg := range l.segments {
			totalW += len(seg.text)
		}
		if totalW == 0 {
			y += lineHeight
			continue
		}

		x := padX
		for _, segment := range l.segments {
			if segment.color != nil {
				d.Src = image.NewUniform(segment.color)
			} else {
				d.Src = image.NewUniform(color.White)
			}
			d.Dot = fixed.P(x, y)
			d.DrawString(segment.text)
			x += len(segment.text) * cellW
		}
		y += lineHeight
	}

	return img
}

type textSegment struct {
	text  string
	color color.Color
}

type textLine struct {
	segments []textSegment
	bold     bool
}

func runtimeInfo() []textLine {
	green := color.RGBA{0, 200, 80, 255}
	white := color.RGBA{220, 220, 220, 255}
	dim := color.RGBA{140, 140, 140, 255}
	red := color.RGBA{255, 80, 80, 255}
	yellow := color.RGBA{255, 200, 50, 255}
	cyan := color.RGBA{80, 200, 220, 255}

	return []textLine{
		{segments: []textSegment{{"logo", green}}, bold: true},
		{segments: []textSegment{{"", nil}}},
		{segments: []textSegment{
			{"os:   ", dim},
			{runtime.GOOS, white},
		}},
		{segments: []textSegment{
			{"arch: ", dim},
			{runtime.GOARCH, white},
		}},
		{segments: []textSegment{
			{"go:   ", dim},
			{runtime.Version(), white},
		}},
		{segments: []textSegment{{"", nil}}},
		{segments: []textSegment{
			{"\u2588\u2588", red},
			{"\u2588\u2588", yellow},
			{"\u2588\u2588", green},
			{"\u2588\u2588", cyan},
			{"\u2588\u2588", color.RGBA{80, 80, 255, 255}},
			{"\u2588\u2588", color.RGBA{200, 80, 200, 255}},
		}},
	}
}

func stitchSideBySide(left, right image.Image) image.Image {
	lb := left.Bounds()
	rb := right.Bounds()

	w := lb.Dx() + rb.Dx()
	h := lb.Dy()
	if rb.Dy() > h {
		h = rb.Dy()
	}

	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(dst, lb, left, lb.Min, draw.Src)

	rOffX := lb.Dx()
	draw.Draw(dst, rb.Add(image.Pt(rOffX, 0)), rb, rb.Min, draw.Src)

	return dst
}
