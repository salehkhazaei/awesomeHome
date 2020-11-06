package webcam

import (
	"bytes"
	"errors"
	"github.com/blackjack/webcam"
)

func PrepareCamera(cam *webcam.Webcam, formatString string, sizeString string) (webcam.PixelFormat, uint32, uint32, error) {
	// select pixel format
	format := GetFormat(cam, formatString)
	if format == 0 {
		return 0, 0, 0, errors.New("no format found, exiting")
	}

	// select frame size
	size := GetFrameSize(cam, format, sizeString)
	if size == nil {
		return 0, 0, 0, errors.New("no matching frame size, exiting")
	}

	f, w, h, err := cam.SetImageFormat(format, uint32(size.MaxWidth), uint32(size.MaxHeight))
	if err != nil {
		return 0, 0, 0, err
	}

	// start streaming
	err = cam.StartStreaming()
	if err != nil {
		return 0, 0, 0, err
	}

	return f, w, h, nil
}

func (s *WebcamService) OpenCamera(dev string, fmtstr string, szstr string, fps bool) error {
	cam, err := webcam.Open(dev)
	if err != nil {
		return err
	}

	format, w, h, err := PrepareCamera(cam, fmtstr, szstr)
	if err != nil {
		return err
	}

	go s.HandleFrameRequests(cam, format, w, h)

	return nil
}

func (s *WebcamService) ReadImageFromCamera(cam *webcam.Webcam, format webcam.PixelFormat, w, h uint32) (*bytes.Buffer, error) {
	timeout := uint32(5) //5 seconds

	err := cam.WaitForFrame(timeout)
	if err != nil {
		return nil, err
	}

	frame, err := cam.ReadFrame()
	if err != nil {
		return nil, err
	}

	imgBuf, err := encodeToImageFrame(frame, w, h, format)
	if err != nil {
		return nil, err
	}

	return imgBuf, nil
}

func (s *WebcamService) StopCamera(cam *webcam.Webcam) error {
	return cam.StopStreaming()
}

func (s *WebcamService) ResumeCamera(cam *webcam.Webcam) error {
	return cam.StartStreaming()
}
