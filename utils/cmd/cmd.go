//https://github.com/wenchengyao/smpNVR
package cmd
import (
	// "log"
	// "syscall"
	"context"
	"os"
	"os/exec"
	"strings"
)

type C struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	Cmd    *exec.Cmd
	Done    chan bool
}
func New(cmdLine string) C {
	var c C
	doneCh :=make (chan bool)
	args := strings.Fields(cmdLine)
	c.Done = doneCh
	c.Ctx, c.Cancel = context.WithCancel(context.Background())
	c.Cmd = exec.CommandContext(c.Ctx, args[0], args[1:]...)
	// c.Cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	c.Cmd.Stdout = os.Stdin
	return c
}
// Run runs ffmpeg with given set of arguments, optional callback will be used to report progress (current duration,
// total duration). Callback total duration can be 0 if unable to automatically detect.
func (c *C) Run() error {
	err := c.Cmd.Start()
	if err != nil{
		return err
	}
	err = c.Cmd.Wait()
	if err != nil{
		c.Done <-false
		return err
	}
	c.Done <-true
	// for {
	// 	select {
	// 		case <-c.Ctx.Done():
	// 			log.Println("cancel")
	// 			return 
	// 	}
	// }
	return nil
}
// func (c *C) RunThenClose(ch chan int) error {
// 	err := c.Cmd.Start()
// 	c.cancel()
// 	//tell go i take it
// 	return err
// }
func (c *C) Close() {
	c.Cancel()
}
