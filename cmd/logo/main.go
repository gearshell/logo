package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/gearshell/logo"
)

func main() {
	size := flag.String("size", "160x160", "resize image to given dimensions (e.g. 160x160)")
	flag.Parse()

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
