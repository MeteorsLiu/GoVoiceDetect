package govoicedetect

import (
	"io"
	"log"
	"os"

	webrtcvad "github.com/baabaaox/go-webrtcvad"
)

const (
	FRAME_WIDTH            float64 = 4096.0
	MAX_REGION_SIZE        float64 = 10.0
	MIN_REGION_SIZE        float64 = 0.5
	VAD_FRAME_DURATION_SEC float64 = 0.02
	DEFAULT_RATE                   = 16000
	MAX_CONCURRENT                 = 10
	VAD_FRAME_DURATION             = 20
	VAD_MODE                       = 0
	VERBOSE                        = true
)

type Region struct {
	Start float64
	End   float64
}

type VAD struct {
	file *os.File
	tmp  string
}

func NewRegion(start, end float64) Region {
	return Region{Start: start, End: end}
}

func verbose(c ...any) {
	if VERBOSE {
		log.Println(c...)
	}
}

func NewVad(filename string) (*VAD, error) {
	tmp := os.TempDir()
	verbose("preprocessing the audio file")
	f, err := doWarmUp(tmp, filename)
	if err != nil {
		os.Remove(tmp)
		return nil, err
	}
	verbose("preprocessing done")
	return &VAD{f, tmp}, nil
}

func (v *VAD) prefix(filename string) string {
	return v.tmp + "/" + filename
}

func (v *VAD) Detect() []Region {
	defer os.Remove(v.tmp)
	defer os.Remove(v.file.Name())
	WIDTH := DEFAULT_RATE / 1000 * VAD_FRAME_DURATION * 16 / 8
	frameBuffer := make([]byte, WIDTH)
	frameSize := DEFAULT_RATE / 1000 * VAD_FRAME_DURATION
	chunkDuration := (float64(WIDTH) / float64(DEFAULT_RATE)) / 2.0

	vadInst := webrtcvad.Create()
	defer webrtcvad.Free(vadInst)
	webrtcvad.Init(vadInst)

	err := webrtcvad.SetMode(vadInst, VAD_MODE)
	if err != nil {
		log.Fatal(err)
	}
	var region_start float64
	var elapsed_time float64
	var last_time float64
	var window_size float64
	var active_size float64
	var triggered bool
	var regions []Region

	for {
		n, err := v.file.Read(frameBuffer)
		if n > 0 {
			frameActive, _ := webrtcvad.Process(vadInst, DEFAULT_RATE, frameBuffer, frameSize)
			window_size++
			elapsed_time += chunkDuration
			if frameActive {
				active_size++
			}

			last_time = elapsed_time - region_start
			if triggered {
				if last_time >= MIN_REGION_SIZE {
					regions = append(regions, NewRegion(region_start, elapsed_time))
					region_start = elapsed_time
				}
				if active_size < 10 {
					region_start = elapsed_time
					triggered = false
					window_size = 0
					active_size = 0
				}
			} else {
				if active_size > 10 {
					region_start = elapsed_time
					triggered = true
					window_size = 0
					active_size = 0
				}
			}

		}
		if err != nil {
			if err != io.EOF {
				verbose(err)
			}
			break
		}

	}
	return regions
}
