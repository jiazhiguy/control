package media

import (
	"log"
	"os"
	"time"
	// "context"
	"net/http"
	"io"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)
type Option struct  {
	LoopNum int
	Fast  float64
	Pause  bool
	Volume float64
}
// type audioPanel struct {
// 	sampleRate beep.SampleRate
// 	streamer   beep.StreamSeeker
// 	ctrl       *beep.Ctrl
// 	resampler  *beep.Resampler
// 	volume     *effects.Volume
// }

// func newAudioPanel(sampleRate beep.SampleRate, streamer beep.StreamSeeker) *audioPanel {
// 	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
// 	resampler := beep.ResampleRatio(4, 1, ctrl)
// 	volume := &effects.Volume{Streamer: resampler, Base: 2}
// 	return &audioPanel{sampleRate, streamer, ctrl, resampler, volume}
// }

// func (ap *audioPanel) play() {
// 	speaker.Play(ap.volume)
// }

func PlayMp3(path string, option  Option ,update chan Option) error {
	var f io.ReadCloser
	var err error
	var isLocalFile =true
	if !strings.HasPrefix(path,"http://") {	
		if f, err = os.Open(path) ;err !=nil {
			return err
		}
	}else {
		isLocalFile=false
	}
	// defer f.Close()
	log.Println(isLocalFile)
	if isLocalFile {
		streamer, format, err := mp3.Decode(f)
		if err != nil {
			log.Println(err)
			return 	err
		}
		defer streamer.Close()
		done := make(chan bool)
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		loop := beep.Loop(option.LoopNum, streamer)	
		fast := beep.ResampleRatio(4, option.Fast, loop)
		ctrl := &beep.Ctrl{Streamer: fast, Paused: option.Pause}

		volume := &effects.Volume{
			Streamer: ctrl,
			Base:     2,
			Volume:   option.Volume,
			Silent:   false,
		}
		speaker.Play(beep.Seq(volume, beep.Callback(func() {
			done <- true
		})))
		for {
			select {
			case <-done:
				return nil
			case param:=<-update: {
				log.Println(param)
				speaker.Clear()
				loop := beep.Loop(param.LoopNum, streamer)	
				fast := beep.ResampleRatio(4, param.Fast, loop)
				ctrl := &beep.Ctrl{Streamer: fast, Paused: param.Pause}
				volume := &effects.Volume{
					Streamer: ctrl,
					Base:     2,
					Volume:   param.Volume,
					Silent:   false,
				}
				speaker.Play(beep.Seq(volume, beep.Callback(func() {
					done <- true
				})))
			}
			// case <-time.After(time.Second):
			// 	speaker.Lock()
			// 	log.Println(format.SampleRate.D(streamer.Position()).Round(time.Second))
			// 	speaker.Unlock()
				
			}
		}
	}else {
		i:=1
		for i<=option.LoopNum {
			log.Println(i)
			resp, err := http.Get(path)
			if err != nil {
				return err
			}
			f =resp.Body
			defer f.Close()
			streamer, format, err := mp3.Decode(f)
			if err != nil {
				log.Println(err)
				return 	err
			}
			defer streamer.Close()
			done := make(chan bool)
			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
			loop := beep.Loop(1, streamer)	
			fast := beep.ResampleRatio(4, option.Fast, loop)
			ctrl := &beep.Ctrl{Streamer: fast, Paused: option.Pause}

			volume := &effects.Volume{
				Streamer: ctrl,
				Base:     2,
				Volume:   option.Volume,
				Silent:   false,
			}
			speaker.Play(beep.Seq(volume, beep.Callback(func() {
				done <- true
			})))
	Loop:
			for {
				select {
				// case <-ctx.Done():
				// 	speaker.Clear()
				// 	log.Println("cancel")
				// 	return nil
				case <-done:
					// log.Println("done")
					i=i+1
					// log.Println(i)
					break Loop
					// return nil
				// case param:=<-update: {
				// 	log.Println(param)
				// 	speaker.Clear()
				// 	loop := beep.Loop(param.LoopNum, streamer)	
				// 	fast := beep.ResampleRatio(4, param.Fast, loop)
				// 	ctrl := &beep.Ctrl{Streamer: fast, Paused: param.Pause}
				// 	volume := &effects.Volume{
				// 		Streamer: ctrl,
				// 		Base:     2,
				// 		Volume:   param.Volume,
				// 		Silent:   false,
				// 	}
				// 	speaker.Play(beep.Seq(volume, beep.Callback(func() {
				// 		done <- true
				// 	})))
					// speaker.Play(volume)

				// }
				// case <-time.After(time.Second):
				// 	speaker.Lock()
				// 	log.Println(format.SampleRate.D(streamer.Position()).Round(time.Second))
				// 	speaker.Unlock()
					
				}
			}
			log.Println("next")
	
		}

	}
	return nil
}
