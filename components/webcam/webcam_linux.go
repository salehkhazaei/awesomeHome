// +build linux

package webcam

func (s *WebcamService) HandleFrameRequests(cam *webcam.Webcam, format webcam.PixelFormat, w, h uint32) {
	defer cam.Close()

	var stopCount = 0
	for {
		if len(s.imageChannelMap) == 0 {
			// no requests
			time.Sleep(500 * time.Millisecond)
			stopCount++

			if stopCount == 4 {
				err := s.StopCamera(cam)
				if err != nil {
					fmt.Printf("error in stopping camera: %v\n", err)
					stopCount--
				}
			}
			continue
		}

		if stopCount > 0 {
			err := s.ResumeCamera(cam)
			if err != nil {
				fmt.Printf("error in resuming camera: %v\n", err)
				continue
			}

			stopCount = 0
		}

		fmt.Printf("sending frame for %d clients\n", len(s.imageChannelMap))

		imgBuf, err := s.ReadImageFromCamera(cam, format, w, h)
		if err != nil {
			fmt.Printf("error in reading image from webcam: %v\n", err)
			continue
		}

		for _, channel := range s.imageChannelMap {
			select {
			case channel <- imgBuf:
			default:
			}
		}

		time.Sleep(50 * time.Millisecond)
	}
}
