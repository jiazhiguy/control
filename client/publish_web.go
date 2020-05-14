package main

import (
    "fmt"
    "context"
    "log"
    "io"
    "time"
    "strconv"
    "net/http"
    "strings"
    "os/signal"
    "syscall"
    "os"
    "errors"
    "flag"
    "encoding/json"

    "github.com/pion/webrtc/v2"
    sdp "grpc-demo/utils/signal"
    message "grpc-demo/utils/message"
    gin "github.com/gin-gonic/gin"
    pb "grpc-demo/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/keepalive" 
    "gopkg.in/tylerb/graceful.v1"
    "github.com/robfig/cron"
    "grpc-demo/utils/cron"
    // "grpc-demo/utils"
    // "github.com/zserge/webview"
 
)
// 
//ffmpeg |_|-i|_|rtsp://admin:admin123@192.168.2.241/h264/ch1/main/av_stream|_|-vframes|_|1|_|-y|_|-f|_|image2
//ffmpeg|_|-re|_|-i|_|rtsp://admin:admin123@192.168.2.241/h264/ch1/main/av_stream|_|-c|_|copy|_|-f|_|flv|_|rtmp://47.99.78.179:1935/live/movie
var cronHub *cron.Cron
var serverIp string
var clientId,topic string
func main() {
    const httpPort = ":9090"
    const remoteServerIp = "47.99.78.179:8123"
    const localServerIp = "127.0.0.1:8123"

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
    
    fmt.Print("自定义设备ID:")
    fmt.Scanln(&clientId)
    fmt.Print("自定义主题:")
    fmt.Scanln(&topic)
    go func(){
        webrtcConn(clientId,topic,"4","http://"+strings.Split(serverIp,":")[0]+httpPort+"/message")

    }()
    cronHub = cron.New()
    cronHub.Start() 
    // go func () {
    //     cronHub.Start() 
    //     select{}
    // }()
    // gin.SetMode(gin.ReleaseMode)
    router := gin.Default()
    router.Use(Cors())
    router.Static("/static","dist/static")   // 添加资源路径
    router.StaticFS("/down", http.Dir("./tmp"))
    router.StaticFile("/", "./dist/index.html")  //前端接口
    // cmdOut :=make(chan *message.Response,5)
    // cmdIn := make(chan *message.Cmd)
    // go PubishServer(cmdIn,cmdOut)
    router.POST("/cron", func(c *gin.Context) {
        tip := c.PostForm("tip")
        if tip == ""{
            c.JSON(200,gin.H{
                "message":"paraments is null",
            })
            return
        }
        log.Printf("tip:%s",tip)
        cmdms,err :=parse(c)
        log.Printf("cmdms1:%+v\n",cmdms)
        if err != nil{
            c.JSON(200,gin.H{
                "message":err.Error(),
            })
            log.Println(err)
            return 
        }
        log.Printf("cmdms2:%+v\n",cmdms)
        err = cronHub.AddFunc(tip,func(){
            cronutil.HttpPost(cmdms,"http://"+strings.Split(serverIp,":")[0]+httpPort+"/message")
        })
        if err != nil { 
            log.Println(err)           
            c.JSON(200,gin.H{
                "message":"add cron fail",
            })
            return 
        } 
        c.JSON(200,gin.H{
            "message":cmdms,
        })
        return
    })
    router.POST("/message", func(c *gin.Context) {
        log.Println("*******HANDEL.message*******")
        log.Println(" ")
        cmdOut :=make(chan *message.Response,5)
        cmdIn := make(chan *message.Cmd)
        go PubishServer(cmdIn,cmdOut)
        var cmdms *message.Cmd
        cmdms,err := parse(c)
        // log.Printf("^^^%+v",c)
        if err !=nil {
            c.JSON(200,gin.H{
                "message":err.Error(),
            })
            return 
        }
        if cmdms == nil {
            c.JSON(200,gin.H{
                "message":"fail",
            })
            return
        }
        cmdIn <- cmdms
        loop: for {
                    select{
                    case data := <-cmdOut:
                        log.Println(data)
                        c.JSON(200,gin.H{
                            "message":data.Msg,
                        })
                        break loop
                        return
                }
            }
        // }()
        // loop:
        // for {
        //     select{
        //         case data :=<-cmdOut:{
        //             log.Println(data.Cmdmessage == cmdms)
        //             //     for  {
        //                     // if (data==Response{}){
        //                     //     continue
        //                     // }
        //                     log.Println("^^^^^^^^^^from clint^^^^^^^^\n")
        //                     // if data.cmdmessage == cmdms&&c!=nil{
        //                         c.JSON(200,gin.H{
        //                             "message":data.Msg,
        //                         })
        //                         // break//for中只能发送1次，连续发送报错？？
        //                     // }
        //             //     }
        //             //     return   
        //             // }()
                    
        //             break loop
        //             return
        //         }
        //         // case <-time.After(10*time.Second):{
        //         //     log.Println("timeout")
        //         //     if c!=nil{
        //         //         c.JSON(200,gin.H{
        //         //             "message":"timeout/device offline",
        //         //         })
        //         //         // break loop
        //         //         return
        //         //     }


        //         // }
        //     }
        //     }
    })
    fmt.Println("")
    fmt.Println("-----------使用方法----------")
    fmt.Println("")
    fmt.Println("使用Google Chrome浏览器,打开网页 http://localhost:9090/")
    // server := &http.Server{
    //     Addr:           ":9090",
    //     Handler:        router,
    //     // ReadTimeout:    10 * time.Second,
    //     // WriteTimeout:   10 * time.Second,
    //     // MaxHeaderBytes: 1 << 20,
    // }
    // gracefulExitWeb(server)
    graceful.Run(httpPort,10*time.Second,router)
    // router.Run(":9090")
}
func PubishServer(cmdIn chan *message.Cmd,cmdOut chan *message.Response) {

    fmt.Println("**1.消息发布端和接收端设备ID和主题填写一致**")
    fmt.Println("**2.先开启接受端填写参数，再通过发布端发布指令**")
    fmt.Println("**3.发布端地址MP3地址可以是url,也可以是本地MP3,本地文件和接受端放在一起,地址为 ./文件名**")
    // const serverIp ="localhost:8123"
    //********************connect*********************
    // const serverIp ="47.99.78.179:8123"
    // cmdline := "ffmpeg -re -i D:/21.mp4 -c copy -f flv rtmp://localhost:1935/live/movie"
    var kacp = keepalive.ClientParameters{
        Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
        Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
        PermitWithoutStream: true,             // send pings even without active streams
    }
    reconnect := make (chan bool)
    for {
        conn, err := grpc.Dial(serverIp, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
        // conn, err := grpc.Dial(serverIp, grpc.WithInsecure())
        if err != nil {
            defer conn.Close()
            reconnect <-true
            // log.Fatal(err)
        }
        client := pb.NewPubsubServiceClient(conn)
        //******************* 发布和订阅内容处理***********************
        signal :=make (chan *pb.Channel)
        send := make (chan *pb.PulishMessage)
        response :=make (chan *pb.SubscribeResult)
        //处理订阅消息go
        go func(){
            loop:
            for {
                select{
                case v :=<-signal:
                    // fmt.Println(v)
                    err := Subscribe(client,v,response)
                    if err != nil{
                        reconnect <- true
                        break loop
                        
                        return
                    }
                }
            }
        }()
        //处理发布消息go
        go func(){
            loop:
            for {
                select{
                case v :=<-send:
                   err := Publish(client,v)
                   if err != nil {
                        reconnect <- true
                        break loop
                        return
                   }
                }
            } 
        }()
        fmt.Println("==================消息指令发布端==================")
        go func(){
            for {
                select {
                case cmdInput := <-cmdIn:{
                    clientId :=cmdInput.ClientId
                    topic :=cmdInput.Topic
                    optype :=cmdInput.Type
                    switch optype{
                        case "1":{
                            if ms ,ok := cmdInput.Message.(message.Msg_1);ok{
                                path := ms.Uri
                                loop := ms.LoopNum
                                timestamp := ms.Timestamp
                                newTopic :=fmt.Sprintf("%s_%d",topic,timestamp)
                                channel := &pb.Channel{Name: clientId,Topic: newTopic}
                                signal <-channel
                                go func (cmdInput *message.Cmd,cmdOut chan *message.Response) {
                                    log.Printf("***%+v***\n",cmdInput)
                                    for responseMessage :=range response{
                                        log.Printf("^^^^%+v^^^^^\n",cmdInput)
                                        fmt.Println("-----接受端返回信息:-----\n"+responseMessage.Msg)
                                        cmdOut <- &message.Response {
                                            cmdInput,
                                            responseMessage.Msg,
                                        }
                                    }
                                }(cmdInput,cmdOut)
                                send <- &pb.PulishMessage{
                                        Topic: &pb.Channel{Name: clientId,Topic: topic},
                                        Result: &pb.SubscribeResult{Msg: path,Loop: int32(loop),Fast:1,Pause:false,Volume:0,Type:0,Timestamp:timestamp},
                                    }
                            }else{
                                cmdOut <- &message.Response{
                                    cmdInput,
                                    "input paraments error",
                                }
                            }
                        }
                        case "2":{
                            if ms ,ok := cmdInput.Message.(message.Msg_2);ok {
                                cmdPause := ms.Operate
                                timestamp := ms.Timestamp
                                liveCmd := ms.LiveCmd
                                newTopic :=fmt.Sprintf("%s_%d",topic,timestamp)
                                channel := &pb.Channel{Name: clientId,Topic: newTopic}

                                signal <-channel
                                go func (cmdInput *message.Cmd,cmdOut chan *message.Response) {
                                    // log.Printf("***%+v***\n",cmdInput)
                                    for responseMessage :=range response{
                                        // log.Printf("^^^^%+v^^^^^\n",cmdInput)//为什么这两处cmdinput不相等
                                        fmt.Println("-----接受端返回信息:-----\n"+responseMessage.Msg)
                                        cmdOut <- &message.Response{
                                            cmdInput,
                                            responseMessage.Msg,
                                        } 
                                    }
                                }(cmdInput,cmdOut)
                                sendcmd := &pb.PulishMessage{
                                            Topic: &pb.Channel{Name: clientId,Topic: topic},
                                            Result: &pb.SubscribeResult{Msg: liveCmd,Type:2,Pause:cmdPause,Timestamp:timestamp},
                                        }
                                log.Printf("%+v",sendcmd)
                                send <- sendcmd

                            }else{
                                cmdOut <- &message.Response{
                                    cmdInput,
                                    "input paraments error",
                                }
                            }
                        }
                        case "3":{
                            if ms ,ok := cmdInput.Message.(message.Msg_3);ok {
                                shotCmd := ms.ShotCmd
                                timestamp := ms.Timestamp
                                // path := ms.Uri
                                newTopic :=fmt.Sprintf("%s_%d",topic,timestamp)
                                channel := &pb.Channel{Name: clientId,Topic: newTopic}

                                signal <-channel
                                go func (cmdInput *message.Cmd,cmdOut chan *message.Response) {
                                    // log.Printf("***%+v***\n",cmdInput)
                                    for responseMessage :=range response{
                                        // log.Printf("^^^^%+v^^^^^\n",cmdInput)//为什么这两处cmdinput不相等
                                        fmt.Println("-----接受端返回信息:-----\n"+responseMessage.Msg)
                                        cmdOut <- &message.Response{
                                            cmdInput,
                                            responseMessage.Msg,
                                        } 
                                    }
                                }(cmdInput,cmdOut)
                                sendcmd := &pb.PulishMessage{
                                            Topic: &pb.Channel{Name: clientId,Topic: topic},
                                            Result: &pb.SubscribeResult{Msg:shotCmd,Type:3,Timestamp:timestamp},
                                        }
                                log.Printf("%+v",sendcmd)
                                send <- sendcmd

                            }else{
                                cmdOut <- &message.Response{
                                    cmdInput,
                                    "input paraments error",
                                }
                            }
                        }
                        case "4":{
                            log.Println("CASE 4 .....")
                            if ms ,ok := cmdInput.Message.(message.Msg_4);ok {
                                offer := ms.Offer
                                timestamp := ms.Timestamp
                                newTopic :=fmt.Sprintf("%s_%d",topic,timestamp)
                                channel := &pb.Channel{Name: clientId,Topic: newTopic}

                                signal <-channel
                                go func (cmdInput *message.Cmd,cmdOut chan *message.Response) {
                                    // log.Printf("***%+v***\n",cmdInput)
                                    for responseMessage :=range response{
                                        if responseMessage == nil {
                                            continue
                                        }
                                        log.Printf("^^^^%+v^^^^^\n",responseMessage)//为什么这两处cmdinput不相等
                                        cmdOut <- &message.Response{
                                            cmdInput,
                                            responseMessage.Msg,
                                        } 
                                        fmt.Printf("-----接受端返回信息:%d-----\n",len(responseMessage.Msg))
                                    }
                                }(cmdInput,cmdOut)
                                sendcmd := &pb.PulishMessage{
                                            Topic: &pb.Channel{Name: clientId,Topic: topic},
                                            Result: &pb.SubscribeResult{Msg:offer,Type:4,Timestamp:timestamp},
                                        }
                                // log.Printf("%+v",sendcmd)
                                send <- sendcmd

                            }else{
                                cmdOut <- &message.Response{
                                    cmdInput,
                                    "input paraments error",
                                }
                            }
                        }
                    }
                    }

                }
            }

        }() 
        <-reconnect 
        time.Sleep(5*time.Second)
        log.Println("reconnect**")
    }
}
//发布内容
func  Publish(client pb.PubsubServiceClient,msg  *pb.PulishMessage) error {
    _, err := client.Publish(context.Background(),msg )
    if err != nil {
        log.Println(err)
        return err
        // log.Fatal(err)
    }
    return nil
}
//根据channel内容完成订阅
func  Subscribe(client pb.PubsubServiceClient,channel *pb.Channel,c chan *pb.SubscribeResult) error {
    stream, err := client.Subscribe(
        context.Background(), channel,
    )
    if err != nil {
      
        log.Println(err)
        return err
    }
    go func(){
        for {
            reply, err := stream.Recv()
            if err != nil {
                if err == io.EOF {
                    break
                    return
                }
            }
            c<-reply
        }
    }()
    return nil
}
func Cors() gin.HandlerFunc {
    return func(c *gin.Context) {
        method := c.Request.Method      //请求方法
        origin := c.Request.Header.Get("Origin")        //请求头部
        var headerKeys []string                             // 声明请求头keys
        for k, _ := range c.Request.Header {
            headerKeys = append(headerKeys, k)
        }
        headerStr := strings.Join(headerKeys, ", ")
        if headerStr != "" {
            headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
        } else {
            headerStr = "access-control-allow-origin, access-control-allow-headers"
        }
        if origin != "" {
            c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
            c.Header("Access-Control-Allow-Origin", "*")        // 这是允许访问所有域
            c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")      //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
            //  header的类型
            c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
            //  允许跨域设置                                                                                                      可以返回其他子段
            c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")      // 跨域关键设置 让浏览器可以解析
            c.Header("Access-Control-Max-Age", "172800")        // 缓存请求信息 单位为秒
            c.Header("Access-Control-Allow-Credentials", "false")       //  跨域请求是否需要带cookie信息 默认设置为true
            c.Set("content-type", "application/json")       // 设置返回格式是json
        }
        //放行所有OPTIONS方法
        if method == "OPTIONS" {
            c.JSON(http.StatusOK, "Options Request!")
        }
        // 处理请求
        c.Next()
    }
}
func gracefulExitWeb(server *http.Server) {
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
    sig := <-ch

    fmt.Println("got a signal", sig)
    now := time.Now()
    cxt, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()
    err := server.Shutdown(cxt)
    if err != nil{
        fmt.Println("err", err)
    }
    // 看看实际退出所耗费的时间
    fmt.Println("------exited--------", time.Since(now))
}
func parse(c *gin.Context) (*message.Cmd,error){
    var ms interface{}
    initPath := "./Lame_Drivers_-_01_-_Frozen_Egg.mp3"
    clientId := c.PostForm("id")
    topic := c.PostForm("topic")
    category := c.PostForm("type")
    // log.Println(clientId)
    // log.Println(topic)
    // log.Println(category)
    if clientId==""||topic==""||(category!="1"&&category!="2"&&category!="3"&&category!="4"&&category!="5") {
        return nil,errors.New("paraments is wrong")
    }
    if category == "1"{
        var loopNum int
        var err error
        loopString := c.PostForm("times")
        if loopString==""{
            loopNum = 1
        }else{
            if loopNum,err = strconv.Atoi(loopString) ;err != nil {
                return nil,errors.New("paraments of 'times' is a number")
            }  
        }
        path := c.PostForm("path")
        if path == ""{
            path =initPath
        }
        ms = message.NewMsg_1(loopNum,path)
       
    }
    if category == "2"{
        cmd := c.PostForm("cmd")
        if cmd == ""{
            return nil,errors.New("paraments of 'cmd' is null")
        }
        startString := c.PostForm("open")
        startBool:= false
        if startString == "false"{
            startBool=true
        }
        ms = message.NewMsg_2(cmd,startBool)
       
    }
    if category == "3"{
        cmd := c.PostForm("cmd")
        if cmd ==""{
            return nil,errors.New("paraments of 'cmd' is null")
        }
        ms = message.NewMsg_3(cmd)
    }
    if category == "4"{
        offer :=c.PostForm("offer")
        if offer == ""{
            return nil,errors.New("paraments of 'offer' is null")
        }
        ms = message.NewMsg_4(offer) 
    }
    cmdms := &message.Cmd{
        ClientId:clientId,
        Topic:topic,
        Type:category,
        Message:ms,
    }
    return cmdms,nil
}
func webrtcConn(clientId,topic,category ,httpPort string){
    peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
        ICEServers: []webrtc.ICEServer{
            {
                // URLs: []string{"stun:stun.l.google.com:19302"},
                URLs: []string{"stun:47.99.78.179:3478"},
            },
        },

    })
    if err != nil {
        panic(err)
    }
    log.Println(peerConnection.ConnectionState())
    // Create an offer to send to the browser
    offer, err := peerConnection.CreateOffer(nil)
    if err != nil {
        panic(err)
    }
    // Sets the LocalDescription, and starts our UDP listeners
    err = peerConnection.SetLocalDescription(offer)
    if err != nil {
        panic(err)
    }
    offerBase64 ,err:=sdp.Encode(offer) 
    if err != nil {
        panic(err)
    }
    // Exchange the offer for the answer
    ms := message.NewMsg_4(offerBase64)
    cmdms := &message.Cmd{
        ClientId:clientId,
        Topic:topic,
        Type:category,
        Message:ms,
    }
    info ,_:= cronutil.HttpPost(cmdms,httpPort)
    type message struct{
        Message string
    }
    infoParse := message{}
    err = json.Unmarshal([]byte(info), &infoParse)
    if err!=nil{
        fmt.Println(err)
    }
    answerBase64 :=strings.Split(infoParse.Message,":")[1]
    //除去空白,解析出的字符串前面多了个空白？
    answerBase64 =strings.TrimSpace(answerBase64)
    answerJson := webrtc.SessionDescription{}
    sdp.Decode(answerBase64,&answerJson)
    // Apply the answer as the remote description
    err = peerConnection.SetRemoteDescription(answerJson)
    if err != nil {
        panic(err)
    }

    // Block forever
    // select {}
}
