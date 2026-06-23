package compilesvg

import (
	"encoding/binary"
	"testing"
)

func TestImageDimensions_png(t *testing.T) {
	w, h := ImageDimensions(tinyPNG, "image/png")
	if w != 1 || h != 1 {
		t.Fatalf("png dims: got %dx%d want 1x1", w, h)
	}
}

func TestImageDimensions_webpVP8X(t *testing.T) {
	body := buildWebPVP8X(720, 12800)
	w, h := ImageDimensions(body, "image/webp")
	if w != 720 || h != 12800 {
		t.Fatalf("webp dims: got %dx%d want 720x12800", w, h)
	}
}

func TestImageDimensions_webpSniff(t *testing.T) {
	body := buildWebPVP8X(800, 6000)
	w, h := ImageDimensions(body, "")
	if w != 800 || h != 6000 {
		t.Fatalf("webp sniff dims: got %dx%d want 800x6000", w, h)
	}
}

func TestImageDimensions_avif_ispe(t *testing.T) {
	body := buildMinimalAVIF(1080, 15000)
	w, h := ImageDimensions(body, "image/avif")
	if w != 1080 || h != 15000 {
		t.Fatalf("avif dims: got %dx%d want 1080x15000", w, h)
	}
}

func TestImageDimensions_unknownDefaults(t *testing.T) {
	w, h := ImageDimensions([]byte("not an image"), "text/plain")
	if w != defaultWidth || h != defaultHeight {
		t.Fatalf("defaults: got %dx%d want %dx%d", w, h, defaultWidth, defaultHeight)
	}
}

func buildWebPVP8X(width, height int) []byte {
	// VP8X extended header: canvas size is stored as width-1, height-1 (24-bit LE).
	vp8x := make([]byte, 8+10)
	copy(vp8x[0:4], "VP8X")
	binary.LittleEndian.PutUint32(vp8x[4:8], 10)
	vp8x[8] = 0x02 // alpha + exif flags optional
	wm1 := width - 1
	hm1 := height - 1
	vp8x[12] = byte(wm1)
	vp8x[13] = byte(wm1 >> 8)
	vp8x[14] = byte(wm1 >> 16)
	vp8x[15] = byte(hm1)
	vp8x[16] = byte(hm1 >> 8)
	vp8x[17] = byte(hm1 >> 16)

	riffSize := 4 + len(vp8x) // "WEBP" + chunks
	out := make([]byte, 8+4+riffSize)
	copy(out[0:4], "RIFF")
	binary.LittleEndian.PutUint32(out[4:8], uint32(riffSize))
	copy(out[8:12], "WEBP")
	copy(out[12:], vp8x)
	return out
}

func buildMinimalAVIF(width, height uint32) []byte {
	ispePayload := make([]byte, 12)
	binary.BigEndian.PutUint32(ispePayload[4:8], width)
	binary.BigEndian.PutUint32(ispePayload[8:12], height)
	ispe := wrapBox("ispe", ispePayload)

	ipco := wrapBox("ipco", ispe)
	iprpPayload := make([]byte, 4) // FullBox version+flags
	iprpPayload = append(iprpPayload, ipco...)
	iprp := wrapBox("iprp", iprpPayload)
	metaPayload := make([]byte, 4)
	metaPayload = append(metaPayload, iprp...)
	meta := wrapBox("meta", metaPayload)

	ftypPayload := append([]byte("avif"), 0, 0, 0, 0)
	ftypPayload = append(ftypPayload, []byte("avif")...)
	ftyp := wrapBox("ftyp", ftypPayload)

	return append(ftyp, meta...)
}

func wrapBox(boxType string, payload []byte) []byte {
	size := 8 + len(payload)
	buf := make([]byte, size)
	binary.BigEndian.PutUint32(buf[0:4], uint32(size))
	copy(buf[4:8], boxType)
	copy(buf[8:], payload)
	return buf
}

