package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	pb "github.com/ronitboddu/pi/pb/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(1024*1024*1024),
		grpc.MaxCallSendMsgSize(1024*1024*1024)))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProcessClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	temp, err := readLines("pi.txt")
	check(err)
	details, err := c.GetDetails(ctx, &pb.Textfile{FileName: "pi.txt", Data: temp}, grpc.MaxCallSendMsgSize(1024*1024*1024),
		grpc.MaxCallRecvMsgSize(1024*1024*1024))
	check(err)
	log.Printf(`Total: %d`, details.GetTotalCount())
	for key, val := range details.GetCount() {
		log.Println(key + " : " + strconv.Itoa(int(val)))
	}
}

func readLines(filename string) (string, error) {
	var lines = ""
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return lines, err
	}
	buf := bytes.NewBuffer(file)
	for {
		line, err := buf.ReadString('\n')
		if len(line) == 0 {
			if err != nil {
				if err == io.EOF {
					break
				}
				return lines, err
			}
		}
		lines += line
		if err != nil && err != io.EOF {
			return lines, err
		}
	}
	return lines, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
