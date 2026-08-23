package captor

import (
	"bytes"
	"image/png"
	"github.com/kbinani/screenshot"
)

func CaptureScreen() ([]byte, error) {
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err!= nil { return nil, err }
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err!= nil { return nil, err }
	return buf.Bytes(), nil
}
