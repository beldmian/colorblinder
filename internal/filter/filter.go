package filter

import (
	"bytes"
	"colorblinder/pkg/metrics"
	"log"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

func StartProcess(id string, url string) error {
	var errorBuf bytes.Buffer
	if err := os.Mkdir("/tmp/"+id, os.ModePerm); err != nil {
		return err
	}
	filters := make([]string, 0)
	filters = append(filters, "-filter_complex", "photosensitivity=30:0.6:20")

	cmdSlice := []string{
		"-re", "-i", url,
	}
	cmdSlice = append(cmdSlice, filters...)
	cmdSlice = append(cmdSlice,
		"-f", "dash", "-seg_duration", "1", "-adaptation_sets", "id=0,streams=v  id=1,streams=a",
		"-ldash", "1", "-preset", "veryfast",
		"/tmp/"+id+"/file.mpd",
	)

	cmd := exec.Command("ffmpeg", cmdSlice...)
	cmd.Stderr = &errorBuf
	metrics.ActiveFilterers.Inc()
	if err := cmd.Run(); err != nil {
		return errors.WithMessage(err, errorBuf.String())
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Println("cleanup error: ", err)
		}
	}("/tmp/" + id)
	metrics.ActiveFilterers.Dec()
	return nil
}
