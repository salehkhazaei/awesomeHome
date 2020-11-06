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
	for {
		if len(s.imageChannelMap) == 0 {
			// no requests
			time.Sleep(500 * time.Millisecond)
			continue
		}

		fmt.Printf("sending frame for %d clients\n", len(s.imageChannelMap))

		imgBuf, err := s.ReadImageFromCamera(cam, format, w, h)
		if err != nil {
			fmt.Printf("error in reading image from webcam: %v\n", err)
			continue
		}

		for _, channel := range s.imageChannelMap {
			channel <- imgBuf
		}
	}
}

func (s *WebcamService) HandleHttp(w http.ResponseWriter, r *http.Request) {
	log.Println("connect from", r.RemoteAddr, r.URL)

	rid := rand.Int63()
	imageChan := make(chan *bytes.Buffer)

	s.imageChannelMap[rid] = imageChan
	<-imageChan

	defer func() {
		delete(s.imageChannelMap, rid)
	}()

	const boundary = `frame`
	w.Header().Set("Content-Type", `multipart/x-mixed-replace;boundary=`+boundary)
	multipartWriter := multipart.NewWriter(w)
	multipartWriter.SetBoundary(boundary)
	for {
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
}
