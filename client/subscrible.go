package main

import (
    "fmt"
    "context"
    "io"
    "log"
    "sync"
    "time"
    "strings"
    "os"
    "strconv"
    "image/png"
    // "encoding/base64"
    "bytes"
    "flag"
    "os/exec"
    "regexp"

    myutil "grpc-demo/utils/my_utils"
    job "grpc-demo/utils/jobs"
    pb "grpc-demo/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/keepalive"
    util "grpc-demo/utils"
    "grpc-demo/utils/play"
    "grpc-demo/utils/cmd"
    "grpc-demo/utils/ips"
    "github.com/pkg/errors"
    "grpc-demo/utils/webrtc"
    "grpc-demo/client/models"
    qrcode "github.com/skip2/go-qrcode"

)
type Context struct {
    Ctx  context.Context
    Cancel context.CancelFunc
}
type Unit struct {
    *pb.Channel
    *pb.SubscribeResult
}
type Map struct {
    Mu *sync.Mutex
    Unit map[string]context.CancelFunc
}
func main() {
    var serverIp ,randomTag,webPort string
    var clientId,topic string
    const initRemoteServerIp = "47.99.78.179:8123"
    const initWebPort = "9090"
    const localServerIp = "127.0.0.1:8123"
    flag.StringVar(&serverIp, "s", "local", "默认服务器地址为局域网内服务器（外网服务器，参数-s 设置为remote@serverIp）")
    flag.StringVar(&randomTag, "t", "random", "默认id和topic随机生成（自定义，参数-t 设置为[id:topic],id和topic均为4-6位字母数字组成）")
    flag.StringVar(&webPort, "p", initWebPort, "默认Web网页端口 9090")
    flag.Parse()
    fmt.Println("**1.消息发布端和接收端设备ID和主题填写一致**")
    fmt.Println("**2.先开启接受端填写参数，再通过发布端发布指令**")
    fmt.Println("**3.发布端地址MP3地址可是url,也可为本地MP3（该文件和接受端放在一起,地址为 ./文件名）**")
    fmt.Println("")
    isMatch,_ :=regexp.MatchString(`[\d]+`,webPort)
    if !isMatch {
        fmt.Println("参数错误（参数：-p 格式:xxxx ,xxxx为数字）")
        return
    }
    if randomTag == "random"{
        salt := "random"
        stamp :=int(time.Now().Unix())
        randomHash,_ :=util.Hash(strconv.Itoa(stamp)+salt)
        fmt.Printf(randomHash)
        clientId = randomHash[0:6]
        topic = randomHash[len(randomHash)-6:]
    }else{
        r,_:=regexp.Compile(`\[[\w]{4,6}:[\w]{4,6}\]`)
        isMatch:=r.MatchString(randomTag)
        if isMatch {
            fmt.Println(randomTag)
            s := strings.Split(randomTag,":")
            clientId = s[0][1:]
            topic = s[1][:len(s[1])-1]
        }else{
            fmt.Println("参数错误（参数：-t 格式：[id:topic],id和topic均为4-6位字母数字组成）")
            return
        }
    }
    if serverIp == "local"{
        serverIp = localServerIp
        fmt.Println("------------------------本地服务器模式------------------------")
        fmt.Println("")
        ipAdrress := ips.GetIPs()
        if ipAdrress[0] !=""{
        }
        RenderStringQr("http://"+ipAdrress[0]+":"+webPort)
    }else{
        s := strings.Split(serverIp,"@")
        fmt.Printf(s[0])
        if s[0] == "remote"{
            if len(s) ==1 {
                serverIp =  initRemoteServerIp
            }
            if len(s)==2{           
                isMatch:=IsIP(s[1])
                if isMatch {
                    serverIp = s[1]
                }else{
                    fmt.Println("serverIp 格式为 xxxx.xxxx.xxxx.xxxx:xxxx")
                    return
                }
            }
            fmt.Println("------------------------远程服务器模式------------------------")
            RenderStringQr(strings.Split(serverIp,":")[0]+":"+webPort)
        }else{
            fmt.Println("参数错误（参数：-s 格式：remote|serverIp")
            return
        }
    }
    var cmdCancleMap=Map{
        new(sync.Mutex),
        make(map[string]context.CancelFunc),
    }
        var kacp = keepalive.ClientParameters{
        Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
        Timeout:             20*time.Second,      // wait 1 second for ping ack before considering the connection dead
        PermitWithoutStream: true,             // send pings even without active streams
    }
    reconnect := make (chan bool)
    for{
        sdpchan := make (chan string)
        answerchan := make (chan string)
        signal := make (chan Unit)
        reback := make (chan *pb.PulishMessage)
        go webrtcserver.Run(sdpchan,answerchan,reback)//开启webrtc服务端
        log.Println("-----Start grpc Client-----")
        conn, err := grpc.Dial(serverIp, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
        if err != nil {
             log.Println(err)
             defer conn.Close()
             reconnect <-true
        }
        
        client := pb.NewPubsubServiceClient(conn)
        fmt.Println("===============消息指令接受端==============="+"[设备ID:"+clientId+" 主题:"+topic+"]")
     
        upContent :=Context{}

    // for {
        //处理订阅消息go
        go func(){
            for {
                select{
                case v :=<-signal:
                    fmt.Println(v.Type)
                    switch v.Type {
                        case 0:go playsoundWork(v,upContent,reback)
                        case 2:go cmdWork(v,cmdCancleMap,reback)
                        case 3:go shotWork(v,reback)
                        case 4:go sdpwork(&v,sdpchan,answerchan,reback)
                    }
                }
            }   
        }()
        //处理发布消息go
        go func(){
            for {
                select{
                case v :=<-reback:
                    Publish(client,v)
                }
            }   
        }()
        //从服务器获取订阅的参数，并完成订阅，并向订阅处理信息通道（signal）发送数据
        c:=make (chan *pb.Channel)
        go func () {
            err :=GetChanneklList(client,&pb.Channel{Name:clientId,Topic:topic},c)
            if err !=nil{
                log. Println(err)
                reconnect <-true
            }
        }()
        go func(){
            for channelMessage := range c {
                go func(channelMessage *pb.Channel){
                    ch1:=make (chan *pb.SubscribeResult)
                    go func (ch1 chan *pb.SubscribeResult) {
                         err := Subscribe(client,channelMessage,ch1)
                         if err !=nil{
                            reconnect <- true
                         }
                    }(ch1)
                    for resultMessage := range ch1 {
                        signal<-Unit{channelMessage,resultMessage}
                    }
                }(channelMessage)
            }

        }()
        <- reconnect
        time.Sleep(5*time.Second)
    }
}
//根据客户端获取该用户channel列表
func  GetChanneklList(client pb.PubsubServiceClient, agent *pb.Channel , c chan *pb.Channel) error{
    stream, err :=client.GetChanneklList(context.Background(),agent)
    if err !=nil{
        return err
    }
    for {
        reply, err := stream.Recv()
        if err != nil {
            if err == io.EOF {
                break
            }
            log.Println(err)
            return err
        }
        c <- reply
    }
    return nil  
}
//根据channel内容完成订阅
func  Subscribe(client pb.PubsubServiceClient,channel *pb.Channel,c chan *pb.SubscribeResult) error{
    stream, err := client.Subscribe(
        context.Background(), channel,
    )
    if err != nil {
        log.Println(err)
        return err
        //log.Fatal(err)
    }
    errChan := make(chan error)
    go func(ch chan error){
        for {
            reply, err := stream.Recv()
            if err != nil {
                if err == io.EOF {
                    break
                }
                 log.Println("!_!")
                log.Println(err)
                ch<-err
                return 
            }
            c<-reply
        }
    }(errChan)
    return <-errChan
}
//发布内容
func  Publish(client pb.PubsubServiceClient,msg  *pb.PulishMessage){
    _, err := client.Publish(context.Background(),msg )
    if err != nil {
        log.Println(err)
        // log.Fatal(err)
    }
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
func cmdWork(v Unit,cmdCancleMap Map,c chan *pb.PulishMessage) {
    log.Println("cmdWork")
    var cmdHash string
    var err error
    errorCh := make (chan error)
    newTopic := fmt.Sprintf("%s_%d",v.Topic,v.Timestamp)
    go func(){
        for {
            select {
            case err :=<-errorCh:
                log.Println(err)
                errString :=errors.Wrap(err, v.Name+"|"+newTopic).Error()
                c <-&pb.PulishMessage {
                    Topic: &pb.Channel{Name: v.Name,Topic: newTopic},
                    Result: &pb.SubscribeResult{Msg:errString},
                }
            }
        }
    }()
    cmdline :=v.Msg
    if v.Msg == ""{
        errorCh <- errors.New("cmd is null")
        return
    }
    cmdline = strings.Replace(cmdline,"|_|"," ",-1)
    log.Println(cmdline)
    if cmdHash,err = util.Hash(cmdline);err!=nil{
        log.Println(err)
        return
    }
    log.Println(v.Pause)
    go func(){
        stop :=v.Pause
        if stop {
            cmdCancleMap.Mu.Lock()

            cancel, ok := cmdCancleMap.Unit[cmdHash]
            if ok {
                cancel()
                log.Println(".... stop ...")
                errorCh <- errors.New("live steam stop...")
                
            }else{
                log.Println("cmd is not valid")
                errorCh <- errors.New("cmd of stopping is not valid")
            }
            cmdCancleMap.Mu.Unlock()
        }else{
            cmdCancleMap.Mu.Lock()
            cancel, ok := cmdCancleMap.Unit[cmdHash]
            if ok {
                log.Println("the same cmd of before stop ")
                cancel()
            }
            cmdCancleMap.Mu.Unlock()
            // cmdline := "ffmpeg -re -i D:/21.mp4 -c copy -f flv rtmp://localhost:1935/live/movie"
            cmd := cmd.New(cmdline)
            cmdCancleMap.Unit[cmdHash] =cmd.Cancel
            log.Printf("%+v",cmdCancleMap)
            errorCh <- errors.New("live steam Start...")
            err := cmd.Run()
            if err != nil {
                fmt.Println(err)
                errorCh<-err
                cmdCancleMap.Mu.Lock()
                delete(cmdCancleMap.Unit,cmdHash)
                cmdCancleMap.Mu.Unlock()
            }
        }
    }()
    
}
func playsoundWork(v Unit,upContent Context,c chan *pb.PulishMessage) {
    path := v.Msg
    option := media.Option{LoopNum:int(v.Loop),Fast:float64(v.Fast),Pause:v.Pause,Volume:float64(v.Volume)}
    upCancel :=upContent.Cancel
    newTopic :=fmt.Sprintf("%s_%d",v.Topic,v.Timestamp)
    // formatTimeStr :=time.Unix(v.Timestamp,0).Format("2006-01-02 15:04:05")
    errorCh := make (chan error)

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
    go func(){
        for {
            select {
            case err :=<-errorCh:
                errString :=errors.Wrap(err, v.Name+"|"+newTopic).Error()
                c <-&pb.PulishMessage {
                    Topic: &pb.Channel{Name: v.Name,Topic: newTopic},
                    Result: &pb.SubscribeResult{Msg:errString},
                }
            }
        }
    }()
    errorCh <- errors.New("play sound Start...")
}
func shotWork(v Unit,c chan *pb.PulishMessage) {
    go func () {
        initPath,err := myutil.GetCurrentDirectory()
        if err!=nil{
            log.Println(err)
        }
        //生成本地文件保存路径
        // parentPath:=initPath[:strings.LastIndex(initPath, "/")]
        shotCmd := v.Msg
        time_now := time.Now()
        year, month, day := time.Now().Date()
        fileName := fmt.Sprintf("%d_%d_%d",year,month,day)
        subfileName,_:= util.Hash(shotCmd)
        direct:=initPath+"/public/images/"+fileName+"/"+subfileName
        // direct:=initPath+"/public/images/"+fileName
        Exist,_:=myutil.PathExists(direct)
        if (!Exist){
            err=os.MkdirAll(direct,os.ModePerm)
            if err!=nil{
               fmt.Println(err)
            }
        }
        savepath := direct  
        savename := "image_"+strconv.FormatInt(time_now.UnixNano(),10)+".png"
        newTopic := fmt.Sprintf("%s_%d",v.Topic,v.Timestamp)
        errorCh := make (chan error)
        doneCh := make (chan bool)
        
        go func () {
            for ms :=range errorCh {
                log.Println(ms)
                errString :=errors.Wrap(ms, v.Name+"|"+newTopic).Error()
                c <-&pb.PulishMessage {
                    Topic: &pb.Channel{Name: v.Name,Topic: newTopic},
                    Result: &pb.SubscribeResult{Msg:errString},
                }
            }
        }()
        go func (){
            loop:for{
                    select {
                        case <-doneCh :
                            log.Printf("shot end")
                            imageBytes ,err:= myutil.FindFile(savepath,savename)
                            if err != nil {
                                log.Println(err)
                                return
                            }
                            //缩略图保存到数据库
                            img_thumbnail := myutil.ImgResize(imageBytes)
                            buf := new(bytes.Buffer)
                            err = png.Encode(buf, img_thumbnail)
                            if err !=nil{
                                return
                            }
                            img_thumbnail_byte := buf.Bytes()
                            imageBytes = img_thumbnail_byte
                            save := models.Picture{
                                Name:shotCmd,
                                Data:imageBytes,
                                CreatedAt: time_now,
                            }
                            models.Datachan <- save
                            // //返回给客户端
                            // sourcestring := base64.StdEncoding.EncodeToString(img_thumbnail_byte)
                            // errorCh <- errors.New(sourcestring)
                            break loop
                            return
                        case <-time.After(30*time.Second):
                            log.Println("shot picture lost")
                            return
                            break
                    }
                }

        }()
        if shotCmd == ""{
            return
        }else{
            if shotCmd != "screen" {
                shotCmd = fmt.Sprintf("ffmpeg |_|-i|_|%s|_|-vframes|_|1|_|-y|_|-f|_|image2",shotCmd)
            }
            job.Shot(shotCmd,savepath,savename,doneCh,errorCh)
        }
    }()
    //返回给请求端
    // if imgbase64,err :=myutil.ImgToBase64(localfilePath) ;err!=nil{
    //     // picture :=models.Picture{
    //     //     Data:imgbase64,
    //     // }
    //     // models.SavePicture(&picture)
    //     // var pictures []models.Picture
    //     // err := models.Db.AllByIndex("CreatedAt", &pictures)
    //     // if err !=nil {
    //     //     log.Println(err)
    //     // }
    //     // log.Printf("***%+v***",pictures)
    //     errorCh <- errors.New(err.Error())
    // }else{
    //     errorCh <- errors.New(imgbase64)
    // }
    

    //截图完成上传值minio服务器
    // select {
    //     case endflag:=<-endChan:
    //         if endflag{
    //             utils.PutObject(minioC,savePath+"/"+savename+".jpg",BucketName,savename)
    //             fmt.Println("this is end1")
    //         } 
    // }
    // time.Sleep(time.Second)
    log.Println("Done")
}
func sdpwork (v *Unit,sdpchan,answerchan chan string,c chan *pb.PulishMessage){
    sdpchan <-v.Msg
    newTopic := fmt.Sprintf("%s_%d",v.Topic,v.Timestamp)
    loop:
        for {
            select{
            case answer := <-answerchan:
                errString :=errors.Wrap(errors.New(answer), v.Name+"|"+newTopic).Error()
                c <-&pb.PulishMessage {
                    Topic: &pb.Channel{Name: v.Name,Topic: newTopic},
                    Result: &pb.SubscribeResult{Msg:errString},
                }
                break loop //不退出循环，newTopic的值和loop外值不一样？
            }
        }
   
}
func write2Stdout(content string, level qrcode.RecoveryLevel, size int) error {
    var q *qrcode.QRCode
    q, err := qrcode.New(content, level)
    if err != nil {
        return err
    }
    qrstring := q.ToSmallString (false)
    fmt.Println(qrstring)
    // if _,err :=os.Stdout.Write([]byte(qrstring));err!=nil{
    //     return err
    // }
    return nil
}
func RenderString(s string) {
    q, err := qrcode.New(s, qrcode.Medium)
    if err != nil {
        panic(err)
    }
    fmt.Println(q.ToSmallString(false))
}
func RenderStringQr(s string) {
    err := qrcode.WriteFile(s, qrcode.Medium, 200, "qr.png")
    if err != nil {
        panic(err)
    }
    cmd:=exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", "qr.png")
    cmd.Start()
}
func IsIP(ip string)(b bool){
    m,_ :=regexp.MatchString(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}:[\d]+$`,ip)
    if !m { 
        return false 
    }else{ 
        return true 
    } 
 }