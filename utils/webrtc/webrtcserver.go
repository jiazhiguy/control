package webrtcserver

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"math/rand"
	"sync"
	"image/png"
	"bytes"

	"github.com/pion/webrtc/v2"
	"grpc-demo/utils/signal"
	 pb "grpc-demo/proto"
	 message "grpc-demo/utils/message"
	 job "grpc-demo/utils/jobs"
	 myutil "grpc-demo/utils/my_utils"
	 stream "grpc-demo/utils/streams"
	 "grpc-demo/client/models"
)
type TrackPool struct {
	sync.RWMutex
	Map map[string] *webrtc.Track
}

var rtspserver *stream.WebrtcServer
var tracksHub =&TrackPool{
	Map:map[string]*webrtc.Track{},
}
var upContent =job.Context{}
func init(){
	go func(){
		rtspserver =stream.NewWebrtcServer()
		rtspserver.Run()
	}()
}
func Run (sdpChan,answerchan chan string,reback chan *pb.PulishMessage){
	// rand.Seed(time.Now().UTC().UnixNano())
	for {	fmt.Println("")
			fmt.Println("Curl an base64 SDP to start sendonly peer connection")
			fmt.Println("")
			recvOnlyOffer := webrtc.SessionDescription{}
			signal.Decode(<-sdpChan, &recvOnlyOffer)
			// Create a new PeerConnection
			// peerConnection, err := api.NewPeerConnection(peerConnectionConfig)
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
			peerConnection.OnDataChannel(func(datachannel *webrtc.DataChannel){
				log.Println("<<<<<<<<")
				label := datachannel.Label()
				log.Println(label)
				if label == "control"{
					log.Println("control")
					datachannel.OnMessage(func(Data webrtc.DataChannelMessage){
						log.Printf("----%s----",Data)
					    cmdms := &message.Cmd{}
				        err :=signal.Decode(string(Data.Data), &cmdms)
				        if err !=nil{
				        	log.Println(err)
				        	return
				        }
						log.Printf("%+v",cmdms)
						switch cmdms.Type{
							case "1":
								var loop =1
								log.Println("CASE 1")
								fmt.Printf("%T", cmdms.Message)
								if ms ,ok := cmdms.Message.(map[string]interface{}) ;ok{
									loop = int(ms["loop"].(float64))
									path := ms["path"].(string)
									errorCh := make (chan error)
									go func () {
										for ms :=range errorCh {
											log.Println(ms)
										}
									}()
									job.Sound(path,loop,upContent,errorCh)
								}
							case "2":
								log.Println("CASE 2")
							
							case "3":
								go func () {
									log.Println("CASE 3")
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
								    time_now := time.Now()
									savepath := direct	
								    savename := "image_"+strconv.FormatInt(time_now.UnixNano(),10)+".png"
									errorCh := make (chan error)
									doneCh := make (chan bool)
									shotCmd :=""
									// timestamp :=0
									if ms ,ok := cmdms.Message.(map[string]interface{}) ;ok{
										shotCmd = ms["shot"].(string)
									}
									go func () {
										for ms :=range errorCh {
											log.Println(ms)
										}
									}()
									go func (){
										loop:for{
												select {
													case <-doneCh :
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
														//发送缩略图大小
														fileSize := len(imageBytes)
														log.Println("Create shot channek")
														sendchannel ,_:= peerConnection.CreateDataChannel("shot",&webrtc.DataChannelInit{})
														data ,_:= signal.Encode(message.FileMeta{FileSize:fileSize})
														sendchannel.SendText(data)
														//切片发送图片数据
														chunk_len:= 64000;
														sendMax := fileSize/chunk_len
														index :=1; 
														for index<=sendMax+1 {
															lastlen :=index*chunk_len
															if index*chunk_len>fileSize{
																lastlen =fileSize
															}
															sendchannel.Send(imageBytes[(index-1)*chunk_len:lastlen])
															index = index+1
														}
														sendchannel.Close()
														break loop
														return
													case <-time.After(30*time.Second):
														log.Println("shot picture lost")
														return
														break
												}
											}

									}()
									// log.Printf("%+v",cmdms.Message)
									// fmt.Printf("%T", cmdms.Message)
									// shotCmd :=""
									// // timestamp :=0
									// if ms ,ok := cmdms.Message.(map[string]interface{}) ;ok{
									// 	shotCmd = ms["shot"].(string)
									// }
									if shotCmd == ""{
										return
									}else{
										if shotCmd != "screen" {
											shotCmd = fmt.Sprintf("ffmpeg |_|-i|_|%s|_|-vframes|_|1|_|-y|_|-f|_|image2",shotCmd)
										}
										// log.Println(shotCmd)
										job.Shot(shotCmd,savepath,savename,doneCh,errorCh)
									}
	
								}()
							case "5":
								rand.Seed(time.Now().Unix())
								type info struct{
									Url string `json:"url"`
									// Offer string `json:"offer"`
								}
								response := info{}
								signal.Decode(string(Data.Data), &response)
								rtspUrl := response.Url
								log.Println("b0",len(peerConnection.GetReceivers()))
								go func(){
									log.Println("CASE 5")
									var hasStream bool
									videoTrack, err := peerConnection.NewTrack(webrtc.DefaultPayloadTypeH264, rand.Uint32(), "video", rtspUrl)
									if err != nil {
										log.Println(err)
									}
									tracksHub.Lock()
									if track,ok := tracksHub.Map[rtspUrl];ok{
										videoTrack = track
										hasStream = true
									}
									tracksHub.Unlock()
									// _ ,err = peerConnection.AddTransceiverFromTrack(videoTrack,webrtc.RtpTransceiverInit{
									// 	Direction:webrtc.RTPTransceiverDirectionSendonly,
									// })
									log.Println("b1",len(peerConnection.GetReceivers()))
							        _ , err = peerConnection.AddTransceiverFromTrack(videoTrack, webrtc.RtpTransceiverInit{
									Direction: webrtc.RTPTransceiverDirectionSendrecv,})
									log.Println("b2",len(peerConnection.GetReceivers()))
									// _, err = peerConnection.AddTrack(videoTrack)
									if err != nil {
										panic(err)
									}
									sdpchannel ,_:= peerConnection.CreateDataChannel("sdp",&webrtc.DataChannelInit{})
									sdpchannel.SendText("ready")
									sdpchannel.OnMessage(func(msg webrtc.DataChannelMessage){
										var offer webrtc.SessionDescription
									    signal.Decode(string(msg.Data),&offer)
										if err := peerConnection.SetRemoteDescription(offer); err != nil {
													log.Println(err)
													// panic(err)
										}
										log.Println("b3",len(peerConnection.GetReceivers()))
										answer, err := peerConnection.CreateAnswer(nil)
										if err != nil {
											log.Println(err)
											// panic(err)
										} else if err = peerConnection.SetLocalDescription(answer); err != nil {
											log.Println(err)
											// panic(err)
										}
										log.Println("b4",len(peerConnection.GetReceivers()))
										log.Printf("---%+v---\n",peerConnection.GetSenders())
										answerBase ,err:= signal.Encode(answer) 
										if err!=nil{
											log.Println(err)
											return
										}
										sdpchannel.SendText(answerBase)
										// sdpchannel.Close()
										if(!hasStream){
								    		go func(){
												conected := make(chan bool, 10)
												if(peerConnection.ICEConnectionState() == webrtc.ICEConnectionStateConnected){
													log.Println("START stream")
													conected <- true
												}
												<-conected
												tracksHub.Lock()
												tracksHub.Map[rtspUrl]=videoTrack
												tracksHub.Unlock()
												rtspserver.InStream(rtspUrl,videoTrack,stream.Deal)

											}()
										}
	
									})
								}()
						}
					})
				}
				log.Println(">>>>>>>>")
			})
			log.Println(peerConnection.GetSenders())
			log.Println(peerConnection.GetReceivers())
			// Set the remote SessionDescription
			err = peerConnection.SetRemoteDescription(recvOnlyOffer)
			// log.Printf("sdp offer%+v",recvOnlyOffer)
			if err != nil {
				log.Println(err)
				// panic(err)
			}
			// Create answer
			answer, err := peerConnection.CreateAnswer(nil)
			if err != nil {
				log.Println(err)
				// panic(err)
			}

			// Sets the LocalDescription, and starts our UDP listeners
			err = peerConnection.SetLocalDescription(answer)
			if err != nil {
				log.Println(err)
			}
			answer_encode ,_ := signal.Encode(answer)
           // Apply the answer as the remote description
			answerchan <- answer_encode
			// go func(){
			// 	conected := make(chan bool, 10)
			// 	log.Println("Client-----")
			// 	log.Printf("%+v",peerConnection)
			//     peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
			// 		fmt.Printf("Connection State has changed %s \n", connectionState.String())
			// 		if connectionState != webrtc.ICEConnectionStateConnected {
			// 			log.Println("Client Close Exit")
			// 			// control <- true
			// 			// return
			// 		}
			// 		if connectionState == webrtc.ICEConnectionStateConnected {
			// 			log.Println("connected")
						
			// 			conected <- true
			// 		}
			// 	})
			// 	<-conected
			// 	rtspserver.InStream("rtsp://admin:admin123@192.168.2.241/h264/ch1/main/av_stream",videoTrack,stream.Deal)
			// }()
			// Get the LocalDescription and take it to base64 so we can paste in browser

			
		}
		select{}
	}


// func addVideoTrack(peerConnection *webrtc.peerConnection,conected chan bool){
// 	// control := make(chan bool, 10)
// 	// conected := make(chan bool, 10)
// 	// defer peerConnection.Close()
// 	peerConnection.add
// 	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
// 		fmt.Printf("Connection State has changed %s \n", connectionState.String())
// 		if connectionState != webrtc.ICEConnectionStateConnected {
// 			log.Println("Client Close Exit")
// 			// control <- true
// 			// return
// 		}
// 		if connectionState == webrtc.ICEConnectionStateConnected {
// 			conected <- true
// 		}
// 	})
// 	// <-conected
// }