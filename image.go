package logo

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/dolmen-go/kittyimg"
)

//go:embed 160.png
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

// RenderString reads an image from r and returns Kitty protocol output.
func RenderImage(r io.Reader) (string, error) {
	var buf bytes.Buffer
	if err := kittyimg.Transcode(&buf, r); err != nil {
		return "", err
	}
	return buf.String(), nil
}
