package logo

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"

	"github.com/dolmen-go/kittyimg"
	"golang.org/x/image/draw"
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
