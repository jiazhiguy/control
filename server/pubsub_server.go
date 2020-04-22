package main

import (
    "context"
    "log"
    "net"
    // "strings"
    "time"
    "reflect"
    // "net/http"

    "grpc-demo/utils"
    pb "grpc-demo/proto"
    "github.com/golang/protobuf/ptypes/empty"
    "google.golang.org/grpc"
    "google.golang.org/grpc/keepalive"
    // pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
    port = ":8123"
)
type pubSubServer struct {
    pub *local_pubsub.Publisher
}
func NewPubsubService() *pubSubServer {
    return &pubSubServer{
        pub: local_pubsub.NewPublisher(100*time.Millisecond, 10),
    }
}
func (p *pubSubServer) Publish(ctx context.Context, req *pb.PulishMessage) (*empty.Empty, error) {
    // p.pub.Publish(req.Name)
    log.Println("Publish")
    p.pub.Publish(req)
    return &empty.Empty{}, nil
}
func (p *pubSubServer) Subscribe(channel *pb.Channel, stream pb.PubsubService_SubscribeServer) error {
    // ch := p.pub.SubscribeTopic(func(v interface{}) bool {
    //     if key, ok := v.(string); ok {
    //         if strings.HasPrefix(key, channel.Name) {
    //             return true
    //         }
    //     }
    //     return false
    // })
    log.Println(channel)
    ch := p.pub.SubscribeTopic(func(v interface{}) bool {
        // if key, ok := v.(string); ok {
        //     if strings.HasPrefix(key, channel.Name) {
        //         return true
        //     }
        // }
        log.Println(v)

        if key ,ok := v.(*pb.PulishMessage) ;ok{
            if reflect.DeepEqual(key.Topic,channel) {
                return true
            }
        }
        return false
    })
    for v := range ch {
        if err := stream.Send(v.(*pb.PulishMessage).Result); err != nil {
            return err
        }
    }
    // return NewPublisher
    return nil
}
func (p *pubSubServer) GetChanneklList(client *pb.Channel, stream pb.PubsubService_GetChanneklListServer) error {
    // GetTopicById(client)
    // channelList:=[]string{"golang","jz2","jz3"}
    // if client.Id=="test"{
    //     channelList=[]string{"golang","jz4","jz5"}
    // }
    
    // for _,v :=range channelList {
    //     if err :=stream.Send(&pb.Channel{Name: v.ClientId,Topic :v.Topic});err !=nil {
    //         return err
    //     }
    // }
    if err :=stream.Send(&pb.Channel{Name: client.Name,Topic :client.Topic});err !=nil {
            return err
        }
    return nil
}
func GetTopicById(client *pb.Client) {
    log.Printf("clientId:%s\n",client.Id)
    url:="http://localhost:8000/api/user/userList"
    data:=local_pubsub.Get(url)
    log.Printf("%+v\n",data)
}
func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v\n", err)
    } else{
        log.Printf("listening on port %s\n", port)
    }
    var kaep = keepalive.EnforcementPolicy{
        MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
        PermitWithoutStream: true,            // Allow pings even when there are no active streams
    }

    var kasp = keepalive.ServerParameters{
        MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
        //MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
        MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
        Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
        Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
    }

    s:= grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
    // s := grpc.NewServer(
    //     grpc.KeepaliveParams(keepalive.ServerParameters{
    //         MaxConnectionIdle: 5 * time.Minute,           // <--- This fixes it!
    //     }),
    // )
    pubSubServer:=NewPubsubService()
    pb.RegisterPubsubServiceServer(s,pubSubServer)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}


