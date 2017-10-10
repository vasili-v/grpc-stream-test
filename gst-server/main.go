package main

//go:generate bash -c "mkdir -p $GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/stream && protoc -I $GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/ $GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/stream.proto --go_out=plugins=grpc:$GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/stream && ls $GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/stream"

import (
	"fmt"
	"io"
	"net"
	"sync"

	"google.golang.org/grpc"

	pb "github.com/vasili-v/grpc-stream-test/gst-server/stream"
)

func handler(in *pb.Request) *pb.Response {
	out := pb.Response{
		Id:      in.Id,
		Payload: make([]byte, len(in.Payload)),
	}

	for i := range out.Payload {
		out.Payload[i] ^= 0x55
	}

	return &out
}

type server struct{}

func (s *server) Test(stream pb.Stream_TestServer) error {
	fmt.Println("got new stream")

	ch := make(chan *pb.Response)
	th := make(chan int, limit)
	go func() {
		defer close(ch)

		var wg sync.WaitGroup
		defer wg.Wait()

		for {
			in, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("stream depleted")
				return
			}

			if err != nil {
				fmt.Printf("receiving error: %s\n", err)
				return
			}

			wg.Add(1)
			th <- 0
			go func(in *pb.Request) {
				defer func() {
					<-th
					wg.Done()
				}()

				ch <- handler(in)
			}(in)
		}
	}()

	for out := range ch {
		err := stream.Send(out)
		if err != nil {
			fmt.Printf("sending error: %s\n", err)
			return err
		}
	}

	return nil
}

func main() {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	opts := []grpc.ServerOption{}
	if compression {
		opts = append(opts,
			grpc.RPCCompressor(grpc.NewGZIPCompressor()),
			grpc.RPCDecompressor(grpc.NewGZIPDecompressor()),
		)
	}

	p := grpc.NewServer(opts...)
	pb.RegisterStreamServer(p, &server{})
	err = p.Serve(ln)
	if err != nil {
		panic(err)
	}
}
