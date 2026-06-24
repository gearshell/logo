package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"golang.org/x/image/draw"

	"github.com/gearshell/logo"
)

func main() {
	size := flag.String("size", "240x240", "resize image to given dimensions (e.g. 160x160)")
	fetch := flag.Bool("fetch", false, "display system info next to the image (interactive TUI)")
	flag.Parse()

	if *fetch {
		img := logo.Image()
		if *size != "" {
			w, h, err := parseSize(*size)
			if err != nil {
				fmt.Printf("invalid size %q: expected WxH (e.g. 160x160)\n", *size)
				flag.Usage()
				return
			}
			img = resizeImage(img, w, h)
		}
		m := newFetchModel(img, 4242)
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *size == "" {
		fmt.Println(logo.Render())
		return
	}

	w, h, err := parseSize(*size)
	if err != nil {
		fmt.Printf("invalid size %q: expected WxH (e.g. 160x160)\n", *size)
		flag.Usage()
		return
	}
	fmt.Println(logo.RenderWithSize(w, h))
}

func resizeImage(src image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

func parseSize(s string) (int, int, error) {
	parts := strings.SplitN(s, "x", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("missing 'x' separator")
	}
	w, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width: %v", err)
	}
	h, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height: %v", err)
	}
	if w <= 0 || h <= 0 {
		return 0, 0, fmt.Errorf("dimensions must be positive")
	}
	return w, h, nil
}
