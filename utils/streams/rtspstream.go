package webrtcserver

import (
	// "fmt"
	// "io"
	"time"
	"log"
	"sync"

	// "encoding/json"
	// "github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"

	// "strings"
	// "encoding/json"
	// "math/rand"

	"github.com/deepch/vdk/codec/h264parser"

	"github.com/deepch/vdk/format/rtsp"
	"github.com/deepch/vdk/av"
	"github.com/pion/webrtc/v2/pkg/media"
	
)
type FileInfo struct {
	Name string	 
	Size int64
}
type StreamsMutex struct {
	sync.RWMutex
	Map map[string]*RtspMessage
}
type SimplePeer struct {
	Conn *webrtc.PeerConnection
	Offer string
	Answer string
	TrackChan chan *webrtc.Track
	Tracks []*webrtc.Track
}
type WebrtcServer struct {
	Url chan string
	//LatestStream  *RtspMessage
	Streams *StreamsMutex
	Liveon chan string
}
type RtspMessage struct{
	Id string
	Codecs []av.CodecData
	Packet chan av.Packet
	Session  *rtsp.Client
}
func NewWebrtcServer() *WebrtcServer{
	url :=make(chan string,10)
	liveon := make(chan string)
	streamPool := &StreamsMutex{
		Map:map[string]*RtspMessage{},
	}
	return &WebrtcServer{
		Url:url,
		Streams:streamPool,
		Liveon:liveon,
	}
}
func (self *WebrtcServer) Run(){
	go func(){
		for {
			select {
			case url :=<-self.Url:
				go rtspStreams(url,self)
			}
		}
	}()
}
func (self *WebrtcServer) hasStream(url string)(rp *RtspMessage,ok bool){
	 rp,ok =self.Streams.Map[url]
	 return
}
func (self *WebrtcServer) InStream (url string,track  *webrtc.Track,fn func (out *RtspMessage,videoTrack  *webrtc.Track)){
	self.Url <-url

	// start :=make(chan *RtspMessage)
	go func	(){
		for url :=range self.Liveon{
			go func(url string){
				// log.Printf("$$$$$$$$track%+v$$$$$$$$$\n",track)
				fn(self.Streams.Map[url],track)
			}(url)
		} 
		// for	{
		// 	// self.Streams.RLock()
		// 	// _,ok:=self.Streams.Map[url]
		// 	// self.Streams.RLock()
		// 	// if !ok{
		// 	// 	continue
		// 	// }
		// 	select{
		// 	case <-self.Streams.Map[url].Packet:
		// 		log.Printf("Packet......")
		// 		go fn(self.Streams.Map[url],track)	
		// 		// go func(){
		// 		// 	for {
		// 		// 		select{
		// 		// 		case track_ := <-track:
		// 		// 			fn(self.Streams[url],track_)	
		// 		// 		}
		// 		// 	}

		// 		// }()
		// 	}
		// 	// } 
		// }

	}()

}
func rtspStreams(url string ,server * WebrtcServer) {
	log.Println("server run")
	log.Printf("create rtspStreams:%s\n",url)
	// for k, v := range Config.Streams {
		// go func(url string ,server * WebrtcServer) {
			// mschan :=server.LastStream
			name := url
			// server.Streams.Lock()
			// server.Streams.Map[url]=ms
			// server.Streams.Unlock()
			retrynum := 5
			for {
				log.Println(name, "connect", url)
				rtsp.DebugRtsp = true
				session, err := rtsp.Dial(url)
				if err != nil {
					log.Println(name, err)
					time.Sleep(5 * time.Second)
					retrynum = retrynum -1
					if(retrynum == 0) {
						break
					}
					continue
				}
				session.RtpKeepAliveTimeout = time.Duration(10 * time.Second)
				if err != nil {
					log.Println(name, err)
					time.Sleep(5 * time.Second)
					continue
				}
				codec, err := session.Streams()
				if err != nil {
					log.Println(name, err)
					time.Sleep(5 * time.Second)

					continue
				}
				ms :=&RtspMessage{
					Id:url,
					Codecs:codec,
					Session:session,
				}
				ms.Packet=make(chan av.Packet)
				server.Streams.Lock()
				server.Streams.Map[url]=ms
				server.Streams.Unlock()
				log.Printf("---%+v----\n",server.Streams)		
				// Config.coAd(name, codec)
				// done :=make (chan bool)

				// go func(done chan bool,url string,server * WebrtcServer){
					fisrtLive := true
					for {
						pkt, err := session.ReadPacket()//长时间运行runtime error: makeslice: len out of range位置format/rtsp/client.go:443 
						if err != nil {
							log.Println(name, err)
							// server.Streams.Lock()
							delete(server.Streams.Map,url)
							// server.Streams.Unlock()
							// done <-true
							break
						}
						// server.Streams.RLock()
						if stream,ok:=server.Streams.Map[url];ok{
							if fisrtLive {
								server.Liveon <-url
								fisrtLive = false
							}
							
							stream.Packet <-pkt
						}
						// server.Streams.RUnlock()
						// Config.cast(name, pkt)
					}
				// }(done,url,server)
				// <-done
				session.Close()
				log.Println(name, "reconnect wait 5s")
				time.Sleep(5 * time.Second)
			}
		// }(url,server)
	// }
}
// videoTrack *webrtc.Track
func Deal(out *RtspMessage,track  *webrtc.Track){
 		
	// go func() {
		ms :=out
		codecs := ms.Codecs

		if codecs == nil {
			log.Println("No Codec Info")
			return
		}
			//sps, pps := []byte{}, []byte{}
		sps := codecs[0].(h264parser.CodecData).SPS()
		pps := codecs[0].(h264parser.CodecData).PPS()
		control := make(chan bool, 10)
		
		// conected := make(chan bool, 10)
		// defer peerConnection.Close()
		// peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		// 	fmt.Printf("Connection State has changed %s \n", connectionState.String())
		// 	if connectionState != webrtc.ICEConnectionStateConnected {
		// 		log.Println("Client Close Exit")
		// 		control <- true
		// 		return
		// 	}
		// 	if connectionState == webrtc.ICEConnectionStateConnected {
		// 		conected <- true
		// 	}
		// })
		// <-conected
		// cuuid, ch := Config.clAd(suuid)
		// defer Config.clDe(suuid, cuuid)
		var pre uint32
		var start bool
		for {
			select {
			case <-control:
				return
			case pck := <-ms.Packet:
				// log.Println("pck:",pck)
			// case pck := <-ch:
				if pck.IsKeyFrame {
					start = true
				}
				if !start {
					continue
				}
				if pck.IsKeyFrame {
					pck.Data = append([]byte("\000\000\001"+string(sps)+"\000\000\001"+string(pps)+"\000\000\001"), pck.Data[4:]...)

				} else {
					pck.Data = pck.Data[4:]
				}
				var ts uint32
				if pre != 0 {
					ts = uint32(timeToTs(pck.Time)) - pre
				}
				go func(){
					// log.Println("THIS IS  WriteSample ")
					// select{	
					// case videoTrack := <-track:
					
					videoTrack := track

					// log.Println("THIS IS  WriteSample ")
						err := videoTrack.WriteSample(media.Sample{Data: pck.Data, Samples: uint32(ts)})
						pre = uint32(timeToTs(pck.Time))
						if err != nil {
							return
						}
					// }
				}()

			}
		}
	// }()
}
func timeToTs(tm time.Duration) int64 {
	return int64(tm * time.Duration(90000) / time.Second)
}



