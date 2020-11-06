package webcam

import (
	"bytes"
	"errors"
	"github.com/blackjack/webcam"
	"image"
	"image/jpeg"
)

func createYUYVImage(frame []byte, w, h uint32) image.Image {
	yuyv := image.NewYCbCr(image.Rect(0, 0, int(w), int(h)), image.YCbCrSubsampleRatio422)
	for i := range yuyv.Cb {
		ii := i * 4
		yuyv.Y[i*2] = frame[ii]
		yuyv.Y[i*2+1] = frame[ii+2]
		yuyv.Cb[i] = frame[ii+1]
		yuyv.Cr[i] = frame[ii+3]

	}
	return yuyv
}

func encodeToImageFrame(frame []byte, w, h uint32, format webcam.PixelFormat) (*bytes.Buffer, error) {
	var img image.Image
	switch format {
	case V4L2_PIX_FMT_YUYV:
		img = createYUYVImage(frame, w, h)
	default:
		return nil, errors.New("invalid format")
	}

	//convert to jpeg
	buf := &bytes.Buffer{}
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return nil, err
	}

	return buf, nil
}
