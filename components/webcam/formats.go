package webcam

import "github.com/blackjack/webcam"

const (
	V4L2_PIX_FMT_PJPG = 0x47504A50
	V4L2_PIX_FMT_YUYV = 0x56595559
)

var supportedFormats = map[webcam.PixelFormat]bool{
	V4L2_PIX_FMT_PJPG: true,
	V4L2_PIX_FMT_YUYV: true,
}

func GetFormat(cam *webcam.Webcam, formatString string) webcam.PixelFormat {
	formatDesc := cam.GetSupportedFormats()
	var format webcam.PixelFormat

	for f, s := range formatDesc {
		if formatString == "" {
			if supportedFormats[f] {
				format = f
				break
			}

		} else if formatString == s {
			if !supportedFormats[f] {
				return 0
			}

			format = f
			break
		}
	}

	return format
}
