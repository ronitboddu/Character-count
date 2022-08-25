package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/ronitboddu/pi/pb/proto"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type piServer struct {
	pb.UnimplementedProcessServer
}

var count = 0

func (s *piServer) GetDetails(ctx context.Context, in *pb.Textfile) (*pb.Details, error) {
	log.Printf("Received: %v", in.GetFileName())
	var mapChar = make(map[string]int32)
	var data = in.GetData()
	var length = len(data)
	var index = 0
	for {
		wg.Add(1)
		if index < length && length <= index+10 {
			//m.Lock()
			go countDigits(index, length, data, mapChar, ch)
			count += <-ch
			break
		}
		//m.Lock()
		go countDigits(index, index+10, data, mapChar, ch)
		count += <-ch
		index += 10
	}
	wg.Wait()

	return &pb.Details{TotalCount: int32(count), Count: mapChar}, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var wg = sync.WaitGroup{}
var m = sync.RWMutex{}
var ch = make(chan int)

func main() {
	lis, err := net.Listen("tcp", port)
	check(err)

	s := grpc.NewServer(grpc.MaxRecvMsgSize(1024*1024*1024), grpc.MaxSendMsgSize(1024*1024*1024))
	pb.RegisterProcessServer(s, &piServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func countDigits(start int, end int, temp string, mapChar map[string]int32, ch chan<- int) {
	var total = 0
	for i := start; i < end; i++ {
		mapChar[string(temp[i])] += 1
		total += 1
	}
	//m.Unlock()
	ch <- total
	wg.Done()
}
