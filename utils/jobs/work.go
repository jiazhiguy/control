package job
import (
	"log"
	"errors"
	"strings"
	"os"
	"fmt"
    "context"
	// "image/jpeg"
	"image/png"
    "net/http"
    "io/ioutil"
    "mime/multipart"
    "bytes"
    "time"

	"github.com/vova616/screenshot"
    "grpc-demo/utils/play"
	"grpc-demo/utils/cmd"
)
type Context struct {
    Ctx  context.Context
    Cancel context.CancelFunc
}
var (
    HttpClient = &http.Client{
        Timeout: 3 * time.Second,
    }
)
func Shot(cmdline string ,savepath,savename string,done chan bool,errorCh chan error)  {
	log.Println("shot")
    if cmdline == ""{
        errorCh <- errors.New("shot cmd is null")
        return
    }
    cmdline = strings.Replace(cmdline,"|_|"," ",-1)
    localfilePath := savepath+"/"+savename
    if cmdline=="screen" {//截屏
        go func(done chan bool){
            // errorCh <- errors.New("shot picture Start...")
            img, err := screenshot.CaptureScreen()
            if err != nil {
               errorCh<-err
            }
            f, err := os.Create(localfilePath)
            if err != nil {
                fmt.Println(err)
            }
            err = png.Encode(f,img)
            if err != nil {
                fmt.Println(err)
            }
            f.Close()
            done <- true
        }(done)
    }else{//调用ffmpeg进行截图
        newcmd :=cmdline +" "+ localfilePath
        go func(newcmd string){
            cmd := cmd.New(newcmd)
            cmd.Done=done
            // errorCh <- errors.New("shot picture Start...")
            err := cmd.Run()
            if err != nil {
                fmt.Println(err)
                errorCh<-err
            }
        }(newcmd)
    }
    <-done
    // errorCh <-errors.New("shot picture end...")
    // log.Println("Done")
    return
    //截图完成上传值minio服务器
    // select {
    //     case endflag:=<-endChan:
    //         if endflag{
    //             utils.PutObject(minioC,savePath+"/"+savename+".jpg",BucketName,savename)
    //             fmt.Println("this is end1")
    //         } 
    // }
    // time.Sleep(time.Second)
}
func PlaySound(ctx context.Context,path string,option media.Option,ch chan error) (chan media.Option){
   log.Println("playing ["+path+"]")
   update:=make(chan media.Option)
   go func () {
       if err :=media.PlayMp3(path,option,update); err != nil {
            fmt.Println(err)
            ch <-err
       }
       // for {
       //     select{
       //     case <-ctx.Done():
       //      log.Println("Cancel")
       //      return
       //     }
       // }

   }()
   return update
}
func Sound(path string,loop int,upContent Context,errorCh chan error) {
    option := media.Option{LoopNum:int(loop),Fast:float64(1),Pause:false,Volume:float64(0)}
    upCancel :=upContent.Cancel
    if upCancel == nil{
        ctx, cancel := context.WithCancel(context.Background())
        PlaySound(ctx,path,option,errorCh)
        upContent=Context{ctx,cancel}
    }else{
        upContent.Cancel()
        _ctx, _cancel := context.WithCancel(context.Background())
        PlaySound(_ctx,path,option,errorCh)
        upContent=Context{_ctx,_cancel}
    }
    errorCh <- errors.New("play sound Start...")
}

func UploadFile(url string, params map[string]string, nameField, fileName string, fileData []byte) ([]byte, error) {
    body := new(bytes.Buffer)
    body_writer := multipart.NewWriter(body)
    //写入文件
    formFile, err := body_writer.CreateFormFile(nameField, fileName)
    if err != nil {
        return nil, err
    }
    formFile.Write(fileData)
    // 其他参数列表写入body
    for key, val := range params {
        _ = body_writer.WriteField(key, val)
    }
    err = body_writer.Close()
    if err != nil {
        return nil, err
    }
    req, err := http.NewRequest("POST", url, body)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", body_writer.FormDataContentType())
    req.Header.Add("X-No-Sign","yes")
    resp, err := HttpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    content, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    return content, nil
}