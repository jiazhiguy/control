package main

import (
    "fmt"
    "context"
    "log"
    "strconv"
    "io"
    "time"
    // "strings"

    pb "grpc-demo/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/keepalive"
)
type Unit struct {
    *pb.Channel
    *pb.SubscribeResult
}
func main() {
    fmt.Println("**1.消息发布端和接收端设备ID和主题填写一致**")
    fmt.Println("**2.先开启接受端填写参数，再通过发布端发布指令**")
    fmt.Println("**3.发布端地址MP3地址可以是url,也可以是本地MP3,本地文件和接受端放在一起,地址为 ./文件名**")
    // const serverIp ="localhost:8123"
    //********************connect*********************
    const serverIp ="47.99.78.179:8123"
    // cmdline := "ffmpeg -re -i D:/21.mp4 -c copy -f flv rtmp://localhost:1935/live/movie"
    var initLoop = 1
    var clientId,topic,path,loop,optype,liveCmd,start string

    var kacp = keepalive.ClientParameters{
        Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
        Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
        PermitWithoutStream: true,             // send pings even without active streams
    }
    conn, err := grpc.Dial(serverIp, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
    // conn, err := grpc.Dial(serverIp, grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    client := pb.NewPubsubServiceClient(conn)


    //******************* 发布和订阅内容处理***********************
    signal :=make (chan *pb.Channel)
    send := make (chan *pb.PulishMessage)
    response :=make (chan *pb.SubscribeResult)
    //处理订阅消息go
    go func(){
        for {
            select{
            case v :=<-signal:
                // fmt.Println(v)
                Subscribe(client,v,response)
            }
        }
    }()
    //处理发布消息go
    go func(){
        for {
            select{
            case v :=<-send:
                Publish(client,v)
            }
        } 
    }()
     //*********************input **********************
    // fmt.Print("用户名:")
    // fmt.Scanln(&usernanme)
    // fmt.Print("登录密码:")
    // fmt.Scanln(&password)
    fmt.Println("==================消息指令发布端==================")
       for {
        fmt.Println(" ")
        fmt.Print("设备ID:")
        fmt.Scanln(&clientId)
        fmt.Print("主题:")
        fmt.Scanln(&topic)
        fmt.Printf("操作类型(1代表播放音乐，2代表视频推流):")
        fmt.Scanln(&optype)
        switch optype {
        case "1":{
            fmt.Print("MP3播放地址:")
            fmt.Scanln(&path)
            fmt.Print("循环播放次数:")
            fmt.Scanln(&loop)
            if clientId ==""{
                clientId= "jz4"
            }
            if topic ==""{
                topic= "playsound"
            }
            if loopNum,err := strconv.Atoi(loop) ; err != nil{
                log.Fatal("循环次数必须为有效数字")
            } else{
                 initLoop =loopNum
            }
            if path == ""{
                path ="./Lame_Drivers_-_01_-_Frozen_Egg.mp3"
            }
            timestamp := time.Now().Unix()
            newTopic :=fmt.Sprintf("%s_%d",topic,timestamp)
            channel := &pb.Channel{Name: clientId,Topic: newTopic}
            signal <-channel
            go func () {
                for responseMessage :=range response{
                    fmt.Println("-----接受端返回信息:-----\n"+responseMessage.Msg)
                }
            }()
            send <- &pb.PulishMessage{
                    Topic: &pb.Channel{Name: clientId,Topic: topic},
                    Result: &pb.SubscribeResult{Msg: path,Loop: int32(initLoop),Fast:1,Pause:false,Volume:0,Type:0,Timestamp:timestamp},
                }
            // time.Sleep(4*time.Second)
            path=""
            loop="1"
        }
        case "2":{
            cmdPause :=false
            //ffmpeg|_|-re|_|-i|_|rtsp://admin:admin123@192.168.2.241/h264/ch1/main/av_stream|_|-c|_|copy|_|-f|_|flv|_|rtmp://47.99.78.179:1935/live/movie
            fmt.Print("推流指令:")
            fmt.Scanln(&liveCmd)
            fmt.Print("开启或关闭(1为关闭):")
            fmt.Scanln(&start)
            if start=="1"{
                cmdPause =true
            }
            // newLiveCmd :=strings.Replace(liveCmd ," ","|_|",-1)
            log.Println(liveCmd)
            if clientId ==""{
                clientId= "jz4"
            }
            if topic ==""{
                topic= "playsound"
            }
            timestamp := time.Now().Unix()
            newTopic :=fmt.Sprintf("%s_%d",topic,timestamp)
            channel := &pb.Channel{Name: clientId,Topic: newTopic}

            signal <-channel
            go func () {
                for responseMessage :=range response{
                    fmt.Println("-----接受端返回信息:-----\n"+responseMessage.Msg)
                }
            }()
            sendcmd := &pb.PulishMessage{
                        Topic: &pb.Channel{Name: clientId,Topic: topic},
                        Result: &pb.SubscribeResult{Msg: liveCmd,Type:2,Pause:cmdPause,Timestamp:timestamp},
                    }
            log.Printf("%+v",sendcmd)
            send <- sendcmd
            // time.Sleep(4*time.Second)
        }

        }
        clientId = ""
        topic = ""
        start =" "
        time.Sleep(4*time.Second)
    }
}
//发布内容
func  Publish(client pb.PubsubServiceClient,msg  *pb.PulishMessage){
    _, err := client.Publish(context.Background(),msg )
    if err != nil {
        log.Println(err)
        // log.Fatal(err)
    }
}
//根据channel内容完成订阅
func  Subscribe(client pb.PubsubServiceClient,channel *pb.Channel,c chan *pb.SubscribeResult){
    // fmt.Printf("channel:%+v\n",channel)
    stream, err := client.Subscribe(
        context.Background(), channel,
    )
    if err != nil {
        log.Fatal(err)
    }
    go func(){
        for {
            reply, err := stream.Recv()
            if err != nil {
                if err == io.EOF {
                    break
                }
                log.Fatal(err)
            }
            c<-reply
        }
    }()
}