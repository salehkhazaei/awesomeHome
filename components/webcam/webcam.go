package webcam

import (
	"bytes"
	"fmt"
	"github.com/blackjack/webcam"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"time"
)

type WebcamService struct {
	requestChan     chan interface{}
	imageChannelMap map[int64]chan *bytes.Buffer
}

func NewWebcamService() *WebcamService {
	return &WebcamService{
		imageChannelMap: make(map[int64]chan *bytes.Buffer),
	}
}

func (s *WebcamService) Init(dev string, fmtstr string, szstr string, fps bool) error {
	return s.OpenCamera(dev, fmtstr, szstr, fps)
}

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
					continue
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
	}
}

func (s *WebcamService) HandleHttp(w http.ResponseWriter, r *http.Request) {
	log.Println("connect from", r.RemoteAddr, r.URL)
	startTime := time.Now()

	rid := rand.Int63()
	imageChan := make(chan *bytes.Buffer)
	s.imageChannelMap[rid] = imageChan

	log.Printf("image channel added %d", rid)
	<-imageChan

	defer func() {
		log.Printf("image channel removed %d", rid)
		delete(s.imageChannelMap, rid)
	}()

	const boundary = `frame`
	w.Header().Set("Content-Type", `multipart/x-mixed-replace;boundary=`+boundary)
	multipartWriter := multipart.NewWriter(w)
	multipartWriter.SetBoundary(boundary)
	for {
		if time.Now().After(startTime.Add(30 * time.Second)) {
			break
		}

		img := <-imageChan
		image := img.Bytes()
		iw, err := multipartWriter.CreatePart(textproto.MIMEHeader{
			"Content-type":   []string{"image/jpeg"},
			"Content-length": []string{strconv.Itoa(len(image))},
		})
		if err != nil {
			log.Println(err)
			return
		}
		_, err = iw.Write(image)
		if err != nil {
			log.Println(err)
			return
		}
	}

	multipartWriter.Close()
}
