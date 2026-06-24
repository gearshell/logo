package main

import (
	"fmt"
	"image"
	"image/color"
	"runtime"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/NimbleMarkets/ntcharts/v2/picture"
)

type fetchModel struct {
	pic    picture.Model
	info   []infoLine
	ready  bool
	width  int
	height int
	imgW   int
	imgH   int
}

type infoLine struct {
	segments []segment
}

type segment struct {
	text  string
	color color.Color
}

var (
	white   = color.RGBA{220, 220, 220, 255}
	dim     = color.RGBA{120, 120, 120, 255}
	green   = color.RGBA{0, 200, 80, 255}
	red     = color.RGBA{255, 80, 80, 255}
	yellow  = color.RGBA{255, 200, 50, 255}
	cyan    = color.RGBA{80, 200, 220, 255}
	magenta = color.RGBA{200, 80, 200, 255}
	blue    = color.RGBA{80, 80, 255, 255}
)

func newFetchModel(img image.Image, kittyID int) fetchModel {
	picture.ForceKittyCapability(picture.KittyCapabilitySupported)

	m := fetchModel{
		info: []infoLine{
			{segments: []segment{{"logo", green}}},
			{segments: []segment{}},
			{segments: []segment{
				{"os:   ", dim},
				{runtime.GOOS, white},
			}},
			{segments: []segment{
				{"arch: ", dim},
				{runtime.GOARCH, white},
			}},
			{segments: []segment{
				{"go:   ", dim},
				{runtime.Version(), white},
			}},
			{segments: []segment{}},
			{segments: []segment{
				{"\u2588\u2588", red},
				{"\u2588\u2588", yellow},
				{"\u2588\u2588", green},
				{"\u2588\u2588", cyan},
				{"\u2588\u2588", blue},
				{"\u2588\u2588", magenta},
			}},
		},
	}

	m.pic = picture.NewWithConfig(picture.Config{
		KittyID:    kittyID,
		Background: color.Black,
		Fit:        picture.FitContain,
	})
	if img != nil {
		bounds := img.Bounds()
		m.imgW = bounds.Dx()
		m.imgH = bounds.Dy()
		m.pic.SetImage(img)
	}
	return m
}

func (m fetchModel) Init() tea.Cmd {
	return m.pic.Init()
}

type quitMsg struct{}

func (m fetchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case quitMsg:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		picCols := 30
		picRows := picCols * m.imgH / m.imgW / 2
		if picRows < 1 {
			picRows = 1
		}
		cmds := []tea.Cmd{}
		if c := m.pic.SetSize(picCols, picRows); c != nil {
			cmds = append(cmds, c)
		}
		if m.pic.Mode() != picture.PictureKitty {
			if c := m.pic.Toggle(); c != nil {
				cmds = append(cmds, c)
			}
		}
		cmds = append(cmds, tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg {
			return quitMsg{}
		}))
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	if cmd := m.pic.Update(msg); cmd != nil {
		return m, cmd
	}
	return m, nil
}

func (m fetchModel) View() tea.View {
	if !m.ready || m.width == 0 {
		return tea.NewView("loading...")
	}

	picView := m.pic.View().Content

	picLines := strings.Split(picView, "\n")
	picWidth := 0
	for _, l := range picLines {
		if w := stringWidth(l); w > picWidth {
			picWidth = w
		}
	}

	var textLines []string
	for _, info := range m.info {
		if len(info.segments) == 0 {
			textLines = append(textLines, "")
			continue
		}
		var sb strings.Builder
		for _, seg := range info.segments {
			r, g, b, _ := seg.color.RGBA()
			sb.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m",
				r>>8, g>>8, b>>8, seg.text))
		}
		textLines = append(textLines, sb.String())
	}

	maxPicLines := len(picLines)
	maxTextLines := len(textLines)
	rows := maxPicLines
	if maxTextLines > rows {
		rows = maxTextLines
	}

	for len(picLines) < rows {
		picLines = append(picLines, "")
	}
	for len(textLines) < rows {
		textLines = append(textLines, "")
	}

	var out strings.Builder
	for i := 0; i < rows; i++ {
		pl := picLines[i]
		tl := textLines[i]
		out.WriteString(pl)
		out.WriteString("  ")
		out.WriteString(tl)
		if i < rows-1 {
			out.WriteByte('\n')
		}
	}

	return tea.NewView(out.String())
}

func stringWidth(s string) int {
	w := 0
	inEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z') {
				inEscape = false
			}
			continue
		}
		w++
	}
	return w
}
