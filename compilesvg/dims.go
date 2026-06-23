package compilesvg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"

	"golang.org/x/image/webp"
)

const (
	defaultWidth  = 1080
	defaultHeight = 1920
)

// ImageDimensions returns width and height from image bytes, or defaults on failure.
func ImageDimensions(body []byte, mime string) (w, h int) {
	if len(body) == 0 {
		return defaultWidth, defaultHeight
	}
	if w, h, ok := dimensionsFromStandard(body); ok {
		return w, h
	}
	m := strings.ToLower(strings.TrimSpace(mime))
	if m == "image/webp" || isWebP(body) {
		if w, h, ok := webpDimensions(body); ok {
			return w, h
		}
	}
	if m == "image/avif" || isAVIF(body) {
		if w, h, ok := avifDimensions(body); ok {
			return w, h
		}
	}
	return defaultWidth, defaultHeight
}

func dimensionsFromStandard(body []byte) (w, h int, ok bool) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(body))
	if err != nil || cfg.Width <= 0 || cfg.Height <= 0 {
		return 0, 0, false
	}
	return cfg.Width, cfg.Height, true
}

func webpDimensions(body []byte) (w, h int, ok bool) {
	cfg, err := webp.DecodeConfig(bytes.NewReader(body))
	if err == nil && cfg.Width > 0 && cfg.Height > 0 {
		return cfg.Width, cfg.Height, true
	}
	return webpDimensionsFromChunks(body)
}

// webpDimensionsFromChunks reads VP8X (preferred) or VP8L chunk headers without full decode.
func webpDimensionsFromChunks(body []byte) (w, h int, ok bool) {
	if !isWebP(body) {
		return 0, 0, false
	}
	off := 12 // after RIFF header
	for off+8 <= len(body) {
		if off+8 > len(body) {
			break
		}
		chunk := string(body[off : off+4])
		size := int(binary.LittleEndian.Uint32(body[off+4 : off+8]))
		dataStart := off + 8
		dataEnd := dataStart + size
		if size < 0 || dataEnd > len(body) {
			break
		}
		switch chunk {
		case "VP8X":
			if size >= 10 && dataEnd <= len(body) {
				w := int(body[dataStart+4]) | int(body[dataStart+5])<<8 | int(body[dataStart+6])<<16
				h := int(body[dataStart+7]) | int(body[dataStart+8])<<8 | int(body[dataStart+9])<<16
				if w > 0 && h > 0 {
					return w + 1, h + 1, true
				}
			}
		case "VP8L":
			if size >= 5 && dataEnd <= len(body) {
				b0 := body[dataStart]
				b1 := body[dataStart+1]
				b2 := body[dataStart+2]
				b3 := body[dataStart+3]
				w := 1 + int(b0) + int(b1&0x3F)<<8
				h := 1 + int(b1>>6) + int(b2&0xF)<<4 + int(b3&0xFC)<<2
				if w > 0 && h > 0 {
					return w, h, true
				}
			}
		}
		pad := size & 1
		off = dataEnd + pad
	}
	return 0, 0, false
}

func isWebP(body []byte) bool {
	return len(body) >= 12 &&
		body[0] == 'R' && body[1] == 'I' && body[2] == 'F' && body[3] == 'F' &&
		body[8] == 'W' && body[9] == 'E' && body[10] == 'B' && body[11] == 'P'
}

func isAVIF(body []byte) bool {
	if len(body) < 12 {
		return false
	}
	// ISO BMFF: size + "ftyp" + major brand
	if string(body[4:8]) != "ftyp" {
		return false
	}
	brand := string(body[8:12])
	return brand == "avif" || brand == "avis" || brand == "mif1" || brand == "MA1A"
}

// avifDimensions reads width/height from the ispe (Image Spatial Extents) box.
func avifDimensions(body []byte) (w, h int, ok bool) {
	var found bool
	walkBoxes(body, func(boxType string, payload []byte) bool {
		if boxType != "ispe" || len(payload) < 12 {
			return true
		}
		width := binary.BigEndian.Uint32(payload[4:8])
		height := binary.BigEndian.Uint32(payload[8:12])
		if width > 0 && height > 0 {
			w, h = int(width), int(height)
			found = true
			return false
		}
		return true
	})
	return w, h, found
}

type boxWalker func(boxType string, payload []byte) (continueWalk bool)

func walkBoxes(data []byte, visit boxWalker) {
	var walk func([]byte) bool
	walk = func(buf []byte) bool {
		for off := 0; off+8 <= len(buf); {
			size32 := int(binary.BigEndian.Uint32(buf[off : off+4]))
			boxType := string(buf[off+4 : off+8])
			header := 8
			size := size32
			if size32 == 0 {
				break
			}
			if size32 == 1 {
				if off+16 > len(buf) {
					break
				}
				size = int(binary.BigEndian.Uint64(buf[off+8 : off+16]))
				header = 16
			}
			if size < header || off+size > len(buf) {
				break
			}
			payloadStart := off + header
			payloadEnd := off + size
			payload := buf[payloadStart:payloadEnd]

			cont := visit(boxType, payload)
			if !cont {
				return false
			}

			switch boxType {
			case "meta", "iprp", "moov", "trak", "mdia", "minf", "stbl", "stsd":
				child := payload
				if len(payload) >= 4 {
					child = payload[4:] // FullBox version+flags
				}
				if !walk(child) {
					return false
				}
			case "ipco":
				if !walk(payload) {
					return false
				}
			}

			off += size
		}
		return true
	}
	walk(data)
}

// FormatInt returns decimal string for SVG attributes.
func FormatInt(n int) string {
	return fmt.Sprintf("%d", n)
}

// IsImageMIME reports whether mime is an image type.
func IsImageMIME(mime string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(mime)), "image/")
}
