package window

import (
	"lamp/config"
	"log"
	"os"

	xfont "golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const FontSize = 15.0

var (
	CharW, CharH int
	Face         xfont.Face
)

func InitFont() {
	fontPaths := []string{
		"/System/Library/Fonts/Menlo.ttc",
		"/Library/Fonts/Courier New.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
	}
	var fontData []byte
	for _, p := range fontPaths {
		data, err := os.ReadFile(p)
		if err == nil {
			fontData = data
			break
		}
	}
	if fontData == nil {
		log.Fatal("no monospace font found")
	}

	collection, err := opentype.ParseCollection(fontData)
	var f *opentype.Font
	if err != nil {
		f, err = opentype.Parse(fontData)
		if err != nil {
			log.Fatal("failed to parse font:", err)
		}
	} else {
		f, err = collection.Font(0)
		if err != nil {
			log.Fatal("failed to get font from collection:", err)
		}
	}

	const dpi = 144
	Face, err = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    FontSize,
		DPI:     dpi,
		Hinting: xfont.HintingFull,
	})
	if err != nil {
		log.Fatal("failed to create font face:", err)
	}

	metrics := Face.Metrics()
	CharH = (metrics.Ascent + metrics.Descent).Ceil()
	advance, ok := Face.GlyphAdvance('M')
	if !ok || advance == 0 {
		advance, _ = Face.GlyphAdvance('A')
	}
	CharW = advance.Ceil()
	if CharW == 0 {
		CharW = CharH / 2
	}

	log.Printf("Font metrics: CharW=%d CharH=%d", CharW, CharH)
	_ = config.Cols // ensure import used
}
