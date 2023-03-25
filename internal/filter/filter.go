package filter

import (
	"bytes"
	"colorblinder/pkg/metrics"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Filter struct {
	ID               string
	RGBAOverlay      [4]int
	StartSecond      int
	IsPhotosensitive bool
}

func NewFilter(RGBAOverlay [4]int, StartSecond int, IsPhotosensitive bool) Filter {
	return Filter{
		ID:               uuid.NewString(),
		RGBAOverlay:      RGBAOverlay,
		StartSecond:      StartSecond,
		IsPhotosensitive: IsPhotosensitive,
	}
}

func (f Filter) StartProcess(url string) error {
	var errorBuf bytes.Buffer
	photosensitivityFilter := ""
	if f.IsPhotosensitive {
		photosensitivityFilter = ",photosensitivity=30:0.6:20"
	}
	if err := os.Mkdir("/tmp/"+f.ID, os.ModePerm); err != nil {
		return err
	}

	dataCmd := exec.Command("ffprobe", "-v", "error",
		"-select_streams", "v", "-show_entries", "stream=index",
		"-of", "compact=p=0:nk=1", url,
	)
	dataCountCommand := exec.Command("wc", "-w")
	r, w := io.Pipe()
	var dataBuf bytes.Buffer
	dataCmd.Stdout = w
	dataCountCommand.Stdin = r
	dataCountCommand.Stdout = &dataBuf
	dataCmd.Stderr = &errorBuf
	dataCountCommand.Stderr = &errorBuf
	if err := dataCmd.Start(); err != nil {
		return errors.WithMessage(err, errorBuf.String())
	}
	if err := dataCountCommand.Start(); err != nil {
		return errors.WithMessage(err, errorBuf.String())
	}
	dataCmd.Wait()
	w.Close()
	dataCountCommand.Wait()
	videoStreams, err := strconv.Atoi(dataBuf.String()[:dataBuf.Len()-1])
	if err != nil {
		return err
	}
	filters := []string{}
	videoStreams = 2
	for i := 0; i < videoStreams/2; i++ {
		filters = append(filters, "-filter_complex")
		filters = append(filters,
			fmt.Sprintf("[1:v][0:v:%d]scale2ref[1v][0v],[0v][1v]blend=shortest=1:all_mode=softlight:all_opacity=1", i)+photosensitivityFilter)
	}

	cmdSlice := []string{
		"-re", "-i", url, "-f", "lavfi", "-i",
		fmt.Sprintf("color=color=#%02x%02x%02x%02d", f.RGBAOverlay[0], f.RGBAOverlay[1], f.RGBAOverlay[2], f.RGBAOverlay[3]),
	}
	cmdSlice = append(cmdSlice, filters...)
	cmdSlice = append(cmdSlice,
		"-f", "dash", "-seg_duration", "1", "-adaptation_sets", "id=0,streams=v  id=1,streams=a",
		"/tmp/"+f.ID+"/file.mpd",
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
	}("/tmp/" + f.ID)
	metrics.ActiveFilterers.Dec()
	return nil
}
