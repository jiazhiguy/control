package message 
import(
	"time"
    
)
type Cmd struct {
    ClientId string `json:"id"`
    Topic string `json:"topic"`
    Type string `json:"type"`
    Message interface{} `json:"message"`
}
//播放音乐
type Msg_1 struct {
    LoopNum int `json:"loop"`//循环播放次数
    Uri string `json:"path"`//播放地址
    Timestamp int64`json:"time"`
}
//推流
type Msg_2 struct {
    LiveCmd string//指令
    Operate bool //开启和关闭
    Timestamp int64
}
//截图
type Msg_3 struct {
    ShotCmd string `json:"shot"`
    // Uri:path,
    Timestamp int64 `json:"time"`
}
//webrtc sdp
type Msg_4 struct {
    Offer string
    Timestamp int64
}
//返回信号
type Response struct {
    Cmdmessage *Cmd
    Msg string
}
type FileMeta struct {
	FileSize int `json:"filesize"`
}
func NewMsg_1(num int,path string) Msg_1{
    return Msg_1{
        LoopNum:num,
        Uri:path,
        Timestamp:time.Now().Unix(),
    }
}
func NewMsg_2(cmd string,start bool) Msg_2{
    return Msg_2{
        LiveCmd:cmd,
        Operate:start,
        Timestamp:time.Now().Unix(),
    }
}
func NewMsg_3(cmd string) Msg_3{
    return Msg_3{
        ShotCmd:cmd,
        // Uri:path,
        Timestamp:time.Now().Unix(),
    }
}
func NewMsg_4(offer string) Msg_4{
    return Msg_4{
        Offer:offer,
        // Uri:path,
        Timestamp:time.Now().Unix(),
    }
}