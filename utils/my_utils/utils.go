package utils
import (
    "runtime"
    "encoding/base64"
    "path/filepath"
    "os"
    "strings"
    "errors"
    "io/ioutil"
    "fmt"
    "image"
    "bytes"
    "image/jpeg"
    "image/png"
    "log"


    "github.com/nfnt/resize"
)


func CurrentFile() (string,error) {
    _, file, _, ok := runtime.Caller(1)
    if !ok {
        // panic(errors.New("Can not get current file info"))
        return "",errors.New("Can not get current file info")
    }
    return file,nil
}
func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}
func ImgToBase64(filepath string) (string,error){
    n :=0
    if f,err := os.Open(filepath);err!=nil {
        return "",err
    }else{
        sourcebuffer := make([]byte, 500000)
        if n,err = f.Read(sourcebuffer) ;err !=nil{
            return "",err
        }
        sourcestring := base64.StdEncoding.EncodeToString(sourcebuffer[:n])
        return sourcestring,nil
    }
}

func substr(s string, pos, length int) string {
    runes := []rune(s)
    l := pos + length
    if l > len(runes) {
        l = len(runes)
    }
    return string(runes[pos:l])
}

func GtParentDirectory(dirctory string) string {
    return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func GetCurrentDirectory() (string ,error){
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        return "",err
    }
    return strings.Replace(dir, "\\", "/", -1),nil
}
func FindFile(savepath,savename string) ([]byte,error) {
    path := savepath+"/"+savename

    // path,_ := getCurrentDirectory()
    // path :="E:/workbench/go/src/grpc-demo/client"
    // path =path+"/public/images/"+"image_1583985334765801700.jpg"
    // return "findFile"+hash
    f, err := os.Open(path)
    if err != nil {
        fmt.Println("read file fail", err)
        return []byte{},err
    }
    defer f.Close()
    fd, err := ioutil.ReadAll(f)
    if err != nil {
        fmt.Println("read to fd fail", err)
        return []byte{},err
    }
    return fd,nil
}
func ImgResize(byteData []byte) image.Image{
    reader := bytes.NewReader(byteData)
    img, err:= jpeg.Decode(reader)
    if err != nil {
        // log.Printf("jpg decode:%s",err)
        reader := bytes.NewReader(byteData)
        img, err = png.Decode(reader)
        if err != nil {
            log.Printf("png decode:%s",err)
        }
    }
    if img!=nil{
        // img = resize.Resize(1000, 0, img, resize.Lanczos3)
        img = resize.Resize(200, 200, img, resize.Lanczos3)
    }

    return img
}