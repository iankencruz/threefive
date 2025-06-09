package imaging

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

func ConvertAndResizeToWebP(buf []byte) ([]byte, []byte, []byte, error) {
	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, nil, nil, err
	}

	// Resize versions
	thumbImg := imaging.Resize(img, 300, 0, imaging.Lanczos)
	mediumImg := imaging.Resize(img, 1200, 0, imaging.Lanczos)

	var originalBuf, thumbBuf, mediumBuf bytes.Buffer
	opt := &webp.Options{Lossless: false, Quality: 85}

	if err := webp.Encode(&originalBuf, img, opt); err != nil {
		return nil, nil, nil, err
	}
	if err := webp.Encode(&thumbBuf, thumbImg, opt); err != nil {
		return nil, nil, nil, err
	}
	if err := webp.Encode(&mediumBuf, mediumImg, opt); err != nil {
		return nil, nil, nil, err
	}

	return originalBuf.Bytes(), thumbBuf.Bytes(), mediumBuf.Bytes(), nil
}
