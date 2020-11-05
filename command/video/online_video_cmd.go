package video

import (
	"ir.skhf/awesomeHome/process"
	"net/http"
	"strings"
)

var supportedFormats []string = []string{"3gp", "asf", "avi", "dvr-ms", "flv", "mkv", "midi", "quicktime",
	"mp4", "ogg", "ogm", "wav", "mpeg", "aiff", "raw", "mxf", "vob", "rm", "vcd", "svcd", "dvb", "heif",
	"avif", "aac", "ac3", "alac", "amr", "dts", "xm", "flac", "it", "mace", "mod", "mp3", "opus", "pls",
	"qcp", "qdm2", "qdmc", "realaudio", "speex", "s3m", "tta", "vorbis", "wavpack", "wma"}

type OnlineVideoCmd struct {
	Url string
}

func (cmd *OnlineVideoCmd) Detect(commandStr string) bool {
	filename := fetchFilename(commandStr)
	if len(filename) > 0 {
		for _, extension := range supportedFormats {
			if len(filename) >= len(extension) {
				if filename[len(filename)-len(extension):] == extension {
					cmd.Url = commandStr
					return true
				}
			}
		}
	}

	contentType := fetchContentType(commandStr)
	if len(contentType) > 0 {
		videoStr := "video/"
		if len(contentType) >= len(videoStr) {
			if contentType[:len(videoStr)] == videoStr {
				cmd.Url = commandStr
				return true
			}
		}

		contentTypeArr := strings.Split(contentType, "/")
		lastPart := contentTypeArr[len(contentTypeArr)-1]
		for _, extension := range supportedFormats {
			if lastPart == extension {
				cmd.Url = commandStr
				return true
			}
		}
	}

	return false
}

func (cmd *OnlineVideoCmd) Run(processService *process.ProcessService) error {
	return processService.VLC(cmd.Url)
}

func fetchContentType(url string) string {
	res, err := http.Head(url)
	if err != nil {
		return ""
	}

	if res.StatusCode == 200 {
		contentType := ""
		for key, value := range res.Header {
			if strings.ToLower(key) == "content-type" {
				contentType = value[0]
			}
		}

		if len(contentType) <= 0 {
			return ""
		}

		return contentType
	}
	return ""
}

func fetchFilename(url string) string {
	res, err := http.Head(url)
	if err != nil {
		return ""
	}

	if res.StatusCode == 200 {
		contentDisposision := ""
		for key, value := range res.Header {
			if strings.ToLower(key) == "content-disposition" {
				contentDisposision = value[0]
			}
		}

		if len(contentDisposision) <= 0 {
			return ""
		}

		if !strings.Contains(contentDisposision, "filename") {
			return ""
		}

		cells := strings.Split(contentDisposision, ";")
		for _, cell := range cells {
			keyValue := strings.Split(cell, "=")
			if keyValue[0] == "filename" {
				return keyValue[1]
			}
		}
	}
	return ""
}
