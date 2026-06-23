package webconnector

import (
	"bytes"
	"fmt"
	"os"
)

var allowedImageMIMEs = map[string]struct{}{
	"image/png":     {},
	"image/jpeg":    {},
	"image/webp":    {},
	"image/gif":     {},
	"image/avif":    {},
	"image/svg+xml": {},
}

type ImageAsset struct {
	Bytes []byte
	MIME  string
}

func ReadAllowedImage(path string) (ImageAsset, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return ImageAsset{}, err
	}
	mime := sniffImageMIME(body)
	if !IsAllowedImageMIME(mime) {
		return ImageAsset{}, &MIMEError{MIME: mime}
	}
	return ImageAsset{Bytes: body, MIME: mime}, nil
}

func IsAllowedImageMIME(mime string) bool {
	_, ok := allowedImageMIMEs[mime]
	return ok
}

type MIMEError struct {
	MIME string
}

func (e *MIMEError) Error() string {
	return fmt.Sprintf("unsupported media type: %s", e.MIME)
}

func sniffImageMIME(b []byte) string {
	sample := b
	const maxSniff = 512
	if len(sample) > maxSniff {
		sample = sample[:maxSniff]
	}
	trim := bytes.TrimSpace(sample)
	if len(trim) >= 5 && bytes.HasPrefix(trim, []byte("<?xml")) {
		return "image/svg+xml"
	}
	if len(trim) >= 4 && bytes.HasPrefix(trim, []byte("<svg")) {
		return "image/svg+xml"
	}
	switch {
	case len(sample) >= 3 && bytes.Equal(sample[:3], []byte{0xFF, 0xD8, 0xFF}):
		return "image/jpeg"
	case len(sample) >= 8 && bytes.Equal(sample[:8], []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}):
		return "image/png"
	case len(sample) >= 6 && string(sample[:6]) == "GIF87a":
		return "image/gif"
	case len(sample) >= 6 && string(sample[:6]) == "GIF89a":
		return "image/gif"
	case len(sample) >= 12 && string(sample[0:4]) == "RIFF" && string(sample[8:12]) == "WEBP":
		return "image/webp"
	case len(sample) >= 12 && string(sample[4:8]) == "ftyp":
		brand := string(sample[8:12])
		if brand == "avif" || brand == "avis" || brand == "mif1" {
			return "image/avif"
		}
	}
	return "application/octet-stream"
}
