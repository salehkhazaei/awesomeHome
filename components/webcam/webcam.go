package webcam

import (
	"bytes"
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
		if time.Now().After(startTime.Add(60 * time.Second)) {
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
