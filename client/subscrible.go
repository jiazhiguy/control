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
    // "image/jpeg"

    
    myutil "grpc-demo/utils/my_utils"
    // "github.com/vova616/screenshot"
    job "grpc-demo/utils/jobs"
    pb "grpc-demo/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/keepalive"
    util "grpc-demo/utils"
    "grpc-demo/utils/play"
    "grpc-demo/utils/cmd"
    "github.com/pkg/errors"
    "grpc-demo/utils/webrtc"
    "flag"
    // "grpc-demo/client/models"

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
    var serverIp string
    var clientId,topic string
    const remoteServerIp = "47.99.78.179:8123"
    const localServerIp = "127.0.0.1:8123"
    // sdpchan := make (chan string)
    // answerchan := make (chan string)
    // signal := make (chan Unit)
    // reback := make (chan *pb.PulishMessage)
    // grpc客户端 完成交互信号的接受和处理
    // const serverIp ="localhost:8123"
    // const serverIp = "47.99.78.179:8123"
    flag.StringVar(&serverIp, "s", "local", "服务器地址默认为局域网内服务器（外网服务器，参数-s 设置为remote）")
    flag.Parse()
    if serverIp == "remote"{
        serverIp = remoteServerIp
        fmt.Println("------------------------远程服务器模式------------------------")
        fmt.Println("")
     }
     if serverIp == "local"{
        serverIp = localServerIp
        fmt.Println("------------------------本地服务器模式------------------------")
        fmt.Println("")
    }
    
    fmt.Println("**1.消息发布端和接收端设备ID和主题填写一致**")
    fmt.Println("**2.先开启接受端填写参数，再通过发布端发布指令**")
    fmt.Println("**3.发布端地址MP3地址可以是url,也可以是本地MP3,本地文件和接受端放在一起,地址为 ./文件名**")
    fmt.Print("自定义设备ID:")
    fmt.Scanln(&clientId)
    fmt.Print("自定义主题:")
    fmt.Scanln(&topic)


    var cmdCancleMap=Map{
        new(sync.Mutex),
        make(map[string]context.CancelFunc),
    }
    // var cmdCancleMap = make(map[string]context.CancelFunc)
    // var syncLock sync.RWMutex
    // ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
    // defer cancel()
    // conn, err := grpc.DialContext(ctx, serverIp, grpc.WithBlock(), grpc.WithInsecure())
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
        // conn, err := grpc.Dial(serverIp, grpc.WithInsecure())
        if err != nil {
             log.Println(err)
             defer conn.Close()
             reconnect <-true
            // log.Fatal(err)
        }
        
        client := pb.NewPubsubServiceClient(conn)

        // fmt.Print("用户名:")
        // fmt.Scanln(&usernanme)
        // fmt.Print("登录密码:")
        // fmt.Scanln(&password)
        // fmt.Print("设备ID:")
        // fmt.Scanln(&clientId)
        // fmt.Print("主题:")
        // fmt.Scanln(&topic)
        // if clientId =="" {
        //     clientId = "jz101"
        // }
        // if topic ==""{
        //     topic= "playsound"
        // }
        fmt.Println("===============消息指令接受端==============="+"[设备ID:"+clientId+" 主题:"+topic+"]")
     
        upContent :=Context{}
        // path := "../Lame_Drivers_-_01_-_Frozen_Egg.mp3"http://47.99.78.179:9000/temp/095889a1aa6db7b345acf3cf5b20fcc9?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=TEST%2F20200213%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20200213T124322Z&X-Amz-Expires=6000&X-Amz-SignedHeaders=host&response-content-disposition=attachment%3B%20filename%3D095889a1aa6db7b345acf3cf5b20fcc9&X-Amz-Signature=7b82622a43727855d0da46cc6504a0079af4b6270d7784182c5504d81efd500c
        // path := "http://47.99.78.179:9000/temp/095889a1aa6db7b345acf3cf5b20fcc9?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=TEST%2F20200214%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20200214T021216Z&X-Amz-Expires=6000&X-Amz-SignedHeaders=host&response-content-disposition=attachment%3B%20filename%3D095889a1aa6db7b345acf3cf5b20fcc9&X-Amz-Signature=d705ce908ed8d8ccee430aae0129b5d3ee1b2b87bce54e31589d450fd2c12155"
        // option := media.Option{LoopNum:1,Fast:1,Pause:false,Volume:0}
        
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
                // close(c)
                break
            }
            log.Println(err)
            return err
            // log.Fatal(err)
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
                // log.Fatal(err)
            }
            c<-reply
        }
    }(errChan)
    return <-errChan
    // return nil
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
            log.Printf("%+v",cmdCancleMap)
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
    log.Println("^^^",newTopic)
    go func(){
        for {
            select {
            case err :=<-errorCh:
                log.Println(">>>",newTopic)
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
        direct:=initPath+"/public/images"
        fmt.Println(direct)
        Exist,_:=myutil.PathExists(direct)
        if (!Exist){
            err=os.MkdirAll(direct,os.ModePerm)
            if err!=nil{
               fmt.Println(err)
            }
        }
        savepath := direct  
        savename := "image_"+strconv.FormatInt(time.Now().UnixNano(),10)+".png"
        newTopic := fmt.Sprintf("%s_%d",v.Topic,v.Timestamp)
        errorCh := make (chan error)
        doneCh := make (chan bool)
        // shotCmd :=""
        shotCmd := v.Msg
        // timestamp :=0
        // if ms ,ok := cmdms.Message.(map[string]interface{}) ;ok{
        //     shotCmd = ms["shot"].(string)
        // }
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
            
            log.Println("~~~~~***~~",newTopic)
            select{
            case answer := <-answerchan:
                // newTopic := fmt.Sprintf("%s_%d",v.Topic,v.Timestamp)
                log.Println("~~~~~~~~~~",newTopic)
                errString :=errors.Wrap(errors.New(answer), v.Name+"|"+newTopic).Error()
                // log.Printf("%s",errString)
                c <-&pb.PulishMessage {
                    Topic: &pb.Channel{Name: v.Name,Topic: newTopic},
                    Result: &pb.SubscribeResult{Msg:errString},
                }
                break loop //不退出循环，newTopic的值和loop外值不一样？
            }
        }
   
}
