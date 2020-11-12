// +build linux

package webcam

import (
	"github.com/blackjack/webcam"
	"sort"
)

type FrameSizes []webcam.FrameSize

func (slice FrameSizes) Len() int {
	return len(slice)
}

//For sorting purposes
func (slice FrameSizes) Less(i, j int) bool {
	ls := slice[i].MaxWidth * slice[i].MaxHeight
	rs := slice[j].MaxWidth * slice[j].MaxHeight
	return ls < rs
}

//For sorting purposes
func (slice FrameSizes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func GetFrameSize(cam *webcam.Webcam, format webcam.PixelFormat, sizeString string) *webcam.FrameSize {
	frames := FrameSizes(cam.GetSupportedFrameSizes(format))
	sort.Sort(frames)

	var size *webcam.FrameSize
	if sizeString == "" {
		size = &frames[len(frames)-1]
	} else {
		for _, f := range frames {
			if sizeString == f.GetString() {
				size = &f
			}
		}
	}
	return size
}
