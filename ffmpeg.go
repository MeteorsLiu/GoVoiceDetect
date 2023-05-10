package govoicedetect

import (
	"errors"
	"os"
	"os/exec"
)

func exists(cmd string) (string, bool) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return "", false
	}
	return path, true
}

func doWarmUp(tmpDir, filename string) (*os.File, error) {
	if cmd, ok := exists("ffmpeg"); ok {
		audio, err := os.CreateTemp(tmpDir, "*.pcm")
		if err != nil {
			return nil, err
		}
		if err := exec.Command(cmd, "-y", "-i", filename, "-f", "s16le", "-ar", "16000", "-ac", "1", "-acodec", "pcm_s16le", audio.Name()).Run(); err != nil {
			return nil, err
		}
		return audio, nil
	}
	return nil, errors.New("please install ffmpeg")
}
