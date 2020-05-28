package cronutil
import (
	// "github.com/robfig/cron/v3"
	 message "grpc-demo/utils/message"
	 "net/http"
	 "io/ioutil"
	 "strings"
	 "fmt"
)
func HttpPost(cmd *message.Cmd,host string) (string,error){
    // fmt.Println("****post发送操作指令******")
    clientId := cmd.ClientId
    topic := cmd.Topic
    category := cmd.Type
    ms := cmd.Message
    // fmt.Printf("---%+v---",ms)
    var sendParament string ="" 
    switch category{
    case "1":
        m,_ := ms.(message.Msg_1)
        // fmt.Printf("%+v",m)
        sendParament = fmt.Sprintf("id=%s&&topic=%s&&type=%s&&times=%d&&path=%s",clientId,topic,category,m.LoopNum,m.Uri)
        fmt.Println(sendParament)
    case "2":
        m,_ := ms.(message.Msg_2)
        // fmt.Printf("%+v",m)
        startString := "true"
        if m.Operate{
            startString = "false"
        }
        sendParament = fmt.Sprintf("id=%s&&topic=%s&&type=%s&&open=%d&&cmd=%s",clientId,topic,category,startString,m.LiveCmd)
        fmt.Println(sendParament)
    case "3":
        m,_ := ms.(message.Msg_3)
        // fmt.Printf("%+v",m)
        sendParament = fmt.Sprintf("id=%s&&topic=%s&&type=%s&&cmd=%s",clientId,topic,category,m.ShotCmd)
        fmt.Println(sendParament)
  
    case "4":
        m,_ := ms.(message.Msg_4)
        // fmt.Printf("%+v\n",m)
        sendParament = fmt.Sprintf("id=%s&&topic=%s&&type=%s&&offer=%s",clientId,topic,category,m.Offer)
        // fmt.Println("sendParament:",sendParament)
    }
    resp, err := http.Post(host,
        "application/x-www-form-urlencoded",
        strings.NewReader(sendParament))
    if err != nil {
        fmt.Println(err)
        return "",err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "",err
    }
    // fmt.Println("----",string(body))
    return string(body),nil
    
}
