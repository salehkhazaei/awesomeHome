package webcam

import (
	"bytes"
	"fmt"
	"github.com/blackjack/webcam"
	"image"
	"image/jpeg"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"sort"
	"strconv"
	"time"
)

type WebcamService struct {
	imageChan chan *bytes.Buffer
}

func NewWebcamService() *WebcamService {
	return &WebcamService{
		imageChan: make(chan *bytes.Buffer),
	}
}

func (s *WebcamService) Init(dev *string, fmtstr *string, szstr *string, fps *bool) {
	cam, err := webcam.Open(*dev)
	if err != nil {
		panic(err.Error())
	}
	defer cam.Close()

	// select pixel format
	formatDesc := cam.GetSupportedFormats()
	var format webcam.PixelFormat
FMT:
	for f, s := range formatDesc {
		if *fmtstr == "" {
			if supportedFormats[f] {
				format = f
				break FMT
			}

		} else if *fmtstr == s {
			if !supportedFormats[f] {
				log.Println(formatDesc[f], "format is not supported, exiting")
				return
			}
			format = f
			break
		}
	}

	if format == 0 {
		log.Println("No format found, exiting")
		return
	}

	// select frame size
	frames := FrameSizes(cam.GetSupportedFrameSizes(format))
	sort.Sort(frames)

	var size *webcam.FrameSize
	if *szstr == "" {
		size = &frames[len(frames)-1]
	} else {
		for _, f := range frames {
			if *szstr == f.GetString() {
				size = &f
			}
		}
	}
	if size == nil {
		log.Println("No matching frame size, exiting")
		return
	}

	f, w, h, err := cam.SetImageFormat(format, uint32(size.MaxWidth), uint32(size.MaxHeight))
	if err != nil {
		log.Println("SetImageFormat return error", err)
		return
	}

	// start streaming
	err = cam.StartStreaming()
	if err != nil {
		log.Println(err)
		return
	}

	var (
		fi   = make(chan []byte)
		back = make(chan struct{})
	)

	go encodeToImage(cam, back, fi, s.imageChan, w, h, f)

	timeout := uint32(5) //5 seconds
	start := time.Now()
	var fr time.Duration

	for {
		err = cam.WaitForFrame(timeout)
		if err != nil {
			log.Println(err)
			return
		}

		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			log.Println(err)
			continue
		default:
			log.Println(err)
			return
		}

		frame, err := cam.ReadFrame()
		if err != nil {
			log.Println(err)
			return
		}
		if len(frame) != 0 {

			// print framerate info every 10 seconds
			fr++
			if *fps {
				if d := time.Since(start); d > time.Second*10 {
					fmt.Println(float64(fr)/(float64(d)/float64(time.Second)), "fps")
					start = time.Now()
					fr = 0
				}
			}

			select {
			case fi <- frame:
				<-back
			default:
			}
		}
	}
}

func (s *WebcamService) HandleHttp(w http.ResponseWriter, r *http.Request) {
	log.Println("connect from", r.RemoteAddr, r.URL)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	//remove stale image
	<-s.imageChan
	const boundary = `frame`
	w.Header().Set("Content-Type", `multipart/x-mixed-replace;boundary=`+boundary)
	multipartWriter := multipart.NewWriter(w)
	multipartWriter.SetBoundary(boundary)
	for {
		img := <-s.imageChan
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

func encodeToImage(wc *webcam.Webcam, back chan struct{}, fi chan []byte, li chan *bytes.Buffer, w, h uint32, format webcam.PixelFormat) {

	var (
		frame []byte
		img   image.Image
	)
	for {
		bframe := <-fi
		// copy frame
		if len(frame) < len(bframe) {
			frame = make([]byte, len(bframe))
		}
		copy(frame, bframe)
		back <- struct{}{}

		switch format {
		case V4L2_PIX_FMT_YUYV:
			yuyv := image.NewYCbCr(image.Rect(0, 0, int(w), int(h)), image.YCbCrSubsampleRatio422)
			for i := range yuyv.Cb {
				ii := i * 4
				yuyv.Y[i*2] = frame[ii]
				yuyv.Y[i*2+1] = frame[ii+2]
				yuyv.Cb[i] = frame[ii+1]
				yuyv.Cr[i] = frame[ii+3]

			}
			img = yuyv
		default:
			log.Fatal("invalid format ?")
		}
		//convert to jpeg
		buf := &bytes.Buffer{}
		if err := jpeg.Encode(buf, img, nil); err != nil {
			log.Fatal(err)
			return
		}

		const N = 50
		// broadcast image up to N ready clients
		nn := 0
	FOR:
		for ; nn < N; nn++ {
			select {
			case li <- buf:
			default:
				break FOR
			}
		}
		if nn == 0 {
			li <- buf
		}

	}
}
